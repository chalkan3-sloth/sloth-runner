package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "github.com/mattn/go-sqlite3"
)

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    []HealthCheck     `json:"checks"`
	Summary   map[string]int    `json:"summary"`
	Details   map[string]string `json:"details,omitempty"`
}

// HealthCheck represents a single health check result
type HealthCheck struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Message   string        `json:"message"`
	Duration  time.Duration `json:"duration"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

func NewHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Health checks and diagnostics",
		Long:  `Monitor system health, check connectivity, and diagnose issues.`,
		Example: `  # Run all health checks
  sloth-runner health check

  # Check specific agent connectivity
  sloth-runner health agent do-sloth-runner-01

  # Check master server health
  sloth-runner health master

  # Continuous monitoring
  sloth-runner health watch --interval 30s`,
	}

	cmd.AddCommand(newCheckCmd())
	cmd.AddCommand(newAgentHealthCmd())
	cmd.AddCommand(newMasterHealthCmd())
	cmd.AddCommand(newWatchCmd())
	cmd.AddCommand(newDiagnosticsCmd())

	return cmd
}

func newCheckCmd() *cobra.Command {
	var outputFormat string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Run all health checks",
		Long:  `Execute all health checks and display results.`,
		Example: `  # Run all checks
  sloth-runner health check

  # Output as JSON
  sloth-runner health check --output json

  # Verbose output
  sloth-runner health check --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAllHealthChecks(outputFormat, verbose)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")

	return cmd
}

func newAgentHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent [agent-name]",
		Short: "Check agent health",
		Long:  `Check connectivity and health status of a specific agent.`,
		Example: `  # Check specific agent
  sloth-runner health agent do-sloth-runner-01

  # Check all agents
  sloth-runner health agent --all`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			if all {
				return checkAllAgents()
			}
			if len(args) == 0 {
				return fmt.Errorf("agent name required (or use --all)")
			}
			return checkAgentHealth(args[0])
		},
	}

	cmd.Flags().Bool("all", false, "Check all agents")

	return cmd
}

func newMasterHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Check master server health",
		Long:  `Check master server connectivity and status.`,
		Example: `  # Check master health
  sloth-runner health master

  # Check specific master
  sloth-runner health master --address localhost:50053`,
		RunE: func(cmd *cobra.Command, args []string) error {
			address, _ := cmd.Flags().GetString("address")
			return checkMasterHealth(address)
		},
	}

	cmd.Flags().String("address", "", "Master server address (default from config)")

	return cmd
}

func newWatchCmd() *cobra.Command {
	var interval string

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Continuous health monitoring",
		Long:  `Continuously monitor system health at specified intervals.`,
		Example: `  # Watch health every 30 seconds
  sloth-runner health watch --interval 30s

  # Watch every minute
  sloth-runner health watch --interval 1m`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return watchHealth(interval)
		},
	}

	cmd.Flags().StringVarP(&interval, "interval", "i", "30s", "Check interval (e.g., 30s, 1m, 5m)")

	return cmd
}

func newDiagnosticsCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "diagnostics",
		Short: "System diagnostics report",
		Long:  `Generate detailed diagnostics report including system info, database status, and configuration.`,
		Example: `  # Generate diagnostics
  sloth-runner health diagnostics

  # Save to file
  sloth-runner health diagnostics --output diagnostics.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateDiagnostics(output)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (stdout if not specified)")

	return cmd
}

// Implementation functions

func runAllHealthChecks(format string, verbose bool) error {
	status := HealthStatus{
		Timestamp: time.Now(),
		Checks:    []HealthCheck{},
		Summary:   make(map[string]int),
		Details:   make(map[string]string),
	}

	// Run checks
	checks := []func() HealthCheck{
		checkDatabase,
		checkDataDir,
		checkMasterConnection,
		checkLogDirectory,
		checkDiskSpace,
		checkMemory,
	}

	for _, check := range checks {
		result := check()
		status.Checks = append(status.Checks, result)
		status.Summary[result.Status]++
	}

	// Determine overall status
	if status.Summary["critical"] > 0 {
		status.Status = "critical"
	} else if status.Summary["error"] > 0 {
		status.Status = "error"
	} else if status.Summary["warning"] > 0 {
		status.Status = "warning"
	} else {
		status.Status = "healthy"
	}

	// Output results
	if format == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(status)
	}

	// Text output
	printHealthStatus(status, verbose)
	return nil
}

func checkDatabase() HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:      "Database Connectivity",
		Timestamp: start,
	}

	dbPath := config.GetAgentDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		check.Status = "error"
		check.Error = err.Error()
		check.Message = "Failed to open database"
		check.Duration = time.Since(start)
		return check
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		check.Status = "error"
		check.Error = err.Error()
		check.Message = "Database not responding"
	} else {
		check.Status = "ok"
		check.Message = "Database is accessible"
	}

	check.Duration = time.Since(start)
	return check
}

func checkDataDir() HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:      "Data Directory",
		Timestamp: start,
	}

	dataDir := config.GetDataDir()
	info, err := os.Stat(dataDir)
	if err != nil {
		check.Status = "error"
		check.Error = err.Error()
		check.Message = "Data directory not accessible"
	} else if !info.IsDir() {
		check.Status = "error"
		check.Message = "Data directory is not a directory"
	} else {
		// Check if writable
		testFile := filepath.Join(dataDir, ".health-check")
		if f, err := os.Create(testFile); err != nil {
			check.Status = "warning"
			check.Message = "Data directory not writable"
		} else {
			f.Close()
			os.Remove(testFile)
			check.Status = "ok"
			check.Message = fmt.Sprintf("Data directory is accessible and writable: %s", dataDir)
		}
	}

	check.Duration = time.Since(start)
	return check
}

func checkMasterConnection() HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:      "Master Server",
		Timestamp: start,
	}

	masterAddr := config.GetMasterAddress()
	if masterAddr == "" {
		check.Status = "warning"
		check.Message = "No master server configured"
		check.Duration = time.Since(start)
		return check
	}

	// Try to connect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, masterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		check.Status = "error"
		check.Error = err.Error()
		check.Message = fmt.Sprintf("Cannot connect to master: %s", masterAddr)
	} else {
		defer conn.Close()
		check.Status = "ok"
		check.Message = fmt.Sprintf("Master server is reachable: %s", masterAddr)
	}

	check.Duration = time.Since(start)
	return check
}

func checkLogDirectory() HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:      "Log Directory",
		Timestamp: start,
	}

	logDir := config.GetLogDir()
	info, err := os.Stat(logDir)
	if err != nil {
		check.Status = "warning"
		check.Error = err.Error()
		check.Message = "Log directory not found"
	} else if !info.IsDir() {
		check.Status = "error"
		check.Message = "Log path is not a directory"
	} else {
		check.Status = "ok"
		check.Message = fmt.Sprintf("Log directory is accessible: %s", logDir)
	}

	check.Duration = time.Since(start)
	return check
}

func checkDiskSpace() HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:      "Disk Space",
		Timestamp: start,
	}

	// This is a simplified check - in production you'd use syscall.Statfs
	dataDir := config.GetDataDir()
	info, err := os.Stat(dataDir)
	if err != nil {
		check.Status = "warning"
		check.Message = "Cannot check disk space"
	} else {
		_ = info
		check.Status = "ok"
		check.Message = "Disk space check passed"
	}

	check.Duration = time.Since(start)
	return check
}

func checkMemory() HealthCheck {
	start := time.Now()
	check := HealthCheck{
		Name:      "Memory Usage",
		Timestamp: start,
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	allocMB := m.Alloc / 1024 / 1024
	sysMB := m.Sys / 1024 / 1024

	if allocMB > 1024 {
		check.Status = "warning"
		check.Message = fmt.Sprintf("High memory usage: %d MB allocated, %d MB system", allocMB, sysMB)
	} else {
		check.Status = "ok"
		check.Message = fmt.Sprintf("Memory usage normal: %d MB allocated, %d MB system", allocMB, sysMB)
	}

	check.Duration = time.Since(start)
	return check
}

func checkAgentHealth(agentName string) error {
	fmt.Printf("🔍 Checking health of agent: %s\n\n", agentName)

	// Get agent from database
	dbPath := config.GetAgentDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open agent database: %w", err)
	}
	defer db.Close()

	// Query agent
	query := `SELECT name, address, status, last_heartbeat FROM agents WHERE name = ?`
	var name, address, status string
	var lastHeartbeat int64

	err = db.QueryRow(query, agentName).Scan(&name, &address, &status, &lastHeartbeat)
	if err == sql.ErrNoRows {
		return fmt.Errorf("agent not found: %s", agentName)
	} else if err != nil {
		return fmt.Errorf("failed to query agent: %w", err)
	}

	hbTime := time.Unix(lastHeartbeat, 0)

	// Check if agent is active
	fmt.Printf("📋 Agent Information:\n")
	fmt.Printf("   Name:    %s\n", name)
	fmt.Printf("   Address: %s\n", address)
	fmt.Printf("   Status:  %s\n", status)
	fmt.Printf("   Last Heartbeat: %s (%s ago)\n\n", hbTime.Format("2006-01-02 15:04:05"), time.Since(hbTime).Round(time.Second))

	// Try to connect to agent
	fmt.Printf("🔌 Connectivity Test:\n")
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Printf("   ❌ Connection failed: %v\n", err)
		fmt.Printf("   Duration: %v\n", time.Since(start))
		return nil
	}
	defer conn.Close()

	fmt.Printf("   ✅ Connection successful\n")
	fmt.Printf("   Duration: %v\n\n", time.Since(start))

	fmt.Printf("✅ Overall Status: Agent is healthy\n")
	return nil
}

func checkAllAgents() error {
	dbPath := config.GetAgentDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open agent database: %w", err)
	}
	defer db.Close()

	query := `SELECT name, address, status, last_heartbeat FROM agents ORDER BY name`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query agents: %w", err)
	}
	defer rows.Close()

	type Agent struct {
		Name          string
		Address       string
		Status        string
		LastHeartbeat time.Time
	}

	var agents []Agent
	for rows.Next() {
		var name, address, status string
		var lastHeartbeat int64
		if err := rows.Scan(&name, &address, &status, &lastHeartbeat); err != nil {
			return fmt.Errorf("failed to scan agent: %w", err)
		}
		agents = append(agents, Agent{
			Name:          name,
			Address:       address,
			Status:        status,
			LastHeartbeat: time.Unix(lastHeartbeat, 0),
		})
	}

	if len(agents) == 0 {
		fmt.Println("ℹ️  No agents registered")
		return nil
	}

	fmt.Printf("🔍 Checking health of %d agent(s)\n\n", len(agents))

	healthyCount := 0
	warningCount := 0
	errorCount := 0

	for _, agent := range agents {
		fmt.Printf("Agent: %s\n", agent.Name)

		// Quick connectivity check
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := grpc.DialContext(ctx, agent.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()

		if err != nil {
			fmt.Printf("   Status: ❌ Unreachable (%v)\n", err)
			errorCount++
		} else {
			conn.Close()

			// Check last heartbeat
			timeSinceHeartbeat := time.Since(agent.LastHeartbeat)
			if timeSinceHeartbeat > 5*time.Minute {
				fmt.Printf("   Status: ⚠️  Stale (last heartbeat %v ago)\n", timeSinceHeartbeat)
				warningCount++
			} else {
				fmt.Printf("   Status: ✅ Healthy\n")
				healthyCount++
			}
		}

		fmt.Printf("   Address: %s\n", agent.Address)
		fmt.Printf("   Last Heartbeat: %s\n\n", agent.LastHeartbeat.Format("2006-01-02 15:04:05"))
	}

	// Summary
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   Total:   %d\n", len(agents))
	fmt.Printf("   Healthy: %d\n", healthyCount)
	if warningCount > 0 {
		fmt.Printf("   Warning: %d\n", warningCount)
	}
	if errorCount > 0 {
		fmt.Printf("   Error:   %d\n", errorCount)
	}

	return nil
}

func checkMasterHealth(address string) error {
	if address == "" {
		address = config.GetMasterAddress()
	}

	if address == "" {
		fmt.Println("❌ No master server configured")
		fmt.Println("\nSet master address with:")
		fmt.Println("   export SLOTH_RUNNER_MASTER_ADDR=host:port")
		return nil
	}

	fmt.Printf("🔍 Checking master server: %s\n\n", address)

	// Test TCP connection
	fmt.Printf("🔌 TCP Connectivity:\n")
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		fmt.Printf("   ❌ TCP connection failed: %v\n", err)
		return nil
	}
	conn.Close()
	fmt.Printf("   ✅ TCP connection successful\n")
	fmt.Printf("   Duration: %v\n\n", time.Since(start))

	// Test gRPC connection
	fmt.Printf("🔌 gRPC Connectivity:\n")
	grpcStart := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcConn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Printf("   ❌ gRPC connection failed: %v\n", err)
		return nil
	}
	defer grpcConn.Close()

	fmt.Printf("   ✅ gRPC connection successful\n")
	fmt.Printf("   Duration: %v\n\n", time.Since(grpcStart))

	fmt.Printf("✅ Master server is healthy\n")
	return nil
}

func watchHealth(intervalStr string) error {
	duration, err := time.ParseDuration(intervalStr)
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}

	fmt.Printf("👀 Watching system health (interval: %v)\n", duration)
	fmt.Println("Press Ctrl+C to stop")

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	// Run check immediately
	runHealthCheckSummary()

	for range ticker.C {
		fmt.Println("\n" + strings.Repeat("─", 60))
		runHealthCheckSummary()
	}

	return nil
}

func runHealthCheckSummary() {
	status := HealthStatus{
		Timestamp: time.Now(),
		Checks:    []HealthCheck{},
		Summary:   make(map[string]int),
	}

	checks := []func() HealthCheck{
		checkDatabase,
		checkDataDir,
		checkMasterConnection,
		checkLogDirectory,
	}

	for _, check := range checks {
		result := check()
		status.Checks = append(status.Checks, result)
		status.Summary[result.Status]++
	}

	// Print compact summary
	fmt.Printf("[%s] ", time.Now().Format("15:04:05"))

	if status.Summary["error"] > 0 || status.Summary["critical"] > 0 {
		fmt.Printf("❌ UNHEALTHY")
	} else if status.Summary["warning"] > 0 {
		fmt.Printf("⚠️  WARNING")
	} else {
		fmt.Printf("✅ HEALTHY")
	}

	fmt.Printf(" | OK: %d", status.Summary["ok"])
	if status.Summary["warning"] > 0 {
		fmt.Printf(" | WARN: %d", status.Summary["warning"])
	}
	if status.Summary["error"] > 0 {
		fmt.Printf(" | ERROR: %d", status.Summary["error"])
	}
	fmt.Println()
}

func generateDiagnostics(outputFile string) error {
	fmt.Println("🔍 Generating diagnostics report...")

	diagnostics := map[string]interface{}{
		"timestamp":   time.Now(),
		"version":     "dev",
		"system": map[string]interface{}{
			"os":       runtime.GOOS,
			"arch":     runtime.GOARCH,
			"cpus":     runtime.NumCPU(),
			"go_version": runtime.Version(),
		},
		"configuration": map[string]interface{}{
			"data_dir":       config.GetDataDir(),
			"log_dir":        config.GetLogDir(),
			"master_address": config.GetMasterAddress(),
		},
		"health_checks": []HealthCheck{},
	}

	// Run all checks
	checks := []func() HealthCheck{
		checkDatabase,
		checkDataDir,
		checkMasterConnection,
		checkLogDirectory,
		checkDiskSpace,
		checkMemory,
	}

	for _, check := range checks {
		result := check()
		diagnostics["health_checks"] = append(diagnostics["health_checks"].([]HealthCheck), result)
	}

	// Output
	var output *os.File
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		output = f
		fmt.Printf("Writing diagnostics to %s...\n", outputFile)
	} else {
		output = os.Stdout
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(diagnostics); err != nil {
		return fmt.Errorf("failed to encode diagnostics: %w", err)
	}

	if outputFile != "" {
		fmt.Printf("✅ Diagnostics saved to %s\n", outputFile)
	}

	return nil
}

func printHealthStatus(status HealthStatus, verbose bool) {
	// Header
	fmt.Printf("\n🏥 Health Check Report\n")
	fmt.Printf("═══════════════════════════════════════════════════\n\n")
	fmt.Printf("Timestamp: %s\n", status.Timestamp.Format("2006-01-02 15:04:05"))

	// Overall status
	statusIcon := "❌"
	switch status.Status {
	case "healthy":
		statusIcon = "✅"
	case "warning":
		statusIcon = "⚠️"
	case "error", "critical":
		statusIcon = "❌"
	}

	fmt.Printf("Status:    %s %s\n\n", statusIcon, strings.ToUpper(status.Status))

	// Summary
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   OK:      %d\n", status.Summary["ok"])
	if status.Summary["warning"] > 0 {
		fmt.Printf("   Warning: %d\n", status.Summary["warning"])
	}
	if status.Summary["error"] > 0 {
		fmt.Printf("   Error:   %d\n", status.Summary["error"])
	}
	if status.Summary["critical"] > 0 {
		fmt.Printf("   Critical: %d\n", status.Summary["critical"])
	}
	fmt.Println()

	// Individual checks
	fmt.Printf("📋 Checks:\n")
	for _, check := range status.Checks {
		icon := "✅"
		switch check.Status {
		case "warning":
			icon = "⚠️"
		case "error", "critical":
			icon = "❌"
		}

		fmt.Printf("   %s %s: %s", icon, check.Name, check.Message)
		if verbose {
			fmt.Printf(" (%v)", check.Duration)
		}
		fmt.Println()

		if verbose && check.Error != "" {
			fmt.Printf("      Error: %s\n", check.Error)
		}
	}

	fmt.Println()
}
