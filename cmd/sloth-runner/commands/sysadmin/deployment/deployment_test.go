package deployment

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewDeploymentCmd(t *testing.T) {
	cmd := NewDeploymentCmd()
	if cmd == nil {
		t.Fatal("NewDeploymentCmd() returned nil")
	}
	if cmd.Use != "deployment" {
		t.Errorf("Expected Use='deployment', got '%s'", cmd.Use)
	}

	// Test alias
	hasDeployAlias := false
	for _, alias := range cmd.Aliases {
		if alias == "deploy" {
			hasDeployAlias = true
			break
		}
	}
	if !hasDeployAlias {
		t.Error("Missing 'deploy' alias")
	}
}

func TestDeploymentSubcommands(t *testing.T) {
	cmd := NewDeploymentCmd()
	expected := []string{"deploy", "rollback"}

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

func TestDeploymentDeployCommand(t *testing.T) {
	cmd := NewDeploymentCmd()
	var deployCmd *cobra.Command
	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == "deploy" {
			deployCmd = subcmd
			break
		}
	}

	if deployCmd == nil {
		t.Fatal("deploy command not found")
	}

	// Test that command has a Run function
	if deployCmd.Run == nil {
		t.Error("deploy command has no Run function")
	}

	// Test that command can be executed without panicking
	if deployCmd.Run != nil {
		// Just verify it doesn't panic - pterm writes to os.Stdout directly
		deployCmd.Run(deployCmd, []string{})
	}
}

func BenchmarkDeploymentCmdCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewDeploymentCmd()
	}
}
