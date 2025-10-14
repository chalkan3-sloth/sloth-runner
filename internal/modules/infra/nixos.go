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

	// Service management functions
	L.SetField(nixosModule, "enable_service", L.NewFunction(nixosEnableService))
	L.SetField(nixosModule, "disable_service", L.NewFunction(nixosDisableService))

	// System options functions
	L.SetField(nixosModule, "set_hostname", L.NewFunction(nixosSetHostname))
	L.SetField(nixosModule, "set_timezone", L.NewFunction(nixosSetTimezone))
	L.SetField(nixosModule, "set_locale", L.NewFunction(nixosSetLocale))
	L.SetField(nixosModule, "enable_firewall", L.NewFunction(nixosEnableFirewall))
	L.SetField(nixosModule, "add_firewall_port", L.NewFunction(nixosAddFirewallPort))

	// Nix store management
	L.SetField(nixosModule, "collect_garbage", L.NewFunction(nixosCollectGarbage))
	L.SetField(nixosModule, "optimize_store", L.NewFunction(nixosOptimizeStore))

	// Channel management
	L.SetField(nixosModule, "update_channels", L.NewFunction(nixosUpdateChannels))
	L.SetField(nixosModule, "list_channels", L.NewFunction(nixosListChannels))

	// Boot configuration
	L.SetField(nixosModule, "set_bootloader", L.NewFunction(nixosSetBootloader))

	// Import management
	L.SetField(nixosModule, "add_import", L.NewFunction(nixosAddImport))
	L.SetField(nixosModule, "remove_import", L.NewFunction(nixosRemoveImport))

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

// ============================================================================
// Service Management Functions
// ============================================================================

// nixosEnableService enables a service in configuration.nix (idempotent)
func nixosEnableService(L *lua.LState) int {
	params := L.CheckTable(1)

	serviceName := getStringField(L, params, "service", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("service name is required"))
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

	// Check if service is already enabled
	enablePattern := fmt.Sprintf(`services\.%s\.enable\s*=\s*true`, regexp.QuoteMeta(serviceName))
	if matched, _ := regexp.MatchString(enablePattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("service %s already enabled (idempotent)", serviceName)))
		return 2
	}

	// Check if service exists but is disabled
	disablePattern := fmt.Sprintf(`services\.%s\.enable\s*=\s*false`, regexp.QuoteMeta(serviceName))
	re := regexp.MustCompile(disablePattern)
	if re.MatchString(configStr) {
		// Change false to true
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf("services.%s.enable = true", serviceName))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("service %s enabled", serviceName)))
		return 2
	}

	// Add new service enable line
	serviceConfig := fmt.Sprintf("  services.%s.enable = true;\n", serviceName)
	newContent := addLineToConfig(configStr, serviceConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("service %s enabled", serviceName)))
	return 2
}

// nixosDisableService disables a service in configuration.nix (idempotent)
func nixosDisableService(L *lua.LState) int {
	params := L.CheckTable(1)

	serviceName := getStringField(L, params, "service", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("service name is required"))
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

	// Check if service is already disabled
	disablePattern := fmt.Sprintf(`services\.%s\.enable\s*=\s*false`, regexp.QuoteMeta(serviceName))
	if matched, _ := regexp.MatchString(disablePattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("service %s already disabled (idempotent)", serviceName)))
		return 2
	}

	// Check if service exists and is enabled
	enablePattern := fmt.Sprintf(`services\.%s\.enable\s*=\s*true`, regexp.QuoteMeta(serviceName))
	re := regexp.MustCompile(enablePattern)
	if re.MatchString(configStr) {
		// Change true to false
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf("services.%s.enable = false", serviceName))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("service %s disabled", serviceName)))
		return 2
	}

	// Service doesn't exist in config, nothing to disable
	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("service %s not in configuration (idempotent)", serviceName)))
	return 2
}

// ============================================================================
// System Options Functions
// ============================================================================

