package luainterface

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/crypto/ssh"
)

// RegisterSSHModule registers the SSH module with comprehensive functionality
func RegisterSSHModule(L *lua.LState) {
	mod := L.NewTable()

	// Connection functions
	L.SetField(mod, "connect", L.NewFunction(sshConnect))
	L.SetField(mod, "disconnect", L.NewFunction(sshDisconnect))
	L.SetField(mod, "exec", L.NewFunction(sshExec))

	// File transfer functions
	L.SetField(mod, "upload", L.NewFunction(sshUpload))
	L.SetField(mod, "download", L.NewFunction(sshDownload))
	L.SetField(mod, "upload_dir", L.NewFunction(sshUploadDir))
	L.SetField(mod, "download_dir", L.NewFunction(sshDownloadDir))

	// SFTP functions
	L.SetField(mod, "exists", L.NewFunction(sshExists))
	L.SetField(mod, "stat", L.NewFunction(sshStat))
	L.SetField(mod, "mkdir", L.NewFunction(sshMkdir))
	L.SetField(mod, "remove", L.NewFunction(sshRemove))
	L.SetField(mod, "rename", L.NewFunction(sshRename))
	L.SetField(mod, "chmod", L.NewFunction(sshChmod))
	L.SetField(mod, "chown", L.NewFunction(sshChown))
	L.SetField(mod, "list_dir", L.NewFunction(sshListDir))

	// Key management
	L.SetField(mod, "load_private_key", L.NewFunction(sshLoadPrivateKey))
	L.SetField(mod, "generate_keypair", L.NewFunction(sshGenerateKeypair))

	// Tunneling
	L.SetField(mod, "create_tunnel", L.NewFunction(sshCreateTunnel))
	L.SetField(mod, "close_tunnel", L.NewFunction(sshCloseTunnel))

	// Agent forwarding
	L.SetField(mod, "enable_agent_forward", L.NewFunction(sshEnableAgentForward))

	// Local authorized_keys management (idempotent)
	L.SetField(mod, "add_authorized_key", L.NewFunction(sshAddAuthorizedKey))
	L.SetField(mod, "remove_authorized_key", L.NewFunction(sshRemoveAuthorizedKey))
	L.SetField(mod, "list_authorized_keys", L.NewFunction(sshListAuthorizedKeys))
	L.SetField(mod, "key_exists", L.NewFunction(sshKeyExists))

	L.SetGlobal("ssh", mod)
}

// sshConnect connects to an SSH server
// ssh.connect(host, user, options)
func sshConnect(L *lua.LState) int {
	host := L.CheckString(1)
	user := L.CheckString(2)
	options := L.OptTable(3, L.NewTable())

	port := getTableInt(options, "port", 22)
	password := getTableString(options, "password", "")
	keyPath := getTableString(options, "key_path", "")
	timeout := getTableInt(options, "timeout", 30)

	var authMethods []ssh.AuthMethod

	// Password authentication
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// Key-based authentication
	if keyPath != "" {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to read key: %v", err)))
			return 2
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to parse key: %v", err)))
			return 2
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key validation
		Timeout:         time.Duration(timeout) * time.Second,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to connect: %v", err)))
		return 2
	}

	// Store connection in a user data
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("ssh_connection"))

	L.Push(ud)
	return 1
}

// sshDisconnect closes an SSH connection
func sshDisconnect(L *lua.LState) int {
	ud := L.CheckUserData(1)
	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	if err := client.Close(); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to disconnect: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshExec executes a command on the remote server
func sshExec(L *lua.LState) int {
	ud := L.CheckUserData(1)
	command := L.CheckString(2)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	session, err := client.NewSession()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create session: %v", err)))
		return 2
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
		} else {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to execute: %v", err)))
			return 2
		}
	}

	result := L.NewTable()
	L.SetField(result, "stdout", lua.LString(stdout.String()))
	L.SetField(result, "stderr", lua.LString(stderr.String()))
	L.SetField(result, "exit_code", lua.LNumber(exitCode))
	L.SetField(result, "success", lua.LBool(exitCode == 0))

	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// sshUpload uploads a file to the remote server
