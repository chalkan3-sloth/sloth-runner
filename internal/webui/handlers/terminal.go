package handlers

import (
	"fmt"
	"io"
	"net/http"
	osexec "os/exec"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TerminalHandler handles web terminal sessions
type TerminalHandler struct {
	sessions map[string]*TerminalSession
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

// TerminalSession represents a terminal session
type TerminalSession struct {
	ID       string
	AgentID  string
	cmd      *osexec.Cmd
	stdin    io.WriteCloser
	conn     *websocket.Conn
	done     chan struct{}
}

// NewTerminalHandler creates a new terminal handler
func NewTerminalHandler() *TerminalHandler {
	return &TerminalHandler{
		sessions: make(map[string]*TerminalSession),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// CreateSession creates a new terminal session
func (h *TerminalHandler) CreateSession(c *gin.Context) {
	var req struct{
		AgentID string `json:"agent_id"`
		Command string `json:"command"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := uuid.New().String()

	session := &TerminalSession{
		ID:      sessionID,
		AgentID: req.AgentID,
		done:    make(chan struct{}),
	}

	h.mu.Lock()
	h.sessions[sessionID] = session
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"message":    "Terminal session created",
	})
}

// ConnectTerminal handles WebSocket connection for terminal
func (h *TerminalHandler) ConnectTerminal(c *gin.Context) {
	sessionID := c.Param("id")

	h.mu.RLock()
	session, exists := h.sessions[sessionID]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	session.conn = conn

	// Start shell process
	session.cmd = osexec.Command("sh", "-i")
	session.stdin, _ = session.cmd.StdinPipe()
	stdout, _ := session.cmd.StdoutPipe()
	stderr, _ := session.cmd.StderrPipe()

	if err := session.cmd.Start(); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %s\n", err.Error())))
		conn.Close()
		return
	}

	// Read from stdout
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				return
			}
			conn.WriteMessage(websocket.TextMessage, buf[:n])
		}
	}()

	// Read from stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				return
			}
			conn.WriteMessage(websocket.TextMessage, buf[:n])
		}
	}()

	// Read from WebSocket and write to stdin
	go func() {
		defer func() {
			session.stdin.Close()
			session.cmd.Process.Kill()
			conn.Close()
			close(session.done)
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			session.stdin.Write(message)
		}
	}()

	<-session.done
}

// CloseSession closes a terminal session
func (h *TerminalHandler) CloseSession(c *gin.Context) {
	sessionID := c.Param("id")

	h.mu.Lock()
	session, exists := h.sessions[sessionID]
	if exists {
		if session.cmd != nil && session.cmd.Process != nil {
			session.cmd.Process.Kill()
		}
		if session.conn != nil {
			session.conn.Close()
		}
		delete(h.sessions, sessionID)
	}
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session closed"})
}

// ListSessions lists all active terminal sessions
func (h *TerminalHandler) ListSessions(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sessions := make([]gin.H, 0, len(h.sessions))
	for id, session := range h.sessions {
		sessions = append(sessions, gin.H{
			"id":       id,
			"agent_id": session.AgentID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}
