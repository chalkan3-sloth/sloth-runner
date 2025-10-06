package core

import (
	"fmt"
	"time"

	"github.com/yuin/gopher-lua"
)

// EventDispatcher is a function type for dispatching events
// This avoids import cycles by using dependency injection
type EventDispatcher func(eventType string, data map[string]interface{}) error

var globalEventDispatcher EventDispatcher

// SetGlobalEventDispatcher sets the global event dispatcher
// This should be called at application startup
func SetGlobalEventDispatcher(dispatcher EventDispatcher) {
	globalEventDispatcher = dispatcher
}

// EventModule provides event dispatching functionality for workflows
type EventModule struct {
	info CoreModuleInfo
}

// NewEventModule creates a new event module
func NewEventModule() *EventModule {
	info := CoreModuleInfo{
		Name:        "event",
		Version:     "1.0.0",
		Description: "Event dispatching system for triggering hooks from workflows",
		Author:      "Sloth Runner Team",
		Category:    "core",
		Dependencies: []string{},
	}

	return &EventModule{
		info: info,
	}
}

// Info returns module information
func (e *EventModule) Info() CoreModuleInfo {
	return e.info
}

// Loader returns the Lua loader function
func (e *EventModule) Loader(L *lua.LState) int {
	eventTable := L.NewTable()

	// Event functions
	L.SetFuncs(eventTable, map[string]lua.LGFunction{
		"dispatch":         e.luaDispatch,
		"dispatch_custom":  e.luaDispatchCustom,
		"dispatch_file":    e.luaDispatchFile,
	})

	L.Push(eventTable)
	return 1
}

// luaDispatch dispatches a generic event to hooks
// Usage: event.dispatch(event_type, data_table)
func (e *EventModule) luaDispatch(L *lua.LState) int {
	if globalEventDispatcher == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("event dispatcher not initialized"))
		return 2
	}

	eventType := L.CheckString(1)
	dataTable := L.CheckTable(2)

	// Convert Lua table to Go map
	data := make(map[string]interface{})
	dataTable.ForEach(func(k, v lua.LValue) {
		key := k.String()
		data[key] = e.luaValueToGo(v)
	})

	// Add timestamp to data
	data["timestamp"] = time.Now().Unix()

	// Dispatch event using global dispatcher
	if err := globalEventDispatcher(eventType, data); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to dispatch event: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// luaDispatchCustom dispatches a custom event with a simple message
// Usage: event.dispatch_custom(event_name, message)
func (e *EventModule) luaDispatchCustom(L *lua.LState) int {
	if globalEventDispatcher == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("event dispatcher not initialized"))
		return 2
	}

	eventName := L.CheckString(1)
	message := L.CheckString(2)

	// Create custom event data
	data := map[string]interface{}{
		"name":      eventName,
		"message":   message,
		"timestamp": time.Now().Unix(),
	}

	// Dispatch using "custom" event type
	if err := globalEventDispatcher("custom", data); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to dispatch custom event: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// luaDispatchFile dispatches a file event
// Usage: event.dispatch_file(event_type, file_path, [watcher_name])
func (e *EventModule) luaDispatchFile(L *lua.LState) int {
	if globalEventDispatcher == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("event dispatcher not initialized"))
		return 2
	}

	eventType := L.CheckString(1) // created, modified, deleted
	filePath := L.CheckString(2)
	watcherName := ""
	if L.GetTop() >= 3 {
		watcherName = L.CheckString(3)
	}

	// Map event type to file event type
	var fileEventType string
	switch eventType {
	case "created":
		fileEventType = "file.created"
	case "modified":
		fileEventType = "file.modified"
	case "deleted":
		fileEventType = "file.deleted"
	default:
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("invalid file event type: %s (must be created, modified, or deleted)", eventType)))
		return 2
	}

	// Create file event data
	data := map[string]interface{}{
		"path":      filePath,
		"watcher":   watcherName,
		"timestamp": time.Now().Unix(),
	}

	// Dispatch file event
	if err := globalEventDispatcher(fileEventType, data); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to dispatch file event: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// luaValueToGo converts a Lua value to a Go value
func (e *EventModule) luaValueToGo(val lua.LValue) interface{} {
	switch v := val.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		// Check if it's an array or a map
		isArray := true
		length := 0
		v.ForEach(func(k, _ lua.LValue) {
			if k.Type() != lua.LTNumber {
				isArray = false
			}
			length++
		})

		if isArray && length > 0 {
			// Convert to array
			arr := make([]interface{}, 0, length)
			v.ForEach(func(_, val lua.LValue) {
				arr = append(arr, e.luaValueToGo(val))
			})
			return arr
		} else {
			// Convert to map
			m := make(map[string]interface{})
			v.ForEach(func(k, val lua.LValue) {
				m[k.String()] = e.luaValueToGo(val)
			})
			return m
		}
	default:
		return v.String()
	}
}
