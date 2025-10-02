package ui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestNewServer(t *testing.T) {
	server := NewServer()
	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.clients == nil {
		t.Error("clients map not initialized")
	}
	if server.broadcast == nil {
		t.Error("broadcast channel not initialized")
	}
}

func TestHandleStatus(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/api/status", nil)
	rr := httptest.NewRecorder()

	server.handleStatus(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestHandleTasks(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/api/tasks", nil)
	rr := httptest.NewRecorder()

	server.handleTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestHandleAgents(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/api/agents", nil)
	rr := httptest.NewRecorder()

	server.handleAgents(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestHandleIndex(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	server.handleIndex(rr, req)

	// Should redirect or serve content
	if status := rr.Code; status != http.StatusOK && status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestRouterSetup(t *testing.T) {
	server := NewServer()
	router := mux.NewRouter()

	// Setup routes similar to Start method
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tasks", server.handleTasks).Methods("GET", "POST")
	api.HandleFunc("/agents", server.handleAgents).Methods("GET", "POST")
	api.HandleFunc("/status", server.handleStatus).Methods("GET")

	tests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/api/status", http.StatusOK},
		{"GET", "/api/tasks", http.StatusOK},
		{"GET", "/api/agents", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
			}
		})
	}
}

func TestAPIResponse(t *testing.T) {
	tests := []struct {
		name     string
		response APIResponse
		wantJSON string
	}{
		{
			name: "success response",
			response: APIResponse{
				Success: true,
				Data:    map[string]string{"key": "value"},
			},
			wantJSON: `{"success":true,"data":{"key":"value"}}`,
		},
		{
			name: "error response",
			response: APIResponse{
				Success: false,
				Error:   "test error",
			},
			wantJSON: `{"success":false,"error":"test error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}

			var decoded APIResponse
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if decoded.Success != tt.response.Success {
				t.Errorf("Success mismatch: got %v want %v", decoded.Success, tt.response.Success)
			}
		})
	}
}

func TestTaskStatus(t *testing.T) {
	task := TaskStatus{
		ID:     "task-1",
		Name:   "test-task",
		Status: "running",
		Type:   "local",
		Agent:  "agent-1",
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	var decoded TaskStatus
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	if decoded.ID != task.ID {
		t.Errorf("ID mismatch: got %v want %v", decoded.ID, task.ID)
	}
	if decoded.Name != task.Name {
		t.Errorf("Name mismatch: got %v want %v", decoded.Name, task.Name)
	}
	if decoded.Status != task.Status {
		t.Errorf("Status mismatch: got %v want %v", decoded.Status, task.Status)
	}
}

func TestAgentStatus(t *testing.T) {
	agent := AgentStatus{
		Name:    "test-agent",
		Address: "localhost:50051",
		Status:  "active",
	}

	data, err := json.Marshal(agent)
	if err != nil {
		t.Fatalf("Failed to marshal agent: %v", err)
	}

	var decoded AgentStatus
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal agent: %v", err)
	}

	if decoded.Name != agent.Name {
		t.Errorf("Name mismatch: got %v want %v", decoded.Name, agent.Name)
	}
	if decoded.Address != agent.Address {
		t.Errorf("Address mismatch: got %v want %v", decoded.Address, agent.Address)
	}
	if decoded.Status != agent.Status {
		t.Errorf("Status mismatch: got %v want %v", decoded.Status, agent.Status)
	}
}

func TestServerStop(t *testing.T) {
	server := NewServer()
	
	// Start server in goroutine
	go func() {
		// Use a high port that's likely available
		err := server.Start(59999)
		if err != nil && err != http.ErrServerClosed {
			t.Logf("Server start error (expected): %v", err)
		}
	}()

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Stop the server
	err := server.Stop()
	if err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}
}

func TestBroadcastTaskUpdate(t *testing.T) {
	server := NewServer()
	
	// This should not panic
	server.BroadcastTaskUpdate("task-1", "running", "test output")
}

func TestBroadcastAgentUpdate(t *testing.T) {
	server := NewServer()
	
	// This should not panic
	server.BroadcastAgentUpdate("agent-1", "active")
}

func TestBroadcastConsoleOutput(t *testing.T) {
	server := NewServer()
	
	// This should not panic
	server.BroadcastConsoleOutput("test console message")
}

func TestHandleStopAllTasks(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("POST", "/api/tasks/stop-all", nil)
	rr := httptest.NewRecorder()

	server.handleStopAllTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestHandleTasksPost(t *testing.T) {
	server := NewServer()

	taskJSON := `{"name":"test-task","type":"shell"}`
	req := httptest.NewRequest("POST", "/api/tasks", strings.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.handleTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestHandleTasksPostInvalid(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("POST", "/api/tasks", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.handleTasks(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}
}

func TestHandleAgentsPost(t *testing.T) {
	server := NewServer()

	agentJSON := `{"name":"test-agent","address":"localhost:50051"}`
	req := httptest.NewRequest("POST", "/api/agents", strings.NewReader(agentJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.handleAgents(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestHandleAgentsPostInvalid(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("POST", "/api/agents", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.handleAgents(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response APIResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}
}

func TestBroadcastMessage(t *testing.T) {
	server := NewServer()
	
	// Test broadcasting with no clients
	server.broadcastMessage(map[string]interface{}{
		"type":    "test",
		"message": "test message",
	})
	
	// Should not panic
}

func TestServerStopWithoutStart(t *testing.T) {
	server := NewServer()
	
	// Should not panic when stopping without starting
	err := server.Stop()
	if err != nil {
		t.Errorf("Unexpected error when stopping unstarted server: %v", err)
	}
}
