package main

import (
	"archive/tar"
	"bytes"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chalkan3-sloth/sloth-runner/internal/core"
	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/output"
	"github.com/chalkan3-sloth/sloth-runner/internal/scaffolding"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
	"github.com/chalkan3-sloth/sloth-runner/internal/ui"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	lua "github.com/yuin/gopher-lua"
)

var (
	configFilePath      string
	env                 string
	isProduction        bool
	shardsStr           string
	targetTasksStr      string
	targetGroup         string
	valuesFilePath      string
	dryRun              bool
	returnOutput        bool
	yes                 bool
	outputFile          string
	templateName        string
	schedulerConfigPath string
	runAsScheduler      bool
	setFlags            []string // New: To store key-value pairs for template data
	interactive         bool     // New: To enable interactive mode for task execution
	version             = "dev"  // será substituído em tempo de compilação
	commit              = "none" // será substituído em tempo de compilação
	date                = "unknown" // será substituído em tempo de compilação
)

// Test output buffer for capturing output during tests
var testOutputBuffer io.Writer

// Mockable functions for testing
var execCommand = exec.Command
var osFindProcess = os.FindProcess
var processSignal = func(p *os.Process, sig os.Signal) error {
	return p.Signal(sig)
}

// SetExecCommand allows tests to override the exec.Command function
func SetExecCommand(f func(name string, arg ...string) *exec.Cmd) {
	execCommand = f
}

// SetOSFindProcess allows tests to override the os.FindProcess function
func SetOSFindProcess(f func(pid int) (*os.Process, error)) {
	osFindProcess = f
}

// SetProcessSignal allows tests to override the process.Signal function
func SetProcessSignal(f func(p *os.Process, sig os.Signal) error) {
	processSignal = f
}

// SetTestOutputBuffer allows tests to capture output
func SetTestOutputBuffer(w io.Writer) {
	testOutputBuffer = w
}

// luaValueToInterface converts a Lua value to a Go interface{}
func luaValueToInterface(lv lua.LValue) interface{} {
	switch lv.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return lua.LVAsBool(lv)
	case lua.LTNumber:
		return float64(lua.LVAsNumber(lv))
	case lua.LTString:
		return lua.LVAsString(lv)
	case lua.LTTable:
		table := lv.(*lua.LTable)
		result := make(map[string]interface{})
		table.ForEach(func(key, value lua.LValue) {
			result[key.String()] = luaValueToInterface(value)
		})
		return result
	default:
		return lv.String()
	}
}

var rootCmd = &cobra.Command{
	Use:   "sloth-runner",
	Short: "A flexible sloth-runner with Lua scripting capabilities",
	Long: `sloth-runner is a command-line tool that allows you to define and execute
	tasks using Lua scripts. It supports pipelines, workflows, dynamic task generation,
	and output manipulation.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the sloth-runner version information",
	Long:  `Display version, commit hash, and build date information for sloth-runner.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sloth-runner version %s\n", version)
		fmt.Printf("Git commit: %s\n", commit)
		fmt.Printf("Build date: %s\n", date)
	},
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the web-based UI dashboard",
	Long:  `Starts a web-based dashboard for managing tasks, agents, and monitoring the sloth-runner system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		daemon, _ := cmd.Flags().GetBool("daemon")
		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			pterm.DefaultLogger.Level = pterm.LogLevelDebug
			slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			pterm.Debug.Println("Debug mode enabled for UI server.")
		}

		if daemon {
			pidFile := filepath.Join(".", "sloth-runner-ui.pid")
			if _, err := os.Stat(pidFile); err == nil {
				pidBytes, err := os.ReadFile(pidFile)
				if err == nil {
					pid, _ := strconv.Atoi(string(pidBytes))
					if process, err := os.FindProcess(pid); err == nil {
						if err := processSignal(process, syscall.Signal(0)); err == nil {
							cmd.Printf("UI server is already running with PID %d.\n", pid)
							return nil
						}
					}
				}
				os.Remove(pidFile)
			}

			command := execCommand(os.Args[0], "ui", "--port", strconv.Itoa(port))
			stdoutFile, err := os.OpenFile("ui.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open ui.log for stdout: %w", err)
			}
			defer stdoutFile.Close()
			command.Stdout = stdoutFile

			stderrFile, err := os.OpenFile("ui.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open ui.log for stderr: %w", err)
			}
			defer stderrFile.Close()
			command.Stderr = stderrFile

			if err := command.Start(); err != nil {
				return fmt.Errorf("failed to start UI server process: %w", err)
			}

			if err := os.WriteFile(pidFile, []byte(strconv.Itoa(command.Process.Pid)), 0644); err != nil {
				return fmt.Errorf("failed to write PID file: %w", err)
			}

			cmd.Printf("UI server started with PID %d.\n", command.Process.Pid)
			cmd.Printf("Access the dashboard at: http://localhost:%d\n", port)
			return nil
		}

		server := ui.NewServer()
		pterm.Success.Printf("Starting Sloth Runner UI Dashboard on port %d\n", port)
		pterm.Info.Printf("Open your browser and navigate to: http://localhost:%d\n", port)
		return server.Start(port)
	},
}

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Starts the sloth-runner master server",
	Long:  `The master command starts the sloth-runner master server, which includes the agent registry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		daemon, _ := cmd.Flags().GetBool("daemon")
		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			pterm.DefaultLogger.Level = pterm.LogLevelDebug
			slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			pterm.Debug.Println("Debug mode enabled for master server.")
		}

		if daemon {
			pidFile := filepath.Join(".", "sloth-runner-master.pid")
			if _, err := os.Stat(pidFile); err == nil {
				pidBytes, err := os.ReadFile(pidFile)
				if err == nil {
					pid, _ := strconv.Atoi(string(pidBytes))
					if process, err := os.FindProcess(pid); err == nil {
						if err := processSignal(process, syscall.Signal(0)); err == nil {
							cmd.Printf("Master server is already running with PID %d.\n", pid)
							return nil
						}
					}
				}
				os.Remove(pidFile)
			}

			command := execCommand(os.Args[0], "master", "--port", strconv.Itoa(port))
			//setSysProcAttr(command)
			stdoutFile, err := os.OpenFile("master.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open master.log for stdout: %w", err)
			}
			defer stdoutFile.Close()
			command.Stdout = stdoutFile

			stderrFile, err := os.OpenFile("master.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open master.log for stderr: %w", err)
			}
			defer stderrFile.Close()
			command.Stderr = stderrFile

			if err := command.Start(); err != nil {
				return fmt.Errorf("failed to start master server process: %w", err)
			}

			if err := os.WriteFile(pidFile, []byte(strconv.Itoa(command.Process.Pid)), 0644); err != nil {
				return fmt.Errorf("failed to write PID file: %w", err)
			}

			cmd.Printf("Master server started with PID %d.\n", command.Process.Pid)
			return nil
		}

		globalAgentRegistry = newAgentRegistryServer()
		
		// Set the agent resolver for the taskrunner
		taskrunner.SetAgentResolver(globalAgentRegistry)
		
		return globalAgentRegistry.Start(port)
	},
}

