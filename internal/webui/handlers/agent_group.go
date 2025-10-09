package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentGroupHandler handles agent group operations
type AgentGroupHandler struct {
	db *AgentDBWrapper
}

// NewAgentGroupHandler creates a new agent group handler
func NewAgentGroupHandler(db *AgentDBWrapper) *AgentGroupHandler {
	return &AgentGroupHandler{db: db}
}

// AgentGroup represents a group of agents
type AgentGroup struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	AgentNames  []string          `json:"agent_names"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   int64             `json:"created_at"`
	AgentCount  int               `json:"agent_count"`
}

// List returns all agent groups
func (h *AgentGroupHandler) List(c *gin.Context) {
	groups, err := h.db.ListAgentGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list agent groups: " + err.Error()})
		return
	}

	if groups == nil {
		groups = []*AgentGroup{}
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// Get returns a specific agent group
func (h *AgentGroupHandler) Get(c *gin.Context) {
	groupID := c.Param("name")

	group, err := h.db.GetAgentGroup(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent group not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// Create creates a new agent group
func (h *AgentGroupHandler) Create(c *gin.Context) {
	var req struct {
		GroupName   string            `json:"group_name"`
		Description string            `json:"description"`
		AgentNames  []string          `json:"agent_names"`
		Tags        map[string]string `json:"tags"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.GroupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group name is required"})
		return
	}

	// Create the group
	group := &AgentGroup{
		ID:          req.GroupName, // Using name as ID for simplicity
		Name:        req.GroupName,
		Description: req.Description,
		AgentNames:  req.AgentNames,
		Tags:        req.Tags,
	}

	if err := h.db.CreateAgentGroup(c.Request.Context(), group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Group created successfully",
		"group_id": req.GroupName,
	})
}

// Delete deletes an agent group
func (h *AgentGroupHandler) Delete(c *gin.Context) {
	groupID := c.Param("name")

	if err := h.db.DeleteAgentGroup(c.Request.Context(), groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group deleted successfully",
	})
}

// AddAgents adds agents to a group
func (h *AgentGroupHandler) AddAgents(c *gin.Context) {
	groupID := c.Param("name")

	var req struct {
		AgentNames []string `json:"agent_names"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.AgentNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No agent names provided"})
		return
	}

	if err := h.db.AddAgentsToGroup(c.Request.Context(), groupID, req.AgentNames); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add agents: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Agents added to group",
	})
}

// RemoveAgents removes agents from a group
func (h *AgentGroupHandler) RemoveAgents(c *gin.Context) {
	groupID := c.Param("name")

	var req struct {
		AgentNames []string `json:"agent_names"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.AgentNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No agent names provided"})
		return
	}

	if err := h.db.RemoveAgentsFromGroup(c.Request.Context(), groupID, req.AgentNames); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove agents: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Agents removed from group",
	})
}

// GetAggregatedMetrics returns aggregated metrics for a group
func (h *AgentGroupHandler) GetAggregatedMetrics(c *gin.Context) {
	groupID := c.Param("name")

	metrics, err := h.db.GetGroupAggregatedMetrics(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get metrics: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
