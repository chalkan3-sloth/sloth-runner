package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/chalkan3/sloth-runner/internal/luainterface"
	"github.com/pterm/pterm"
	lua "github.com/yuin/gopher-lua"
)

// Enhanced demo showing improved DSL and capabilities
func main() {
	pterm.DefaultHeader.WithFullWidth().Println("🦥 Sloth Runner Enhanced Demo")
	
	// Create enhanced Lua environment
	L := lua.NewState()
	defer L.Close()
	
	// Setup enhanced Lua environment with modern DSL
	setupEnhancedLuaEnvironment(L)
	
	// Load enhanced example
	examplePath := filepath.Join("examples", "enhanced_demo.lua")
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		pterm.Warning.Printf("Demo file not found: %s, using inline demo\n", examplePath)
		runInlineDemo(L)
		return
	}
	
	pterm.Info.Printf("Loading enhanced pipeline: %s\n", examplePath)
	
	// Execute enhanced example
	if err := L.DoFile(examplePath); err != nil {
		pterm.Error.Printf("Failed to execute enhanced pipeline: %v\n", err)
		return
	}
	
	// Show enhanced features
	demonstrateEnhancedFeatures(L)
	
	pterm.Success.Println("Enhanced demo completed successfully!")
}

// setupEnhancedLuaEnvironment sets up the enhanced Lua environment
func setupEnhancedLuaEnvironment(L *lua.LState) {
	// Setup import function
	luainterface.OpenImport(L, "examples/enhanced_modern_pipeline.lua")
	
	// Register enhanced DSL functions
	registerEnhancedDSL(L)
	
	// Register enhanced modules
	registerEnhancedModules(L)
	
	pterm.Info.Println("Enhanced Lua environment initialized")
}

