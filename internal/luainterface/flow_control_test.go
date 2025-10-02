package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Lua control flow structures
func TestLuaConditionals(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		expected float64
	}{
		{
			name: "if true branch",
			script: `
				local result = 0
				if true then
					result = 10
				end
				return result
			`,
			expected: 10,
		},
		{
			name: "if false branch",
			script: `
				local result = 0
				if false then
					result = 10
				end
				return result
			`,
			expected: 0,
		},
		{
			name: "if else",
			script: `
				local result = 0
				if false then
					result = 10
				else
					result = 20
				end
				return result
			`,
			expected: 20,
		},
		{
			name: "if elseif else",
			script: `
				local x = 2
				local result = 0
				if x == 1 then
					result = 10
				elseif x == 2 then
					result = 20
				else
					result = 30
				end
				return result
			`,
			expected: 20,
		},
		{
			name: "nested if",
			script: `
				local a = true
				local b = true
				local result = 0
				if a then
					if b then
						result = 100
					end
				end
				return result
			`,
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			err := L.DoString(tt.script)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)
			assert.Equal(t, tt.expected, float64(result.(lua.LNumber)))
		})
	}
}

func TestLuaLoops(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		expected float64
	}{
		{
			name: "for loop",
			script: `
				local sum = 0
				for i = 1, 10 do
					sum = sum + i
				end
				return sum
			`,
			expected: 55,
		},
		{
			name: "while loop",
			script: `
				local sum = 0
				local i = 1
				while i <= 5 do
					sum = sum + i
					i = i + 1
				end
				return sum
			`,
			expected: 15,
		},
		{
			name: "repeat until",
			script: `
				local sum = 0
				local i = 1
				repeat
					sum = sum + i
					i = i + 1
				until i > 5
				return sum
			`,
			expected: 15,
		},
		{
			name: "break in loop",
			script: `
				local sum = 0
				for i = 1, 10 do
					if i == 5 then
						break
					end
					sum = sum + i
				end
				return sum
			`,
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			err := L.DoString(tt.script)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)
			assert.Equal(t, tt.expected, float64(result.(lua.LNumber)))
		})
	}
}

func TestLuaFunctions(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		expected float64
	}{
		{
			name: "simple function",
			script: `
				local function add(a, b)
					return a + b
				end
				return add(5, 3)
			`,
			expected: 8,
		},
		{
			name: "recursive function",
			script: `
				local function factorial(n)
					if n <= 1 then
						return 1
					end
					return n * factorial(n - 1)
				end
				return factorial(5)
			`,
			expected: 120,
		},
		{
			name: "closure",
			script: `
				local function counter()
					local count = 0
					return function()
						count = count + 1
						return count
					end
				end
				local c = counter()
				c()
				c()
				return c()
			`,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			err := L.DoString(tt.script)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)
			assert.Equal(t, tt.expected, float64(result.(lua.LNumber)))
		})
	}
}

func TestLuaTables(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		validate func(t *testing.T, L *lua.LState)
	}{
		{
			name: "table creation",
			script: `
				local t = {a = 1, b = 2, c = 3}
				return t.a + t.b + t.c
			`,
			validate: func(t *testing.T, L *lua.LState) {
				result := L.Get(-1)
				L.Pop(1)
				assert.Equal(t, 6.0, float64(result.(lua.LNumber)))
			},
		},
		{
			name: "array table",
			script: `
				local arr = {10, 20, 30}
				return arr[1] + arr[2] + arr[3]
			`,
			validate: func(t *testing.T, L *lua.LState) {
				result := L.Get(-1)
				L.Pop(1)
				assert.Equal(t, 60.0, float64(result.(lua.LNumber)))
			},
		},
		{
			name: "table.insert",
			script: `
				local t = {}
				table.insert(t, 1)
				table.insert(t, 2)
				table.insert(t, 3)
				return #t
			`,
			validate: func(t *testing.T, L *lua.LState) {
				result := L.Get(-1)
				L.Pop(1)
				assert.Equal(t, 3.0, float64(result.(lua.LNumber)))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			err := L.DoString(tt.script)
			require.NoError(t, err)

			tt.validate(t, L)
		})
	}
}

func TestLuaStringOperations(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		expected string
	}{
		{
			name:     "concatenation",
			script:   `return "hello" .. " " .. "world"`,
			expected: "hello world",
		},
		{
			name:     "string.upper",
			script:   `return string.upper("hello")`,
			expected: "HELLO",
		},
		{
			name:     "string.lower",
			script:   `return string.lower("WORLD")`,
			expected: "world",
		},
		{
			name:     "string.sub",
			script:   `return string.sub("hello", 1, 2)`,
			expected: "he",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			err := L.DoString(tt.script)
			require.NoError(t, err)

			result := L.Get(-1)
			L.Pop(1)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestLuaErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		script      string
		shouldError bool
	}{
		{
			name:        "syntax error",
			script:      `local x = `,
			shouldError: true,
		},
		{
			name: "pcall success",
			script: `
				local function safe()
					return "ok"
				end
				local status, result = pcall(safe)
				return result
			`,
			shouldError: false,
		},
		{
			name: "pcall catch error",
			script: `
				local function unsafe()
					error("test error")
				end
				local status, result = pcall(unsafe)
				if status then
					return "no error"
				else
					return "caught"
				end
			`,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			err := L.DoString(tt.script)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
