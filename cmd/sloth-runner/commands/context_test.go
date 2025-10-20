package commands

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
)

// Test NewAppContext
func TestNewAppContext(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc123", "2025-01-01")

	if ctx == nil {
		t.Error("Expected non-nil context")
	}

	if ctx.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", ctx.Version)
	}

	if ctx.Commit != "abc123" {
		t.Errorf("Expected commit 'abc123', got '%s'", ctx.Commit)
	}

	if ctx.Date != "2025-01-01" {
		t.Errorf("Expected date '2025-01-01', got '%s'", ctx.Date)
	}
}

func TestNewAppContext_DefaultSurveyAsker(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc123", "2025-01-01")

	if ctx.SurveyAsker == nil {
		t.Error("Expected SurveyAsker to be initialized")
	}

	_, ok := ctx.SurveyAsker.(*taskrunner.DefaultSurveyAsker)
	if !ok {
		t.Error("Expected DefaultSurveyAsker instance")
	}
}

func TestNewAppContext_DefaultOutputWriter(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc123", "2025-01-01")

	if ctx.OutputWriter != os.Stdout {
		t.Error("Expected OutputWriter to be os.Stdout")
	}
}

func TestNewAppContext_DefaultExecCommand(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc123", "2025-01-01")

	if ctx.ExecCommand == nil {
		t.Error("Expected ExecCommand to be set")
	}

	// Test that it creates commands
	cmd := ctx.ExecCommand("echo", "test")
	if cmd == nil {
		t.Error("Expected non-nil command")
	}
}

func TestNewAppContext_DefaultOsFindProcess(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc123", "2025-01-01")

	if ctx.OsFindProcess == nil {
		t.Error("Expected OsFindProcess to be set")
	}

	// Test that it works with current process
	currentPID := os.Getpid()
	process, err := ctx.OsFindProcess(currentPID)
	if err != nil {
		t.Errorf("Expected no error finding current process: %v", err)
	}

	if process == nil {
		t.Error("Expected non-nil process")
	}
}

func TestNewAppContext_DefaultProcessSignal(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc123", "2025-01-01")

	if ctx.ProcessSignal == nil {
		t.Error("Expected ProcessSignal to be set")
	}
}

func TestNewAppContext_TestModeDefault(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc123", "2025-01-01")

	if ctx.TestMode {
		t.Error("Expected TestMode to be false by default")
	}
}

// Test AppContext fields
func TestAppContext_VersionFields(t *testing.T) {
	ctx := &AppContext{
		Version: "2.0.0",
		Commit:  "def456",
		Date:    "2025-02-01",
	}

	if ctx.Version != "2.0.0" {
		t.Error("Expected version to be set")
	}

	if ctx.Commit != "def456" {
		t.Error("Expected commit to be set")
	}

	if ctx.Date != "2025-02-01" {
		t.Error("Expected date to be set")
	}
}

func TestAppContext_AgentRegistry(t *testing.T) {
	ctx := &AppContext{
		AgentRegistry: "test-registry",
	}

	if ctx.AgentRegistry == nil {
		t.Error("Expected AgentRegistry to be set")
	}

	registry, ok := ctx.AgentRegistry.(string)
	if !ok || registry != "test-registry" {
		t.Error("Expected AgentRegistry to be 'test-registry'")
	}
}

func TestAppContext_CustomSurveyAsker(t *testing.T) {
	ctx := &AppContext{
		SurveyAsker: &taskrunner.DefaultSurveyAsker{},
	}

	if ctx.SurveyAsker == nil {
		t.Error("Expected SurveyAsker to be set")
	}
}

func TestAppContext_CustomOutputWriter(t *testing.T) {
	var buf bytes.Buffer
	ctx := &AppContext{
		OutputWriter: &buf,
	}

	if ctx.OutputWriter == nil {
		t.Error("Expected OutputWriter to be set")
	}

	// Test writing
	ctx.OutputWriter.Write([]byte("test"))
	if buf.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", buf.String())
	}
}

func TestAppContext_TestMode(t *testing.T) {
	ctx := &AppContext{
		TestMode: true,
	}

	if !ctx.TestMode {
		t.Error("Expected TestMode to be true")
	}
}

func TestAppContext_CustomExecCommand(t *testing.T) {
	callCount := 0
	ctx := &AppContext{
		ExecCommand: func(name string, arg ...string) *exec.Cmd {
			callCount++
			return exec.Command(name, arg...)
		},
	}

	cmd := ctx.ExecCommand("echo", "test")
	if cmd == nil {
		t.Error("Expected non-nil command")
	}

	if callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", callCount)
	}
}

