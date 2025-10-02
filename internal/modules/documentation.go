package modules

// ModuleDoc represents documentation for a module
type ModuleDoc struct {
	Name        string
	Description string
	Functions   []FunctionDoc
}

// FunctionDoc represents documentation for a function
type FunctionDoc struct {
	Name        string
	Description string
	Example     string
	Parameters  string
	Returns     string
}

// GetAllModuleDocs returns documentation for all available modules
func GetAllModuleDocs() []ModuleDoc {
	return []ModuleDoc{
		// Core Modules
		{
			Name:        "pkg",
			Description: "Package management for multiple Linux distributions",
			Functions: []FunctionDoc{
				{
					Name:        "pkg.install",
					Description: "Install one or more packages",
					Parameters:  "{packages = {...}, target = 'agent_name'}",
					Example: `pkg.install({
    packages = {"nginx", "curl"},
    target = "web-server"
})`,
				},
				{
					Name:        "pkg.remove",
					Description: "Remove one or more packages",
					Parameters:  "{packages = {...}, target = 'agent_name'}",
					Example: `pkg.remove({
    packages = {"apache2"},
    target = "web-server"
})`,
				},
				{
					Name:        "pkg.update",
					Description: "Update package cache",
					Parameters:  "{target = 'agent_name'}",
					Example: `pkg.update({
    target = "web-server"
})`,
				},
				{
					Name:        "pkg.upgrade",
					Description: "Upgrade all packages",
					Parameters:  "{target = 'agent_name'}",
					Example: `pkg.upgrade({
    target = "web-server"
})`,
				},
				{
					Name:        "pkg.is_installed",
					Description: "Check if a package is installed",
					Parameters:  "{package = 'name', target = 'agent_name'}",
					Example: `local installed = pkg.is_installed({
    package = "nginx",
    target = "web-server"
})`,
				},
			},
		},
		{
			Name:        "systemd",
			Description: "Systemd service management",
			Functions: []FunctionDoc{
				{
					Name:        "systemd.enable",
					Description: "Enable a service to start on boot",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `systemd.enable({
    service = "nginx",
    target = "web-server"
})`,
				},
				{
					Name:        "systemd.disable",
					Description: "Disable a service from starting on boot",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `systemd.disable({
    service = "apache2",
    target = "web-server"
})`,
				},
				{
					Name:        "systemd.start",
					Description: "Start a service",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `systemd.start({
    service = "nginx",
    target = "web-server"
})`,
				},
				{
					Name:        "systemd.stop",
					Description: "Stop a service",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `systemd.stop({
    service = "nginx",
    target = "web-server"
})`,
				},
				{
					Name:        "systemd.restart",
					Description: "Restart a service",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `systemd.restart({
    service = "nginx",
    target = "web-server"
})`,
				},
				{
					Name:        "systemd.reload",
					Description: "Reload a service configuration",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `systemd.reload({
    service = "nginx",
    target = "web-server"
})`,
				},
				{
					Name:        "systemd.is_active",
					Description: "Check if a service is active",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `local active = systemd.is_active({
    service = "nginx",
    target = "web-server"
})`,
				},
				{
					Name:        "systemd.is_enabled",
					Description: "Check if a service is enabled",
					Parameters:  "{service = 'name', target = 'agent_name'}",
					Example: `local enabled = systemd.is_enabled({
    service = "nginx",
    target = "web-server"
})`,
				},
			},
		},
		{
			Name:        "user",
			Description: "Linux user management",
			Functions: []FunctionDoc{
				{
					Name:        "user.create",
					Description: "Create a new user",
					Parameters:  "options (table): username (required), password, uid, gid, home, shell, groups, comment, system, create_home, no_create_home, expiry",
					Returns:     "boolean (success), string (message)",
					Example: `-- New syntax (recommended)
local success, msg = user.create({
    username = "deploy",
    password = "securepassword",
    home = "/home/deploy",
    shell = "/bin/bash",
    groups = {"docker", "sudo"},
    comment = "Deployment User",
    create_home = true
})
if success then
    print("User created: " .. msg)
else
    print("Failed to create user: " .. msg)
end

-- Old syntax (still supported)
local success, msg = user.create("deploy", {
    password = "securepassword",
    home = "/home/deploy",
    shell = "/bin/bash"
})`,
				},
				{
					Name:        "user.delete",
					Description: "Delete a user",
					Parameters:  "username (string), remove_home (boolean, optional)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.delete("olduser", true)`,
				},
				{
					Name:        "user.exists",
					Description: "Check if a user exists",
					Parameters:  "username (string)",
					Returns:     "boolean (exists), string (message)",
					Example: `local exists, msg = user.exists("deploy")
if exists then
    print("User exists")
end`,
				},
				{
					Name:        "user.modify",
					Description: "Modify user properties",
					Parameters:  "username (string), options (table): uid, gid, home, move_home, shell, groups, comment, expiry, lock, unlock",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.modify("deploy", {
    shell = "/bin/zsh",
    groups = "docker,sudo,www-data"
})`,
				},
				{
					Name:        "user.add_to_group",
					Description: "Add user to a group",
					Parameters:  "username (string), group (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.add_to_group("deploy", "docker")`,
				},
				{
					Name:        "user.remove_from_group",
					Description: "Remove user from a group",
					Parameters:  "username (string), group (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.remove_from_group("deploy", "sudo")`,
				},
				{
					Name:        "user.get_info",
					Description: "Get user information",
					Parameters:  "username (string)",
					Returns:     "table (info: username, uid, gid, home, shell, comment) or nil, string (error message)",
					Example: `local info, err = user.get_info("deploy")
if info then
    print("UID: " .. info.uid)
    print("Home: " .. info.home)
    print("Shell: " .. info.shell)
end`,
				},
				{
					Name:        "user.set_password",
					Description: "Set user password",
					Parameters:  "username (string), password (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.set_password("deploy", "newsecurepassword")`,
				},
				{
					Name:        "user.list",
					Description: "List all users",
					Parameters:  "system_only (boolean, optional)",
					Returns:     "table (array of user info tables) or nil, string (error message)",
					Example: `local users, err = user.list(false)
for i, u in ipairs(users) do
    print(u.username .. " - UID: " .. u.uid)
end`,
				},
				{
					Name:        "user.lock",
					Description: "Lock a user account",
					Parameters:  "username (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.lock("tempuser")`,
				},
				{
					Name:        "user.unlock",
					Description: "Unlock a user account",
					Parameters:  "username (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.unlock("tempuser")`,
				},
				{
					Name:        "user.is_locked",
					Description: "Check if a user account is locked",
					Parameters:  "username (string)",
					Returns:     "boolean (locked), string (message)",
					Example: `local locked, msg = user.is_locked("deploy")`,
				},
				{
					Name:        "user.expire_password",
					Description: "Expire a user's password",
					Parameters:  "username (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.expire_password("deploy")`,
				},
				{
					Name:        "user.set_shell",
					Description: "Set the user's shell",
					Parameters:  "username (string), shell (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.set_shell("deploy", "/bin/zsh")`,
				},
				{
					Name:        "user.set_home",
					Description: "Set the user's home directory",
					Parameters:  "username (string), home_dir (string), move_files (boolean, optional)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.set_home("deploy", "/opt/deploy", true)`,
				},
				{
					Name:        "user.get_uid",
					Description: "Get the UID of a user",
					Parameters:  "username (string)",
					Returns:     "number (uid) or nil, string (error message)",
					Example: `local uid, err = user.get_uid("deploy")
if uid then print("UID: " .. uid) end`,
				},
				{
					Name:        "user.get_gid",
					Description: "Get the primary GID of a user",
					Parameters:  "username (string)",
					Returns:     "number (gid) or nil, string (error message)",
					Example: `local gid, err = user.get_gid("deploy")
if gid then print("GID: " .. gid) end`,
				},
				{
					Name:        "user.get_groups",
					Description: "Get all groups a user belongs to",
					Parameters:  "username (string)",
					Returns:     "table (array of group names) or nil, string (error message)",
					Example: `local groups, err = user.get_groups("deploy")
for i, g in ipairs(groups) do
    print(g)
end`,
				},
				{
					Name:        "user.set_primary_group",
					Description: "Set the user's primary group",
					Parameters:  "username (string), group (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.set_primary_group("deploy", "developers")`,
				},
				{
					Name:        "user.get_home",
					Description: "Get the user's home directory",
					Parameters:  "username (string)",
					Returns:     "string (home directory) or nil, string (error message)",
					Example: `local home, err = user.get_home("deploy")
if home then print("Home: " .. home) end`,
				},
				{
					Name:        "user.get_shell",
					Description: "Get the user's shell",
					Parameters:  "username (string)",
					Returns:     "string (shell) or nil, string (error message)",
					Example: `local shell, err = user.get_shell("deploy")
if shell then print("Shell: " .. shell) end`,
				},
				{
					Name:        "user.get_comment",
					Description: "Get the user's comment/GECOS field",
					Parameters:  "username (string)",
					Returns:     "string (comment) or nil, string (error message)",
					Example: `local comment, err = user.get_comment("deploy")
if comment then print("Comment: " .. comment) end`,
				},
				{
					Name:        "user.set_comment",
					Description: "Set the user's comment/GECOS field",
					Parameters:  "username (string), comment (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.set_comment("deploy", "Deployment User - Updated")`,
				},
				{
					Name:        "user.is_system_user",
					Description: "Check if a user is a system user (UID < 1000)",
					Parameters:  "username (string)",
					Returns:     "boolean (is_system) or nil, string (error message)",
					Example: `local is_system, err = user.is_system_user("deploy")
if is_system ~= nil then print("System user: " .. tostring(is_system)) end`,
				},
				{
					Name:        "user.get_current",
					Description: "Get the current user",
					Parameters:  "none",
					Returns:     "table (user info: username, uid, gid, home, shell) or nil, string (error message)",
					Example: `local current, err = user.get_current()
if current then
    print("Current user: " .. current.username)
end`,
				},
				{
					Name:        "user.group_create",
					Description: "Create a new group",
					Parameters:  "groupname (string), options (table, optional): gid, system",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.group_create("developers", {gid = "5000"})`,
				},
				{
					Name:        "user.group_delete",
					Description: "Delete a group",
					Parameters:  "groupname (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.group_delete("oldgroup")`,
				},
				{
					Name:        "user.group_exists",
					Description: "Check if a group exists",
					Parameters:  "groupname (string)",
					Returns:     "boolean (exists), string (message)",
					Example: `local exists, msg = user.group_exists("developers")`,
				},
				{
					Name:        "user.group_get_info",
					Description: "Get group information",
					Parameters:  "groupname (string)",
					Returns:     "table (info: name, gid, members) or nil, string (error message)",
					Example: `local info, err = user.group_get_info("developers")
if info then
    print("GID: " .. info.gid)
    print("Members: " .. table.concat(info.members, ", "))
end`,
				},
				{
					Name:        "user.group_list",
					Description: "List all groups",
					Parameters:  "none",
					Returns:     "table (array of group info tables) or nil, string (error message)",
					Example: `local groups, err = user.group_list()
for i, g in ipairs(groups) do
    print(g.name .. " - GID: " .. g.gid)
end`,
				},
				{
					Name:        "user.group_get_gid",
					Description: "Get the GID of a group",
					Parameters:  "groupname (string)",
					Returns:     "number (gid) or nil, string (error message)",
					Example: `local gid, err = user.group_get_gid("developers")
if gid then print("GID: " .. gid) end`,
				},
				{
					Name:        "user.group_members",
					Description: "Get all members of a group",
					Parameters:  "groupname (string)",
					Returns:     "table (array of usernames) or nil, string (error message)",
					Example: `local members, err = user.group_members("developers")
for i, m in ipairs(members) do
    print(m)
end`,
				},
				{
					Name:        "user.group_add_member",
					Description: "Add a member to a group",
					Parameters:  "groupname (string), username (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.group_add_member("developers", "deploy")`,
				},
				{
					Name:        "user.group_remove_member",
					Description: "Remove a member from a group",
					Parameters:  "groupname (string), username (string)",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.group_remove_member("developers", "deploy")`,
				},
				{
					Name:        "user.set_expiry",
					Description: "Set when an account expires",
					Parameters:  "username (string), expiry (string) - Format: YYYY-MM-DD",
					Returns:     "boolean (success), string (message)",
					Example: `local success, msg = user.set_expiry("tempuser", "2025-12-31")`,
				},
				{
					Name:        "user.get_last_login",
					Description: "Get the last login time for a user",
					Parameters:  "username (string)",
					Returns:     "string (timestamp) or nil, string (error message)",
					Example: `local last_login, err = user.get_last_login("deploy")
if last_login then print("Last login: " .. last_login) end`,
				},
				{
					Name:        "user.get_failed_logins",
					Description: "Get failed login attempts for a user",
					Parameters:  "username (string)",
					Returns:     "number (count) or nil, string (error message)",
					Example: `local failed, err = user.get_failed_logins("deploy")
if failed then print("Failed logins: " .. failed) end`,
				},
				{
					Name:        "user.validate_username",
					Description: "Validate if a username follows Linux conventions",
					Parameters:  "username (string)",
					Returns:     "boolean (valid), string (message)",
					Example: `local valid, msg = user.validate_username("deploy-user")
if not valid then print("Invalid: " .. msg) end`,
				},
				{
					Name:        "user.is_root",
					Description: "Check if the current user is root",
					Parameters:  "none",
					Returns:     "boolean (is_root) or nil, string (error message)",
					Example: `local is_root, err = user.is_root()
if is_root then print("Running as root") end`,
				},
				{
					Name:        "user.run_as",
					Description: "Run a command as a different user",
					Parameters:  "username (string), command (string)",
					Returns:     "string (output) or nil, string (error message)",
					Example: `local output, err = user.run_as("deploy", "whoami")
if output then print("Result: " .. output) end`,
				},
			},
		},
		{
			Name:        "ssh",
			Description: "SSH key and configuration management",
			Functions: []FunctionDoc{
				{
					Name:        "ssh.generate_keypair",
					Description: "Generate SSH key pair",
					Parameters:  "{path = 'path', type = 'rsa|ed25519', bits = num, comment = 'text', passphrase = 'text', target = 'agent_name'}",
					Example: `ssh.generate_keypair({
    path = "/home/deploy/.ssh/id_ed25519",
    type = "ed25519",
    comment = "deploy@server",
    target = "web-server"
})`,
				},
				{
					Name:        "ssh.add_authorized_key",
					Description: "Add SSH authorized key",
					Parameters:  "{user = 'name', key = 'pubkey', target = 'agent_name'}",
					Example: `ssh.add_authorized_key({
    user = "deploy",
    key = "ssh-ed25519 AAAAC3... user@host",
    target = "web-server"
})`,
				},
				{
					Name:        "ssh.remove_authorized_key",
					Description: "Remove SSH authorized key",
					Parameters:  "{user = 'name', key = 'pubkey', target = 'agent_name'}",
					Example: `ssh.remove_authorized_key({
    user = "deploy",
    key = "ssh-ed25519 AAAAC3... user@host",
    target = "web-server"
})`,
				},
				{
					Name:        "ssh.set_config",
					Description: "Configure SSH client settings",
					Parameters:  "{user = 'name', host = 'hostname', config = {...}, target = 'agent_name'}",
					Example: `ssh.set_config({
    user = "deploy",
    host = "github.com",
    config = {
        HostName = "github.com",
        User = "git",
        IdentityFile = "~/.ssh/id_ed25519"
    },
    target = "web-server"
})`,
				},
			},
		},
		{
			Name:        "file",
			Description: "File and directory operations",
			Functions: []FunctionDoc{
				{
					Name:        "file.copy",
					Description: "Copy files from master to agent",
					Parameters:  "{src = 'path', dest = 'path', mode = '0644', owner = 'user', group = 'group', target = 'agent_name'}",
					Example: `file.copy({
    src = "./config.conf",
    dest = "/etc/app/config.conf",
    mode = "0644",
    owner = "root",
    group = "root",
    target = "web-server"
})`,
				},
				{
					Name:        "file.fetch",
					Description: "Download files from agent to master",
					Parameters:  "{src = 'path', dest = 'path', target = 'agent_name'}",
					Example: `file.fetch({
    src = "/var/log/app.log",
    dest = "./logs/app.log",
    target = "web-server"
})`,
				},
				{
					Name:        "file.template",
					Description: "Render and copy Go template to agent",
					Parameters:  "{src = 'path', dest = 'path', vars = {...}, mode = '0644', owner = 'user', group = 'group', target = 'agent_name'}",
					Example: `file.template({
    src = "./nginx.conf.tmpl",
    dest = "/etc/nginx/nginx.conf",
    vars = {
        port = 8080,
        server_name = "example.com"
    },
    mode = "0644",
    target = "web-server"
})`,
				},
				{
					Name:        "file.set_attributes",
					Description: "Set file attributes (permissions, owner, group)",
					Parameters:  "{path = 'path', mode = '0644', owner = 'user', group = 'group', state = 'file|directory|link', target = 'agent_name'}",
					Example: `file.set_attributes({
    path = "/etc/app/config.conf",
    mode = "0600",
    owner = "app",
    group = "app",
    target = "web-server"
})`,
				},
				{
					Name:        "file.line_in_file",
					Description: "Ensure a line exists in a file",
					Parameters:  "{path = 'path', line = 'text', regex = 'pattern', state = 'present|absent', target = 'agent_name'}",
					Example: `file.line_in_file({
    path = "/etc/hosts",
    line = "192.168.1.10 myserver",
    state = "present",
    target = "web-server"
})`,
				},
				{
					Name:        "file.block_in_file",
					Description: "Insert/update/remove a block of text in a file",
					Parameters:  "{path = 'path', block = 'text', marker = 'text', state = 'present|absent', target = 'agent_name'}",
					Example: `file.block_in_file({
    path = "/etc/nginx/nginx.conf",
    block = [[
server {
    listen 80;
    server_name example.com;
}]],
    marker = "# {mark} ANSIBLE MANAGED BLOCK",
    state = "present",
    target = "web-server"
})`,
				},
				{
					Name:        "file.replace",
					Description: "Replace text in file using regex",
					Parameters:  "{path = 'path', pattern = 'regex', replacement = 'text', target = 'agent_name'}",
					Example: `file.replace({
    path = "/etc/app/config.conf",
    pattern = "port = %d+",
    replacement = "port = 8080",
    target = "web-server"
})`,
				},
				{
					Name:        "file.unarchive",
					Description: "Extract archive files",
					Parameters:  "{src = 'path', dest = 'path', creates = 'path', target = 'agent_name'}",
					Example: `file.unarchive({
    src = "./app.tar.gz",
    dest = "/opt/app",
    creates = "/opt/app/bin/app",
    target = "web-server"
})`,
				},
				{
					Name:        "file.stat",
					Description: "Get file information",
					Parameters:  "{path = 'path', target = 'agent_name'}",
					Example: `local info = file.stat({
    path = "/etc/app/config.conf",
    target = "web-server"
})
print("Size: " .. info.size)`,
				},
			},
		},
		{
			Name:        "http",
			Description: "HTTP client operations",
			Functions: []FunctionDoc{
				{
					Name:        "http.get",
					Description: "Perform HTTP GET request",
					Parameters:  "{url = 'url', headers = {...}}",
					Example: `local response = http.get({
    url = "https://api.example.com/data",
    headers = {
        ["Authorization"] = "Bearer token"
    }
})`,
				},
				{
					Name:        "http.post",
					Description: "Perform HTTP POST request",
					Parameters:  "{url = 'url', body = 'data', headers = {...}}",
					Example: `http.post({
    url = "https://api.example.com/data",
    body = json.encode({name = "test"}),
    headers = {
        ["Content-Type"] = "application/json"
    }
})`,
				},
			},
		},
		{
			Name:        "cmd",
			Description: "Execute shell commands",
			Functions: []FunctionDoc{
				{
					Name:        "cmd.run",
					Description: "Execute a shell command",
					Parameters:  "{command = 'cmd', cwd = 'path', env = {...}}",
					Example: `local result = cmd.run({
    command = "ls -la",
    cwd = "/tmp"
})`,
				},
			},
		},
		{
			Name:        "json",
			Description: "JSON encoding and decoding",
			Functions: []FunctionDoc{
				{
					Name:        "json.encode",
					Description: "Encode Lua table to JSON",
					Parameters:  "table",
					Example: `local jsonStr = json.encode({
    name = "test",
    value = 123
})`,
				},
				{
					Name:        "json.decode",
					Description: "Decode JSON to Lua table",
					Parameters:  "string",
					Example: `local data = json.decode('{"name":"test"}')
print(data.name)`,
				},
			},
		},
		{
			Name:        "yaml",
			Description: "YAML encoding and decoding",
			Functions: []FunctionDoc{
				{
					Name:        "yaml.encode",
					Description: "Encode Lua table to YAML",
					Parameters:  "table",
					Example: `local yamlStr = yaml.encode({
    name = "test",
    items = {1, 2, 3}
})`,
				},
				{
					Name:        "yaml.decode",
					Description: "Decode YAML to Lua table",
					Parameters:  "string",
					Example: `local data = yaml.decode([[
name: test
items:
  - 1
  - 2
]])`,
				},
			},
		},
		{
			Name:        "log",
			Description: "Logging functions",
			Functions: []FunctionDoc{
				{
					Name:        "log.info",
					Description: "Log info message",
					Parameters:  "message",
					Example:     `log.info("Starting deployment")`,
				},
				{
					Name:        "log.warn",
					Description: "Log warning message",
					Parameters:  "message",
					Example:     `log.warn("Service is slow")`,
				},
				{
					Name:        "log.error",
					Description: "Log error message",
					Parameters:  "message",
					Example:     `log.error("Deployment failed")`,
				},
				{
					Name:        "log.debug",
					Description: "Log debug message",
					Parameters:  "message",
					Example:     `log.debug("Variable value: " .. value)`,
				},
			},
		},
		{
			Name:        "crypto",
			Description: "Cryptographic operations",
			Functions: []FunctionDoc{
				{
					Name:        "crypto.hash",
					Description: "Generate hash (md5, sha1, sha256, sha512)",
					Parameters:  "{data = 'text', algorithm = 'sha256'}",
					Example: `local hash = crypto.hash({
    data = "password",
    algorithm = "sha256"
})`,
				},
				{
					Name:        "crypto.encrypt",
					Description: "Encrypt data with AES",
					Parameters:  "{data = 'text', key = 'key'}",
					Example: `local encrypted = crypto.encrypt({
    data = "secret",
    key = "encryption-key"
})`,
				},
				{
					Name:        "crypto.decrypt",
					Description: "Decrypt AES encrypted data",
					Parameters:  "{data = 'encrypted', key = 'key'}",
					Example: `local decrypted = crypto.decrypt({
    data = encrypted,
    key = "encryption-key"
})`,
				},
			},
		},
		{
			Name:        "database",
			Description: "Database operations (PostgreSQL, MySQL, SQLite)",
			Functions: []FunctionDoc{
				{
					Name:        "database.connect",
					Description: "Connect to a database",
					Parameters:  "{driver = 'postgres|mysql|sqlite', dsn = 'connection_string'}",
					Example: `local db = database.connect({
    driver = "postgres",
    dsn = "host=localhost user=admin password=secret dbname=mydb"
})`,
				},
				{
					Name:        "database.query",
					Description: "Execute a query",
					Parameters:  "{db = connection, query = 'sql'}",
					Example: `local rows = database.query({
    db = db,
    query = "SELECT * FROM users"
})`,
				},
				{
					Name:        "database.exec",
					Description: "Execute a statement",
					Parameters:  "{db = connection, query = 'sql'}",
					Example: `database.exec({
    db = db,
    query = "INSERT INTO users (name) VALUES ('John')"
})`,
				},
			},
		},
		{
			Name:        "terraform",
			Description: "Terraform operations",
			Functions: []FunctionDoc{
				{
					Name:        "terraform.init",
					Description: "Initialize Terraform",
					Parameters:  "{dir = 'path'}",
					Example: `terraform.init({
    dir = "./infrastructure"
})`,
				},
				{
					Name:        "terraform.plan",
					Description: "Create Terraform plan",
					Parameters:  "{dir = 'path', vars = {...}}",
					Example: `terraform.plan({
    dir = "./infrastructure",
    vars = {
        region = "us-east-1",
        instance_type = "t2.micro"
    }
})`,
				},
				{
					Name:        "terraform.apply",
					Description: "Apply Terraform changes",
					Parameters:  "{dir = 'path', vars = {...}, auto_approve = bool}",
					Example: `terraform.apply({
    dir = "./infrastructure",
    auto_approve = true
})`,
				},
				{
					Name:        "terraform.destroy",
					Description: "Destroy Terraform resources",
					Parameters:  "{dir = 'path', auto_approve = bool}",
					Example: `terraform.destroy({
    dir = "./infrastructure",
    auto_approve = true
})`,
				},
			},
		},
		{
			Name:        "pulumi",
			Description: "Pulumi operations",
			Functions: []FunctionDoc{
				{
					Name:        "pulumi.up",
					Description: "Deploy Pulumi stack",
					Parameters:  "{dir = 'path', stack = 'name', config = {...}}",
					Example: `pulumi.up({
    dir = "./infrastructure",
    stack = "production"
})`,
				},
				{
					Name:        "pulumi.preview",
					Description: "Preview Pulumi changes",
					Parameters:  "{dir = 'path', stack = 'name'}",
					Example: `pulumi.preview({
    dir = "./infrastructure",
    stack = "production"
})`,
				},
				{
					Name:        "pulumi.destroy",
					Description: "Destroy Pulumi stack",
					Parameters:  "{dir = 'path', stack = 'name'}",
					Example: `pulumi.destroy({
    dir = "./infrastructure",
    stack = "production"
})`,
				},
			},
		},
		{
			Name:        "aws",
			Description: "AWS operations",
			Functions: []FunctionDoc{
				{
					Name:        "aws.ec2_list",
					Description: "List EC2 instances",
					Parameters:  "{region = 'region', filters = {...}}",
					Example: `local instances = aws.ec2_list({
    region = "us-east-1",
    filters = {
        ["tag:Environment"] = "production"
    }
})`,
				},
				{
					Name:        "aws.s3_upload",
					Description: "Upload file to S3",
					Parameters:  "{bucket = 'name', key = 'path', file = 'localpath', region = 'region'}",
					Example: `aws.s3_upload({
    bucket = "my-bucket",
    key = "backup/data.tar.gz",
    file = "./data.tar.gz",
    region = "us-east-1"
})`,
				},
			},
		},
		{
			Name:        "azure",
			Description: "Azure operations",
			Functions: []FunctionDoc{
				{
					Name:        "azure.vm_list",
					Description: "List Azure VMs",
					Parameters:  "{subscription = 'id', resource_group = 'name'}",
					Example: `local vms = azure.vm_list({
    subscription = "sub-id",
    resource_group = "my-rg"
})`,
				},
			},
		},
		{
			Name:        "gcp",
			Description: "Google Cloud Platform operations",
			Functions: []FunctionDoc{
				{
					Name:        "gcp.compute_list",
					Description: "List GCP Compute instances",
					Parameters:  "{project = 'id', zone = 'zone'}",
					Example: `local instances = gcp.compute_list({
    project = "my-project",
    zone = "us-central1-a"
})`,
				},
			},
		},
		{
			Name:        "docker",
			Description: "Docker operations",
			Functions: []FunctionDoc{
				{
					Name:        "docker.build",
					Description: "Build Docker image",
					Parameters:  "{path = 'path', tag = 'name:tag', target = 'agent_name'}",
					Example: `docker.build({
    path = "./app",
    tag = "myapp:latest",
    target = "build-server"
})`,
				},
				{
					Name:        "docker.run",
					Description: "Run Docker container",
					Parameters:  "{image = 'name', name = 'container', ports = {...}, volumes = {...}, env = {...}, target = 'agent_name'}",
					Example: `docker.run({
    image = "nginx:latest",
    name = "web",
    ports = {"80:80"},
    target = "web-server"
})`,
				},
			},
		},
		{
			Name:        "kubernetes",
			Description: "Kubernetes operations",
			Functions: []FunctionDoc{
				{
					Name:        "kubernetes.apply",
					Description: "Apply Kubernetes manifest",
					Parameters:  "{file = 'path', namespace = 'name', target = 'agent_name'}",
					Example: `kubernetes.apply({
    file = "./deployment.yaml",
    namespace = "production",
    target = "k8s-master"
})`,
				},
				{
					Name:        "kubernetes.delete",
					Description: "Delete Kubernetes resources",
					Parameters:  "{file = 'path', namespace = 'name', target = 'agent_name'}",
					Example: `kubernetes.delete({
    file = "./deployment.yaml",
    namespace = "production",
    target = "k8s-master"
})`,
				},
			},
		},
		{
			Name:        "slack",
			Description: "Slack notifications",
			Functions: []FunctionDoc{
				{
					Name:        "slack.send",
					Description: "Send Slack message",
					Parameters:  "{webhook = 'url', channel = 'name', message = 'text', username = 'name'}",
					Example: `slack.send({
    webhook = "https://hooks.slack.com/...",
    channel = "#deployments",
    message = "Deployment completed successfully"
})`,
				},
			},
		},
		{
			Name:        "goroutine",
			Description: "Concurrent execution with goroutines",
			Functions: []FunctionDoc{
				{
					Name:        "goroutine.spawn",
					Description: "Spawn a new goroutine",
					Parameters:  "function",
					Example: `goroutine.spawn(function()
    log.info("Running in parallel")
end)`,
				},
				{
					Name:        "goroutine.wait",
					Description: "Wait for all goroutines to complete",
					Parameters:  "none",
					Example:     `goroutine.wait()`,
				},
			},
		},
		{
			Name:        "incus",
			Description: "Container and VM management with Incus (LXC/LXD fork)",
			Functions: []FunctionDoc{
				{
					Name:        "incus.instance",
					Description: "Create an instance builder for containers or VMs",
					Parameters:  "{name = 'string', image = 'string', type = 'container|virtual-machine', profiles = {...}}",
					Example: `local web = incus.instance({
    name = "web-01",
    image = "ubuntu:22.04",
    profiles = {"default", "web-server"}
})
web:create():start():wait_running()`,
				},
				{
					Name:        "instance:create",
					Description: "Create the instance",
					Parameters:  "none",
					Example:     `instance:create()`,
				},
				{
					Name:        "instance:start",
					Description: "Start the instance",
					Parameters:  "none",
					Example:     `instance:start()`,
				},
				{
					Name:        "instance:stop",
					Description: "Stop the instance (optionally force)",
					Parameters:  "force (boolean, optional)",
					Example: `instance:stop()       -- graceful
instance:stop(true)  -- force`,
				},
				{
					Name:        "instance:restart",
					Description: "Restart the instance",
					Parameters:  "none",
					Example:     `instance:restart()`,
				},
				{
					Name:        "instance:delete",
					Description: "Delete the instance",
					Parameters:  "none",
					Example:     `instance:delete()`,
				},
				{
					Name:        "instance:wait_running",
					Description: "Wait for instance to be running",
					Parameters:  "timeout (number, optional)",
					Example:     `instance:wait_running(120)  -- wait up to 120 seconds`,
				},
				{
					Name:        "instance:exec",
					Description: "Execute command in the instance",
					Parameters:  "command (string), options (table, optional)",
					Example: `instance:exec("apt update && apt upgrade -y")
instance:exec("whoami", {user = "ubuntu"})`,
				},
				{
					Name:        "instance:file_push",
					Description: "Upload file to the instance",
					Parameters:  "local_path (string), remote_path (string)",
					Example:     `instance:file_push("./config.yaml", "/etc/app/config.yaml")`,
				},
				{
					Name:        "instance:file_pull",
					Description: "Download file from the instance",
					Parameters:  "remote_path (string), local_path (string)",
					Example:     `instance:file_pull("/var/log/app.log", "./logs/app.log")`,
				},
				{
					Name:        "instance:set_config",
					Description: "Set instance configuration",
					Parameters:  "{[key] = value, ...}",
					Example: `instance:set_config({
    ["limits.cpu"] = "4",
    ["limits.memory"] = "8GB"
})`,
				},
				{
					Name:        "instance:add_device",
					Description: "Add a device to the instance",
					Parameters:  "name (string), config (table)",
					Example: `instance:add_device("eth0", {
    type = "nic",
    nictype = "bridged",
    parent = "br0"
})`,
				},
				{
					Name:        "instance:delegate_to",
					Description: "Execute instance operations on a specific agent",
					Parameters:  "agent_name (string)",
					Example:     `instance:delegate_to("incus-host-01")`,
				},
				{
					Name:        "incus.network",
					Description: "Create a network builder",
					Parameters:  "{name = 'string', type = 'bridge|macvlan|...'}",
					Example: `incus.network({
    name = "web-dmz",
    type = "bridge"
}):set_config({
    ["ipv4.address"] = "10.0.0.1/24",
    ["ipv4.nat"] = "true"
}):create()`,
				},
				{
					Name:        "network:create",
					Description: "Create the network",
					Parameters:  "none",
					Example:     `network:create()`,
				},
				{
					Name:        "network:attach",
					Description: "Attach network to an instance",
					Parameters:  "instance_name (string)",
					Example:     `network:attach("web-01")`,
				},
				{
					Name:        "network:detach",
					Description: "Detach network from an instance",
					Parameters:  "instance_name (string)",
					Example:     `network:detach("web-01")`,
				},
				{
					Name:        "incus.profile",
					Description: "Create a profile builder",
					Parameters:  "{name = 'string', description = 'string'}",
					Example: `incus.profile({
    name = "web-server",
    description = "Web server profile"
}):set_config({
    ["limits.cpu"] = "4"
}):create()`,
				},
				{
					Name:        "profile:apply",
					Description: "Apply profile to an instance",
					Parameters:  "instance_name (string)",
					Example:     `profile:apply("web-01")`,
				},
				{
					Name:        "incus.storage",
					Description: "Create a storage pool builder",
					Parameters:  "{name = 'string', driver = 'zfs|btrfs|dir|lvm'}",
					Example: `incus.storage({
    name = "ssd-pool",
    driver = "zfs"
}):set_config({
    ["size"] = "100GB"
}):create()`,
				},
				{
					Name:        "incus.snapshot",
					Description: "Create a snapshot builder",
					Parameters:  "{instance = 'string', name = 'string', stateful = boolean}",
					Example: `incus.snapshot({
    instance = "web-01",
    name = "backup-20241002",
    stateful = true
}):create()`,
				},
				{
					Name:        "snapshot:restore",
					Description: "Restore a snapshot",
					Parameters:  "none",
					Example:     `snapshot:restore()`,
				},
				{
					Name:        "incus.list",
					Description: "List Incus resources",
					Parameters:  "type (string): 'instances', 'networks', 'profiles', 'storage-pools'",
					Example:     `local instances = incus.list("instances")`,
				},
				{
					Name:        "incus.info",
					Description: "Get information about a resource",
					Parameters:  "type (string), name (string)",
					Example:     `local info = incus.info("instance", "web-01")`,
				},
				{
					Name:        "incus.exec",
					Description: "Execute command in an instance (standalone)",
					Parameters:  "instance (string), command (string), options (table, optional)",
					Example:     `incus.exec("web-01", "systemctl status nginx")`,
				},
				{
					Name:        "incus.delete",
					Description: "Delete a resource (standalone)",
					Parameters:  "type (string), name (string)",
					Example:     `incus.delete("instance", "old-container")`,
				},
			},
		},
	}
}
