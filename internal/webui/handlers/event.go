package handlers

import (
	"net/http"
	"strconv"

	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/gin-gonic/gin"
)

// EventHandler handles event operations
type EventHandler struct {
	hookRepo *HookRepoWrapper
}

// NewEventHandler creates a new event handler
func NewEventHandler(hookRepo *HookRepoWrapper) *EventHandler {
	return &EventHandler{hookRepo: hookRepo}
}

// List returns all events
func (h *EventHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	eventQueue := h.hookRepo.GetEventQueue()
	// Use ListEvents with empty filters to get all events
	events, err := eventQueue.ListEvents("", "", limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// ListPending returns pending events
func (h *EventHandler) ListPending(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	eventQueue := h.hookRepo.GetEventQueue()
	events, err := eventQueue.GetPendingEvents(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// Get returns an event by ID
func (h *EventHandler) Get(c *gin.Context) {
	id := c.Param("id")

	eventQueue := h.hookRepo.GetEventQueue()
	event, err := eventQueue.GetEvent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// ListByAgent returns events for a specific agent
func (h *EventHandler) ListByAgent(c *gin.Context) {
	agent := c.Query("agent")
	if agent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent parameter is required"})
		return
	}

	eventType := c.DefaultQuery("type", "")
	status := c.DefaultQuery("status", "")
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	eventQueue := h.hookRepo.GetEventQueue()
	events, err := eventQueue.ListEventsByAgent(agent, hooks.EventType(eventType), hooks.EventStatus(status), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate statistics
	stats := map[string]int{
		"total":      len(events),
		"pending":    0,
		"processing": 0,
		"completed":  0,
		"failed":     0,
	}

	for _, event := range events {
		switch event.Status {
		case "pending":
			stats["pending"]++
		case "processing":
			stats["processing"]++
		case "completed":
			stats["completed"]++
		case "failed":
			stats["failed"]++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"stats":  stats,
	})
}

// Retry retries a failed event
func (h *EventHandler) Retry(c *gin.Context) {
	id := c.Param("id")

	eventQueue := h.hookRepo.GetEventQueue()
	// Update event status to pending to retry
	if err := eventQueue.UpdateEventStatus(id, "pending", ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event queued for retry"})
}

// ListHookExecutionsByAgent returns hook executions for a specific agent
func (h *EventHandler) ListHookExecutionsByAgent(c *gin.Context) {
	agent := c.Query("agent")
	if agent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent parameter is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	eventQueue := h.hookRepo.GetEventQueue()
	executions, err := eventQueue.GetHookExecutionsByAgent(agent, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate statistics
	stats := map[string]int{
		"total":     len(executions),
		"success":   0,
		"failed":    0,
	}

	for _, exec := range executions {
		if exec.Success {
			stats["success"]++
		} else {
			stats["failed"]++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"executions": executions,
		"stats":      stats,
	})
}

// GetStatistics returns event statistics
func (h *EventHandler) GetStatistics(c *gin.Context) {
	eventQueue := h.hookRepo.GetEventQueue()
	events, err := eventQueue.ListEvents("", "", 1000) // Get more for better stats
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stats := map[string]interface{}{
		"total":      len(events),
		"pending":    0,
		"processing": 0,
		"completed":  0,
		"failed":     0,
		"by_type":    make(map[string]int),
	}

	byType := stats["by_type"].(map[string]int)

	for _, event := range events {
		// Count by status
		switch event.Status {
		case "pending":
			stats["pending"] = stats["pending"].(int) + 1
		case "processing":
			stats["processing"] = stats["processing"].(int) + 1
		case "completed":
			stats["completed"] = stats["completed"].(int) + 1
		case "failed":
			stats["failed"] = stats["failed"].(int) + 1
		}

		// Count by type
		byType[string(event.Type)]++
	}

	c.JSON(http.StatusOK, stats)
}

// ListByType lists events by type
func (h *EventHandler) ListByType(c *gin.Context) {
	eventType := c.Query("type")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type parameter is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	eventQueue := h.hookRepo.GetEventQueue()
	events, err := eventQueue.ListEvents(hooks.EventType(eventType), "", limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
		"type":   eventType,
	})
}

// ListByStatus lists events by status
func (h *EventHandler) ListByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	eventQueue := h.hookRepo.GetEventQueue()
	events, err := eventQueue.ListEvents("", hooks.EventStatus(status), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
		"status": status,
	})
}

// GetRecentActivity returns recent event activity
func (h *EventHandler) GetRecentActivity(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	eventQueue := h.hookRepo.GetEventQueue()
	events, err := eventQueue.ListEvents("", "", limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Group events by minute for activity chart
	activity := make(map[string]int)
	for _, event := range events {
		minute := event.Timestamp.Format("2006-01-02 15:04")
		activity[minute]++
	}

	c.JSON(http.StatusOK, gin.H{
		"events":   events,
		"activity": activity,
	})
}