// nixosSetHostname sets the system hostname in configuration.nix (idempotent)
func nixosSetHostname(L *lua.LState) int {
	params := L.CheckTable(1)

	hostname := getStringField(L, params, "hostname", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if hostname == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("hostname is required"))
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

	// Check if hostname is already set
	currentPattern := fmt.Sprintf(`networking\.hostName\s*=\s*"%s"`, regexp.QuoteMeta(hostname))
	if matched, _ := regexp.MatchString(currentPattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("hostname already set to %s (idempotent)", hostname)))
		return 2
	}

	// Check if hostname exists and update it
	hostnamePattern := `networking\.hostName\s*=\s*"[^"]*"`
	re := regexp.MustCompile(hostnamePattern)
	if re.MatchString(configStr) {
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf(`networking.hostName = "%s"`, hostname))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("hostname set to %s", hostname)))
		return 2
	}

	// Add new hostname line
	hostnameConfig := fmt.Sprintf("  networking.hostName = \"%s\";\n", hostname)
	newContent := addLineToConfig(configStr, hostnameConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("hostname set to %s", hostname)))
	return 2
}

// nixosSetTimezone sets the system timezone in configuration.nix (idempotent)
func nixosSetTimezone(L *lua.LState) int {
	params := L.CheckTable(1)

	timezone := getStringField(L, params, "timezone", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if timezone == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("timezone is required"))
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

	// Check if timezone is already set
	currentPattern := fmt.Sprintf(`time\.timeZone\s*=\s*"%s"`, regexp.QuoteMeta(timezone))
	if matched, _ := regexp.MatchString(currentPattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("timezone already set to %s (idempotent)", timezone)))
		return 2
	}

	// Check if timezone exists and update it
	timezonePattern := `time\.timeZone\s*=\s*"[^"]*"`
	re := regexp.MustCompile(timezonePattern)
	if re.MatchString(configStr) {
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf(`time.timeZone = "%s"`, timezone))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("timezone set to %s", timezone)))
		return 2
	}

	// Add new timezone line
	timezoneConfig := fmt.Sprintf("  time.timeZone = \"%s\";\n", timezone)
	newContent := addLineToConfig(configStr, timezoneConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("timezone set to %s", timezone)))
	return 2
}

// nixosSetLocale sets the system locale in configuration.nix (idempotent)
func nixosSetLocale(L *lua.LState) int {
	params := L.CheckTable(1)

	locale := getStringField(L, params, "locale", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if locale == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("locale is required"))
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

	// Check if locale is already set
	currentPattern := fmt.Sprintf(`i18n\.defaultLocale\s*=\s*"%s"`, regexp.QuoteMeta(locale))
	if matched, _ := regexp.MatchString(currentPattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("locale already set to %s (idempotent)", locale)))
		return 2
	}

	// Check if locale exists and update it
	localePattern := `i18n\.defaultLocale\s*=\s*"[^"]*"`
	re := regexp.MustCompile(localePattern)
	if re.MatchString(configStr) {
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf(`i18n.defaultLocale = "%s"`, locale))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("locale set to %s", locale)))
		return 2
	}

	// Add new locale line
	localeConfig := fmt.Sprintf("  i18n.defaultLocale = \"%s\";\n", locale)
	newContent := addLineToConfig(configStr, localeConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("locale set to %s", locale)))
	return 2
}

// nixosEnableFirewall enables or disables the firewall in configuration.nix (idempotent)
func nixosEnableFirewall(L *lua.LState) int {
	params := L.CheckTable(1)

	enable := getBoolField(L, params, "enable", true)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)
	enableStr := "true"
	if !enable {
		enableStr = "false"
	}

	// Check if firewall is already set correctly
	currentPattern := fmt.Sprintf(`networking\.firewall\.enable\s*=\s*%s`, enableStr)
	if matched, _ := regexp.MatchString(currentPattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("firewall already set to %s (idempotent)", enableStr)))
		return 2
	}

	// Check if firewall setting exists and update it
	firewallPattern := `networking\.firewall\.enable\s*=\s*(true|false)`
	re := regexp.MustCompile(firewallPattern)
	if re.MatchString(configStr) {
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf("networking.firewall.enable = %s", enableStr))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("firewall set to %s", enableStr)))
		return 2
	}

	// Add new firewall line
	firewallConfig := fmt.Sprintf("  networking.firewall.enable = %s;\n", enableStr)
	newContent := addLineToConfig(configStr, firewallConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("firewall set to %s", enableStr)))
	return 2
}

