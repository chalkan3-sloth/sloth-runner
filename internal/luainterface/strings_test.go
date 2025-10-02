package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestStringsModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterStringModule(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "strings.upper()",
			script: `
				local result = strings.upper("hello world")
				assert(result == "HELLO WORLD", "upper() should convert to uppercase")
			`,
		},
		{
			name: "strings.lower()",
			script: `
				local result = strings.lower("HELLO WORLD")
				assert(result == "hello world", "lower() should convert to lowercase")
			`,
		},
		{
			name: "strings.trim()",
			script: `
				local result = strings.trim("  hello  ")
				assert(result == "hello", "trim() should remove leading and trailing spaces")
			`,
		},
// 		{
// 			name: "strings.trim_prefix()",
// 			script: `
// 				local result = strings.trim_prefix("hello world", "hello ")
// 				assert(result == "world", "trim_prefix() should remove prefix")
// 			`,
// 		},
// 		{
// 			name: "strings.trim_suffix()",
// 			script: `
// 				local result = strings.trim_suffix("hello world", " world")
// 				assert(result == "hello", "trim_suffix() should remove suffix")
// 			`,
// 		},
		{
			name: "strings.split()",
			script: `
				local parts = strings.split("a,b,c", ",")
				assert(#parts == 3, "split() should return 3 parts")
				assert(parts[1] == "a", "first part should be 'a'")
				assert(parts[2] == "b", "second part should be 'b'")
				assert(parts[3] == "c", "third part should be 'c'")
			`,
		},
		{
			name: "strings.join()",
			script: `
				local result = strings.join({["1"]="a", ["2"]="b", ["3"]="c"}, ",")
				assert(type(result) == "string", "join() should return a string")
			`,
		},
		{
			name: "strings.contains()",
			script: `
				local result = strings.contains("hello world", "world")
				assert(result == true, "contains() should return true when substring exists")
				local result2 = strings.contains("hello world", "xyz")
				assert(result2 == false, "contains() should return false when substring doesn't exist")
			`,
		},
		{
			name: "strings.starts_with()",
			script: `
				local result = strings.starts_with("hello world", "hello")
				assert(result == true, "has_prefix() should return true when prefix exists")
				local result2 = strings.starts_with("hello world", "world")
				assert(result2 == false, "has_prefix() should return false when prefix doesn't exist")
			`,
		},
		{
			name: "strings.ends_with()",
			script: `
				local result = strings.ends_with("hello world", "world")
				assert(result == true, "has_suffix() should return true when suffix exists")
				local result2 = strings.ends_with("hello world", "hello")
				assert(result2 == false, "has_suffix() should return false when suffix doesn't exist")
			`,
		},
		{
			name: "strings.replace()",
			script: `
				local result = strings.replace("hello world", "world", "there")
				assert(result == "hello there", "replace() should replace substring")
			`,
		},
		{
			name: "strings.replace_regex()",
			script: `
				local result = strings.replace_regex("foo bar foo", "foo", "baz")
				assert(result == "baz bar baz", "replace_all() should replace all occurrences")
			`,
		},
// 		{
// 			name: "strings.repeat()",
// 			script: `
// 				local result = strings.repeat("abc", 3)
// 				assert(result == "abcabcabc", "repeat() should repeat string n times")
// 			`,
// 		},
// 		{
// 			name: "strings.index()",
// 			script: `
// 				local pos = strings.index("hello world", "world")
// 				assert(pos >= 0, "index() should return position of substring")
// 			`,
// 		},
// 		{
// 			name: "strings.last_index()",
// 			script: `
// 				local pos = strings.last_index("hello world world", "world")
// 				assert(pos >= 0, "last_index() should return last position of substring")
// 			`,
// 		},
// 		{
// 			name: "strings.count()",
// 			script: `
// 				local cnt = strings.count("hello world world", "world")
// 				assert(cnt == 2, "count() should return number of occurrences")
// 			`,
// 		},
		{
			name: "strings.title()",
			script: `
				local result = strings.title("hello world")
				assert(result == "Hello World", "title() should capitalize first letter of each word")
			`,
		},
// 		{
// 			name: "strings.reverse()",
// 			script: `
// 				local result = strings.reverse("hello")
// 				assert(result == "olleh", "reverse() should reverse the string")
// 			`,
// 		},
// 		{
// 			name: "strings.len()",
// 			script: `
// 				local length = strings.len("hello")
// 				assert(length == 5, "len() should return string length")
// 			`,
// 		},
// 		{
// 			name: "strings.substr()",
// 			script: `
// 				local result = strings.substr("hello world", 0, 5)
// 				assert(result == "hello", "substr() should extract substring")
// 			`,
// 		},
		{
			name: "strings.match()",
			script: `
				local result = strings.match("test123", "[0-9]+")
				assert(type(result) == "table", "match() should return a table for matching pattern")
				assert(result[1] == "123", "match() should return the matched string")
			`,
		},
		{
			name: "strings.match_all()",
			script: `
				local matches = strings.match_all("test 123 foo 456", "[0-9]+")
				assert(type(matches) == "table", "find_all() should return a table")
			`,
		},
		{
			name: "strings.pad_left()",
			script: `
				local result = strings.pad_left("test", 10, " ")
				assert(#result == 10, "pad_left() should pad to specified length")
			`,
		},
		{
			name: "strings.pad_right()",
			script: `
				local result = strings.pad_right("test", 10, " ")
				assert(#result == 10, "pad_right() should pad to specified length")
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

func TestStringsEdgeCases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterStringModule(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "empty string",
			script: `
				local result = strings.upper("")
				assert(result == "", "upper() should handle empty string")
			`,
		},
		{
			name: "split with empty result",
			script: `
				local parts = strings.split("", ",")
				assert(type(parts) == "table", "split() should return a table even for empty string")
			`,
		},
		{
			name: "replace not found",
			script: `
				local result = strings.replace("hello", "xyz", "abc")
				assert(result == "hello", "replace() should return original if substring not found")
			`,
		},
		{
			name: "count not found",
			script: `
				local cnt = strings.count("hello", "xyz")
				assert(cnt == 0, "count() should return 0 if substring not found")
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
