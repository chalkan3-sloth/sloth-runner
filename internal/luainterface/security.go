package luainterface

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// SecurityModule provides security scanning and compliance features
type SecurityModule struct{}

// NewSecurityModule creates a new security module
func NewSecurityModule() *SecurityModule {
	return &SecurityModule{}
}

// RegisterSecurityModule registers the security module with the Lua state
func RegisterSecurityModule(L *lua.LState) {
	module := NewSecurityModule()
	
	// Create the security table
	securityTable := L.NewTable()
	
	// Port scanning and network security
	L.SetField(securityTable, "scan_ports", L.NewFunction(module.luaScanPorts))
	L.SetField(securityTable, "check_ssl_cert", L.NewFunction(module.luaCheckSSLCert))
	L.SetField(securityTable, "check_http_headers", L.NewFunction(module.luaCheckHTTPHeaders))
	
	// File and directory security
	L.SetField(securityTable, "audit_permissions", L.NewFunction(module.luaAuditPermissions))
	L.SetField(securityTable, "find_suid_files", L.NewFunction(module.luaFindSUIDFiles))
	L.SetField(securityTable, "check_file_integrity", L.NewFunction(module.luaCheckFileIntegrity))
	
	// Password and authentication
	L.SetField(securityTable, "password_strength", L.NewFunction(module.luaPasswordStrength))
	L.SetField(securityTable, "check_weak_passwords", L.NewFunction(module.luaCheckWeakPasswords))
	
	// System security
	L.SetField(securityTable, "firewall_status", L.NewFunction(module.luaFirewallStatus))
	L.SetField(securityTable, "selinux_status", L.NewFunction(module.luaSELinuxStatus))
	L.SetField(securityTable, "check_updates", L.NewFunction(module.luaCheckUpdates))
	
	// Vulnerability scanning
	L.SetField(securityTable, "vulnerability_scan", L.NewFunction(module.luaVulnerabilityScan))
	L.SetField(securityTable, "check_open_ports", L.NewFunction(module.luaCheckOpenPorts))
	L.SetField(securityTable, "malware_scan", L.NewFunction(module.luaMalwareScan))
	
	// Compliance and hardening
	L.SetField(securityTable, "security_baseline", L.NewFunction(module.luaSecurityBaseline))
	L.SetField(securityTable, "cis_benchmark", L.NewFunction(module.luaCISBenchmark))
	
	// Register the security table globally
	L.SetGlobal("security", securityTable)
}

