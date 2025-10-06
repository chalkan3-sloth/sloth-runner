package taskrunner

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	lua "github.com/yuin/gopher-lua"
)

// executeOnAgent handles execution of a task on a remote agent via gRPC
func (tr *TaskRunner) executeOnAgent(ctx context.Context, t *types.Task, agentAddress string, session *types.SharedSession, groupName string) error {
	// Connect to the agent
	pterm.DefaultBox.
		WithTitle("🔗 Agent Connection").
		WithTitleTopLeft().
		WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).
		Printfln("Task:  %s\nAgent: %s", pterm.Cyan(t.Name), pterm.Yellow(agentAddress))

	conn, err := grpc.Dial(agentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ CONNECTION FAILED").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Agent: %s\nTask:  %s\n\nError: %v\n\n"+
					"💡 Troubleshooting:\n"+
					"  • Check agent status: systemctl status sloth-runner-agent\n"+
					"  • Verify agent address is correct\n"+
					"  • Check network: ping <agent-host>\n"+
					"  • Verify firewall rules",
				pterm.Yellow(agentAddress),
				pterm.Cyan(t.Name),
				err,
			)
		pterm.Println()

		slog.Error("Failed to connect to agent",
			"agent_address", agentAddress,
			"task", t.Name,
			"error", err)
		return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to connect to agent %s: %w", agentAddress, err)}
	}
	defer conn.Close()
	c := pb.NewAgentClient(conn)

	// Create a tarball of the workspace
	var buf bytes.Buffer
	if err := createTar(session.Workdir, &buf); err != nil {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ WORKSPACE ERROR").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Task:      %s\nWorkspace: %s\n\nError: %v",
				pterm.Cyan(t.Name),
				pterm.Gray(session.Workdir),
				err,
			)
		pterm.Println()

		slog.Error("Failed to create workspace tarball",
			"task", t.Name,
			"error", err)
		return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to create workspace tarball: %w", err)}
	}

	// Generate a script compatible with agent execution (without delegate_to)
	agentScript := tr.generateAgentScript(t, groupName)

	pterm.Info.Printfln("📤 Sending task to agent...")

	// Send the task and workspace to the agent
	r, err := c.ExecuteTask(ctx, &pb.ExecuteTaskRequest{
		TaskName:  t.Name,
		TaskGroup: groupName,
		LuaScript: agentScript,
		Workspace: buf.Bytes(),
		User:      t.User,
	})
	if err != nil {
		pterm.Error.Println("═════════════════════════════════════════════════════════════════════════════════════")
		pterm.Error.Printfln("❌ FAILED TO SEND/EXECUTE TASK ON AGENT")
		pterm.Error.Println("═════════════════════════════════════════════════════════════════════════════════════")
		pterm.Error.Printfln("Agent Address: %s", agentAddress)
		pterm.Error.Printfln("Task Name    : %s", t.Name)
		pterm.Error.Printfln("Group Name   : %s", groupName)
		pterm.Error.Printfln("Error        : %v", err)
		pterm.Error.Println("═════════════════════════════════════════════════════════════════════════════════════")

		slog.Error("Failed to send task to agent",
			"agent_address", agentAddress,
			"task", t.Name,
			"error", err)
		return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to execute task on agent %s: %w", agentAddress, err)}
	}

	if !r.GetSuccess() {
		// Parse and display agent error clearly
		agentError := r.GetOutput()

		// Display error in a very visible way using pterm
		pterm.Error.Println("═════════════════════════════════════════════════════════════════════════════════════")
		pterm.Error.Printfln("REMOTE AGENT EXECUTION FAILED")
		pterm.Error.Println("═════════════════════════════════════════════════════════════════════════════════════")
		pterm.Error.Printfln("Agent Address: %s", agentAddress)
		pterm.Error.Println("─────────────────────────────────────────────────────────────────────────────────────")

		// If the error already has formatting from agent, print it directly
		if strings.Contains(agentError, "╔═══") || strings.Contains(agentError, "║") {
			// Agent already formatted the error nicely, print as-is
			fmt.Println(agentError)
		} else {
			// Format it ourselves
			pterm.Error.Println("ERROR OUTPUT FROM AGENT:")
			pterm.Error.Println("─────────────────────────────────────────────────────────────────────────────────────")
			for _, line := range strings.Split(agentError, "\n") {
				if line != "" {
					pterm.Error.Println(line)
				}
			}
		}
		pterm.Error.Println("═════════════════════════════════════════════════════════════════════════════════════")

		// Also log for debugging
		slog.Error("Agent execution failed",
			"task", t.Name,
			"group", groupName,
			"agent", agentAddress,
			"error", agentError)

		// Include the actual error from the agent in the returned error
		return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("agent execution failed on %s:\n%s", agentAddress, agentError)}
	}

	pterm.DefaultBox.
		WithTitle("✅ SUCCESS").
		WithTitleTopLeft().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
		Printfln("Task:  %s\nAgent: %s", pterm.Cyan(t.Name), pterm.Yellow(agentAddress))
	pterm.Println()

	// Extract the updated workspace
	if err := extractTar(bytes.NewReader(r.GetWorkspace()), session.Workdir); err != nil {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ WORKSPACE EXTRACTION FAILED").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Agent:     %s\nTask:      %s\nWorkspace: %s\n\nError: %v",
				pterm.Yellow(agentAddress),
				pterm.Cyan(t.Name),
				pterm.Gray(session.Workdir),
				err,
			)
		pterm.Println()

		slog.Error("Failed to extract updated workspace from agent",
			"agent_address", agentAddress,
			"task", t.Name,
			"error", err)
		return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to extract updated workspace from agent %s: %w", agentAddress, err)}
	}

	pterm.Info.Printfln("📥 Workspace synchronized")
	return nil
}

