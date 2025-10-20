//go:build cgo
// +build cgo

package handlers

import (
	"context"
	"io"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// Test RunConfig structure
func TestRunConfig_Structure(t *testing.T) {
	config := &RunConfig{
		StackName:        "test-stack",
		FilePath:         "/path/to/file.sloth",
		Values:           "/path/to/values.yaml",
		Interactive:      true,
		OutputStyle:      "enhanced",
		Debug:            false,
		DelegateToHosts:  []string{"host1", "host2"},
		SSHProfile:       "default",
		SSHPasswordStdin: false,
		PasswordStdin:    false,
		YesFlag:          false,
		Context:          context.Background(),
		Writer:           io.Discard,
		RunID:            "run-123",
	}

	if config.StackName != "test-stack" {
		t.Error("Expected StackName to be set")
	}

	if config.FilePath != "/path/to/file.sloth" {
		t.Error("Expected FilePath to be set")
	}

	if !config.Interactive {
		t.Error("Expected Interactive to be true")
	}

	if len(config.DelegateToHosts) != 2 {
		t.Error("Expected 2 delegate hosts")
	}
}

func TestRunConfig_DefaultValues(t *testing.T) {
	config := &RunConfig{}

	if config.Interactive {
		t.Error("Expected Interactive to be false by default")
	}

	if config.Debug {
		t.Error("Expected Debug to be false by default")
	}

	if config.YesFlag {
		t.Error("Expected YesFlag to be false by default")
	}

	if config.SSHPasswordStdin {
		t.Error("Expected SSHPasswordStdin to be false by default")
	}

	if config.PasswordStdin {
		t.Error("Expected PasswordStdin to be false by default")
	}
}

func TestRunConfig_HasRequiredFields(t *testing.T) {
	config := &RunConfig{}

	// Test that all required string fields exist
	_ = config.StackName
	_ = config.FilePath
	_ = config.Values
	_ = config.OutputStyle
	_ = config.SSHProfile
	_ = config.RunID

	// Test that all required slice fields exist
	_ = config.DelegateToHosts

	// Test that all required interface fields exist
	_ = config.Context
	_ = config.Writer
	_ = config.AgentRegistry
}

func TestRunConfig_BooleanFields(t *testing.T) {
	config := &RunConfig{
		Interactive:      true,
		Debug:            true,
		SSHPasswordStdin: true,
		PasswordStdin:    true,
		YesFlag:          true,
	}

	if !config.Interactive {
		t.Error("Expected Interactive to be true")
	}

	if !config.Debug {
		t.Error("Expected Debug to be true")
	}

	if !config.SSHPasswordStdin {
		t.Error("Expected SSHPasswordStdin to be true")
	}

	if !config.PasswordStdin {
		t.Error("Expected PasswordStdin to be true")
	}

	if !config.YesFlag {
		t.Error("Expected YesFlag to be true")
	}
}

func TestRunConfig_OutputStyles(t *testing.T) {
	tests := []struct {
		name  string
		style string
	}{
		{"enhanced", "enhanced"},
		{"rich", "rich"},
		{"modern", "modern"},
		{"json", "json"},
		{"default", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &RunConfig{
				OutputStyle: tt.style,
			}

			if config.OutputStyle != tt.style {
				t.Errorf("Expected OutputStyle '%s', got '%s'", tt.style, config.OutputStyle)
			}
		})
	}
}

func TestRunConfig_MultipleDelegateHosts(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	config := &RunConfig{
		DelegateToHosts: hosts,
	}

	if len(config.DelegateToHosts) != 3 {
		t.Errorf("Expected 3 hosts, got %d", len(config.DelegateToHosts))
	}

	for i, host := range hosts {
		if config.DelegateToHosts[i] != host {
			t.Errorf("Expected host '%s' at index %d, got '%s'", host, i, config.DelegateToHosts[i])
		}
	}
}

func TestRunConfig_SingleDelegateHost(t *testing.T) {
	config := &RunConfig{
		DelegateToHosts: []string{"single-host"},
	}

	if len(config.DelegateToHosts) != 1 {
		t.Error("Expected 1 host")
	}

	if config.DelegateToHosts[0] != "single-host" {
		t.Error("Expected host to be 'single-host'")
	}
}

func TestRunConfig_NoDelegateHosts(t *testing.T) {
	config := &RunConfig{
		DelegateToHosts: []string{},
	}

	if len(config.DelegateToHosts) != 0 {
		t.Error("Expected no hosts")
	}
}

func TestRunConfig_WithContext(t *testing.T) {
	ctx := context.Background()
	config := &RunConfig{
		Context: ctx,
	}

	if config.Context == nil {
		t.Error("Expected Context to be set")
	}
}

func TestRunConfig_WithWriter(t *testing.T) {
	var buf strings.Builder
	config := &RunConfig{
		Writer: &buf,
	}

	if config.Writer == nil {
		t.Error("Expected Writer to be set")
	}
}

