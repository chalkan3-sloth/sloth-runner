package agent

import (
	"context"
	"fmt"
	"io"
	"strings"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
)

// moduleCheck represents a tool or module to check
type moduleCheck struct {
	Name        string
	Command     string
	Description string
}

// ModulesCheckOptions contains options for checking agent modules
type ModulesCheckOptions struct {
	AgentName string
	Writer    io.Writer
}

// ModulesCheckResult contains the results of module checking
type ModulesCheckResult struct {
	Available []moduleCheck
	Missing   []moduleCheck
}

// getDefaultModules returns the standard list of modules to check
func getDefaultModules() []moduleCheck {
	return []moduleCheck{
		{"Incus", "incus", "Container and VM management (LXC/LXD successor)"},
		{"Terraform", "terraform", "Infrastructure as Code provisioning"},
		{"Pulumi", "pulumi", "Modern Infrastructure as Code with programming languages"},
		{"AWS CLI", "aws", "Amazon Web Services command-line interface"},
		{"Azure CLI", "az", "Microsoft Azure command-line interface"},
		{"Google Cloud SDK", "gcloud", "Google Cloud Platform command-line interface"},
		{"kubectl", "kubectl", "Kubernetes command-line tool"},
		{"Docker", "docker", "Container platform"},
		{"Ansible", "ansible", "IT automation and configuration management"},
		{"Git", "git", "Version control system"},
		{"Helm", "helm", "Kubernetes package manager"},
		{"systemctl", "systemctl", "systemd service manager"},
		{"curl", "curl", "HTTP client for data transfer"},
		{"jq", "jq", "JSON processor"},
	}
}

// checkAgentModulesWithClient checks module availability using an injected client (testable)
func checkAgentModulesWithClient(ctx context.Context, client AgentRegistryClient, opts ModulesCheckOptions) error {
	pterm.DefaultHeader.WithFullWidth().Printf("Module Availability Check - Agent: %s", opts.AgentName)
	fmt.Fprintln(opts.Writer)

	modules := getDefaultModules()
	result := &ModulesCheckResult{
		Available: []moduleCheck{},
		Missing:   []moduleCheck{},
	}

	// Check each module
	spinner, _ := pterm.DefaultSpinner.Start("Checking modules on agent...")

	for _, mod := range modules {
		isAvailable, err := checkModuleAvailability(ctx, client, opts.AgentName, mod)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to check module %s: %v", mod.Name, err))
			return err
		}

		if isAvailable {
			result.Available = append(result.Available, mod)
		} else {
			result.Missing = append(result.Missing, mod)
		}
	}

	spinner.Success("Module check completed")
	fmt.Fprintln(opts.Writer)

	// Format and display results
	return formatModulesResults(result, modules, opts.Writer)
}

// checkModuleAvailability checks if a single module is available on the agent (testable)
func checkModuleAvailability(ctx context.Context, client AgentRegistryClient, agentName string, mod moduleCheck) (bool, error) {
	checkCmd := fmt.Sprintf("command -v %s >/dev/null 2>&1 && echo 'found' || echo 'not found'", mod.Command)

	stream, err := client.ExecuteCommand(ctx, &pb.ExecuteCommandRequest{
		AgentName: agentName,
		Command:   checkCmd,
	})

	if err != nil {
		return false, fmt.Errorf("failed to execute command: %w", err)
	}

	var output strings.Builder
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, fmt.Errorf("error receiving stream: %w", err)
		}
		output.WriteString(resp.GetStdoutChunk())
		output.WriteString(resp.GetStderrChunk())
	}

	result := strings.TrimSpace(output.String())
	return result == "found", nil
}

// formatModulesResults formats and displays the module check results (testable)
func formatModulesResults(result *ModulesCheckResult, allModules []moduleCheck, w io.Writer) error {
	// Display available modules
	if len(result.Available) > 0 {
		pterm.Success.Println("✅ Available Modules:")
		fmt.Fprintln(w)
		for _, mod := range result.Available {
			fmt.Fprintf(w, "  %s %s\n", pterm.Green("✓"), pterm.Cyan(mod.Name))
			fmt.Fprintf(w, "    %s\n", pterm.Gray(mod.Description))
		}
		fmt.Fprintln(w)
	}

	// Display missing modules
	if len(result.Missing) > 0 {
		pterm.Warning.Println("❌ Missing Modules:")
		fmt.Fprintln(w)
		for _, mod := range result.Missing {
			fmt.Fprintf(w, "  %s %s\n", pterm.Red("✗"), pterm.Cyan(mod.Name))
			fmt.Fprintf(w, "    %s\n", pterm.Gray(mod.Description))
		}
		fmt.Fprintln(w)

		pterm.Info.Println("ℹ️  Information:")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  Missing modules are optional but required if you want to use their")
		fmt.Fprintln(w, "  corresponding Lua functions in your tasks.")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  For example:")
		fmt.Fprintln(w, "    - To use incus.instance() functions, install Incus")
		fmt.Fprintln(w, "    - To use terraform.init() functions, install Terraform")
		fmt.Fprintln(w, "    - To use aws.s3() functions, install AWS CLI")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  Install the tools you need based on your infrastructure requirements.")
	}

	// Summary
	fmt.Fprintln(w)
	pterm.DefaultBox.WithTitle("Summary").WithTitleTopCenter().Println(
		fmt.Sprintf("Available: %s  |  Missing: %s  |  Total: %s",
			pterm.Green(fmt.Sprintf("%d", len(result.Available))),
			pterm.Red(fmt.Sprintf("%d", len(result.Missing))),
			pterm.Cyan(fmt.Sprintf("%d", len(allModules))),
		),
	)

	return nil
}
