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

	// Advanced configuration functions with table support
	L.SetField(nixosModule, "configure_service", L.NewFunction(nixosConfigureService))
	L.SetField(nixosModule, "configure_user", L.NewFunction(nixosConfigureUser))
	L.SetField(nixosModule, "configure_networking", L.NewFunction(nixosConfigureNetworking))
	L.SetField(nixosModule, "configure_system", L.NewFunction(nixosConfigureSystem))
	L.SetField(nixosModule, "configure_environment", L.NewFunction(nixosConfigureEnvironment))

	// Systemd management
	L.SetField(nixosModule, "create_systemd_service", L.NewFunction(nixosCreateSystemdService))
	L.SetField(nixosModule, "create_systemd_timer", L.NewFunction(nixosCreateSystemdTimer))
	L.SetField(nixosModule, "create_systemd_mount", L.NewFunction(nixosCreateSystemdMount))

	// Storage management
	L.SetField(nixosModule, "configure_zfs", L.NewFunction(nixosConfigureZFS))
	L.SetField(nixosModule, "configure_filesystem", L.NewFunction(nixosConfigureFilesystem))

	// Container management
	L.SetField(nixosModule, "configure_container", L.NewFunction(nixosConfigureContainer))
	L.SetField(nixosModule, "configure_docker", L.NewFunction(nixosConfigureDocker))

	// Advanced networking
	L.SetField(nixosModule, "configure_vlan", L.NewFunction(nixosConfigureVLAN))
	L.SetField(nixosModule, "configure_bridge", L.NewFunction(nixosConfigureBridge))
	L.SetField(nixosModule, "configure_vpn", L.NewFunction(nixosConfigureVPN))

	// Security and secrets
	L.SetField(nixosModule, "configure_security", L.NewFunction(nixosConfigureSecurity))

	// Generation management
	L.SetField(nixosModule, "list_generations", L.NewFunction(nixosListGenerations))
	L.SetField(nixosModule, "rollback", L.NewFunction(nixosRollback))
	L.SetField(nixosModule, "switch_generation", L.NewFunction(nixosSwitchGeneration))
	L.SetField(nixosModule, "delete_generations", L.NewFunction(nixosDeleteGenerations))

	// Declarative containers and VMs
	L.SetField(nixosModule, "test_config", L.NewFunction(nixosTestConfig))

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
// Advanced Configuration Functions
// ============================================================================

// nixosConfigureService configures a service with all options in a single call
// Accepts a configuration table with all service settings
func nixosConfigureService(L *lua.LState) int {
	params := L.CheckTable(1)

	serviceName := getStringField(L, params, "service", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")
	enable := getBoolField(L, params, "enable", true)

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

	// Build service configuration
	var serviceConfig strings.Builder
	serviceConfig.WriteString(fmt.Sprintf("  services.%s = {\n", serviceName))
	serviceConfig.WriteString(fmt.Sprintf("    enable = %t;\n", enable))

	// Process all service-specific settings from the settings table
	settingsTable := L.GetField(params, "settings")
	if tbl, ok := settingsTable.(*lua.LTable); ok {
		tbl.ForEach(func(key, value lua.LValue) {
			keyStr := lua.LVAsString(key)
			if keyStr == "" {
				return
			}

			// Generate Nix config based on value type
			nixValue := luaValueToNixString(L, value)
			serviceConfig.WriteString(fmt.Sprintf("    %s = %s;\n", keyStr, nixValue))
		})
	}

	serviceConfig.WriteString("  };\n")

	// Check if service block already exists
	servicePattern := fmt.Sprintf(`services\.%s\s*=\s*\{[^}]*\};`, regexp.QuoteMeta(serviceName))
	re := regexp.MustCompile(servicePattern)

	var newContent string
	if re.MatchString(configStr) {
		// Replace existing service config
		newContent = re.ReplaceAllString(configStr, strings.TrimSpace(serviceConfig.String()))
	} else {
		// Add new service config
		newContent = addLineToConfig(configStr, serviceConfig.String())
	}

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("service %s configured successfully", serviceName)))
	return 2
}

