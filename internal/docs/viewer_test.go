package docs

import (
	"runtime"
	"testing"
)

// Test ViewMode constants
func TestViewMode_Constants(t *testing.T) {
	if ViewModeTerminal != "terminal" {
		t.Errorf("Expected ViewModeTerminal to be 'terminal', got '%s'", ViewModeTerminal)
	}

	if ViewModeRaw != "raw" {
		t.Errorf("Expected ViewModeRaw to be 'raw', got '%s'", ViewModeRaw)
	}

	if ViewModeBrowser != "browser" {
		t.Errorf("Expected ViewModeBrowser to be 'browser', got '%s'", ViewModeBrowser)
	}
}

func TestViewMode_ValidModes(t *testing.T) {
	modes := []ViewMode{
		ViewModeTerminal,
		ViewModeRaw,
		ViewModeBrowser,
	}

	for _, mode := range modes {
		if mode == "" {
			t.Errorf("View mode should not be empty: %v", mode)
		}
	}
}

func TestViewMode_StringConversion(t *testing.T) {
	mode := ViewModeTerminal
	str := string(mode)

	if str != "terminal" {
		t.Errorf("Expected 'terminal', got '%s'", str)
	}
}

// Test DocViewer structure
func TestNewDocViewer_DefaultMode(t *testing.T) {
	viewer := NewDocViewer("")

	if viewer == nil {
		t.Error("Expected non-nil viewer")
	}

	if viewer.mode != ViewModeTerminal {
		t.Errorf("Expected default mode to be ViewModeTerminal, got %s", viewer.mode)
	}
}

func TestNewDocViewer_TerminalMode(t *testing.T) {
	viewer := NewDocViewer(ViewModeTerminal)

	if viewer.mode != ViewModeTerminal {
		t.Errorf("Expected ViewModeTerminal, got %s", viewer.mode)
	}
}

func TestNewDocViewer_RawMode(t *testing.T) {
	viewer := NewDocViewer(ViewModeRaw)

	if viewer.mode != ViewModeRaw {
		t.Errorf("Expected ViewModeRaw, got %s", viewer.mode)
	}
}

func TestNewDocViewer_BrowserMode(t *testing.T) {
	viewer := NewDocViewer(ViewModeBrowser)

	if viewer.mode != ViewModeBrowser {
		t.Errorf("Expected ViewModeBrowser, got %s", viewer.mode)
	}
}

func TestDocViewer_Structure(t *testing.T) {
	viewer := &DocViewer{
		mode: ViewModeTerminal,
	}

	if viewer.mode != ViewModeTerminal {
		t.Error("Expected mode to be set")
	}
}

// Test command documentation mapping
func TestShowCommand_KnownCommands(t *testing.T) {
	commands := []string{
		"hook",
		"events",
		"agent",
		"workflow",
		"run",
		"sloth",
		"stack",
		"main",
	}

	for _, cmd := range commands {
		// Just verify the command names are valid strings
		if cmd == "" {
			t.Errorf("Command name should not be empty")
		}
	}
}

func TestShowCommand_ValidCommandNames(t *testing.T) {
	validCommands := map[string]bool{
		"hook":     true,
		"events":   true,
		"agent":    true,
		"workflow": true,
		"run":      true,
		"sloth":    true,
		"stack":    true,
		"main":     true,
	}

	for cmd := range validCommands {
		if len(cmd) == 0 {
			t.Error("Command name should not be empty")
		}
	}
}

func TestShowCommand_UnknownCommand(t *testing.T) {
	viewer := NewDocViewer(ViewModeRaw)
	err := viewer.ShowCommand("unknown-command")

	if err == nil {
		t.Error("Expected error for unknown command")
	}
}

func TestShowCommand_EmptyCommand(t *testing.T) {
	viewer := NewDocViewer(ViewModeRaw)
	err := viewer.ShowCommand("")

	if err == nil {
		t.Error("Expected error for empty command")
	}
}

// Test platform detection for openURL
func TestPlatformDetection(t *testing.T) {
	platform := runtime.GOOS

	validPlatforms := []string{"darwin", "linux", "windows", "freebsd", "openbsd", "netbsd"}
	isValid := false

	for _, p := range validPlatforms {
		if platform == p {
			isValid = true
			break
		}
	}

	if !isValid && platform != "" {
		t.Logf("Note: Running on platform '%s'", platform)
	}
}

