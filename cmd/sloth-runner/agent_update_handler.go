package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	pb "github.com/chalkan3-sloth/sloth-runner/proto"
)

// UpdateAgent handles the agent update request
func (s *agentServer) UpdateAgent(ctx context.Context, in *pb.UpdateAgentRequest) (*pb.UpdateAgentResponse, error) {
	slog.Info("Received agent update request", "target_version", in.GetTargetVersion(), "force", in.GetForce())

	// Get current version
	currentVersion := version
	oldVersion := currentVersion

	// Fetch latest release info
	targetVer := in.GetTargetVersion()
	if targetVer == "" {
		targetVer = "latest"
	}

	release, err := fetchLatestReleaseInfo(targetVer)
	if err != nil {
		slog.Error("Failed to fetch release info", "error", err)
		return &pb.UpdateAgentResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to fetch release info: %v", err),
		}, nil
	}

	latestVersion := release.TagName

	// Check if update is needed
	if latestVersion == currentVersion && !in.GetForce() {
		slog.Info("Already running latest version", "version", currentVersion)
		return &pb.UpdateAgentResponse{
			Success:    true,
			Message:    "Already running the latest version",
			OldVersion: currentVersion,
			NewVersion: currentVersion,
		}, nil
	}

	slog.Info("Updating agent", "from", currentVersion, "to", latestVersion)

	// Download and install
	if err := downloadAndInstallUpdate(release); err != nil {
		slog.Error("Failed to update agent", "error", err)
		return &pb.UpdateAgentResponse{
			Success: false,
			Message: fmt.Sprintf("Update failed: %v", err),
		}, nil
	}

	// Schedule agent shutdown to allow update script to replace binary
	if !in.GetSkipRestart() {
		go func() {
			time.Sleep(3 * time.Second)
			slog.Info("Shutting down agent for update...")
			os.Exit(0)
		}()
	}

	return &pb.UpdateAgentResponse{
		Success:    true,
		Message:    "Agent update prepared. Shutting down for binary replacement and restart.",
		OldVersion: oldVersion,
		NewVersion: latestVersion,
	}, nil
}

// Release represents a GitHub release (duplicate from agent_update.go but needed here)
type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func fetchLatestReleaseInfo(version string) (*ReleaseInfo, error) {
	var url string
	if version != "" && version != "latest" {
		url = fmt.Sprintf("https://api.github.com/repos/chalkan3-sloth/sloth-runner/releases/tags/%s", version)
	} else {
		url = "https://api.github.com/repos/chalkan3-sloth/sloth-runner/releases/latest"
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

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadAndInstallUpdate(release *ReleaseInfo) error {
	// Determine OS and architecture
	osName := runtime.GOOS
	archName := runtime.GOARCH

	// Find the appropriate asset
	var downloadURL string
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, osName) &&
			strings.Contains(asset.Name, archName) &&
			strings.HasSuffix(asset.Name, ".tar.gz") {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no suitable release found for %s/%s", osName, archName)
	}

	slog.Info("Downloading update", "url", downloadURL)

	// Download
	tmpFile := "/tmp/sloth-runner-update.tar.gz"
	if err := downloadFileToPath(downloadURL, tmpFile); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(tmpFile)

	// Extract
	slog.Info("Extracting update")
	extractedBinary := "/tmp/sloth-runner-new"
	if err := extractBinaryFromTar(tmpFile, extractedBinary); err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}
	defer os.Remove(extractedBinary)

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %w", err)
	}

	// Backup current binary
	backupPath := currentExe + ".backup"
	slog.Info("Creating backup", "path", backupPath)
	if err := copyFileSimple(currentExe, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Create update script that will run after agent stops
	updateScript := fmt.Sprintf(`#!/bin/bash
# Wait for agent to stop
sleep 2

# Replace binary
cp %s %s.backup
mv %s %s
chmod +x %s

# Restart service or agent
if systemctl is-active --quiet sloth-runner-agent 2>/dev/null; then
    systemctl restart sloth-runner-agent
elif systemctl is-active --quiet sloth-agent 2>/dev/null; then
    systemctl restart sloth-agent
else
    # Try to restart via nohup if no systemd service
    cd /home/igor && nohup %s agent start --name lady-arch --master 192.168.1.29:50053 --port 50051 --bind-address 0.0.0.0 --report-address 192.168.1.16:50052 --telemetry --metrics-port 9090 > agent.log 2>&1 &
fi

# Cleanup
rm -f %s %s.backup /tmp/agent-update.sh
`, currentExe, currentExe, extractedBinary, currentExe, currentExe, currentExe, extractedBinary, currentExe)

	updateScriptPath := "/tmp/agent-update.sh"
	if err := os.WriteFile(updateScriptPath, []byte(updateScript), 0755); err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Run update script in background and shutdown agent
	go func() {
		cmd := exec.Command("bash", updateScriptPath)
		cmd.Start()
	}()

	slog.Info("Update prepared. Agent will shutdown and restart with new version.")
	return nil
}

func downloadFileToPath(url, dest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinaryFromTar(tarPath, dest string) error {
	// For simplicity, use tar command
	cmd := exec.Command("tar", "-xzf", tarPath, "-C", "/tmp")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Find the binary
	binary := "/tmp/sloth-runner"
	if _, err := os.Stat(binary); os.IsNotExist(err) {
		return fmt.Errorf("binary not found in archive")
	}

	// Move to dest
	return os.Rename(binary, dest)
}

func copyFileSimple(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0755)
}

func restartAgentService() {
	slog.Info("Attempting to restart agent service")

	// Try to restart via systemctl
	cmd := exec.Command("systemctl", "restart", "sloth-runner-agent")
	if err := cmd.Run(); err != nil {
		slog.Warn("Failed to restart via systemctl", "error", err)

		// Try alternative service names
		for _, serviceName := range []string{"sloth-runner-agent-*", "sloth-agent"} {
			cmd = exec.Command("systemctl", "restart", serviceName)
			if err := cmd.Run(); err == nil {
				slog.Info("Restarted service", "name", serviceName)
				return
			}
		}

		slog.Warn("Could not restart service automatically")
	} else {
		slog.Info("Service restarted successfully")
	}
}