// nixosAddFirewallPort adds a port to the firewall in configuration.nix (idempotent)
func nixosAddFirewallPort(L *lua.LState) int {
	params := L.CheckTable(1)

	port := getIntField(L, params, "port", 0)
	protocol := getStringField(L, params, "protocol", "tcp") // tcp or udp
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if port == 0 {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("port is required"))
		return 2
	}

	if protocol != "tcp" && protocol != "udp" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("protocol must be 'tcp' or 'udp'"))
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
	portField := fmt.Sprintf("allowed%sPorts", strings.ToUpper(protocol[:1])+protocol[1:])

	// Check if port is already in the list
	portPattern := fmt.Sprintf(`networking\.firewall\.%s\s*=\s*\[[^\]]*\b%d\b[^\]]*\]`, portField, port)
	if matched, _ := regexp.MatchString(portPattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("port %d/%s already allowed (idempotent)", port, protocol)))
		return 2
	}

	// Check if the ports list exists
	portsPattern := fmt.Sprintf(`(networking\.firewall\.%s\s*=\s*\[)([^\]]*)(])`, portField)
	re := regexp.MustCompile(portsPattern)
	if re.MatchString(configStr) {
		// Add port to existing list
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf("${1}${2} %d${3}", port))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("port %d/%s added to firewall", port, protocol)))
		return 2
	}

	// Create new ports list
	portsConfig := fmt.Sprintf("  networking.firewall.%s = [ %d ];\n", portField, port)
	newContent := addLineToConfig(configStr, portsConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("port %d/%s added to firewall", port, protocol)))
	return 2
}

// ============================================================================
// Nix Store Management Functions
// ============================================================================

// nixosCollectGarbage runs nix-collect-garbage to free disk space
func nixosCollectGarbage(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())

	deleteOlderThan := getStringField(L, params, "delete_older_than", "") // e.g., "30d"
	useSudo := getBoolField(L, params, "use_sudo", true)

	// Build command
	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nix-collect-garbage")
	if deleteOlderThan != "" {
		cmdArgs = append(cmdArgs, "-d", "--delete-older-than", deleteOlderThan)
	} else {
		cmdArgs = append(cmdArgs, "-d")
	}

	// Execute command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("garbage collection failed: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("garbage collection completed successfully")))
	return 2
}

// nixosOptimizeStore runs nix-store --optimize to deduplicate store files
func nixosOptimizeStore(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())

	useSudo := getBoolField(L, params, "use_sudo", true)

	// Build command
	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nix-store", "--optimize")

	// Execute command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("store optimization failed: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("store optimization completed successfully"))
	return 2
}

// ============================================================================
// Channel Management Functions
// ============================================================================

// nixosUpdateChannels runs nix-channel --update to update all channels
func nixosUpdateChannels(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())

	useSudo := getBoolField(L, params, "use_sudo", true)

	// Build command
	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nix-channel", "--update")

	// Execute command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("channel update failed: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("channels updated successfully"))
	return 2
}

// nixosListChannels runs nix-channel --list to list all channels
func nixosListChannels(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())

	useSudo := getBoolField(L, params, "use_sudo", false)

	// Build command
	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nix-channel", "--list")

	// Execute command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to list channels: %v", err)))
		return 2
	}

	L.Push(lua.LString(string(output)))
	L.Push(lua.LNil)
	return 2
}

// ============================================================================
// Boot Configuration Functions
// ============================================================================

// nixosSetBootloader sets the bootloader in configuration.nix (idempotent)
func nixosSetBootloader(L *lua.LState) int {
	params := L.CheckTable(1)

	bootloader := getStringField(L, params, "bootloader", "") // "systemd-boot" or "grub"
	device := getStringField(L, params, "device", "")         // for grub, e.g., "/dev/sda"
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if bootloader == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("bootloader is required (systemd-boot or grub)"))
		return 2
	}

	if bootloader != "systemd-boot" && bootloader != "grub" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("bootloader must be 'systemd-boot' or 'grub'"))
		return 2
	}

	if bootloader == "grub" && device == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("device is required for grub bootloader"))
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

	var bootloaderConfig string
	if bootloader == "systemd-boot" {
		bootloaderConfig = "  boot.loader.systemd-boot.enable = true;\n  boot.loader.efi.canTouchEfiVariables = true;\n"

		// Check if already set
		if strings.Contains(configStr, "boot.loader.systemd-boot.enable = true") {
			L.Push(lua.LBool(true))
			L.Push(lua.LString("systemd-boot already configured (idempotent)"))
			return 2
		}
	} else if bootloader == "grub" {
		bootloaderConfig = fmt.Sprintf("  boot.loader.grub.enable = true;\n  boot.loader.grub.device = \"%s\";\n", device)

		// Check if already set
		grubPattern := fmt.Sprintf(`boot\.loader\.grub\.device\s*=\s*"%s"`, regexp.QuoteMeta(device))
		if matched, _ := regexp.MatchString(grubPattern, configStr); matched {
			L.Push(lua.LBool(true))
			L.Push(lua.LString(fmt.Sprintf("grub already configured with device %s (idempotent)", device)))
			return 2
		}
	}

	// Remove any existing bootloader config
	configStr = removeBootloaderConfig(configStr)

	// Add new bootloader config
	newContent := addLineToConfig(configStr, bootloaderConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("bootloader set to %s", bootloader)))
	return 2
}

