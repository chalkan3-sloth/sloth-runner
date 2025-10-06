package data

import (
	"encoding/json"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

// ParseJSON parses JSON string into Lua table
func ParseJSON(L *lua.LState) int {
	jsonStr := L.CheckString(1)

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(goValueToLua(L, data))
	return 1
}

// ToJSON converts Lua table to JSON string
func ToJSON(L *lua.LState) int {
	value := L.CheckAny(1)
	goValue := luaToGoValue(L, value)

	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(string(jsonBytes)))
	return 1
}

// ParseYAML parses YAML string into Lua table
func ParseYAML(L *lua.LState) int {
	yamlStr := L.CheckString(1)

	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &data); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Convert map[interface{}]interface{} to map[string]interface{}
	data = normalizeYAMLValue(data)

	L.Push(goValueToLua(L, data))
	return 1
}

// ToYAML converts Lua table to YAML string
func ToYAML(L *lua.LState) int {
	value := L.CheckAny(1)
	goValue := luaToGoValue(L, value)

	yamlBytes, err := yaml.Marshal(goValue)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(string(yamlBytes)))
	return 1
}

// Loader returns the data module loader
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"parse_json": ParseJSON,
		"to_json":    ToJSON,
		"parse_yaml": ParseYAML,
		"to_yaml":    ToYAML,
	})
	L.Push(mod)
	return 1
}

// Open registers the data module and loads it globally
func Open(L *lua.LState) {
	L.PreloadModule("data", Loader)
	if err := L.DoString(`data = require("data")`); err != nil {
		panic(err)
	}
}

// Helper functions

func goValueToLua(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		tbl := L.NewTable()
		for i, item := range v {
			tbl.RawSetInt(i+1, goValueToLua(L, item))
		}
		return tbl
	case map[string]interface{}:
		tbl := L.NewTable()
		for k, val := range v {
			tbl.RawSetString(k, goValueToLua(L, val))
		}
		return tbl
	default:
		return lua.LString(value.(string))
	}
}

func luaToGoValue(L *lua.LState, value lua.LValue) interface{} {
	switch v := value.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		maxN := v.MaxN()
		if maxN > 0 {
			// Array
			arr := make([]interface{}, 0, maxN)
			for i := 1; i <= maxN; i++ {
				arr = append(arr, luaToGoValue(L, v.RawGetInt(i)))
			}
			return arr
		}
		// Map
		m := make(map[string]interface{})
		v.ForEach(func(key, val lua.LValue) {
			m[key.String()] = luaToGoValue(L, val)
		})
		return m
	default:
		return v.String()
	}
}

func normalizeYAMLValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, val := range v {
			m[k.(string)] = normalizeYAMLValue(val)
		}
		return m
	case []interface{}:
		arr := make([]interface{}, len(v))
		for i, val := range v {
			arr[i] = normalizeYAMLValue(val)
		}
		return arr
	default:
		return v
	}
}
