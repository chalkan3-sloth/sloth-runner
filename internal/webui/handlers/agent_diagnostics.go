package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AgentDiagnosticsHandler handles advanced agent diagnostics and troubleshooting
type AgentDiagnosticsHandler struct {
	db          *AgentDBWrapper
	agentClient *services.AgentClient
}

// NewAgentDiagnosticsHandler creates a new diagnostics handler
func NewAgentDiagnosticsHandler(db *AgentDBWrapper, agentClient *services.AgentClient) *AgentDiagnosticsHandler {
	return &AgentDiagnosticsHandler{
		db:          db,
		agentClient: agentClient,
	}
}

// GetDetailedMetrics returns comprehensive system metrics for troubleshooting
// GET /api/v1/agents/:name/metrics/detailed
func (h *AgentDiagnosticsHandler) GetDetailedMetrics(c *gin.Context) {
	agentName := c.Param("name")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent name is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Get agent from database
	agent, err := h.db.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Agent not found: %v", err)})
		return
	}

	// Connect to agent using address from database
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Failed to connect to agent: %v", err),
		})
		return
	}
	resp, err := client.GetDetailedMetrics(ctx, &pb.DetailedMetricsRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unimplemented {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Agent does not support detailed metrics (update required)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get detailed metrics: %v", err),
		})
		return
	}

	// Convert protobuf response to JSON-friendly format
	metrics := map[string]interface{}{
		"timestamp": resp.Timestamp,
		"cpu": map[string]interface{}{
			"core_count":     resp.Cpu.CoreCount,
			"per_core_usage": resp.Cpu.PerCoreUsage,
			"user_time":      resp.Cpu.UserTime,
			"system_time":    resp.Cpu.SystemTime,
			"idle_time":      resp.Cpu.IdleTime,
			"iowait_time":    resp.Cpu.IowaitTime,
			"model_name":     resp.Cpu.ModelName,
			"mhz":            resp.Cpu.Mhz,
		},
		"memory": map[string]interface{}{
			"total_bytes":      resp.Memory.TotalBytes,
			"available_bytes":  resp.Memory.AvailableBytes,
			"used_bytes":       resp.Memory.UsedBytes,
			"free_bytes":       resp.Memory.FreeBytes,
			"cached_bytes":     resp.Memory.CachedBytes,
			"buffers_bytes":    resp.Memory.BuffersBytes,
			"swap_total_bytes": resp.Memory.SwapTotalBytes,
			"swap_used_bytes":  resp.Memory.SwapUsedBytes,
			"swap_free_bytes":  resp.Memory.SwapFreeBytes,
			"percent":          resp.Memory.Percent,
			"swap_percent":     resp.Memory.SwapPercent,
		},
		"disk": map[string]interface{}{
			"partitions":        convertDiskPartitions(resp.Disk.Partitions),
			"read_bytes_total":  resp.Disk.ReadBytesTotal,
			"write_bytes_total": resp.Disk.WriteBytesTotal,
			"read_count":        resp.Disk.ReadCount,
			"write_count":       resp.Disk.WriteCount,
			"read_time_ms":      resp.Disk.ReadTimeMs,
			"write_time_ms":     resp.Disk.WriteTimeMs,
		},
		"network": map[string]interface{}{
			"interfaces":          convertNetworkInterfaces(resp.Network.Interfaces),
			"bytes_sent_total":    resp.Network.BytesSentTotal,
			"bytes_recv_total":    resp.Network.BytesRecvTotal,
			"packets_sent_total":  resp.Network.PacketsSentTotal,
			"packets_recv_total":  resp.Network.PacketsRecvTotal,
			"errors_in":           resp.Network.ErrorsIn,
			"errors_out":          resp.Network.ErrorsOut,
			"drops_in":            resp.Network.DropsIn,
			"drops_out":           resp.Network.DropsOut,
			"active_connections":  resp.Network.ActiveConnections,
			"listening_ports":     resp.Network.ListeningPorts,
		},
		"load_avg_1min":   resp.LoadAvg_1Min,
		"load_avg_5min":   resp.LoadAvg_5Min,
		"load_avg_15min":  resp.LoadAvg_15Min,
		"uptime_seconds":  resp.UptimeSeconds,
		"process_count":   resp.ProcessCount,
		"thread_count":    resp.ThreadCount,
		"kernel_version":  resp.KernelVersion,
		"os_version":      resp.OsVersion,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetRecentLogs returns recent log entries for troubleshooting
// GET /api/v1/agents/:name/logs?max_lines=100&level=error
func (h *AgentDiagnosticsHandler) GetRecentLogs(c *gin.Context) {
	agentName := c.Param("name")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent name is required"})
		return
	}

	maxLines, _ := strconv.ParseInt(c.DefaultQuery("max_lines", "100"), 10, 32)
	levelFilter := c.Query("level")
	sourceFilter := c.Query("source")
	sinceTimestamp, _ := strconv.ParseInt(c.DefaultQuery("since", "0"), 10, 64)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	agent, err := h.db.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Agent not found: %v", err)})
		return
	}

	// Connect to agent using address from database
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Failed to connect to agent: %v", err),
		})
		return
	}

	resp, err := client.GetRecentLogs(ctx, &pb.RecentLogsRequest{
		MaxLines:       int32(maxLines),
		LevelFilter:    levelFilter,
		SourceFilter:   sourceFilter,
		SinceTimestamp: sinceTimestamp,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unimplemented {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Agent does not support log retrieval (update required)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get logs: %v", err),
		})
		return
	}

	logs := make([]map[string]interface{}, len(resp.Logs))
	for i, log := range resp.Logs {
		logs[i] = map[string]interface{}{
			"timestamp": log.Timestamp,
			"level":     log.Level,
			"message":   log.Message,
			"source":    log.Source,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":        logs,
		"total_count": resp.TotalCount,
		"has_more":    resp.HasMore,
	})
}