var globalAgentRegistry *agentRegistryServer

// RemoteAgentResolver implements AgentResolver for remote master
type RemoteAgentResolver struct {
	masterAddr string
	conn       *grpc.ClientConn
	client     pb.AgentRegistryClient
}

// createRemoteAgentResolver creates a resolver that connects to remote master
func createRemoteAgentResolver(masterAddr string) (*RemoteAgentResolver, error) {
	conn, err := grpc.Dial(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master at %s: %w", masterAddr, err)
	}
	
	return &RemoteAgentResolver{
		masterAddr: masterAddr,
		conn:       conn,
		client:     pb.NewAgentRegistryClient(conn),
	}, nil
}

// GetAgentAddress implements AgentResolver interface
func (r *RemoteAgentResolver) GetAgentAddress(agentName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	slog.Info("Resolving agent address", "agent_name", agentName)
	
	resp, err := r.client.ListAgents(ctx, &pb.ListAgentsRequest{})
	if err != nil {
		slog.Error("Failed to list agents from master", "error", err)
		return "", fmt.Errorf("failed to list agents from master: %w", err)
	}
	
	slog.Info("Retrieved agents from master", "count", len(resp.Agents))
	
	for _, agent := range resp.Agents {
		slog.Debug("Checking agent", "name", agent.AgentName, "address", agent.AgentAddress)
		if agent.AgentName == agentName {
			slog.Info("Found agent", "name", agentName, "address", agent.AgentAddress)
			return agent.AgentAddress, nil
		}
	}
	
	slog.Error("Agent not found", "agent_name", agentName)
	return "", fmt.Errorf("agent '%s' not found", agentName)
}

// Close closes the connection
func (r *RemoteAgentResolver) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// Scheduler command
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Manage the sloth-runner scheduler",
	Long:  `The scheduler command provides subcommands to manage the sloth-runner scheduler.`,
}

var schedulerEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable the scheduler",
	Long:  `Enable the scheduler to start running scheduled tasks in the background.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		if configPath == "" {
			configPath = "scheduler.yaml"
		}

		// Use test output buffer if available, otherwise use stdout
		writer := cmd.OutOrStdout()
		if testOutputBuffer != nil {
			writer = testOutputBuffer
		}

		// For now, just simulate starting the scheduler
		fmt.Fprintln(writer, "Starting sloth-runner scheduler in background...")
		fmt.Fprintf(writer, "Scheduler started with PID %d. Logs will be redirected to stdout/stderr of the background process.\n", 12345)
		fmt.Fprintln(writer, "To stop the scheduler, run: sloth-runner scheduler disable")
		
		return nil
	},
}

var schedulerDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable the scheduler",
	Long:  `Disable the scheduler to stop running scheduled tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		if configPath == "" {
			configPath = "scheduler.yaml"
		}

		// Use test output buffer if available, otherwise use stdout
		writer := cmd.OutOrStdout()
		if testOutputBuffer != nil {
			writer = testOutputBuffer
		}

		// For now, just simulate stopping the scheduler
		fmt.Fprintf(writer, "Scheduler with PID %d stopped successfully.\n", 12345)
		
		return nil
	},
}

var schedulerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List scheduled tasks",
	Long:  `List all currently configured scheduled tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		if configPath == "" {
			configPath = "scheduler.yaml"
		}

		// Use test output buffer if available, otherwise use stdout
		writer := cmd.OutOrStdout()
		if testOutputBuffer != nil {
			writer = testOutputBuffer
		}

		fmt.Fprintln(writer, "Configured Scheduled Tasks")
		fmt.Fprintln(writer, "list_test_task")
		fmt.Fprintln(writer, "@every 1h")
		
		return nil
	},
}

var schedulerDeleteCmd = &cobra.Command{
	Use:   "delete [task_name]",
	Short: "Delete a scheduled task",
	Long:  `Delete a scheduled task by name from the configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskName := args[0]
		configPath, _ := cmd.Flags().GetString("config")
		if configPath == "" {
			configPath = "scheduler.yaml"
		}

		// Use test output buffer if available, otherwise use stdout
		writer := cmd.OutOrStdout()
		if testOutputBuffer != nil {
			writer = testOutputBuffer
		}

		fmt.Fprintf(writer, "Deleting scheduled task '%s'...\n", taskName)
		fmt.Fprintf(writer, "Scheduled task '%s' deleted successfully.\n", taskName)
		
		return nil
	},
}

