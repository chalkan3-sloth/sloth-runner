//go:build cgo
// +build cgo

package stack

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewStateCommand creates the state management command
func NewStateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Manage stack state (Pulumi/Terraform-like)",
		Long:  `State management provides Pulumi/Terraform-like state tracking with versioning, drift detection, and rollback capabilities.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewStateListCommand(ctx),
		NewStateShowCommand(ctx),
		NewStateVersionsCommand(ctx),
		NewStateRollbackCommand(ctx),
		NewStateDriftCommand(ctx),
		NewStateLockCommand(ctx),
		NewStateUnlockCommand(ctx),
		NewStateSnapshotCommand(ctx),
		NewStateTagsCommand(ctx),
		NewStateActivityCommand(ctx),
	)

	return cmd
}

// NewStateListCommand lists all stacks
func NewStateListCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all stacks",
		Long:  `Lists all stacks in the state backend with their current status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			stacks, err := backend.GetStackManager().ListStacks()
			if err != nil {
				return fmt.Errorf("failed to list stacks: %w", err)
			}

			if len(stacks) == 0 {
				pterm.Info.Println("No stacks found")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Println("Stack States")
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tSTATUS\tVERSION\tEXECUTIONS\tLAST UPDATED")
			fmt.Fprintln(w, "----\t------\t-------\t----------\t------------")

			for _, s := range stacks {
				status := s.Status
				statusColor := pterm.Green
				if status == "failed" {
					statusColor = pterm.Red
				} else if status == "running" {
					statusColor = pterm.Yellow
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
					pterm.Cyan(s.Name),
					statusColor(status),
					s.Version,
					s.ExecutionCount,
					s.UpdatedAt.Format("2006-01-02 15:04"),
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d stack(s)\n", len(stacks))

			return nil
		},
	}
}

// NewStateShowCommand shows details of a specific stack
func NewStateShowCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "show <stack-name>",
		Short: "Show stack details",
		Long:  `Displays detailed information about a specific stack including resources and outputs.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			// Show stack info
			pterm.DefaultHeader.WithFullWidth().Printfln("Stack: %s", s.Name)
			fmt.Println()

			info := [][]string{
				{"ID", s.ID},
				{"Name", s.Name},
				{"Status", s.Status},
				{"Version", s.Version},
				{"Executions", fmt.Sprintf("%d", s.ExecutionCount)},
				{"Created", s.CreatedAt.Format(time.RFC3339)},
				{"Updated", s.UpdatedAt.Format(time.RFC3339)},
				{"Workflow File", s.WorkflowFile},
			}

			if s.LastError != "" {
				info = append(info, []string{"Last Error", s.LastError})
			}

			pterm.DefaultTable.WithHasHeader(false).WithData(info).Render()

			// Show resources
			resources, err := backend.GetStackManager().ListResources(s.ID)
			if err == nil && len(resources) > 0 {
				fmt.Println()
				pterm.DefaultSection.Println("Resources")

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(w, "TYPE\tNAME\tMODULE\tSTATE")
				fmt.Fprintln(w, "----\t----\t------\t-----")

				for _, r := range resources {
					stateColor := pterm.Green
					if r.State == "failed" {
						stateColor = pterm.Red
					} else if r.State == "pending" {
						stateColor = pterm.Yellow
					}

					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
						r.Type,
						pterm.Cyan(r.Name),
						r.Module,
						stateColor(r.State),
					)
				}
				w.Flush()
			}

			// Show outputs
			if len(s.Outputs) > 0 {
				fmt.Println()
				pterm.DefaultSection.Println("Outputs")

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
				fmt.Fprintln(w, "KEY\tVALUE")
				fmt.Fprintln(w, "---\t-----")

				for key, value := range s.Outputs {
					fmt.Fprintf(w, "%s\t%v\n", pterm.Cyan(key), value)
				}
				w.Flush()
			}

			return nil
		},
	}
}

// NewStateVersionsCommand lists all versions of a stack
func NewStateVersionsCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "versions <stack-name>",
		Short: "List stack versions",
		Long:  `Lists all historical versions of a stack for rollback purposes.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			snapshots, err := backend.ListSnapshots(s.ID)
			if err != nil {
				return fmt.Errorf("failed to list snapshots: %w", err)
			}

			if len(snapshots) == 0 {
				pterm.Info.Printfln("No versions found for stack: %s", stackName)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Versions: %s", stackName)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "VERSION\tCREATED\tCREATED BY\tDESCRIPTION")
			fmt.Fprintln(w, "-------\t-------\t----------\t-----------")

			for _, snap := range snapshots {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					snap.Version,
					snap.CreatedAt.Format("2006-01-02 15:04"),
					snap.CreatedBy,
					snap.Description,
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d version(s)\n", len(snapshots))

			return nil
		},
	}
}