// nixosConfigureUser configures a complete user with all settings including SSH, packages, etc
func nixosConfigureUser(L *lua.LState) int {
	params := L.CheckTable(1)

	username := getStringField(L, params, "username", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if username == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("username is required"))
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

	// Build complete user configuration
	var userConfig strings.Builder
	userConfig.WriteString(fmt.Sprintf("  users.users.%s = {\n", username))

	// Basic user settings
	if isNormalUser := getBoolField(L, params, "is_normal_user", true); isNormalUser {
		userConfig.WriteString("    isNormalUser = true;\n")
	} else {
		userConfig.WriteString("    isSystemUser = true;\n")
	}

	if uid := getIntField(L, params, "uid", 0); uid > 0 {
		userConfig.WriteString(fmt.Sprintf("    uid = %d;\n", uid))
	}

	if home := getStringField(L, params, "home", ""); home != "" {
		userConfig.WriteString(fmt.Sprintf("    home = \"%s\";\n", home))
	} else {
		userConfig.WriteString(fmt.Sprintf("    home = \"/home/%s\";\n", username))
	}

	if description := getStringField(L, params, "description", ""); description != "" {
		userConfig.WriteString(fmt.Sprintf("    description = \"%s\";\n", description))
	}

	if shell := getStringField(L, params, "shell", ""); shell != "" {
		userConfig.WriteString(fmt.Sprintf("    shell = %s;\n", shell))
	}

	// Groups
	var groups []string
	groupsTable := L.GetField(params, "groups")
	if tbl, ok := groupsTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				groups = append(groups, string(str))
			}
		})
	}
	if len(groups) > 0 {
		userConfig.WriteString("    extraGroups = [ ")
		for _, group := range groups {
			userConfig.WriteString(fmt.Sprintf("\"%s\" ", group))
		}
		userConfig.WriteString("];\n")
	}

	// SSH keys
	var sshKeys []string
	sshKeysTable := L.GetField(params, "ssh_keys")
	if tbl, ok := sshKeysTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				sshKeys = append(sshKeys, string(str))
			}
		})
	}
	// Single ssh_key parameter
	if singleKey := getStringField(L, params, "ssh_key", ""); singleKey != "" {
		sshKeys = append(sshKeys, singleKey)
	}

	if len(sshKeys) > 0 {
		userConfig.WriteString("    openssh.authorizedKeys.keys = [\n")
		for _, key := range sshKeys {
			userConfig.WriteString(fmt.Sprintf("      \"%s\"\n", strings.TrimSpace(key)))
		}
		userConfig.WriteString("    ];\n")
	}

	// User-specific packages
	var packages []string
	packagesTable := L.GetField(params, "packages")
	if tbl, ok := packagesTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				packages = append(packages, string(str))
			}
		})
	}
	if len(packages) > 0 {
		userConfig.WriteString("    packages = with pkgs; [\n")
		for _, pkg := range packages {
			userConfig.WriteString(fmt.Sprintf("      %s\n", pkg))
		}
		userConfig.WriteString("    ];\n")
	}

	// Initial password (hashed)
	if hashedPassword := getStringField(L, params, "hashed_password", ""); hashedPassword != "" {
		userConfig.WriteString(fmt.Sprintf("    hashedPassword = \"%s\";\n", hashedPassword))
	}

	// Initial password file
	if passwordFile := getStringField(L, params, "password_file", ""); passwordFile != "" {
		userConfig.WriteString(fmt.Sprintf("    passwordFile = \"%s\";\n", passwordFile))
	}

	// createHome
	if createHome := getBoolField(L, params, "create_home", true); !createHome {
		userConfig.WriteString("    createHome = false;\n")
	}

	// Additional custom settings
	settingsTable := L.GetField(params, "settings")
	if tbl, ok := settingsTable.(*lua.LTable); ok {
		tbl.ForEach(func(key, value lua.LValue) {
			keyStr := lua.LVAsString(key)
			if keyStr == "" {
				return
			}
			nixValue := luaValueToNixString(L, value)
			userConfig.WriteString(fmt.Sprintf("    %s = %s;\n", keyStr, nixValue))
		})
	}

	userConfig.WriteString("  };\n")

	// Check if user already exists
	userPattern := fmt.Sprintf(`users\.users\.%s\s*=\s*\{[^}]*\};`, regexp.QuoteMeta(username))
	re := regexp.MustCompile(userPattern)

	var newContent string
	if re.MatchString(configStr) {
		// Replace existing user config
		newContent = re.ReplaceAllString(configStr, strings.TrimSpace(userConfig.String()))
	} else {
		// Add new user config
		newContent = insertUserIntoConfig(configStr, userConfig.String())
	}

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("user %s configured successfully", username)))
	return 2
}

