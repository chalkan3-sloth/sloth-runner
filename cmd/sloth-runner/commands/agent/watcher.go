package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	agentInternal "github.com/chalkan3-sloth/sloth-runner/internal/agent"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewWatcherCommand creates the parent watcher command
func NewWatcherCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watcher",
		Short: "Manage watchers on agents",
		Long:  `Commands to list, get, create, and delete watchers running on agents`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add all watcher subcommands
	cmd.AddCommand(
		NewWatcherListCommand(ctx),
		NewWatcherGetCommand(ctx),
		NewWatcherDeleteCommand(ctx),
		NewWatcherCreateCommand(ctx),
	)

	return cmd
}

// NewWatcherListCommand lists all watchers on an agent
func NewWatcherListCommand(ctx *commands.AppContext) *cobra.Command {
	var (
		agentName string
		format    string
		watcherType string
	)

	cmd := &cobra.Command{
		Use:   "list [agent-name]",
		Short: "List all watchers on an agent",
		Long:  `Lists all active watchers running on a specific agent`,
		Example: `  # List all watchers on agent 'my-agent'
  sloth-runner agent watcher list my-agent

  # List watchers in JSON format
  sloth-runner agent watcher list my-agent --format json

  # List only file watchers
  sloth-runner agent watcher list my-agent --type file`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				agentName = args[0]
			}

			if agentName == "" {
				return fmt.Errorf("agent name is required")
			}

			// Get agent address
			agentAddr, err := resolveAgentAddress(agentName)
			if err != nil {
				return fmt.Errorf("failed to resolve agent address: %w", err)
			}

			// Connect to agent
			conn, err := createGRPCConnection(agentAddr)
			if err != nil {
				return fmt.Errorf("failed to connect to agent: %w", err)
			}
			defer conn.Close()

			client := pb.NewAgentClient(conn)

			// List watchers
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.ListWatchers(ctx, &pb.ListWatchersRequest{})
			if err != nil {
				return fmt.Errorf("failed to list watchers: %w", err)
			}

			// Filter by type if specified
			watchers := resp.Watchers
			if watcherType != "" {
				filtered := []*pb.WatcherConfig{}
				for _, w := range watchers {
					if w.Type == watcherType {
						filtered = append(filtered, w)
					}
				}
				watchers = filtered
			}

			// Display results
			switch format {
			case "json":
				data, _ := json.MarshalIndent(watchers, "", "  ")
				fmt.Println(string(data))
			case "table":
				displayWatcherTable(watchers)
			default:
				displayWatcherTable(watchers)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&agentName, "agent", "a", "", "Agent name")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json)")
	cmd.Flags().StringVarP(&watcherType, "type", "t", "", "Filter by watcher type (file, cpu, memory, process, etc)")

	return cmd
}

// NewWatcherGetCommand gets details of a specific watcher
func NewWatcherGetCommand(ctx *commands.AppContext) *cobra.Command {
	var (
		agentName string
		format    string
	)

	cmd := &cobra.Command{
		Use:   "get <watcher-id> [agent-name]",
		Short: "Get details of a specific watcher",
		Long:  `Gets detailed information about a specific watcher by ID`,
		Example: `  # Get watcher details
  sloth-runner agent watcher get abc123 my-agent

  # Get watcher details in JSON format
  sloth-runner agent watcher get abc123 my-agent --format json`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			watcherID := args[0]

			if len(args) > 1 {
				agentName = args[1]
			}

			if agentName == "" {
				return fmt.Errorf("agent name is required")
			}

			// Get agent address
			agentAddr, err := resolveAgentAddress(agentName)
			if err != nil {
				return fmt.Errorf("failed to resolve agent address: %w", err)
			}

			// Connect to agent
			conn, err := createGRPCConnection(agentAddr)
			if err != nil {
				return fmt.Errorf("failed to connect to agent: %w", err)
			}
			defer conn.Close()

			client := pb.NewAgentClient(conn)

			// Get watcher
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.GetWatcher(ctx, &pb.GetWatcherRequest{
				WatcherId: watcherID,
			})
			if err != nil {
				return fmt.Errorf("failed to get watcher: %w", err)
			}

			// Display result
			switch format {
			case "json":
				data, _ := json.MarshalIndent(resp.Watcher, "", "  ")
				fmt.Println(string(data))
			default:
				displayWatcherDetails(resp.Watcher)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&agentName, "agent", "a", "", "Agent name")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json)")

	return cmd
}

