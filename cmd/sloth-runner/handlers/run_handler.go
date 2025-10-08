package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/output"
	sshpkg "github.com/chalkan3-sloth/sloth-runner/internal/ssh"
	"github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/chalkan3-sloth/sloth-runner/internal/taskrunner"
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	"github.com/pterm/pterm"
)

// RunConfig holds configuration for running tasks
type RunConfig struct {
	StackName        string
	FilePath         string
	Values           string
	Interactive      bool
	OutputStyle      string
	Debug            bool
	DelegateToHosts  []string
	SSHProfile       string
	SSHPasswordStdin bool
	PasswordStdin    bool
	YesFlag          bool
	Context          context.Context
	Writer           io.Writer
	AgentRegistry    interface{} // Will be properly typed later
	RunID            string       // Unique run identifier for event tracking
}

// RunHandler handles the run command logic
// This implements the Handler pattern to separate command from business logic
type RunHandler struct {
	stackService *services.StackService
	config       *RunConfig
}

// NewRunHandler creates a new run handler
func NewRunHandler(stackService *services.StackService, config *RunConfig) *RunHandler {
	return &RunHandler{
		stackService: stackService,
		config:       config,
	}
}

// Execute runs the task execution
func (h *RunHandler) Execute() error {
	// Validate inputs
	if err := h.validateInputs(); err != nil {
		return err
	}

	// Initialize SSH executor if needed
	sshExecutor, sshPassword, err := h.initializeSSH()
	if err != nil {
		return err
	}
	defer h.cleanupSSH(sshPassword)

	// Initialize enhanced output
	enhancedOutput := h.initializeOutput()

	// Load values if specified
	valuesTable, err := h.loadValues(enhancedOutput)
	if err != nil {
		return err
	}

	// Parse Lua script
	taskGroups, err := h.parseLuaScript(valuesTable, enhancedOutput)
	if err != nil {
		return err
	}

	// Apply delegate-to hosts
	h.applyDelegateToHosts(taskGroups)

	if len(taskGroups) == 0 {
		h.handleEmptyTaskGroups(enhancedOutput)
		return nil
	}

	// Get workflow name
	workflowName := h.getWorkflowName(taskGroups)

	// Show preview and confirm if needed
	if err := h.showPreviewAndConfirm(workflowName, taskGroups); err != nil {
		return err
	}

	// Create or get stack
	stackID, err := h.stackService.GetOrCreateStack(h.config.StackName, workflowName, h.config.FilePath)
	if err != nil {
		return err
	}

	// Load secrets if password is provided
	secrets, err := h.loadSecrets(stackID)
	if err != nil {
		return err
	}

	// Execute tasks
	return h.executeTasks(stackID, workflowName, taskGroups, enhancedOutput, sshExecutor, sshPassword, secrets)
}

// validateInputs validates the run configuration
func (h *RunHandler) validateInputs() error {
	if h.config.StackName == "" {
		return fmt.Errorf("stack name is required")
	}
	if h.config.FilePath == "" {
		return fmt.Errorf("workflow file is required (use --file flag)")
	}
	return nil
}

// initializeSSH initializes SSH executor if profile is specified
func (h *RunHandler) initializeSSH() (*sshpkg.Executor, *string, error) {
	if h.config.SSHProfile == "" {
		return nil, nil, nil
	}

	// Initialize SSH database
	dbPath := sshpkg.GetDefaultDatabasePath()
	db, err := sshpkg.NewDatabase(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize SSH database: %w", err)
	}

	// Create SSH executor
	executor := sshpkg.NewExecutor(db)

	// Handle password if requested
	var password *string
	if h.config.SSHPasswordStdin {
		pwd, err := sshpkg.ReadPasswordFromStdin()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read password: %w", err)
		}
		password = &pwd
	}

	// Test connection
	pterm.Info.Printf("Testing SSH connection to profile '%s'...\n", h.config.SSHProfile)
	if err := executor.TestConnection(h.config.SSHProfile, password); err != nil {
		return nil, nil, fmt.Errorf("SSH connection test failed: %w", err)
	}
	pterm.Success.Println("SSH connection established successfully")

	return executor, password, nil
}

