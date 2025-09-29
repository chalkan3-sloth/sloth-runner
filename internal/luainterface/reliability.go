package luainterface

import (
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/reliability"
	lua "github.com/yuin/gopher-lua"
)

// ReliabilityModule provides reliability patterns for Lua scripts
type ReliabilityModule struct {
	circuitBreakerManager *reliability.CircuitBreakerManager
	defaultRetryConfig    reliability.RetryConfig
}

// NewReliabilityModule creates a new reliability module
func NewReliabilityModule() *ReliabilityModule {
	return &ReliabilityModule{
		circuitBreakerManager: reliability.NewCircuitBreakerManager(),
		defaultRetryConfig:    reliability.DefaultRetryConfig(),
	}
}

// Loader is the module loader function for Lua
func (rm *ReliabilityModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"retry":                rm.luaRetry,
		"retry_with_config":    rm.luaRetryWithConfig,
		"circuit_breaker":      rm.luaCircuitBreaker,
		"get_circuit_stats":    rm.luaGetCircuitStats,
		"reset_circuit":        rm.luaResetCircuit,
		"list_circuits":        rm.luaListCircuits,
	})

	// Add constants
	constants := L.NewTable()
	constants.RawSetString("FIXED_DELAY", lua.LNumber(reliability.FixedDelay))
	constants.RawSetString("EXPONENTIAL_BACKOFF", lua.LNumber(reliability.ExponentialBackoff))
	constants.RawSetString("LINEAR_BACKOFF", lua.LNumber(reliability.LinearBackoff))
	constants.RawSetString("CUSTOM_BACKOFF", lua.LNumber(reliability.CustomBackoff))
	mod.RawSetString("strategy", constants)

	L.Push(mod)
	return 1
}

// ReliabilityLoader is the global loader function
func ReliabilityLoader(L *lua.LState) int {
	return NewReliabilityModule().Loader(L)
}