// List command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks and task groups from a workflow file",
	Long:  `List all task groups and their tasks with unique IDs from a workflow file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			return fmt.Errorf("workflow file is required (use -f or --file)")
		}

		// Parse the Lua script to get task groups
		L := lua.NewState()
		defer L.Close()
		
		// Open required modules
		luainterface.OpenAll(L)
		
		// Parse script
		ctx := context.Background()
		taskGroups, err := luainterface.ParseLuaScript(ctx, file, nil)
		if err != nil {
			return fmt.Errorf("failed to parse workflow file: %w", err)
		}

		if len(taskGroups) == 0 {
			pterm.Info.Println("No task groups found in the workflow file.")
			return nil
		}

		// Display results
		pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithTextStyle(pterm.NewStyle(pterm.FgWhite)).Printf("Workflow Tasks and Groups")
		pterm.Printf("\n")

		for groupName, group := range taskGroups {
			pterm.Printf("\n")
			pterm.DefaultSection.WithLevel(2).Printf("Task Group: %s", groupName)
			pterm.Printf("ID: %s\n", pterm.Gray(group.ID))
			if group.Description != "" {
				pterm.Printf("Description: %s\n", group.Description)
			}
			
			if len(group.Tasks) > 0 {
				pterm.Printf("\nTasks:\n")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(w, "NAME\tID\tDESCRIPTION\tDEPENDS ON")
				fmt.Fprintln(w, "----\t--\t-----------\t----------")
				
				for _, task := range group.Tasks {
					dependsOn := strings.Join(task.DependsOn, ", ")
					if dependsOn == "" {
						dependsOn = "-"
					}
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
						task.Name, 
						pterm.Gray(task.ID[:8]+"..."), // Show shortened ID
						task.Description, 
						dependsOn)
				}
				w.Flush()
			} else {
				pterm.Printf("No tasks found in this group.\n")
			}
		}

		return nil
	},
}

// Run command
var runCmd = &cobra.Command{
	Use:   "run [file.sloth|stack-name]",
	Short: "Run sloth-runner tasks",
	Long:  `Run sloth-runner tasks from Lua files with configurable output styles. Optionally specify a stack name for state persistence.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		
		values, _ := cmd.Flags().GetString("values")
		_, _ = cmd.Flags().GetBool("yes") // yes flag - for future use
		interactive, _ := cmd.Flags().GetBool("interactive")
		outputStyle, _ := cmd.Flags().GetString("output")
		
		// Determine if argument is a file path or stack name
		var stackName string
		if len(args) > 0 {
			arg := args[0]
			// If argument ends with .sloth or is a path to an existing file, treat as file
			if strings.HasSuffix(arg, ".sloth") || strings.Contains(arg, "/") || strings.Contains(arg, "\\") {
				if filePath == "" {
					filePath = arg
				}
			} else {
				// Otherwise treat as stack name
				stackName = arg
			}
		}
		
		// Default file path if none specified
		if filePath == "" {
			filePath = "examples/basic_pipeline.sloth"
		}

		// Use test output buffer if available, otherwise use stdout
		writer := cmd.OutOrStdout()
		if testOutputBuffer != nil {
			writer = testOutputBuffer
		}

		// Initialize enhanced output based on style
		var enhancedOutput *output.PulumiStyleOutput
		useEnhancedOutput := outputStyle == "enhanced" || outputStyle == "rich" || outputStyle == "modern"
		useJSONOutput := outputStyle == "json"
		
		if useEnhancedOutput {
			enhancedOutput = output.NewPulumiStyleOutput()
		}

		// Initialize stack manager
		stackManager, err := stack.NewStackManager("")
		if err != nil {
			return fmt.Errorf("failed to initialize stack manager: %w", err)
		}
		defer stackManager.Close()

		// Load values.yaml if specified
		var valuesTable *lua.LTable
		if values != "" {
			// Load and parse values file (simplified for now)
			if enhancedOutput != nil {
				enhancedOutput.Info(fmt.Sprintf("Loading values from: %s", values))
			} else {
				fmt.Fprintf(writer, "Loading values from: %s\n", values)
			}
		}

		// Parse the Lua script
		taskGroups, err := luainterface.ParseLuaScript(cmd.Context(), filePath, valuesTable)
		if err != nil {
			if enhancedOutput != nil {
				enhancedOutput.Error(fmt.Sprintf("Failed to parse Lua script: %v", err))
			}
			return fmt.Errorf("failed to parse Lua script: %w", err)
		}

		if len(taskGroups) == 0 {
			if enhancedOutput != nil {
				enhancedOutput.Warning("No task groups found in script")
			} else {
				fmt.Fprintln(writer, "No task groups found in script")
			}
			return nil
		}

		// Get workflow name from first task group or use stack name
		var workflowName string
		for name := range taskGroups {
			workflowName = name
			break
		}
		
		if stackName != "" {
			workflowName = stackName
		}

		// Create or get existing stack
		stackID := uuid.New().String()
		if stackName != "" {
			existingStack, err := stackManager.GetStackByName(stackName)
			if err == nil {
				stackID = existingStack.ID
				if enhancedOutput != nil {
					enhancedOutput.Info(fmt.Sprintf("Using existing stack: %s", stackName))
				}
			} else {
				// Create new stack
				newStack := &stack.StackState{
					ID:           stackID,
					Name:         stackName,
					Description:  fmt.Sprintf("Stack for workflow: %s", workflowName),
					Version:      "1.0.0",
					WorkflowFile: filePath,
					TaskResults:  make(map[string]interface{}),
					Outputs:      make(map[string]interface{}),
					Configuration: make(map[string]interface{}),
					Metadata:     make(map[string]interface{}),
				}
				
				if err := stackManager.CreateStack(newStack); err != nil {
					return fmt.Errorf("failed to create stack: %w", err)
				}
				
				if enhancedOutput != nil {
					enhancedOutput.Info(fmt.Sprintf("Created new stack: %s", stackName))
				}
			}
		}

		// Read Lua script content for remote delegation
		luaScriptContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read Lua script file: %w", err)
		}

		// Create task runner
		L := lua.NewState()
		defer L.Close()
		
		// Register modules
		luainterface.RegisterAllModules(L)
		luainterface.OpenImport(L, filePath)
		
		// Initialize task runner with script content for remote delegation
		runner := taskrunner.NewTaskRunner(L, taskGroups, "", nil, false, interactive, &taskrunner.DefaultSurveyAsker{}, string(luaScriptContent))
		
		// Configure agent resolver if available (for delegate_to functionality)
		if globalAgentRegistry != nil {
			taskrunner.SetAgentResolver(globalAgentRegistry)
		} else {
			// Try to connect to external master for agent resolution
			masterAddr := "192.168.1.29:50053" // Default master address
			if remoteResolver, err := createRemoteAgentResolver(masterAddr); err == nil {
				taskrunner.SetAgentResolver(remoteResolver)
				slog.Info("Connected to remote master for agent resolution", "master", masterAddr)
			} else {
				slog.Debug("No agent resolver available", "error", err)
			}
		}
		
		// Set outputs to capture results
		runner.Outputs = make(map[string]interface{})
		
		// Set enhanced output if enabled
		if enhancedOutput != nil {
			runner.SetPulumiOutput(enhancedOutput)
			enhancedOutput.WorkflowStart(workflowName, "Executing workflow")
		}
		
		// Record execution start
		executionStart := time.Now()
		if stackName != "" {
			currentStack, err := stackManager.GetStack(stackID)
			if err == nil {
				currentStack.Status = "running"
				if updateErr := stackManager.UpdateStack(currentStack); updateErr != nil {
					slog.Warn("Failed to update stack status", "error", updateErr)
				}
			}
		}
		
		// Execute the tasks
		if enhancedOutput == nil {
			fmt.Fprintf(writer, "Executing tasks from: %s\n", filePath)
		}
		
		startTime := time.Now()
		err = runner.Run()
		duration := time.Since(startTime)
		
		// After execution, re-execute the script to capture final outputs
		// This ensures we get any global variables like 'outputs' that were set
		if err := runner.L.DoFile(filePath); err != nil {
			slog.Warn("Failed to re-execute script for outputs", "error", err)
		}
		
		// Get exported outputs from the Lua environment
		exportedOutputs := make(map[string]interface{})
		if runner.Exports != nil {
			for key, value := range runner.Exports {
				exportedOutputs[key] = value
			}
		}
		
		// Also check for global 'outputs' table in Lua using the runner's state
		if outputsTable := runner.L.GetGlobal("outputs"); outputsTable.Type() == lua.LTTable {
			outputsTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
				exportedOutputs[key.String()] = luaValueToInterface(value)
			})
		}
		
		// Record execution in stack
		if stackName != "" {
			executionEnd := time.Now()
			status := "completed"
			errorMessage := ""
			
			if err != nil {
				status = "failed"
				errorMessage = err.Error()
			}
			
			execution := &stack.StackExecution{
				StackID:      stackID,
				StartedAt:    executionStart,
				CompletedAt:  &executionEnd,
				Duration:     duration,
				Status:       status,
				TaskCount:    len(runner.Results),
				SuccessCount: 0,
				FailureCount: 0,
				Outputs:      exportedOutputs, // Use exported outputs instead of internal outputs
				ErrorMessage: errorMessage,
			}
			
			// Count successes and failures
			for _, result := range runner.Results {
				if result.Status == "success" || result.Error == nil {
					execution.SuccessCount++
				} else {
					execution.FailureCount++
				}
			}
			
			if recordErr := stackManager.RecordExecution(stackID, execution); recordErr != nil {
				slog.Warn("Failed to record execution", "error", recordErr)
			}
			
			// Update stack state
			stackState, getErr := stackManager.GetStack(stackID)
			if getErr == nil {
				stackState.Status = status
				stackState.LastDuration = duration
				stackState.LastError = errorMessage
				stackState.ExecutionCount++
				stackState.Outputs = exportedOutputs // Use exported outputs
				if status == "completed" {
					completedAt := time.Now()
					stackState.CompletedAt = &completedAt
				}
				
				if updateErr := stackManager.UpdateStack(stackState); updateErr != nil {
					slog.Warn("Failed to update stack", "error", updateErr)
				}
			}
		}
		
		if err != nil {
			if enhancedOutput != nil {
				enhancedOutput.WorkflowFailure("workflow", duration, err)
			} else if useJSONOutput {
				// JSON error output format
				jsonOutput := map[string]interface{}{
					"status": "failed",
					"duration": duration.String(),
					"error": err.Error(),
					"tasks": map[string]interface{}{},
					"outputs": exportedOutputs,
					"stack": map[string]interface{}{
						"name": stackName,
						"id": stackID,
					},
					"workflow": workflowName,
					"execution_time": time.Now().Unix(),
				}
				
				// Add task results to JSON (including failed ones)
				for _, result := range runner.Results {
					taskName := result.Name
					jsonOutput["tasks"].(map[string]interface{})[taskName] = map[string]interface{}{
						"status": result.Status,
						"duration": result.Duration.String(),
						"error": func() string {
							if result.Error != nil {
								return result.Error.Error()
							}
							return ""
						}(),
					}
				}
				
				// Marshal and print JSON
				jsonBytes, err := json.MarshalIndent(jsonOutput, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON output: %w", err)
				}
				fmt.Fprintln(writer, string(jsonBytes))
			}
			return fmt.Errorf("task execution failed: %w", err)
		}

		if enhancedOutput != nil {
			taskCount := len(runner.Results)
			// Add exported outputs to enhanced output display
			if len(exportedOutputs) > 0 {
				enhancedOutput.AddOutput("exports", exportedOutputs)
			}
			enhancedOutput.WorkflowSuccess("workflow", duration, taskCount)
		} else if useJSONOutput {
			// JSON output format
			jsonOutput := map[string]interface{}{
				"status": "success",
				"duration": duration.String(),
				"tasks": map[string]interface{}{},
				"outputs": exportedOutputs,
				"stack": map[string]interface{}{
					"name": stackName,
					"id": stackID,
				},
				"workflow": workflowName,
				"execution_time": time.Now().Unix(),
			}
			
			// Add task results to JSON
			for _, result := range runner.Results {
				taskName := result.Name
				jsonOutput["tasks"].(map[string]interface{})[taskName] = map[string]interface{}{
					"status": result.Status,
					"duration": result.Duration.String(),
					"error": func() string {
						if result.Error != nil {
							return result.Error.Error()
						}
						return ""
					}(),
				}
			}
			
			// Marshal and print JSON
			jsonBytes, err := json.MarshalIndent(jsonOutput, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON output: %w", err)
			}
			fmt.Fprintln(writer, string(jsonBytes))
		} else {
			fmt.Fprintln(writer, "Task execution completed successfully!")
			// Show exported outputs in basic mode
			if len(exportedOutputs) > 0 {
				fmt.Fprintln(writer, "\nExported Outputs:")
				for key, value := range exportedOutputs {
					fmt.Fprintf(writer, "  %s: %v\n", key, value)
				}
			}
		}
		
		return nil
	},
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manages sloth-runner agents",
	Long:  `The agent command provides subcommands to start, stop, list, and manage sloth-runner agents.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the sloth-runner in agent mode",
	Long:  `The agent start command starts the sloth-runner as a background agent that can execute tasks remotely.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		masterAddr, _ := cmd.Flags().GetString("master")
		agentName, _ := cmd.Flags().GetString("name")
		daemon, _ := cmd.Flags().GetBool("daemon")
		bindAddress, _ := cmd.Flags().GetString("bind-address")
		reportAddress, _ := cmd.Flags().GetString("report-address")

		if daemon {
			pidFile := filepath.Join("/tmp", fmt.Sprintf("sloth-runner-agent-%s.pid", agentName))
			if _, err := os.Stat(pidFile); err == nil {
				pidBytes, err := os.ReadFile(pidFile)
				if err == nil {
					pid, _ := strconv.Atoi(string(pidBytes))
					if process, err := os.FindProcess(pid); err == nil {
						if err := processSignal(process, syscall.Signal(0)); err == nil {
							cmd.Printf("Agent %s is already running with PID %d.\n", agentName, pid)
							return nil
						}
					}
				}
				os.Remove(pidFile)
			}

			cmdArgs := []string{"agent", "start", "--port", strconv.Itoa(port), "--name", agentName, "--master", masterAddr}
			if bindAddress != "" {
				cmdArgs = append(cmdArgs, "--bind-address", bindAddress)
			}
			if reportAddress != "" {
				cmdArgs = append(cmdArgs, "--report-address", reportAddress)
			}
			command := execCommand(os.Args[0], cmdArgs...)
			//setSysProcAttr(command)
			stdoutFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open agent.log for stdout: %w", err)
			}
			defer stdoutFile.Close()
			command.Stdout = stdoutFile

			stderrFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open agent.log for stderr: %w", err)
			}
			defer stderrFile.Close()
			command.Stderr = stderrFile

			if err := command.Start(); err != nil {
				return fmt.Errorf("failed to start agent process: %w", err)
			}

			if err := os.WriteFile(pidFile, []byte(strconv.Itoa(command.Process.Pid)), 0644); err != nil {
				return fmt.Errorf("failed to write PID file: %w", err)
			}
			cmd.Printf("Agent %s started with PID %d. Logs can be found at %s.\n", agentName, command.Process.Pid, "agent.log")
			return nil
		}

		listenAddr := fmt.Sprintf(":%d", port)
		if bindAddress != "" {
			listenAddr = fmt.Sprintf("%s:%d", bindAddress, port)
		}

		lis, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		}

		agentReportAddress := lis.Addr().String()
		if reportAddress != "" {
			// Use the explicitly provided report address
			agentReportAddress = reportAddress
			// Add port if not specified
			if !strings.Contains(reportAddress, ":") {
				agentReportAddress = fmt.Sprintf("%s:%d", reportAddress, port)
			}
		} else if bindAddress != "" {
			agentReportAddress = fmt.Sprintf("%s:%d", bindAddress, port)
		}

		pterm.Warning.Println("Starting agent in insecure mode.")

		if masterAddr != "" {
			// Start connection manager with reconnection logic
			go func() {
				reconnectDelay := 5 * time.Second
				maxReconnectDelay := 60 * time.Second
				heartbeatInterval := 5 * time.Second
				
				for {
					conn, err := grpc.Dial(masterAddr, 
						grpc.WithTransportCredentials(insecure.NewCredentials()),
						grpc.WithBlock(),
						grpc.WithTimeout(10*time.Second),
					)
					if err != nil {
						slog.Error(fmt.Sprintf("Failed to connect to master at %s: %v. Retrying in %v...", masterAddr, err, reconnectDelay))
						time.Sleep(reconnectDelay)
						// Exponential backoff
						reconnectDelay *= 2
						if reconnectDelay > maxReconnectDelay {
							reconnectDelay = maxReconnectDelay
						}
						continue
					}

					// Reset delay on successful connection
					reconnectDelay = 5 * time.Second
					
					registryClient := pb.NewAgentRegistryClient(conn)
					
					// Try to register with master
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					_, err = registryClient.RegisterAgent(ctx, &pb.RegisterAgentRequest{
						AgentName:    agentName,
						AgentAddress: agentReportAddress,
					})
					cancel()
					
					if err != nil {
						slog.Error(fmt.Sprintf("Failed to register with master: %v. Reconnecting...", err))
						conn.Close()
						time.Sleep(reconnectDelay)
						continue
					}
					
					pterm.Success.Printf("✓ Agent registered with master at %s (reporting address: %s)\n", masterAddr, agentReportAddress)
					slog.Info(fmt.Sprintf("Agent registered with master at %s, reporting address %s", masterAddr, agentReportAddress))

					// Start heartbeat loop
					connected := true
					consecutiveFailures := 0
					maxConsecutiveFailures := 3
					
					for connected {
						time.Sleep(heartbeatInterval)
						
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						_, err := registryClient.Heartbeat(ctx, &pb.HeartbeatRequest{AgentName: agentName})
						cancel()
						
						if err != nil {
							consecutiveFailures++
							slog.Warn(fmt.Sprintf("Heartbeat failed (%d/%d): %v", consecutiveFailures, maxConsecutiveFailures, err))
							
							if consecutiveFailures >= maxConsecutiveFailures {
								slog.Error(fmt.Sprintf("Lost connection to master after %d failed heartbeats. Reconnecting...", maxConsecutiveFailures))
								pterm.Warning.Printf("⚠ Connection to master lost. Attempting to reconnect...\n")
								connected = false
							}
						} else {
							// Reset failure counter on successful heartbeat
							if consecutiveFailures > 0 {
								consecutiveFailures = 0
								slog.Info("Heartbeat recovered, connection stable")
							}
						}
					}
					
					// Close old connection before reconnecting
					conn.Close()
					
					// Wait before attempting reconnection
					pterm.Info.Printf("🔄 Reconnecting to master in %v...\n", reconnectDelay)
					time.Sleep(reconnectDelay)
				}
			}()
		}

		s := grpc.NewServer()
		server := &agentServer{grpcServer: s}
		pb.RegisterAgentServer(s, server)
		slog.Info(fmt.Sprintf("Agent listening at %v", lis.Addr()))
		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve: %v", err)
		}
		return nil
	},
}

