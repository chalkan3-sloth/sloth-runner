package core

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/sloth"
	lua "github.com/yuin/gopher-lua"
)

// SlothModule provides automation capabilities for sloth-runner itself
type SlothModule struct{}

// RegisterSlothModule registers the sloth module in the Lua state
func RegisterSlothModule(L *lua.LState) {
	// Create sloth table
	slothTable := L.NewTable()

	// Agent management functions
	agentTable := L.NewTable()
	L.SetField(agentTable, "install", L.NewFunction(slothAgentInstall))
	L.SetField(agentTable, "update", L.NewFunction(slothAgentUpdate))
	L.SetField(agentTable, "list", L.NewFunction(slothAgentList))
	L.SetField(agentTable, "get", L.NewFunction(slothAgentGet))
	L.SetField(agentTable, "delete", L.NewFunction(slothAgentDelete))
	L.SetField(agentTable, "start", L.NewFunction(slothAgentStart))
	L.SetField(agentTable, "stop", L.NewFunction(slothAgentStop))
	L.SetField(slothTable, "agent", agentTable)

	// Sloth (workflow) management functions
	workflowTable := L.NewTable()
	L.SetField(workflowTable, "add", L.NewFunction(slothWorkflowAdd))
	L.SetField(workflowTable, "list", L.NewFunction(slothWorkflowList))
	L.SetField(workflowTable, "get", L.NewFunction(slothWorkflowGet))
	L.SetField(workflowTable, "remove", L.NewFunction(slothWorkflowRemove))
	L.SetField(workflowTable, "activate", L.NewFunction(slothWorkflowActivate))
	L.SetField(workflowTable, "deactivate", L.NewFunction(slothWorkflowDeactivate))
	L.SetField(slothTable, "workflow", workflowTable)

	// SSH management functions
	sshTable := L.NewTable()
	L.SetField(sshTable, "add", L.NewFunction(slothSSHAdd))
	L.SetField(sshTable, "list", L.NewFunction(slothSSHList))
	L.SetField(sshTable, "remove", L.NewFunction(slothSSHRemove))
	L.SetField(slothTable, "ssh", sshTable)

	// Stack management
	stackTable := L.NewTable()
	L.SetField(stackTable, "list", L.NewFunction(slothStackList))
	L.SetField(stackTable, "get", L.NewFunction(slothStackGet))
	L.SetField(stackTable, "delete", L.NewFunction(slothStackDelete))
	L.SetField(slothTable, "stack", stackTable)

	// Run workflow
	L.SetField(slothTable, "run", L.NewFunction(slothRun))

	// Set global
	L.SetGlobal("sloth", slothTable)
}

// Helper function to run sloth-runner CLI commands
func runSlothCommand(args ...string) (string, error) {
	// Find sloth-runner in PATH or use current executable
	executable, err := exec.LookPath("sloth-runner")
	if err != nil {
		// Fallback to current executable
		executable, err = os.Executable()
		if err != nil {
			return "", fmt.Errorf("failed to find sloth-runner: %w", err)
		}
	}

	cmd := exec.Command(executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	output := stdout.String()
	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%s: %s", err.Error(), stderr.String())
		}
		return output, err
	}

	return output, nil
}

