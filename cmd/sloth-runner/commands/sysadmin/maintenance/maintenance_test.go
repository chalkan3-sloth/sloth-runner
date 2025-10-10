package maintenance

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewMaintenanceCmd(t *testing.T) {
	cmd := NewMaintenanceCmd()
	if cmd == nil {
		t.Fatal("NewMaintenanceCmd() returned nil")
	}
	if cmd.Use != "maintenance" {
		t.Errorf("Expected Use='maintenance', got '%s'", cmd.Use)
	}
}

func TestMaintenanceSubcommands(t *testing.T) {
	cmd := NewMaintenanceCmd()
	expected := []string{"clean-logs", "optimize-db", "cleanup"}

	for _, exp := range expected {
		found := false
		for _, subcmd := range cmd.Commands() {
			if subcmd.Use == exp {
				found = true
				if subcmd.Short == "" {
					t.Errorf("Subcommand %s has no description", exp)
				}
				break
			}
		}
		if !found {
			t.Errorf("Missing subcommand: %s", exp)
		}
	}
}

func TestMaintenanceCleanLogsCommand(t *testing.T) {
	cmd := NewMaintenanceCmd()
	var cleanCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "clean-logs" {
			cleanCmd = subcmd
			break
		}
	}

	if cleanCmd == nil {
		t.Fatal("clean-logs command not found")
	}

	// Test that command has a Run function
	if cleanCmd.Run == nil {
		t.Error("clean-logs command has no Run function")
	}

	// Test that command can be executed without panicking
	if cleanCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		cleanCmd.Run(cleanCmd, []string{})
	}
}

func BenchmarkMaintenanceCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewMaintenanceCmd()
	}
}