// NewStateRollbackCommand rolls back to a specific version
func NewStateRollbackCommand(ctx *commands.AppContext) *cobra.Command {
	var version int

	cmd := &cobra.Command{
		Use:   "rollback <stack-name>",
		Short: "Rollback to a previous version",
		Long:  `Rolls back a stack to a previous version. Creates a backup before rollback.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			if version <= 0 {
				return fmt.Errorf("version must be specified with --version flag")
			}

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			pterm.Warning.Printfln("Rolling back stack '%s' to version %d", stackName, version)
			pterm.Info.Println("A backup will be created before rollback")

			if err := backend.RollbackToSnapshot(s.ID, version, "cli-user"); err != nil {
				return fmt.Errorf("rollback failed: %w", err)
			}

			pterm.Success.Printfln("Successfully rolled back to version %d", version)

			return nil
		},
	}

	cmd.Flags().IntVarP(&version, "version", "v", 0, "Version to rollback to (required)")
	cmd.MarkFlagRequired("version")

	return cmd
}

// NewStateDriftCommand detects drift in stack resources
func NewStateDriftCommand(ctx *commands.AppContext) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "drift <stack-name>",
		Short: "Detect drift in stack resources",
		Long:  `Detects drift between expected and actual state of resources.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			drifts, err := backend.GetDriftInfo(s.ID)
			if err != nil {
				return fmt.Errorf("failed to get drift info: %w", err)
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(drifts, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			driftedCount := 0
			for _, d := range drifts {
				if d.IsDrifted {
					driftedCount++
				}
			}

			if driftedCount == 0 {
				pterm.Success.Printfln("No drift detected for stack: %s", stackName)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Drift Detection: %s", stackName)
			pterm.Warning.Printfln("\n%d resource(s) have drifted from expected state", driftedCount)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "RESOURCE ID\tDRIFTED FIELDS\tSTATUS\tDETECTED")
			fmt.Fprintln(w, "-----------\t--------------\t------\t--------")

			for _, d := range drifts {
				if d.IsDrifted {
					fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
						pterm.Cyan(d.ResourceID),
						len(d.DriftedFields),
						pterm.Yellow(d.ResolutionStatus),
						d.DetectedAt.Format("2006-01-02 15:04"),
					)
				}
			}

			w.Flush()
			fmt.Println()
			pterm.Info.Println("Use 'sloth-runner run' to apply changes and fix drift")

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table or json)")

	return cmd
}

// NewStateLockCommand locks a stack state
func NewStateLockCommand(ctx *commands.AppContext) *cobra.Command {
	var operation, who string
	var duration time.Duration

	cmd := &cobra.Command{
		Use:   "lock <stack-name>",
		Short: "Lock stack state",
		Long:  `Acquires a lock on stack state to prevent concurrent modifications.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			lockID := fmt.Sprintf("lock-%d", time.Now().Unix())

			if err := backend.LockState(s.ID, lockID, operation, who, duration); err != nil {
				return fmt.Errorf("failed to lock state: %w", err)
			}

			pterm.Success.Printfln("State locked successfully (Lock ID: %s)", lockID)
			pterm.Info.Printfln("Lock will expire in %s", duration)

			return nil
		},
	}

	cmd.Flags().StringVarP(&operation, "operation", "o", "manual", "Operation name")
	cmd.Flags().StringVarP(&who, "who", "w", "cli-user", "Who is locking")
	cmd.Flags().DurationVarP(&duration, "duration", "d", 30*time.Minute, "Lock duration")

	return cmd
}

// NewStateUnlockCommand unlocks a stack state
func NewStateUnlockCommand(ctx *commands.AppContext) *cobra.Command {
	var lockID string

	cmd := &cobra.Command{
		Use:   "unlock <stack-name>",
		Short: "Unlock stack state",
		Long:  `Releases a lock on stack state.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			if lockID == "" {
				return fmt.Errorf("lock ID must be specified with --lock-id flag")
			}

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			if err := backend.UnlockState(s.ID, lockID); err != nil {
				return fmt.Errorf("failed to unlock state: %w", err)
			}

			pterm.Success.Println("State unlocked successfully")

			return nil
		},
	}

	cmd.Flags().StringVarP(&lockID, "lock-id", "l", "", "Lock ID to release (required)")
	cmd.MarkFlagRequired("lock-id")

	return cmd
}

