//go:build linux
// +build linux

package agent

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	TIOCGPTN   = 0x80045430
	TIOCSPTLCK = 0x40045431
)

// startPty creates a PTY and starts the command attached to it
func startPty(cmd *exec.Cmd) (*os.File, error) {
	// Open PTY master
	ptmx, err := openPty()
	if err != nil {
		return nil, fmt.Errorf("failed to open PTY: %w", err)
	}

	// Get PTY slave name
	sname, err := ptsname(ptmx)
	if err != nil {
		ptmx.Close()
		return nil, fmt.Errorf("failed to get PTY name: %w", err)
	}

	// Unlock PTY
	if err := unlockpt(ptmx); err != nil {
		ptmx.Close()
		return nil, fmt.Errorf("failed to unlock PTY: %w", err)
	}

	// Open PTY slave
	pts, err := os.OpenFile(sname, os.O_RDWR, 0)
	if err != nil {
		ptmx.Close()
		return nil, fmt.Errorf("failed to open PTY slave: %w", err)
	}

	// Set up the command to use the PTY
	cmd.Stdin = pts
	cmd.Stdout = pts
	cmd.Stderr = pts
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		pts.Close()
		ptmx.Close()
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Close slave side in parent process
	pts.Close()

	return ptmx, nil
}

// openPty opens a new PTY master
func openPty() (*os.File, error) {
	return os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
}

// ptsname returns the name of the slave PTY
func ptsname(f *os.File) (string, error) {
	var n uintptr
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != 0 {
		return "", err
	}
	return fmt.Sprintf("/dev/pts/%d", n), nil
}

// unlockpt unlocks the slave PTY
func unlockpt(f *os.File) error {
	var u uintptr
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
	if err != 0 {
		return err
	}
	return nil
}
