# Arquitetura Modular do Sloth Runner

## Visão Geral

O Sloth Runner foi refatorado para seguir princípios de design modular e aplicar design patterns reconhecidos pela indústria. Isso torna o código mais fácil de manter, testar e estender.

## Estrutura de Diretórios

```
cmd/sloth-runner/
├── commands/           # Comandos Cobra organizados
│   ├── context.go     # AppContext (Dependency Injection)
│   ├── root.go        # Root command
│   ├── version.go     # Version command
│   ├── run.go         # Run command
│   ├── agent/         # Comandos do agente
│   ├── stack/         # Comandos de stack
│   └── scheduler/     # Comandos de scheduler
├── handlers/          # Business logic handlers
│   └── run_handler.go # Handler para comando run
├── services/          # Serviços reutilizáveis
│   └── stack_service.go # Serviço de gerenciamento de stack
├── repositories/      # Camada de acesso a dados (futuro)
└── main.go           # Entry point (mínimo)
```

## Design Patterns Aplicados

### 1. Dependency Injection Pattern

**Arquivo**: `commands/context.go`

O `AppContext` encapsula todas as dependências compartilhadas entre comandos:

```go
type AppContext struct {
    Version       string
    Commit        string
    Date          string
    AgentRegistry interface{}
    SurveyAsker   taskrunner.SurveyAsker
    OutputWriter  io.Writer
    TestMode      bool
    // ...
}
```

**Benefícios:**
- ✅ Facilita testes (mock de dependências)
- ✅ Reduz acoplamento
- ✅ Torna dependências explícitas
- ✅ Permite configuração centralizada

### 2. Factory Pattern

**Arquivo**: `commands/*.go`

Cada comando é criado através de uma factory function:

```go
func NewRunCommand(ctx *AppContext) *cobra.Command {
    // Cria e configura o comando
    return cmd
}
```

**Benefícios:**
- ✅ Encapsula lógica de criação
- ✅ Permite configuração consistente
- ✅ Facilita testes unitários
- ✅ Suporta diferentes configurações

### 3. Handler Pattern (Command Handler)

**Arquivo**: `handlers/run_handler.go`

Separa a lógica de negócio do framework Cobra:

```go
type RunHandler struct {
    stackService *services.StackService
    config       *RunConfig
}

func (h *RunHandler) Execute() error {
    // Lógica de negócio aqui
}
```

**Benefícios:**
- ✅ Lógica de negócio independente do framework
- ✅ Facilita testes unitários
- ✅ Código reutilizável
- ✅ Single Responsibility Principle

### 4. Service Layer Pattern

**Arquivo**: `services/stack_service.go`

Encapsula operações de negócio relacionadas:

```go
type StackService struct {
    manager *stack.StackManager
}

func (s *StackService) GetOrCreateStack(...) (string, error) {
    // Lógica de serviço
}
```

**Benefícios:**
- ✅ Reutilização de lógica
- ✅ Transações e coordenação
- ✅ Abstração da camada de dados
- ✅ Testabilidade

### 5. Strategy Pattern (Futuro)

**Planejado para**: Executores (Local, SSH, Agent)

```go
type Executor interface {
    Execute(task Task) error
}

type LocalExecutor struct {}
type SSHExecutor struct {}
type AgentExecutor struct {}
```

**Benefícios:**
- ✅ Diferentes estratégias de execução
- ✅ Fácil adicionar novos executores
- ✅ Open/Closed Principle

### 6. Repository Pattern (Futuro)

**Planejado para**: Acesso a dados (DB, API, etc)

```go
type StackRepository interface {
    Get(id string) (*Stack, error)
    Create(stack *Stack) error
    Update(stack *Stack) error
    Delete(id string) error
}
```

**Benefícios:**
- ✅ Abstração de persistência
- ✅ Facilita testes com mocks
- ✅ Troca de backend transparente

## Fluxo de Execução

### Antes (Monolítico)

```
main.go (3462 linhas)
  └─> runCmd (500+ linhas)
      └─> Toda lógica inline
```

### Depois (Modular)

