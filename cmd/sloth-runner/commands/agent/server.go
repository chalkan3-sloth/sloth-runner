package agent

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	agentInternal "github.com/chalkan3-sloth/sloth-runner/internal/agent"
	"github.com/chalkan3-sloth/sloth-runner/internal/core"
	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/creack/pty"
	"github.com/pterm/pterm"
	"github.com/yuin/gopher-lua"
	"google.golang.org/grpc"
)

// agentServer implements the gRPC agent server with optimizations
type agentServer struct {
	pb.UnimplementedAgentServer
	grpcServer *grpc.Server

	// Optimizations: cached metrics
	cachedMetrics     *CachedMetrics
	metricsCache      sync.RWMutex
	lastMetricsUpdate time.Time

	// Additional caches for network and disk info
	cachedNetwork     []NetworkInterfaceInfo
	networkCache      sync.RWMutex
	lastNetworkUpdate time.Time

	cachedDisk        []DiskPartitionInfo
	diskCache         sync.RWMutex
	lastDiskUpdate    time.Time

	// Event worker for sending events to master
	eventWorker       interface{} // Will be *agentInternal.EventWorker, using interface{} to avoid import cycle
	watcherManager    interface{} // Will be *agentInternal.EventWatcherManager, using interface{} to avoid import cycle
}

// CachedMetrics holds cached resource usage data
type CachedMetrics struct {
	CPUPercent      float64
	MemoryPercent   float64
	MemoryUsedBytes uint64
	MemoryTotal     uint64
	DiskPercent     float64
	DiskUsed        uint64
	DiskTotal       uint64
	LoadAvg         [3]float64
	ProcessCount    uint32
	Uptime          uint64
	NetworkRxBytes  uint64
	NetworkTxBytes  uint64
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

// UpdateAgent updates the agent binary to a new version
func (s *agentServer) UpdateAgent(ctx context.Context, in *pb.UpdateAgentRequest) (*pb.UpdateAgentResponse, error) {
	slog.Info("Agent update requested", "version", in.TargetVersion, "force", in.Force, "skip_restart", in.SkipRestart)

	// Get current version
	currentVersion := getCurrentAgentVersion()
	slog.Info("Current agent version", "version", currentVersion)

	// Determine target version
	targetVersion := in.TargetVersion
	if targetVersion == "" || targetVersion == "latest" {
		latest, err := getLatestReleaseVersion()
		if err != nil {
			return &pb.UpdateAgentResponse{
				Success:    false,
				Message:    fmt.Sprintf("Failed to fetch latest version: %v", err),
				OldVersion: currentVersion,
			}, nil
		}
		targetVersion = latest
	}

	slog.Info("Target version determined", "version", targetVersion)

	// Check if already on target version
	if !in.Force && currentVersion == targetVersion {
		return &pb.UpdateAgentResponse{
			Success:    true,
			Message:    "Already on target version",
			OldVersion: currentVersion,
			NewVersion: targetVersion,
		}, nil
	}

	// Download new binary
	newBinaryPath, err := downloadAgentBinary(targetVersion)
	if err != nil {
		return &pb.UpdateAgentResponse{
			Success:    false,
			Message:    fmt.Sprintf("Failed to download new version: %v", err),
			OldVersion: currentVersion,
		}, nil
	}
	defer os.Remove(newBinaryPath)

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return &pb.UpdateAgentResponse{
			Success:    false,
			Message:    fmt.Sprintf("Failed to get current executable path: %v", err),
			OldVersion: currentVersion,
		}, nil
	}

	// Check if running as systemd service
	isSystemd := isRunningAsSystemd()
	slog.Info("Update mode detected", "systemd", isSystemd)

	if isSystemd {
		// Running as systemd - update binary and let systemd restart us
		slog.Info("Updating binary for systemd service")

		// Copy new binary to current location
		if err := copyFile(newBinaryPath, currentExe); err != nil {
			return &pb.UpdateAgentResponse{
				Success:    false,
				Message:    fmt.Sprintf("Failed to update binary: %v", err),
				OldVersion: currentVersion,
			}, nil
		}

		// Make it executable
		if err := os.Chmod(currentExe, 0755); err != nil {
			return &pb.UpdateAgentResponse{
				Success:    false,
				Message:    fmt.Sprintf("Failed to set permissions: %v", err),
				OldVersion: currentVersion,
			}, nil
		}

		slog.Info("Agent binary updated", "old", currentVersion, "new", targetVersion)

		if !in.SkipRestart {
			// Exit and let systemd restart us with the new binary
			slog.Info("Exiting to allow systemd restart...")
			go func() {
				time.Sleep(500 * time.Millisecond)
				os.Exit(0)
			}()
		}
	} else {
		// Not running as systemd - use update script
		updateScript := fmt.Sprintf(`#!/bin/bash
# Agent auto-update script
sleep 2
cp -f %s %s || exit 1
chmod +x %s || exit 1
# Restart the agent with original arguments
%s %s &
# Clean up
rm -f $0
`, newBinaryPath, currentExe, currentExe, currentExe, strings.Join(os.Args[1:], " "))

		scriptPath := "/tmp/sloth-agent-update.sh"
		if err := os.WriteFile(scriptPath, []byte(updateScript), 0755); err != nil {
			return &pb.UpdateAgentResponse{
				Success:    false,
				Message:    fmt.Sprintf("Failed to create update script: %v", err),
				OldVersion: currentVersion,
			}, nil
		}

		slog.Info("Agent binary update prepared", "old", currentVersion, "new", targetVersion)

		// Launch update script and exit
		if !in.SkipRestart {
			slog.Info("Launching update script and exiting...")
			go func() {
				time.Sleep(1 * time.Second)
				// Execute update script in background
				cmd := exec.Command("bash", scriptPath)
				cmd.Start()
				// Exit current process to allow binary replacement
				os.Exit(0)
			}()
		}
	}

	return &pb.UpdateAgentResponse{
		Success:    true,
		Message:    "Agent update initiated - binary will be replaced and agent restarted",
		OldVersion: currentVersion,
		NewVersion: targetVersion,
	}, nil
}

// isRunningAsSystemd detects if the agent is running as a systemd service
func isRunningAsSystemd() bool {
	// Check for INVOCATION_ID environment variable (set by systemd)
	if os.Getenv("INVOCATION_ID") != "" {
		return true
	}

	// Check parent process name
	ppid := os.Getppid()
	if ppid == 1 {
		// Parent is init/systemd (PID 1)
		return true
	}

	// Check if parent process is systemd
	cmdPath := fmt.Sprintf("/proc/%d/comm", ppid)
	data, err := os.ReadFile(cmdPath)
	if err == nil {
		parentName := strings.TrimSpace(string(data))
		if parentName == "systemd" {
			return true
		}
	}

	return false
}

// getCurrentAgentVersion returns the current agent version
func getCurrentAgentVersion() string {
	cmd := exec.Command("sloth-runner", "version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	// Parse version from output (format: "sloth-runner version vX.Y.Z")
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return parts[2]
	}
	return "unknown"
}

