package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestEnvModule(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{
			name: "env.get existing",
			script: `
				local path = env.get("PATH")
				assert(path ~= nil, "PATH should exist")
				assert(type(path) == "string", "PATH should be a string")
			`,
		},
		{
			name: "env.get with default",
			script: `
				local val = env.get("NONEXISTENT_VAR_XYZ", "default_value")
				assert(val == "default_value", "should return default value")
			`,
		},
		{
			name: "env.set and get",
			script: `
				env.set("TEST_VAR_123", "test_value")
				local val = env.get("TEST_VAR_123")
				assert(val == "test_value", "should get the set value")
			`,
		},
		{
			name: "env.unset",
			script: `
				env.set("TO_UNSET", "value")
				local before = env.get("TO_UNSET")
				assert(before == "value", "value should be set")
				env.unset("TO_UNSET")
				local after = env.get("TO_UNSET")
				assert(after == nil or after == "", "value should be unset")
			`,
		},
		{
			name: "env.list",
			script: `
				local envs = env.list()
				assert(type(envs) == "table", "should return a table")
				assert(envs["PATH"] ~= nil, "should contain PATH")
			`,
		},
		{
			name: "env.expand simple",
			script: `
				env.set("MYVAR", "world")
				local result = env.expand("hello $MYVAR")
				assert(result == "hello world", "should expand variable")
			`,
		},
		{
			name: "env.expand braces",
			script: `
				env.set("VAR1", "foo")
				local result = env.expand("prefix ${VAR1} suffix")
				assert(result == "prefix foo suffix", "should expand variable with braces")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			loadEnvModule(L)

			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestEnvModuleComplexExpansion(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	loadEnvModule(L)

	script := `
		env.set("A", "value_a")
		env.set("B", "value_b")
		local result = env.expand("$A and ${B} together")
		assert(type(result) == "string", "should return a string")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
