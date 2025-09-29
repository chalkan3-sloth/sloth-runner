package luainterface

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"
)

// DataModule provides data processing functionality for Lua scripts
type DataModule struct{}

// NewDataModule creates a new data processing module
func NewDataModule() *DataModule {
	return &DataModule{}
}

// RegisterDataModule registers the data module with the Lua state
func RegisterDataModule(L *lua.LState) {
	module := NewDataModule()
	
	// Create the data table
	dataTable := L.NewTable()
	
	// JSON functions
	L.SetField(dataTable, "json_encode", L.NewFunction(module.luaJSONEncode))
	L.SetField(dataTable, "json_decode", L.NewFunction(module.luaJSONDecode))
	L.SetField(dataTable, "json_pretty", L.NewFunction(module.luaJSONPretty))
	L.SetField(dataTable, "json_validate", L.NewFunction(module.luaJSONValidate))
	
	// XML functions
	L.SetField(dataTable, "xml_encode", L.NewFunction(module.luaXMLEncode))
	L.SetField(dataTable, "xml_decode", L.NewFunction(module.luaXMLDecode))
	L.SetField(dataTable, "xml_to_json", L.NewFunction(module.luaXMLToJSON))
	
	// YAML functions
	L.SetField(dataTable, "yaml_encode", L.NewFunction(module.luaYAMLEncode))
	L.SetField(dataTable, "yaml_decode", L.NewFunction(module.luaYAMLDecode))
	L.SetField(dataTable, "yaml_to_json", L.NewFunction(module.luaYAMLToJSON))
	L.SetField(dataTable, "json_to_yaml", L.NewFunction(module.luaJSONToYAML))
	
	// CSV functions
	L.SetField(dataTable, "csv_parse", L.NewFunction(module.luaCSVParse))
	L.SetField(dataTable, "csv_generate", L.NewFunction(module.luaCSVGenerate))
	
	// Data transformation functions
	L.SetField(dataTable, "deep_merge", L.NewFunction(module.luaDeepMerge))
	L.SetField(dataTable, "flatten", L.NewFunction(module.luaFlatten))
	L.SetField(dataTable, "unflatten", L.NewFunction(module.luaUnflatten))
	L.SetField(dataTable, "get_path", L.NewFunction(module.luaGetPath))
	L.SetField(dataTable, "set_path", L.NewFunction(module.luaSetPath))
	
	// Register the data table globally
	L.SetGlobal("data", dataTable)
}

// JSON functions
func (d *DataModule) luaJSONEncode(L *lua.LState) int {
	lv := L.CheckAny(1)
	
	data, err := dataLuaValueToGo(lv)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(jsonBytes)))
	return 1
}

func (d *DataModule) luaJSONDecode(L *lua.LState) int {
	jsonStr := L.CheckString(1)
	
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	lv := dataGoValueToLua(L, data)
	L.Push(lv)
	return 1
}

func (d *DataModule) luaJSONPretty(L *lua.LState) int {
	lv := L.CheckAny(1)
	
	data, err := dataLuaValueToGo(lv)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(jsonBytes)))
	return 1
}

func (d *DataModule) luaJSONValidate(L *lua.LState) int {
	jsonStr := L.CheckString(1)
	
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// XML functions
func (d *DataModule) luaXMLEncode(L *lua.LState) int {
	lv := L.CheckAny(1)
	rootName := L.OptString(2, "root")
	
	data, err := dataLuaValueToGo(lv)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Simple XML generation for basic structures
	xmlStr := fmt.Sprintf("<%s>%s</%s>", rootName, interfaceToXML(data), rootName)
	
	L.Push(lua.LString(xmlStr))
	return 1
}

func (d *DataModule) luaXMLDecode(L *lua.LState) int {
	xmlStr := L.CheckString(1)
	
	// Simple XML parsing - for complex XML, users should use dedicated XML libraries
	var data interface{}
	err := xml.Unmarshal([]byte(xmlStr), &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	lv := dataGoValueToLua(L, data)
	L.Push(lv)
	return 1
}

func (d *DataModule) luaXMLToJSON(L *lua.LState) int {
	xmlStr := L.CheckString(1)
	
	var data interface{}
	err := xml.Unmarshal([]byte(xmlStr), &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(jsonBytes)))
	return 1
}

// YAML functions
func (d *DataModule) luaYAMLEncode(L *lua.LState) int {
	lv := L.CheckAny(1)
	
	data, err := dataLuaValueToGo(lv)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(yamlBytes)))
	return 1
}

func (d *DataModule) luaYAMLDecode(L *lua.LState) int {
	yamlStr := L.CheckString(1)
	
	var data interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	lv := dataGoValueToLua(L, data)
	L.Push(lv)
	return 1
}

func (d *DataModule) luaYAMLToJSON(L *lua.LState) int {
	yamlStr := L.CheckString(1)
	
	var data interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(jsonBytes)))
	return 1
}

func (d *DataModule) luaJSONToYAML(L *lua.LState) int {
	jsonStr := L.CheckString(1)
	
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(yamlBytes)))
	return 1
}

// CSV functions
func (d *DataModule) luaCSVParse(L *lua.LState) int {
	csvStr := L.CheckString(1)
	delimiter := L.OptString(2, ",")
	
	lines := strings.Split(csvStr, "\n")
	result := L.NewTable()
	
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		fields := strings.Split(line, delimiter)
		row := L.NewTable()
		
		for j, field := range fields {
			row.RawSetInt(j+1, lua.LString(strings.TrimSpace(field)))
		}
		
		result.RawSetInt(i+1, row)
	}
	
	L.Push(result)
	return 1
}

