package infra

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// NFSSMBModule provides NFS and Samba file sharing management
type NFSSMBModule struct {
	L *lua.LState
}

// NewNFSSMBModule creates a new NFS/SMB module instance
func NewNFSSMBModule(L *lua.LState) *NFSSMBModule {
	return &NFSSMBModule{L: L}
}

// Register registers the NFS/SMB module with the Lua state
func (m *NFSSMBModule) Register(L *lua.LState) {
	// NFS table
	nfsTable := L.NewTable()
	L.SetField(nfsTable, "export", L.NewFunction(m.nfsExport))
	L.SetField(nfsTable, "unexport", L.NewFunction(m.nfsUnexport))
	L.SetField(nfsTable, "list_exports", L.NewFunction(m.nfsListExports))
	L.SetField(nfsTable, "reload", L.NewFunction(m.nfsReload))
	L.SetField(nfsTable, "is_exported", L.NewFunction(m.nfsIsExported))
	L.SetGlobal("nfs", nfsTable)

	// SMB/Samba table
	smbTable := L.NewTable()
	L.SetField(smbTable, "share", L.NewFunction(m.smbShare))
	L.SetField(smbTable, "unshare", L.NewFunction(m.smbUnshare))
	L.SetField(smbTable, "list_shares", L.NewFunction(m.smbListShares))
	L.SetField(smbTable, "reload", L.NewFunction(m.smbReload))
	L.SetField(smbTable, "is_shared", L.NewFunction(m.smbIsShared))
	L.SetField(smbTable, "add_user", L.NewFunction(m.smbAddUser))
	L.SetField(smbTable, "remove_user", L.NewFunction(m.smbRemoveUser))
	L.SetGlobal("smb", smbTable)
}

func (m *NFSSMBModule) execCommand(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func (m *NFSSMBModule) needsSudo() bool {
	output, err := m.execCommand("id", "-u")
	if err != nil {
		return true
	}
	return output != "0"
}

func (m *NFSSMBModule) prependSudo(args []string) []string {
	if m.needsSudo() {
		return append([]string{"sudo"}, args...)
	}
	return args
}

// ============================
// NFS Functions
// ============================

const nfsExportsFile = "/etc/exports"

// nfsExport exports a directory via NFS
func (m *NFSSMBModule) nfsExport(L *lua.LState) int {
	options := L.CheckTable(1)

	path := options.RawGetString("path").String()
	clients := options.RawGetString("clients").String()
	opts := options.RawGetString("options").String()

	if path == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Export path is required"))
		return 2
	}

	if clients == "" {
		clients = "*" // Default to all clients
	}

	if opts == "" {
		opts = "rw,sync,no_subtree_check" // Default options
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Export path does not exist: %s", path)))
		return 2
	}

	// Read current exports
	content, err := os.ReadFile(nfsExportsFile)
	if err != nil && !os.IsNotExist(err) {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read exports file: %s", err)))
		return 2
	}

	currentExports := string(content)
	exportLine := fmt.Sprintf("%s %s(%s)", path, clients, opts)

	// Check if export already exists (idempotency)
	lines := strings.Split(currentExports, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, path+" ") || trimmedLine == path {
			// Export exists, check if it matches
			if strings.Contains(trimmedLine, clients) && strings.Contains(trimmedLine, opts) {
				L.Push(lua.LTrue)
				L.Push(lua.LString(fmt.Sprintf("NFS export for %s already exists (idempotent)", path)))
				return 2
			}

			// Update existing export
			newLines := []string{}
			for _, l := range lines {
				if strings.TrimSpace(l) == trimmedLine {
					newLines = append(newLines, exportLine)
				} else {
					newLines = append(newLines, l)
				}
			}

			newContent := strings.Join(newLines, "\n")
			if err := m.writeFile(nfsExportsFile, newContent); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(fmt.Sprintf("Failed to update exports file: %s", err)))
				return 2
			}

			// Reload exports
			if err := m.reloadNFS(); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(fmt.Sprintf("Failed to reload NFS exports: %s", err)))
				return 2
			}

			L.Push(lua.LTrue)
			L.Push(lua.LString(fmt.Sprintf("NFS export for %s updated", path)))
			return 2
		}
	}

	// Add new export
	newContent := currentExports
	if !strings.HasSuffix(newContent, "\n") && newContent != "" {
		newContent += "\n"
	}
	newContent += exportLine + "\n"

	if err := m.writeFile(nfsExportsFile, newContent); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to write exports file: %s", err)))
		return 2
	}

	// Reload exports
	if err := m.reloadNFS(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reload NFS exports: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("NFS export for %s added", path)))
	return 2
}

