# Test Examples and Patterns for Sloth Runner

Este documento contém exemplos práticos de como testar diferentes partes do Sloth Runner.

## Exemplo 1: Testando Módulos Lua

### Testando um módulo simples

```go
package luainterface

import (
	"testing"
	lua "github.com/yuin/gopher-lua"
)

func TestSystemModule_Hostname(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	// Register o módulo
	RegisterSystemModule(L)
	
	// Execute um script Lua que usa a função
	script := `
		local system = require('system')
		local hostname = system.hostname()
		return hostname
	`
	
	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}
	
	// Verifique o resultado
	result := L.Get(-1)
	if result.Type() != lua.LTString {
		t.Errorf("Expected string result, got %v", result.Type())
	}
	
	hostname := result.String()
	if hostname == "" {
		t.Error("Expected non-empty hostname")
	}
}
```

### Testando com Mock de comandos Shell

```go
package luainterface

import (
	"testing"
	lua "github.com/yuin/gopher-lua"
)

// Mock executor for testing
type MockShellExecutor struct {
	output   string
	exitCode int
	err      error
}

func (m *MockShellExecutor) Execute(cmd string, args ...string) (string, int, error) {
	return m.output, m.exitCode, m.err
}

func TestPkgModule_Install_Success(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	// Inject mock executor
	mockExec := &MockShellExecutor{
		output:   "Package installed successfully",
		exitCode: 0,
		err:      nil,
	}
	SetShellExecutor(mockExec) // Você precisaria criar esta função
	defer ResetShellExecutor()
	
	RegisterPkgModule(L)
	
	script := `
		local pkg = require('pkg')
		local success, msg = pkg.install('test-package')
		return success, msg
	`
	
	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script failed: %v", err)
	}
	
	// Verify success
	success := L.Get(-2).(lua.LBool)
	if !bool(success) {
		t.Error("Expected pkg.install to succeed")
	}
}
```

## Exemplo 2: Testando TaskRunner

```go
package taskrunner

import (
	"testing"
	"time"
	
	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	lua "github.com/yuin/gopher-lua"
)

func TestTaskRunner_SimpleTask(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	// Create a simple task
	task := &types.Task{
		ID:   "test-task-1",
		Name: "test-task",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LBool(true))
			L.Push(lua.LString("Task completed"))
			L.Push(L.NewTable())
			return 3
		}),
	}
	
	// Create task group
	taskGroup := &types.TaskGroup{
		ID:    "test-group",
		Tasks: []types.Task{*task},
	}
	
	// Create runner
	runner := NewTaskRunner(L, taskGroup, nil)
	
	// Execute
	err := runner.Run()
	if err != nil {
		t.Fatalf("Task execution failed: %v", err)
	}
}

func TestTaskRunner_TaskWithDependencies(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	// Task 1 (no dependencies)
	task1 := &types.Task{
		ID:   "task1",
		Name: "first-task",
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			output := L.NewTable()
			output.RawSetString("result", lua.LString("value1"))
			L.Push(lua.LBool(true))
			L.Push(lua.LString("Task 1 done"))
			L.Push(output)
			return 3
		}),
	}
	
	// Task 2 (depends on task1)
	task2 := &types.Task{
		ID:        "task2",
		Name:      "second-task",
		DependsOn: []string{"first-task"},
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			// Get dependency output
			depInput := L.Get(1).(*lua.LTable)
			result := depInput.RawGetString("result")
			
			// Verify we got the expected value
			if result.String() != "value1" {
				L.RaiseError("Expected 'value1' from dependency")
			}
			
			L.Push(lua.LBool(true))
			L.Push(lua.LString("Task 2 done"))
			L.Push(L.NewTable())
			return 3
		}),
	}
	
	taskGroup := &types.TaskGroup{
		ID:    "test-group",
		Tasks: []types.Task{*task1, *task2},
	}
	
	runner := NewTaskRunner(L, taskGroup, nil)
	err := runner.Run()
	if err != nil {
		t.Fatalf("Task execution failed: %v", err)
	}
}

func TestTaskRunner_TaskRetry(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	
	attempts := 0
	task := &types.Task{
		ID:      "retry-task",
		Name:    "task-with-retry",
		Retries: 3,
		CommandFunc: L.NewFunction(func(L *lua.LState) int {
			attempts++
			if attempts < 3 {
				// Fail first 2 attempts
				L.Push(lua.LBool(false))
				L.Push(lua.LString("Failed attempt"))
			} else {
				// Succeed on 3rd attempt
				L.Push(lua.LBool(true))
				L.Push(lua.LString("Success"))
			}
			L.Push(L.NewTable())
			return 3
		}),
	}
	
	taskGroup := &types.TaskGroup{
		ID:    "test-group",
		Tasks: []types.Task{*task},
	}
	
	runner := NewTaskRunner(L, taskGroup, nil)
	err := runner.Run()
	if err != nil {
		t.Fatalf("Task should succeed after retries: %v", err)
	}
	
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}
```

