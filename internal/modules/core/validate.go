package core

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/gopher-lua"
)

// ValidateModule provides data validation and sanitization
type ValidateModule struct {
	info CoreModuleInfo
}

// NewValidateModule creates a new validation module
func NewValidateModule() *ValidateModule {
	info := CoreModuleInfo{
		Name:        "validate",
		Version:     "1.0.0",
		Description: "Data validation and sanitization utilities",
		Author:      "Sloth Runner Team",
		Category:    "core",
		Dependencies: []string{},
	}
	
	return &ValidateModule{
		info: info,
	}
}

// Info returns module information
func (v *ValidateModule) Info() CoreModuleInfo {
	return v.info
}

// Loader returns the Lua loader function
func (m *ValidateModule) Loader(L *lua.LState) int {
	validateTable := L.NewTable()
	
	// Validation functions
	L.SetFuncs(validateTable, map[string]lua.LGFunction{
		"email":      m.luaValidateEmail,
		"url":        m.luaValidateURL,
		"ip":         m.luaValidateIP,
		"regex":      m.luaValidateRegex,
		"length":     m.luaValidateLength,
		"range":      m.luaValidateRange,
		"required":   m.luaValidateRequired,
		"schema":     m.luaValidateSchema,
		"sanitize":   m.luaSanitize,
	})
	
	L.Push(validateTable)
	return 1
}

// luaValidateEmail validates email addresses
func (m *ValidateModule) luaValidateEmail(L *lua.LState) int {
	email := L.CheckString(1)
	
	_, err := mail.ParseAddress(email)
	isValid := err == nil
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	if !isValid {
		result.RawSetString("error", lua.LString(err.Error()))
	}
	
	L.Push(result)
	return 1
}

// luaValidateURL validates URLs
func (m *ValidateModule) luaValidateURL(L *lua.LState) int {
	urlStr := L.CheckString(1)
	
	parsedURL, err := url.Parse(urlStr)
	isValid := err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	
	if isValid {
		result.RawSetString("scheme", lua.LString(parsedURL.Scheme))
		result.RawSetString("host", lua.LString(parsedURL.Host))
		result.RawSetString("path", lua.LString(parsedURL.Path))
		if parsedURL.RawQuery != "" {
			result.RawSetString("query", lua.LString(parsedURL.RawQuery))
		}
	} else if err != nil {
		result.RawSetString("error", lua.LString(err.Error()))
	} else {
		result.RawSetString("error", lua.LString("invalid URL format"))
	}
	
	L.Push(result)
	return 1
}

// luaValidateIP validates IP addresses
func (m *ValidateModule) luaValidateIP(L *lua.LState) int {
	ip := L.CheckString(1)
	version := L.OptString(2, "any") // "any", "v4", "v6"
	
	ipv4Regex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	ipv6Regex := regexp.MustCompile(`^([0-9a-fA-F]{0,4}:){2,7}[0-9a-fA-F]{0,4}$`)
	
	var isValid bool
	var detectedVersion string
	
	switch version {
	case "v4":
		isValid = ipv4Regex.MatchString(ip)
		detectedVersion = "v4"
	case "v6":
		isValid = ipv6Regex.MatchString(ip)
		detectedVersion = "v6"
	default: // "any"
		if ipv4Regex.MatchString(ip) {
			isValid = true
			detectedVersion = "v4"
		} else if ipv6Regex.MatchString(ip) {
			isValid = true
			detectedVersion = "v6"
		}
	}
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	if isValid {
		result.RawSetString("version", lua.LString(detectedVersion))
	} else {
		result.RawSetString("error", lua.LString("invalid IP address format"))
	}
	
	L.Push(result)
	return 1
}

// luaValidateRegex validates against a regular expression
func (m *ValidateModule) luaValidateRegex(L *lua.LState) int {
	value := L.CheckString(1)
	pattern := L.CheckString(2)
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		errorTable := L.NewTable()
		errorTable.RawSetString("error", lua.LString("invalid regex pattern"))
		errorTable.RawSetString("message", lua.LString(err.Error()))
		L.Push(errorTable)
		return 1
	}
	
	isValid := regex.MatchString(value)
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	
	if isValid {
		// Extract matches if any
		matches := regex.FindStringSubmatch(value)
		if len(matches) > 1 {
			matchesTable := L.NewTable()
			for i, match := range matches[1:] { // Skip full match
				matchesTable.RawSetInt(i+1, lua.LString(match))
			}
			result.RawSetString("matches", matchesTable)
		}
	}
	
	L.Push(result)
	return 1
}

// luaValidateLength validates string length
func (m *ValidateModule) luaValidateLength(L *lua.LState) int {
	value := L.CheckString(1)
	options := L.CheckTable(2)
	
	length := len(value)
	var isValid = true
	var errors []string
	
	if minVal := options.RawGetString("min"); minVal != lua.LNil {
		if min, err := strconv.Atoi(minVal.String()); err == nil {
			if length < min {
				isValid = false
				errors = append(errors, fmt.Sprintf("minimum length is %d", min))
			}
		}
	}
	
	if maxVal := options.RawGetString("max"); maxVal != lua.LNil {
		if max, err := strconv.Atoi(maxVal.String()); err == nil {
			if length > max {
				isValid = false
				errors = append(errors, fmt.Sprintf("maximum length is %d", max))
			}
		}
	}
	
	if exactVal := options.RawGetString("exact"); exactVal != lua.LNil {
		if exact, err := strconv.Atoi(exactVal.String()); err == nil {
			if length != exact {
				isValid = false
				errors = append(errors, fmt.Sprintf("exact length must be %d", exact))
			}
		}
	}
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	result.RawSetString("length", lua.LNumber(length))
	
	if !isValid {
		errorsTable := L.NewTable()
		for i, err := range errors {
			errorsTable.RawSetInt(i+1, lua.LString(err))
		}
		result.RawSetString("errors", errorsTable)
	}
	
	L.Push(result)
	return 1
}

