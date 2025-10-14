package infra

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// RegisterNixOSModule registers the NixOS management module in the Lua state
func RegisterNixOSModule(L *lua.LState) {
	// Create nixos module table
	nixosModule := L.NewTable()

	// User management functions
	L.SetField(nixosModule, "add_user", L.NewFunction(nixosAddUser))
	L.SetField(nixosModule, "remove_user", L.NewFunction(nixosRemoveUser))
	L.SetField(nixosModule, "user_exists", L.NewFunction(nixosUserExists))

	// SSH key management functions
	L.SetField(nixosModule, "add_ssh_key", L.NewFunction(nixosAddSSHKey))
	L.SetField(nixosModule, "remove_ssh_key", L.NewFunction(nixosRemoveSSHKey))

	// System management functions
	L.SetField(nixosModule, "rebuild", L.NewFunction(nixosRebuild))
	L.SetField(nixosModule, "get_config", L.NewFunction(nixosGetConfig))
	L.SetField(nixosModule, "backup_config", L.NewFunction(nixosBackupConfig))

	// Package management functions
	L.SetField(nixosModule, "add_package", L.NewFunction(nixosAddPackage))
	L.SetField(nixosModule, "remove_package", L.NewFunction(nixosRemovePackage))

	// Set as global
	L.SetGlobal("nixos", nixosModule)
}

// nixosAddUser adds a user to NixOS configuration.nix (idempotent)
// Usage: local success, msg = nixos.add_user({username = "user", groups = {"wheel"}, description = "User"})
func nixosAddUser(L *lua.LState) int {
	params := L.CheckTable(1)

	username := getStringField(L, params, "username", "")
	description := getStringField(L, params, "description", "")
	shell := getStringField(L, params, "shell", "/bin/bash")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")
	isNormalUser := getBoolField(L, params, "is_normal_user", true)

	if username == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("username is required"))
		return 2
	}

	// Get groups
	var groups []string
	groupsTable := L.GetField(params, "groups")
	if tbl, ok := groupsTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				groups = append(groups, string(str))
			}
		})
	}

	// Check if user already exists in config
	exists, err := userExistsInConfig(configPath, username)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to check config: %v", err)))
		return 2
	}

	if exists {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("user %s already exists in configuration (idempotent)", username)))
		return 2
	}

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	// Build user configuration
	userConfig := buildUserConfig(username, description, shell, groups, isNormalUser)

	// Insert user into config
	newContent := insertUserIntoConfig(string(content), userConfig)

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("user %s added to configuration successfully", username)))
	return 2
}

// nixosRemoveUser removes a user from NixOS configuration.nix (idempotent)
func nixosRemoveUser(L *lua.LState) int {
	params := L.CheckTable(1)

	username := getStringField(L, params, "username", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if username == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("username is required"))
		return 2
	}

	// Check if user exists in config
	exists, err := userExistsInConfig(configPath, username)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to check config: %v", err)))
		return 2
	}

	if !exists {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("user %s not in configuration (idempotent)", username)))
		return 2
	}

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	// Remove user from config
	newContent := removeUserFromConfig(string(content), username)

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("user %s removed from configuration", username)))
	return 2
}

// nixosUserExists checks if a user exists in configuration.nix
func nixosUserExists(L *lua.LState) int {
	params := L.CheckTable(1)

	username := getStringField(L, params, "username", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if username == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("username is required"))
		return 2
	}

	exists, err := userExistsInConfig(configPath, username)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("error checking user: %v", err)))
		return 2
	}

	L.Push(lua.LBool(exists))
	L.Push(lua.LNil)
	return 2
}

// nixosAddSSHKey adds an SSH key to a user in configuration.nix (idempotent)
func nixosAddSSHKey(L *lua.LState) int {
	params := L.CheckTable(1)

	username := getStringField(L, params, "username", "")
	sshKey := getStringField(L, params, "key", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if username == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("username is required"))
		return 2
	}

	if sshKey == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("key is required"))
		return 2
	}

	// Normalize the key
	sshKey = strings.TrimSpace(sshKey)

	// Check if key already exists
	hasKey, err := userHasSSHKey(configPath, username, sshKey)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to check SSH key: %v", err)))
		return 2
	}

	if hasKey {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("SSH key already present for user %s (idempotent)", username)))
		return 2
	}

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	// Add SSH key to user
	newContent := addSSHKeyToUser(string(content), username, sshKey)

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("SSH key added to user %s successfully", username)))
	return 2
}