// nfsUnexport removes an NFS export
func (m *NFSSMBModule) nfsUnexport(L *lua.LState) int {
	path := L.CheckString(1)

	// Read current exports
	content, err := os.ReadFile(nfsExportsFile)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read exports file: %s", err)))
		return 2
	}

	currentExports := string(content)
	lines := strings.Split(currentExports, "\n")

	// Check if export exists
	found := false
	newLines := []string{}
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, path+" ") || trimmedLine == path {
			found = true
			continue // Skip this line
		}
		newLines = append(newLines, line)
	}

	if !found {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("NFS export for %s does not exist (idempotent)", path)))
		return 2
	}

	// Write updated exports
	newContent := strings.Join(newLines, "\n")
	if err := m.writeFile(nfsExportsFile, newContent); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to write exports file: %s", err)))
		return 2
	}

	// Reload exports
	if err := m.reloadNFS(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reload NFS exports: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("NFS export for %s removed", path)))
	return 2
}

// nfsListExports lists all NFS exports
func (m *NFSSMBModule) nfsListExports(L *lua.LState) int {
	// Read exports file
	content, err := os.ReadFile(nfsExportsFile)
	if err != nil {
		if os.IsNotExist(err) {
			L.Push(L.NewTable())
			L.Push(lua.LNil)
			return 2
		}
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to read exports file: %s", err)))
		return 2
	}

	exportsTable := L.NewTable()
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		// Parse export line: /path client(options)
		parts := strings.Fields(trimmedLine)
		if len(parts) >= 2 {
			exportInfo := L.NewTable()
			exportInfo.RawSetString("path", lua.LString(parts[0]))

			// Parse clients and options
			clientOpts := strings.Join(parts[1:], " ")
			exportInfo.RawSetString("clients_options", lua.LString(clientOpts))

			exportsTable.Append(exportInfo)
		}
	}

	L.Push(exportsTable)
	L.Push(lua.LNil)
	return 2
}

// nfsReload reloads NFS exports
func (m *NFSSMBModule) nfsReload(L *lua.LState) int {
	if err := m.reloadNFS(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reload NFS exports: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString("NFS exports reloaded"))
	return 2
}

// nfsIsExported checks if a path is exported
func (m *NFSSMBModule) nfsIsExported(L *lua.LState) int {
	path := L.CheckString(1)

	content, err := os.ReadFile(nfsExportsFile)
	if err != nil {
		L.Push(lua.LFalse)
		return 1
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, path+" ") || trimmedLine == path {
			L.Push(lua.LTrue)
			return 1
		}
	}

	L.Push(lua.LFalse)
	return 1
}

// reloadNFS reloads NFS exports
func (m *NFSSMBModule) reloadNFS() error {
	args := m.prependSudo([]string{"exportfs", "-ra"})
	_, err := m.execCommand(args...)
	return err
}

// ============================
// Samba Functions
// ============================

const smbConfigFile = "/etc/samba/smb.conf"

