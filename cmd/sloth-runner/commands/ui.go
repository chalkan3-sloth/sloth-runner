package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/chalkan3-sloth/sloth-runner/internal/webui"
	"github.com/spf13/cobra"
)

// NewUICommand creates the UI command
func NewUICommand(ctx *AppContext) *cobra.Command {
	var (
		port       int
		debug      bool
		enableAuth bool
		username   string
		password   string
	)

	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Start the web-based UI dashboard",
		Long: `Starts a web-based dashboard for managing agents, workflows, hooks, and monitoring the sloth-runner system.

The UI provides:
  - Real-time dashboard with system statistics
  - Agent management and monitoring
  - Workflow (sloth) management
  - Hook configuration and execution history
  - Event queue monitoring
  - SSH profile management
  - Secrets overview (read-only)

All data persists in SQLite databases and updates in real-time via WebSockets.`,
		Example: `  # Start UI on default port 8080
  sloth-runner ui

  # Start UI on custom port
  sloth-runner ui --port 3000

  # Enable debug logging
  sloth-runner ui --debug

  # Enable basic authentication
  sloth-runner ui --auth --username admin --password secret`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUIServer(port, debug, enableAuth, username, password)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for the UI server")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug logging")
	cmd.Flags().BoolVar(&enableAuth, "auth", false, "Enable HTTP basic authentication")
	cmd.Flags().StringVar(&username, "username", "admin", "Username for basic auth")
	cmd.Flags().StringVar(&password, "password", "", "Password for basic auth (required if --auth is enabled)")

	return cmd
}

func runUIServer(port int, debug bool, enableAuth bool, username string, password string) error {
	// Validate auth config
	if enableAuth && password == "" {
		return fmt.Errorf("password is required when authentication is enabled")
	}

	// Get database paths
	agentDBPath := filepath.Join(".sloth-cache", "agents.db")
	slothDBPath := "/etc/sloth-runner/sloths.db"
	hookDBPath := filepath.Join(".sloth-cache", "hooks.db")

	homeDir, _ := os.UserHomeDir()
	secretsDBPath := filepath.Join(homeDir, ".sloth-runner", "secrets.db")
	sshDBPath := filepath.Join(homeDir, ".sloth-runner", "ssh_profiles.db")

	// Create server config
	cfg := &webui.Config{
		Port:          port,
		Debug:         debug,
		AgentDBPath:   agentDBPath,
		SlothDBPath:   slothDBPath,
		HookDBPath:    hookDBPath,
		SecretsDBPath: secretsDBPath,
		SSHDBPath:     sshDBPath,
		EnableAuth:    enableAuth,
		Username:      username,
		Password:      password,
	}

	// Create server
	server, err := webui.NewServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to create UI server: %w", err)
	}

	// Listen for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Printf("🚀 Sloth Runner UI starting on port %d", port)
		log.Printf("📊 Dashboard: http://localhost:%d", port)
		if enableAuth {
			log.Printf("🔐 Authentication: enabled (user: %s)", username)
		}

		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		log.Println("\n⚠️  Shutdown signal received")
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}

	// Graceful shutdown
	log.Println("🛑 Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("✅ Server stopped gracefully")
	return nil
}