// getLatestReleaseVersion fetches the latest release version from GitHub
func getLatestReleaseVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/chalkan3-sloth/sloth-runner/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// downloadAgentBinary downloads the agent binary for the current platform
func downloadAgentBinary(version string) (string, error) {
	// Determine platform and architecture
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Map architectures
	if arch == "amd64" {
		arch = "amd64"
	} else if arch == "arm64" {
		arch = "arm64"
	}

	// Construct download URL
	filename := fmt.Sprintf("sloth-runner_%s_%s_%s.tar.gz", version, platform, arch)
	url := fmt.Sprintf("https://github.com/chalkan3-sloth/sloth-runner/releases/download/%s/%s", version, filename)

	slog.Info("Downloading new agent binary", "url", url)

	// Download the tarball
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "sloth-update-")
	if err != nil {
		return "", err
	}

	// Save tarball to temp file
	tarPath := filepath.Join(tmpDir, filename)
	tarFile, err := os.Create(tarPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	if _, err := io.Copy(tarFile, resp.Body); err != nil {
		tarFile.Close()
		os.RemoveAll(tmpDir)
		return "", err
	}
	tarFile.Close()

	// Extract binary from tarball
	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	// Extract using tar command
	cmd := exec.Command("tar", "-xzf", tarPath, "-C", extractDir)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to extract tarball: %w", err)
	}

	// Find the binary (should be sloth-runner)
	binaryPath := filepath.Join(extractDir, "sloth-runner")
	if _, err := os.Stat(binaryPath); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("binary not found in tarball: %w", err)
	}

	return binaryPath, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}