// slothAgentInstall installs a new agent on a remote host
// Usage: sloth.agent.install({name = "agent-name", ssh_host = "192.168.1.10", master = "192.168.1.29:50053", ...})
func slothAgentInstall(L *lua.LState) int {
	opts := L.CheckTable(1)

	name := getStringField(L, opts, "name", "")
	sshHost := getStringField(L, opts, "ssh_host", "")
	sshUser := getStringField(L, opts, "ssh_user", "root")
	sshPort := getStringField(L, opts, "ssh_port", "22")
	sshKey := getStringField(L, opts, "ssh_key", "")
	master := getStringField(L, opts, "master", "localhost:50051")
	bindAddr := getStringField(L, opts, "bind_address", "0.0.0.0")
	port := getStringField(L, opts, "port", "50051")
	reportAddr := getStringField(L, opts, "report_address", "")

	if name == "" {
		return luaError(L, "name is required")
	}
	if sshHost == "" {
		return luaError(L, "ssh_host is required")
	}

	// Check if agent already exists (idempotency)
	checkArgs := []string{"agent", "list", "--master", master}
	output, err := runSlothCommand(checkArgs...)
	if err == nil && strings.Contains(output, name) {
		// Agent already exists
		result := L.NewTable()
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "message", lua.LString(fmt.Sprintf("Agent '%s' already exists", name)))
		L.SetField(result, "name", lua.LString(name))
		L.Push(result)
		return 1
	}

	// Build install command
	args := []string{
		"agent", "install", name,
		"--ssh-host", sshHost,
		"--ssh-user", sshUser,
		"--ssh-port", sshPort,
		"--master", master,
		"--bind-address", bindAddr,
		"--port", port,
	}

	if sshKey != "" {
		args = append(args, "--ssh-key", sshKey)
	}
	if reportAddr != "" {
		args = append(args, "--report-address", reportAddr)
	}

	// Install agent
	output, err = runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to install agent: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Agent '%s' installed successfully", name)))
	L.SetField(result, "name", lua.LString(name))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothAgentUpdate updates an agent to the latest version
// Usage: sloth.agent.update({name = "agent-name", master = "192.168.1.29:50053", version = "latest"})
func slothAgentUpdate(L *lua.LState) int {
	opts := L.CheckTable(1)

	name := getStringField(L, opts, "name", "")
	master := getStringField(L, opts, "master", "localhost:50051")
	version := getStringField(L, opts, "version", "latest")

	if name == "" {
		return luaError(L, "name is required")
	}

	// Get current agent info
	getArgs := []string{"agent", "get", name, "--master", master}
	output, err := runSlothCommand(getArgs...)
	if err != nil {
		return luaError(L, fmt.Sprintf("agent not found: %v", err))
	}

	// Check current version (parse from output)
	// If already on target version, return unchanged
	if version != "latest" && strings.Contains(output, version) {
		result := L.NewTable()
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "message", lua.LString(fmt.Sprintf("Agent '%s' already on version %s", name, version)))
		L.Push(result)
		return 1
	}

	// Update agent
	args := []string{"agent", "update", name, "--master", master}
	if version != "" && version != "latest" {
		args = append(args, "--version", version)
	}

	output, err = runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to update agent: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Agent '%s' updated successfully", name)))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothAgentList lists all agents
// Usage: sloth.agent.list({master = "192.168.1.29:50053"})
func slothAgentList(L *lua.LState) int {
	opts := L.CheckTable(1)
	master := getStringField(L, opts, "master", "localhost:50051")

	args := []string{"agent", "list", "--master", master}
	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to list agents: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(false))
	L.SetField(result, "message", lua.LString("Agents listed successfully"))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothAgentGet gets agent details
// Usage: sloth.agent.get({name = "agent-name", master = "192.168.1.29:50053"})
func slothAgentGet(L *lua.LState) int {
	opts := L.CheckTable(1)
	name := getStringField(L, opts, "name", "")
	master := getStringField(L, opts, "master", "localhost:50051")

	if name == "" {
		return luaError(L, "name is required")
	}

	args := []string{"agent", "get", name, "--master", master}
	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to get agent: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(false))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Agent '%s' details retrieved", name)))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothAgentDelete deletes an agent
// Usage: sloth.agent.delete({name = "agent-name", master = "192.168.1.29:50053"})
func slothAgentDelete(L *lua.LState) int {
	opts := L.CheckTable(1)
	name := getStringField(L, opts, "name", "")
	master := getStringField(L, opts, "master", "localhost:50051")

	if name == "" {
		return luaError(L, "name is required")
	}

	// Check if agent exists (idempotency)
	checkArgs := []string{"agent", "list", "--master", master}
	output, err := runSlothCommand(checkArgs...)
	if err != nil || !strings.Contains(output, name) {
		result := L.NewTable()
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "message", lua.LString(fmt.Sprintf("Agent '%s' does not exist", name)))
		L.Push(result)
		return 1
	}

	args := []string{"agent", "delete", name, "--master", master, "--yes"}
	output, err = runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to delete agent: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Agent '%s' deleted successfully", name)))
	L.Push(result)
	return 1
}