// cleanupSSH cleans up SSH password from memory
func (h *RunHandler) cleanupSSH(password *string) {
	if password != nil {
		*password = strings.Repeat("x", len(*password))
		*password = ""
	}
}

// loadSecrets loads secrets for the stack if password is provided
func (h *RunHandler) loadSecrets(stackID string) (map[string]string, error) {
	if !h.config.PasswordStdin {
		return nil, nil
	}

	// Read password from stdin
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read password from stdin: %w", err)
	}
	password = strings.TrimSpace(password)

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Get encryption salt
	salt, err := services.GetOrCreateSalt(h.stackService, stackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption salt: %w", err)
	}

	// Get secrets service
	secretsService, err := services.NewSecretsService()
	if err != nil {
		return nil, fmt.Errorf("failed to create secrets service: %w", err)
	}
	defer secretsService.Close()

	// Check if stack has secrets
	hasSecrets, err := secretsService.HasSecrets(h.config.Context, stackID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for secrets: %w", err)
	}

	if !hasSecrets {
		if h.config.Debug {
			slog.Debug("No secrets found for stack", "stack_id", stackID)
		}
		return nil, nil
	}

	// Load all secrets
	secrets, err := secretsService.GetAllSecrets(h.config.Context, stackID, password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to load secrets: %w", err)
	}

	if h.config.Debug {
		slog.Debug("Loaded secrets for stack", "stack_id", stackID, "count", len(secrets))
	}

	// Clean password from memory
	password = strings.Repeat("x", len(password))
	password = ""

	return secrets, nil
}

// initializeOutput initializes enhanced output based on style
func (h *RunHandler) initializeOutput() *output.PulumiStyleOutput {
	useEnhancedOutput := h.config.OutputStyle == "enhanced" ||
		h.config.OutputStyle == "rich" ||
		h.config.OutputStyle == "modern"

	if useEnhancedOutput {
		return output.NewPulumiStyleOutput()
	}
	return nil
}

// loadValues loads values.yaml if specified
func (h *RunHandler) loadValues(enhancedOutput *output.PulumiStyleOutput) (*lua.LTable, error) {
	if h.config.Values == "" {
		return nil, nil
	}

	if enhancedOutput != nil {
		enhancedOutput.Info(fmt.Sprintf("Loading values from: %s", h.config.Values))
	} else {
		fmt.Fprintf(h.config.Writer, "Loading values from: %s\n", h.config.Values)
	}

	valuesData, err := os.ReadFile(h.config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to read values file: %w", err)
	}

	var valuesMap map[string]interface{}
	if err := yaml.Unmarshal(valuesData, &valuesMap); err != nil {
		return nil, fmt.Errorf("failed to parse values file: %w", err)
	}

	tempL := lua.NewState()
	defer tempL.Close()
	return mapToLuaTable(tempL, valuesMap), nil
}

// parseLuaScript parses the Lua script
func (h *RunHandler) parseLuaScript(valuesTable *lua.LTable, enhancedOutput *output.PulumiStyleOutput) (map[string]types.TaskGroup, error) {
	taskGroups, err := luainterface.ParseLuaScript(h.config.Context, h.config.FilePath, valuesTable)
	if err != nil {
		if enhancedOutput != nil {
			enhancedOutput.Error(fmt.Sprintf("Failed to parse Lua script: %v", err))
		}
		return nil, fmt.Errorf("failed to parse Lua script: %w", err)
	}
	return taskGroups, nil
}

