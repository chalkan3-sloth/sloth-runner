package core

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestValidateModule_Info(t *testing.T) {
	module := NewValidateModule()
	info := module.Info()

	if info.Name != "validate" {
		t.Errorf("Expected module name 'validate', got '%s'", info.Name)
	}

	if info.Category != "core" {
		t.Errorf("Expected category 'core', got '%s'", info.Category)
	}
}

func TestValidateModule_Email(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid email", "user@example.com", true},
		{"Valid with subdomain", "user@mail.example.com", true},
		{"Valid with plus", "user+tag@example.com", true},
		{"Invalid missing @", "userexample.com", false},
		{"Invalid missing domain", "user@", false},
		{"Invalid missing local", "@example.com", false},
		{"Invalid double @", "user@@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			module := NewValidateModule()
			L.PreloadModule("validate", module.Loader)

			code := `
				local validate = require("validate")
				local result = validate.email("` + tt.email + `")
				return result.valid
			`

			if err := L.DoString(code); err != nil {
				t.Fatalf("Lua execution failed: %v", err)
			}

			result := L.Get(-1)
			isValid := bool(result.(lua.LBool))

			if isValid != tt.expected {
				t.Errorf("Expected valid=%v, got %v", tt.expected, isValid)
			}
		})
	}
}

func TestValidateModule_URL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"Valid HTTP", "http://example.com", true},
		{"Valid HTTPS", "https://example.com", true},
		{"Valid with path", "https://example.com/path/to/page", true},
		{"Valid with query", "https://example.com?key=value", true},
		{"Invalid no scheme", "example.com", false},
		{"Invalid no host", "http://", false},
		{"Invalid malformed", "ht!tp://example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			module := NewValidateModule()
			L.PreloadModule("validate", module.Loader)

			code := `
				local validate = require("validate")
				local result = validate.url("` + tt.url + `")
				return result.valid
			`

			if err := L.DoString(code); err != nil {
				t.Fatalf("Lua execution failed: %v", err)
			}

			result := L.Get(-1)
			isValid := bool(result.(lua.LBool))

			if isValid != tt.expected {
				t.Errorf("Expected valid=%v, got %v", tt.expected, isValid)
			}
		})
	}
}

func TestValidateModule_IP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		version  string
		expected bool
	}{
		{"Valid IPv4", "192.168.1.1", "any", true},
		{"Valid IPv4 with v4", "10.0.0.1", "v4", true},
		{"Valid IPv4 localhost", "127.0.0.1", "any", true},
		{"Invalid IPv4 format", "192.168.1", "v4", false},
		{"Valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "v6", true},
		{"Valid IPv6 short", "::1", "v6", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			module := NewValidateModule()
			L.PreloadModule("validate", module.Loader)

			code := `
				local validate = require("validate")
				local result = validate.ip("` + tt.ip + `", "` + tt.version + `")
				return result.valid
			`

			if err := L.DoString(code); err != nil {
				t.Fatalf("Lua execution failed: %v", err)
			}

			result := L.Get(-1)
			isValid := bool(result.(lua.LBool))

			if isValid != tt.expected {
				t.Errorf("Expected valid=%v, got %v", tt.expected, isValid)
			}
		})
	}
}

