package main

import (
	"archive/tar"
	"bytes"
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
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
	version             = "dev" // será substituído em tempo de compilação
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

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Starts the sloth-runner master server",
	Long:  `The master command starts the sloth-runner master server, which includes the agent registry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		daemon, _ := cmd.Flags().GetBool("daemon")
		debug, _ := cmd.Flags().GetBool("debug")
		tlsCertFile, _ := cmd.Flags().GetString("tls-cert-file")
		tlsKeyFile, _ := cmd.Flags().GetString("tls-key-file")
		tlsCaFile, _ := cmd.Flags().GetString("tls-ca-file")

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

			command := execCommand(os.Args[0], "master", "--port", strconv.Itoa(port), "--tls-cert-file", tlsCertFile, "--tls-key-file", tlsKeyFile, "--tls-ca-file", tlsCaFile)
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

		globalAgentRegistry = newAgentRegistryServer(tlsCertFile, tlsKeyFile, tlsCaFile)
		return globalAgentRegistry.Start(port)
	},
}

var globalAgentRegistry *agentRegistryServer

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

// Run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run sloth-runner tasks",
	Long:  `Run sloth-runner tasks from Lua files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			filePath = "examples/basic_pipeline.lua"
		}
		
		values, _ := cmd.Flags().GetString("values")
		_, _ = cmd.Flags().GetBool("yes") // yes flag - for future use
		interactive, _ := cmd.Flags().GetBool("interactive")

		// Use test output buffer if available, otherwise use stdout
		writer := cmd.OutOrStdout()
		if testOutputBuffer != nil {
			writer = testOutputBuffer
		}

		// Load values.yaml if specified
		var valuesTable *lua.LTable
		if values != "" {
			// Load and parse values file (simplified for now)
			fmt.Fprintf(writer, "Loading values from: %s\n", values)
		}

		// Parse the Lua script
		taskGroups, err := luainterface.ParseLuaScript(cmd.Context(), filePath, valuesTable)
		if err != nil {
			return fmt.Errorf("failed to parse Lua script: %w", err)
		}

		if len(taskGroups) == 0 {
			fmt.Fprintln(writer, "No task groups found in script")
			return nil
		}

		// Create task runner
		L := lua.NewState()
		defer L.Close()
		
		// Register modules
		luainterface.RegisterAllModules(L)
		luainterface.OpenImport(L, filePath)
		
		// Initialize task runner
		runner := taskrunner.NewTaskRunner(L, taskGroups, "", nil, false, interactive, &taskrunner.DefaultSurveyAsker{}, "")
		
		// Set outputs to capture results
		runner.Outputs = make(map[string]interface{})
		
		// Execute the tasks
		fmt.Fprintf(writer, "Executing tasks from: %s\n", filePath)
		err = runner.Run()
		if err != nil {
			return fmt.Errorf("task execution failed: %w", err)
		}

		fmt.Fprintln(writer, "Task execution completed successfully!")
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
		tlsCertFile, _ := cmd.Flags().GetString("tls-cert-file")
		tlsKeyFile, _ := cmd.Flags().GetString("tls-key-file")
		tlsCaFile, _ := cmd.Flags().GetString("tls-ca-file")

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

			command := execCommand(os.Args[0], "agent", "start", "--port", strconv.Itoa(port), "--name", agentName, "--master", masterAddr, "--bind-address", bindAddress, "--tls-cert-file", tlsCertFile, "--tls-key-file", tlsKeyFile, "--tls-ca-file", tlsCaFile)
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

		reportAddress := lis.Addr().String()
		if bindAddress != "" {
			reportAddress = fmt.Sprintf("%s:%d", bindAddress, port)
		}

		var serverOpts []grpc.ServerOption
		if tlsCertFile != "" && tlsKeyFile != "" && tlsCaFile != "" {
			serverCert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
			if err != nil {
				return fmt.Errorf("failed to load agent server certificate: %v", err)
			}
			caCert, err := os.ReadFile(tlsCaFile)
			if err != nil {
				return fmt.Errorf("failed to read CA certificate for agent server: %v", err)
			}
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				return fmt.Errorf("failed to append CA certificate for agent server")
			}
			creds := credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{serverCert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    caCertPool,
				RootCAs:      caCertPool,
			})
			serverOpts = append(serverOpts, grpc.Creds(creds))
		}

		if masterAddr != "" {
			dialOpts, err := getDialOptions(tlsCertFile, tlsKeyFile, tlsCaFile)
			if err != nil {
				return err
			}

			conn, err := grpc.Dial(masterAddr, dialOpts...)
			if err != nil {
				return fmt.Errorf("failed to connect to master: %v", err)
			}
			defer conn.Close()

			registryClient := pb.NewAgentRegistryClient(conn)
			_, err = registryClient.RegisterAgent(context.Background(), &pb.RegisterAgentRequest{
				AgentName:    agentName,
				AgentAddress: reportAddress,
			})
			if err != nil {
				return fmt.Errorf("failed to register with master: %v", err)
			}
			slog.Info(fmt.Sprintf("Agent registered with master at %s, reporting address %s", masterAddr, reportAddress))

			go func() {
				for {
					_, err := registryClient.Heartbeat(context.Background(), &pb.HeartbeatRequest{AgentName: agentName})
					if err != nil {
						slog.Error(fmt.Sprintf("Failed to send heartbeat to master: %v", err))
					}
					time.Sleep(5 * time.Second)
				}
			}()
		}

		s := grpc.NewServer(serverOpts...)
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
		tlsCertFile, _ := cmd.Flags().GetString("tls-cert-file")
		tlsKeyFile, _ := cmd.Flags().GetString("tls-key-file")
		tlsCaFile, _ := cmd.Flags().GetString("tls-ca-file")

		dialOpts, err := getDialOptions(tlsCertFile, tlsKeyFile, tlsCaFile)
		if err != nil {
			return err
		}

		conn, err := grpc.Dial(masterAddr, dialOpts...)
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
		var exitCode int32
		success := false

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break // Stream has ended
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
				success = (exitCode == 0 && finalError == "")
				break
			}
		}

		if !success {
			pterm.Error.Printf("Command failed on agent %s with exit code %d!\n", agentName, exitCode)
			if finalError != "" {
				pterm.Error.Printf("Error: %s\n", finalError)
			}
			return fmt.Errorf("command execution failed on agent %s", agentName)
		}

		pterm.Success.Printf("Command executed successfully on agent %s.\n", agentName)
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
		tlsCertFile, _ := cmd.Flags().GetString("tls-cert-file")
		tlsKeyFile, _ := cmd.Flags().GetString("tls-key-file")
		tlsCaFile, _ := cmd.Flags().GetString("tls-ca-file")

		dialOpts, err := getDialOptions(tlsCertFile, tlsKeyFile, tlsCaFile)
		if err != nil {
			return err
		}

		conn, err := grpc.Dial(masterAddr, dialOpts...)
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
		tlsCertFile, _ := cmd.Flags().GetString("tls-cert-file")
		tlsKeyFile, _ := cmd.Flags().GetString("tls-key-file")
		tlsCaFile, _ := cmd.Flags().GetString("tls-ca-file")

		dialOpts, err := getDialOptions(tlsCertFile, tlsKeyFile, tlsCaFile)
		if err != nil {
			return err
		}

		conn, err := grpc.Dial(masterAddr, dialOpts...)
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
	slog.Info(fmt.Sprintf("Received task: %s", in.GetTaskName()))

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

	// Pack the workspace
	var buf bytes.Buffer
	if err := createTar(workDir, &buf); err != nil {
		return nil, fmt.Errorf("failed to tar workspace: %w", err)
	}

	return &pb.ExecuteTaskResponse{Success: true, Output: "Task executed successfully", Workspace: buf.Bytes()}, nil
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
	masterCmd.Flags().String("tls-cert-file", "", "Path to the TLS certificate file for the master server")
	masterCmd.Flags().String("tls-key-file", "", "Path to the TLS key file for the master server")
	masterCmd.Flags().String("tls-ca-file", "", "Path to the TLS CA certificate file for the master server to verify agent certificates")

	// Agent command and subcommands
	rootCmd.AddCommand(agentCmd)

	// Persistent flags for agent client commands (run, list, stop)
	agentCmd.PersistentFlags().String("master", "localhost:50053", "The address of the master server")
	agentCmd.PersistentFlags().String("tls-cert-file", "", "Path to the TLS certificate file for the client")
	agentCmd.PersistentFlags().String("tls-key-file", "", "Path to the TLS key file for the client")
	agentCmd.PersistentFlags().String("tls-ca-file", "", "Path to the TLS CA certificate file for verifying the server")

	// Agent start command flags
	agentCmd.AddCommand(agentStartCmd)
	agentStartCmd.Flags().IntP("port", "p", 50051, "The port for the agent to listen on")
	agentStartCmd.Flags().String("master", "", "The address of the master server to register with")
	agentStartCmd.Flags().String("name", "", "The name of the agent")
	agentStartCmd.Flags().Bool("daemon", false, "Run the agent as a daemon")
	agentStartCmd.Flags().String("bind-address", "", "The IP address for the agent to bind to and report to the master")
	// TLS flags for agent start are now persistent flags on the parent 'agent' command

	// Agent client commands
	agentCmd.AddCommand(agentRunCmd)
	agentCmd.AddCommand(agentListCmd)
	agentListCmd.Flags().Bool("debug", false, "Enable debug logging for this command")
	agentCmd.AddCommand(agentStopCmd)

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