func sshUpload(L *lua.LState) int {
	ud := L.CheckUserData(1)
	localPath := L.CheckString(2)
	remotePath := L.CheckString(3)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	// Open local file
	srcFile, err := os.Open(localPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to open local file: %v", err)))
		return 2
	}
	defer srcFile.Close()

	// Create remote file
	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create remote file: %v", err)))
		return 2
	}
	defer dstFile.Close()

	// Copy file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to copy file: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshDownload downloads a file from the remote server
func sshDownload(L *lua.LState) int {
	ud := L.CheckUserData(1)
	remotePath := L.CheckString(2)
	localPath := L.CheckString(3)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	// Open remote file
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to open remote file: %v", err)))
		return 2
	}
	defer srcFile.Close()

	// Create local file
	dstFile, err := os.Create(localPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create local file: %v", err)))
		return 2
	}
	defer dstFile.Close()

	// Copy file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to copy file: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshUploadDir uploads a directory recursively
func sshUploadDir(L *lua.LState) int {
	ud := L.CheckUserData(1)
	localPath := L.CheckString(2)
	remotePath := L.CheckString(3)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	err = filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(localPath, path)
		if err != nil {
			return err
		}

		remoteFilePath := filepath.Join(remotePath, relPath)

		if info.IsDir() {
			return sftpClient.MkdirAll(remoteFilePath)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := sftpClient.Create(remoteFilePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to upload directory: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshDownloadDir downloads a directory recursively
func sshDownloadDir(L *lua.LState) int {
	ud := L.CheckUserData(1)
	remotePath := L.CheckString(2)
	localPath := L.CheckString(3)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	walker := sftpClient.Walk(remotePath)
	for walker.Step() {
		if walker.Err() != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("walk error: %v", walker.Err())))
			return 2
		}

		relPath, err := filepath.Rel(remotePath, walker.Path())
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to get relative path: %v", err)))
			return 2
		}

		localFilePath := filepath.Join(localPath, relPath)

		if walker.Stat().IsDir() {
			if err := os.MkdirAll(localFilePath, 0755); err != nil {
				L.Push(lua.LBool(false))
				L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
				return 2
			}
			continue
		}

		srcFile, err := sftpClient.Open(walker.Path())
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to open remote file: %v", err)))
			return 2
		}

		dstFile, err := os.Create(localFilePath)
		if err != nil {
			srcFile.Close()
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to create local file: %v", err)))
			return 2
		}

		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to copy file: %v", err)))
			return 2
		}
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshExists checks if a file or directory exists on the remote server
func sshExists(L *lua.LState) int {
	ud := L.CheckUserData(1)
	path := L.CheckString(2)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		return 1
	}
	defer sftpClient.Close()

	_, err = sftpClient.Stat(path)
	L.Push(lua.LBool(err == nil))
	return 1
}