## Exemplo 3: Testando GitOps Manager

```go
package gitops

import (
	"testing"
	"os/exec"
)

// Mock Git operations
type MockGitClient struct {
	cloneFunc  func(url, dest string) error
	pullFunc   func(dir string) error
	commitFunc func(dir, message string) error
}

func (m *MockGitClient) Clone(url, dest string) error {
	if m.cloneFunc != nil {
		return m.cloneFunc(url, dest)
	}
	return nil
}

func TestGitOpsManager_Sync(t *testing.T) {
	tmpDir := t.TempDir()
	
	mockGit := &MockGitClient{
		cloneFunc: func(url, dest string) error {
			// Simulate successful clone
			return nil
		},
	}
	
	manager := NewGitOpsManager(&Config{
		RepoURL:  "https://github.com/test/repo.git",
		Branch:   "main",
		LocalDir: tmpDir,
	})
	manager.SetGitClient(mockGit) // Você precisaria adicionar este método
	
	err := manager.Sync()
	if err != nil {
		t.Errorf("Sync failed: %v", err)
	}
}
```

## Exemplo 4: Testando Stack Manager

```go
package stack

import (
	"testing"
	"time"
)

func TestStackManager_CreateStack(t *testing.T) {
	// Use in-memory database for testing
	dbPath := ":memory:"
	manager, err := NewStackManager(dbPath)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()
	
	// Create a stack
	stack := &StackState{
		ID:          "test-stack-1",
		Name:        "test-stack",
		Description: "Test stack description",
		Version:     "1.0.0",
		Status:      "created",
		CreatedAt:   time.Now(),
	}
	
	err = manager.SaveStack(stack)
	if err != nil {
		t.Fatalf("Failed to save stack: %v", err)
	}
	
	// Retrieve the stack
	retrieved, err := manager.GetStack("test-stack-1")
	if err != nil {
		t.Fatalf("Failed to retrieve stack: %v", err)
	}
	
	if retrieved.Name != stack.Name {
		t.Errorf("Expected name %s, got %s", stack.Name, retrieved.Name)
	}
}

func TestStackManager_ListStacks(t *testing.T) {
	dbPath := ":memory:"
	manager, err := NewStackManager(dbPath)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()
	
	// Create multiple stacks
	for i := 0; i < 5; i++ {
		stack := &StackState{
			ID:   fmt.Sprintf("stack-%d", i),
			Name: fmt.Sprintf("test-stack-%d", i),
		}
		manager.SaveStack(stack)
	}
	
	// List stacks
	stacks, err := manager.ListStacks()
	if err != nil {
		t.Fatalf("Failed to list stacks: %v", err)
	}
	
	if len(stacks) != 5 {
		t.Errorf("Expected 5 stacks, got %d", len(stacks))
	}
}
```

## Exemplo 5: Testando Output Formatter

```go
package output

import (
	"bytes"
	"testing"
	"io"
	"os"
)

func TestPulumiStyleOutput_TaskStart(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	output := NewPulumiStyleOutput()
	output.TaskStart("test-task", "Test task description")
	
	w.Close()
	os.Stdout = old
	
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	result := buf.String()
	if !strings.Contains(result, "test-task") {
		t.Error("Output should contain task name")
	}
}

func TestPulumiStyleOutput_Indentation(t *testing.T) {
	output := NewPulumiStyleOutput()
	
	initialIndent := output.indent
	output.Indent()
	if output.indent != initialIndent+1 {
		t.Error("Indent should increase")
	}
	
	output.Unindent()
	if output.indent != initialIndent {
		t.Error("Unindent should decrease")
	}
	
	// Test that unindent doesn't go negative
	output.Unindent()
	output.Unindent()
	if output.indent < 0 {
		t.Error("Indent should not be negative")
	}
}
```

