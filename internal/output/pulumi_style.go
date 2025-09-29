package output

import (
	"fmt"
	"strings"
	"time"
	"os"

	"github.com/pterm/pterm"
)

// PulumiStyleOutput provides Pulumi-like output formatting
type PulumiStyleOutput struct {
	spinner    *pterm.SpinnerPrinter
	progressBar *pterm.ProgressbarPrinter
	indent     int
	outputs    map[string]interface{}
}

// NewPulumiStyleOutput creates a new Pulumi-style output formatter
func NewPulumiStyleOutput() *PulumiStyleOutput {
	return &PulumiStyleOutput{
		outputs: make(map[string]interface{}),
		indent:  0,
	}
}

// StartOperation starts a new operation with spinner
func (p *PulumiStyleOutput) StartOperation(operation string) {
	indentStr := strings.Repeat("  ", p.indent)
	p.spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("%s%s", indentStr, operation))
}

// StopOperation stops the current operation
func (p *PulumiStyleOutput) StopOperation(success bool, message string) {
	if p.spinner != nil {
		if success {
			p.spinner.Success(message)
		} else {
			p.spinner.Fail(message)
		}
	}
}

// UpdateOperation updates the spinner text
func (p *PulumiStyleOutput) UpdateOperation(text string) {
	if p.spinner != nil {
		indentStr := strings.Repeat("  ", p.indent)
		p.spinner.UpdateText(fmt.Sprintf("%s%s", indentStr, text))
	}
}

// StartProgress starts a progress bar
func (p *PulumiStyleOutput) StartProgress(total int, title string) {
	indentStr := strings.Repeat("  ", p.indent)
	p.progressBar, _ = pterm.DefaultProgressbar.WithTotal(total).WithTitle(fmt.Sprintf("%s%s", indentStr, title)).Start()
}

// UpdateProgress updates the progress bar
func (p *PulumiStyleOutput) UpdateProgress() {
	if p.progressBar != nil {
		p.progressBar.Increment()
	}
}

// StopProgress stops the progress bar
func (p *PulumiStyleOutput) StopProgress() {
	if p.progressBar != nil {
		p.progressBar.Stop()
	}
}

// Indent increases indentation level
func (p *PulumiStyleOutput) Indent() {
	p.indent++
}

// Unindent decreases indentation level
func (p *PulumiStyleOutput) Unindent() {
	if p.indent > 0 {
		p.indent--
	}
}

// TaskStart displays task start information
func (p *PulumiStyleOutput) TaskStart(taskName, description string) {
	indentStr := strings.Repeat("  ", p.indent)
	
	// Create a styled header for the task
	header := pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithTextStyle(pterm.NewStyle(pterm.FgBlack))
	header.Printf("%sTask: %s", indentStr, taskName)
	
	if description != "" {
		pterm.Printf("%s%s %s\n", indentStr, pterm.Gray("│"), pterm.LightBlue(description))
	}
	pterm.Printf("%s%s\n", indentStr, pterm.Gray("└─"))
}

// TaskSuccess displays task success
func (p *PulumiStyleOutput) TaskSuccess(taskName string, duration time.Duration, output interface{}) {
	indentStr := strings.Repeat("  ", p.indent)
	durationStr := duration.Truncate(time.Millisecond).String()
	
	pterm.Printf("%s%s %s %s %s\n", 
		indentStr,
		pterm.Green("✓"),
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(taskName),
		pterm.Gray(fmt.Sprintf("(%s)", durationStr)),
		pterm.Green("completed"))
	
	// Store output for later display
	if output != nil {
		p.outputs[taskName] = output
	}
}

// TaskFailure displays task failure
func (p *PulumiStyleOutput) TaskFailure(taskName string, duration time.Duration, err error) {
	indentStr := strings.Repeat("  ", p.indent)
	durationStr := duration.Truncate(time.Millisecond).String()
	
	pterm.Printf("%s%s %s %s %s\n", 
		indentStr,
		pterm.Red("✗"),
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(taskName),
		pterm.Gray(fmt.Sprintf("(%s)", durationStr)),
		pterm.Red("failed"))
	
	// Show error details
	if err != nil {
		errorLines := strings.Split(err.Error(), "\n")
		for _, line := range errorLines {
			if line != "" {
				pterm.Printf("%s  %s %s\n", indentStr, pterm.Red("│"), pterm.LightRed(line))
			}
		}
	}
}

