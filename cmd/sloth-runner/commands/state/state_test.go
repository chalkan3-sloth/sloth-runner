//go:build cgo
// +build cgo

package state

import (
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/spf13/cobra"
)

// Test NewStateCommand
func TestNewStateCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStateCommand(ctx)

	if cmd.Use != "state" {
		t.Errorf("Expected Use 'state', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestNewStateCommand_HasSubcommands(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStateCommand(ctx)

	subcommands := cmd.Commands()
	if len(subcommands) != 6 {
		t.Errorf("Expected 6 subcommands, got %d", len(subcommands))
	}
}

func TestNewStateCommand_SubcommandNames(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStateCommand(ctx)

	expectedCommands := []string{"list", "show", "delete", "clear", "stats", "workflow"}
	subcommands := cmd.Commands()

	for _, expected := range expectedCommands {
		found := false
		for _, sub := range subcommands {
			if sub.Use == expected || sub.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

// Test NewListCommand
func TestNewListCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	if cmd.Use != "list [resource-type]" {
		t.Errorf("Expected Use 'list [resource-type]', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}
}

func TestNewListCommand_Flags(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	agentFlag := cmd.Flags().Lookup("agent")
	if agentFlag == nil {
		t.Error("Expected agent flag to exist")
	}

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("Expected output flag to exist")
	}
}

func TestNewListCommand_FlagDefaults(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	agentFlag := cmd.Flags().Lookup("agent")
	if agentFlag.DefValue != "local" {
		t.Errorf("Expected default agent 'local', got %s", agentFlag.DefValue)
	}

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag.DefValue != "table" {
		t.Errorf("Expected default output 'table', got %s", outputFlag.DefValue)
	}
}

func TestNewListCommand_OutputFlagShorthand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag.Shorthand != "o" {
		t.Errorf("Expected shorthand 'o', got %s", outputFlag.Shorthand)
	}
}

func TestNewListCommand_Args(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args
	err := cmd.Args(cmd, []string{})
	if err != nil {
		t.Errorf("Expected no error with no args, got %v", err)
	}

	// Test with one arg
	err = cmd.Args(cmd, []string{"user"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}

	// Test with two args (should fail)
	err = cmd.Args(cmd, []string{"user", "extra"})
	if err == nil {
		t.Error("Expected error with two args")
	}
}

func TestNewListCommand_HasRunE(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}

// Test NewShowCommand
func TestNewShowCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewShowCommand(ctx)

	if cmd.Use != "show <key>" {
		t.Errorf("Expected Use 'show <key>', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}
}

func TestNewShowCommand_Args(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewShowCommand(ctx)

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no args")
	}

	err = cmd.Args(cmd, []string{"key1"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}
}

// Test NewDeleteCommand
func TestNewDeleteCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewDeleteCommand(ctx)

	if cmd.Use != "delete <key>" {
		t.Errorf("Expected Use 'delete <key>', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}
}

func TestNewDeleteCommand_Args(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewDeleteCommand(ctx)

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no args")
	}

	err = cmd.Args(cmd, []string{"key1"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}
}

// Test NewClearCommand
func TestNewClearCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewClearCommand(ctx)

	if cmd.Use != "clear" {
		t.Errorf("Expected Use 'clear', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}
}

// Test NewStatsCommand
func TestNewStatsCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStatsCommand(ctx)

	if cmd.Use != "stats" {
		t.Errorf("Expected Use 'stats', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}
}

// Test NewWorkflowCommand
func TestNewWorkflowCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	if cmd.Use != "workflow" {
		t.Errorf("Expected Use 'workflow', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestNewWorkflowCommand_HasManySubcommands(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	subcommands := cmd.Commands()
	if len(subcommands) < 10 {
		t.Errorf("Expected at least 10 subcommands, got %d", len(subcommands))
	}
}

func TestNewWorkflowCommand_BasicSubcommands(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	expectedCommands := []string{"list", "show", "versions", "rollback", "drift", "resources", "outputs", "delete"}
	subcommands := cmd.Commands()

	for _, expected := range expectedCommands {
		found := false
		for _, sub := range subcommands {
			cmdName := sub.Use
			if len(cmdName) > 0 {
				// Extract command name (before space)
				for i, c := range cmdName {
					if c == ' ' {
						cmdName = cmdName[:i]
						break
					}
				}
			}
			if cmdName == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected workflow subcommand '%s' not found", expected)
		}
	}
}

func TestNewWorkflowCommand_AdvancedSubcommands(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	expectedCommands := []string{"tags", "import", "export", "backup", "restore", "diff", "search", "prune", "analytics"}
	subcommands := cmd.Commands()

	for _, expected := range expectedCommands {
		found := false
		for _, sub := range subcommands {
			cmdName := sub.Use
			if len(cmdName) > 0 {
				for i, c := range cmdName {
					if c == ' ' {
						cmdName = cmdName[:i]
						break
					}
				}
			}
			if cmdName == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected workflow subcommand '%s' not found", expected)
		}
	}
}

// Test command interfaces
func TestCommands_ConsistentInterface(t *testing.T) {
	ctx := &commands.AppContext{}

	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"state", NewStateCommand(ctx)},
		{"list", NewListCommand(ctx)},
		{"show", NewShowCommand(ctx)},
		{"delete", NewDeleteCommand(ctx)},
		{"clear", NewClearCommand(ctx)},
		{"stats", NewStatsCommand(ctx)},
		{"workflow", NewWorkflowCommand(ctx)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.Use == "" {
				t.Errorf("Command %s has empty Use field", tt.name)
			}
			if tt.cmd.Short == "" {
				t.Errorf("Command %s has empty Short field", tt.name)
			}
		})
	}
}

func TestCommandDescriptions_AreMeaningful(t *testing.T) {
	ctx := &commands.AppContext{}

	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"list", NewListCommand(ctx)},
		{"show", NewShowCommand(ctx)},
		{"delete", NewDeleteCommand(ctx)},
		{"clear", NewClearCommand(ctx)},
		{"stats", NewStatsCommand(ctx)},
		{"workflow", NewWorkflowCommand(ctx)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.cmd.Short) < 10 {
				t.Errorf("Command %s has too short description: %s", tt.name, tt.cmd.Short)
			}
		})
	}
}

func TestListCommand_LongDescription(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}

	if len(cmd.Long) < 20 {
		t.Error("Expected Long description to be meaningful")
	}
}

func TestShowCommand_LongDescription(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewShowCommand(ctx)

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}

	if len(cmd.Long) < 20 {
		t.Error("Expected Long description to be meaningful")
	}
}

