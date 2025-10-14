//go:build unix || darwin || linux
// +build unix darwin linux

package agent

import (
	"os"
	"os/signal"
	"syscall"
)

// setupShellSignals sets up signal handling for interactive shell (Unix-specific)
func setupShellSignals(sigChan chan os.Signal) {
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGWINCH)
}

// isSigwinch checks if the signal is SIGWINCH (window resize)
func isSigwinch(sig os.Signal) bool {
	return sig == syscall.SIGWINCH
}