// NewWatcherDeleteCommand deletes a watcher
func NewWatcherDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	var (
		agentName string
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "delete <watcher-id> [agent-name]",
		Short: "Delete a watcher",
		Long:  `Removes a watcher from an agent, stopping its monitoring`,
		Example: `  # Delete a watcher
  sloth-runner agent watcher delete abc123 my-agent

  # Delete without confirmation
  sloth-runner agent watcher delete abc123 my-agent --force`,
		Aliases: []string{"remove", "rm"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			watcherID := args[0]

			if len(args) > 1 {
				agentName = args[1]
			}

			if agentName == "" {
				return fmt.Errorf("agent name is required")
			}

			// Confirm deletion if not forced
			if !force {
				result, _ := pterm.DefaultInteractiveConfirm.Show(
					fmt.Sprintf("Are you sure you want to delete watcher '%s' on agent '%s'?", watcherID, agentName),
				)
				if !result {
					pterm.Info.Println("Deletion cancelled")
					return nil
				}
			}

			// Get agent address
			agentAddr, err := resolveAgentAddress(agentName)
			if err != nil {
				return fmt.Errorf("failed to resolve agent address: %w", err)
			}

			// Connect to agent
			conn, err := createGRPCConnection(agentAddr)
			if err != nil {
				return fmt.Errorf("failed to connect to agent: %w", err)
			}
			defer conn.Close()

			client := pb.NewAgentClient(conn)

			// Delete watcher
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err = client.RemoveWatcher(ctx, &pb.RemoveWatcherRequest{
				WatcherId: watcherID,
			})
			if err != nil {
				return fmt.Errorf("failed to delete watcher: %w", err)
			}

			pterm.Success.Printf("Watcher '%s' deleted from agent '%s'\n", watcherID, agentName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&agentName, "agent", "a", "", "Agent name")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")

	return cmd
}

// NewWatcherCreateCommand creates a new watcher
func NewWatcherCreateCommand(ctx *commands.AppContext) *cobra.Command {
	var (
		agentName     string
		watcherType   string
		filePath      string
		processName   string
		port          int
		serviceName   string
		threshold     float64
		interval      string
		conditions    []string
		checkHash     bool
	)

	cmd := &cobra.Command{
		Use:   "create [agent-name]",
		Short: "Create a new watcher on an agent",
		Long:  `Creates and registers a new watcher on the specified agent`,
		Example: `  # Create a file watcher
  sloth-runner agent watcher create my-agent --type file --path /tmp/test.txt --when created,changed,deleted

  # Create a CPU watcher
  sloth-runner agent watcher create my-agent --type cpu --threshold 80 --when above

  # Create a memory watcher
  sloth-runner agent watcher create my-agent --type memory --threshold 75 --when above

  # Create a process watcher
  sloth-runner agent watcher create my-agent --type process --process nginx --when created,deleted`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				agentName = args[0]
			}

			if agentName == "" {
				return fmt.Errorf("agent name is required")
			}

			if watcherType == "" {
				return fmt.Errorf("watcher type is required (--type)")
			}

			if len(conditions) == 0 {
				return fmt.Errorf("at least one condition is required (--when)")
			}

			// Parse interval
			var intervalDuration time.Duration
			if interval != "" {
				var err error
				intervalDuration, err = time.ParseDuration(interval)
				if err != nil {
					return fmt.Errorf("invalid interval: %w", err)
				}
			} else {
				intervalDuration = 5 * time.Second // default
			}

			// Build watcher config
			config := &agentInternal.WatcherConfig{
				ID:         generateWatcherID(),
				Type:       agentInternal.WatcherType(watcherType),
				Interval:   intervalDuration,
				Conditions: parseConditions(conditions),
			}

			// Type-specific configuration
			switch watcherType {
			case "file", "directory":
				if filePath == "" {
					return fmt.Errorf("file path is required for file/directory watcher (--path)")
				}
				config.FilePath = filePath
				config.CheckHash = checkHash

			case "process":
				if processName == "" {
					return fmt.Errorf("process name is required for process watcher (--process)")
				}
				config.ProcessName = processName

			case "port":
				if port == 0 {
					return fmt.Errorf("port is required for port watcher (--port)")
				}
				config.Port = port

			case "service":
				if serviceName == "" {
					return fmt.Errorf("service name is required for service watcher (--service)")
				}
				config.ServiceName = serviceName

			case "cpu":
				if threshold == 0 {
					return fmt.Errorf("threshold is required for CPU watcher (--threshold)")
				}
				config.CPUThreshold = threshold

			case "memory":
				if threshold == 0 {
					return fmt.Errorf("threshold is required for memory watcher (--threshold)")
				}
				config.MemoryThreshold = threshold

			default:
				return fmt.Errorf("unsupported watcher type: %s", watcherType)
			}

			// Get agent address
			agentAddr, err := resolveAgentAddress(agentName)
			if err != nil {
				return fmt.Errorf("failed to resolve agent address: %w", err)
			}

			// Connect to agent
			conn, err := createGRPCConnection(agentAddr)
			if err != nil {
				return fmt.Errorf("failed to connect to agent: %w", err)
			}
			defer conn.Close()

			client := pb.NewAgentClient(conn)

			// Register watcher
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err = client.RegisterWatcher(ctx, watcherConfigToProto(config))
			if err != nil {
				return fmt.Errorf("failed to create watcher: %w", err)
			}

			pterm.Success.Printf("Watcher '%s' created on agent '%s'\n", config.ID, agentName)
			pterm.Info.Printf("Type: %s, Interval: %s, Conditions: %v\n", config.Type, intervalDuration, conditions)

			return nil
		},
	}

	cmd.Flags().StringVarP(&agentName, "agent", "a", "", "Agent name")
	cmd.Flags().StringVarP(&watcherType, "type", "t", "", "Watcher type (file, cpu, memory, process, port, service)")
	cmd.Flags().StringVarP(&filePath, "path", "p", "", "File/directory path (for file/directory watchers)")
	cmd.Flags().StringVar(&processName, "process", "", "Process name (for process watchers)")
	cmd.Flags().IntVar(&port, "port", 0, "Port number (for port watchers)")
	cmd.Flags().StringVar(&serviceName, "service", "", "Service name (for service watchers)")
	cmd.Flags().Float64Var(&threshold, "threshold", 0, "Threshold value (for CPU/memory watchers)")
	cmd.Flags().StringVarP(&interval, "interval", "i", "5s", "Check interval (e.g., 5s, 1m, 10m)")
	cmd.Flags().StringSliceVarP(&conditions, "when", "w", []string{}, "Conditions to trigger (comma-separated: created,changed,deleted,above,below)")
	cmd.Flags().BoolVar(&checkHash, "check-hash", false, "Check file hash for changes (file watchers only)")

	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("when")

	return cmd
}

