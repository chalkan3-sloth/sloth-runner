package deployment

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewDeploymentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployment",
		Short:   "Deployment and rollback management",
		Long:    `Manage deployments, rollbacks, and version control across agents.`,
		Aliases: []string{"deploy"},
		Example: `  # Deploy new version
  sloth-runner sysadmin deployment deploy --version v6.25.0 --agents all

  # Rollback to previous version
  sloth-runner sysadmin deployment rollback --agent do-sloth-runner-01

  # Show deployment history
  sloth-runner sysadmin deployment history --agent do-sloth-runner-01`,
	}

	// deploy command
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy new version to agents",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := cmd.Flags().GetString("version")
			agents, _ := cmd.Flags().GetStringSlice("agents")
			strategy, _ := cmd.Flags().GetString("strategy")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			healthCheck, _ := cmd.Flags().GetBool("health-check")

			if version == "" {
				pterm.Error.Println("Version is required (use --version)")
				return
			}

			if len(agents) == 0 {
				pterm.Error.Println("At least one agent is required (use --agents)")
				return
			}

			if err := runDeploy(version, agents, strategy, dryRun, healthCheck); err != nil {
				pterm.Error.Printf("Deployment failed: %v\n", err)
			}
		},
	}
	deployCmd.Flags().StringP("version", "v", "", "Version to deploy")
	deployCmd.Flags().StringSliceP("agents", "a", []string{}, "Agents to deploy to (comma-separated)")
	deployCmd.Flags().StringP("strategy", "s", "direct", "Deployment strategy (direct, rolling, canary, blue-green)")
	deployCmd.Flags().Bool("dry-run", false, "Simulate deployment without making changes")
	deployCmd.Flags().Bool("health-check", true, "Perform health check after deployment")
	cmd.AddCommand(deployCmd)

	// rollback command
	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback to previous version",
		Run: func(cmd *cobra.Command, args []string) {
			agent, _ := cmd.Flags().GetString("agent")
			version, _ := cmd.Flags().GetString("version")
			verify, _ := cmd.Flags().GetBool("verify")

			if agent == "" {
				pterm.Error.Println("Agent is required (use --agent)")
				return
			}

			if err := runRollback(agent, version, verify); err != nil {
				pterm.Error.Printf("Rollback failed: %v\n", err)
			}
		},
	}
	rollbackCmd.Flags().StringP("agent", "a", "", "Agent to rollback")
	rollbackCmd.Flags().StringP("version", "v", "", "Version to rollback to (default: previous version)")
	rollbackCmd.Flags().Bool("verify", true, "Verify agent health after rollback")
	cmd.AddCommand(rollbackCmd)

	// history command
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show deployment history",
		Run: func(cmd *cobra.Command, args []string) {
			agent, _ := cmd.Flags().GetString("agent")
			limit, _ := cmd.Flags().GetInt("limit")

			if err := runHistory(agent, limit); err != nil {
				pterm.Error.Printf("Failed to get history: %v\n", err)
			}
		},
	}
	historyCmd.Flags().StringP("agent", "a", "", "Filter by agent")
	historyCmd.Flags().IntP("limit", "l", 10, "Limit number of records")
	cmd.AddCommand(historyCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the deployment command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showDeploymentDocs()
		},
	})

	return cmd
}

