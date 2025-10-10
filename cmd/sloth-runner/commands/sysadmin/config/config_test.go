package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewConfigCmd(t *testing.T) {
	cmd := NewConfigCmd()

	if cmd == nil {
		t.Fatal("NewConfigCmd() returned nil")
	}

	if cmd.Use != "config" {
		t.Errorf("Expected Use='config', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}

	if cmd.Example == "" {
		t.Error("Example is empty")
	}
}

func TestConfigSubcommands(t *testing.T) {
	cmd := NewConfigCmd()

	expectedSubcommands := []string{
		"validate",
		"diff",
		"export",
		"import",
		"set",
		"get",
		"reset",
	}

	for _, expected := range expectedSubcommands {
		found := false
		for _, subcmd := range cmd.Commands() {
			if subcmd.Use == expected {
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

func TestConfigValidateCommand(t *testing.T) {
	cmd := NewConfigCmd()

	var validateCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "validate" {
			validateCmd = subcmd
			break
		}
	}

	if validateCmd == nil {
		t.Fatal("validate subcommand not found")
	}

	if validateCmd.Run == nil {
		t.Error("validate command has no Run function")
	}

	// Test execution
	if validateCmd.Run != nil {
		validateCmd.Run(validateCmd, []string{})
	}
}

func TestConfigDiffCommand(t *testing.T) {
	cmd := NewConfigCmd()

	var diffCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "diff" {
			diffCmd = subcmd
			break
		}
	}

	if diffCmd == nil {
		t.Fatal("diff subcommand not found")
	}

	if diffCmd.Run == nil {
		t.Error("diff command has no Run function")
	}

	// Test execution
	if diffCmd.Run != nil {
		diffCmd.Run(diffCmd, []string{})
	}
}

func TestConfigSetCommand(t *testing.T) {
	cmd := NewConfigCmd()

	var setCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "set" {
			setCmd = subcmd
			break
		}
	}

	if setCmd == nil {
		t.Fatal("set subcommand not found")
	}

	if setCmd.Run == nil {
		t.Error("set command has no Run function")
	}

	// Test execution
	if setCmd.Run != nil {
		setCmd.Run(setCmd, []string{})
	}
}

func TestConfigGetCommand(t *testing.T) {
	cmd := NewConfigCmd()

	var getCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "get" {
			getCmd = subcmd
			break
		}
	}

	if getCmd == nil {
		t.Fatal("get subcommand not found")
	}

	if getCmd.Run == nil {
		t.Error("get command has no Run function")
	}

	// Test execution
	if getCmd.Run != nil {
		getCmd.Run(getCmd, []string{})
	}
}

func TestConfigHelpText(t *testing.T) {
	cmd := NewConfigCmd()

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("help command failed: %v", err)
	}

	output := buf.String()

	requiredStrings := []string{
		"config",
		"validate",
		"diff",
		"export",
		"import",
		"set",
		"get",
		"reset",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(output, required) {
			t.Errorf("Help text missing '%s'", required)
		}
	}
}

func TestConfigDescriptionKeywords(t *testing.T) {
	cmd := NewConfigCmd()

	fullDescription := strings.ToLower(cmd.Short + " " + cmd.Long)

	requiredKeywords := []string{
		"config",
		"manage",
		"validate",
	}

	for _, keyword := range requiredKeywords {
		if !strings.Contains(fullDescription, keyword) {
			t.Errorf("Description missing keyword: %s", keyword)
		}
	}
}

func TestConfigCommandCount(t *testing.T) {
	cmd := NewConfigCmd()

	// Should have exactly 7 subcommands
	expectedCount := 7
	if len(cmd.Commands()) != expectedCount {
		t.Errorf("Expected %d subcommands, got %d", expectedCount, len(cmd.Commands()))
	}
}

func BenchmarkConfigCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewConfigCmd()
	}
}
