package luainterface

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/gitops"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	lua "github.com/yuin/gopher-lua"
)

// registerGitOpsModule registers the GitOps module for Lua
func (li *LuaInterface) registerGitOpsModule() {
	li.L.PreloadModule("gitops", li.loadGitOpsModule)
}

// loadGitOpsModule loads the GitOps module into Lua
func (li *LuaInterface) loadGitOpsModule(L *lua.LState) int {
	// Create GitOps module table
	gitopsTable := L.NewTable()
	
	// Workflow management functions
	L.SetField(gitopsTable, "workflow", L.NewFunction(li.luaGitOpsWorkflow))
	L.SetField(gitopsTable, "create_workflow", L.NewFunction(li.luaGitOpsCreateWorkflow))
	L.SetField(gitopsTable, "sync_workflow", L.NewFunction(li.luaGitOpsSyncWorkflow))
	L.SetField(gitopsTable, "get_workflow_status", L.NewFunction(li.luaGitOpsGetWorkflowStatus))
	L.SetField(gitopsTable, "list_workflows", L.NewFunction(li.luaGitOpsListWorkflows))
	
	// Repository management functions
	L.SetField(gitopsTable, "register_repository", L.NewFunction(li.luaGitOpsRegisterRepository))
	
	// Diff and preview functions
	L.SetField(gitopsTable, "generate_diff", L.NewFunction(li.luaGitOpsGenerateDiff))
	L.SetField(gitopsTable, "preview_changes", L.NewFunction(li.luaGitOpsPreviewChanges))
	
	// Rollback functions
	L.SetField(gitopsTable, "rollback_workflow", L.NewFunction(li.luaGitOpsRollbackWorkflow))
	
	// Utility functions
	L.SetField(gitopsTable, "start_auto_sync", L.NewFunction(li.luaGitOpsStartAutoSync))
	L.SetField(gitopsTable, "stop_auto_sync", L.NewFunction(li.luaGitOpsStopAutoSync))
	
	L.Push(gitopsTable)
	return 1
}

// luaGitOpsWorkflow creates a GitOps workflow (simplified API)
func (li *LuaInterface) luaGitOpsWorkflow(L *lua.LState) int {
	config := L.CheckTable(1)
	
	// Extract configuration
	repo := li.getStringFieldFromTable(config, "repo")
	branch := li.getStringFieldFromTable(config, "branch")
	if branch == "" {
		branch = "main"
	}
	
	autoSync := li.getBoolFieldFromTable(config, "auto_sync")
	diffPreview := li.getBoolFieldFromTable(config, "diff_preview")
	rollbackOnFailure := li.getBoolFieldFromTable(config, "rollback_on_failure")
	
	if repo == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("repository URL is required"))
		return 2
	}
	
	// Get GitOps manager
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	ctx := context.Background()
	
	// Generate unique IDs
	repoID := li.generateID("repo")
	workflowID := li.generateID("workflow")
	
	// Register repository if not exists
	repository := &gitops.Repository{
		ID:           repoID,
		URL:          repo,
		Branch:       branch,
		PollInterval: 30 * time.Second,
	}
	
	if err := manager.RegisterRepository(ctx, repository); err != nil {
		slog.Warn("Failed to register repository (may already exist)", "error", err)
	}
	
	// Create workflow
	workflow := &gitops.Workflow{
		ID:                workflowID,
		Name:              "lua-generated-workflow",
		Repository:        repoID,
		TargetPath:        ".",
		AutoSync:          autoSync,
		DiffPreview:       diffPreview,
		RollbackOnFailure: rollbackOnFailure,
		SyncPolicy: gitops.SyncPolicy{
			AutoPrune: true,
			Retry: gitops.RetryPolicy{
				Limit: 3,
				Backoff: gitops.BackoffPolicy{
					Duration: 10 * time.Second,
					Factor:   2.0,
				},
			},
			HealthCheck: gitops.HealthCheck{
				Enabled: true,
				Timeout: 5 * time.Minute,
			},
		},
	}
	
	if err := manager.CreateWorkflow(ctx, workflow); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Return workflow info
	result := L.NewTable()
	L.SetField(result, "workflow_id", lua.LString(workflowID))
	L.SetField(result, "repository_id", lua.LString(repoID))
	L.SetField(result, "status", lua.LString("created"))
	L.SetField(result, "auto_sync", lua.LBool(autoSync))
	L.SetField(result, "diff_preview", lua.LBool(diffPreview))
	L.SetField(result, "rollback_on_failure", lua.LBool(rollbackOnFailure))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaGitOpsCreateWorkflow creates a detailed GitOps workflow
