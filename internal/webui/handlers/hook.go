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