// nixosConfigureNetworking configures complete networking settings
func nixosConfigureNetworking(L *lua.LState) int {
	params := L.CheckTable(1)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	// Build networking configuration
	var networkingConfig strings.Builder
	networkingConfig.WriteString("  networking = {\n")

	// Hostname
	if hostname := getStringField(L, params, "hostname", ""); hostname != "" {
		networkingConfig.WriteString(fmt.Sprintf("    hostName = \"%s\";\n", hostname))
	}

	// Domain
	if domain := getStringField(L, params, "domain", ""); domain != "" {
		networkingConfig.WriteString(fmt.Sprintf("    domain = \"%s\";\n", domain))
	}

	// Enable NetworkManager
	if enableNetworkManager := getBoolField(L, params, "enable_network_manager", false); enableNetworkManager {
		networkingConfig.WriteString("    networkmanager.enable = true;\n")
	}

	// DHCP
	if useDHCP := getBoolField(L, params, "use_dhcp", false); useDHCP {
		networkingConfig.WriteString("    useDHCP = true;\n")
	}

	// Firewall settings
	firewallTable := L.GetField(params, "firewall")
	if tbl, ok := firewallTable.(*lua.LTable); ok {
		networkingConfig.WriteString("    firewall = {\n")

		// Enable firewall
		if enable := getBoolFieldFromTable(L, tbl, "enable", true); !enable {
			networkingConfig.WriteString("      enable = false;\n")
		} else {
			networkingConfig.WriteString("      enable = true;\n")
		}

		// TCP ports
		var tcpPorts []int
		tcpPortsTable := L.GetField(tbl, "tcp_ports")
		if portsTbl, ok := tcpPortsTable.(*lua.LTable); ok {
			portsTbl.ForEach(func(_, value lua.LValue) {
				if num, ok := value.(lua.LNumber); ok {
					tcpPorts = append(tcpPorts, int(num))
				}
			})
		}
		if len(tcpPorts) > 0 {
			networkingConfig.WriteString("      allowedTCPPorts = [ ")
			for _, port := range tcpPorts {
				networkingConfig.WriteString(fmt.Sprintf("%d ", port))
			}
			networkingConfig.WriteString("];\n")
		}

		// UDP ports
		var udpPorts []int
		udpPortsTable := L.GetField(tbl, "udp_ports")
		if portsTbl, ok := udpPortsTable.(*lua.LTable); ok {
			portsTbl.ForEach(func(_, value lua.LValue) {
				if num, ok := value.(lua.LNumber); ok {
					udpPorts = append(udpPorts, int(num))
				}
			})
		}
		if len(udpPorts) > 0 {
			networkingConfig.WriteString("      allowedUDPPorts = [ ")
			for _, port := range udpPorts {
				networkingConfig.WriteString(fmt.Sprintf("%d ", port))
			}
			networkingConfig.WriteString("];\n")
		}

		// Allowed TCP port ranges
		tcpRangesTable := L.GetField(tbl, "tcp_port_ranges")
		if rangesTbl, ok := tcpRangesTable.(*lua.LTable); ok {
			hasRanges := false
			var ranges []string
			rangesTbl.ForEach(func(_, value lua.LValue) {
				if rangeTbl, ok := value.(*lua.LTable); ok {
					from := getIntField(L, rangeTbl, "from", 0)
					to := getIntField(L, rangeTbl, "to", 0)
					if from > 0 && to > 0 {
						ranges = append(ranges, fmt.Sprintf("{ from = %d; to = %d; }", from, to))
						hasRanges = true
					}
				}
			})
			if hasRanges {
				networkingConfig.WriteString("      allowedTCPPortRanges = [\n")
				for _, r := range ranges {
					networkingConfig.WriteString(fmt.Sprintf("        %s\n", r))
				}
				networkingConfig.WriteString("      ];\n")
			}
		}

		networkingConfig.WriteString("    };\n")
	}

	// Custom networking settings
	settingsTable := L.GetField(params, "settings")
	if tbl, ok := settingsTable.(*lua.LTable); ok {
		tbl.ForEach(func(key, value lua.LValue) {
			keyStr := lua.LVAsString(key)
			if keyStr == "" {
				return
			}
			nixValue := luaValueToNixString(L, value)
			networkingConfig.WriteString(fmt.Sprintf("    %s = %s;\n", keyStr, nixValue))
		})
	}

	networkingConfig.WriteString("  };\n")

	// Replace or add networking block
	networkingPattern := `networking\s*=\s*\{[^}]*\};`
	re := regexp.MustCompile(networkingPattern)

	var newContent string
	if re.MatchString(configStr) {
		newContent = re.ReplaceAllString(configStr, strings.TrimSpace(networkingConfig.String()))
	} else {
		newContent = addLineToConfig(configStr, networkingConfig.String())
	}

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("networking configured successfully"))
	return 2
}

