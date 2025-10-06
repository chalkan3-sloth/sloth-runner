package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/chalkan3-sloth/sloth-runner/internal/sloth"
)

// SlothHandler handles workflow/sloth operations
type SlothHandler struct {
	repo *SlothRepoWrapper
}

// NewSlothHandler creates a new sloth handler
func NewSlothHandler(repo *SlothRepoWrapper) *SlothHandler {
	return &SlothHandler{repo: repo}
}

// List returns all sloths
func (h *SlothHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	activeOnly := c.Query("active") == "true"

	sloths, err := h.repo.List(ctx, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sloths": sloths})
}

// Get returns a sloth by name
func (h *SlothHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	s, err := h.repo.Get(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// CreateSlothRequest represents a create sloth request
type CreateSlothRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	FilePath    string `json:"file_path" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Tags        string `json:"tags"`
}

// Create creates a new sloth
func (h *SlothHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateSlothRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate file hash
	hash := sha256.Sum256([]byte(req.Content))
	fileHash := fmt.Sprintf("%x", hash)

	now := time.Now()
	s := &sloth.Sloth{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		FilePath:    req.FilePath,
		Content:     req.Content,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		UsageCount:  0,
		Tags:        req.Tags,
		FileHash:    fileHash,
	}

	if err := h.repo.Create(ctx, s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, s)
}

// Update updates a sloth
func (h *SlothHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	var req CreateSlothRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing sloth
	existing, err := h.repo.Get(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	// Calculate file hash
	hash := sha256.Sum256([]byte(req.Content))
	fileHash := fmt.Sprintf("%x", hash)

	existing.Description = req.Description
	existing.FilePath = req.FilePath
	existing.Content = req.Content
	existing.Tags = req.Tags
	existing.FileHash = fileHash
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// Delete deletes a sloth
func (h *SlothHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if err := h.repo.Delete(ctx, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow deleted successfully"})
}

// Activate activates a sloth
func (h *SlothHandler) Activate(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if err := h.repo.SetActive(ctx, name, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow activated successfully"})
}

// Deactivate deactivates a sloth
func (h *SlothHandler) Deactivate(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if err := h.repo.SetActive(ctx, name, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow deactivated successfully"})
}

// Run executes a sloth (placeholder for future implementation)
func (h *SlothHandler) Run(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Workflow execution not implemented in UI yet"})
}
