package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	agentInternal "github.com/chalkan3-sloth/sloth-runner/internal/agent"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// GetAgentInfoOptions contains options for getting agent info
type GetAgentInfoOptions struct {
	AgentName    string
	OutputFormat string
	Writer       io.Writer
}

// getAgentInfoWithClient retrieves agent info using an injected client (testable)
func getAgentInfoWithClient(ctx context.Context, client AgentRegistryClient, opts GetAgentInfoOptions) error {
	resp, err := client.GetAgentInfo(ctx, &pb.GetAgentInfoRequest{
		AgentName: opts.AgentName,
	})
	if err != nil {
		return fmt.Errorf("failed to get agent info: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("%s\n\n%s", pterm.Red("✗ Failed to Get Agent Info"), resp.Message)
	}

	agent := resp.GetAgentInfo()

	if opts.OutputFormat == "json" {
		return formatAgentInfoJSON(agent, opts.Writer)
	}

	return formatAgentInfoText(agent, opts.Writer)
}

// formatAgentInfoJSON formats agent info as JSON (testable)
func formatAgentInfoJSON(agent *pb.AgentInfo, w io.Writer) error {
	output := map[string]interface{}{
		"agent_name":          agent.GetAgentName(),
		"agent_address":       agent.GetAgentAddress(),
		"status":              agent.GetStatus(),
		"last_heartbeat":      agent.GetLastHeartbeat(),
		"last_info_collected": agent.GetLastInfoCollected(),
	}

	// Parse and include system info if available
	if agent.GetSystemInfoJson() != "" {
		var sysInfo map[string]interface{}
		if err := json.Unmarshal([]byte(agent.GetSystemInfoJson()), &sysInfo); err == nil {
			output["system_info"] = sysInfo
		} else {
			output["system_info"] = agent.GetSystemInfoJson()
		}
	} else {
		output["system_info"] = nil
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	fmt.Fprintln(w, string(jsonOutput))
	return nil
}

// formatAgentInfoText formats agent info as human-readable text (testable)
func formatAgentInfoText(agent *pb.AgentInfo, w io.Writer) error {
	// Header
	pterm.DefaultHeader.WithFullWidth().WithWriter(w).Println(fmt.Sprintf("Agent Information: %s", agent.GetAgentName()))
	fmt.Fprintln(w)

	// Basic info
	pterm.Info.WithWriter(w).Println("Basic Information:")
	fmt.Fprintf(w, "  Name:         %s\n", pterm.Cyan(agent.GetAgentName()))
	fmt.Fprintf(w, "  Address:      %s\n", pterm.Cyan(agent.GetAgentAddress()))

	status := agent.GetStatus()
	if status == "Active" {
		fmt.Fprintf(w, "  Status:       %s\n", pterm.Green(status))
	} else {
		fmt.Fprintf(w, "  Status:       %s\n", pterm.Red(status))
	}

	if agent.GetLastHeartbeat() > 0 {
		fmt.Fprintf(w, "  Last Heartbeat: %s\n", pterm.Yellow(time.Unix(agent.GetLastHeartbeat(), 0).Format(time.RFC3339)))
	} else {
		fmt.Fprintf(w, "  Last Heartbeat: %s\n", pterm.Gray("Never"))
	}

	if agent.GetLastInfoCollected() > 0 {
		fmt.Fprintf(w, "  Last Info:     %s\n", pterm.Yellow(time.Unix(agent.GetLastInfoCollected(), 0).Format(time.RFC3339)))
	} else {
		fmt.Fprintf(w, "  Last Info:     %s\n", pterm.Gray("Not collected"))
	}

	fmt.Fprintln(w)

	// System info
	if agent.GetSystemInfoJson() != "" {
		sysInfo, err := agentInternal.FromJSON(agent.GetSystemInfoJson())
		if err != nil {
			pterm.Warning.WithWriter(w).Printf("Failed to parse system info: %v\n", err)
		} else {
			formatSystemInfo(sysInfo, w)
		}
	} else {
		pterm.Warning.WithWriter(w).Println("No system information available")
		fmt.Fprintln(w, "  System info has not been collected yet.")
		fmt.Fprintln(w, "  Wait for the agent to send a heartbeat with system info.")
	}

	return nil
}

// formatSystemInfo formats system information (testable)
func formatSystemInfo(sysInfo *agentInternal.SystemInfo, w io.Writer) {
	pterm.Info.WithWriter(w).Println("System Information:")
	fmt.Fprintf(w, "  Hostname:      %s\n", pterm.Cyan(sysInfo.Hostname))
	fmt.Fprintf(w, "  Platform:      %s %s\n", pterm.Cyan(sysInfo.Platform), pterm.Gray(sysInfo.PlatformVersion))
	fmt.Fprintf(w, "  Architecture:  %s\n", pterm.Cyan(sysInfo.Architecture))
	fmt.Fprintf(w, "  CPUs:          %s\n", pterm.Cyan(fmt.Sprintf("%d", sysInfo.CPUs)))
	fmt.Fprintf(w, "  Kernel:        %s %s\n", pterm.Cyan(sysInfo.Kernel), pterm.Gray(sysInfo.KernelVersion))

	if sysInfo.Virtualization != "none" {
		fmt.Fprintf(w, "  Virtualization: %s\n", pterm.Magenta(sysInfo.Virtualization))
	}

	fmt.Fprintf(w, "  Uptime:        %s\n", pterm.Yellow(fmt.Sprintf("%d seconds", sysInfo.Uptime)))

	if len(sysInfo.LoadAverage) == 3 {
		fmt.Fprintf(w, "  Load Average:  %s\n", pterm.Cyan(fmt.Sprintf("%.2f, %.2f, %.2f",
			sysInfo.LoadAverage[0], sysInfo.LoadAverage[1], sysInfo.LoadAverage[2])))
	}

	// Memory info
	if sysInfo.Memory != nil {
		formatMemoryInfo(sysInfo.Memory, w)
	}

	// Disk info
	if len(sysInfo.Disk) > 0 {
		formatDiskInfoList(sysInfo.Disk, w)
	}

	// Network info
	if len(sysInfo.Network) > 0 {
		formatNetworkInfoList(sysInfo.Network, w)
	}

	// Package info
	if sysInfo.Packages != nil && sysInfo.Packages.Manager != "" {
		formatPackageInfo(sysInfo.Packages, w)
	}

	// Services
	if len(sysInfo.Services) > 0 {
		formatServicesInfoList(sysInfo.Services, w)
	}
}

// formatMemoryInfo formats memory information (testable)
func formatMemoryInfo(memory *agentInternal.MemoryInfo, w io.Writer) {
	fmt.Fprintln(w)
	pterm.Info.WithWriter(w).Println("Memory Information:")
	fmt.Fprintf(w, "  Total:        %s\n", pterm.Cyan(formatBytes(memory.Total)))
	fmt.Fprintf(w, "  Used:         %s (%.1f%%)\n",
		pterm.Yellow(formatBytes(memory.Used)), memory.UsedPercent)
	fmt.Fprintf(w, "  Available:    %s\n", pterm.Green(formatBytes(memory.Available)))
	fmt.Fprintf(w, "  Free:         %s\n", pterm.Cyan(formatBytes(memory.Free)))
	if memory.Cached > 0 {
		fmt.Fprintf(w, "  Cached:       %s\n", pterm.Cyan(formatBytes(memory.Cached)))
	}
}

// formatDiskInfoList formats disk information (testable)
func formatDiskInfoList(disks []*agentInternal.DiskInfo, w io.Writer) {
	fmt.Fprintln(w)
	pterm.Info.WithWriter(w).Println("Disk Information:")
	for _, disk := range disks {
		if disk.Total > 0 {
			fmt.Fprintf(w, "  %s (%s):\n", pterm.Cyan(disk.Mountpoint), pterm.Gray(disk.Device))
			fmt.Fprintf(w, "    Total:  %s\n", pterm.Cyan(formatBytes(disk.Total)))
			fmt.Fprintf(w, "    Used:   %s (%.1f%%)\n", pterm.Yellow(formatBytes(disk.Used)), disk.UsedPercent)
			fmt.Fprintf(w, "    Free:   %s\n", pterm.Green(formatBytes(disk.Free)))
		}
	}
}

// formatNetworkInfoList formats network information (testable)
func formatNetworkInfoList(network []*agentInternal.NetworkInfo, w io.Writer) {
	fmt.Fprintln(w)
	pterm.Info.WithWriter(w).Println("Network Interfaces:")
	for _, iface := range network {
		if iface.Name != "lo" && len(iface.Addresses) > 0 {
			status := pterm.Red("DOWN")
			if iface.IsUp {
				status = pterm.Green("UP")
			}
			fmt.Fprintf(w, "  %s [%s]:\n", pterm.Cyan(iface.Name), status)
			if iface.MAC != "" {
				fmt.Fprintf(w, "    MAC:        %s\n", pterm.Gray(iface.MAC))
			}
			for _, addr := range iface.Addresses {
				fmt.Fprintf(w, "    Address:    %s\n", pterm.Yellow(addr))
			}
		}
	}
}

// formatPackageInfo formats package information (testable)
func formatPackageInfo(packages *agentInternal.PackageInfo, w io.Writer) {
	fmt.Fprintln(w)
	pterm.Info.WithWriter(w).Println("Package Information:")
	fmt.Fprintf(w, "  Manager:      %s\n", pterm.Cyan(packages.Manager))
	fmt.Fprintf(w, "  Installed:    %s\n", pterm.Cyan(fmt.Sprintf("%d packages", packages.InstalledCount)))
	if packages.UpdatesAvailable > 0 {
		fmt.Fprintf(w, "  Updates:      %s\n", pterm.Yellow(fmt.Sprintf("%d available", packages.UpdatesAvailable)))
	} else {
		fmt.Fprintf(w, "  Updates:      %s\n", pterm.Green("System is up to date"))
	}
}

// formatServicesInfoList formats services information (testable)
func formatServicesInfoList(services []agentInternal.ServiceInfo, w io.Writer) {
	fmt.Fprintln(w)
	pterm.Info.WithWriter(w).Printf("Running Services: %d total\n", len(services))
	fmt.Fprintln(w, "  (showing first 10)")
	count := len(services)
	if count > 10 {
		count = 10
	}
	for i := 0; i < count; i++ {
		fmt.Fprintf(w, "  - %s (%s)\n", pterm.Cyan(services[i].Name), pterm.Gray(services[i].Status))
	}
	if len(services) > 10 {
		fmt.Fprintf(w, "  ... and %d more\n", len(services)-10)
	}
}