// slothAgentStart starts a local agent
// Usage: sloth.agent.start({bind_address = "0.0.0.0", port = 50051, master = "192.168.1.29:50053"})
func slothAgentStart(L *lua.LState) int {
	opts := L.CheckTable(1)
	bindAddr := getStringField(L, opts, "bind_address", "0.0.0.0")
	port := getStringField(L, opts, "port", "50051")
	master := getStringField(L, opts, "master", "localhost:50051")
	reportAddr := getStringField(L, opts, "report_address", "")
	daemon := getBoolField(L, opts, "daemon", true)

	args := []string{
		"agent", "start",
		"--bind-address", bindAddr,
		"--port", port,
		"--master", master,
	}

	if reportAddr != "" {
		args = append(args, "--report-address", reportAddr)
	}

	daemonStr := "false"
	if daemon {
		daemonStr = "true"
	}
	args = append(args, "--daemon", daemonStr)

	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to start agent: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString("Agent started successfully"))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothAgentStop stops a local agent
// Usage: sloth.agent.stop({})
func slothAgentStop(L *lua.LState) int {
	L.CheckTable(1)

	args := []string{"agent", "stop"}
	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to stop agent: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString("Agent stopped successfully"))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothWorkflowAdd adds a workflow to the database
// Usage: sloth.workflow.add({name = "workflow-name", file = "/path/to/file.sloth", description = "...", active = true})
func slothWorkflowAdd(L *lua.LState) int {
	opts := L.CheckTable(1)

	name := getStringField(L, opts, "name", "")
	file := getStringField(L, opts, "file", "")
	description := getStringField(L, opts, "description", "")
	active := getBoolField(L, opts, "active", true)

	if name == "" {
		return luaError(L, "name is required")
	}
	if file == "" {
		return luaError(L, "file is required")
	}

	// Check if workflow already exists (idempotency)
	homeDir, _ := os.UserHomeDir()
	dbPath := homeDir + "/.sloth-runner/sloths.db"
	repo, err := sloth.NewSQLiteRepository(dbPath)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to open database: %v", err))
	}
	defer repo.Close()

	ctx := context.Background()
	existing, err := repo.GetByName(ctx, name)
	if err == nil && existing != nil {
		result := L.NewTable()
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' already exists", name)))
		L.Push(result)
		return 1
	}

	// Build command
	args := []string{"sloth", "add", name, "--file", file}
	if description != "" {
		args = append(args, "--description", description)
	}
	activeStr := "true"
	if !active {
		activeStr = "false"
	}
	args = append(args, "--active="+activeStr)

	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to add workflow: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' added successfully", name)))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothWorkflowList lists all workflows
// Usage: sloth.workflow.list({active_only = false})
func slothWorkflowList(L *lua.LState) int {
	opts := L.CheckTable(1)
	activeOnly := getBoolField(L, opts, "active_only", false)

	args := []string{"sloth", "list"}
	if activeOnly {
		args = append(args, "--active")
	}

	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to list workflows: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(false))
	L.SetField(result, "message", lua.LString("Workflows listed successfully"))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothWorkflowGet gets workflow details
// Usage: sloth.workflow.get({name = "workflow-name"})
func slothWorkflowGet(L *lua.LState) int {
	opts := L.CheckTable(1)
	name := getStringField(L, opts, "name", "")

	if name == "" {
		return luaError(L, "name is required")
	}

	args := []string{"sloth", "get", name}
	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to get workflow: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(false))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' details retrieved", name)))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothWorkflowRemove removes a workflow
