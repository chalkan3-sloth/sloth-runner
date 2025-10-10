package sysadmin

import (
	"strings"
	"testing"
)

func TestNewSysadminCmd(t *testing.T) {
	cmd := NewSysadminCmd()

	if cmd == nil {
		t.Fatal("NewSysadminCmd() returned nil")
	}

	if cmd.Use != "sysadmin" {
		t.Errorf("Expected Use='sysadmin', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}
}

func TestSysadminSubcommands(t *testing.T) {
	cmd := NewSysadminCmd()

	expectedSubcommands := []string{
		"logs",
		"health",
		"debug",
		"backup",
		"config",
		"deployment",
		"maintenance",
		"network",
		"packages",
		"performance",
		"resources",
		"security",
		"services",
	}

	for _, expected := range expectedSubcommands {
		found := false
		for _, subcmd := range cmd.Commands() {
			if subcmd.Use == expected || strings.HasPrefix(subcmd.Use, expected+" ") {
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

func TestSysadminCommandCount(t *testing.T) {
	cmd := NewSysadminCmd()

	// Should have at least 13 subcommands (excluding built-in like completion, help)
	if len(cmd.Commands()) < 13 {
		t.Errorf("Expected at least 13 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestSysadminDescriptionKeywords(t *testing.T) {
	cmd := NewSysadminCmd()

	fullDescription := strings.ToLower(cmd.Short + " " + cmd.Long)

	requiredKeywords := []string{
		"system",
		"administration",
		"monitor",
		"backup",
		"performance",
		"network",
	}

	for _, keyword := range requiredKeywords {
		if !strings.Contains(fullDescription, keyword) {
			t.Errorf("Description missing keyword: %s", keyword)
		}
	}
}

func TestSysadminExampleText(t *testing.T) {
	cmd := NewSysadminCmd()

	if cmd.Example == "" {
		t.Error("No example text provided")
	}

	// Check for key examples
	examples := []string{
		"logs",
		"health",
		"backup",
		"performance",
		"network",
	}

	for _, example := range examples {
		if !strings.Contains(cmd.Example, example) {
			t.Errorf("Example missing command: %s", example)
		}
	}
}

func BenchmarkSysadminCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewSysadminCmd()
	}
}
