package agent

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/core"
	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/yuin/gopher-lua"
	"google.golang.org/grpc"
)

// agentServer implements the gRPC agent server
type agentServer struct {
	pb.UnimplementedAgentServer
	grpcServer *grpc.Server
}

// RunCommand executes a shell command and streams output
func (s *agentServer) RunCommand(in *pb.RunCommandRequest, stream pb.Agent_RunCommandServer) error {
	slog.Info(fmt.Sprintf("Executing command on agent: %s", in.GetCommand()))

	var cmd *exec.Cmd

	// If user is specified and not root, run command as that user
	if in.GetUser() != "" && in.GetUser() != "root" {
		// Use sudo to run as specific user
		cmd = exec.Command("sudo", "-u", in.GetUser(), "bash", "-c", in.GetCommand())
		slog.Info(fmt.Sprintf("Running command as user: %s", in.GetUser()))
	} else {
		// Run as current user (typically root for agent)
		cmd = exec.Command("bash", "-c", in.GetCommand())
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Goroutine to stream stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			chunk := scanner.Text()
			stream.Send(&pb.StreamOutputResponse{StdoutChunk: chunk + "\n"})
		}
	}()

	// Goroutine to stream stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			chunk := scanner.Text()
			stream.Send(&pb.StreamOutputResponse{StderrChunk: chunk + "\n"})
		}
	}()

	err = cmd.Wait()

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			stream.Send(&pb.StreamOutputResponse{Error: err.Error(), Finished: true})
			return err
		}
	}

	stream.Send(&pb.StreamOutputResponse{Finished: true, ExitCode: int32(exitCode)})
	return nil
}

// Shutdown gracefully stops the agent server
func (s *agentServer) Shutdown(ctx context.Context, in *pb.ShutdownRequest) (*pb.ShutdownResponse, error) {
	slog.Info("Shutting down agent server")
	go func() {
		time.Sleep(1 * time.Second)
		s.grpcServer.GracefulStop()
	}()
	return &pb.ShutdownResponse{}, nil
}