func (d *DataModule) luaCSVGenerate(L *lua.LState) int {
	table := L.CheckTable(1)
	delimiter := L.OptString(2, ",")
	
	var lines []string
	
	table.ForEach(func(key, value lua.LValue) {
		if row, ok := value.(*lua.LTable); ok {
			var fields []string
			row.ForEach(func(k, v lua.LValue) {
				fields = append(fields, lua.LVAsString(v))
			})
			lines = append(lines, strings.Join(fields, delimiter))
		}
	})
	
	L.Push(lua.LString(strings.Join(lines, "\n")))
	return 1
}

// Data transformation functions
func (d *DataModule) luaDeepMerge(L *lua.LState) int {
	// Simple implementation - merge two tables
	table1 := L.CheckTable(1)
	table2 := L.CheckTable(2)
	
	result := L.NewTable()
	
	// Copy table1
	table1.ForEach(func(key, value lua.LValue) {
		result.RawSet(key, value)
	})
	
	// Merge table2
	table2.ForEach(func(key, value lua.LValue) {
		result.RawSet(key, value)
	})
	
	L.Push(result)
	return 1
}

func (d *DataModule) luaFlatten(L *lua.LState) int {
	table := L.CheckTable(1)
	separator := L.OptString(2, ".")
	
	result := L.NewTable()
	flattenTable(table, result, "", separator)
	
	L.Push(result)
	return 1
}

func (d *DataModule) luaUnflatten(L *lua.LState) int {
	table := L.CheckTable(1)
	separator := L.OptString(2, ".")
	
	result := L.NewTable()
	
	table.ForEach(func(key, value lua.LValue) {
		keyStr := lua.LVAsString(key)
		parts := strings.Split(keyStr, separator)
		
		current := result
		for i, part := range parts {
			if i == len(parts)-1 {
				current.RawSetString(part, value)
			} else {
				next := current.RawGetString(part)
				if next == lua.LNil {
					next = L.NewTable()
					current.RawSetString(part, next)
				}
				current = next.(*lua.LTable)
			}
		}
	})
	
	L.Push(result)
	return 1
}

func (d *DataModule) luaGetPath(L *lua.LState) int {
	table := L.CheckTable(1)
	path := L.CheckString(2)
	separator := L.OptString(3, ".")
	
	parts := strings.Split(path, separator)
	current := lua.LValue(table)
	
	for _, part := range parts {
		if tbl, ok := current.(*lua.LTable); ok {
			current = tbl.RawGetString(part)
			if current == lua.LNil {
				L.Push(lua.LNil)
				return 1
			}
		} else {
			L.Push(lua.LNil)
			return 1
		}
	}
	
	L.Push(current)
	return 1
}

func (d *DataModule) luaSetPath(L *lua.LState) int {
	table := L.CheckTable(1)
	path := L.CheckString(2)
	value := L.CheckAny(3)
	separator := L.OptString(4, ".")
	
	parts := strings.Split(path, separator)
	current := table
	
	for i, part := range parts {
		if i == len(parts)-1 {
			current.RawSetString(part, value)
		} else {
			next := current.RawGetString(part)
			if next == lua.LNil {
				next = L.NewTable()
				current.RawSetString(part, next)
			}
			current = next.(*lua.LTable)
		}
	}
	
	L.Push(table)
	return 1
}

// Helper functions for data module
func dataLuaValueToGo(lv lua.LValue) (interface{}, error) {
	switch v := lv.(type) {
	case lua.LBool:
		return bool(v), nil
	case lua.LString:
		return string(v), nil
	case lua.LNumber:
		return float64(v), nil
	case *lua.LTable:
		result := make(map[string]interface{})
		var err error
		v.ForEach(func(key, value lua.LValue) {
			keyStr := lua.LVAsString(key)
			result[keyStr], err = dataLuaValueToGo(value)
		})
		return result, err
	case *lua.LNilType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported Lua type: %T", v)
	}
}

func dataGoValueToLua(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}
	
	switch v := value.(type) {
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case map[string]interface{}:
		table := L.NewTable()
		for key, val := range v {
			table.RawSetString(key, dataGoValueToLua(L, val))
		}
		return table
	case []interface{}:
		table := L.NewTable()
		for i, val := range v {
			table.RawSetInt(i+1, dataGoValueToLua(L, val))
		}
		return table
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

func interfaceToXML(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return strconv.FormatBool(v)
	case map[string]interface{}:
		var result strings.Builder
		for key, value := range v {
			result.WriteString(fmt.Sprintf("<%s>%s</%s>", key, interfaceToXML(value), key))
		}
		return result.String()
	case []interface{}:
		var result strings.Builder
		for _, value := range v {
			result.WriteString(fmt.Sprintf("<item>%s</item>", interfaceToXML(value)))
		}
		return result.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func flattenTable(source *lua.LTable, target *lua.LTable, prefix, separator string) {
	source.ForEach(func(key, value lua.LValue) {
		keyStr := lua.LVAsString(key)
		newKey := keyStr
		if prefix != "" {
			newKey = prefix + separator + keyStr
		}
		
		if subtable, ok := value.(*lua.LTable); ok {
			flattenTable(subtable, target, newKey, separator)
		} else {
			target.RawSetString(newKey, value)
		}
	})
}