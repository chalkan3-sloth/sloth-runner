package performance

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewPerformanceCmd(t *testing.T) {
	cmd := NewPerformanceCmd()

	if cmd == nil {
		t.Fatal("NewPerformanceCmd() returned nil")
	}

	if cmd.Use != "performance" {
		t.Errorf("Expected Use='performance', got '%s'", cmd.Use)
	}

	// Test aliases
	if len(cmd.Aliases) == 0 {
		t.Error("No aliases defined")
	}

	hasPerf := false
	for _, alias := range cmd.Aliases {
		if alias == "perf" {
			hasPerf = true
			break
		}
	}
	if !hasPerf {
		t.Error("Missing 'perf' alias")
	}
}

func TestPerformanceSubcommands(t *testing.T) {
	cmd := NewPerformanceCmd()

	tests := []struct {
		name     string
		expected string
	}{
		{"show command", "show"},
		{"monitor command", "monitor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, subcmd := range cmd.Commands() {
				if subcmd.Use == tt.expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected subcommand '%s' not found", tt.expected)
			}
		})
	}
}

func TestPerformanceShowCommand(t *testing.T) {
	cmd := NewPerformanceCmd()

	var showCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "show" {
			showCmd = subcmd
			break
		}
	}

	if showCmd == nil {
		t.Fatal("show subcommand not found")
	}

	// Test that command has a Run function
	if showCmd.Run == nil {
		t.Error("show command has no Run function")
	}

	// Test that command can be executed without panicking
	if showCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		showCmd.Run(showCmd, []string{})
	}
}

func TestPerformanceMonitorCommand(t *testing.T) {
	cmd := NewPerformanceCmd()

	var monitorCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "monitor" {
			monitorCmd = subcmd
			break
		}
	}

	if monitorCmd == nil {
		t.Fatal("monitor subcommand not found")
	}

	// Test that command has a Run function
	if monitorCmd.Run == nil {
		t.Error("monitor command has no Run function")
	}

	// Test that command can be executed without panicking
	if monitorCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		monitorCmd.Run(monitorCmd, []string{})
	}
}

func TestPerformanceDescriptions(t *testing.T) {
	cmd := NewPerformanceCmd()

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}

	requiredKeywords := []string{"performance", "monitor", "CPU", "memory"}
	description := cmd.Short + " " + cmd.Long

	for _, keyword := range requiredKeywords {
		if !strings.Contains(strings.ToLower(description), strings.ToLower(keyword)) {
			t.Errorf("Description missing keyword: %s", keyword)
		}
	}
}

func BenchmarkPerformanceCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPerformanceCmd()
	}
}