// registerEnhancedDSL registers the modern DSL functions
func registerEnhancedDSL(L *lua.LState) {
	// Core system functions
	coreTable := L.NewTable()
	coreTable.RawSetString("stats", L.NewFunction(func(L *lua.LState) int {
		statsTable := L.NewTable()
		statsTable.RawSetString("uptime_seconds", lua.LNumber(time.Since(time.Now()).Seconds()))
		statsTable.RawSetString("memory_alloc", lua.LNumber(1024*1024)) // 1MB
		statsTable.RawSetString("worker_active", lua.LNumber(4))
		statsTable.RawSetString("tasks_executed", lua.LNumber(10))
		L.Push(statsTable)
		return 1
	}))
	
	coreTable.RawSetString("submit", L.NewFunction(func(L *lua.LState) int {
		taskFunc := L.CheckFunction(1)
		context := L.OptString(2, "lua_task")
		
		pterm.Info.Printf("Submitting task to core worker pool: %s\n", context)
		
		// Simulate task execution
		go func() {
			L.CallByParam(lua.P{
				Fn:      taskFunc,
				NRet:    0,
				Protect: true,
			})
		}()
		
		L.Push(lua.LBool(true))
		return 1
	}))
	
	L.SetGlobal("core", coreTable)
	
	// Async operations
	asyncTable := L.NewTable()
	asyncTable.RawSetString("parallel", L.NewFunction(func(L *lua.LState) int {
		tasks := L.CheckTable(1)
		maxWorkers := L.OptInt(2, 4)
		
		pterm.Info.Printf("Executing parallel tasks with %d workers\n", maxWorkers)
		
		results := L.NewTable()
		tasks.ForEach(func(key, value lua.LValue) {
			if taskFunc, ok := value.(*lua.LFunction); ok {
				pterm.Debug.Printf("Executing parallel task: %s\n", lua.LVAsString(key))
				
				// Simulate task execution
				L.CallByParam(lua.P{
					Fn:      taskFunc,
					NRet:    1,
					Protect: true,
				})
				
				if L.GetTop() > 0 {
					result := L.Get(-1)
					results.RawSet(key, result)
					L.Pop(1)
				}
			}
		})
		
		L.Push(results)
		L.Push(lua.LNil) // No errors
		return 2
	}))
	
	asyncTable.RawSetString("sleep", L.NewFunction(func(L *lua.LState) int {
		ms := L.CheckInt(1)
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return 0
	}))
	
	L.SetGlobal("async", asyncTable)
	
	// Performance monitoring
	perfTable := L.NewTable()
	perfTable.RawSetString("measure", L.NewFunction(func(L *lua.LState) int {
		taskFunc := L.CheckFunction(1)
		name := L.OptString(2, "unnamed_task")
		
		start := time.Now()
		pterm.Debug.Printf("Starting performance measurement: %s\n", name)
		
		L.CallByParam(lua.P{
			Fn:      taskFunc,
			NRet:    1,
			Protect: true,
		})
		
		duration := time.Since(start)
		pterm.Info.Printf("Task '%s' completed in %v\n", name, duration)
		
		var result lua.LValue = lua.LNil
		if L.GetTop() > 0 {
			result = L.Get(-1)
			L.Pop(1)
		}
		
		L.Push(result)
		L.Push(lua.LNumber(duration.Milliseconds()))
		L.Push(lua.LNil) // No error
		return 3
	}))
	
	perfTable.RawSetString("memory", L.NewFunction(func(L *lua.LState) int {
		memTable := L.NewTable()
		memTable.RawSetString("current_mb", lua.LNumber(64))
		memTable.RawSetString("peak_mb", lua.LNumber(128))
		memTable.RawSetString("usage_percent", lua.LNumber(25.5))
		L.Push(memTable)
		return 1
	}))
	
	L.SetGlobal("perf", perfTable)
	
	// Flow control
	flowTable := L.NewTable()
	flowTable.RawSetString("circuit_breaker", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		taskFunc := L.CheckFunction(2)
		
		pterm.Info.Printf("Executing with circuit breaker: %s\n", name)
		
		L.CallByParam(lua.P{
			Fn:      taskFunc,
			NRet:    1,
			Protect: true,
		})
		
		var result lua.LValue = lua.LNil
		if L.GetTop() > 0 {
			result = L.Get(-1)
			L.Pop(1)
		}
		
		L.Push(result)
		L.Push(lua.LNil) // No error
		return 2
	}))
	
	flowTable.RawSetString("rate_limit", L.NewFunction(func(L *lua.LState) int {
		rps := L.CheckInt(1)
		taskFunc := L.CheckFunction(2)
		
		pterm.Debug.Printf("Rate limiting to %d RPS\n", rps)
		
		if rps > 0 {
			time.Sleep(time.Second / time.Duration(rps))
		}
		
		L.CallByParam(lua.P{
			Fn:      taskFunc,
			NRet:    1,
			Protect: true,
		})
		
		var result lua.LValue = lua.LNil
		if L.GetTop() > 0 {
			result = L.Get(-1)
			L.Pop(1)
		}
		
		L.Push(result)
		L.Push(lua.LNil) // No error
		return 2
	}))
	
	L.SetGlobal("flow", flowTable)
	
	// Error handling
	errorTable := L.NewTable()
	errorTable.RawSetString("try", L.NewFunction(func(L *lua.LState) int {
		tryFunc := L.CheckFunction(1)
		catchFunc := L.OptFunction(2, nil)
		finallyFunc := L.OptFunction(3, nil)
		
		pterm.Debug.Println("Executing try-catch block")
		
		var result lua.LValue = lua.LNil
		var caught lua.LValue = lua.LNil
		
		// Execute try block
		err := L.CallByParam(lua.P{
			Fn:      tryFunc,
			NRet:    1,
			Protect: true,
		})
		
		if err != nil {
			caught = lua.LString(err.Error())
			pterm.Warning.Printf("Caught error: %v\n", err)
			
			if catchFunc != nil {
				L.CallByParam(lua.P{
					Fn:      catchFunc,
					NRet:    0,
					Protect: true,
				}, caught)
			}
		} else if L.GetTop() > 0 {
			result = L.Get(-1)
			L.Pop(1)
		}
		
		// Execute finally block
		if finallyFunc != nil {
			L.CallByParam(lua.P{
				Fn:      finallyFunc,
				NRet:    0,
				Protect: true,
			})
		}
		
		L.Push(result)
		L.Push(caught)
		return 2
	}))
	
	errorTable.RawSetString("retry", L.NewFunction(func(L *lua.LState) int {
		taskFunc := L.CheckFunction(1)
		maxAttempts := L.OptInt(2, 3)
		initialDelayMs := L.OptInt(3, 1000)
		
		pterm.Info.Printf("Retrying task up to %d times\n", maxAttempts)
		
		var result lua.LValue = lua.LNil
		var lastError error
		
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			err := L.CallByParam(lua.P{
				Fn:      taskFunc,
				NRet:    1,
				Protect: true,
			})
			
			if err == nil {
				if L.GetTop() > 0 {
					result = L.Get(-1)
					L.Pop(1)
				}
				break
			}
			
			lastError = err
			pterm.Warning.Printf("Attempt %d failed: %v\n", attempt, err)
			
			if attempt < maxAttempts {
				delay := time.Duration(initialDelayMs*attempt) * time.Millisecond
				pterm.Debug.Printf("Waiting %v before retry\n", delay)
				time.Sleep(delay)
			}
		}
		
		L.Push(result)
		if lastError != nil {
			L.Push(lua.LString(lastError.Error()))
		} else {
			L.Push(lua.LNil)
		}
		return 2
	}))
	
	L.SetGlobal("error", errorTable)
}

