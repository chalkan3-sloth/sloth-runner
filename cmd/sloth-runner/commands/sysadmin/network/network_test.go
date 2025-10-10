package network

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewNetworkCmd(t *testing.T) {
	cmd := NewNetworkCmd()

	if cmd == nil {
		t.Fatal("NewNetworkCmd() returned nil")
	}

	if cmd.Use != "network" {
		t.Errorf("Expected Use='network', got '%s'", cmd.Use)
	}

	// Test aliases
	expectedAliases := []string{"net"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}
}

func TestNetworkSubcommands(t *testing.T) {
	cmd := NewNetworkCmd()

	tests := []struct {
		name     string
		expected string
	}{
		{"ping command", "ping"},
		{"port-check command", "port-check"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, subcmd := range cmd.Commands() {
				if subcmd.Use == tt.expected {
					found = true

					// Verify Short description exists
					if subcmd.Short == "" {
						t.Errorf("Subcommand '%s' has empty Short description", tt.expected)
					}

					break
				}
			}
			if !found {
				t.Errorf("Expected subcommand '%s' not found", tt.expected)
			}
		})
	}
}

func TestNetworkPingCommand(t *testing.T) {
	cmd := NewNetworkCmd()

	var pingCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "ping" {
			pingCmd = subcmd
			break
		}
	}

	if pingCmd == nil {
		t.Fatal("ping subcommand not found")
	}

	// Test that command has a Run function
	if pingCmd.Run == nil {
		t.Error("ping command has no Run function")
	}

	// Test that command can be executed without panicking
	if pingCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		pingCmd.Run(pingCmd, []string{})
	}
}

func TestNetworkPortCheckCommand(t *testing.T) {
	cmd := NewNetworkCmd()

	var portCheckCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "port-check" {
			portCheckCmd = subcmd
			break
		}
	}

	if portCheckCmd == nil {
		t.Fatal("port-check subcommand not found")
	}

	// Test that command has a Run function
	if portCheckCmd.Run == nil {
		t.Error("port-check command has no Run function")
	}

	// Test that command can be executed without panicking
	if portCheckCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		portCheckCmd.Run(portCheckCmd, []string{})
	}
}

func TestNetworkCommandDescriptions(t *testing.T) {
	cmd := NewNetworkCmd()

	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Long description is empty")
	}

	// Check for key networking terms
	fullText := strings.ToLower(cmd.Short + " " + cmd.Long)
	keywords := []string{"network", "diagnostic", "connectivity", "latency"}

	for _, keyword := range keywords {
		if !strings.Contains(fullText, keyword) {
			t.Errorf("Description missing keyword: %s", keyword)
		}
	}
}

func TestNetworkCommandCount(t *testing.T) {
	cmd := NewNetworkCmd()

	expectedCount := 2 // ping and port-check
	actualCount := len(cmd.Commands())

	if actualCount != expectedCount {
		t.Errorf("Expected %d subcommands, got %d", expectedCount, actualCount)
	}
}

func BenchmarkNetworkCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewNetworkCmd()
	}
}

func BenchmarkNetworkPingExecution(b *testing.B) {
	cmd := NewNetworkCmd()
	var pingCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "ping" {
			pingCmd = subcmd
			break
		}
	}

	buf := new(bytes.Buffer)
	pingCmd.SetOut(buf)
	pingCmd.SetErr(buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pingCmd.Execute()
	}
}
