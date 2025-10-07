//go:build darwin
// +build darwin

package agent

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	TIOCPTYGNAME = 0x40807453
	TIOCPTYGRANT = 0x20007454
	TIOCPTYUNLK  = 0x20007452
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

	// Grant and unlock PTY
	if err := grantpt(ptmx); err != nil {
		ptmx.Close()
		return nil, fmt.Errorf("failed to grant PTY: %w", err)
	}

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
	var buf [128]byte
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), TIOCPTYGNAME, uintptr(unsafe.Pointer(&buf[0])))
	if err != 0 {
		return "", err
	}
	// Find null terminator
	n := 0
	for n < len(buf) && buf[n] != 0 {
		n++
	}
	return string(buf[:n]), nil
}

// grantpt grants access to the slave PTY
func grantpt(f *os.File) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), TIOCPTYGRANT, 0)
	if err != 0 {
		return err
	}
	return nil
}

// unlockpt unlocks the slave PTY
func unlockpt(f *os.File) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), TIOCPTYUNLK, 0)
	if err != 0 {
		return err
	}
	return nil
}
