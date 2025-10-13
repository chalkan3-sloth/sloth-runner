package infra

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestFirewallModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create and register firewall module
	module := NewFirewallModule(nil)
	module.Register(L)
	RegisterRuleMetatable(L)

	// Test that firewall module is registered
	firewallMod := L.GetGlobal("firewall")
	if firewallMod.Type() == lua.LTNil {
		t.Fatal("firewall module not registered")
	}
}

func TestFirewallDetectBackend(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)

	// Test detect function
	err := L.DoString(`
		backend = firewall.detect()
	`)
	if err != nil {
		t.Fatalf("Failed to call firewall.detect(): %v", err)
	}

	backend := L.GetGlobal("backend")
	if backend.Type() != lua.LTString {
		t.Fatal("firewall.detect() should return a string")
	}

	backendStr := backend.String()
	validBackends := []string{"firewalld", "ufw", "iptables", "nftables", "unknown"}
	isValid := false
	for _, valid := range validBackends {
		if backendStr == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		t.Fatalf("Invalid backend detected: %s", backendStr)
	}
}

func TestFirewallRuleBuilder(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)
	RegisterRuleMetatable(L)

	// Test rule builder fluent API
	err := L.DoString(`
		rule = firewall.rule()
			:allow()
			:port(22)
			:protocol("tcp")
			:from("192.168.1.0/24")
			:zone("public")
			:comment("SSH access from LAN")
	`)
	if err != nil {
		t.Fatalf("Failed to create rule with fluent API: %v", err)
	}

	rule := L.GetGlobal("rule")
	if rule.Type() != lua.LTUserData {
		t.Fatal("firewall.rule() should return userdata")
	}

	// Verify we can extract the rule
	ud, ok := rule.(*lua.LUserData)
	if !ok {
		t.Fatal("Expected userdata")
	}

	firewallRule, ok := ud.Value.(*FirewallRule)
	if !ok {
		t.Fatal("Expected FirewallRule")
	}

	// Verify rule properties
	if firewallRule.action != "allow" {
		t.Errorf("Expected action 'allow', got '%s'", firewallRule.action)
	}

	if firewallRule.port != 22 {
		t.Errorf("Expected port 22, got %d", firewallRule.port)
	}

	if firewallRule.protocol != "tcp" {
		t.Errorf("Expected protocol 'tcp', got '%s'", firewallRule.protocol)
	}

	if firewallRule.sourceIP != "192.168.1.0/24" {
		t.Errorf("Expected sourceIP '192.168.1.0/24', got '%s'", firewallRule.sourceIP)
	}

	if firewallRule.zone != "public" {
		t.Errorf("Expected zone 'public', got '%s'", firewallRule.zone)
	}

	if firewallRule.comment != "SSH access from LAN" {
		t.Errorf("Expected comment 'SSH access from LAN', got '%s'", firewallRule.comment)
	}
}

func TestFirewallRuleDeny(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)
	RegisterRuleMetatable(L)

	// Test deny rule
	err := L.DoString(`
		rule = firewall.rule()
			:deny()
			:port(23)
			:protocol("tcp")
	`)
	if err != nil {
		t.Fatalf("Failed to create deny rule: %v", err)
	}

	rule := L.GetGlobal("rule")
	ud, _ := rule.(*lua.LUserData)
	firewallRule, _ := ud.Value.(*FirewallRule)

	if firewallRule.action != "deny" {
		t.Errorf("Expected action 'deny', got '%s'", firewallRule.action)
	}
}

func TestFirewallRulePortRange(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)
	RegisterRuleMetatable(L)

	// Test port range
	err := L.DoString(`
		rule = firewall.rule()
			:allow()
			:port_range("8000:9000")
			:protocol("tcp")
	`)
	if err != nil {
		t.Fatalf("Failed to create rule with port range: %v", err)
	}

	rule := L.GetGlobal("rule")
	ud, _ := rule.(*lua.LUserData)
	firewallRule, _ := ud.Value.(*FirewallRule)

	if firewallRule.portRange != "8000:9000" {
		t.Errorf("Expected portRange '8000:9000', got '%s'", firewallRule.portRange)
	}
}

func TestFirewallRuleDirection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)
	RegisterRuleMetatable(L)

	// Test direction
	err := L.DoString(`
		rule = firewall.rule()
			:allow()
			:port(80)
			:direction("out")
	`)
	if err != nil {
		t.Fatalf("Failed to create rule with direction: %v", err)
	}

	rule := L.GetGlobal("rule")
	ud, _ := rule.(*lua.LUserData)
	firewallRule, _ := ud.Value.(*FirewallRule)

	if firewallRule.direction != "out" {
		t.Errorf("Expected direction 'out', got '%s'", firewallRule.direction)
	}
}

