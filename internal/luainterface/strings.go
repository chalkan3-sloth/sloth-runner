package luainterface

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	lua "github.com/yuin/gopher-lua"
)

// StringModule provides string processing functionality for Lua scripts
type StringModule struct{}

// NewStringModule creates a new string processing module
func NewStringModule() *StringModule {
	return &StringModule{}
}

// RegisterStringModule registers the string module with the Lua state
func RegisterStringModule(L *lua.LState) {
	module := NewStringModule()
	
	// Create the strings table
	stringsTable := L.NewTable()
	
	// String manipulation functions
	L.SetField(stringsTable, "trim", L.NewFunction(module.luaStringTrim))
	L.SetField(stringsTable, "upper", L.NewFunction(module.luaStringUpper))
	L.SetField(stringsTable, "lower", L.NewFunction(module.luaStringLower))
	L.SetField(stringsTable, "title", L.NewFunction(module.luaStringTitle))
	L.SetField(stringsTable, "split", L.NewFunction(module.luaStringSplit))
	L.SetField(stringsTable, "join", L.NewFunction(module.luaStringJoin))
	L.SetField(stringsTable, "replace", L.NewFunction(module.luaStringReplace))
	L.SetField(stringsTable, "contains", L.NewFunction(module.luaStringContains))
	L.SetField(stringsTable, "starts_with", L.NewFunction(module.luaStringStartsWith))
	L.SetField(stringsTable, "ends_with", L.NewFunction(module.luaStringEndsWith))
	
	// Regular expressions
	L.SetField(stringsTable, "match", L.NewFunction(module.luaStringMatch))
	L.SetField(stringsTable, "match_all", L.NewFunction(module.luaStringMatchAll))
	L.SetField(stringsTable, "replace_regex", L.NewFunction(module.luaStringReplaceRegex))
	
	// Encoding/Decoding
	L.SetField(stringsTable, "base64_encode", L.NewFunction(module.luaStringBase64Encode))
	L.SetField(stringsTable, "base64_decode", L.NewFunction(module.luaStringBase64Decode))
	L.SetField(stringsTable, "url_encode", L.NewFunction(module.luaStringURLEncode))
	L.SetField(stringsTable, "url_decode", L.NewFunction(module.luaStringURLDecode))
	
	// Hashing
	L.SetField(stringsTable, "md5", L.NewFunction(module.luaStringMD5))
	L.SetField(stringsTable, "sha1", L.NewFunction(module.luaStringSHA1))
	L.SetField(stringsTable, "sha256", L.NewFunction(module.luaStringSHA256))
	
	// Validation
	L.SetField(stringsTable, "is_email", L.NewFunction(module.luaStringIsEmail))
	L.SetField(stringsTable, "is_url", L.NewFunction(module.luaStringIsURL))
	L.SetField(stringsTable, "is_numeric", L.NewFunction(module.luaStringIsNumeric))
	L.SetField(stringsTable, "is_alpha", L.NewFunction(module.luaStringIsAlpha))
	L.SetField(stringsTable, "is_alphanumeric", L.NewFunction(module.luaStringIsAlphanumeric))
	
	// Formatting
	L.SetField(stringsTable, "pad_left", L.NewFunction(module.luaStringPadLeft))
	L.SetField(stringsTable, "pad_right", L.NewFunction(module.luaStringPadRight))
	L.SetField(stringsTable, "truncate", L.NewFunction(module.luaStringTruncate))
	
	// Set the strings module in global scope
	L.SetGlobal("strings", stringsTable)
}

// String manipulation functions

