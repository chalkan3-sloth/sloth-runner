-- Security Module Examples

print("🔐 SECURITY MODULE SHOWCASE")
print("=" .. string.rep("=", 40))

-- 1. Network Security
print("\n🌐 Network Security Scanning:")

-- Port scanning
local port_scan = security.scan_ports("127.0.0.1", {
    start = 20, 
    ["end"] = 100, 
    timeout = 1
})

if port_scan then
    print("🔍 Port scan results for localhost:")
    print("   Open ports:", port_scan.open_count or 0)
    print("   Total scanned:", port_scan.total_scanned or 0)
    
    if port_scan.open and #port_scan.open > 0 then
        print("   Open services found:")
        for i = 1, math.min(#port_scan.open, 3) do
            local port_info = port_scan.open[i]
            print("     - Port " .. port_info.port .. " (" .. port_info.service .. ")")
        end
    end
end

-- SSL Certificate check
local ssl_check = security.check_ssl_cert("https://google.com")
if ssl_check then
    if ssl_check.valid then
        print("🔒 SSL Certificate Analysis:")
        print("   Domain: google.com")
        print("   Valid:", ssl_check.verified and "Yes" or "No")
        print("   Expires in:", math.floor(ssl_check.expires_in_days or 0) .. " days")
        
        if ssl_check.issues and #ssl_check.issues > 0 then
            print("   Security Issues:")
            for i = 1, #ssl_check.issues do
                print("     ⚠️ " .. ssl_check.issues[i])
            end
        else
            print("   ✅ No SSL issues detected")
        end
    else
        print("❌ SSL Certificate check failed:", ssl_check.error or "Unknown error")
    end
end

-- HTTP Security Headers
print("\n🛡️ HTTP Security Headers:")
print("   Use security.check_http_headers(url) to analyze web security")
print("   Checks for HSTS, CSP, X-Frame-Options, etc.")

-- 2. File System Security
print("\n📁 File System Security:")

-- Check current directory permissions
local audit_result = security.audit_permissions(".")
if audit_result then
    print("🔍 Permission audit of current directory:")
    print("   Issues found:", audit_result.total_issues or 0)
    
    if audit_result.issues and #audit_result.issues > 0 then
        print("   Security concerns:")
        for i = 1, math.min(#audit_result.issues, 3) do
            print("     ⚠️ " .. audit_result.issues[i])
        end
    else
        print("   ✅ No permission issues in current directory")
    end
end

-- Find SUID files (may require permissions)
print("🔒 SUID/SGID file detection available")
print("   Use security.find_suid_files(path) to scan for privileged files")

-- File integrity check
local integrity = security.check_file_integrity("/etc/hosts")
if integrity then
    print("🧾 File integrity check (/etc/hosts):")
    print("   Exists:", integrity.exists and "Yes" or "No")
    if integrity.exists then
        print("   Size:", integrity.size or 0, "bytes")
        print("   Modified:", os.date("%Y-%m-%d %H:%M:%S", integrity.modified or 0))
        if integrity.sha256 then
            print("   SHA256:", string.sub(integrity.sha256, 1, 16) .. "...")
        end
    end
end

-- 3. Password Security
print("\n🔑 Password Security:")

-- Test various password strengths
local passwords = {
    {pwd = "123456", desc = "Weak password"},
    {pwd = "password123", desc = "Common pattern"},
    {pwd = "MyStr0ng@Pass!", desc = "Strong password"},
    {pwd = "Sup3r$3cur3P@ssw0rd2024!", desc = "Very strong password"}
}

for _, test in ipairs(passwords) do
    local strength = security.password_strength(test.pwd)
    if strength then
        print("🔐 " .. test.desc .. ":")
        print("   Score: " .. (strength.score or 0) .. "/100")
        print("   Strength: " .. (strength.strength or "unknown"))
        print("   Length: " .. (strength.length or 0) .. " characters")
        print("   Character variety: " .. (strength.character_variety or 0) .. "/4")
        
        if strength.issues and #strength.issues > 0 then
            print("   Issues: " .. #strength.issues .. " found")
        end
    end
    print()
end

-- 4. System Security
print("\n🖥️ System Security Status:")

-- Firewall status
local firewall = security.firewall_status()
if firewall then
    print("🔥 Firewall Status:")
    print("   iptables active:", firewall.iptables_active and "Yes" or "No")
    
    if firewall.ufw_available then
        print("   UFW active:", firewall.ufw_active and "Yes" or "No")
    end
    
    if firewall.firewalld_active ~= nil then
        print("   firewalld active:", firewall.firewalld_active and "Yes" or "No")
    end
end

-- SELinux status
local selinux = security.selinux_status()
if selinux then
    print("🛡️ SELinux Status:")
    print("   Available:", selinux.available and "Yes" or "No")
    if selinux.available then
        print("   Status:", selinux.status or "Unknown")
        print("   Enforcing:", selinux.enforcing and "Yes" or "No")
    end
end

-- 5. Vulnerability Assessment
print("\n🔍 Vulnerability Assessment:")

-- Basic vulnerability scan
local vuln_scan = security.vulnerability_scan("localhost", "basic")
if vuln_scan then
    print("🦠 Vulnerability scan results:")
    print("   Target:", vuln_scan.target or "Unknown")
    print("   Scan type:", vuln_scan.scan_type or "Unknown")
    print("   Vulnerabilities found:", vuln_scan.vulnerability_count or 0)
    
    if vuln_scan.vulnerabilities and #vuln_scan.vulnerabilities > 0 then
        print("   Issues identified:")
        for i = 1, #vuln_scan.vulnerabilities do
            print("     ⚠️ " .. vuln_scan.vulnerabilities[i])
        end
    else
        print("   ✅ No vulnerabilities detected in basic scan")
    end
end

-- Check open ports on system
local open_ports = security.check_open_ports()
if open_ports then
    print("🔌 Open ports on system:")
    print("   Listening ports found:", open_ports.total_ports or 0)
    
    if open_ports.listening_ports and #open_ports.listening_ports > 0 then
        print("   Services listening:")
        for i = 1, math.min(#open_ports.listening_ports, 3) do
            local port = open_ports.listening_ports[i]
            print("     - " .. (port.protocol or "tcp") .. " " .. (port.address or "unknown"))
        end
    end
end

-- 6. Security Baseline Check
print("\n📊 Security Baseline Assessment:")

local baseline = security.security_baseline()
if baseline then
    print("🎯 Security compliance check:")
    print("   Total checks:", baseline.total_checks or 0)
    print("   Passed checks:", baseline.passed_checks or 0)
    print("   Compliance rate:", string.format("%.1f%%", baseline.compliance_percentage or 0))
    
    if baseline.checks and #baseline.checks > 0 then
        print("   Check results:")
        for i = 1, math.min(#baseline.checks, 3) do
            local check = baseline.checks[i]
            local status = check.passed and "✅" or "❌"
            print("     " .. status .. " " .. (check.name or "Unknown check"))
        end
    end
end

-- 7. Malware Detection
print("\n🦠 Malware Detection:")

-- Basic malware scan of current directory
local malware_scan = security.malware_scan(".")
if malware_scan then
    print("🔍 Malware scan of current directory:")
    print("   Suspicious files found:", malware_scan.suspicious_count or 0)
    print("   Scan path:", malware_scan.scan_path or "Unknown")
    
    if malware_scan.suspicious_files and #malware_scan.suspicious_files > 0 then
        print("   Suspicious files:")
        for i = 1, math.min(#malware_scan.suspicious_files, 3) do
            print("     ⚠️ " .. malware_scan.suspicious_files[i])
        end
    else
        print("   ✅ No suspicious files detected")
    end
end

-- 8. Advanced Security Features
print("\n🚀 Advanced Security Features:")

print("🔬 Advanced capabilities available:")
print("   • CIS Benchmark compliance checking")
print("   • Custom security rule validation")
print("   • Integration with external security tools")
print("   • Automated security reporting")
print("   • Security event correlation")

-- Security recommendations
print("\n📋 Security Recommendations:")
print("🔐 Security best practices:")
print("   • Enable and configure firewall")
print("   • Use strong passwords with complexity requirements")
print("   • Keep system and packages updated")
print("   • Monitor file integrity of critical files")
print("   • Regular vulnerability assessments")
print("   • Implement proper access controls")
print("   • Enable security auditing and logging")

print("\n✅ Security module demonstration completed!")
print("🛡️ Comprehensive security scanning and compliance tools ready!")