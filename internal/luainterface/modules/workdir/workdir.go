package workdir

import (
	"log/slog"
	"os"

	lua "github.com/yuin/gopher-lua"
)

// Get returns the current workdir from task context or working directory
func Get(L *lua.LState) int {
	// Get current workdir from task context
	taskContext := L.GetGlobal("__task_context")
	if taskContext.Type() == lua.LTTable {
		workdir := taskContext.(*lua.LTable).RawGetString("workdir")
		if workdir.Type() == lua.LTString {
			L.Push(workdir)
			return 1
		}
	}

	// Fallback to current working directory
	if cwd, err := os.Getwd(); err == nil {
		L.Push(lua.LString(cwd))
	} else {
		L.Push(lua.LString("/tmp"))
	}
	return 1
}

// Cleanup cleans up the workdir (currently disabled)
func Cleanup(L *lua.LState) int {
	// Get workdir path (optional argument)
	var workdirPath string
	if L.GetTop() >= 1 {
		workdirPath = L.CheckString(1)
	} else {
		// Get from context
		taskContext := L.GetGlobal("__task_context")
		if taskContext.Type() == lua.LTTable {
			workdir := taskContext.(*lua.LTable).RawGetString("workdir")
			if workdir.Type() == lua.LTString {
				workdirPath = workdir.String()
			}
		}
	}

	if workdirPath == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("no workdir specified"))
		return 2
	}

	slog.Warn("Manual workdir cleanup is disabled. Workdir preserved.", "workdir", workdirPath)
	// if err := os.RemoveAll(workdirPath); err != nil {
	// 	L.Push(lua.LBool(false))
	// 	L.Push(lua.LString(err.Error()))
	// 	return 2
	// }

	L.Push(lua.LBool(true))
	L.Push(lua.LString("workdir cleanup is disabled"))
	return 2
}

// Exists checks if workdir exists
func Exists(L *lua.LState) int {
	// Get workdir path (optional argument)
	var workdirPath string
	if L.GetTop() >= 1 {
		workdirPath = L.CheckString(1)
	} else {
		// Get from context
		taskContext := L.GetGlobal("__task_context")
		if taskContext.Type() == lua.LTTable {
			workdir := taskContext.(*lua.LTable).RawGetString("workdir")
			if workdir.Type() == lua.LTString {
				workdirPath = workdir.String()
			}
		}
	}

	if workdirPath == "" {
		L.Push(lua.LBool(false))
		return 1
	}

	// Check if directory exists
	if _, err := os.Stat(workdirPath); err == nil {
		L.Push(lua.LBool(true))
	} else {
		L.Push(lua.LBool(false))
	}
	return 1
}

// Create creates a new workdir
func Create(L *lua.LState) int {
	// Get workdir path (required argument)
	workdirPath := L.CheckString(1)

	// Create directory with all parent directories
	if err := os.MkdirAll(workdirPath, 0755); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("workdir created successfully"))
	return 2
}

// Loader returns the workdir module loader
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get":     Get,
		"cleanup": Cleanup,
		"exists":  Exists,
		"create":  Create,
	})
	L.Push(mod)
	return 1
}

// Open registers the workdir module as a global
func Open(L *lua.LState) {
	// Register as global module (like exec, fs, etc.)
	workdirMt := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get":     Get,
		"cleanup": Cleanup,
		"exists":  Exists,
		"create":  Create,
	})
	L.SetGlobal("workdir", workdirMt)
}

// CreateRuntimeWorkdirObjectWithColonSupport creates workdir object supporting this.workdir.method() syntax (exported)
func CreateRuntimeWorkdirObjectWithColonSupport(L *lua.LState, workdirPath string) *lua.LUserData {
	return createRuntimeWorkdirObjectWithColonSupport(L, workdirPath)
}

