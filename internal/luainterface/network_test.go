package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestNetworkModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterNetworkModule(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "network.hostname()",
			script: `
				local hostname = network.hostname()
				assert(type(hostname) == "string", "hostname() should return a string")
				assert(#hostname > 0, "hostname() should return non-empty string")
			`,
		},
		{
			name: "network.lookup_ip()",
			script: `
				local ips = network.lookup_ip("localhost")
				assert(type(ips) == "table", "lookup_ip() should return a table")
			`,
		},
		{
			name: "network.lookup_host()",
			script: `
				local addrs = network.lookup_host("localhost")
				assert(type(addrs) == "table", "lookup_host() should return a table")
			`,
		},
		{
			name: "network.resolve_tcp()",
			script: `
				local addr = network.resolve_tcp("localhost:80")
				assert(type(addr) == "string" or type(addr) == "nil", "resolve_tcp() should return string or nil")
			`,
		},
		{
			name: "network.get_interfaces()",
			script: `
				local ifaces = network.get_interfaces()
				assert(type(ifaces) == "table", "get_interfaces() should return a table")
			`,
		},
		{
			name: "network.is_port_open()",
			script: `
				-- This will likely fail, but should not error
				local open = network.is_port_open("localhost", 99999)
				assert(type(open) == "boolean", "is_port_open() should return a boolean")
			`,
		},
		{
			name: "network.get_public_ip()",
			script: `
				-- May fail if no internet, but should not crash
				local ip = network.get_public_ip()
				assert(type(ip) == "string" or ip == nil, "get_public_ip() should return string or nil")
			`,
		},
		{
			name: "network.get_local_ips()",
			script: `
				local ips = network.get_local_ips()
				assert(type(ips) == "table", "get_local_ips() should return a table")
			`,
		},
		{
			name: "network.parse_url()",
			script: `
				local url = network.parse_url("http://example.com:8080/path?query=value")
				assert(type(url) == "table", "parse_url() should return a table")
				assert(url.scheme == "http", "scheme should be 'http'")
				assert(url.host == "example.com:8080", "host should be 'example.com:8080'")
				assert(url.path == "/path", "path should be '/path'")
			`,
		},
		{
			name: "network.parse_cidr()",
			script: `
				local cidr = network.parse_cidr("192.168.1.0/24")
				assert(type(cidr) == "table", "parse_cidr() should return a table")
			`,
		},
		{
			name: "network.ip_in_cidr()",
			script: `
				local result = network.ip_in_cidr("192.168.1.100", "192.168.1.0/24")
				assert(type(result) == "boolean", "ip_in_cidr() should return a boolean")
			`,
		},
		{
			name: "network.is_ipv4()",
			script: `
				assert(network.is_ipv4("192.168.1.1") == true, "should detect valid IPv4")
				assert(network.is_ipv4("invalid") == false, "should reject invalid IPv4")
			`,
		},
		{
			name: "network.is_ipv6()",
			script: `
				assert(network.is_ipv6("::1") == true, "should detect valid IPv6")
				assert(network.is_ipv6("192.168.1.1") == false, "should reject IPv4 as IPv6")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestNetworkURLParsing(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterNetworkModule(L)

	script := `
		local url1 = network.parse_url("https://user:pass@example.com:443/path?q=1#frag")
		assert(url1.scheme == "https", "scheme should be 'https'")
		assert(url1.host == "example.com:443", "host should include port")
		assert(url1.path == "/path", "path should be '/path'")
		
		local url2 = network.parse_url("ftp://ftp.example.com/file.txt")
		assert(url2.scheme == "ftp", "scheme should be 'ftp'")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestNetworkIPValidation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterNetworkModule(L)

	script := `
		-- IPv4 tests
		assert(network.is_ipv4("0.0.0.0") == true, "0.0.0.0 is valid IPv4")
		assert(network.is_ipv4("255.255.255.255") == true, "255.255.255.255 is valid IPv4")
		assert(network.is_ipv4("256.1.1.1") == false, "256.1.1.1 is invalid IPv4")
		assert(network.is_ipv4("1.1.1") == false, "1.1.1 is invalid IPv4")
		
		-- IPv6 tests
		assert(network.is_ipv6("::1") == true, "::1 is valid IPv6")
		assert(network.is_ipv6("2001:db8::1") == true, "2001:db8::1 is valid IPv6")
		assert(network.is_ipv6("gggg::1") == false, "gggg::1 is invalid IPv6")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestNetworkCIDR(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterNetworkModule(L)

	script := `
		-- Test CIDR matching
		assert(network.ip_in_cidr("192.168.1.50", "192.168.1.0/24") == true, "IP should be in CIDR range")
		assert(network.ip_in_cidr("192.168.2.50", "192.168.1.0/24") == false, "IP should not be in CIDR range")
		assert(network.ip_in_cidr("10.0.0.1", "10.0.0.0/8") == true, "IP should be in large CIDR range")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
