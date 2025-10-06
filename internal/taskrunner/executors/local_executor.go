package executors

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	lua "github.com/yuin/gopher-lua"
)

// LocalExecutor executes tasks locally using Lua scripts
type LocalExecutor struct {
	TaskRunner            interface{} // Reference to parent TaskRunner
	ExecuteSuccessHandler func(L *lua.LState, t *types.Task, ctx context.Context, output *lua.LTable)
	ExecuteFailureHandler func(L *lua.LState, t *types.Task, ctx context.Context, errorMsg string)
}

// TaskExecutionError provides a more context-rich error for task failures
type TaskExecutionError struct {
	TaskName string
	Err      error
}

func (e *TaskExecutionError) Error() string {
	return fmt.Sprintf("task '%s' failed: %v", e.TaskName, e.Err)
}

// Execute runs the task locally using Lua
func (le *LocalExecutor) Execute(
	ctx context.Context,
	task *types.Task,
	mainL *lua.LState,
	inputFromDependencies *lua.LTable,
	session *types.SharedSession,
	groupName string,
	mu *sync.Mutex,
	completedTasks map[string]bool,
	taskOutputs map[string]*lua.LTable,
	runningTasks map[string]bool,
	results *[]types.TaskResult,
) (taskErr error) {
	startTime := time.Now()

	// Create new Lua state for this task
	L := lua.NewState()
	defer L.Close()
	luainterface.OpenAll(L)

	localInputFromDependencies := luainterface.CopyTable(inputFromDependencies, L)
	task.Output = L.NewTable()

	// Defer cleanup and result tracking
	defer func() {
		if r := recover(); r != nil {
			taskErr = &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("panic: %v", r)}
		}

		duration := time.Since(startTime)
		status := "Success"
		if taskErr != nil {
			status = "Failed"
		}

		mu.Lock()
		*results = append(*results, types.TaskResult{
			Name:     task.Name,
			Status:   status,
			Duration: duration,
			Error:    taskErr,
		})
		taskOutputs[task.Name] = luainterface.CopyTable(task.Output, mainL)
		completedTasks[task.Name] = true
		delete(runningTasks, task.Name)
		mu.Unlock()
	}()

	// Execute pre_exec hook
	if task.PreExec != nil {
		success, msg, _, err := luainterface.ExecuteLuaFunction(L, task.PreExec, task.Params, localInputFromDependencies, 2, ctx)
		if err != nil {
			return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("error executing pre_exec hook: %w", err)}
		} else if !success {
			return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("pre-execution hook failed: %s", msg)}
		}
	}

	// Execute command function
	if task.CommandFunc != nil {
		if task.Params == nil {
			task.Params = make(map[string]string)
		}
		task.Params["task_name"] = task.Name
		task.Params["group_name"] = groupName

		// Use task workdir if defined, otherwise use session workdir
		taskWorkdir := session.Workdir
		if task.Workdir != "" {
			taskWorkdir = task.Workdir
			// Create workdir if it doesn't exist
			if err := os.MkdirAll(taskWorkdir, 0755); err != nil {
				return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("failed to create workdir %s: %w", taskWorkdir, err)}
			}
		}
		task.Params["workdir"] = taskWorkdir

		var sessionUD *lua.LUserData
		if session != nil {
			sessionUD = L.NewUserData()
			sessionUD.Value = session
			L.SetMetatable(sessionUD, L.GetTypeMetatable("session"))
		}

		success, msg, outputTable, err := luainterface.ExecuteLuaFunction(L, task.CommandFunc, task.Params, localInputFromDependencies, 3, ctx, sessionUD)
		if err != nil {
			// Execute OnFailure handler if command function has error
			if task.OnFailure != nil && le.ExecuteFailureHandler != nil {
				le.ExecuteFailureHandler(L, task, ctx, fmt.Sprintf("error executing command function: %v", err))
			}
			return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("error executing command function: %w", err)}
		} else if !success {
			// Execute OnFailure handler if command function returns false
			if task.OnFailure != nil && le.ExecuteFailureHandler != nil {
				le.ExecuteFailureHandler(L, task, ctx, msg)
			}
			return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("command function returned failure: %s", msg)}
		} else if outputTable != nil {
			task.Output = outputTable
			// Execute OnSuccess handler if command was successful
			if task.OnSuccess != nil && le.ExecuteSuccessHandler != nil {
				le.ExecuteSuccessHandler(L, task, ctx, outputTable)
			}
		} else {
			// Execute OnSuccess handler even if no output table
			if task.OnSuccess != nil && le.ExecuteSuccessHandler != nil {
				le.ExecuteSuccessHandler(L, task, ctx, L.NewTable())
			}
		}
	}

	// Execute post_exec hook
	if task.PostExec != nil {
		var postExecSecondArg lua.LValue = task.Output
		if task.Output == nil {
			postExecSecondArg = L.NewTable()
		}
		success, msg, _, err := luainterface.ExecuteLuaFunction(L, task.PostExec, task.Params, postExecSecondArg, 2, ctx)
		if err != nil {
			return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("error executing post_exec hook: %w", err)}
		} else if !success {
			return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("post-execution hook failed: %s", msg)}
		}
	}

	return nil
}
