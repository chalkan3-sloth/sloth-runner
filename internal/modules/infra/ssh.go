package infra

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// RegisterSSHModule registers the SSH management module in the Lua state
func RegisterSSHModule(L *lua.LState) {
	// Create ssh module table
	sshModule := L.NewTable()

	// Register functions
	L.SetField(sshModule, "add_authorized_key", L.NewFunction(sshAddAuthorizedKey))
	L.SetField(sshModule, "remove_authorized_key", L.NewFunction(sshRemoveAuthorizedKey))
	L.SetField(sshModule, "list_authorized_keys", L.NewFunction(sshListAuthorizedKeys))
	L.SetField(sshModule, "ensure_ssh_dir", L.NewFunction(sshEnsureSSHDir))
	L.SetField(sshModule, "key_exists", L.NewFunction(sshKeyExists))

	// Set as global
	L.SetGlobal("ssh", sshModule)
}

// Helper function to get string field from table
func getStringField(L *lua.LState, tbl *lua.LTable, key, defaultValue string) string {
	lv := L.GetField(tbl, key)
	if str, ok := lv.(lua.LString); ok {
		return string(str)
	}
	return defaultValue
}

// Helper function to get bool field from table
func getBoolField(L *lua.LState, tbl *lua.LTable, key string, defaultValue bool) bool {
	lv := L.GetField(tbl, key)
	if b, ok := lv.(lua.LBool); ok {
		return bool(b)
	}
	return defaultValue
}

// sshAddAuthorizedKey adds an SSH public key to a user's authorized_keys file (idempotent)
// Usage: local success, msg = ssh.add_authorized_key({user = "username", key = "ssh-ed25519 AAAA..."})
func sshAddAuthorizedKey(L *lua.LState) int {
	params := L.CheckTable(1)

	user := getStringField(L, params, "user", "")
	publicKey := getStringField(L, params, "key", "")
	comment := getStringField(L, params, "comment", "")

	if user == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("user is required"))
		return 2
	}

	if publicKey == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("key is required"))
		return 2
	}

	// Normalize the key (trim whitespace)
	publicKey = strings.TrimSpace(publicKey)

	// Add comment if provided
	keyLine := publicKey
	if comment != "" {
		keyLine = publicKey + " " + comment
	}

	// Get home directory
	homeDir := filepath.Join("/home", user)
	if user == "root" {
		homeDir = "/root"
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	authKeysFile := filepath.Join(sshDir, "authorized_keys")

	// Ensure .ssh directory exists
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create .ssh directory: %v", err)))
		return 2
	}

	// Set ownership on .ssh directory
	if err := setOwnership(sshDir, user); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to set ownership on .ssh directory: %v", err)))
		return 2
	}

	// IDEMPOTENCY CHECK: Check if key already exists
	keyExists, err := checkKeyExists(authKeysFile, publicKey)
	if err != nil && !os.IsNotExist(err) {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to check existing keys: %v", err)))
		return 2
	}

	if keyExists {
		L.Push(lua.LBool(true))
		L.Push(lua.LString("SSH key already present (idempotent)"))
		return 2
	}

	// Open file for appending (create if doesn't exist)
	f, err := os.OpenFile(authKeysFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to open authorized_keys: %v", err)))
		return 2
	}
	defer f.Close()

	// Append the key
	if _, err := f.WriteString(keyLine + "\n"); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write key: %v", err)))
		return 2
	}

	// Set ownership on authorized_keys
	if err := setOwnership(authKeysFile, user); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to set ownership on authorized_keys: %v", err)))
		return 2
	}

	// Ensure correct permissions
	if err := os.Chmod(authKeysFile, 0600); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to set permissions: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("SSH key added successfully"))
	return 2
}

// sshRemoveAuthorizedKey removes an SSH public key from authorized_keys (idempotent)
// Usage: local success, msg = ssh.remove_authorized_key({user = "username", key = "ssh-ed25519 AAAA..."})
func sshRemoveAuthorizedKey(L *lua.LState) int {
	params := L.CheckTable(1)

	user := getStringField(L, params, "user", "")
	publicKey := getStringField(L, params, "key", "")

	if user == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("user is required"))
		return 2
	}

	if publicKey == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("key is required"))
		return 2
	}

	// Normalize the key
	publicKey = strings.TrimSpace(publicKey)

	// Get file path
	homeDir := filepath.Join("/home", user)
	if user == "root" {
		homeDir = "/root"
	}
	authKeysFile := filepath.Join(homeDir, ".ssh", "authorized_keys")

	// IDEMPOTENCY CHECK: If file doesn't exist, nothing to do
	if _, err := os.Stat(authKeysFile); os.IsNotExist(err) {
		L.Push(lua.LBool(true))
		L.Push(lua.LString("authorized_keys doesn't exist (idempotent)"))
		return 2
	}

	// Check if key exists
	keyExists, err := checkKeyExists(authKeysFile, publicKey)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read authorized_keys: %v", err)))
		return 2
	}

	if !keyExists {
		L.Push(lua.LBool(true))
		L.Push(lua.LString("SSH key not present (idempotent)"))
		return 2
	}

	// Read file
	content, err := os.ReadFile(authKeysFile)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read file: %v", err)))
		return 2
	}

	// Filter out the key
	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			newLines = append(newLines, line)
			continue
		}
		// Check if this line contains the key
		if !strings.Contains(line, publicKey) {
			newLines = append(newLines, line)
		}
	}

	// Write back
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(authKeysFile, []byte(newContent), 0600); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write file: %v", err)))
		return 2
	}

	// Set ownership
	if err := setOwnership(authKeysFile, user); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to set ownership: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("SSH key removed successfully"))
	return 2
}

