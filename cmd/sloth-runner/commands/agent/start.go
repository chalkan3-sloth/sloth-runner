package agent

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	agentInternal "github.com/chalkan3-sloth/sloth-runner/internal/agent"
	"github.com/chalkan3-sloth/sloth-runner/internal/telemetry"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// NewStartCommand creates the agent start command
func NewStartCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
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
			telemetryEnabled, _ := cmd.Flags().GetBool("telemetry")
			metricsPort, _ := cmd.Flags().GetInt("metrics-port")

			return startAgent(ctx, port, masterAddr, agentName, daemon, bindAddress, reportAddress, telemetryEnabled, metricsPort)
		},
	}

	cmd.Flags().Int("port", 50052, "Port for the agent to listen on")
	cmd.Flags().String("master", "localhost:50051", "Address of the master server")
	cmd.Flags().String("name", "default-agent", "Name of the agent")
	cmd.Flags().Bool("daemon", false, "Run agent as daemon in background")
	cmd.Flags().String("bind-address", "", "Address to bind to (default: all interfaces)")
	cmd.Flags().String("report-address", "", "Address to report to master (if different from bind)")
	cmd.Flags().Bool("telemetry", false, "Enable telemetry and metrics server")
	cmd.Flags().Int("metrics-port", 9090, "Port for metrics server")

	return cmd
}

