package secrets

import (
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewRemoveCommand creates the remove secret command
func NewRemoveCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <secret-name> --stack <stack-name>",
		Short: "Remove a secret from a stack",
		Long: `Remove an encrypted secret from a stack.

Examples:
  # Remove a specific secret
  sloth-runner secrets remove api_key --stack my-app

  # Remove all secrets from a stack
  sloth-runner secrets remove --stack my-app --all`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName, _ := cmd.Flags().GetString("stack")
			removeAll, _ := cmd.Flags().GetBool("all")

			if stackName == "" {
				return fmt.Errorf("--stack is required")
			}

			if !removeAll && len(args) == 0 {
				return fmt.Errorf("secret name is required (or use --all to remove all secrets)")
			}

			// Get stack service
			stackService, err := services.NewStackService()
			if err != nil {
				return fmt.Errorf("failed to create stack service: %w", err)
			}
			defer stackService.Close()

			// Verify stack exists
			stack, err := stackService.GetStack(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			// Get secrets service
			secretsService, err := services.NewSecretsService()
			if err != nil {
				return fmt.Errorf("failed to create secrets service: %w", err)
			}
			defer secretsService.Close()

			if removeAll {
				// Remove all secrets
				err = secretsService.RemoveAllSecrets(cmd.Context(), stack.ID)
				if err != nil {
					return err
				}
				pterm.Success.Printf("All secrets removed from stack '%s'\n", stackName)
			} else {
				// Remove specific secret
				secretName := args[0]
				err = secretsService.RemoveSecret(cmd.Context(), stack.ID, secretName)
				if err != nil {
					return err
				}
				pterm.Success.Printf("Secret '%s' removed from stack '%s'\n", secretName, stackName)
			}

			return nil
		},
	}

	cmd.Flags().String("stack", "", "Stack name (required)")
	cmd.Flags().Bool("all", false, "Remove all secrets from the stack")
	cmd.MarkFlagRequired("stack")

	return cmd
}
