# Testing Guide - Sloth Runner

## Overview

Este guia documenta a estratégia de testes e o trabalho realizado para alcançar cobertura de testes no Sloth Runner.

## Status Atual

```
Current Test Coverage: 9.0%
Target Coverage: 80%
Tests Created: 8 new test files
Tests Passing: ✅ All tests passing
```

## Documentação

Este projeto de testes inclui:

1. **TEST_COVERAGE_REPORT.md** - Relatório detalhado do status atual e roadmap para 80%
2. **TEST_EXAMPLES.md** - Exemplos práticos de como testar diferentes partes do código
3. **TESTING_README.md** (este arquivo) - Guia rápido de uso

## Quick Start

### Executar todos os testes

```bash
go test ./...
```

### Executar testes com coverage

```bash
go test ./... -coverprofile=coverage.out -coverpkg=./...
go tool cover -func=coverage.out | tail -1
```

### Ver relatório HTML

```bash
go tool cover -html=coverage.out
```

### Usar o script helper

```bash
# Ver ajuda
./scripts/test-coverage.sh help

# Executar testes com coverage
./scripts/test-coverage.sh coverage

# Ver relatório detalhado
./scripts/test-coverage.sh report

# Abrir relatório HTML
./scripts/test-coverage.sh html

# Verificar threshold
./scripts/test-coverage.sh threshold 80

# Encontrar arquivos sem testes
./scripts/test-coverage.sh untested
```

## Módulos com Testes Completos

### ✅ 100% Coverage

- **internal/types** - Tipos e estruturas base
- **internal/runner** - Executor simples de tasks

### ✅ Boa Coverage (>60%)

- **internal/agent** (69%) - Health monitoring e métricas
- **internal/scheduler** (82.7%) - Agendamento de tasks
- **internal/state** (61.1%) - Gerenciamento de estado

### 🟡 Coverage Moderado (40-60%)

- **internal/core** (42.8%) - Primitivas de concorrência e core system
- **internal/reliability** (42.3%) - Circuit breakers e retries
- **internal/modules** (41.3%) - Sistema de módulos

### 🔴 Necessita Testes (<40%)

- **internal/taskrunner** (30.1%) - Executor principal de workflows
- **internal/luainterface** (1.1%) - Interface Lua (muito grande)
- **cmd/sloth-runner** (11.3%) - CLI commands
- Vários módulos sem testes (0%)

## Arquivos de Teste Criados

### 1. internal/types/types_test.go
```
✅ TestGenerateTaskID
✅ TestGenerateTaskGroupID  
✅ TestGenerateTaskID_Uniqueness
✅ TestGenerateTaskGroupID_Uniqueness
✅ TestTask_Struct
✅ TestTaskGroup_Struct
✅ TestTaskResult_Struct
✅ TestPythonVenv_Struct
```

### 2. internal/runner/simple_runner_test.go
```
✅ TestRunSingleTask_Success
✅ TestRunSingleTask_Failure
✅ TestRunSingleTask_WithOutput
✅ TestRunSingleTask_NilCommandFunc
✅ TestRunSingleTask_WithParams
✅ TestRunSingleTask_TimingAccuracy
```

### 3. internal/agent/health_test.go
```
✅ TestNewHealthMonitor
✅ TestHealthMonitor_SetThresholds
✅ TestHealthMonitor_AddCollector
✅ TestHealthMonitor_CollectMetrics
✅ TestHealthMonitor_CollectSystemMetrics
✅ TestHealthMonitor_CollectMemoryMetrics
✅ TestHealthMonitor_CollectCPUMetrics
✅ TestHealthMonitor_CollectDiskMetrics
✅ TestHealthMonitor_DetermineHealthStatus_*
✅ TestHealthMonitor_StartPeriodicCollection
✅ TestTaskMetricsCollector_*
✅ TestHealthMonitor_ConcurrentAccess
```

### 4. internal/agent/metrics_test.go
```
✅ TestNewMetricsCollector
✅ TestMetricsCollector_CollectSystemMetrics
✅ TestMetricsCollector_CollectRuntimeMetrics
✅ TestMetricsCollector_GetSnapshot
✅ TestMetricsCollector_UpdateTaskMetrics
✅ TestMetricsCollector_CustomMetrics
✅ TestMetricsCollector_HTTPHandler
✅ TestSystemMetrics_JSONMarshaling
✅ TestRuntimeMetrics_JSONMarshaling
✅ TestTaskMetrics_JSONMarshaling
✅ TestMetricsCollector_Enable_Disable
✅ TestMetricsCollector_ConcurrentAccess
✅ TestMetricsSnapshot_CompleteData
✅ BenchmarkMetricsCollector_*
```

