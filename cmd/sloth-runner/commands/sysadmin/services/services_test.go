package services

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewServicesCmd(t *testing.T) {
	cmd := NewServicesCmd()

	if cmd == nil {
		t.Fatal("NewServicesCmd() returned nil")
	}

	if cmd.Use != "services" {
		t.Errorf("Expected Use='services', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}
}

func TestServicesAliases(t *testing.T) {
	cmd := NewServicesCmd()

	expectedAliases := []string{"service", "svc"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	for i, alias := range expectedAliases {
		if i >= len(cmd.Aliases) || cmd.Aliases[i] != alias {
			t.Errorf("Expected alias '%s' at position %d", alias, i)
		}
	}
}

func TestServicesSubcommands(t *testing.T) {
	cmd := NewServicesCmd()

	expectedSubcommands := []string{
		"list",
		"status",
		"start",
		"stop",
		"restart",
		"reload",
		"enable",
		"disable",
		"logs",
	}

	for _, expected := range expectedSubcommands {
		found := false
		for _, subcmd := range cmd.Commands() {
			if strings.HasPrefix(subcmd.Use, expected) {
				found = true

				// Verify each subcommand has description
				if subcmd.Short == "" {
					t.Errorf("Subcommand '%s' has no Short description", expected)
				}

				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

func TestServicesListCommand(t *testing.T) {
	cmd := NewServicesCmd()

	var listCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "list" {
			listCmd = subcmd
			break
		}
	}

	if listCmd == nil {
		t.Fatal("list subcommand not found")
	}

	if listCmd.Run == nil {
		t.Error("list command has no Run function")
	}

	// Test execution
	if listCmd.Run != nil {
		listCmd.Run(listCmd, []string{})
	}
}

func TestServicesStatusCommand(t *testing.T) {
	cmd := NewServicesCmd()

	var statusCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "status") {
			statusCmd = subcmd
			break
		}
	}

	if statusCmd == nil {
		t.Fatal("status subcommand not found")
	}

	if statusCmd.Run == nil {
		t.Error("status command has no Run function")
	}

	// Test execution
	if statusCmd.Run != nil {
		statusCmd.Run(statusCmd, []string{})
	}
}

func TestServicesStartCommand(t *testing.T) {
	cmd := NewServicesCmd()

	var startCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "start") {
			startCmd = subcmd
			break
		}
	}

	if startCmd == nil {
		t.Fatal("start subcommand not found")
	}

	if startCmd.Run == nil {
		t.Error("start command has no Run function")
	}

	// Test execution
	if startCmd.Run != nil {
		startCmd.Run(startCmd, []string{})
	}
}

func TestServicesRestartCommand(t *testing.T) {
	cmd := NewServicesCmd()

	var restartCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "restart") {
			restartCmd = subcmd
			break
		}
	}

	if restartCmd == nil {
		t.Fatal("restart subcommand not found")
	}

	if restartCmd.Run == nil {
		t.Error("restart command has no Run function")
	}

	// Test execution
	if restartCmd.Run != nil {
		restartCmd.Run(restartCmd, []string{})
	}
}

func TestServicesDescriptionKeywords(t *testing.T) {
	cmd := NewServicesCmd()

	fullDescription := strings.ToLower(cmd.Short + " " + cmd.Long)

	requiredKeywords := []string{
		"service",
		"systemd",
		"manage",
	}

	for _, keyword := range requiredKeywords {
		if !strings.Contains(fullDescription, keyword) {
			t.Errorf("Description missing keyword: %s", keyword)
		}
	}
}

func TestServicesCommandCount(t *testing.T) {
	cmd := NewServicesCmd()

	// Should have exactly 9 subcommands
	expectedCount := 9
	if len(cmd.Commands()) != expectedCount {
		t.Errorf("Expected %d subcommands, got %d", expectedCount, len(cmd.Commands()))
	}
}

func BenchmarkServicesCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewServicesCmd()
	}
}
