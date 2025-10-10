package luamodules

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/agent"
	"github.com/google/uuid"
	"github.com/yuin/gopher-lua"
)

// EventRegisterModule provides Lua API for registering event watchers
type EventRegisterModule struct{}

// AgentResolver interface for resolving agent names to addresses
type AgentResolver interface {
	ResolveAgent(name string) (string, error)
}

// NewEventRegisterModule creates a new event register module
func NewEventRegisterModule() *EventRegisterModule {
	return &EventRegisterModule{}
}

// Load loads the event.register module into Lua state
func (m *EventRegisterModule) Load(L *lua.LState) int {
	// Create event.register table
	registerTable := L.NewTable()

	// Register functions
	L.SetField(registerTable, "file", L.NewFunction(m.registerFileWatcher))
	L.SetField(registerTable, "process", L.NewFunction(m.registerProcessWatcher))
	L.SetField(registerTable, "port", L.NewFunction(m.registerPortWatcher))
	L.SetField(registerTable, "service", L.NewFunction(m.registerServiceWatcher))
	L.SetField(registerTable, "cpu", L.NewFunction(m.registerCPUWatcher))
	L.SetField(registerTable, "memory", L.NewFunction(m.registerMemoryWatcher))
	L.SetField(registerTable, "custom", L.NewFunction(m.registerCustomWatcher))

	L.Push(registerTable)
	return 1
}

// registerFileWatcher registers a file watcher on the current agent
// Usage: watcher.register.file({when = {'changed', 'created', 'deleted'}, file_path = '/path/to/file', check_hash = true, interval = '5s'})
// The watcher will be registered on the agent where this code executes (local or delegate_to)
func (m *EventRegisterModule) registerFileWatcher(L *lua.LState) int {
	config := L.CheckTable(1)

	// Extract configuration
	watcherConfig := map[string]interface{}{
		"id":   uuid.New().String(),
		"type": "file",
	}

	// Get file_path (required)
	filePath := L.GetField(config, "file_path")
	if filePath == lua.LNil {
		L.RaiseError("file_path is required for file watcher")
		return 0
	}
	watcherConfig["file_path"] = filePath.String()

	// Get when conditions (required)
	when := L.GetField(config, "when")
	if when == lua.LNil {
		L.RaiseError("when is required (e.g., {'changed', 'created', 'deleted'})")
		return 0
	}

	conditions := []string{}
	if whenTable, ok := when.(*lua.LTable); ok {
		whenTable.ForEach(func(_, value lua.LValue) {
			conditions = append(conditions, value.String())
		})
	} else {
		conditions = append(conditions, when.String())
	}
	watcherConfig["conditions"] = conditions

	// Optional: check_hash
	if checkHash := L.GetField(config, "check_hash"); checkHash != lua.LNil {
		watcherConfig["check_hash"] = lua.LVAsBool(checkHash)
	}

	// Optional: recursive
	if recursive := L.GetField(config, "recursive"); recursive != lua.LNil {
		watcherConfig["recursive"] = lua.LVAsBool(recursive)
	}

	// Optional: interval
	if interval := L.GetField(config, "interval"); interval != lua.LNil {
		watcherConfig["interval"] = interval.String()
	}

	// Store in global registry
	// Watchers are automatically registered on the agent where this code executes
	m.storeWatcher(L, watcherConfig)

	// Return watcher ID
	L.Push(lua.LString(watcherConfig["id"].(string)))
	return 1
}