// GetActiveConnections returns active network connections
// GET /api/v1/agents/:name/connections?state=ESTABLISHED
func (h *AgentDiagnosticsHandler) GetActiveConnections(c *gin.Context) {
	agentName := c.Param("name")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent name is required"})
		return
	}

	stateFilter := c.Query("state")
	includeLocal, _ := strconv.ParseBool(c.DefaultQuery("include_local", "false"))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	agent, err := h.db.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Agent not found: %v", err)})
		return
	}

	// Connect to agent using address from database
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Failed to connect to agent: %v", err),
		})
		return
	}

	resp, err := client.GetActiveConnections(ctx, &pb.ConnectionsRequest{
		StateFilter:  stateFilter,
		IncludeLocal: includeLocal,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unimplemented {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Agent does not support connection monitoring (update required)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get connections: %v", err),
		})
		return
	}

	connections := make([]map[string]interface{}, len(resp.Connections))
	for i, conn := range resp.Connections {
		connections[i] = map[string]interface{}{
			"local_addr":      conn.LocalAddr,
			"local_port":      conn.LocalPort,
			"remote_addr":     conn.RemoteAddr,
			"remote_port":     conn.RemotePort,
			"state":           conn.State,
			"pid":             conn.Pid,
			"process_name":    conn.ProcessName,
			"bytes_sent":      conn.BytesSent,
			"bytes_recv":      conn.BytesRecv,
			"established_at":  conn.EstablishedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"connections":       connections,
		"total_established": resp.TotalEstablished,
		"total_listening":   resp.TotalListening,
		"total_time_wait":   resp.TotalTimeWait,
		"total_all":         resp.TotalAll,
	})
}

// GetSystemErrors returns recent system errors
// GET /api/v1/agents/:name/errors?max_errors=50
func (h *AgentDiagnosticsHandler) GetSystemErrors(c *gin.Context) {
	agentName := c.Param("name")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent name is required"})
		return
	}

	maxErrors, _ := strconv.ParseInt(c.DefaultQuery("max_errors", "50"), 10, 32)
	sinceTimestamp, _ := strconv.ParseInt(c.DefaultQuery("since", "0"), 10, 64)
	includeWarnings, _ := strconv.ParseBool(c.DefaultQuery("include_warnings", "true"))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	agent, err := h.db.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Agent not found: %v", err)})
		return
	}

	// Connect to agent using address from database
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Failed to connect to agent: %v", err),
		})
		return
	}

	resp, err := client.GetSystemErrors(ctx, &pb.SystemErrorsRequest{
		MaxErrors:       int32(maxErrors),
		SinceTimestamp:  sinceTimestamp,
		IncludeWarnings: includeWarnings,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unimplemented {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Agent does not support error tracking (update required)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get errors: %v", err),
		})
		return
	}

	errors := make([]map[string]interface{}, len(resp.Errors))
	for i, e := range resp.Errors {
		errors[i] = map[string]interface{}{
			"timestamp":        e.Timestamp,
			"severity":         e.Severity,
			"source":           e.Source,
			"message":          e.Message,
			"stack_trace":      e.StackTrace,
			"context":          e.Context,
			"occurrence_count": e.OccurrenceCount,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errors":             errors,
		"total_errors":       resp.TotalErrors,
		"total_warnings":     resp.TotalWarnings,
		"most_common_error":  resp.MostCommonError,
	})
}

