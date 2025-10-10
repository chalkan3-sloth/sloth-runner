package security

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewSecurityCmd(t *testing.T) {
	cmd := NewSecurityCmd()
	if cmd == nil {
		t.Fatal("NewSecurityCmd() returned nil")
	}
	if cmd.Use != "security" {
		t.Errorf("Expected Use='security', got '%s'", cmd.Use)
	}
}

func TestSecuritySubcommands(t *testing.T) {
	cmd := NewSecurityCmd()
	expected := []string{"audit", "scan"}

	for _, exp := range expected {
		found := false
		for _, subcmd := range cmd.Commands() {
			if subcmd.Use == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing subcommand: %s", exp)
		}
	}
}

func TestSecurityAuditCommand(t *testing.T) {
	cmd := NewSecurityCmd()
	var auditCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "audit" {
			auditCmd = subcmd
			break
		}
	}

	if auditCmd == nil {
		t.Fatal("audit command not found")
	}

	// Test that command has a Run function
	if auditCmd.Run == nil {
		t.Error("audit command has no Run function")
	}

	// Test that command can be executed without panicking
	if auditCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		auditCmd.Run(auditCmd, []string{})
	}
}

func BenchmarkSecurityCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewSecurityCmd()
	}
}