// nixosConfigureSystem configures complete system settings
func nixosConfigureSystem(L *lua.LState) int {
	params := L.CheckTable(1)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)
	modified := false

	// Timezone
	if timezone := getStringField(L, params, "timezone", ""); timezone != "" {
		timezoneConfig := fmt.Sprintf("time.timeZone = \"%s\"", timezone)
		timezonePattern := `time\.timeZone\s*=\s*"[^"]*"`
		re := regexp.MustCompile(timezonePattern)
		if re.MatchString(configStr) {
			configStr = re.ReplaceAllString(configStr, timezoneConfig)
		} else {
			configStr = addLineToConfig(configStr, fmt.Sprintf("  %s;\n", timezoneConfig))
		}
		modified = true
	}

	// Locale
	if locale := getStringField(L, params, "locale", ""); locale != "" {
		localeConfig := fmt.Sprintf("i18n.defaultLocale = \"%s\"", locale)
		localePattern := `i18n\.defaultLocale\s*=\s*"[^"]*"`
		re := regexp.MustCompile(localePattern)
		if re.MatchString(configStr) {
			configStr = re.ReplaceAllString(configStr, localeConfig)
		} else {
			configStr = addLineToConfig(configStr, fmt.Sprintf("  %s;\n", localeConfig))
		}
		modified = true
	}

	// Console settings
	consoleTable := L.GetField(params, "console")
	if tbl, ok := consoleTable.(*lua.LTable); ok {
		var consoleConfig strings.Builder
		consoleConfig.WriteString("  console = {\n")

		if keyMap := getStringFieldFromTable(L, tbl, "keymap", ""); keyMap != "" {
			consoleConfig.WriteString(fmt.Sprintf("    keyMap = \"%s\";\n", keyMap))
		}
		if font := getStringFieldFromTable(L, tbl, "font", ""); font != "" {
			consoleConfig.WriteString(fmt.Sprintf("    font = \"%s\";\n", font))
		}

		consoleConfig.WriteString("  };\n")

		consolePattern := `console\s*=\s*\{[^}]*\};`
		re := regexp.MustCompile(consolePattern)
		if re.MatchString(configStr) {
			configStr = re.ReplaceAllString(configStr, strings.TrimSpace(consoleConfig.String()))
		} else {
			configStr = addLineToConfig(configStr, consoleConfig.String())
		}
		modified = true
	}

	// Boot settings
	bootTable := L.GetField(params, "boot")
	if tbl, ok := bootTable.(*lua.LTable); ok {
		// Remove existing boot config
		configStr = removeBootloaderConfig(configStr)

		var bootConfig strings.Builder
		bootConfig.WriteString("  boot = {\n")

		// Bootloader
		bootloaderTable := L.GetField(tbl, "loader")
		if loaderTbl, ok := bootloaderTable.(*lua.LTable); ok {
			bootConfig.WriteString("    loader = {\n")

			// systemd-boot
			systemdBootTable := L.GetField(loaderTbl, "systemd_boot")
			if sdBootTbl, ok := systemdBootTable.(*lua.LTable); ok {
				bootConfig.WriteString("      systemd-boot = {\n")
				if enable := getBoolFieldFromTable(L, sdBootTbl, "enable", false); enable {
					bootConfig.WriteString("        enable = true;\n")
				}
				bootConfig.WriteString("      };\n")
			}

			// GRUB
			grubTable := L.GetField(loaderTbl, "grub")
			if grubTbl, ok := grubTable.(*lua.LTable); ok {
				bootConfig.WriteString("      grub = {\n")
				if enable := getBoolFieldFromTable(L, grubTbl, "enable", false); enable {
					bootConfig.WriteString("        enable = true;\n")
				}
				if device := getStringFieldFromTable(L, grubTbl, "device", ""); device != "" {
					bootConfig.WriteString(fmt.Sprintf("        device = \"%s\";\n", device))
				}
				bootConfig.WriteString("      };\n")
			}

			// EFI
			efiTable := L.GetField(loaderTbl, "efi")
			if efiTbl, ok := efiTable.(*lua.LTable); ok {
				bootConfig.WriteString("      efi = {\n")
				if canTouch := getBoolFieldFromTable(L, efiTbl, "can_touch_efi_variables", false); canTouch {
					bootConfig.WriteString("        canTouchEfiVariables = true;\n")
				}
				bootConfig.WriteString("      };\n")
			}

			bootConfig.WriteString("    };\n")
		}

		// Kernel params
		var kernelParams []string
		kernelParamsTable := L.GetField(tbl, "kernel_params")
		if paramsTbl, ok := kernelParamsTable.(*lua.LTable); ok {
			paramsTbl.ForEach(func(_, value lua.LValue) {
				if str, ok := value.(lua.LString); ok {
					kernelParams = append(kernelParams, string(str))
				}
			})
		}
		if len(kernelParams) > 0 {
			bootConfig.WriteString("    kernelParams = [ ")
			for _, param := range kernelParams {
				bootConfig.WriteString(fmt.Sprintf("\"%s\" ", param))
			}
			bootConfig.WriteString("];\n")
		}

		bootConfig.WriteString("  };\n")

		bootPattern := `boot\s*=\s*\{[^}]*\};`
		re := regexp.MustCompile(bootPattern)
		if re.MatchString(configStr) {
			configStr = re.ReplaceAllString(configStr, strings.TrimSpace(bootConfig.String()))
		} else {
			configStr = addLineToConfig(configStr, bootConfig.String())
		}
		modified = true
	}

	if !modified {
		L.Push(lua.LBool(true))
		L.Push(lua.LString("no system configuration changes specified"))
		return 2
	}

	// Write back to config
	if err := os.WriteFile(configPath, []byte(configStr), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("system configured successfully"))
	return 2
}

