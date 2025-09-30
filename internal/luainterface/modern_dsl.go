package luainterface

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/core"
	lua "github.com/yuin/gopher-lua"
)

// ModernDSL provides a fluent, modern syntax for task definition
type ModernDSL struct {
core         *core.GlobalCore
logger       *slog.Logger
taskRegistry *TaskRegistry
builders     map[string]*TaskBuilder
mu           sync.RWMutex
}

// TaskRegistry maintains task definitions and metadata
type TaskRegistry struct {
tasks       map[string]*TaskDefinition
groups      map[string]*TaskGroup
templates   map[string]*TaskTemplate
validators  map[string]TaskValidator
mu          sync.RWMutex
}

// WorkflowBuilder provides fluent API for workflow construction
type WorkflowBuilder struct {
	name        string
	description string
	version     string
	tasks       []*TaskDefinition
	config      map[string]interface{}
	metadata    map[string]interface{}
	onComplete  *lua.LFunction
	onStart     *lua.LFunction
}

// TaskBuilder provides fluent API for task construction
type TaskBuilder struct {
definition *TaskDefinition
context    *BuildContext
chain      []BuildStep
}

// TaskDefinition represents a complete task definition
type TaskDefinition struct {
Name            string                 `json:"name"`
Description     string                 `json:"description"`
Version         string                 `json:"version"`
Tags            []string               `json:"tags"`
Category        string                 `json:"category"`
Workdir         string                 `json:"workdir"` // ✅ Added workdir field

// Execution properties
Command         interface{}            `json:"command"`
Dependencies    []Dependency           `json:"dependencies"`
Parameters      map[string]Parameter   `json:"parameters"`
Environment     map[string]string      `json:"environment"`

// Behavior
Async           bool                   `json:"async"`
Timeout         time.Duration          `json:"timeout"`
Retries         RetryConfig            `json:"retries"`

// Hooks and lifecycle
PreHooks        []Hook                 `json:"pre_hooks"`
PostHooks       []Hook                 `json:"post_hooks"`
OnSuccess       []Hook                 `json:"on_success"`
OnFailure       []Hook                 `json:"on_failure"`
Cleanup         []Hook                 `json:"cleanup"`

// Advanced features
Conditions      []Condition            `json:"conditions"`
Outputs         []Output               `json:"outputs"`
Artifacts       []Artifact             `json:"artifacts"`
Resources       ResourceRequirements   `json:"resources"`
Security        SecurityPolicy         `json:"security"`

// Orchestration
Delegation      DelegationConfig       `json:"delegation"`
Saga            SagaConfig            `json:"saga"`
Circuit         CircuitConfig         `json:"circuit"`

// Metadata
CreatedAt       time.Time             `json:"created_at"`
UpdatedAt       time.Time             `json:"updated_at"`
Author          string                `json:"author"`
Metadata        map[string]interface{} `json:"metadata"`
}

// Supporting types for modern DSL
type Dependency struct {
Name      string            `json:"name"`
Type      DependencyType    `json:"type"`
Optional  bool             `json:"optional"`
Condition string           `json:"condition"`
Timeout   time.Duration    `json:"timeout"`
}

type Parameter struct {
Type        ParameterType    `json:"type"`
Required    bool            `json:"required"`
Default     interface{}     `json:"default"`
Validation  []Validator     `json:"validation"`
Description string          `json:"description"`
}

type Hook struct {
Name      string           `json:"name"`
Type      HookType         `json:"type"`
Command   interface{}      `json:"command"`
Condition string          `json:"condition"`
Async     bool            `json:"async"`
}

type Condition struct {
Name       string          `json:"name"`
Expression string          `json:"expression"`
Type       ConditionType   `json:"type"`
}

type Output struct {
Name        string         `json:"name"`
Type        OutputType     `json:"type"`
Path        string         `json:"path"`
Transform   string         `json:"transform"`
Persistent  bool           `json:"persistent"`
}

type Artifact struct {
Name        string         `json:"name"`
Path        string         `json:"path"`
Type        ArtifactType   `json:"type"`
Compression bool           `json:"compression"`
Retention   time.Duration  `json:"retention"`
}

// Enums
type DependencyType int
type ParameterType int
type HookType int
type ConditionType int
type OutputType int
type ArtifactType int

const (
DependencyTask DependencyType = iota
DependencyResource
DependencyService
DependencyData
)

const (
ParamString ParameterType = iota
ParamInt
ParamFloat
ParamBool
ParamArray
ParamObject
)

const (
HookShell HookType = iota
HookLua
HookHTTP
HookGRPC
)

const (
ConditionShell ConditionType = iota
ConditionLua
ConditionExpression
)

const (
OutputString OutputType = iota
OutputJSON
OutputFile
OutputArtifact
)