// restartAgent restarts the agent service
func restartAgent() {
	slog.Info("Restarting agent process...")

	// Try systemctl restart first
	if err := exec.Command("systemctl", "restart", "sloth-runner-agent").Run(); err == nil {
		slog.Info("Agent restarted via systemctl")
		return
	}

	// If systemctl fails, try to restart current process
	currentExe, err := os.Executable()
	if err != nil {
		slog.Error("Failed to get executable path for restart", "error", err)
		return
	}

	// Get current process arguments
	args := os.Args[1:]

	cmd := exec.Command(currentExe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		slog.Error("Failed to restart agent", "error", err)
		return
	}

	slog.Info("Agent restart initiated, exiting current process...")
	os.Exit(0)
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

	// Store EventWatcherManager in Lua state for watcher registration
	if s.watcherManager != nil {
		ud := L.NewUserData()
		ud.Value = s.watcherManager
		L.SetGlobal("__WATCHER_MANAGER__", ud)
		slog.Info("EventWatcherManager attached to Lua state for watcher registration")
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

	// Extract and register watchers after task execution
	slog.Info("Checking for watchers to register", "watcherManager_nil", s.watcherManager == nil)
	if s.watcherManager != nil {
		watchers, extractErr := luainterface.GetRegisteredWatchers(L)
		slog.Info("Extracted watchers from Lua state",
			"count", len(watchers),
			"error", extractErr,
			"has_watchers", len(watchers) > 0)

		if extractErr == nil && len(watchers) > 0 {
			slog.Info("Processing watchers for registration", "count", len(watchers))

			// Type assert watcherManager to its proper type
			if watcherMgr, ok := s.watcherManager.(*agentInternal.EventWatcherManager); ok {
				for _, watcherConfigMap := range watchers {
					// Convert map to WatcherConfig struct
					config, convertErr := convertMapToWatcherConfig(watcherConfigMap)
					if convertErr != nil {
						slog.Warn("Failed to convert watcher config", "error", convertErr)
						continue
					}

					if regErr := watcherMgr.RegisterWatcher(config); regErr == nil {
						slog.Info("Successfully registered watcher with manager",
							"watcher_id", config.ID,
							"type", config.Type)
						pterm.Success.Printf("✓ Watcher %s started on agent\n", config.ID)
					} else {
						slog.Warn("Failed to register watcher with manager",
							"watcher_id", config.ID,
							"type", config.Type,
							"error", regErr)
					}
				}
			} else {
				slog.Warn("watcherManager type assertion failed - watchers not registered")
			}
		} else if extractErr != nil {
			slog.Warn("Failed to extract watchers from Lua state", "error", extractErr)
		} else {
			slog.Info("No watchers found in Lua state")
		}
	} else {
		slog.Warn("watcherManager is nil - cannot register watchers")
	}

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

// convertMapToWatcherConfig converts a map[string]interface{} from Lua to WatcherConfig
func convertMapToWatcherConfig(m map[string]interface{}) (agentInternal.WatcherConfig, error) {
	config := agentInternal.WatcherConfig{}

	// Extract basic fields
	if id, ok := m["id"].(string); ok {
		config.ID = id
	}
	if typeStr, ok := m["type"].(string); ok {
		config.Type = agentInternal.WatcherType(typeStr)
	}

	// Extract conditions (when field)
	if when, ok := m["when"].([]interface{}); ok {
		for _, c := range when {
			if condStr, ok := c.(string); ok {
				config.Conditions = append(config.Conditions, agentInternal.EventCondition(condStr))
			}
		}
	}

	// Extract file/directory fields
	if path, ok := m["file_path"].(string); ok {
		config.FilePath = path
	}
	if recursive, ok := m["recursive"].(bool); ok {
		config.Recursive = recursive
	}
	if checkHash, ok := m["check_hash"].(bool); ok {
		config.CheckHash = checkHash
	}
	if pattern, ok := m["pattern"].(string); ok {
		config.Pattern = pattern
	}

	// Extract process fields
	if processName, ok := m["process_name"].(string); ok {
		config.ProcessName = processName
	}
	if pid, ok := m["pid"].(float64); ok {
		config.PID = int(pid)
	}

	// Extract port/network fields
	if port, ok := m["port"].(float64); ok {
		config.Port = int(port)
	}
	if protocol, ok := m["protocol"].(string); ok {
		config.Protocol = protocol
	}

	// Extract interval (parse duration string)
	if interval, ok := m["interval"].(string); ok {
		if d, err := time.ParseDuration(interval); err == nil {
			config.Interval = d
		}
	}

	return config, nil
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

// GetResourceUsage returns current resource usage of the agent (optimized with caching)
func (s *agentServer) GetResourceUsage(ctx context.Context, in *pb.ResourceUsageRequest) (*pb.ResourceUsageResponse, error) {
	// Check cache first (60 second TTL - increased from 30s for better performance)
	s.metricsCache.RLock()
	if time.Since(s.lastMetricsUpdate) < 60*time.Second && s.cachedMetrics != nil {
		cached := s.cachedMetrics
		s.metricsCache.RUnlock()

		return &pb.ResourceUsageResponse{
			CpuPercent:       cached.CPUPercent,
			MemoryPercent:    cached.MemoryPercent,
			MemoryUsedBytes:  cached.MemoryUsedBytes,
			MemoryTotalBytes: cached.MemoryTotal,
			DiskPercent:      cached.DiskPercent,
			DiskUsedBytes:    cached.DiskUsed,
			DiskTotalBytes:   cached.DiskTotal,
			ProcessCount:     cached.ProcessCount,
			LoadAvg_1Min:     cached.LoadAvg[0],
			LoadAvg_5Min:     cached.LoadAvg[1],
			LoadAvg_15Min:    cached.LoadAvg[2],
			UptimeSeconds:    cached.Uptime,
			NetworkRxBytes:   cached.NetworkRxBytes,
			NetworkTxBytes:   cached.NetworkTxBytes,
		}, nil
	}
	s.metricsCache.RUnlock()

	// Update cache async (non-blocking for subsequent requests)
	go s.updateMetricsCache()

	// For first request or expired cache, get fresh data
	cpuPercent := getCPUPercent()
	memInfo, _ := getMemoryInfo()
	diskInfo, _ := getDiskUsage("/")
	loadAvg := getLoadAverage()
	processCount := getProcessCount()
	uptime := getSystemUptime()
	networkRx, networkTx := getNetworkBytes()

	// Calculate memory percent
	memPercent := float64(memInfo.Used) / float64(memInfo.Total) * 100

	response := &pb.ResourceUsageResponse{
		CpuPercent:       cpuPercent,
		MemoryPercent:    memPercent,
		MemoryUsedBytes:  memInfo.Used,
		MemoryTotalBytes: memInfo.Total,
		DiskPercent:      diskInfo.UsedPercent,
		DiskUsedBytes:    diskInfo.Used,
		DiskTotalBytes:   diskInfo.Total,
		ProcessCount:     uint32(processCount),
		LoadAvg_1Min:     loadAvg[0],
		LoadAvg_5Min:     loadAvg[1],
		LoadAvg_15Min:    loadAvg[2],
		UptimeSeconds:    uptime,
		NetworkRxBytes:   networkRx,
		NetworkTxBytes:   networkTx,
	}

	return response, nil
}

// updateMetricsCache updates the metrics cache asynchronously
func (s *agentServer) updateMetricsCache() {
	s.metricsCache.Lock()
	defer s.metricsCache.Unlock()

	// Get memory info using platform-independent function
	memInfo, _ := getMemoryInfo()

	// Update cache
	s.cachedMetrics = &CachedMetrics{
		CPUPercent:      getCPUPercent(),
		MemoryPercent:   float64(memInfo.Used) / float64(memInfo.Total) * 100,
		MemoryUsedBytes: memInfo.Used,
		MemoryTotal:     memInfo.Total,
		ProcessCount:    uint32(getProcessCount()),
		Uptime:          getSystemUptime(),
	}

	// Get disk (optional, more expensive)
	if diskInfo, err := getDiskUsage("/"); err == nil {
		s.cachedMetrics.DiskPercent = diskInfo.UsedPercent
		s.cachedMetrics.DiskUsed = diskInfo.Used
		s.cachedMetrics.DiskTotal = diskInfo.Total
	}

	// Get load average
	loadAvg := getLoadAverage()
	s.cachedMetrics.LoadAvg = [3]float64{loadAvg[0], loadAvg[1], loadAvg[2]}

	// Get network bytes
	networkRx, networkTx := getNetworkBytes()
	s.cachedMetrics.NetworkRxBytes = networkRx
	s.cachedMetrics.NetworkTxBytes = networkTx

	s.lastMetricsUpdate = time.Now()
}

// GetProcessList returns list of running processes
func (s *agentServer) GetProcessList(ctx context.Context, in *pb.ProcessListRequest) (*pb.ProcessListResponse, error) {
	processes, err := getProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	// Pre-allocate exact size to avoid re-allocations
	pbProcesses := make([]*pb.ProcessInfo, len(processes))
	for i, p := range processes {
		pbProcesses[i] = &pb.ProcessInfo{
			Pid:           int32(p.PID),
			Name:          p.Name,
			Status:        p.Status,
			CpuPercent:    p.CPUPercent,
			MemoryPercent: p.MemoryPercent,
			MemoryBytes:   p.MemoryBytes,
			User:          p.User,
			Command:       p.Command,
			StartedAt:     p.StartedAt,
		}
	}

	return &pb.ProcessListResponse{Processes: pbProcesses}, nil
}

// GetNetworkInfo returns network interface information (with 120s cache)
func (s *agentServer) GetNetworkInfo(ctx context.Context, in *pb.NetworkInfoRequest) (*pb.NetworkInfoResponse, error) {
	// Check cache first (120 second TTL - network info changes rarely, increased from 60s)
	s.networkCache.RLock()
	if s.cachedNetwork != nil && time.Since(s.lastNetworkUpdate) < 120*time.Second {
		interfaces := s.cachedNetwork
		s.networkCache.RUnlock()

		hostname, _ := os.Hostname()
		pbInterfaces := make([]*pb.NetworkInterface, 0, len(interfaces))
		for _, iface := range interfaces {
			pbInterfaces = append(pbInterfaces, &pb.NetworkInterface{
				Name:        iface.Name,
				IpAddresses: iface.IPAddresses,
				MacAddress:  iface.MACAddress,
				BytesSent:   iface.BytesSent,
				BytesRecv:   iface.BytesRecv,
				IsUp:        iface.IsUp,
			})
		}

		return &pb.NetworkInfoResponse{
			Interfaces: pbInterfaces,
			Hostname:   hostname,
		}, nil
	}
	s.networkCache.RUnlock()

	// Get fresh data
	interfaces, err := getNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Update cache
	s.networkCache.Lock()
	s.cachedNetwork = interfaces
	s.lastNetworkUpdate = time.Now()
	s.networkCache.Unlock()

	hostname, _ := os.Hostname()

	pbInterfaces := make([]*pb.NetworkInterface, 0, len(interfaces))
	for _, iface := range interfaces {
		pbInterfaces = append(pbInterfaces, &pb.NetworkInterface{
			Name:        iface.Name,
			IpAddresses: iface.IPAddresses,
			MacAddress:  iface.MACAddress,
			BytesSent:   iface.BytesSent,
			BytesRecv:   iface.BytesRecv,
			IsUp:        iface.IsUp,
		})
	}

	return &pb.NetworkInfoResponse{
		Interfaces: pbInterfaces,
		Hostname:   hostname,
	}, nil
}

// GetDiskInfo returns disk partition information (with 300s cache)
func (s *agentServer) GetDiskInfo(ctx context.Context, in *pb.DiskInfoRequest) (*pb.DiskInfoResponse, error) {
	// Check cache first (300 second TTL - disk info changes very rarely, increased from 60s)
	s.diskCache.RLock()
	if s.cachedDisk != nil && time.Since(s.lastDiskUpdate) < 300*time.Second {
		partitions := s.cachedDisk
		s.diskCache.RUnlock()

		pbPartitions := make([]*pb.DiskPartition, 0, len(partitions))
		var totalRead, totalWrite uint64

		for _, part := range partitions {
			pbPartitions = append(pbPartitions, &pb.DiskPartition{
				Device:     part.Device,
				Mountpoint: part.Mountpoint,
				Fstype:     part.FSType,
				TotalBytes: part.TotalBytes,
				UsedBytes:  part.UsedBytes,
				FreeBytes:  part.FreeBytes,
				Percent:    part.Percent,
			})
			totalRead += part.IOReadBytes
			totalWrite += part.IOWriteBytes
		}

		return &pb.DiskInfoResponse{
			Partitions:        pbPartitions,
			TotalIoReadBytes:  totalRead,
			TotalIoWriteBytes: totalWrite,
		}, nil
	}
	s.diskCache.RUnlock()

	// Get fresh data
	partitions, err := getDiskPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	// Update cache
	s.diskCache.Lock()
	s.cachedDisk = partitions
	s.lastDiskUpdate = time.Now()
	s.diskCache.Unlock()

	pbPartitions := make([]*pb.DiskPartition, 0, len(partitions))
	var totalRead, totalWrite uint64

	for _, part := range partitions {
		pbPartitions = append(pbPartitions, &pb.DiskPartition{
			Device:     part.Device,
			Mountpoint: part.Mountpoint,
			Fstype:     part.FSType,
			TotalBytes: part.TotalBytes,
			UsedBytes:  part.UsedBytes,
			FreeBytes:  part.FreeBytes,
			Percent:    part.Percent,
		})
		totalRead += part.IOReadBytes
		totalWrite += part.IOWriteBytes
	}

	return &pb.DiskInfoResponse{
		Partitions:        pbPartitions,
		TotalIoReadBytes:  totalRead,
		TotalIoWriteBytes: totalWrite,
	}, nil
}

// StreamLogs streams agent logs
func (s *agentServer) StreamLogs(in *pb.StreamLogsRequest, stream pb.Agent_StreamLogsServer) error {
	// For now, stream system journal logs
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	logChan := make(chan LogEntry, 100)
	go collectLogs(logChan)
	
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case entry := <-logChan:
			err := stream.Send(&pb.LogEntry{
				Timestamp: entry.Timestamp,
				Level:     entry.Level,
				Message:   entry.Message,
			})
			if err != nil {
				return err
			}
		}
	}
}

// StreamMetrics streams real-time metrics
func (s *agentServer) StreamMetrics(in *pb.StreamMetricsRequest, stream pb.Agent_StreamMetricsServer) error {
	// Increased from 2s to 5s to reduce CPU/bandwidth usage
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			// Get current metrics
			cpuPercent := getCPUPercent()
			memInfo, _ := getMemoryInfo()
			diskInfo, _ := getDiskUsage("/")
			
			memPercent := float64(memInfo.Used) / float64(memInfo.Total) * 100
			
			err := stream.Send(&pb.MetricsData{
				Timestamp:     time.Now().Unix(),
				CpuPercent:    cpuPercent,
				MemoryPercent: memPercent,
				DiskPercent:   diskInfo.UsedPercent,
			})
			if err != nil {
				return err
			}
		}
	}
}

