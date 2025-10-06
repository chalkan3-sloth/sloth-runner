package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BasicAuth provides HTTP Basic Authentication
type BasicAuth struct {
	username string
	password string
}

// NewBasicAuth creates a new BasicAuth middleware
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

// Middleware returns a Gin middleware function
func (ba *BasicAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for health check
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		user, pass, ok := c.Request.BasicAuth()
		if !ok {
			ba.unauthorized(c)
			return
		}

		// Constant-time comparison to prevent timing attacks
		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(ba.username)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(ba.password)) == 1

		if !userMatch || !passMatch {
			ba.unauthorized(c)
			return
		}

		c.Next()
	}
}

func (ba *BasicAuth) unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", `Basic realm="Sloth Runner UI"`)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
}