// luaRetry provides simple retry functionality
func (rm *ReliabilityModule) luaRetry(L *lua.LState) int {
	maxAttempts := int(L.CheckNumber(1))
	initialDelay := time.Duration(L.CheckNumber(2)) * time.Second
	fn := L.CheckFunction(3)

	retrier := reliability.NewRetrier(reliability.RetryConfig{
		MaxAttempts:  maxAttempts,
		InitialDelay: initialDelay,
		Strategy:     reliability.ExponentialBackoff,
		Multiplier:   2.0,
		Jitter:       true,
	})

	result, err := retrier.Execute(func() (interface{}, error) {
		L.Push(fn)
		err := L.PCall(0, lua.MultRet, nil)
		if err != nil {
			return nil, fmt.Errorf("function execution failed: %v", err)
		}

		// Get return values
		top := L.GetTop()
		if top == 0 {
			return nil, nil
		}

		// First return value is the result, second (if exists) is error
		result := L.Get(1)
		
		if top > 1 {
			if errVal := L.Get(2); errVal != lua.LNil {
				if errStr, ok := errVal.(lua.LString); ok && string(errStr) != "" {
					return result, fmt.Errorf("%s", string(errStr))
				}
			}
		}

		return result, nil
	})

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if result != nil {
		if luaVal, ok := result.(lua.LValue); ok {
			L.Push(luaVal)
		} else {
			L.Push(lua.LString(fmt.Sprintf("%v", result)))
		}
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

// luaRetryWithConfig provides retry with full configuration
func (rm *ReliabilityModule) luaRetryWithConfig(L *lua.LState) int {
	configTable := L.CheckTable(1)
	fn := L.CheckFunction(2)

	// Parse configuration
	config := reliability.DefaultRetryConfig()
	
	if maxAttempts := configTable.RawGetString("max_attempts"); maxAttempts != lua.LNil {
		config.MaxAttempts = int(maxAttempts.(lua.LNumber))
	}
	
	if initialDelay := configTable.RawGetString("initial_delay"); initialDelay != lua.LNil {
		config.InitialDelay = time.Duration(initialDelay.(lua.LNumber)) * time.Second
	}
	
	if maxDelay := configTable.RawGetString("max_delay"); maxDelay != lua.LNil {
		config.MaxDelay = time.Duration(maxDelay.(lua.LNumber)) * time.Second
	}
	
	if strategy := configTable.RawGetString("strategy"); strategy != lua.LNil {
		config.Strategy = reliability.RetryStrategy(strategy.(lua.LNumber))
	}
	
	if multiplier := configTable.RawGetString("multiplier"); multiplier != lua.LNil {
		config.Multiplier = float64(multiplier.(lua.LNumber))
	}
	
	if jitter := configTable.RawGetString("jitter"); jitter != lua.LNil {
		config.Jitter = bool(jitter.(lua.LBool))
	}

	// Add retry callback if provided
	if onRetryFn := configTable.RawGetString("on_retry"); onRetryFn != lua.LNil {
		if retryFunc, ok := onRetryFn.(*lua.LFunction); ok {
			config.OnRetry = func(attempt int, delay time.Duration, err error) {
				L.Push(retryFunc)
				L.Push(lua.LNumber(attempt))
				L.Push(lua.LNumber(delay.Seconds()))
				L.Push(lua.LString(err.Error()))
				L.PCall(3, 0, nil)
			}
		}
	}

	retrier := reliability.NewRetrier(config)

	result, err := retrier.Execute(func() (interface{}, error) {
		L.Push(fn)
		err := L.PCall(0, lua.MultRet, nil)
		if err != nil {
			return nil, fmt.Errorf("function execution failed: %v", err)
		}

		top := L.GetTop()
		if top == 0 {
			return nil, nil
		}

		result := L.Get(1)
		
		if top > 1 {
			if errVal := L.Get(2); errVal != lua.LNil {
				if errStr, ok := errVal.(lua.LString); ok && string(errStr) != "" {
					return result, fmt.Errorf("%s", string(errStr))
				}
			}
		}

		return result, nil
	})

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if result != nil {
		if luaVal, ok := result.(lua.LValue); ok {
			L.Push(luaVal)
		} else {
			L.Push(lua.LString(fmt.Sprintf("%v", result)))
		}
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

// luaCircuitBreaker provides circuit breaker functionality
func (rm *ReliabilityModule) luaCircuitBreaker(L *lua.LState) int {
	name := L.CheckString(1)
	configTable := L.CheckTable(2)
	fn := L.CheckFunction(3)

	// Parse circuit breaker configuration
	config := reliability.CircuitBreakerConfig{
		MaxFailures:      5,
		Timeout:          60 * time.Second,
		SuccessThreshold: 1,
	}

	if maxFailures := configTable.RawGetString("max_failures"); maxFailures != lua.LNil {
		config.MaxFailures = int(maxFailures.(lua.LNumber))
	}
	
	if timeout := configTable.RawGetString("timeout"); timeout != lua.LNil {
		config.Timeout = time.Duration(timeout.(lua.LNumber)) * time.Second
	}
	
	if successThreshold := configTable.RawGetString("success_threshold"); successThreshold != lua.LNil {
		config.SuccessThreshold = int(successThreshold.(lua.LNumber))
	}

	// Add state change callback if provided
	if onStateChangeFn := configTable.RawGetString("on_state_change"); onStateChangeFn != lua.LNil {
		if stateChangeFunc, ok := onStateChangeFn.(*lua.LFunction); ok {
			config.OnStateChange = func(from, to reliability.CircuitBreakerState) {
				L.Push(stateChangeFunc)
				L.Push(lua.LString(from.String()))
				L.Push(lua.LString(to.String()))
				L.PCall(2, 0, nil)
			}
		}
	}

	cb := rm.circuitBreakerManager.GetOrCreate(name, config)

	result, err := cb.Execute(func() (interface{}, error) {
		L.Push(fn)
		err := L.PCall(0, lua.MultRet, nil)
		if err != nil {
			return nil, fmt.Errorf("function execution failed: %v", err)
		}

		top := L.GetTop()
		if top == 0 {
			return nil, nil
		}

		result := L.Get(1)
		
		if top > 1 {
			if errVal := L.Get(2); errVal != lua.LNil {
				if errStr, ok := errVal.(lua.LString); ok && string(errStr) != "" {
					return result, fmt.Errorf("%s", string(errStr))
				}
			}
		}

		return result, nil
	})

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if result != nil {
		if luaVal, ok := result.(lua.LValue); ok {
			L.Push(luaVal)
		} else {
			L.Push(lua.LString(fmt.Sprintf("%v", result)))
		}
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

// luaGetCircuitStats returns circuit breaker statistics
func (rm *ReliabilityModule) luaGetCircuitStats(L *lua.LState) int {
	name := L.CheckString(1)

	cb, exists := rm.circuitBreakerManager.Get(name)
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("circuit breaker not found: " + name))
		return 2
	}

	stats := cb.GetStats()
	statsTable := L.NewTable()
	
	statsTable.RawSetString("requests", lua.LNumber(stats.Requests))
	statsTable.RawSetString("total_success", lua.LNumber(stats.TotalSuccess))
	statsTable.RawSetString("total_failures", lua.LNumber(stats.TotalFailures))
	statsTable.RawSetString("consecutive_success", lua.LNumber(stats.ConsecutiveSuccess))
	statsTable.RawSetString("consecutive_failures", lua.LNumber(stats.ConsecutiveFailures))
	statsTable.RawSetString("state", lua.LString(cb.GetState().String()))
	
	if !stats.LastSuccessTime.IsZero() {
		statsTable.RawSetString("last_success_time", lua.LString(stats.LastSuccessTime.Format(time.RFC3339)))
	}
	
	if !stats.LastFailureTime.IsZero() {
		statsTable.RawSetString("last_failure_time", lua.LString(stats.LastFailureTime.Format(time.RFC3339)))
	}

	L.Push(statsTable)
	return 1
}

// luaResetCircuit resets a circuit breaker
func (rm *ReliabilityModule) luaResetCircuit(L *lua.LState) int {
	name := L.CheckString(1)

	cb, exists := rm.circuitBreakerManager.Get(name)
	if !exists {
		L.Push(lua.LFalse)
		L.Push(lua.LString("circuit breaker not found: " + name))
		return 2
	}

	cb.Reset()
	L.Push(lua.LTrue)
	return 1
}

// luaListCircuits returns all circuit breaker names
func (rm *ReliabilityModule) luaListCircuits(L *lua.LState) int {
	names := rm.circuitBreakerManager.List()
	
	table := L.NewTable()
	for i, name := range names {
		table.RawSetInt(i+1, lua.LString(name))
	}
	
	L.Push(table)
	return 1
}