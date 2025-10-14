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

// NewLockCommand creates the state locking command
func NewLockCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage state locks (prevent concurrent modifications)",
		Long:  `Lock and unlock stack state to prevent concurrent modifications (Terraform-like locking).`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewLockAcquireCommand(ctx),
		NewLockReleaseCommand(ctx),
		NewLockStatusCommand(ctx),
		NewLockForceUnlockCommand(ctx),
	)

	return cmd
}

// NewLockAcquireCommand acquires a lock
func NewLockAcquireCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acquire <stack-name>",
		Short: "Acquire a lock on stack state",
		Long:  `Locks the stack state to prevent concurrent modifications.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			reason, _ := cmd.Flags().GetString("reason")
			lockedBy, _ := cmd.Flags().GetString("locked-by")

			if lockedBy == "" {
				lockedBy = "cli-user"
			}

			if reason == "" {
				reason = "Manual lock via CLI"
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

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Acquiring lock for stack '%s'...", stackName))

			err = tracker.LockStateWithEvent(stack.ID, lockedBy, reason)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to acquire lock: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("Lock acquired for stack '%s'", stackName))
			fmt.Println()

			pterm.Info.Printf("Locked by: %s\n", lockedBy)
			pterm.Info.Printf("Reason: %s\n", reason)
			pterm.Warning.Println("⚠ Remember to release the lock when done")

			return nil
		},
	}

	cmd.Flags().String("reason", "", "Reason for locking")
	cmd.Flags().String("locked-by", "", "Lock owner (default: cli-user)")

	return cmd
}

// NewLockReleaseCommand releases a lock
func NewLockReleaseCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <stack-name>",
		Short: "Release a lock on stack state",
		Long:  `Releases an acquired lock to allow modifications.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			unlockedBy, _ := cmd.Flags().GetString("unlocked-by")

			if unlockedBy == "" {
				unlockedBy = "cli-user"
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

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Releasing lock for stack '%s'...", stackName))

			err = tracker.UnlockStateWithEvent(stack.ID, unlockedBy)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to release lock: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("Lock released for stack '%s'", stackName))

			return nil
		},
	}

	cmd.Flags().String("unlocked-by", "", "Unlock performer (default: cli-user)")

	return cmd
}

// NewLockStatusCommand shows lock status
func NewLockStatusCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <stack-name>",
		Short: "Show lock status for a stack",
		Long:  `Displays current lock status and details.`,
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

			isLocked, err := stackService.IsLocked(stack.ID)
			if err != nil {
				return err
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Lock Status: %s", stackName)
			fmt.Println()

			if !isLocked {
				pterm.Success.Println("✓ Stack is unlocked")
				return nil
			}

			pterm.Warning.Println("🔒 Stack is locked")
			fmt.Println()

			lockInfo, err := stackService.GetLockInfo(stack.ID)
			if err != nil {
				pterm.Warning.Printf("Could not retrieve lock details: %v\n", err)
			} else if lockInfo != nil {
				pterm.Info.Printf("Locked by: %s\n", lockInfo.Who)
				pterm.Info.Printf("Locked at: %s\n", lockInfo.CreatedAt.Format("2006-01-02 15:04:05"))
				pterm.Info.Printf("Operation: %s\n", lockInfo.Operation)
				pterm.Info.Printf("Expires: %s\n", lockInfo.ExpiresAt.Format("2006-01-02 15:04:05"))
			}

			return nil
		},
	}

	return cmd
}

// NewLockForceUnlockCommand force releases a lock
func NewLockForceUnlockCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "force-unlock <stack-name>",
		Short: "Force release a lock (use with caution)",
		Long:  `Forces the release of a lock. Use only when the lock holder is unavailable.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			force, _ := cmd.Flags().GetBool("force")

			if !force {
				pterm.Warning.Println("⚠️  WARNING: Force unlocking can cause data corruption if operations are in progress!")
				fmt.Println()

				result, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText(fmt.Sprintf("Are you ABSOLUTELY SURE you want to force unlock stack '%s'?", stackName)).
					Show()

				if !result {
					pterm.Info.Println("Force unlock cancelled")
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

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Force unlocking stack '%s'...", stackName))

			err = tracker.UnlockStateWithEvent(stack.ID, "force-unlock")
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to force unlock: %v", err))
				return err
			}

			spinner.Success(fmt.Sprintf("Stack '%s' force unlocked", stackName))
			pterm.Warning.Println("⚠ Verify stack integrity before proceeding")

			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}
