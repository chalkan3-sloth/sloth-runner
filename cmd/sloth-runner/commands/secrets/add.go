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
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// NewAddCommand creates the add secret command
func NewAddCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <secret-name> --stack <stack-name>",
		Short: "Add an encrypted secret to a stack",
		Long: `Add an encrypted secret to a stack. The secret will be encrypted using AES-256-GCM
with a password-derived key (Argon2). The password is required to decrypt the secret later.

You can provide the secret value in three ways:
1. Interactively (prompted)
2. From a file using --from-file
3. From a YAML file using --from-yaml (for multiple secrets)

Examples:
  # Add secret interactively
  sloth-runner secrets add api_key --stack my-app

  # Add secret from file
  sloth-runner secrets add ssh_key --stack my-app --from-file ~/.ssh/id_rsa

  # Add multiple secrets from YAML
  sloth-runner secrets add --stack my-app --from-yaml secrets.yaml

  # YAML format (secrets.yaml):
  # secrets:
  #   api_key: "sk-abc123"
  #   db_password: "secret123"
  #   aws_access_key: "AKIAIOSFODNN7EXAMPLE"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName, _ := cmd.Flags().GetString("stack")
			fromFile, _ := cmd.Flags().GetString("from-file")
			fromYAML, _ := cmd.Flags().GetString("from-yaml")
			passwordStdin, _ := cmd.Flags().GetBool("password-stdin")

			// Validate arguments
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

			// Get or create encryption salt for stack
			salt, err := services.GetOrCreateSalt(stackService, stack.ID)
			if err != nil {
				return fmt.Errorf("failed to get encryption salt: %w", err)
			}

			// Get password
			password, err := getPassword(passwordStdin)
			if err != nil {
				return fmt.Errorf("failed to get password: %w", err)
			}

			// Get secrets service
			secretsService, err := services.NewSecretsService()
			if err != nil {
				return fmt.Errorf("failed to create secrets service: %w", err)
			}
			defer secretsService.Close()

			// Handle different input methods
			if fromYAML != "" {
				return addSecretsFromYAML(cmd, secretsService, stack.ID, fromYAML, password, salt)
			} else if fromFile != "" {
				if len(args) == 0 {
					return fmt.Errorf("secret name is required when using --from-file")
				}
				return addSecretFromFile(cmd, secretsService, stack.ID, args[0], fromFile, password, salt)
			} else {
				if len(args) == 0 {
					return fmt.Errorf("secret name is required")
				}
				return addSecretInteractive(cmd, secretsService, stack.ID, args[0], password, salt)
			}
		},
	}

	cmd.Flags().String("stack", "", "Stack name (required)")
	cmd.Flags().String("from-file", "", "Read secret value from file")
	cmd.Flags().String("from-yaml", "", "Read multiple secrets from YAML file")
	cmd.Flags().Bool("password-stdin", false, "Read password from stdin")
	cmd.MarkFlagRequired("stack")

	return cmd
}

// getPassword gets password from stdin or prompts user
func getPassword(stdin bool) (string, error) {
	if stdin {
		reader := bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		return strings.TrimSpace(password), nil
	}

	// Prompt for password
	fmt.Print("Enter password to encrypt secrets: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // New line after password input
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	password := string(passwordBytes)
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	// Confirm password
	fmt.Print("Confirm password: ")
	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	if password != string(confirmBytes) {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}

// addSecretInteractive adds a secret interactively
func addSecretInteractive(cmd *cobra.Command, service *services.SecretsService, stackID, name, password string, salt []byte) error {
	fmt.Printf("Enter value for secret '%s': ", name)
	valueBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("failed to read secret value: %w", err)
	}

	value := string(valueBytes)
	if value == "" {
		return fmt.Errorf("secret value cannot be empty")
	}

	// Add secret
	err = service.AddSecret(cmd.Context(), stackID, name, value, password, salt)
	if err != nil {
		return err
	}

	pterm.Success.Printf("Secret '%s' added successfully\n", name)
	return nil
}

// addSecretFromFile adds a secret from a file
func addSecretFromFile(cmd *cobra.Command, service *services.SecretsService, stackID, name, filePath, password string, salt []byte) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	value := string(data)
	if value == "" {
		return fmt.Errorf("file is empty")
	}

	// Add secret
	err = service.AddSecret(cmd.Context(), stackID, name, value, password, salt)
	if err != nil {
		return err
	}

	pterm.Success.Printf("Secret '%s' added from file '%s'\n", name, filePath)
	return nil
}

// addSecretsFromYAML adds multiple secrets from a YAML file
func addSecretsFromYAML(cmd *cobra.Command, service *services.SecretsService, stackID, filePath, password string, salt []byte) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	var config struct {
		Secrets map[string]string `yaml:"secrets"`
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	if len(config.Secrets) == 0 {
		return fmt.Errorf("no secrets found in YAML file (expected 'secrets:' key)")
	}

	// Add each secret
	count := 0
	for name, value := range config.Secrets {
		err = service.AddSecret(cmd.Context(), stackID, name, value, password, salt)
		if err != nil {
			pterm.Error.Printf("Failed to add secret '%s': %v\n", name, err)
			continue
		}
		pterm.Success.Printf("Added secret '%s'\n", name)
		count++
	}

	pterm.Info.Printf("Added %d/%d secrets from '%s'\n", count, len(config.Secrets), filePath)
	return nil
}
