package types

import (
	"io"
	"os/exec"
	"time"

	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

// Task represents a single unit of work in the runner.
type Task struct {
	ID          string // Unique identifier for the task
	Name        string
	Description string
	Workdir     string // ✅ Added individual task workdir support
	User        string // ✅ User to run the task as (default: root)
	CommandFunc *lua.LFunction
	CommandStr  string
	DependsOn   []string
	Artifacts   []string
	Consumes    []string
	NextIfFail  []string
	Params      map[string]string
	Retries     int
	Timeout     string
	Async       bool
	PreExec     *lua.LFunction
	PostExec    *lua.LFunction
	OnSuccess   *lua.LFunction // ✅ Added success handler
	OnFailure   *lua.LFunction // ✅ Added failure handler
	RunIf       string
	RunIfFunc   *lua.LFunction
	AbortIf     string
	AbortIfFunc *lua.LFunction
	Output      *lua.LTable
	DelegateTo  interface{} // Can be string (agent name) or map (inline agent definition)
}

// TaskGroup represents a collection of related tasks.
type TaskGroup struct {
	ID                       string // Unique identifier for the task group
	Description              string
	Tasks                    []Task
	Workdir                  string
	CreateWorkdirBeforeRun   bool
	CleanWorkdirAfterRunFunc *lua.LFunction
	DelegateTo               interface{} `yaml:"delegate_to"` // Can be map[string]Agent or string (default agent)
}

// GenerateTaskID generates a new UUID for a task
func GenerateTaskID() string {
	return uuid.New().String()
}

// GenerateTaskGroupID generates a new UUID for a task group
func GenerateTaskGroupID() string {
	return uuid.New().String()
}

// TaskResult holds the outcome of a single task execution.
type TaskResult struct {
	Name     string
	Status   string
	Duration time.Duration
	Error    error
}

// SharedSession holds data that can be shared between tasks in a group.
type SharedSession struct {
	Workdir string
	Cmd     *exec.Cmd
	Stdin   io.WriteCloser
	Stdout  io.ReadCloser
	Stderr  io.ReadCloser
}

// TaskRunner is the interface for the main task execution engine.
type TaskRunner interface {
	Run() error
	RunTasksParallel(tasks []*Task, input *lua.LTable) ([]TaskResult, error)
}

// Exporter defines an interface for exporting data from a Lua script.
type Exporter interface {
	Export(data map[string]interface{})
}

// PythonVenv represents a Python virtual environment.
type PythonVenv struct {
	Path string
}
