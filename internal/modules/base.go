package modules

import "github.com/yuin/gopher-lua"

// BaseModule provides common functionality for all modules
type BaseModule struct {
	info ModuleInfo
}

// NewBaseModule creates a new base module
func NewBaseModule(info ModuleInfo) *BaseModule {
	return &BaseModule{
		info: info,
	}
}

// Info returns the module information
func (m *BaseModule) Info() ModuleInfo {
	return m.info
}

// Initialize performs basic initialization - can be overridden
func (m *BaseModule) Initialize() error {
	return nil
}

// Cleanup performs basic cleanup - can be overridden
func (m *BaseModule) Cleanup() error {
	return nil
}

// ValidationResult represents the result of input validation
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// ValidateRequired checks if required parameters are present
func ValidateRequired(L *lua.LState, table *lua.LTable, required []string) ValidationResult {
	var result ValidationResult
	result.IsValid = true
	
	for _, field := range required {
		value := table.RawGetString(field)
		if value == lua.LNil || (value.Type() == lua.LTString && value.String() == "") {
			result.IsValid = false
			result.Errors = append(result.Errors, "missing required field: "+field)
		}
	}
	
	return result
}

// CreateErrorResponse creates a standardized error response for Lua
func CreateErrorResponse(L *lua.LState, message string, details ...string) int {
	result := L.NewTable()
	result.RawSetString("success", lua.LBool(false))
	result.RawSetString("error", lua.LString(message))
	
	if len(details) > 0 {
		detailsTable := L.NewTable()
		for i, detail := range details {
			detailsTable.RawSetInt(i+1, lua.LString(detail))
		}
		result.RawSetString("details", detailsTable)
	}
	
	L.Push(result)
	return 1
}

// CreateSuccessResponse creates a standardized success response for Lua
func CreateSuccessResponse(L *lua.LState, data lua.LValue) int {
	result := L.NewTable()
	result.RawSetString("success", lua.LBool(true))
	
	if data != nil {
		result.RawSetString("data", data)
	}
	
	L.Push(result)
	return 1
}

// WrapLuaFunction wraps a function with error handling and validation
func WrapLuaFunction(fn lua.LGFunction, required []string) lua.LGFunction {
	return func(L *lua.LState) int {
		// Validate required parameters if a table is passed
		if len(required) > 0 && L.GetTop() > 0 {
			if table, ok := L.Get(1).(*lua.LTable); ok {
				validation := ValidateRequired(L, table, required)
				if !validation.IsValid {
					return CreateErrorResponse(L, "validation failed", validation.Errors...)
				}
			}
		}
		
		// Call the original function with error handling
		defer func() {
			if r := recover(); r != nil {
				CreateErrorResponse(L, "internal error", r.(string))
			}
		}()
		
		return fn(L)
	}
}