func (s *StringModule) luaStringTrim(L *lua.LState) int {
	str := L.CheckString(1)
	result := strings.TrimSpace(str)
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringUpper(L *lua.LState) int {
	str := L.CheckString(1)
	result := strings.ToUpper(str)
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringLower(L *lua.LState) int {
	str := L.CheckString(1)
	result := strings.ToLower(str)
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringTitle(L *lua.LState) int {
	str := L.CheckString(1)
	result := strings.Title(str)
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringSplit(L *lua.LState) int {
	str := L.CheckString(1)
	separator := L.CheckString(2)
	
	parts := strings.Split(str, separator)
	table := L.NewTable()
	
	for i, part := range parts {
		L.RawSetInt(table, i+1, lua.LString(part))
	}
	
	L.Push(table)
	return 1
}

func (s *StringModule) luaStringJoin(L *lua.LState) int {
	table := L.CheckTable(1)
	separator := L.CheckString(2)
	
	var parts []string
	table.ForEach(func(key, value lua.LValue) {
		parts = append(parts, value.String())
	})
	
	result := strings.Join(parts, separator)
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringReplace(L *lua.LState) int {
	str := L.CheckString(1)
	old := L.CheckString(2)
	new := L.CheckString(3)
	count := L.OptInt(4, -1)
	
	result := strings.Replace(str, old, new, count)
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringContains(L *lua.LState) int {
	str := L.CheckString(1)
	substr := L.CheckString(2)
	
	result := strings.Contains(str, substr)
	L.Push(lua.LBool(result))
	return 1
}

func (s *StringModule) luaStringStartsWith(L *lua.LState) int {
	str := L.CheckString(1)
	prefix := L.CheckString(2)
	
	result := strings.HasPrefix(str, prefix)
	L.Push(lua.LBool(result))
	return 1
}

func (s *StringModule) luaStringEndsWith(L *lua.LState) int {
	str := L.CheckString(1)
	suffix := L.CheckString(2)
	
	result := strings.HasSuffix(str, suffix)
	L.Push(lua.LBool(result))
	return 1
}

// Regular expressions

func (s *StringModule) luaStringMatch(L *lua.LState) int {
	str := L.CheckString(1)
	pattern := L.CheckString(2)
	
	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid regex pattern: " + err.Error()))
		return 2
	}
	
	matches := re.FindStringSubmatch(str)
	if matches == nil {
		L.Push(lua.LNil)
		return 1
	}
	
	table := L.NewTable()
	for i, match := range matches {
		L.RawSetInt(table, i+1, lua.LString(match))
	}
	
	L.Push(table)
	return 1
}

func (s *StringModule) luaStringMatchAll(L *lua.LState) int {
	str := L.CheckString(1)
	pattern := L.CheckString(2)
	
	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid regex pattern: " + err.Error()))
		return 2
	}
	
	allMatches := re.FindAllStringSubmatch(str, -1)
	if allMatches == nil {
		L.Push(L.NewTable())
		return 1
	}
	
	table := L.NewTable()
	for i, matches := range allMatches {
		matchTable := L.NewTable()
		for j, match := range matches {
			L.RawSetInt(matchTable, j+1, lua.LString(match))
		}
		L.RawSetInt(table, i+1, matchTable)
	}
	
	L.Push(table)
	return 1
}

func (s *StringModule) luaStringReplaceRegex(L *lua.LState) int {
	str := L.CheckString(1)
	pattern := L.CheckString(2)
	replacement := L.CheckString(3)
	
	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid regex pattern: " + err.Error()))
		return 2
	}
	
	result := re.ReplaceAllString(str, replacement)
	L.Push(lua.LString(result))
	return 1
}

// Encoding/Decoding

func (s *StringModule) luaStringBase64Encode(L *lua.LState) int {
	str := L.CheckString(1)
	result := base64.StdEncoding.EncodeToString([]byte(str))
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringBase64Decode(L *lua.LState) int {
	str := L.CheckString(1)
	result, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid base64 string: " + err.Error()))
		return 2
	}
	L.Push(lua.LString(string(result)))
	return 1
}

func (s *StringModule) luaStringURLEncode(L *lua.LState) int {
	str := L.CheckString(1)
	// Using Go's url.QueryEscape equivalent
	result := strings.ReplaceAll(str, " ", "%20")
	result = strings.ReplaceAll(result, "+", "%2B")
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringURLDecode(L *lua.LState) int {
	str := L.CheckString(1)
	// Basic URL decoding
	result := strings.ReplaceAll(str, "%20", " ")
	result = strings.ReplaceAll(result, "%2B", "+")
	L.Push(lua.LString(result))
	return 1
}

// Hashing

func (s *StringModule) luaStringMD5(L *lua.LState) int {
	str := L.CheckString(1)
	hash := md5.Sum([]byte(str))
	result := hex.EncodeToString(hash[:])
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringSHA1(L *lua.LState) int {
	str := L.CheckString(1)
	hash := sha1.Sum([]byte(str))
	result := hex.EncodeToString(hash[:])
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringSHA256(L *lua.LState) int {
	str := L.CheckString(1)
	hash := sha256.Sum256([]byte(str))
	result := hex.EncodeToString(hash[:])
	L.Push(lua.LString(result))
	return 1
}

// Validation

func (s *StringModule) luaStringIsEmail(L *lua.LState) int {
	str := L.CheckString(1)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	result := emailRegex.MatchString(str)
	L.Push(lua.LBool(result))
	return 1
}

func (s *StringModule) luaStringIsURL(L *lua.LState) int {
	str := L.CheckString(1)
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	result := urlRegex.MatchString(str)
	L.Push(lua.LBool(result))
	return 1
}

func (s *StringModule) luaStringIsNumeric(L *lua.LState) int {
	str := L.CheckString(1)
	_, err := strconv.ParseFloat(str, 64)
	L.Push(lua.LBool(err == nil))
	return 1
}

func (s *StringModule) luaStringIsAlpha(L *lua.LState) int {
	str := L.CheckString(1)
	for _, r := range str {
		if !unicode.IsLetter(r) {
			L.Push(lua.LBool(false))
			return 1
		}
	}
	L.Push(lua.LBool(len(str) > 0))
	return 1
}

func (s *StringModule) luaStringIsAlphanumeric(L *lua.LState) int {
	str := L.CheckString(1)
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			L.Push(lua.LBool(false))
			return 1
		}
	}
	L.Push(lua.LBool(len(str) > 0))
	return 1
}

// Formatting

func (s *StringModule) luaStringPadLeft(L *lua.LState) int {
	str := L.CheckString(1)
	width := L.CheckInt(2)
	pad := L.OptString(3, " ")
	
	if len(str) >= width {
		L.Push(lua.LString(str))
		return 1
	}
	
	padding := strings.Repeat(pad, (width-len(str))/len(pad)+1)
	result := padding[:width-len(str)] + str
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringPadRight(L *lua.LState) int {
	str := L.CheckString(1)
	width := L.CheckInt(2)
	pad := L.OptString(3, " ")
	
	if len(str) >= width {
		L.Push(lua.LString(str))
		return 1
	}
	
	padding := strings.Repeat(pad, (width-len(str))/len(pad)+1)
	result := str + padding[:width-len(str)]
	L.Push(lua.LString(result))
	return 1
}

func (s *StringModule) luaStringTruncate(L *lua.LState) int {
	str := L.CheckString(1)
	maxLength := L.CheckInt(2)
	suffix := L.OptString(3, "...")
	
	if len(str) <= maxLength {
		L.Push(lua.LString(str))
		return 1
	}
	
	if maxLength <= len(suffix) {
		L.Push(lua.LString(suffix[:maxLength]))
		return 1
	}
	
	result := str[:maxLength-len(suffix)] + suffix
	L.Push(lua.LString(result))
	return 1
}