// nixosConfigureEnvironment configures environment settings including packages and variables
func nixosConfigureEnvironment(L *lua.LState) int {
	params := L.CheckTable(1)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	// Build environment configuration
	var envConfig strings.Builder
	envConfig.WriteString("  environment = {\n")

	// System packages
	var packages []string
	packagesTable := L.GetField(params, "packages")
	if tbl, ok := packagesTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				packages = append(packages, string(str))
			}
		})
	}
	if len(packages) > 0 {
		envConfig.WriteString("    systemPackages = with pkgs; [\n")
		for _, pkg := range packages {
			envConfig.WriteString(fmt.Sprintf("      %s\n", pkg))
		}
		envConfig.WriteString("    ];\n")
	}

	// Environment variables
	var envVars map[string]string
	envVarsTable := L.GetField(params, "variables")
	if tbl, ok := envVarsTable.(*lua.LTable); ok {
		envVars = make(map[string]string)
		tbl.ForEach(func(key, value lua.LValue) {
			keyStr := lua.LVAsString(key)
			valueStr := lua.LVAsString(value)
			if keyStr != "" && valueStr != "" {
				envVars[keyStr] = valueStr
			}
		})
	}
	if len(envVars) > 0 {
		envConfig.WriteString("    variables = {\n")
		for k, v := range envVars {
			envConfig.WriteString(fmt.Sprintf("      %s = \"%s\";\n", k, v))
		}
		envConfig.WriteString("    };\n")
	}

	// System PATH
	var pathsToLink []string
	pathsTable := L.GetField(params, "paths_to_link")
	if tbl, ok := pathsTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				pathsToLink = append(pathsToLink, string(str))
			}
		})
	}
	if len(pathsToLink) > 0 {
		envConfig.WriteString("    pathsToLink = [\n")
		for _, path := range pathsToLink {
			envConfig.WriteString(fmt.Sprintf("      \"%s\"\n", path))
		}
		envConfig.WriteString("    ];\n")
	}

	// Custom environment settings
	settingsTable := L.GetField(params, "settings")
	if tbl, ok := settingsTable.(*lua.LTable); ok {
		tbl.ForEach(func(key, value lua.LValue) {
			keyStr := lua.LVAsString(key)
			if keyStr == "" {
				return
			}
			nixValue := luaValueToNixString(L, value)
			envConfig.WriteString(fmt.Sprintf("    %s = %s;\n", keyStr, nixValue))
		})
	}

	envConfig.WriteString("  };\n")

	// Replace or add environment block
	environmentPattern := `environment\s*=\s*\{[^}]*\};`
	re := regexp.MustCompile(environmentPattern)

	var newContent string
	if re.MatchString(configStr) {
		newContent = re.ReplaceAllString(configStr, strings.TrimSpace(envConfig.String()))
	} else {
		newContent = addLineToConfig(configStr, envConfig.String())
	}

	// Write back to config
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("environment configured successfully"))
	return 2
}

// ============================================================================
// Systemd Management Functions
// ============================================================================

// nixosCreateSystemdService creates a custom systemd service
func nixosCreateSystemdService(L *lua.LState) int {
	params := L.CheckTable(1)

	serviceName := getStringField(L, params, "name", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if serviceName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("service name is required"))
		return 2
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	// Build systemd service configuration
	var serviceConfig strings.Builder
	serviceConfig.WriteString(fmt.Sprintf("  systemd.services.%s = {\n", serviceName))

	if description := getStringField(L, params, "description", ""); description != "" {
		serviceConfig.WriteString(fmt.Sprintf("    description = \"%s\";\n", description))
	}

	if wantedBy := getStringField(L, params, "wanted_by", ""); wantedBy != "" {
		serviceConfig.WriteString(fmt.Sprintf("    wantedBy = [ \"%s\" ];\n", wantedBy))
	}

	if after := getStringField(L, params, "after", ""); after != "" {
		serviceConfig.WriteString(fmt.Sprintf("    after = [ \"%s\" ];\n", after))
	}

	// Service section
	serviceSection := L.GetField(params, "service")
	if tbl, ok := serviceSection.(*lua.LTable); ok {
		serviceConfig.WriteString("    serviceConfig = {\n")

		if execStart := getStringFieldFromTable(L, tbl, "exec_start", ""); execStart != "" {
			serviceConfig.WriteString(fmt.Sprintf("      ExecStart = \"%s\";\n", execStart))
		}
		if restart := getStringFieldFromTable(L, tbl, "restart", ""); restart != "" {
			serviceConfig.WriteString(fmt.Sprintf("      Restart = \"%s\";\n", restart))
		}
		if user := getStringFieldFromTable(L, tbl, "user", ""); user != "" {
			serviceConfig.WriteString(fmt.Sprintf("      User = \"%s\";\n", user))
		}
		if group := getStringFieldFromTable(L, tbl, "group", ""); group != "" {
			serviceConfig.WriteString(fmt.Sprintf("      Group = \"%s\";\n", group))
		}
		if workingDir := getStringFieldFromTable(L, tbl, "working_directory", ""); workingDir != "" {
			serviceConfig.WriteString(fmt.Sprintf("      WorkingDirectory = \"%s\";\n", workingDir))
		}

		serviceConfig.WriteString("    };\n")
	}

	serviceConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, serviceConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("systemd service %s created", serviceName)))
	return 2
}