const (
ArtifactFile ArtifactType = iota
ArtifactDirectory
ArtifactArchive
ArtifactContainer
)

// NewModernDSL creates a new modern DSL instance
func NewModernDSL(globalCore *core.GlobalCore) *ModernDSL {
return &ModernDSL{
core:         globalCore,
logger:       slog.Default(),
taskRegistry: NewTaskRegistry(),
builders:     make(map[string]*TaskBuilder),
}
}

// RegisterModernDSL registers the modern DSL with a Lua state
func (m *ModernDSL) RegisterModernDSL(L *lua.LState) {
	// Setup metatables first
	m.setupMetatables(L)
	
	// Core DSL functions
	m.registerTaskDefinition(L)
	m.registerWorkflowDefinition(L)
	m.registerBuilders(L)
	m.registerUtilities(L)
	m.registerValidators(L)
	m.registerTemplates(L)

	// Advanced features
	m.registerSagaSupport(L)
	m.registerCircuitBreaker(L)
	m.registerResourceManagement(L)
	m.registerSecurityPolicies(L)
}

// setupMetatables creates the metatables for DSL objects
func (m *ModernDSL) setupMetatables(L *lua.LState) {
	// TaskBuilder metatable
	taskBuilderMt := L.NewTypeMetatable("TaskBuilder")
	L.SetField(taskBuilderMt, "__index", L.NewFunction(m.taskBuilderIndex))
	
	// WorkflowBuilder metatable
	workflowBuilderMt := L.NewTypeMetatable("WorkflowBuilder")
	L.SetField(workflowBuilderMt, "__index", L.NewFunction(m.workflowBuilderIndex))
}

// taskBuilderIndex handles method calls on TaskBuilder objects
func (m *ModernDSL) taskBuilderIndex(L *lua.LState) int {
	ud := L.CheckUserData(1)
	key := L.CheckString(2)
	
	builder, ok := ud.Value.(*TaskBuilder)
	if !ok {
		L.ArgError(1, "TaskBuilder expected")
		return 0
	}
	
	switch key {
	case "description":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			desc := L.CheckString(2) // Argument position 2 (1 is self)
			builder.definition.Description = desc
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "command":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			cmd := L.CheckAny(2) // Argument position 2 (1 is self)
			
			// ✅ Store the original function and builder reference
			if cmdFunc, ok := cmd.(*lua.LFunction); ok {
				// Store original function and task definition reference in builder
				builder.definition.Command = cmdFunc
				// We'll handle the 'this' injection during execution
			} else {
				builder.definition.Command = cmd
			}
			
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "timeout":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			timeoutStr := L.CheckString(2) // Argument position 2 (1 is self)
			if duration, err := time.ParseDuration(timeoutStr); err == nil {
				builder.definition.Timeout = duration
			}
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "workdir":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			workdirPath := L.CheckString(2) // Argument position 2 (1 is self)
			builder.definition.Workdir = workdirPath
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "depends_on":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // deps variable - simplified for now
			// Handle dependencies - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "on_success":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			hook := L.CheckFunction(2) // hook function
			
			// Store the success hook in the task definition
			builder.definition.OnSuccess = append(builder.definition.OnSuccess, Hook{
				Name:    "on_success",
				Type:    HookLua,
				Command: hook,
			})
			
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "on_failure", "on_fail":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			hook := L.CheckFunction(2) // hook function
			
			// Store the failure hook in the task definition
			builder.definition.OnFailure = append(builder.definition.OnFailure, Hook{
				Name:    "on_failure",
				Type:    HookLua,
				Command: hook,
			})
			
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "pre_hook":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // hook function - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "post_hook":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // hook function - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "retries":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // retries config - simplified for now
			_ = L.CheckAny(3) // strategy - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "async":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckBool(2) // async flag - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "artifacts":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // artifacts - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "on_timeout":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // timeout handler - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "run_if":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // condition - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "abort_if":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			_ = L.CheckAny(2) // abort condition - simplified for now
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "build":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			// Return the built task definition
			taskUd := L.NewUserData()
			taskUd.Value = builder.definition
			L.Push(taskUd)
			return 1
		}))
	default:
		L.Push(lua.LNil)
	}
	return 1
}