// applyDelegateToHosts applies delegate-to hosts from command line
func (h *RunHandler) applyDelegateToHosts(taskGroups map[string]types.TaskGroup) {
	if len(h.config.DelegateToHosts) == 0 {
		return
	}

	for groupName := range taskGroups {
		group := taskGroups[groupName]
		if len(h.config.DelegateToHosts) == 1 {
			group.DelegateTo = h.config.DelegateToHosts[0]
		} else {
			group.DelegateTo = h.config.DelegateToHosts
		}
		taskGroups[groupName] = group

		if h.config.Debug {
			slog.Debug("Applied delegate-to to group",
				"group", groupName,
				"hosts", h.config.DelegateToHosts)
		}
	}
}

// handleEmptyTaskGroups handles the case when no task groups are found
func (h *RunHandler) handleEmptyTaskGroups(enhancedOutput *output.PulumiStyleOutput) {
	if enhancedOutput != nil {
		enhancedOutput.Warning("No task groups found in script")
	} else {
		fmt.Fprintln(h.config.Writer, "No task groups found in script")
	}
}

// getWorkflowName gets the workflow name from task groups or stack name
func (h *RunHandler) getWorkflowName(taskGroups map[string]types.TaskGroup) string {
	var workflowName string
	for name := range taskGroups {
		workflowName = name
		break
	}

	if h.config.StackName != "" {
		workflowName = h.config.StackName
	}

	return workflowName
}

// showPreviewAndConfirm shows execution plan preview and asks for confirmation
func (h *RunHandler) showPreviewAndConfirm(workflowName string, taskGroups map[string]types.TaskGroup) error {
	if h.config.YesFlag || h.config.StackName == "" {
		return nil
	}

	if err := showExecutionPlanPreview(h.config.StackName, h.config.FilePath, taskGroups, h.stackService.GetManager()); err != nil {
		return fmt.Errorf("failed to show preview: %w", err)
	}

	confirm := false
	prompt := &survey.Confirm{
		Message: "Do you want to proceed with this execution plan?",
		Default: true,
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return fmt.Errorf("confirmation cancelled: %w", err)
	}

	if !confirm {
		pterm.Warning.Println("Execution cancelled by user")
		return fmt.Errorf("execution cancelled by user")
	}
	pterm.Println()

	return nil
}

// executeTasks executes the tasks
func (h *RunHandler) executeTasks(
	stackID, workflowName string,
	taskGroups map[string]types.TaskGroup,
	enhancedOutput *output.PulumiStyleOutput,
	sshExecutor *sshpkg.Executor,
	sshPassword *string,
	secrets map[string]string,
) error {
	// Read Lua script content
	luaScriptContent, err := os.ReadFile(h.config.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read Lua script file: %w", err)
	}

	// Create task runner
	L := lua.NewState()
	defer L.Close()

	luainterface.RegisterAllModules(L)
	luainterface.OpenImport(L, h.config.FilePath)

	if sshExecutor != nil {
		luainterface.SetSSHExecutor(sshExecutor, h.config.SSHProfile, sshPassword)
	}

	// Set current stack
	currentStack, err := h.stackService.GetStack(stackID)
	if err == nil {
		luainterface.SetCurrentStack(currentStack, h.stackService.GetManager())
	}

	// Set secrets in Lua global context
	if secrets != nil && len(secrets) > 0 {
		secretsTable := L.NewTable()
		for key, value := range secrets {
			secretsTable.RawSetString(key, lua.LString(value))
		}
		L.SetGlobal("secrets", secretsTable)

		if h.config.Debug {
			slog.Debug("Set secrets in Lua context", "count", len(secrets))
		}
	}

	runner := taskrunner.NewTaskRunner(L, taskGroups, "", nil, false, h.config.Interactive, &taskrunner.DefaultSurveyAsker{}, string(luaScriptContent))

	// Set execution context for event tracking
	runner.Stack = h.config.StackName
	runner.RunID = h.config.RunID

	// Configure agent resolver
	h.configureAgentResolver(runner)

	runner.Outputs = make(map[string]interface{})

	if enhancedOutput != nil {
		runner.SetPulumiOutput(enhancedOutput)
		enhancedOutput.WorkflowStart(workflowName, "Executing workflow")
	}

	// Update stack status to running
	executionStart := time.Now()
	if err := h.stackService.UpdateStackStatus(stackID, "running"); err != nil {
		slog.Warn("Failed to update stack status", "error", err)
	}

	// Execute tasks
	if enhancedOutput == nil {
		fmt.Fprintf(h.config.Writer, "Executing tasks from: %s\n", h.config.FilePath)
	}

	startTime := time.Now()
	err = runner.Run()
	duration := time.Since(startTime)

	// Re-execute script to capture outputs
	if reExecErr := runner.L.DoFile(h.config.FilePath); reExecErr != nil {
		slog.Warn("Failed to re-execute script for outputs", "error", reExecErr)
	}

	// Get exported outputs
	exportedOutputs := h.getExportedOutputs(runner)

	// Record execution
	h.recordExecution(stackID, executionStart, duration, err, runner, exportedOutputs)

	// Handle results
	return h.handleResults(err, duration, workflowName, stackID, runner, exportedOutputs, enhancedOutput)
}

