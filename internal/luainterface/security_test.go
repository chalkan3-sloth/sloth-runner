package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestSecurityModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSecurityModule(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "security module exists",
			script: `
				assert(type(security) == "table", "security module should exist")
			`,
		},
		{
			name: "security.hash exists",
			script: `
				assert(type(security.hash) == "function", "security.hash should be a function")
			`,
		},
		{
			name: "security.encrypt exists",
			script: `
				assert(type(security.encrypt) == "function", "security.encrypt should be a function")
			`,
		},
		{
			name: "security.decrypt exists",
			script: `
				assert(type(security.decrypt) == "function", "security.decrypt should be a function")
			`,
		},
		{
			name: "security.sign exists",
			script: `
				assert(type(security.sign) == "function", "security.sign should be a function")
			`,
		},
		{
			name: "security.verify exists",
			script: `
				assert(type(security.verify) == "function", "security.verify should be a function")
			`,
		},
		{
			name: "security.generate_key exists",
			script: `
				assert(type(security.generate_key) == "function", "security.generate_key should be a function")
			`,
		},
		{
			name: "security.validate_password exists",
			script: `
				assert(type(security.validate_password) == "function", "security.validate_password should be a function")
			`,
		},
		{
			name: "security.sanitize exists",
			script: `
				assert(type(security.sanitize) == "function", "security.sanitize should be a function")
			`,
		},
		{
			name: "security.escape_html exists",
			script: `
				assert(type(security.escape_html) == "function", "security.escape_html should be a function")
			`,
		},
		{
			name: "security.escape_sql exists",
			script: `
				assert(type(security.escape_sql) == "function", "security.escape_sql should be a function")
			`,
		},
		{
			name: "security.check_permission exists",
			script: `
				assert(type(security.check_permission) == "function", "security.check_permission should be a function")
			`,
		},
		{
			name: "security.rate_limit exists",
			script: `
				assert(type(security.rate_limit) == "function", "security.rate_limit should be a function")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestSecurityModuleAPIs(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSecurityModule(L)

	script := `
		-- Test that core security methods exist
		assert(type(security.hash) == "function", "security.hash should be a function")
		assert(type(security.encrypt) == "function", "security.encrypt should be a function")
		assert(type(security.decrypt) == "function", "security.decrypt should be a function")
		assert(type(security.sign) == "function", "security.sign should be a function")
		assert(type(security.verify) == "function", "security.verify should be a function")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
