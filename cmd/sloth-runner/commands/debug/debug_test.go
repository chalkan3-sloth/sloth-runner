//go:build cgo
// +build cgo

package debug

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestNewDebugCmd tests debug command creation
func TestNewDebugCmd(t *testing.T) {
	cmd := NewDebugCmd()

	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "debug" {
		t.Errorf("Expected Use 'debug', got %s", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("Expected Short description")
	}
	if cmd.Long == "" {
		t.Error("Expected Long description")
	}
}

// TestNewDebugCmd_HasExamples tests debug command has examples
func TestNewDebugCmd_HasExamples(t *testing.T) {
	cmd := NewDebugCmd()

	if cmd.Example == "" {
		t.Error("Expected Example field")
	}
}

// TestNewDebugCmd_Subcommands tests debug command has subcommands
func TestNewDebugCmd_Subcommands(t *testing.T) {
	cmd := NewDebugCmd()

	subcommands := cmd.Commands()
	if len(subcommands) != 3 {
		t.Errorf("Expected 3 subcommands, got %d", len(subcommands))
	}

	expectedSubcommands := map[string]bool{
		"connection": false,
		"agent":      false,
		"workflow":   false,
	}

	for _, subcmd := range subcommands {
		if _, ok := expectedSubcommands[subcmd.Name()]; ok {
			expectedSubcommands[subcmd.Name()] = true
		}
	}

	for name, found := range expectedSubcommands {
		if !found {
			t.Errorf("Expected subcommand '%s' not found", name)
		}
	}
}

// TestNewConnectionCmd tests connection command creation
func TestNewConnectionCmd(t *testing.T) {
	cmd := newConnectionCmd()

	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "connection [agent-name]" {
		t.Errorf("Expected Use 'connection [agent-name]', got %s", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("Expected Short description")
	}
	if cmd.Long == "" {
		t.Error("Expected Long description")
	}
}

// TestNewConnectionCmd_Flags tests connection command flags
func TestNewConnectionCmd_Flags(t *testing.T) {
	cmd := newConnectionCmd()

	timeoutFlag := cmd.Flags().Lookup("timeout")
	if timeoutFlag == nil {
		t.Error("Expected timeout flag")
	}
	if timeoutFlag.DefValue != "5" {
		t.Errorf("Expected default timeout '5', got %s", timeoutFlag.DefValue)
	}

	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("Expected verbose flag")
	}
	if verboseFlag.DefValue != "false" {
		t.Errorf("Expected default verbose 'false', got %s", verboseFlag.DefValue)
	}
}

// TestNewConnectionCmd_FlagShorthands tests connection command flag shorthands
func TestNewConnectionCmd_FlagShorthands(t *testing.T) {
	cmd := newConnectionCmd()

	expectedShorthands := map[string]string{
		"timeout": "t",
		"verbose": "v",
	}

	for flagName, expectedShorthand := range expectedShorthands {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' not found", flagName)
			continue
		}
		if flag.Shorthand != expectedShorthand {
			t.Errorf("Flag %s: expected shorthand '%s', got '%s'", flagName, expectedShorthand, flag.Shorthand)
		}
	}
}

// TestNewConnectionCmd_Args tests connection command arguments
func TestNewConnectionCmd_Args(t *testing.T) {
	cmd := newConnectionCmd()

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error for no arguments, got nil")
	}

	err = cmd.Args(cmd, []string{"agent1"})
	if err != nil {
		t.Errorf("Expected no error for one argument, got %v", err)
	}

	err = cmd.Args(cmd, []string{"agent1", "agent2"})
	if err == nil {
		t.Error("Expected error for two arguments, got nil")
	}
}

// TestNewAgentCmd tests agent command creation
func TestNewAgentCmd(t *testing.T) {
	cmd := newAgentCmd()

	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "agent [agent-name]" {
		t.Errorf("Expected Use 'agent [agent-name]', got %s", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("Expected Short description")
	}
	if cmd.Long == "" {
		t.Error("Expected Long description")
	}
}

// TestNewAgentCmd_Flags tests agent command flags
func TestNewAgentCmd_Flags(t *testing.T) {
	cmd := newAgentCmd()

	fullFlag := cmd.Flags().Lookup("full")
	if fullFlag == nil {
		t.Error("Expected full flag")
	}
	if fullFlag.DefValue != "false" {
		t.Errorf("Expected default full 'false', got %s", fullFlag.DefValue)
	}
}

// TestNewAgentCmd_FlagShorthand tests agent command flag shorthand
func TestNewAgentCmd_FlagShorthand(t *testing.T) {
	cmd := newAgentCmd()

	fullFlag := cmd.Flags().Lookup("full")
	if fullFlag == nil {
		t.Fatal("Expected full flag")
	}

	if fullFlag.Shorthand != "f" {
		t.Errorf("Expected shorthand 'f', got %s", fullFlag.Shorthand)
	}
}

// TestNewAgentCmd_Args tests agent command arguments
func TestNewAgentCmd_Args(t *testing.T) {
	cmd := newAgentCmd()

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error for no arguments, got nil")
	}

	err = cmd.Args(cmd, []string{"agent1"})
	if err != nil {
		t.Errorf("Expected no error for one argument, got %v", err)
	}

	err = cmd.Args(cmd, []string{"agent1", "agent2"})
	if err == nil {
		t.Error("Expected error for two arguments, got nil")
	}
}