// workflowBuilderIndex handles method calls on WorkflowBuilder objects
func (m *ModernDSL) workflowBuilderIndex(L *lua.LState) int {
	ud := L.CheckUserData(1)
	key := L.CheckString(2)
	
	builder, ok := ud.Value.(*WorkflowBuilder)
	if !ok {
		L.ArgError(1, "WorkflowBuilder expected")
		return 0
	}
	
	switch key {
	case "description":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			desc := L.CheckString(2) // Argument position 2 (1 is self)
			builder.description = desc
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "version":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			version := L.CheckString(2) // Argument position 2 (1 is self)
			builder.version = version
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "tasks":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			tasksArg := L.CheckTable(2) // Argument position 2 (1 is self)
			
			// Convert Lua table to task definitions
			builder.tasks = []*TaskDefinition{}
			tasksArg.ForEach(func(_, taskValue lua.LValue) {
				if taskValue.Type() == lua.LTUserData {
					taskUD := taskValue.(*lua.LUserData)
					if taskDef, ok := taskUD.Value.(*TaskDefinition); ok {
						builder.tasks = append(builder.tasks, taskDef)
					}
				}
			})
			
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "config":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			configArg := L.CheckTable(2) // Argument position 2 (1 is self)
			
			// Convert Lua table to config map
			builder.config = make(map[string]interface{})
			configArg.ForEach(func(key, value lua.LValue) {
				if key.Type() == lua.LTString {
					keyStr := key.String()
					switch value.Type() {
					case lua.LTString:
						builder.config[keyStr] = value.String()
					case lua.LTNumber:
						builder.config[keyStr] = float64(value.(lua.LNumber))
					case lua.LTBool:
						builder.config[keyStr] = bool(value.(lua.LBool))
					}
				}
			})
			
			L.Push(ud) // Return self for chaining
			return 1
		}))
	case "on_complete":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			onCompleteFunc := L.CheckFunction(2) // Argument position 2 (1 is self)
			builder.onComplete = onCompleteFunc
			
			// Build and register immediately when on_complete is called
			// This supports the syntax with the trailing }))
			m.buildAndRegisterWorkflow(L, builder)
			
			return 0 // Don't return anything to end the chain
		}))
	case "on_start":
		L.Push(L.NewFunction(func(L *lua.LState) int {
			onStartFunc := L.CheckFunction(2) // Argument position 2 (1 is self)
			builder.onStart = onStartFunc
			L.Push(ud) // Return self for chaining
			return 1
		}))
	default:
		L.Push(lua.LNil)
	}
	return 1
}

// buildAndRegisterWorkflow finalizes a workflow definition and registers it
func (m *ModernDSL) buildAndRegisterWorkflow(L *lua.LState, builder *WorkflowBuilder) {
	// Get existing TaskDefinitions table or create it
	taskDefs := L.GetGlobal("TaskDefinitions")
	if taskDefs.Type() != lua.LTTable {
		taskDefs = L.NewTable()
		L.SetGlobal("TaskDefinitions", taskDefs)
	}
	
	// Create group structure compatible with legacy parser
	groupTable := L.NewTable()
	
	// Set description
	if builder.description != "" {
		groupTable.RawSetString("description", lua.LString(builder.description))
	}
	
	// Set version
	if builder.version != "" {
		groupTable.RawSetString("version", lua.LString(builder.version))
	}
	
	// Convert Modern DSL tasks to legacy format
	if len(builder.tasks) > 0 {
		tasksTable := L.NewTable()
		
		for i, taskDef := range builder.tasks {
			// Create legacy task structure
			legacyTask := L.NewTable()
			legacyTask.RawSetString("name", lua.LString(taskDef.Name))
			legacyTask.RawSetString("description", lua.LString(taskDef.Description))
			
			// Set workdir if specified
			if taskDef.Workdir != "" {
				legacyTask.RawSetString("workdir", lua.LString(taskDef.Workdir))
			}
			
			// Convert command
			if taskDef.Command != nil {
				if luaValue, ok := taskDef.Command.(lua.LValue); ok {
					legacyTask.RawSetString("command", luaValue)
				}
			}
			
			// Convert timeout
			if taskDef.Timeout > 0 {
				legacyTask.RawSetString("timeout", lua.LString(taskDef.Timeout.String()))
			}
			
			// Convert hooks
			if len(taskDef.OnSuccess) > 0 {
				if hook := taskDef.OnSuccess[0]; hook.Command != nil {
					if luaValue, ok := hook.Command.(lua.LValue); ok {
						legacyTask.RawSetString("on_success", luaValue)
					}
				}
			}
			
			if len(taskDef.OnFailure) > 0 {
				if hook := taskDef.OnFailure[0]; hook.Command != nil {
					if luaValue, ok := hook.Command.(lua.LValue); ok {
						legacyTask.RawSetString("on_failure", luaValue)
					}
				}
			}
			
			tasksTable.RawSetInt(i+1, legacyTask)
		}
		
		groupTable.RawSetString("tasks", tasksTable)
	}
	
	// Set config
	if len(builder.config) > 0 {
		configTable := L.NewTable()
		for key, value := range builder.config {
			switch v := value.(type) {
			case string:
				configTable.RawSetString(key, lua.LString(v))
			case float64:
				configTable.RawSetString(key, lua.LNumber(v))
			case bool:
				configTable.RawSetString(key, lua.LBool(v))
			}
		}
		groupTable.RawSetString("config", configTable)
	}
	
	// Set on_complete handler
	if builder.onComplete != nil {
		groupTable.RawSetString("on_complete", builder.onComplete)
	}
	
	// Set on_start handler
	if builder.onStart != nil {
		groupTable.RawSetString("on_start", builder.onStart)
	}
	
	// Add to TaskDefinitions
	taskDefs.(*lua.LTable).RawSetString(builder.name, groupTable)
}