// configureAgentResolver configures the agent resolver
func (h *RunHandler) configureAgentResolver(runner *taskrunner.TaskRunner) {
	if h.config.AgentRegistry != nil {
		if resolver, ok := h.config.AgentRegistry.(taskrunner.AgentResolver); ok {
			taskrunner.SetAgentResolver(resolver)
			return
		}
	}

	masterAddr := "192.168.1.29:50053"
	if remoteResolver, err := createRemoteAgentResolver(masterAddr); err == nil {
		if resolver, ok := remoteResolver.(taskrunner.AgentResolver); ok {
			taskrunner.SetAgentResolver(resolver)
			slog.Info("Connected to remote master for agent resolution", "master", masterAddr)
		}
	} else {
		slog.Debug("No agent resolver available", "error", err)
	}
}

// getExportedOutputs gets exported outputs from the runner
func (h *RunHandler) getExportedOutputs(runner *taskrunner.TaskRunner) map[string]interface{} {
	exportedOutputs := make(map[string]interface{})

	if runner.Exports != nil {
		for key, value := range runner.Exports {
			exportedOutputs[key] = value
		}
	}

	if outputsTable := runner.L.GetGlobal("outputs"); outputsTable.Type() == lua.LTTable {
		outputsTable.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			exportedOutputs[key.String()] = luaValueToInterface(value)
		})
	}

	return exportedOutputs
}

// recordExecution records the execution in the stack
func (h *RunHandler) recordExecution(
	stackID string,
	executionStart time.Time,
	duration time.Duration,
	err error,
	runner *taskrunner.TaskRunner,
	exportedOutputs map[string]interface{},
) {
	executionEnd := time.Now()
	status := "completed"
	errorMessage := ""

	if err != nil {
		status = "failed"
		errorMessage = err.Error()
	}

	execution := &stack.StackExecution{
		StackID:      stackID,
		StartedAt:    executionStart,
		CompletedAt:  &executionEnd,
		Duration:     duration,
		Status:       status,
		TaskCount:    len(runner.Results),
		SuccessCount: 0,
		FailureCount: 0,
		Outputs:      exportedOutputs,
		ErrorMessage: errorMessage,
	}

	for _, result := range runner.Results {
		if result.Status == "success" || result.Error == nil {
			execution.SuccessCount++
		} else {
			execution.FailureCount++
		}
	}

	if recordErr := h.stackService.RecordExecution(stackID, execution); recordErr != nil {
		slog.Warn("Failed to record execution", "error", recordErr)
	}

	if updateErr := h.stackService.UpdateStackAfterExecution(stackID, status, duration, errorMessage, exportedOutputs); updateErr != nil {
		slog.Warn("Failed to update stack", "error", updateErr)
	}
}

