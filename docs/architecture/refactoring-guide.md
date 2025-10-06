# Guia de Refatoração - Sloth Runner

## 📊 Status Atual da Refatoração

### ✅ Concluído

1. **Estrutura Base Modular**
   - ✅ `commands/context.go` - Dependency Injection
   - ✅ `commands/root.go` - Root command
   - ✅ `commands/version.go` - Version command
   - ✅ `commands/run.go` - Run command refatorado
   - ✅ `handlers/run_handler.go` - Handler do run
   - ✅ `services/stack_service.go` - Serviço de stack
   - ✅ `services/agent_service.go` - Serviço de agent

2. **Comandos Agent (Parcial)**
   - ✅ `commands/agent/agent.go` - Parent command
   - ✅ `commands/agent/list.go` - List agents
   - ✅ `commands/agent/stop.go` - Stop agent
   - ✅ `commands/agent/delete.go` - Delete agent
   - ⏳ get, start, run, modules, metrics (pendentes)

### ⏳ Pendente

| Arquivo Original | Linhas | Comando | Status |
|-----------------|--------|---------|--------|
| `main.go` | 3462 | agent start | ⏳ |
| `main.go` | 3462 | agent get | ⏳ |
| `main.go` | 3462 | agent run | ⏳ |
| `main.go` | 3462 | agent modules | ⏳ |
| `main.go` | 3462 | agent metrics | ⏳ |
| `main.go` | 3462 | stack * | ⏳ |
| `main.go` | 3462 | scheduler * | ⏳ |
| `main.go` | 3462 | state * | ⏳ |
| `main.go` | 3462 | ssh * | ⏳ |
| `luainterface.go` | 1793 | Módulos Lua | ⏳ |
| `taskrunner.go` | 1573 | Task execution | ⏳ |
| `user.go` | 1669 | User module | ⏳ |

## 🎯 Estratégia de Refatoração

### Fase 1: Comandos CLI (Prioridade Alta)

**Objetivo**: Extrair todos os comandos do `main.go`

```
main.go (3462 linhas)
├── commands/run.go (✅ Feito)
├── commands/list.go
├── commands/ui.go
├── commands/master.go
├── commands/agent/
│   ├── start.go (⏳)
│   ├── stop.go (✅)
│   ├── list.go (✅)
│   ├── delete.go (✅)
│   ├── get.go (⏳)
│   ├── run.go (⏳)
│   ├── modules.go (⏳)
│   ├── metrics.go (⏳)
│   └── update.go (existe em agent_update.go)
├── commands/stack/
│   ├── new.go
│   ├── list.go
│   ├── show.go
│   ├── delete.go
│   └── history.go
├── commands/scheduler/
│   ├── enable.go
│   ├── disable.go
│   ├── list.go
│   └── delete.go
└── commands/state/
    ├── list.go
    ├── get.go
    ├── set.go
    └── delete.go
```

### Fase 2: Módulos Lua (Prioridade Média)

**Objetivo**: Modularizar `internal/luainterface/luainterface.go` (1793 linhas)

```
internal/luainterface/
├── luainterface.go (core, ~300 linhas)
├── modules/
│   ├── task.go (task, task_group)
│   ├── pipeline.go (pipeline functions)
│   ├── template.go (template functions)
│   ├── env.go (environment functions)
│   ├── import.go (import functionality)
│   └── validation.go (validation functions)
```

### Fase 3: Task Runner (Prioridade Média)

**Objetivo**: Modularizar `internal/taskrunner/taskrunner.go` (1573 linhas)

```
internal/taskrunner/
├── taskrunner.go (core, ~300 linhas)
├── executor/
│   ├── local.go (execução local)
│   ├── remote.go (execução remota)
│   ├── ssh.go (execução SSH)
│   └── agent.go (execução via agent)
├── output/
│   ├── formatter.go (formatação de saída)
│   ├── json.go (JSON output)
│   └── enhanced.go (enhanced output)
└── result/
    ├── collector.go (coleta de resultados)
    └── aggregator.go (agregação)
```

### Fase 4: Módulos de Usuário (Prioridade Baixa)

**Objetivo**: Modularizar arquivos grandes em `internal/luainterface/`

```
internal/luainterface/
├── user.go (1669 linhas) → user/
│   ├── user.go (~300 linhas)
│   ├── group.go
│   ├── sudo.go
│   └── validation.go
├── modern_dsl.go (1619 linhas) → dsl/
│   ├── parser.go
│   ├── executor.go
│   └── validator.go
└── terraform_advanced.go (1511 linhas) → terraform/
    ├── core.go
    ├── advanced.go
    └── helpers.go
```

## 📋 Templates para Refatoração

### Template 1: Comando Simples

