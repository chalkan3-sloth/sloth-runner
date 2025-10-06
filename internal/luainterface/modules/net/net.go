package net

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// HttpGet performs HTTP GET request
func HttpGet(L *lua.LState) int {
	url := L.CheckString(1)

	resp, err := http.Get(url)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNumber(0))
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 4
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNumber(resp.StatusCode))
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 4
	}

	headersTable := L.NewTable()
	for name, values := range resp.Header {
		headerValues := L.NewTable()
		for i, val := range values {
			headerValues.RawSetInt(i+1, lua.LString(val))
		}
		headersTable.RawSetString(name, headerValues)
	}

	L.Push(lua.LString(string(bodyBytes)))
	L.Push(lua.LNumber(resp.StatusCode))
	L.Push(headersTable)
	L.Push(lua.LNil) // No error
	return 4
}

// HttpPost performs HTTP POST request
func HttpPost(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.CheckString(2)
	headersTable := L.OptTable(3, L.NewTable()) // Optional headers table

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNumber(0))
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 4
	}

	headersTable.ForEach(func(key, value lua.LValue) {
		req.Header.Set(key.String(), value.String())
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNumber(0))
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 4
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNumber(resp.StatusCode))
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 4
	}

	respHeadersTable := L.NewTable()
	for name, values := range resp.Header {
		headerValues := L.NewTable()
		for i, val := range values {
			headerValues.RawSetInt(i+1, lua.LString(val))
		}
		respHeadersTable.RawSetString(name, headerValues)
	}

	L.Push(lua.LString(string(respBodyBytes)))
	L.Push(lua.LNumber(resp.StatusCode))
	L.Push(respHeadersTable)
	L.Push(lua.LNil) // No error
	return 4
}

// Download downloads a file from URL to local path
func Download(L *lua.LState) int {
	url := L.CheckString(1)
	destinationPath := L.CheckString(2)

	resp, err := http.Get(url)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		L.Push(lua.LString(fmt.Sprintf("failed to download file: status code %d", resp.StatusCode)))
		return 1
	}

	out, err := os.Create(destinationPath)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// Loader returns the net module loader
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"http_get":  HttpGet,
		"http_post": HttpPost,
		"download":  Download,
	})
	L.Push(mod)
	return 1
}

// Open registers the net module and loads it globally
func Open(L *lua.LState) {
	L.PreloadModule("net", Loader)
	if err := L.DoString(`net = require("net")`); err != nil {
		panic(err)
	}
}
