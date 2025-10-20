package executors

import (
	"context"
	"sync"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	lua "github.com/yuin/gopher-lua"
)

// Test TaskExecutor interface
func TestTaskExecutor_InterfaceExists(t *testing.T) {
	// Verify interface can be referenced
	var _ TaskExecutor

	// This test ensures the interface exists
}

type mockExecutor struct {
	executeFunc func(
		ctx context.Context,
		task *types.Task,
		L *lua.LState,
		inputFromDependencies *lua.LTable,
		session *types.SharedSession,
		groupName string,
		mu *sync.Mutex,
		completedTasks map[string]bool,
		taskOutputs map[string]*lua.LTable,
		runningTasks map[string]bool,
	) error
}

func (m *mockExecutor) Execute(
	ctx context.Context,
	task *types.Task,
	L *lua.LState,
	inputFromDependencies *lua.LTable,
	session *types.SharedSession,
	groupName string,
	mu *sync.Mutex,
	completedTasks map[string]bool,
	taskOutputs map[string]*lua.LTable,
	runningTasks map[string]bool,
) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, task, L, inputFromDependencies, session, groupName, mu, completedTasks, taskOutputs, runningTasks)
	}
	return nil
}

func TestTaskExecutor_MockImplementation(t *testing.T) {
	// MockExecutor should implement TaskExecutor interface
	var _ TaskExecutor = (*mockExecutor)(nil)

	mock := &mockExecutor{
		executeFunc: func(
			ctx context.Context,
			task *types.Task,
			L *lua.LState,
			inputFromDependencies *lua.LTable,
			session *types.SharedSession,
			groupName string,
			mu *sync.Mutex,
			completedTasks map[string]bool,
			taskOutputs map[string]*lua.LTable,
			runningTasks map[string]bool,
		) error {
			return nil
		},
	}

	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	err := mock.Execute(ctx, &types.Task{}, L, L.NewTable(), &types.SharedSession{}, "test", &sync.Mutex{}, make(map[string]bool), make(map[string]*lua.LTable), make(map[string]bool))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test ExecutionContext struct
func TestExecutionContext_Creation(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	task := &types.Task{Name: "test-task"}
	inputTable := L.NewTable()
	session := &types.SharedSession{}
	mu := &sync.Mutex{}
	completedTasks := make(map[string]bool)
	taskOutputs := make(map[string]*lua.LTable)
	runningTasks := make(map[string]bool)

	execCtx := &ExecutionContext{
		Ctx:                   ctx,
		Task:                  task,
		L:                     L,
		InputFromDependencies: inputTable,
		Session:               session,
		GroupName:             "test-group",
		Mu:                    mu,
		CompletedTasks:        completedTasks,
		TaskOutputs:           taskOutputs,
		RunningTasks:          runningTasks,
	}

	if execCtx.Task.Name != "test-task" {
		t.Error("Expected Task to be set")
	}

	if execCtx.GroupName != "test-group" {
		t.Error("Expected GroupName to be set")
	}
}

func TestExecutionContext_ZeroValue(t *testing.T) {
	var execCtx ExecutionContext

	if execCtx.Ctx != nil {
		t.Error("Expected zero Ctx")
	}

	if execCtx.Task != nil {
		t.Error("Expected zero Task")
	}

	if execCtx.GroupName != "" {
		t.Error("Expected zero GroupName")
	}
}

func TestExecutionContext_AllFieldsSet(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	task := &types.Task{Name: "task1"}
	inputTable := L.NewTable()
	session := &types.SharedSession{}
	mu := &sync.Mutex{}
	completedTasks := map[string]bool{"task0": true}
	taskOutputs := make(map[string]*lua.LTable)
	runningTasks := map[string]bool{"task1": true}

	execCtx := &ExecutionContext{
		Ctx:                   ctx,
		Task:                  task,
		L:                     L,
		InputFromDependencies: inputTable,
		Session:               session,
		GroupName:             "group1",
		Mu:                    mu,
		CompletedTasks:        completedTasks,
		TaskOutputs:           taskOutputs,
		RunningTasks:          runningTasks,
	}

	if execCtx.Ctx == nil {
		t.Error("Expected Ctx to be set")
	}

	if execCtx.Task == nil {
		t.Error("Expected Task to be set")
	}

	if execCtx.L == nil {
		t.Error("Expected L to be set")
	}

	if execCtx.InputFromDependencies == nil {
		t.Error("Expected InputFromDependencies to be set")
	}

	if execCtx.Session == nil {
		t.Error("Expected Session to be set")
	}

	if execCtx.Mu == nil {
		t.Error("Expected Mu to be set")
	}

	if execCtx.CompletedTasks == nil {
		t.Error("Expected CompletedTasks to be set")
	}

	if execCtx.TaskOutputs == nil {
		t.Error("Expected TaskOutputs to be set")
	}

	if execCtx.RunningTasks == nil {
		t.Error("Expected RunningTasks to be set")
	}
}

// Test NewExecutionContext function
func TestNewExecutionContext(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	task := &types.Task{Name: "new-task"}
	inputTable := L.NewTable()
	session := &types.SharedSession{}
	mu := &sync.Mutex{}
	completedTasks := make(map[string]bool)
	taskOutputs := make(map[string]*lua.LTable)
	runningTasks := make(map[string]bool)

	execCtx := NewExecutionContext(
		ctx,
		task,
		L,
		inputTable,
		session,
		"new-group",
		mu,
		completedTasks,
		taskOutputs,
		runningTasks,
	)

	if execCtx == nil {
		t.Error("Expected non-nil ExecutionContext")
	}

	if execCtx.Task.Name != "new-task" {
		t.Errorf("Expected task name 'new-task', got '%s'", execCtx.Task.Name)
	}

	if execCtx.GroupName != "new-group" {
		t.Errorf("Expected group name 'new-group', got '%s'", execCtx.GroupName)
	}
}

func TestNewExecutionContext_MultipleInstances(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	task1 := &types.Task{Name: "task1"}
	task2 := &types.Task{Name: "task2"}

	execCtx1 := NewExecutionContext(
		ctx, task1, L, L.NewTable(), &types.SharedSession{}, "group1",
		&sync.Mutex{}, make(map[string]bool), make(map[string]*lua.LTable), make(map[string]bool),
	)

	execCtx2 := NewExecutionContext(
		ctx, task2, L, L.NewTable(), &types.SharedSession{}, "group2",
		&sync.Mutex{}, make(map[string]bool), make(map[string]*lua.LTable), make(map[string]bool),
	)

	if execCtx1 == execCtx2 {
		t.Error("Expected different instances")
	}

	if execCtx1.Task.Name == execCtx2.Task.Name {
		t.Error("Expected different task names")
	}

	if execCtx1.GroupName == execCtx2.GroupName {
		t.Error("Expected different group names")
	}
}

func TestNewExecutionContext_NilValues(t *testing.T) {
	// Test creating ExecutionContext with nil values (should not panic)
	execCtx := NewExecutionContext(
		nil,       // ctx
		nil,       // task
		nil,       // L
		nil,       // inputTable
		nil,       // session
		"",        // groupName
		nil,       // mu
		nil,       // completedTasks
		nil,       // taskOutputs
		nil,       // runningTasks
	)

	if execCtx == nil {
		t.Error("Expected non-nil ExecutionContext even with nil parameters")
	}

	if execCtx.Ctx != nil {
		t.Error("Expected Ctx to be nil")
	}

	if execCtx.Task != nil {
		t.Error("Expected Task to be nil")
	}
}

// Test ExecutionContext with context cancellation
func TestExecutionContext_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	L := lua.NewState()
	defer L.Close()

	execCtx := NewExecutionContext(
		ctx,
		&types.Task{Name: "cancellable"},
		L,
		L.NewTable(),
		&types.SharedSession{},
		"group",
		&sync.Mutex{},
		make(map[string]bool),
		make(map[string]*lua.LTable),
		make(map[string]bool),
	)

	// Context should not be cancelled initially
	select {
	case <-execCtx.Ctx.Done():
		t.Error("Context should not be cancelled initially")
	default:
		// Good
	}

	// Cancel context
	cancel()

	// Context should now be cancelled
	select {
	case <-execCtx.Ctx.Done():
		// Good - context is cancelled
	default:
		t.Error("Context should be cancelled")
	}
}

