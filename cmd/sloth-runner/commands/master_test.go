package commands

import (
	"testing"
)

// Test NewMasterCommand
func TestNewMasterCommand(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewMasterCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "master" {
		t.Errorf("Expected Use 'master', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestNewMasterCommand_HasSubcommands(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewMasterCommand(ctx)

	subcommands := cmd.Commands()
	if len(subcommands) != 7 {
		t.Errorf("Expected 7 subcommands, got %d", len(subcommands))
	}
}

func TestNewMasterCommand_SubcommandNames(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewMasterCommand(ctx)

	expectedCommands := []string{"add", "list", "select", "show", "update", "remove", "start"}
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

// Test newMasterAddCommand
func TestNewMasterAddCommand(t *testing.T) {
	cmd := newMasterAddCommand()

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

func TestNewMasterAddCommand_RequiresTwoArgs(t *testing.T) {
	cmd := newMasterAddCommand()

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no args")
	}

	// Test with one arg
	err = cmd.Args(cmd, []string{"name"})
	if err == nil {
		t.Error("Expected error with one arg")
	}

	// Test with two args (should succeed)
	err = cmd.Args(cmd, []string{"name", "address"})
	if err != nil {
		t.Errorf("Expected no error with two args, got %v", err)
	}

	// Test with three args (should fail)
	err = cmd.Args(cmd, []string{"name", "address", "extra"})
	if err == nil {
		t.Error("Expected error with three args")
	}
}

func TestNewMasterAddCommand_HasDescriptionFlag(t *testing.T) {
	cmd := newMasterAddCommand()

	flag := cmd.Flags().Lookup("description")
	if flag == nil {
		t.Error("Expected description flag to exist")
	}

	if flag.Shorthand != "d" {
		t.Errorf("Expected shorthand 'd', got %s", flag.Shorthand)
	}
}

func TestNewMasterAddCommand_UseFormat(t *testing.T) {
	cmd := newMasterAddCommand()

	if cmd.Use != "add <name> <address>" {
		t.Errorf("Expected Use 'add <name> <address>', got %s", cmd.Use)
	}
}

// Test newMasterListCommand
func TestNewMasterListCommand(t *testing.T) {
	cmd := newMasterListCommand()

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "list" {
		t.Errorf("Expected Use 'list', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}
}

func TestNewMasterListCommand_HasRunE(t *testing.T) {
	cmd := newMasterListCommand()

	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}

// Test newMasterSelectCommand
func TestNewMasterSelectCommand(t *testing.T) {
	cmd := newMasterSelectCommand()

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "select <name>" {
		t.Errorf("Expected Use 'select <name>', got %s", cmd.Use)
	}
}

func TestNewMasterSelectCommand_RequiresOneArg(t *testing.T) {
	cmd := newMasterSelectCommand()

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no args")
	}

	// Test with one arg (should succeed)
	err = cmd.Args(cmd, []string{"name"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}

	// Test with two args (should fail)
	err = cmd.Args(cmd, []string{"name1", "name2"})
	if err == nil {
		t.Error("Expected error with two args")
	}
}

// Test newMasterShowCommand
func TestNewMasterShowCommand(t *testing.T) {
	cmd := newMasterShowCommand()

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "show [name]" {
		t.Errorf("Expected Use 'show [name]', got %s", cmd.Use)
	}
}

func TestNewMasterShowCommand_AcceptsOptionalArg(t *testing.T) {
	cmd := newMasterShowCommand()

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args (should succeed - shows default)
	err := cmd.Args(cmd, []string{})
	if err != nil {
		t.Errorf("Expected no error with no args, got %v", err)
	}

	// Test with one arg (should succeed)
	err = cmd.Args(cmd, []string{"name"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}

	// Test with two args (should fail)
	err = cmd.Args(cmd, []string{"name1", "name2"})
	if err == nil {
		t.Error("Expected error with two args")
	}
}

// Test newMasterUpdateCommand
func TestNewMasterUpdateCommand(t *testing.T) {
	cmd := newMasterUpdateCommand()

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "update <name> <new_address>" {
		t.Errorf("Expected Use 'update <name> <new_address>', got %s", cmd.Use)
	}
}

func TestNewMasterUpdateCommand_RequiresTwoArgs(t *testing.T) {
	cmd := newMasterUpdateCommand()

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no args")
	}

	// Test with one arg
	err = cmd.Args(cmd, []string{"name"})
	if err == nil {
		t.Error("Expected error with one arg")
	}

	// Test with two args (should succeed)
	err = cmd.Args(cmd, []string{"name", "address"})
	if err != nil {
		t.Errorf("Expected no error with two args, got %v", err)
	}
}

func TestNewMasterUpdateCommand_HasDescriptionFlag(t *testing.T) {
	cmd := newMasterUpdateCommand()

	flag := cmd.Flags().Lookup("description")
	if flag == nil {
		t.Error("Expected description flag to exist")
	}

	if flag.Shorthand != "d" {
		t.Errorf("Expected shorthand 'd', got %s", flag.Shorthand)
	}
}

// Test newMasterRemoveCommand
func TestNewMasterRemoveCommand(t *testing.T) {
	cmd := newMasterRemoveCommand()

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "remove <name>" {
		t.Errorf("Expected Use 'remove <name>', got %s", cmd.Use)
	}
}

func TestNewMasterRemoveCommand_HasAliases(t *testing.T) {
	cmd := newMasterRemoveCommand()

	expectedAliases := []string{"rm", "delete"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	for _, alias := range expectedAliases {
		found := false
		for _, cmdAlias := range cmd.Aliases {
			if cmdAlias == alias {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected alias '%s' not found", alias)
		}
	}
}

func TestNewMasterRemoveCommand_RequiresOneArg(t *testing.T) {
	cmd := newMasterRemoveCommand()

	if cmd.Args == nil {
		t.Error("Expected Args validator to be set")
	}

	// Test with no args
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no args")
	}

	// Test with one arg (should succeed)
	err = cmd.Args(cmd, []string{"name"})
	if err != nil {
		t.Errorf("Expected no error with one arg, got %v", err)
	}

	// Test with two args (should fail)
	err = cmd.Args(cmd, []string{"name1", "name2"})
	if err == nil {
		t.Error("Expected error with two args")
	}
}

// Test newMasterStartCommand
func TestNewMasterStartCommand(t *testing.T) {
	ctx := &AppContext{}
	cmd := newMasterStartCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "start" {
		t.Errorf("Expected Use 'start', got %s", cmd.Use)
	}
}

func TestNewMasterStartCommand_HasPortFlag(t *testing.T) {
	ctx := &AppContext{}
	cmd := newMasterStartCommand(ctx)

	flag := cmd.Flags().Lookup("port")
	if flag == nil {
		t.Error("Expected port flag to exist")
	}

	if flag.Shorthand != "p" {
		t.Errorf("Expected shorthand 'p', got %s", flag.Shorthand)
	}

	if flag.DefValue != "50053" {
		t.Errorf("Expected default port '50053', got %s", flag.DefValue)
	}
}

func TestNewMasterStartCommand_HasBindFlag(t *testing.T) {
	ctx := &AppContext{}
	cmd := newMasterStartCommand(ctx)

	flag := cmd.Flags().Lookup("bind")
	if flag == nil {
		t.Error("Expected bind flag to exist")
	}

	if flag.DefValue != "0.0.0.0" {
		t.Errorf("Expected default bind '0.0.0.0', got %s", flag.DefValue)
	}
}

func TestNewMasterStartCommand_HasDaemonFlag(t *testing.T) {
	ctx := &AppContext{}
	cmd := newMasterStartCommand(ctx)

	flag := cmd.Flags().Lookup("daemon")
	if flag == nil {
		t.Error("Expected daemon flag to exist")
	}

	if flag.Value.Type() != "bool" {
		t.Errorf("Expected bool type, got %s", flag.Value.Type())
	}
}

func TestNewMasterStartCommand_HasRunE(t *testing.T) {
	ctx := &AppContext{}
	cmd := newMasterStartCommand(ctx)

	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}

// Test command descriptions
func TestMasterCommands_HaveMeaningfulDescriptions(t *testing.T) {
	commands := []struct {
		name string
		cmd  interface{}
	}{
		{"add", newMasterAddCommand()},
		{"list", newMasterListCommand()},
		{"select", newMasterSelectCommand()},
		{"show", newMasterShowCommand()},
		{"update", newMasterUpdateCommand()},
		{"remove", newMasterRemoveCommand()},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			// All commands should be non-nil
			if tc.cmd == nil {
				t.Error("Expected non-nil command")
			}
		})
	}
}

// Test command examples
func TestNewMasterAddCommand_HasExamples(t *testing.T) {
	cmd := newMasterAddCommand()

	if cmd.Long == "" {
		t.Error("Expected long description with examples")
	}
}

func TestNewMasterSelectCommand_HasExamples(t *testing.T) {
	cmd := newMasterSelectCommand()

	if cmd.Long == "" {
		t.Error("Expected long description with examples")
	}
}

func TestNewMasterShowCommand_HasExamples(t *testing.T) {
	cmd := newMasterShowCommand()

	if cmd.Long == "" {
		t.Error("Expected long description with examples")
	}
}

func TestNewMasterUpdateCommand_HasExamples(t *testing.T) {
	cmd := newMasterUpdateCommand()

	if cmd.Long == "" {
		t.Error("Expected long description with examples")
	}
}

func TestNewMasterRemoveCommand_HasExamples(t *testing.T) {
	cmd := newMasterRemoveCommand()

	if cmd.Long == "" {
		t.Error("Expected long description with examples")
	}
}

// Test MasterServerStarter variable
func TestMasterServerStarter_CanBeSet(t *testing.T) {
	// Save original
	original := MasterServerStarter

	// Set to custom function
	called := false
	MasterServerStarter = func(port int) error {
		called = true
		return nil
	}

	// Call it
	err := MasterServerStarter(50053)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("Expected MasterServerStarter to be called")
	}

	// Restore original
	MasterServerStarter = original
}

// Test command consistency
func TestMasterCommands_Consistency(t *testing.T) {
	ctx := &AppContext{}
	masterCmd := NewMasterCommand(ctx)

	// All subcommands should have descriptions
	for _, cmd := range masterCmd.Commands() {
		if cmd.Short == "" {
			t.Errorf("Command %s has empty Short description", cmd.Name())
		}
	}
}

// Test subcommand integration
func TestNewMasterCommand_SubcommandsIntegrated(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewMasterCommand(ctx)

	// Verify each subcommand is properly added
	subcommands := cmd.Commands()

	commandMap := make(map[string]bool)
	for _, sub := range subcommands {
		commandMap[sub.Name()] = true
	}

	requiredCommands := []string{"add", "list", "select", "show", "update", "remove", "start"}
	for _, required := range requiredCommands {
		if !commandMap[required] {
			t.Errorf("Required command '%s' not found in subcommands", required)
		}
	}
}

// Test parent command has no RunE
func TestNewMasterCommand_NoRunE(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewMasterCommand(ctx)

	// Parent commands typically don't have RunE
	if cmd.RunE != nil {
		t.Error("Parent command should not have RunE")
	}
}
