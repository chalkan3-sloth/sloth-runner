package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestNewSystemModule(t *testing.T) {
	module := NewSystemModule()
	if module == nil {
		t.Fatal("NewSystemModule returned nil")
	}
}

func TestRegisterSystemModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Check if system table is registered
	systemTable := L.GetGlobal("system")
	if systemTable == lua.LNil {
		t.Fatal("system table not registered")
	}

	// Check if system is a table
	if systemTable.Type() != lua.LTTable {
		t.Errorf("system should be a table, got %v", systemTable.Type())
	}
}

func TestSystemCPUCount(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Execute Lua code to test cpu_count
	err := L.DoString(`
		count = system.cpu_count()
		assert(count > 0, "CPU count should be greater than 0")
	`)
	if err != nil {
		t.Fatalf("Error executing Lua code: %v", err)
	}
}

func TestSystemMemoryInfo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Execute Lua code to test memory_info
	err := L.DoString(`
		info, err = system.memory_info()
		assert(info ~= nil, "memory_info should return a value")
		if info then
			assert(info.total ~= nil, "memory_info should have total field")
		end
	`)
	if err != nil {
		t.Fatalf("Error executing Lua code: %v", err)
	}
}

func TestSystemHostInfo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Execute Lua code to test host_info
	err := L.DoString(`
		info, err = system.host_info()
		assert(info ~= nil, "host_info should return a value")
		if info then
			assert(info.os ~= nil, "host_info should have os field")
			assert(info.hostname ~= nil, "host_info should have hostname field")
		end
	`)
	if err != nil {
		t.Fatalf("Error executing Lua code: %v", err)
	}
}

func TestSystemUptime(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Execute Lua code to test uptime
	err := L.DoString(`
		uptime, err = system.uptime()
		assert(uptime ~= nil, "uptime should return a value")
		-- uptime may return a table or number depending on implementation
		if type(uptime) == "number" then
			assert(uptime > 0, "uptime should be greater than 0")
		end
	`)
	if err != nil {
		t.Fatalf("Error executing Lua code: %v", err)
	}
}

func TestSystemEnvironment(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Execute Lua code to test environment
	err := L.DoString(`
		env = system.environment()
		assert(env ~= nil, "environment should return a value")
		assert(type(env) == "table", "environment should return a table")
	`)
	if err != nil {
		t.Fatalf("Error executing Lua code: %v", err)
	}
}

func TestSystemDiskUsage(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Execute Lua code to test disk_usage
	err := L.DoString(`
		usage, err = system.disk_usage("/")
		-- Just check it doesn't crash
		-- Results may vary by platform
	`)
	if err != nil {
		t.Fatalf("Error executing Lua code: %v", err)
	}
}

func TestSystemPerformanceSnapshot(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	// Execute Lua code to test performance_snapshot
	err := L.DoString(`
		snapshot, err = system.performance_snapshot()
		-- Just check it doesn't crash
	`)
	if err != nil {
		t.Fatalf("Error executing Lua code: %v", err)
	}
}

func TestSystemModuleFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSystemModule(L)

	functions := []string{
		"cpu_info",
		"cpu_usage",
		"cpu_count",
		"load_average",
		"memory_info",
		"memory_usage",
		"swap_info",
		"disk_usage",
		"disk_io",
		"disk_partitions",
		"network_interfaces",
		"network_stats",
		"network_connections",
		"processes",
		"process_info",
		"kill_process",
		"host_info",
		"uptime",
		"environment",
		"users",
		"performance_snapshot",
		"system_health",
	}

	systemTable := L.GetGlobal("system")
	if systemTable.Type() != lua.LTTable {
		t.Fatal("system should be a table")
	}

	table := systemTable.(*lua.LTable)
	for _, funcName := range functions {
		fn := L.GetField(table, funcName)
		if fn.Type() != lua.LTFunction {
			t.Errorf("system.%s should be a function, got %v", funcName, fn.Type())
		}
	}
}
