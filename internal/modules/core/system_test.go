package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestSystemModule_Exists tests the exists function with various conditions
func TestSystemModule_Exists(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	mod := NewSystemModule()
	mod.Loader(L)
	systemTable := L.Get(-1).(*lua.LTable)

	tests := []struct {
		name           string
		setupPath      func() string
		cleanupPath    func(string)
		expectExists   bool
	}{
		{
			name: "existing file returns true",
			setupPath: func() string {
				tmpFile, err := os.CreateTemp("", "test_exists_*.txt")
				require.NoError(t, err)
				tmpFile.Close()
				return tmpFile.Name()
			},
			cleanupPath: func(path string) {
				os.Remove(path)
			},
			expectExists: true,
		},
		{
			name: "non-existing file returns false",
			setupPath: func() string {
				return "/tmp/this_file_definitely_does_not_exist_12345.txt"
			},
			cleanupPath: func(path string) {},
			expectExists: false,
		},
		{
			name: "existing directory returns true",
			setupPath: func() string {
				tmpDir, err := os.MkdirTemp("", "test_exists_dir_*")
				require.NoError(t, err)
				return tmpDir
			},
			cleanupPath: func(path string) {
				os.RemoveAll(path)
			},
			expectExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupPath()
			defer tt.cleanupPath(path)

			// Call system.exists(path)
			existsFunc := systemTable.RawGetString("exists").(*lua.LFunction)
			L.Push(existsFunc)
			L.Push(lua.LString(path))
			
			err := L.PCall(1, 1, nil)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)

			assert.Equal(t, lua.LBool(tt.expectExists), result)
		})
	}
}

// TestSystemModule_Mkdir tests the mkdir function with various conditions
func TestSystemModule_Mkdir(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	mod := NewSystemModule()
	mod.Loader(L)
	systemTable := L.Get(-1).(*lua.LTable)

	tests := []struct {
		name          string
		setupDir      func() string
		cleanupDir    func(string)
		expectSuccess bool
	}{
		{
			name: "create new directory succeeds",
			setupDir: func() string {
				tmpDir := filepath.Join(os.TempDir(), "test_mkdir_new_dir")
				os.RemoveAll(tmpDir) // Ensure it doesn't exist
				return tmpDir
			},
			cleanupDir: func(path string) {
				os.RemoveAll(path)
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupDir()
			defer tt.cleanupDir(path)

			// Call system.mkdir(path)
			mkdirFunc := systemTable.RawGetString("mkdir").(*lua.LFunction)
			L.Push(mkdirFunc)
			L.Push(lua.LString(path))
			
			err := L.PCall(1, 1, nil)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)

			if tt.expectSuccess {
				assert.Equal(t, lua.LTrue, result)
				// Verify directory was created
				_, err := os.Stat(path)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, lua.LFalse, result)
			}
		})
	}
}

// TestSystemModule_Exec tests the exec function with various conditions
func TestSystemModule_Exec(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	mod := NewSystemModule()
	mod.Loader(L)
	systemTable := L.Get(-1).(*lua.LTable)

	tests := []struct {
		name          string
		command       string
		args          []string
		expectSuccess bool
		checkOutput   func(string) bool
	}{
		{
			name:          "successful command returns success",
			command:       "echo",
			args:          []string{"hello"},
			expectSuccess: true,
			checkOutput: func(output string) bool {
				return len(output) > 0
			},
		},
		{
			name:          "failing command returns failure",
			command:       "false",
			args:          []string{},
			expectSuccess: false,
			checkOutput: func(output string) bool {
				return true // Any output is acceptable
			},
		},
		{
			name:          "command with exit 0 returns success",
			command:       "sh",
			args:          []string{"-c", "exit 0"},
			expectSuccess: true,
			checkOutput: func(output string) bool {
				return true
			},
		},
		{
			name:          "command with exit 1 returns failure",
			command:       "sh",
			args:          []string{"-c", "exit 1"},
			expectSuccess: false,
			checkOutput: func(output string) bool {
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call system.exec(command, args)
			execFunc := systemTable.RawGetString("exec").(*lua.LFunction)
			L.Push(execFunc)
			L.Push(lua.LString(tt.command))
			
			// Create args table
			argsTable := L.NewTable()
			for _, arg := range tt.args {
				argsTable.Append(lua.LString(arg))
			}
			L.Push(argsTable)
			
			err := L.PCall(2, 1, nil)
			require.NoError(t, err)

			result := L.Get(-1).(*lua.LTable)
			L.Pop(1)

			success := result.RawGetString("success")
			output := result.RawGetString("output")

			assert.Equal(t, lua.LBool(tt.expectSuccess), success)
			if tt.checkOutput != nil {
				assert.True(t, tt.checkOutput(output.String()))
			}
		})
	}
}

// TestSystemModule_Env tests the env function
func TestSystemModule_Env(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	mod := NewSystemModule()
	mod.Loader(L)
	systemTable := L.Get(-1).(*lua.LTable)

	tests := []struct {
		name        string
		envKey      string
		envValue    string
		setupEnv    bool
		expectNil   bool
	}{
		{
			name:      "existing env var returns value",
			envKey:    "TEST_ENV_VAR_EXISTS",
			envValue:  "test_value_123",
			setupEnv:  true,
			expectNil: false,
		},
		{
			name:      "non-existing env var returns nil",
			envKey:    "TEST_ENV_VAR_DOES_NOT_EXIST_XYZ",
			envValue:  "",
			setupEnv:  false,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			// Call system.env(key)
			envFunc := systemTable.RawGetString("env").(*lua.LFunction)
			L.Push(envFunc)
			L.Push(lua.LString(tt.envKey))
			
			err := L.PCall(1, 1, nil)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)

			if tt.expectNil {
				assert.Equal(t, lua.LNil, result)
			} else {
				assert.Equal(t, lua.LString(tt.envValue), result)
			}
		})
	}
}

// TestSystemModule_Which tests the which function with various conditions
func TestSystemModule_Which(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	mod := NewSystemModule()
	mod.Loader(L)
	systemTable := L.Get(-1).(*lua.LTable)

	tests := []struct {
		name          string
		command       string
		expectFound   bool
	}{
		{
			name:        "existing command returns path",
			command:     "sh",
			expectFound: true,
		},
		{
			name:        "non-existing command returns nil",
			command:     "this_command_definitely_does_not_exist_12345",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call system.which(command)
			whichFunc := systemTable.RawGetString("which").(*lua.LFunction)
			L.Push(whichFunc)
			L.Push(lua.LString(tt.command))
			
			err := L.PCall(1, 1, nil)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)

			if tt.expectFound {
				assert.NotEqual(t, lua.LNil, result)
				assert.NotEqual(t, lua.LString(""), result)
			} else {
				assert.Equal(t, lua.LNil, result)
			}
		})
	}
}