// ExecuteTask executes a complete Lua task with workspace
func (s *agentServer) ExecuteTask(ctx context.Context, in *pb.ExecuteTaskRequest) (*pb.ExecuteTaskResponse, error) {
	slog.Info(fmt.Sprintf("Received task: %s from group: %s", in.GetTaskName(), in.GetTaskGroup()))

	// Create a temporary directory for the workspace
	workDir, err := os.MkdirTemp("", "sloth-runner-agent-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	// Unpack the workspace
	if err := extractTarData(bytes.NewReader(in.GetWorkspace()), workDir); err != nil {
		return nil, fmt.Errorf("failed to untar workspace: %w", err)
	}

	// Create temporary Lua script file
	scriptPath := filepath.Join(workDir, "task.lua")
	if err := os.WriteFile(scriptPath, []byte(in.GetLuaScript()), 0644); err != nil {
		return nil, fmt.Errorf("failed to write lua script: %w", err)
	}

	// Log the script content for debugging
	slog.Debug("Agent script content", "script_path", scriptPath, "script_size", len(in.GetLuaScript()))

	// Create new Lua state for execution
	L := lua.NewState()
	defer L.Close()

	// Set task user as global variable for modules to use
	if in.GetUser() != "" {
		L.SetGlobal("__TASK_USER__", lua.LString(in.GetUser()))
		slog.Info("Task will run with user context", "user", in.GetUser(), "task", in.GetTaskName())
	}

	// Register all modules
	luainterface.RegisterAllModules(L)
	luainterface.OpenImport(L, scriptPath)

	// Ensure Modern DSL is registered for task/workflow definitions
	globalCore := core.GetGlobalCore()
	if globalCore == nil {
		// Initialize a minimal core for agent execution
		logger := slog.Default()
		config := core.DefaultCoreConfig()
		if err := core.InitializeGlobalCore(config, logger); err == nil {
			globalCore = core.GetGlobalCore()
		}
	}
	if globalCore != nil {
		modernDSL := luainterface.NewModernDSL(globalCore)
		modernDSL.RegisterModernDSL(L)
	}

	// Parse the Lua script to get task definitions
	taskGroups, err := luainterface.ParseLuaScript(ctx, scriptPath, nil)
	if err != nil {
		slog.Error("Failed to parse lua script on agent", "error", err, "script_path", scriptPath)
		return nil, fmt.Errorf("failed to load task definitions: %w", err)
	}

	// Verify we have task groups
	if len(taskGroups) == 0 {
		slog.Error("No task groups found on agent", "script_path", scriptPath)
		return nil, fmt.Errorf("expected 'TaskDefinitions' to be a table, got nil")
	}

	// Remove delegate_to from all tasks to prevent recursive delegation
	for groupName, group := range taskGroups {
		slog.Info("Agent checking group for delegate_to", "group", groupName, "task_count", len(group.Tasks))
		for i, task := range group.Tasks {
			if task.DelegateTo != nil {
				slog.Info("Removing delegate_to from task on agent", "task", task.Name, "group", groupName, "delegate_to", task.DelegateTo)
				task.DelegateTo = nil
				group.Tasks[i] = task
			}
		}
	}

	slog.Info("Agent parsed task groups", "count", len(taskGroups))

	// Change to the workspace directory so file operations work correctly
	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	if err := os.Chdir(workDir); err != nil {
		return nil, fmt.Errorf("failed to change to workspace directory: %w", err)
	}
	defer os.Chdir(originalDir)

	// Filter the task group to contain only the specified task
	targetTaskName := in.GetTaskName()
	if targetTaskName != "" && targetTaskName != "nil" {
		slog.Info("Agent filtering for specific task", "task", targetTaskName, "group", in.GetTaskGroup())
		if group, exists := taskGroups[in.GetTaskGroup()]; exists {
			var filteredTasks []types.Task
			for _, task := range group.Tasks {
				if task.Name == targetTaskName {
					filteredTasks = append(filteredTasks, task)
					slog.Info("Found target task for agent execution", "task", task.Name)
					break
				}
			}

			if len(filteredTasks) > 0 {
				group.Tasks = filteredTasks
				taskGroups[in.GetTaskGroup()] = group
				slog.Info("Agent will execute only the delegated task", "task", targetTaskName)
			}
		}
	}

	// Create task runner
	runner := taskrunner.NewTaskRunner(L, taskGroups, in.GetTaskGroup(), nil, false, false, &taskrunner.DefaultSurveyAsker{}, in.GetLuaScript())

	// Execute the specific task group
	slog.Info("Agent executing task group", "group", in.GetTaskGroup())
	err = runner.Run()

	// Pack the updated workspace
	var buf bytes.Buffer
	if err := createTarData(workDir, &buf); err != nil {
		return nil, fmt.Errorf("failed to tar workspace: %w", err)
	}

	// Return response based on execution result
	if err != nil {
		slog.Error("AGENT TASK EXECUTION FAILED", "task", in.GetTaskName(), "group", in.GetTaskGroup(), "error", err)

		errorMsg := err.Error()
		var errorDetails strings.Builder

		errorDetails.WriteString(fmt.Sprintf("╔═══════════════════════════════════════════════════════════════════════════════════\n"))
		errorDetails.WriteString(fmt.Sprintf("║ ❌ AGENT EXECUTION FAILURE\n"))
		errorDetails.WriteString(fmt.Sprintf("╠═══════════════════════════════════════════════════════════════════════════════════\n"))
		errorDetails.WriteString(fmt.Sprintf("║ Task  : %s\n", in.GetTaskName()))
		errorDetails.WriteString(fmt.Sprintf("║ Group : %s\n", in.GetTaskGroup()))
		errorDetails.WriteString(fmt.Sprintf("╠═══════════════════════════════════════════════════════════════════════════════════\n"))
		errorDetails.WriteString(fmt.Sprintf("║ 🔴 ERROR:\n"))
		errorDetails.WriteString(fmt.Sprintf("║   %s\n", errorMsg))
		errorDetails.WriteString(fmt.Sprintf("╚═══════════════════════════════════════════════════════════════════════════════════\n"))

		return &pb.ExecuteTaskResponse{
			Success:   false,
			Output:    errorDetails.String(),
			Workspace: buf.Bytes(),
		}, nil
	}

	slog.Info("Agent task execution succeeded", "task", in.GetTaskName(), "group", in.GetTaskGroup())
	return &pb.ExecuteTaskResponse{
		Success:   true,
		Output:    fmt.Sprintf("Task '%s' executed successfully on agent", in.GetTaskName()),
		Workspace: buf.Bytes(),
	}, nil
}

// createTarData creates a tarball from source directory
func createTarData(source string, writer io.Writer) error {
	tw := tar.NewWriter(writer)
	defer tw.Close()
	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(file[len(source):])
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !fi.IsDir() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})
}

// extractTarData extracts a tarball to destination directory
func extractTarData(reader io.Reader, dest string) error {
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
}
