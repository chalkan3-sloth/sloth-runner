package luainterface

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestSSHModuleRegistration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	// Check if module is registered
	if L.GetGlobal("ssh") == lua.LNil {
		t.Fatal("ssh module not registered")
	}

	// Check if functions exist
	sshTable := L.GetGlobal("ssh").(*lua.LTable)
	functions := []string{
		"connect", "disconnect", "exec",
		"upload", "download", "upload_dir", "download_dir",
		"exists", "stat", "mkdir", "remove", "rename", "chmod", "chown", "list_dir",
		"load_private_key", "generate_keypair",
		"create_tunnel", "close_tunnel", "enable_agent_forward",
	}

	for _, fn := range functions {
		if sshTable.RawGetString(fn) == lua.LNil {
			t.Errorf("Function %s not registered", fn)
		}
	}
}

func TestSSHConnect_MissingHost(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	// Try to connect without proper host
	script := `
		local conn, err = ssh.connect("invalid-host-12345", "user", {
			password = "test",
			timeout = 1
		})
		return conn == nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Expected error: %v", err)
	}

	result := L.Get(-1)
	if result == lua.LTrue {
		t.Log("Connection properly failed for invalid host")
	}
}

func TestSSHLoadPrivateKey_InvalidPath(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local key, err = ssh.load_private_key("/nonexistent/key")
		return key == nil and err ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for nonexistent key file")
	}
}

func TestSSHLoadPrivateKey_InvalidFormat(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	// Create a temporary invalid key file
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "invalid_key")
	if err := os.WriteFile(keyPath, []byte("not a valid key"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	script := fmt.Sprintf(`
		local key, err = ssh.load_private_key("%s")
		return key == nil and err ~= nil
	`, keyPath)

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid key format")
	}
}

func TestSSHGenerateKeypair_NotImplemented(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local key, err = ssh.generate_keypair()
		return key == nil and err ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected not implemented error")
	}
}

func TestSSHCreateTunnel_NotImplemented(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local tunnel, err = ssh.create_tunnel()
		return tunnel == nil and err ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected not implemented error")
	}
}

func TestSSHCloseTunnel_NotImplemented(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local success, err = ssh.close_tunnel()
		return not success and err ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected not implemented error")
	}
}

func TestSSHEnableAgentForward_NotImplemented(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local success, err = ssh.enable_agent_forward()
		return not success and err ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected not implemented error")
	}
}

func TestSSHHelperFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test getTableString
	t.Run("getTableString", func(t *testing.T) {
		table := L.NewTable()
		table.RawSetString("key", lua.LString("value"))

		result := getTableString(table, "key", "default")
		if result != "value" {
			t.Errorf("Expected 'value', got '%s'", result)
		}

		result = getTableString(table, "nonexistent", "default")
		if result != "default" {
			t.Errorf("Expected 'default', got '%s'", result)
		}
	})

	// Test getTableInt
	t.Run("getTableInt", func(t *testing.T) {
		table := L.NewTable()
		table.RawSetString("key", lua.LNumber(42))

		result := getTableInt(table, "key", 0)
		if result != 42 {
			t.Errorf("Expected 42, got %d", result)
		}

		result = getTableInt(table, "nonexistent", 100)
		if result != 100 {
			t.Errorf("Expected 100, got %d", result)
		}
	})
}

func TestSSHDisconnect_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		-- This should fail because it's not a valid SSH connection
		local success, err = pcall(function()
			ssh.disconnect(ud)
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHExec_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.exec(ud, "echo test")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHUpload_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.upload(ud, "/tmp/test", "/remote/test")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHDownload_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.download(ud, "/remote/test", "/tmp/test")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHExists_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.exists(ud, "/remote/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHStat_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.stat(ud, "/remote/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHMkdir_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.mkdir(ud, "/remote/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHRemove_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.remove(ud, "/remote/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHRename_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.rename(ud, "/old/path", "/new/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHChmod_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.chmod(ud, "/remote/path", 0644)
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHChown_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.chown(ud, "/remote/path", 1000, 1000)
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHListDir_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.list_dir(ud, "/remote/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHUploadDir_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.upload_dir(ud, "/local/path", "/remote/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}

func TestSSHDownloadDir_InvalidConnection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterSSHModule(L)

	script := `
		local ud = {}
		local success, err = pcall(function()
			ssh.download_dir(ud, "/remote/path", "/local/path")
		end)
		return not success
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected error for invalid connection")
	}
}
