package infra

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// FirewallModule provides firewall management functionality
type FirewallModule struct {
	agentClient interface{}
}

// NewFirewallModule creates a new firewall module instance
func NewFirewallModule(agentClient interface{}) *FirewallModule {
	return &FirewallModule{
		agentClient: agentClient,
	}
}

// Register registers the firewall module and its functions in the Lua state
func (m *FirewallModule) Register(L *lua.LState) {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"rule":         m.createRule,
		"detect":       m.detectBackend,
		"enable":       m.enableFirewall,
		"disable":      m.disableFirewall,
		"status":       m.getStatus,
		"list":         m.listRules,
		"flush":        m.flushRules,
		"allow_port":   m.allowPort,
		"deny_port":    m.denyPort,
		"allow_from":   m.allowFrom,
		"deny_from":    m.denyFrom,
		"save":         m.saveRules,
		"reload":       m.reloadFirewall,
	})

	L.SetGlobal("firewall", mod)
}

// ============================================================================
// Firewall Backend Detection
// ============================================================================

// FirewallBackend represents the detected firewall backend
type FirewallBackend string

const (
	BackendFirewalld FirewallBackend = "firewalld"
	BackendUFW       FirewallBackend = "ufw"
	BackendIPTables  FirewallBackend = "iptables"
	BackendNFTables  FirewallBackend = "nftables"
	BackendUnknown   FirewallBackend = "unknown"
)

// detectBackend detects which firewall backend is available
func (m *FirewallModule) detectBackend(L *lua.LState) int {
	backend := detectFirewallBackend(L)
	L.Push(lua.LString(string(backend)))
	return 1
}

// detectFirewallBackend detects the firewall backend
func detectFirewallBackend(L *lua.LState) FirewallBackend {
	// Check for firewalld
	if isCommandAvailable(L, "firewall-cmd") {
		return BackendFirewalld
	}

	// Check for ufw
	if isCommandAvailable(L, "ufw") {
		return BackendUFW
	}

	// Check for nftables
	if isCommandAvailable(L, "nft") {
		return BackendNFTables
	}

	// Check for iptables
	if isCommandAvailable(L, "iptables") {
		return BackendIPTables
	}

	return BackendUnknown
}

// isCommandAvailable checks if a command is available
func isCommandAvailable(L *lua.LState, cmd string) bool {
	checkCmd := fmt.Sprintf("command -v %s >/dev/null 2>&1", cmd)
	_, err := executeFirewallCommand(L, checkCmd)
	return err == nil
}

// ============================================================================
// Firewall Rule Builder (Fluent API)
// ============================================================================

// FirewallRule represents a firewall rule configuration builder
type FirewallRule struct {
	L           *lua.LState
	agentClient interface{}
	backend     FirewallBackend
	action      string // "allow" or "deny"
	protocol    string // "tcp", "udp", "icmp", or "any"
	port        int
	portRange   string // "8000:9000"
	sourceIP    string
	destIP      string
	zone        string // for firewalld
	chain       string // for iptables/nftables
	direction   string // "in", "out", or "both"
	comment     string
}

// createRule creates a new firewall rule builder instance
func (m *FirewallModule) createRule(L *lua.LState) int {
	rule := &FirewallRule{
		L:           L,
		agentClient: m.agentClient,
		backend:     detectFirewallBackend(L),
		protocol:    "tcp",
		zone:        "public",
		chain:       "INPUT",
		direction:   "in",
	}

	ud := L.NewUserData()
	ud.Value = rule
	L.SetMetatable(ud, L.GetTypeMetatable("firewall_rule"))
	L.Push(ud)
	return 1
}

// RegisterRuleMetatable registers the firewall rule metatable with Lua methods
func RegisterRuleMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("firewall_rule")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"allow":      ruleAllow,
		"deny":       ruleDeny,
		"protocol":   ruleProtocol,
		"port":       rulePort,
		"port_range": rulePortRange,
		"from":       ruleFrom,
		"to":         ruleTo,
		"zone":       ruleZone,
		"chain":      ruleChain,
		"direction":  ruleDirection,
		"comment":    ruleComment,
		"apply":      ruleApply,
		"remove":     ruleRemove,
		"exists":     ruleExists,
	}))
}

// checkFirewallRule extracts FirewallRule from Lua userdata
func checkFirewallRule(L *lua.LState, n int) *FirewallRule {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*FirewallRule); ok {
		return v
	}
	L.ArgError(n, "FirewallRule expected")
	return nil
}