// smbShare creates a Samba share
func (m *NFSSMBModule) smbShare(L *lua.LState) int {
	options := L.CheckTable(1)

	name := options.RawGetString("name").String()
	path := options.RawGetString("path").String()
	comment := options.RawGetString("comment").String()
	writeable := true
	browseable := true
	guestOk := false

	// Parse boolean options
	if options.RawGetString("writeable").Type() == lua.LTBool {
		writeable = bool(options.RawGetString("writeable").(lua.LBool))
	}
	if options.RawGetString("browseable").Type() == lua.LTBool {
		browseable = bool(options.RawGetString("browseable").(lua.LBool))
	}
	if options.RawGetString("guest_ok").Type() == lua.LTBool {
		guestOk = bool(options.RawGetString("guest_ok").(lua.LBool))
	}

	validUsers := options.RawGetString("valid_users").String()

	if name == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Share name is required"))
		return 2
	}

	if path == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Share path is required"))
		return 2
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Share path does not exist: %s", path)))
		return 2
	}

	// Read current config
	content, err := os.ReadFile(smbConfigFile)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read Samba config: %s", err)))
		return 2
	}

	currentConfig := string(content)

	// Check if share already exists
	shareHeader := fmt.Sprintf("[%s]", name)
	if strings.Contains(currentConfig, shareHeader) {
		// Share exists, check if idempotent
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Samba share '%s' already exists", name)))
		return 2
	}

	// Create share section
	shareConfig := fmt.Sprintf("\n[%s]\n", name)
	shareConfig += fmt.Sprintf("   path = %s\n", path)
	if comment != "" {
		shareConfig += fmt.Sprintf("   comment = %s\n", comment)
	}
	shareConfig += fmt.Sprintf("   writeable = %s\n", boolToYesNo(writeable))
	shareConfig += fmt.Sprintf("   browseable = %s\n", boolToYesNo(browseable))
	shareConfig += fmt.Sprintf("   guest ok = %s\n", boolToYesNo(guestOk))
	if validUsers != "" {
		shareConfig += fmt.Sprintf("   valid users = %s\n", validUsers)
	}

	// Append to config
	newContent := currentConfig + shareConfig

	if err := m.writeFile(smbConfigFile, newContent); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to write Samba config: %s", err)))
		return 2
	}

	// Test config
	args := m.prependSudo([]string{"testparm", "-s"})
	if _, err := m.execCommand(args...); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Invalid Samba config: %s", err)))
		return 2
	}

	// Reload Samba
	if err := m.reloadSamba(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reload Samba: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Samba share '%s' created", name)))
	return 2
}

// smbUnshare removes a Samba share
func (m *NFSSMBModule) smbUnshare(L *lua.LState) int {
	name := L.CheckString(1)

	// Read current config
	content, err := os.ReadFile(smbConfigFile)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to read Samba config: %s", err)))
		return 2
	}

	currentConfig := string(content)
	shareHeader := fmt.Sprintf("[%s]", name)

	if !strings.Contains(currentConfig, shareHeader) {
		L.Push(lua.LTrue)
		L.Push(lua.LString(fmt.Sprintf("Samba share '%s' does not exist (idempotent)", name)))
		return 2
	}

	// Remove share section
	lines := strings.Split(currentConfig, "\n")
	newLines := []string{}
	inShareSection := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == shareHeader {
			inShareSection = true
			continue
		}

		if inShareSection && strings.HasPrefix(trimmedLine, "[") {
			inShareSection = false
		}

		if !inShareSection {
			newLines = append(newLines, line)
		}
	}

	newContent := strings.Join(newLines, "\n")

	if err := m.writeFile(smbConfigFile, newContent); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to write Samba config: %s", err)))
		return 2
	}

	// Reload Samba
	if err := m.reloadSamba(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reload Samba: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Samba share '%s' removed", name)))
	return 2
}

// smbListShares lists all Samba shares
func (m *NFSSMBModule) smbListShares(L *lua.LState) int {
	// Read config file
	content, err := os.ReadFile(smbConfigFile)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to read Samba config: %s", err)))
		return 2
	}

	sharesTable := L.NewTable()
	lines := strings.Split(string(content), "\n")

	var currentShare string
	shareInfo := L.NewTable()

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") || strings.HasPrefix(trimmedLine, ";") {
			continue
		}

		// Check for share header
		if strings.HasPrefix(trimmedLine, "[") && strings.HasSuffix(trimmedLine, "]") {
			// Save previous share
			if currentShare != "" && currentShare != "global" {
				shareInfo.RawSetString("name", lua.LString(currentShare))
				sharesTable.Append(shareInfo)
			}

			// Start new share
			currentShare = strings.TrimSuffix(strings.TrimPrefix(trimmedLine, "["), "]")
			shareInfo = L.NewTable()
			continue
		}

		// Parse share properties
		if currentShare != "" && strings.Contains(trimmedLine, "=") {
			parts := strings.SplitN(trimmedLine, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				shareInfo.RawSetString(key, lua.LString(value))
			}
		}
	}

	// Add last share
	if currentShare != "" && currentShare != "global" {
		shareInfo.RawSetString("name", lua.LString(currentShare))
		sharesTable.Append(shareInfo)
	}

	L.Push(sharesTable)
	L.Push(lua.LNil)
	return 2
}

