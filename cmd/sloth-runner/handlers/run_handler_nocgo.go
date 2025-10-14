//go:build !cgo
// +build !cgo

package handlers

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RunHandler stub for non-CGO builds
type RunHandler struct{}

// NewRunHandler returns a stub
func NewRunHandler() *RunHandler {
	return &RunHandler{}
}

// Handle returns an error for non-CGO builds
func (h *RunHandler) Handle(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("the 'run' command with stack/state management requires CGO support (SQLite).\n" +
		"This binary was compiled without CGO support for maximum portability.\n\n" +
		"To use stack/state management features:\n" +
		"  1. Download the CGO-enabled version for Linux (amd64/arm64)\n" +
		"  2. Or compile from source with CGO_ENABLED=1\n\n" +
		"This binary can still be used for basic command execution without state management.")
}

// HandleGroup stub
func (h *RunHandler) HandleGroup(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("stack/state management requires CGO support (SQLite). Please use a CGO-enabled build")
}

// Close stub
func (h *RunHandler) Close() error {
	return nil
}