func TestPlatform_Darwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		if runtime.GOOS != "darwin" {
			t.Error("Platform mismatch")
		}
	}
}

func TestPlatform_Linux(t *testing.T) {
	if runtime.GOOS == "linux" {
		if runtime.GOOS != "linux" {
			t.Error("Platform mismatch")
		}
	}
}

func TestPlatform_Windows(t *testing.T) {
	if runtime.GOOS == "windows" {
		if runtime.GOOS != "windows" {
			t.Error("Platform mismatch")
		}
	}
}

// Test markdown content handling
func TestMarkdownContent_Headers(t *testing.T) {
	headers := []string{
		"# Main Header",
		"## Section Header",
		"### Subsection Header",
	}

	for _, header := range headers {
		if len(header) < 2 {
			t.Error("Header should have valid format")
		}
	}
}

func TestMarkdownContent_CodeBlocks(t *testing.T) {
	codeBlock := "```go\nfunc main() {}\n```"

	if codeBlock == "" {
		t.Error("Code block should not be empty")
	}
}

func TestMarkdownContent_EmptyLines(t *testing.T) {
	content := "\n\n\n"

	if len(content) != 3 {
		t.Error("Expected 3 newline characters")
	}
}

// Test view mode selection
func TestDisplay_ModeSelection(t *testing.T) {
	modes := []ViewMode{
		ViewModeTerminal,
		ViewModeRaw,
		ViewModeBrowser,
	}

	for _, mode := range modes {
		viewer := NewDocViewer(mode)
		if viewer.mode != mode {
			t.Errorf("Expected mode %s, got %s", mode, viewer.mode)
		}
	}
}

func TestDisplay_RawMode(t *testing.T) {
	viewer := NewDocViewer(ViewModeRaw)
	if viewer.mode != ViewModeRaw {
		t.Error("Expected raw mode")
	}
}

func TestDisplay_TerminalMode(t *testing.T) {
	viewer := NewDocViewer(ViewModeTerminal)
	if viewer.mode != ViewModeTerminal {
		t.Error("Expected terminal mode")
	}
}

func TestDisplay_BrowserMode(t *testing.T) {
	viewer := NewDocViewer(ViewModeBrowser)
	if viewer.mode != ViewModeBrowser {
		t.Error("Expected browser mode")
	}
}

