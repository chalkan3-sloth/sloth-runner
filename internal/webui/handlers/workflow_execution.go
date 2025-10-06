package handlers

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	osexec "os/exec"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WorkflowExecutionHandler handles workflow execution operations
type WorkflowExecutionHandler struct {
	executions map[string]*WorkflowExecution
	mu         sync.RWMutex
	wsHub      *WebSocketHub
}

// WorkflowExecution represents a running workflow
type WorkflowExecution struct {
	ID         string    `json:"id"`
	WorkflowID string    `json:"workflow_id"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Logs       []string  `json:"logs"`
	ExitCode   *int      `json:"exit_code,omitempty"`
	Error      string    `json:"error,omitempty"`
	cmd        *osexec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkflowExecutionHandler creates a new workflow execution handler
func NewWorkflowExecutionHandler(wsHub *WebSocketHub) *WorkflowExecutionHandler {
	return &WorkflowExecutionHandler{
		executions: make(map[string]*WorkflowExecution),
		wsHub:      wsHub,
	}
}

// ExecuteWorkflow starts a workflow execution
func (h *WorkflowExecutionHandler) ExecuteWorkflow(c *gin.Context) {
	var req struct {
		WorkflowName string            `json:"workflow_name" binding:"required"`
		DelegateTo   string            `json:"delegate_to"`
		Variables    map[string]string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	executionID := uuid.New().String()
	ctx, cancel := context.WithCancel(context.Background())

	execution := &WorkflowExecution{
		ID:         executionID,
		WorkflowID: req.WorkflowName,
		Status:     "running",
		StartTime:  time.Now(),
		Logs:       make([]string, 0),
		ctx:        ctx,
		cancel:     cancel,
	}

	h.mu.Lock()
	h.executions[executionID] = execution
	h.mu.Unlock()

	// Start execution in goroutine
	go h.runWorkflow(execution, req.WorkflowName, req.DelegateTo, req.Variables)

	c.JSON(http.StatusOK, gin.H{
		"execution_id": executionID,
		"status":       "started",
	})
}

// runWorkflow executes the workflow
func (h *WorkflowExecutionHandler) runWorkflow(exec *WorkflowExecution, workflowName, delegateTo string, variables map[string]string) {
	// Build command
	args := []string{"run", workflowName}

	if delegateTo != "" {
		args = append(args, "--delegate-to", delegateTo)
	}

	for key, value := range variables {
		args = append(args, "--var", fmt.Sprintf("%s=%s", key, value))
	}

	cmd := osexec.CommandContext(exec.ctx, "sloth-runner", args...)
	exec.cmd = cmd

	// Capture stdout and stderr
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		exec.Status = "failed"
		exec.Error = err.Error()
		h.notifyUpdate(exec)
		return
	}

	// Read output
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			exec.Logs = append(exec.Logs, line)
			h.notifyLog(exec.ID, line)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			exec.Logs = append(exec.Logs, "[ERROR] " + line)
			h.notifyLog(exec.ID, "[ERROR] "+line)
		}
	}()

	// Wait for completion
	err := cmd.Wait()
	endTime := time.Now()
	exec.EndTime = &endTime

	if err != nil {
		exec.Status = "failed"
		exec.Error = err.Error()
		if exitErr, ok := err.(*osexec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			exec.ExitCode = &exitCode
		}
	} else {
		exec.Status = "completed"
		exitCode := 0
		exec.ExitCode = &exitCode
	}

	h.notifyUpdate(exec)
}

// GetExecution returns execution details
func (h *WorkflowExecutionHandler) GetExecution(c *gin.Context) {
	id := c.Param("id")

	h.mu.RLock()
	exec, exists := h.executions[id]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Execution not found"})
		return
	}

	c.JSON(http.StatusOK, exec)
}

// ListExecutions returns all executions
func (h *WorkflowExecutionHandler) ListExecutions(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	executions := make([]*WorkflowExecution, 0, len(h.executions))
	for _, exec := range h.executions {
		executions = append(executions, exec)
	}

	c.JSON(http.StatusOK, gin.H{"executions": executions})
}

// CancelExecution cancels a running execution
func (h *WorkflowExecutionHandler) CancelExecution(c *gin.Context) {
	id := c.Param("id")

	h.mu.RLock()
	exec, exists := h.executions[id]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Execution not found"})
		return
	}

	if exec.Status != "running" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Execution is not running"})
		return
	}

	exec.cancel()
	exec.Status = "cancelled"
	endTime := time.Now()
	exec.EndTime = &endTime

	h.notifyUpdate(exec)

	c.JSON(http.StatusOK, gin.H{"message": "Execution cancelled"})
}

// GetExecutionLogs streams logs for an execution
func (h *WorkflowExecutionHandler) GetExecutionLogs(c *gin.Context) {
	id := c.Param("id")

	h.mu.RLock()
	exec, exists := h.executions[id]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Execution not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": exec.Logs})
}

// notifyUpdate sends execution update via WebSocket
func (h *WorkflowExecutionHandler) notifyUpdate(exec *WorkflowExecution) {
	if h.wsHub != nil {
		h.wsHub.Broadcast("workflow_execution_update", exec)
	}
}

// notifyLog sends log line via WebSocket
func (h *WorkflowExecutionHandler) notifyLog(executionID, logLine string) {
	if h.wsHub != nil {
		h.wsHub.Broadcast("workflow_execution_log", gin.H{
			"execution_id": executionID,
			"log":          logLine,
			"timestamp":    time.Now().Unix(),
		})
	}
}
