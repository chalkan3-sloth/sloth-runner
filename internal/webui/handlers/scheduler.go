package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SchedulerHandler handles scheduled workflow executions
type SchedulerHandler struct {
	schedules map[string]*Schedule
	wsHub     *WebSocketHub
}

// Schedule represents a scheduled workflow execution
type Schedule struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	WorkflowName string            `json:"workflow_name"`
	CronExpr     string            `json:"cron_expr"`
	Enabled      bool              `json:"enabled"`
	DelegateTo   string            `json:"delegate_to,omitempty"`
	Variables    map[string]string `json:"variables,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	LastRun      *time.Time        `json:"last_run,omitempty"`
	NextRun      *time.Time        `json:"next_run,omitempty"`
	RunCount     int               `json:"run_count"`
}

// NewSchedulerHandler creates a new scheduler handler
func NewSchedulerHandler(wsHub *WebSocketHub) *SchedulerHandler {
	return &SchedulerHandler{
		schedules: make(map[string]*Schedule),
		wsHub:     wsHub,
	}
}

// ListSchedules returns all schedules
func (h *SchedulerHandler) ListSchedules(c *gin.Context) {
	schedules := make([]*Schedule, 0, len(h.schedules))
	for _, s := range h.schedules {
		schedules = append(schedules, s)
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// GetSchedule returns a specific schedule
func (h *SchedulerHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, exists := h.schedules[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// CreateSchedule creates a new schedule
func (h *SchedulerHandler) CreateSchedule(c *gin.Context) {
	var req struct {
		Name         string            `json:"name" binding:"required"`
		WorkflowName string            `json:"workflow_name" binding:"required"`
		CronExpr     string            `json:"cron_expr" binding:"required"`
		Enabled      bool              `json:"enabled"`
		DelegateTo   string            `json:"delegate_to"`
		Variables    map[string]string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule := &Schedule{
		ID:           uuid.New().String(),
		Name:         req.Name,
		WorkflowName: req.WorkflowName,
		CronExpr:     req.CronExpr,
		Enabled:      req.Enabled,
		DelegateTo:   req.DelegateTo,
		Variables:    req.Variables,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		RunCount:     0,
	}

	h.schedules[schedule.ID] = schedule

	c.JSON(http.StatusCreated, schedule)
}

// UpdateSchedule updates a schedule
func (h *SchedulerHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, exists := h.schedules[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	var req struct {
		Name         string            `json:"name"`
		WorkflowName string            `json:"workflow_name"`
		CronExpr     string            `json:"cron_expr"`
		Enabled      *bool             `json:"enabled"`
		DelegateTo   string            `json:"delegate_to"`
		Variables    map[string]string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		schedule.Name = req.Name
	}
	if req.WorkflowName != "" {
		schedule.WorkflowName = req.WorkflowName
	}
	if req.CronExpr != "" {
		schedule.CronExpr = req.CronExpr
	}
	if req.Enabled != nil {
		schedule.Enabled = *req.Enabled
	}
	schedule.DelegateTo = req.DelegateTo
	schedule.Variables = req.Variables
	schedule.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, schedule)
}

// DeleteSchedule deletes a schedule
func (h *SchedulerHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")

	if _, exists := h.schedules[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	delete(h.schedules, id)

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}

// EnableSchedule enables a schedule
func (h *SchedulerHandler) EnableSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, exists := h.schedules[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	schedule.Enabled = true
	schedule.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, gin.H{"message": "Schedule enabled"})
}

// DisableSchedule disables a schedule
func (h *SchedulerHandler) DisableSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, exists := h.schedules[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	schedule.Enabled = false
	schedule.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, gin.H{"message": "Schedule disabled"})
}

// TriggerSchedule manually triggers a schedule execution
func (h *SchedulerHandler) TriggerSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, exists := h.schedules[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	// TODO: Trigger actual workflow execution
	now := time.Now()
	schedule.LastRun = &now
	schedule.RunCount++

	c.JSON(http.StatusOK, gin.H{
		"message": "Schedule triggered",
		"schedule": schedule,
	})
}
