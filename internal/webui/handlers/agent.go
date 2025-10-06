package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AgentHandler handles agent operations
type AgentHandler struct {
	db *AgentDBWrapper
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(db *AgentDBWrapper) *AgentHandler {
	return &AgentHandler{db: db}
}

// AgentMetrics represents detailed agent metrics
type AgentMetrics struct {
	CPUPercent      float64                `json:"cpu_percent"`
	MemoryPercent   float64                `json:"memory_percent"`
	MemoryUsedBytes uint64                 `json:"memory_used_bytes"`
	MemoryTotalBytes uint64                `json:"memory_total_bytes"`
	DiskPercent     float64                `json:"disk_percent"`
	DiskUsedBytes   uint64                 `json:"disk_used_bytes"`
	DiskTotalBytes  uint64                 `json:"disk_total_bytes"`
	ProcessCount    int                    `json:"process_count"`
	LoadAvg1Min     float64                `json:"load_avg_1min"`
	LoadAvg5Min     float64                `json:"load_avg_5min"`
	LoadAvg15Min    float64                `json:"load_avg_15min"`
	UptimeSeconds   uint64                 `json:"uptime_seconds"`
	NetworkInterfaces []NetworkInterface   `json:"network_interfaces,omitempty"`
	DiskPartitions  []DiskPartition        `json:"disk_partitions,omitempty"`
	Processes       []ProcessInfo          `json:"processes,omitempty"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name        string   `json:"name"`
	IPAddresses []string `json:"ip_addresses"`
	MACAddress  string   `json:"mac_address"`
	BytesSent   uint64   `json:"bytes_sent"`
	BytesRecv   uint64   `json:"bytes_recv"`
	IsUp        bool     `json:"is_up"`
}

// DiskPartition represents a disk partition
type DiskPartition struct {
	Device     string  `json:"device"`
	Mountpoint string  `json:"mountpoint"`
	FSType     string  `json:"fstype"`
	TotalBytes uint64  `json:"total_bytes"`
	UsedBytes  uint64  `json:"used_bytes"`
	FreeBytes  uint64  `json:"free_bytes"`
	Percent    float64 `json:"percent"`
}

// ProcessInfo represents a running process
type ProcessInfo struct {
	PID           int     `json:"pid"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryBytes   uint64  `json:"memory_bytes"`
	User          string  `json:"user"`
	Command       string  `json:"command"`
	StartedAt     int64   `json:"started_at"`
}

// List returns all agents
func (h *AgentHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	agents, err := h.db.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// Get returns an agent by name
func (h *AgentHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// Delete removes an agent
func (h *AgentHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if err := h.db.DeleteAgent(ctx, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

// GetStats returns agent statistics
func (h *AgentHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.db.GetStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Update updates an agent (placeholder for future implementation)
func (h *AgentHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// GetResourceUsage returns detailed resource usage for an agent
func (h *AgentHandler) GetResourceUsage(c *gin.Context) {
	name := c.Param("name")

	// TODO: Get real metrics from agent via gRPC
	// For now, return mock data
	metrics := generateMockMetrics(name)

	c.JSON(http.StatusOK, metrics)
}

// GetProcessList returns list of processes running on agent
func (h *AgentHandler) GetProcessList(c *gin.Context) {
	name := c.Param("name")

	// TODO: Get real process list from agent via gRPC
	processes := generateMockProcesses(name)

	c.JSON(http.StatusOK, gin.H{
		"processes":   processes,
		"total_count": len(processes),
	})
}

// GetNetworkInfo returns network information for an agent
func (h *AgentHandler) GetNetworkInfo(c *gin.Context) {
	name := c.Param("name")

	// TODO: Get real network info from agent via gRPC
	interfaces := generateMockNetworkInterfaces(name)

	c.JSON(http.StatusOK, gin.H{
		"interfaces": interfaces,
		"hostname":   name,
	})
}

// GetDiskInfo returns disk information for an agent
func (h *AgentHandler) GetDiskInfo(c *gin.Context) {
	name := c.Param("name")

	// TODO: Get real disk info from agent via gRPC
	partitions := generateMockDiskPartitions(name)

	c.JSON(http.StatusOK, gin.H{
		"partitions":          partitions,
		"total_io_read_bytes": uint64(rand.Int63n(1000000000000)),
		"total_io_write_bytes": uint64(rand.Int63n(500000000000)),
	})
}

// ExecuteCommand executes a command on an agent
func (h *AgentHandler) ExecuteCommand(c *gin.Context) {
	name := c.Param("name")

	var req struct {
		Command string `json:"command"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: Execute command via gRPC and stream output
	// For now, return mock output
	c.Header("Content-Type", "text/plain")
	c.Header("Transfer-Encoding", "chunked")

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Simulate command execution
	fmt.Fprintf(w, "$ %s\n", req.Command)
	flusher.Flush()

	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(w, "Executing on agent: %s\n", name)
	flusher.Flush()

	time.Sleep(200 * time.Millisecond)
	fmt.Fprintf(w, "Output: Command executed successfully\n")
	flusher.Flush()

	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(w, "Exit code: 0\n")
	flusher.Flush()
}

// RestartAgent restarts an agent
func (h *AgentHandler) RestartAgent(c *gin.Context) {
	name := c.Param("name")

	// TODO: Restart agent via gRPC
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Agent %s restart initiated", name),
	})
}

// ShutdownAgent shuts down an agent
func (h *AgentHandler) ShutdownAgent(c *gin.Context) {
	name := c.Param("name")

	// TODO: Shutdown agent via gRPC
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Agent %s shutdown initiated", name),
	})
}

// StreamLogs streams logs from an agent
func (h *AgentHandler) StreamLogs(c *gin.Context) {
	name := c.Param("name")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// TODO: Stream real logs from agent
	// For now, generate mock log entries
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	logLines := []string{
		fmt.Sprintf("[INFO] Agent %s - System check completed", name),
		fmt.Sprintf("[DEBUG] Agent %s - Heartbeat sent to master", name),
		fmt.Sprintf("[INFO] Agent %s - Resource monitoring active", name),
		fmt.Sprintf("[DEBUG] Agent %s - Task queue: 0 pending", name),
		fmt.Sprintf("[INFO] Agent %s - All systems operational", name),
	}

	for i, line := range logLines {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			logEntry := map[string]interface{}{
				"timestamp": time.Now().Unix(),
				"level":     "INFO",
				"message":   line,
				"source":    name,
			}
			data, _ := json.Marshal(logEntry)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			if i >= len(logLines)-1 {
				return
			}
		}
	}
}

// StreamMetrics streams metrics from an agent
func (h *AgentHandler) StreamMetrics(c *gin.Context) {
	name := c.Param("name")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// TODO: Stream real metrics from agent
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			metrics := generateMockMetrics(name)
			metricsData := map[string]interface{}{
				"timestamp":      time.Now().Unix(),
				"cpu_percent":    metrics.CPUPercent,
				"memory_percent": metrics.MemoryPercent,
				"disk_percent":   metrics.DiskPercent,
			}
			data, _ := json.Marshal(metricsData)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// BulkExecute executes a command on multiple agents
func (h *AgentHandler) BulkExecute(c *gin.Context) {
	var req struct {
		AgentNames []string `json:"agent_names"`
		GroupName  string   `json:"group_name"`
		Command    string   `json:"command"`
		Parallel   bool     `json:"parallel"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Transfer-Encoding", "chunked")

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// TODO: Execute on real agents
	for _, agentName := range req.AgentNames {
		time.Sleep(100 * time.Millisecond)

		result := map[string]interface{}{
			"agent_name":       agentName,
			"success":          true,
			"output":           fmt.Sprintf("Command executed on %s", agentName),
			"exit_code":        0,
			"execution_time_ms": rand.Intn(500) + 100,
		}

		data, _ := json.Marshal(result)
		fmt.Fprintf(w, "%s\n", data)
		flusher.Flush()
	}
}

// GetMultipleStatus returns status of multiple agents
func (h *AgentHandler) GetMultipleStatus(c *gin.Context) {
	var req struct {
		AgentNames []string `json:"agent_names"`
		GroupName  string   `json:"group_name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: Get real status from agents
	statuses := make([]map[string]interface{}, 0)

	for _, name := range req.AgentNames {
		metrics := generateMockMetrics(name)
		status := map[string]interface{}{
			"agent_name":     name,
			"status":         "connected",
			"cpu_percent":    metrics.CPUPercent,
			"memory_percent": metrics.MemoryPercent,
			"last_heartbeat": time.Now().Unix(),
			"healthy":        true,
		}
		statuses = append(statuses, status)
	}

	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}

// Mock data generators
func generateMockMetrics(agentName string) *AgentMetrics {
	rand.Seed(time.Now().UnixNano() + int64(len(agentName)))

	totalMemory := uint64(16 * 1024 * 1024 * 1024) // 16 GB
	usedMemory := uint64(float64(totalMemory) * (rand.Float64()*0.5 + 0.2))

	totalDisk := uint64(500 * 1024 * 1024 * 1024) // 500 GB
	usedDisk := uint64(float64(totalDisk) * (rand.Float64()*0.6 + 0.2))

	return &AgentMetrics{
		CPUPercent:       rand.Float64()*60 + 10,
		MemoryPercent:    float64(usedMemory) / float64(totalMemory) * 100,
		MemoryUsedBytes:  usedMemory,
		MemoryTotalBytes: totalMemory,
		DiskPercent:      float64(usedDisk) / float64(totalDisk) * 100,
		DiskUsedBytes:    usedDisk,
		DiskTotalBytes:   totalDisk,
		ProcessCount:     rand.Intn(200) + 50,
		LoadAvg1Min:      rand.Float64()*2 + 0.5,
		LoadAvg5Min:      rand.Float64()*2 + 0.3,
		LoadAvg15Min:     rand.Float64()*2 + 0.2,
		UptimeSeconds:    uint64(rand.Intn(86400*30) + 86400),
	}
}

func generateMockProcesses(agentName string) []ProcessInfo {
	processes := []ProcessInfo{
		{PID: 1, Name: "systemd", Status: "running", CPUPercent: 0.1, MemoryPercent: 0.5, MemoryBytes: 80000000, User: "root", Command: "/sbin/systemd", StartedAt: time.Now().Add(-24 * time.Hour).Unix()},
		{PID: 123, Name: "sloth-runner", Status: "running", CPUPercent: 2.5, MemoryPercent: 3.2, MemoryBytes: 512000000, User: "sloth", Command: "/usr/bin/sloth-runner agent start", StartedAt: time.Now().Add(-6 * time.Hour).Unix()},
		{PID: 456, Name: "nginx", Status: "running", CPUPercent: 1.2, MemoryPercent: 1.8, MemoryBytes: 288000000, User: "www-data", Command: "nginx: master process", StartedAt: time.Now().Add(-12 * time.Hour).Unix()},
		{PID: 789, Name: "postgres", Status: "running", CPUPercent: 3.1, MemoryPercent: 5.4, MemoryBytes: 864000000, User: "postgres", Command: "/usr/lib/postgresql/14/bin/postgres", StartedAt: time.Now().Add(-48 * time.Hour).Unix()},
		{PID: 1011, Name: "docker", Status: "running", CPUPercent: 0.8, MemoryPercent: 2.1, MemoryBytes: 336000000, User: "root", Command: "/usr/bin/dockerd", StartedAt: time.Now().Add(-72 * time.Hour).Unix()},
	}

	return processes
}

func generateMockNetworkInterfaces(agentName string) []NetworkInterface {
	return []NetworkInterface{
		{
			Name:        "eth0",
			IPAddresses: []string{"192.168.1.100", "fe80::1"},
			MACAddress:  "00:1a:2b:3c:4d:5e",
			BytesSent:   uint64(rand.Int63n(10000000000)),
			BytesRecv:   uint64(rand.Int63n(20000000000)),
			IsUp:        true,
		},
		{
			Name:        "lo",
			IPAddresses: []string{"127.0.0.1", "::1"},
			MACAddress:  "00:00:00:00:00:00",
			BytesSent:   uint64(rand.Int63n(1000000000)),
			BytesRecv:   uint64(rand.Int63n(1000000000)),
			IsUp:        true,
		},
	}
}

func generateMockDiskPartitions(agentName string) []DiskPartition {
	return []DiskPartition{
		{
			Device:     "/dev/sda1",
			Mountpoint: "/",
			FSType:     "ext4",
			TotalBytes: 500 * 1024 * 1024 * 1024,
			UsedBytes:  250 * 1024 * 1024 * 1024,
			FreeBytes:  250 * 1024 * 1024 * 1024,
			Percent:    50.0,
		},
		{
			Device:     "/dev/sda2",
			Mountpoint: "/home",
			FSType:     "ext4",
			TotalBytes: 1000 * 1024 * 1024 * 1024,
			UsedBytes:  600 * 1024 * 1024 * 1024,
			FreeBytes:  400 * 1024 * 1024 * 1024,
			Percent:    60.0,
		},
	}
}