func showDeploymentDocs() {
	title := "SLOTH-RUNNER SYSADMIN DEPLOYMENT(1)"
	description := "sloth-runner sysadmin deployment - Deployment and rollback management"
	synopsis := "sloth-runner sysadmin deployment [subcommand] [options]"

	options := [][]string{
		{"deploy", "Deploy new version to agents with rolling updates, canary deployments, or blue-green strategies."},
		{"rollback", "Rollback to previous version with one-click rollback, version history, and safety checks."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"Deploy to production",
			"sloth-runner sysadmin deployment deploy --env production --strategy rolling",
			"Deploys to production using rolling update strategy",
		},
		{
			"Canary deployment",
			"sloth-runner sysadmin deploy deploy --version v2.0.0 --canary 10%",
			"Deploys to 10% of agents first, then gradually increases",
		},
		{
			"Blue-green deployment",
			"sloth-runner sysadmin deployment deploy --strategy blue-green --all-agents",
			"Deploys to parallel infrastructure, then switches traffic",
		},
		{
			"Quick rollback",
			"sloth-runner sysadmin deploy rollback --version v1.2.3",
			"Rolls back to version 1.2.3 immediately",
		},
		{
			"Rollback with verification",
			"sloth-runner sysadmin deployment rollback --agent web-01 --verify",
			"Rolls back and verifies the agent is healthy",
		},
		{
			"Deployment history",
			"sloth-runner sysadmin deploy history --agent web-01 --limit 10",
			"Shows last 10 deployments on web-01",
		},
	}

	seeAlso := []string{
		"sloth-runner agent update - Update agents",
		"sloth-runner sysadmin backup - Backup and restore",
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
	pterm.FgGray.Println("For more information: sloth-runner sysadmin deployment --help")
}

// runDeploy executa deployment
func runDeploy(version string, agents []string, strategyStr string, dryRun bool, healthCheck bool) error {
	manager := NewDeploymentManager()

	// Converte strategy string para tipo
	var strategy DeploymentStrategy
	switch strategyStr {
	case "rolling":
		strategy = StrategyRolling
	case "canary":
		strategy = StrategyCanary
	case "blue-green":
		strategy = StrategyBlueGreen
	default:
		strategy = StrategyDirect
	}

	pterm.DefaultHeader.WithFullWidth().Println("Deployment")
	pterm.Println()

	// Info box
	pterm.DefaultBox.WithTitle("Deployment Configuration").WithTitleTopCenter().Println(
		fmt.Sprintf("Version: %s\nAgents: %v\nStrategy: %s\nDry Run: %v",
			version, agents, strategy, dryRun),
	)
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Deploying...")

	options := DeployOptions{
		Version:      version,
		Agents:       agents,
		Strategy:     strategy,
		DryRun:       dryRun,
		HealthCheck:  healthCheck,
		BackupBefore: true,
	}

	result, err := manager.Deploy(options)
	if err != nil {
		spinner.Fail("Deployment failed")
		return err
	}

	if result.Success {
		spinner.Success("✅ Deployment completed successfully")
	} else {
		spinner.Warning("⚠️  Deployment completed with errors")
	}

	pterm.Println()

	// Results table
	tableData := pterm.TableData{
		{"Metric", "Value"},
		{"Version Deployed", result.Version},
		{"Agents Updated", fmt.Sprintf("%d", len(result.AgentsUpdated))},
		{"Agents Failed", fmt.Sprintf("%d", len(result.AgentsFailed))},
		{"Duration", result.Duration.String()},
		{"Status", map[bool]string{true: "Success", false: "Partial Failure"}[result.Success]},
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	if len(result.AgentsUpdated) > 0 {
		pterm.Success.Println("Successfully updated agents:")
		for _, agent := range result.AgentsUpdated {
			pterm.Println("  • " + agent)
		}
		pterm.Println()
	}

	if len(result.AgentsFailed) > 0 {
		pterm.Error.Println("Failed agents:")
		for _, agent := range result.AgentsFailed {
			pterm.Println("  • " + agent)
		}
		pterm.Println()
	}

	pterm.Info.Println(result.Message)

	return nil
}

// runRollback executa rollback
func runRollback(agent string, version string, verify bool) error {
	manager := NewDeploymentManager()

	pterm.DefaultHeader.WithFullWidth().Println("Rollback")
	pterm.Println()

	targetMsg := "previous successful version"
	if version != "" {
		targetMsg = version
	}

	pterm.Warning.Printf("Rolling back agent '%s' to %s\n", agent, targetMsg)
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Rolling back...")

	options := RollbackOptions{
		Agent:   agent,
		Version: version,
		Verify:  verify,
	}

	result, err := manager.Rollback(options)
	if err != nil {
		spinner.Fail("Rollback failed")
		return err
	}

	spinner.Success("✅ Rollback completed successfully")
	pterm.Println()

	// Results
	pterm.Success.Printf("Agent '%s' rolled back to version %s\n", agent, result.Version)
	pterm.Info.Printf("Duration: %s\n", result.Duration)

	if verify {
		pterm.Success.Println("Health check passed ✓")
	}

	return nil
}

// runHistory mostra histórico de deployments
func runHistory(agent string, limit int) error {
	manager := NewDeploymentManager()

	pterm.DefaultHeader.WithFullWidth().Println("Deployment History")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading deployment history...")

	history, err := manager.GetHistory(agent, limit)
	if err != nil {
		spinner.Fail("Failed to load history")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Loaded %d deployment records", len(history)))
	pterm.Println()

	if len(history) == 0 {
		pterm.Info.Println("No deployment history found")
		return nil
	}

	// History table
	tableData := pterm.TableData{
		{"Timestamp", "Agent", "Version", "Strategy", "Status", "Duration"},
	}

	for _, record := range history {
		status := "✓ Success"
		if !record.Success {
			status = "✗ Failed"
		}

		tableData = append(tableData, []string{
			record.Timestamp.Format("2006-01-02 15:04:05"),
			record.Agent,
			record.Version,
			string(record.Strategy),
			status,
			record.Duration.String(),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Info.Printf("Showing %d of %d total records\n", len(history), limit)

	return nil
}
