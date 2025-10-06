package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MasterServerStarter is a function type that starts the master server
// This will be injected from main package
var MasterServerStarter func(port int) error

// NewMasterCommand creates the master command
func NewMasterCommand(ctx *AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Starts the sloth-runner master server",
		Long:  `The master command starts the sloth-runner master server, which includes the agent registry.`,
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

	return cmd
}
