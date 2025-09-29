package taskrunner

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/AlecAivazis/survey/v2"

	"github.com/chalkan3-sloth/sloth-runner/internal/core"
	"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	lua "github.com/yuin/gopher-lua"
)

type SurveyAsker interface {
	AskOne(survey.Prompt, interface{}, ...survey.AskOpt) error
}

type DefaultSurveyAsker struct{}

func (d *DefaultSurveyAsker) AskOne(p survey.Prompt, r interface{}, o ...survey.AskOpt) error {
	return survey.AskOne(p, r, o...)
}

// executeShellCondition executes a shell command and returns true if it succeeds (exit code 0).
func executeShellCondition(command string) (bool, error) {
	if command == "" {
		return false, fmt.Errorf("command cannot be empty")
	}
	cmd := luainterface.ExecCommand("bash", "-c", command)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Command executed but returned a non-zero exit code.
			return false, nil
		}
		// Other error (e.g., command not found).
		return false, err
	}
	// Command succeeded.
	return true, nil
}

// TaskExecutionError provides a more context-rich error for task failures.
type TaskExecutionError struct {
	TaskName string
	Err      error
}

func (e *TaskExecutionError) Error() string {
	return fmt.Sprintf("task '%s' failed: %v", e.TaskName, e.Err)
}

type TaskRunner struct {
	L           *lua.LState
	TaskGroups  map[string]types.TaskGroup
	TargetGroup string
	TargetTasks []string
	Results     []types.TaskResult
	Outputs     map[string]interface{}
	Exports     map[string]interface{}
	DryRun      bool
	Interactive bool
	surveyAsker SurveyAsker
	LuaScript   string
	
	// Core integration
	globalCore *core.GlobalCore
	logger     *slog.Logger
}

func NewTaskRunner(L *lua.LState, groups map[string]types.TaskGroup, targetGroup string, targetTasks []string, dryRun bool, interactive bool, asker SurveyAsker, luaScript string) *TaskRunner {
	// Initialize or get existing global core
	globalCore := core.GetGlobalCore()
	if globalCore == nil {
		logger := slog.Default()
		config := core.DefaultCoreConfig()
		var err error
		if err = core.InitializeGlobalCore(config, logger); err != nil {
			slog.Error("Failed to initialize global core", "error", err)
		}
		globalCore = core.GetGlobalCore()
	}
	
	return &TaskRunner{
		L:           L,
		TaskGroups:  groups,
		TargetGroup: targetGroup,
		TargetTasks: targetTasks,
		Outputs:     make(map[string]interface{}),
		Exports:     make(map[string]interface{}),
		DryRun:      dryRun,
		Interactive: interactive,
		surveyAsker: asker,
		LuaScript:   luaScript,
		globalCore:  globalCore,
		logger:      slog.Default(),
	}
}

func (tr *TaskRunner) Export(data map[string]interface{}) {
	for key, value := range data {
		tr.Exports[key] = value
	}
}

