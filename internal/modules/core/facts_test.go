package core

import (
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/internal/agent"
	lua "github.com/yuin/gopher-lua"
)

func TestFactsModule_Register(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	// Verify module is registered
	lv := L.GetGlobal("facts")
	if lv.Type() != lua.LTTable {
		t.Errorf("Expected facts to be a table, got %s", lv.Type())
	}

	table := lv.(*lua.LTable)
	
	// Verify all functions are registered
	functions := []string{
		"get_all", "get_hostname", "get_platform", "get_memory",
		"get_disk", "get_network", "get_packages", "get_package",
		"get_services", "get_service", "get_users", "get_user",
		"get_processes", "get_mounts", "get_uptime", "get_load",
		"get_kernel", "query",
	}

	for _, fn := range functions {
		lv := table.RawGetString(fn)
		if lv.Type() != lua.LTFunction {
			t.Errorf("Expected %s to be a function, got %s", fn, lv.Type())
		}
	}
}

func TestFactsModule_GetHostname_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	// Call without agent parameter
	err := L.DoString(`
		local result, err = facts.get_hostname({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetPlatform_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_platform({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetMemory_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_memory({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetDisk_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_disk({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetNetwork_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_network({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetPackages_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_packages({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetPackage_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_package({ name = "nginx" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetPackage_MissingName(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_package({ agent = "test-agent" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "package name is required") then
			error("Expected 'package name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetServices_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_services({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetService_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_service({ name = "nginx" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetService_MissingName(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_service({ agent = "test-agent" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "service name is required") then
			error("Expected 'service name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetUsers_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_users({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetUser_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_user({ username = "root" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetUser_MissingUsername(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_user({ agent = "test-agent" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "username is required") then
			error("Expected 'username is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetProcesses_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_processes({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetMounts_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_mounts({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetUptime_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_uptime({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetLoad_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_load({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_GetKernel_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.get_kernel({})
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_Query_MissingAgent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.query({ path = "$.memory.total" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "agent name is required") then
			error("Expected 'agent name is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_Query_MissingPath(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	module.Register(L)

	err := L.DoString(`
		local result, err = facts.query({ agent = "test-agent" })
		if result then
			error("Expected nil result")
		end
		if not err or not string.find(err, "query path is required") then
			error("Expected 'query path is required' error")
		end
	`)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFactsModule_DiskToLuaTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	
	disk := &agent.DiskInfo{
		Device:      "/dev/sda1",
		Mountpoint:  "/",
		Fstype:      "ext4",
		Total:       1000000,
		Used:        500000,
		Free:        500000,
		UsedPercent: 50.0,
	}

	table := module.diskToLuaTable(L, disk)
	
	if table.Type() != lua.LTTable {
		t.Errorf("Expected table, got %s", table.Type())
	}

	// Verify fields
	device := table.RawGetString("device")
	if str, ok := device.(lua.LString); !ok || string(str) != "/dev/sda1" {
		t.Errorf("Expected device '/dev/sda1', got %v", device)
	}
}

func TestFactsModule_NetworkToLuaTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFactsModule("localhost:50053")
	
	network := &agent.NetworkInfo{
		Name:      "eth0",
		MAC:       "00:11:22:33:44:55",
		MTU:       1500,
		IsUp:      true,
		Speed:     1000,
		Addresses: []string{"192.168.1.10", "fe80::1"},
	}

	table := module.networkToLuaTable(L, network)
	
	if table.Type() != lua.LTTable {
		t.Errorf("Expected table, got %s", table.Type())
	}

	// Verify fields
	name := table.RawGetString("name")
	if str, ok := name.(lua.LString); !ok || string(str) != "eth0" {
		t.Errorf("Expected name 'eth0', got %v", name)
	}

	isUp := table.RawGetString("is_up")
	if b, ok := isUp.(lua.LBool); !ok || !bool(b) {
		t.Errorf("Expected is_up true, got %v", isUp)
	}
}

func TestFactsModule_GetStringField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("key1", lua.LString("value1"))
	table.RawSetString("key2", lua.LNumber(123))

	// Test existing string key
	result := getStringField(L, table, "key1", "default")
	if result != "value1" {
		t.Errorf("Expected 'value1', got '%s'", result)
	}

	// Test non-existent key (should return default)
	result = getStringField(L, table, "key3", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}

	// Test non-string value (should return default)
	result = getStringField(L, table, "key2", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}

// Mock tests for integration (would need mock gRPC client)
func TestFactsModule_Integration_Placeholder(t *testing.T) {
	// This is a placeholder for integration tests
	// In a real scenario, you would:
	// 1. Create a mock gRPC client
	// 2. Set up expected responses
	// 3. Test the full flow of each function
	
	t.Skip("Integration tests require mock gRPC client setup")
}
