package runner

import (
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	lua "github.com/yuin/gopher-lua"
)

func TestRunSingleTask_Success(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a simple task with a command function that returns success
	task := &types.Task{
		Name: "test-task",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LBool(true))
			L.Push(lua.LString("success"))
			L.Push(L.NewTable())
			return 3
		}),
		Params: make(map[string]string),
	}

	success, msg, output, duration, err := RunSingleTask(L, task)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !success {
		t.Error("Expected success to be true")
	}
	if msg != "success" {
		t.Errorf("Expected message 'success', got '%s'", msg)
	}
	if output == nil {
		t.Error("Expected output table to not be nil")
	}
	if duration == 0 {
		t.Error("Expected duration to be greater than 0")
	}
}

func TestRunSingleTask_Failure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a task that returns failure
	task := &types.Task{
		Name: "failing-task",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("task failed"))
			L.Push(L.NewTable())
			return 3
		}),
		Params: make(map[string]string),
	}

	success, msg, output, duration, err := RunSingleTask(L, task)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if success {
		t.Error("Expected success to be false")
	}
	if msg != "task failed" {
		t.Errorf("Expected message 'task failed', got '%s'", msg)
	}
	if output == nil {
		t.Error("Expected output table to not be nil")
	}
	if duration == 0 {
		t.Error("Expected duration to be greater than 0")
	}
}

func TestRunSingleTask_WithOutput(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a task that returns data in the output table
	task := &types.Task{
		Name: "task-with-output",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			output := L.NewTable()
			output.RawSetString("result", lua.LString("test-value"))
			output.RawSetString("count", lua.LNumber(42))
			
			L.Push(lua.LBool(true))
			L.Push(lua.LString("completed"))
			L.Push(output)
			return 3
		}),
		Params: make(map[string]string),
	}

	success, msg, output, _, err := RunSingleTask(L, task)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !success {
		t.Error("Expected success to be true")
	}
	if msg != "completed" {
		t.Errorf("Expected message 'completed', got '%s'", msg)
	}
	if output == nil {
		t.Fatal("Expected output table to not be nil")
	}

	// Verify output table contents
	result := output.RawGetString("result")
	if str, ok := result.(lua.LString); !ok || string(str) != "test-value" {
		t.Errorf("Expected output.result to be 'test-value', got '%v'", result)
	}

	count := output.RawGetString("count")
	if num, ok := count.(lua.LNumber); !ok || float64(num) != 42 {
		t.Errorf("Expected output.count to be 42, got '%v'", count)
	}
}

func TestRunSingleTask_NilCommandFunc(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	task := &types.Task{
		Name:        "nil-func-task",
		CommandFunc: nil,
		Params:      make(map[string]string),
	}

	success, msg, output, _, err := RunSingleTask(L, task)

	// Should handle nil function gracefully
	if err == nil {
		t.Error("Expected error for nil command function")
	}
	if success {
		t.Error("Expected success to be false with nil function")
	}
	_ = msg
	_ = output
}

func TestRunSingleTask_WithParams(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a task that uses params
	task := &types.Task{
		Name: "task-with-params",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LBool(true))
			L.Push(lua.LString("ok"))
			L.Push(L.NewTable())
			return 3
		}),
		Params: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	success, _, _, _, err := RunSingleTask(L, task)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !success {
		t.Error("Expected success to be true")
	}
}

func TestRunSingleTask_TimingAccuracy(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a task that completes quickly
	task := &types.Task{
		Name: "quick-task",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LBool(true))
			L.Push(lua.LString("done"))
			L.Push(L.NewTable())
			return 3
		}),
		Params: make(map[string]string),
	}

	_, _, _, duration, err := RunSingleTask(L, task)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if duration <= 0 {
		t.Error("Expected positive duration")
	}
	// Duration should be relatively quick (less than 1 second for this simple task)
	if duration.Seconds() > 1.0 {
		t.Logf("Warning: Task took longer than expected: %v", duration)
	}
}
