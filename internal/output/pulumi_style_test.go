package output

import (
	"errors"
	"testing"
	"time"
)

func TestNewPulumiStyleOutput(t *testing.T) {
	output := NewPulumiStyleOutput()
	if output == nil {
		t.Fatal("NewPulumiStyleOutput returned nil")
	}
	if output.outputs == nil {
		t.Error("outputs map not initialized")
	}
	if output.indent != 0 {
		t.Errorf("expected indent 0, got %d", output.indent)
	}
}

func TestIndentation(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	if output.indent != 0 {
		t.Errorf("expected initial indent 0, got %d", output.indent)
	}
	
	output.Indent()
	if output.indent != 1 {
		t.Errorf("expected indent 1, got %d", output.indent)
	}
	
	output.Indent()
	if output.indent != 2 {
		t.Errorf("expected indent 2, got %d", output.indent)
	}
	
	output.Unindent()
	if output.indent != 1 {
		t.Errorf("expected indent 1 after unindent, got %d", output.indent)
	}
	
	output.Unindent()
	if output.indent != 0 {
		t.Errorf("expected indent 0, got %d", output.indent)
	}
	
	// Should not go below 0
	output.Unindent()
	if output.indent != 0 {
		t.Errorf("indent should not go below 0, got %d", output.indent)
	}
}

func TestAddAndGetOutputs(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	// Test adding string output
	output.AddOutput("task1", "test output")
	
	outputs := output.GetOutputs()
	if len(outputs) != 1 {
		t.Errorf("expected 1 output, got %d", len(outputs))
	}
	
	if outputs["task1"] != "test output" {
		t.Errorf("expected 'test output', got %v", outputs["task1"])
	}
	
	// Test adding map output
	mapOutput := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	output.AddOutput("task2", mapOutput)
	
	outputs = output.GetOutputs()
	if len(outputs) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(outputs))
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool // just check it doesn't panic
	}{
		{"string", "test", true},
		{"int", 42, true},
		{"int32", int32(42), true},
		{"int64", int64(42), true},
		{"float32", float32(3.14), true},
		{"float64", 3.14, true},
		{"bool_true", true, true},
		{"bool_false", false, true},
		{"nil", nil, true},
		{"slice", []string{"a", "b"}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValue(tt.value)
			if result == "" {
				t.Errorf("formatValue returned empty string for %v", tt.value)
			}
		})
	}
}

func TestTaskOperations(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	// Test TaskStart
	output.TaskStart("test-task", "test description")
	
	// Test TaskSuccess
	duration := 100 * time.Millisecond
	output.TaskSuccess("test-task", duration, "success output")
	
	outputs := output.GetOutputs()
	if outputs["test-task"] != "success output" {
		t.Errorf("expected 'success output', got %v", outputs["test-task"])
	}
	
	// Test TaskFailure
	err := errors.New("test error")
	output.TaskFailure("failed-task", duration, err)
}

func TestWorkflowOperations(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	// Test WorkflowStart
	output.WorkflowStart("test-workflow", "test workflow description")
	
	// Test WorkflowSuccess
	duration := 200 * time.Millisecond
	output.WorkflowSuccess("test-workflow", duration, 5)
	
	// Test WorkflowFailure
	err := errors.New("workflow error")
	output.WorkflowFailure("failed-workflow", duration, err)
}

func TestMessageOperations(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	// Test Info
	output.Info("info message")
	
	// Test Warning
	output.Warning("warning message")
	
	// Test Error
	output.Error("error message")
	
	// Test Debug
	output.Debug("debug message")
}

func TestOperationSpinner(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	// Test StartOperation
	output.StartOperation("test operation")
	
	// Test UpdateOperation
	output.UpdateOperation("updated operation")
	
	// Test StopOperation with success
	output.StopOperation(true, "operation succeeded")
	
	// Test StopOperation with failure
	output.StartOperation("failing operation")
	output.StopOperation(false, "operation failed")
}

func TestProgressBar(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	// Test StartProgress
	output.StartProgress(10, "test progress")
	
	// Test UpdateProgress
	output.UpdateProgress()
	output.UpdateProgress()
	
	// Test StopProgress
	output.StopProgress()
}