// Port scanning and network security
func (s *SecurityModule) luaScanPorts(L *lua.LState) int {
	host := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	// Parse options
	startPort := 1
	endPort := 1024
	timeout := 2
	
	if startVal := options.RawGetString("start"); startVal != lua.LNil {
		startPort = int(startVal.(lua.LNumber))
	}
	if endVal := options.RawGetString("end"); endVal != lua.LNil {
		endPort = int(endVal.(lua.LNumber))
	}
	if timeoutVal := options.RawGetString("timeout"); timeoutVal != lua.LNil {
		timeout = int(timeoutVal.(lua.LNumber))
	}
	
	result := L.NewTable()
	openPorts := L.NewTable()
	closedPorts := L.NewTable()
	
	openCount := 1
	closedCount := 1
	
	for port := startPort; port <= endPort; port++ {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Second)
		if err == nil {
			conn.Close()
			
			portInfo := L.NewTable()
			L.SetField(portInfo, "port", lua.LNumber(port))
			L.SetField(portInfo, "service", lua.LString(getServiceName(port)))
			L.SetField(portInfo, "status", lua.LString("open"))
			
			openPorts.RawSetInt(openCount, portInfo)
			openCount++
		} else {
			closedPorts.RawSetInt(closedCount, lua.LNumber(port))
			closedCount++
		}
	}
	
	L.SetField(result, "open", openPorts)
	L.SetField(result, "closed", closedPorts)
	L.SetField(result, "total_scanned", lua.LNumber(endPort-startPort+1))
	L.SetField(result, "open_count", lua.LNumber(openCount-1))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaCheckSSLCert(L *lua.LState) int {
	url := L.CheckString(1)
	
	// Extract hostname from URL
	hostname := strings.TrimPrefix(url, "https://")
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.Split(hostname, "/")[0]
	hostname = strings.Split(hostname, ":")[0]
	
	// Connect and get certificate
	conn, err := tls.Dial("tcp", hostname+":443", &tls.Config{})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer conn.Close()
	
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		L.Push(lua.LNil)
		L.Push(lua.LString("no certificates found"))
		return 2
	}
	
	cert := certs[0]
	result := L.NewTable()
	
	L.SetField(result, "subject", lua.LString(cert.Subject.String()))
	L.SetField(result, "issuer", lua.LString(cert.Issuer.String()))
	L.SetField(result, "not_before", lua.LNumber(cert.NotBefore.Unix()))
	L.SetField(result, "not_after", lua.LNumber(cert.NotAfter.Unix()))
	L.SetField(result, "serial_number", lua.LString(cert.SerialNumber.String()))
	
	// Check if certificate is valid
	now := time.Now()
	L.SetField(result, "valid", lua.LBool(now.After(cert.NotBefore) && now.Before(cert.NotAfter)))
	L.SetField(result, "expires_in_days", lua.LNumber(cert.NotAfter.Sub(now).Hours()/24))
	
	// Check for common vulnerabilities
	issues := L.NewTable()
	issueCount := 1
	
	// Check expiration
	if cert.NotAfter.Sub(now).Hours() < 24*30 { // Less than 30 days
		issues.RawSetInt(issueCount, lua.LString("Certificate expires soon"))
		issueCount++
	}
	
	// Check weak signature algorithm
	if cert.SignatureAlgorithm.String() == "SHA1-RSA" {
		issues.RawSetInt(issueCount, lua.LString("Weak signature algorithm: SHA1"))
		issueCount++
	}
	
	L.SetField(result, "issues", issues)
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaCheckHTTPHeaders(L *lua.LState) int {
	url := L.CheckString(1)
	
	resp, err := http.Get(url)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	
	result := L.NewTable()
	headers := L.NewTable()
	security_issues := L.NewTable()
	issueCount := 1
	
	// Collect all headers
	for name, values := range resp.Header {
		headers.RawSetString(name, lua.LString(strings.Join(values, ", ")))
	}
	L.SetField(result, "headers", headers)
	
	// Check for security headers
	securityHeaders := []string{
		"Strict-Transport-Security",
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Content-Security-Policy",
		"Referrer-Policy",
	}
	
	missing := L.NewTable()
	missingCount := 1
	
	for _, header := range securityHeaders {
		if resp.Header.Get(header) == "" {
			missing.RawSetInt(missingCount, lua.LString(header))
			missingCount++
			security_issues.RawSetInt(issueCount, lua.LString("Missing security header: "+header))
			issueCount++
		}
	}
	
	L.SetField(result, "missing_security_headers", missing)
	L.SetField(result, "security_issues", security_issues)
	L.SetField(result, "status_code", lua.LNumber(resp.StatusCode))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// File and directory security
func (s *SecurityModule) luaAuditPermissions(L *lua.LState) int {
	path := L.CheckString(1)
	
	result := L.NewTable()
	issues := L.NewTable()
	issueCount := 1
	
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking
		}
		
		mode := info.Mode()
		perm := mode.Perm()
		
		// Check for world-writable files
		if perm&0002 != 0 {
			issues.RawSetInt(issueCount, lua.LString("World-writable file: "+filePath))
			issueCount++
		}
		
		// Check for SUID/SGID files
		if mode&os.ModeSetuid != 0 {
			issues.RawSetInt(issueCount, lua.LString("SUID file: "+filePath))
			issueCount++
		}
		if mode&os.ModeSetgid != 0 {
			issues.RawSetInt(issueCount, lua.LString("SGID file: "+filePath))
			issueCount++
		}
		
		return nil
	})
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.SetField(result, "issues", issues)
	L.SetField(result, "total_issues", lua.LNumber(issueCount-1))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaFindSUIDFiles(L *lua.LState) int {
	path := L.OptString(1, "/")
	
	cmd := exec.Command("find", path, "-type", "f", "-perm", "-4000", "-o", "-perm", "-2000")
	output, err := cmd.Output()
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	for i, file := range files {
		if file != "" {
			result.RawSetInt(i+1, lua.LString(file))
		}
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaCheckFileIntegrity(L *lua.LState) int {
	path := L.CheckString(1)
	
	result := L.NewTable()
	
	// Basic file integrity checks
	if _, err := os.Stat(path); os.IsNotExist(err) {
		L.SetField(result, "exists", lua.LBool(false))
		L.Push(lua.LTrue)
		L.Push(result)
		return 2
	}
	
	L.SetField(result, "exists", lua.LBool(true))
	
	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.SetField(result, "size", lua.LNumber(info.Size()))
	L.SetField(result, "mode", lua.LString(info.Mode().String()))
	L.SetField(result, "modified", lua.LNumber(info.ModTime().Unix()))
	
	// Calculate file hash if it's a regular file
	if info.Mode().IsRegular() {
		cmd := exec.Command("sha256sum", path)
		output, err := cmd.Output()
		if err == nil {
			hash := strings.Fields(string(output))[0]
			L.SetField(result, "sha256", lua.LString(hash))
		}
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Password and authentication
func (s *SecurityModule) luaPasswordStrength(L *lua.LState) int {
	password := L.CheckString(1)
	
	result := L.NewTable()
	score := 0
	issues := L.NewTable()
	issueCount := 1
	
	// Length check
	if len(password) >= 12 {
		score += 25
	} else if len(password) >= 8 {
		score += 15
	} else {
		issues.RawSetInt(issueCount, lua.LString("Password too short (minimum 8 characters)"))
		issueCount++
	}
	
	// Character variety checks
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	
	variety := 0
	if hasLower {
		variety++
		score += 10
	} else {
		issues.RawSetInt(issueCount, lua.LString("Missing lowercase letters"))
		issueCount++
	}
	
	if hasUpper {
		variety++
		score += 10
	} else {
		issues.RawSetInt(issueCount, lua.LString("Missing uppercase letters"))
		issueCount++
	}
	
	if hasDigit {
		variety++
		score += 10
	} else {
		issues.RawSetInt(issueCount, lua.LString("Missing digits"))
		issueCount++
	}
	
	if hasSpecial {
		variety++
		score += 15
	} else {
		issues.RawSetInt(issueCount, lua.LString("Missing special characters"))
		issueCount++
	}
	
	// Bonus for variety
	if variety >= 3 {
		score += 10
	}
	if variety == 4 {
		score += 10
	}
	
	// Common password patterns
	commonPasswords := []string{"password", "123456", "qwerty", "admin", "letmein", "welcome"}
	for _, common := range commonPasswords {
		if strings.Contains(strings.ToLower(password), common) {
			score -= 30
			issues.RawSetInt(issueCount, lua.LString("Contains common password pattern"))
			issueCount++
			break
		}
	}
	
	// Determine strength
	var strength string
	if score >= 80 {
		strength = "strong"
	} else if score >= 60 {
		strength = "medium"
	} else if score >= 40 {
		strength = "weak"
	} else {
		strength = "very weak"
	}
	
	L.SetField(result, "score", lua.LNumber(score))
	L.SetField(result, "strength", lua.LString(strength))
	L.SetField(result, "issues", issues)
	L.SetField(result, "length", lua.LNumber(len(password)))
	L.SetField(result, "character_variety", lua.LNumber(variety))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaCheckWeakPasswords(L *lua.LState) int {
	_ = L.OptString(1, "/etc/shadow") // passwordFile unused in placeholder
	
	// This is a placeholder - in practice, you'd need appropriate permissions
	// and would check against known weak password hashes
	result := L.NewTable()
	L.SetField(result, "message", lua.LString("Weak password checking requires elevated privileges"))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// System security
func (s *SecurityModule) luaFirewallStatus(L *lua.LState) int {
	result := L.NewTable()
	
	// Check iptables
	cmd := exec.Command("iptables", "-L", "-n")
	output, err := cmd.Output()
	if err == nil {
		L.SetField(result, "iptables_active", lua.LBool(true))
		L.SetField(result, "iptables_rules", lua.LString(string(output)))
	} else {
		L.SetField(result, "iptables_active", lua.LBool(false))
	}
	
	// Check ufw (Ubuntu)
	cmd = exec.Command("ufw", "status")
	output, err = cmd.Output()
	if err == nil {
		status := string(output)
		L.SetField(result, "ufw_available", lua.LBool(true))
		L.SetField(result, "ufw_active", lua.LBool(strings.Contains(status, "Status: active")))
		L.SetField(result, "ufw_status", lua.LString(status))
	}
	
	// Check firewalld (RedHat/CentOS)
	cmd = exec.Command("firewall-cmd", "--state")
	output, err = cmd.Output()
	if err == nil {
		L.SetField(result, "firewalld_active", lua.LBool(strings.TrimSpace(string(output)) == "running"))
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaSELinuxStatus(L *lua.LState) int {
	result := L.NewTable()
	
	// Check SELinux status
	cmd := exec.Command("getenforce")
	output, err := cmd.Output()
	if err == nil {
		status := strings.TrimSpace(string(output))
		L.SetField(result, "available", lua.LBool(true))
		L.SetField(result, "status", lua.LString(status))
		L.SetField(result, "enforcing", lua.LBool(status == "Enforcing"))
	} else {
		L.SetField(result, "available", lua.LBool(false))
		L.SetField(result, "status", lua.LString("Not available"))
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaCheckUpdates(L *lua.LState) int {
	result := L.NewTable()
	
	// Check for package updates based on the system
	if _, err := exec.LookPath("apt"); err == nil {
		// Debian/Ubuntu
		cmd := exec.Command("apt", "list", "--upgradable")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			updates := L.NewTable()
			updateCount := 1
			
			for _, line := range lines {
				if strings.Contains(line, "upgradable") {
					updates.RawSetInt(updateCount, lua.LString(line))
					updateCount++
				}
			}
			
			L.SetField(result, "available_updates", updates)
			L.SetField(result, "update_count", lua.LNumber(updateCount-1))
		}
	} else if _, err := exec.LookPath("yum"); err == nil {
		// RedHat/CentOS
		cmd := exec.Command("yum", "check-update")
		output, err := cmd.Output()
		// yum check-update returns non-zero when updates are available
		L.SetField(result, "updates_available", lua.LBool(err != nil))
		L.SetField(result, "output", lua.LString(string(output)))
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Vulnerability scanning
func (s *SecurityModule) luaVulnerabilityScan(L *lua.LState) int {
	target := L.CheckString(1)
	scanType := L.OptString(2, "basic")
	
	result := L.NewTable()
	
	// Basic vulnerability checks
	if scanType == "basic" || scanType == "all" {
		// Check for common vulnerabilities
		vulns := L.NewTable()
		vulnCount := 1
		
		// Check SSH configuration
		if _, err := os.Stat("/etc/ssh/sshd_config"); err == nil {
			content, err := os.ReadFile("/etc/ssh/sshd_config")
			if err == nil {
				config := string(content)
				if strings.Contains(config, "PermitRootLogin yes") {
					vulns.RawSetInt(vulnCount, lua.LString("SSH root login enabled"))
					vulnCount++
				}
				if strings.Contains(config, "PasswordAuthentication yes") {
					vulns.RawSetInt(vulnCount, lua.LString("SSH password authentication enabled"))
					vulnCount++
				}
			}
		}
		
		L.SetField(result, "vulnerabilities", vulns)
		L.SetField(result, "vulnerability_count", lua.LNumber(vulnCount-1))
	}
	
	L.SetField(result, "scan_type", lua.LString(scanType))
	L.SetField(result, "target", lua.LString(target))
	L.SetField(result, "timestamp", lua.LNumber(time.Now().Unix()))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaCheckOpenPorts(L *lua.LState) int {
	// Check locally open ports
	cmd := exec.Command("netstat", "-tuln")
	output, err := cmd.Output()
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	ports := L.NewTable()
	portCount := 1
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "LISTEN") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				portInfo := L.NewTable()
				L.SetField(portInfo, "protocol", lua.LString(fields[0]))
				L.SetField(portInfo, "address", lua.LString(fields[3]))
				L.SetField(portInfo, "state", lua.LString(fields[5]))
				
				ports.RawSetInt(portCount, portInfo)
				portCount++
			}
		}
	}
	
	L.SetField(result, "listening_ports", ports)
	L.SetField(result, "total_ports", lua.LNumber(portCount-1))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaMalwareScan(L *lua.LState) int {
	path := L.CheckString(1)
	
	result := L.NewTable()
	
	// Basic malware detection patterns
	suspiciousFiles := L.NewTable()
	suspiciousCount := 1
	
	// Look for suspicious file patterns
	patterns := []string{
		"*.exe", "*.scr", "*.com", "*.pif", "*.vbs", "*.js",
	}
	
	for _, pattern := range patterns {
		cmd := exec.Command("find", path, "-name", pattern, "-type", "f")
		output, err := cmd.Output()
		if err == nil {
			files := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, file := range files {
				if file != "" {
					suspiciousFiles.RawSetInt(suspiciousCount, lua.LString(file))
					suspiciousCount++
				}
			}
		}
	}
	
	L.SetField(result, "suspicious_files", suspiciousFiles)
	L.SetField(result, "suspicious_count", lua.LNumber(suspiciousCount-1))
	L.SetField(result, "scan_path", lua.LString(path))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Compliance and hardening
func (s *SecurityModule) luaSecurityBaseline(L *lua.LState) int {
	result := L.NewTable()
	checks := L.NewTable()
	checkCount := 1
	passed := 0
	total := 0
	
	// System checks
	systemChecks := []struct {
		name string
		test func() bool
	}{
		{"Root login disabled", func() bool {
			content, err := os.ReadFile("/etc/ssh/sshd_config")
			if err != nil {
				return false
			}
			return !strings.Contains(string(content), "PermitRootLogin yes")
		}},
		{"Firewall active", func() bool {
			cmd := exec.Command("ufw", "status")
			output, _ := cmd.Output()
			return strings.Contains(string(output), "Status: active")
		}},
		{"Password complexity enabled", func() bool {
			_, err := os.Stat("/etc/security/pwquality.conf")
			return err == nil
		}},
	}
	
	for _, check := range systemChecks {
		total++
		result := check.test()
		if result {
			passed++
		}
		
		checkTable := L.NewTable()
		L.SetField(checkTable, "name", lua.LString(check.name))
		L.SetField(checkTable, "passed", lua.LBool(result))
		
		checks.RawSetInt(checkCount, checkTable)
		checkCount++
	}
	
	L.SetField(result, "checks", checks)
	L.SetField(result, "total_checks", lua.LNumber(total))
	L.SetField(result, "passed_checks", lua.LNumber(passed))
	L.SetField(result, "compliance_percentage", lua.LNumber(float64(passed)/float64(total)*100))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (s *SecurityModule) luaCISBenchmark(L *lua.LState) int {
	benchmark := L.OptString(1, "ubuntu")
	
	result := L.NewTable()
	L.SetField(result, "benchmark", lua.LString(benchmark))
	L.SetField(result, "message", lua.LString("CIS Benchmark implementation requires specialized tools"))
	
	// This would typically integrate with tools like:
	// - CIS-CAT
	// - OpenSCAP
	// - Custom compliance scripts
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Helper functions
func getServiceName(port int) string {
	services := map[int]string{
		21:   "ftp",
		22:   "ssh",
		23:   "telnet",
		25:   "smtp",
		53:   "dns",
		80:   "http",
		110:  "pop3",
		143:  "imap",
		443:  "https",
		993:  "imaps",
		995:  "pop3s",
		3306: "mysql",
		5432: "postgresql",
		6379: "redis",
		27017: "mongodb",
	}
	
	if service, exists := services[port]; exists {
		return service
	}
	return "unknown"
}