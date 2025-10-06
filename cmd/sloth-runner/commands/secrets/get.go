package secrets

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewGetCommand creates the get secret command
func NewGetCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <secret-name> --stack <stack-name>",
		Short: "Get and decrypt a secret value",
		Long: `Get and decrypt a secret value from a stack. Requires the encryption password.

WARNING: This will display the decrypted secret value in plain text.

Examples:
  # Get secret value (password prompt)
  sloth-runner secrets get api_key --stack my-app

  # Get secret value (password from stdin)
  echo 'mypassword' | sloth-runner secrets get api_key --stack my-app --password-stdin -`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName, _ := cmd.Flags().GetString("stack")
			passwordStdin, _ := cmd.Flags().GetBool("password-stdin")
			secretName := args[0]

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

			// Get encryption salt
			salt, err := services.GetOrCreateSalt(stackService, stack.ID)
			if err != nil {
				return fmt.Errorf("failed to get encryption salt: %w", err)
			}

			// Get password
			password, err := getDecryptPassword(passwordStdin)
			if err != nil {
				return fmt.Errorf("failed to get password: %w", err)
			}

			// Get secrets service
			secretsService, err := services.NewSecretsService()
			if err != nil {
				return fmt.Errorf("failed to create secrets service: %w", err)
			}
			defer secretsService.Close()

			// Get and decrypt secret
			value, err := secretsService.GetSecret(cmd.Context(), stack.ID, secretName, password, salt)
			if err != nil {
				return err
			}

			// Display decrypted value
			fmt.Println(value)

			return nil
		},
	}

	cmd.Flags().String("stack", "", "Stack name (required)")
	cmd.Flags().Bool("password-stdin", false, "Read password from stdin")
	cmd.MarkFlagRequired("stack")

	return cmd
}

// getDecryptPassword gets password for decryption
func getDecryptPassword(stdin bool) (string, error) {
	if stdin {
		reader := bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		return strings.TrimSpace(password), nil
	}

	// Prompt for password
	fmt.Print("Enter password to decrypt secret: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // New line after password input
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	password := string(passwordBytes)
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}
