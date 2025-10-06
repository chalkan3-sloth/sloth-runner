package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentHandler handles agent operations
type AgentHandler struct {
	db *AgentDBWrapper
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(db *AgentDBWrapper) *AgentHandler {
	return &AgentHandler{db: db}
}

// List returns all agents
func (h *AgentHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	agents, err := h.db.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// Get returns an agent by name
func (h *AgentHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	agent, err := h.db.GetAgent(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// Delete removes an agent
func (h *AgentHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if err := h.db.DeleteAgent(ctx, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

// GetStats returns agent statistics
func (h *AgentHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.db.GetStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Update updates an agent (placeholder for future implementation)
func (h *AgentHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}
