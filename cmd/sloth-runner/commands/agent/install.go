package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// NewInstallCommand creates the agent install command
func NewInstallCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <agent_name>",
		Short: "Install and bootstrap a new agent on a remote host",
		Long: `Connects to a remote host via SSH and performs complete agent bootstrap:
  - Downloads latest sloth-runner binary
  - Installs to /usr/local/bin/sloth-runner
  - Creates systemd service
  - Enables and starts the agent service`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName := args[0]

			// Get SSH flags
			sshHost, _ := cmd.Flags().GetString("ssh-host")
			sshUser, _ := cmd.Flags().GetString("ssh-user")
			sshPort, _ := cmd.Flags().GetInt("ssh-port")
			sshKey, _ := cmd.Flags().GetString("ssh-key")

			// Get agent configuration flags
			masterAddr, _ := cmd.Flags().GetString("master")
			bindAddress, _ := cmd.Flags().GetString("bind-address")
			port, _ := cmd.Flags().GetInt("port")
			reportAddress, _ := cmd.Flags().GetString("report-address")

			return installAgent(agentName, InstallOptions{
				SSHHost:       sshHost,
				SSHUser:       sshUser,
				SSHPort:       sshPort,
				SSHKey:        sshKey,
				MasterAddr:    masterAddr,
				BindAddress:   bindAddress,
				Port:          port,
				ReportAddress: reportAddress,
			})
		},
	}

	// SSH connection flags
	cmd.Flags().String("ssh-host", "", "SSH host address (required)")
	cmd.Flags().String("ssh-user", "root", "SSH username")
	cmd.Flags().Int("ssh-port", 22, "SSH port")
	cmd.Flags().String("ssh-key", "", "SSH private key path (default: ~/.ssh/id_rsa)")

	// Agent configuration flags
	cmd.Flags().String("master", "localhost:50051", "Master server address")
	cmd.Flags().String("bind-address", "0.0.0.0", "Agent bind address")
	cmd.Flags().Int("port", 50051, "Agent port")
	cmd.Flags().String("report-address", "", "Address agent reports to master (optional)")

	cmd.MarkFlagRequired("ssh-host")
	cmd.MarkFlagRequired("master")

	return cmd
}

// InstallOptions contains all configuration for agent installation
type InstallOptions struct {
	SSHHost       string
	SSHUser       string
	SSHPort       int
	SSHKey        string
	MasterAddr    string
	BindAddress   string
	Port          int
	ReportAddress string
}

func installAgent(agentName string, opts InstallOptions) error {
	// Create SSH client
	pterm.Info.Printf("Connecting to %s@%s:%d...\n", opts.SSHUser, opts.SSHHost, opts.SSHPort)

	sshClient, err := createSSHClient(opts)
	if err != nil {
		pterm.Error.Printf("Failed to connect via SSH: %v\n", err)
		return err
	}
	defer sshClient.Close()

	pterm.Success.Println("SSH connection established")

	// Detect platform and architecture
	pterm.Info.Println("Detecting remote platform...")
	platform, arch, err := detectPlatform(sshClient)
	if err != nil {
		pterm.Error.Printf("Failed to detect platform: %v\n", err)
		return err
	}
	pterm.Success.Printf("Detected: %s/%s\n", platform, arch)

	// Get latest version
	pterm.Info.Println("Fetching latest sloth-runner version...")
	latestVersion, err := getLatestReleaseVersion()
	if err != nil {
		pterm.Error.Printf("Failed to fetch latest version: %v\n", err)
		return err
	}
	pterm.Success.Printf("Latest version: %s\n", latestVersion)

	// Download and install binary
	pterm.Info.Println("Downloading and installing sloth-runner binary...")
	if err := downloadAndInstallBinary(sshClient, latestVersion, platform, arch); err != nil {
		pterm.Error.Printf("Failed to install binary: %v\n", err)
		return err
	}
	pterm.Success.Println("Binary installed to /usr/local/bin/sloth-runner")

	// Create systemd service
	pterm.Info.Println("Creating systemd service...")
	if err := createSystemdService(sshClient, agentName, opts); err != nil {
		pterm.Error.Printf("Failed to create systemd service: %v\n", err)
		return err
	}
	pterm.Success.Println("Systemd service created")

	// Enable and start service
	pterm.Info.Println("Enabling and starting agent service...")
	if err := enableAndStartService(sshClient, agentName); err != nil {
		pterm.Error.Printf("Failed to start service: %v\n", err)
		return err
	}
	pterm.Success.Println("Agent service started")

	// Verify agent is running
	pterm.Info.Println("Verifying agent is running...")
	time.Sleep(3 * time.Second) // Give agent time to start

	if err := verifyAgentRunning(sshClient, agentName); err != nil {
		pterm.Warning.Printf("Agent may not be running: %v\n", err)
	} else {
		pterm.Success.Println("Agent is running!")
	}

	pterm.Println()
	pterm.Success.Printf("✅ Agent '%s' installed successfully!\n", agentName)
	pterm.Info.Printf("   Service: sloth-runner-agent-%s\n", agentName)
	pterm.Info.Printf("   Address: %s:%d\n", opts.SSHHost, opts.Port)
	pterm.Info.Printf("   Master: %s\n", opts.MasterAddr)
	pterm.Println()
	pterm.Info.Println("Check agent status with:")
	pterm.Printf("  systemctl status sloth-runner-agent-%s\n", agentName)
	pterm.Println()

	return nil
}