// sshListAuthorizedKeys lists all authorized keys for a user
// Usage: local keys, err = ssh.list_authorized_keys({user = "username"})
func sshListAuthorizedKeys(L *lua.LState) int {
	params := L.CheckTable(1)

	user := getStringField(L, params, "user", "")

	if user == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("user is required"))
		return 2
	}

	// Get file path
	homeDir := filepath.Join("/home", user)
	if user == "root" {
		homeDir = "/root"
	}
	authKeysFile := filepath.Join(homeDir, ".ssh", "authorized_keys")

	// Check if file exists
	if _, err := os.Stat(authKeysFile); os.IsNotExist(err) {
		// Return empty table
		L.Push(L.NewTable())
		L.Push(lua.LNil)
		return 2
	}

	// Read file
	file, err := os.Open(authKeysFile)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to open file: %v", err)))
		return 2
	}
	defer file.Close()

	// Parse keys
	keys := L.NewTable()
	scanner := bufio.NewScanner(file)
	index := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			keys.RawSetInt(index, lua.LString(line))
			index++
		}
	}

	if err := scanner.Err(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("error reading file: %v", err)))
		return 2
	}

	L.Push(keys)
	L.Push(lua.LNil)
	return 2
}

// sshEnsureSSHDir ensures .ssh directory exists with correct permissions
// Usage: local success, msg = ssh.ensure_ssh_dir({user = "username"})
func sshEnsureSSHDir(L *lua.LState) int {
	params := L.CheckTable(1)

	user := getStringField(L, params, "user", "")

	if user == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("user is required"))
		return 2
	}

	// Get home directory
	homeDir := filepath.Join("/home", user)
	if user == "root" {
		homeDir = "/root"
	}

	sshDir := filepath.Join(homeDir, ".ssh")

	// Check if directory already exists
	if info, err := os.Stat(sshDir); err == nil {
		if !info.IsDir() {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(".ssh exists but is not a directory"))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(".ssh directory already exists (idempotent)"))
		return 2
	}

	// Create directory
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create .ssh directory: %v", err)))
		return 2
	}

	// Set ownership
	if err := setOwnership(sshDir, user); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to set ownership: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(".ssh directory created successfully"))
	return 2
}

// sshKeyExists checks if a key exists in authorized_keys
// Usage: local exists, err = ssh.key_exists({user = "username", key = "ssh-ed25519 AAAA..."})
func sshKeyExists(L *lua.LState) int {
	params := L.CheckTable(1)

	user := getStringField(L, params, "user", "")
	publicKey := getStringField(L, params, "key", "")

	if user == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("user is required"))
		return 2
	}

	if publicKey == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("key is required"))
		return 2
	}

	// Normalize the key
	publicKey = strings.TrimSpace(publicKey)

	// Get file path
	homeDir := filepath.Join("/home", user)
	if user == "root" {
		homeDir = "/root"
	}
	authKeysFile := filepath.Join(homeDir, ".ssh", "authorized_keys")

	// Check if key exists
	exists, err := checkKeyExists(authKeysFile, publicKey)
	if err != nil {
		if os.IsNotExist(err) {
			L.Push(lua.LBool(false))
			L.Push(lua.LNil)
			return 2
		}
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("error checking key: %v", err)))
		return 2
	}

	L.Push(lua.LBool(exists))
	L.Push(lua.LNil)
	return 2
}

// Helper functions

// checkKeyExists checks if a public key exists in the authorized_keys file
func checkKeyExists(authKeysFile, publicKey string) (bool, error) {
	file, err := os.Open(authKeysFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Normalize the search key
	searchKey := strings.TrimSpace(publicKey)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Check if this line contains the key (may have comment at end)
		if strings.Contains(line, searchKey) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

// setOwnership sets the ownership of a file/directory to user:user
func setOwnership(path, user string) error {
	if user == "root" || user == "" {
		return nil // No need to change ownership for root
	}

	cmd := exec.Command("chown", user+":"+user, path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("chown failed: %v", err)
	}
	return nil
}
