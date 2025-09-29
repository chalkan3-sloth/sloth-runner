package ui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	httpServer *http.Server
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
}

type TaskStatus struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Agent  string `json:"agent,omitempty"`
}

type AgentStatus struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (s *Server) Start(port int) error {
	router := mux.NewRouter()

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	// Serve index.html at root
	router.HandleFunc("/", s.handleIndex).Methods("GET")
	
	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tasks", s.handleTasks).Methods("GET", "POST")
	api.HandleFunc("/tasks/stop-all", s.handleStopAllTasks).Methods("POST")
	api.HandleFunc("/agents", s.handleAgents).Methods("GET", "POST")
	api.HandleFunc("/status", s.handleStatus).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", s.handleWebSocket)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start broadcast handler
	go s.handleBroadcast()

	slog.Info(fmt.Sprintf("Starting UI server on http://localhost:%d", port))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	content, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Failed to read index.html", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "GET":
		// Mock data for now - in real implementation, this would query actual tasks
		tasks := []TaskStatus{
			{ID: "1", Name: "Build Project", Status: "completed", Type: "shell"},
			{ID: "2", Name: "Run Tests", Status: "running", Type: "lua"},
			{ID: "3", Name: "Deploy to Staging", Status: "pending", Type: "pipeline"},
			{ID: "4", Name: "Backup Database", Status: "failed", Type: "shell"},
		}
		
		response := APIResponse{
			Success: true,
			Data:    map[string]interface{}{"tasks": tasks},
		}
		json.NewEncoder(w).Encode(response)
		
	case "POST":
		var task TaskStatus
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			response := APIResponse{
				Success: false,
				Error:   "Invalid task data",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		
		// In real implementation, this would create and execute the task
		task.ID = fmt.Sprintf("%d", time.Now().Unix())
		task.Status = "pending"
		
		response := APIResponse{
			Success: true,
			Data:    task,
		}
		json.NewEncoder(w).Encode(response)
		
		// Broadcast task creation
		s.broadcastMessage(map[string]interface{}{
			"type":    "task_update",
			"taskId":  task.ID,
			"status":  task.Status,
			"message": fmt.Sprintf("Task '%s' created", task.Name),
		})
	}
}

func (s *Server) handleStopAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// In real implementation, this would stop all running tasks
	response := APIResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "All tasks stopped"},
	}
	json.NewEncoder(w).Encode(response)
	
	// Broadcast stop all message
	s.broadcastMessage(map[string]interface{}{
		"type":    "console_output",
		"message": "All running tasks have been stopped",
	})
}

func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "GET":
		// Mock data for now - in real implementation, this would query actual agents
		agents := []AgentStatus{
			{Name: "local-agent", Address: "localhost:50051", Status: "online"},
			{Name: "prod-server-1", Address: "192.168.1.100:50051", Status: "online"},
			{Name: "test-server", Address: "192.168.1.101:50051", Status: "offline"},
		}
		
		response := APIResponse{
			Success: true,
			Data:    map[string]interface{}{"agents": agents},
		}
		json.NewEncoder(w).Encode(response)
		
	case "POST":
		var agent AgentStatus
		if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
			response := APIResponse{
				Success: false,
				Error:   "Invalid agent data",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		
		// In real implementation, this would register the agent
		agent.Status = "online"
		
		response := APIResponse{
			Success: true,
			Data:    agent,
		}
		json.NewEncoder(w).Encode(response)
		
		// Broadcast agent addition
		s.broadcastMessage(map[string]interface{}{
			"type":      "agent_update",
			"agentName": agent.Name,
			"status":    agent.Status,
		})
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	status := map[string]interface{}{
		"server":    "running",
		"version":   "1.0.0",
		"uptime":    time.Now().Format(time.RFC3339),
		"connected": len(s.clients),
	}
	
	response := APIResponse{
		Success: true,
		Data:    status,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	s.clients[conn] = true
	slog.Info("New WebSocket client connected", "remote_addr", r.RemoteAddr)

	// Send welcome message
	welcomeMsg := map[string]interface{}{
		"type":    "console_output",
		"message": "Connected to Sloth Runner WebSocket",
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		slog.Error("Failed to send welcome message", "error", err)
		delete(s.clients, conn)
		return
	}

	// Listen for messages from client
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("WebSocket error", "error", err)
			}
			break
		}
		
		// Handle client messages if needed
		slog.Debug("Received WebSocket message", "message", msg)
	}

	delete(s.clients, conn)
	slog.Info("WebSocket client disconnected", "remote_addr", r.RemoteAddr)
}

func (s *Server) handleBroadcast() {
	for {
		select {
		case message := <-s.broadcast:
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					slog.Error("Failed to write WebSocket message", "error", err)
					client.Close()
					delete(s.clients, client)
				}
			}
		}
	}
}

func (s *Server) broadcastMessage(message map[string]interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		slog.Error("Failed to marshal broadcast message", "error", err)
		return
	}
	
	select {
	case s.broadcast <- data:
	default:
		slog.Warn("Broadcast channel is full, dropping message")
	}
}

// BroadcastTaskUpdate sends a task status update to all connected clients
func (s *Server) BroadcastTaskUpdate(taskID, status, output string) {
	s.broadcastMessage(map[string]interface{}{
		"type":   "task_update",
		"taskId": taskID,
		"status": status,
		"output": output,
	})
}

// BroadcastAgentUpdate sends an agent status update to all connected clients
func (s *Server) BroadcastAgentUpdate(agentName, status string) {
	s.broadcastMessage(map[string]interface{}{
		"type":      "agent_update",
		"agentName": agentName,
		"status":    status,
	})
}

// BroadcastConsoleOutput sends console output to all connected clients
func (s *Server) BroadcastConsoleOutput(message string) {
	s.broadcastMessage(map[string]interface{}{
		"type":    "console_output",
		"message": message,
	})
}