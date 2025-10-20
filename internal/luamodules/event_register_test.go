package luamodules

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// Test EventRegisterModule creation
func TestNewEventRegisterModule(t *testing.T) {
	module := NewEventRegisterModule()

	if module == nil {
		t.Error("Expected non-nil module")
	}
}

func TestNewEventRegisterModule_MultipleInstances(t *testing.T) {
	module1 := NewEventRegisterModule()
	module2 := NewEventRegisterModule()

	if module1 == nil || module2 == nil {
		t.Error("Expected both modules to be created")
	}

	// Modules should be separate instances
	if module1 == module2 {
		t.Error("Expected different instances")
	}
}

// Test Load method
func TestEventRegisterModule_Load(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()
	result := module.Load(L)

	if result != 1 {
		t.Errorf("Expected Load to return 1, got %d", result)
	}

	// Check that a table was pushed to the stack
	stackTop := L.GetTop()
	if stackTop == 0 {
		t.Error("Expected table on stack")
	}

	value := L.Get(-1)
	if value.Type() != lua.LTTable {
		t.Errorf("Expected table, got %s", value.Type())
	}
}

func TestEventRegisterModule_LoadFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()
	module.Load(L)

	table := L.Get(-1).(*lua.LTable)

	expectedFunctions := []string{
		"file",
		"process",
		"port",
		"service",
		"cpu",
		"memory",
		"custom",
	}

	for _, funcName := range expectedFunctions {
		value := L.GetField(table, funcName)
		if value.Type() != lua.LTFunction {
			t.Errorf("Expected %s to be a function, got %s", funcName, value.Type())
		}
	}
}

func TestEventRegisterModule_LoadedTableStructure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()
	module.Load(L)

	table := L.Get(-1).(*lua.LTable)

	// Verify all expected functions are present
	functionCount := 0
	table.ForEach(func(key, value lua.LValue) {
		if value.Type() == lua.LTFunction {
			functionCount++
		}
	})

	if functionCount != 7 {
		t.Errorf("Expected 7 functions, got %d", functionCount)
	}
}

// Test AgentResolver interface
func TestAgentResolver_InterfaceExists(t *testing.T) {
	// Verify interface can be referenced
	var _ AgentResolver

	// This test ensures the interface exists and can be referenced
}

type mockAgentResolver struct {
	resolveFunc func(string) (string, error)
}

func (m *mockAgentResolver) ResolveAgent(name string) (string, error) {
	if m.resolveFunc != nil {
		return m.resolveFunc(name)
	}
	return "localhost:50051", nil
}

func TestAgentResolver_MockImplementation(t *testing.T) {
	// MockAgentResolver should implement AgentResolver interface
	var _ AgentResolver = (*mockAgentResolver)(nil)

	mock := &mockAgentResolver{
		resolveFunc: func(name string) (string, error) {
			return "test:1234", nil
		},
	}

	addr, err := mock.ResolveAgent("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if addr != "test:1234" {
		t.Errorf("Expected 'test:1234', got %s", addr)
	}
}

// Test extractCommonFields
func TestEventRegisterModule_ExtractCommonFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()

	config := L.NewTable()
	L.SetField(config, "interval", lua.LString("10s"))
	L.SetField(config, "agent", lua.LString("test-agent"))

	watcherConfig := make(map[string]interface{})
	module.extractCommonFields(L, config, watcherConfig)

	if interval, ok := watcherConfig["interval"].(string); !ok || interval != "10s" {
		t.Error("Expected interval to be extracted")
	}

	if agent, ok := watcherConfig["agent"].(string); !ok || agent != "test-agent" {
		t.Error("Expected agent to be extracted")
	}
}

func TestEventRegisterModule_ExtractCommonFields_OnlyInterval(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()

	config := L.NewTable()
	L.SetField(config, "interval", lua.LString("5s"))

	watcherConfig := make(map[string]interface{})
	module.extractCommonFields(L, config, watcherConfig)

	if interval, ok := watcherConfig["interval"].(string); !ok || interval != "5s" {
		t.Error("Expected interval to be extracted")
	}

	if _, ok := watcherConfig["agent"]; ok {
		t.Error("Expected agent to not be present")
	}
}

func TestEventRegisterModule_ExtractCommonFields_Empty(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()

	config := L.NewTable()
	watcherConfig := make(map[string]interface{})

	module.extractCommonFields(L, config, watcherConfig)

	if len(watcherConfig) != 0 {
		t.Errorf("Expected empty watcherConfig, got %d fields", len(watcherConfig))
	}
}

// Test storeWatcher
func TestEventRegisterModule_StoreWatcher(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()

	watcherConfig := map[string]interface{}{
		"id":   "test-watcher-123",
		"type": "file",
		"path": "/tmp/test",
	}

	module.storeWatcher(L, watcherConfig)

	// Check that _WATCHERS global was created
	watchersTable := L.GetGlobal("_WATCHERS")
	if watchersTable == lua.LNil {
		t.Error("Expected _WATCHERS global to be created")
	}

	// Check that watcher was stored
	if watchersTable.Type() != lua.LTTable {
		t.Error("Expected _WATCHERS to be a table")
	}

	watcherData := L.GetField(watchersTable.(*lua.LTable), "test-watcher-123")
	if watcherData == lua.LNil {
		t.Error("Expected watcher to be stored")
	}
}