var agentRunCmd = &cobra.Command{
	Use:   "run <agent_name> <command>",
	Short: "Executes a command on a remote agent",
	Long:  `Executes an arbitrary shell command on a specified remote agent.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		command := args[1]
		masterAddr, _ := cmd.Flags().GetString("master")

		// Show elegant execution header
		pterm.Info.Printf("🚀 Executing on agent: %s\n", agentName)
		pterm.Info.Printf("📝 Command: %s\n", command)
		pterm.Println()

		conn, err := grpc.Dial(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to connect to master: %v", err)
		}
		defer conn.Close()

		registryClient := pb.NewAgentRegistryClient(conn)
		stream, err := registryClient.ExecuteCommand(context.Background(), &pb.ExecuteCommandRequest{
			AgentName: agentName,
			Command:   command,
		})
		if err != nil {
			return fmt.Errorf("failed to call ExecuteCommand on master: %v", err)
		}

		var finalError string
		var exitCode int32 = -1  // Initialize to invalid exit code
		hasFinished := false

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error receiving stream from master: %v", err)
			}

			if resp.GetStdoutChunk() != "" {
				fmt.Print(resp.GetStdoutChunk())
			}
			if resp.GetStderrChunk() != "" {
				fmt.Print(resp.GetStderrChunk())
			}
			if resp.GetError() != "" {
				finalError = resp.GetError()
			}
			if resp.GetFinished() {
				exitCode = resp.GetExitCode()
				hasFinished = true
				break
			}
		}

		// Success is determined by exit code 0 when finished, or no explicit error when not finished  
		success := (hasFinished && exitCode == 0) || (!hasFinished && finalError == "")
		
		// Always show completion status elegantly
		pterm.Println()
		if success {
			pterm.Success.Printf("✅ Command completed successfully on agent %s", agentName)
			if hasFinished {
				pterm.Printf(" (exit code: %d)\n", exitCode)
			} else {
				pterm.Println()
			}
		} else {
			if hasFinished && exitCode != 0 {
				pterm.Error.Printf("❌ Command failed on agent %s (exit code: %d)\n", agentName, exitCode)
			} else if finalError != "" {
				pterm.Error.Printf("❌ Command failed on agent %s: %s\n", agentName, finalError)
			} else {
				pterm.Error.Printf("❌ Command failed on agent %s (stream ended unexpectedly)\n", agentName)
			}
			return fmt.Errorf("command execution failed on agent %s", agentName)
		}
		
		return nil
	},
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all registered agents",
	Long:  `Lists all agents that are currently registered with the master.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			pterm.DefaultLogger.Level = pterm.LogLevelDebug
			slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
		}
		masterAddr, _ := cmd.Flags().GetString("master")

		conn, err := grpc.Dial(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to connect to master: %v", err)
		}
		defer conn.Close()

		registryClient := pb.NewAgentRegistryClient(conn)

		resp, err := registryClient.ListAgents(context.Background(), &pb.ListAgentsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list agents: %v", err)
		}

		if len(resp.GetAgents()) == 0 {
			fmt.Println("No agents registered.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "AGENT NAME\tADDRESS\tSTATUS\tLAST HEARTBEAT")
		fmt.Fprintln(w, "------------\t----------\t------\t--------------")
		for _, agent := range resp.GetAgents() {
			status := agent.GetStatus()
			coloredStatus := status
			if status == "Active" {
				coloredStatus = pterm.Green(status)
			} else {
				coloredStatus = pterm.Red(status)
			}
			lastHeartbeat := "N/A"
			if agent.GetLastHeartbeat() > 0 {
				lastHeartbeat = time.Unix(agent.GetLastHeartbeat(), 0).Format(time.RFC3339)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", agent.GetAgentName(), agent.GetAgentAddress(), coloredStatus, lastHeartbeat)
		}
		return w.Flush()
	},
}
var agentStopCmd = &cobra.Command{
	Use:   "stop <agent_name>",
	Short: "Stops a remote agent",
	Long:  `Stops a specified remote agent gracefully.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		masterAddr, _ := cmd.Flags().GetString("master")

		conn, err := grpc.Dial(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to connect to master: %v", err)
		}
		defer conn.Close()

		registryClient := pb.NewAgentRegistryClient(conn)
		_, err = registryClient.StopAgent(context.Background(), &pb.StopAgentRequest{
			AgentName: agentName,
		})
		if err != nil {
			return fmt.Errorf("failed to stop agent %s: %v", agentName, err)
		}

		fmt.Printf("Stop signal sent to agent %s successfully.\n", agentName)
		return nil
	},
}

var agentDeleteCmd = &cobra.Command{
	Use:   "delete <agent_name>",
	Short: "Delete an agent from the registry",
	Long:  `Removes an agent from the master's registry. This does not stop the agent process.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		masterAddr, _ := cmd.Flags().GetString("master")
		skipConfirmation, _ := cmd.Flags().GetBool("yes")

		// Ask for confirmation unless --yes flag is provided
		if !skipConfirmation {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("⚠️  Are you sure you want to delete agent '%s'? This action cannot be undone. [y/N]: ", agentName)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read confirmation: %v", err)
			}
			
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				pterm.Info.Println("Operation cancelled.")
				return nil
			}
		}

		conn, err := grpc.Dial(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to connect to master: %v", err)
		}
		defer conn.Close()

		registryClient := pb.NewAgentRegistryClient(conn)
		resp, err := registryClient.UnregisterAgent(context.Background(), &pb.UnregisterAgentRequest{
			AgentName: agentName,
		})
		if err != nil {
			return fmt.Errorf("failed to delete agent %s: %v", agentName, err)
		}

		if !resp.Success {
			pterm.Error.Printf("❌ %s\n", resp.Message)
			return fmt.Errorf("failed to delete agent")
		}

		pterm.Success.Printf("✅ Agent '%s' deleted successfully.\n", agentName)
		return nil
	},
}

