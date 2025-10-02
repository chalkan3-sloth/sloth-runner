package luainterface

import (
	"net/http"
	"net/http/httptest"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestHTTPModuleBasic(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/test":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "hello world"}`))
		case "/json":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok", "data": {"value": 123}}`))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "internal error"}`))
		case "/headers":
			w.Header().Set("X-Custom-Header", "test-value")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"headers": "ok"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "http.get()",
			script: `
				local resp = http.get("` + server.URL + `/test")
				assert(type(resp) == "table", "get() should return a table")
				assert(resp.status_code == 200, "status_code should be 200")
				assert(type(resp.body) == "string", "body should be a string")
			`,
		},
		{
			name: "http.get() with JSON",
			script: `
				local resp = http.get("` + server.URL + `/json")
				assert(resp.status_code == 200, "status_code should be 200")
				assert(type(resp.body) == "string", "body should be a string")
			`,
		},
		{
			name: "http.post()",
			script: `
				local data = '{"key":"value"}'
				local resp = http.post("` + server.URL + `/test", data)
				assert(type(resp) == "table", "post() should return a table")
			`,
		},
		{
			name: "http.request() with custom options",
			script: `
				local options = {
					url = "` + server.URL + `/test",
					method = "GET",
					headers = {["X-Test"] = "test"}
				}
				local resp = http.request(options)
				assert(type(resp) == "table", "request() should return a table")
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

func TestHTTPClient(t *testing.T) {
	t.Skip("HTTP client methods not yet implemented")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "ok"}`))
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	script := `
		local client = http.new_client({
			timeout = 30,
			follow_redirects = true
		})
		assert(client ~= nil, "new_client() should return a client")
		
		local resp = client:get("` + server.URL + `")
		assert(type(resp) == "table", "client:get() should return a table")
		assert(resp.status == 200, "status should be 200")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestHTTPStatusChecks(t *testing.T) {
	t.Skip("HTTP status check helpers not yet implemented")
	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	script := `
		assert(http.is_success(200) == true, "200 is success")
		assert(http.is_success(299) == true, "299 is success")
		assert(http.is_success(300) == false, "300 is not success")
		assert(http.is_success(404) == false, "404 is not success")
		
		assert(http.is_redirect(301) == true, "301 is redirect")
		assert(http.is_redirect(302) == true, "302 is redirect")
		assert(http.is_redirect(200) == false, "200 is not redirect")
		
		assert(http.is_client_error(404) == true, "404 is client error")
		assert(http.is_client_error(400) == true, "400 is client error")
		assert(http.is_client_error(500) == false, "500 is not client error")
		
		assert(http.is_server_error(500) == true, "500 is server error")
		assert(http.is_server_error(503) == true, "503 is server error")
		assert(http.is_server_error(404) == false, "404 is not server error")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestHTTPURLEncoding(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	script := `
		local encoded = http.url_encode("hello world")
		assert(type(encoded) == "string", "url_encode() should return a string")
		assert(encoded == "hello+world" or encoded == "hello%20world", "should encode spaces")
		
		local decoded = http.url_decode(encoded)
		assert(decoded == "hello world", "url_decode() should decode correctly")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestHTTPBuildURL(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	script := `
		local url = http.build_url("https://example.com", "/path", {
			key1 = "value1",
			key2 = "value2"
		})
		assert(type(url) == "string", "build_url() should return a string")
		assert(url:find("example.com") ~= nil, "URL should contain domain")
		assert(url:find("key1") ~= nil, "URL should contain query parameter")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestHTTPParseJSON(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	script := `
		local json_str = '{"name": "test", "value": 123, "nested": {"key": "val"}}'
		local data = http.parse_json(json_str)
		assert(type(data) == "table", "parse_json() should return a table")
		assert(data.name == "test", "should parse string field")
		assert(data.value == 123, "should parse number field")
		assert(type(data.nested) == "table", "should parse nested object")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestHTTPToJSON(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	script := `
		local data = {
			name = "test",
			value = 123,
			nested = {key = "val"}
		}
		local json_str = http.to_json(data)
		assert(type(json_str) == "string", "to_json() should return a string")
		assert(json_str:find("test") ~= nil, "JSON should contain data")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestHTTPDownload(t *testing.T) {
t.Skip("Module not yet registered globally - needs refactoring")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file content"))
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	RegisterHTTPModule(L)

	script := `
		local success = http.download("` + server.URL + `", "/tmp/test-download.txt")
		assert(type(success) == "boolean", "download() should return a boolean")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