// executeLocally handles execution of a task locally using Lua
func (tr *TaskRunner) executeLocally(ctx context.Context, t *types.Task, inputFromDependencies *lua.LTable, session *types.SharedSession, groupName string) error {
	L := lua.NewState()
	defer L.Close()
	luainterface.OpenAll(L)

	localInputFromDependencies := luainterface.CopyTable(inputFromDependencies, L)
	t.Output = L.NewTable()

	// Execute pre_exec hook
	if t.PreExec != nil {
		success, msg, _, err := luainterface.ExecuteLuaFunction(L, t.PreExec, t.Params, localInputFromDependencies, 2, ctx)
		if err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("error executing pre_exec hook: %w", err)}
		} else if !success {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("pre-execution hook failed: %s", msg)}
		}
	}

	// Execute command function
	if t.CommandFunc != nil {
		if t.Params == nil {
			t.Params = make(map[string]string)
		}
		t.Params["task_name"] = t.Name
		t.Params["group_name"] = groupName

		// Use task workdir if defined, otherwise use session workdir
		taskWorkdir := session.Workdir
		if t.Workdir != "" {
			taskWorkdir = t.Workdir
			// Create workdir if it doesn't exist
			if err := os.MkdirAll(taskWorkdir, 0755); err != nil {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to create workdir %s: %w", taskWorkdir, err)}
			}
		}
		t.Params["workdir"] = taskWorkdir

		var sessionUD *lua.LUserData
		if session != nil {
			sessionUD = L.NewUserData()
			sessionUD.Value = session
			L.SetMetatable(sessionUD, L.GetTypeMetatable("session"))
		}

		success, msg, outputTable, err := luainterface.ExecuteLuaFunction(L, t.CommandFunc, t.Params, localInputFromDependencies, 3, ctx, sessionUD)
		if err != nil {
			// Execute OnFailure handler if command function has error
			if t.OnFailure != nil {
				tr.executeFailureHandler(L, t, ctx, fmt.Sprintf("error executing command function: %v", err))
			}
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("error executing command function: %w", err)}
		} else if !success {
			// Execute OnFailure handler if command function returns false
			if t.OnFailure != nil {
				tr.executeFailureHandler(L, t, ctx, msg)
			}
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("command function returned failure: %s", msg)}
		} else if outputTable != nil {
			t.Output = outputTable
			// Execute OnSuccess handler if command was successful
			if t.OnSuccess != nil {
				tr.executeSuccessHandler(L, t, ctx, outputTable)
			}
		} else {
			// Execute OnSuccess handler even if no output table
			if t.OnSuccess != nil {
				tr.executeSuccessHandler(L, t, ctx, L.NewTable())
			}
		}
	}

	// Execute post_exec hook
	if t.PostExec != nil {
		var postExecSecondArg lua.LValue = t.Output
		if t.Output == nil {
			postExecSecondArg = L.NewTable()
		}
		success, msg, _, err := luainterface.ExecuteLuaFunction(L, t.PostExec, t.Params, postExecSecondArg, 2, ctx)
		if err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("error executing post_exec hook: %w", err)}
		} else if !success {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("post-execution hook failed: %s", msg)}
		}
	}

	return nil
}

// createTar creates a tarball from the source directory
func createTar(source string, writer io.Writer) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, file)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})
}

// extractTar extracts a tarball to the destination directory
func extractTar(reader io.Reader, dest string) error {
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
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
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
