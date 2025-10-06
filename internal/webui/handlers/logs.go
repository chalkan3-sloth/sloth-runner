package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// LogsHandler handles log viewing and streaming
type LogsHandler struct {
	wsHub *WebSocketHub
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(wsHub *WebSocketHub) *LogsHandler {
	return &LogsHandler{
		wsHub: wsHub,
	}
}

// ListLogFiles returns available log files
func (h *LogsHandler) ListLogFiles(c *gin.Context) {
	logDir := "/var/log/sloth-runner"

	// Fallback to local logs if system logs not accessible
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		homeDir, _ := os.UserHomeDir()
		logDir = filepath.Join(homeDir, ".sloth-runner", "logs")
	}

	files, err := os.ReadDir(logDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"files": []string{}})
		return
	}

	logFiles := make([]gin.H, 0)
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".log") || strings.HasSuffix(file.Name(), ".txt")) {
			info, _ := file.Info()
			logFiles = append(logFiles, gin.H{
				"name":         file.Name(),
				"size":         info.Size(),
				"modified_at":  info.ModTime().Unix(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"files": logFiles})
}

// GetLogFile returns log file content
func (h *LogsHandler) GetLogFile(c *gin.Context) {
	filename := c.Param("filename")

	// Validate filename (prevent directory traversal)
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	logDir := "/var/log/sloth-runner"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		homeDir, _ := os.UserHomeDir()
		logDir = filepath.Join(homeDir, ".sloth-runner", "logs")
	}

	filePath := filepath.Join(logDir, filename)

	// Get parameters
	tail := c.DefaultQuery("tail", "100")
	search := c.Query("search")

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log file not found"})
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Filter by search term if provided
		if search != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(search)) {
			continue
		}

		lines = append(lines, line)
	}

	// Return tail lines if requested
	if tail != "all" {
		// Parse tail count
		var tailCount int
		if _, err := fmt.Sscanf(tail, "%d", &tailCount); err == nil {
			if len(lines) > tailCount {
				lines = lines[len(lines)-tailCount:]
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"lines":    lines,
		"total":    len(lines),
	})
}

// StreamLogs streams log updates via Server-Sent Events
func (h *LogsHandler) StreamLogs(c *gin.Context) {
	filename := c.Param("filename")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// TODO: Implement actual log tailing for file: filename
	c.String(http.StatusOK, fmt.Sprintf("data: Log streaming for %s not yet implemented\n\n", filename))
}
