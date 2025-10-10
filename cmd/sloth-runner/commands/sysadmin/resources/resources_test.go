package resources

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewResourcesCmd(t *testing.T) {
	cmd := NewResourcesCmd()

	if cmd == nil {
		t.Fatal("NewResourcesCmd() returned nil")
	}

	if cmd.Use != "resources" {
		t.Errorf("Expected Use='resources', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}
}

func TestResourcesAliases(t *testing.T) {
	cmd := NewResourcesCmd()

	expectedAliases := []string{"resource", "res"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	for i, alias := range expectedAliases {
		if i >= len(cmd.Aliases) || cmd.Aliases[i] != alias {
			t.Errorf("Expected alias '%s' at position %d", alias, i)
		}
	}
}

func TestResourcesSubcommands(t *testing.T) {
	cmd := NewResourcesCmd()

	expectedSubcommands := []string{
		"overview",
		"cpu",
		"memory",
		"disk",
		"io",
		"network",
		"check",
		"history",
		"top",
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

func TestResourcesOverviewCommand(t *testing.T) {
	cmd := NewResourcesCmd()

	var overviewCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "overview" {
			overviewCmd = subcmd
			break
		}
	}

	if overviewCmd == nil {
		t.Fatal("overview subcommand not found")
	}

	if overviewCmd.Run == nil {
		t.Error("overview command has no Run function")
	}

	// Test execution
	if overviewCmd.Run != nil {
		overviewCmd.Run(overviewCmd, []string{})
	}
}

func TestResourcesCpuCommand(t *testing.T) {
	cmd := NewResourcesCmd()

	var cpuCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "cpu" {
			cpuCmd = subcmd
			break
		}
	}

	if cpuCmd == nil {
		t.Fatal("cpu subcommand not found")
	}

	if cpuCmd.Run == nil {
		t.Error("cpu command has no Run function")
	}

	// Test execution
	if cpuCmd.Run != nil {
		cpuCmd.Run(cpuCmd, []string{})
	}
}

func TestResourcesMemoryCommand(t *testing.T) {
	cmd := NewResourcesCmd()

	var memoryCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "memory" {
			memoryCmd = subcmd
			break
		}
	}

	if memoryCmd == nil {
		t.Fatal("memory subcommand not found")
	}

	if memoryCmd.Run == nil {
		t.Error("memory command has no Run function")
	}

	// Test execution
	if memoryCmd.Run != nil {
		memoryCmd.Run(memoryCmd, []string{})
	}
}

func TestResourcesDiskCommand(t *testing.T) {
	cmd := NewResourcesCmd()

	var diskCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "disk" {
			diskCmd = subcmd
			break
		}
	}

	if diskCmd == nil {
		t.Fatal("disk subcommand not found")
	}

	if diskCmd.Run == nil {
		t.Error("disk command has no Run function")
	}

	// Test execution
	if diskCmd.Run != nil {
		diskCmd.Run(diskCmd, []string{})
	}
}

func TestResourcesCheckCommand(t *testing.T) {
	cmd := NewResourcesCmd()

	var checkCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "check" {
			checkCmd = subcmd
			break
		}
	}

	if checkCmd == nil {
		t.Fatal("check subcommand not found")
	}

	if checkCmd.Run == nil {
		t.Error("check command has no Run function")
	}

	// Test execution
	if checkCmd.Run != nil {
		checkCmd.Run(checkCmd, []string{})
	}
}

func TestResourcesDescriptionKeywords(t *testing.T) {
	cmd := NewResourcesCmd()

	fullDescription := strings.ToLower(cmd.Short + " " + cmd.Long)

	requiredKeywords := []string{
		"resource",
		"cpu",
		"memory",
	}

	for _, keyword := range requiredKeywords {
		if !strings.Contains(fullDescription, keyword) {
			t.Errorf("Description missing keyword: %s", keyword)
		}
	}
}

func TestResourcesCommandCount(t *testing.T) {
	cmd := NewResourcesCmd()

	// Should have exactly 9 subcommands
	expectedCount := 9
	if len(cmd.Commands()) != expectedCount {
		t.Errorf("Expected %d subcommands, got %d", expectedCount, len(cmd.Commands()))
	}
}

func BenchmarkResourcesCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewResourcesCmd()
	}
}
