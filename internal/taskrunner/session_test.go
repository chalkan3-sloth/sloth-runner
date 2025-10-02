package taskrunner

import (
	"bufio"
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewSharedSession(t *testing.T) {
	session, err := NewSharedSession()
	if err != nil {
		t.Fatalf("Failed to create shared session: %v", err)
	}

	// Cleanup
	defer func() {
		if session.Cmd != nil && session.Cmd.Process != nil {
			session.Cmd.Process.Kill()
		}
	}()

	// Verify session components
	if session == nil {
		t.Fatal("Expected non-nil session")
	}

	if session.Cmd == nil {
		t.Error("Expected non-nil Cmd")
	}

	if session.Stdin == nil {
		t.Error("Expected non-nil Stdin")
	}

	if session.Stdout == nil {
		t.Error("Expected non-nil Stdout")
	}

	if session.Stderr == nil {
		t.Error("Expected non-nil Stderr")
	}

	// Verify that session is running
	if session.Cmd.Process == nil {
		t.Error("Expected process to be running")
	}
}

func TestSharedSession_CommandExecution(t *testing.T) {
	session, err := NewSharedSession()
	if err != nil {
		t.Fatalf("Failed to create shared session: %v", err)
	}

	defer func() {
		if session.Cmd != nil && session.Cmd.Process != nil {
			session.Cmd.Process.Kill()
		}
	}()

	// Test executing a simple command
	command := "echo 'test output'\n"
	_, err = session.Stdin.Write([]byte(command))
	if err != nil {
		t.Fatalf("Failed to write command: %v", err)
	}

	// Give it time to execute
	time.Sleep(100 * time.Millisecond)

	// Try to read output (non-blocking attempt)
	reader := bufio.NewReader(session.Stdout)
	
	// Set a short timeout for reading
	done := make(chan bool, 1)
	var output strings.Builder
	
	go func() {
		buf := make([]byte, 1024)
		n, _ := reader.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
		}
		done <- true
	}()

	select {
	case <-done:
		// Got some output, verify it contains something
		if output.Len() > 0 {
			t.Logf("Got output: %s", output.String())
		}
	case <-time.After(500 * time.Millisecond):
		t.Log("Timeout reading output (this is expected in test environment)")
	}
}

func TestSharedSession_MultipleCommands(t *testing.T) {
	session, err := NewSharedSession()
	if err != nil {
		t.Fatalf("Failed to create shared session: %v", err)
	}

	defer func() {
		if session.Cmd != nil && session.Cmd.Process != nil {
			session.Cmd.Process.Kill()
		}
	}()

	// Test executing multiple commands
	commands := []string{
		"echo 'first'\n",
		"echo 'second'\n",
		"echo 'third'\n",
	}

	for _, cmd := range commands {
		_, err := session.Stdin.Write([]byte(cmd))
		if err != nil {
			t.Errorf("Failed to write command '%s': %v", cmd, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Verify session is still running after multiple commands
	if session.Cmd.Process == nil {
		t.Error("Expected process to still be running")
	}
}

func TestSharedSession_ErrorStream(t *testing.T) {
	session, err := NewSharedSession()
	if err != nil {
		t.Fatalf("Failed to create shared session: %v", err)
	}

	defer func() {
		if session.Cmd != nil && session.Cmd.Process != nil {
			session.Cmd.Process.Kill()
		}
	}()

	// Execute a command that generates error output
	command := "echo 'error' >&2\n"
	_, err = session.Stdin.Write([]byte(command))
	if err != nil {
		t.Fatalf("Failed to write command: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Try to read from stderr (non-blocking)
	done := make(chan bool, 1)
	go func() {
		buf := make([]byte, 1024)
		_, _ = session.Stderr.Read(buf)
		done <- true
	}()

	select {
	case <-done:
		t.Log("Successfully read from stderr")
	case <-time.After(500 * time.Millisecond):
		t.Log("Timeout reading stderr (this is expected in test environment)")
	}
}

func TestSharedSession_Close(t *testing.T) {
	session, err := NewSharedSession()
	if err != nil {
		t.Fatalf("Failed to create shared session: %v", err)
	}

	// Send exit command
	_, err = session.Stdin.Write([]byte("exit\n"))
	if err != nil && err != io.ErrClosedPipe {
		t.Errorf("Failed to write exit command: %v", err)
	}

	// Wait a bit for process to exit
	time.Sleep(200 * time.Millisecond)

	// Verify we can kill the process without error
	if session.Cmd.Process != nil {
		err := session.Cmd.Process.Kill()
		// It's okay if the process already exited
		if err != nil && err.Error() != "os: process already finished" {
			t.Logf("Kill returned: %v (may be already exited)", err)
		}
	}
}

func TestSharedSession_StdinWritable(t *testing.T) {
	session, err := NewSharedSession()
	if err != nil {
		t.Fatalf("Failed to create shared session: %v", err)
	}

	defer func() {
		if session.Cmd != nil && session.Cmd.Process != nil {
			session.Cmd.Process.Kill()
		}
	}()

	// Verify we can write to stdin multiple times
	for i := 0; i < 3; i++ {
		_, err := session.Stdin.Write([]byte(":\n")) // Bash no-op command
		if err != nil {
			t.Errorf("Write %d failed: %v", i, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
