package commands

import (
	"io"
	"os"
	"os/exec"

	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
)

// AppContext holds shared dependencies for all commands
// This implements Dependency Injection pattern
type AppContext struct {
	// Version information
	Version string
	Commit  string
	Date    string

	// Global instances
	AgentRegistry interface{} // Will be properly typed later
	SurveyAsker   taskrunner.SurveyAsker

	// IO
	OutputWriter io.Writer
	TestMode     bool

	// For testing
	ExecCommand   func(name string, arg ...string) *exec.Cmd
	OsFindProcess func(pid int) (*os.Process, error)
	ProcessSignal func(p *os.Process, sig os.Signal) error
}

// NewAppContext creates a new application context with default values
func NewAppContext(version, commit, date string) *AppContext {
	return &AppContext{
		Version:       version,
		Commit:        commit,
		Date:          date,
		SurveyAsker:   &taskrunner.DefaultSurveyAsker{},
		OutputWriter:  os.Stdout,
		ExecCommand:   exec.Command,
		OsFindProcess: os.FindProcess,
		ProcessSignal: func(p *os.Process, sig os.Signal) error {
			return p.Signal(sig)
		},
	}
}
