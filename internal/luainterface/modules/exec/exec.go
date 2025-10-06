package exec

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// ExecCommand is a variable that can be overridden for testing
var ExecCommand = exec.Command

// SSH execution helpers (will be set by luainterface package)
var (
	IsSSHExecutionEnabled func() bool
	GetSSHProfile         func() string
	ExecuteCommandWithSSH func(string) (string, error)
)

// Run executes a command locally or via SSH
func Run(L *lua.LState) int {
	commandStr := L.CheckString(1)
	opts := L.OptTable(2, L.NewTable())

	ctx := L.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Check if SSH execution is enabled
	if IsSSHExecutionEnabled != nil && IsSSHExecutionEnabled() {
		slog.Debug("executing command via SSH", "source", "lua", "command", commandStr, "profile", GetSSHProfile())

		// For SSH execution, we need to handle workdir by prepending cd command
		if workdir := opts.RawGetString("workdir"); workdir.Type() == lua.LTString {
			commandStr = fmt.Sprintf("cd %s && %s", workdir.String(), commandStr)
		}

		// For SSH, environment variables need to be exported in the command
		if envTbl := opts.RawGetString("env"); envTbl.Type() == lua.LTTable {
			var envExports []string
			envTbl.(*lua.LTable).ForEach(func(key, value lua.LValue) {
				envExports = append(envExports, fmt.Sprintf("export %s='%s'", key.String(), value.String()))
			})
			if len(envExports) > 0 {
				commandStr = strings.Join(envExports, "; ") + "; " + commandStr
			}
		}

		// Execute via SSH
		output, err := ExecuteCommandWithSSH(commandStr)

		// SSH combines stdout and stderr
		stdoutStr := output
		stderrStr := ""
		if err != nil {
			stderrStr = err.Error()
		}

		if stdoutStr != "" {
			slog.Info(stdoutStr, "source", "lua-ssh", "stream", "stdout", "profile", GetSSHProfile())
		}
		if stderrStr != "" {
			slog.Warn(stderrStr, "source", "lua-ssh", "stream", "stderr", "profile", GetSSHProfile())
		}

		L.Push(lua.LString(stdoutStr))
		L.Push(lua.LString(stderrStr))
		L.Push(lua.LBool(err != nil))
		return 3
	}

	// Local execution
	slog.Debug("executing command locally", "source", "lua", "command", commandStr)

	cmd := ExecCommand("bash", "-c", commandStr)

	// Start with a minimal, controlled environment
	cmd.Env = []string{
		"PATH=/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin", // Set a default PATH
		"HOME=" + os.Getenv("HOME"),                         // Keep HOME if it exists
	}

	// Set workdir from options
	if workdir := opts.RawGetString("workdir"); workdir.Type() == lua.LTString {
		cmd.Dir = workdir.String()
	}

	// Add environment variables from options
	if envTbl := opts.RawGetString("env"); envTbl.Type() == lua.LTTable {
		envMap := make(map[string]string)
		for _, envVar := range cmd.Env {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}
		envTbl.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			envMap[key.String()] = value.String()
		})
		cmd.Env = []string{}
		for k, v := range envMap {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if stdoutStr != "" {
		slog.Info(stdoutStr, "source", "lua", "stream", "stdout")
	}
	if stderrStr != "" {
		slog.Warn(stderrStr, "source", "lua", "stream", "stderr")
	}

	L.Push(lua.LString(stdoutStr))
	L.Push(lua.LString(stderrStr))
	L.Push(lua.LBool(err != nil))
	return 3
}

// Loader returns the exec module loader
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"run": Run,
	})
	L.Push(mod)
	return 1
}

// Open registers the exec module and loads it globally
func Open(L *lua.LState) {
	L.PreloadModule("exec", Loader)
	if err := L.DoString(`exec = require("exec")`); err != nil {
		panic(err)
	}
}