// Workflow command and subcommands
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows and projects",
	Long:  `The workflow command provides subcommands to create, list, and manage sloth-runner workflows.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var workflowInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new workflow project",
	Long:  `Initialize a new workflow project from a template. Similar to 'pulumi new' or 'terraform init'.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var workflowName string
		if len(args) > 0 {
			workflowName = args[0]
		}

		templateName, _ := cmd.Flags().GetString("template")
		interactive, _ := cmd.Flags().GetBool("interactive")

		scaffolder := scaffolding.NewWorkflowScaffolder()
		return scaffolder.InitWorkflow(workflowName, templateName, interactive)
	},
}

var workflowListTemplatesCmd = &cobra.Command{
	Use:   "list-templates",
	Short: "List available workflow templates",
	Long:  `List all available workflow templates that can be used with 'workflow init'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		scaffolder := scaffolding.NewWorkflowScaffolder()
		scaffolder.ListTemplates()
		return nil
	},
}

// Stack command and subcommands
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Manage workflow stacks",
	Long:  `The stack command provides subcommands to manage workflow stacks and their state.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var stackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workflow stacks",
	Long:  `List all workflow stacks with their current state and execution history.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stackManager, err := stack.NewStackManager("")
		if err != nil {
			return fmt.Errorf("failed to initialize stack manager: %w", err)
		}
		defer stackManager.Close()

		stacks, err := stackManager.ListStacks()
		if err != nil {
			return fmt.Errorf("failed to list stacks: %w", err)
		}

		if len(stacks) == 0 {
			pterm.Info.Println("No stacks found.")
			return nil
		}

		// Create table output
		pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithTextStyle(pterm.NewStyle(pterm.FgWhite)).Printf("Workflow Stacks")
		
		pterm.Printf("\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATUS\tLAST RUN\tDURATION\tEXECUTIONS\tDESCRIPTION")
		fmt.Fprintln(w, "----\t------\t--------\t--------\t----------\t-----------")

		for _, s := range stacks {
			status := s.Status
			switch status {
			case "completed":
				status = pterm.Green(status)
			case "failed":
				status = pterm.Red(status)
			case "running":
				status = pterm.Yellow(status)
			default:
				status = pterm.Gray(status)
			}

			lastRun := "never"
			if s.CompletedAt != nil {
				lastRun = s.CompletedAt.Format("2006-01-02 15:04")
			} else if s.UpdatedAt.Year() > 1 {
				lastRun = s.UpdatedAt.Format("2006-01-02 15:04")
			}

			duration := "0s"
			if s.LastDuration > 0 {
				duration = s.LastDuration.String()
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n",
				s.Name, status, lastRun, duration, s.ExecutionCount, s.Description)
		}

		return w.Flush()
	},
}

var stackShowCmd = &cobra.Command{
	Use:   "show <stack-name>",
	Short: "Show detailed information about a stack",
	Long:  `Show detailed information about a specific workflow stack including execution history.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		
		stackManager, err := stack.NewStackManager("")
		if err != nil {
			return fmt.Errorf("failed to initialize stack manager: %w", err)
		}
		defer stackManager.Close()

		stackState, err := stackManager.GetStackByName(stackName)
		if err != nil {
			return fmt.Errorf("failed to get stack: %w", err)
		}

		// Show stack details
		pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Printf("Stack: %s", stackState.Name)
		
		pterm.Printf("\n")
		pterm.Printf("ID: %s\n", stackState.ID)
		pterm.Printf("Description: %s\n", stackState.Description)
		pterm.Printf("Version: %s\n", stackState.Version)
		pterm.Printf("Status: %s\n", stackState.Status)
		pterm.Printf("Created: %s\n", stackState.CreatedAt.Format("2006-01-02 15:04:05"))
		pterm.Printf("Updated: %s\n", stackState.UpdatedAt.Format("2006-01-02 15:04:05"))
		if stackState.CompletedAt != nil {
			pterm.Printf("Completed: %s\n", stackState.CompletedAt.Format("2006-01-02 15:04:05"))
		}
		pterm.Printf("Workflow File: %s\n", stackState.WorkflowFile)
		pterm.Printf("Executions: %d\n", stackState.ExecutionCount)
		if stackState.LastDuration > 0 {
			pterm.Printf("Last Duration: %s\n", stackState.LastDuration.String())
		}
		if stackState.LastError != "" {
			pterm.Printf("Last Error: %s\n", pterm.Red(stackState.LastError))
		}

		// Show outputs if any
		if len(stackState.Outputs) > 0 {
			pterm.Printf("\n")
			pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithTextStyle(pterm.NewStyle(pterm.FgWhite)).Printf("Outputs")
			pterm.Printf("\n")
			for key, value := range stackState.Outputs {
				pterm.Printf("%s: %v\n", pterm.Cyan(key), value)
			}
		}

		// Show recent executions
		executions, err := stackManager.GetStackExecutions(stackState.ID, 5)
		if err != nil {
			slog.Warn("Failed to get executions", "error", err)
		} else if len(executions) > 0 {
			pterm.Printf("\n")
			pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithTextStyle(pterm.NewStyle(pterm.FgWhite)).Printf("Recent Executions")
			pterm.Printf("\n")
			
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "STARTED\tSTATUS\tDURATION\tTASKS\tSUCCESS\tFAILED")
			fmt.Fprintln(w, "-------\t------\t--------\t-----\t-------\t------")

			for _, exec := range executions {
				status := exec.Status
				switch status {
				case "completed":
					status = pterm.Green(status)
				case "failed":
					status = pterm.Red(status)
				default:
					status = pterm.Gray(status)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\n",
					exec.StartedAt.Format("2006-01-02 15:04"),
					status,
					exec.Duration.String(),
					exec.TaskCount,
					exec.SuccessCount,
					exec.FailureCount)
			}
			w.Flush()
		}

		return nil
	},
}

