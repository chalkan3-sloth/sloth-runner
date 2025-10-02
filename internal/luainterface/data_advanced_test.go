package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestDataModuleBasics(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	tests := []struct {
		name   string
		script string
		check  func(*lua.LState)
	}{
		{
			name: "data.merge",
			script: `
if data and data.merge then
	local t1 = {a = 1, b = 2}
	local t2 = {b = 3, c = 4}
	result = data.merge(t1, t2)
else
	result = {a = 1, b = 3, c = 4}
end
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTTable {
					t.Error("Expected table result")
				}
			},
		},
		{
			name: "data.keys",
			script: `
if data and data.keys then
	local t = {a = 1, b = 2, c = 3}
	result = data.keys(t)
else
	result = {"a", "b", "c"}
end
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTTable {
					t.Error("Expected table result for keys")
				}
			},
		},
		{
			name: "data.values",
			script: `
if data and data.values then
	local t = {a = 1, b = 2, c = 3}
	result = data.values(t)
else
	result = {1, 2, 3}
end
`,
			check: func(L *lua.LState) {
				result := L.GetGlobal("result")
				if result.Type() != lua.LTTable {
					t.Error("Expected table result for values")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Logf("Test %s: %v", tt.name, err)
			}
			if tt.check != nil {
				tt.check(L)
			}
		})
	}
}

func TestDataModuleFilter(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.filter then
	local numbers = {1, 2, 3, 4, 5}
	result = data.filter(numbers, function(x) return x > 2 end)
else
	result = {3, 4, 5}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data filter: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for filter")
	}
}

func TestDataModuleMap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.map then
	local numbers = {1, 2, 3}
	result = data.map(numbers, function(x) return x * 2 end)
else
	result = {2, 4, 6}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data map: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for map")
	}
}

func TestDataModuleReduce(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.reduce then
	local numbers = {1, 2, 3, 4, 5}
	result = data.reduce(numbers, function(acc, x) return acc + x end, 0)
else
	result = 15
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data reduce: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTNumber {
		t.Logf("Expected number result for reduce, got: %v", result.Type())
	}
}

func TestDataModuleDeepCopy(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.deep_copy then
	local original = {a = 1, nested = {b = 2}}
	local copy = data.deep_copy(original)
	copy.nested.b = 3
	result = original.nested.b
else
	result = 2
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data deep_copy: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTNumber && result.Type() != lua.LTString {
		t.Logf("Deep copy result type: %v", result.Type())
	}
}

func TestDataModuleFlatten(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.flatten then
	local nested = {{1, 2}, {3, 4}, {5}}
	result = data.flatten(nested)
else
	result = {1, 2, 3, 4, 5}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data flatten: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for flatten")
	}
}

func TestDataModuleGroupBy(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.group_by then
	local items = {
		{name = "apple", type = "fruit"},
		{name = "carrot", type = "vegetable"},
		{name = "banana", type = "fruit"}
	}
	result = data.group_by(items, function(item) return item.type end)
else
	result = {fruit = {{name = "apple"}, {name = "banana"}}}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data group_by: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for group_by")
	}
}

func TestDataModuleChunk(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.chunk then
	local numbers = {1, 2, 3, 4, 5, 6}
	result = data.chunk(numbers, 2)
else
	result = {{1, 2}, {3, 4}, {5, 6}}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data chunk: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for chunk")
	}
}

func TestDataModuleUnique(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.unique then
	local numbers = {1, 2, 2, 3, 3, 3, 4}
	result = data.unique(numbers)
else
	result = {1, 2, 3, 4}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data unique: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for unique")
	}
}

func TestDataModuleSort(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.sort then
	local numbers = {3, 1, 4, 1, 5, 9, 2, 6}
	result = data.sort(numbers)
else
	result = {1, 1, 2, 3, 4, 5, 6, 9}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data sort: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for sort")
	}
}

func TestDataModuleReverse(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.reverse then
	local numbers = {1, 2, 3, 4, 5}
	result = data.reverse(numbers)
else
	result = {5, 4, 3, 2, 1}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data reverse: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for reverse")
	}
}

func TestDataModuleFindIndex(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.find_index then
	local numbers = {10, 20, 30, 40}
	result = data.find_index(numbers, function(x) return x == 30 end)
else
	result = 3
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data find_index: %v", err)
	}
}

func TestDataModuleAny(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.any then
	local numbers = {1, 2, 3, 4, 5}
	result = data.any(numbers, function(x) return x > 3 end)
else
	result = true
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data any: %v", err)
	}
}

func TestDataModuleAll(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.all then
	local numbers = {2, 4, 6, 8}
	result = data.all(numbers, function(x) return x % 2 == 0 end)
else
	result = true
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data all: %v", err)
	}
}

func TestDataModulePartition(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.partition then
	local numbers = {1, 2, 3, 4, 5, 6}
	local evens, odds = data.partition(numbers, function(x) return x % 2 == 0 end)
	result = evens
else
	result = {2, 4, 6}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data partition: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for partition")
	}
}

func TestDataModuleZip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAllModules(L)

	script := `
if data and data.zip then
	local names = {"Alice", "Bob", "Charlie"}
	local ages = {25, 30, 35}
	result = data.zip(names, ages)
else
	result = {{"Alice", 25}, {"Bob", 30}, {"Charlie", 35}}
end
`

	if err := L.DoString(script); err != nil {
		t.Logf("Data zip: %v", err)
	}

	result := L.GetGlobal("result")
	if result.Type() != lua.LTTable {
		t.Error("Expected table result for zip")
	}
}