// ruleAllow sets the rule action to allow (fluent method)
func ruleAllow(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	rule.action = "allow"
	L.Push(L.Get(1))
	return 1
}

// ruleDeny sets the rule action to deny (fluent method)
func ruleDeny(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	rule.action = "deny"
	L.Push(L.Get(1))
	return 1
}

// ruleProtocol sets the protocol (fluent method)
func ruleProtocol(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	protocol := L.CheckString(2)
	rule.protocol = strings.ToLower(protocol)
	L.Push(L.Get(1))
	return 1
}

// rulePort sets the port (fluent method)
func rulePort(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	port := L.CheckInt(2)
	rule.port = port
	L.Push(L.Get(1))
	return 1
}

// rulePortRange sets a port range (fluent method)
func rulePortRange(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	portRange := L.CheckString(2)
	rule.portRange = portRange
	L.Push(L.Get(1))
	return 1
}

// ruleFrom sets the source IP/CIDR (fluent method)
func ruleFrom(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	sourceIP := L.CheckString(2)
	rule.sourceIP = sourceIP
	L.Push(L.Get(1))
	return 1
}

// ruleTo sets the destination IP/CIDR (fluent method)
func ruleTo(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	destIP := L.CheckString(2)
	rule.destIP = destIP
	L.Push(L.Get(1))
	return 1
}

// ruleZone sets the zone (for firewalld) (fluent method)
func ruleZone(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	zone := L.CheckString(2)
	rule.zone = zone
	L.Push(L.Get(1))
	return 1
}

// ruleChain sets the chain (for iptables/nftables) (fluent method)
func ruleChain(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	chain := L.CheckString(2)
	rule.chain = strings.ToUpper(chain)
	L.Push(L.Get(1))
	return 1
}

// ruleDirection sets the direction (fluent method)
func ruleDirection(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	direction := L.CheckString(2)
	rule.direction = strings.ToLower(direction)
	L.Push(L.Get(1))
	return 1
}

// ruleComment sets a comment for the rule (fluent method)
func ruleComment(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)
	comment := L.CheckString(2)
	rule.comment = comment
	L.Push(L.Get(1))
	return 1
}

// ruleApply applies the firewall rule (action method - idempotent)
func ruleApply(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)

	// Check if rule already exists (idempotent)
	exists, err := rule.checkExists()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to check if rule exists: %v", err)))
		return 2
	}

	if exists {
		L.Push(lua.LString("Rule already exists (idempotent)"))
		L.Push(lua.LNil)
		return 2
	}

	// Build and execute command based on backend
	cmd, err := rule.buildAddCommand()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to build command: %v", err)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to apply rule: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Rule applied successfully: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// ruleRemove removes the firewall rule (action method - idempotent)
func ruleRemove(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)

	// Check if rule exists (idempotent)
	exists, err := rule.checkExists()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to check if rule exists: %v", err)))
		return 2
	}

	if !exists {
		L.Push(lua.LString("Rule does not exist (idempotent)"))
		L.Push(lua.LNil)
		return 2
	}

	// Build and execute command based on backend
	cmd, err := rule.buildRemoveCommand()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to build command: %v", err)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to remove rule: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Rule removed successfully: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// ruleExists checks if the firewall rule exists (query method)
func ruleExists(L *lua.LState) int {
	rule := checkFirewallRule(L, 1)

	exists, err := rule.checkExists()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("failed to check if rule exists: %v", err)))
		return 2
	}

	L.Push(lua.LBool(exists))
	L.Push(lua.LNil)
	return 2
}

// ============================================================================
// Rule Command Builders
// ============================================================================

// buildAddCommand builds the command to add the rule based on the backend
func (r *FirewallRule) buildAddCommand() (string, error) {
	switch r.backend {
	case BackendFirewalld:
		return r.buildFirewalldAddCommand()
	case BackendUFW:
		return r.buildUFWAddCommand()
	case BackendIPTables:
		return r.buildIPTablesAddCommand()
	case BackendNFTables:
		return r.buildNFTablesAddCommand()
	default:
		return "", fmt.Errorf("unsupported firewall backend: %s", r.backend)
	}
}