// createSSHClient creates an SSH client connection
func createSSHClient(opts InstallOptions) (*ssh.Client, error) {
	// Determine SSH key path
	keyPath := opts.SSHKey
	if keyPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		keyPath = filepath.Join(home, ".ssh", "id_rsa")
	}

	// Read private key
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	// Configure SSH client
	config := &ssh.ClientConfig{
		User: opts.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Make this configurable
		Timeout:         30 * time.Second,
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", opts.SSHHost, opts.SSHPort)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}

	return client, nil
}

// runSSHCommand executes a command on the remote host
func runSSHCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// detectPlatform detects the OS and architecture of the remote host
func detectPlatform(client *ssh.Client) (string, string, error) {
	// Get OS
	osOutput, err := runSSHCommand(client, "uname -s")
	if err != nil {
		return "", "", err
	}
	platform := strings.ToLower(strings.TrimSpace(osOutput))

	// Get architecture
	archOutput, err := runSSHCommand(client, "uname -m")
	if err != nil {
		return "", "", err
	}
	arch := strings.TrimSpace(archOutput)

	// Normalize architecture names
	switch arch {
	case "x86_64":
		arch = "amd64"
	case "aarch64":
		arch = "arm64"
	case "armv7l":
		arch = "arm"
	}

	return platform, arch, nil
}

// downloadAndInstallBinary downloads and installs the sloth-runner binary
func downloadAndInstallBinary(client *ssh.Client, version, platform, arch string) error {
	// Construct download URL
	filename := fmt.Sprintf("sloth-runner_%s_%s_%s.tar.gz", version, platform, arch)
	url := fmt.Sprintf("https://github.com/chalkan3-sloth/sloth-runner/releases/download/%s/%s", version, filename)

	// Download and extract in one command
	installScript := fmt.Sprintf(`
set -e
cd /tmp
curl -sL "%s" -o sloth-runner.tar.gz
tar -xzf sloth-runner.tar.gz
chmod +x sloth-runner
mv -f sloth-runner /usr/local/bin/sloth-runner
rm -f sloth-runner.tar.gz
/usr/local/bin/sloth-runner version
`, url)

	_, err := runSSHCommand(client, installScript)
	return err
}

// createSystemdService creates the systemd service file
func createSystemdService(client *ssh.Client, agentName string, opts InstallOptions) error {
	// Determine report address
	reportAddr := opts.ReportAddress
	if reportAddr == "" {
		reportAddr = fmt.Sprintf("%s:%d", opts.SSHHost, opts.Port)
	}

	// Create systemd service content
	serviceContent := fmt.Sprintf(`[Unit]
Description=Sloth Runner Agent - %s
After=network.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/sloth-runner agent start \
  --name %s \
  --bind-address %s \
  --port %d \
  --master %s \
  --report-address %s \
  --daemon=false
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=sloth-runner-agent-%s

[Install]
WantedBy=multi-user.target
`, agentName, agentName, opts.BindAddress, opts.Port, opts.MasterAddr, reportAddr, agentName)

	// Create service file
	servicePath := fmt.Sprintf("/etc/systemd/system/sloth-runner-agent-%s.service", agentName)
	createServiceScript := fmt.Sprintf(`
cat > %s << 'SLOTH_SERVICE_EOF'
%s
SLOTH_SERVICE_EOF
chmod 644 %s
`, servicePath, serviceContent, servicePath)

	_, err := runSSHCommand(client, createServiceScript)
	return err
}

// enableAndStartService enables and starts the systemd service
func enableAndStartService(client *ssh.Client, agentName string) error {
	serviceName := fmt.Sprintf("sloth-runner-agent-%s", agentName)

	commands := []string{
		"systemctl daemon-reload",
		fmt.Sprintf("systemctl enable %s", serviceName),
		fmt.Sprintf("systemctl start %s", serviceName),
	}

	for _, cmd := range commands {
		if _, err := runSSHCommand(client, cmd); err != nil {
			return err
		}
	}

	return nil
}

// verifyAgentRunning verifies that the agent service is running
func verifyAgentRunning(client *ssh.Client, agentName string) error {
	serviceName := fmt.Sprintf("sloth-runner-agent-%s", agentName)
	output, err := runSSHCommand(client, fmt.Sprintf("systemctl is-active %s", serviceName))
	if err != nil {
		return err
	}

	if !strings.Contains(output, "active") {
		return fmt.Errorf("service is not active: %s", output)
	}

	return nil
}
