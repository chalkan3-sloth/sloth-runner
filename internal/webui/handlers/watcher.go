package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
	"github.com/gin-gonic/gin"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
)

// WatcherHandler handles watcher operations
type WatcherHandler struct {
	agentDB     *AgentDBWrapper
	agentClient *services.AgentClient
}

// NewWatcherHandler creates a new watcher handler
func NewWatcherHandler(agentDB *AgentDBWrapper, agentClient *services.AgentClient) *WatcherHandler {
	return &WatcherHandler{
		agentDB:     agentDB,
		agentClient: agentClient,
	}
}

// WatcherRequest represents a watcher creation/update request
type WatcherRequest struct {
	Type       string   `json:"type" binding:"required"`
	Conditions []string `json:"conditions" binding:"required"`
	Interval   string   `json:"interval" binding:"required"`

	// File watcher specific
	FilePath  string `json:"file_path,omitempty"`
	CheckHash bool   `json:"check_hash,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`

	// Process watcher specific
	ProcessName string `json:"process_name,omitempty"`
	PID         int32  `json:"pid,omitempty"`

	// Port watcher specific
	Port     int32  `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`

	// Threshold watchers
	CPUThreshold    float64 `json:"cpu_threshold,omitempty"`
	MemoryThreshold float64 `json:"memory_threshold,omitempty"`
	DiskThreshold   float64 `json:"disk_threshold,omitempty"`
}

// ListByAgent lists all watchers for an agent
func (h *WatcherHandler) ListByAgent(c *gin.Context) {
	agentName := c.Param("agent")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent parameter is required"})
		return
	}

	// Get agent info
	ctx := context.Background()
	agent, err := h.agentDB.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get agent client
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to connect to agent: %v", err)})
		return
	}

	// List watchers
	resp, err := client.ListWatchers(context.Background(), &pb.ListWatchersRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list watchers: %v", err)})
		return
	}

	// Convert to response format
	watchers := make([]map[string]interface{}, 0, len(resp.Watchers))
	for _, w := range resp.Watchers {
		watcher := map[string]interface{}{
			"id":         w.Id,
			"type":       w.Type,
			"conditions": w.Conditions,
			"interval":   w.Interval,
		}

		// Add type-specific fields
		if w.FilePath != "" {
			watcher["file_path"] = w.FilePath
			watcher["check_hash"] = w.CheckHash
			watcher["recursive"] = w.Recursive
		}
		if w.ProcessName != "" {
			watcher["process_name"] = w.ProcessName
		}
		if w.Pid != 0 {
			watcher["pid"] = w.Pid
		}
		if w.Port != 0 {
			watcher["port"] = w.Port
			watcher["protocol"] = w.Protocol
		}
		if w.CpuThreshold != 0 {
			watcher["cpu_threshold"] = w.CpuThreshold
		}
		if w.MemoryThreshold != 0 {
			watcher["memory_threshold"] = w.MemoryThreshold
		}
		if w.DiskThreshold != 0 {
			watcher["disk_threshold"] = w.DiskThreshold
		}

		watchers = append(watchers, watcher)
	}

	c.JSON(http.StatusOK, gin.H{
		"agent":    agentName,
		"watchers": watchers,
		"count":    len(watchers),
	})
}

// GetByAgent gets a specific watcher from an agent
func (h *WatcherHandler) GetByAgent(c *gin.Context) {
	agentName := c.Param("agent")
	watcherID := c.Param("id")

	if agentName == "" || watcherID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent and id parameters are required"})
		return
	}

	// Get agent info
	ctx := context.Background()
	agent, err := h.agentDB.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get agent client
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to connect to agent: %v", err)})
		return
	}

	// Get watcher
	resp, err := client.GetWatcher(context.Background(), &pb.GetWatcherRequest{
		WatcherId: watcherID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get watcher: %v", err)})
		return
	}

	if !resp.Found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Watcher not found"})
		return
	}

	w := resp.Watcher
	watcher := map[string]interface{}{
		"id":         w.Id,
		"type":       w.Type,
		"conditions": w.Conditions,
		"interval":   w.Interval,
	}

	// Add type-specific fields
	if w.FilePath != "" {
		watcher["file_path"] = w.FilePath
		watcher["check_hash"] = w.CheckHash
		watcher["recursive"] = w.Recursive
	}
	if w.ProcessName != "" {
		watcher["process_name"] = w.ProcessName
	}
	if w.Pid != 0 {
		watcher["pid"] = w.Pid
	}
	if w.Port != 0 {
		watcher["port"] = w.Port
		watcher["protocol"] = w.Protocol
	}
	if w.CpuThreshold != 0 {
		watcher["cpu_threshold"] = w.CpuThreshold
	}
	if w.MemoryThreshold != 0 {
		watcher["memory_threshold"] = w.MemoryThreshold
	}
	if w.DiskThreshold != 0 {
		watcher["disk_threshold"] = w.DiskThreshold
	}

	c.JSON(http.StatusOK, watcher)
}