func startAgent(ctx *commands.AppContext, port int, masterAddr, agentName string, daemon bool, bindAddress, reportAddress string, telemetryEnabled bool, metricsPort int) error {
	// Apply runtime optimizations for reduced resource usage
	configureAgentRuntimeOptimizations()

	if daemon {
		pidFile := filepath.Join("/tmp", fmt.Sprintf("sloth-runner-agent-%s.pid", agentName))
		if _, err := os.Stat(pidFile); err == nil {
			pidBytes, err := os.ReadFile(pidFile)
			if err == nil {
				pid, _ := strconv.Atoi(string(pidBytes))
				if process, err := os.FindProcess(pid); err == nil {
					if err := process.Signal(syscall.Signal(0)); err == nil {
						fmt.Printf("Agent %s is already running with PID %d.\n", agentName, pid)
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
		if telemetryEnabled {
			cmdArgs = append(cmdArgs, "--telemetry")
		}
		if metricsPort != 9090 {
			cmdArgs = append(cmdArgs, "--metrics-port", strconv.Itoa(metricsPort))
		}

		command := exec.Command(os.Args[0], cmdArgs...)
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
		fmt.Printf("Agent %s started with PID %d. Logs can be found at %s.\n", agentName, command.Process.Pid, "agent.log")
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
		agentReportAddress = reportAddress
		if !strings.Contains(reportAddress, ":") {
			agentReportAddress = fmt.Sprintf("%s:%d", reportAddress, port)
		}
	} else if bindAddress != "" {
		agentReportAddress = fmt.Sprintf("%s:%d", bindAddress, port)
	}

	pterm.Warning.Println("Starting agent in insecure mode.")

	if masterAddr != "" {
		// Start connection manager with reconnection logic
		go startMasterConnection(ctx, masterAddr, agentName, agentReportAddress)
	}

	// Initialize telemetry server
	telemetryServer := telemetry.InitGlobal(metricsPort, telemetryEnabled)
	if telemetryEnabled {
		if err := telemetryServer.Start(); err != nil {
			slog.Error("Failed to start telemetry server", "error", err)
		} else {
			telemetry.SetAgentInfo(ctx.Version, runtime.GOOS, runtime.GOARCH)
			pterm.Success.Printf("✓ Telemetry server started at %s\n", telemetryServer.GetEndpoint())
		}
	}

	// Create optimized gRPC server
	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(10),     // Limit concurrent streams
		grpc.ReadBufferSize(8192),         // 8KB read buffer (vs 32KB default)
		grpc.WriteBufferSize(8192),        // 8KB write buffer (vs 32KB default)
		grpc.MaxRecvMsgSize(4 * 1024 * 1024), // 4MB max message
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second, // Close idle connections faster
			MaxConnectionAge:      30 * time.Second, // Recycle connections
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  10 * time.Second, // Ping every 10s
			Timeout:               3 * time.Second,
		}),
	}

	s := grpc.NewServer(opts...)
	server := &agentServer{
		grpcServer:    s,
		cachedMetrics: &CachedMetrics{},
	}
	pb.RegisterAgentServer(s, server)

	pterm.Success.Printf("✓ Agent '%s' listening at %v\n", agentName, lis.Addr())
	pterm.Info.Println("Optimizations enabled: 30s metrics cache, batched DB writes, process list caching")

	slog.Info(fmt.Sprintf("Agent listening at %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

// configureAgentRuntimeOptimizations applies Go runtime optimizations for agents
func configureAgentRuntimeOptimizations() {
	// Limit to 2 CPU cores for agent (reduces CPU usage)
	runtime.GOMAXPROCS(2)

	// More aggressive GC (50% vs 100% default)
	debug.SetGCPercent(50)

	// Set memory limit to 100MB (agents should stay well under this)
	debug.SetMemoryLimit(100 * 1024 * 1024)

	slog.Info("Agent runtime optimizations applied",
		"max_procs", 2,
		"gc_percent", 50,
		"memory_limit_mb", 100)
}

func startMasterConnection(ctx *commands.AppContext, masterAddr, agentName, agentReportAddress string) {
	reconnectDelay := 5 * time.Second
	maxReconnectDelay := 60 * time.Second
	heartbeatInterval := 5 * time.Second

	for {
		// Create connection context with timeout
		connCtx, connCancel := context.WithTimeout(context.Background(), 10*time.Second)
		conn, err := grpc.DialContext(connCtx, masterAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		connCancel()

		if err != nil {
			slog.Error(fmt.Sprintf("Failed to connect to master at %s: %v. Retrying in %v...", masterAddr, err, reconnectDelay))
			pterm.Warning.Printf("⚠ Cannot connect to master at %s. Retrying in %v...\n", masterAddr, reconnectDelay)
			time.Sleep(reconnectDelay)
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
		regCtx, regCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err = registryClient.RegisterAgent(regCtx, &pb.RegisterAgentRequest{
			AgentName:    agentName,
			AgentAddress: agentReportAddress,
		})
		regCancel()

		if err != nil {
			slog.Error(fmt.Sprintf("Failed to register with master: %v. Reconnecting...", err))
			pterm.Warning.Printf("⚠ Failed to register with master: %v\n", err)
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
		heartbeatCounter := 0
		sysInfoCollectInterval := 12 // Collect system info every 12 heartbeats (60 seconds)

		for connected {
			time.Sleep(heartbeatInterval)
			heartbeatCounter++

			// Collect system info periodically (every minute)
			var sysInfoJSON string
			if heartbeatCounter%sysInfoCollectInterval == 0 {
				if sysInfo, err := agentInternal.CollectSystemInfo(); err == nil {
					if jsonStr, err := sysInfo.ToJSON(); err == nil {
						sysInfoJSON = jsonStr
						slog.Debug("System info collected and will be sent with heartbeat")
					}
				}
			}

			hbCtx, hbCancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := registryClient.Heartbeat(hbCtx, &pb.HeartbeatRequest{
				AgentName:      agentName,
				SystemInfoJson: sysInfoJSON,
				Version:        ctx.Version,
			})
			hbCancel()

			if err != nil {
				consecutiveFailures++
				slog.Warn(fmt.Sprintf("Heartbeat failed (%d/%d): %v", consecutiveFailures, maxConsecutiveFailures, err))

				if consecutiveFailures >= maxConsecutiveFailures {
					slog.Error(fmt.Sprintf("Lost connection to master after %d failed heartbeats. Reconnecting...", maxConsecutiveFailures))
					pterm.Warning.Printf("⚠ Connection to master lost. Attempting to reconnect...\n")
					connected = false
				}
			} else {
				if consecutiveFailures > 0 {
					consecutiveFailures = 0
					slog.Info("Heartbeat recovered, connection stable")
					pterm.Success.Printf("✓ Connection to master recovered\n")
				}
			}
		}

		// Close old connection before reconnecting
		conn.Close()
		slog.Info("Closed connection to master, preparing to reconnect")

		// Wait before attempting reconnection
		pterm.Info.Printf("🔄 Reconnecting to master in %v...\n", reconnectDelay)
		time.Sleep(reconnectDelay)
	}
}