// registerProcessWatcher registers a process watcher
// Usage: event.register.process({when = {'created', 'deleted'}, process_name = 'nginx', interval = '10s'})
func (m *EventRegisterModule) registerProcessWatcher(L *lua.LState) int {
	config := L.CheckTable(1)

	watcherConfig := map[string]interface{}{
		"id":   uuid.New().String(),
		"type": "process",
	}

	// Get process_name (required)
	processName := L.GetField(config, "process_name")
	if processName == lua.LNil {
		L.RaiseError("process_name is required for process watcher")
		return 0
	}
	watcherConfig["process_name"] = processName.String()

	// Get when conditions (required)
	when := L.GetField(config, "when")
	if when == lua.LNil {
		L.RaiseError("when is required (e.g., {'created', 'deleted'})")
		return 0
	}

	conditions := []string{}
	if whenTable, ok := when.(*lua.LTable); ok {
		whenTable.ForEach(func(_, value lua.LValue) {
			conditions = append(conditions, value.String())
		})
	} else {
		conditions = append(conditions, when.String())
	}
	watcherConfig["conditions"] = conditions

	// Optional: interval
	if interval := L.GetField(config, "interval"); interval != lua.LNil {
		watcherConfig["interval"] = interval.String()
	}

	m.storeWatcher(L, watcherConfig)

	L.Push(lua.LString(watcherConfig["id"].(string)))
	return 1
}

// registerPortWatcher registers a port watcher
// Usage: event.register.port({when = {'created', 'deleted'}, port = 8080, interval = '5s'})
func (m *EventRegisterModule) registerPortWatcher(L *lua.LState) int {
	config := L.CheckTable(1)

	watcherConfig := map[string]interface{}{
		"id":   uuid.New().String(),
		"type": "port",
	}

	// Get port (required)
	port := L.GetField(config, "port")
	if port == lua.LNil {
		L.RaiseError("port is required for port watcher")
		return 0
	}
	watcherConfig["port"] = int(port.(lua.LNumber))

	// Get when conditions (required)
	when := L.GetField(config, "when")
	if when == lua.LNil {
		L.RaiseError("when is required (e.g., {'created', 'deleted'})")
		return 0
	}

	conditions := []string{}
	if whenTable, ok := when.(*lua.LTable); ok {
		whenTable.ForEach(func(_, value lua.LValue) {
			conditions = append(conditions, value.String())
		})
	} else {
		conditions = append(conditions, when.String())
	}
	watcherConfig["conditions"] = conditions

	// Optional: interval
	if interval := L.GetField(config, "interval"); interval != lua.LNil {
		watcherConfig["interval"] = interval.String()
	}

	m.storeWatcher(L, watcherConfig)

	L.Push(lua.LString(watcherConfig["id"].(string)))
	return 1
}

// registerServiceWatcher registers a service watcher
// Usage: event.register.service({when = {'changed'}, service_name = 'nginx', interval = '10s'})
func (m *EventRegisterModule) registerServiceWatcher(L *lua.LState) int {
	config := L.CheckTable(1)

	watcherConfig := map[string]interface{}{
		"id":   uuid.New().String(),
		"type": "service",
	}

	// Get service_name (required)
	serviceName := L.GetField(config, "service_name")
	if serviceName == lua.LNil {
		L.RaiseError("service_name is required for service watcher")
		return 0
	}
	watcherConfig["service_name"] = serviceName.String()

	// Get when conditions (required)
	when := L.GetField(config, "when")
	if when == lua.LNil {
		L.RaiseError("when is required (e.g., {'changed'})")
		return 0
	}

	conditions := []string{}
	if whenTable, ok := when.(*lua.LTable); ok {
		whenTable.ForEach(func(_, value lua.LValue) {
			conditions = append(conditions, value.String())
		})
	} else {
		conditions = append(conditions, when.String())
	}
	watcherConfig["conditions"] = conditions

	// Optional: interval
	if interval := L.GetField(config, "interval"); interval != lua.LNil {
		watcherConfig["interval"] = interval.String()
	}

	m.storeWatcher(L, watcherConfig)

	L.Push(lua.LString(watcherConfig["id"].(string)))
	return 1
}