func TestAppContext_CustomOsFindProcess(t *testing.T) {
	callCount := 0
	ctx := &AppContext{
		OsFindProcess: func(pid int) (*os.Process, error) {
			callCount++
			return os.FindProcess(pid)
		},
	}

	currentPID := os.Getpid()
	process, err := ctx.OsFindProcess(currentPID)
	if err != nil {
		t.Errorf("Expected no error: %v", err)
	}

	if process == nil {
		t.Error("Expected non-nil process")
	}

	if callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", callCount)
	}
}

func TestAppContext_CustomProcessSignal(t *testing.T) {
	callCount := 0
	ctx := &AppContext{
		ProcessSignal: func(p *os.Process, sig os.Signal) error {
			callCount++
			return nil // Don't actually signal
		},
	}

	currentPID := os.Getpid()
	process, err := os.FindProcess(currentPID)
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}

	err = ctx.ProcessSignal(process, os.Interrupt)
	if err != nil {
		t.Errorf("Expected no error: %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", callCount)
	}
}

// Test version information combinations
func TestNewAppContext_EmptyVersion(t *testing.T) {
	ctx := NewAppContext("", "", "")

	if ctx.Version != "" {
		t.Error("Expected empty version")
	}

	if ctx.Commit != "" {
		t.Error("Expected empty commit")
	}

	if ctx.Date != "" {
		t.Error("Expected empty date")
	}
}

func TestNewAppContext_LongVersionInfo(t *testing.T) {
	longVersion := "1.2.3-beta.4+build.1234567890"
	longCommit := "abcdef1234567890abcdef1234567890abcdef12"
	longDate := "2025-01-01T12:34:56Z"

	ctx := NewAppContext(longVersion, longCommit, longDate)

	if ctx.Version != longVersion {
		t.Errorf("Expected long version to be preserved")
	}

	if ctx.Commit != longCommit {
		t.Errorf("Expected long commit to be preserved")
	}

	if ctx.Date != longDate {
		t.Errorf("Expected long date to be preserved")
	}
}

// Test multiple context creation
func TestNewAppContext_Multiple(t *testing.T) {
	ctx1 := NewAppContext("1.0.0", "abc", "2025-01-01")
	ctx2 := NewAppContext("2.0.0", "def", "2025-02-01")

	if ctx1.Version == ctx2.Version {
		t.Error("Expected different versions")
	}

	if ctx1.Commit == ctx2.Commit {
		t.Error("Expected different commits")
	}

	if ctx1.Date == ctx2.Date {
		t.Error("Expected different dates")
	}
}

// Test context isolation
func TestAppContext_Isolation(t *testing.T) {
	ctx1 := NewAppContext("1.0.0", "abc", "2025-01-01")
	ctx2 := NewAppContext("2.0.0", "def", "2025-02-01")

	// Modifying one should not affect the other
	ctx1.TestMode = true
	if ctx2.TestMode {
		t.Error("Modifying ctx1.TestMode should not affect ctx2")
	}
}

// Test dependency injection
func TestAppContext_DependencyInjection(t *testing.T) {
	var buf bytes.Buffer
	mockSurvey := &taskrunner.DefaultSurveyAsker{}

	ctx := &AppContext{
		SurveyAsker:  mockSurvey,
		OutputWriter: &buf,
		TestMode:     true,
	}

	if ctx.SurveyAsker != mockSurvey {
		t.Error("Expected injected SurveyAsker")
	}

	if ctx.OutputWriter != &buf {
		t.Error("Expected injected OutputWriter")
	}

	if !ctx.TestMode {
		t.Error("Expected TestMode to be true")
	}
}

// Test interface compliance
func TestAppContext_Interface(t *testing.T) {
	ctx := NewAppContext("1.0.0", "abc", "2025-01-01")

	// Verify all fields are accessible
	_ = ctx.Version
	_ = ctx.Commit
	_ = ctx.Date
	_ = ctx.AgentRegistry
	_ = ctx.SurveyAsker
	_ = ctx.OutputWriter
	_ = ctx.TestMode
	_ = ctx.ExecCommand
	_ = ctx.OsFindProcess
	_ = ctx.ProcessSignal
}

// Test default values consistency
func TestNewAppContext_DefaultsConsistency(t *testing.T) {
	ctx1 := NewAppContext("1.0.0", "abc", "2025-01-01")
	ctx2 := NewAppContext("1.0.0", "abc", "2025-01-01")

	// Both should have same defaults (but different instances)
	if ctx1 == ctx2 {
		t.Error("Expected different context instances")
	}

	// Both should use os.Stdout
	if ctx1.OutputWriter != os.Stdout || ctx2.OutputWriter != os.Stdout {
		t.Error("Expected both to use os.Stdout")
	}

	// Both should have TestMode false
	if ctx1.TestMode || ctx2.TestMode {
		t.Error("Expected both to have TestMode false")
	}
}
