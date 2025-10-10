package backup

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// stripANSI removes ANSI color codes from a string
func stripANSI(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(s, "")
}

func TestNewBackupCmd(t *testing.T) {
	cmd := NewBackupCmd()

	if cmd == nil {
		t.Fatal("NewBackupCmd() returned nil")
	}

	if cmd.Use != "backup" {
		t.Errorf("Expected Use='backup', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}
}

func TestBackupSubcommands(t *testing.T) {
	cmd := NewBackupCmd()

	expectedSubcommands := []string{"create", "restore"}

	for _, expected := range expectedSubcommands {
		found := false
		for _, subcmd := range cmd.Commands() {
			if subcmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

func TestBackupCreateCommand(t *testing.T) {
	cmd := NewBackupCmd()

	var createCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "create" {
			createCmd = subcmd
			break
		}
	}

	if createCmd == nil {
		t.Fatal("create subcommand not found")
	}

	// Test that command has a Run function
	if createCmd.Run == nil {
		t.Error("create command has no Run function")
	}

	// Test that command can be executed without panicking
	if createCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		createCmd.Run(createCmd, []string{})
	}
}

func TestBackupRestoreCommand(t *testing.T) {
	cmd := NewBackupCmd()

	var restoreCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "restore" {
			restoreCmd = subcmd
			break
		}
	}

	if restoreCmd == nil {
		t.Fatal("restore subcommand not found")
	}

	// Test that command has a Run function
	if restoreCmd.Run == nil {
		t.Error("restore command has no Run function")
	}

	// Test that command can be executed without panicking
	if restoreCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		restoreCmd.Run(restoreCmd, []string{})
	}
}

func TestBackupHelpText(t *testing.T) {
	cmd := NewBackupCmd()

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("help command failed: %v", err)
	}

	output := buf.String()

	requiredStrings := []string{
		"backup",
		"create",
		"restore",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(output, required) {
			t.Errorf("Help text missing '%s'", required)
		}
	}
}

func BenchmarkBackupCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewBackupCmd()
	}
}

func BenchmarkBackupCreateExecution(b *testing.B) {
	cmd := NewBackupCmd()
	var createCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "create" {
			createCmd = subcmd
			break
		}
	}

	buf := new(bytes.Buffer)
	createCmd.SetOut(buf)
	createCmd.SetErr(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = createCmd.Execute()
	}
}
