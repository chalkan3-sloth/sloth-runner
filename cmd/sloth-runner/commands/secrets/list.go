//go:build cgo
// +build cgo

package secrets

import (
	"fmt"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/spf13/cobra"
)

// NewListCommand creates the list secrets command
func NewListCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list --stack <stack-name>",
		Short: "List all secrets for a stack",
		Long: `List all secrets for a stack. Only secret names are shown, values remain encrypted.

Example:
  sloth-runner secrets list --stack my-app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName, _ := cmd.Flags().GetString("stack")

			if stackName == "" {
				return fmt.Errorf("--stack is required")
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

			// List secrets
			secrets, err := secretsService.ListSecrets(cmd.Context(), stack.ID)
			if err != nil {
				return err
			}

			if len(secrets) == 0 {
				fmt.Printf("No secrets found for stack '%s'\n", stackName)
				return nil
			}

			// Display table
			writer := cmd.OutOrStdout()
			tw := tabwriter.NewWriter(writer, 0, 0, 3, ' ', 0)
			fmt.Fprintln(tw, "NAME\tCREATED\tUPDATED")
			fmt.Fprintln(tw, "----\t-------\t-------")

			for _, secret := range secrets {
				fmt.Fprintf(tw, "%s\t%s\t%s\n",
					secret.Name,
					secret.CreatedAt.Format("2006-01-02 15:04:05"),
					secret.UpdatedAt.Format("2006-01-02 15:04:05"),
				)
			}

			return tw.Flush()
		},
	}

	cmd.Flags().String("stack", "", "Stack name (required)")
	cmd.MarkFlagRequired("stack")

	return cmd
}