func (m *ModernDSL) registerTaskDefinition(L *lua.LState) {
	// task() - main task builder function
	L.SetGlobal("task", L.NewFunction(m.taskBuilderFunc))
	
	// chain() - sequential task chain
	L.SetGlobal("chain", L.NewFunction(m.chainBuilderFunc))
	
	// parallel() - parallel task execution
	L.SetGlobal("parallel", L.NewFunction(m.parallelBuilderFunc))
	
	// when() - conditional execution
	L.SetGlobal("when", L.NewFunction(m.conditionalBuilderFunc))
}

func (m *ModernDSL) registerWorkflowDefinition(L *lua.LState) {
	// Create workflow namespace
	workflowMt := L.NewTable()
	L.SetField(workflowMt, "define", L.NewFunction(m.workflowDefineFunc))
	L.SetField(workflowMt, "parallel", L.NewFunction(m.workflowParallelFunc))
	L.SetField(workflowMt, "sequence", L.NewFunction(m.workflowSequenceFunc))
	L.SetField(workflowMt, "conditional", L.NewFunction(m.workflowConditionalFunc))
	L.SetGlobal("workflow", workflowMt)
}

func (m *ModernDSL) registerBuilders(L *lua.LState) {
	// async namespace
	asyncMt := L.NewTable()
	L.SetField(asyncMt, "parallel", L.NewFunction(m.asyncParallelFunc))
	L.SetField(asyncMt, "sequence", L.NewFunction(m.asyncSequenceFunc))
	L.SetField(asyncMt, "timeout", L.NewFunction(m.asyncTimeoutFunc))
	L.SetGlobal("async", asyncMt)
	
	// perf namespace for performance monitoring
	perfMt := L.NewTable()
	L.SetField(perfMt, "measure", L.NewFunction(m.perfMeasureFunc))
	L.SetField(perfMt, "stats", L.NewFunction(m.perfStatsFunc))
	L.SetGlobal("perf", perfMt)
	
	// core namespace for core system access
	coreMt := L.NewTable()
	L.SetField(coreMt, "stats", L.NewFunction(m.coreStatsFunc))
	L.SetField(coreMt, "resources", L.NewFunction(m.coreResourcesFunc))
	L.SetGlobal("core", coreMt)
}

func (m *ModernDSL) registerUtilities(L *lua.LState) {
	// utils namespace
	utilsMt := L.NewTable()
	L.SetField(utilsMt, "config", L.NewFunction(m.utilsConfigFunc))
	L.SetField(utilsMt, "secret", L.NewFunction(m.utilsSecretFunc))
	L.SetField(utilsMt, "env", L.NewFunction(m.utilsEnvFunc))
	L.SetGlobal("utils", utilsMt)
	
	// workdir namespace for workdir management
	workdirMt := L.NewTable()
	L.SetField(workdirMt, "get", L.NewFunction(m.workdirGetFunc))
	L.SetField(workdirMt, "cleanup", L.NewFunction(m.workdirCleanupFunc))
	L.SetField(workdirMt, "exists", L.NewFunction(m.workdirExistsFunc))
	L.SetField(workdirMt, "create", L.NewFunction(m.workdirCreateFunc))
	L.SetGlobal("workdir", workdirMt)
}

func (m *ModernDSL) registerValidators(L *lua.LState) {
	// validate namespace
	validateMt := L.NewTable()
	L.SetField(validateMt, "required", L.NewFunction(m.validateRequiredFunc))
	L.SetField(validateMt, "type", L.NewFunction(m.validateTypeFunc))
	L.SetField(validateMt, "range", L.NewFunction(m.validateRangeFunc))
	L.SetGlobal("validate", validateMt)
}

func (m *ModernDSL) registerTemplates(L *lua.LState) {
	// template namespace
	templateMt := L.NewTable()
	L.SetField(templateMt, "render", L.NewFunction(m.templateRenderFunc))
	L.SetField(templateMt, "load", L.NewFunction(m.templateLoadFunc))
	L.SetGlobal("template", templateMt)
}

