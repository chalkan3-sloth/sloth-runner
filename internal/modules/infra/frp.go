package infra

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
	lua "github.com/yuin/gopher-lua"
)

// FrpModule provides FRP (Fast Reverse Proxy) management functionality
type FrpModule struct {
	agentClient interface{}
}

// NewFrpModule creates a new FRP module instance
func NewFrpModule(agentClient interface{}) *FrpModule {
	return &FrpModule{
		agentClient: agentClient,
	}
}

// Register registers the FRP module and its functions in the Lua state
func (m *FrpModule) Register(L *lua.LState) {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"server":  m.createServer,
		"client":  m.createClient,
		"install": m.install,
	})

	L.SetGlobal("frp", mod)
}

// ============================================================================
// FRP Server Builder
// ============================================================================

// FrpServer represents an FRP server configuration builder
type FrpServer struct {
	L           *lua.LState
	agentClient interface{}
	name        string
	config      map[string]interface{}
	configPath  string
	target      string
	version     string
}

// createServer creates a new FRP server builder instance
func (m *FrpModule) createServer(L *lua.LState) int {
	name := "frps"
	if L.GetTop() >= 1 {
		name = L.CheckString(1)
	}

	server := &FrpServer{
		L:           L,
		agentClient: m.agentClient,
		name:        name,
		config:      make(map[string]interface{}),
		configPath:  "/etc/frp/frps.toml",
		version:     "latest",
	}

	// Set default configuration
	server.config["bindPort"] = 7000
	server.config["vhostHTTPPort"] = 80
	server.config["vhostHTTPSPort"] = 443

	ud := L.NewUserData()
	ud.Value = server
	L.SetMetatable(ud, L.GetTypeMetatable("frp_server"))
	L.Push(ud)
	return 1
}

// RegisterServerMetatable registers the FRP server metatable with Lua methods
func RegisterServerMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("frp_server")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"config":       serverConfig,
		"config_path":  serverConfigPath,
		"version":      serverVersion,
		"delegate_to":  serverDelegateTo,
		"save_config":  serverSaveConfig,
		"load_config":  serverLoadConfig,
		"start":        serverStart,
		"stop":         serverStop,
		"restart":      serverRestart,
		"status":       serverStatus,
		"install":      serverInstall,
		"enable":       serverEnable,
		"disable":      serverDisable,
	}))
}

// checkFrpServer extracts FrpServer from Lua userdata
func checkFrpServer(L *lua.LState, n int) *FrpServer {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*FrpServer); ok {
		return v
	}
	L.ArgError(n, "FrpServer expected")
	return nil
}

// serverConfig sets configuration from a Lua table (fluent method)
func serverConfig(L *lua.LState) int {
	server := checkFrpServer(L, 1)
	config := L.CheckTable(2)

	config.ForEach(func(k, v lua.LValue) {
		key := k.String()
		server.config[key] = luaValueToGo(v)
	})

	L.Push(L.Get(1)) // Return self for chaining
	L.Push(lua.LNil)
	return 2
}

