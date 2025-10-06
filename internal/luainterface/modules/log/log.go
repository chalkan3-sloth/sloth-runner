package log

import (
	"fmt"
	"log/slog"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// Info logs an info message
func Info(L *lua.LState) int {
	message := L.CheckString(1)
	slog.Info(message, "source", "lua")
	return 0
}

// Warn logs a warning message
func Warn(L *lua.LState) int {
	message := L.CheckString(1)
	slog.Warn(message, "source", "lua")
	return 0
}

// Error logs an error message
func Error(L *lua.LState) int {
	message := L.CheckString(1)
	slog.Error(message, "source", "lua")
	return 0
}

// Debug logs a debug message
func Debug(L *lua.LState) int {
	message := L.CheckString(1)
	slog.Debug(message, "source", "lua")
	return 0
}

// Print prints directly to stdout without formatting
func Print(L *lua.LState) int {
	message := L.CheckString(1)
	fmt.Println(message)
	return 0
}

// Table logs complex table/object data in a structured format
func Table(L *lua.LState) int {
	value := L.CheckAny(1)
	title := L.OptString(2, "Table Contents")

	// Convert value to a formatted string representation
	formatted := formatLuaValueForLog(value, 0)

	// Log with structured format
	slog.Info(fmt.Sprintf("📊 %s:\n%s", title, formatted), "source", "lua")
	return 0
}

// formatLuaValueForLog formats a Lua value for pretty logging
func formatLuaValueForLog(value lua.LValue, indent int) string {
	indentStr := strings.Repeat("  ", indent)

	switch v := value.(type) {
	case lua.LString:
		return fmt.Sprintf("\"%s\"", v.String())
	case lua.LNumber:
		return v.String()
	case lua.LBool:
		if bool(v) {
			return "true"
		}
		return "false"
	case *lua.LTable:
		var parts []string
		parts = append(parts, "{")

		// Check if it's an array-like table
		isArray := true
		maxIndex := 0
		v.ForEach(func(key, _ lua.LValue) {
			if key.Type() == lua.LTNumber {
				if idx := int(lua.LVAsNumber(key)); idx > maxIndex {
					maxIndex = idx
				}
			} else {
				isArray = false
			}
		})

		if isArray && maxIndex > 0 {
			// Format as array
			for i := 1; i <= maxIndex; i++ {
				val := v.RawGetInt(i)
				if val != lua.LNil {
					formatted := formatLuaValueForLog(val, indent+1)
					parts = append(parts, fmt.Sprintf("%s  [%d] = %s", indentStr, i, formatted))
				}
			}
		} else {
			// Format as object
			v.ForEach(func(key, val lua.LValue) {
				keyStr := key.String()
				formatted := formatLuaValueForLog(val, indent+1)
				parts = append(parts, fmt.Sprintf("%s  %s = %s", indentStr, keyStr, formatted))
			})
		}

		parts = append(parts, indentStr+"}")
		return strings.Join(parts, "\n")

	default:
		if value == lua.LNil {
			return "nil"
		}
		return fmt.Sprintf("<%s>", value.Type().String())
	}
}

// Loader returns the log module loader
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"info":  Info,
		"warn":  Warn,
		"error": Error,
		"debug": Debug,
		"print": Print,
		"table": Table,
	})
	L.Push(mod)
	return 1
}

// Open registers the log module and loads it globally
func Open(L *lua.LState) {
	L.PreloadModule("log", Loader)
	if err := L.DoString(`log = require("log")`); err != nil {
		panic(err)
	}
}
