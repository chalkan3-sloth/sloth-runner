package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chalkan3/sloth-runner/internal/modules"
	"github.com/yuin/gopher-lua"
)

// HTTPModule provides HTTP client functionality
type HTTPModule struct {
	*modules.BaseModule
	client *http.Client
}

// NewHTTPModule creates a new HTTP module
func NewHTTPModule() *HTTPModule {
	info := modules.ModuleInfo{
		Name:        "http",
		Version:     "1.0.0",
		Description: "HTTP client with advanced features including retries, timeouts, and response validation",
		Author:      "Sloth Runner Team",
		Category:    "core",
		Dependencies: []string{},
	}
	
	return &HTTPModule{
		BaseModule: modules.NewBaseModule(info),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Loader returns the Lua loader function
func (m *HTTPModule) Loader(L *lua.LState) int {
	httpTable := L.NewTable()
	
	// HTTP methods
	L.SetFuncs(httpTable, map[string]lua.LGFunction{
		"get":     modules.WrapLuaFunction(m.luaGet, []string{"url"}),
		"post":    modules.WrapLuaFunction(m.luaPost, []string{"url"}),
		"put":     modules.WrapLuaFunction(m.luaPut, []string{"url"}),
		"delete":  modules.WrapLuaFunction(m.luaDelete, []string{"url"}),
		"patch":   modules.WrapLuaFunction(m.luaPatch, []string{"url"}),
		"request": modules.WrapLuaFunction(m.luaRequest, []string{"method", "url"}),
		"client":  m.luaCreateClient,
	})
	
	L.Push(httpTable)
	return 1
}

// HTTPRequest represents an HTTP request configuration
type HTTPRequest struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        string
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
	ValidateTLS bool
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
	Elapsed    time.Duration
}

// parseRequestOptions parses Lua table into HTTPRequest
func (m *HTTPModule) parseRequestOptions(L *lua.LState, table *lua.LTable) HTTPRequest {
	req := HTTPRequest{
		Headers:     make(map[string]string),
		Timeout:     30 * time.Second,
		MaxRetries:  0,
		RetryDelay:  1 * time.Second,
		ValidateTLS: true,
	}
	
	if url := table.RawGetString("url"); url != lua.LNil {
		req.URL = url.String()
	}
	
	if method := table.RawGetString("method"); method != lua.LNil {
		req.Method = strings.ToUpper(method.String())
	}
	
	if body := table.RawGetString("body"); body != lua.LNil {
		req.Body = body.String()
	}
	
	if timeout := table.RawGetString("timeout"); timeout != lua.LNil {
		if t, err := strconv.Atoi(timeout.String()); err == nil {
			req.Timeout = time.Duration(t) * time.Second
		}
	}
	
	if retries := table.RawGetString("max_retries"); retries != lua.LNil {
		if r, err := strconv.Atoi(retries.String()); err == nil {
			req.MaxRetries = r
		}
	}
	
	if delay := table.RawGetString("retry_delay"); delay != lua.LNil {
		if d, err := strconv.Atoi(delay.String()); err == nil {
			req.RetryDelay = time.Duration(d) * time.Second
		}
	}
	
	// Parse headers
	if headers := table.RawGetString("headers"); headers != lua.LNil {
		if headersTable, ok := headers.(*lua.LTable); ok {
			headersTable.ForEach(func(key, value lua.LValue) {
				req.Headers[key.String()] = value.String()
			})
		}
	}
	
	// Parse JSON body
	if jsonData := table.RawGetString("json"); jsonData != lua.LNil {
		if jsonTable, ok := jsonData.(*lua.LTable); ok {
			data := m.luaTableToGoValue(jsonTable)
			if jsonBytes, err := json.Marshal(data); err == nil {
				req.Body = string(jsonBytes)
				req.Headers["Content-Type"] = "application/json"
			}
		}
	}
	
	return req
}

// executeRequest performs the HTTP request with retries
func (m *HTTPModule) executeRequest(req HTTPRequest) (*HTTPResponse, error) {
	start := time.Now()
	
	var lastErr error
	for attempt := 0; attempt <= req.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(req.RetryDelay)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
		defer cancel()
		
		httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, strings.NewReader(req.Body))
		if err != nil {
			lastErr = err
			continue
		}
		
		// Set headers
		for key, value := range req.Headers {
			httpReq.Header.Set(key, value)
		}
		
		resp, err := m.client.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}
		
		// Build response headers
		headers := make(map[string]string)
		for key, values := range resp.Header {
			headers[key] = strings.Join(values, ", ")
		}
		
		response := &HTTPResponse{
			StatusCode: resp.StatusCode,
			Headers:    headers,
			Body:       string(body),
			Elapsed:    time.Since(start),
		}
		
		return response, nil
	}
	
	return nil, fmt.Errorf("request failed after %d attempts: %w", req.MaxRetries+1, lastErr)
}

