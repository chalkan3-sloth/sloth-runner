package docs

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pterm/pterm"
)

//go:embed content/*.md
var docsFS embed.FS

// ViewMode represents how to display documentation
type ViewMode string

const (
	ViewModeTerminal ViewMode = "terminal" // Display in terminal with pager
	ViewModeRaw      ViewMode = "raw"      // Display raw markdown
	ViewModeBrowser  ViewMode = "browser"  // Open in browser
)

// DocViewer handles documentation display
type DocViewer struct {
	mode ViewMode
}

// NewDocViewer creates a new documentation viewer
func NewDocViewer(mode ViewMode) *DocViewer {
	if mode == "" {
		mode = ViewModeTerminal
	}
	return &DocViewer{mode: mode}
}

// ShowCommand displays documentation for a specific command
func (v *DocViewer) ShowCommand(command string) error {
	// Map command names to documentation files
	docFiles := map[string]string{
		"hook":     "content/hook.md",
		"events":   "content/events.md",
		"agent":    "content/agent.md",
		"workflow": "content/workflow.md",
		"run":      "content/run.md",
		"sloth":    "content/sloth.md",
		"stack":    "content/stack.md",
		"main":     "content/sloth-runner.md",
	}

	docFile, exists := docFiles[command]
	if !exists {
		return fmt.Errorf("no documentation found for command: %s", command)
	}

	// Read documentation content from embedded filesystem
	content, err := docsFS.ReadFile(docFile)
	if err != nil {
		return fmt.Errorf("failed to read documentation: %w", err)
	}

	return v.display(string(content))
}

// display renders and displays the documentation based on the view mode
func (v *DocViewer) display(content string) error {
	switch v.mode {
	case ViewModeRaw:
		fmt.Println(content)
		return nil

	case ViewModeBrowser:
		return v.openInBrowser(content)

	case ViewModeTerminal:
		return v.displayInTerminal(content)

	default:
		return fmt.Errorf("unsupported view mode: %s", v.mode)
	}
}

// displayInTerminal renders markdown and displays in terminal with pager
func (v *DocViewer) displayInTerminal(content string) error {
	// Try to use glamour for nice markdown rendering
	rendered, err := v.renderMarkdown(content)
	if err != nil {
		// Fallback to raw content if rendering fails
		rendered = content
	}

	// Try to use pager (less, more) if available
	if v.usePager(rendered) == nil {
		return nil
	}

	// Fallback: just print to stdout
	fmt.Println(rendered)
	return nil
}

// renderMarkdown uses glamour to render markdown with syntax highlighting
func (v *DocViewer) renderMarkdown(content string) (string, error) {
	// Configure glamour renderer with terminal-friendly settings
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return "", err
	}

	return r.Render(content)
}

// usePager attempts to display content using system pager (less/more)
func (v *DocViewer) usePager(content string) error {
	// Determine which pager to use
	pager := os.Getenv("PAGER")
	if pager == "" {
		// Try common pagers
		if _, err := exec.LookPath("less"); err == nil {
			pager = "less"
		} else if _, err := exec.LookPath("more"); err == nil {
			pager = "more"
		} else {
			return fmt.Errorf("no pager available")
		}
	}

	// Special handling for less to enable colors and exit on EOF
	pagerArgs := []string{}
	if strings.Contains(pager, "less") {
		pagerArgs = []string{"-R", "-F", "-X"}
	}

	cmd := exec.Command(pager, pagerArgs...)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// openInBrowser converts markdown to HTML and opens in default browser
func (v *DocViewer) openInBrowser(content string) error {
	// Convert markdown to HTML
	rendered, err := glamour.RenderWithEnvironmentConfig(content)
	if err != nil {
		return fmt.Errorf("failed to render markdown: %w", err)
	}

	// Create temporary HTML file
	tmpFile, err := os.CreateTemp("", "sloth-docs-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Write HTML content
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Sloth Runner Documentation</title>
    <style>
        body { max-width: 1000px; margin: 0 auto; padding: 20px; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
        pre { background: #f4f4f4; padding: 10px; border-radius: 5px; overflow-x: auto; }
        code { background: #f4f4f4; padding: 2px 5px; border-radius: 3px; }
    </style>
</head>
<body>
%s
</body>
</html>`, rendered)

	if _, err := tmpFile.WriteString(html); err != nil {
		return fmt.Errorf("failed to write HTML: %w", err)
	}

	// Open in default browser
	return openURL(tmpFile.Name())
}

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// ShowWithPterm displays documentation using pterm for nice formatting
func ShowWithPterm(content string) {
	// Split content into sections
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			// Main heading
			pterm.DefaultHeader.WithFullWidth().Println(strings.TrimPrefix(line, "# "))
		} else if strings.HasPrefix(line, "## ") {
			// Section heading
			pterm.DefaultSection.Println(strings.TrimPrefix(line, "## "))
		} else if strings.HasPrefix(line, "### ") {
			// Subsection
			pterm.FgLightCyan.Println(strings.TrimPrefix(line, "### "))
		} else if strings.HasPrefix(line, "```") {
			// Code block marker - skip
			continue
		} else if line != "" {
			// Regular content
			fmt.Println(line)
		} else {
			// Empty line
			fmt.Println()
		}
	}
}