// Test ExecutionContext with mutex
func TestExecutionContext_MutexLocking(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	mu := &sync.Mutex{}
	execCtx := NewExecutionContext(
		ctx,
		&types.Task{},
		L,
		L.NewTable(),
		&types.SharedSession{},
		"group",
		mu,
		make(map[string]bool),
		make(map[string]*lua.LTable),
		make(map[string]bool),
	)

	// Lock and unlock should work
	execCtx.Mu.Lock()
	execCtx.Mu.Unlock()

	// Should be same mutex instance
	if execCtx.Mu != mu {
		t.Error("Expected same mutex instance")
	}
}

// Test ExecutionContext maps manipulation
func TestExecutionContext_CompletedTasksMap(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	completedTasks := make(map[string]bool)
	execCtx := NewExecutionContext(
		ctx,
		&types.Task{},
		L,
		L.NewTable(),
		&types.SharedSession{},
		"group",
		&sync.Mutex{},
		completedTasks,
		make(map[string]*lua.LTable),
		make(map[string]bool),
	)

	// Initially empty
	if len(execCtx.CompletedTasks) != 0 {
		t.Error("Expected empty CompletedTasks map")
	}

	// Add task
	execCtx.CompletedTasks["task1"] = true

	// Should be reflected in original map
	if !completedTasks["task1"] {
		t.Error("Expected task1 to be marked completed in original map")
	}
}

