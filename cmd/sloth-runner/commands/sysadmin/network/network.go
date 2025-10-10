package network

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Short:   "Network diagnostics and monitoring",
		Long:    `Diagnose network connectivity, measure latency, check ports, and monitor network performance.`,
		Aliases: []string{"net"},
	}

	// ping command
	pingCmd := &cobra.Command{
		Use:   "ping [host]",
		Short: "Ping host to test connectivity",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			count, _ := cmd.Flags().GetInt("count")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			useTCP, _ := cmd.Flags().GetBool("tcp")
			tcpPort, _ := cmd.Flags().GetInt("tcp-port")

			if err := runPing(host, count, timeout, useTCP, tcpPort); err != nil {
				pterm.Error.Printf("Failed to ping %s: %v\n", host, err)
			}
		},
	}
	pingCmd.Flags().IntP("count", "c", 4, "Number of pings to send")
	pingCmd.Flags().DurationP("timeout", "t", 10*time.Second, "Timeout for ping operation")
	pingCmd.Flags().Bool("tcp", false, "Use TCP ping instead of ICMP")
	pingCmd.Flags().Int("tcp-port", 80, "Port to use for TCP ping")
	cmd.AddCommand(pingCmd)

	// port-check command
	portCheckCmd := &cobra.Command{
		Use:   "port-check [host]",
		Short: "Check if port(s) are open on a host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			port, _ := cmd.Flags().GetInt("port")
			ports, _ := cmd.Flags().GetIntSlice("ports")
			timeout, _ := cmd.Flags().GetDuration("timeout")

			// If --ports is specified, use that; otherwise use --port
			if len(ports) == 0 && port > 0 {
				ports = []int{port}
			}

			if len(ports) == 0 {
				pterm.Error.Println("Please specify at least one port with --port or --ports")
				return
			}

			if err := runPortCheck(host, ports, timeout); err != nil {
				pterm.Error.Printf("Failed to check ports on %s: %v\n", host, err)
			}
		},
	}
	portCheckCmd.Flags().IntP("port", "p", 0, "Single port to check")
	portCheckCmd.Flags().IntSliceP("ports", "P", []int{}, "Multiple ports to check (comma-separated)")
	portCheckCmd.Flags().DurationP("timeout", "t", 5*time.Second, "Timeout for each port check")
	cmd.AddCommand(portCheckCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the network command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showNetworkDocs()
		},
	})

	return cmd
}

