package debug

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	_ "github.com/mattn/go-sqlite3"
)

func NewDebugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug and troubleshoot issues",
		Long:  `Debug workflows, agents, and connections to identify and resolve issues.`,
		Example: `  # Debug agent connectivity
  sloth-runner sysadmin debug connection web-01

  # Debug workflow execution
  sloth-runner sysadmin debug workflow latest

  # Get agent diagnostics
  sloth-runner sysadmin debug agent web-01`,
	}

	cmd.AddCommand(newConnectionCmd())
	cmd.AddCommand(newAgentCmd())
	cmd.AddCommand(newWorkflowCmd())

	return cmd
}

func newConnectionCmd() *cobra.Command {
	var timeout int
	var verbose bool

	cmd := &cobra.Command{
		Use:   "connection [agent-name]",
		Short: "Debug connection to an agent",
		Long:  `Test and debug connectivity to a specific agent, including ping, latency, and gRPC connection.`,
		Example: `  # Test connection to agent
  sloth-runner sysadmin debug connection web-01

  # Verbose output with timeout
  sloth-runner sysadmin debug connection web-01 --verbose --timeout 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugConnection(args[0], timeout, verbose)
		},
	}

	cmd.Flags().IntVarP(&timeout, "timeout", "t", 5, "Connection timeout in seconds")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func newAgentCmd() *cobra.Command {
	var full bool

	cmd := &cobra.Command{
		Use:   "agent [agent-name]",
		Short: "Debug agent configuration and status",
		Long:  `Get detailed diagnostics about an agent including configuration, status, and system info.`,
		Example: `  # Basic agent debug
  sloth-runner sysadmin debug agent web-01

  # Full diagnostics
  sloth-runner sysadmin debug agent web-01 --full`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugAgent(args[0], full)
		},
	}

	cmd.Flags().BoolVarP(&full, "full", "f", false, "Full diagnostics including system info")

	return cmd
}

func newWorkflowCmd() *cobra.Command {
	var last int

	cmd := &cobra.Command{
		Use:   "workflow [workflow-name|latest]",
		Short: "Debug workflow execution",
		Long:  `Analyze workflow execution, show task details, and identify bottlenecks.`,
		Example: `  # Debug latest workflow
  sloth-runner sysadmin debug workflow latest

  # Debug specific workflow
  sloth-runner sysadmin debug workflow deploy-prod

  # Debug last 5 workflows
  sloth-runner sysadmin debug workflow latest --last 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugWorkflow(args[0], last)
		},
	}

	cmd.Flags().IntVarP(&last, "last", "n", 1, "Number of workflow executions to show")

	return cmd
}

// Implementation functions