// CreateForAgent creates a new watcher on an agent
func (h *WatcherHandler) CreateForAgent(c *gin.Context) {
	agentName := c.Param("agent")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent parameter is required"})
		return
	}

	var req WatcherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get agent info
	ctx := context.Background()
	agent, err := h.agentDB.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get agent client
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to connect to agent: %v", err)})
		return
	}

	// Create watcher config
	config := &pb.WatcherConfig{
		Type:       req.Type,
		Conditions: req.Conditions,
		Interval:   req.Interval,

		FilePath:  req.FilePath,
		CheckHash: req.CheckHash,
		Recursive: req.Recursive,

		ProcessName: req.ProcessName,
		Pid:         req.PID,

		Port:     req.Port,
		Protocol: req.Protocol,

		CpuThreshold:    req.CPUThreshold,
		MemoryThreshold: req.MemoryThreshold,
		DiskThreshold:   req.DiskThreshold,
	}

	// Register watcher
	resp, err := client.RegisterWatcher(context.Background(), &pb.RegisterWatcherRequest{
		Config: config,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to register watcher: %v", err)})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"message":    resp.Message,
		"watcher_id": resp.WatcherId,
	})
}

// DeleteFromAgent deletes a watcher from an agent
func (h *WatcherHandler) DeleteFromAgent(c *gin.Context) {
	agentName := c.Param("agent")
	watcherID := c.Param("id")

	if agentName == "" || watcherID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent and id parameters are required"})
		return
	}

	// Get agent info
	ctx := context.Background()
	agent, err := h.agentDB.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get agent client
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to connect to agent: %v", err)})
		return
	}

	// Remove watcher
	resp, err := client.RemoveWatcher(context.Background(), &pb.RemoveWatcherRequest{
		WatcherId: watcherID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove watcher: %v", err)})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": resp.Message,
	})
}

// ListAllWatchers lists watchers across all agents
func (h *WatcherHandler) ListAllWatchers(c *gin.Context) {
	ctx := context.Background()
	agents, err := h.agentDB.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allWatchers := make(map[string]interface{})
	totalCount := 0

	for _, agent := range agents {
		// Get agent client
		client, err := h.agentClient.GetClient(ctx, agent.Address)
		if err != nil {
			continue // Skip failed connections
		}

		// List watchers
		resp, err := client.ListWatchers(context.Background(), &pb.ListWatchersRequest{})
		if err != nil {
			continue // Skip failed requests
		}

		// Convert to response format
		watchers := make([]map[string]interface{}, 0, len(resp.Watchers))
		for _, w := range resp.Watchers {
			watcher := map[string]interface{}{
				"id":         w.Id,
				"type":       w.Type,
				"conditions": w.Conditions,
				"interval":   w.Interval,
			}

			// Add type-specific fields
			if w.FilePath != "" {
				watcher["file_path"] = w.FilePath
				watcher["check_hash"] = w.CheckHash
			}
			if w.ProcessName != "" {
				watcher["process_name"] = w.ProcessName
			}
			if w.Port != 0 {
				watcher["port"] = w.Port
			}
			if w.CpuThreshold != 0 {
				watcher["cpu_threshold"] = w.CpuThreshold
			}
			if w.MemoryThreshold != 0 {
				watcher["memory_threshold"] = w.MemoryThreshold
			}

			watchers = append(watchers, watcher)
		}

		allWatchers[agent.Name] = watchers
		totalCount += len(watchers)
	}

	c.JSON(http.StatusOK, gin.H{
		"watchers":    allWatchers,
		"total_count": totalCount,
		"agents":      len(allWatchers),
	})
}

// GetStatistics returns watcher statistics across all agents
func (h *WatcherHandler) GetStatistics(c *gin.Context) {
	ctx := context.Background()
	agents, err := h.agentDB.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stats := map[string]interface{}{
		"total_watchers": 0,
		"by_type": map[string]int{
			"file":    0,
			"cpu":     0,
			"memory":  0,
			"process": 0,
			"port":    0,
			"service": 0,
		},
		"by_agent": make(map[string]int),
	}

	for _, agent := range agents {
		// Get agent client
		client, err := h.agentClient.GetClient(ctx, agent.Address)
		if err != nil {
			continue
		}

		resp, err := client.ListWatchers(context.Background(), &pb.ListWatchersRequest{})
		if err != nil {
			continue
		}

		count := len(resp.Watchers)
		stats["total_watchers"] = stats["total_watchers"].(int) + count
		stats["by_agent"].(map[string]int)[agent.Name] = count

		// Count by type
		byType := stats["by_type"].(map[string]int)
		for _, w := range resp.Watchers {
			byType[w.Type]++
		}
	}

	c.JSON(http.StatusOK, stats)
}