### 5. internal/core/core_test.go (enhanced)
```
✅ TestGlobalCore
✅ TestCircuitBreaker
✅ TestResourcePool
✅ TestWorkerPool
✅ TestErrorHandling
✅ TestSafeMap
✅ TestSemaphore (new)
✅ TestRWCounter (new)
✅ TestRateLimiter (new)
✅ TestSafeMapForEach (new)
✅ TestWorkerPoolSubmitWithTimeout (new)
✅ TestSemaphoreWithTimeout (new)
✅ TestDefaultCoreConfig (new)
✅ TestGlobalCoreMetrics (new)
✅ TestConcurrentSafeMap (new)
✅ BenchmarkWorkerPool
✅ BenchmarkSafeMap
```

## Comandos Úteis

### Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com verbose
go test -v ./...

# Executar testes de um pacote específico
go test -v ./internal/core

# Executar teste específico
go test -v -run TestGlobalCore ./internal/core

# Executar testes com timeout
go test -timeout 30s ./...
```

### Coverage

```bash
# Gerar coverage
go test ./... -coverprofile=coverage.out -coverpkg=./...

# Ver summary
go tool cover -func=coverage.out | tail -1

# Ver por função
go tool cover -func=coverage.out

# Gerar HTML
go tool cover -html=coverage.out -o coverage.html

# Coverage de pacote específico
go test -cover ./internal/core
```

### Race Detection

```bash
# Detectar race conditions
go test -race ./...

# Com verbose
go test -race -v ./...
```

### Benchmarks

```bash
# Executar todos os benchmarks
go test -bench=. ./...

# Com análise de memória
go test -bench=. -benchmem ./...

# Benchmark específico
go test -bench=BenchmarkSafeMap ./internal/core
```

### Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=. ./internal/core

# Memory profile
go test -memprofile=mem.prof -bench=. ./internal/core

# Analisar profile
go tool pprof cpu.prof
```

## Padrões de Teste

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    // Setup
    setup := createTestSetup()
    defer cleanup()
    
    // Execute
    result := functionUnderTest(input)
    
    // Assert
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {"valid", "test", true, false},
        {"invalid", "", false, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr %v, err %v", tt.wantErr, err)
            }
            if got != tt.want {
                t.Errorf("want %v, got %v", tt.want, got)
            }
        })
    }
}
```

### Subtests

```go
func TestModule(t *testing.T) {
    t.Run("initialization", func(t *testing.T) {
        // test init
    })
    
    t.Run("execution", func(t *testing.T) {
        // test exec
    })
}
```

### Parallel Tests

```go
func TestParallel(t *testing.T) {
    t.Parallel()
    
    tests := []struct{
        name string
        // ...
    }{
        // ...
    }
    
    for _, tt := range tests {
        tt := tt // capture
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // test code
        })
    }
}
```

## Mocking Strategies

### Interface-based Mocking

```go
type CommandExecutor interface {
    Execute(cmd string) (string, error)
}

type MockExecutor struct {
    output string
    err    error
}

func (m *MockExecutor) Execute(cmd string) (string, error) {
    return m.output, m.err
}
```

### Using testify/mock

```go
import "github.com/stretchr/testify/mock"

type MockService struct {
    mock.Mock
}

func (m *MockService) DoSomething(arg string) error {
    args := m.Called(arg)
    return args.Error(0)
}

// In test:
mockService := new(MockService)
mockService.On("DoSomething", "test").Return(nil)
```

## Test Helpers

### Temporary Directories

```go
func TestWithTempDir(t *testing.T) {
    tmpDir := t.TempDir() // Auto-cleanup
    // Use tmpDir
}
```

### Setup/Teardown

```go
func setupTest(t *testing.T) (*TestEnv, func()) {
    env := &TestEnv{
        // setup
    }
    
    cleanup := func() {
        // teardown
    }
    
    return env, cleanup
}

func TestWithSetup(t *testing.T) {
    env, cleanup := setupTest(t)
    defer cleanup()
    
    // Use env
}
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $coverage%"
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "Coverage below 80%"
            exit 1
          fi
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

## Próximos Passos

Para alcançar 80% de coverage, siga o roadmap em **TEST_COVERAGE_REPORT.md**:

1. **Fase 1** - Módulos pequenos (output, stack, scaffolding) - +20%
2. **Fase 2** - Core functionality (taskrunner, core) - +30%
3. **Fase 3** - Lua interface & modules - +15%
4. **Fase 4** - GitOps & features avançadas - +6%

Consulte **TEST_EXAMPLES.md** para exemplos práticos de como testar cada tipo de módulo.

## Recursos Adicionais

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Comments](https://golang.org/cmd/go/#hdr-Testing_flags)
- [Testify Documentation](https://github.com/stretchr/testify)

## Suporte

Para dúvidas sobre os testes:

1. Consulte **TEST_EXAMPLES.md** para exemplos práticos
2. Veja **TEST_COVERAGE_REPORT.md** para estratégia geral
3. Use `./scripts/test-coverage.sh help` para comandos úteis

---

**Status**: Base de testes estabelecida (9.0% coverage)
**Próximo Objetivo**: 29% (Fase 1 completa)
**Objetivo Final**: 80% coverage