// buildRemoveCommand builds the command to remove the rule based on the backend
func (r *FirewallRule) buildRemoveCommand() (string, error) {
	switch r.backend {
	case BackendFirewalld:
		return r.buildFirewalldRemoveCommand()
	case BackendUFW:
		return r.buildUFWRemoveCommand()
	case BackendIPTables:
		return r.buildIPTablesRemoveCommand()
	case BackendNFTables:
		return r.buildNFTablesRemoveCommand()
	default:
		return "", fmt.Errorf("unsupported firewall backend: %s", r.backend)
	}
}

// checkExists checks if the rule already exists
func (r *FirewallRule) checkExists() (bool, error) {
	var checkCmd string

	switch r.backend {
	case BackendFirewalld:
		checkCmd = r.buildFirewalldCheckCommand()
	case BackendUFW:
		checkCmd = r.buildUFWCheckCommand()
	case BackendIPTables:
		checkCmd = r.buildIPTablesCheckCommand()
	case BackendNFTables:
		checkCmd = r.buildNFTablesCheckCommand()
	default:
		return false, fmt.Errorf("unsupported firewall backend: %s", r.backend)
	}

	_, err := executeFirewallCommand(r.L, checkCmd)
	return err == nil, nil
}

// ============================================================================
// Firewalld Commands
// ============================================================================

func (r *FirewallRule) buildFirewalldAddCommand() (string, error) {
	parts := []string{"firewall-cmd", "--permanent"}

	if r.zone != "" {
		parts = append(parts, fmt.Sprintf("--zone=%s", r.zone))
	}

	if r.port > 0 {
		portSpec := fmt.Sprintf("%d/%s", r.port, r.protocol)
		if r.action == "allow" {
			parts = append(parts, fmt.Sprintf("--add-port=%s", portSpec))
		} else {
			parts = append(parts, fmt.Sprintf("--remove-port=%s", portSpec))
		}
	} else if r.portRange != "" {
		portSpec := fmt.Sprintf("%s/%s", r.portRange, r.protocol)
		if r.action == "allow" {
			parts = append(parts, fmt.Sprintf("--add-port=%s", portSpec))
		} else {
			parts = append(parts, fmt.Sprintf("--remove-port=%s", portSpec))
		}
	}

	if r.sourceIP != "" {
		if r.action == "allow" {
			parts = append(parts, fmt.Sprintf("--add-source=%s", r.sourceIP))
		} else {
			parts = append(parts, fmt.Sprintf("--remove-source=%s", r.sourceIP))
		}
	}

	cmd := strings.Join(parts, " ")
	return cmd + " && firewall-cmd --reload", nil
}

func (r *FirewallRule) buildFirewalldRemoveCommand() (string, error) {
	parts := []string{"firewall-cmd", "--permanent"}

	if r.zone != "" {
		parts = append(parts, fmt.Sprintf("--zone=%s", r.zone))
	}

	if r.port > 0 {
		portSpec := fmt.Sprintf("%d/%s", r.port, r.protocol)
		parts = append(parts, fmt.Sprintf("--remove-port=%s", portSpec))
	} else if r.portRange != "" {
		portSpec := fmt.Sprintf("%s/%s", r.portRange, r.protocol)
		parts = append(parts, fmt.Sprintf("--remove-port=%s", portSpec))
	}

	if r.sourceIP != "" {
		parts = append(parts, fmt.Sprintf("--remove-source=%s", r.sourceIP))
	}

	cmd := strings.Join(parts, " ")
	return cmd + " && firewall-cmd --reload", nil
}

func (r *FirewallRule) buildFirewalldCheckCommand() string {
	parts := []string{"firewall-cmd"}

	if r.zone != "" {
		parts = append(parts, fmt.Sprintf("--zone=%s", r.zone))
	}

	if r.port > 0 {
		portSpec := fmt.Sprintf("%d/%s", r.port, r.protocol)
		parts = append(parts, fmt.Sprintf("--query-port=%s", portSpec))
	} else if r.sourceIP != "" {
		parts = append(parts, fmt.Sprintf("--query-source=%s", r.sourceIP))
	}

	return strings.Join(parts, " ")
}

// ============================================================================
// UFW Commands
// ============================================================================