func (tr *TaskRunner) executeTaskWithRetries(t *types.Task, inputFromDependencies *lua.LTable, mu *sync.Mutex, completedTasks map[string]bool, taskOutputs map[string]*lua.LTable, runningTasks map[string]bool, session *types.SharedSession, groupName string) error {
	// Use core for error recovery and performance tracking
	return tr.globalCore.ExecuteWithRecovery(func() error {
		// AbortIf check with circuit breaker for external commands
		if t.AbortIfFunc != nil {
			shouldAbort, _, _, err := luainterface.ExecuteLuaFunction(tr.L, t.AbortIfFunc, t.Params, inputFromDependencies, 1, nil)
			if err != nil {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to execute abort_if function: %w", err)}
			}
			if shouldAbort {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("execution aborted by abort_if function")}
			}
		} else if t.AbortIf != "" {
			// Execute shell command directly (simplified)
			shouldAbort, err := executeShellCondition(t.AbortIf)
			if err != nil {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to execute abort_if condition: %w", err)}
			}
			if shouldAbort {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("execution aborted by abort_if condition")}
			}
		}

		// RunIf check
		if t.RunIfFunc != nil {
			shouldRun, _, _, err := luainterface.ExecuteLuaFunction(tr.L, t.RunIfFunc, t.Params, inputFromDependencies, 1, nil)
			if err != nil {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to execute run_if function: %w", err)}
			}
			if !shouldRun {
				pterm.Info.Printf("Skipping task '%s' due to run_if function condition.\n", t.Name)
				mu.Lock()
				tr.Results = append(tr.Results, types.TaskResult{
					Name:   t.Name,
					Status: "Skipped",
				})
				completedTasks[t.Name] = true
				delete(runningTasks, t.Name)
				mu.Unlock()
				return nil
			}
		} else if t.RunIf != "" {
			// Execute shell command directly (simplified)
			shouldRun, err := executeShellCondition(t.RunIf)
			
			if err != nil {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to execute run_if condition: %w", err)}
			}
			if !shouldRun {
				pterm.Info.Printf("Skipping task '%s' due to run_if condition.\n", t.Name)
				mu.Lock()
				tr.Results = append(tr.Results, types.TaskResult{
					Name:   t.Name,
					Status: "Skipped",
				})
				completedTasks[t.Name] = true
				delete(runningTasks, t.Name)
				mu.Unlock()
				return nil
			}
		}

		var taskErr error
		maxRetries := t.Retries
		if maxRetries < 0 {
			maxRetries = 0
		}

		for i := 0; i <= maxRetries; i++ {
			if i > 0 {
				backoffDelay := time.Duration(i) * time.Second
				if i > 3 {
					backoffDelay = time.Duration(i*i) * time.Second // Exponential backoff
				}
				pterm.Warning.Printf("Task '%s' failed. Retrying in %v (%d/%d)...\n", t.Name, backoffDelay, i, maxRetries)
				time.Sleep(backoffDelay)
			}

			slog.Info("starting task", "task", t.Name, "attempt", i+1, "retries", maxRetries)

			var ctx context.Context
			var cancel context.CancelFunc

			if t.Timeout != "" {
				timeout, err := time.ParseDuration(t.Timeout)
				if err != nil {
					return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("invalid timeout duration: %w", err)}
				}
				ctx, cancel = context.WithTimeout(context.Background(), timeout)
			} else {
				// Use default timeout from core config
				defaultTimeout := tr.globalCore.Config.TimeoutDefault
				ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
			}
			defer cancel()

			taskErr = tr.runTask(ctx, t, inputFromDependencies, mu, completedTasks, taskOutputs, runningTasks, session, groupName)

			if taskErr == nil {
				slog.Info("task finished", "task", t.Name, "status", "success")
				return nil // Success
			}

			// Log retry attempt for monitoring
			tr.logger.Warn("task retry", "task", t.Name, "attempt", i+1, "error", taskErr)
		}

		slog.Error("task failed", "task", t.Name, "retries", maxRetries, "err", taskErr)
		return taskErr // Final failure
	}, fmt.Sprintf("task_%s_%s", groupName, t.Name))
}

func (tr *TaskRunner) runTask(ctx context.Context, t *types.Task, inputFromDependencies *lua.LTable, mu *sync.Mutex, completedTasks map[string]bool, taskOutputs map[string]*lua.LTable, runningTasks map[string]bool, session *types.SharedSession, groupName string) (taskErr error) {
	startTime := time.Now()

	var agentAddress string

	// Determine agent address from task's DelegateTo or group's DelegateTo
	if t.DelegateTo != nil {
		switch v := t.DelegateTo.(type) {
		case string:
			agentAddress = v // Direct address
		case map[string]interface{}:
			if addr, ok := v["address"].(string); ok {
				agentAddress = addr
			} else {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("invalid agent definition in task delegate_to: missing address")}
			}
		default:
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("invalid type for task delegate_to: %T", v)}
		}
	} else if tr.TaskGroups[groupName].DelegateTo != nil {
		switch v := tr.TaskGroups[groupName].DelegateTo.(type) {
		case string:
			agentAddress = v // Direct address
		case map[string]interface{}:
			if addr, ok := v["address"].(string); ok {
				agentAddress = addr
			} else {
				return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("invalid agent definition in group delegate_to: missing address")}
			}
		default:
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("invalid type for group delegate_to: %T", v)}
		}
	}

	if agentAddress != "" {
		// Connect to the agent
		conn, err := grpc.Dial(agentAddress, grpc.WithInsecure())
		if err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to connect to agent %s: %w", agentAddress, err)}
		}
		defer conn.Close()
		c := pb.NewAgentClient(conn)

		// Create a tarball of the workspace
		var buf bytes.Buffer
		if err := createTar(session.Workdir, &buf); err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to create workspace tarball: %w", err)}
		}

		// Send the task and workspace to the agent
		r, err := c.ExecuteTask(ctx, &pb.ExecuteTaskRequest{
			TaskName:    t.Name,
			TaskGroup:   groupName,
			LuaScript:   tr.LuaScript,
			Workspace:   buf.Bytes(),
		})
		if err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to execute task on agent %s: %w", agentAddress, err)}
		}

		if !r.GetSuccess() {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("task failed on agent %s: %s", agentAddress, r.GetOutput())}
		}

		// Extract the updated workspace
		if err := extractTar(bytes.NewReader(r.GetWorkspace()), session.Workdir); err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("failed to extract updated workspace from agent %s: %w", agentAddress, err)}
		}

		return nil
	}

	L := lua.NewState()
	defer L.Close()
	luainterface.OpenAll(L)
	
	localInputFromDependencies := luainterface.CopyTable(inputFromDependencies, L)
	
