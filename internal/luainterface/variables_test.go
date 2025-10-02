package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestVariablesModule(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{
			name: "variables basic set and get",
			script: `
				variables.set("key1", "value1")
				local val = variables.get("key1")
				assert(val == "value1", "should get the set value")
			`,
		},
		{
			name: "variables get with default",
			script: `
				local val = variables.get("nonexistent", "default")
				assert(val == "default", "should return default value")
			`,
		},
		{
			name: "variables set multiple types",
			script: `
				variables.set("str", "string")
				variables.set("num", 42)
				variables.set("bool", true)
				assert(variables.get("str") == "string", "string value")
				assert(variables.get("num") == 42, "number value")
				assert(variables.get("bool") == true, "boolean value")
			`,
		},
		{
			name: "variables has",
			script: `
				variables.set("exists", "yes")
				assert(variables.has("exists") == true, "should return true for existing key")
				assert(variables.has("not_exists") == false, "should return false for non-existing key")
			`,
		},
		{
			name: "variables delete",
			script: `
				variables.set("to_delete", "value")
				assert(variables.has("to_delete") == true, "key should exist")
				variables.delete("to_delete")
				assert(variables.has("to_delete") == false, "key should be deleted")
			`,
		},
		{
			name: "variables clear",
			script: `
				variables.set("k1", "v1")
				variables.set("k2", "v2")
				variables.clear()
				assert(variables.has("k1") == false, "all keys should be cleared")
				assert(variables.has("k2") == false, "all keys should be cleared")
			`,
		},
		{
			name: "variables keys",
			script: `
				variables.clear()
				variables.set("key1", "val1")
				variables.set("key2", "val2")
				local keys = variables.keys()
				assert(type(keys) == "table", "keys should return a table")
			`,
		},
		{
			name: "variables all",
			script: `
				variables.clear()
				variables.set("a", 1)
				variables.set("b", 2)
				local all = variables.all()
				assert(type(all) == "table", "all should return a table")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			loadVariablesModule(L)

			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestVariablesTableValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	loadVariablesModule(L)

	script := `
		variables.set("table_var", {x = 10, y = 20})
		local val = variables.get("table_var")
		assert(type(val) == "table", "should return a table")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