func (m *ModernDSL) registerSagaSupport(L *lua.LState) {
	// saga namespace for transaction management
	sagaMt := L.NewTable()
	L.SetField(sagaMt, "begin", L.NewFunction(m.sagaBeginFunc))
	L.SetField(sagaMt, "compensate", L.NewFunction(m.sagaCompensateFunc))
	L.SetField(sagaMt, "commit", L.NewFunction(m.sagaCommitFunc))
	L.SetGlobal("saga", sagaMt)
}

func (m *ModernDSL) registerCircuitBreaker(L *lua.LState) {
	// circuit namespace
	circuitMt := L.NewTable()
	L.SetField(circuitMt, "protect", L.NewFunction(m.circuitProtectFunc))
	L.SetField(circuitMt, "status", L.NewFunction(m.circuitStatusFunc))
	L.SetGlobal("circuit", circuitMt)
}

func (m *ModernDSL) registerResourceManagement(L *lua.LState) {
	// resource namespace
	resourceMt := L.NewTable()
	L.SetField(resourceMt, "allocate", L.NewFunction(m.resourceAllocateFunc))
	L.SetField(resourceMt, "release", L.NewFunction(m.resourceReleaseFunc))
	L.SetField(resourceMt, "usage", L.NewFunction(m.resourceUsageFunc))
	L.SetGlobal("resource", resourceMt)
}

func (m *ModernDSL) registerSecurityPolicies(L *lua.LState) {
	// security namespace
	securityMt := L.NewTable()
	L.SetField(securityMt, "sandbox", L.NewFunction(m.securitySandboxFunc))
	L.SetField(securityMt, "policy", L.NewFunction(m.securityPolicyFunc))
	L.SetGlobal("security", securityMt)
}