// registerEnhancedModules registers enhanced utility modules
func registerEnhancedModules(L *lua.LState) {
	// Utilities
	utilsTable := L.NewTable()
	utilsTable.RawSetString("config", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		defaultValue := L.OptString(2, "")
		
		// Simulate config retrieval
		value := os.Getenv(key)
		if value == "" {
			value = defaultValue
		}
		
		L.Push(lua.LString(value))
		return 1
	}))
	
	utilsTable.RawSetString("secret", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		
		// Simulate secret retrieval
		pterm.Info.Printf("Retrieving secret: %s\n", key)
		L.Push(lua.LString("***SECRET***"))
		L.Push(lua.LNil)
		return 2
	}))
	
	L.SetGlobal("utils", utilsTable)
	
	// Task functions
	taskTable := L.NewTable()
	taskTable.RawSetString("checkpoint", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		_ = L.OptTable(2, nil) // state parameter, currently unused
		
		pterm.Info.Printf("Creating checkpoint: %s\n", name)
		
		L.Push(lua.LString(name))
		return 1
	}))
	
	L.SetGlobal("task", taskTable)
	
	// Workflow functions
	workflowTable := L.NewTable()
	workflowTable.RawSetString("define", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		_ = L.CheckTable(2) // definition parameter, currently unused
		
		pterm.Info.Printf("Defining workflow: %s\n", name)
		
		L.Push(lua.LBool(true))
		L.Push(lua.LNil)
		return 2
	}))
	
	workflowTable.RawSetString("parallel", L.NewFunction(func(L *lua.LState) int {
		tasks := L.CheckTable(1)
		options := L.OptTable(2, nil)
		
		pterm.Info.Println("Creating parallel workflow configuration")
		
		parallelConfig := L.NewTable()
		parallelConfig.RawSetString("type", lua.LString("parallel"))
		parallelConfig.RawSetString("tasks", tasks)
		if options != nil {
			parallelConfig.RawSetString("options", options)
		}
		
		L.Push(parallelConfig)
		return 1
	}))
	
	L.SetGlobal("workflow", workflowTable)
	
	// Logging
	logTable := L.NewTable()
	logTable.RawSetString("info", L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		data := L.OptTable(2, nil)
		
		if data != nil {
			pterm.Info.Printf("%s (with data)\n", message)
		} else {
			pterm.Info.Println(message)
		}
		return 0
	}))
	
	logTable.RawSetString("warn", L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		pterm.Warning.Println(message)
		return 0
	}))
	
	logTable.RawSetString("error", L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		pterm.Error.Println(message)
		return 0
	}))
	
	logTable.RawSetString("debug", L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		pterm.Debug.Println(message)
		return 0
	}))
	
	L.SetGlobal("log", logTable)
}

