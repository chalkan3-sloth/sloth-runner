package core

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestNewNotificationModule(t *testing.T) {
	module := NewNotificationModule()
	
	if module == nil {
		t.Fatal("Expected module to be created")
	}
	
	if module.info.Name != "notify" {
		t.Errorf("Expected name 'notify', got '%s'", module.info.Name)
	}
	
	if module.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestNotificationModuleLoader(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	
	L.PreloadModule("notify", module.Loader)
	
	// Load the module
	if err := L.DoString(`notify = require("notify")`); err != nil {
		t.Fatalf("Failed to load module: %v", err)
	}
	
	// Check if functions are available
	functions := []string{"slack", "discord", "email", "webhook", "teams", "telegram"}
	for _, fn := range functions {
		if err := L.DoString(`assert(notify.` + fn + ` ~= nil, "` + fn + ` function not found")`); err != nil {
			t.Errorf("Function %s not found: %v", fn, err)
		}
	}
}

func TestNotificationSlack(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		body, _ := io.ReadAll(r.Body)
		var msg SlackMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("Failed to parse message: %v", err)
		}
		
		if msg.Text != "Test message" {
			t.Errorf("Expected text 'Test message', got '%s'", msg.Text)
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	script := `
		notify = require("notify")
		local success, err = notify.slack("` + server.URL + `", {
			text = "Test message",
			username = "Test Bot",
			channel = "#general"
		})
		assert(success == true, "Expected success")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestNotificationSlackWithAttachments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var msg SlackMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("Failed to parse message: %v", err)
		}
		
		if len(msg.Attachments) == 0 {
			t.Error("Expected attachments")
		}
		
		if msg.Attachments[0].Title != "Alert" {
			t.Errorf("Expected title 'Alert', got '%s'", msg.Attachments[0].Title)
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	script := `
		notify = require("notify")
		local success = notify.slack("` + server.URL + `", {
			text = "Test",
			attachments = {
				{
					color = "danger",
					title = "Alert",
					text = "Something went wrong"
				}
			}
		})
		assert(success == true)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestNotificationDiscord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var msg DiscordMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("Failed to parse message: %v", err)
		}
		
		if msg.Content != "Discord test" {
			t.Errorf("Expected content 'Discord test', got '%s'", msg.Content)
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	script := `
		notify = require("notify")
		local success = notify.discord("` + server.URL + `", {
			content = "Discord test",
			username = "Test Bot"
		})
		assert(success == true)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestNotificationDiscordWithEmbeds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var msg DiscordMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("Failed to parse message: %v", err)
		}
		
		if len(msg.Embeds) == 0 {
			t.Error("Expected embeds")
		}
		
		if msg.Embeds[0].Title != "Embed Title" {
			t.Errorf("Expected title 'Embed Title', got '%s'", msg.Embeds[0].Title)
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	script := `
		notify = require("notify")
		local success = notify.discord("` + server.URL + `", {
			embeds = {
				{
					title = "Embed Title",
					description = "Embed Description",
					color = 16711680
				}
			}
		})
		assert(success == true)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestNotificationWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("Failed to parse payload: %v", err)
		}
		
		if payload["message"] != "Test webhook" {
			t.Errorf("Expected message 'Test webhook', got '%v'", payload["message"])
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	script := `
		notify = require("notify")
		local success = notify.webhook("` + server.URL + `", {
			message = "Test webhook",
			level = "info"
		})
		assert(success == true)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestNotificationTeams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var msg map[string]interface{}
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("Failed to parse message: %v", err)
		}
		
		if msg["text"] != "Teams message" {
			t.Errorf("Expected text 'Teams message', got '%v'", msg["text"])
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	script := `
		notify = require("notify")
		local success = notify.teams("` + server.URL + `", {
			text = "Teams message",
			title = "Test Alert"
		})
		assert(success == true)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestNotificationTelegram(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var msg map[string]interface{}
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("Failed to parse message: %v", err)
		}
		
		if msg["text"] != "Telegram message" {
			t.Errorf("Expected text 'Telegram message', got '%v'", msg["text"])
		}
		
		if msg["chat_id"] != "123456" {
			t.Errorf("Expected chat_id '123456', got '%v'", msg["chat_id"])
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	// Replace the Telegram API URL with our test server
	script := `
		notify = require("notify")
		local url = "` + server.URL + `/botTOKEN/sendMessage"
		local success = notify.webhook(url, {
			chat_id = "123456",
			text = "Telegram message"
		})
		assert(success == true)
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestNotificationFailure(t *testing.T) {
	// Server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	L.PreloadModule("notify", module.Loader)
	
	script := `
		notify = require("notify")
		local success, err = notify.slack("` + server.URL + `", {
			text = "Test"
		})
		assert(success == false, "Expected failure")
		assert(err ~= nil, "Expected error message")
	`
	
	if err := L.DoString(script); err != nil {
		t.Fatalf("Script failed: %v", err)
	}
}

func TestLuaValueToGo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	module := NewNotificationModule()
	
	tests := []struct {
		name     string
		lua      string
		expected interface{}
	}{
		{"string", `"test"`, "test"},
		{"number", `42`, 42.0},
		{"bool", `true`, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(`value = ` + tt.lua); err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}
			
			value := L.GetGlobal("value")
			result := module.luaValueToGo(value)
			
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNotificationModuleInfo(t *testing.T) {
	module := NewNotificationModule()
	info := module.Info()
	
	if info.Name != "notify" {
		t.Errorf("Expected name 'notify', got '%s'", info.Name)
	}
	
	if info.Category != "core" {
		t.Errorf("Expected category 'core', got '%s'", info.Category)
	}
	
	if info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", info.Version)
	}
}
