package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentGroupAdvancedHandler handles advanced agent group operations
type AgentGroupAdvancedHandler struct {
	db *AgentDBWrapper
}

// NewAgentGroupAdvancedHandler creates a new advanced handler
func NewAgentGroupAdvancedHandler(db *AgentDBWrapper) *AgentGroupAdvancedHandler {
	return &AgentGroupAdvancedHandler{db: db}
}

// ExecuteBulkOperation executes a bulk operation on a group
func (h *AgentGroupAdvancedHandler) ExecuteBulkOperation(c *gin.Context) {
	var req BulkOperation

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.GroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group ID is required"})
		return
	}

	if req.Operation == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Operation is required"})
		return
	}

	result, err := h.db.ExecuteBulkOperation(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute bulk operation: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListTemplates returns all group templates
func (h *AgentGroupAdvancedHandler) ListTemplates(c *gin.Context) {
	templates, err := h.db.ListGroupTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list templates: " + err.Error()})
		return
	}

	if templates == nil {
		templates = []*GroupTemplate{}
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetTemplate returns a specific template
func (h *AgentGroupAdvancedHandler) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")

	template, err := h.db.GetGroupTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CreateTemplate creates a new group template
func (h *AgentGroupAdvancedHandler) CreateTemplate(c *gin.Context) {
	var template GroupTemplate

	if err := c.BindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if template.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Template name is required"})
		return
	}

	if err := h.db.CreateGroupTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Template created successfully",
		"template_id": template.ID,
	})
}

// DeleteTemplate deletes a template
func (h *AgentGroupAdvancedHandler) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")

	if err := h.db.DeleteGroupTemplate(c.Request.Context(), templateID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Template deleted successfully",
	})
}

// ApplyTemplate applies a template to create a group
func (h *AgentGroupAdvancedHandler) ApplyTemplate(c *gin.Context) {
	templateID := c.Param("id")

	var req struct {
		GroupName   string            `json:"group_name"`
		Description string            `json:"description"`
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

	groupID, err := h.db.ApplyGroupTemplate(c.Request.Context(), templateID, req.GroupName, req.Description, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Template applied successfully",
		"group_id": groupID,
	})
}

// SetGroupHierarchy sets the parent-child relationship
func (h *AgentGroupAdvancedHandler) SetGroupHierarchy(c *gin.Context) {
	groupID := c.Param("name")

	var req struct {
		ParentID string `json:"parent_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.db.SetGroupHierarchy(c.Request.Context(), groupID, req.ParentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set hierarchy: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hierarchy set successfully",
	})
}

// GetGroupHierarchy returns the hierarchy for a group
func (h *AgentGroupAdvancedHandler) GetGroupHierarchy(c *gin.Context) {
	groupID := c.Param("name")

	hierarchy, err := h.db.GetGroupHierarchy(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hierarchy not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, hierarchy)
}

// RemoveGroupHierarchy removes a group from hierarchy
func (h *AgentGroupAdvancedHandler) RemoveGroupHierarchy(c *gin.Context) {
	groupID := c.Param("name")

	if err := h.db.RemoveGroupHierarchy(c.Request.Context(), groupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove hierarchy: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hierarchy removed successfully",
	})
}

// GetGroupChildren returns child groups
func (h *AgentGroupAdvancedHandler) GetGroupChildren(c *gin.Context) {
	groupID := c.Param("name")

	children, err := h.db.GetGroupChildren(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get children: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"children": children})
}

// ListAutoDiscoveryConfigs returns all auto-discovery configurations
func (h *AgentGroupAdvancedHandler) ListAutoDiscoveryConfigs(c *gin.Context) {
	configs, err := h.db.ListAutoDiscoveryConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list configs: " + err.Error()})
		return
	}

	if configs == nil {
		configs = []*AutoDiscoveryConfig{}
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// GetAutoDiscoveryConfig returns a specific config
func (h *AgentGroupAdvancedHandler) GetAutoDiscoveryConfig(c *gin.Context) {
	configID := c.Param("id")

	config, err := h.db.GetAutoDiscoveryConfig(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// CreateAutoDiscoveryConfig creates a new auto-discovery config
func (h *AgentGroupAdvancedHandler) CreateAutoDiscoveryConfig(c *gin.Context) {
	var config AutoDiscoveryConfig

	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if config.TargetGroup == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target group is required"})
		return
	}

	if err := h.db.CreateAutoDiscoveryConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create config: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Config created successfully",
		"config_id": config.ID,
	})
}

// UpdateAutoDiscoveryConfig updates a config
func (h *AgentGroupAdvancedHandler) UpdateAutoDiscoveryConfig(c *gin.Context) {
	configID := c.Param("id")

	var config AutoDiscoveryConfig
	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	config.ID = configID

	if err := h.db.UpdateAutoDiscoveryConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Config updated successfully",
	})
}

// DeleteAutoDiscoveryConfig deletes a config
func (h *AgentGroupAdvancedHandler) DeleteAutoDiscoveryConfig(c *gin.Context) {
	configID := c.Param("id")

	if err := h.db.DeleteAutoDiscoveryConfig(c.Request.Context(), configID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete config: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Config deleted successfully",
	})
}

// RunAutoDiscovery manually triggers auto-discovery
func (h *AgentGroupAdvancedHandler) RunAutoDiscovery(c *gin.Context) {
	configID := c.Param("id")

	count, err := h.db.RunAutoDiscovery(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run discovery: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Auto-discovery completed",
		"agents_added": count,
	})
}

// ListWebhooks returns all webhook configurations
func (h *AgentGroupAdvancedHandler) ListWebhooks(c *gin.Context) {
	webhooks, err := h.db.ListWebhooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list webhooks: " + err.Error()})
		return
	}

	if webhooks == nil {
		webhooks = []*WebhookConfig{}
	}

	c.JSON(http.StatusOK, gin.H{"webhooks": webhooks})
}

// GetWebhook returns a specific webhook
func (h *AgentGroupAdvancedHandler) GetWebhook(c *gin.Context) {
	webhookID := c.Param("id")

	webhook, err := h.db.GetWebhook(c.Request.Context(), webhookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

// CreateWebhook creates a new webhook
func (h *AgentGroupAdvancedHandler) CreateWebhook(c *gin.Context) {
	var webhook WebhookConfig

	if err := c.BindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if webhook.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook URL is required"})
		return
	}

	if err := h.db.CreateWebhook(c.Request.Context(), &webhook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create webhook: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Webhook created successfully",
		"webhook_id": webhook.ID,
	})
}

// UpdateWebhook updates a webhook
func (h *AgentGroupAdvancedHandler) UpdateWebhook(c *gin.Context) {
	webhookID := c.Param("id")

	var webhook WebhookConfig
	if err := c.BindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	webhook.ID = webhookID

	if err := h.db.UpdateWebhook(c.Request.Context(), &webhook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update webhook: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Webhook updated successfully",
	})
}

// DeleteWebhook deletes a webhook
func (h *AgentGroupAdvancedHandler) DeleteWebhook(c *gin.Context) {
	webhookID := c.Param("id")

	if err := h.db.DeleteWebhook(c.Request.Context(), webhookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete webhook: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Webhook deleted successfully",
	})
}

// GetWebhookLogs returns webhook delivery logs
func (h *AgentGroupAdvancedHandler) GetWebhookLogs(c *gin.Context) {
	webhookID := c.Param("id")

	logs, err := h.db.GetWebhookLogs(c.Request.Context(), webhookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs: " + err.Error()})
		return
	}

	if logs == nil {
		logs = []*WebhookLog{}
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
