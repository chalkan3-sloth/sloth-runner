package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestContextModule(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{
			name: "context.get_value",
			script: `
				context.set_value("test_key", "test_value")
				local val = context.get_value("test_key")
				assert(val == "test_value", "context.get_value should return the set value")
			`,
		},
		{
			name: "context.set_value",
			script: `
				context.set_value("key1", "value1")
				context.set_value("key2", 123)
				context.set_value("key3", true)
				local v1 = context.get_value("key1")
				local v2 = context.get_value("key2")
				local v3 = context.get_value("key3")
				assert(v1 == "value1", "string value should match")
				assert(v2 == 123, "number value should match")
				assert(v3 == true, "boolean value should match")
			`,
		},
		{
			name: "context.has_value",
			script: `
				context.set_value("exists", "yes")
				assert(context.has_value("exists") == true, "should return true for existing key")
				assert(context.has_value("not_exists") == false, "should return false for non-existing key")
			`,
		},
		{
			name: "context.delete_value",
			script: `
				context.set_value("to_delete", "value")
				assert(context.has_value("to_delete") == true, "key should exist")
				context.delete_value("to_delete")
				assert(context.has_value("to_delete") == false, "key should be deleted")
			`,
		},
		{
			name: "context.clear",
			script: `
				context.set_value("key1", "val1")
				context.set_value("key2", "val2")
				context.clear()
				assert(context.has_value("key1") == false, "all keys should be cleared")
				assert(context.has_value("key2") == false, "all keys should be cleared")
			`,
		},
		{
			name: "context.keys",
			script: `
				context.clear()
				context.set_value("k1", "v1")
				context.set_value("k2", "v2")
				local keys = context.keys()
				assert(type(keys) == "table", "keys should return a table")
				assert(#keys == 2, "should have 2 keys")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			loadContextModule(L)

			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestContextTableValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	loadContextModule(L)

	script := `
		context.set_value("table_key", {a = 1, b = 2})
		local val = context.get_value("table_key")
		assert(type(val) == "table", "should return a table")
		assert(val.a == 1, "table value should be preserved")
		assert(val.b == 2, "table value should be preserved")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