func showNetworkDocs() {
	title := "SLOTH-RUNNER SYSADMIN NETWORK(1)"
	description := "sloth-runner sysadmin network - Network diagnostics and monitoring"
	synopsis := "sloth-runner sysadmin network [subcommand] [options]"

	options := [][]string{
		{"ping", "Test connectivity with agents. Measures latency and packet loss."},
		{"port-check", "Check if a port is open on remote agent. Includes service detection."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Test connectivity",
			"sloth-runner sysadmin network ping --agent web-01",
			"Tests network connectivity and measures round-trip time",
		},
		{
			"Ping all agents",
			"sloth-runner sysadmin net ping --all-agents",
			"Tests connectivity to all registered agents",
		},
		{
			"Check single port",
			"sloth-runner sysadmin network port-check --agent web-01 --port 80",
			"Checks if port 80 is accessible on web-01",
		},
		{
			"Check multiple ports",
			"sloth-runner sysadmin net port-check --agent db-01 --ports 5432,3306,6379",
			"Checks PostgreSQL, MySQL, and Redis ports",
		},
		{
			"Network diagnostics",
			"sloth-runner sysadmin network ping --agent web-01 --count 10 --interval 1s",
			"Sends 10 pings with 1 second interval",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin performance - Performance monitoring",
		"sloth-runner sysadmin health - Health checks",
	}

	showDocs(title, description, synopsis, options, examples, seeAlso)
}

// showDocs displays formatted documentation similar to man pages
func showDocs(title, description, synopsis string, options [][]string, examples [][]string, seeAlso []string) {
	// Header
	pterm.DefaultHeader.WithFullWidth().Println(title)
	fmt.Println()

	// Name and Description
	pterm.DefaultSection.Println("NAME")
	fmt.Printf("    %s\n\n", description)

	// Synopsis
	if synopsis != "" {
		pterm.DefaultSection.Println("SYNOPSIS")
		fmt.Printf("    %s\n\n", synopsis)
	}

	// Options
	if len(options) > 0 {
		pterm.DefaultSection.Println("OPTIONS")
		for _, opt := range options {
			if len(opt) >= 2 {
				pterm.FgCyan.Printf("    %s\n", opt[0])
				fmt.Printf("        %s\n\n", opt[1])
			}
		}
	}

	// Examples
	if len(examples) > 0 {
		pterm.DefaultSection.Println("EXAMPLES")
		for i, ex := range examples {
			if len(ex) >= 2 {
				pterm.FgYellow.Printf("    Example %d: %s\n", i+1, ex[0])
				pterm.FgGreen.Printf("    $ %s\n", ex[1])
				if len(ex) >= 3 {
					fmt.Printf("        %s\n", ex[2])
				}
				fmt.Println()
			}
		}
	}

	// See Also
	if len(seeAlso) > 0 {
		pterm.DefaultSection.Println("SEE ALSO")
		for _, item := range seeAlso {
			fmt.Printf("    • %s\n", item)
		}
		fmt.Println()
	}

	// Footer
	pterm.FgGray.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	pterm.FgGray.Println("Documentation generated for sloth-runner sysadmin v2.0")
	pterm.FgGray.Println("For more information: sloth-runner sysadmin network --help")
}

// runPing executa o comando ping
func runPing(host string, count int, timeout time.Duration, useTCP bool, tcpPort int) error {
	diag := NewDiagnostics()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Pinging %s...", host))

	var result *PingResult
	var err error

	if useTCP {
		result, err = TCPPing(host, tcpPort, count, timeout)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to TCP ping %s", host))
			return err
		}
		spinner.Success(fmt.Sprintf("✅ TCP ping completed to %s:%d", host, tcpPort))
	} else {
		result, err = diag.Ping(host, count, timeout)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to ping %s", host))
			return err
		}
		spinner.Success(fmt.Sprintf("✅ Ping completed to %s", host))
	}

	// Display results
	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println(fmt.Sprintf("Ping Results: %s", host))
	pterm.Println()

	// Statistics table
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Packets Sent", fmt.Sprintf("%d", result.PacketsSent)},
		{"Packets Received", fmt.Sprintf("%d", result.PacketsRecv)},
		{"Packet Loss", fmt.Sprintf("%.1f%%", result.PacketLoss)},
	}

	if result.PacketsRecv > 0 {
		tableData = append(tableData,
			[]string{"Min RTT", result.MinRTT.String()},
			[]string{"Avg RTT", result.AvgRTT.String()},
			[]string{"Max RTT", result.MaxRTT.String()},
		)
		if result.StdDevRTT > 0 {
			tableData = append(tableData, []string{"StdDev RTT", result.StdDevRTT.String()})
		}
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Status indicator
	if result.PacketLoss == 0 {
		pterm.Success.Printf("Excellent connectivity: 0%% packet loss\n")
	} else if result.PacketLoss < 10 {
		pterm.Warning.Printf("Minor packet loss: %.1f%%\n", result.PacketLoss)
	} else if result.PacketLoss < 50 {
		pterm.Warning.Printf("Significant packet loss: %.1f%%\n", result.PacketLoss)
	} else {
		pterm.Error.Printf("Critical packet loss: %.1f%%\n", result.PacketLoss)
	}

	return nil
}

// runPortCheck executa verificação de portas
func runPortCheck(host string, ports []int, timeout time.Duration) error {
	diag := NewDiagnostics()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking %d port(s) on %s...", len(ports), host))

	results, err := diag.CheckPorts(host, ports, timeout)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to check ports on %s", host))
		return err
	}

	// Count open ports
	openCount := 0
	for _, r := range results {
		if r.Open {
			openCount++
		}
	}

	spinner.Success(fmt.Sprintf("✅ Port check completed: %d/%d ports open", openCount, len(ports)))
	pterm.Println()

	// Display results
	pterm.DefaultHeader.WithFullWidth().Println(fmt.Sprintf("Port Scan Results: %s", host))
	pterm.Println()

	// Results table
	tableData := pterm.TableData{
		{"Port", "Status", "Service", "Latency"},
	}

	for _, r := range results {
		status := "❌ Closed"
		statusColor := pterm.FgRed

		if r.Open {
			status = "✅ Open"
			statusColor = pterm.FgGreen
		}

		tableData = append(tableData, []string{
			fmt.Sprintf("%d", r.Port),
			statusColor.Sprint(status),
			r.Service,
			r.Latency.String(),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	// Summary
	if openCount == len(ports) {
		pterm.Success.Printf("All %d port(s) are open\n", len(ports))
	} else if openCount == 0 {
		pterm.Error.Printf("All %d port(s) are closed or filtered\n", len(ports))
	} else {
		pterm.Info.Printf("%d of %d port(s) are open\n", openCount, len(ports))
	}

	return nil
}
