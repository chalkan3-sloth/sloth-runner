package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command
// This uses the Factory Pattern to create commands
func NewRootCommand(ctx *AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sloth-runner",
		Short: "A flexible sloth-runner with Lua scripting capabilities",
		Long: `sloth-runner is a command-line tool that allows you to define and execute
tasks using Lua scripts. It supports pipelines, workflows, dynamic task generation,
and output manipulation.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if --version flag is set
			versionFlag, _ := cmd.Flags().GetBool("version")
			if versionFlag {
				fmt.Printf("sloth-runner version %s\n", ctx.Version)
				fmt.Printf("Git commit: %s\n", ctx.Commit)
				fmt.Printf("Build date: %s\n", ctx.Date)
				return
			}
			cmd.Help()
		},
	}

	// Add persistent flags
	cmd.PersistentFlags().BoolP("version", "V", false, "Show version information")

	return cmd
}
