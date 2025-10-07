package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
	"github.com/gin-gonic/gin"
)

// AgentHandler handles agent operations
type AgentHandler struct {
	db          *AgentDBWrapper
	agentClient *services.AgentClient
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(db *AgentDBWrapper, agentClient *services.AgentClient) *AgentHandler {
	return &AgentHandler{
		db:          db,
		agentClient: agentClient,
	}
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

// List returns all agents with enriched metrics
func (h *AgentHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	agents, err := h.db.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enrich agents with metrics from system_info
	enrichedAgents := make([]map[string]interface{}, 0, len(agents))
	for _, agent := range agents {
		enriched := map[string]interface{}{
			"id":                 agent.ID,
			"name":               agent.Name,
			"address":            agent.Address,
			"status":             agent.Status,
			"last_heartbeat":     agent.LastHeartbeat,
			"registered_at":      agent.RegisteredAt,
			"updated_at":         agent.UpdatedAt,
			"last_info_collected": agent.LastInfoCollected,
			"version":            agent.Version,
		}

		// Parse system_info JSON to extract metrics
		if agent.SystemInfo != "" {
			var sysInfo map[string]interface{}
			if err := json.Unmarshal([]byte(agent.SystemInfo), &sysInfo); err == nil {
				// Extract memory metrics
				if memory, ok := sysInfo["memory"].(map[string]interface{}); ok {
					if usedPercent, ok := memory["used_percent"].(float64); ok {
						enriched["memory_percent"] = usedPercent
					}
				}

				// Extract CPU count
				if cpus, ok := sysInfo["cpus"].(float64); ok {
					enriched["cpu_count"] = int(cpus)
				}

				// Extract disk metrics (use first disk)
				if disks, ok := sysInfo["disk"].([]interface{}); ok && len(disks) > 0 {
					if disk, ok := disks[0].(map[string]interface{}); ok {
						if usedPercent, ok := disk["used_percent"].(float64); ok {
							enriched["disk_percent"] = usedPercent
						}
					}
				}

				// Extract load average
				if loadAvg, ok := sysInfo["load_average"].([]interface{}); ok && len(loadAvg) > 0 {
					if load1, ok := loadAvg[0].(float64); ok {
						enriched["load_avg"] = load1
					}
				}

				// Extract uptime
				if uptime, ok := sysInfo["uptime"].(float64); ok {
					enriched["uptime"] = int64(uptime)
				}

				// Extract hostname
				if hostname, ok := sysInfo["hostname"].(string); ok {
					enriched["hostname"] = hostname
				}
			}
		}

		enrichedAgents = append(enrichedAgents, enriched)
	}

	c.JSON(http.StatusOK, gin.H{"agents": enrichedAgents})
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

// GetMetricsHistory returns historical metrics for an agent
func (h *AgentHandler) GetMetricsHistory(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	limit := 50 // Last 50 data points
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}

	history, err := h.db.GetMetricsHistory(ctx, name, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// GetResourceUsage returns detailed resource usage for an agent
func (h *AgentHandler) GetResourceUsage(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get real resource usage via gRPC
	resp, err := h.agentClient.GetResourceUsage(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get resource usage: %v", err)})
		return
	}

	metrics := &AgentMetrics{
		CPUPercent:       resp.CpuPercent,
		MemoryPercent:    resp.MemoryPercent,
		MemoryUsedBytes:  resp.MemoryUsedBytes,
		MemoryTotalBytes: resp.MemoryTotalBytes,
		DiskPercent:      resp.DiskPercent,
		DiskUsedBytes:    resp.DiskUsedBytes,
		DiskTotalBytes:   resp.DiskTotalBytes,
		ProcessCount:     int(resp.ProcessCount),
		LoadAvg1Min:      resp.LoadAvg_1Min,
		LoadAvg5Min:      resp.LoadAvg_5Min,
		LoadAvg15Min:     resp.LoadAvg_15Min,
		UptimeSeconds:    resp.UptimeSeconds,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetProcessList returns list of processes running on agent
func (h *AgentHandler) GetProcessList(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get real process list via gRPC
	resp, err := h.agentClient.GetProcessList(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get process list: %v", err)})
		return
	}

	processes := make([]ProcessInfo, 0, len(resp.Processes))
	for _, p := range resp.Processes {
		processes = append(processes, ProcessInfo{
			PID:           int(p.Pid),
			Name:          p.Name,
			Status:        p.Status,
			CPUPercent:    p.CpuPercent,
			MemoryPercent: p.MemoryPercent,
			MemoryBytes:   p.MemoryBytes,
			User:          p.User,
			Command:       p.Command,
			StartedAt:     p.StartedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"processes":   processes,
		"total_count": len(processes),
	})
}

// GetNetworkInfo returns network information for an agent
func (h *AgentHandler) GetNetworkInfo(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get real network info via gRPC
	resp, err := h.agentClient.GetNetworkInfo(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get network info: %v", err)})
		return
	}

	interfaces := make([]NetworkInterface, 0, len(resp.Interfaces))
	for _, iface := range resp.Interfaces {
		interfaces = append(interfaces, NetworkInterface{
			Name:        iface.Name,
			IPAddresses: iface.IpAddresses,
			MACAddress:  iface.MacAddress,
			BytesSent:   iface.BytesSent,
			BytesRecv:   iface.BytesRecv,
			IsUp:        iface.IsUp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"interfaces": interfaces,
		"hostname":   resp.Hostname,
	})
}

// GetDiskInfo returns disk information for an agent
func (h *AgentHandler) GetDiskInfo(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get real disk info via gRPC
	resp, err := h.agentClient.GetDiskInfo(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get disk info: %v", err)})
		return
	}

	partitions := make([]DiskPartition, 0, len(resp.Partitions))
	for _, part := range resp.Partitions {
		partitions = append(partitions, DiskPartition{
			Device:     part.Device,
			Mountpoint: part.Mountpoint,
			FSType:     part.Fstype,
			TotalBytes: part.TotalBytes,
			UsedBytes:  part.UsedBytes,
			FreeBytes:  part.FreeBytes,
			Percent:    part.Percent,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"partitions":           partitions,
		"total_io_read_bytes":  resp.TotalIoReadBytes,
		"total_io_write_bytes": resp.TotalIoWriteBytes,
	})
}

// ExecuteCommand executes a command on an agent
func (h *AgentHandler) ExecuteCommand(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	var req struct {
		Command string `json:"command"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Execute command via gRPC and stream output
	stream, err := h.agentClient.RunCommand(ctx, agent.Address, req.Command)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to execute command: %v\n", err)
		return
	}

	c.Header("Content-Type", "text/plain")
	c.Header("Transfer-Encoding", "chunked")

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Stream output from agent
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "\nError: %v\n", err)
			flusher.Flush()
			break
		}

		if resp.StdoutChunk != "" {
			fmt.Fprint(w, resp.StdoutChunk)
			flusher.Flush()
		}
		if resp.StderrChunk != "" {
			fmt.Fprint(w, resp.StderrChunk)
			flusher.Flush()
		}
		if resp.Finished {
			fmt.Fprintf(w, "\nExit code: %d\n", resp.ExitCode)
			flusher.Flush()
			break
		}
	}
}

// RestartAgent restarts an agent
func (h *AgentHandler) RestartAgent(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Restart agent via gRPC
	err = h.agentClient.RestartService(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to restart agent: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Agent %s restart initiated", name),
	})
}

// ShutdownAgent shuts down an agent
func (h *AgentHandler) ShutdownAgent(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Shutdown agent via gRPC
	err = h.agentClient.Shutdown(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to shutdown agent: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Agent %s shutdown initiated", name),
	})
}

// StreamLogs streams logs from an agent
func (h *AgentHandler) StreamLogs(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Stream real logs from agent via gRPC
	stream, err := h.agentClient.StreamLogs(ctx, agent.Address)
	if err != nil {
		logEntry := map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"level":     "ERROR",
			"message":   fmt.Sprintf("Failed to stream logs: %v", err),
			"source":    name,
		}
		data, _ := json.Marshal(logEntry)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			logEntry, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errorEntry := map[string]interface{}{
					"timestamp": time.Now().Unix(),
					"level":     "ERROR",
					"message":   fmt.Sprintf("Stream error: %v", err),
					"source":    name,
				}
				data, _ := json.Marshal(errorEntry)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
				return
			}

			entry := map[string]interface{}{
				"timestamp": logEntry.Timestamp,
				"level":     logEntry.Level,
				"message":   logEntry.Message,
				"source":    name,
			}
			data, _ := json.Marshal(entry)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// StreamMetrics streams metrics from an agent
func (h *AgentHandler) StreamMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Stream real metrics from agent via gRPC
	stream, err := h.agentClient.StreamMetrics(ctx, agent.Address)
	if err != nil {
		errorData := map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"error":     fmt.Sprintf("Failed to stream metrics: %v", err),
		}
		data, _ := json.Marshal(errorData)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricsData, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errorData := map[string]interface{}{
					"timestamp": time.Now().Unix(),
					"error":     fmt.Sprintf("Stream error: %v", err),
				}
				data, _ := json.Marshal(errorData)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
				return
			}

			metrics := map[string]interface{}{
				"timestamp":      metricsData.Timestamp,
				"cpu_percent":    metricsData.CpuPercent,
				"memory_percent": metricsData.MemoryPercent,
				"disk_percent":   metricsData.DiskPercent,
			}
			data, _ := json.Marshal(metrics)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// BulkExecute executes a command on multiple agents
func (h *AgentHandler) BulkExecute(c *gin.Context) {
	ctx := c.Request.Context()

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

	// Execute on real agents
	executeOnAgent := func(agentName string) {
		startTime := time.Now()

		// Get agent from database
		agent, err := h.db.GetAgent(ctx, agentName)
		if err != nil {
			result := map[string]interface{}{
				"agent_name":        agentName,
				"success":           false,
				"error":             fmt.Sprintf("Agent not found: %v", err),
				"execution_time_ms": time.Since(startTime).Milliseconds(),
			}
			data, _ := json.Marshal(result)
			fmt.Fprintf(w, "%s\n", data)
			flusher.Flush()
			return
		}

		// Execute command via gRPC
		stream, err := h.agentClient.RunCommand(ctx, agent.Address, req.Command)
		if err != nil {
			result := map[string]interface{}{
				"agent_name":        agentName,
				"success":           false,
				"error":             fmt.Sprintf("Failed to execute: %v", err),
				"execution_time_ms": time.Since(startTime).Milliseconds(),
			}
			data, _ := json.Marshal(result)
			fmt.Fprintf(w, "%s\n", data)
			flusher.Flush()
			return
		}

		// Collect output
		var output string
		var exitCode int32
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				result := map[string]interface{}{
					"agent_name":        agentName,
					"success":           false,
					"error":             fmt.Sprintf("Stream error: %v", err),
					"execution_time_ms": time.Since(startTime).Milliseconds(),
				}
				data, _ := json.Marshal(result)
				fmt.Fprintf(w, "%s\n", data)
				flusher.Flush()
				return
			}

			output += resp.StdoutChunk + resp.StderrChunk
			if resp.Finished {
				exitCode = resp.ExitCode
				break
			}
		}

		result := map[string]interface{}{
			"agent_name":        agentName,
			"success":           exitCode == 0,
			"output":            output,
			"exit_code":         exitCode,
			"execution_time_ms": time.Since(startTime).Milliseconds(),
		}
		data, _ := json.Marshal(result)
		fmt.Fprintf(w, "%s\n", data)
		flusher.Flush()
	}

	if req.Parallel {
		// Execute in parallel
		var wg sync.WaitGroup
		for _, agentName := range req.AgentNames {
			wg.Add(1)
			go func(name string) {
				defer wg.Done()
				executeOnAgent(name)
			}(agentName)
		}
		wg.Wait()
	} else {
		// Execute sequentially
		for _, agentName := range req.AgentNames {
			executeOnAgent(agentName)
		}
	}
}

// GetMultipleStatus returns status of multiple agents
func (h *AgentHandler) GetMultipleStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		AgentNames []string `json:"agent_names"`
		GroupName  string   `json:"group_name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	statuses := make([]map[string]interface{}, 0, len(req.AgentNames))

	for _, name := range req.AgentNames {
		// Get agent from database
		agent, err := h.db.GetAgent(ctx, name)
		if err != nil {
			statuses = append(statuses, map[string]interface{}{
				"agent_name": name,
				"status":     "unknown",
				"error":      "Agent not found",
				"healthy":    false,
			})
			continue
		}

		// Get real metrics from agent
		metrics, err := h.agentClient.GetResourceUsage(ctx, agent.Address)
		if err != nil {
			statuses = append(statuses, map[string]interface{}{
				"agent_name":     name,
				"status":         "disconnected",
				"error":          err.Error(),
				"last_heartbeat": agent.LastHeartbeat,
				"healthy":        false,
			})
			continue
		}

		status := map[string]interface{}{
			"agent_name":     name,
			"status":         "connected",
			"cpu_percent":    metrics.CpuPercent,
			"memory_percent": metrics.MemoryPercent,
			"last_heartbeat": agent.LastHeartbeat,
			"healthy":        true,
		}
		statuses = append(statuses, status)
	}

	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}