func TestEventRegisterModule_StoreWatcher_MultipleWatchers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()

	// Store first watcher
	watcher1 := map[string]interface{}{
		"id":   "watcher-1",
		"type": "cpu",
	}
	module.storeWatcher(L, watcher1)

	// Store second watcher
	watcher2 := map[string]interface{}{
		"id":   "watcher-2",
		"type": "memory",
	}
	module.storeWatcher(L, watcher2)

	// Verify both are stored
	watchersTable := L.GetGlobal("_WATCHERS").(*lua.LTable)

	w1 := L.GetField(watchersTable, "watcher-1")
	w2 := L.GetField(watchersTable, "watcher-2")

	if w1 == lua.LNil || w2 == lua.LNil {
		t.Error("Expected both watchers to be stored")
	}
}

// Test EventRegisterModule struct
func TestEventRegisterModule_StructCreation(t *testing.T) {
	module := &EventRegisterModule{}

	if module == nil {
		t.Error("Expected non-nil module")
	}
}

func TestEventRegisterModule_ZeroValue(t *testing.T) {
	var module EventRegisterModule

	// Zero value should be usable
	L := lua.NewState()
	defer L.Close()

	result := module.Load(L)
	if result != 1 {
		t.Errorf("Expected Load to return 1 with zero value module, got %d", result)
	}
}

// Test module initialization patterns
func TestEventRegisterModule_LoadedMultipleTimes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()

	// Load module multiple times
	result1 := module.Load(L)
	result2 := module.Load(L)
	result3 := module.Load(L)

	if result1 != 1 || result2 != 1 || result3 != 1 {
		t.Error("Expected Load to consistently return 1")
	}

	// Should have 3 tables on stack
	if L.GetTop() != 3 {
		t.Errorf("Expected 3 items on stack, got %d", L.GetTop())
	}
}

func TestEventRegisterModule_ConcurrentLoad(t *testing.T) {
	module := NewEventRegisterModule()

	done := make(chan bool, 10)

	// Load module concurrently in separate Lua states
	for i := 0; i < 10; i++ {
		go func() {
			L := lua.NewState()
			defer L.Close()

			result := module.Load(L)
			if result != 1 {
				t.Errorf("Expected result 1, got %d", result)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test function registration names
func TestEventRegisterModule_FunctionNames(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()
	module.Load(L)

	table := L.Get(-1).(*lua.LTable)

	expectedNames := []string{
		"file",
		"process",
		"port",
		"service",
		"cpu",
		"memory",
		"custom",
	}

	for _, name := range expectedNames {
		value := L.GetField(table, name)
		if value == lua.LNil {
			t.Errorf("Expected function '%s' to be registered", name)
		}
	}
}

func TestEventRegisterModule_NoExtraFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()
	module.Load(L)

	table := L.Get(-1).(*lua.LTable)

	// Count all keys in table
	keyCount := 0
	table.ForEach(func(key, value lua.LValue) {
		keyCount++
	})

	expectedCount := 7 // file, process, port, service, cpu, memory, custom
	if keyCount != expectedCount {
		t.Errorf("Expected exactly %d functions, got %d", expectedCount, keyCount)
	}
}

// Test module state independence
func TestEventRegisterModule_StateIndependence(t *testing.T) {
	L1 := lua.NewState()
	defer L1.Close()

	L2 := lua.NewState()
	defer L2.Close()

	module := NewEventRegisterModule()

	// Load in both states
	module.Load(L1)
	module.Load(L2)

	// Store watcher in L1
	watcher1 := map[string]interface{}{
		"id":   "state1-watcher",
		"type": "cpu",
	}
	module.storeWatcher(L1, watcher1)

	// Verify L2 doesn't have L1's watcher
	watchers2 := L2.GetGlobal("_WATCHERS")
	if watchers2 != lua.LNil {
		table := watchers2.(*lua.LTable)
		value := L2.GetField(table, "state1-watcher")
		if value != lua.LNil {
			t.Error("Expected L2 to not have L1's watcher")
		}
	}
}

func TestEventRegisterModule_EmptyConfigStorage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewEventRegisterModule()

	// Store watcher with minimal config
	minimalConfig := map[string]interface{}{
		"id":   "minimal-watcher",
		"type": "test",
	}

	module.storeWatcher(L, minimalConfig)

	watchersTable := L.GetGlobal("_WATCHERS")
	if watchersTable == lua.LNil {
		t.Error("Expected _WATCHERS to be created")
	}
}

// Test extractCommonFields with various inputs
func TestEventRegisterModule_ExtractCommonFields_VariousIntervals(t *testing.T) {
	testCases := []struct {
		name     string
		interval string
	}{
		{"1 second", "1s"},
		{"30 seconds", "30s"},
		{"1 minute", "1m"},
		{"5 minutes", "5m"},
		{"1 hour", "1h"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			module := NewEventRegisterModule()

			config := L.NewTable()
			L.SetField(config, "interval", lua.LString(tc.interval))

			watcherConfig := make(map[string]interface{})
			module.extractCommonFields(L, config, watcherConfig)

			if interval, ok := watcherConfig["interval"].(string); !ok || interval != tc.interval {
				t.Errorf("Expected interval %s, got %v", tc.interval, watcherConfig["interval"])
			}
		})
	}
}

func TestEventRegisterModule_ExtractCommonFields_AgentNames(t *testing.T) {
	agentNames := []string{
		"local-agent",
		"remote-agent-1",
		"production-server",
		"agent.example.com",
		"192.168.1.100",
	}

	for _, agentName := range agentNames {
		t.Run(agentName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			module := NewEventRegisterModule()

			config := L.NewTable()
			L.SetField(config, "agent", lua.LString(agentName))

			watcherConfig := make(map[string]interface{})
			module.extractCommonFields(L, config, watcherConfig)

			if agent, ok := watcherConfig["agent"].(string); !ok || agent != agentName {
				t.Errorf("Expected agent %s, got %v", agentName, watcherConfig["agent"])
			}
		})
	}
}