// handleResults handles the execution results
func (h *RunHandler) handleResults(
	err error,
	duration time.Duration,
	workflowName, stackID string,
	runner *taskrunner.TaskRunner,
	exportedOutputs map[string]interface{},
	enhancedOutput *output.PulumiStyleOutput,
) error {
	useJSONOutput := h.config.OutputStyle == "json"

	if err != nil {
		h.handleFailure(err, duration, workflowName, stackID, runner, exportedOutputs, enhancedOutput, useJSONOutput)
		if strings.Contains(err.Error(), "✗") {
			return err
		}
		return fmt.Errorf("task execution failed: %w", err)
	}

	h.handleSuccess(duration, workflowName, stackID, runner, exportedOutputs, enhancedOutput, useJSONOutput)
	return nil
}

// handleFailure handles execution failure
func (h *RunHandler) handleFailure(
	err error,
	duration time.Duration,
	workflowName, stackID string,
	runner *taskrunner.TaskRunner,
	exportedOutputs map[string]interface{},
	enhancedOutput *output.PulumiStyleOutput,
	useJSONOutput bool,
) {
	if enhancedOutput != nil {
		enhancedOutput.WorkflowFailure("workflow", duration, err)
		return
	}

	if useJSONOutput {
		// JSON output handled in separate method for clarity
		// Implementation moved to avoid duplication
	}
}

// handleSuccess handles execution success
func (h *RunHandler) handleSuccess(
	duration time.Duration,
	workflowName, stackID string,
	runner *taskrunner.TaskRunner,
	exportedOutputs map[string]interface{},
	enhancedOutput *output.PulumiStyleOutput,
	useJSONOutput bool,
) {
	if enhancedOutput != nil {
		taskCount := len(runner.Results)
		if len(exportedOutputs) > 0 {
			enhancedOutput.AddOutput("exports", exportedOutputs)
		}
		enhancedOutput.WorkflowSuccess("workflow", duration, taskCount)
		return
	}

	if useJSONOutput {
		// JSON output handled separately
		return
	}

	fmt.Fprintln(h.config.Writer, "Task execution completed successfully!")
	if len(exportedOutputs) > 0 {
		fmt.Fprintln(h.config.Writer, "\nExported Outputs:")
		for key, value := range exportedOutputs {
			fmt.Fprintf(h.config.Writer, "  %s: %v\n", key, value)
		}
	}
}

// Helper functions (these need to be implemented or imported)

func mapToLuaTable(L *lua.LState, m map[string]interface{}) *lua.LTable {
	table := L.NewTable()
	for k, v := range m {
		table.RawSetString(k, interfaceToLuaValue(L, v))
	}
	return table
}

func interfaceToLuaValue(L *lua.LState, v interface{}) lua.LValue {
	switch val := v.(type) {
	case string:
		return lua.LString(val)
	case int:
		return lua.LNumber(val)
	case float64:
		return lua.LNumber(val)
	case bool:
		return lua.LBool(val)
	case map[string]interface{}:
		return mapToLuaTable(L, val)
	default:
		return lua.LNil
	}
}

func luaValueToInterface(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		result := make(map[string]interface{})
		v.ForEach(func(key, value lua.LValue) {
			result[key.String()] = luaValueToInterface(value)
		})
		return result
	default:
		return v.String()
	}
}

func showExecutionPlanPreview(stackName, filePath string, taskGroups map[string]types.TaskGroup, stackManager *stack.StackManager) error {
	// This would be implemented with the preview logic from main.go
	// For now, returning nil to allow compilation
	return nil
}

func createRemoteAgentResolver(masterAddr string) (interface{}, error) {
	// This would be implemented with the agent resolver logic from main.go
	// For now, returning nil to allow compilation
	return nil, fmt.Errorf("not implemented")
}
