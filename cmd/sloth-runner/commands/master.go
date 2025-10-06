package commands

import (
	"github.com/spf13/cobra"
)

// NewMasterCommand creates the master command
// TODO: Extract logic from main.go masterCmd (lines ~341-469)
// Complex command - starts gRPC server, agent registry
func NewMasterCommand(ctx *AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Starts the sloth-runner master server",
		Long:  `The master command starts the sloth-runner master server, which includes the agent registry.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement
			// - Start gRPC server
			// - Agent registry management
			// - Daemon mode support
			return nil
		},
	}

	cmd.Flags().IntP("port", "p", 50053, "Port for the master gRPC server")
	cmd.Flags().String("bind", "0.0.0.0", "Address to bind the master server")
	cmd.Flags().Bool("daemon", false, "Run master server as daemon")

	return cmd
}