// RestartService restarts the agent service
func (s *agentServer) RestartService(ctx context.Context, in *pb.RestartServiceRequest) (*pb.RestartServiceResponse, error) {
	slog.Info("Agent restart requested")

	go func() {
		time.Sleep(1 * time.Second)
		// Get current executable and args
		exe, _ := os.Executable()
		args := os.Args[1:]

		// Start new instance
		cmd := exec.Command(exe, args...)
		cmd.Start()

		// Stop current instance
		s.grpcServer.GracefulStop()
		os.Exit(0)
	}()

	return &pb.RestartServiceResponse{
		Success: true,
		Message: "Agent restart initiated",
	}, nil
}

// GetDetailedMetrics returns detailed system metrics
func (s *agentServer) GetDetailedMetrics(ctx context.Context, in *pb.DetailedMetricsRequest) (*pb.DetailedMetricsResponse, error) {
	// Get CPU details
	data, _ := os.ReadFile("/proc/cpuinfo")
	cpuInfo := string(data)

	// Count CPU cores
	coreCount := int32(0)
	modelName := ""
	cpuMhz := 0.0
	for _, line := range strings.Split(cpuInfo, "\n") {
		if strings.HasPrefix(line, "processor") {
			coreCount++
		}
		if strings.HasPrefix(line, "model name") && modelName == "" {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				modelName = strings.TrimSpace(parts[1])
			}
		}
		if strings.HasPrefix(line, "cpu MHz") && cpuMhz == 0 {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				cpuMhz, _ = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			}
		}
	}

	// Get per-core CPU usage from /proc/stat
	perCoreUsage := []float64{}
	statData, _ := os.ReadFile("/proc/stat")
	for _, line := range strings.Split(string(statData), "\n") {
		if strings.HasPrefix(line, "cpu") && !strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				var total, idle uint64
				for i := 1; i < len(fields); i++ {
					val, _ := strconv.ParseUint(fields[i], 10, 64)
					total += val
					if i == 4 {
						idle = val
					}
				}
				if total > 0 {
					usage := float64(total-idle) / float64(total) * 100.0
					perCoreUsage = append(perCoreUsage, usage)
				}
			}
		}
	}

	// Parse CPU times
	firstLine := strings.Split(string(statData), "\n")[0]
	cpuFields := strings.Fields(firstLine)
	userTime, _ := strconv.ParseFloat(cpuFields[1], 64)
	systemTime, _ := strconv.ParseFloat(cpuFields[3], 64)
	idleTime, _ := strconv.ParseFloat(cpuFields[4], 64)
	iowaitTime, _ := strconv.ParseFloat(cpuFields[5], 64)

	cpuDetail := &pb.CPUDetail{
		CoreCount:    coreCount,
		PerCoreUsage: perCoreUsage,
		UserTime:     userTime,
		SystemTime:   systemTime,
		IdleTime:     idleTime,
		IowaitTime:   iowaitTime,
		ModelName:    modelName,
		Mhz:          cpuMhz,
	}

	// Get memory details
	memInfo, _ := getMemoryInfo()
	memData, _ := os.ReadFile("/proc/meminfo")
	var cached, buffers, swapTotal, swapFree uint64

	for _, line := range strings.Split(string(memData), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024 // KB to bytes

		if strings.HasPrefix(fields[0], "Cached:") {
			cached = value
		} else if strings.HasPrefix(fields[0], "Buffers:") {
			buffers = value
		} else if strings.HasPrefix(fields[0], "SwapTotal:") {
			swapTotal = value
		} else if strings.HasPrefix(fields[0], "SwapFree:") {
			swapFree = value
		}
	}

	swapUsed := swapTotal - swapFree
	swapPercent := 0.0
	if swapTotal > 0 {
		swapPercent = float64(swapUsed) / float64(swapTotal) * 100.0
	}

	memoryDetail := &pb.MemoryDetail{
		TotalBytes:     memInfo.Total,
		AvailableBytes: memInfo.Free,
		UsedBytes:      memInfo.Used,
		FreeBytes:      memInfo.Free,
		CachedBytes:    cached,
		BuffersBytes:   buffers,
		SwapTotalBytes: swapTotal,
		SwapUsedBytes:  swapUsed,
		SwapFreeBytes:  swapFree,
		Percent:        float64(memInfo.Used) / float64(memInfo.Total) * 100.0,
		SwapPercent:    swapPercent,
	}

	// Get disk details
	partitions, _ := getDiskPartitions()
	pbPartitions := make([]*pb.DiskPartition, 0, len(partitions))
	var totalRead, totalWrite uint64

	for _, part := range partitions {
		pbPartitions = append(pbPartitions, &pb.DiskPartition{
			Device:     part.Device,
			Mountpoint: part.Mountpoint,
			Fstype:     part.FSType,
			TotalBytes: part.TotalBytes,
			UsedBytes:  part.UsedBytes,
			FreeBytes:  part.FreeBytes,
			Percent:    part.Percent,
		})
		totalRead += part.IOReadBytes
		totalWrite += part.IOWriteBytes
	}

	diskDetail := &pb.DiskDetail{
		Partitions:      pbPartitions,
		ReadBytesTotal:  totalRead,
		WriteBytesTotal: totalWrite,
		ReadCount:       0,
		WriteCount:      0,
		ReadTimeMs:      0,
		WriteTimeMs:     0,
	}

	// Get network details
	interfaces, _ := getNetworkInterfaces()
	pbInterfaces := make([]*pb.NetworkInterface, 0, len(interfaces))
	var bytesSentTotal, bytesRecvTotal uint64

	for _, iface := range interfaces {
		pbInterfaces = append(pbInterfaces, &pb.NetworkInterface{
			Name:        iface.Name,
			IpAddresses: iface.IPAddresses,
			MacAddress:  iface.MACAddress,
			BytesSent:   iface.BytesSent,
			BytesRecv:   iface.BytesRecv,
			IsUp:        iface.IsUp,
		})
		bytesSentTotal += iface.BytesSent
		bytesRecvTotal += iface.BytesRecv
	}

	networkDetail := &pb.NetworkDetail{
		Interfaces:        pbInterfaces,
		BytesSentTotal:    bytesSentTotal,
		BytesRecvTotal:    bytesRecvTotal,
		PacketsSentTotal:  0,
		PacketsRecvTotal:  0,
		ErrorsIn:          0,
		ErrorsOut:         0,
		DropsIn:           0,
		DropsOut:          0,
		ActiveConnections: 0,
		ListeningPorts:    0,
	}

	// Get load average and other metrics
	loadAvg := getLoadAverage()
	uptime := getSystemUptime()
	processCount := getProcessCount()

	// Get kernel and OS version
	kernelData, _ := os.ReadFile("/proc/version")
	kernelVersion := strings.TrimSpace(string(kernelData))

	osData, _ := os.ReadFile("/etc/os-release")
	osVersion := "Unknown"
	for _, line := range strings.Split(string(osData), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			osVersion = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			break
		}
	}

	return &pb.DetailedMetricsResponse{
		Timestamp:     time.Now().Unix(),
		Cpu:           cpuDetail,
		Memory:        memoryDetail,
		Disk:          diskDetail,
		Network:       networkDetail,
		LoadAvg_1Min:  loadAvg[0],
		LoadAvg_5Min:  loadAvg[1],
		LoadAvg_15Min: loadAvg[2],
		UptimeSeconds: uptime,
		ProcessCount:  int32(processCount),
		ThreadCount:   0,
		KernelVersion: kernelVersion,
		OsVersion:     osVersion,
	}, nil
}