func (r *FirewallRule) buildUFWAddCommand() (string, error) {
	parts := []string{"ufw"}

	if r.action == "allow" {
		parts = append(parts, "allow")
	} else {
		parts = append(parts, "deny")
	}

	if r.sourceIP != "" {
		parts = append(parts, "from", r.sourceIP)
	}

	if r.port > 0 {
		parts = append(parts, "to", "any", "port", fmt.Sprintf("%d", r.port))
		if r.protocol != "any" {
			parts = append(parts, "proto", r.protocol)
		}
	} else if r.portRange != "" {
		parts = append(parts, "to", "any", "port", r.portRange)
		if r.protocol != "any" {
			parts = append(parts, "proto", r.protocol)
		}
	}

	if r.comment != "" {
		parts = append(parts, "comment", fmt.Sprintf("'%s'", r.comment))
	}

	return strings.Join(parts, " "), nil
}

func (r *FirewallRule) buildUFWRemoveCommand() (string, error) {
	addCmd, err := r.buildUFWAddCommand()
	if err != nil {
		return "", err
	}
	// UFW uses "delete" prefix to remove rules
	return "ufw delete " + strings.TrimPrefix(addCmd, "ufw "), nil
}

func (r *FirewallRule) buildUFWCheckCommand() string {
	// Check if rule exists by listing and grepping
	pattern := ""
	if r.port > 0 {
		pattern = fmt.Sprintf("%d/%s", r.port, r.protocol)
	}
	return fmt.Sprintf("ufw status numbered | grep -q '%s'", pattern)
}

// ============================================================================
// IPTables Commands
// ============================================================================

func (r *FirewallRule) buildIPTablesAddCommand() (string, error) {
	parts := []string{"iptables", "-A", r.chain}

	if r.protocol != "any" && r.protocol != "" {
		parts = append(parts, "-p", r.protocol)
	}

	if r.sourceIP != "" {
		parts = append(parts, "-s", r.sourceIP)
	}

	if r.destIP != "" {
		parts = append(parts, "-d", r.destIP)
	}

	if r.port > 0 {
		parts = append(parts, "--dport", fmt.Sprintf("%d", r.port))
	} else if r.portRange != "" {
		parts = append(parts, "--dport", r.portRange)
	}

	if r.action == "allow" {
		parts = append(parts, "-j", "ACCEPT")
	} else {
		parts = append(parts, "-j", "DROP")
	}

	if r.comment != "" {
		parts = append(parts, "-m", "comment", "--comment", fmt.Sprintf("\"%s\"", r.comment))
	}

	return strings.Join(parts, " "), nil
}

func (r *FirewallRule) buildIPTablesRemoveCommand() (string, error) {
	// Replace -A with -D to delete the rule
	addCmd, err := r.buildIPTablesAddCommand()
	if err != nil {
		return "", err
	}
	return strings.Replace(addCmd, "-A ", "-D ", 1), nil
}

func (r *FirewallRule) buildIPTablesCheckCommand() string {
	// Check if rule exists using -C (check) flag
	addCmd, _ := r.buildIPTablesAddCommand()
	return strings.Replace(addCmd, "-A ", "-C ", 1)
}

// ============================================================================
// NFTables Commands
// ============================================================================

func (r *FirewallRule) buildNFTablesAddCommand() (string, error) {
	// NFTables has a different syntax, simplified version
	action := "accept"
	if r.action == "deny" {
		action = "drop"
	}

	rule := fmt.Sprintf("nft add rule inet filter %s", strings.ToLower(r.chain))

	if r.protocol != "any" && r.protocol != "" {
		rule += fmt.Sprintf(" %s", r.protocol)
	}

	if r.port > 0 {
		rule += fmt.Sprintf(" dport %d", r.port)
	}

	if r.sourceIP != "" {
		rule += fmt.Sprintf(" ip saddr %s", r.sourceIP)
	}

	rule += fmt.Sprintf(" %s", action)

	return rule, nil
}

func (r *FirewallRule) buildNFTablesRemoveCommand() (string, error) {
	// NFTables removal is complex, would need to find handle first
	// Simplified: flush and recreate without this rule
	return "", fmt.Errorf("nftables rule removal not fully implemented, use flush and recreate")
}

func (r *FirewallRule) buildNFTablesCheckCommand() string {
	// Check by listing rules and grepping
	pattern := ""
	if r.port > 0 {
		pattern = fmt.Sprintf("dport %d", r.port)
	}
	return fmt.Sprintf("nft list ruleset | grep -q '%s'", pattern)
}

// ============================================================================
// Convenience Functions
// ============================================================================