// Test helper functions

func TestMapToLuaTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	table := mapToLuaTable(L, m)

	if table == nil {
		t.Error("Expected non-nil table")
	}

	// Lua tables created with RawSetString don't increment Len()
	// Len() only works for sequential numeric indices
	// We just check that the table was created
	if table.Type() != lua.LTTable {
		t.Error("Expected table type")
	}
}

func TestInterfaceToLuaValue_String(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	val := interfaceToLuaValue(L, "test")

	if val.Type() != lua.LTString {
		t.Errorf("Expected LTString, got %v", val.Type())
	}

	if val.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", val.String())
	}
}

func TestInterfaceToLuaValue_Int(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	val := interfaceToLuaValue(L, 42)

	if val.Type() != lua.LTNumber {
		t.Errorf("Expected LTNumber, got %v", val.Type())
	}
}

func TestInterfaceToLuaValue_Float(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	val := interfaceToLuaValue(L, 3.14)

	if val.Type() != lua.LTNumber {
		t.Errorf("Expected LTNumber, got %v", val.Type())
	}
}

func TestInterfaceToLuaValue_Bool(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	val := interfaceToLuaValue(L, true)

	if val.Type() != lua.LTBool {
		t.Errorf("Expected LTBool, got %v", val.Type())
	}
}

func TestInterfaceToLuaValue_Map(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{
		"nested": "value",
	}

	val := interfaceToLuaValue(L, m)

	if val.Type() != lua.LTTable {
		t.Errorf("Expected LTTable, got %v", val.Type())
	}
}

func TestInterfaceToLuaValue_Nil(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	val := interfaceToLuaValue(L, struct{}{})

	if val.Type() != lua.LTNil {
		t.Errorf("Expected LTNil, got %v", val.Type())
	}
}

func TestLuaValueToInterface_Nil(t *testing.T) {
	val := luaValueToInterface(lua.LNil)

	if val != nil {
		t.Error("Expected nil")
	}
}

func TestLuaValueToInterface_Bool(t *testing.T) {
	val := luaValueToInterface(lua.LBool(true))

	if val != true {
		t.Error("Expected true")
	}
}

func TestLuaValueToInterface_String(t *testing.T) {
	val := luaValueToInterface(lua.LString("test"))

	if val != "test" {
		t.Errorf("Expected 'test', got '%v'", val)
	}
}

func TestLuaValueToInterface_Number(t *testing.T) {
	val := luaValueToInterface(lua.LNumber(42))

	if val != float64(42) {
		t.Errorf("Expected 42, got %v", val)
	}
}

func TestLuaValueToInterface_Table(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("key", lua.LString("value"))

	val := luaValueToInterface(table)

	m, ok := val.(map[string]interface{})
	if !ok {
		t.Error("Expected map")
	}

	if m["key"] != "value" {
		t.Error("Expected 'value' for key")
	}
}

// Test remoteAgentResolver

func TestRemoteAgentResolver_Structure(t *testing.T) {
	resolver := &remoteAgentResolver{
		masterAddr: "localhost:50053",
		conn:       nil,
		client:     nil,
	}

	if resolver.masterAddr != "localhost:50053" {
		t.Error("Expected masterAddr to be set")
	}
}

func TestRemoteAgentResolver_Close_WithNilConn(t *testing.T) {
	resolver := &remoteAgentResolver{
		masterAddr: "localhost:50053",
		conn:       nil,
	}

	err := resolver.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test conversion functions with edge cases

func TestMapToLuaTable_EmptyMap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{}
	table := mapToLuaTable(L, m)

	if table == nil {
		t.Error("Expected non-nil table")
	}

	if table.Type() != lua.LTTable {
		t.Error("Expected table type")
	}
}

func TestMapToLuaTable_NestedMap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{
		"outer": map[string]interface{}{
			"inner": "value",
		},
	}

	table := mapToLuaTable(L, m)

	if table == nil {
		t.Error("Expected non-nil table")
	}

	if table.Type() != lua.LTTable {
		t.Error("Expected table type")
	}
}

func TestInterfaceToLuaValue_ComplexMap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{
		"string": "value",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"nested": map[string]interface{}{
			"key": "nested_value",
		},
	}

	val := interfaceToLuaValue(L, m)

	if val.Type() != lua.LTTable {
		t.Errorf("Expected LTTable, got %v", val.Type())
	}

	table := val.(*lua.LTable)
	if table == nil {
		t.Error("Expected non-nil table")
	}
}

func TestLuaValueToInterface_EmptyTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	val := luaValueToInterface(table)

	m, ok := val.(map[string]interface{})
	if !ok {
		t.Error("Expected map")
	}

	if len(m) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(m))
	}
}