// WorkflowStart displays workflow start information
func (p *PulumiStyleOutput) WorkflowStart(workflowName, description string) {
	// Clear screen and show banner
	pterm.DefaultCenter.Printf("%s\n\n", pterm.LightCyan("🦥 Sloth Runner"))
	
	// Workflow header
	pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Printf("Workflow: %s", workflowName)
	
	if description != "" {
		pterm.Printf("%s\n", pterm.Cyan(description))
	}
	
	pterm.Printf("Started at: %s\n\n", pterm.Gray(time.Now().Format("2006-01-02 15:04:05")))
}

// WorkflowSuccess displays workflow completion
func (p *PulumiStyleOutput) WorkflowSuccess(workflowName string, duration time.Duration, taskCount int) {
	durationStr := duration.Truncate(time.Millisecond).String()
	
	pterm.Printf("\n")
	pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Printf("Workflow Completed Successfully")
	
	pterm.Printf("%s %s\n", pterm.Green("✓"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(workflowName))
	pterm.Printf("Duration: %s\n", pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(durationStr))
	pterm.Printf("Tasks executed: %s\n", pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(fmt.Sprintf("%d", taskCount)))
	
	// Show outputs if any
	p.showOutputs()
}

// WorkflowFailure displays workflow failure
func (p *PulumiStyleOutput) WorkflowFailure(workflowName string, duration time.Duration, err error) {
	durationStr := duration.Truncate(time.Millisecond).String()
	
	pterm.Printf("\n")
	pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgRed)).WithTextStyle(pterm.NewStyle(pterm.FgWhite)).Printf("Workflow Failed")
	
	pterm.Printf("%s %s\n", pterm.Red("✗"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(workflowName))
	pterm.Printf("Duration: %s\n", pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(durationStr))
	
	if err != nil {
		pterm.Printf("Error: %s\n", pterm.Red(err.Error()))
	}
}

// showOutputs displays captured outputs in Pulumi style
func (p *PulumiStyleOutput) showOutputs() {
	if len(p.outputs) == 0 {
		return
	}
	
	pterm.Printf("\n")
	pterm.DefaultHeader.WithFullWidth(false).WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithTextStyle(pterm.NewStyle(pterm.FgWhite)).Printf("Outputs")
	
	for taskName, output := range p.outputs {
		pterm.Printf("\n%s %s:\n", pterm.Cyan("├─"), pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.Bold)).Sprint(taskName))
		
		// Format output based on type
		switch v := output.(type) {
		case map[string]interface{}:
			for key, value := range v {
				pterm.Printf("  %s %s: %s\n", pterm.Gray("│"), pterm.Cyan(key), formatValue(value))
			}
		case string:
			if v != "" {
				pterm.Printf("  %s %s\n", pterm.Gray("│"), pterm.LightCyan(v))
			}
		default:
			pterm.Printf("  %s %v\n", pterm.Gray("│"), v)
		}
	}
	pterm.Printf("%s\n", pterm.Gray("└─"))
}

// formatValue formats a value for display
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return pterm.LightGreen(fmt.Sprintf("%q", v))
	case int, int32, int64:
		return pterm.LightYellow(fmt.Sprintf("%d", v))
	case float32, float64:
		return pterm.LightYellow(fmt.Sprintf("%.2f", v))
	case bool:
		if v {
			return pterm.Green("true")
		}
		return pterm.Red("false")
	default:
		return pterm.LightMagenta(fmt.Sprintf("%v", v))
	}
}

// Info displays an info message
func (p *PulumiStyleOutput) Info(message string) {
	indentStr := strings.Repeat("  ", p.indent)
	pterm.Printf("%s%s %s\n", indentStr, pterm.Blue("ℹ"), message)
}

// Warning displays a warning message
func (p *PulumiStyleOutput) Warning(message string) {
	indentStr := strings.Repeat("  ", p.indent)
	pterm.Printf("%s%s %s\n", indentStr, pterm.Yellow("⚠"), pterm.Yellow(message))
}

// Error displays an error message
func (p *PulumiStyleOutput) Error(message string) {
	indentStr := strings.Repeat("  ", p.indent)
	pterm.Printf("%s%s %s\n", indentStr, pterm.Red("✗"), pterm.Red(message))
}

// Debug displays a debug message
func (p *PulumiStyleOutput) Debug(message string) {
	if os.Getenv("SLOTH_DEBUG") == "true" {
		indentStr := strings.Repeat("  ", p.indent)
		pterm.Printf("%s%s %s\n", indentStr, pterm.Gray("🐛"), pterm.Gray(message))
	}
}

// AddOutput adds an output value for a task
func (p *PulumiStyleOutput) AddOutput(taskName string, output interface{}) {
	p.outputs[taskName] = output
}

// GetOutputs returns all captured outputs
func (p *PulumiStyleOutput) GetOutputs() map[string]interface{} {
	return p.outputs
}