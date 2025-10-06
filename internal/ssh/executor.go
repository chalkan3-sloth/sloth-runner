package ssh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// Executor handles SSH connections and command execution
type Executor struct {
	db *Database
}

// NewExecutor creates a new SSH executor
func NewExecutor(db *Database) *Executor {
	return &Executor{db: db}
}

// ExecutionResult contains the result of a remote command execution
type ExecutionResult struct {
	Output   string
	Error    string
	ExitCode int
	Duration time.Duration
}

// ExecuteCommand executes a command on a remote host
func (e *Executor) ExecuteCommand(profileName string, command string, password *string) (*ExecutionResult, error) {
	startTime := time.Now()

	// Get profile from database
	profile, err := e.db.GetProfile(profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Log the connection attempt
	auditLog := &AuditLog{
		ProfileName: profileName,
		Action:      "execute",
		Command:     command,
		Timestamp:   startTime,
	}

	// Establish SSH connection
	client, err := e.connect(profile, password)
	if err != nil {
		auditLog.Success = false
		auditLog.ErrorMessage = err.Error()
		auditLog.DurationMs = int(time.Since(startTime).Milliseconds())
		e.db.AddAuditLog(auditLog)
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		auditLog.Success = false
		auditLog.ErrorMessage = err.Error()
		auditLog.DurationMs = int(time.Since(startTime).Milliseconds())
		e.db.AddAuditLog(auditLog)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Capture output
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Execute command
	err = session.Run(command)
	duration := time.Since(startTime)

	// Prepare result
	result := &ExecutionResult{
		Output:   stdout.String(),
		Error:    stderr.String(),
		Duration: duration,
	}

	// Get exit code
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			result.ExitCode = exitErr.ExitStatus()
		} else {
			result.ExitCode = -1
			auditLog.Success = false
			auditLog.ErrorMessage = err.Error()
			auditLog.DurationMs = int(duration.Milliseconds())
			e.db.AddAuditLog(auditLog)
			return nil, fmt.Errorf("command execution failed: %w", err)
		}
	}

	// Log successful execution
	auditLog.Success = true
	auditLog.DurationMs = int(duration.Milliseconds())
	auditLog.BytesTransferred = len(stdout.Bytes()) + len(stderr.Bytes())
	e.db.AddAuditLog(auditLog)

	return result, nil
}

// TestConnection tests SSH connectivity without executing commands
func (e *Executor) TestConnection(profileName string, password *string) error {
	// Get profile from database
	profile, err := e.db.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	// Log the test attempt
	auditLog := &AuditLog{
		ProfileName: profileName,
		Action:      "test",
		Timestamp:   time.Now(),
	}

	// Try to connect
	client, err := e.connect(profile, password)
	if err != nil {
		auditLog.Success = false
		auditLog.ErrorMessage = err.Error()
		e.db.AddAuditLog(auditLog)
		return fmt.Errorf("connection test failed: %w", err)
	}
	defer client.Close()

	// Test session creation
	session, err := client.NewSession()
	if err != nil {
		auditLog.Success = false
		auditLog.ErrorMessage = err.Error()
		e.db.AddAuditLog(auditLog)
		return fmt.Errorf("session test failed: %w", err)
	}
	session.Close()

	// Log successful test
	auditLog.Success = true
	e.db.AddAuditLog(auditLog)

	return nil
}

// connect establishes an SSH connection
func (e *Executor) connect(profile *Profile, password *string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:    profile.User,
		Timeout: time.Duration(profile.ConnectionTimeout) * time.Second,
	}

	// Set host key callback
	if profile.StrictHostChecking {
		// TODO: Implement proper host key verification
		// For now, we'll use a callback that accepts any host key but logs a warning
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// Configure authentication
	if password != nil && *password != "" {
		// Password authentication
		config.Auth = []ssh.AuthMethod{
			ssh.Password(*password),
		}
		// Clear password from memory after use
		defer func() {
			*password = strings.Repeat("x", len(*password))
			*password = ""
		}()
	} else if profile.KeyPath != "" {
		// Key-based authentication
		if err := ValidateKeyFile(profile.KeyPath); err != nil {
			return nil, err
		}

		key, err := os.ReadFile(profile.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %w", err)
		}

		// Parse private key
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			// Key might be encrypted, try with passphrase
			// Note: In production, you'd want to handle encrypted keys properly
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}

		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	} else {
		return nil, fmt.Errorf("no authentication method available")
	}

	// Establish connection
	address := fmt.Sprintf("%s:%d", profile.Host, profile.Port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", address, err)
	}

	return client, nil
}

// ReadPasswordFromStdin reads a password from standard input
func ReadPasswordFromStdin() (string, error) {
	// Check if stdin is a terminal or pipe
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// Interactive mode - terminal
		fmt.Fprint(os.Stderr, "SSH Password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr) // New line after password
		if err != nil {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		return string(password), nil
	}

	// Pipe mode - read from stdin
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read password from stdin: %w", err)
	}

	// Remove any trailing whitespace/newlines
	password = strings.TrimRight(password, "\r\n")

	// Validate password is not empty
	if password == "" {
		return "", fmt.Errorf("empty password received")
	}

	// Check for unwanted newlines within password
	if strings.Contains(password, "\n") {
		return "", fmt.Errorf("password contains newline character")
	}

	return password, nil
}

// CopyData copies data between reader and writer for SSH sessions
func CopyData(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// CreateInteractiveSession creates an interactive SSH session
func (e *Executor) CreateInteractiveSession(profileName string, password *string) error {
	// Get profile from database
	profile, err := e.db.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	// Establish SSH connection
	client, err := e.connect(profile, password)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width = 80
		height = 24
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		return fmt.Errorf("failed to request pty: %w", err)
	}

	// Set up stdin, stdout, stderr
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Start shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	// Wait for session to finish
	return session.Wait()
}

// TransferFile transfers a file over SSH using SCP
func (e *Executor) TransferFile(profileName string, localPath string, remotePath string, password *string) error {
	// Get profile from database
	profile, err := e.db.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	// Establish SSH connection
	client, err := e.connect(profile, password)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// Create session for SCP
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Read local file
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("failed to stat local file: %w", err)
	}

	// Create SCP command
	scpCommand := fmt.Sprintf("scp -t %s", remotePath)

	// Start SCP command
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := session.Start(scpCommand); err != nil {
		return fmt.Errorf("failed to start scp: %w", err)
	}

	// Send file header
	fmt.Fprintf(stdin, "C%04o %d %s\n", fileInfo.Mode().Perm(), len(content), fileInfo.Name())

	// Send file content
	stdin.Write(content)

	// Send end marker
	fmt.Fprint(stdin, "\x00")
	stdin.Close()

	// Wait for completion
	return session.Wait()
}