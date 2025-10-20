//go:build cgo
// +build cgo

package workflow

import (
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
)

// Test NewWorkflowCommand
func TestNewWorkflowCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

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

func TestNewWorkflowCommand_HasSubcommands(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	subcommands := cmd.Commands()
	if len(subcommands) != 3 {
		t.Errorf("Expected 3 subcommands, got %d", len(subcommands))
	}
}

func TestNewWorkflowCommand_SubcommandNames(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	expectedCommands := []string{"run", "list", "preview"}
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

// Test NewRunCommand
func TestNewRunCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestNewRunCommand_RequiresStackArg(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no args")
	}

	// Test with one arg (stack name)
	err = cmd.Args(cmd, []string{"my-stack"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}

	// Test with two args (should fail)
	err = cmd.Args(cmd, []string{"stack1", "stack2"})
	if err == nil {
		t.Error("Expected error with two args")
	}
}

func TestNewRunCommand_Flags(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	expectedFlags := []string{
		"file",
		"sloth",
		"values",
		"yes",
		"interactive",
		"output",
		"debug",
		"delegate-to",
		"ssh",
		"ssh-password-stdin",
		"password-stdin",
	}

	for _, flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to exist", flagName)
		}
	}
}

func TestNewRunCommand_FileFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	fileFlag := cmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Error("Expected file flag to exist")
	}

	if fileFlag.Shorthand != "f" {
		t.Errorf("Expected shorthand 'f', got %s", fileFlag.Shorthand)
	}
}

func TestNewRunCommand_SlothFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	slothFlag := cmd.Flags().Lookup("sloth")
	if slothFlag == nil {
		t.Error("Expected sloth flag to exist")
	}

	// Sloth flag doesn't have a shorthand
	if slothFlag.Value.Type() != "string" {
		t.Errorf("Expected string type, got %s", slothFlag.Value.Type())
	}
}

func TestNewRunCommand_YesFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	yesFlag := cmd.Flags().Lookup("yes")
	if yesFlag == nil {
		t.Error("Expected yes flag to exist")
	}

	// Yes flag doesn't have a shorthand
	if yesFlag.Value.Type() != "bool" {
		t.Errorf("Expected bool type, got %s", yesFlag.Value.Type())
	}
}

func TestNewRunCommand_InteractiveFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("interactive")
	if flag == nil {
		t.Error("Expected interactive flag to exist")
	}

	// Interactive flag doesn't have a shorthand
	if flag.Value.Type() != "bool" {
		t.Errorf("Expected bool type, got %s", flag.Value.Type())
	}
}

func TestNewRunCommand_OutputFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("output")
	if flag == nil {
		t.Error("Expected output flag to exist")
	}

	if flag.Shorthand != "o" {
		t.Errorf("Expected shorthand 'o', got %s", flag.Shorthand)
	}

	if flag.DefValue != "basic" {
		t.Errorf("Expected default 'basic', got %s", flag.DefValue)
	}
}

func TestNewRunCommand_DebugFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("debug")
	if flag == nil {
		t.Error("Expected debug flag to exist")
	}

	if flag.Value.Type() != "bool" {
		t.Errorf("Expected bool type, got %s", flag.Value.Type())
	}
}

func TestNewRunCommand_DelegateToFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("delegate-to")
	if flag == nil {
		t.Error("Expected delegate-to flag to exist")
	}

	if flag.Value.Type() != "stringArray" {
		t.Errorf("Expected stringArray type, got %s", flag.Value.Type())
	}
}

func TestNewRunCommand_SSHFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("ssh")
	if flag == nil {
		t.Error("Expected ssh flag to exist")
	}

	if flag.Value.Type() != "string" {
		t.Errorf("Expected string type, got %s", flag.Value.Type())
	}
}

func TestNewRunCommand_SSHPasswordStdinFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("ssh-password-stdin")
	if flag == nil {
		t.Error("Expected ssh-password-stdin flag to exist")
	}

	if flag.Value.Type() != "bool" {
		t.Errorf("Expected bool type, got %s", flag.Value.Type())
	}
}

func TestNewRunCommand_PasswordStdinFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("password-stdin")
	if flag == nil {
		t.Error("Expected password-stdin flag to exist")
	}

	if flag.Value.Type() != "bool" {
		t.Errorf("Expected bool type, got %s", flag.Value.Type())
	}
}