// nixosRemoveSSHKey removes an SSH key from a user in configuration.nix (idempotent)
func nixosRemoveSSHKey(L *lua.LState) int {
	params := L.CheckTable(1)

	username := getStringField(L, params, "username", "")
	sshKey := getStringField(L, params, "key", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if username == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("username is required"))
		return 2
	}

	if sshKey == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("key is required"))
		return 2
	}

	// Normalize the key
	sshKey = strings.TrimSpace(sshKey)

	// Check if key exists
	hasKey, err := userHasSSHKey(configPath, username, sshKey)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to check SSH key: %v", err)))
		return 2
	}

	if !hasKey {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("SSH key not present for user %s (idempotent)", username)))
		return 2
	}

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	// Remove SSH key from user
	newContent := removeSSHKeyFromUser(string(content), username, sshKey)

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("SSH key removed from user %s", username)))
	return 2
}

// nixosRebuild executes nixos-rebuild switch
func nixosRebuild(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())

	action := getStringField(L, params, "action", "switch")
	upgrade := getBoolField(L, params, "upgrade", false)
	useSudo := getBoolField(L, params, "use_sudo", true)

	// Build command
	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nixos-rebuild", action)
	if upgrade {
		cmdArgs = append(cmdArgs, "--upgrade")
	}

	// Execute command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("nixos-rebuild failed: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("nixos-rebuild %s completed successfully", action)))
	return 2
}

// nixosGetConfig reads the configuration.nix file
func nixosGetConfig(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())

	configPath := getStringField(L, params, "path", "/etc/nixos/configuration.nix")

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	L.Push(lua.LString(string(content)))
	L.Push(lua.LNil)
	return 2
}

// nixosBackupConfig creates a backup of configuration.nix
func nixosBackupConfig(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())

	configPath := getStringField(L, params, "path", "/etc/nixos/configuration.nix")
	backupPath := getStringField(L, params, "backup_path", configPath+".backup")

	// Read original config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	// Write backup
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write backup: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("configuration backed up to %s", backupPath)))
	return 2
}

// nixosAddPackage adds a package to system packages (idempotent)
func nixosAddPackage(L *lua.LState) int {
	params := L.CheckTable(1)

	packageName := getStringField(L, params, "package", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if packageName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("package name is required"))
		return 2
	}

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	// Check if package already exists
	if strings.Contains(configStr, "pkgs."+packageName) || strings.Contains(configStr, packageName) {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("package %s already in configuration (idempotent)", packageName)))
		return 2
	}

	// Add package to environment.systemPackages
	newContent := addPackageToConfig(configStr, packageName)

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("package %s added to configuration", packageName)))
	return 2
}

// nixosRemovePackage removes a package from system packages (idempotent)
func nixosRemovePackage(L *lua.LState) int {
	params := L.CheckTable(1)

	packageName := getStringField(L, params, "package", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if packageName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("package name is required"))
		return 2
	}

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	// Check if package exists
	if !strings.Contains(configStr, "pkgs."+packageName) && !strings.Contains(configStr, packageName) {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("package %s not in configuration (idempotent)", packageName)))
		return 2
	}

	// Remove package from config
	newContent := removePackageFromConfig(configStr, packageName)

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("package %s removed from configuration", packageName)))
	return 2
}

// Helper functions

func userExistsInConfig(configPath, username string) (bool, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return false, err
	}

	// Look for users.users.<username> = {
	pattern := fmt.Sprintf(`users\.users\.%s\s*=\s*{`, regexp.QuoteMeta(username))
	matched, err := regexp.MatchString(pattern, string(content))
	return matched, err
}

func buildUserConfig(username, description, shell string, groups []string, isNormalUser bool) string {
	var config strings.Builder

	config.WriteString(fmt.Sprintf("  users.users.%s = {\n", username))
	if isNormalUser {
		config.WriteString("    isNormalUser = true;\n")
	} else {
		config.WriteString("    isSystemUser = true;\n")
	}
	config.WriteString(fmt.Sprintf("    home = \"/home/%s\";\n", username))
	if description != "" {
		config.WriteString(fmt.Sprintf("    description = \"%s\";\n", description))
	}
	if shell != "" {
		config.WriteString(fmt.Sprintf("    shell = %s;\n", shell))
	}
	if len(groups) > 0 {
		config.WriteString("    extraGroups = [ ")
		for _, group := range groups {
			config.WriteString(fmt.Sprintf("\"%s\" ", group))
		}
		config.WriteString("];\n")
	}
	config.WriteString("  };\n")

	return config.String()
}

