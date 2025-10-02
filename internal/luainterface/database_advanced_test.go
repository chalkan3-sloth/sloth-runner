package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestDatabaseModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDatabaseModule(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "db module exists",
			script: `
				assert(type(db) == "table", "db module should exist")
			`,
		},
		{
			name: "db.connect exists",
			script: `
				assert(type(db.connect) == "function", "db.connect should be a function")
			`,
		},
		{
			name: "db.close exists",
			script: `
				assert(type(db.close) == "function", "db.close should be a function")
			`,
		},
		{
			name: "db.query exists",
			script: `
				assert(type(db.query) == "function", "db.query should be a function")
			`,
		},
		{
			name: "db.exec exists",
			script: `
				assert(type(db.exec) == "function", "db.exec should be a function")
			`,
		},
		{
			name: "db.prepare exists",
			script: `
				assert(type(db.prepare) == "function", "db.prepare should be a function")
			`,
		},
		{
			name: "db.begin exists",
			script: `
				assert(type(db.begin) == "function", "db.begin should be a function")
			`,
		},
		{
			name: "db.commit exists",
			script: `
				assert(type(db.commit) == "function", "db.commit should be a function")
			`,
		},
		{
			name: "db.rollback exists",
			script: `
				assert(type(db.rollback) == "function", "db.rollback should be a function")
			`,
		},
		{
			name: "db.ping exists",
			script: `
				assert(type(db.ping) == "function", "db.ping should be a function")
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

func TestDatabaseModuleAPIs(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterDatabaseModule(L)

	script := `
		-- Test that core database methods exist
		assert(type(db.connect) == "function", "db.connect should be a function")
		assert(type(db.query) == "function", "db.query should be a function")
		assert(type(db.exec) == "function", "db.exec should be a function")
		assert(type(db.close) == "function", "db.close should be a function")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