// luaValidateRange validates numeric ranges
func (m *ValidateModule) luaValidateRange(L *lua.LState) int {
	valueStr := L.CheckString(1)
	options := L.CheckTable(2)
	
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		errorTable := L.NewTable()
		errorTable.RawSetString("error", lua.LString("invalid numeric value"))
		errorTable.RawSetString("message", lua.LString(err.Error()))
		L.Push(errorTable)
		return 1
	}
	
	var isValid = true
	var errors []string
	
	if minVal := options.RawGetString("min"); minVal != lua.LNil {
		if min, err := strconv.ParseFloat(minVal.String(), 64); err == nil {
			if value < min {
				isValid = false
				errors = append(errors, fmt.Sprintf("minimum value is %g", min))
			}
		}
	}
	
	if maxVal := options.RawGetString("max"); maxVal != lua.LNil {
		if max, err := strconv.ParseFloat(maxVal.String(), 64); err == nil {
			if value > max {
				isValid = false
				errors = append(errors, fmt.Sprintf("maximum value is %g", max))
			}
		}
	}
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	result.RawSetString("value", lua.LNumber(value))
	
	if !isValid {
		errorsTable := L.NewTable()
		for i, err := range errors {
			errorsTable.RawSetInt(i+1, lua.LString(err))
		}
		result.RawSetString("errors", errorsTable)
	}
	
	L.Push(result)
	return 1
}

// luaValidateRequired validates required fields
func (m *ValidateModule) luaValidateRequired(L *lua.LState) int {
	data := L.CheckTable(1)
	requiredFields := L.CheckTable(2)
	
	var missing []string
	
	requiredFields.ForEach(func(_, fieldName lua.LValue) {
		field := fieldName.String()
		value := data.RawGetString(field)
		if value == lua.LNil || (value.Type() == lua.LTString && strings.TrimSpace(value.String()) == "") {
			missing = append(missing, field)
		}
	})
	
	isValid := len(missing) == 0
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	
	if !isValid {
		missingTable := L.NewTable()
		for i, field := range missing {
			missingTable.RawSetInt(i+1, lua.LString(field))
		}
		result.RawSetString("missing", missingTable)
	}
	
	L.Push(result)
	return 1
}

// luaValidateSchema validates data against a schema
func (m *ValidateModule) luaValidateSchema(L *lua.LState) int {
	data := L.CheckTable(1)
	schema := L.CheckTable(2)
	
	var errors []string
	isValid := true
	
	schema.ForEach(func(fieldName, rules lua.LValue) {
		field := fieldName.String()
		value := data.RawGetString(field)
		
		if rulesTable, ok := rules.(*lua.LTable); ok {
			// Check if required
			if required := rulesTable.RawGetString("required"); required != lua.LNil && required.String() == "true" {
				if value == lua.LNil || (value.Type() == lua.LTString && strings.TrimSpace(value.String()) == "") {
					errors = append(errors, fmt.Sprintf("field '%s' is required", field))
					isValid = false
				}
			}
			
			// Skip other validations if field is empty and not required
			if value == lua.LNil || (value.Type() == lua.LTString && value.String() == "") {
				return
			}
			
			// Type validation
			if expectedType := rulesTable.RawGetString("type"); expectedType != lua.LNil {
				if !m.validateType(value, expectedType.String()) {
					errors = append(errors, fmt.Sprintf("field '%s' must be of type %s", field, expectedType.String()))
					isValid = false
				}
			}
		}
	})
	
	result := L.NewTable()
	result.RawSetString("valid", lua.LBool(isValid))
	
	if !isValid {
		errorsTable := L.NewTable()
		for i, err := range errors {
			errorsTable.RawSetInt(i+1, lua.LString(err))
		}
		result.RawSetString("errors", errorsTable)
	}
	
	L.Push(result)
	return 1
}

// luaSanitize sanitizes input strings
func (m *ValidateModule) luaSanitize(L *lua.LState) int {
	input := L.CheckString(1)
	method := L.OptString(2, "html")
	
	var sanitized string
	
	switch method {
	case "html":
		sanitized = m.sanitizeHTML(input)
	case "sql":
		sanitized = m.sanitizeSQL(input)
	case "trim":
		sanitized = strings.TrimSpace(input)
	case "lower":
		sanitized = strings.ToLower(strings.TrimSpace(input))
	case "upper":
		sanitized = strings.ToUpper(strings.TrimSpace(input))
	default:
		sanitized = input
	}
	
	result := L.NewTable()
	result.RawSetString("original", lua.LString(input))
	result.RawSetString("sanitized", lua.LString(sanitized))
	result.RawSetString("method", lua.LString(method))
	
	L.Push(result)
	return 1
}

// Helper functions
func (m *ValidateModule) validateType(value lua.LValue, expectedType string) bool {
	switch expectedType {
	case "string":
		return value.Type() == lua.LTString
	case "number":
		return value.Type() == lua.LTNumber
	case "boolean":
		return value.Type() == lua.LTBool
	case "table":
		return value.Type() == lua.LTTable
	default:
		return false
	}
}

func (m *ValidateModule) sanitizeHTML(input string) string {
	// Basic HTML escaping
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#x27;",
	)
	return replacer.Replace(input)
}

func (m *ValidateModule) sanitizeSQL(input string) string {
	// Basic SQL injection prevention
	input = strings.ReplaceAll(input, "'", "''")
	input = strings.ReplaceAll(input, "\"", "\"\"")
	return input
}