// TestNewWorkflowCmd tests workflow command creation
func TestNewWorkflowCmd(t *testing.T) {
	cmd := newWorkflowCmd()

	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "workflow [workflow-name|latest]" {
		t.Errorf("Expected Use 'workflow [workflow-name|latest]', got %s", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("Expected Short description")
	}
	if cmd.Long == "" {
		t.Error("Expected Long description")
	}
}

// TestNewWorkflowCmd_Flags tests workflow command flags
func TestNewWorkflowCmd_Flags(t *testing.T) {
	cmd := newWorkflowCmd()

	lastFlag := cmd.Flags().Lookup("last")
	if lastFlag == nil {
		t.Error("Expected last flag")
	}
	if lastFlag.DefValue != "1" {
		t.Errorf("Expected default last '1', got %s", lastFlag.DefValue)
	}
}

// TestNewWorkflowCmd_FlagShorthand tests workflow command flag shorthand
func TestNewWorkflowCmd_FlagShorthand(t *testing.T) {
	cmd := newWorkflowCmd()

	lastFlag := cmd.Flags().Lookup("last")
	if lastFlag == nil {
		t.Fatal("Expected last flag")
	}

	if lastFlag.Shorthand != "n" {
		t.Errorf("Expected shorthand 'n', got %s", lastFlag.Shorthand)
	}
}

// TestNewWorkflowCmd_Args tests workflow command arguments
func TestNewWorkflowCmd_Args(t *testing.T) {
	cmd := newWorkflowCmd()

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error for no arguments, got nil")
	}

	err = cmd.Args(cmd, []string{"latest"})
	if err != nil {
		t.Errorf("Expected no error for 'latest' argument, got %v", err)
	}

	err = cmd.Args(cmd, []string{"workflow-name"})
	if err != nil {
		t.Errorf("Expected no error for workflow name argument, got %v", err)
	}

	err = cmd.Args(cmd, []string{"workflow1", "workflow2"})
	if err == nil {
		t.Error("Expected error for two arguments, got nil")
	}
}

// TestCommands_ConsistentInterface tests all commands have consistent interface
func TestCommands_ConsistentInterface(t *testing.T) {
	cmds := map[string]*cobra.Command{
		"debug":      NewDebugCmd(),
		"connection": newConnectionCmd(),
		"agent":      newAgentCmd(),
		"workflow":   newWorkflowCmd(),
	}

	for name, cmd := range cmds {
		t.Run(name, func(t *testing.T) {
			if cmd.Use == "" {
				t.Errorf("Command %s missing Use field", name)
			}
			if cmd.Short == "" {
				t.Errorf("Command %s missing Short description", name)
			}
			// Only non-parent commands need RunE
			if name != "debug" && cmd.RunE == nil {
				t.Errorf("Command %s missing RunE function", name)
			}
		})
	}
}

// TestCommandDescriptions_AreMeaningful tests command descriptions
func TestCommandDescriptions_AreMeaningful(t *testing.T) {
	commands := map[string]*cobra.Command{
		"debug":      NewDebugCmd(),
		"connection": newConnectionCmd(),
		"agent":      newAgentCmd(),
		"workflow":   newWorkflowCmd(),
	}

	for name, cmd := range commands {
		t.Run(name, func(t *testing.T) {
			if len(cmd.Short) < 10 {
				t.Errorf("Command %s has too short description: %s", name, cmd.Short)
			}
		})
	}
}

// TestConnectionCmd_HasExamples tests that connection command has examples
func TestConnectionCmd_HasExamples(t *testing.T) {
	cmd := newConnectionCmd()

	if cmd.Example == "" {
		t.Error("Expected Example field")
	}
}

// TestAgentCmd_HasExamples tests that agent command has examples
func TestAgentCmd_HasExamples(t *testing.T) {
	cmd := newAgentCmd()

	if cmd.Example == "" {
		t.Error("Expected Example field")
	}
}

// TestWorkflowCmd_HasExamples tests that workflow command has examples
func TestWorkflowCmd_HasExamples(t *testing.T) {
	cmd := newWorkflowCmd()

	if cmd.Example == "" {
		t.Error("Expected Example field")
	}
}