func insertUserIntoConfig(content, userConfig string) string {
	// Try to find users.users section
	usersPattern := regexp.MustCompile(`(users\.users\s*=\s*{[^}]*)(};)`)
	if usersPattern.MatchString(content) {
		// Insert before closing };
		return usersPattern.ReplaceAllString(content, "${1}\n"+userConfig+"${2}")
	}

	// If users.users doesn't exist, create it
	// Find a good place to insert (before the closing } of the main config)
	lines := strings.Split(content, "\n")
	var result strings.Builder
	inserted := false

	for i, line := range lines {
		result.WriteString(line)
		result.WriteString("\n")

		// Insert before the last closing brace
		if !inserted && i > 0 && strings.TrimSpace(line) == "}" && i == len(lines)-2 {
			result.WriteString("\n  users.users = {\n")
			result.WriteString(userConfig)
			result.WriteString("  };\n")
			inserted = true
		}
	}

	if !inserted {
		// Fallback: append at the end
		return content + "\n  users.users = {\n" + userConfig + "  };\n"
	}

	return result.String()
}

func removeUserFromConfig(content, username string) string {
	// Remove the entire user block
	pattern := fmt.Sprintf(`\s*users\.users\.%s\s*=\s*{[^}]*};\s*`, regexp.QuoteMeta(username))
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(content, "")
}

func userHasSSHKey(configPath, username, sshKey string) (bool, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return false, err
	}

	// Normalize the key for comparison
	sshKey = strings.TrimSpace(sshKey)

	// Look for the key in the user's openssh.authorizedKeys.keys
	return strings.Contains(string(content), sshKey), nil
}

func addSSHKeyToUser(content, username, sshKey string) string {
	// Check if user already has openssh.authorizedKeys.keys section
	userPattern := fmt.Sprintf(`(users\.users\.%s\s*=\s*{[^}]*)(};)`, regexp.QuoteMeta(username))
	re := regexp.MustCompile(userPattern)

	if re.MatchString(content) {
		// Check if openssh.authorizedKeys.keys already exists
		if strings.Contains(content, username) && strings.Contains(content, "openssh.authorizedKeys.keys") {
			// Add to existing keys array
			keysPattern := fmt.Sprintf(`(users\.users\.%s[^}]*openssh\.authorizedKeys\.keys\s*=\s*\[)([^\]]*)(];)`,
				regexp.QuoteMeta(username))
			keysRe := regexp.MustCompile(keysPattern)
			return keysRe.ReplaceAllString(content, fmt.Sprintf(`${1}${2}      "%s"%s${3}`, sshKey, "\n"))
		} else {
			// Add openssh.authorizedKeys.keys section
			sshConfig := fmt.Sprintf("    openssh.authorizedKeys.keys = [\n      \"%s\"\n    ];\n", sshKey)
			return re.ReplaceAllString(content, "${1}\n"+sshConfig+"  ${2}")
		}
	}

	return content
}

func removeSSHKeyFromUser(content, username, sshKey string) string {
	// Remove the specific SSH key line
	keyPattern := fmt.Sprintf(`\s*"%s"\s*`, regexp.QuoteMeta(sshKey))
	re := regexp.MustCompile(keyPattern)
	return re.ReplaceAllString(content, "")
}

func addPackageToConfig(content, packageName string) string {
	// Try to find environment.systemPackages section
	packagesPattern := regexp.MustCompile(`(environment\.systemPackages\s*=\s*with\s+pkgs;\s*\[)([^\]]*)(];)`)
	if packagesPattern.MatchString(content) {
		// Add package to existing list
		return packagesPattern.ReplaceAllString(content, fmt.Sprintf("${1}${2}    %s\n  ${3}", packageName))
	}

	// If section doesn't exist, create it
	return content + fmt.Sprintf("\n  environment.systemPackages = with pkgs; [\n    %s\n  ];\n", packageName)
}

func removePackageFromConfig(content, packageName string) string {
	// Remove package line
	patterns := []string{
		fmt.Sprintf(`\s*pkgs\.%s\s*`, regexp.QuoteMeta(packageName)),
		fmt.Sprintf(`\s*%s\s*`, regexp.QuoteMeta(packageName)),
	}

	result := content
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, "")
	}

	return result
}
