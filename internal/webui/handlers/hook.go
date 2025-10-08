package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
)

// HookHandler handles hook operations
type HookHandler struct {
	repo *HookRepoWrapper
}

// NewHookHandler creates a new hook handler
func NewHookHandler(repo *HookRepoWrapper) *HookHandler {
	return &HookHandler{repo: repo}
}

// List returns all hooks
func (h *HookHandler) List(c *gin.Context) {
	hookList, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hooks": hookList})
}

// Get returns a hook by ID
func (h *HookHandler) Get(c *gin.Context) {
	id := c.Param("id")

	hook, err := h.repo.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hook not found"})
		return
	}

	c.JSON(http.StatusOK, hook)
}

// CreateHookRequest represents a create hook request
type CreateHookRequest struct {
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description"`
	EventType   hooks.EventType  `json:"event_type" binding:"required"`
	FilePath    string           `json:"file_path" binding:"required"`
	Stack       string           `json:"stack"`
	Enabled     bool             `json:"enabled"`
}

// Create creates a new hook
func (h *HookHandler) Create(c *gin.Context) {
	var req CreateHookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hook := &hooks.Hook{
		Name:        req.Name,
		Description: req.Description,
		EventType:   req.EventType,
		FilePath:    req.FilePath,
		Stack:       req.Stack,
		Enabled:     req.Enabled,
		RunCount:    0,
	}

	if err := h.repo.Add(hook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hook)
}

// Update updates a hook
func (h *HookHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateHookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing hook
	hook, err := h.repo.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hook not found"})
		return
	}

	hook.Name = req.Name
	hook.Description = req.Description
	hook.EventType = req.EventType
	hook.FilePath = req.FilePath
	hook.Stack = req.Stack
	hook.Enabled = req.Enabled

	if err := h.repo.Update(hook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hook)
}

// Delete deletes a hook
func (h *HookHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hook deleted successfully"})
}

// Enable enables a hook
func (h *HookHandler) Enable(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Enable(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hook enabled successfully"})
}

// Disable disables a hook
func (h *HookHandler) Disable(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Disable(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hook disabled successfully"})
}

// GetHistory returns execution history for a hook
func (h *HookHandler) GetHistory(c *gin.Context) {
	id := c.Param("id")
	limit := 50 // Default limit

	history, err := h.repo.GetExecutionHistory(id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// GetStatistics returns statistics for all hooks
func (h *HookHandler) GetStatistics(c *gin.Context) {
	hookList, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total":       len(hookList),
		"enabled":     0,
		"disabled":    0,
		"total_runs":  int64(0),
		"by_event_type": make(map[string]int),
		"by_stack":    make(map[string]int),
	}

	byEventType := stats["by_event_type"].(map[string]int)
	byStack := stats["by_stack"].(map[string]int)

	for _, hook := range hookList {
		stats["total_runs"] = stats["total_runs"].(int64) + hook.RunCount

		if hook.Enabled {
			stats["enabled"] = stats["enabled"].(int) + 1
		} else {
			stats["disabled"] = stats["disabled"].(int) + 1
		}

		// Count by event type
		eventType := string(hook.EventType)
		byEventType[eventType]++

		// Count by stack
		stack := hook.Stack
		if stack == "" {
			stack = "default"
		}
		byStack[stack]++
	}

	c.JSON(http.StatusOK, stats)
}

// ListByEventType lists hooks by event type
func (h *HookHandler) ListByEventType(c *gin.Context) {
	eventType := c.Query("event_type")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_type parameter is required"})
		return
	}

	hookList, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by event type
	filtered := make([]*hooks.Hook, 0)
	for _, hook := range hookList {
		if string(hook.EventType) == eventType {
			filtered = append(filtered, hook)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"hooks":      filtered,
		"count":      len(filtered),
		"event_type": eventType,
	})
}

// ListByStack lists hooks by stack
func (h *HookHandler) ListByStack(c *gin.Context) {
	stack := c.Query("stack")
	if stack == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stack parameter is required"})
		return
	}

	hookList, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by stack
	filtered := make([]*hooks.Hook, 0)
	for _, hook := range hookList {
		if hook.Stack == stack {
			filtered = append(filtered, hook)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"hooks": filtered,
		"count": len(filtered),
		"stack": stack,
	})
}

// GetExecutionStats returns execution statistics for a hook
func (h *HookHandler) GetExecutionStats(c *gin.Context) {
	id := c.Param("id")

	history, err := h.repo.GetExecutionHistory(id, 1000) // Get more for better stats
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	successCount := 0
	failedCount := 0
	totalDuration := int64(0)

	for _, result := range history {
		if result.Success {
			successCount++
		} else {
			failedCount++
		}
		totalDuration += result.Duration.Milliseconds()
	}

	avgDuration := int64(0)
	if len(history) > 0 {
		avgDuration = totalDuration / int64(len(history))
	}

	stats := map[string]interface{}{
		"total_executions": len(history),
		"success_count":    successCount,
		"failed_count":     failedCount,
		"success_rate":     0.0,
		"avg_duration_ms":  avgDuration,
	}

	if len(history) > 0 {
		stats["success_rate"] = float64(successCount) / float64(len(history)) * 100
	}

	c.JSON(http.StatusOK, stats)
}