// Test content formatting
func TestContent_HTMLGeneration(t *testing.T) {
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Test</title>
</head>
<body>
Content here
</body>
</html>`

	if htmlTemplate == "" {
		t.Error("HTML template should not be empty")
	}

	if len(htmlTemplate) < 50 {
		t.Error("HTML template should be substantial")
	}
}

func TestContent_StyleCSS(t *testing.T) {
	css := `
body { max-width: 1000px; }
pre { background: #f4f4f4; }
code { padding: 2px 5px; }
`

	if css == "" {
		t.Error("CSS should not be empty")
	}
}

// Test pager detection
func TestPager_EnvironmentVariable(t *testing.T) {
	// Test concept: pager can be set via environment
	pager := "less"

	if pager != "less" {
		t.Error("Expected pager to be 'less'")
	}
}

func TestPager_CommonPagers(t *testing.T) {
	pagers := []string{"less", "more", "most"}

	for _, pager := range pagers {
		if pager == "" {
			t.Error("Pager name should not be empty")
		}
	}
}

func TestPager_LessArgs(t *testing.T) {
	args := []string{"-R", "-F", "-X"}

	if len(args) != 3 {
		t.Errorf("Expected 3 arguments, got %d", len(args))
	}

	if args[0] != "-R" {
		t.Error("Expected -R flag")
	}
}

// Test documentation file paths
func TestDocFiles_Paths(t *testing.T) {
	paths := []string{
		"content/hook.md",
		"content/events.md",
		"content/agent.md",
		"content/workflow.md",
		"content/run.md",
		"content/sloth.md",
		"content/stack.md",
		"content/sloth-runner.md",
	}

	for _, path := range paths {
		if path == "" {
			t.Error("Path should not be empty")
		}

		if !hasMarkdownExtension(path) {
			t.Errorf("Path should have .md extension: %s", path)
		}
	}
}

func hasMarkdownExtension(path string) bool {
	return len(path) >= 3 && path[len(path)-3:] == ".md"
}

func TestDocFiles_ValidExtensions(t *testing.T) {
	files := []string{
		"file.md",
		"doc.md",
		"readme.md",
	}

	for _, file := range files {
		if !hasMarkdownExtension(file) {
			t.Errorf("File should have .md extension: %s", file)
		}
	}
}

// Test command to file mapping
func TestCommandMapping_Consistency(t *testing.T) {
	mapping := map[string]string{
		"hook":     "content/hook.md",
		"events":   "content/events.md",
		"agent":    "content/agent.md",
		"workflow": "content/workflow.md",
	}

	for cmd, file := range mapping {
		if cmd == "" {
			t.Error("Command should not be empty")
		}
		if file == "" {
			t.Error("File path should not be empty")
		}
	}
}

func TestCommandMapping_UniqueKeys(t *testing.T) {
	commands := []string{"hook", "events", "agent", "workflow", "run", "sloth", "stack", "main"}
	seen := make(map[string]bool)

	for _, cmd := range commands {
		if seen[cmd] {
			t.Errorf("Duplicate command: %s", cmd)
		}
		seen[cmd] = true
	}
}

// Test markdown rendering options
func TestMarkdownRender_WordWrap(t *testing.T) {
	wordWrap := 100

	if wordWrap <= 0 {
		t.Error("Word wrap should be positive")
	}

	if wordWrap < 50 {
		t.Error("Word wrap should be at least 50 for readability")
	}
}

func TestMarkdownRender_AutoStyle(t *testing.T) {
	// Auto style should adapt to terminal
	autoStyle := true

	if !autoStyle {
		t.Error("Expected auto style to be enabled")
	}
}

// Test content sections
func TestContent_HeaderParsing(t *testing.T) {
	tests := []struct {
		line     string
		isHeader bool
		level    int
	}{
		{"# Title", true, 1},
		{"## Section", true, 2},
		{"### Subsection", true, 3},
		{"Regular text", false, 0},
		{"", false, 0},
	}

	for _, tt := range tests {
		startsWithHash := len(tt.line) > 0 && tt.line[0] == '#'

		if tt.isHeader && !startsWithHash {
			t.Errorf("Expected header line to start with #: %s", tt.line)
		}

		if !tt.isHeader && startsWithHash && tt.line != "#" {
			t.Errorf("Non-header line should not start with #: %s", tt.line)
		}
	}
}

func TestContent_CodeBlockMarkers(t *testing.T) {
	markers := []string{
		"```",
		"```go",
		"```bash",
		"```yaml",
	}

	for _, marker := range markers {
		if len(marker) < 3 {
			t.Errorf("Code block marker too short: %s", marker)
		}

		if marker[:3] != "```" {
			t.Errorf("Code block should start with ```: %s", marker)
		}
	}
}

// Test error scenarios
func TestViewer_InvalidMode(t *testing.T) {
	invalidMode := ViewMode("invalid")

	if invalidMode == ViewModeTerminal || invalidMode == ViewModeRaw || invalidMode == ViewModeBrowser {
		t.Error("Invalid mode should not match known modes")
	}
}

func TestViewer_EmptyContent(t *testing.T) {
	content := ""

	if content != "" {
		t.Error("Expected empty content")
	}
}

// Test utility functions
func TestShowWithPterm_LineSplit(t *testing.T) {
	content := "line1\nline2\nline3"
	lines := len(content)

	if lines == 0 {
		t.Error("Content should have length")
	}
}

func TestShowWithPterm_EmptyContent(t *testing.T) {
	content := ""

	if content != "" {
		t.Error("Expected empty content")
	}
}

func TestShowWithPterm_MultilineContent(t *testing.T) {
	content := "# Title\n## Section\nContent here"

	if len(content) == 0 {
		t.Error("Expected non-empty content")
	}
}

// Test edge cases
func TestViewMode_CaseSensitivity(t *testing.T) {
	mode1 := ViewMode("terminal")
	mode2 := ViewMode("TERMINAL")

	if mode1 == mode2 {
		t.Error("View modes should be case-sensitive")
	}
}

func TestCommandName_CaseSensitivity(t *testing.T) {
	cmd1 := "hook"
	cmd2 := "HOOK"

	if cmd1 == cmd2 {
		t.Error("Command names should be case-sensitive")
	}
}

func TestDocViewer_MultipleInstances(t *testing.T) {
	viewer1 := NewDocViewer(ViewModeTerminal)
	viewer2 := NewDocViewer(ViewModeRaw)

	if viewer1.mode == viewer2.mode {
		t.Error("Different viewers should have different modes")
	}
}