// GetRecentLogs returns recent system logs
func (s *agentServer) GetRecentLogs(ctx context.Context, in *pb.RecentLogsRequest) (*pb.RecentLogsResponse, error) {
	maxLines := in.MaxLines
	if maxLines == 0 {
		maxLines = 100
	}

	// Build journalctl command
	args := []string{"-n", strconv.Itoa(int(maxLines)), "-o", "json"}

	// Add level filter
	if in.LevelFilter != "" {
		args = append(args, "-p", in.LevelFilter)
	}

	// Add source filter
	if in.SourceFilter != "" {
		args = append(args, "-u", in.SourceFilter)
	}

	// Add since timestamp
	if in.SinceTimestamp > 0 {
		args = append(args, "--since", fmt.Sprintf("@%d", in.SinceTimestamp))
	}

	cmd := exec.Command("journalctl", args...)
	output, err := cmd.Output()
	if err != nil {
		// Fallback to reading /var/log/syslog
		data, err := os.ReadFile("/var/log/syslog")
		if err != nil {
			return &pb.RecentLogsResponse{
				Logs:       []*pb.LogEntry{},
				TotalCount: 0,
				HasMore:    false,
			}, nil
		}

		lines := strings.Split(string(data), "\n")
		logs := []*pb.LogEntry{}
		count := 0

		for i := len(lines) - 1; i >= 0 && count < int(maxLines); i-- {
			if lines[i] == "" {
				continue
			}
			logs = append([]*pb.LogEntry{{
				Timestamp: time.Now().Unix(),
				Level:     "INFO",
				Message:   lines[i],
				Source:    "syslog",
			}}, logs...)
			count++
		}

		return &pb.RecentLogsResponse{
			Logs:       logs,
			TotalCount: int32(count),
			HasMore:    len(lines) > int(maxLines),
		}, nil
	}

	// Parse journalctl output
	logs := []*pb.LogEntry{}
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}

		// Simple parsing (journalctl -o json returns one JSON per line)
		logs = append(logs, &pb.LogEntry{
			Timestamp: time.Now().Unix(),
			Level:     "INFO",
			Message:   line,
			Source:    "journalctl",
		})
	}

	return &pb.RecentLogsResponse{
		Logs:       logs,
		TotalCount: int32(len(logs)),
		HasMore:    false,
	}, nil
}

// GetActiveConnections returns active network connections
func (s *agentServer) GetActiveConnections(ctx context.Context, in *pb.ConnectionsRequest) (*pb.ConnectionsResponse, error) {
	connections := []*pb.ConnectionInfo{}
	totalEstablished := int32(0)
	totalListening := int32(0)
	totalTimeWait := int32(0)

	// Read TCP connections
	tcpData, err := os.ReadFile("/proc/net/tcp")
	if err == nil {
		lines := strings.Split(string(tcpData), "\n")
		for i, line := range lines {
			if i == 0 || line == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}

			// Parse local address and port
			localParts := strings.Split(fields[1], ":")
			remoteParts := strings.Split(fields[2], ":")

			if len(localParts) != 2 || len(remoteParts) != 2 {
				continue
			}

			localPort, _ := strconv.ParseInt(localParts[1], 16, 64)
			remotePort, _ := strconv.ParseInt(remoteParts[1], 16, 64)

			// Parse state
			stateHex, _ := strconv.ParseInt(fields[3], 16, 64)
			state := "UNKNOWN"
			switch stateHex {
			case 1:
				state = "ESTABLISHED"
				totalEstablished++
			case 2:
				state = "SYN_SENT"
			case 3:
				state = "SYN_RECV"
			case 6:
				state = "TIME_WAIT"
				totalTimeWait++
			case 10:
				state = "LISTEN"
				totalListening++
			}

			// Apply state filter
			if in.StateFilter != "" && state != in.StateFilter {
				continue
			}

			// Parse local address (hex IP)
			localIP := parseHexIP(localParts[0])
			remoteIP := parseHexIP(remoteParts[0])

			// Skip localhost if not included
			if !in.IncludeLocal && (localIP == "127.0.0.1" || remoteIP == "127.0.0.1") {
				continue
			}

			connections = append(connections, &pb.ConnectionInfo{
				LocalAddr:  localIP,
				LocalPort:  uint32(localPort),
				RemoteAddr: remoteIP,
				RemotePort: uint32(remotePort),
				State:      state,
				Pid:        0,
			})
		}
	}

	// Read UDP connections
	udpData, err := os.ReadFile("/proc/net/udp")
	if err == nil {
		lines := strings.Split(string(udpData), "\n")
		for i, line := range lines {
			if i == 0 || line == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}

			localParts := strings.Split(fields[1], ":")
			if len(localParts) != 2 {
				continue
			}

			localPort, _ := strconv.ParseInt(localParts[1], 16, 64)
			localIP := parseHexIP(localParts[0])

			if !in.IncludeLocal && localIP == "127.0.0.1" {
				continue
			}

			connections = append(connections, &pb.ConnectionInfo{
				LocalAddr:  localIP,
				LocalPort:  uint32(localPort),
				RemoteAddr: "0.0.0.0",
				RemotePort: 0,
				State:      "UDP",
			})
		}
	}

	return &pb.ConnectionsResponse{
		Connections:      connections,
		TotalEstablished: totalEstablished,
		TotalListening:   totalListening,
		TotalTimeWait:    totalTimeWait,
		TotalAll:         int32(len(connections)),
	}, nil
}