func debugConnection(agentName string, timeout int, verbose bool) error {
	fmt.Printf("🔍 Debugging connection to agent: %s\n\n", agentName)

	// Get agent from database
	dbPath := config.GetAgentDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var address, status string
	var lastHeartbeat int64
	err = db.QueryRow(`
		SELECT address, status, last_heartbeat
		FROM agents
		WHERE name = ?
	`, agentName).Scan(&address, &status, &lastHeartbeat)

	if err == sql.ErrNoRows {
		return fmt.Errorf("agent not found: %s", agentName)
	}
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	fmt.Printf("📋 Agent Info:\n")
	fmt.Printf("  Name:     %s\n", agentName)
	fmt.Printf("  Address:  %s\n", address)
	fmt.Printf("  Status:   %s\n", status)
	fmt.Printf("  Last HB:  %s\n", time.Unix(lastHeartbeat, 0).Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Parse address
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid address format: %s", address)
	}

	// Test 1: TCP Connection
	fmt.Printf("🔌 Test 1: TCP Connection to %s:%s\n", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("  ❌ FAILED: %v\n", err)
		fmt.Printf("  Duration: %v\n\n", duration)
		return fmt.Errorf("TCP connection failed")
	}
	conn.Close()
	fmt.Printf("  ✅ SUCCESS\n")
	fmt.Printf("  Duration: %v\n\n", duration)

	// Test 2: DNS Resolution
	if verbose {
		fmt.Printf("🌐 Test 2: DNS Resolution\n")
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Printf("  ❌ FAILED: %v\n\n", err)
		} else {
			fmt.Printf("  ✅ SUCCESS\n")
			for _, ip := range ips {
				fmt.Printf("  IP: %s\n", ip)
			}
			fmt.Println()
		}
	}

	// Test 3: gRPC Connection
	fmt.Printf("🔗 Test 3: gRPC Connection\n")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	start = time.Now()
	grpcConn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Duration(timeout)*time.Second))
	duration = time.Since(start)

	if err != nil {
		fmt.Printf("  ❌ FAILED: %v\n", err)
		fmt.Printf("  Duration: %v\n\n", duration)
		return fmt.Errorf("gRPC connection failed")
	}
	defer grpcConn.Close()
	fmt.Printf("  ✅ SUCCESS\n")
	fmt.Printf("  Duration: %v\n\n", duration)

	// Test 4: Agent RPC Call
	fmt.Printf("💓 Test 4: Agent RPC Call\n")
	client := pb.NewAgentClient(grpcConn)

	start = time.Now()
	_, err = client.GetResourceUsage(ctx, &pb.ResourceUsageRequest{})
	duration = time.Since(start)

	if err != nil {
		fmt.Printf("  ❌ FAILED: %v\n", err)
		fmt.Printf("  Duration: %v\n\n", duration)
	} else {
		fmt.Printf("  ✅ SUCCESS\n")
		fmt.Printf("  Duration: %v\n\n", duration)
	}

	// Test 5: Latency Test (3 pings)
	if verbose {
		fmt.Printf("⏱️  Test 5: Latency Test (3 pings)\n")
		var total time.Duration
		for i := 1; i <= 3; i++ {
			start = time.Now()
			_, err = client.GetResourceUsage(ctx, &pb.ResourceUsageRequest{})
			duration = time.Since(start)
			total += duration

			if err != nil {
				fmt.Printf("  Ping %d: ❌ FAILED (%v)\n", i, err)
			} else {
				fmt.Printf("  Ping %d: ✅ %v\n", i, duration)
			}
		}
		avg := total / 3
		fmt.Printf("  Average: %v\n\n", avg)
	}

	fmt.Println("✅ All tests completed successfully!")
	return nil
}

