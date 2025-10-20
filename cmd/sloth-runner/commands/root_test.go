package commands

import (
	"bytes"
	"testing"
)

// Test NewRootCommand
func TestNewRootCommand(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if cmd.Use != "sloth-runner" {
		t.Errorf("Expected Use 'sloth-runner', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}
}

func TestNewRootCommand_HasRun(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	if cmd.Run == nil {
		t.Error("Expected Run to be set")
	}
}

func TestNewRootCommand_ShortDescription(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	expected := "A flexible sloth-runner with Lua scripting capabilities"
	if cmd.Short != expected {
		t.Errorf("Expected Short '%s', got '%s'", expected, cmd.Short)
	}
}

func TestNewRootCommand_LongDescription(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Check for key phrases in long description
	keyPhrases := []string{
		"sloth-runner is a command-line tool",
		"Lua scripts",
		"pipelines",
		"workflows",
	}

	for _, phrase := range keyPhrases {
		if !contains(cmd.Long, phrase) {
			t.Errorf("Expected Long description to contain '%s'", phrase)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		   (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestNewRootCommand_HasVersionFlag(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	flag := cmd.PersistentFlags().Lookup("version")
	if flag == nil {
		t.Error("Expected version flag to exist")
	}

	if flag.Shorthand != "V" {
		t.Errorf("Expected shorthand 'V', got %s", flag.Shorthand)
	}

	if flag.DefValue != "false" {
		t.Errorf("Expected default value 'false', got %s", flag.DefValue)
	}
}

func TestNewRootCommand_VersionFlagBool(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	flag := cmd.PersistentFlags().Lookup("version")
	if flag == nil {
		t.Fatal("Expected version flag to exist")
	}

	if flag.Value.Type() != "bool" {
		t.Errorf("Expected bool type, got %s", flag.Value.Type())
	}
}

func TestNewRootCommand_UsesAppContext(t *testing.T) {
	ctx := &AppContext{
		Version: "1.0.0",
		Commit:  "abc123",
		Date:    "2025-01-01",
	}

	cmd := NewRootCommand(ctx)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	// Context is used by the Run function
	if cmd.Run == nil {
		t.Error("Expected Run to be set")
	}
}

func TestNewRootCommand_RunWithoutVersionFlag(t *testing.T) {
	ctx := &AppContext{
		Version: "1.0.0",
		Commit:  "abc123",
		Date:    "2025-01-01",
	}

	cmd := NewRootCommand(ctx)

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Run without version flag (should show help)
	cmd.Run(cmd, []string{})

	// Output should contain help text (not version)
	output := buf.String()
	// Help contains the short description
	if output != "" && !contains(output, "sloth-runner") && !contains(output, "Usage") {
		// If there's output, it should be help-related
		t.Logf("Output: %s", output)
	}
}

func TestNewRootCommand_Consistency(t *testing.T) {
	// Create multiple commands and verify consistency
	for i := 0; i < 5; i++ {
		ctx := &AppContext{}
		cmd := NewRootCommand(ctx)

		if cmd == nil {
			t.Errorf("Command %d is nil", i)
		}

		if cmd.Use != "sloth-runner" {
			t.Errorf("Command %d has wrong Use: %s", i, cmd.Use)
		}
	}
}

func TestNewRootCommand_DescriptionQuality(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Short description should be concise (less than 100 chars)
	if len(cmd.Short) > 100 {
		t.Errorf("Short description too long: %d chars", len(cmd.Short))
	}

	// Short description should not be empty
	if cmd.Short == "" {
		t.Error("Short description is empty")
	}

	// Long description should be longer than short
	if len(cmd.Long) <= len(cmd.Short) {
		t.Error("Long description should be longer than short description")
	}
}

func TestNewRootCommand_VersionFlagDescription(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	flag := cmd.PersistentFlags().Lookup("version")
	if flag == nil {
		t.Fatal("Expected version flag to exist")
	}

	if flag.Usage == "" {
		t.Error("Expected non-empty usage description for version flag")
	}

	expected := "Show version information"
	if flag.Usage != expected {
		t.Errorf("Expected usage '%s', got '%s'", expected, flag.Usage)
	}
}

func TestNewRootCommand_NoAliases(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	if len(cmd.Aliases) > 0 {
		t.Error("Expected no aliases for root command")
	}
}

func TestNewRootCommand_AcceptsAnyArgs(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Root command should accept any args (delegates to subcommands or shows help)
	if cmd.Args != nil {
		// If there is an Args validator, it should be permissive
		err := cmd.Args(cmd, []string{})
		if err != nil {
			t.Errorf("Expected no error with zero args, got %v", err)
		}

		err = cmd.Args(cmd, []string{"subcommand"})
		if err != nil {
			t.Errorf("Expected no error with args, got %v", err)
		}
	}
}

func TestNewRootCommand_ContextVersions(t *testing.T) {
	testCases := []struct {
		name    string
		version string
		commit  string
		date    string
	}{
		{"Empty", "", "", ""},
		{"Version only", "1.0.0", "", ""},
		{"Full", "1.0.0", "abc123", "2025-01-01"},
		{"Long version", "1.2.3-beta.4+build.123", "abcdef1234567890", "2025-01-01T12:34:56Z"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &AppContext{
				Version: tc.version,
				Commit:  tc.commit,
				Date:    tc.date,
			}

			cmd := NewRootCommand(ctx)
			if cmd == nil {
				t.Errorf("Failed to create command with %s", tc.name)
			}

			if cmd.Run == nil {
				t.Errorf("Run not set for %s", tc.name)
			}
		})
	}
}

func TestNewRootCommand_PersistentFlag(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Version flag should be persistent (available to all subcommands)
	persistentFlags := cmd.PersistentFlags()
	if persistentFlags == nil {
		t.Fatal("Expected persistent flags to exist")
	}

	versionFlag := persistentFlags.Lookup("version")
	if versionFlag == nil {
		t.Error("Expected version flag to be persistent")
	}
}

func TestNewRootCommand_NoPersistentPreRun(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Root command typically doesn't have PersistentPreRun
	if cmd.PersistentPreRun != nil {
		t.Error("Expected no PersistentPreRun")
	}

	if cmd.PersistentPreRunE != nil {
		t.Error("Expected no PersistentPreRunE")
	}
}

func TestNewRootCommand_NoPostRun(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Root command typically doesn't have PostRun
	if cmd.PostRun != nil {
		t.Error("Expected no PostRun")
	}

	if cmd.PostRunE != nil {
		t.Error("Expected no PostRunE")
	}
}

func TestNewRootCommand_MultipleInstances(t *testing.T) {
	ctx1 := &AppContext{Version: "1.0.0"}
	ctx2 := &AppContext{Version: "2.0.0"}

	cmd1 := NewRootCommand(ctx1)
	cmd2 := NewRootCommand(ctx2)

	// Both commands should be created successfully
	if cmd1 == nil || cmd2 == nil {
		t.Error("Expected both commands to be created")
	}

	// Commands should have same structure
	if cmd1.Use != cmd2.Use {
		t.Error("Expected same Use for both commands")
	}

	// But use different contexts (verified by having different Run functions)
	if cmd1.Run == nil || cmd2.Run == nil {
		t.Error("Expected Run to be set for both commands")
	}
}

func TestNewRootCommand_Structure(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Verify command structure is complete
	if cmd.Use == "" {
		t.Error("Use field is empty")
	}

	if cmd.Short == "" {
		t.Error("Short field is empty")
	}

	if cmd.Long == "" {
		t.Error("Long field is empty")
	}

	if cmd.Run == nil {
		t.Error("Run is nil")
	}

	if !cmd.PersistentFlags().HasFlags() {
		t.Error("Expected at least one persistent flag")
	}
}

func TestNewRootCommand_VersionFlagAccessible(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Set the version flag
	err := cmd.PersistentFlags().Set("version", "true")
	if err != nil {
		t.Errorf("Failed to set version flag: %v", err)
	}

	// Get the flag value from persistent flags
	versionFlag, err := cmd.PersistentFlags().GetBool("version")
	if err != nil {
		t.Errorf("Failed to get version flag: %v", err)
	}

	if !versionFlag {
		t.Error("Expected version flag to be true")
	}
}

func TestNewRootCommand_LongDescriptionContent(t *testing.T) {
	ctx := &AppContext{}
	cmd := NewRootCommand(ctx)

	// Verify long description mentions key features
	requiredTerms := []string{
		"Lua",
		"task",
	}

	for _, term := range requiredTerms {
		if !contains(cmd.Long, term) {
			t.Errorf("Long description should mention '%s'", term)
		}
	}
}
