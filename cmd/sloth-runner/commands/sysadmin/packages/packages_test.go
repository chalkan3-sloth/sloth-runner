package packages

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewPackagesCmd(t *testing.T) {
	cmd := NewPackagesCmd()

	if cmd == nil {
		t.Fatal("NewPackagesCmd() returned nil")
	}

	if cmd.Use != "packages" {
		t.Errorf("Expected Use='packages', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}
}

func TestPackagesAliases(t *testing.T) {
	cmd := NewPackagesCmd()

	expectedAliases := []string{"package", "pkg"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	for i, alias := range expectedAliases {
		if i >= len(cmd.Aliases) || cmd.Aliases[i] != alias {
			t.Errorf("Expected alias '%s' at position %d", alias, i)
		}
	}
}

func TestPackagesSubcommands(t *testing.T) {
	cmd := NewPackagesCmd()

	expectedSubcommands := []string{
		"list",
		"search",
		"install",
		"remove",
		"update",
		"upgrade",
		"check-updates",
		"info",
		"history",
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

func TestPackagesListCommand(t *testing.T) {
	cmd := NewPackagesCmd()

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

func TestPackagesInstallCommand(t *testing.T) {
	cmd := NewPackagesCmd()

	var installCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "install") {
			installCmd = subcmd
			break
		}
	}

	if installCmd == nil {
		t.Fatal("install subcommand not found")
	}

	if installCmd.Run == nil {
		t.Error("install command has no Run function")
	}

	// Test execution
	if installCmd.Run != nil {
		installCmd.Run(installCmd, []string{})
	}
}

func TestPackagesUpgradeCommand(t *testing.T) {
	cmd := NewPackagesCmd()

	var upgradeCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if strings.HasPrefix(subcmd.Use, "upgrade") {
			upgradeCmd = subcmd
			break
		}
	}

	if upgradeCmd == nil {
		t.Fatal("upgrade subcommand not found")
	}

	if upgradeCmd.Run == nil {
		t.Error("upgrade command has no Run function")
	}

	// Test execution
	if upgradeCmd.Run != nil {
		upgradeCmd.Run(upgradeCmd, []string{})
	}
}

func TestPackagesDescriptionKeywords(t *testing.T) {
	cmd := NewPackagesCmd()

	fullDescription := strings.ToLower(cmd.Short + " " + cmd.Long)

	requiredKeywords := []string{
		"package",
		"install",
		"update",
	}

	for _, keyword := range requiredKeywords {
		if !strings.Contains(fullDescription, keyword) {
			t.Errorf("Description missing keyword: %s", keyword)
		}
	}
}

func TestPackagesCommandCount(t *testing.T) {
	cmd := NewPackagesCmd()

	// Should have exactly 9 subcommands
	expectedCount := 9
	if len(cmd.Commands()) != expectedCount {
		t.Errorf("Expected %d subcommands, got %d", expectedCount, len(cmd.Commands()))
	}
}

func BenchmarkPackagesCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPackagesCmd()
	}
}