func TestValidateModule_Regex(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		
		-- Test valid pattern match
		local result1 = validate.regex("test123", "^[a-z]+[0-9]+$")
		if not result1.valid then
			error("Expected pattern to match")
		end
		
		-- Test invalid pattern match
		local result2 = validate.regex("123test", "^[a-z]+[0-9]+$")
		if result2.valid then
			error("Expected pattern not to match")
		end
		
		-- Test with capture groups
		local result3 = validate.regex("John Doe", "^([A-Z][a-z]+) ([A-Z][a-z]+)$")
		if not result3.valid then
			error("Expected pattern with groups to match")
		end
		
		if not result3.matches then
			error("Expected matches to be present")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestValidateModule_Length(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		
		-- Test minimum length
		local result1 = validate.length("hello", {min = 3})
		if not result1.valid then
			error("Expected to pass min length")
		end
		
		local result2 = validate.length("hi", {min = 3})
		if result2.valid then
			error("Expected to fail min length")
		end
		
		-- Test maximum length
		local result3 = validate.length("hello", {max = 10})
		if not result3.valid then
			error("Expected to pass max length")
		end
		
		local result4 = validate.length("hello world test", {max = 10})
		if result4.valid then
			error("Expected to fail max length")
		end
		
		-- Test exact length
		local result5 = validate.length("12345", {exact = 5})
		if not result5.valid then
			error("Expected to pass exact length")
		end
		
		local result6 = validate.length("1234", {exact = 5})
		if result6.valid then
			error("Expected to fail exact length")
		end
		
		-- Test min and max together
		local result7 = validate.length("hello", {min = 3, max = 10})
		if not result7.valid then
			error("Expected to pass min-max range")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestValidateModule_Range(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		
		-- Test minimum value
		local result1 = validate.range("10", {min = "5"})
		if not result1.valid then
			error("Expected to pass min value")
		end
		
		local result2 = validate.range("3", {min = "5"})
		if result2.valid then
			error("Expected to fail min value")
		end
		
		-- Test maximum value
		local result3 = validate.range("10", {max = "20"})
		if not result3.valid then
			error("Expected to pass max value")
		end
		
		local result4 = validate.range("25", {max = "20"})
		if result4.valid then
			error("Expected to fail max value")
		end
		
		-- Test range
		local result5 = validate.range("15", {min = "10", max = "20"})
		if not result5.valid then
			error("Expected to pass range")
		end
		
		local result6 = validate.range("5", {min = "10", max = "20"})
		if result6.valid then
			error("Expected to fail range (below min)")
		end
		
		-- Test float values
		local result7 = validate.range("3.14", {min = "3", max = "4"})
		if not result7.valid then
			error("Expected to pass float range")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestValidateModule_Required(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		
		-- Test all fields present
		local data1 = {
			name = "John",
			email = "john@example.com",
			age = "30"
		}
		local result1 = validate.required(data1, {"name", "email", "age"})
		if not result1.valid then
			error("Expected all required fields to be present")
		end
		
		-- Test missing field
		local data2 = {
			name = "John",
			age = "30"
		}
		local result2 = validate.required(data2, {"name", "email", "age"})
		if result2.valid then
			error("Expected to fail with missing field")
		end
		
		if not result2.missing then
			error("Expected missing field list")
		end
		
		-- Test empty string as missing
		local data3 = {
			name = "John",
			email = "",
			age = "30"
		}
		local result3 = validate.required(data3, {"name", "email", "age"})
		if result3.valid then
			error("Expected to fail with empty field")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestValidateModule_Schema(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		
		local schema = {
			name = {
				required = "true",
				type = "string"
			},
			age = {
				required = "true",
				type = "number"
			},
			email = {
				required = "false",
				type = "string"
			}
		}
		
		-- Test valid data
		local data1 = {
			name = "John",
			age = 30,
			email = "john@example.com"
		}
		local result1 = validate.schema(data1, schema)
		if not result1.valid then
			error("Expected valid data to pass schema validation")
		end
		
		-- Test missing required field
		local data2 = {
			name = "John"
		}
		local result2 = validate.schema(data2, schema)
		if result2.valid then
			error("Expected to fail with missing required field")
		end
		
		-- Test wrong type
		local data3 = {
			name = "John",
			age = "not a number"
		}
		local result3 = validate.schema(data3, schema)
		if result3.valid then
			error("Expected to fail with wrong type")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestValidateModule_Sanitize(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		
		-- Test HTML sanitization
		local result1 = validate.sanitize("<script>alert('xss')</script>", "html")
		if not string.find(result1.sanitized, "&lt;script&gt;") then
			error("Expected HTML to be escaped")
		end
		
		-- Test SQL sanitization
		local result2 = validate.sanitize("'; DROP TABLE users; --", "sql")
		if not string.find(result2.sanitized, "''") then
			error("Expected SQL to be escaped")
		end
		
		-- Test trim
		local result3 = validate.sanitize("  hello world  ", "trim")
		if result3.sanitized ~= "hello world" then
			error("Expected whitespace to be trimmed: got '" .. result3.sanitized .. "'")
		end
		
		-- Test lowercase
		local result4 = validate.sanitize("HELLO WORLD", "lower")
		if result4.sanitized ~= "hello world" then
			error("Expected lowercase conversion")
		end
		
		-- Test uppercase
		local result5 = validate.sanitize("hello world", "upper")
		if result5.sanitized ~= "HELLO WORLD" then
			error("Expected uppercase conversion")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestValidateModule_InvalidRegexPattern(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		local result = validate.regex("test", "[invalid(")
		
		if not result.error then
			error("Expected error for invalid regex pattern")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestValidateModule_InvalidNumericValue(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewValidateModule()
	L.PreloadModule("validate", module.Loader)

	code := `
		local validate = require("validate")
		local result = validate.range("not a number", {min = "0", max = "100"})
		
		if not result.error then
			error("Expected error for invalid numeric value")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}
