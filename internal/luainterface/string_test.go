package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestStringModule(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{
			name: "string.trim",
			script: `
				local s = string.trim("  hello  ")
				assert(s == "hello", "should trim whitespace")
			`,
		},
		{
			name: "string.trim_prefix",
			script: `
				local s = string.trim_prefix("prefix_hello", "prefix_")
				assert(s == "hello", "should trim prefix")
			`,
		},
		{
			name: "string.trim_suffix",
			script: `
				local s = string.trim_suffix("hello_suffix", "_suffix")
				assert(s == "hello", "should trim suffix")
			`,
		},
		{
			name: "string.split",
			script: `
				local parts = string.split("a,b,c", ",")
				assert(type(parts) == "table", "should return a table")
				assert(#parts == 3, "should have 3 parts")
				assert(parts[1] == "a", "first part should be 'a'")
				assert(parts[2] == "b", "second part should be 'b'")
				assert(parts[3] == "c", "third part should be 'c'")
			`,
		},
		{
			name: "string.join",
			script: `
				local s = string.join({{"a", "b", "c"}}, ",")
				assert(type(s) == "string", "should return a string")
			`,
		},
		{
			name: "string.contains",
			script: `
				assert(string.contains("hello world", "world") == true, "should contain 'world'")
				assert(string.contains("hello world", "xyz") == false, "should not contain 'xyz'")
			`,
		},
		{
			name: "string.has_prefix",
			script: `
				assert(string.has_prefix("hello", "hel") == true, "should have prefix")
				assert(string.has_prefix("hello", "wor") == false, "should not have prefix")
			`,
		},
		{
			name: "string.has_suffix",
			script: `
				assert(string.has_suffix("hello", "lo") == true, "should have suffix")
				assert(string.has_suffix("hello", "ab") == false, "should not have suffix")
			`,
		},
		{
			name: "string.to_upper",
			script: `
				local s = string.to_upper("hello")
				assert(s == "HELLO", "should convert to uppercase")
			`,
		},
		{
			name: "string.to_lower",
			script: `
				local s = string.to_lower("HELLO")
				assert(s == "hello", "should convert to lowercase")
			`,
		},
		{
			name: "string.replace",
			script: `
				local s = string.replace("hello world", "world", "universe")
				assert(s == "hello universe", "should replace text")
			`,
		},
		{
			name: "string.replace_all",
			script: `
				local s = string.replace_all("aaa", "a", "b")
				assert(s == "bbb", "should replace all occurrences")
			`,
		},
		{
			name: "string.count",
			script: `
				local c = string.count("aaa", "a")
				assert(c == 3, "should count 3 occurrences")
			`,
		},
		{
			name: "string.repeat_str",
			script: `
				local s = string.repeat_str("ab", 3)
				assert(s == "ababab", "should repeat string 3 times")
			`,
		},
		{
			name: "string.index",
			script: `
				local i = string.index("hello", "l")
				assert(i == 2, "should find index of 'l'")
			`,
		},
		{
			name: "string.last_index",
			script: `
				local i = string.last_index("hello", "l")
				assert(i == 3, "should find last index of 'l'")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			loadStringModule(L)

			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestStringModuleRegex(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{
			name: "string.match",
			script: `
				local matches = string.match("hello123world456", "[0-9]+")
				assert(type(matches) == "table", "should return a table")
			`,
		},
		{
			name: "string.match_all",
			script: `
				local matches = string.match_all("hello123world456", "[0-9]+")
				assert(type(matches) == "table", "should return a table")
			`,
		},
		{
			name: "string.replace_regex",
			script: `
				local s = string.replace_regex("hello123", "[0-9]+", "XXX")
				assert(type(s) == "string", "should return a string")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			loadStringModule(L)

			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}