// Helper functions

func displayWatcherTable(watchers []*pb.WatcherConfig) {
	if len(watchers) == 0 {
		pterm.Info.Println("No watchers found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTYPE\tCONDITIONS\tINTERVAL\tTARGET")
	fmt.Fprintln(w, "--\t----\t----------\t--------\t------")

	for _, watcher := range watchers {
		target := getWatcherTarget(watcher)
		conditions := strings.Join(watcher.Conditions, ",")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			truncateString(watcher.Id, 12),
			watcher.Type,
			conditions,
			watcher.Interval,
			truncateString(target, 40),
		)
	}

	w.Flush()
	fmt.Printf("\nTotal: %d watchers\n", len(watchers))
}

func displayWatcherDetails(watcher *pb.WatcherConfig) {
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
		WithFullWidth().
		Println("Watcher Details")

	data := [][]string{
		{"ID", watcher.Id},
		{"Type", watcher.Type},
		{"Conditions", strings.Join(watcher.Conditions, ", ")},
		{"Interval", watcher.Interval},
		{"Target", getWatcherTarget(watcher)},
	}

	// Add type-specific fields
	if watcher.FilePath != "" {
		data = append(data, []string{"File Path", watcher.FilePath})
		data = append(data, []string{"Check Hash", fmt.Sprintf("%v", watcher.CheckHash)})
	}
	if watcher.ProcessName != "" {
		data = append(data, []string{"Process Name", watcher.ProcessName})
	}
	if watcher.Port != 0 {
		data = append(data, []string{"Port", fmt.Sprintf("%d", watcher.Port)})
	}
	if watcher.CpuThreshold != 0 {
		data = append(data, []string{"CPU Threshold", fmt.Sprintf("%.2f%%", watcher.CpuThreshold)})
	}
	if watcher.MemoryThreshold != 0 {
		data = append(data, []string{"Memory Threshold", fmt.Sprintf("%.2f%%", watcher.MemoryThreshold)})
	}

	pterm.DefaultTable.WithHasHeader(false).WithData(data).Render()
}