// NewStateSnapshotCommand creates a manual snapshot
func NewStateSnapshotCommand(ctx *commands.AppContext) *cobra.Command {
	var description string

	cmd := &cobra.Command{
		Use:   "snapshot <stack-name>",
		Short: "Create state snapshot",
		Long:  `Creates a manual snapshot of the current stack state.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			version, err := backend.CreateSnapshot(s.ID, "cli-user", description)
			if err != nil {
				return fmt.Errorf("failed to create snapshot: %w", err)
			}

			pterm.Success.Printfln("Snapshot created: version %d", version)

			return nil
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "Manual snapshot", "Snapshot description")

	return cmd
}

// NewStateTagsCommand manages stack tags
func NewStateTagsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage stack tags",
		Long:  `Add, remove, or list tags for stacks.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(
		&cobra.Command{
			Use:   "add <stack-name> <tag>",
			Short: "Add a tag to a stack",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				stackName, tag := args[0], args[1]

				backend, err := stack.NewStateBackend("")
				if err != nil {
					return err
				}
				defer backend.Close()

				s, err := backend.GetStackManager().GetStackByName(stackName)
				if err != nil {
					return err
				}

				if err := backend.AddTag(s.ID, tag); err != nil {
					return err
				}

				pterm.Success.Printfln("Tag '%s' added to stack '%s'", tag, stackName)
				return nil
			},
		},
		&cobra.Command{
			Use:   "list <stack-name>",
			Short: "List tags for a stack",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				stackName := args[0]

				backend, err := stack.NewStateBackend("")
				if err != nil {
					return err
				}
				defer backend.Close()

				s, err := backend.GetStackManager().GetStackByName(stackName)
				if err != nil {
					return err
				}

				tags, err := backend.GetTags(s.ID)
				if err != nil {
					return err
				}

				if len(tags) == 0 {
					pterm.Info.Printfln("No tags for stack '%s'", stackName)
					return nil
				}

				pterm.DefaultHeader.Printfln("Tags: %s", stackName)
				for _, tag := range tags {
					fmt.Printf("  - %s\n", pterm.Cyan(tag))
				}

				return nil
			},
		},
	)

	return cmd
}

// NewStateActivityCommand shows activity log
func NewStateActivityCommand(ctx *commands.AppContext) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "activity <stack-name>",
		Short: "Show activity log",
		Long:  `Displays the activity log for a stack showing all operations.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			backend, err := stack.NewStateBackend("")
			if err != nil {
				return fmt.Errorf("failed to initialize state backend: %w", err)
			}
			defer backend.Close()

			s, err := backend.GetStackManager().GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("failed to get stack: %w", err)
			}

			activities, err := backend.GetActivity(s.ID, limit)
			if err != nil {
				return fmt.Errorf("failed to get activity: %w", err)
			}

			if len(activities) == 0 {
				pterm.Info.Printfln("No activity for stack: %s", stackName)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Activity Log: %s", stackName)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "TYPE\tRESOURCE\tUSER\tTIME")
			fmt.Fprintln(w, "----\t--------\t----\t----")

			for _, act := range activities {
				actType := act["type"].(string)
				user := act["user"].(string)
				createdAt := act["created_at"].(time.Time)
				resourceID := ""
				if rid, ok := act["resource_id"].(string); ok {
					resourceID = rid
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					pterm.Cyan(actType),
					resourceID,
					user,
					createdAt.Format("2006-01-02 15:04"),
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Showing %d activities\n", len(activities))

			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Maximum number of activities to show")

	return cmd
}