func (li *LuaInterface) luaGitOpsCreateWorkflow(L *lua.LState) int {
	config := L.CheckTable(1)
	
	// Extract detailed configuration
	workflowID := li.getStringFieldFromTable(config, "id")
	name := li.getStringFieldFromTable(config, "name")
	repositoryID := li.getStringFieldFromTable(config, "repository")
	targetPath := li.getStringFieldFromTable(config, "target_path")
	
	if workflowID == "" {
		workflowID = li.generateID("workflow")
	}
	if targetPath == "" {
		targetPath = "."
	}
	
	// Get GitOps manager
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	// Create workflow
	workflow := &gitops.Workflow{
		ID:                workflowID,
		Name:              name,
		Repository:        repositoryID,
		TargetPath:        targetPath,
		AutoSync:          li.getBoolFieldFromTable(config, "auto_sync"),
		DiffPreview:       li.getBoolFieldFromTable(config, "diff_preview"),
		RollbackOnFailure: li.getBoolFieldFromTable(config, "rollback_on_failure"),
	}
	
	ctx := context.Background()
	if err := manager.CreateWorkflow(ctx, workflow); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(workflowID))
	return 1
}

// luaGitOpsSyncWorkflow manually triggers a workflow sync
func (li *LuaInterface) luaGitOpsSyncWorkflow(L *lua.LState) int {
	workflowID := L.CheckString(1)
	
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	ctx := context.Background()
	if err := manager.SyncWorkflow(ctx, workflowID); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// luaGitOpsGetWorkflowStatus gets the status of a workflow
func (li *LuaInterface) luaGitOpsGetWorkflowStatus(L *lua.LState) int {
	workflowID := L.CheckString(1)
	
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	workflow, err := manager.GetWorkflowStatus(workflowID)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Convert workflow to Lua table
	result := L.NewTable()
	L.SetField(result, "id", lua.LString(workflow.ID))
	L.SetField(result, "name", lua.LString(workflow.Name))
	L.SetField(result, "status", lua.LString(string(workflow.Status)))
	L.SetField(result, "auto_sync", lua.LBool(workflow.AutoSync))
	L.SetField(result, "repository", lua.LString(workflow.Repository))
	
	if workflow.LastSyncResult != nil {
		syncResult := L.NewTable()
		L.SetField(syncResult, "id", lua.LString(workflow.LastSyncResult.ID))
		L.SetField(syncResult, "status", lua.LString(string(workflow.LastSyncResult.Status)))
		L.SetField(syncResult, "start_time", lua.LString(workflow.LastSyncResult.StartTime.Format(time.RFC3339)))
		L.SetField(syncResult, "commit_hash", lua.LString(workflow.LastSyncResult.CommitHash))
		L.SetField(syncResult, "message", lua.LString(workflow.LastSyncResult.Message))
		
		// Add metrics
		metrics := L.NewTable()
		L.SetField(metrics, "duration", lua.LString(workflow.LastSyncResult.Metrics.Duration.String()))
		L.SetField(metrics, "resources_processed", lua.LNumber(workflow.LastSyncResult.Metrics.ResourcesProcessed))
		L.SetField(metrics, "resources_applied", lua.LNumber(workflow.LastSyncResult.Metrics.ResourcesApplied))
		L.SetField(syncResult, "metrics", metrics)
		
		L.SetField(result, "last_sync_result", syncResult)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaGitOpsListWorkflows lists all workflows
func (li *LuaInterface) luaGitOpsListWorkflows(L *lua.LState) int {
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(L.NewTable())
		return 1
	}
	
	workflows := manager.ListWorkflows()
	
	result := L.NewTable()
	i := 1
	for _, workflow := range workflows {
		workflowTable := L.NewTable()
		L.SetField(workflowTable, "id", lua.LString(workflow.ID))
		L.SetField(workflowTable, "name", lua.LString(workflow.Name))
		L.SetField(workflowTable, "status", lua.LString(string(workflow.Status)))
		L.SetField(workflowTable, "repository", lua.LString(workflow.Repository))
		
		result.RawSetInt(i, workflowTable)
		i++
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaGitOpsRegisterRepository registers a new repository
func (li *LuaInterface) luaGitOpsRegisterRepository(L *lua.LState) int {
	config := L.CheckTable(1)
	
	repoID := li.getStringFieldFromTable(config, "id")
	url := li.getStringFieldFromTable(config, "url")
	branch := li.getStringFieldFromTable(config, "branch")
	
	if repoID == "" {
		repoID = li.generateID("repo")
	}
	if branch == "" {
		branch = "main"
	}
	
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	repository := &gitops.Repository{
		ID:           repoID,
		URL:          url,
		Branch:       branch,
		PollInterval: 30 * time.Second,
	}
	
	ctx := context.Background()
	if err := manager.RegisterRepository(ctx, repository); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(repoID))
	return 1
}

// luaGitOpsGenerateDiff generates a diff preview for a workflow
func (li *LuaInterface) luaGitOpsGenerateDiff(L *lua.LState) int {
	workflowID := L.CheckString(1)
	
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	ctx := context.Background()
	diff, err := manager.GenerateDiff(ctx, workflowID)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Convert diff to Lua table
	result := L.NewTable()
	L.SetField(result, "workflow_id", lua.LString(diff.WorkflowID))
	L.SetField(result, "generated_at", lua.LString(diff.GeneratedAt.Format(time.RFC3339)))
	
	// Summary
	summary := L.NewTable()
	L.SetField(summary, "total_changes", lua.LNumber(diff.Summary.TotalChanges))
	L.SetField(summary, "created_resources", lua.LNumber(diff.Summary.CreatedResources))
	L.SetField(summary, "updated_resources", lua.LNumber(diff.Summary.UpdatedResources))
	L.SetField(summary, "deleted_resources", lua.LNumber(diff.Summary.DeletedResources))
	L.SetField(summary, "conflict_count", lua.LNumber(diff.Summary.ConflictCount))
	L.SetField(result, "summary", summary)
	
	// Changes
	changes := L.NewTable()
	for i, change := range diff.Changes {
		changeTable := L.NewTable()
		L.SetField(changeTable, "type", lua.LString(string(change.Type)))
		L.SetField(changeTable, "resource", lua.LString(change.Resource.Kind+"/"+change.Resource.Name))
		L.SetField(changeTable, "diff", lua.LString(change.Diff))
		L.SetField(changeTable, "impact", lua.LString(string(change.Impact)))
		changes.RawSetInt(i+1, changeTable)
	}
	L.SetField(result, "changes", changes)
	
	// Warnings
	warnings := L.NewTable()
	for i, warning := range diff.Warnings {
		warnings.RawSetInt(i+1, lua.LString(warning))
	}
	L.SetField(result, "warnings", warnings)
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaGitOpsPreviewChanges is an alias for generate_diff
func (li *LuaInterface) luaGitOpsPreviewChanges(L *lua.LState) int {
	return li.luaGitOpsGenerateDiff(L)
}

// luaGitOpsRollbackWorkflow rolls back a workflow
func (li *LuaInterface) luaGitOpsRollbackWorkflow(L *lua.LState) int {
	workflowID := L.CheckString(1)
	reason := L.OptString(2, "Manual rollback")
	
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	ctx := context.Background()
	if err := manager.RollbackWorkflow(ctx, workflowID, reason); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// luaGitOpsStartAutoSync starts auto-sync for all workflows
func (li *LuaInterface) luaGitOpsStartAutoSync(L *lua.LState) int {
	manager := li.getGitOpsManager()
	if manager == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("GitOps manager not available"))
		return 2
	}
	
	ctx := context.Background()
	manager.StartAutoSync(ctx)
	
	L.Push(lua.LBool(true))
	return 1
}

// luaGitOpsStopAutoSync stops auto-sync (placeholder)
func (li *LuaInterface) luaGitOpsStopAutoSync(L *lua.LState) int {
	// In real implementation, this would stop the auto-sync controller
	L.Push(lua.LBool(true))
	return 1
}

// Helper methods
func (li *LuaInterface) getGitOpsManager() *gitops.GitOpsManager {
	if li.gitopsManager == nil {
		// Create state manager if not exists
		if li.stateManager == nil {
			stateManager, err := state.NewStateManager("")
			if err != nil {
				slog.Error("Failed to create state manager for GitOps", "error", err)
				return nil
			}
			li.stateManager = stateManager
		}
		
		li.gitopsManager = gitops.NewGitOpsManager(li.stateManager)
	}
	return li.gitopsManager
}

func (li *LuaInterface) getStringFieldFromTable(table *lua.LTable, field string) string {
	value := table.RawGetString(field)
	if str, ok := value.(lua.LString); ok {
		return string(str)
	}
	return ""
}

func (li *LuaInterface) getBoolFieldFromTable(table *lua.LTable, field string) bool {
	value := table.RawGetString(field)
	if boolean, ok := value.(lua.LBool); ok {
		return bool(boolean)
	}
	return false
}

func (li *LuaInterface) generateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().Unix())
}