package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsHistoryHandler handles persistent metrics history operations
type MetricsHistoryHandler struct {
	metricsDB *metrics.MetricsDB
}

// NewMetricsHistoryHandler creates a new metrics history handler
func NewMetricsHistoryHandler(db *metrics.MetricsDB) *MetricsHistoryHandler {
	return &MetricsHistoryHandler{
		metricsDB: db,
	}
}

// GetAgentMetricsHistory returns historical metrics for a specific agent
// GET /api/v1/agents/:name/metrics/persistent
// Query params:
//   - duration: Time range (1h, 6h, 24h, 7d) - default: 1h
//   - maxPoints: Maximum number of data points to return - default: 60
func (h *MetricsHistoryHandler) GetAgentMetricsHistory(c *gin.Context) {
	ctx := c.Request.Context()
	agentName := c.Param("name")

	// Parse duration parameter
	durationStr := c.DefaultQuery("duration", "1h")
	duration, err := parseDuration(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid duration format. Use: 1h, 6h, 24h, 7d",
		})
		return
	}

	// Parse maxPoints parameter
	maxPoints := 60
	if maxPointsStr := c.Query("maxPoints"); maxPointsStr != "" {
		if mp, err := strconv.Atoi(maxPointsStr); err == nil && mp > 0 {
			maxPoints = mp
		}
	}

	// Calculate time range
	endTime := time.Now().Unix()
	startTime := endTime - int64(duration.Seconds())

	// Fetch metrics from database
	datapoints, err := h.metricsDB.GetMetricsHistory(ctx, agentName, startTime, endTime, maxPoints)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch metrics history",
		})
		return
	}

	// Return empty array if no data
	if datapoints == nil {
		datapoints = []metrics.MetricPoint{}
	}

	c.JSON(http.StatusOK, gin.H{
		"agent_name": agentName,
		"start_time": startTime,
		"end_time":   endTime,
		"duration":   durationStr,
		"datapoints": datapoints,
	})
}

// GetAllAgentsMetrics returns current (latest) metrics for all agents
// GET /api/v1/metrics/all
func (h *MetricsHistoryHandler) GetAllAgentsMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all agent names from metrics DB
	agentNames, err := h.metricsDB.GetAgentNames(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch agent names",
		})
		return
	}

	// Fetch latest metrics for each agent
	allMetrics := make(map[string]interface{})
	for _, agentName := range agentNames {
		latestMetric, err := h.metricsDB.GetLatestMetric(ctx, agentName)
		if err != nil {
			// Skip this agent if there's an error
			continue
		}
		if latestMetric != nil {
			allMetrics[agentName] = latestMetric
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"agents":    allMetrics,
		"timestamp": time.Now().Unix(),
	})
}

// GetMetricsStats returns statistics about metrics collection
// GET /api/v1/metrics/stats
func (h *MetricsHistoryHandler) GetMetricsStats(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all agent names
	agentNames, err := h.metricsDB.GetAgentNames(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch metrics statistics",
		})
		return
	}

	// Calculate total data points across all agents
	totalPoints := 0
	oldestTimestamp := int64(0)
	newestTimestamp := int64(0)

	for _, agentName := range agentNames {
		// Get metrics for the last 30 days to calculate stats
		endTime := time.Now().Unix()
		startTime := endTime - (30 * 24 * 60 * 60) // 30 days ago

		datapoints, err := h.metricsDB.GetMetricsHistory(ctx, agentName, startTime, endTime, 0)
		if err != nil {
			continue
		}

		totalPoints += len(datapoints)

		// Track oldest and newest timestamps
		if len(datapoints) > 0 {
			if oldestTimestamp == 0 || datapoints[0].Timestamp < oldestTimestamp {
				oldestTimestamp = datapoints[0].Timestamp
			}
			if datapoints[len(datapoints)-1].Timestamp > newestTimestamp {
				newestTimestamp = datapoints[len(datapoints)-1].Timestamp
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_agents":      len(agentNames),
		"total_data_points": totalPoints,
		"oldest_metric":     oldestTimestamp,
		"newest_metric":     newestTimestamp,
		"collection_active": true,
	})
}

// parseDuration parses duration strings like "1h", "6h", "24h", "7d"
func parseDuration(durationStr string) (time.Duration, error) {
	switch durationStr {
	case "1h":
		return time.Hour, nil
	case "6h":
		return 6 * time.Hour, nil
	case "12h":
		return 12 * time.Hour, nil
	case "24h":
		return 24 * time.Hour, nil
	case "7d":
		return 7 * 24 * time.Hour, nil
	case "30d":
		return 30 * 24 * time.Hour, nil
	default:
		// Try to parse as Go duration
		return time.ParseDuration(durationStr)
	}
}
