# Test Coverage Report - Sloth Runner

## Current Status

**Overall Coverage: 9.0%**

## Summary of Work Completed

### New Test Files Created

1. **internal/types/types_test.go** ✅ 100% coverage
   - Tests for all struct types
   - Tests for ID generation functions
   - Complete validation of data structures

2. **internal/runner/simple_runner_test.go** ✅ 100% coverage
   - Tests for RunSingleTask function
   - Success and failure scenarios
   - Output validation
   - Parameter handling
   - Timing accuracy tests

3. **internal/agent/health_test.go** ✅ 69% coverage
   - HealthMonitor initialization
   - Metrics collection (system, memory, CPU, disk)
   - Health status determination
   - Threshold management
   - Collector pattern
   - Concurrent access safety
   - Periodic collection

4. **internal/agent/metrics_test.go** ✅ 69% coverage
   - MetricsCollector initialization
   - System metrics collection
   - Runtime metrics collection
   - Task metrics tracking
   - Custom metrics management
   - HTTP handlers
   - JSON marshaling/unmarshaling
   - Concurrent access tests
   - Benchmarks

### Enhanced Existing Tests

5. **internal/core/core_test.go** ✅ 42.8% coverage (up from 34.1%)
   - Added tests for Semaphore
   - Added tests for RWCounter
   - Added tests for RateLimiter
   - Added tests for SafeMap operations (ForEach, concurrent access)
   - Added tests for WorkerPool with timeout
   - Added tests for concurrent operations
   - Added benchmarks

## Coverage by Package

| Package | Coverage | Status | Notes |
|---------|----------|--------|-------|
| internal/types | 100.0% | ✅ Complete | All functions tested |
| internal/runner | 100.0% | ✅ Complete | All functions tested |
| internal/agent | 69.0% | 🟡 Good | Health & metrics modules covered |
| internal/core | 42.8% | 🟡 Moderate | Concurrency primitives tested |
| internal/scheduler | 82.7% | ✅ Excellent | Pre-existing tests |
| internal/state | 61.1% | 🟡 Good | Pre-existing tests |
| internal/reliability | 42.3% | 🟡 Moderate | Pre-existing tests |
| internal/modules | 41.3% | 🟡 Moderate | Pre-existing tests |
| internal/taskrunner | 30.1% | 🔴 Needs work | Large module, complex |
| cmd/sloth-runner | 11.3% | 🔴 Needs work | CLI commands |
| internal/luainterface | 1.1% | 🔴 Needs work | Very large (1590+ lines) |
| internal/gitops | 0.0% | 🔴 Not tested | Needs tests |
| internal/output | 0.0% | 🔴 Not tested | Needs tests |
| internal/repl | 0.0% | 🔴 Not tested | Interactive, hard to test |
| internal/scaffolding | 0.0% | 🔴 Not tested | Needs tests |
| internal/stack | 0.0% | 🔴 Not tested | Needs tests |
| internal/ui | 0.0% | 🔴 Not tested | Needs tests |
| internal/ai | 0.0% | 🔴 Not tested | Needs tests |

## Roadmap to 80% Coverage

### Phase 1: Low-Hanging Fruit (Target: +20%, Total: ~29%)
**Estimated Time: 4-6 hours**

Priority modules (small, high impact):

1. **internal/output** (~300 lines)
   - Test PulumiStyleOutput struct
   - Test formatting functions
   - Test indentation logic
   - Test progress bars

2. **internal/stack** (~500 lines)
   - Test StackManager CRUD operations
   - Test database operations
   - Test state management

3. **internal/scaffolding** (~700 lines)
   - Test template rendering
   - Test workflow initialization
   - Test file generation

4. **internal/repl** (~100 lines)
   - Basic initialization tests
   - Executor function tests (with mocks)
   - Completer logic tests

### Phase 2: Core Functionality (Target: +30%, Total: ~59%)
**Estimated Time: 8-12 hours**

1. **internal/taskrunner** (1109 lines) - Priority High
   - Test task execution flow
   - Test dependency resolution
   - Test error handling
   - Test hooks (PreExec, PostExec, OnSuccess, OnFailure)
   - Test async task execution
   - Test retry logic

2. **internal/core** (improve from 42.8% to 70%)
   - Add performance tests
   - Add utils tests
   - Add error handling tests
   - Add cache tests

