package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewUpdateCommand creates the agent update command
func NewUpdateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <agent_name>",
		Short: "Update an agent to the latest version",
		Long:  `Updates the specified agent to the latest available version from GitHub releases.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]
			version, _ := cmd.Flags().GetString("version")
			restart, _ := cmd.Flags().GetBool("restart")
			local, _ := cmd.Flags().GetBool("local")

			// Get master address (supports both names and addresses)
			masterAddr := getMasterAddress(cmd)

			return updateAgent(agentName, masterAddr, version, restart, local)
		},
	}

	cmd.Flags().String("master", "", "Master server address (if empty, uses local database)")
	cmd.Flags().String("version", "latest", "Version to update to (default: latest)")
	cmd.Flags().Bool("restart", true, "Restart agent service after update")
	cmd.Flags().Bool("local", false, "Force using local database instead of master server")

	return cmd
}

func updateAgent(agentName, masterAddr, version string, restart bool, local bool) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create connection factory
	factory := NewDefaultConnectionFactory()

	// If --local flag is set, use local database
	if local {
		return updateAgentLocal(ctx, factory, agentName, version, restart)
	}

	// If no master address, use local database
	if masterAddr == "" {
		pterm.Info.Println("No master address configured, using local database")
		return updateAgentLocal(ctx, factory, agentName, version, restart)
	}

	// Try to connect to master server
	registryClient, cleanup, err := factory.CreateRegistryClient(masterAddr)
	if err != nil {
		pterm.Warning.Printf("Could not connect to master at %s, using local database\n", masterAddr)
		return updateAgentLocal(ctx, factory, agentName, version, restart)
	}
	defer cleanup()

	// Create agent client factory function
	agentClientFactory := func(addr string) (AgentClient, func(), error) {
		return factory.CreateAgentClient(addr)
	}

	// Use refactored function with injected clients
	opts := UpdateAgentOptions{
		AgentName:     agentName,
		TargetVersion: version,
		Restart:       restart,
		Writer:        os.Stdout,
	}

	_, err = updateAgentWithClients(ctx, registryClient, agentClientFactory, opts)

	// Track operation
	trackAgentUpdate(agentName, version, err == nil)

	return err
}

// updateAgentLocal updates an agent without using the master server
// It connects directly to the agent using address from local database
func updateAgentLocal(ctx context.Context, factory ConnectionFactory, agentName, version string, restart bool) error {
	pterm.DefaultHeader.WithFullWidth().Printf("Agent Update - %s (Local Mode)", agentName)
	fmt.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Looking up agent in local database...")

	// Get agent address from local database
	agentAddress, err := getAgentAddressFromLocalDB(agentName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to find agent: %v", err))
		return err
	}

	spinner.UpdateText(fmt.Sprintf("Connecting to agent at %s...", agentAddress))

	// Connect directly to agent
	agentClient, cleanup, err := factory.CreateAgentClient(agentAddress)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to connect to agent: %v", err))
		return fmt.Errorf("failed to connect to agent at %s: %w", agentAddress, err)
	}
	defer cleanup()

	spinner.UpdateText("Initiating agent update...")

	// Call UpdateAgent
	opts := UpdateAgentOptions{
		AgentName:     agentName,
		TargetVersion: version,
		Restart:       restart,
		Writer:        os.Stdout,
	}

	result, err := performAgentUpdate(ctx, agentClient, opts)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to initiate update: %v", err))
		return err
	}

	if !result.Success {
		spinner.Fail("Update failed")
		pterm.Error.Println(result.Message)
		return fmt.Errorf("update failed: %s", result.Message)
	}

	spinner.Success("Update completed successfully")
	fmt.Println()

	// Display update summary
	displayUpdateSummary(result, opts, os.Stdout)

	pterm.Success.Printf("✅ Agent %s updated successfully\n", agentName)

	// Track operation
	trackAgentUpdate(agentName, version, true)

	return nil
}
