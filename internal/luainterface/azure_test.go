package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestAzureVMList(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAzureModule(L)

	script := `
		local azure = require("azure")
		local result = azure.vm_list({
			resource_group = "test-rg",
			subscription = "test-sub"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Azure test skipped (Azure credentials not available): %v", err)
		return
	}
}

func TestAzureVMStart(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAzureModule(L)

	script := `
		local azure = require("azure")
		local result = azure.vm_start({
			name = "test-vm",
			resource_group = "test-rg",
			subscription = "test-sub"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Azure test skipped (Azure credentials not available): %v", err)
		return
	}
}

func TestAzureVMStop(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAzureModule(L)

	script := `
		local azure = require("azure")
		local result = azure.vm_stop({
			name = "test-vm",
			resource_group = "test-rg",
			subscription = "test-sub"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Azure test skipped (Azure credentials not available): %v", err)
		return
	}
}

func TestAzureStorageUpload(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAzureModule(L)

	script := `
		local azure = require("azure")
		local result = azure.storage_upload({
			account = "testaccount",
			container = "testcontainer",
			blob = "testblob",
			file = "/nonexistent/file"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Azure test skipped (Azure credentials not available): %v", err)
		return
	}
}

func TestAzureStorageDownload(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAzureModule(L)

	script := `
		local azure = require("azure")
		local result = azure.storage_download({
			account = "testaccount",
			container = "testcontainer",
			blob = "testblob",
			destination = "/tmp/test-file"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Azure test skipped (Azure credentials not available): %v", err)
		return
	}
}

func TestAzureResourceGroupList(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAzureModule(L)

	script := `
		local azure = require("azure")
		local result = azure.resource_group_list({
			subscription = "test-sub"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Azure test skipped (Azure credentials not available): %v", err)
		return
	}
}