func debugAgent(agentName string, full bool) error {
	fmt.Printf("🔍 Agent Diagnostics: %s\n\n", agentName)

	// Get agent from database
	dbPath := config.GetAgentDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var address, status, version, systemInfo string
	var lastHeartbeat, registeredAt, updatedAt, lastInfoCollected int64

	err = db.QueryRow(`
		SELECT address, status, last_heartbeat, registered_at, updated_at,
		       last_info_collected, system_info, version
		FROM agents
		WHERE name = ?
	`, agentName).Scan(&address, &status, &lastHeartbeat, &registeredAt,
		&updatedAt, &lastInfoCollected, &systemInfo, &version)

	if err == sql.ErrNoRows {
		return fmt.Errorf("agent not found: %s", agentName)
	}
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	// Basic Info
	fmt.Printf("📋 Basic Information:\n")
	fmt.Printf("  Name:             %s\n", agentName)
	fmt.Printf("  Address:          %s\n", address)
	fmt.Printf("  Status:           %s\n", status)
	fmt.Printf("  Version:          %s\n", version)
	fmt.Printf("  Registered At:    %s\n", time.Unix(registeredAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  Last Updated:     %s\n", time.Unix(updatedAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  Last Heartbeat:   %s (%s ago)\n",
		time.Unix(lastHeartbeat, 0).Format("2006-01-02 15:04:05"),
		time.Since(time.Unix(lastHeartbeat, 0)).Round(time.Second))
	fmt.Println()

	// Connection Status
	now := time.Now().Unix()
	hbAge := now - lastHeartbeat

	fmt.Printf("🔌 Connection Status:\n")
	if hbAge < 60 {
		fmt.Printf("  ✅ HEALTHY - Recent heartbeat\n")
	} else if hbAge < 300 {
		fmt.Printf("  ⚠️  WARNING - Heartbeat is %d seconds old\n", hbAge)
	} else {
		fmt.Printf("  ❌ CRITICAL - No heartbeat for %d seconds\n", hbAge)
	}
	fmt.Println()

	// System Info
	if full && systemInfo != "" {
		fmt.Printf("💻 System Information:\n")

		var sysInfo map[string]interface{}
		if err := json.Unmarshal([]byte(systemInfo), &sysInfo); err == nil {
			// Pretty print system info
			prettyJSON, _ := json.MarshalIndent(sysInfo, "  ", "  ")
			fmt.Printf("%s\n\n", string(prettyJSON))
		} else {
			fmt.Printf("  Error parsing system info: %v\n\n", err)
		}
	}

	// Recommendations
	fmt.Printf("💡 Recommendations:\n")
	if hbAge > 300 {
		fmt.Printf("  • Check if agent service is running\n")
		fmt.Printf("  • Verify network connectivity\n")
		fmt.Printf("  • Check agent logs for errors\n")
	}
	if status == "Inactive" {
		fmt.Printf("  • Agent is marked as inactive\n")
		fmt.Printf("  • Consider restarting the agent\n")
	}
	if version == "" || version == "unknown" {
		fmt.Printf("  • Agent version is unknown\n")
		fmt.Printf("  • Consider updating the agent\n")
	}

	return nil
}

func debugWorkflow(workflowName string, last int) error {
	fmt.Printf("🔍 Workflow Debug: %s\n\n", workflowName)

	// Get execution history from database
	dbPath := config.GetHistoryDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer db.Close()

	query := `
		SELECT id, workflow_name, start_time, end_time, status, error_message,
		       group_name, agent_name, tasks_total, tasks_success, tasks_failed
		FROM executions
		WHERE 1=1
	`
	args := []interface{}{}

	if workflowName != "latest" {
		query += ` AND workflow_name = ?`
		args = append(args, workflowName)
	}

	query += ` ORDER BY start_time DESC LIMIT ?`
	args = append(args, last)

	rows, err := db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var tasksTotal, tasksSuccess, tasksFailed sql.NullInt64
		var name, status string
		var startTime, endTime int64
		var id, errorMsg, groupName, agentName sql.NullString

		err := rows.Scan(&id, &name, &startTime, &endTime, &status, &errorMsg,
			&groupName, &agentName, &tasksTotal, &tasksSuccess, &tasksFailed)
		if err != nil {
			fmt.Printf("Error reading row: %v\n", err)
			continue
		}

		count++
		fmt.Printf("📊 Execution #%d:\n", count)
		if id.Valid {
			fmt.Printf("  ID:           %s\n", id.String)
		}
		fmt.Printf("  Workflow:     %s\n", name)
		if groupName.Valid && groupName.String != "" {
			fmt.Printf("  Group:        %s\n", groupName.String)
		}
		if agentName.Valid && agentName.String != "" {
			fmt.Printf("  Agent:        %s\n", agentName.String)
		}
		fmt.Printf("  Status:       %s\n", status)
		fmt.Printf("  Start Time:   %s\n", time.Unix(startTime, 0).Format("2006-01-02 15:04:05"))

		if endTime > 0 {
			duration := time.Unix(endTime, 0).Sub(time.Unix(startTime, 0))
			fmt.Printf("  End Time:     %s\n", time.Unix(endTime, 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("  Duration:     %v\n", duration.Round(time.Millisecond))
		} else {
			fmt.Printf("  End Time:     (still running)\n")
		}

		if tasksTotal.Valid && tasksTotal.Int64 > 0 {
			fmt.Printf("  Tasks:        %d total, %d success, %d failed\n",
				tasksTotal.Int64, tasksSuccess.Int64, tasksFailed.Int64)
		}

		if errorMsg.Valid && errorMsg.String != "" {
			fmt.Printf("  Error:        %s\n", errorMsg.String)
		}
		fmt.Println()
	}

	if count == 0 {
		fmt.Printf("No workflow executions found")
		if workflowName != "latest" {
			fmt.Printf(" for: %s", workflowName)
		}
		fmt.Println()
	}

	return nil
}
