package types

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateTaskID(t *testing.T) {
	id := GenerateTaskID()
	if id == "" {
		t.Error("Expected non-empty task ID")
	}
	
	// Validate it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("Generated ID is not a valid UUID: %s", id)
	}
}

func TestGenerateTaskGroupID(t *testing.T) {
	id := GenerateTaskGroupID()
	if id == "" {
		t.Error("Expected non-empty task group ID")
	}
	
	// Validate it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("Generated ID is not a valid UUID: %s", id)
	}
}

func TestGenerateTaskID_Uniqueness(t *testing.T) {
	id1 := GenerateTaskID()
	id2 := GenerateTaskID()
	
	if id1 == id2 {
		t.Error("Expected unique task IDs, got identical IDs")
	}
}

func TestGenerateTaskGroupID_Uniqueness(t *testing.T) {
	id1 := GenerateTaskGroupID()
	id2 := GenerateTaskGroupID()
	
	if id1 == id2 {
		t.Error("Expected unique task group IDs, got identical IDs")
	}
}

func TestTask_Struct(t *testing.T) {
	task := Task{
		ID:          "test-id",
		Name:        "test-task",
		Description: "Test Description",
		Workdir:     "/tmp/test",
		DependsOn:   []string{"dep1", "dep2"},
		Artifacts:   []string{"artifact1"},
		Params:      map[string]string{"key": "value"},
		Retries:     3,
		Timeout:     "5m",
		Async:       true,
	}
	
	if task.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", task.ID)
	}
	if task.Name != "test-task" {
		t.Errorf("Expected Name 'test-task', got '%s'", task.Name)
	}
	if task.Workdir != "/tmp/test" {
		t.Errorf("Expected Workdir '/tmp/test', got '%s'", task.Workdir)
	}
	if len(task.DependsOn) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(task.DependsOn))
	}
	if task.Retries != 3 {
		t.Errorf("Expected Retries 3, got %d", task.Retries)
	}
	if !task.Async {
		t.Error("Expected Async to be true")
	}
}

func TestTaskGroup_Struct(t *testing.T) {
	taskGroup := TaskGroup{
		ID:                       "group-id",
		Description:              "Test Group",
		Tasks:                    []Task{},
		Workdir:                  "/tmp/group",
		CreateWorkdirBeforeRun:   true,
	}
	
	if taskGroup.ID != "group-id" {
		t.Errorf("Expected ID 'group-id', got '%s'", taskGroup.ID)
	}
	if taskGroup.Description != "Test Group" {
		t.Errorf("Expected Description 'Test Group', got '%s'", taskGroup.Description)
	}
	if !taskGroup.CreateWorkdirBeforeRun {
		t.Error("Expected CreateWorkdirBeforeRun to be true")
	}
}

func TestTaskResult_Struct(t *testing.T) {
	result := TaskResult{
		Name:     "test-task",
		Status:   "success",
		Error:    nil,
	}
	
	if result.Name != "test-task" {
		t.Errorf("Expected Name 'test-task', got '%s'", result.Name)
	}
	if result.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", result.Status)
	}
	if result.Error != nil {
		t.Error("Expected no error")
	}
}

func TestPythonVenv_Struct(t *testing.T) {
	venv := PythonVenv{
		Path: "/path/to/venv",
	}
	
	if venv.Path != "/path/to/venv" {
		t.Errorf("Expected Path '/path/to/venv', got '%s'", venv.Path)
	}
}
