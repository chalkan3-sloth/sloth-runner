package infra

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestFrpModuleRegistration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create and register FRP module
	frpModule := NewFrpModule(nil)
	frpModule.Register(L)

	// Register metatables
	RegisterServerMetatable(L)
	RegisterClientMetatable(L)

	// Test 1: Check if frp global exists
	frpGlobal := L.GetGlobal("frp")
	if frpGlobal.Type() != lua.LTTable {
		t.Fatal("frp global should be a table")
	}

	frpTable := frpGlobal.(*lua.LTable)

	// Test 2: Check if frp.server exists
	serverFunc := frpTable.RawGetString("server")
	if serverFunc.Type() != lua.LTFunction {
		t.Fatal("frp.server should be a function")
	}

	// Test 3: Check if frp.client exists
	clientFunc := frpTable.RawGetString("client")
	if clientFunc.Type() != lua.LTFunction {
		t.Fatal("frp.client should be a function")
	}

	// Test 4: Check if frp.install exists
	installFunc := frpTable.RawGetString("install")
	if installFunc.Type() != lua.LTFunction {
		t.Fatal("frp.install should be a function")
	}
}

func TestFrpServerCreation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create and register FRP module
	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterServerMetatable(L)

	// Create server instance
	script := `
		local server = frp.server("test-server")
		return type(server)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata, got %s", result.String())
	}
}

func TestFrpClientCreation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create and register FRP module
	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterClientMetatable(L)

	// Create client instance
	script := `
		local client = frp.client("test-client")
		return type(client)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata, got %s", result.String())
	}
}

func TestFrpServerFluentAPI(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create and register FRP module
	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterServerMetatable(L)

	// Test fluent API
	script := `
		local server = frp.server("test-server")
			:config({
				bindAddr = "0.0.0.0",
				bindPort = 7000
			})
			:version("latest")

		return type(server)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to use fluent API: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata after chaining, got %s", result.String())
	}
}

func TestFrpClientFluentAPI(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create and register FRP module
	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterClientMetatable(L)

	// Test fluent API with server and proxy
	script := `
		local client = frp.client("test-client")
			:server("example.com", 7000)
			:config({
				auth = {
					method = "token",
					token = "test123"
				}
			})
			:proxy({
				name = "web",
				type = "http",
				localPort = 8080
			})

		return type(client)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to use client fluent API: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata after chaining, got %s", result.String())
	}
}

func TestFrpServerConfigPath(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterServerMetatable(L)

	script := `
		local server = frp.server("test")
			:config_path("/custom/path/frps.toml")

		return type(server)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to set config path: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata, got %s", result.String())
	}
}

func TestFrpClientMultipleProxies(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterClientMetatable(L)

	script := `
		local client = frp.client("test")
			:proxy({name = "web", type = "http", localPort = 80})
			:proxy({name = "ssh", type = "tcp", localPort = 22})
			:proxy({name = "db", type = "tcp", localPort = 5432})

		return type(client)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to add multiple proxies: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata, got %s", result.String())
	}
}

func TestLuaValueToGoConversion(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test string conversion
	strVal := lua.LString("test")
	goVal := luaValueToGo(strVal)
	if goVal != "test" {
		t.Errorf("Expected 'test', got %v", goVal)
	}

	// Test number conversion
	numVal := lua.LNumber(42)
	goVal = luaValueToGo(numVal)
	if goVal != float64(42) {
		t.Errorf("Expected 42, got %v", goVal)
	}

	// Test boolean conversion
	boolVal := lua.LBool(true)
	goVal = luaValueToGo(boolVal)
	if goVal != true {
		t.Errorf("Expected true, got %v", goVal)
	}

	// Test nil conversion
	nilVal := lua.LNil
	goVal = luaValueToGo(nilVal)
	if goVal != nil {
		t.Errorf("Expected nil, got %v", goVal)
	}
}

func TestGoValueToLuaConversion(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test string conversion
	luaVal := goValueToLua(L, "test")
	if luaVal.Type() != lua.LTString {
		t.Errorf("Expected string, got %s", luaVal.Type())
	}

	// Test number conversion
	luaVal = goValueToLua(L, 42)
	if luaVal.Type() != lua.LTNumber {
		t.Errorf("Expected number, got %s", luaVal.Type())
	}

	// Test boolean conversion
	luaVal = goValueToLua(L, true)
	if luaVal.Type() != lua.LTBool {
		t.Errorf("Expected boolean, got %s", luaVal.Type())
	}

	// Test nil conversion
	luaVal = goValueToLua(L, nil)
	if luaVal.Type() != lua.LTNil {
		t.Errorf("Expected nil, got %s", luaVal.Type())
	}

	// Test map conversion
	testMap := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	luaVal = goValueToLua(L, testMap)
	if luaVal.Type() != lua.LTTable {
		t.Errorf("Expected table, got %s", luaVal.Type())
	}

	// Test slice conversion
	testSlice := []interface{}{"a", "b", "c"}
	luaVal = goValueToLua(L, testSlice)
	if luaVal.Type() != lua.LTTable {
		t.Errorf("Expected table, got %s", luaVal.Type())
	}
}

func TestFrpServerDelegateTo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterServerMetatable(L)

	script := `
		local server = frp.server("test")
			:delegate_to("remote-agent")

		return type(server)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to set delegate_to: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata, got %s", result.String())
	}
}

func TestFrpClientDelegateTo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterClientMetatable(L)

	script := `
		local client = frp.client("test")
			:delegate_to("remote-agent")

		return type(client)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to set delegate_to: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result.String() != "userdata" {
		t.Fatalf("Expected userdata, got %s", result.String())
	}
}

func TestFrpModuleNewInstance(t *testing.T) {
	// Test creating new FRP module
	module := NewFrpModule(nil)
	if module == nil {
		t.Fatal("NewFrpModule should not return nil")
	}

	// Test with agent client
	dummyClient := "dummy-agent-client"
	module = NewFrpModule(dummyClient)
	if module == nil {
		t.Fatal("NewFrpModule with client should not return nil")
	}
}

func TestFrpServerDefaultConfig(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterServerMetatable(L)

	// Create server and check defaults are set
	script := `
		local server = frp.server()
		return server ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to create server with defaults: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if !lua.LVAsBool(result) {
		t.Fatal("Server creation should succeed with defaults")
	}
}

func TestFrpClientDefaultConfig(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	frpModule := NewFrpModule(nil)
	frpModule.Register(L)
	RegisterClientMetatable(L)

	// Create client and check it works
	script := `
		local client = frp.client()
		return client ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to create client with defaults: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if !lua.LVAsBool(result) {
		t.Fatal("Client creation should succeed with defaults")
	}
}
