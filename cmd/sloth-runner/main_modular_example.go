package main

// Este é um EXEMPLO de como o main.go ficaria após a refatoração completa
// Demonstra a simplicidade e clareza da nova arquitetura modular

/*
import (
	"fmt"
	"os"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
)

// Version information (set by build flags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// 1. Create application context with dependencies
	ctx := commands.NewAppContext(version, commit, date)

	// 2. Initialize global services if needed
	// ctx.AgentRegistry = initializeAgentRegistry()

	// 3. Create root command with all subcommands
	rootCmd := commands.NewRootCommand(ctx)

	// 4. Add all subcommands using factory functions
	rootCmd.AddCommand(
		commands.NewVersionCommand(ctx),
		commands.NewRunCommand(ctx),
		commands.NewListCommand(ctx),
		commands.NewUICommand(ctx),
		commands.NewMasterCommand(ctx),
		// Agent commands
		commands.NewAgentCommand(ctx),
		// Stack commands
		commands.NewStackCommand(ctx),
		// Scheduler commands
		commands.NewSchedulerCommand(ctx),
		// State commands
		commands.NewStateCommand(ctx),
		// SSH commands
		commands.NewSSHCommand(ctx),
		// Modules commands
		commands.NewModulesCommand(ctx),
	)

	// 5. Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// COMPARAÇÃO: Antes vs Depois

// ===== ANTES (main.go com 3462 linhas) =====
// - Toda lógica inline
// - 37 comandos definidos no mesmo arquivo
// - Lógica de negócio misturada com CLI
// - Difícil de testar
// - Difícil de manter
// - Acoplamento alto

// ===== DEPOIS (main.go com ~40 linhas) =====
// - Apenas entry point
// - Comandos em arquivos separados
// - Lógica de negócio em handlers
// - Fácil de testar
// - Fácil de manter
// - Baixo acoplamento
// - Design patterns aplicados:
//   * Dependency Injection (AppContext)
//   * Factory Pattern (NewXXXCommand)
//   * Handler Pattern (handlers/)
//   * Service Layer (services/)

// ESTRUTURA DE DIRETÓRIOS APÓS REFATORAÇÃO:
//
// cmd/sloth-runner/
// ├── main.go (40 linhas - apenas entry point)
// ├── commands/
// │   ├── context.go           (Dependency Injection)
// │   ├── root.go              (Root command)
// │   ├── version.go           (Version command)
// │   ├── run.go               (Run command)
// │   ├── list.go              (List command)
// │   ├── ui.go                (UI command)
// │   ├── master.go            (Master command)
// │   ├── agent/
// │   │   ├── agent.go         (Agent parent command)
// │   │   ├── start.go         (Agent start)
// │   │   ├── stop.go          (Agent stop)
// │   │   ├── list.go          (Agent list)
// │   │   ├── delete.go        (Agent delete)
// │   │   ├── get.go           (Agent get)
// │   │   ├── run.go           (Agent run)
// │   │   ├── modules.go       (Agent modules)
// │   │   ├── metrics.go       (Agent metrics)
// │   │   └── update.go        (Agent update)
// │   ├── stack/
// │   │   ├── stack.go         (Stack parent command)
// │   │   ├── new.go           (Stack new)
// │   │   ├── list.go          (Stack list)
// │   │   ├── show.go          (Stack show)
// │   │   ├── delete.go        (Stack delete)
// │   │   └── history.go       (Stack history)
// │   ├── scheduler/
// │   │   ├── scheduler.go     (Scheduler parent command)
// │   │   ├── enable.go        (Scheduler enable)
// │   │   ├── disable.go       (Scheduler disable)
// │   │   ├── list.go          (Scheduler list)
// │   │   └── delete.go        (Scheduler delete)
// │   ├── state/
// │   │   ├── state.go         (State parent command)
// │   │   ├── list.go          (State list)
// │   │   ├── get.go           (State get)
// │   │   ├── set.go           (State set)
// │   │   └── delete.go        (State delete)
// │   ├── ssh/
// │   │   ├── ssh.go           (SSH parent command)
// │   │   ├── add.go           (SSH add)
// │   │   ├── list.go          (SSH list)
// │   │   ├── test.go          (SSH test)
// │   │   ├── exec.go          (SSH exec)
// │   │   └── delete.go        (SSH delete)
// │   └── modules/
// │       ├── modules.go       (Modules parent command)
// │       └── list.go          (Modules list)
// ├── handlers/
// │   ├── run_handler.go       (Business logic para run)
// │   ├── agent_handler.go     (Business logic para agent)
// │   ├── stack_handler.go     (Business logic para stack)
// │   └── ...
// ├── services/
// │   ├── stack_service.go     (Serviço de stack)
// │   ├── agent_service.go     (Serviço de agent)
// │   ├── scheduler_service.go (Serviço de scheduler)
// │   └── ...
// └── repositories/ (futuro)
//     ├── stack_repository.go
//     ├── agent_repository.go
//     └── ...

// EXEMPLO DE TESTE COM A NOVA ARQUITETURA:

func ExampleTestRunCommand() {
	// Arrange
	ctx := commands.NewAppContext("test", "abc123", "2024-01-01")
	ctx.TestMode = true
	ctx.OutputWriter = &bytes.Buffer{}

	// Act
	cmd := commands.NewRunCommand(ctx)
	cmd.SetArgs([]string{"test-stack", "--file", "test.sloth", "--yes"})
	err := cmd.Execute()

	// Assert
	// assert.NoError(t, err)
	// assert.Contains(t, output.String(), "success")
	_ = err
}

// BENEFÍCIOS DA NOVA ARQUITETURA:

// 1. TESTABILIDADE
//    - Handlers podem ser testados sem CLI
//    - Serviços podem ser mockados
//    - Dependências injetadas são fáceis de substituir

// 2. MANUTENIBILIDADE
//    - Cada arquivo tem < 200 linhas
//    - Fácil encontrar código específico
//    - Mudanças isoladas não afetam outros componentes

// 3. EXTENSIBILIDADE
//    - Novos comandos seguem padrão claro
//    - Novos serviços são independentes
//    - Novos handlers são autônomos

// 4. LEGIBILIDADE
//    - Estrutura previsível
//    - Responsabilidades claras
//    - Código auto-documentado

// 5. PROFISSIONALISMO
//    - Segue best practices da indústria
//    - Aplica design patterns reconhecidos
//    - Código enterprise-grade
*/
