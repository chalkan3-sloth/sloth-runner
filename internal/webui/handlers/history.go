package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/chalkan3-sloth/sloth-runner/internal/execution"
	"github.com/gin-gonic/gin"
)

// ListExecutionsHandler handles GET /api/v1/executions
func ListExecutionsHandler(c *gin.Context) {
	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open history database"})
		return
	}
	defer db.Close()

	// Parse query parameters
	filters := make(map[string]interface{})

	if workflow := c.Query("workflow"); workflow != "" {
		filters["workflow"] = workflow
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if agent := c.Query("agent"); agent != "" {
		filters["agent"] = agent
	}
	if group := c.Query("group"); group != "" {
		filters["group"] = group
	}
	if since := c.Query("since"); since != "" {
		sinceTime, err := parseDurationParam(since)
		if err == nil {
			filters["since"] = sinceTime
		}
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	executions, err := db.ListExecutions(filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list executions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"executions": executions})
}

// GetExecutionHandler handles GET /api/v1/executions/:id
func GetExecutionHandler(c *gin.Context) {
	id := c.Param("id")

	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open history database"})
		return
	}
	defer db.Close()

	exec, err := db.GetExecution(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Execution not found"})
		return
	}

	tasks, err := db.GetTaskExecutions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task executions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"execution": exec,
		"tasks":     tasks,
	})
}

// GetExecutionStatsHandler handles GET /api/v1/executions/stats
func GetExecutionStatsHandler(c *gin.Context) {
	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open history database"})
		return
	}
	defer db.Close()

	filters := make(map[string]interface{})

	if workflow := c.Query("workflow"); workflow != "" {
		filters["workflow"] = workflow
	}
	if since := c.Query("since"); since != "" {
		sinceTime, err := parseDurationParam(since)
		if err == nil {
			filters["since"] = sinceTime
		}
	}

	stats, err := db.GetStatistics(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteOldExecutionsHandler handles DELETE /api/v1/executions/cleanup
func DeleteOldExecutionsHandler(c *gin.Context) {
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	db, err := execution.NewHistoryDB(config.GetHistoryDBPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open history database"})
		return
	}
	defer db.Close()

	deleted, err := db.DeleteOldExecutions(req.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete executions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"deleted": deleted,
		"message": "Old executions deleted successfully",
	})
}

func parseDurationParam(s string) (int64, error) {
	var value int
	var unit string
	_, err := sscanf(s, "%d%s", &value, &unit)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	var since time.Time

	switch unit {
	case "h":
		since = now.Add(-time.Duration(value) * time.Hour)
	case "d":
		since = now.AddDate(0, 0, -value)
	case "w":
		since = now.AddDate(0, 0, -value*7)
	case "m":
		since = now.AddDate(0, -value, 0)
	default:
		// Default to days if no unit specified
		since = now.AddDate(0, 0, -value)
	}

	return since.Unix(), nil
}

// Simple sscanf implementation for parsing duration strings
func sscanf(s, format string, a ...interface{}) (int, error) {
	var value int
	var unit string

	// Simple parser for "123h", "7d", etc.
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		value = value*10 + int(s[i]-'0')
		i++
	}

	if i < len(s) {
		unit = s[i:]
	}

	if len(a) >= 2 {
		if ptr, ok := a[0].(*int); ok {
			*ptr = value
		}
		if ptr, ok := a[1].(*string); ok {
			*ptr = unit
		}
	}

	return 2, nil
}