t.Output = L.NewTable()

	defer func() {
		if r := recover(); r != nil {
			taskErr = &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("panic: %v", r)}
		}

		duration := time.Since(startTime)
		status := "Success"
		if taskErr != nil {
			status = "Failed"
		}

		mu.Lock()
		tr.Results = append(tr.Results, types.TaskResult{
			Name:     t.Name,
			Status:   status,
			Duration: duration,
			Error:    taskErr,
		})
		taskOutputs[t.Name] = luainterface.CopyTable(t.Output, tr.L)
		completedTasks[t.Name] = true
		delete(runningTasks, t.Name)
		mu.Unlock()
	}()

	if t.PreExec != nil {
		success, msg, _, err := luainterface.ExecuteLuaFunction(L, t.PreExec, t.Params, localInputFromDependencies, 2, ctx)
		if err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("error executing pre_exec hook: %w", err)}
		} else if !success {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("pre-execution hook failed: %s", msg)}
		}
	}

	if t.CommandFunc != nil {
		if t.Params == nil {
			t.Params = make(map[string]string)
		}
		t.Params["task_name"] = t.Name
		t.Params["group_name"] = groupName
		t.Params["workdir"] = session.Workdir

		var sessionUD *lua.LUserData
		if session != nil {
			sessionUD = L.NewUserData()
			sessionUD.Value = session
			L.SetMetatable(sessionUD, L.GetTypeMetatable("session"))
		}

		success, msg, outputTable, err := luainterface.ExecuteLuaFunction(L, t.CommandFunc, t.Params, localInputFromDependencies, 3, ctx, sessionUD)
		if err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("error executing command function: %w", err)}
		} else if !success {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("command function returned failure: %s", msg)}
		} else if outputTable != nil {
			t.Output = outputTable
		}
	}

	if t.PostExec != nil {
		var postExecSecondArg lua.LValue = t.Output
		if t.Output == nil {
			postExecSecondArg = L.NewTable()
		}
		success, msg, _, err := luainterface.ExecuteLuaFunction(L, t.PostExec, t.Params, postExecSecondArg, 2, ctx)
		if err != nil {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("error executing post_exec hook: %w", err)}
		} else if !success {
			return &TaskExecutionError{TaskName: t.Name, Err: fmt.Errorf("post-execution hook failed: %s", msg)}
		}
	}

	return nil
}