// registerCPUWatcher registers a CPU usage watcher on the current agent
// Usage: watcher.register.cpu({threshold = 80, when = {'above'}, interval = '5s'})
// The watcher will be registered on the agent where this code executes (local or delegate_to)
func (m *EventRegisterModule) registerCPUWatcher(L *lua.LState) int {
	config := L.CheckTable(1)

	watcherConfig := map[string]interface{}{
		"id":   uuid.New().String(),
		"type": "cpu",
	}

	// Get threshold (required)
	threshold := L.GetField(config, "threshold")
	if threshold == lua.LNil {
		L.RaiseError("threshold is required for CPU watcher")
		return 0
	}
	watcherConfig["cpu_threshold"] = float64(threshold.(lua.LNumber))

	// Get when conditions (required)
	when := L.GetField(config, "when")
	if when == lua.LNil {
		L.RaiseError("when is required (e.g., {'above', 'below'})")
		return 0
	}

	conditions := []string{}
	if whenTable, ok := when.(*lua.LTable); ok {
		whenTable.ForEach(func(_, value lua.LValue) {
			conditions = append(conditions, value.String())
		})
	} else {
		conditions = append(conditions, when.String())
	}
	watcherConfig["conditions"] = conditions

	// Extract common fields (interval, agent)
	m.extractCommonFields(L, config, watcherConfig)

	m.storeWatcher(L, watcherConfig)

	L.Push(lua.LString(watcherConfig["id"].(string)))
	return 1
}

// registerMemoryWatcher registers a memory usage watcher on the current agent
// Usage: watcher.register.memory({threshold = 90, when = {'above'}, interval = '5s'})
// The watcher will be registered on the agent where this code executes (local or delegate_to)
func (m *EventRegisterModule) registerMemoryWatcher(L *lua.LState) int {
	config := L.CheckTable(1)

	watcherConfig := map[string]interface{}{
		"id":   uuid.New().String(),
		"type": "memory",
	}

	// Get threshold (required)
	threshold := L.GetField(config, "threshold")
	if threshold == lua.LNil {
		L.RaiseError("threshold is required for memory watcher")
		return 0
	}
	watcherConfig["memory_threshold"] = float64(threshold.(lua.LNumber))

	// Get when conditions (required)
	when := L.GetField(config, "when")
	if when == lua.LNil {
		L.RaiseError("when is required (e.g., {'above', 'below'})")
		return 0
	}

	conditions := []string{}
	if whenTable, ok := when.(*lua.LTable); ok {
		whenTable.ForEach(func(_, value lua.LValue) {
			conditions = append(conditions, value.String())
		})
	} else {
		conditions = append(conditions, when.String())
	}
	watcherConfig["conditions"] = conditions

	// Extract common fields (interval, agent)
	m.extractCommonFields(L, config, watcherConfig)

	m.storeWatcher(L, watcherConfig)

	L.Push(lua.LString(watcherConfig["id"].(string)))
	return 1
}

// registerCustomWatcher registers a custom watcher with a Lua callback
// Usage: event.register.custom({when = 'check', check = function() return true, {data} end, interval = '30s'})
func (m *EventRegisterModule) registerCustomWatcher(L *lua.LState) int {
	config := L.CheckTable(1)

	watcherConfig := map[string]interface{}{
		"id":   uuid.New().String(),
		"type": "custom",
	}

	// Get check function (required)
	checkFunc := L.GetField(config, "check")
	if checkFunc == lua.LNil || checkFunc.Type() != lua.LTFunction {
		L.RaiseError("check function is required for custom watcher")
		return 0
	}
	// Note: Custom watcher with Lua function callbacks not yet supported
	// Will be implemented in future version
	watcherConfig["check_func"] = "lua_function"

	// Optional: conditions
	when := L.GetField(config, "when")
	conditions := []string{"check"}
	if when != lua.LNil {
		if whenTable, ok := when.(*lua.LTable); ok {
			conditions = []string{}
			whenTable.ForEach(func(_, value lua.LValue) {
				conditions = append(conditions, value.String())
			})
		} else {
			conditions = []string{when.String()}
		}
	}
	watcherConfig["conditions"] = conditions

	// Optional: interval
	if interval := L.GetField(config, "interval"); interval != lua.LNil {
		watcherConfig["interval"] = interval.String()
	}

	m.storeWatcher(L, watcherConfig)

	L.Push(lua.LString(watcherConfig["id"].(string)))
	return 1
}