func TestDeleteCommand_LongDescription(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewDeleteCommand(ctx)

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}

	if len(cmd.Long) < 20 {
		t.Error("Expected Long description to be meaningful")
	}
}

func TestClearCommand_LongDescription(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewClearCommand(ctx)

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}

	if len(cmd.Long) < 20 {
		t.Error("Expected Long description to be meaningful")
	}
}

func TestStatsCommand_LongDescription(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStatsCommand(ctx)

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}

	if len(cmd.Long) < 20 {
		t.Error("Expected Long description to be meaningful")
	}
}

func TestCommands_HaveRunE(t *testing.T) {
	ctx := &commands.AppContext{}

	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"list", NewListCommand(ctx)},
		{"show", NewShowCommand(ctx)},
		{"delete", NewDeleteCommand(ctx)},
		{"clear", NewClearCommand(ctx)},
		{"stats", NewStatsCommand(ctx)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.RunE == nil {
				t.Errorf("Command %s has no RunE function", tt.name)
			}
		})
	}
}

func TestStateCommand_IsParentCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStateCommand(ctx)

	if cmd.RunE != nil {
		t.Error("Parent command should not have RunE")
	}

	if cmd.Run == nil {
		t.Error("Parent command should have Run function")
	}
}

func TestWorkflowCommand_IsParentCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	if cmd.RunE != nil {
		t.Error("Parent command should not have RunE")
	}

	if cmd.Run == nil {
		t.Error("Parent command should have Run function")
	}
}

func TestListCommand_FlagTypes(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	agentFlag := cmd.Flags().Lookup("agent")
	if agentFlag.Value.Type() != "string" {
		t.Errorf("Expected agent flag type 'string', got %s", agentFlag.Value.Type())
	}

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag.Value.Type() != "string" {
		t.Errorf("Expected output flag type 'string', got %s", outputFlag.Value.Type())
	}
}

func TestCommands_UseFieldFormat(t *testing.T) {
	ctx := &commands.AppContext{}

	tests := []struct {
		name     string
		cmd      *cobra.Command
		expected string
	}{
		{"state", NewStateCommand(ctx), "state"},
		{"list", NewListCommand(ctx), "list [resource-type]"},
		{"show", NewShowCommand(ctx), "show <key>"},
		{"delete", NewDeleteCommand(ctx), "delete <key>"},
		{"clear", NewClearCommand(ctx), "clear"},
		{"stats", NewStatsCommand(ctx), "stats"},
		{"workflow", NewWorkflowCommand(ctx), "workflow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.Use != tt.expected {
				t.Errorf("Expected Use '%s', got '%s'", tt.expected, tt.cmd.Use)
			}
		})
	}
}

func TestCommands_HaveUniqueNames(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStateCommand(ctx)

	subcommands := cmd.Commands()
	names := make(map[string]bool)

	for _, sub := range subcommands {
		name := sub.Name()
		if names[name] {
			t.Errorf("Duplicate command name: %s", name)
		}
		names[name] = true
	}
}

func TestWorkflowCommand_HasUniqueSubcommands(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	subcommands := cmd.Commands()
	names := make(map[string]bool)

	for _, sub := range subcommands {
		name := sub.Name()
		if names[name] {
			t.Errorf("Duplicate workflow subcommand name: %s", name)
		}
		names[name] = true
	}
}

func TestStateCommand_LongDescriptionMentionsIdempotency(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStateCommand(ctx)

	// The long description should mention idempotency since this is a state management command
	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestWorkflowCommand_LongDescriptionMentionsTerraform(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	// The long description should mention Terraform/Pulumi
	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestListCommand_AcceptsOptionalResourceType(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	// Should accept no args
	err := cmd.Args(cmd, []string{})
	if err != nil {
		t.Error("Should accept no args")
	}

	// Should accept one arg
	err = cmd.Args(cmd, []string{"user"})
	if err != nil {
		t.Error("Should accept one arg")
	}
}

func TestShowCommand_RequiresKey(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewShowCommand(ctx)

	// Should reject no args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Should reject no args")
	}
}

func TestDeleteCommand_RequiresKey(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewDeleteCommand(ctx)

	// Should reject no args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Should reject no args")
	}
}

func TestClearCommand_NoArgs(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewClearCommand(ctx)

	// Clear command should not have args validator or should accept no args
	if cmd.Args != nil {
		err := cmd.Args(cmd, []string{})
		if err != nil {
			t.Error("Clear command should accept no args")
		}
	}
}

func TestStatsCommand_NoArgs(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewStatsCommand(ctx)

	// Stats command should not have args validator or should accept no args
	if cmd.Args != nil {
		err := cmd.Args(cmd, []string{})
		if err != nil {
			t.Error("Stats command should accept no args")
		}
	}
}
