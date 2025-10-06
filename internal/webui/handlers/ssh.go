package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/chalkan3-sloth/sloth-runner/internal/ssh"
)

// SSHHandler handles SSH profile operations
type SSHHandler struct {
	db *SSHDBWrapper
}

// NewSSHHandler creates a new SSH handler
func NewSSHHandler(db *SSHDBWrapper) *SSHHandler {
	return &SSHHandler{db: db}
}

// List returns all SSH profiles
func (h *SSHHandler) List(c *gin.Context) {
	profiles, err := h.db.ListProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

// Get returns an SSH profile by name
func (h *SSHHandler) Get(c *gin.Context) {
	name := c.Param("name")

	profile, err := h.db.GetProfile(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// CreateSSHProfileRequest represents a create SSH profile request
type CreateSSHProfileRequest struct {
	Name               string `json:"name" binding:"required"`
	Host               string `json:"host" binding:"required"`
	User               string `json:"user" binding:"required"`
	Port               int    `json:"port"`
	KeyPath            string `json:"key_path"`
	Description        string `json:"description"`
	ConnectionTimeout  int    `json:"connection_timeout"`
	KeepaliveInterval  int    `json:"keepalive_interval"`
	StrictHostChecking bool   `json:"strict_host_checking"`
}

// Create creates a new SSH profile
func (h *SSHHandler) Create(c *gin.Context) {
	var req CreateSSHProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.Port == 0 {
		req.Port = 22
	}
	if req.ConnectionTimeout == 0 {
		req.ConnectionTimeout = 30
	}
	if req.KeepaliveInterval == 0 {
		req.KeepaliveInterval = 60
	}

	profile := &ssh.Profile{
		Name:               req.Name,
		Host:               req.Host,
		User:               req.User,
		Port:               req.Port,
		KeyPath:            req.KeyPath,
		Description:        req.Description,
		ConnectionTimeout:  req.ConnectionTimeout,
		KeepaliveInterval:  req.KeepaliveInterval,
		StrictHostChecking: req.StrictHostChecking,
	}

	if err := h.db.AddProfile(profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// Update updates an SSH profile
func (h *SSHHandler) Update(c *gin.Context) {
	name := c.Param("name")

	var req CreateSSHProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	updates["host"] = req.Host
	updates["user"] = req.User
	updates["port"] = req.Port
	updates["key_path"] = req.KeyPath
	updates["description"] = req.Description
	updates["connection_timeout"] = req.ConnectionTimeout
	updates["keepalive_interval"] = req.KeepaliveInterval
	updates["strict_host_checking"] = req.StrictHostChecking

	if err := h.db.UpdateProfile(name, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// Delete removes an SSH profile
func (h *SSHHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	if err := h.db.RemoveProfile(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

// GetAuditLogs returns audit logs for a profile
func (h *SSHHandler) GetAuditLogs(c *gin.Context) {
	name := c.Param("name")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.db.GetAuditLogs(name, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
