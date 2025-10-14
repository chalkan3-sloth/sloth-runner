//go:build windows
// +build windows

package agent

import (
	"os"
	"os/signal"
	"syscall"
)

// setupShellSignals sets up signal handling for interactive shell (Windows-specific)
// Note: Windows doesn't support SIGWINCH (terminal resize signal)
func setupShellSignals(sigChan chan os.Signal) {
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
}

// isSigwinch checks if the signal is SIGWINCH (window resize)
// On Windows, this always returns false as SIGWINCH doesn't exist
func isSigwinch(sig os.Signal) bool {
	return false
}