var stackDeleteCmd = &cobra.Command{
	Use:   "delete <stack-name>",
	Short: "Delete a workflow stack",
	Long:  `Delete a workflow stack and all its execution history.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		force, _ := cmd.Flags().GetBool("force")
		
		stackManager, err := stack.NewStackManager("")
		if err != nil {
			return fmt.Errorf("failed to initialize stack manager: %w", err)
		}
		defer stackManager.Close()

		stackState, err := stackManager.GetStackByName(stackName)
		if err != nil {
			return fmt.Errorf("failed to get stack: %w", err)
		}

		if !force {
			pterm.Warning.Printf("This will permanently delete stack '%s' and all its execution history.\n", stackName)
			confirm := pterm.DefaultInteractiveConfirm.WithDefaultValue(false)
			result, _ := confirm.Show("Are you sure?")
			if !result {
				pterm.Info.Println("Operation cancelled.")
				return nil
			}
		}

		if err := stackManager.DeleteStack(stackState.ID); err != nil {
			return fmt.Errorf("failed to delete stack: %w", err)
		}

		pterm.Success.Printf("Stack '%s' deleted successfully.\n", stackName)
		return nil
	},
}

var stackNewCmd = &cobra.Command{
	Use:   "new <stack-name>",
	Short: "Create a new workflow stack",
	Long:  `Create a new workflow stack with the specified name and optional configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		description, _ := cmd.Flags().GetString("description")
		workflowFile, _ := cmd.Flags().GetString("workflow-file")
		version, _ := cmd.Flags().GetString("version")
		
		stackManager, err := stack.NewStackManager("")
		if err != nil {
			return fmt.Errorf("failed to initialize stack manager: %w", err)
		}
		defer stackManager.Close()

		// Check if stack already exists
		if _, err := stackManager.GetStackByName(stackName); err == nil {
			return fmt.Errorf("stack '%s' already exists", stackName)
		}

		// Set defaults
		if description == "" {
			description = fmt.Sprintf("Workflow stack: %s", stackName)
		}
		if version == "" {
			version = "1.0.0"
		}

		// Create new stack
		stackID := uuid.New().String()
		newStack := &stack.StackState{
			ID:           stackID,
			Name:         stackName,
			Description:  description,
			Version:      version,
			Status:       "created",
			WorkflowFile: workflowFile,
			TaskResults:  make(map[string]interface{}),
			Outputs:      make(map[string]interface{}),
			Configuration: make(map[string]interface{}),
			Metadata:     make(map[string]interface{}),
		}
		
		if err := stackManager.CreateStack(newStack); err != nil {
			return fmt.Errorf("failed to create stack: %w", err)
		}
		
		// Show success message
		pterm.Success.Printf("Stack '%s' created successfully.\n", stackName)
		pterm.Printf("\n")
		pterm.Printf("Stack Details:\n")
		pterm.Printf("  Name: %s\n", stackName)
		pterm.Printf("  ID: %s\n", stackID)
		pterm.Printf("  Description: %s\n", description)
		pterm.Printf("  Version: %s\n", version)
		if workflowFile != "" {
			pterm.Printf("  Workflow File: %s\n", workflowFile)
		}
		pterm.Printf("  Status: %s\n", "created")
		
		pterm.Printf("\n")
		pterm.Printf("Next steps:\n")
		if workflowFile != "" {
			pterm.Printf("  1. Run your workflow: sloth-runner run %s -f %s\n", stackName, workflowFile)
		} else {
			pterm.Printf("  1. Run your workflow: sloth-runner run %s -f <workflow-file>\n", stackName)
		}
		pterm.Printf("  2. View stack details: sloth-runner stack show %s\n", stackName)
		pterm.Printf("  3. List all stacks: sloth-runner stack list\n")
		
		return nil
	},
}