// TestCommands_HaveRunE tests all non-parent commands have RunE
func TestCommands_HaveRunE(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"connection", newConnectionCmd()},
		{"agent", newAgentCmd()},
		{"workflow", newWorkflowCmd()},
	}

	for _, c := range commands {
		t.Run(c.name, func(t *testing.T) {
			if c.cmd.RunE == nil {
				t.Errorf("Command %s missing RunE function", c.name)
			}
			if c.cmd.Run != nil {
				t.Errorf("Command %s should use RunE, not Run", c.name)
			}
		})
	}
}

// TestCommands_UseFieldFormat tests Use field follows conventions
func TestCommands_UseFieldFormat(t *testing.T) {
	commands := []struct {
		name       string
		cmd        *cobra.Command
		expectArgs bool
	}{
		{"debug", NewDebugCmd(), false},
		{"connection", newConnectionCmd(), true},
		{"agent", newAgentCmd(), true},
		{"workflow", newWorkflowCmd(), true},
	}

	for _, c := range commands {
		t.Run(c.name, func(t *testing.T) {
			if c.cmd.Use == "" {
				t.Errorf("Command %s has empty Use field", c.name)
			}

			// Commands expecting args should have < or [ in Use
			if c.expectArgs {
				hasArg := false
				for _, ch := range c.cmd.Use {
					if ch == '<' || ch == '[' {
						hasArg = true
						break
					}
				}
				if !hasArg {
					t.Errorf("Command %s expects args but Use doesn't show them: %s", c.name, c.cmd.Use)
				}
			}
		})
	}
}

// TestCommands_HaveUniqueNames tests command names are unique
func TestCommands_HaveUniqueNames(t *testing.T) {
	debugCmd := NewDebugCmd()

	subcommands := debugCmd.Commands()
	names := make(map[string]bool)

	for _, cmd := range subcommands {
		if names[cmd.Name()] {
			t.Errorf("Duplicate command name: %s", cmd.Name())
		}
		names[cmd.Name()] = true
	}
}

// TestConnectionCmd_TimeoutFlagType tests timeout flag is int
func TestConnectionCmd_TimeoutFlagType(t *testing.T) {
	cmd := newConnectionCmd()

	timeoutFlag := cmd.Flags().Lookup("timeout")
	if timeoutFlag == nil {
		t.Fatal("Expected timeout flag")
	}

	if timeoutFlag.Value.Type() != "int" {
		t.Errorf("Expected timeout flag type 'int', got %s", timeoutFlag.Value.Type())
	}
}

// TestConnectionCmd_VerboseFlagType tests verbose flag is bool
func TestConnectionCmd_VerboseFlagType(t *testing.T) {
	cmd := newConnectionCmd()

	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Fatal("Expected verbose flag")
	}

	if verboseFlag.Value.Type() != "bool" {
		t.Errorf("Expected verbose flag type 'bool', got %s", verboseFlag.Value.Type())
	}
}

// TestAgentCmd_FullFlagType tests full flag is bool
func TestAgentCmd_FullFlagType(t *testing.T) {
	cmd := newAgentCmd()

	fullFlag := cmd.Flags().Lookup("full")
	if fullFlag == nil {
		t.Fatal("Expected full flag")
	}

	if fullFlag.Value.Type() != "bool" {
		t.Errorf("Expected full flag type 'bool', got %s", fullFlag.Value.Type())
	}
}

// TestWorkflowCmd_LastFlagType tests last flag is int
func TestWorkflowCmd_LastFlagType(t *testing.T) {
	cmd := newWorkflowCmd()

	lastFlag := cmd.Flags().Lookup("last")
	if lastFlag == nil {
		t.Fatal("Expected last flag")
	}

	if lastFlag.Value.Type() != "int" {
		t.Errorf("Expected last flag type 'int', got %s", lastFlag.Value.Type())
	}
}

// TestConnectionCmd_RequiresAgentName tests connection requires agent name
func TestConnectionCmd_RequiresAgentName(t *testing.T) {
	cmd := newConnectionCmd()

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error when no agent name provided")
	}
}

// TestAgentCmd_RequiresAgentName tests agent requires agent name
func TestAgentCmd_RequiresAgentName(t *testing.T) {
	cmd := newAgentCmd()

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error when no agent name provided")
	}
}

// TestWorkflowCmd_RequiresWorkflowName tests workflow requires workflow name
func TestWorkflowCmd_RequiresWorkflowName(t *testing.T) {
	cmd := newWorkflowCmd()

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error when no workflow name provided")
	}
}