// nixosCreateSystemdTimer creates a systemd timer
func nixosCreateSystemdTimer(L *lua.LState) int {
	params := L.CheckTable(1)

	timerName := getStringField(L, params, "name", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if timerName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("timer name is required"))
		return 2
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var timerConfig strings.Builder
	timerConfig.WriteString(fmt.Sprintf("  systemd.timers.%s = {\n", timerName))
	timerConfig.WriteString("    wantedBy = [ \"timers.target\" ];\n")

	if onCalendar := getStringField(L, params, "on_calendar", ""); onCalendar != "" {
		timerConfig.WriteString("    timerConfig = {\n")
		timerConfig.WriteString(fmt.Sprintf("      OnCalendar = \"%s\";\n", onCalendar))
		if persistent := getBoolField(L, params, "persistent", false); persistent {
			timerConfig.WriteString("      Persistent = true;\n")
		}
		timerConfig.WriteString("    };\n")
	}

	timerConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, timerConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("systemd timer %s created", timerName)))
	return 2
}

// nixosCreateSystemdMount creates a systemd mount unit
func nixosCreateSystemdMount(L *lua.LState) int {
	params := L.CheckTable(1)

	mountName := getStringField(L, params, "name", "")
	what := getStringField(L, params, "what", "")
	where := getStringField(L, params, "where", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if mountName == "" || what == "" || where == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name, what, and where are required"))
		return 2
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var mountConfig strings.Builder
	mountConfig.WriteString(fmt.Sprintf("  systemd.mounts = [{\n"))
	mountConfig.WriteString(fmt.Sprintf("    what = \"%s\";\n", what))
	mountConfig.WriteString(fmt.Sprintf("    where = \"%s\";\n", where))
	if fsType := getStringField(L, params, "type", ""); fsType != "" {
		mountConfig.WriteString(fmt.Sprintf("    type = \"%s\";\n", fsType))
	}
	if options := getStringField(L, params, "options", ""); options != "" {
		mountConfig.WriteString(fmt.Sprintf("    options = \"%s\";\n", options))
	}
	mountConfig.WriteString("  }];\n")

	newContent := addLineToConfig(configStr, mountConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("systemd mount %s created", mountName)))
	return 2
}

// ============================================================================
// Storage Management Functions
// ============================================================================

// nixosConfigureZFS configures ZFS support and pools
func nixosConfigureZFS(L *lua.LState) int {
	params := L.CheckTable(1)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var zfsConfig strings.Builder
	zfsConfig.WriteString("  boot.supportedFilesystems = [ \"zfs\" ];\n")

	if hostId := getStringField(L, params, "host_id", ""); hostId != "" {
		zfsConfig.WriteString(fmt.Sprintf("  networking.hostId = \"%s\";\n", hostId))
	}

	zfsConfig.WriteString("  services.zfs = {\n")
	zfsConfig.WriteString("    autoScrub.enable = true;\n")
	if autoSnapshot := getBoolField(L, params, "auto_snapshot", false); autoSnapshot {
		zfsConfig.WriteString("    autoSnapshot.enable = true;\n")
	}
	zfsConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, zfsConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("ZFS configured successfully"))
	return 2
}

// nixosConfigureFilesystem configures advanced filesystem options
func nixosConfigureFilesystem(L *lua.LState) int {
	params := L.CheckTable(1)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var fsConfig strings.Builder

	// Filesystems support
	var supportedFS []string
	supportedFSTable := L.GetField(params, "supported")
	if tbl, ok := supportedFSTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				supportedFS = append(supportedFS, string(str))
			}
		})
	}

	if len(supportedFS) > 0 {
		fsConfig.WriteString("  boot.supportedFilesystems = [ ")
		for _, fs := range supportedFS {
			fsConfig.WriteString(fmt.Sprintf("\"%s\" ", fs))
		}
		fsConfig.WriteString("];\n")
	}

	newContent := addLineToConfig(configStr, fsConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("filesystem configuration updated"))
	return 2
}

// ============================================================================
// Container Management Functions
// ============================================================================

// nixosConfigureContainer creates a declarative NixOS container
func nixosConfigureContainer(L *lua.LState) int {
	params := L.CheckTable(1)

	containerName := getStringField(L, params, "name", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if containerName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("container name is required"))
		return 2
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var containerConfig strings.Builder
	containerConfig.WriteString(fmt.Sprintf("  containers.%s = {\n", containerName))
	containerConfig.WriteString("    autoStart = true;\n")

	if privateNetwork := getBoolField(L, params, "private_network", false); privateNetwork {
		containerConfig.WriteString("    privateNetwork = true;\n")
		if hostAddress := getStringField(L, params, "host_address", ""); hostAddress != "" {
			containerConfig.WriteString(fmt.Sprintf("    hostAddress = \"%s\";\n", hostAddress))
		}
		if localAddress := getStringField(L, params, "local_address", ""); localAddress != "" {
			containerConfig.WriteString(fmt.Sprintf("    localAddress = \"%s\";\n", localAddress))
		}
	}

	containerConfig.WriteString("    config = { config, pkgs, ... }: {\n")
	containerConfig.WriteString("      services.openssh.enable = true;\n")
	containerConfig.WriteString("    };\n")
	containerConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, containerConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("container %s configured", containerName)))
	return 2
}

