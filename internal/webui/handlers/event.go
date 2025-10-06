package handlers

import (
	"net/http"
	"strconv"

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