// allowPort is a convenience function to allow a port
func (m *FirewallModule) allowPort(L *lua.LState) int {
	port := L.CheckInt(1)
	protocol := "tcp"
	if L.GetTop() >= 2 {
		protocol = L.CheckString(2)
	}

	rule := &FirewallRule{
		L:        L,
		backend:  detectFirewallBackend(L),
		action:   "allow",
		protocol: protocol,
		port:     port,
		zone:     "public",
		chain:    "INPUT",
	}

	// Check if exists (idempotent)
	exists, _ := rule.checkExists()
	if exists {
		L.Push(lua.LString(fmt.Sprintf("Port %d/%s already allowed (idempotent)", port, protocol)))
		L.Push(lua.LNil)
		return 2
	}

	cmd, err := rule.buildAddCommand()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to build command: %v", err)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to allow port: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Port %d/%s allowed successfully", port, protocol)))
	L.Push(lua.LNil)
	return 2
}

// denyPort is a convenience function to deny a port
func (m *FirewallModule) denyPort(L *lua.LState) int {
	port := L.CheckInt(1)
	protocol := "tcp"
	if L.GetTop() >= 2 {
		protocol = L.CheckString(2)
	}

	rule := &FirewallRule{
		L:        L,
		backend:  detectFirewallBackend(L),
		action:   "deny",
		protocol: protocol,
		port:     port,
		zone:     "public",
		chain:    "INPUT",
	}

	// Check if exists (idempotent)
	exists, _ := rule.checkExists()
	if exists {
		L.Push(lua.LString(fmt.Sprintf("Port %d/%s already denied (idempotent)", port, protocol)))
		L.Push(lua.LNil)
		return 2
	}

	cmd, err := rule.buildAddCommand()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to build command: %v", err)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to deny port: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Port %d/%s denied successfully", port, protocol)))
	L.Push(lua.LNil)
	return 2
}

// allowFrom is a convenience function to allow traffic from an IP/CIDR
func (m *FirewallModule) allowFrom(L *lua.LState) int {
	sourceIP := L.CheckString(1)

	rule := &FirewallRule{
		L:        L,
		backend:  detectFirewallBackend(L),
		action:   "allow",
		sourceIP: sourceIP,
		zone:     "public",
		chain:    "INPUT",
	}

	// Check if exists (idempotent)
	exists, _ := rule.checkExists()
	if exists {
		L.Push(lua.LString(fmt.Sprintf("Source %s already allowed (idempotent)", sourceIP)))
		L.Push(lua.LNil)
		return 2
	}

	cmd, err := rule.buildAddCommand()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to build command: %v", err)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to allow source: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Source %s allowed successfully", sourceIP)))
	L.Push(lua.LNil)
	return 2
}

// denyFrom is a convenience function to deny traffic from an IP/CIDR
func (m *FirewallModule) denyFrom(L *lua.LState) int {
	sourceIP := L.CheckString(1)

	rule := &FirewallRule{
		L:        L,
		backend:  detectFirewallBackend(L),
		action:   "deny",
		sourceIP: sourceIP,
		zone:     "public",
		chain:    "INPUT",
	}

	// Check if exists (idempotent)
	exists, _ := rule.checkExists()
	if exists {
		L.Push(lua.LString(fmt.Sprintf("Source %s already denied (idempotent)", sourceIP)))
		L.Push(lua.LNil)
		return 2
	}

	cmd, err := rule.buildAddCommand()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to build command: %v", err)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to deny source: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Source %s denied successfully", sourceIP)))
	L.Push(lua.LNil)
	return 2
}

// ============================================================================
// Firewall Management Functions
// ============================================================================