// luaGet implements HTTP GET
func (m *HTTPModule) luaGet(L *lua.LState) int {
	table := L.CheckTable(1)
	req := m.parseRequestOptions(L, table)
	req.Method = "GET"
	
	return m.performRequest(L, req)
}

// luaPost implements HTTP POST
func (m *HTTPModule) luaPost(L *lua.LState) int {
	table := L.CheckTable(1)
	req := m.parseRequestOptions(L, table)
	req.Method = "POST"
	
	return m.performRequest(L, req)
}

// luaPut implements HTTP PUT
func (m *HTTPModule) luaPut(L *lua.LState) int {
	table := L.CheckTable(1)
	req := m.parseRequestOptions(L, table)
	req.Method = "PUT"
	
	return m.performRequest(L, req)
}

// luaDelete implements HTTP DELETE
func (m *HTTPModule) luaDelete(L *lua.LState) int {
	table := L.CheckTable(1)
	req := m.parseRequestOptions(L, table)
	req.Method = "DELETE"
	
	return m.performRequest(L, req)
}

// luaPatch implements HTTP PATCH
func (m *HTTPModule) luaPatch(L *lua.LState) int {
	table := L.CheckTable(1)
	req := m.parseRequestOptions(L, table)
	req.Method = "PATCH"
	
	return m.performRequest(L, req)
}

// luaRequest implements generic HTTP request
func (m *HTTPModule) luaRequest(L *lua.LState) int {
	table := L.CheckTable(1)
	req := m.parseRequestOptions(L, table)
	
	return m.performRequest(L, req)
}

// performRequest executes the HTTP request and returns Lua response
func (m *HTTPModule) performRequest(L *lua.LState, req HTTPRequest) int {
	resp, err := m.executeRequest(req)
	if err != nil {
		return modules.CreateErrorResponse(L, "HTTP request failed", err.Error())
	}
	
	// Create response table
	responseTable := L.NewTable()
	responseTable.RawSetString("status_code", lua.LNumber(resp.StatusCode))
	responseTable.RawSetString("body", lua.LString(resp.Body))
	responseTable.RawSetString("elapsed_ms", lua.LNumber(resp.Elapsed.Milliseconds()))
	
	// Add headers
	headersTable := L.NewTable()
	for key, value := range resp.Headers {
		headersTable.RawSetString(key, lua.LString(value))
	}
	responseTable.RawSetString("headers", headersTable)
	
	// Try to parse JSON response
	var jsonData interface{}
	if err := json.Unmarshal([]byte(resp.Body), &jsonData); err == nil {
		responseTable.RawSetString("json", m.goValueToLua(L, jsonData))
	}
	
	return modules.CreateSuccessResponse(L, responseTable)
}

// luaCreateClient creates a custom HTTP client
func (m *HTTPModule) luaCreateClient(L *lua.LState) int {
	table := L.OptTable(1, L.NewTable())
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	if timeout := table.RawGetString("timeout"); timeout != lua.LNil {
		if t, err := strconv.Atoi(timeout.String()); err == nil {
			client.Timeout = time.Duration(t) * time.Second
		}
	}
	
	// Create client userdata
	ud := L.NewUserData()
	ud.Value = client
	L.Push(ud)
	return 1
}

// Helper functions
func (m *HTTPModule) luaTableToGoValue(table *lua.LTable) interface{} {
	result := make(map[string]interface{})
	table.ForEach(func(key, value lua.LValue) {
		switch v := value.(type) {
		case lua.LBool:
			result[key.String()] = bool(v)
		case lua.LNumber:
			result[key.String()] = float64(v)
		case lua.LString:
			result[key.String()] = string(v)
		case *lua.LTable:
			result[key.String()] = m.luaTableToGoValue(v)
		}
	})
	return result
}

func (m *HTTPModule) goValueToLua(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case map[string]interface{}:
		table := L.NewTable()
		for key, val := range v {
			table.RawSetString(key, m.goValueToLua(L, val))
		}
		return table
	case []interface{}:
		table := L.NewTable()
		for i, val := range v {
			table.RawSetInt(i+1, m.goValueToLua(L, val))
		}
		return table
	default:
		return lua.LNil
	}
}