func TestExecutionContext_TaskOutputsMap(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	taskOutputs := make(map[string]*lua.LTable)
	execCtx := NewExecutionContext(
		ctx,
		&types.Task{},
		L,
		L.NewTable(),
		&types.SharedSession{},
		"group",
		&sync.Mutex{},
		make(map[string]bool),
		taskOutputs,
		make(map[string]bool),
	)

	// Initially empty
	if len(execCtx.TaskOutputs) != 0 {
		t.Error("Expected empty TaskOutputs map")
	}

	// Add output
	outputTable := L.NewTable()
	execCtx.TaskOutputs["task1"] = outputTable

	// Should be reflected in original map
	if taskOutputs["task1"] != outputTable {
		t.Error("Expected task1 output to be in original map")
	}
}

func TestExecutionContext_RunningTasksMap(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	runningTasks := make(map[string]bool)
	execCtx := NewExecutionContext(
		ctx,
		&types.Task{},
		L,
		L.NewTable(),
		&types.SharedSession{},
		"group",
		&sync.Mutex{},
		make(map[string]bool),
		make(map[string]*lua.LTable),
		runningTasks,
	)

	// Initially empty
	if len(execCtx.RunningTasks) != 0 {
		t.Error("Expected empty RunningTasks map")
	}

	// Add running task
	execCtx.RunningTasks["task1"] = true

	// Should be reflected in original map
	if !runningTasks["task1"] {
		t.Error("Expected task1 to be marked running in original map")
	}
}

// Test ExecutionContext with SharedSession
func TestExecutionContext_SharedSession(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	session := &types.SharedSession{}
	execCtx := NewExecutionContext(
		ctx,
		&types.Task{},
		L,
		L.NewTable(),
		session,
		"group",
		&sync.Mutex{},
		make(map[string]bool),
		make(map[string]*lua.LTable),
		make(map[string]bool),
	)

	if execCtx.Session != session {
		t.Error("Expected same session instance")
	}
}

// Test ExecutionContext with different group names
func TestExecutionContext_GroupNames(t *testing.T) {
	groupNames := []string{
		"default",
		"production",
		"staging",
		"development",
		"group-with-dashes",
		"group_with_underscores",
	}

	for _, groupName := range groupNames {
		ctx := context.Background()
		L := lua.NewState()
		defer L.Close()

		execCtx := NewExecutionContext(
			ctx,
			&types.Task{},
			L,
			L.NewTable(),
			&types.SharedSession{},
			groupName,
			&sync.Mutex{},
			make(map[string]bool),
			make(map[string]*lua.LTable),
			make(map[string]bool),
		)

		if execCtx.GroupName != groupName {
			t.Errorf("Expected group name '%s', got '%s'", groupName, execCtx.GroupName)
		}
	}
}