// Usage: sloth.workflow.remove({name = "workflow-name"})
func slothWorkflowRemove(L *lua.LState) int {
	opts := L.CheckTable(1)
	name := getStringField(L, opts, "name", "")

	if name == "" {
		return luaError(L, "name is required")
	}

	// Check if exists (idempotency)
	homeDir, _ := os.UserHomeDir()
	dbPath := homeDir + "/.sloth-runner/sloths.db"
	repo, err := sloth.NewSQLiteRepository(dbPath)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to open database: %v", err))
	}
	defer repo.Close()

	ctx := context.Background()
	existing, err := repo.GetByName(ctx, name)
	if err != nil || existing == nil {
		result := L.NewTable()
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' does not exist", name)))
		L.Push(result)
		return 1
	}

	args := []string{"sloth", "remove", name, "--yes"}
	_, err = runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to remove workflow: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' removed successfully", name)))
	L.Push(result)
	return 1
}

// slothWorkflowActivate activates a workflow
// Usage: sloth.workflow.activate({name = "workflow-name"})
func slothWorkflowActivate(L *lua.LState) int {
	opts := L.CheckTable(1)
	name := getStringField(L, opts, "name", "")

	if name == "" {
		return luaError(L, "name is required")
	}

	// Check current status (idempotency)
	homeDir, _ := os.UserHomeDir()
	dbPath := homeDir + "/.sloth-runner/sloths.db"
	repo, err := sloth.NewSQLiteRepository(dbPath)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to open database: %v", err))
	}
	defer repo.Close()

	ctx := context.Background()
	existing, err := repo.GetByName(ctx, name)
	if err != nil {
		return luaError(L, fmt.Sprintf("workflow not found: %v", err))
	}

	if existing.IsActive {
		result := L.NewTable()
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' is already active", name)))
		L.Push(result)
		return 1
	}

	args := []string{"sloth", "activate", name}
	_, err = runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to activate workflow: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' activated successfully", name)))
	L.Push(result)
	return 1
}

// slothWorkflowDeactivate deactivates a workflow
// Usage: sloth.workflow.deactivate({name = "workflow-name"})
func slothWorkflowDeactivate(L *lua.LState) int {
	opts := L.CheckTable(1)
	name := getStringField(L, opts, "name", "")

	if name == "" {
		return luaError(L, "name is required")
	}

	// Check current status (idempotency)
	homeDir, _ := os.UserHomeDir()
	dbPath := homeDir + "/.sloth-runner/sloths.db"
	repo, err := sloth.NewSQLiteRepository(dbPath)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to open database: %v", err))
	}
	defer repo.Close()

	ctx := context.Background()
	existing, err := repo.GetByName(ctx, name)
	if err != nil {
		return luaError(L, fmt.Sprintf("workflow not found: %v", err))
	}

	if !existing.IsActive {
		result := L.NewTable()
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' is already inactive", name)))
		L.Push(result)
		return 1
	}

	args := []string{"sloth", "deactivate", name}
	_, err = runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to deactivate workflow: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Workflow '%s' deactivated successfully", name)))
	L.Push(result)
	return 1
}