// extractCommonFields extracts common fields from Lua config table
func (m *EventRegisterModule) extractCommonFields(L *lua.LState, config *lua.LTable, watcherConfig map[string]interface{}) {
	// Optional: interval
	if interval := L.GetField(config, "interval"); interval != lua.LNil {
		watcherConfig["interval"] = interval.String()
	}

	// Optional: agent (for remote registration)
	if agent := L.GetField(config, "agent"); agent != lua.LNil {
		watcherConfig["agent"] = agent.String()
	}
}

// storeWatcher stores watcher config in Lua global registry AND registers with EventWatcherManager
func (m *EventRegisterModule) storeWatcher(L *lua.LState, config map[string]interface{}) {
	// Always store locally in _WATCHERS for backward compatibility
	watchersTable := L.GetGlobal("_WATCHERS")
	if watchersTable == lua.LNil {
		watchersTable = L.NewTable()
		L.SetGlobal("_WATCHERS", watchersTable)
	}

	// Convert config to JSON for storage
	configJSON, err := json.Marshal(config)
	if err != nil {
		const errMsg = "failed to marshal watcher config: %v"
		L.RaiseError(errMsg, err)
		return
	}

	// Store in table with ID as key
	L.SetField(watchersTable.(*lua.LTable), config["id"].(string), lua.LString(configJSON))

	// NEW: Register directly with EventWatcherManager if available
	// This ensures watchers persist with infinite lifecycle
	watcherManagerUD := L.GetGlobal("__WATCHER_MANAGER__")
	if watcherManagerUD.Type() == lua.LTUserData {
		if mgr, ok := watcherManagerUD.(*lua.LUserData).Value.(*agent.EventWatcherManager); ok {
			// Convert map config to agent.WatcherConfig struct
			watcherConfig := &agent.WatcherConfig{
				ID:   config["id"].(string),
				Type: agent.WatcherType(config["type"].(string)),
			}

			// Parse interval
			if interval, ok := config["interval"].(string); ok {
				duration, err := time.ParseDuration(interval)
				if err == nil {
					watcherConfig.Interval = duration
				}
			}

			// Parse conditions (when)
			if when, ok := config["when"].([]interface{}); ok {
				conditions := make([]agent.EventCondition, 0, len(when))
				for _, e := range when {
					if eventStr, ok := e.(string); ok {
						conditions = append(conditions, agent.EventCondition(eventStr))
					}
				}
				watcherConfig.Conditions = conditions
			}

			// Copy type-specific fields
			switch watcherConfig.Type {
			case agent.WatcherTypeFile, agent.WatcherTypeDirectory:
				if path, ok := config["path"].(string); ok {
					watcherConfig.FilePath = path
				}
				if checkHash, ok := config["check_hash"].(bool); ok {
					watcherConfig.CheckHash = checkHash
				}

			case agent.WatcherTypeProcess:
				if processName, ok := config["process_name"].(string); ok {
					watcherConfig.ProcessName = processName
				}

			case agent.WatcherTypePort:
				if port, ok := config["port"].(int); ok {
					watcherConfig.Port = port
				}
				if protocol, ok := config["protocol"].(string); ok {
					watcherConfig.Protocol = protocol
				}

			case agent.WatcherTypeService:
				if serviceName, ok := config["service_name"].(string); ok {
					watcherConfig.ServiceName = serviceName
				}

			case agent.WatcherTypeCPU:
				if threshold, ok := config["threshold"].(float64); ok {
					watcherConfig.CPUThreshold = threshold
				}

			case agent.WatcherTypeMemory:
				if threshold, ok := config["threshold"].(float64); ok {
					watcherConfig.MemoryThreshold = threshold
				}
			}

			// Register with manager (infinite lifecycle)
			slog.Info("Registering watcher with EventWatcherManager",
				"watcher_id", watcherConfig.ID,
				"type", watcherConfig.Type,
				"conditions", watcherConfig.Conditions)

			if err := mgr.RegisterWatcher(*watcherConfig); err != nil {
				slog.Error("Failed to register watcher", "error", err, "watcher_id", watcherConfig.ID)
				const errMsg = "failed to register watcher with manager: %v"
				L.RaiseError(errMsg, err)
			} else {
				slog.Info("Watcher successfully registered with infinite lifecycle", "watcher_id", watcherConfig.ID)
			}
		}
	}
}

