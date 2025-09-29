package core

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/yuin/gopher-lua"
)

// SystemModule provides system operations
type SystemModule struct {
	info CoreModuleInfo
}

// NewSystemModule creates a new system module
func NewSystemModule() *SystemModule {
	info := CoreModuleInfo{
		Name:         "system",
		Version:      "1.0.0",
		Description:  "System operations including process management, file operations, and system information",
		Author:       "Sloth Runner Team",
		Category:     "core",
		Dependencies: []string{},
	}

	return &SystemModule{
		info: info,
	}
}

// Info returns module information
func (s *SystemModule) Info() CoreModuleInfo {
	return s.info
}

// Loader loads the system module into Lua
func (s *SystemModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"exec":        s.luaExec,
		"exec_async":  s.luaExecAsync,
		"kill":        s.luaKill,
		"exists":      s.luaExists,
		"mkdir":       s.luaMkdir,
		"rmdir":       s.luaRmdir,
		"copy":        s.luaCopy,
		"move":        s.luaMove,
		"chmod":       s.luaChmod,
		"stat":        s.luaStat,
		"env":         s.luaEnv,
		"setenv":      s.luaSetenv,
		"unsetenv":    s.luaUnsetenv,
		"hostname":    s.luaHostname,
		"platform":    s.luaPlatform,
		"arch":        s.luaArch,
		"cpu_count":   s.luaCPUCount,
		"memory":      s.luaMemory,
		"uptime":      s.luaUptime,
		"processes":   s.luaProcesses,
		"pwd":         s.luaPwd,
		"cd":          s.luaCd,
		"which":       s.luaWhich,
		"temp_dir":    s.luaTempDir,
		"home_dir":    s.luaHomeDir,
	})

	L.Push(mod)
	return 1
}

// luaExec executes a command synchronously
func (s *SystemModule) luaExec(L *lua.LState) int {
	command := L.CheckString(1)
	var args []string
	
	if L.GetTop() > 1 {
		if argsTable := L.CheckTable(2); argsTable != nil {
			argsTable.ForEach(func(_, value lua.LValue) {
				if str, ok := value.(lua.LString); ok {
					args = append(args, string(str))
				}
			})
		}
	}

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	result := L.NewTable()
	result.RawSetString("output", lua.LString(string(output)))
	result.RawSetString("success", lua.LBool(err == nil))
	
	if err != nil {
		result.RawSetString("error", lua.LString(err.Error()))
		if exitError, ok := err.(*exec.ExitError); ok {
			result.RawSetString("exit_code", lua.LNumber(exitError.ExitCode()))
		}
	} else {
		result.RawSetString("exit_code", lua.LNumber(0))
	}

	L.Push(result)
	return 1
}

// luaExecAsync executes a command asynchronously
func (s *SystemModule) luaExecAsync(L *lua.LState) int {
	command := L.CheckString(1)
	var args []string
	
	if L.GetTop() > 1 {
		if argsTable := L.CheckTable(2); argsTable != nil {
			argsTable.ForEach(func(_, value lua.LValue) {
				if str, ok := value.(lua.LString); ok {
					args = append(args, string(str))
				}
			})
		}
	}

	cmd := exec.Command(command, args...)
	err := cmd.Start()

	result := L.NewTable()
	if err != nil {
		result.RawSetString("success", lua.LBool(false))
		result.RawSetString("error", lua.LString(err.Error()))
		result.RawSetString("pid", lua.LNumber(0))
	} else {
		result.RawSetString("success", lua.LBool(true))
		result.RawSetString("pid", lua.LNumber(cmd.Process.Pid))
		
		// Start a goroutine to wait for completion and cleanup
		go func() {
			cmd.Wait()
		}()
	}

	L.Push(result)
	return 1
}