// createRuntimeWorkdirObjectWithColonSupport creates workdir object supporting this:workdir:method() syntax
func createRuntimeWorkdirObjectWithColonSupport(L *lua.LState, workdirPath string) *lua.LUserData {
	workdirUD := L.NewUserData()
	workdirUD.Value = workdirPath

	// Create metatable for workdir object with colon syntax support
	workdirMt := L.NewTypeMetatable("RuntimeWorkdirColonSupport")
	L.SetField(workdirMt, "__index", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)

		workdirPath, ok := ud.Value.(string)
		if !ok {
			L.ArgError(1, "RuntimeWorkdirColonSupport expected")
			return 0
		}

		switch key {
		case "get":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath != "" {
					L.Push(lua.LString(workdirPath))
				} else {
					if cwd, err := os.Getwd(); err == nil {
						L.Push(lua.LString(cwd))
					} else {
						L.Push(lua.LString("/tmp"))
					}
				}
				return 1
			}))
		case "ensure":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}

				// Remove existing directory
				os.RemoveAll(workdirPath)

				// Create new directory
				if err := os.MkdirAll(workdirPath, 0755); err != nil {
					L.Push(lua.LBool(false))
					return 1
				}

				L.Push(lua.LBool(true))
				return 1
			}))
		case "exists":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}

				if _, err := os.Stat(workdirPath); err == nil {
					L.Push(lua.LBool(true))
				} else {
					L.Push(lua.LBool(false))
				}
				return 1
			}))
		case "cleanup":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}

				slog.Warn("Manual workdir cleanup is disabled. Workdir preserved.", "workdir", workdirPath)
				// if err := os.RemoveAll(workdirPath); err != nil {
				// 	L.Push(lua.LBool(false))
				// 	return 1
				// }

				L.Push(lua.LBool(true))
				return 1
			}))
		case "recreate":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}

				// Remove and recreate
				os.RemoveAll(workdirPath)
				if err := os.MkdirAll(workdirPath, 0755); err != nil {
					L.Push(lua.LBool(false))
					return 1
				}

				L.Push(lua.LBool(true))
				return 1
			}))
		default:
			L.Push(lua.LNil)
		}
		return 1
	}))

	L.SetMetatable(workdirUD, workdirMt)
	return workdirUD
}

// CreateRuntimeWorkdirObject creates workdir object for runtime execution
func CreateRuntimeWorkdirObject(L *lua.LState, workdirPath string) *lua.LUserData {
	workdirUD := L.NewUserData()
	workdirUD.Value = workdirPath

	// Create metatable for workdir object
	workdirMt := L.NewTypeMetatable("RuntimeWorkdir")
	L.SetField(workdirMt, "__index", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)

		workdirPath, ok := ud.Value.(string)
		if !ok {
			L.ArgError(1, "RuntimeWorkdir expected")
			return 0
		}

		switch key {
		case "get":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath != "" {
					L.Push(lua.LString(workdirPath))
				} else {
					if cwd, err := os.Getwd(); err == nil {
						L.Push(lua.LString(cwd))
					} else {
						L.Push(lua.LString("/tmp"))
					}
				}
				return 1
			}))
		case "exists":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}

				if _, err := os.Stat(workdirPath); err == nil {
					L.Push(lua.LBool(true))
				} else {
					L.Push(lua.LBool(false))
				}
				return 1
			}))
		case "ensure":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					return 1
				}

				// Remove existing directory
				os.RemoveAll(workdirPath)

				// Create new directory
				if err := os.MkdirAll(workdirPath, 0755); err != nil {
					L.Push(lua.LBool(false))
					return 1
				}

				L.Push(lua.LBool(true))
				return 1
			}))
		case "cleanup":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					L.Push(lua.LString("no workdir specified"))
					return 2
				}

				slog.Warn("Manual workdir cleanup is disabled. Workdir preserved.", "workdir", workdirPath)
				// if err := os.RemoveAll(workdirPath); err != nil {
				// 	L.Push(lua.LBool(false))
				// 	L.Push(lua.LString(err.Error()))
				// 	return 2
				// }

				L.Push(lua.LBool(true))
				L.Push(lua.LString("workdir cleanup is disabled"))
				return 2
			}))
		case "recreate":
			L.Push(L.NewFunction(func(L *lua.LState) int {
				if workdirPath == "" {
					L.Push(lua.LBool(false))
					L.Push(lua.LString("no workdir specified"))
					return 2
				}

				// Remove and recreate
				os.RemoveAll(workdirPath)
				if err := os.MkdirAll(workdirPath, 0755); err != nil {
					L.Push(lua.LBool(false))
					L.Push(lua.LString(err.Error()))
					return 2
				}

				L.Push(lua.LBool(true))
				L.Push(lua.LString("workdir recreated successfully"))
				return 2
			}))
		default:
			L.Push(lua.LNil)
		}
		return 1
	}))

	L.SetMetatable(workdirUD, workdirMt)
	return workdirUD
}