// nixosConfigureDocker configures Docker with advanced settings
func nixosConfigureDocker(L *lua.LState) int {
	params := L.CheckTable(1)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var dockerConfig strings.Builder
	dockerConfig.WriteString("  virtualisation.docker = {\n")
	dockerConfig.WriteString("    enable = true;\n")

	if rootless := getBoolField(L, params, "rootless", false); rootless {
		dockerConfig.WriteString("    rootless = {\n")
		dockerConfig.WriteString("      enable = true;\n")
		dockerConfig.WriteString("      setSocketVariable = true;\n")
		dockerConfig.WriteString("    };\n")
	}

	if autoPrune := getBoolField(L, params, "auto_prune", false); autoPrune {
		dockerConfig.WriteString("    autoPrune.enable = true;\n")
	}

	dockerConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, dockerConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("Docker configured successfully"))
	return 2
}

// ============================================================================
// Advanced Networking Functions
// ============================================================================

// nixosConfigureVLAN configures a VLAN interface
func nixosConfigureVLAN(L *lua.LState) int {
	params := L.CheckTable(1)

	vlanName := getStringField(L, params, "name", "")
	id := getIntField(L, params, "id", 0)
	interface_ := getStringField(L, params, "interface", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if vlanName == "" || id == 0 || interface_ == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("name, id, and interface are required"))
		return 2
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var vlanConfig strings.Builder
	vlanConfig.WriteString(fmt.Sprintf("  networking.vlans.%s = {\n", vlanName))
	vlanConfig.WriteString(fmt.Sprintf("    id = %d;\n", id))
	vlanConfig.WriteString(fmt.Sprintf("    interface = \"%s\";\n", interface_))
	vlanConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, vlanConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("VLAN %s configured", vlanName)))
	return 2
}

// nixosConfigureBridge configures a network bridge
func nixosConfigureBridge(L *lua.LState) int {
	params := L.CheckTable(1)

	bridgeName := getStringField(L, params, "name", "")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	if bridgeName == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("bridge name is required"))
		return 2
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var bridgeConfig strings.Builder
	bridgeConfig.WriteString(fmt.Sprintf("  networking.bridges.%s = {\n", bridgeName))

	var interfaces []string
	interfacesTable := L.GetField(params, "interfaces")
	if tbl, ok := interfacesTable.(*lua.LTable); ok {
		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				interfaces = append(interfaces, string(str))
			}
		})
	}

	if len(interfaces) > 0 {
		bridgeConfig.WriteString("    interfaces = [ ")
		for _, iface := range interfaces {
			bridgeConfig.WriteString(fmt.Sprintf("\"%s\" ", iface))
		}
		bridgeConfig.WriteString("];\n")
	}

	bridgeConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, bridgeConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("bridge %s configured", bridgeName)))
	return 2
}

// nixosConfigureVPN configures VPN (WireGuard or OpenVPN)
func nixosConfigureVPN(L *lua.LState) int {
	params := L.CheckTable(1)

	vpnType := getStringField(L, params, "type", "wireguard")
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var vpnConfig strings.Builder

	if vpnType == "wireguard" {
		interfaceName := getStringField(L, params, "interface", "wg0")
		vpnConfig.WriteString(fmt.Sprintf("  networking.wireguard.interfaces.%s = {\n", interfaceName))

		if privateKeyFile := getStringField(L, params, "private_key_file", ""); privateKeyFile != "" {
			vpnConfig.WriteString(fmt.Sprintf("    privateKeyFile = \"%s\";\n", privateKeyFile))
		}

		if listenPort := getIntField(L, params, "listen_port", 0); listenPort > 0 {
			vpnConfig.WriteString(fmt.Sprintf("    listenPort = %d;\n", listenPort))
		}

		vpnConfig.WriteString("  };\n")
	}

	newContent := addLineToConfig(configStr, vpnConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("%s VPN configured", vpnType)))
	return 2
}

// ============================================================================
// Security Functions
// ============================================================================

// nixosConfigureSecurity configures security options
func nixosConfigureSecurity(L *lua.LState) int {
	params := L.CheckTable(1)
	configPath := getStringField(L, params, "config_path", "/etc/nixos/configuration.nix")

	content, err := os.ReadFile(configPath)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	configStr := string(content)

	var securityConfig strings.Builder
	securityConfig.WriteString("  security = {\n")

	if sudo := getBoolField(L, params, "sudo_wheel_needs_password", true); !sudo {
		securityConfig.WriteString("    sudo.wheelNeedsPassword = false;\n")
	}

	if polkit := getBoolField(L, params, "polkit_enable", false); polkit {
		securityConfig.WriteString("    polkit.enable = true;\n")
	}

	if apparmor := getBoolField(L, params, "apparmor_enable", false); apparmor {
		securityConfig.WriteString("    apparmor.enable = true;\n")
	}

	securityConfig.WriteString("  };\n")

	newContent := addLineToConfig(configStr, securityConfig.String())

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("security settings configured"))
	return 2
}

