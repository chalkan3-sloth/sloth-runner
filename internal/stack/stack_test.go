package stack

import (
"fmt"
"path/filepath"
"testing"
"time"
)

func TestNewStackManager(t *testing.T) {
tmpDir := t.TempDir()
dbPath := filepath.Join(tmpDir, "test_stacks.db")

sm, err := NewStackManager(dbPath)
if err != nil {
t.Fatalf("Failed to create stack manager: %v", err)
}
defer sm.Close()

if sm == nil {
t.Fatal("NewStackManager returned nil")
}

if sm.db == nil {
t.Error("database not initialized")
}

if sm.path != dbPath {
t.Errorf("Expected path %s, got %s", dbPath, sm.path)
}
}

func TestStackStateFields(t *testing.T) {
now := time.Now()
state := &StackState{
ID:             "stack-1",
Name:           "test-stack",
Description:    "Test stack description",
Version:        "1.0.0",
Status:         "created",
CreatedAt:      now,
UpdatedAt:      now,
WorkflowFile:   "/path/to/workflow.sloth",
TaskResults:    make(map[string]interface{}),
Outputs:        make(map[string]interface{}),
Configuration:  make(map[string]interface{}),
Metadata:       make(map[string]interface{}),
ExecutionCount: 0,
}

if state.ID != "stack-1" {
t.Errorf("Expected ID 'stack-1', got '%s'", state.ID)
}

if state.Name != "test-stack" {
t.Errorf("Expected Name 'test-stack', got '%s'", state.Name)
}

if state.Status != "created" {
t.Errorf("Expected Status 'created', got '%s'", state.Status)
}

if state.Version != "1.0.0" {
t.Errorf("Expected Version '1.0.0', got '%s'", state.Version)
}
}

func TestStackManagerCreateStack(t *testing.T) {
tmpDir := t.TempDir()
dbPath := filepath.Join(tmpDir, "test_stacks.db")

sm, err := NewStackManager(dbPath)
if err != nil {
t.Fatalf("Failed to create stack manager: %v", err)
}
defer sm.Close()

state := &StackState{
ID:            "test-stack-1",
Name:          "TestStack",
Description:   "Test description",
Version:       "1.0.0",
Status:        "created",
WorkflowFile:  "/test/workflow.sloth",
TaskResults:   make(map[string]interface{}),
Outputs:       make(map[string]interface{}),
Configuration: make(map[string]interface{}),
Metadata:      make(map[string]interface{}),
}

err = sm.CreateStack(state)
if err != nil {
t.Fatalf("Failed to create stack: %v", err)
}

// Verify stack was created
retrieved, err := sm.GetStack("test-stack-1")
if err != nil {
t.Fatalf("Failed to retrieve stack: %v", err)
}

if retrieved.ID != state.ID {
t.Errorf("Expected ID '%s', got '%s'", state.ID, retrieved.ID)
}

if retrieved.Name != state.Name {
t.Errorf("Expected Name '%s', got '%s'", state.Name, retrieved.Name)
}
}

func TestStackManagerUpdateStack(t *testing.T) {
tmpDir := t.TempDir()
dbPath := filepath.Join(tmpDir, "test_stacks.db")

sm, err := NewStackManager(dbPath)
if err != nil {
t.Fatalf("Failed to create stack manager: %v", err)
}
defer sm.Close()

// Create initial stack
state := &StackState{
ID:            "test-stack-2",
Name:          "TestStack",
Description:   "Initial description",
Version:       "1.0.0",
Status:        "created",
WorkflowFile:  "/test/workflow.sloth",
TaskResults:   make(map[string]interface{}),
Outputs:       make(map[string]interface{}),
Configuration: make(map[string]interface{}),
Metadata:      make(map[string]interface{}),
}

err = sm.CreateStack(state)
if err != nil {
t.Fatalf("Failed to create stack: %v", err)
}

// Update stack
state.Status = "running"
state.Description = "Updated description"

err = sm.UpdateStack(state)
if err != nil {
t.Fatalf("Failed to update stack: %v", err)
}

// Verify update
retrieved, err := sm.GetStack("test-stack-2")
if err != nil {
t.Fatalf("Failed to retrieve stack: %v", err)
}

if retrieved.Status != "running" {
t.Errorf("Expected Status 'running', got '%s'", retrieved.Status)
}

if retrieved.Description != "Updated description" {
t.Errorf("Expected Description 'Updated description', got '%s'", retrieved.Description)
}
}

func TestStackManagerListStacks(t *testing.T) {
tmpDir := t.TempDir()
dbPath := filepath.Join(tmpDir, "test_stacks.db")

sm, err := NewStackManager(dbPath)
if err != nil {
t.Fatalf("Failed to create stack manager: %v", err)
}
defer sm.Close()

// Create multiple stacks
for i := 1; i <= 3; i++ {
state := &StackState{
ID:            fmt.Sprintf("stack-%d", i),
Name:          fmt.Sprintf("Stack%d", i),
Description:   fmt.Sprintf("Description %d", i),
Version:       "1.0.0",
Status:        "created",
WorkflowFile:  "/test/workflow.sloth",
TaskResults:   make(map[string]interface{}),
Outputs:       make(map[string]interface{}),
Configuration: make(map[string]interface{}),
Metadata:      make(map[string]interface{}),
}

err = sm.CreateStack(state)
if err != nil {
t.Fatalf("Failed to create stack %d: %v", i, err)
}
}

// List stacks
stacks, err := sm.ListStacks()
if err != nil {
t.Fatalf("Failed to list stacks: %v", err)
}

if len(stacks) < 3 {
t.Errorf("Expected at least 3 stacks, got %d", len(stacks))
}
}

func TestStackManagerDeleteStack(t *testing.T) {
tmpDir := t.TempDir()
dbPath := filepath.Join(tmpDir, "test_stacks.db")

sm, err := NewStackManager(dbPath)
if err != nil {
t.Fatalf("Failed to create stack manager: %v", err)
}
defer sm.Close()

// Create stack
state := &StackState{
ID:            "stack-to-delete",
Name:          "DeleteStack",
Description:   "Will be deleted",
Version:       "1.0.0",
Status:        "created",
WorkflowFile:  "/test/workflow.sloth",
TaskResults:   make(map[string]interface{}),
Outputs:       make(map[string]interface{}),
Configuration: make(map[string]interface{}),
Metadata:      make(map[string]interface{}),
}

err = sm.CreateStack(state)
if err != nil {
t.Fatalf("Failed to create stack: %v", err)
}

// Delete stack
err = sm.DeleteStack("stack-to-delete")
if err != nil {
t.Fatalf("Failed to delete stack: %v", err)
}

// Verify deletion
_, err = sm.GetStack("stack-to-delete")
if err == nil {
t.Error("Expected error when getting deleted stack, got nil")
}
}

func TestStackManagerClose(t *testing.T) {
tmpDir := t.TempDir()
dbPath := filepath.Join(tmpDir, "test_stacks.db")

sm, err := NewStackManager(dbPath)
if err != nil {
t.Fatalf("Failed to create stack manager: %v", err)
}

err = sm.Close()
if err != nil {
t.Errorf("Failed to close stack manager: %v", err)
}
}
