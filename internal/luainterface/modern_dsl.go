package luainterface

import (
	"log/slog"
	"sync"
	"time"

	"github.com/chalkan3/sloth-runner/internal/core"
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

// Core registration methods
func (m *ModernDSL) registerTaskDefinition(L *lua.LState) {
	// task() - main task builder function
	L.SetGlobal("task", L.NewFunction(m.taskBuilderFunc))
	
	// workflow() - workflow definition function
	L.SetGlobal("workflow", L.NewFunction(m.workflowBuilderFunc))
	
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
	// Implement workflow builder
	return 0
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
func (m *ModernDSL) workflowDefineFunc(L *lua.LState) int      { return 0 }
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
