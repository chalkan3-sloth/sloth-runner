package stack

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewSnapshotCommand creates the snapshot management command
func NewSnapshotCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage stack snapshots (Terraform/Pulumi-like state versioning)",
		Long:  `Create, list, restore, and manage stack snapshots for version control and rollback.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewSnapshotCreateCommand(ctx),
		NewSnapshotListCommand(ctx),
		NewSnapshotShowCommand(ctx),
		NewSnapshotRestoreCommand(ctx),
		NewSnapshotDeleteCommand(ctx),
		NewSnapshotCompareCommand(ctx),
	)

	return cmd
}

// NewSnapshotCreateCommand creates a snapshot
func NewSnapshotCreateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <stack-name>",
		Short: "Create a new snapshot of a stack",
		Long:  `Creates a point-in-time snapshot of the stack state for later rollback.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			description, _ := cmd.Flags().GetString("description")
			creator, _ := cmd.Flags().GetString("creator")

			if creator == "" {
				creator = "cli-user"
			}

			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			// Get stack
			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Creating snapshot for stack '%s'...", stackName))

			version, err := tracker.CreateSnapshotWithEvent(stack.ID, creator, description)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to create snapshot: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("Snapshot created successfully (version %d)", version))
			fmt.Println()

			pterm.Info.Printf("Stack: %s\n", pterm.Cyan(stackName))
			pterm.Info.Printf("Version: %s\n", pterm.Green(fmt.Sprintf("%d", version)))
			if description != "" {
				pterm.Info.Printf("Description: %s\n", description)
			}

			return nil
		},
	}

	cmd.Flags().StringP("description", "d", "", "Snapshot description")
	cmd.Flags().String("creator", "", "Creator name (default: cli-user)")

	return cmd
}

// NewSnapshotListCommand lists snapshots
func NewSnapshotListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <stack-name>",
		Short: "List all snapshots for a stack",
		Long:  `Lists all available snapshots with their versions and metadata.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			outputFormat, _ := cmd.Flags().GetString("output")

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			snapshots, err := stackService.ListSnapshots(stack.ID)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(snapshots, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if len(snapshots) == 0 {
				pterm.Info.Printf("No snapshots found for stack '%s'\n", stackName)
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Snapshots: %s", stackName)
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "VERSION\tCREATED\tCREATOR\tDESCRIPTION")
			fmt.Fprintln(w, "-------\t-------\t-------\t-----------")

			for _, snap := range snapshots {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					snap.Version,
					snap.CreatedAt.Format("2006-01-02 15:04:05"),
					snap.CreatedBy,
					snap.Description,
				)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d snapshot(s)\n", len(snapshots))

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format (table or json)")

	return cmd
}

// NewSnapshotShowCommand shows snapshot details
func NewSnapshotShowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <stack-name> <version>",
		Short: "Show detailed information about a snapshot",
		Long:  `Displays detailed information about a specific snapshot version.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			version := args[1]
			outputFormat, _ := cmd.Flags().GetString("output")

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			var versionNum int
			fmt.Sscanf(version, "%d", &versionNum)

			snapshot, err := stackService.GetSnapshot(stack.ID, versionNum)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				data, err := json.MarshalIndent(snapshot, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Snapshot Details: %s (v%d)", stackName, versionNum)
			fmt.Println()

			pterm.Info.Printf("Version: %s\n", pterm.Green(fmt.Sprintf("%d", snapshot.Version)))
			pterm.Info.Printf("Created: %s\n", snapshot.CreatedAt.Format("2006-01-02 15:04:05"))
			pterm.Info.Printf("Creator: %s\n", snapshot.CreatedBy)
			pterm.Info.Printf("Description: %s\n", snapshot.Description)
			fmt.Println()

			// Show state summary
			if snapshot.StackState != nil {
				pterm.DefaultSection.Println("Stack State Summary")
				pterm.Info.Printf("Name: %s\n", snapshot.StackState.Name)
				pterm.Info.Printf("State Version: %s\n", snapshot.StackState.Version)
				pterm.Info.Printf("Status: %s\n", snapshot.StackState.Status)
				pterm.Info.Printf("Resources: %d\n", len(snapshot.Resources))
			}

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format (table or json)")

	return cmd
}

// NewSnapshotRestoreCommand restores a snapshot
func NewSnapshotRestoreCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore <stack-name> <version>",
		Short: "Restore a stack to a previous snapshot",
		Long:  `Rolls back the stack state to a previous snapshot version.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			version := args[1]
			force, _ := cmd.Flags().GetBool("force")
			performer, _ := cmd.Flags().GetString("performer")

			if performer == "" {
				performer = "cli-user"
			}

			if !force {
				result, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText(fmt.Sprintf("Are you sure you want to restore stack '%s' to version %s?", stackName, version)).
					Show()

				if !result {
					pterm.Info.Println("Restore cancelled")
					return nil
				}
			}

			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			var versionNum int
			fmt.Sscanf(version, "%d", &versionNum)

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Restoring stack '%s' to version %d...", stackName, versionNum))

			err = tracker.RollbackToSnapshotWithEvent(stack.ID, versionNum, performer)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to restore snapshot: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("Stack '%s' restored to version %d successfully", stackName, versionNum))

			pterm.Warning.Println("⚠ Stack state has been rolled back. Review changes before proceeding.")

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	cmd.Flags().String("performer", "", "Performer name (default: cli-user)")

	return cmd
}

// NewSnapshotDeleteCommand deletes a snapshot
func NewSnapshotDeleteCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <stack-name> <version>",
		Short: "Delete a specific snapshot",
		Long:  `Deletes a snapshot version. Cannot delete the current version.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			version := args[1]
			force, _ := cmd.Flags().GetBool("force")

			if !force {
				result, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText(fmt.Sprintf("Are you sure you want to delete snapshot version %s?", version)).
					Show()

				if !result {
					pterm.Info.Println("Deletion cancelled")
					return nil
				}
			}

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			_, err = stackService.GetStack(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			var versionNum int
			fmt.Sscanf(version, "%d", &versionNum)

			// TODO: Implement DeleteSnapshot in backend using stackID and versionNum
			pterm.Warning.Println("Snapshot deletion not yet implemented")

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}

// NewSnapshotCompareCommand compares two snapshots
func NewSnapshotCompareCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare <stack-name> <version1> <version2>",
		Short: "Compare two snapshots",
		Long:  `Shows differences between two snapshot versions.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			version1 := args[1]
			version2 := args[2]

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			var v1, v2 int
			fmt.Sscanf(version1, "%d", &v1)
			fmt.Sscanf(version2, "%d", &v2)

			snap1, err := stackService.GetSnapshot(stack.ID, v1)
			if err != nil {
				return fmt.Errorf("snapshot v%d not found: %w", v1, err)
			}

			snap2, err := stackService.GetSnapshot(stack.ID, v2)
			if err != nil {
				return fmt.Errorf("snapshot v%d not found: %w", v2, err)
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Snapshot Comparison: %s", stackName)
			fmt.Println()

			pterm.Info.Printf("Version %d: %s (created %s by %s)\n",
				snap1.Version,
				snap1.Description,
				snap1.CreatedAt.Format("2006-01-02 15:04:05"),
				snap1.CreatedBy,
			)

			pterm.Info.Printf("Version %d: %s (created %s by %s)\n",
				snap2.Version,
				snap2.Description,
				snap2.CreatedAt.Format("2006-01-02 15:04:05"),
				snap2.CreatedBy,
			)

			fmt.Println()
			pterm.Info.Println("Detailed diff analysis coming soon...")

			return nil
		},
	}

	return cmd
}
