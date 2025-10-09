package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/db"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/events"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/group"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/hook"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/scheduler"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/secrets"
	slothcmd "github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/sloth"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/stack"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/state"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/workflow"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	coremodules "github.com/chalkan3-sloth/sloth-runner/internal/modules/core"
	"github.com/pterm/pterm"
)

var (
	// Build variables (set via ldflags at build time)
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Set up structured logging with pterm
	slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))

	// Execute the CLI
	if err := Execute(); err != nil {
		// Print formatted errors
		if strings.Contains(err.Error(), "✗") {
			// Already formatted, just print it
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			// Log using slog for unformatted errors
			slog.Error("execution failed", "err", err)
		}
		os.Exit(1)
	}
}

// InitializeHookSystem initializes the global hook dispatcher
// This should be called only by commands that need hook support
func InitializeHookSystem() error {
	if err := hooks.InitializeGlobalDispatcher(); err != nil {
		return err
	}

	// Wire up the dispatcher to the event module
	dispatcher := hooks.GetGlobalDispatcher()
	if dispatcher != nil {
		dispatcherFunc := dispatcher.CreateEventDispatcherFunc()
		coremodules.SetGlobalEventDispatcher(dispatcherFunc)
	}

	return nil
}

func Execute() error {
	// Initialize master server starter function
	commands.MasterServerStarter = func(port int) error {
		server := newAgentRegistryServer()
		return server.Start(port)
	}

	// Create application context with build info
	ctx := commands.NewAppContext(version, commit, date)

	// Create root command
	rootCmd := commands.NewRootCommand(ctx)
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	// Add version command
	versionCmd := commands.NewVersionCommand(ctx)
	rootCmd.AddCommand(versionCmd)

	// Add run command
	runCmd := commands.NewRunCommand(ctx)
	rootCmd.AddCommand(runCmd)

	// Add agent command and subcommands
	agentCmd := agent.NewAgentCommand(ctx)
	rootCmd.AddCommand(agentCmd)

	// Add group command and subcommands
	groupCmd := group.NewGroupCmd()
	rootCmd.AddCommand(groupCmd)

	// Add stack command and subcommands
	stackCmd := stack.NewStackCommand(ctx)
	rootCmd.AddCommand(stackCmd)

	// Add sloth command and subcommands
	slothCmd := slothcmd.NewSlothCommand(ctx)
	rootCmd.AddCommand(slothCmd)

	// Add scheduler command and subcommands
	schedulerCmd := scheduler.NewSchedulerCommand(ctx)
	rootCmd.AddCommand(schedulerCmd)

	// Add state command and subcommands
	stateCmd := state.NewStateCommand(ctx)
	rootCmd.AddCommand(stateCmd)

	// Add workflow command and subcommands
	workflowCmd := workflow.NewWorkflowCommand(ctx)
	rootCmd.AddCommand(workflowCmd)

	// Add secrets command and subcommands
	secretsCmd := secrets.NewSecretsCommand(ctx)
	rootCmd.AddCommand(secretsCmd)

	// Add hook command and subcommands
	hookCmd := hook.NewHookCommand(ctx)
	rootCmd.AddCommand(hookCmd)

	// Add events command and subcommands
	eventsCmd := events.NewEventsCommand(ctx)
	rootCmd.AddCommand(eventsCmd)

	// Add db command and subcommands
	dbCmd := db.NewDBCommand(ctx)
	rootCmd.AddCommand(dbCmd)

	// Add other root commands (kept for backward compatibility)
	listCmd := commands.NewListCommand(ctx)
	rootCmd.AddCommand(listCmd)

	masterCmd := commands.NewMasterCommand(ctx)
	rootCmd.AddCommand(masterCmd)

	uiCmd := commands.NewUICommand(ctx)
	rootCmd.AddCommand(uiCmd)

	// Add modules command (list available Lua modules)
	rootCmd.AddCommand(modulesCmd)

	// Execute root command
	return rootCmd.Execute()
}