// demonstrateEnhancedFeatures shows the enhanced features in action
func demonstrateEnhancedFeatures(L *lua.LState) {
	pterm.DefaultHeader.WithFullWidth().Println("🚀 Enhanced Features Demo")
	
	// Demonstrate core stats
	pterm.DefaultSection.Println("Core System Statistics")
	if err := L.DoString(`
		local stats = core.stats()
		log.info("Core Statistics Retrieved", stats)
		log.info("Memory usage: " .. stats.memory_alloc .. " bytes")
		log.info("Active workers: " .. stats.worker_active)
	`); err != nil {
		pterm.Error.Printf("Failed to demonstrate core stats: %v\n", err)
	}
	
	// Demonstrate parallel execution
	pterm.DefaultSection.Println("Parallel Task Execution")
	if err := L.DoString(`
		local results, errors = async.parallel({
			task1 = function()
				log.info("Executing task 1")
				async.sleep(100)
				return "result1"
			end,
			task2 = function()
				log.info("Executing task 2")
				async.sleep(150)
				return "result2"
			end,
			task3 = function()
				log.info("Executing task 3")
				async.sleep(80)
				return "result3"
			end
		}, 3)
		
		log.info("Parallel execution completed")
		if errors then
			log.error("Some tasks failed")
		else
			log.info("All tasks completed successfully")
		end
	`); err != nil {
		pterm.Error.Printf("Failed to demonstrate parallel execution: %v\n", err)
	}
	
	// Demonstrate performance monitoring
	pterm.DefaultSection.Println("Performance Monitoring")
	if err := L.DoString(`
		local result, duration, err = perf.measure(function()
			log.info("Performing monitored task")
			async.sleep(200)
			return "monitored_result"
		end, "demo_task")
		
		log.info("Task completed in " .. duration .. " ms")
		
		local memory = perf.memory()
		log.info("Memory usage: " .. memory.usage_percent .. "%")
	`); err != nil {
		pterm.Error.Printf("Failed to demonstrate performance monitoring: %v\n", err)
	}
	
	// Demonstrate error handling
	pterm.DefaultSection.Println("Advanced Error Handling")
	if err := L.DoString(`
		local result, caught = error.try(
			function()
				log.info("Attempting operation that might fail")
				-- Simulate operation
				return "success"
			end,
			function(err)
				log.warn("Caught error: " .. err)
			end,
			function()
				log.info("Finally block executed")
			end
		)
		
		log.info("Try-catch completed")
		
		-- Demonstrate retry
		local retry_result, retry_err = error.retry(function()
			log.info("Retryable operation")
			return "retry_success"
		end, 3, 1000)
		
		log.info("Retry operation completed")
	`); err != nil {
		pterm.Error.Printf("Failed to demonstrate error handling: %v\n", err)
	}
	
	// Demonstrate circuit breaker
	pterm.DefaultSection.Println("Circuit Breaker Pattern")
	if err := L.DoString(`
		local cb_result, cb_err = flow.circuit_breaker("demo_service", function()
			log.info("Operation protected by circuit breaker")
			return "protected_result"
		end)
		
		if cb_err then
			log.error("Circuit breaker prevented execution")
		else
			log.info("Circuit breaker allowed execution")
		end
	`); err != nil {
		pterm.Error.Printf("Failed to demonstrate circuit breaker: %v\n", err)
	}
}

func init() {
	// Configure pterm for better output
	pterm.DefaultLogger.Level = pterm.LogLevelInfo
	pterm.EnableColor()
}

// runInlineDemo runs an inline demo if file is not found
func runInlineDemo(L *lua.LState) {
	pterm.Info.Println("Running inline enhanced features demo")
	demonstrateEnhancedFeatures(L)
}