// slothSSHAdd adds an SSH profile
// Usage: sloth.ssh.add({name = "profile-name", host = "192.168.1.10", user = "root", port = 22, key = "..."})
func slothSSHAdd(L *lua.LState) int {
	opts := L.CheckTable(1)

	name := getStringField(L, opts, "name", "")
	host := getStringField(L, opts, "host", "")
	user := getStringField(L, opts, "user", "root")
	port := getStringField(L, opts, "port", "22")
	key := getStringField(L, opts, "key", "")

	if name == "" {
		return luaError(L, "name is required")
	}
	if host == "" {
		return luaError(L, "host is required")
	}

	args := []string{"ssh", "add", name, "--host", host, "--user", user, "--port", port}
	if key != "" {
		args = append(args, "--key", key)
	}

	output, err := runSlothCommand(args...)
	if err != nil {
		// Check if already exists
		if strings.Contains(err.Error(), "already exists") {
			result := L.NewTable()
			L.SetField(result, "changed", lua.LBool(false))
			L.SetField(result, "message", lua.LString(fmt.Sprintf("SSH profile '%s' already exists", name)))
			L.Push(result)
			return 1
		}
		return luaError(L, fmt.Sprintf("failed to add SSH profile: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("SSH profile '%s' added successfully", name)))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothSSHList lists SSH profiles
// Usage: sloth.ssh.list({})
func slothSSHList(L *lua.LState) int {
	L.CheckTable(1)

	args := []string{"ssh", "list"}
	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to list SSH profiles: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(false))
	L.SetField(result, "message", lua.LString("SSH profiles listed successfully"))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothSSHRemove removes an SSH profile
// Usage: sloth.ssh.remove({name = "profile-name"})
func slothSSHRemove(L *lua.LState) int {
	opts := L.CheckTable(1)
	name := getStringField(L, opts, "name", "")

	if name == "" {
		return luaError(L, "name is required")
	}

	args := []string{"ssh", "remove", name, "--yes"}
	_, err := runSlothCommand(args...)
	if err != nil {
		// Check if doesn't exist
		if strings.Contains(err.Error(), "not found") {
			result := L.NewTable()
			L.SetField(result, "changed", lua.LBool(false))
			L.SetField(result, "message", lua.LString(fmt.Sprintf("SSH profile '%s' does not exist", name)))
			L.Push(result)
			return 1
		}
		return luaError(L, fmt.Sprintf("failed to remove SSH profile: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("SSH profile '%s' removed successfully", name)))
	L.Push(result)
	return 1
}

// slothStackList lists workflow stacks
// Usage: sloth.stack.list({})
func slothStackList(L *lua.LState) int {
	L.CheckTable(1)

	args := []string{"stack", "list"}
	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to list stacks: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(false))
	L.SetField(result, "message", lua.LString("Stacks listed successfully"))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothStackGet gets stack details
// Usage: sloth.stack.get({id = "stack-id"})
func slothStackGet(L *lua.LState) int {
	opts := L.CheckTable(1)
	id := getStringField(L, opts, "id", "")

	if id == "" {
		return luaError(L, "id is required")
	}

	args := []string{"stack", "show", id}
	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to get stack: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(false))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Stack '%s' details retrieved", id)))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// slothStackDelete deletes a stack
// Usage: sloth.stack.delete({id = "stack-id"})
func slothStackDelete(L *lua.LState) int {
	opts := L.CheckTable(1)
	id := getStringField(L, opts, "id", "")

	if id == "" {
		return luaError(L, "id is required")
	}

	args := []string{"stack", "delete", id, "--yes"}
	_, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to delete stack: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString(fmt.Sprintf("Stack '%s' deleted successfully", id)))
	L.Push(result)
	return 1
}

// slothRun executes a workflow
// Usage: sloth.run({task = "task-name", file = "/path/to/file.sloth", sloth = "saved-workflow", ...})
func slothRun(L *lua.LState) int {
	opts := L.CheckTable(1)

	task := getStringField(L, opts, "task", "")
	file := getStringField(L, opts, "file", "")
	slothName := getStringField(L, opts, "sloth", "")
	delegateTo := getStringField(L, opts, "delegate_to", "")
	values := getStringField(L, opts, "values", "")

	args := []string{"run"}

	if task != "" {
		args = append(args, task)
	}

	if file != "" {
		args = append(args, "--file", file)
	} else if slothName != "" {
		args = append(args, "--sloth", slothName)
	} else {
		return luaError(L, "either 'file' or 'sloth' is required")
	}

	if delegateTo != "" {
		args = append(args, "--delegate-to", delegateTo)
	}

	if values != "" {
		args = append(args, "--values", values)
	}

	args = append(args, "--yes")

	output, err := runSlothCommand(args...)
	if err != nil {
		return luaError(L, fmt.Sprintf("failed to run workflow: %v", err))
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "message", lua.LString("Workflow executed successfully"))
	L.SetField(result, "output", lua.LString(output))
	L.Push(result)
	return 1
}

// luaError returns an error to Lua
func luaError(L *lua.LState, msg string) int {
	L.Push(lua.LNil)
	L.Push(lua.LString(msg))
	return 2
}