func getWatcherTarget(watcher *pb.WatcherConfig) string {
	switch watcher.Type {
	case "file", "directory":
		return watcher.FilePath
	case "process":
		return watcher.ProcessName
	case "port":
		return fmt.Sprintf("port %d", watcher.Port)
	case "cpu":
		return fmt.Sprintf("%.1f%%", watcher.CpuThreshold)
	case "memory":
		return fmt.Sprintf("%.1f%%", watcher.MemoryThreshold)
	}
	return "-"
}

func parseConditions(conditions []string) []agentInternal.EventCondition {
	result := make([]agentInternal.EventCondition, 0, len(conditions))
	for _, c := range conditions {
		result = append(result, agentInternal.EventCondition(c))
	}
	return result
}

func generateWatcherID() string {
	return fmt.Sprintf("watcher-%d", time.Now().UnixNano())
}

func watcherConfigToProto(config *agentInternal.WatcherConfig) *pb.RegisterWatcherRequest {
	// Convert conditions
	conditions := make([]string, 0, len(config.Conditions))
	for _, c := range config.Conditions {
		conditions = append(conditions, string(c))
	}

	protoConfig := &pb.WatcherConfig{
		Id:         config.ID,
		Type:       string(config.Type),
		Conditions: conditions,
		Interval:   config.Interval.String(),

		// Type-specific fields
		FilePath:    config.FilePath,
		CheckHash:   config.CheckHash,
		Recursive:   config.Recursive,

		ProcessName: config.ProcessName,
		Pid:         int32(config.PID),

		Port:     int32(config.Port),
		Protocol: config.Protocol,

		CpuThreshold:    config.CPUThreshold,
		MemoryThreshold: config.MemoryThreshold,
		DiskThreshold:   config.DiskThreshold,
	}

	return &pb.RegisterWatcherRequest{
		Config: protoConfig,
	}
}

func resolveAgentAddress(agentName string) (string, error) {
	// Get master address from environment or default
	masterAddr := os.Getenv("SLOTH_RUNNER_MASTER_ADDR")
	if masterAddr == "" {
		masterAddr = "localhost:50053"
	}

	// Connect to master
	conn, err := grpc.Dial(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("failed to connect to master: %w", err)
	}
	defer conn.Close()

	client := pb.NewAgentRegistryClient(conn)

	// Get agent info
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetAgentInfo(ctx, &pb.GetAgentInfoRequest{
		AgentName: agentName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get agent info: %w", err)
	}

	if resp.AgentInfo == nil {
		return "", fmt.Errorf("agent '%s' not found", agentName)
	}

	return resp.AgentInfo.AgentAddress, nil
}
