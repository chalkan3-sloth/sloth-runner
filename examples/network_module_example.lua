-- Network Module Examples

print("🌐 NETWORK MODULE SHOWCASE")
print("=" .. string.rep("=", 40))

-- 1. Connectivity Tests
print("\n📡 Connectivity Tests:")

-- Ping test
local ping_result = network.ping("google.com", {count = 3, timeout = 5})
if ping_result.success then
    print("✅ Ping to google.com successful")
    print("   Statistics:", ping_result.statistics or "N/A")
else
    print("❌ Ping failed:", ping_result.output)
end

-- Port check
local port_open, msg = network.port_check("google.com", 80, 3)
print("🔌 Port 80 on google.com:", port_open and "OPEN" or "CLOSED")

-- Port scan
local scan_result = network.port_scan("127.0.0.1", {22, 80, 443, 8080})
print("🔍 Port scan results:")
print("   Open ports:", #scan_result.open)
print("   Total scanned:", scan_result.total_scanned)

-- 2. DNS Operations
print("\n🌐 DNS Operations:")

-- DNS lookup
local ips = network.dns_lookup("google.com", "A")
if ips then
    print("🔍 Google.com IP addresses:")
    for i = 1, #ips do
        print("   -", ips[i])
    end
end

-- MX lookup
local mx_records = network.mx_lookup("google.com")
if mx_records then
    print("📧 Google.com MX records:")
    for i = 1, #mx_records do
        print("   - " .. mx_records[i].host .. " (priority: " .. mx_records[i].priority .. ")")
    end
end

-- Reverse DNS
local reverse = network.reverse_dns("8.8.8.8")
if reverse then
    print("🔄 Reverse DNS for 8.8.8.8:", reverse[1] or "No PTR record")
end

-- 3. Network Information
print("\n💻 Network Information:")

-- Network interfaces
local interfaces = network.interfaces()
if interfaces then
    print("🔗 Network interfaces found:", #interfaces)
    for i = 1, math.min(#interfaces, 3) do  -- Show first 3
        local iface = interfaces[i]
        print("   - " .. iface.name .. " (" .. iface.hardware_addr .. ")")
        if iface.addresses then
            for j = 1, math.min(#iface.addresses, 2) do
                print("     Address: " .. iface.addresses[j])
            end
        end
    end
end

-- Local IP
local local_ip = network.local_ip()
if local_ip then
    print("🏠 Local IP address:", local_ip)
end

-- Public IP (commented out to avoid external calls in example)
-- local public_ip = network.public_ip()
-- if public_ip then
--     print("🌍 Public IP address:", public_ip)
-- end

-- 4. Network Utilities
print("\n🛠️ Network Utilities:")

-- SSL Certificate check
local ssl_check = network.ssl_check("https://google.com")
if ssl_check.valid then
    print("🔒 SSL Certificate for google.com:")
    print("   Valid:", ssl_check.verified and "Yes" or "No")
    print("   Expires in days:", math.floor(ssl_check.expires_in_days or 0))
    if ssl_check.issues and #ssl_check.issues > 0 then
        print("   Issues found:", #ssl_check.issues)
    end
end

-- Latency test
local latency = network.latency_test("google.com", 3)
print("⚡ Latency test to google.com:")
print("   Average latency:", latency.average_ms .. "ms")
print("   Success rate:", latency.success_rate .. "%")

-- 5. Advanced Features
print("\n🚀 Advanced Features:")

-- Traceroute (may take time, commented for example)
-- local traceroute = network.traceroute("8.8.8.8", 10)
-- if traceroute.success then
--     print("🗺️ Traceroute to 8.8.8.8 completed")
--     print("   Hops found:", #traceroute.hops)
-- end

-- WHOIS lookup (may be slow, commented for example)
-- local whois_info = network.whois("google.com")
-- if whois_info then
--     print("📋 WHOIS information available for google.com")
-- end

-- HTTP headers security check
print("🛡️ Security Features:")
print("   HTTP headers check, SSL validation, and security scanning available")
print("   Use network.check_http_headers(url) for web security analysis")

print("\n✅ Network module demonstration completed!")
print("🔧 All network diagnostics and connectivity tools are ready to use.")