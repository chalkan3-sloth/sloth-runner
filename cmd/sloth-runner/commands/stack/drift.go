//go:build cgo
// +build cgo

package stack

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewDriftCommand creates the drift detection command
func NewDriftCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect and manage state drift",
		Long:  `Detect when actual state differs from declared state (Terraform-like drift detection).`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewDriftDetectCommand(ctx),
		NewDriftShowCommand(ctx),
		NewDriftFixCommand(ctx),
	)

	return cmd
}

// NewDriftDetectCommand detects drift
func NewDriftDetectCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detect <stack-name>",
		Short: "Detect state drift for a stack",
		Long:  `Scans the stack to detect if actual state differs from declared state.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

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

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Detecting drift for stack '%s'...", stackName))

			hasDrift, drifts, err := tracker.DetectDriftWithEvent(stack.ID)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to detect drift: %v", err))
				return err
			}

			if !hasDrift {
				spinner.Success("No drift detected - state is in sync")
				return nil
			}

			spinner.Warning(fmt.Sprintf("Drift detected: %d resource(s) have drifted", len(drifts)))
			fmt.Println()

			pterm.DefaultSection.Println("Drifted Resources")
			for _, drift := range drifts {
				pterm.Warning.Printf("  • %s\n", drift)
			}

			fmt.Println()
			pterm.Info.Println("Use 'sloth-runner stack drift fix' to automatically repair drift")

			return nil
		},
	}

	return cmd
}

// NewDriftShowCommand shows drift details
func NewDriftShowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <stack-name>",
		Short: "Show detailed drift information",
		Long:  `Displays detailed information about detected drift.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			driftInfos, err := stackService.GetDriftInfo(stack.ID)
			if err != nil {
				return err
			}

			if len(driftInfos) == 0 {
				pterm.Info.Println("No drift information available. Run 'drift detect' first.")
				return nil
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Drift Information: %s", stackName)
			fmt.Println()

			hasDrift := false
			for _, driftInfo := range driftInfos {
				if driftInfo.IsDrifted && driftInfo.ResolutionStatus == "pending" {
					hasDrift = true
				}
			}

			pterm.Info.Printf("Total Drift Records: %d\n", len(driftInfos))
			pterm.Info.Printf("Has Pending Drift: %v\n", hasDrift)
			fmt.Println()

			if hasDrift {
				pterm.DefaultSection.Println("Drifted Resources")
				for _, driftInfo := range driftInfos {
					if driftInfo.IsDrifted && driftInfo.ResolutionStatus == "pending" {
						pterm.Warning.Printf("  • Resource: %s\n", driftInfo.ResourceID)
						pterm.Warning.Printf("    Detected: %s\n", driftInfo.DetectedAt.Format("2006-01-02 15:04:05"))
						pterm.Warning.Printf("    Drifted Fields: %v\n", driftInfo.DriftedFields)
						fmt.Println()
					}
				}
			} else {
				pterm.Success.Println("No pending drift detected")
			}

			return nil
		},
	}

	return cmd
}

// NewDriftFixCommand auto-fixes drift
func NewDriftFixCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fix <stack-name>",
		Short: "Automatically fix detected drift",
		Long:  `Attempts to automatically repair state drift by reconciling differences.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			force, _ := cmd.Flags().GetBool("force")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			if !force && !dryRun {
				result, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText(fmt.Sprintf("Are you sure you want to fix drift for stack '%s'?", stackName)).
					Show()

				if !result {
					pterm.Info.Println("Drift fix cancelled")
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

			if dryRun {
				pterm.Info.Println("DRY RUN: Would fix drift for stack")
				return nil
			}

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fixing drift for stack '%s'...", stackName))

			// TODO: Implement auto-fix drift in backend using stackID
			err = fmt.Errorf("auto-fix drift not yet implemented")
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to fix drift: %v", err))
				return err
			}

			spinner.Success("Drift fixed successfully")

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	cmd.Flags().Bool("dry-run", false, "Show what would be fixed without making changes")

	return cmd
}