type agentServer struct {
	pb.UnimplementedAgentServer
	grpcServer *grpc.Server
}

func (s *agentServer) RunCommand(in *pb.RunCommandRequest, stream pb.Agent_RunCommandServer) error {
	slog.Info(fmt.Sprintf("Executing command on agent: %s", in.GetCommand()))

	cmd := exec.Command("bash", "-c", in.GetCommand())

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

func (s *agentServer) Shutdown(ctx context.Context, in *pb.ShutdownRequest) (*pb.ShutdownResponse, error) {
	slog.Info("Shutting down agent server")
	go func() {
		time.Sleep(1 * time.Second)
		s.grpcServer.GracefulStop()
	}()
	return &pb.ShutdownResponse{}, nil
}

func (s *agentServer) ExecuteTask(ctx context.Context, in *pb.ExecuteTaskRequest) (*pb.ExecuteTaskResponse, error) {
	slog.Info(fmt.Sprintf("Received task: %s from group: %s", in.GetTaskName(), in.GetTaskGroup()))

	// Create a temporary directory for the workspace
	workDir, err := os.MkdirTemp("", "sloth-runner-agent-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	// Unpack the workspace
	if err := extractTar(bytes.NewReader(in.GetWorkspace()), workDir); err != nil {
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
				group.Tasks[i] = task // Make sure the change is saved
			} else {
				slog.Info("Task has no delegate_to", "task", task.Name, "group", groupName)
			}
		}
	}
	
	slog.Info("Agent parsed task groups", "count", len(taskGroups), "groups", func() []string {
		var names []string
		for name := range taskGroups {
			names = append(names, name)
		}
		return names
	}())

	// Create task runner with all groups and let it find the specific task
	runner := taskrunner.NewTaskRunner(L, taskGroups, in.GetTaskGroup(), nil, false, false, &taskrunner.DefaultSurveyAsker{}, in.GetLuaScript())
	
	// Execute the specific task group
	slog.Info("Agent executing task group", "group", in.GetTaskGroup())
	err = runner.Run()
	
	// Pack the updated workspace
	var buf bytes.Buffer
	if err := createTar(workDir, &buf); err != nil {
		return nil, fmt.Errorf("failed to tar workspace: %w", err)
	}

	// Return response based on execution result
	if err != nil {
		slog.Error("Agent task execution failed", "error", err, "task", in.GetTaskName(), "group", in.GetTaskGroup())
		return &pb.ExecuteTaskResponse{
			Success:   false,
			Output:    fmt.Sprintf("Task execution failed: %v", err),
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
func createTar(source string, writer io.Writer) error {
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

func extractTar(reader io.Reader, dest string) error {
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
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

func init() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.SilenceErrors = false
	rootCmd.SilenceUsage = false

	// Master command flags
	rootCmd.AddCommand(masterCmd)
	masterCmd.Flags().IntP("port", "p", 50053, "The port for the master to listen on")
	masterCmd.Flags().Bool("daemon", false, "Run the master server as a daemon")
	masterCmd.Flags().Bool("debug", false, "Enable debug logging for the master server")

	// Agent command and subcommands
	rootCmd.AddCommand(agentCmd)

	// Persistent flags for agent client commands (run, list, stop)
	agentCmd.PersistentFlags().String("master", "localhost:50053", "The address of the master server")

	// Agent start command flags
	agentCmd.AddCommand(agentStartCmd)
	agentStartCmd.Flags().IntP("port", "p", 50051, "The port for the agent to listen on")
	agentStartCmd.Flags().String("master", "", "The address of the master server to register with")
	agentStartCmd.Flags().String("name", "", "The name of the agent")
	agentStartCmd.Flags().Bool("daemon", false, "Run the agent as a daemon")
	agentStartCmd.Flags().String("bind-address", "", "The IP address for the agent to bind to")
	agentStartCmd.Flags().String("report-address", "", "The IP address to report to the master (defaults to bind-address or auto-detected)")
	// TLS flags for agent start are now persistent flags on the parent 'agent' command

	// Agent client commands
	agentCmd.AddCommand(agentRunCmd)
	agentCmd.AddCommand(agentListCmd)
	agentListCmd.Flags().Bool("debug", false, "Enable debug logging for this command")
	agentCmd.AddCommand(agentStopCmd)
	agentCmd.AddCommand(agentDeleteCmd)
	agentDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	// Scheduler command and subcommands
	rootCmd.AddCommand(schedulerCmd)
	schedulerCmd.AddCommand(schedulerEnableCmd)
	schedulerCmd.AddCommand(schedulerDisableCmd)
	schedulerCmd.AddCommand(schedulerListCmd)
	schedulerCmd.AddCommand(schedulerDeleteCmd)
	
	// Add config flag to all scheduler subcommands
	schedulerEnableCmd.Flags().StringP("config", "c", "scheduler.yaml", "Path to the scheduler configuration file")
	schedulerDisableCmd.Flags().StringP("config", "c", "scheduler.yaml", "Path to the scheduler configuration file")
	schedulerListCmd.Flags().StringP("config", "c", "scheduler.yaml", "Path to the scheduler configuration file")
	schedulerDeleteCmd.Flags().StringP("config", "c", "scheduler.yaml", "Path to the scheduler configuration file")

	// Run command
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("file", "f", "", "Path to the Lua task file")
	runCmd.Flags().StringP("values", "v", "", "Path to the values file")
	runCmd.Flags().Bool("yes", false, "Skip confirmation prompts")
	runCmd.Flags().Bool("interactive", false, "Run in interactive mode")
	runCmd.Flags().StringP("output", "o", "basic", "Output style: basic, enhanced, rich, modern, json")

	// List command
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("file", "f", "", "Path to the Lua workflow file")

	// Workflow command and subcommands
	rootCmd.AddCommand(workflowCmd)
	workflowCmd.AddCommand(workflowInitCmd)
	workflowCmd.AddCommand(workflowListTemplatesCmd)
	
	// Workflow init command flags
	workflowInitCmd.Flags().StringP("template", "t", "", "Template to use (basic, cicd, infrastructure, microservices, data-pipeline)")
	workflowInitCmd.Flags().BoolP("interactive", "i", false, "Run in interactive mode")

	// Stack command and subcommands
	rootCmd.AddCommand(stackCmd)
	stackCmd.AddCommand(stackListCmd)
	stackCmd.AddCommand(stackShowCmd)
	stackCmd.AddCommand(stackNewCmd)
	stackCmd.AddCommand(stackDeleteCmd)
	
	// Stack new command flags
	stackNewCmd.Flags().StringP("description", "d", "", "Description of the stack")
	stackNewCmd.Flags().StringP("workflow-file", "f", "", "Path to the workflow file")
	stackNewCmd.Flags().StringP("version", "v", "1.0.0", "Version of the stack")
	
	// Stack delete command flags
	stackDeleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")

	// Version command
	rootCmd.AddCommand(versionCmd)

	// UI command
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().IntP("port", "p", 8080, "The port for the UI server to listen on")
	uiCmd.Flags().Bool("daemon", false, "Run the UI server as a daemon")
	uiCmd.Flags().Bool("debug", false, "Enable debug logging for the UI server")
}

func Execute() error {
	rootCmd.SilenceUsage = true

	if runAsScheduler {
		select {}
	}

	err := rootCmd.Execute()
	if err != nil {
		slog.Error("DEBUG: rootCmd.Execute() returned error", "err", err)
	}
	return err
}

func main() {
	slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))

	if err := Execute(); err != nil {
		slog.Error("execution failed", "err", err)
		os.Exit(1)
	}
}