// Test ExecutionContext with different tasks
func TestExecutionContext_DifferentTasks(t *testing.T) {
	taskNames := []string{"task1", "task2", "task3", "build", "deploy", "test"}

	for _, taskName := range taskNames {
		ctx := context.Background()
		L := lua.NewState()
		defer L.Close()

		task := &types.Task{Name: taskName}
		execCtx := NewExecutionContext(
			ctx,
			task,
			L,
			L.NewTable(),
			&types.SharedSession{},
			"group",
			&sync.Mutex{},
			make(map[string]bool),
			make(map[string]*lua.LTable),
			make(map[string]bool),
		)

		if execCtx.Task.Name != taskName {
			t.Errorf("Expected task name '%s', got '%s'", taskName, execCtx.Task.Name)
		}
	}
}

// Test ExecutionContext concurrent access
func TestExecutionContext_ConcurrentMapAccess(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	execCtx := NewExecutionContext(
		ctx,
		&types.Task{},
		L,
		L.NewTable(),
		&types.SharedSession{},
		"group",
		&sync.Mutex{},
		make(map[string]bool),
		make(map[string]*lua.LTable),
		make(map[string]bool),
	)

	done := make(chan bool, 10)

	// Multiple goroutines modifying maps with mutex protection
	for i := 0; i < 10; i++ {
		go func(index int) {
			execCtx.Mu.Lock()
			taskName := string(rune('a' + index))
			execCtx.CompletedTasks[taskName] = true
			execCtx.RunningTasks[taskName] = false
			execCtx.Mu.Unlock()
			done <- true
		}(i)
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all tasks were added
	if len(execCtx.CompletedTasks) != 10 {
		t.Errorf("Expected 10 completed tasks, got %d", len(execCtx.CompletedTasks))
	}
}

// Test ExecutionContext with input dependencies
func TestExecutionContext_InputDependencies(t *testing.T) {
	ctx := context.Background()
	L := lua.NewState()
	defer L.Close()

	inputTable := L.NewTable()
	L.SetField(inputTable, "key1", lua.LString("value1"))
	L.SetField(inputTable, "key2", lua.LNumber(42))

	execCtx := NewExecutionContext(
		ctx,
		&types.Task{},
		L,
		inputTable,
		&types.SharedSession{},
		"group",
		&sync.Mutex{},
		make(map[string]bool),
		make(map[string]*lua.LTable),
		make(map[string]bool),
	)

	// Verify input table is accessible
	value1 := L.GetField(execCtx.InputFromDependencies, "key1")
	if value1.String() != "value1" {
		t.Errorf("Expected 'value1', got '%s'", value1.String())
	}

	value2 := L.GetField(execCtx.InputFromDependencies, "key2")
	if float64(value2.(lua.LNumber)) != 42 {
		t.Errorf("Expected 42, got %v", value2)
	}
}

// Test ExecutionContext initialization patterns
func TestExecutionContext_InitializationPatterns(t *testing.T) {
	testCases := []struct {
		name      string
		groupName string
		taskName  string
	}{
		{"Simple", "group1", "task1"},
		{"With dashes", "prod-group", "deploy-task"},
		{"With underscores", "dev_group", "build_task"},
		{"Mixed", "staging-env_1", "test_suite-1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			L := lua.NewState()
			defer L.Close()

			task := &types.Task{Name: tc.taskName}
			execCtx := NewExecutionContext(
				ctx,
				task,
				L,
				L.NewTable(),
				&types.SharedSession{},
				tc.groupName,
				&sync.Mutex{},
				make(map[string]bool),
				make(map[string]*lua.LTable),
				make(map[string]bool),
			)

			if execCtx.GroupName != tc.groupName {
				t.Errorf("Expected group '%s', got '%s'", tc.groupName, execCtx.GroupName)
			}

			if execCtx.Task.Name != tc.taskName {
				t.Errorf("Expected task '%s', got '%s'", tc.taskName, execCtx.Task.Name)
			}
		})
	}
}
