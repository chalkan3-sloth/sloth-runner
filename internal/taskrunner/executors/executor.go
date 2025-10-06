package executors

import (
	"context"
	"sync"

	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	lua "github.com/yuin/gopher-lua"
)

// TaskExecutor defines the interface for executing tasks
type TaskExecutor interface {
	Execute(
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

// ExecutionContext holds common execution context
type ExecutionContext struct {
	Ctx                   context.Context
	Task                  *types.Task
	L                     *lua.LState
	InputFromDependencies *lua.LTable
	Session               *types.SharedSession
	GroupName             string
	Mu                    *sync.Mutex
	CompletedTasks        map[string]bool
	TaskOutputs           map[string]*lua.LTable
	RunningTasks          map[string]bool
}

// NewExecutionContext creates a new execution context
func NewExecutionContext(
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
) *ExecutionContext {
	return &ExecutionContext{
		Ctx:                   ctx,
		Task:                  task,
		L:                     L,
		InputFromDependencies: inputFromDependencies,
		Session:               session,
		GroupName:             groupName,
		Mu:                    mu,
		CompletedTasks:        completedTasks,
		TaskOutputs:           taskOutputs,
		RunningTasks:          runningTasks,
	}
}