// enableFirewall enables the firewall service
func (m *FirewallModule) enableFirewall(L *lua.LState) int {
	backend := detectFirewallBackend(L)

	var cmd string
	switch backend {
	case BackendFirewalld:
		cmd = "systemctl enable --now firewalld"
	case BackendUFW:
		cmd = "ufw enable"
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("enable not supported for backend: %s", backend)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to enable firewall: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Firewall enabled: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// disableFirewall disables the firewall service
func (m *FirewallModule) disableFirewall(L *lua.LState) int {
	backend := detectFirewallBackend(L)

	var cmd string
	switch backend {
	case BackendFirewalld:
		cmd = "systemctl disable --now firewalld"
	case BackendUFW:
		cmd = "ufw disable"
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("disable not supported for backend: %s", backend)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to disable firewall: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Firewall disabled: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// getStatus gets the firewall status
func (m *FirewallModule) getStatus(L *lua.LState) int {
	backend := detectFirewallBackend(L)

	var cmd string
	switch backend {
	case BackendFirewalld:
		cmd = "firewall-cmd --state"
	case BackendUFW:
		cmd = "ufw status"
	case BackendIPTables:
		cmd = "iptables -L -n -v"
	case BackendNFTables:
		cmd = "nft list ruleset"
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("status not supported for backend: %s", backend)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LString(result)) // Still return output even on error
		L.Push(lua.LNil)
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// listRules lists all firewall rules
func (m *FirewallModule) listRules(L *lua.LState) int {
	backend := detectFirewallBackend(L)

	var cmd string
	switch backend {
	case BackendFirewalld:
		cmd = "firewall-cmd --list-all"
	case BackendUFW:
		cmd = "ufw status numbered"
	case BackendIPTables:
		cmd = "iptables -L -n --line-numbers"
	case BackendNFTables:
		cmd = "nft list ruleset"
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("list not supported for backend: %s", backend)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to list rules: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// flushRules flushes all firewall rules
func (m *FirewallModule) flushRules(L *lua.LState) int {
	backend := detectFirewallBackend(L)

	var cmd string
	switch backend {
	case BackendFirewalld:
		cmd = "firewall-cmd --reload"
	case BackendUFW:
		cmd = "ufw --force reset"
	case BackendIPTables:
		cmd = "iptables -F"
	case BackendNFTables:
		cmd = "nft flush ruleset"
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("flush not supported for backend: %s", backend)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to flush rules: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Rules flushed: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// saveRules saves the current firewall rules
func (m *FirewallModule) saveRules(L *lua.LState) int {
	backend := detectFirewallBackend(L)

	var cmd string
	switch backend {
	case BackendFirewalld:
		cmd = "firewall-cmd --runtime-to-permanent"
	case BackendUFW:
		L.Push(lua.LString("UFW rules are automatically saved"))
		L.Push(lua.LNil)
		return 2
	case BackendIPTables:
		cmd = "iptables-save > /etc/iptables/rules.v4"
	case BackendNFTables:
		cmd = "nft list ruleset > /etc/nftables.conf"
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("save not supported for backend: %s", backend)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to save rules: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Rules saved: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// reloadFirewall reloads the firewall
func (m *FirewallModule) reloadFirewall(L *lua.LState) int {
	backend := detectFirewallBackend(L)

	var cmd string
	switch backend {
	case BackendFirewalld:
		cmd = "firewall-cmd --reload"
	case BackendUFW:
		cmd = "ufw reload"
	case BackendIPTables:
		cmd = "iptables-restore < /etc/iptables/rules.v4"
	case BackendNFTables:
		cmd = "nft -f /etc/nftables.conf"
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("reload not supported for backend: %s", backend)))
		return 2
	}

	result, err := executeFirewallCommand(L, cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to reload firewall: %v - %s", err, result)))
		return 2
	}

	L.Push(lua.LString(fmt.Sprintf("Firewall reloaded: %s", result)))
	L.Push(lua.LNil)
	return 2
}

// ============================================================================
// Helper Functions
// ============================================================================

// executeFirewallCommand executes a firewall command using exec module or fallback
func executeFirewallCommand(L *lua.LState, cmd string) (string, error) {
	// Try to use exec module if available
	execMod := L.GetGlobal("exec")
	if execMod.Type() != lua.LTNil {
		runFunc := L.GetField(execMod, "run")
		if runFunc.Type() == lua.LTFunction {
			L.Push(runFunc)
			L.Push(lua.LString(cmd))
			if err := L.PCall(1, 2, nil); err == nil {
				result := L.Get(-2)
				errValue := L.Get(-1)
				L.Pop(2)

				if errValue.Type() == lua.LTNil {
					if resultTbl, ok := result.(*lua.LTable); ok {
						stdout := L.GetField(resultTbl, "stdout")
						success := L.GetField(resultTbl, "success")

						if success.Type() == lua.LTBool && bool(success.(lua.LBool)) {
							return stdout.String(), nil
						}

						stderr := L.GetField(resultTbl, "stderr")
						return stdout.String(), fmt.Errorf("%s", stderr.String())
					}
				}
			}
		}
	}

	// Fallback to local execution
	return executeCommandLocalFirewall(cmd)
}

// executeCommandLocalFirewall executes a command locally (fallback)
func executeCommandLocalFirewall(cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	execCmd := exec.CommandContext(ctx, "bash", "-c", cmd)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v", err)
	}
	return string(output), nil
}