// GetPerformanceHistory returns historical performance data
// GET /api/v1/agents/:name/performance/history?duration=60&data_points=60
func (h *AgentDiagnosticsHandler) GetPerformanceHistory(c *gin.Context) {
	agentName := c.Param("name")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent name is required"})
		return
	}

	durationMinutes, _ := strconv.ParseInt(c.DefaultQuery("duration", "60"), 10, 32)
	dataPoints, _ := strconv.ParseInt(c.DefaultQuery("data_points", "60"), 10, 32)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	agent, err := h.db.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Agent not found: %v", err)})
		return
	}

	// Connect to agent using address from database
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Failed to connect to agent: %v", err),
		})
		return
	}

	resp, err := client.GetPerformanceHistory(ctx, &pb.PerformanceHistoryRequest{
		DurationMinutes: int32(durationMinutes),
		DataPoints:      int32(dataPoints),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unimplemented {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Agent does not support performance history (update required)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get performance history: %v", err),
		})
		return
	}

	snapshots := make([]map[string]interface{}, len(resp.Snapshots))
	for i, s := range resp.Snapshots {
		snapshots[i] = convertPerformanceSnapshot(s)
	}

	c.JSON(http.StatusOK, gin.H{
		"snapshots": snapshots,
		"avg":       convertPerformanceSnapshot(resp.Avg),
		"min":       convertPerformanceSnapshot(resp.Min),
		"max":       convertPerformanceSnapshot(resp.Max),
	})
}

// DiagnoseHealth performs health diagnostic check
// GET /api/v1/agents/:name/health/diagnose?deep_check=true
func (h *AgentDiagnosticsHandler) DiagnoseHealth(c *gin.Context) {
	agentName := c.Param("name")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent name is required"})
		return
	}

	includeSuggestions, _ := strconv.ParseBool(c.DefaultQuery("include_suggestions", "true"))
	deepCheck, _ := strconv.ParseBool(c.DefaultQuery("deep_check", "false"))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	agent, err := h.db.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Agent not found: %v", err)})
		return
	}

	// Connect to agent using address from database
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Failed to connect to agent: %v", err),
		})
		return
	}

	resp, err := client.DiagnoseHealth(ctx, &pb.HealthDiagnosticRequest{
		IncludeSuggestions: includeSuggestions,
		DeepCheck:          deepCheck,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unimplemented {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Agent does not support health diagnostics (update required)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to diagnose health: %v", err),
		})
		return
	}

	issues := make([]map[string]interface{}, len(resp.Issues))
	for i, issue := range resp.Issues {
		issues[i] = map[string]interface{}{
			"category":      issue.Category,
			"severity":      issue.Severity,
			"description":   issue.Description,
			"current_value": issue.CurrentValue,
			"threshold":     issue.Threshold,
			"suggestions":   issue.Suggestions,
			"auto_fixable":  issue.AutoFixable,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"overall_status":  resp.OverallStatus,
		"health_score":    resp.HealthScore,
		"issues":          issues,
		"summary":         resp.Summary,
		"check_timestamp": resp.CheckTimestamp,
		"total_warnings":  resp.TotalWarnings,
		"total_errors":    resp.TotalErrors,
	})
}

// Helper functions

func convertDiskPartitions(partitions []*pb.DiskPartition) []map[string]interface{} {
	result := make([]map[string]interface{}, len(partitions))
	for i, p := range partitions {
		result[i] = map[string]interface{}{
			"device":     p.Device,
			"mountpoint": p.Mountpoint,
			"fstype":     p.Fstype,
			"total_bytes": p.TotalBytes,
			"used_bytes":  p.UsedBytes,
			"free_bytes":  p.FreeBytes,
			"percent":     p.Percent,
		}
	}
	return result
}

func convertNetworkInterfaces(interfaces []*pb.NetworkInterface) []map[string]interface{} {
	result := make([]map[string]interface{}, len(interfaces))
	for i, iface := range interfaces {
		result[i] = map[string]interface{}{
			"name":         iface.Name,
			"ip_addresses": iface.IpAddresses,
			"mac_address":  iface.MacAddress,
			"bytes_sent":   iface.BytesSent,
			"bytes_recv":   iface.BytesRecv,
			"packets_sent": iface.PacketsSent,
			"packets_recv": iface.PacketsRecv,
			"is_up":        iface.IsUp,
		}
	}
	return result
}

func convertPerformanceSnapshot(s *pb.PerformanceSnapshot) map[string]interface{} {
	if s == nil {
		return nil
	}
	return map[string]interface{}{
		"timestamp":                s.Timestamp,
		"cpu_percent":              s.CpuPercent,
		"memory_percent":           s.MemoryPercent,
		"disk_percent":             s.DiskPercent,
		"network_throughput_mbps":  s.NetworkThroughputMbps,
		"load_avg":                 s.LoadAvg,
		"active_connections":       s.ActiveConnections,
		"process_count":            s.ProcessCount,
	}
}