// parseHexIP converts hex IP address to dotted decimal
func parseHexIP(hexIP string) string {
	if len(hexIP) != 8 {
		return "0.0.0.0"
	}

	// Parse as little-endian
	ip := make([]byte, 4)
	for i := 0; i < 4; i++ {
		val, _ := strconv.ParseUint(hexIP[i*2:(i+1)*2], 16, 8)
		ip[3-i] = byte(val)
	}

	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// GetSystemErrors returns system errors from logs
func (s *agentServer) GetSystemErrors(ctx context.Context, in *pb.SystemErrorsRequest) (*pb.SystemErrorsResponse, error) {
	maxErrors := in.MaxErrors
	if maxErrors == 0 {
		maxErrors = 50
	}

	errors := []*pb.SystemError{}
	totalErrors := int32(0)
	totalWarnings := int32(0)
	errorCounts := make(map[string]int32)

	// Try to get errors from journalctl
	args := []string{"-p", "err", "-n", strconv.Itoa(int(maxErrors)), "-o", "short"}
	if in.SinceTimestamp > 0 {
		args = append(args, "--since", fmt.Sprintf("@%d", in.SinceTimestamp))
	}

	cmd := exec.Command("journalctl", args...)
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}

			severity := "error"
			if strings.Contains(strings.ToLower(line), "warning") {
				severity = "warning"
				totalWarnings++
			} else {
				totalErrors++
			}

			errors = append(errors, &pb.SystemError{
				Timestamp:       time.Now().Unix(),
				Severity:        severity,
				Source:          "journalctl",
				Message:         line,
				OccurrenceCount: 1,
			})

			// Count error types
			errorCounts[line]++
		}
	}

	// Try to get kernel errors from dmesg
	dmesgCmd := exec.Command("dmesg", "-l", "err,warn", "-T")
	dmesgOutput, err := dmesgCmd.Output()
	if err == nil {
		lines := strings.Split(string(dmesgOutput), "\n")
		count := 0
		for _, line := range lines {
			if line == "" || count >= int(maxErrors) {
				continue
			}

			severity := "error"
			if strings.Contains(strings.ToLower(line), "warn") {
				severity = "warning"
				totalWarnings++
			} else {
				totalErrors++
			}

			errors = append(errors, &pb.SystemError{
				Timestamp:       time.Now().Unix(),
				Severity:        severity,
				Source:          "dmesg",
				Message:         line,
				OccurrenceCount: 1,
			})

			errorCounts[line]++
			count++
		}
	}

	// Include warnings if requested
	if in.IncludeWarnings {
		warnArgs := []string{"-p", "warning", "-n", "20", "-o", "short"}
		if in.SinceTimestamp > 0 {
			warnArgs = append(warnArgs, "--since", fmt.Sprintf("@%d", in.SinceTimestamp))
		}

		warnCmd := exec.Command("journalctl", warnArgs...)
		warnOutput, _ := warnCmd.Output()

		for _, line := range strings.Split(string(warnOutput), "\n") {
			if line == "" {
				continue
			}

			errors = append(errors, &pb.SystemError{
				Timestamp:       time.Now().Unix(),
				Severity:        "warning",
				Source:          "journalctl",
				Message:         line,
				OccurrenceCount: 1,
			})
			totalWarnings++
		}
	}

	// Find most common error
	mostCommon := ""
	maxCount := int32(0)
	for msg, count := range errorCounts {
		if count > maxCount {
			maxCount = count
			mostCommon = msg
		}
	}

	return &pb.SystemErrorsResponse{
		Errors:           errors,
		TotalErrors:      totalErrors,
		TotalWarnings:    totalWarnings,
		MostCommonError:  mostCommon,
	}, nil
}

// GetPerformanceHistory returns historical performance metrics
func (s *agentServer) GetPerformanceHistory(ctx context.Context, in *pb.PerformanceHistoryRequest) (*pb.PerformanceHistoryResponse, error) {
	duration := in.DurationMinutes
	if duration == 0 {
		duration = 60
	}

	dataPoints := in.DataPoints
	if dataPoints == 0 {
		dataPoints = 30
	}

	// Calculate interval between samples
	intervalSeconds := (duration * 60) / dataPoints

	snapshots := []*pb.PerformanceSnapshot{}
	var sumCPU, sumMemory, sumDisk, sumLoad float64
	var minCPU, minMemory, minDisk = 100.0, 100.0, 100.0
	var maxCPU, maxMemory, maxDisk float64

	// Collect samples
	for i := int32(0); i < dataPoints; i++ {
		cpuPercent := getCPUPercent()
		memInfo, _ := getMemoryInfo()
		diskInfo, _ := getDiskUsage("/")
		loadAvg := getLoadAverage()
		processCount := getProcessCount()

		memPercent := float64(memInfo.Used) / float64(memInfo.Total) * 100.0

		snapshot := &pb.PerformanceSnapshot{
			Timestamp:             time.Now().Unix() - int64((dataPoints-i-1)*intervalSeconds),
			CpuPercent:            cpuPercent,
			MemoryPercent:         memPercent,
			DiskPercent:           diskInfo.UsedPercent,
			NetworkThroughputMbps: 0.0,
			LoadAvg:               loadAvg[0],
			ActiveConnections:     0,
			ProcessCount:          int32(processCount),
		}

		snapshots = append(snapshots, snapshot)

		// Update aggregates
		sumCPU += cpuPercent
		sumMemory += memPercent
		sumDisk += diskInfo.UsedPercent
		sumLoad += loadAvg[0]

		if cpuPercent < minCPU {
			minCPU = cpuPercent
		}
		if cpuPercent > maxCPU {
			maxCPU = cpuPercent
		}
		if memPercent < minMemory {
			minMemory = memPercent
		}
		if memPercent > maxMemory {
			maxMemory = memPercent
		}
		if diskInfo.UsedPercent < minDisk {
			minDisk = diskInfo.UsedPercent
		}
		if diskInfo.UsedPercent > maxDisk {
			maxDisk = diskInfo.UsedPercent
		}

		// Sleep between samples (only if not last iteration)
		if i < dataPoints-1 {
			time.Sleep(time.Duration(intervalSeconds) * time.Second)
		}
	}

	count := float64(len(snapshots))

	return &pb.PerformanceHistoryResponse{
		Snapshots: snapshots,
		Avg: &pb.PerformanceSnapshot{
			CpuPercent:    sumCPU / count,
			MemoryPercent: sumMemory / count,
			DiskPercent:   sumDisk / count,
			LoadAvg:       sumLoad / count,
		},
		Min: &pb.PerformanceSnapshot{
			CpuPercent:    minCPU,
			MemoryPercent: minMemory,
			DiskPercent:   minDisk,
		},
		Max: &pb.PerformanceSnapshot{
			CpuPercent:    maxCPU,
			MemoryPercent: maxMemory,
			DiskPercent:   maxDisk,
		},
	}, nil
}