3. **cmd/sloth-runner** (improve from 11.3% to 50%)
   - Test command parsing
   - Test flag handling
   - Test agent commands
   - Test master commands
   - Use CLI testing patterns

### Phase 3: Lua Interface & Modules (Target: +15%, Total: ~74%)
**Estimated Time: 12-16 hours**

1. **internal/luainterface** (1590+ lines split across many files)
   Strategy: Test each module file separately
   
   - terraform_advanced.go (1511 lines)
   - modern_dsl.go (1370 lines)  
   - kubernetes.go (945 lines)
   - helm.go (908 lines)
   - pkg.go (789 lines)
   - security.go (756 lines)
   - system.go (731 lines)
   - infratest_module.go (720 lines)
   
   For each file:
   - Test module registration
   - Test main exported functions
   - Test error handling
   - Mock external dependencies (shell commands, APIs)

2. **internal/modules** (improve to 70%)
   - Test core module functions
   - Test module registration

### Phase 4: GitOps & Advanced Features (Target: +6%, Total: ~80%)
**Estimated Time: 6-8 hours**

1. **internal/gitops** (4 files, ~500 lines each)
   - Test diff engine
   - Test rollback engine
   - Test sync controller
   - Test manager
   - Mock Git operations

2. **internal/ai** 
   - Test AI integration
   - Mock AI API calls

3. **internal/ui**
   - Test UI server
   - Test HTTP handlers

## Testing Best Practices for This Project

### 1. Use Table-Driven Tests
```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test", "TEST", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr %v, got %v", tt.wantErr, err)
            }
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### 2. Mock External Dependencies

For shell commands, use interfaces:
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

### 3. Use Temporary Directories

```go
func TestFileOperations(t *testing.T) {
    tmpDir := t.TempDir() // Auto-cleanup
    // Use tmpDir for file operations
}
```

### 4. Test Concurrency with Race Detector

```bash
go test -race ./...
```

### 5. Use Subtests for Organization

```go
func TestModule(t *testing.T) {
    t.Run("initialization", func(t *testing.T) {
        // test initialization
    })
    
    t.Run("execution", func(t *testing.T) {
        // test execution
    })
}
```

## Tools and Commands

### Run tests with coverage
```bash
go test ./... -coverprofile=coverage.out -coverpkg=./...
```

### View coverage report
```bash
go tool cover -html=coverage.out
```

### View coverage by function
```bash
go tool cover -func=coverage.out
```

### Run tests with race detector
```bash
go test -race ./...
```

### Run specific package tests
```bash
go test -v ./internal/core
```

### Run specific test
```bash
go test -v -run TestGlobalCore ./internal/core
```

### Benchmark tests
```bash
go test -bench=. ./internal/core
```

## Integration Test Strategy

For reaching 80% coverage efficiently, consider:

1. **Integration tests** that exercise multiple components
2. **E2E tests** for CLI commands
3. **Property-based testing** for complex logic
4. **Fuzz testing** for input validation

## Automated Testing

Consider setting up:

1. **GitHub Actions** for CI/CD
   ```yaml
   - name: Run tests
     run: go test -race -coverprofile=coverage.out ./...
   
   - name: Coverage threshold
     run: |
       coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
       if (( $(echo "$coverage < 80" | bc -l) )); then
         echo "Coverage $coverage% is below 80%"
         exit 1
       fi
   ```

2. **Pre-commit hooks** to run tests locally

3. **Coverage badges** in README.md

## Notes

- Some modules like `internal/repl` are harder to test due to interactive nature
- The `internal/luainterface` module is very large and split across many files
- Mock external dependencies (git, kubectl, terraform, etc.) to avoid external system requirements
- Use `t.Parallel()` for tests that can run in parallel
- Some proto-generated files don't need manual tests

## Conclusion

The foundation has been laid with comprehensive tests for:
- Core data types (100% coverage)
- Core utilities and concurrency primitives (42.8% coverage)
- Agent health and metrics monitoring (69% coverage)
- Task runner basics (100% coverage)

The remaining work to reach 80% is well-defined and can be tackled incrementally following the phased approach above. Each phase builds on the previous one and targets high-impact areas first.

**Next immediate steps:**
1. Complete Phase 1 (output, stack, scaffolding) - Quick wins
2. Improve taskrunner coverage - Critical path
3. Add Lua module tests with proper mocking - Largest impact
