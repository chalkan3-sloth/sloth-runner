package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// MasterServerStarter is a function type that starts the master server
// This will be injected from main package
var MasterServerStarter func(port int) error

// NewMasterCommand creates the master command
func NewMasterCommand(ctx *AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Manage the sloth-runner master server",
		Long:  `The master command manages the sloth-runner master server, which includes the agent registry.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")

			if MasterServerStarter == nil {
				return fmt.Errorf("master server starter not initialized")
			}

			return MasterServerStarter(port)
		},
	}

	cmd.Flags().IntP("port", "p", 50053, "Port for the master gRPC server")
	cmd.Flags().String("bind", "0.0.0.0", "Address to bind the master server")
	cmd.Flags().Bool("daemon", false, "Run master server as daemon")

	// Add subcommands
	cmd.AddCommand(newMasterSelectCommand())
	cmd.AddCommand(newMasterShowCommand())

	return cmd
}

// newMasterSelectCommand creates the master select command
func newMasterSelectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "select <master_address>",
		Short: "Set the default master server address",
		Long:  `Sets the default master server address (e.g., 192.168.1.29:50053) and saves it to configuration file.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			masterAddr := args[0]

			// Create config directory if it doesn't exist
			if err := config.EnsureDataDir(); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			// Get config file path
			configPath := filepath.Join(config.GetDataDir(), "master.conf")

			// Write master address to config file
			if err := os.WriteFile(configPath, []byte(masterAddr), 0644); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			pterm.Success.Printf("Master server address set to: %s\n", masterAddr)
			pterm.Info.Printf("Configuration saved to: %s\n", configPath)
			pterm.Info.Println("💡 This address will now be used by default for all commands")

			return nil
		},
	}
}

// newMasterShowCommand creates the master show command
func newMasterShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the current master server address",
		Long:  `Displays the currently configured master server address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			masterAddr := config.GetMasterAddress()

			if masterAddr == "" || masterAddr == "localhost:50051" {
				pterm.Warning.Println("No master server configured (using default: localhost:50051)")
				pterm.Info.Println("💡 Use 'sloth-runner master select <address>' to set a master server")
			} else {
				pterm.Success.Printf("Current master server: %s\n", masterAddr)
			}

			return nil
		},
	}
}