// DiagnoseHealth performs health diagnostic and calculates health score
func (s *agentServer) DiagnoseHealth(ctx context.Context, in *pb.HealthDiagnosticRequest) (*pb.HealthDiagnosticResponse, error) {
	issues := []*pb.HealthIssue{}
	healthScore := int32(100)
	totalWarnings := int32(0)
	totalErrors := int32(0)

	// Check CPU usage
	cpuPercent := getCPUPercent()
	if cpuPercent > 90 {
		healthScore -= 20
		totalErrors++
		issues = append(issues, &pb.HealthIssue{
			Category:     "cpu",
			Severity:     "critical",
			Description:  "CPU usage is critically high",
			CurrentValue: fmt.Sprintf("%.2f%%", cpuPercent),
			Threshold:    "90%",
			Suggestions:  []string{"Identify and terminate resource-intensive processes", "Consider scaling up CPU resources", "Check for runaway processes"},
			AutoFixable:  false,
		})
	} else if cpuPercent > 75 {
		healthScore -= 10
		totalWarnings++
		issues = append(issues, &pb.HealthIssue{
			Category:     "cpu",
			Severity:     "warning",
			Description:  "CPU usage is elevated",
			CurrentValue: fmt.Sprintf("%.2f%%", cpuPercent),
			Threshold:    "75%",
			Suggestions:  []string{"Monitor CPU-intensive processes", "Consider load balancing"},
			AutoFixable:  false,
		})
	}

	// Check memory usage
	memInfo, _ := getMemoryInfo()
	memPercent := float64(memInfo.Used) / float64(memInfo.Total) * 100.0
	if memPercent > 90 {
		healthScore -= 20
		totalErrors++
		issues = append(issues, &pb.HealthIssue{
			Category:     "memory",
			Severity:     "critical",
			Description:  "Memory usage is critically high",
			CurrentValue: fmt.Sprintf("%.2f%%", memPercent),
			Threshold:    "90%",
			Suggestions:  []string{"Free up memory by closing unused applications", "Increase swap space", "Add more RAM", "Check for memory leaks"},
			AutoFixable:  false,
		})
	} else if memPercent > 80 {
		healthScore -= 10
		totalWarnings++
		issues = append(issues, &pb.HealthIssue{
			Category:     "memory",
			Severity:     "warning",
			Description:  "Memory usage is high",
			CurrentValue: fmt.Sprintf("%.2f%%", memPercent),
			Threshold:    "80%",
			Suggestions:  []string{"Monitor memory usage", "Consider adding more RAM"},
			AutoFixable:  false,
		})
	}

	// Check disk usage
	diskInfo, _ := getDiskUsage("/")
	if diskInfo.UsedPercent > 90 {
		healthScore -= 20
		totalErrors++
		issues = append(issues, &pb.HealthIssue{
			Category:     "disk",
			Severity:     "critical",
			Description:  "Disk space is critically low",
			CurrentValue: fmt.Sprintf("%.2f%%", diskInfo.UsedPercent),
			Threshold:    "90%",
			Suggestions:  []string{"Clean up old logs and temporary files", "Remove unused packages", "Expand disk space", "Archive old data"},
			AutoFixable:  true,
		})
	} else if diskInfo.UsedPercent > 80 {
		healthScore -= 10
		totalWarnings++
		issues = append(issues, &pb.HealthIssue{
			Category:     "disk",
			Severity:     "warning",
			Description:  "Disk space is running low",
			CurrentValue: fmt.Sprintf("%.2f%%", diskInfo.UsedPercent),
			Threshold:    "80%",
			Suggestions:  []string{"Review and clean up unnecessary files", "Monitor disk usage"},
			AutoFixable:  true,
		})
	}

	// Check load average
	loadAvg := getLoadAverage()
	cpuCount := float64(runtime.NumCPU())
	loadRatio := loadAvg[0] / cpuCount

	if loadRatio > 2.0 {
		healthScore -= 15
		totalErrors++
		issues = append(issues, &pb.HealthIssue{
			Category:     "cpu",
			Severity:     "critical",
			Description:  "System load is very high",
			CurrentValue: fmt.Sprintf("%.2f (%.0fx CPU count)", loadAvg[0], loadRatio),
			Threshold:    "2x CPU count",
			Suggestions:  []string{"Check for hung processes", "Reduce concurrent workloads", "Increase CPU resources"},
			AutoFixable:  false,
		})
	} else if loadRatio > 1.5 {
		healthScore -= 5
		totalWarnings++
		issues = append(issues, &pb.HealthIssue{
			Category:     "cpu",
			Severity:     "warning",
			Description:  "System load is elevated",
			CurrentValue: fmt.Sprintf("%.2f (%.1fx CPU count)", loadAvg[0], loadRatio),
			Threshold:    "1.5x CPU count",
			Suggestions:  []string{"Monitor system load", "Review running processes"},
			AutoFixable:  false,
		})
	}

	// Deep check: check for system errors if requested
	if in.DeepCheck {
		errCmd := exec.Command("journalctl", "-p", "err", "-n", "10", "--since", "1 hour ago")
		errOutput, err := errCmd.Output()
		if err == nil {
			errorLines := strings.Split(string(errOutput), "\n")
			errorCount := 0
			for _, line := range errorLines {
				if line != "" {
					errorCount++
				}
			}

			if errorCount > 10 {
				healthScore -= 10
				totalWarnings++
				issues = append(issues, &pb.HealthIssue{
					Category:     "process",
					Severity:     "warning",
					Description:  fmt.Sprintf("System has %d recent errors in logs", errorCount),
					CurrentValue: fmt.Sprintf("%d errors", errorCount),
					Threshold:    "10 errors/hour",
					Suggestions:  []string{"Review system logs for recurring issues", "Check application health"},
					AutoFixable:  false,
				})
			}
		}
	}

	// Ensure health score doesn't go negative
	if healthScore < 0 {
		healthScore = 0
	}

	// Determine overall status
	overallStatus := "healthy"
	if healthScore < 50 {
		overallStatus = "unhealthy"
	} else if healthScore < 80 {
		overallStatus = "degraded"
	}

	summary := map[string]string{
		"cpu_usage":    fmt.Sprintf("%.2f%%", cpuPercent),
		"memory_usage": fmt.Sprintf("%.2f%%", memPercent),
		"disk_usage":   fmt.Sprintf("%.2f%%", diskInfo.UsedPercent),
		"load_avg":     fmt.Sprintf("%.2f", loadAvg[0]),
	}

	return &pb.HealthDiagnosticResponse{
		OverallStatus:  overallStatus,
		HealthScore:    healthScore,
		Issues:         issues,
		Summary:        summary,
		CheckTimestamp: time.Now().Unix(),
		TotalWarnings:  totalWarnings,
		TotalErrors:    totalErrors,
	}, nil
}