// ============================================================================
// Import Management Functions
// ============================================================================

// nixosAddImport adds an import to the imports list in configuration.nix (idempotent)
func nixosAddImport(L *lua.LState) int {
	params := L.CheckTable(1)

	importPath := getStringField(L, params, "import", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if importPath == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("import path is required"))
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

	// Check if import already exists
	importPattern := regexp.QuoteMeta(importPath)
	if matched, _ := regexp.MatchString(importPattern, configStr); matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("import %s already present (idempotent)", importPath)))
		return 2
	}

	// Check if imports list exists
	importsPattern := `(imports\s*=\s*\[)([^\]]*)(];)`
	re := regexp.MustCompile(importsPattern)
	if re.MatchString(configStr) {
		// Add to existing imports list
		newContent := re.ReplaceAllString(configStr, fmt.Sprintf("${1}${2}    %s\n  ${3}", importPath))
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
			return 2
		}
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("import %s added", importPath)))
		return 2
	}

	// Create imports list
	importsConfig := fmt.Sprintf("  imports = [\n    %s\n  ];\n", importPath)
	newContent := addLineToConfig(configStr, importsConfig)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("import %s added", importPath)))
	return 2
}

// nixosRemoveImport removes an import from the imports list in configuration.nix (idempotent)
func nixosRemoveImport(L *lua.LState) int {
	params := L.CheckTable(1)

	importPath := getStringField(L, params, "import", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if importPath == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("import path is required"))
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

	// Check if import exists
	importPattern := regexp.QuoteMeta(importPath)
	if matched, _ := regexp.MatchString(importPattern, configStr); !matched {
		L.Push(lua.LBool(true))
		L.Push(lua.LString(fmt.Sprintf("import %s not present (idempotent)", importPath)))
		return 2
	}

	// Remove the import line
	linePattern := fmt.Sprintf(`\s*%s\s*`, regexp.QuoteMeta(importPath))
	re := regexp.MustCompile(linePattern)
	newContent := re.ReplaceAllString(configStr, "")

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("import %s removed", importPath)))
	return 2
}

// ============================================================================
// Additional Helper Functions
// ============================================================================

// addLineToConfig adds a configuration line to the appropriate place in the config
func addLineToConfig(content, line string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder
	inserted := false

	for i, l := range lines {
		result.WriteString(l)
		result.WriteString("\n")

		// Insert before the last closing brace
		if !inserted && i > 0 && strings.TrimSpace(l) == "}" && i == len(lines)-2 {
			result.WriteString(line)
			inserted = true
		}
	}

	if !inserted {
		// Fallback: append near the end
		return content + "\n" + line
	}

	return result.String()
}

// removeBootloaderConfig removes existing bootloader configuration
func removeBootloaderConfig(content string) string {
	// Remove systemd-boot config
	content = regexp.MustCompile(`\s*boot\.loader\.systemd-boot\.enable\s*=\s*(true|false);\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*boot\.loader\.efi\.canTouchEfiVariables\s*=\s*(true|false);\s*`).ReplaceAllString(content, "")

	// Remove grub config
	content = regexp.MustCompile(`\s*boot\.loader\.grub\.enable\s*=\s*(true|false);\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*boot\.loader\.grub\.device\s*=\s*"[^"]*";\s*`).ReplaceAllString(content, "")

	return content
}

// getIntField gets an integer field from a Lua table
func getIntField(L *lua.LState, table *lua.LTable, key string, defaultValue int) int {
	value := L.GetField(table, key)
	if num, ok := value.(lua.LNumber); ok {
		return int(num)
	}
	return defaultValue
}