func TestFirewallRuleChain(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)
	RegisterRuleMetatable(L)

	// Test chain
	err := L.DoString(`
		rule = firewall.rule()
			:allow()
			:port(443)
			:chain("OUTPUT")
	`)
	if err != nil {
		t.Fatalf("Failed to create rule with chain: %v", err)
	}

	rule := L.GetGlobal("rule")
	ud, _ := rule.(*lua.LUserData)
	firewallRule, _ := ud.Value.(*FirewallRule)

	if firewallRule.chain != "OUTPUT" {
		t.Errorf("Expected chain 'OUTPUT', got '%s'", firewallRule.chain)
	}
}

func TestBuildFirewalldCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendFirewalld,
		action:   "allow",
		protocol: "tcp",
		port:     22,
		zone:     "public",
	}

	// Test add command
	addCmd, err := rule.buildFirewalldAddCommand()
	if err != nil {
		t.Fatalf("Failed to build firewalld add command: %v", err)
	}

	expectedAdd := "firewall-cmd --permanent --zone=public --add-port=22/tcp && firewall-cmd --reload"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}

	// Test remove command
	removeCmd, err := rule.buildFirewalldRemoveCommand()
	if err != nil {
		t.Fatalf("Failed to build firewalld remove command: %v", err)
	}

	expectedRemove := "firewall-cmd --permanent --zone=public --remove-port=22/tcp && firewall-cmd --reload"
	if removeCmd != expectedRemove {
		t.Errorf("Expected remove command:\n%s\nGot:\n%s", expectedRemove, removeCmd)
	}

	// Test check command
	checkCmd := rule.buildFirewalldCheckCommand()
	expectedCheck := "firewall-cmd --zone=public --query-port=22/tcp"
	if checkCmd != expectedCheck {
		t.Errorf("Expected check command:\n%s\nGot:\n%s", expectedCheck, checkCmd)
	}
}

func TestBuildFirewalldSourceCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendFirewalld,
		action:   "allow",
		sourceIP: "192.168.1.0/24",
		zone:     "public",
	}

	// Test add command with source
	addCmd, err := rule.buildFirewalldAddCommand()
	if err != nil {
		t.Fatalf("Failed to build firewalld add command with source: %v", err)
	}

	expectedAdd := "firewall-cmd --permanent --zone=public --add-source=192.168.1.0/24 && firewall-cmd --reload"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}
}

func TestBuildUFWCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendUFW,
		action:   "allow",
		protocol: "tcp",
		port:     22,
	}

	// Test add command
	addCmd, err := rule.buildUFWAddCommand()
	if err != nil {
		t.Fatalf("Failed to build UFW add command: %v", err)
	}

	expectedAdd := "ufw allow to any port 22 proto tcp"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}

	// Test remove command
	removeCmd, err := rule.buildUFWRemoveCommand()
	if err != nil {
		t.Fatalf("Failed to build UFW remove command: %v", err)
	}

	expectedRemove := "ufw delete allow to any port 22 proto tcp"
	if removeCmd != expectedRemove {
		t.Errorf("Expected remove command:\n%s\nGot:\n%s", expectedRemove, removeCmd)
	}
}

func TestBuildUFWSourceCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendUFW,
		action:   "allow",
		protocol: "tcp",
		port:     22,
		sourceIP: "192.168.1.0/24",
	}

	// Test add command with source
	addCmd, err := rule.buildUFWAddCommand()
	if err != nil {
		t.Fatalf("Failed to build UFW add command with source: %v", err)
	}

	expectedAdd := "ufw allow from 192.168.1.0/24 to any port 22 proto tcp"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}
}

func TestBuildIPTablesCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendIPTables,
		action:   "allow",
		protocol: "tcp",
		port:     22,
		chain:    "INPUT",
	}

	// Test add command
	addCmd, err := rule.buildIPTablesAddCommand()
	if err != nil {
		t.Fatalf("Failed to build iptables add command: %v", err)
	}

	expectedAdd := "iptables -A INPUT -p tcp --dport 22 -j ACCEPT"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}

	// Test remove command
	removeCmd, err := rule.buildIPTablesRemoveCommand()
	if err != nil {
		t.Fatalf("Failed to build iptables remove command: %v", err)
	}

	expectedRemove := "iptables -D INPUT -p tcp --dport 22 -j ACCEPT"
	if removeCmd != expectedRemove {
		t.Errorf("Expected remove command:\n%s\nGot:\n%s", expectedRemove, removeCmd)
	}

	// Test check command
	checkCmd := rule.buildIPTablesCheckCommand()
	expectedCheck := "iptables -C INPUT -p tcp --dport 22 -j ACCEPT"
	if checkCmd != expectedCheck {
		t.Errorf("Expected check command:\n%s\nGot:\n%s", expectedCheck, checkCmd)
	}
}

func TestBuildIPTablesSourceCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendIPTables,
		action:   "allow",
		protocol: "tcp",
		port:     22,
		sourceIP: "192.168.1.0/24",
		chain:    "INPUT",
	}

	// Test add command with source
	addCmd, err := rule.buildIPTablesAddCommand()
	if err != nil {
		t.Fatalf("Failed to build iptables add command with source: %v", err)
	}

	expectedAdd := "iptables -A INPUT -p tcp -s 192.168.1.0/24 --dport 22 -j ACCEPT"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}
}

func TestBuildIPTablesDenyCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendIPTables,
		action:   "deny",
		protocol: "tcp",
		port:     23,
		chain:    "INPUT",
	}

	// Test deny command
	addCmd, err := rule.buildIPTablesAddCommand()
	if err != nil {
		t.Fatalf("Failed to build iptables deny command: %v", err)
	}

	expectedAdd := "iptables -A INPUT -p tcp --dport 23 -j DROP"
	if addCmd != expectedAdd {
		t.Errorf("Expected deny command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}
}

func TestBuildIPTablesWithComment(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendIPTables,
		action:   "allow",
		protocol: "tcp",
		port:     22,
		chain:    "INPUT",
		comment:  "SSH access",
	}

	// Test command with comment
	addCmd, err := rule.buildIPTablesAddCommand()
	if err != nil {
		t.Fatalf("Failed to build iptables command with comment: %v", err)
	}

	expectedAdd := "iptables -A INPUT -p tcp --dport 22 -j ACCEPT -m comment --comment \"SSH access\""
	if addCmd != expectedAdd {
		t.Errorf("Expected command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}
}

func TestBuildNFTablesCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendNFTables,
		action:   "allow",
		protocol: "tcp",
		port:     22,
		chain:    "INPUT",
	}

	// Test add command
	addCmd, err := rule.buildNFTablesAddCommand()
	if err != nil {
		t.Fatalf("Failed to build nftables add command: %v", err)
	}

	expectedAdd := "nft add rule inet filter input tcp dport 22 accept"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}
}

func TestBuildNFTablesSourceCommands(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	rule := &FirewallRule{
		L:        L,
		backend:  BackendNFTables,
		action:   "allow",
		protocol: "tcp",
		port:     22,
		sourceIP: "192.168.1.0/24",
		chain:    "INPUT",
	}

	// Test add command with source
	addCmd, err := rule.buildNFTablesAddCommand()
	if err != nil {
		t.Fatalf("Failed to build nftables add command with source: %v", err)
	}

	expectedAdd := "nft add rule inet filter input tcp dport 22 ip saddr 192.168.1.0/24 accept"
	if addCmd != expectedAdd {
		t.Errorf("Expected add command:\n%s\nGot:\n%s", expectedAdd, addCmd)
	}
}

func TestFirewallBackendDetection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	backend := detectFirewallBackend(L)

	// Should return a valid backend type
	validBackends := []FirewallBackend{
		BackendFirewalld,
		BackendUFW,
		BackendIPTables,
		BackendNFTables,
		BackendUnknown,
	}

	isValid := false
	for _, valid := range validBackends {
		if backend == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		t.Errorf("Invalid backend detected: %s", backend)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)

	// Test that convenience functions are registered
	functions := []string{
		"allow_port",
		"deny_port",
		"allow_from",
		"deny_from",
	}

	for _, funcName := range functions {
		firewallMod := L.GetGlobal("firewall")
		funcValue := L.GetField(firewallMod, funcName)
		if funcValue.Type() != lua.LTFunction {
			t.Errorf("Function %s not registered or not a function", funcName)
		}
	}
}

func TestManagementFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewFirewallModule(nil)
	module.Register(L)

	// Test that management functions are registered
	functions := []string{
		"enable",
		"disable",
		"status",
		"list",
		"flush",
		"save",
		"reload",
	}

	for _, funcName := range functions {
		firewallMod := L.GetGlobal("firewall")
		funcValue := L.GetField(firewallMod, funcName)
		if funcValue.Type() != lua.LTFunction {
			t.Errorf("Function %s not registered or not a function", funcName)
		}
	}
}
