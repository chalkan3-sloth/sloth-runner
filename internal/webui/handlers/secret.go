package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecretHandler handles secret operations
type SecretHandler struct {
	svc *SecretsServiceWrapper
}

// NewSecretHandler creates a new secret handler
func NewSecretHandler(svc *SecretsServiceWrapper) *SecretHandler {
	return &SecretHandler{svc: svc}
}

// List returns all secrets for a stack (names only, not values)
func (h *SecretHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	stack := c.Param("stack")

	secrets, err := h.svc.ListSecrets(ctx, stack)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secrets": secrets})
}

// Add adds a new secret (placeholder - requires encryption implementation)
func (h *SecretHandler) Add(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Secret addition requires encryption keys - use CLI",
	})
}

// Delete removes a secret
func (h *SecretHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Secret deletion requires encryption keys - use CLI",
	})
}