// Run executes the task groups and tasks defined in the TaskRunner.
// It orchestrates the entire execution process, including:
// - Filtering task groups if a target group is specified.
// - Setting up and cleaning up work directories.
// - Resolving the correct task execution order based on dependencies.
// - Displaying a real-time progress bar using pterm.
// - Executing each task sequentially, respecting dependency statuses.
// - Collecting results and outputs.
// - Rendering a final summary table.
func (tr *TaskRunner) Run() error {
	if len(tr.TaskGroups) == 0 {
		slog.Warn("No task groups defined.")
		return nil
	}

	var allGroupErrors []error

	filteredGroups := make(map[string]types.TaskGroup)
	if tr.TargetGroup != "" {
		if group, ok := tr.TaskGroups[tr.TargetGroup]; ok {
			filteredGroups[tr.TargetGroup] = group
		} else {
			return fmt.Errorf("task group '%s' not found", tr.TargetGroup)
		}
	} else {
		filteredGroups = tr.TaskGroups
	}

	for groupName, group := range filteredGroups {
		slog.Info("starting group", "group", groupName, "description", group.Description)

		var workdir string
		var err error
		if group.Workdir != "" {
			workdir = group.Workdir
		} else if group.CreateWorkdirBeforeRun {
			uuid, err := uuid.NewRandom()
			if err != nil {
				return fmt.Errorf("failed to generate UUID for workdir: %w", err)
			}
			workdir = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s", groupName, uuid.String()))
			if err := os.RemoveAll(workdir); err != nil {
				return fmt.Errorf("failed to clean fixed workdir %s: %w", workdir, err)
			}
		} else {
			workdir, err = ioutil.TempDir(os.TempDir(), groupName+"-*")
			if err != nil {
				return fmt.Errorf("failed to create ephemeral workdir: %w", err)
			}
		}

		if err := os.MkdirAll(workdir, 0755); err != nil {
			return fmt.Errorf("failed to create workdir %s: %w", workdir, err)
		}

		artifactsBaseDir := "artifacts" // Persistent artifacts directory in project root
		artifactsGroupDir := filepath.Join(artifactsBaseDir, groupName)
		artifactsTaskRunDir := filepath.Join(artifactsGroupDir, time.Now().Format("20060102-150405")) // Timestamped directory for each run

		if err := os.MkdirAll(artifactsTaskRunDir, 0755); err != nil {
			return fmt.Errorf("failed to create persistent artifacts directory %s: %w", artifactsTaskRunDir, err)
		}
		artifactsDir := artifactsTaskRunDir // Use this as the destination for artifacts

		session := &types.SharedSession{
			Workdir: workdir,
		}

		taskMap := make(map[string]*types.Task)
		for i := range group.Tasks {
			taskMap[group.Tasks[i].Name] = &group.Tasks[i]
		}

		tasksToRun, err := tr.resolveTasksToRun(taskMap, tr.TargetTasks)
		if err != nil {
			return err
		}

		executionOrder, err := tr.getExecutionOrder(tasksToRun)
		if err != nil {
			return err
		}

		p, _ := pterm.DefaultProgressbar.WithTotal(len(executionOrder)).WithTitle("Executing tasks").Start()
		defer p.Stop()

		var mu sync.Mutex
		completedTasks := make(map[string]bool)
		taskOutputs := make(map[string]*lua.LTable)
		runningTasks := make(map[string]bool)
		taskStatus := make(map[string]string)
		var groupErrors []error

		for _, taskName := range executionOrder {
			p.UpdateTitle("Executing task: " + taskName)
			task := taskMap[taskName]
			runningTasks[task.Name] = true

			// Dependency checks
			skip := false
			for _, depName := range task.DependsOn {
				if status, ok := taskStatus[depName]; !ok || (status != "Success" && status != "Skipped") {
					slog.Warn("Skipping task due to dependency failure", "task", task.Name, "dependency", depName, "dep_status", taskStatus[depName])
					skip = true
					break
				}
			}
			if skip {
				taskStatus[task.Name] = "Skipped"
				p.Increment()
				continue
			}

			// Consume artifacts 
			for _, artifactName := range task.Consumes {
				srcPath := filepath.Join(artifactsDir, artifactName)
				destPath := filepath.Join(workdir, artifactName)
				
				if err := copyFile(srcPath, destPath); err != nil {
					slog.Error("Failed to consume artifact", "task", task.Name, "artifact", artifactName, "error", err)
					groupErrors = append(groupErrors, err)
					taskStatus[task.Name] = "Failed"
					skip = true
					break
				}
				
				slog.Info("Consumed artifact", "task", task.Name, "artifact", artifactName)
			}
			if skip {
				p.Increment()
				continue
			}

			inputFromDependencies := tr.L.NewTable()
			for _, depName := range task.DependsOn {
				if output, ok := taskOutputs[depName]; ok {
					inputFromDependencies.RawSetString(depName, output)
				}
			}

			if tr.Interactive {
				action := ""
				prompt := &survey.Select{
					Message: fmt.Sprintf("Task: %s (%s)", task.Name, task.Description),
					Options: []string{"run", "skip", "abort", "continue"},
					Default: "run",
				}
				tr.surveyAsker.AskOne(prompt, &action)

				switch action {
				case "skip":
					pterm.Info.Printf("Skipping task '%s' by user choice.\n", task.Name)
					taskStatus[task.Name] = "Skipped"
					p.Increment()
					continue
				case "abort":
					pterm.Warning.Println("Aborting execution by user choice.")
					return fmt.Errorf("execution aborted by user")
				case "continue":
					tr.Interactive = false // Disable interactive mode for subsequent tasks
				}
			}

			err := tr.executeTaskWithRetries(task, inputFromDependencies, &mu, completedTasks, taskOutputs, runningTasks, session, groupName)
			if err != nil {
				groupErrors = append(groupErrors, err)
				taskStatus[task.Name] = "Failed"
			} else {
				taskStatus[task.Name] = "Success"

				// Produce artifacts
				for _, artifactPattern := range task.Artifacts {
					matches, err := filepath.Glob(filepath.Join(workdir, artifactPattern))
					if err != nil {
						slog.Error("Invalid artifact pattern", "task", task.Name, "pattern", artifactPattern, "error", err)
						continue
					}
					for _, match := range matches {
						destPath := filepath.Join(artifactsDir, filepath.Base(match))
						if err := copyFile(match, destPath); err != nil {
							slog.Error("Failed to produce artifact", "task", task.Name, "artifact", match, "error", err)
						} else {
							slog.Info("Produced artifact", "task", task.Name, "artifact", destPath)
						}
					}
				}
			}
			p.Increment()
		}

		groupHadSuccess := len(groupErrors) == 0
		if !groupHadSuccess {
			allGroupErrors = append(allGroupErrors, fmt.Errorf("task group '%s' encountered errors", groupName))
		}

		mu.Lock()
		for name, outputTable := range taskOutputs {
			tr.Outputs[name] = luainterface.LuaTableToGoMap(tr.L, outputTable)
		}
		mu.Unlock()

		shouldClean := true
		if group.CleanWorkdirAfterRunFunc != nil {
			L := lua.NewState()
			defer L.Close()
			luainterface.OpenAll(L)

			resultTable := L.NewTable()
			resultTable.RawSetString("success", lua.LBool(groupHadSuccess))
			if !groupHadSuccess && len(groupErrors) > 0 {
				resultTable.RawSetString("error", lua.LString(groupErrors[0].Error()))
			}
			// Find the output of the last task to run
			if len(executionOrder) > 0 {
				lastTaskName := executionOrder[len(executionOrder)-1]
				if output, ok := taskOutputs[lastTaskName]; ok {
					resultTable.RawSetString("output", output)
				}
			}

			success, _, _, err := luainterface.ExecuteLuaFunction(L, group.CleanWorkdirAfterRunFunc, nil, resultTable, 1, context.Background(), lua.LNil, lua.LString("clean_workdir_after_run"))
			if err != nil {
				slog.Error("Error executing clean_workdir_after_run", "group", groupName, "err", err)
			} else {
				shouldClean = success
			}
		}

		if shouldClean {
			slog.Info("Cleaning up workdir", "group", groupName, "workdir", workdir)
			os.RemoveAll(workdir)
		} else {
			slog.Warn("Workdir preserved", "group", groupName, "workdir", workdir)
		}
	}

	pterm.DefaultSection.Println("Execution Summary")
	tableData := pterm.TableData{{"Task", "Status", "Duration", "Error"}}
	for _, result := range tr.Results {
		status := pterm.Green(result.Status)
		errStr := ""
		if result.Error != nil {
			status = pterm.Red(result.Status)
			errStr = result.Error.Error()
		} else if result.Status == "Skipped" {
			status = pterm.Yellow(result.Status)
		} else if result.Status == "DryRun" {
			status = pterm.Cyan(result.Status)
		}
		tableData = append(tableData, []string{result.Name, status, result.Duration.String(), errStr})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	if len(allGroupErrors) > 0 {
		return fmt.Errorf("one or more task groups failed")
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func (tr *TaskRunner) getExecutionOrder(tasksToRun []*types.Task) ([]string, error) {
	taskMap := make(map[string]*types.Task)
	for _, task := range tasksToRun {
		taskMap[task.Name] = task
	}

	var order []string
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	var visit func(taskName string) error
	visit = func(taskName string) error {
		if recursionStack[taskName] {
			return fmt.Errorf("circular dependency detected: %s", taskName)
		}
		if visited[taskName] {
			return nil
		}

		recursionStack[taskName] = true
		visited[taskName] = true

		task := taskMap[taskName]
		depNames := task.DependsOn
		sort.Strings(depNames)

		for _, depName := range depNames {
			if _, ok := taskMap[depName]; !ok {
				continue
			}
			if err := visit(depName); err != nil {
				return err
			}
		}

		order = append(order, taskName)
		delete(recursionStack, taskName)
		return nil
	}

	var taskNames []string
	for _, task := range tasksToRun {
		taskNames = append(taskNames, task.Name)
	}
	sort.Strings(taskNames)

	for _, taskName := range taskNames {
		if !visited[taskName] {
			if err := visit(taskName); err != nil {
				return nil, err
			}
		}
	}

	return order, nil
}

func (tr *TaskRunner) resolveTasksToRun(originalTaskMap map[string]*types.Task, targetTasks []string) ([]*types.Task, error) {
	if len(targetTasks) == 0 {
		var allTasks []*types.Task
		for _, task := range originalTaskMap {
			allTasks = append(allTasks, task)
		}
		return allTasks, nil
	}

	resolved := make(map[string]*types.Task)
	queue := make([]string, 0, len(targetTasks))
	visited := make(map[string]bool)

	for _, taskName := range targetTasks {
		if !visited[taskName] {
			queue = append(queue, taskName)
			visited[taskName] = true
		}
	}

	head := 0
	for head < len(queue) {
		currentTaskName := queue[head]
		head++

		currentTask, ok := originalTaskMap[currentTaskName]
		if !ok {
			return nil, fmt.Errorf("task '%s' not found in group", currentTaskName)
		}
		resolved[currentTaskName] = currentTask

		for _, depName := range currentTask.DependsOn {
			if !visited[depName] {
				visited[depName] = true
				queue = append(queue, depName)
			}
		}
	}

	var result []*types.Task
	for _, task := range resolved {
		result = append(result, task)
	}
	return result, nil
}

// RunTasksParallel executes a slice of tasks concurrently and waits for them to complete.
func (tr *TaskRunner) RunTasksParallel(tasks []*types.Task, input *lua.LTable) ([]types.TaskResult, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	resultsChan := make(chan types.TaskResult, len(tasks))
	errChan := make(chan error, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t *types.Task) {
			defer wg.Done()

			completed := make(map[string]bool)
			outputs := make(map[string]*lua.LTable)
			running := make(map[string]bool)

			var taskMu sync.Mutex

			err := tr.executeTaskWithRetries(t, input, &taskMu, completed, outputs, running, nil, "")

			mu.Lock()
			var result types.TaskResult
			for i := len(tr.Results) - 1; i >= 0; i-- {
				if tr.Results[i].Name == t.Name {
					result = tr.Results[i]
					break
				}
			}
			mu.Unlock()

			if err != nil {
				errChan <- err
			}
			resultsChan <- result

		}(task)
	}

	wg.Wait()
	close(resultsChan)
	close(errChan)

	var allErrors []error
	for err := range errChan {
		allErrors = append(allErrors, err)
	}

	var results []types.TaskResult
	for result := range resultsChan {
		results = append(results, result)
	}

	if len(allErrors) > 0 {
		return results, fmt.Errorf("encountered %d errors during parallel execution: %v", len(allErrors), allErrors)
	}

	return results, nil
}

// createTar function to create a tarball of a directory
func createTar(source string, writer io.Writer) error {
	tw := tar.NewWriter(writer)
	defer tw.Close()
	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(file[len(source):])
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !fi.IsDir() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})
}

// extractTar function to extract a tarball to a directory
func extractTar(reader io.Reader, dest string) error {
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}
