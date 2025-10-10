package stack

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewOperationsCommand creates the operations tracking command
func NewOperationsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operations",
		Short: "View all tracked operations across the system",
		Long:  `Operations command provides a unified view of all operations tracked by the state system.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewOperationsListCommand(ctx),
		NewOperationsStatsCommand(ctx),
		NewOperationsSearchCommand(ctx),
		NewOperationsDashboardCommand(ctx),
	)

	return cmd
}

// NewOperationsListCommand lists operations by type
func NewOperationsListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <operation-type>",
		Short: "List operations of a specific type",
		Long:  `Lists all tracked operations of a specific type (workflow, agent, scheduler, secret, hook, sloth, backup, deployment).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opTypeStr := args[0]
			limit, _ := cmd.Flags().GetInt("limit")
			outputFormat, _ := cmd.Flags().GetString("output")

			// Parse operation type
			opType := parseOperationType(opTypeStr)
			if opType == "" {
				return fmt.Errorf("invalid operation type: %s", opTypeStr)
			}

			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			history, err := tracker.GetOperationHistory(opType, limit)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(history, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(history) == 0 {
				pterm.Info.Printfln("No operations found for type: %s", opTypeStr)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Operations: %s", opTypeStr)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "RESOURCE\tTYPE\tSTATUS\tCREATED\tSTATE")
			fmt.Fprintln(w, "--------\t----\t------\t-------\t-----")

			for _, res := range history {
				statusColor := pterm.Green
				if res.State == "failed" {
					statusColor = pterm.Red
				} else if res.State == "running" {
					statusColor = pterm.Yellow
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					pterm.Cyan(res.Name),
					res.Type,
					statusColor(res.State),
					res.CreatedAt.Format("2006-01-02 15:04"),
					res.State,
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d operation(s)\n", len(history))

			return nil
		},
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of operations to show")
	cmd.Flags().StringP("output", "o", "table", "Output format (table or json)")

	return cmd
}

// NewOperationsStatsCommand shows statistics for operations
func NewOperationsStatsCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show statistics for all operations",
		Long:  `Displays comprehensive statistics for all operation types in the system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			allStats, err := tracker.GetAllOperationStats()
			if err != nil {
				return err
			}

			pterm.DefaultHeader.WithFullWidth().Println("Operation Statistics")
			fmt.Println()

			if len(allStats) == 0 {
				pterm.Info.Println("No operations tracked yet")
				return nil
			}

			// Create summary table
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "OPERATION TYPE\tTOTAL\tCOMPLETED\tFAILED\tRUNNING\tPENDING")
			fmt.Fprintln(w, "--------------\t-----\t---------\t------\t-------\t-------")

			for opType, stats := range allStats {
				fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\n",
					opType,
					stats["total"].(int),
					stats["completed"].(int),
					stats["failed"].(int),
					stats["running"].(int),
					stats["pending"].(int),
				)
			}

			w.Flush()
			fmt.Println()

			return nil
		},
	}
}

// NewOperationsSearchCommand searches operations
func NewOperationsSearchCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search operations by criteria",
		Long:  `Searches for operations matching specified criteria.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opType, _ := cmd.Flags().GetString("type")
			status, _ := cmd.Flags().GetString("status")
			outputFormat, _ := cmd.Flags().GetString("output")

			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			criteria := make(map[string]interface{})
			if opType != "" {
				criteria["type"] = opType
			}
			if status != "" {
				criteria["status"] = status
			}

			results, err := tracker.SearchOperations(criteria)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(results, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(results) == 0 {
				pterm.Info.Println("No operations found matching criteria")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Println("Search Results")
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "RESOURCE\tTYPE\tSTATUS\tCREATED")
			fmt.Fprintln(w, "--------\t----\t------\t-------")

			for _, res := range results {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					pterm.Cyan(res.Name),
					res.Type,
					res.State,
					res.CreatedAt.Format("2006-01-02 15:04:05"),
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Found: %d operation(s)\n", len(results))

			return nil
		},
	}

	cmd.Flags().StringP("type", "t", "", "Filter by operation type")
	cmd.Flags().StringP("status", "s", "", "Filter by status")
	cmd.Flags().StringP("output", "o", "table", "Output format (table or json)")

	return cmd
}

// NewOperationsDashboardCommand shows a comprehensive dashboard
func NewOperationsDashboardCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "dashboard",
		Short: "Show comprehensive operations dashboard",
		Long:  `Displays a comprehensive dashboard of all system operations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			allStats, err := tracker.GetAllOperationStats()
			if err != nil {
				return err
			}

			pterm.DefaultHeader.WithFullWidth().Println("SLOTH-RUNNER OPERATIONS DASHBOARD")
			fmt.Println()

			if len(allStats) == 0 {
				pterm.Info.Println("No operations tracked yet")
				return nil
			}

			// Calculate totals
			totalOps := 0
			totalCompleted := 0
			totalFailed := 0
			totalRunning := 0

			for _, stats := range allStats {
				totalOps += stats["total"].(int)
				totalCompleted += stats["completed"].(int)
				totalFailed += stats["failed"].(int)
				totalRunning += stats["running"].(int)
			}

			// Summary box
			pterm.DefaultBox.WithTitle("SYSTEM SUMMARY").WithTitleTopCenter().Println(
				fmt.Sprintf("Total Operations: %d\nCompleted: %s\nFailed: %s\nRunning: %s",
					totalOps,
					pterm.Green(fmt.Sprintf("%d", totalCompleted)),
					pterm.Red(fmt.Sprintf("%d", totalFailed)),
					pterm.Yellow(fmt.Sprintf("%d", totalRunning)),
				),
			)
			fmt.Println()

			// Breakdown by operation type
			pterm.DefaultSection.Println("Operations Breakdown")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "TYPE\tTOTAL\tCOMPLETED\tFAILED\tSUCCESS RATE")
			fmt.Fprintln(w, "----\t-----\t---------\t------\t------------")

			for opType, stats := range allStats {
				total := stats["total"].(int)
				completed := stats["completed"].(int)
				failed := stats["failed"].(int)

				successRate := float64(0)
				if total > 0 {
					successRate = float64(completed) / float64(total) * 100
				}

				fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%.1f%%\n",
					opType,
					total,
					pterm.Green(fmt.Sprintf("%d", completed)),
					pterm.Red(fmt.Sprintf("%d", failed)),
					successRate,
				)
			}

			w.Flush()
			fmt.Println()

			pterm.Info.Println("Use 'sloth-runner stack operations list <type>' for detailed view")

			return nil
		},
	}
}

// parseOperationType converts string to OperationType
func parseOperationType(s string) stack.OperationType {
	switch s {
	case "workflow", "workflows":
		return stack.OpWorkflowExecution
	case "agent", "agents":
		return stack.OpAgentRegistration
	case "scheduler", "scheduled":
		return stack.OpScheduledExecution
	case "secret", "secrets":
		return stack.OpSecretCreate
	case "hook", "hooks":
		return stack.OpHookRegister
	case "sloth", "sloths":
		return stack.OpSlothAdd
	case "backup", "backups":
		return stack.OpBackup
	case "deployment", "deployments":
		return stack.OpDeployment
	default:
		return ""
	}
}