// luaKill kills a process by PID
func (s *SystemModule) luaKill(L *lua.LState) int {
	pid := int(L.CheckNumber(1))
	signal := syscall.SIGTERM // default signal
	
	if L.GetTop() > 1 {
		signalStr := L.CheckString(2)
		switch strings.ToUpper(signalStr) {
		case "TERM", "SIGTERM":
			signal = syscall.SIGTERM
		case "KILL", "SIGKILL":
			signal = syscall.SIGKILL
		case "INT", "SIGINT":
			signal = syscall.SIGINT
		case "HUP", "SIGHUP":
			signal = syscall.SIGHUP
		}
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = process.Signal(signal)
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaExists checks if a file or directory exists
func (s *SystemModule) luaExists(L *lua.LState) int {
	path := L.CheckString(1)
	_, err := os.Stat(path)
	L.Push(lua.LBool(err == nil))
	return 1
}

// luaMkdir creates a directory
func (s *SystemModule) luaMkdir(L *lua.LState) int {
	path := L.CheckString(1)
	recursive := false
	
	if L.GetTop() > 1 {
		recursive = bool(L.CheckBool(2))
	}

	var err error
	if recursive {
		err = os.MkdirAll(path, 0755)
	} else {
		err = os.Mkdir(path, 0755)
	}

	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaRmdir removes a directory
func (s *SystemModule) luaRmdir(L *lua.LState) int {
	path := L.CheckString(1)
	recursive := false
	
	if L.GetTop() > 1 {
		recursive = bool(L.CheckBool(2))
	}

	var err error
	if recursive {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaCopy copies a file
func (s *SystemModule) luaCopy(L *lua.LState) int {
	src := L.CheckString(1)
	dst := L.CheckString(2)

	sourceFile, err := os.Open(src)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer destFile.Close()

	_, err = sourceFile.WriteTo(destFile)
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaMove moves/renames a file
func (s *SystemModule) luaMove(L *lua.LState) int {
	src := L.CheckString(1)
	dst := L.CheckString(2)
	
	err := os.Rename(src, dst)
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaChmod changes file permissions
func (s *SystemModule) luaChmod(L *lua.LState) int {
	path := L.CheckString(1)
	modeStr := L.CheckString(2)
	
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("Invalid mode: " + err.Error()))
		return 2
	}
	
	err = os.Chmod(path, os.FileMode(mode))
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaStat gets file information
func (s *SystemModule) luaStat(L *lua.LState) int {
	path := L.CheckString(1)
	
	info, err := os.Stat(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	result := L.NewTable()
	result.RawSetString("name", lua.LString(info.Name()))
	result.RawSetString("size", lua.LNumber(info.Size()))
	result.RawSetString("mode", lua.LString(fmt.Sprintf("%o", info.Mode())))
	result.RawSetString("mod_time", lua.LNumber(info.ModTime().Unix()))
	result.RawSetString("is_dir", lua.LBool(info.IsDir()))

	L.Push(result)
	return 1
}

// luaEnv gets environment variable
func (s *SystemModule) luaEnv(L *lua.LState) int {
	key := L.CheckString(1)
	value := os.Getenv(key)
	if value == "" {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(value))
	}
	return 1
}

// luaSetenv sets environment variable
func (s *SystemModule) luaSetenv(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)
	
	err := os.Setenv(key, value)
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaUnsetenv unsets environment variable
func (s *SystemModule) luaUnsetenv(L *lua.LState) int {
	key := L.CheckString(1)
	
	err := os.Unsetenv(key)
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaHostname gets the hostname
func (s *SystemModule) luaHostname(L *lua.LState) int {
	hostname, err := os.Hostname()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(hostname))
	return 1
}

// luaPlatform gets the platform (OS)
func (s *SystemModule) luaPlatform(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOOS))
	return 1
}

// luaArch gets the architecture
func (s *SystemModule) luaArch(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOARCH))
	return 1
}

// luaCPUCount gets the number of CPUs
func (s *SystemModule) luaCPUCount(L *lua.LState) int {
	L.Push(lua.LNumber(runtime.NumCPU()))
	return 1
}

// luaMemory gets memory statistics (simplified)
func (s *SystemModule) luaMemory(L *lua.LState) int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	result := L.NewTable()
	result.RawSetString("alloc", lua.LNumber(m.Alloc))
	result.RawSetString("total_alloc", lua.LNumber(m.TotalAlloc))
	result.RawSetString("sys", lua.LNumber(m.Sys))
	result.RawSetString("num_gc", lua.LNumber(m.NumGC))
	
	L.Push(result)
	return 1
}

// luaUptime gets system uptime (simplified - process uptime)
func (s *SystemModule) luaUptime(L *lua.LState) int {
	// This is a simplified version - actual system uptime would require OS-specific code
	L.Push(lua.LNumber(time.Since(time.Now().Add(-time.Hour)).Seconds())) // placeholder
	return 1
}

// luaProcesses lists running processes (simplified)
func (s *SystemModule) luaProcesses(L *lua.LState) int {
	// This is a placeholder - actual process listing would require OS-specific code
	result := L.NewTable()
	L.Push(result)
	return 1
}

// luaPwd gets current working directory
func (s *SystemModule) luaPwd(L *lua.LState) int {
	pwd, err := os.Getwd()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(pwd))
	return 1
}

// luaCd changes current working directory
func (s *SystemModule) luaCd(L *lua.LState) int {
	dir := L.CheckString(1)
	
	err := os.Chdir(dir)
	L.Push(lua.LBool(err == nil))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// luaWhich finds executable in PATH
func (s *SystemModule) luaWhich(L *lua.LState) int {
	command := L.CheckString(1)
	
	path, err := exec.LookPath(command)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(path))
	return 1
}

// luaTempDir gets temporary directory
func (s *SystemModule) luaTempDir(L *lua.LState) int {
	L.Push(lua.LString(os.TempDir()))
	return 1
}

// luaHomeDir gets user home directory
func (s *SystemModule) luaHomeDir(L *lua.LState) int {
	home, err := os.UserHomeDir()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(home))
	return 1
}