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

// NewValidateCommand creates the validation command
func NewValidateCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate and repair stack state",
		Long:  `Validate stack state integrity and automatically repair issues.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewValidateCheckCommand(ctx),
		NewValidateRepairCommand(ctx),
		NewValidateAllCommand(ctx),
	)

	return cmd
}

// NewValidateCheckCommand validates a stack
func NewValidateCheckCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check <stack-name>",
		Short: "Validate stack state integrity",
		Long:  `Checks for issues like orphaned dependencies, missing resources, and invalid metadata.`,
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

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Validating stack '%s'...", stackName))

			valid, issues, err := tracker.ValidateState(stack.ID)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to validate: %v", err))
				return err
			}

			if valid {
				spinner.Success("✓ Stack is valid - no issues found")
				return nil
			}

			spinner.Warning(fmt.Sprintf("⚠ Validation found %d issue(s)", len(issues)))
			fmt.Println()

			pterm.DefaultSection.Println("Validation Issues")
			for i, issue := range issues {
				pterm.Error.Printf("  %d. %s\n", i+1, issue)
			}

			fmt.Println()
			pterm.Info.Println("Use 'sloth-runner stack validate repair' to fix issues")

			return nil
		},
	}

	return cmd
}

// NewValidateRepairCommand repairs stack issues
func NewValidateRepairCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair <stack-name>",
		Short: "Automatically repair stack issues",
		Long:  `Attempts to automatically fix detected validation issues.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			force, _ := cmd.Flags().GetBool("force")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			if !force && !dryRun {
				result, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText(fmt.Sprintf("Are you sure you want to repair stack '%s'?", stackName)).
					Show()

				if !result {
					pterm.Info.Println("Repair cancelled")
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

			// First validate to find issues
			spinner, _ := pterm.DefaultSpinner.Start("Detecting issues...")
			valid, issues, err := tracker.ValidateState(stack.ID)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to validate: %v", err))
				return err
			}

			if valid {
				spinner.Success("✓ Stack is already valid - nothing to repair")
				return nil
			}

			spinner.UpdateText(fmt.Sprintf("Found %d issue(s)", len(issues)))

			if dryRun {
				fmt.Println()
				pterm.Info.Println("DRY RUN: Would repair the following issues:")
				for i, issue := range issues {
					pterm.Info.Printf("  %d. %s\n", i+1, issue)
				}
				return nil
			}

			spinner.UpdateText("Repairing issues...")

			// TODO: Implement auto-repair logic
			err = fmt.Errorf("auto-repair not yet implemented")
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to repair: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("✓ Repaired %d issue(s)", len(issues)))

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	cmd.Flags().Bool("dry-run", false, "Show what would be repaired without making changes")

	return cmd
}

// NewValidateAllCommand validates all stacks
func NewValidateAllCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Validate all stacks",
		Long:  `Runs validation on all stacks and reports issues.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tracker, err := services.GetGlobalStateTracker()
			if err != nil {
				return err
			}

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stacks, err := stackService.ListStacks()
			if err != nil {
				return err
			}

			pterm.DefaultHeader.WithFullWidth().Println("Validating All Stacks")
			fmt.Println()

			totalStacks := len(stacks)
			validStacks := 0
			invalidStacks := 0
			totalIssues := 0

			for i, stack := range stacks {
				spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("[%d/%d] Validating %s...", i+1, totalStacks, stack.Name))

				valid, issues, err := tracker.ValidateState(stack.ID)
				if err != nil {
					spinner.Fail(fmt.Sprintf("Error: %v", err))
					continue
				}

				if valid {
					spinner.Success(fmt.Sprintf("✓ %s", stack.Name))
					validStacks++
				} else {
					spinner.Warning(fmt.Sprintf("⚠ %s (%d issue(s))", stack.Name, len(issues)))
					invalidStacks++
					totalIssues += len(issues)
				}
			}

			fmt.Println()
			pterm.DefaultBox.WithTitle("VALIDATION SUMMARY").WithTitleTopCenter().Println(
				fmt.Sprintf("Total Stacks: %d\nValid: %s\nInvalid: %s\nTotal Issues: %s",
					totalStacks,
					pterm.Green(fmt.Sprintf("%d", validStacks)),
					pterm.Red(fmt.Sprintf("%d", invalidStacks)),
					pterm.Yellow(fmt.Sprintf("%d", totalIssues)),
				),
			)

			if invalidStacks > 0 {
				fmt.Println()
				pterm.Info.Println("Use 'sloth-runner stack validate check <stack-name>' for details")
			}

			return nil
		},
	}

	return cmd
}