// smbReload reloads Samba configuration
func (m *NFSSMBModule) smbReload(L *lua.LState) int {
	if err := m.reloadSamba(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to reload Samba: %s", err)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString("Samba configuration reloaded"))
	return 2
}

// smbIsShared checks if a share exists
func (m *NFSSMBModule) smbIsShared(L *lua.LState) int {
	name := L.CheckString(1)

	content, err := os.ReadFile(smbConfigFile)
	if err != nil {
		L.Push(lua.LFalse)
		return 1
	}

	shareHeader := fmt.Sprintf("[%s]", name)
	L.Push(lua.LBool(strings.Contains(string(content), shareHeader)))
	return 1
}

// smbAddUser adds a Samba user
func (m *NFSSMBModule) smbAddUser(L *lua.LState) int {
	username := L.CheckString(1)
	password := L.CheckString(2)

	// Add user with smbpasswd
	args := m.prependSudo([]string{"sh", "-c", fmt.Sprintf("echo -e '%s\\n%s' | smbpasswd -a -s %s", password, password, username)})
	output, err := m.execCommand(args...)
	if err != nil {
		// Check if user already exists
		if strings.Contains(output, "already exists") {
			L.Push(lua.LTrue)
			L.Push(lua.LString(fmt.Sprintf("Samba user '%s' already exists (idempotent)", username)))
			return 2
		}

		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to add Samba user: %s\n%s", err, output)))
		return 2
	}

	// Enable user
	enableArgs := m.prependSudo([]string{"smbpasswd", "-e", username})
	m.execCommand(enableArgs...)

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Samba user '%s' added", username)))
	return 2
}

// smbRemoveUser removes a Samba user
func (m *NFSSMBModule) smbRemoveUser(L *lua.LState) int {
	username := L.CheckString(1)

	args := m.prependSudo([]string{"smbpasswd", "-x", username})
	output, err := m.execCommand(args...)
	if err != nil {
		// Check if user doesn't exist
		if strings.Contains(output, "Failed to find") || strings.Contains(output, "not found") {
			L.Push(lua.LTrue)
			L.Push(lua.LString(fmt.Sprintf("Samba user '%s' does not exist (idempotent)", username)))
			return 2
		}

		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove Samba user: %s\n%s", err, output)))
		return 2
	}

	L.Push(lua.LTrue)
	L.Push(lua.LString(fmt.Sprintf("Samba user '%s' removed", username)))
	return 2
}

// reloadSamba reloads Samba service
func (m *NFSSMBModule) reloadSamba() error {
	// Try systemctl first
	args := m.prependSudo([]string{"systemctl", "reload", "smbd"})
	_, err := m.execCommand(args...)
	if err == nil {
		return nil
	}

	// Try service command
	args = m.prependSudo([]string{"service", "smbd", "reload"})
	_, err = m.execCommand(args...)
	return err
}

// writeFile writes content to a file with sudo if needed
func (m *NFSSMBModule) writeFile(filename, content string) error {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "sloth-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Ensure directory exists
	dir := filepath.Dir(filename)
	mkdirArgs := m.prependSudo([]string{"mkdir", "-p", dir})
	m.execCommand(mkdirArgs...)

	// Copy temp file to target with sudo
	cpArgs := m.prependSudo([]string{"cp", tmpFile.Name(), filename})
	_, err = m.execCommand(cpArgs...)
	return err
}

// boolToYesNo converts bool to yes/no string
func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