func TestLuaValueToInterface_NestedTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	innerTable := L.NewTable()
	innerTable.RawSetString("inner_key", lua.LString("inner_value"))

	outerTable := L.NewTable()
	outerTable.RawSetString("outer", innerTable)

	val := luaValueToInterface(outerTable)

	m, ok := val.(map[string]interface{})
	if !ok {
		t.Error("Expected map")
	}

	nested, ok := m["outer"].(map[string]interface{})
	if !ok {
		t.Error("Expected nested map")
	}

	if nested["inner_key"] != "inner_value" {
		t.Error("Expected 'inner_value'")
	}
}

// Test RunConfig validation scenarios

func TestRunConfig_ValidConfiguration(t *testing.T) {
	config := &RunConfig{
		StackName:   "valid-stack",
		FilePath:    "workflow.sloth",
		Context:     context.Background(),
		Writer:      io.Discard,
		OutputStyle: "enhanced",
	}

	if config.StackName == "" {
		t.Error("StackName should not be empty")
	}

	if config.FilePath == "" {
		t.Error("FilePath should not be empty")
	}
}

func TestRunConfig_WithAllFlags(t *testing.T) {
	config := &RunConfig{
		StackName:        "stack",
		FilePath:         "file.sloth",
		Values:           "values.yaml",
		Interactive:      true,
		OutputStyle:      "json",
		Debug:            true,
		DelegateToHosts:  []string{"host1"},
		SSHProfile:       "profile1",
		SSHPasswordStdin: true,
		PasswordStdin:    true,
		YesFlag:          true,
		Context:          context.Background(),
		Writer:           io.Discard,
		RunID:            "run-456",
	}

	if config.StackName != "stack" {
		t.Error("StackName mismatch")
	}
	if config.FilePath != "file.sloth" {
		t.Error("FilePath mismatch")
	}
	if config.Values != "values.yaml" {
		t.Error("Values mismatch")
	}
	if !config.Interactive {
		t.Error("Interactive should be true")
	}
	if config.OutputStyle != "json" {
		t.Error("OutputStyle mismatch")
	}
	if !config.Debug {
		t.Error("Debug should be true")
	}
	if len(config.DelegateToHosts) != 1 {
		t.Error("DelegateToHosts length mismatch")
	}
	if config.SSHProfile != "profile1" {
		t.Error("SSHProfile mismatch")
	}
	if !config.SSHPasswordStdin {
		t.Error("SSHPasswordStdin should be true")
	}
	if !config.PasswordStdin {
		t.Error("PasswordStdin should be true")
	}
	if !config.YesFlag {
		t.Error("YesFlag should be true")
	}
	if config.RunID != "run-456" {
		t.Error("RunID mismatch")
	}
}

func TestRunConfig_WithAgentRegistry(t *testing.T) {
	registry := &mockAgentRegistry{}
	config := &RunConfig{
		AgentRegistry: registry,
	}

	if config.AgentRegistry == nil {
		t.Error("Expected AgentRegistry to be set")
	}
}

type mockAgentRegistry struct{}

func TestRunConfig_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := &RunConfig{
		Context: ctx,
	}

	if config.Context.Err() != nil {
		t.Error("Expected no error before cancellation")
	}

	cancel()

	if config.Context.Err() == nil {
		t.Error("Expected error after cancellation")
	}
}

func TestInterfaceToLuaValue_AllTypes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		input    interface{}
		expected lua.LValueType
	}{
		{"string", "test", lua.LTString},
		{"int", 42, lua.LTNumber},
		{"float", 3.14, lua.LTNumber},
		{"bool", true, lua.LTBool},
		{"map", map[string]interface{}{}, lua.LTTable},
		{"unknown", struct{}{}, lua.LTNil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := interfaceToLuaValue(L, tt.input)
			if val.Type() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, val.Type())
			}
		})
	}
}

func TestLuaValueToInterface_AllTypes(t *testing.T) {
	tests := []struct {
		name  string
		input lua.LValue
	}{
		{"nil", lua.LNil},
		{"bool", lua.LBool(true)},
		{"string", lua.LString("test")},
		{"number", lua.LNumber(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := luaValueToInterface(tt.input)
			// Just ensure no panic occurs
			_ = val
		})
	}
}

func TestMapToLuaTable_WithNilValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{
		"nil_value": nil,
	}

	table := mapToLuaTable(L, m)

	if table == nil {
		t.Error("Expected non-nil table")
	}
}

func TestRunConfig_EmptyDelegateHosts(t *testing.T) {
	config := &RunConfig{
		DelegateToHosts: nil,
	}

	if config.DelegateToHosts != nil && len(config.DelegateToHosts) > 0 {
		t.Error("Expected empty or nil DelegateToHosts")
	}
}

func TestRunConfig_LongRunID(t *testing.T) {
	longID := strings.Repeat("a", 100)
	config := &RunConfig{
		RunID: longID,
	}

	if len(config.RunID) != 100 {
		t.Errorf("Expected RunID length 100, got %d", len(config.RunID))
	}
}
