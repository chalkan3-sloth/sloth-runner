package agent

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewListCommand creates the agent list command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all registered agents",
		Long:  `Lists all agents that are currently registered with the master.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				pterm.DefaultLogger.Level = pterm.LogLevelDebug
				slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))
			}

			masterAddr, _ := cmd.Flags().GetString("master")

			// Create agent service
			agentService := services.NewAgentService(masterAddr)

			// List agents
			agents, err := agentService.ListAgents()
			if err != nil {
				return fmt.Errorf("failed to list agents: %w", err)
			}

			if len(agents) == 0 {
				fmt.Println("No agents registered.")
				return nil
			}

			// Display agents in table format
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "AGENT NAME\tADDRESS\tSTATUS\tLAST HEARTBEAT\tLAST INFO COLLECTED")
			fmt.Fprintln(w, "------------\t----------\t------\t--------------\t-------------------")

			for _, agent := range agents {
				status := agent.GetStatus()
				coloredStatus := status
				if status == "Active" {
					coloredStatus = pterm.Green(status)
				} else {
					coloredStatus = pterm.Red(status)
				}

				lastHeartbeat := "N/A"
				if agent.GetLastHeartbeat() > 0 {
					lastHeartbeat = time.Unix(agent.GetLastHeartbeat(), 0).Format(time.RFC3339)
				}

				lastInfoCollected := "Never"
				if agent.GetLastInfoCollected() > 0 {
					lastInfoCollected = time.Unix(agent.GetLastInfoCollected(), 0).Format(time.RFC3339)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					agent.GetAgentName(),
					agent.GetAgentAddress(),
					coloredStatus,
					lastHeartbeat,
					lastInfoCollected)
			}

			return w.Flush()
		},
	}

	cmd.Flags().String("master", "192.168.1.29:50053", "Master registry address")
	cmd.Flags().Bool("debug", false, "Enable debug logging")

	return cmd
}