// Task builder function implementations
func (m *ModernDSL) taskBuilderFunc(L *lua.LState) int {
	name := L.CheckString(1)
	
	builder := &TaskBuilder{
		definition: &TaskDefinition{
			Name:      name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		context: &BuildContext{},
		chain:   []BuildStep{},
	}
	
	m.mu.Lock()
	m.builders[name] = builder
	m.mu.Unlock()
	
	// Return task builder userdata
	ud := L.NewUserData()
	ud.Value = builder
	L.SetMetatable(ud, L.GetTypeMetatable("TaskBuilder"))
	L.Push(ud)
	return 1
}

func (m *ModernDSL) workflowBuilderFunc(L *lua.LState) int {
	// This function should not be called directly anymore since we register workflow as a table
	// in registerWorkflowDefinition. Return nil to indicate error.
	L.Push(lua.LNil)
	return 1
}

func (m *ModernDSL) chainBuilderFunc(L *lua.LState) int {
	// Implement chain builder
	return 0
}

func (m *ModernDSL) parallelBuilderFunc(L *lua.LState) int {
	// Implement parallel builder
	return 0
}

func (m *ModernDSL) conditionalBuilderFunc(L *lua.LState) int {
	// Implement conditional builder
	return 0
}

// Additional supporting types
type TaskGroup struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Tasks       []*TaskDefinition      `json:"tasks"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type TaskTemplate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Template    string                 `json:"template"`
	Parameters  map[string]Parameter   `json:"parameters"`
}

type TaskValidator interface {
	Validate(task *TaskDefinition) error
}

type BuildContext struct {
	Variables   map[string]interface{} `json:"variables"`
	Environment map[string]string      `json:"environment"`
	Workspace   string                 `json:"workspace"`
}

type BuildStep struct {
	Type    string      `json:"type"`
	Config  interface{} `json:"config"`
	Applied bool        `json:"applied"`
}

type RetryConfig struct {
	MaxAttempts int           `json:"max_attempts"`
	Delay       time.Duration `json:"delay"`
	Backoff     string        `json:"backoff"`
	Jitter      bool          `json:"jitter"`
}

type ResourceRequirements struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	Disk   string `json:"disk"`
}

type SecurityPolicy struct {
	Sandbox    bool     `json:"sandbox"`
	AllowedOps []string `json:"allowed_ops"`
	DeniedOps  []string `json:"denied_ops"`
}

type DelegationConfig struct {
	Agent   string            `json:"agent"`
	Filters map[string]string `json:"filters"`
}

type SagaConfig struct {
	Enabled      bool              `json:"enabled"`
	Compensation map[string]string `json:"compensation"`
}

type CircuitConfig struct {
	Enabled     bool          `json:"enabled"`
	Threshold   int           `json:"threshold"`
	Timeout     time.Duration `json:"timeout"`
	ResetTime   time.Duration `json:"reset_time"`
}

type Validator struct {
	Type   string      `json:"type"`
	Rule   string      `json:"rule"`
	Config interface{} `json:"config"`
}

// NewTaskRegistry creates a new task registry
func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		tasks:      make(map[string]*TaskDefinition),
		groups:     make(map[string]*TaskGroup),
		templates:  make(map[string]*TaskTemplate),
		validators: make(map[string]TaskValidator),
	}
}

// Placeholder implementations for DSL functions
func (m *ModernDSL) workflowDefineFunc(L *lua.LState) int {
	workflowName := L.CheckString(1)
	
	// Check if second argument is a table (legacy syntax) or if it's missing (fluent syntax)
	if L.GetTop() >= 2 && L.Get(2).Type() == lua.LTTable {
		// Legacy table-based syntax
		workflowConfig := L.CheckTable(2)
		
		// Get existing TaskDefinitions table or create it
		taskDefs := L.GetGlobal("TaskDefinitions")
		if taskDefs.Type() != lua.LTTable {
			taskDefs = L.NewTable()
			L.SetGlobal("TaskDefinitions", taskDefs)
		}
		
		// Create group structure compatible with legacy parser
		groupTable := L.NewTable()
		
		// Set description
		if desc := workflowConfig.RawGetString("description"); desc != lua.LNil {
			groupTable.RawSetString("description", desc)
		}
		
		// Convert Modern DSL tasks to legacy format
		if tasksValue := workflowConfig.RawGetString("tasks"); tasksValue.Type() == lua.LTTable {
			tasksTable := L.NewTable()
			taskIndex := 1
			
			tasksValue.(*lua.LTable).ForEach(func(_, taskValue lua.LValue) {
				if taskValue.Type() == lua.LTUserData {
					taskUD := taskValue.(*lua.LUserData)
					if taskDef, ok := taskUD.Value.(*TaskDefinition); ok {
						// Create legacy task structure
						legacyTask := L.NewTable()
						legacyTask.RawSetString("name", lua.LString(taskDef.Name))
						legacyTask.RawSetString("description", lua.LString(taskDef.Description))
						
						// Convert command
						if taskDef.Command != nil {
							if luaValue, ok := taskDef.Command.(lua.LValue); ok {
								legacyTask.RawSetString("command", luaValue)
							}
						}
						
						tasksTable.RawSetInt(taskIndex, legacyTask)
						taskIndex++
					}
				}
			})
			
			groupTable.RawSetString("tasks", tasksTable)
		}
		
		// Add group to TaskDefinitions
		taskDefs.(*lua.LTable).RawSetString(workflowName, groupTable)
		
		return 0
	} else {
		// New fluent syntax - return WorkflowBuilder
		builder := &WorkflowBuilder{
			name:     workflowName,
			config:   make(map[string]interface{}),
			metadata: make(map[string]interface{}),
		}
		
		// Return workflow builder userdata
		ud := L.NewUserData()
		ud.Value = builder
		L.SetMetatable(ud, L.GetTypeMetatable("WorkflowBuilder"))
		L.Push(ud)
		return 1
	}
}
func (m *ModernDSL) workflowParallelFunc(L *lua.LState) int    { return 0 }
func (m *ModernDSL) workflowSequenceFunc(L *lua.LState) int    { return 0 }
func (m *ModernDSL) workflowConditionalFunc(L *lua.LState) int { return 0 }
func (m *ModernDSL) asyncParallelFunc(L *lua.LState) int       { return 0 }
func (m *ModernDSL) asyncSequenceFunc(L *lua.LState) int       { return 0 }
func (m *ModernDSL) asyncTimeoutFunc(L *lua.LState) int        { return 0 }
func (m *ModernDSL) perfMeasureFunc(L *lua.LState) int         { return 0 }
func (m *ModernDSL) perfStatsFunc(L *lua.LState) int           { return 0 }
func (m *ModernDSL) coreStatsFunc(L *lua.LState) int           { return 0 }
func (m *ModernDSL) coreResourcesFunc(L *lua.LState) int       { return 0 }
func (m *ModernDSL) utilsConfigFunc(L *lua.LState) int         { return 0 }
func (m *ModernDSL) utilsSecretFunc(L *lua.LState) int         { return 0 }
func (m *ModernDSL) utilsEnvFunc(L *lua.LState) int            { return 0 }
func (m *ModernDSL) validateRequiredFunc(L *lua.LState) int    { return 0 }
func (m *ModernDSL) validateTypeFunc(L *lua.LState) int        { return 0 }
func (m *ModernDSL) validateRangeFunc(L *lua.LState) int       { return 0 }
func (m *ModernDSL) templateRenderFunc(L *lua.LState) int      { return 0 }
func (m *ModernDSL) templateLoadFunc(L *lua.LState) int        { return 0 }
func (m *ModernDSL) sagaBeginFunc(L *lua.LState) int           { return 0 }
func (m *ModernDSL) sagaCompensateFunc(L *lua.LState) int      { return 0 }
func (m *ModernDSL) sagaCommitFunc(L *lua.LState) int          { return 0 }
func (m *ModernDSL) circuitProtectFunc(L *lua.LState) int      { return 0 }
func (m *ModernDSL) circuitStatusFunc(L *lua.LState) int       { return 0 }
func (m *ModernDSL) resourceAllocateFunc(L *lua.LState) int    { return 0 }
func (m *ModernDSL) resourceReleaseFunc(L *lua.LState) int     { return 0 }
func (m *ModernDSL) resourceUsageFunc(L *lua.LState) int       { return 0 }
func (m *ModernDSL) securitySandboxFunc(L *lua.LState) int     { return 0 }
func (m *ModernDSL) securityPolicyFunc(L *lua.LState) int      { return 0 }

// ✅ createTaskThisObject creates the 'this' object for task functions
func (m *ModernDSL) createTaskThisObject(L *lua.LState, taskDef *TaskDefinition) *lua.LUserData {
	// Create 'this' userdata
	thisUD := L.NewUserData()
	thisUD.Value = taskDef
	
	// Create metatable for 'this' object
	thisMt := L.NewTypeMetatable("TaskThis")
	L.SetField(thisMt, "__index", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)
		
		taskDef, ok := ud.Value.(*TaskDefinition)
		if !ok {
			L.ArgError(1, "TaskThis expected")
			return 0
		}
		
		switch key {
		case "name":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				L.Push(lua.LString(taskDef.Name))
				return 1
			}))
		case "workdir":
			// Return workdir object with methods - usando dois pontos
			workdirObj := m.createWorkdirObjectWithColons(L, taskDef)
			L.Push(workdirObj)
		default:
			L.Push(lua.LNil)
		}
		return 1
	}))
	
	L.SetMetatable(thisUD, thisMt)
	return thisUD
}

// ✅ createWorkdirObjectWithColons creates the workdir object with colon methods (this:workdir:method)
func (m *ModernDSL) createWorkdirObjectWithColons(L *lua.LState, taskDef *TaskDefinition) *lua.LUserData {
	workdirUD := L.NewUserData()
	workdirUD.Value = taskDef
	
	// Create metatable for workdir object with colon syntax support
	workdirMt := L.NewTypeMetatable("TaskWorkdirColons")
	L.SetField(workdirMt, "__index", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)
		
		taskDef, ok := ud.Value.(*TaskDefinition)
		if !ok {
			L.ArgError(1, "TaskWorkdirColons expected")
			return 0
		}
		
		switch key {
		case "get":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if taskDef.Workdir != "" {
					L.Push(lua.LString(taskDef.Workdir))
				} else {
					if cwd, err := os.Getwd(); err == nil {
						L.Push(lua.LString(cwd))
					} else {
						L.Push(lua.LString("/tmp"))
					}
				}
				return 1
			}))
		case "ensure":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				workdirPath := taskDef.Workdir
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}
				
				// Remove existing directory
				os.RemoveAll(workdirPath)
				
				// Create new directory
				if err := os.MkdirAll(workdirPath, 0755); err != nil {
					L.Push(lua.LBool(false))
					return 1
				}
				
				L.Push(lua.LBool(true))
				return 1
			}))
		case "exists":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				workdirPath := taskDef.Workdir
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}
				
				if _, err := os.Stat(workdirPath); err == nil {
					L.Push(lua.LBool(true))
				} else {
					L.Push(lua.LBool(false))
				}
				return 1
			}))
		case "cleanup":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				workdirPath := taskDef.Workdir
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}
				
				if err := os.RemoveAll(workdirPath); err != nil {
					L.Push(lua.LBool(false))
					return 1
				}
				
				L.Push(lua.LBool(true))
				return 1
			}))
		case "recreate":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				workdirPath := taskDef.Workdir
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}
				
				// Remove and recreate
				os.RemoveAll(workdirPath)
				if err := os.MkdirAll(workdirPath, 0755); err != nil {
					L.Push(lua.LBool(false))
					return 1
				}
				
				L.Push(lua.LBool(true))
				return 1
			}))
		default:
			L.Push(lua.LNil)
		}
		return 1
	}))
	
	L.SetMetatable(workdirUD, workdirMt)
	return workdirUD
}

// ✅ createWorkdirObject creates the workdir object with methods
func (m *ModernDSL) createWorkdirObject(L *lua.LState, taskDef *TaskDefinition) *lua.LUserData {
	workdirUD := L.NewUserData()
	workdirUD.Value = taskDef
	
	// Create metatable for workdir object
	workdirMt := L.NewTypeMetatable("TaskWorkdir")
	L.SetField(workdirMt, "__index", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)
		
		taskDef, ok := ud.Value.(*TaskDefinition)
		if !ok {
			L.ArgError(1, "TaskWorkdir expected")
			return 0
		}
		
		switch key {
		case "get":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if taskDef.Workdir != "" {
					L.Push(lua.LString(taskDef.Workdir))
				} else {
					if cwd, err := os.Getwd(); err == nil {
						L.Push(lua.LString(cwd))
					} else {
						L.Push(lua.LString("/tmp"))
					}
				}
				return 1
			}))
		case "exists":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				workdirPath := taskDef.Workdir
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}
				
				if _, err := os.Stat(workdirPath); err == nil {
					L.Push(lua.LBool(true))
				} else {
					L.Push(lua.LBool(false))
				}
				return 1
			}))
		case "cleanup":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				workdirPath := taskDef.Workdir
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					L.Push(lua.LString("no workdir specified"))
					return 2
				}
				
				if err := os.RemoveAll(workdirPath); err != nil {
					L.Push(lua.LBool(false))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				
				L.Push(lua.LBool(true))
				L.Push(lua.LString("workdir cleaned up successfully"))
				return 2
			}))
		case "recreate":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				workdirPath := taskDef.Workdir
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					L.Push(lua.LString("no workdir specified"))
					return 2
				}
				
				// Remove and recreate
				os.RemoveAll(workdirPath)
				if err := os.MkdirAll(workdirPath, 0755); err != nil {
					L.Push(lua.LBool(false))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				
				L.Push(lua.LBool(true))
				L.Push(lua.LString("workdir recreated successfully"))
				return 2
			}))
		default:
			L.Push(lua.LNil)
		}
		return 1
	}))
	
	L.SetMetatable(workdirUD, workdirMt)
	return workdirUD
}

// Workdir management functions
func (m *ModernDSL) workdirGetFunc(L *lua.LState) int {
	// Get current workdir from environment or context
	taskContext := L.GetGlobal("__task_context")
	if taskContext.Type() == lua.LTTable {
		workdir := taskContext.(*lua.LTable).RawGetString("workdir")
		if workdir.Type() == lua.LTString {
			L.Push(workdir)
			return 1
		}
	}
	
	// Fallback to current working directory
	if cwd, err := os.Getwd(); err == nil {
		L.Push(lua.LString(cwd))
	} else {
		L.Push(lua.LString("/tmp"))
	}
	return 1
}

func (m *ModernDSL) workdirCleanupFunc(L *lua.LState) int {
	// Get workdir path (optional argument)
	var workdirPath string
	if L.GetTop() >= 1 {
		workdirPath = L.CheckString(1)
	} else {
		// Get from context
		taskContext := L.GetGlobal("__task_context")
		if taskContext.Type() == lua.LTTable {
			workdir := taskContext.(*lua.LTable).RawGetString("workdir")
			if workdir.Type() == lua.LTString {
				workdirPath = workdir.String()
			}
		}
	}
	
	if workdirPath == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("no workdir specified"))
		return 2
	}
	
	// Remove the directory
	if err := os.RemoveAll(workdirPath); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString("workdir cleaned up successfully"))
	return 2
}

func (m *ModernDSL) workdirExistsFunc(L *lua.LState) int {
	// Get workdir path (optional argument)
	var workdirPath string
	if L.GetTop() >= 1 {
		workdirPath = L.CheckString(1)
	} else {
		// Get from context
		taskContext := L.GetGlobal("__task_context")
		if taskContext.Type() == lua.LTTable {
			workdir := taskContext.(*lua.LTable).RawGetString("workdir")
			if workdir.Type() == lua.LTString {
				workdirPath = workdir.String()
			}
		}
	}
	
	if workdirPath == "" {
		L.Push(lua.LBool(false))
		return 1
	}
	
	// Check if directory exists
	if _, err := os.Stat(workdirPath); err == nil {
		L.Push(lua.LBool(true))
	} else {
		L.Push(lua.LBool(false))
	}
	return 1
}

func (m *ModernDSL) workdirCreateFunc(L *lua.LState) int {
	// Get workdir path (required argument)
	workdirPath := L.CheckString(1)
	
	// Create directory with all parent directories
	if err := os.MkdirAll(workdirPath, 0755); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString("workdir created successfully"))
	return 2
}

// Global Modern DSL instance
var globalModernDSL *ModernDSL
var modernDSLOnce sync.Once

// OpenModernDSL initializes the Modern DSL for a Lua state
func OpenModernDSL(L *lua.LState) {
	// Create a singleton ModernDSL instance
	modernDSLOnce.Do(func() {
		globalModernDSL = &ModernDSL{
			taskRegistry: &TaskRegistry{
				tasks:       make(map[string]*TaskDefinition),
				groups:      make(map[string]*TaskGroup),
				templates:   make(map[string]*TaskTemplate),
				validators:  make(map[string]TaskValidator),
			},
			builders:     make(map[string]*TaskBuilder),
		}
	})
	
	// Register the Modern DSL functions
	globalModernDSL.RegisterModernDSL(L)
}