// InteractiveShell provides a bidirectional streaming shell interface with PTY
func (s *agentServer) InteractiveShell(stream pb.Agent_InteractiveShellServer) error {
	slog.Info("Starting interactive shell session")

	// Start a bash shell with PTY for proper interactive behavior
	cmd := exec.Command("/bin/bash", "-i")

	// Set environment variables for a nice interactive shell
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"PS1=\\[\\033[01;32m\\]\\u@\\h\\[\\033[00m\\]:\\[\\033[01;34m\\]\\w\\[\\033[00m\\]\\$ ",
		"HISTFILE=/dev/null", // Don't save history during remote sessions
	)

	// Create a PTY using creack/pty library
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start shell with PTY: %w", err)
	}
	defer ptmx.Close()

	// Channel to signal when to stop
	done := make(chan struct{})
	errChan := make(chan error, 2)

	// Goroutine to read from client and write to PTY
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				slog.Debug("Client closed stream")
				close(done)
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("stream recv error: %w", err)
				return
			}

			if in.Terminate {
				slog.Debug("Client requested termination")
				close(done)
				return
			}

			// Write any data received directly to PTY
			if len(in.StdinData) > 0 {
				if _, err := ptmx.Write(in.StdinData); err != nil {
					slog.Error("Failed to write to PTY", "error", err)
					errChan <- fmt.Errorf("failed to write to PTY: %w", err)
					return
				}
			}
		}
	}()

	// Goroutine to read from PTY and send to client
	go func() {
		buf := make([]byte, 8192)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				if err := stream.Send(&pb.ShellOutput{
					Stdout: buf[:n],
				}); err != nil {
					slog.Debug("Failed to send output (client disconnected)", "error", err)
					return
				}
			}
			if err != nil {
				// EOF or I/O error on PTY means shell exited normally
				if err == io.EOF {
					slog.Debug("Shell exited (EOF)")
				} else {
					slog.Debug("Shell exited (PTY closed)", "error", err)
				}
				close(done)
				return
			}
		}
	}()

	// Wait for shell to exit or error
	select {
	case err := <-errChan:
		slog.Info("Shell session ended with error", "error", err)
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return err
	case <-done:
		slog.Info("Shell session ended normally")
		// Client requested termination or shell exited
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}

	// Wait for process to exit
	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		slog.Debug("Shell process exited", "error", err)
	}

	return stream.Send(&pb.ShellOutput{
		ExitCode:  int32(exitCode),
		Completed: true,
	})
}

// parseDuration parses a duration string with a default fallback
func parseDuration(durationStr string, defaultDuration time.Duration) time.Duration {
	if durationStr == "" {
		return defaultDuration
	}

	d, err := time.ParseDuration(durationStr)
	if err != nil {
		slog.Warn("Failed to parse duration, using default", "input", durationStr, "default", defaultDuration)
		return defaultDuration
	}

	return d
}

// RegisterWatcher registers a new watcher on the agent
func (s *agentServer) RegisterWatcher(ctx context.Context, in *pb.RegisterWatcherRequest) (*pb.RegisterWatcherResponse, error) {
	if s.watcherManager == nil {
		return &pb.RegisterWatcherResponse{
			Success: false,
			Message: "watcher manager not initialized",
		}, nil
	}

	config := in.GetConfig()
	if config == nil {
		return &pb.RegisterWatcherResponse{
			Success: false,
			Message: "watcher config is required",
		}, nil
	}

	// Convert protobuf config to internal WatcherConfig
	watcherConfig := &agentInternal.WatcherConfig{
		ID:         config.GetId(),
		Type:       agentInternal.WatcherType(config.GetType()),
		Conditions: make([]agentInternal.EventCondition, 0),
		Interval:   parseDuration(config.GetInterval(), 10*time.Second),

		// File-specific
		FilePath:  config.GetFilePath(),
		CheckHash: config.GetCheckHash(),
		Recursive: config.GetRecursive(),

		// Process-specific
		ProcessName: config.GetProcessName(),
		PID:         int(config.GetPid()),

		// Port-specific
		Port:     int(config.GetPort()),
		Protocol: config.GetProtocol(),

		// Resource monitoring
		CPUThreshold:    config.GetCpuThreshold(),
		MemoryThreshold: config.GetMemoryThreshold(),
		DiskThreshold:   config.GetDiskThreshold(),
	}

	// Convert conditions
	for _, cond := range config.GetConditions() {
		watcherConfig.Conditions = append(watcherConfig.Conditions, agentInternal.EventCondition(cond))
	}

	// Register the watcher
	if watcherMgr, ok := s.watcherManager.(*agentInternal.EventWatcherManager); ok {
		if err := watcherMgr.RegisterWatcher(*watcherConfig); err != nil {
			return &pb.RegisterWatcherResponse{
				Success: false,
				Message: fmt.Sprintf("failed to register watcher: %v", err),
			}, nil
		}

		slog.Info("Watcher registered via gRPC", "id", watcherConfig.ID, "type", watcherConfig.Type)

		return &pb.RegisterWatcherResponse{
			Success:   true,
			Message:   "watcher registered successfully",
			WatcherId: watcherConfig.ID,
		}, nil
	}

	return &pb.RegisterWatcherResponse{
		Success: false,
		Message: "invalid watcher manager type",
	}, nil
}

// ListWatchers lists all registered watchers on the agent
func (s *agentServer) ListWatchers(ctx context.Context, in *pb.ListWatchersRequest) (*pb.ListWatchersResponse, error) {
	if s.watcherManager == nil {
		return &pb.ListWatchersResponse{
			Watchers: []*pb.WatcherConfig{},
		}, nil
	}

	// Type assert to get access to list method
	if watcherMgr, ok := s.watcherManager.(*agentInternal.EventWatcherManager); ok {
		watchers := watcherMgr.ListWatchers()
		pbWatchers := make([]*pb.WatcherConfig, 0, len(watchers))

		for _, w := range watchers {
			conditions := make([]string, 0, len(w.Conditions))
			for _, c := range w.Conditions {
				conditions = append(conditions, string(c))
			}

			pbWatchers = append(pbWatchers, &pb.WatcherConfig{
				Id:              w.ID,
				Type:            string(w.Type),
				Conditions:      conditions,
				Interval:        w.Interval.String(),
				FilePath:        w.FilePath,
				CheckHash:       w.CheckHash,
				Recursive:       w.Recursive,
				ProcessName:     w.ProcessName,
				Pid:             int32(w.PID),
				Port:            int32(w.Port),
				Protocol:        w.Protocol,
				CpuThreshold:    w.CPUThreshold,
				MemoryThreshold: w.MemoryThreshold,
				DiskThreshold:   w.DiskThreshold,
			})
		}

		return &pb.ListWatchersResponse{
			Watchers: pbWatchers,
		}, nil
	}

	return &pb.ListWatchersResponse{
		Watchers: []*pb.WatcherConfig{},
	}, nil
}

// RemoveWatcher removes a registered watcher from the agent
func (s *agentServer) RemoveWatcher(ctx context.Context, in *pb.RemoveWatcherRequest) (*pb.RemoveWatcherResponse, error) {
	if s.watcherManager == nil {
		return &pb.RemoveWatcherResponse{
			Success: false,
			Message: "watcher manager not initialized",
		}, nil
	}

	watcherID := in.GetWatcherId()
	if watcherID == "" {
		return &pb.RemoveWatcherResponse{
			Success: false,
			Message: "watcher ID is required",
		}, nil
	}

	if watcherMgr, ok := s.watcherManager.(*agentInternal.EventWatcherManager); ok {
		if err := watcherMgr.RemoveWatcher(watcherID); err != nil {
			return &pb.RemoveWatcherResponse{
				Success: false,
				Message: fmt.Sprintf("failed to remove watcher: %v", err),
			}, nil
		}

		slog.Info("Watcher removed via gRPC", "id", watcherID)

		return &pb.RemoveWatcherResponse{
			Success: true,
			Message: "watcher removed successfully",
		}, nil
	}

	return &pb.RemoveWatcherResponse{
		Success: false,
		Message: "invalid watcher manager type",
	}, nil
}
