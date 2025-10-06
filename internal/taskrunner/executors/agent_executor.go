package executors

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	lua "github.com/yuin/gopher-lua"
)

// AgentExecutor executes tasks on remote agents via gRPC
type AgentExecutor struct {
	GenerateAgentScript func(task *types.Task, groupName string) string
	CreateTar           func(source string, writer *bytes.Buffer) error
	ExtractTar          func(reader *bytes.Reader, dest string) error
}

// Execute runs the task on a remote agent via gRPC
func (ae *AgentExecutor) Execute(
	ctx context.Context,
	task *types.Task,
	agentAddress string,
	session *types.SharedSession,
	groupName string,
) error {
	// Connect to the agent
	pterm.DefaultBox.
		WithTitle("🔗 Agent Connection").
		WithTitleTopLeft().
		WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).
		Printfln("Task:  %s\nAgent: %s", pterm.Cyan(task.Name), pterm.Yellow(agentAddress))

	conn, err := grpc.Dial(agentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ CONNECTION FAILED").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Agent: %s\nTask:  %s\n\nError: %v\n\n"+
					"💡 Troubleshooting:\n"+
					"  • Check agent status: systemctl status sloth-runner-agent\n"+
					"  • Verify agent address is correct\n"+
					"  • Check network: ping <agent-host>\n"+
					"  • Verify firewall rules",
				pterm.Yellow(agentAddress),
				pterm.Cyan(task.Name),
				err,
			)
		pterm.Println()

		slog.Error("Failed to connect to agent",
			"agent_address", agentAddress,
			"task", task.Name,
			"error", err)
		return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("failed to connect to agent %s: %w", agentAddress, err)}
	}
	defer conn.Close()
	c := pb.NewAgentClient(conn)

	// Create a tarball of the workspace
	var buf bytes.Buffer
	if err := ae.CreateTar(session.Workdir, &buf); err != nil {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ WORKSPACE ERROR").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Task:      %s\nWorkspace: %s\n\nError: %v",
				pterm.Cyan(task.Name),
				pterm.Gray(session.Workdir),
				err,
			)
		pterm.Println()

		slog.Error("Failed to create workspace tarball",
			"task", task.Name,
			"error", err)
		return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("failed to create workspace tarball: %w", err)}
	}

	// Generate a script compatible with agent execution (without delegate_to)
	agentScript := ae.GenerateAgentScript(task, groupName)

	pterm.Info.Printfln("📤 Sending task to agent...")

	// Send the task and workspace to the agent
	r, err := c.ExecuteTask(ctx, &pb.ExecuteTaskRequest{
		TaskName:  task.Name,
		TaskGroup: groupName,
		LuaScript: agentScript,
		Workspace: buf.Bytes(),
		User:      task.User,
	})
	if err != nil {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ EXECUTION FAILED").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Agent: %s\nTask:  %s\n\nError: %v",
				pterm.Yellow(agentAddress),
				pterm.Cyan(task.Name),
				err,
			)
		pterm.Println()

		slog.Error("Agent execution failed",
			"agent_address", agentAddress,
			"task", task.Name,
			"error", err)
		return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("agent execution failed for %s: %w", agentAddress, err)}
	}

	// Check for execution errors
	if r.GetError() != "" {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ TASK FAILED ON AGENT").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Agent: %s\nTask:  %s\n\n%s",
				pterm.Yellow(agentAddress),
				pterm.Cyan(task.Name),
				pterm.Red(r.GetError()),
			)
		pterm.Println()

		slog.Error("Task failed on agent",
			"agent_address", agentAddress,
			"task", task.Name,
			"error", r.GetError())
		return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("task execution failed on agent %s: %s", agentAddress, r.GetError())}
	}

	// Success display
	pterm.Println()
	pterm.DefaultBox.
		WithTitle("✅ TASK COMPLETED").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
		Printfln("Agent: %s\nTask:  %s", pterm.Yellow(agentAddress), pterm.Cyan(task.Name))
	if r.GetOutput() != "" {
		pterm.Printf("\n%s\n%s\n", pterm.Gray("Output:"), pterm.Gray(r.GetOutput()))
	}
	pterm.Println()

	slog.Info("Task completed successfully on agent",
		"agent_address", agentAddress,
		"task", task.Name)

	// Extract the updated workspace
	if err := ae.ExtractTar(bytes.NewReader(r.GetWorkspace()), session.Workdir); err != nil {
		pterm.Println()
		pterm.DefaultBox.
			WithTitle("❌ WORKSPACE EXTRACTION FAILED").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Printfln(
				"Agent:     %s\nTask:      %s\nWorkspace: %s\n\nError: %v",
				pterm.Yellow(agentAddress),
				pterm.Cyan(task.Name),
				pterm.Gray(session.Workdir),
				err,
			)
		pterm.Println()

		slog.Error("Failed to extract updated workspace from agent",
			"agent_address", agentAddress,
			"task", task.Name,
			"error", err)
		return &TaskExecutionError{TaskName: task.Name, Err: fmt.Errorf("failed to extract updated workspace from agent %s: %w", agentAddress, err)}
	}

	pterm.Info.Printfln("📥 Workspace synchronized")
	return nil
}