## Exemplo 6: Testando com Goroutines

```go
package core

import (
	"sync"
	"testing"
	"time"
)

func TestWorkerPool_ConcurrentExecution(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Close()
	
	var mu sync.Mutex
	counter := 0
	numTasks := 100
	
	// Submit many tasks
	for i := 0; i < numTasks; i++ {
		pool.Submit(func() {
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}
	
	// Wait for completion
	deadline := time.Now().Add(10 * time.Second)
	for counter < numTasks && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}
	
	if counter != numTasks {
		t.Errorf("Expected %d tasks completed, got %d", numTasks, counter)
	}
}
```

## Exemplo 7: Testando com Table-Driven Tests

```go
package system

import "testing"

func TestValidatePackageName(t *testing.T) {
	tests := []struct {
		name      string
		pkgName   string
		wantValid bool
	}{
		{"valid simple", "nginx", true},
		{"valid with dash", "nginx-full", true},
		{"valid with number", "python3", true},
		{"empty name", "", false},
		{"with spaces", "package name", false},
		{"with special chars", "package@123", false},
		{"too long", string(make([]byte, 300)), false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := ValidatePackageName(tt.pkgName)
			if valid != tt.wantValid {
				t.Errorf("ValidatePackageName(%q) = %v, want %v",
					tt.pkgName, valid, tt.wantValid)
			}
		})
	}
}
```

## Exemplo 8: Benchmarking

```go
package core

import "testing"

func BenchmarkSafeMap_Set(b *testing.B) {
	sm := NewSafeMap()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Set("key", i)
	}
}

func BenchmarkSafeMap_Get(b *testing.B) {
	sm := NewSafeMap()
	sm.Set("key", "value")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Get("key")
	}
}

func BenchmarkSafeMap_Concurrent(b *testing.B) {
	sm := NewSafeMap()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%1000)
			sm.Set(key, i)
			sm.Get(key)
			i++
		}
	})
}
```

## Dicas Importantes

### 1. Isolamento de Testes

- Use `t.TempDir()` para diretórios temporários
- Use `:memory:` para bancos de dados SQLite em testes
- Limpe recursos com `defer`

### 2. Mocking

- Crie interfaces para dependências externas
- Use dependency injection
- Considere usar bibliotecas como `testify/mock`

### 3. Testes Paralelos

```go
func TestParallel(t *testing.T) {
	t.Parallel() // Run in parallel with other parallel tests
	
	tests := []struct{
		name string
		// ...
	}{
		// ...
	}
	
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run subtests in parallel
			// test code
		})
	}
}
```

### 4. Test Helpers

```go
// Create helper functions for common setup
func setupTestEnvironment(t *testing.T) (*Environment, func()) {
	env := &Environment{
		TempDir: t.TempDir(),
	}
	
	cleanup := func() {
		// Cleanup code
	}
	
	return env, cleanup
}

func TestWithHelper(t *testing.T) {
	env, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Use env in test
}
```

### 5. Golden Files

Para testes de output complexo:

```go
func TestComplexOutput(t *testing.T) {
	result := GenerateComplexOutput()
	
	goldenFile := "testdata/complex_output.golden"
	
	if *update {
		os.WriteFile(goldenFile, []byte(result), 0644)
	}
	
	expected, _ := os.ReadFile(goldenFile)
	if result != string(expected) {
		t.Errorf("Output doesn't match golden file")
	}
}
```

## Conclusão

Estes exemplos cobrem os padrões mais comuns de teste para o Sloth Runner. Use-os como referência ao criar novos testes.

**Lembre-se:**
- Teste o comportamento, não a implementação
- Um teste deve testar uma coisa
- Testes devem ser independentes
- Use nomes descritivos
- Documente casos de edge
