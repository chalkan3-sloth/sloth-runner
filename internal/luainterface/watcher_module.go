package luainterface

import (
	"encoding/json"
	"log/slog"

	"github.com/chalkan3-sloth/sloth-runner/internal/luamodules"
	lua "github.com/yuin/gopher-lua"
)

// RegisterWatcherModule registers the watcher registration module
func RegisterWatcherModule(L *lua.LState) {
	// Create watcher table
	watcherTable := L.NewTable()

	// Create register sub-table
	registerModule := luamodules.NewEventRegisterModule()
	registerModule.Load(L)
	registerTable := L.Get(-1)
	L.Pop(1)

	// Set register as sub-table of watcher
	L.SetField(watcherTable, "register", registerTable)

	// Set watcher as global
	L.SetGlobal("watcher", watcherTable)

	slog.Info("Watcher module registered successfully")
}

// GetRegisteredWatchers retrieves all registered watchers from Lua state
func GetRegisteredWatchers(L *lua.LState) ([]map[string]interface{}, error) {
	watchersTable := L.GetGlobal("_WATCHERS")
	if watchersTable == lua.LNil {
		return nil, nil
	}

	var watchers []map[string]interface{}

	if tbl, ok := watchersTable.(*lua.LTable); ok {
		tbl.ForEach(func(key, value lua.LValue) {
			if value.Type() == lua.LTString {
				// Parse JSON config
				var config map[string]interface{}
				if err := json.Unmarshal([]byte(value.String()), &config); err == nil {
					watchers = append(watchers, config)
				}
			}
		})
	}

	return watchers, nil
}