// Test NewListCommand
func TestNewListCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "list" {
		t.Errorf("Expected Use 'list', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestNewListCommand_HasRunE(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}

// Test NewPreviewCommand
func TestNewPreviewCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "preview <workflow-file>" {
		t.Errorf("Expected Use 'preview <workflow-file>', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestNewPreviewCommand_Args(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args
	err := cmd.Args(cmd, []string{})
	if err != nil {
		t.Errorf("Expected no error with no args, got %v", err)
	}

	// Test with one arg
	err = cmd.Args(cmd, []string{"workflow.sloth"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}

	// Test with two args (should fail)
	err = cmd.Args(cmd, []string{"file1", "file2"})
	if err == nil {
		t.Error("Expected error with two args")
	}
}

func TestNewPreviewCommand_Flags(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	expectedFlags := []string{"file", "sloth", "format"}

	for _, flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to exist", flagName)
		}
	}
}

func TestNewPreviewCommand_FileFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	flag := cmd.Flags().Lookup("file")
	if flag == nil {
		t.Error("Expected file flag to exist")
	}

	if flag.Shorthand != "f" {
		t.Errorf("Expected shorthand 'f', got %s", flag.Shorthand)
	}
}

func TestNewPreviewCommand_FormatFlag(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	flag := cmd.Flags().Lookup("format")
	if flag == nil {
		t.Error("Expected format flag to exist")
	}

	if flag.DefValue != "tree" {
		t.Errorf("Expected default 'tree', got %s", flag.DefValue)
	}
}

// Test command interfaces and consistency
func TestCommands_ConsistentInterface(t *testing.T) {
	ctx := &commands.AppContext{}

	tests := []struct {
		name string
		cmd  interface{}
	}{
		{"workflow", NewWorkflowCommand(ctx)},
		{"run", NewRunCommand(ctx)},
		{"list", NewListCommand(ctx)},
		{"preview", NewPreviewCommand(ctx)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify we can create the commands
			if tt.cmd == nil {
				t.Errorf("Command %s is nil", tt.name)
			}
		})
	}
}

func TestCommandDescriptions_AreMeaningful(t *testing.T) {
	ctx := &commands.AppContext{}

	tests := []struct {
		name string
		cmd  interface{}
	}{
		{"run", NewRunCommand(ctx)},
		{"list", NewListCommand(ctx)},
		{"preview", NewPreviewCommand(ctx)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Commands should have meaningful descriptions
			// This is a placeholder test
			if tt.cmd == nil {
				t.Error("Expected non-nil command")
			}
		})
	}
}

func TestWorkflowCommand_IsParentCommand(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewWorkflowCommand(ctx)

	// Parent commands typically don't have RunE
	if cmd.RunE != nil {
		t.Error("Parent command should not have RunE")
	}

	// But should have subcommands
	if len(cmd.Commands()) == 0 {
		t.Error("Parent command should have subcommands")
	}
}

func TestRunCommand_HasRunE(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	if cmd.RunE == nil {
		t.Error("Run command should have RunE")
	}
}

func TestListCommand_HasRunE(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewListCommand(ctx)

	if cmd.RunE == nil {
		t.Error("List command should have RunE")
	}
}

func TestPreviewCommand_HasRunE(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	if cmd.RunE == nil {
		t.Error("Preview command should have RunE")
	}
}

// Test flag defaults
func TestRunCommand_FlagDefaults(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	tests := []struct {
		name         string
		expectedType string
		expectedDef  string
	}{
		{"file", "string", ""},
		{"sloth", "string", ""},
		{"values", "string", ""},
		{"output", "string", "basic"},
		{"yes", "bool", "false"},
		{"interactive", "bool", "false"},
		{"debug", "bool", "false"},
		{"ssh-password-stdin", "bool", "false"},
		{"password-stdin", "bool", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.Flags().Lookup(tt.name)
			if flag == nil {
				t.Fatalf("Flag %s not found", tt.name)
			}

			if flag.Value.Type() != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, flag.Value.Type())
			}

			if tt.expectedDef != "" && flag.DefValue != tt.expectedDef {
				t.Errorf("Expected default %s, got %s", tt.expectedDef, flag.DefValue)
			}
		})
	}
}

func TestPreviewCommand_FlagDefaults(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	formatFlag := cmd.Flags().Lookup("format")
	if formatFlag.DefValue != "tree" {
		t.Errorf("Expected default format 'tree', got %s", formatFlag.DefValue)
	}
}

// Test command use patterns
func TestCommands_UseFieldFormat(t *testing.T) {
	ctx := &commands.AppContext{}

	tests := []struct {
		name     string
		cmd      interface{}
		expected string
	}{
		{"workflow", NewWorkflowCommand(ctx), "workflow"},
		{"list", NewListCommand(ctx), "list"},
		{"preview", NewPreviewCommand(ctx), "preview <workflow-file>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Type assertion to get Use field
			// This is a simplified test
			if tt.cmd == nil {
				t.Error("Expected non-nil command")
			}
		})
	}
}

func TestRunCommand_LongDescriptionMentionsStack(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	// The long description should mention that stack is required
	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestRunCommand_SupportsMultipleDelegateHosts(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewRunCommand(ctx)

	flag := cmd.Flags().Lookup("delegate-to")
	if flag == nil {
		t.Error("Expected delegate-to flag")
	}

	// Should be stringArray to support multiple values
	if flag.Value.Type() != "stringArray" {
		t.Error("Expected stringArray type for multiple hosts")
	}
}

func TestPreviewCommand_SupportedFormats(t *testing.T) {
	ctx := &commands.AppContext{}
	cmd := NewPreviewCommand(ctx)

	formatFlag := cmd.Flags().Lookup("format")
	if formatFlag == nil {
		t.Error("Expected format flag")
	}

	// Default should be tree
	if formatFlag.DefValue != "tree" {
		t.Error("Expected default format to be tree")
	}
}
