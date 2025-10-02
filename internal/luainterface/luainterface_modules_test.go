package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLogModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "log.info",
			script: `
log.info("Test info message")
`,
		},
		{
			name: "log.error",
			script: `
log.error("Test error message")
`,
		},
		{
			name: "log.warn",
			script: `
log.warn("Test warning message")
`,
		},
		{
			name: "log.debug",
			script: `
log.debug("Test debug message")
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
			}
		})
	}
}

func TestEnvModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState) error
	}{
		{
			name: "env.get",
			script: `
result = env.get("PATH")
`,
			check: func(L *lua.LState) error {
				result := L.GetGlobal("result")
				if result.Type() == lua.LTNil {
					t.Error("env.get returned nil for PATH")
				}
				return nil
			},
		},
		{
			name: "env.set",
			script: `
env.set("TEST_VAR", "test_value")
result = env.get("TEST_VAR")
`,
			check: func(L *lua.LState) error {
				result := L.GetGlobal("result")
				if result.String() != "test_value" {
					t.Errorf("Expected 'test_value', got: %s", result.String())
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestJsonModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "json.encode",
			script: `
local data = {name = "test", value = 123}
result = json.encode(data)
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTString {
					t.Error("json.encode did not return a string")
				}
			},
		},
		{
			name: "json.decode",
			script: `
local json_str = '{"name":"test","value":123}'
result = json.decode(json_str)
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTTable {
					t.Error("json.decode did not return a table")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestYamlModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "yaml.encode",
			script: `
local data = {name = "test", value = 123}
result = yaml.encode(data)
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTString {
					t.Error("yaml.encode did not return a string")
				}
			},
		},
		{
			name: "yaml.decode",
			script: `
local yaml_str = "name: test\nvalue: 123"
result = yaml.decode(yaml_str)
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTTable {
					t.Error("yaml.decode did not return a table")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestTemplateModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
local tmpl = "Hello, {{.name}}!"
local data = {name = "World"}
result = template.render(tmpl, data)
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute template.render: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got: %s", result.String())
	}
}

func TestMathModuleFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "math.round",
			script: `
result = math.round(3.7)
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.(lua.LNumber) != 4 {
					t.Errorf("Expected 4, got: %v", result)
				}
			},
		},
		{
			name: "math.max",
			script: `
result = math.max(1, 5, 3, 2)
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.(lua.LNumber) != 5 {
					t.Errorf("Expected 5, got: %v", result)
				}
			},
		},
		{
			name: "math.min",
			script: `
result = math.min(1, 5, 3, 2)
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.(lua.LNumber) != 1 {
					t.Errorf("Expected 1, got: %v", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestStringsModuleFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "strings.split",
			script: `
result = strings.split("a,b,c", ",")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTTable {
					t.Error("strings.split did not return a table")
				}
			},
		},
		{
			name: "strings.join",
			script: `
local parts = {"a", "b", "c"}
result = strings.join(parts, ",")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "a,b,c" {
					t.Errorf("Expected 'a,b,c', got: %s", result.String())
				}
			},
		},
		{
			name: "strings.upper",
			script: `
result = strings.upper("hello")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "HELLO" {
					t.Errorf("Expected 'HELLO', got: %s", result.String())
				}
			},
		},
		{
			name: "strings.lower",
			script: `
result = strings.lower("HELLO")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "hello" {
					t.Errorf("Expected 'hello', got: %s", result.String())
				}
			},
		},
		{
			name: "strings.trim",
			script: `
result = strings.trim("  hello  ")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "hello" {
					t.Errorf("Expected 'hello', got: %s", result.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestCryptoModuleBasics(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "crypto.sha256",
			script: `
result = crypto.sha256("hello world")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTString {
					t.Error("crypto.sha256 did not return a string")
				}
			},
		},
		{
			name: "crypto.md5",
			script: `
result = crypto.md5("hello world")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTString {
					t.Error("crypto.md5 did not return a string")
				}
			},
		},
		{
			name: "crypto.base64_encode",
			script: `
result = crypto.base64_encode("hello")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "aGVsbG8=" {
					t.Errorf("Expected 'aGVsbG8=', got: %s", result.String())
				}
			},
		},
		{
			name: "crypto.base64_decode",
			script: `
result = crypto.base64_decode("aGVsbG8=")
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.String() != "hello" {
					t.Errorf("Expected 'hello', got: %s", result.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestTimeModuleBasics(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "time.now",
			script: `
result = time.now()
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() == lua.LTNil {
					t.Error("time.now returned nil")
				}
			},
		},
		{
			name: "time.unix",
			script: `
result = time.unix()
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTNumber {
					t.Error("time.unix did not return a number")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestSystemModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "sys.os",
			script: `
result = sys.os()
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTString {
					t.Error("sys.os did not return a string")
				}
			},
		},
		{
			name: "sys.arch",
			script: `
result = sys.arch()
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTString {
					t.Error("sys.arch did not return a string")
				}
			},
		},
		{
			name: "sys.hostname",
			script: `
result = sys.hostname()
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTString {
					t.Error("sys.hostname did not return a string")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Errorf("Failed to execute %s: %v", tt.name, err)
				return
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestModuleInteraction(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	// Test interaction between modules
	script := `
local data = {name = "test", value = 123}
local json_str = json.encode(data)
local decoded = json.decode(json_str)
result = decoded.name
`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.GetGlobal("result")
	if result.String() != "test" {
		t.Errorf("Expected 'test', got: %s", result.String())
	}
}

func TestErrorHandling(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name: "json.decode invalid",
			script: `
result = json.decode("invalid json")
`,
			wantErr: true,
		},
		{
			name: "yaml.decode invalid",
			script: `
result = yaml.decode("invalid: yaml: data:")
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := L.DoString(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got error=%v", tt.wantErr, err)
			}
		})
	}
}