```
main.go (mínimo)
  └─> NewRootCommand(ctx)
      └─> NewRunCommand(ctx)
          └─> RunHandler.Execute()
              ├─> StackService.GetOrCreateStack()
              ├─> RunHandler.initializeSSH()
              ├─> RunHandler.parseLuaScript()
              ├─> RunHandler.executeTasks()
              └─> StackService.RecordExecution()
```

## Princípios SOLID Aplicados

### Single Responsibility Principle (SRP)
- ✅ Cada classe/struct tem uma única responsabilidade
- ✅ `RunHandler` gerencia execução de tasks
- ✅ `StackService` gerencia stacks
- ✅ `RunCommand` gerencia apenas CLI

### Open/Closed Principle (OCP)
- ✅ Aberto para extensão via interfaces
- ✅ Fechado para modificação (comportamento base)
- ✅ Novos executores podem ser adicionados sem alterar código existente

### Liskov Substitution Principle (LSP)
- ✅ Interfaces podem ser substituídas por implementações
- ✅ Mocks podem substituir serviços reais

### Interface Segregation Principle (ISP)
- ✅ Interfaces pequenas e específicas
- ✅ Clientes não dependem de métodos que não usam

### Dependency Inversion Principle (DIP)
- ✅ Dependência de abstrações, não de implementações
- ✅ AppContext injeta dependências
- ✅ Handlers dependem de interfaces de serviços

## Benefícios da Refatoração

### 1. Testabilidade
- ✅ Handlers podem ser testados sem Cobra
- ✅ Serviços podem ser mockados
- ✅ Lógica de negócio isolada

### 2. Manutenibilidade
- ✅ Código organizado por responsabilidade
- ✅ Arquivos menores e focados
- ✅ Fácil localizar e modificar código

### 3. Extensibilidade
- ✅ Novos comandos seguem mesmo padrão
- ✅ Novos serviços são fáceis de adicionar
- ✅ Novos executores via Strategy Pattern

### 4. Legibilidade
- ✅ Estrutura clara e previsível
- ✅ Nomes descritivos
- ✅ Separação de concerns

## Exemplos de Uso

### Criando um Novo Comando

```go
// 1. Criar comando em commands/
func NewMyCommand(ctx *AppContext) *cobra.Command {
    return &cobra.Command{
        Use: "my-command",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Extrair flags
            // Criar serviços necessários
            // Criar handler
            // Executar handler
            return handler.Execute()
        },
    }
}

// 2. Criar handler em handlers/
type MyHandler struct {
    service *services.MyService
    config  *MyConfig
}

func (h *MyHandler) Execute() error {
    // Lógica de negócio
}

// 3. Criar serviço se necessário em services/
type MyService struct {
    // dependências
}
```

### Testando um Handler

```go
func TestRunHandler_Execute(t *testing.T) {
    // Arrange
    mockStackService := &MockStackService{}
    config := &handlers.RunConfig{
        StackName: "test",
        FilePath:  "test.sloth",
        // ...
    }
    handler := handlers.NewRunHandler(mockStackService, config)

    // Act
    err := handler.Execute()

    // Assert
    assert.NoError(t, err)
    assert.True(t, mockStackService.CreateStackCalled)
}
```

## Próximos Passos

### Curto Prazo
1. ✅ Extrair comando `run`
2. ⏳ Extrair comandos `agent/*`
3. ⏳ Extrair comandos `stack/*`
4. ⏳ Extrair comandos `scheduler/*`

### Médio Prazo
1. ⏳ Implementar Strategy Pattern para executores
2. ⏳ Implementar Repository Pattern
3. ⏳ Adicionar testes unitários
4. ⏳ Adicionar testes de integração

### Longo Prazo
1. ⏳ Métricas e observabilidade
2. ⏳ Plugin system
3. ⏳ API REST/GraphQL
4. ⏳ Web UI

## Referências

- **Clean Architecture** - Robert C. Martin
- **Domain-Driven Design** - Eric Evans
- **Enterprise Integration Patterns** - Gregor Hohpe
- **Refactoring** - Martin Fowler
- **Design Patterns** - Gang of Four

## Conclusão

A arquitetura modular aplicada ao Sloth Runner segue best practices da indústria e torna o código mais profissional, mantível e extensível. Cada padrão foi escolhido para resolver problemas específicos e melhorar a qualidade geral do código.