```go
// commands/agent/COMMAND.go
package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/spf13/cobra"
)

func NewCOMMANDCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "COMMAND",
		Short: "Short description",
		Long:  `Long description`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Extract flags
			flag1, _ := cmd.Flags().GetString("flag1")

			// 2. Create service
			service := services.NewAgentService("")

			// 3. Execute operation
			result, err := service.DoSomething(flag1)
			if err != nil {
				return err
			}

			// 4. Display result
			fmt.Println(result)
			return nil
		},
	}
}
```

### Template 2: Comando com Handler

```go
// commands/agent/COMMAND.go
package agent

import (
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/handlers"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	"github.com/spf13/cobra"
)

func NewCOMMANDCommand(ctx *commands.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "COMMAND",
		Short: "Short description",
		Long:  `Long description`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract flags
			flag1, _ := cmd.Flags().GetString("flag1")

			// Create service
			service, err := services.NewCOMMANDService()
			if err != nil {
				return err
			}
			defer service.Close()

			// Create configuration
			config := &handlers.COMMANDConfig{
				Flag1: flag1,
			}

			// Create and execute handler
			handler := handlers.NewCOMMANDHandler(service, config)
			return handler.Execute()
		},
	}
}

// handlers/COMMAND_handler.go
package handlers

type COMMANDConfig struct {
	Flag1 string
}

type COMMANDHandler struct {
	service *services.COMMANDService
	config  *COMMANDConfig
}

func NewCOMMANDHandler(service *services.COMMANDService, config *COMMANDConfig) *COMMANDHandler {
	return &COMMANDHandler{
		service: service,
		config:  config,
	}
}

func (h *COMMANDHandler) Execute() error {
	// Business logic here
	return nil
}
```

### Template 3: Serviço

```go
// services/DOMAIN_service.go
package services

type DOMAINService struct {
	// dependencies
}

func NewDOMAINService() (*DOMAINService, error) {
	return &DOMAINService{}, nil
}

func (s *DOMAINService) Close() error {
	return nil
}

func (s *DOMAINService) DoSomething(param string) (result interface{}, err error) {
	// Service logic
	return nil, nil
}
```

## 🔄 Processo de Refatoração Passo a Passo

### Para Cada Comando

1. **Identificar o comando no main.go**
   ```bash
   grep -n "var COMMANDCmd = &cobra.Command" cmd/sloth-runner/main.go
   ```

2. **Extrair lógica**
   - Copiar código do RunE
   - Identificar dependências (gRPC, DB, etc)
   - Identificar flags

3. **Criar serviço (se necessário)**
   - Encapsular lógica de negócio reutilizável
   - Abstrair comunicação gRPC/DB
   - Adicionar testes

4. **Criar handler (se complexo)**
   - Separar lógica de negócio do CLI
   - Facilitar testes unitários
   - Melhorar legibilidade

5. **Criar comando**
   - Usar template apropriado
   - Configurar flags
   - Chamar serviço/handler

6. **Testar**
   ```bash
   go build -o sloth-runner-test ./cmd/sloth-runner
   ./sloth-runner-test COMMAND --help
   ./sloth-runner-test COMMAND [args]
   ```

7. **Commit**
   ```bash
   git add .
   git commit -m "refactor: Extract COMMAND to modular structure"
   ```

## 📈 Métricas de Sucesso

### Antes
- ❌ main.go: 3.462 linhas
- ❌ Arquivos > 1.500 linhas: 6 arquivos
- ❌ Cobertura de testes: ~20%
- ❌ Acoplamento: Alto

### Meta Final
- ✅ main.go: < 100 linhas
- ✅ Arquivos > 500 linhas: 0 arquivos
- ✅ Cobertura de testes: > 70%
- ✅ Acoplamento: Baixo
- ✅ Design patterns: 5+
- ✅ SOLID: 100%

## 🚀 Próximos Passos Imediatos

1. ✅ **Comandos Agent Simples** (list, stop, delete) - FEITO
2. ⏳ **Comandos Agent Complexos** (start, get, run, modules, metrics)
3. ⏳ **Comandos Stack** (new, list, show, delete, history)
4. ⏳ **Comandos Scheduler** (enable, disable, list, delete)
5. ⏳ **Comandos State** (list, get, set, delete)
6. ⏳ **Refatorar luainterface.go**
7. ⏳ **Refatorar taskrunner.go**
8. ⏳ **Adicionar testes unitários**

## 📚 Recursos

- [Arquitetura Modular](./README.md)
- [Design Patterns](./modular-design.md)
- [Exemplo Main.go](../../cmd/sloth-runner/main_modular_example.go)

## 💡 Dicas

1. **Comece pelos comandos simples** - list, show, version
2. **Crie serviços para lógica reutilizável** - gRPC, DB
3. **Use handlers para lógica complexa** - múltiplas etapas
4. **Mantenha arquivos < 300 linhas** - fácil de manter
5. **Teste cada comando após extrair** - evita regressions
6. **Commit frequentemente** - facilita rollback
7. **Documente decisões** - ajuda futuros desenvolvedores

---

**Autor**: Claude Code
**Data**: 2025-10-06
**Status**: Em Andamento - 15% Completo
