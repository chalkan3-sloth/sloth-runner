package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	githubAPIURL = "https://api.github.com/repos/chalkan3-sloth/sloth-runner/releases/latest"
	githubReleaseURL = "https://github.com/chalkan3-sloth/sloth-runner/releases/download"
)

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var agentUpdateCmd = &cobra.Command{
	Use:   "update [agent-name]",
	Short: "Update the sloth-runner agent to the latest version",
	Long: `Updates the sloth-runner agent binary to the latest version from GitHub releases.
If an agent name is provided, it will execute the update remotely on the agent host via SSH.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var agentName string
		if len(args) > 0 {
			agentName = args[0]
		}

		force, _ := cmd.Flags().GetBool("force")
		targetVersion, _ := cmd.Flags().GetString("version")
		skipRestart, _ := cmd.Flags().GetBool("skip-restart")

		// If agent name is provided, execute remotely
		if agentName != "" {
			return updateRemoteAgent(agentName, targetVersion, force, skipRestart)
		}

		// Otherwise, update local binary
		// Create a spinner for visual feedback
		spinner, _ := pterm.DefaultSpinner.Start("Checking for updates...")

		// Get current version
		currentVersion := version // From main.go
		spinner.UpdateText(fmt.Sprintf("Current version: %s", currentVersion))

		// Fetch latest release info
		release, err := fetchLatestRelease(targetVersion)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to fetch release info: %v", err))
			return err
		}

		latestVersion := release.TagName
		if latestVersion == currentVersion && !force {
			spinner.Success(fmt.Sprintf("Already running the latest version: %s", currentVersion))
			return nil
		}

		spinner.UpdateText(fmt.Sprintf("New version available: %s → %s", currentVersion, latestVersion))

		// Determine the asset to download based on OS and architecture
		osName := runtime.GOOS
		archName := runtime.GOARCH

		var downloadURL string
		var assetName string
		for _, asset := range release.Assets {
			// Match pattern: sloth-runner_v{VERSION}_{OS}_{ARCH}.tar.gz
			if strings.Contains(asset.Name, osName) &&
			   strings.Contains(asset.Name, archName) &&
			   strings.HasSuffix(asset.Name, ".tar.gz") {
				downloadURL = asset.BrowserDownloadURL
				assetName = asset.Name
				break
			}
		}

		if downloadURL == "" {
			spinner.Fail(fmt.Sprintf("No suitable release found for %s/%s", osName, archName))
			return fmt.Errorf("no suitable release found for %s/%s", osName, archName)
		}

		// Download the new binary
		spinner.UpdateText(fmt.Sprintf("Downloading %s...", assetName))
		tmpFile, err := downloadFile(downloadURL)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to download: %v", err))
			return err
		}
		defer os.Remove(tmpFile)

		// Extract if it's a tar.gz
		var binaryPath string
		if strings.HasSuffix(assetName, ".tar.gz") {
			spinner.UpdateText("Extracting archive...")
			binaryPath, err = extractBinary(tmpFile)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Failed to extract: %v", err))
				return err
			}
			defer os.Remove(binaryPath)
		} else {
			binaryPath = tmpFile
		}

		// Get the current executable path
		currentExe, err := os.Executable()
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to get current executable: %v", err))
			return err
		}

		// Backup current binary
		backupPath := currentExe + ".backup"
		spinner.UpdateText("Creating backup...")
		if err := copyFile(currentExe, backupPath); err != nil {
			spinner.Fail(fmt.Sprintf("Failed to create backup: %v", err))
			return err
		}

		// Replace the binary
		spinner.UpdateText("Installing new version...")
		if err := os.Rename(binaryPath, currentExe); err != nil {
			// Try to copy instead of rename (might be cross-device)
			if err := copyFile(binaryPath, currentExe); err != nil {
				// Restore backup
				copyFile(backupPath, currentExe)
				spinner.Fail(fmt.Sprintf("Failed to install new version: %v", err))
				return err
			}
		}

		// Make sure it's executable
		if err := os.Chmod(currentExe, 0755); err != nil {
			// Restore backup
			copyFile(backupPath, currentExe)
			spinner.Fail(fmt.Sprintf("Failed to set permissions: %v", err))
			return err
		}

		spinner.Success(fmt.Sprintf("Successfully updated to version %s", latestVersion))

		// Check if we need to restart the systemd service
		if agentName != "" && !skipRestart {
			if isSystemdService(agentName) {
				pterm.Info.Printf("Restarting systemd service for agent '%s'...\n", agentName)
				if err := restartSystemdService(agentName); err != nil {
					pterm.Warning.Printf("Failed to restart service: %v\n", err)
					pterm.Info.Println("You may need to restart the service manually:")
					pterm.Printf("  sudo systemctl restart sloth-runner-agent-%s\n", agentName)
				} else {
					pterm.Success.Printf("Service restarted successfully\n")
				}
			} else {
				pterm.Info.Printf("Agent '%s' is not running as a systemd service\n", agentName)
				pterm.Info.Println("Please restart the agent manually if it's running")
			}
		}

		// Clean up backup
		os.Remove(backupPath)

		return nil
	},
}

// fetchLatestRelease fetches the latest release information from GitHub
func fetchLatestRelease(version string) (*Release, error) {
	var url string
	if version != "" && version != "latest" {
		url = fmt.Sprintf("https://api.github.com/repos/chalkan3-sloth/sloth-runner/releases/tags/%s", version)
	} else {
		url = githubAPIURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release info: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// getAssetName returns the appropriate asset name for the current OS and architecture
func getAssetName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go arch names to common release naming
	archMap := map[string]string{
		"amd64": "x86_64",
		"386":   "i386",
		"arm64": "aarch64",
		"arm":   "armv7",
	}

	if mappedArch, ok := archMap[arch]; ok {
		arch = mappedArch
	}

	return fmt.Sprintf("sloth-runner-%s-%s.tar.gz", os, arch)
}

// downloadFile downloads a file from the given URL
func downloadFile(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "sloth-runner-update-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// extractBinary extracts the binary from a tar.gz archive
func extractBinary(archivePath string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	tmpDir, err := os.MkdirTemp("", "sloth-runner-extract-*")
	if err != nil {
		return "", err
	}

	var binaryPath string
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Look for the sloth-runner binary
		if strings.Contains(header.Name, "sloth-runner") && !strings.Contains(header.Name, ".") {
			binaryPath = filepath.Join(tmpDir, "sloth-runner")
			outFile, err := os.Create(binaryPath)
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", err
			}

			// Set executable permissions
			if err := os.Chmod(binaryPath, 0755); err != nil {
				return "", err
			}
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("sloth-runner binary not found in archive")
	}

	return binaryPath, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
}

// updateRemoteAgent updates the agent binary on a remote host via SSH
func updateRemoteAgent(agentName, targetVersion string, force, skipRestart bool) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Connecting to agent '%s'...", agentName))

	// Remote update script
	updateScript := fmt.Sprintf(`#!/bin/bash
set -e

# Download latest release
echo "Fetching latest release info..."
RELEASE_URL="https://api.github.com/repos/chalkan3-sloth/sloth-runner/releases/latest"
RELEASE_INFO=$(curl -s "$RELEASE_URL")
VERSION=$(echo "$RELEASE_INFO" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4)

echo "Latest version: $VERSION"

# Determine architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
esac

# Download URL
DOWNLOAD_URL="https://github.com/chalkan3-sloth/sloth-runner/releases/download/$VERSION/sloth-runner_${VERSION}_linux_${ARCH}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"
curl -L -o /tmp/sloth-runner.tar.gz "$DOWNLOAD_URL"

# Extract
echo "Extracting..."
cd /tmp
tar -xzf sloth-runner.tar.gz

# Backup and replace
echo "Installing..."
sudo cp /usr/local/bin/sloth-runner /usr/local/bin/sloth-runner.backup 2>/dev/null || true
sudo mv sloth-runner /usr/local/bin/sloth-runner
sudo chmod +x /usr/local/bin/sloth-runner

# Cleanup
rm -f /tmp/sloth-runner.tar.gz

echo "Update completed: $VERSION"

# Restart service if exists
if systemctl is-active --quiet sloth-runner-agent-%s; then
    echo "Restarting service..."
    sudo systemctl restart sloth-runner-agent-%s
    echo "Service restarted successfully"
fi
`, agentName, agentName)

	spinner.UpdateText(fmt.Sprintf("Executing update on '%s'...", agentName))

	// Create a wrapper script that will be piped to SSH
	wrapperScript := fmt.Sprintf("incus exec %s -- bash", agentName)

	// Execute via SSH to incus host, piping the update script
	cmd := exec.Command("ssh", "-p", "22", "chalkan3@192.168.1.16", wrapperScript)
	cmd.Stdin = strings.NewReader(updateScript)

	output, err := cmd.CombinedOutput()
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to update agent: %v", err))
		pterm.Error.Println(string(output))
		return err
	}

	spinner.Success(fmt.Sprintf("Agent '%s' updated successfully", agentName))
	pterm.Info.Println(string(output))

	return nil
}

// isSystemdService checks if an agent is running as a systemd service
func isSystemdService(agentName string) bool {
	serviceName := fmt.Sprintf("sloth-runner-agent-%s.service", agentName)
	cmd := exec.Command("systemctl", "status", serviceName)
	err := cmd.Run()
	return err == nil
}

// restartSystemdService restarts the systemd service for an agent
func restartSystemdService(agentName string) error {
	serviceName := fmt.Sprintf("sloth-runner-agent-%s.service", agentName)
	cmd := exec.Command("sudo", "systemctl", "restart", serviceName)
	return cmd.Run()
}

func init() {
	agentUpdateCmd.Flags().Bool("force", false, "Force update even if already on latest version")
	agentUpdateCmd.Flags().String("version", "", "Specific version to install (default: latest)")
	agentUpdateCmd.Flags().Bool("skip-restart", false, "Skip automatic service restart")
}