// sshStat gets file/directory information
func sshStat(L *lua.LState) int {
	ud := L.CheckUserData(1)
	path := L.CheckString(2)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	info, err := sftpClient.Stat(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to stat: %v", err)))
		return 2
	}

	result := L.NewTable()
	L.SetField(result, "name", lua.LString(info.Name()))
	L.SetField(result, "size", lua.LNumber(info.Size()))
	L.SetField(result, "mode", lua.LString(info.Mode().String()))
	L.SetField(result, "is_dir", lua.LBool(info.IsDir()))
	L.SetField(result, "mod_time", lua.LNumber(info.ModTime().Unix()))

	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// sshMkdir creates a directory on the remote server
func sshMkdir(L *lua.LState) int {
	ud := L.CheckUserData(1)
	path := L.CheckString(2)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	err = sftpClient.MkdirAll(path)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshRemove removes a file or directory on the remote server
func sshRemove(L *lua.LState) int {
	ud := L.CheckUserData(1)
	path := L.CheckString(2)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	err = sftpClient.Remove(path)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to remove: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshRename renames a file or directory on the remote server
func sshRename(L *lua.LState) int {
	ud := L.CheckUserData(1)
	oldPath := L.CheckString(2)
	newPath := L.CheckString(3)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	err = sftpClient.Rename(oldPath, newPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to rename: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshChmod changes file permissions on the remote server
func sshChmod(L *lua.LState) int {
	ud := L.CheckUserData(1)
	path := L.CheckString(2)
	mode := L.CheckInt(3)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	err = sftpClient.Chmod(path, os.FileMode(mode))
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to chmod: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshChown changes file ownership on the remote server
func sshChown(L *lua.LState) int {
	ud := L.CheckUserData(1)
	path := L.CheckString(2)
	uid := L.CheckInt(3)
	gid := L.CheckInt(4)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	err = sftpClient.Chown(path, uid, gid)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to chown: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}

// sshListDir lists files in a directory on the remote server
func sshListDir(L *lua.LState) int {
	ud := L.CheckUserData(1)
	path := L.CheckString(2)

	client, ok := ud.Value.(*ssh.Client)
	if !ok {
		L.ArgError(1, "ssh connection expected")
		return 0
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create SFTP client: %v", err)))
		return 2
	}
	defer sftpClient.Close()

	files, err := sftpClient.ReadDir(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to list directory: %v", err)))
		return 2
	}

	result := L.NewTable()
	for i, file := range files {
		fileInfo := L.NewTable()
		L.SetField(fileInfo, "name", lua.LString(file.Name()))
		L.SetField(fileInfo, "size", lua.LNumber(file.Size()))
		L.SetField(fileInfo, "mode", lua.LString(file.Mode().String()))
		L.SetField(fileInfo, "is_dir", lua.LBool(file.IsDir()))
		L.SetField(fileInfo, "mod_time", lua.LNumber(file.ModTime().Unix()))
		result.RawSetInt(i+1, fileInfo)
	}

	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// sshLoadPrivateKey loads a private key from file
func sshLoadPrivateKey(L *lua.LState) int {
	keyPath := L.CheckString(1)

	key, err := os.ReadFile(keyPath)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read key: %v", err)))
		return 2
	}

	_, err = ssh.ParsePrivateKey(key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse key: %v", err)))
		return 2
	}

	L.Push(lua.LString(string(key)))
	return 1
}

// sshGenerateKeypair generates an SSH key pair
func sshGenerateKeypair(L *lua.LState) int {
	// Placeholder for key generation - would need crypto/rsa or similar
	L.Push(lua.LNil)
	L.Push(lua.LString("not implemented"))
	return 2
}

// sshCreateTunnel creates an SSH tunnel
func sshCreateTunnel(L *lua.LState) int {
	// Placeholder for tunnel creation
	L.Push(lua.LNil)
	L.Push(lua.LString("not implemented"))
	return 2
}

// sshCloseTunnel closes an SSH tunnel
func sshCloseTunnel(L *lua.LState) int {
	// Placeholder for tunnel closing
	L.Push(lua.LBool(false))
	L.Push(lua.LString("not implemented"))
	return 2
}

// sshEnableAgentForward enables SSH agent forwarding
func sshEnableAgentForward(L *lua.LState) int {
	// Placeholder for agent forwarding
	L.Push(lua.LBool(false))
	L.Push(lua.LString("not implemented"))
	return 2
}

// Helper function to get string from table
func getTableString(t *lua.LTable, key string, def string) string {
	lv := t.RawGetString(key)
	if s, ok := lv.(lua.LString); ok {
		return string(s)
	}
	return def
}

// Helper function to get int from table
func getTableInt(t *lua.LTable, key string, def int) int {
	lv := t.RawGetString(key)
	if n, ok := lv.(lua.LNumber); ok {
		return int(n)
	}
	return def
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