// serverConfigPath sets the configuration file path (fluent method)
func serverConfigPath(L *lua.LState) int {
	server := checkFrpServer(L, 1)
	path := L.CheckString(2)
	server.configPath = path

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// serverVersion sets the FRP version to install (fluent method)
func serverVersion(L *lua.LState) int {
	server := checkFrpServer(L, 1)
	version := L.CheckString(2)
	server.version = version

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// serverDelegateTo sets the target agent for command execution (fluent method)
func serverDelegateTo(L *lua.LState) int {
	server := checkFrpServer(L, 1)
	target := L.CheckString(2)
	server.target = target

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// serverSaveConfig saves the configuration to a TOML file (action method)
func serverSaveConfig(L *lua.LState) int {
	server := checkFrpServer(L, 1)

	// Convert config map to TOML
	tomlData, err := toml.Marshal(server.config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to marshal TOML: %v", err)))
		return 2
	}

	// Create directory if needed
	dir := filepath.Dir(server.configPath)
	mkdirCmd := fmt.Sprintf("mkdir -p %s", dir)
	_, err = executeCommandWithExec(L, mkdirCmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create config directory: %v", err)))
		return 2
	}

	// Write config file
	writeCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", server.configPath, string(tomlData))
	result, err := executeCommandWithExec(L, writeCmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Configuration saved to %s: %s", server.configPath, result)))
	L.Push(lua.LNil)
	return 2
}

// serverLoadConfig loads configuration from a TOML file (action method)
func serverLoadConfig(L *lua.LState) int {
	server := checkFrpServer(L, 1)

	// Read config file
	readCmd := fmt.Sprintf("cat %s", server.configPath)
	output, err := executeCommandWithExec(L, readCmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	// Parse TOML
	var config map[string]interface{}
	err = toml.Unmarshal([]byte(output), &config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse TOML: %v", err)))
		return 2
	}

	// Convert to Lua table
	table := L.NewTable()
	for k, v := range config {
		table.RawSetString(k, goValueToLua(L, v))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

// serverStart starts the FRP server service (action method)
func serverStart(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl start frps")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to start frps: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP server started: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// serverStop stops the FRP server service (action method)
func serverStop(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl stop frps")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to stop frps: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP server stopped: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// serverRestart restarts the FRP server service (action method)
func serverRestart(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl restart frps")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to restart frps: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP server restarted: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// serverStatus gets the FRP server service status (action method)
func serverStatus(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl status frps")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		// Status command returns non-zero if service is stopped, but output is still useful
		L.Push(lua.LString(result))
		L.Push(lua.LNil)
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// serverInstall installs FRP server binary and sets up systemd service (action method)
func serverInstall(L *lua.LState) int {
	server := checkFrpServer(L, 1)

	// Download and install FRP
	installScript := fmt.Sprintf(`
set -e

# Determine architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Determine OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Get version
VERSION="%s"
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -s https://api.github.com/repos/fatedier/frp/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
else
    VERSION="${VERSION#v}"
fi

echo "Installing FRP version $VERSION for $OS-$ARCH..."

# Download FRP
DOWNLOAD_URL="https://github.com/fatedier/frp/releases/download/v${VERSION}/frp_${VERSION}_${OS}_${ARCH}.tar.gz"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

curl -L "$DOWNLOAD_URL" -o frp.tar.gz
tar -xzf frp.tar.gz
cd frp_${VERSION}_${OS}_${ARCH}

# Install binaries
sudo mkdir -p /usr/local/bin
sudo cp frps /usr/local/bin/
sudo cp frpc /usr/local/bin/
sudo chmod +x /usr/local/bin/frps /usr/local/bin/frpc

# Create config directory
sudo mkdir -p /etc/frp

# Create systemd service for frps
sudo tee /etc/systemd/system/frps.service > /dev/null <<EOF
[Unit]
Description=FRP Server Service
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frps -c /etc/frp/frps.toml

[Install]
WantedBy=multi-user.target
EOF

# Create systemd service for frpc
sudo tee /etc/systemd/system/frpc.service > /dev/null <<EOF
[Unit]
Description=FRP Client Service
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frpc -c /etc/frp/frpc.toml

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
sudo systemctl daemon-reload

# Clean up
cd /
rm -rf "$TMP_DIR"

echo "FRP $VERSION installed successfully!"
/usr/local/bin/frps --version
`, server.version)

	result, err := executeCommandWithExec(L, installScript)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to install FRP: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP installed successfully: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// serverEnable enables FRP server service to start on boot (action method)
func serverEnable(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl enable frps")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to enable frps: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP server enabled: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// serverDisable disables FRP server service from starting on boot (action method)
func serverDisable(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl disable frps")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to disable frps: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP server disabled: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// ============================================================================
// FRP Client Builder
// ============================================================================

// FrpClient represents an FRP client configuration builder
type FrpClient struct {
	L           *lua.LState
	agentClient interface{}
	name        string
	config      map[string]interface{}
	proxies     []map[string]interface{}
	configPath  string
	target      string
	version     string
}

// createClient creates a new FRP client builder instance
func (m *FrpModule) createClient(L *lua.LState) int {
	name := "frpc"
	if L.GetTop() >= 1 {
		name = L.CheckString(1)
	}

	client := &FrpClient{
		L:           L,
		agentClient: m.agentClient,
		name:        name,
		config:      make(map[string]interface{}),
		proxies:     make([]map[string]interface{}, 0),
		configPath:  "/etc/frp/frpc.toml",
		version:     "latest",
	}

	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("frp_client"))
	L.Push(ud)
	return 1
}

// RegisterClientMetatable registers the FRP client metatable with Lua methods
func RegisterClientMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("frp_client")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"config":       clientConfig,
		"server":       clientServer,
		"proxy":        clientProxy,
		"config_path":  clientConfigPath,
		"version":      clientVersion,
		"delegate_to":  clientDelegateTo,
		"save_config":  clientSaveConfig,
		"load_config":  clientLoadConfig,
		"start":        clientStart,
		"stop":         clientStop,
		"restart":      clientRestart,
		"status":       clientStatus,
		"install":      clientInstall,
		"enable":       clientEnable,
		"disable":      clientDisable,
	}))
}

// checkFrpClient extracts FrpClient from Lua userdata
func checkFrpClient(L *lua.LState, n int) *FrpClient {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*FrpClient); ok {
		return v
	}
	L.ArgError(n, "FrpClient expected")
	return nil
}

// clientConfig sets configuration from a Lua table (fluent method)
func clientConfig(L *lua.LState) int {
	client := checkFrpClient(L, 1)
	config := L.CheckTable(2)

	config.ForEach(func(k, v lua.LValue) {
		key := k.String()
		client.config[key] = luaValueToGo(v)
	})

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// clientServer sets the server address and port (fluent method)
func clientServer(L *lua.LState) int {
	client := checkFrpClient(L, 1)
	address := L.CheckString(2)
	port := L.CheckInt(3)

	client.config["serverAddr"] = address
	client.config["serverPort"] = port

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// clientProxy adds a proxy configuration (fluent method)
func clientProxy(L *lua.LState) int {
	client := checkFrpClient(L, 1)
	proxyConfig := L.CheckTable(2)

	proxy := make(map[string]interface{})
	proxyConfig.ForEach(func(k, v lua.LValue) {
		key := k.String()
		proxy[key] = luaValueToGo(v)
	})

	client.proxies = append(client.proxies, proxy)

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// clientConfigPath sets the configuration file path (fluent method)
func clientConfigPath(L *lua.LState) int {
	client := checkFrpClient(L, 1)
	path := L.CheckString(2)
	client.configPath = path

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// clientVersion sets the FRP version to install (fluent method)
func clientVersion(L *lua.LState) int {
	client := checkFrpClient(L, 1)
	version := L.CheckString(2)
	client.version = version

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// clientDelegateTo sets the target agent for command execution (fluent method)
func clientDelegateTo(L *lua.LState) int {
	client := checkFrpClient(L, 1)
	target := L.CheckString(2)
	client.target = target

	L.Push(L.Get(1))
	L.Push(lua.LNil)
	return 2
}

// clientSaveConfig saves the configuration to a TOML file (action method)
func clientSaveConfig(L *lua.LState) int {
	client := checkFrpClient(L, 1)

	// Build complete config with proxies
	fullConfig := make(map[string]interface{})
	for k, v := range client.config {
		fullConfig[k] = v
	}

	// Add proxies
	if len(client.proxies) > 0 {
		fullConfig["proxies"] = client.proxies
	}

	// Convert config map to TOML
	tomlData, err := toml.Marshal(fullConfig)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to marshal TOML: %v", err)))
		return 2
	}

	// Create directory if needed
	dir := filepath.Dir(client.configPath)
	mkdirCmd := fmt.Sprintf("mkdir -p %s", dir)
	_, err = executeCommandWithExec(L, mkdirCmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create config directory: %v", err)))
		return 2
	}

	// Write config file
	writeCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", client.configPath, string(tomlData))
	result, err := executeCommandWithExec(L, writeCmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to write config: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Configuration saved to %s: %s", client.configPath, result)))
	L.Push(lua.LNil)
	return 2
}

// clientLoadConfig loads configuration from a TOML file (action method)
func clientLoadConfig(L *lua.LState) int {
	client := checkFrpClient(L, 1)

	// Read config file
	readCmd := fmt.Sprintf("cat %s", client.configPath)
	output, err := executeCommandWithExec(L, readCmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read config: %v", err)))
		return 2
	}

	// Parse TOML
	var config map[string]interface{}
	err = toml.Unmarshal([]byte(output), &config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse TOML: %v", err)))
		return 2
	}

	// Convert to Lua table
	table := L.NewTable()
	for k, v := range config {
		table.RawSetString(k, goValueToLua(L, v))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

// clientStart starts the FRP client service (action method)
func clientStart(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl start frpc")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to start frpc: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP client started: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// clientStop stops the FRP client service (action method)
func clientStop(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl stop frpc")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to stop frpc: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP client stopped: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// clientRestart restarts the FRP client service (action method)
func clientRestart(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl restart frpc")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to restart frpc: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP client restarted: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// clientStatus gets the FRP client service status (action method)
func clientStatus(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl status frpc")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		// Status command returns non-zero if service is stopped, but output is still useful
		L.Push(lua.LString(result))
		L.Push(lua.LNil)
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// clientInstall installs FRP client binary and sets up systemd service (action method)
func clientInstall(L *lua.LState) int {
	client := checkFrpClient(L, 1)

	// Use the same installation script as server since it installs both binaries
	installScript := fmt.Sprintf(`
set -e

# Determine architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Determine OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Get version
VERSION="%s"
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -s https://api.github.com/repos/fatedier/frp/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
else
    VERSION="${VERSION#v}"
fi

echo "Installing FRP version $VERSION for $OS-$ARCH..."

# Download FRP
DOWNLOAD_URL="https://github.com/fatedier/frp/releases/download/v${VERSION}/frp_${VERSION}_${OS}_${ARCH}.tar.gz"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

curl -L "$DOWNLOAD_URL" -o frp.tar.gz
tar -xzf frp.tar.gz
cd frp_${VERSION}_${OS}_${ARCH}

# Install binaries
sudo mkdir -p /usr/local/bin
sudo cp frps /usr/local/bin/
sudo cp frpc /usr/local/bin/
sudo chmod +x /usr/local/bin/frps /usr/local/bin/frpc

# Create config directory
sudo mkdir -p /etc/frp

# Create systemd service for frps
sudo tee /etc/systemd/system/frps.service > /dev/null <<EOF
[Unit]
Description=FRP Server Service
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frps -c /etc/frp/frps.toml

[Install]
WantedBy=multi-user.target
EOF

# Create systemd service for frpc
sudo tee /etc/systemd/system/frpc.service > /dev/null <<EOF
[Unit]
Description=FRP Client Service
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frpc -c /etc/frp/frpc.toml

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
sudo systemctl daemon-reload

# Clean up
cd /
rm -rf "$TMP_DIR"

echo "FRP $VERSION installed successfully!"
/usr/local/bin/frpc --version
`, client.version)

	result, err := executeCommandWithExec(L, installScript)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to install FRP: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP installed successfully: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// clientEnable enables FRP client service to start on boot (action method)
func clientEnable(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl enable frpc")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to enable frpc: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP client enabled: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// clientDisable disables FRP client service from starting on boot (action method)
func clientDisable(L *lua.LState) int {

	cmd := fmt.Sprintf("systemctl disable frpc")
	result, err := executeCommandWithExec(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to disable frpc: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP client disabled: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// ============================================================================
// Module-level Install Function
// ============================================================================

// install provides a module-level installation function
func (m *FrpModule) install(L *lua.LState) int {
	version := "latest"
	if L.GetTop() >= 1 {
		version = L.CheckString(1)
	}

	installScript := fmt.Sprintf(`
set -e

# Determine architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Determine OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Get version
VERSION="%s"
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -s https://api.github.com/repos/fatedier/frp/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
else
    VERSION="${VERSION#v}"
fi

echo "Installing FRP version $VERSION for $OS-$ARCH..."

# Download FRP
DOWNLOAD_URL="https://github.com/fatedier/frp/releases/download/v${VERSION}/frp_${VERSION}_${OS}_${ARCH}.tar.gz"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

curl -L "$DOWNLOAD_URL" -o frp.tar.gz
tar -xzf frp.tar.gz
cd frp_${VERSION}_${OS}_${ARCH}

# Install binaries
sudo mkdir -p /usr/local/bin
sudo cp frps /usr/local/bin/
sudo cp frpc /usr/local/bin/
sudo chmod +x /usr/local/bin/frps /usr/local/bin/frpc

# Create config directory
sudo mkdir -p /etc/frp

# Create systemd service for frps
sudo tee /etc/systemd/system/frps.service > /dev/null <<EOF
[Unit]
Description=FRP Server Service
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frps -c /etc/frp/frps.toml

[Install]
WantedBy=multi-user.target
EOF

# Create systemd service for frpc
sudo tee /etc/systemd/system/frpc.service > /dev/null <<EOF
[Unit]
Description=FRP Client Service
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/frpc -c /etc/frp/frpc.toml

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
sudo systemctl daemon-reload

# Clean up
cd /
rm -rf "$TMP_DIR"

echo "FRP $VERSION installed successfully!"
/usr/local/bin/frps --version
`, version)

	result, err := executeCommandWithExec(L, installScript)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to install FRP: %v", err)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("FRP installed successfully: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// ============================================================================
// Helper Functions
// ============================================================================

// executeCommandWithExec executes a command using the Lua exec.run() function
// This ensures proper delegation when running on remote agents
func executeCommandWithExec(L *lua.LState, cmd string) (string, error) {
	// Get the global exec module
	execMod := L.GetGlobal("exec")
	if execMod.Type() == lua.LTNil {
		// Fallback to local execution if exec module not available
		return executeCommandLocal(cmd)
	}

	// Get the run function from exec module
	runFunc := L.GetField(execMod, "run")
	if runFunc.Type() != lua.LTFunction {
		return executeCommandLocal(cmd)
	}

	// Call exec.run(cmd)
	L.Push(runFunc)
	L.Push(lua.LString(cmd))
	if err := L.PCall(1, 2, nil); err != nil {
		return "", fmt.Errorf("exec.run failed: %v", err)
	}

	// Get result table
	result := L.Get(-2)
	errValue := L.Get(-1)
	L.Pop(2)

	// Check for error
	if errValue.Type() != lua.LTNil {
		return "", fmt.Errorf("exec.run returned error: %s", errValue.String())
	}

	// Extract output from result
	if resultTbl, ok := result.(*lua.LTable); ok {
		stdout := L.GetField(resultTbl, "stdout")
		success := L.GetField(resultTbl, "success")

		output := stdout.String()

		if success.Type() == lua.LTBool && !bool(success.(lua.LBool)) {
			stderr := L.GetField(resultTbl, "stderr")
			return output, fmt.Errorf("command failed: %s", stderr.String())
		}

		return output, nil
	}

	return "", fmt.Errorf("unexpected result type from exec.run")
}

// executeCommandLocal executes a command locally (fallback)
func executeCommandLocal(cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	execCmd := exec.CommandContext(ctx, "bash", "-c", cmd)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v", err)
	}
	return string(output), nil
}

// luaValueToGo converts a Lua value to a Go value
func luaValueToGo(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		// Check if it's an array or a map
		isArray := true
		maxIndex := 0
		v.ForEach(func(k, _ lua.LValue) {
			if num, ok := k.(lua.LNumber); ok {
				idx := int(num)
				if idx > maxIndex {
					maxIndex = idx
				}
			} else {
				isArray = false
			}
		})

		if isArray && maxIndex > 0 {
			// Convert to slice
			arr := make([]interface{}, maxIndex)
			v.ForEach(func(k, val lua.LValue) {
				if num, ok := k.(lua.LNumber); ok {
					idx := int(num) - 1 // Lua arrays are 1-indexed
					if idx >= 0 && idx < maxIndex {
						arr[idx] = luaValueToGo(val)
					}
				}
			})
			return arr
		}

		// Convert to map
		m := make(map[string]interface{})
		v.ForEach(func(k, val lua.LValue) {
			m[k.String()] = luaValueToGo(val)
		})
		return m
	default:
		return v.String()
	}
}

// goValueToLua converts a Go value to a Lua value
func goValueToLua(L *lua.LState, val interface{}) lua.LValue {
	if val == nil {
		return lua.LNil
	}

	switch v := val.(type) {
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		table := L.NewTable()
		for i, item := range v {
			table.RawSetInt(i+1, goValueToLua(L, item))
		}
		return table
	case map[string]interface{}:
		table := L.NewTable()
		for key, item := range v {
			table.RawSetString(key, goValueToLua(L, item))
		}
		return table
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}
