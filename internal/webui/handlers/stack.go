package handlers

import (
	"net/http"

	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/gin-gonic/gin"
)

// StackHandler handles stack-related HTTP requests
type StackHandler struct {
	manager *stack.StackManager
}

// NewStackHandler creates a new stack handler
func NewStackHandler(stackDBPath string) (*StackHandler, error) {
	manager, err := stack.NewStackManager(stackDBPath)
	if err != nil {
		return nil, err
	}
	return &StackHandler{manager: manager}, nil
}

// List returns all stacks
func (h *StackHandler) List(c *gin.Context) {
	stacks, err := h.manager.ListStacks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stacks": stacks,
		"total":  len(stacks),
	})
}

// Get returns a specific stack
func (h *StackHandler) Get(c *gin.Context) {
	name := c.Param("name")

	stackState, err := h.manager.GetStackByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
		return
	}

	// Mock data for variables and secrets (to be implemented)
	variables := []map[string]interface{}{}
	secrets := []map[string]interface{}{}

	c.JSON(http.StatusOK, gin.H{
		"id":              stackState.ID,
		"name":            stackState.Name,
		"description":     stackState.Description,
		"type":            "custom",
		"active":          stackState.Status != "failed",
		"version":         stackState.Version,
		"status":          stackState.Status,
		"workflow_file":   stackState.WorkflowFile,
		"created_at":      stackState.CreatedAt,
		"updated_at":      stackState.UpdatedAt,
		"execution_count": stackState.ExecutionCount,
		"variables":       variables,
		"secrets":         secrets,
		"variable_count":  len(variables),
		"secret_count":    len(secrets),
	})
}

// Create creates a new stack
func (h *StackHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Active      bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stackState := &stack.StackState{
		Name:          req.Name,
		Description:   req.Description,
		Version:       "1.0.0",
		Status:        "created",
		TaskResults:   make(map[string]interface{}),
		Outputs:       make(map[string]interface{}),
		Configuration: make(map[string]interface{}),
		Metadata:      make(map[string]interface{}),
	}

	if err := h.manager.CreateStack(stackState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Stack created successfully",
		"stack": gin.H{
			"id":   stackState.ID,
			"name": stackState.Name,
		},
	})
}

// Update updates a stack
func (h *StackHandler) Update(c *gin.Context) {
	name := c.Param("name")

	var req struct {
		Description string `json:"description"`
		Type        string `json:"type"`
		Active      bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stackState, err := h.manager.GetStackByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
		return
	}

	stackState.Description = req.Description
	// Note: Type and Active fields don't exist in StackState

	if err := h.manager.UpdateStack(stackState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stack updated successfully"})
}

// Delete deletes a stack
func (h *StackHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	stackState, err := h.manager.GetStackByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
		return
	}

	if err := h.manager.DeleteStack(stackState.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stack deleted successfully"})
}

// AddVariable adds a variable to a stack
func (h *StackHandler) AddVariable(c *gin.Context) {
	stackName := c.Param("name")

	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stackState, err := h.manager.GetStackByName(stackName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
		return
	}

	// Store variable in stack configuration
	stackState.Configuration[req.Key] = req.Value
	if err := h.manager.UpdateStack(stackState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Variable added successfully"})
}

// DeleteVariable deletes a variable from a stack
func (h *StackHandler) DeleteVariable(c *gin.Context) {
	stackName := c.Param("name")
	key := c.Param("key")

	stackState, err := h.manager.GetStackByName(stackName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
		return
	}

	delete(stackState.Configuration, key)
	if err := h.manager.UpdateStack(stackState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variable deleted successfully"})
}