// TestConnectionCmd_DefaultTimeout tests connection has default timeout
func TestConnectionCmd_DefaultTimeout(t *testing.T) {
	cmd := newConnectionCmd()

	timeoutFlag := cmd.Flags().Lookup("timeout")
	if timeoutFlag == nil {
		t.Fatal("Expected timeout flag")
	}

	if timeoutFlag.DefValue != "5" {
		t.Errorf("Expected default timeout '5', got %s", timeoutFlag.DefValue)
	}
}

// TestWorkflowCmd_DefaultLast tests workflow has default last
func TestWorkflowCmd_DefaultLast(t *testing.T) {
	cmd := newWorkflowCmd()

	lastFlag := cmd.Flags().Lookup("last")
	if lastFlag == nil {
		t.Fatal("Expected last flag")
	}

	if lastFlag.DefValue != "1" {
		t.Errorf("Expected default last '1', got %s", lastFlag.DefValue)
	}
}

// TestDebugCmd_LongDescription tests debug command has detailed help
func TestDebugCmd_LongDescription(t *testing.T) {
	cmd := NewDebugCmd()

	if len(cmd.Long) < 30 {
		t.Errorf("Expected detailed Long description, got %d chars", len(cmd.Long))
	}
}

// TestConnectionCmd_LongDescription tests connection command has detailed help
func TestConnectionCmd_LongDescription(t *testing.T) {
	cmd := newConnectionCmd()

	if len(cmd.Long) < 40 {
		t.Errorf("Expected detailed Long description, got %d chars", len(cmd.Long))
	}
}

// TestAgentCmd_LongDescription tests agent command has detailed help
func TestAgentCmd_LongDescription(t *testing.T) {
	cmd := newAgentCmd()

	if len(cmd.Long) < 40 {
		t.Errorf("Expected detailed Long description, got %d chars", len(cmd.Long))
	}
}

// TestWorkflowCmd_LongDescription tests workflow command has detailed help
func TestWorkflowCmd_LongDescription(t *testing.T) {
	cmd := newWorkflowCmd()

	if len(cmd.Long) < 30 {
		t.Errorf("Expected detailed Long description, got %d chars", len(cmd.Long))
	}
}

// TestCommands_AllHaveShortDescriptions tests all commands have short descriptions
func TestCommands_AllHaveShortDescriptions(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"debug", NewDebugCmd()},
		{"connection", newConnectionCmd()},
		{"agent", newAgentCmd()},
		{"workflow", newWorkflowCmd()},
	}

	for _, c := range commands {
		t.Run(c.name, func(t *testing.T) {
			if c.cmd.Short == "" {
				t.Errorf("Command %s missing Short description", c.name)
			}
		})
	}
}

// TestCommands_AllHaveLongDescriptions tests all commands have long descriptions
func TestCommands_AllHaveLongDescriptions(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"debug", NewDebugCmd()},
		{"connection", newConnectionCmd()},
		{"agent", newAgentCmd()},
		{"workflow", newWorkflowCmd()},
	}

	for _, c := range commands {
		t.Run(c.name, func(t *testing.T) {
			if c.cmd.Long == "" {
				t.Errorf("Command %s missing Long description", c.name)
			}
		})
	}
}

// TestCommands_AllHaveExamples tests all commands have examples
func TestCommands_AllHaveExamples(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"debug", NewDebugCmd()},
		{"connection", newConnectionCmd()},
		{"agent", newAgentCmd()},
		{"workflow", newWorkflowCmd()},
	}

	for _, c := range commands {
		t.Run(c.name, func(t *testing.T) {
			if c.cmd.Example == "" {
				t.Errorf("Command %s missing Example", c.name)
			}
		})
	}
}

// TestDebugCmd_IsParentCommand tests debug is a parent command
func TestDebugCmd_IsParentCommand(t *testing.T) {
	cmd := NewDebugCmd()

	// Parent command should have subcommands
	subcommands := cmd.Commands()
	if len(subcommands) == 0 {
		t.Error("Expected parent command to have subcommands")
	}
}

// TestWorkflowCmd_AcceptsLatest tests workflow accepts 'latest' keyword
func TestWorkflowCmd_AcceptsLatest(t *testing.T) {
	cmd := newWorkflowCmd()

	err := cmd.Args(cmd, []string{"latest"})
	if err != nil {
		t.Errorf("Expected workflow to accept 'latest' keyword, got error: %v", err)
	}
}

// TestConnectionCmd_TimeoutRange tests connection timeout is positive
func TestConnectionCmd_TimeoutRange(t *testing.T) {
	cmd := newConnectionCmd()

	timeoutFlag := cmd.Flags().Lookup("timeout")
	if timeoutFlag == nil {
		t.Fatal("Expected timeout flag")
	}

	// Default should be a reasonable positive number
	if timeoutFlag.DefValue != "5" {
		t.Errorf("Expected reasonable default timeout, got %s", timeoutFlag.DefValue)
	}
}
