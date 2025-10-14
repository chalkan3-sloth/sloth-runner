//go:build cgo
// +build cgo

package secrets

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// NewSecretsCommand creates the secrets parent command
func NewSecretsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage encrypted secrets for stacks",
		Long: `Manage encrypted secrets for stacks using strong AES-256-GCM encryption.
Secrets are encrypted per-stack using a password-derived key (Argon2).

Each stack has its own encryption salt, ensuring that secrets cannot be
decrypted without both the correct password and stack context.`,
	}

	// Add subcommands
	cmd.AddCommand(NewAddCommand(ctx))
	cmd.AddCommand(NewListCommand(ctx))
	cmd.AddCommand(NewRemoveCommand(ctx))
	cmd.AddCommand(NewGetCommand(ctx))

	return cmd
}