// ============================================================================
// Generation Management Functions
// ============================================================================

// nixosListGenerations lists all system generations
func nixosListGenerations(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())
	useSudo := getBoolField(L, params, "use_sudo", false)

	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nix-env", "--list-generations", "-p", "/nix/var/nix/profiles/system")

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to list generations: %v", err)))
		return 2
	}

	L.Push(lua.LString(string(output)))
	L.Push(lua.LNil)
	return 2
}

// nixosRollback rolls back to the previous generation
func nixosRollback(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())
	useSudo := getBoolField(L, params, "use_sudo", true)

	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nixos-rebuild", "switch", "--rollback")

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("rollback failed: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("successfully rolled back to previous generation"))
	return 2
}

// nixosSwitchGeneration switches to a specific generation
func nixosSwitchGeneration(L *lua.LState) int {
	params := L.CheckTable(1)

	generation := getIntField(L, params, "generation", 0)
	useSudo := getBoolField(L, params, "use_sudo", true)

	if generation == 0 {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("generation number is required"))
		return 2
	}

	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nix-env", "-p", "/nix/var/nix/profiles/system", "--switch-generation", fmt.Sprintf("%d", generation))

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to switch generation: %v\nOutput: %s", err, string(output))))
		return 2
	}

	// Now activate it
	if useSudo {
		cmdArgs = []string{"sudo", "/nix/var/nix/profiles/system/bin/switch-to-configuration", "switch"}
	} else {
		cmdArgs = []string{"/nix/var/nix/profiles/system/bin/switch-to-configuration", "switch"}
	}

	cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err = cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to activate generation: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("switched to generation %d", generation)))
	return 2
}

// nixosDeleteGenerations deletes old generations
func nixosDeleteGenerations(L *lua.LState) int {
	params := L.CheckTable(1)

	olderThan := getStringField(L, params, "older_than", "")
	useSudo := getBoolField(L, params, "use_sudo", true)

	if olderThan == "" {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("older_than parameter is required (e.g., '30d')"))
		return 2
	}

	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nix-env", "-p", "/nix/var/nix/profiles/system", "--delete-generations", olderThan)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to delete generations: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("deleted generations older than %s", olderThan)))
	return 2
}

// ============================================================================
// Testing Functions
// ============================================================================

// nixosTestConfig tests the configuration without applying it
func nixosTestConfig(L *lua.LState) int {
	params := L.OptTable(1, L.NewTable())
	useSudo := getBoolField(L, params, "use_sudo", true)

	var cmdArgs []string
	if useSudo {
		cmdArgs = append(cmdArgs, "sudo")
	}
	cmdArgs = append(cmdArgs, "nixos-rebuild", "test", "--fast")

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("configuration test failed: %v\nOutput: %s", err, string(output))))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LString("configuration test passed"))
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

// getBoolFieldFromTable gets a boolean field from a Lua table
func getBoolFieldFromTable(L *lua.LState, table *lua.LTable, key string, defaultValue bool) bool {
	value := L.GetField(table, key)
	if b, ok := value.(lua.LBool); ok {
		return bool(b)
	}
	return defaultValue
}

// getStringFieldFromTable gets a string field from a Lua table
func getStringFieldFromTable(L *lua.LState, table *lua.LTable, key string, defaultValue string) string {
	value := L.GetField(table, key)
	if str, ok := value.(lua.LString); ok {
		return string(str)
	}
	return defaultValue
}

// luaValueToNixString converts a Lua value to a Nix configuration string
func luaValueToNixString(L *lua.LState, value lua.LValue) string {
	// Check for nil first
	if value == lua.LNil {
		return "null"
	}

	switch v := value.(type) {
	case lua.LString:
		return fmt.Sprintf("\"%s\"", string(v))
	case lua.LNumber:
		return fmt.Sprintf("%v", v)
	case lua.LBool:
		if bool(v) {
			return "true"
		}
		return "false"
	case *lua.LTable:
		// Check if it's an array or a map
		isArray := true
		v.ForEach(func(key, _ lua.LValue) {
			if _, ok := key.(lua.LNumber); !ok {
				isArray = false
			}
		})

		if isArray {
			// It's an array
			var items []string
			v.ForEach(func(_, val lua.LValue) {
				items = append(items, luaValueToNixString(L, val))
			})
			return fmt.Sprintf("[ %s ]", strings.Join(items, " "))
		} else {
			// It's a map/object
			var attrs []string
			v.ForEach(func(key, val lua.LValue) {
				keyStr := lua.LVAsString(key)
				valStr := luaValueToNixString(L, val)
				attrs = append(attrs, fmt.Sprintf("%s = %s;", keyStr, valStr))
			})
			return fmt.Sprintf("{ %s }", strings.Join(attrs, " "))
		}
	default:
		return "null"
	}
}
