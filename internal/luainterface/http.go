package luainterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// HTTPModule provides HTTP client functionality for Lua scripts
type HTTPModule struct {
	client *http.Client
}

// NewHTTPModule creates a new HTTP module
func NewHTTPModule() *HTTPModule {
	return &HTTPModule{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RegisterHTTPModule registers the HTTP module with the Lua state
func RegisterHTTPModule(L *lua.LState) {
	module := NewHTTPModule()
	
	// Create the http table
	httpTable := L.NewTable()
	
	// Register functions
	L.SetField(httpTable, "get", L.NewFunction(module.luaHTTPGet))
	L.SetField(httpTable, "post", L.NewFunction(module.luaHTTPPost))
	L.SetField(httpTable, "put", L.NewFunction(module.luaHTTPPut))
	L.SetField(httpTable, "delete", L.NewFunction(module.luaHTTPDelete))
	L.SetField(httpTable, "patch", L.NewFunction(module.luaHTTPPatch))
	L.SetField(httpTable, "request", L.NewFunction(module.luaHTTPRequest))
	
	// Set the http module in global scope
	L.SetGlobal("http", httpTable)
}

// luaHTTPGet performs a GET request
func (h *HTTPModule) luaHTTPGet(L *lua.LState) int {
	url := L.CheckString(1)
	headers := L.OptTable(2, nil)
	
	return h.performRequest(L, "GET", url, nil, headers)
}

// luaHTTPPost performs a POST request
func (h *HTTPModule) luaHTTPPost(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.OptString(2, "")
	headers := L.OptTable(3, nil)
	
	return h.performRequest(L, "POST", url, []byte(body), headers)
}

// luaHTTPPut performs a PUT request
func (h *HTTPModule) luaHTTPPut(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.OptString(2, "")
	headers := L.OptTable(3, nil)
	
	return h.performRequest(L, "PUT", url, []byte(body), headers)
}

// luaHTTPDelete performs a DELETE request
func (h *HTTPModule) luaHTTPDelete(L *lua.LState) int {
	url := L.CheckString(1)
	headers := L.OptTable(2, nil)
	
	return h.performRequest(L, "DELETE", url, nil, headers)
}

// luaHTTPPatch performs a PATCH request
func (h *HTTPModule) luaHTTPPatch(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.OptString(2, "")
	headers := L.OptTable(3, nil)
	
	return h.performRequest(L, "PATCH", url, []byte(body), headers)
}

// luaHTTPRequest performs a custom HTTP request
func (h *HTTPModule) luaHTTPRequest(L *lua.LState) int {
	options := L.CheckTable(1)
	
	// Extract options
	method := "GET"
	if methodVal := L.GetField(options, "method"); methodVal != lua.LNil {
		method = methodVal.String()
	}
	
	url := ""
	if urlVal := L.GetField(options, "url"); urlVal != lua.LNil {
		url = urlVal.String()
	}
	
	var body []byte
	if bodyVal := L.GetField(options, "body"); bodyVal != lua.LNil {
		if bodyVal.Type() == lua.LTTable {
			// Convert table to JSON
			bodyData := tableToMap(bodyVal.(*lua.LTable))
			jsonData, err := json.Marshal(bodyData)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString("Failed to marshal JSON body: " + err.Error()))
				return 2
			}
			body = jsonData
		} else {
			body = []byte(bodyVal.String())
		}
	}
	
	var headers *lua.LTable
	if headersVal := L.GetField(options, "headers"); headersVal != lua.LNil {
		headers = headersVal.(*lua.LTable)
	}
	
	var timeout time.Duration = 30 * time.Second
	if timeoutVal := L.GetField(options, "timeout"); timeoutVal != lua.LNil {
		if t, err := time.ParseDuration(timeoutVal.String()); err == nil {
			timeout = t
		}
	}
	
	// Create client with custom timeout
	client := &http.Client{Timeout: timeout}
	h.client = client
	
	return h.performRequest(L, method, url, body, headers)
}

// performRequest executes the HTTP request and returns the response
func (h *HTTPModule) performRequest(L *lua.LState, method, url string, body []byte, headers *lua.LTable) int {
	// Create request
	var req *http.Request
	var err error
	
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("Failed to create request: " + err.Error()))
		return 2
	}
	
	// Set headers
	if headers != nil {
		headers.ForEach(func(key, value lua.LValue) {
			req.Header.Set(key.String(), value.String())
		})
	}
	
	// Set default content type for POST/PUT/PATCH with body
	if body != nil && req.Header.Get("Content-Type") == "" {
		if json.Valid(body) {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "text/plain")
		}
	}
	
	// Perform request
	resp, err := h.client.Do(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("Request failed: " + err.Error()))
		return 2
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("Failed to read response: " + err.Error()))
		return 2
	}
	
	// Create response table
	responseTable := L.NewTable()
	L.SetField(responseTable, "status_code", lua.LNumber(resp.StatusCode))
	L.SetField(responseTable, "status", lua.LString(resp.Status))
	L.SetField(responseTable, "body", lua.LString(string(respBody)))
	
	// Add headers
	headersTable := L.NewTable()
	for key, values := range resp.Header {
		if len(values) == 1 {
			L.SetField(headersTable, key, lua.LString(values[0]))
		} else {
			headerArray := L.NewTable()
			for i, value := range values {
				L.RawSetInt(headerArray, i+1, lua.LString(value))
			}
			L.SetField(headersTable, key, headerArray)
		}
	}
	L.SetField(responseTable, "headers", headersTable)
	
	// Add success flag
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	L.SetField(responseTable, "success", lua.LBool(success))
	
	// Try to parse JSON if response is JSON
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var jsonData interface{}
		if err := json.Unmarshal(respBody, &jsonData); err == nil {
			L.SetField(responseTable, "json", interfaceToLua(L, jsonData))
		}
	}
	
	L.Push(responseTable)
	return 1
}

// Helper function to convert Lua table to Go map
func tableToMap(table *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	table.ForEach(func(key, value lua.LValue) {
		switch v := value.(type) {
		case *lua.LTable:
			result[key.String()] = tableToMap(v)
		case lua.LString:
			result[key.String()] = string(v)
		case lua.LNumber:
			result[key.String()] = float64(v)
		case lua.LBool:
			result[key.String()] = bool(v)
		default:
			result[key.String()] = v.String()
		}
	})
	return result
}

// Helper function to convert Go interface{} to Lua value
func interfaceToLua(L *lua.LState, data interface{}) lua.LValue {
	switch v := data.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		table := L.NewTable()
		for i, item := range v {
			L.RawSetInt(table, i+1, interfaceToLua(L, item))
		}
		return table
	case map[string]interface{}:
		table := L.NewTable()
		for key, value := range v {
			L.SetField(table, key, interfaceToLua(L, value))
		}
		return table
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}