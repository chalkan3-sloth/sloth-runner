package luainterface

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// NetworkModule provides network utilities for Lua scripts
type NetworkModule struct{}

// NewNetworkModule creates a new network module
func NewNetworkModule() *NetworkModule {
	return &NetworkModule{}
}

// RegisterNetworkModule registers the network module with the Lua state
func RegisterNetworkModule(L *lua.LState) {
	module := NewNetworkModule()
	
	// Create the network table
	networkTable := L.NewTable()
	
	// Connectivity tests
	L.SetField(networkTable, "ping", L.NewFunction(module.luaPing))
	L.SetField(networkTable, "port_check", L.NewFunction(module.luaPortCheck))
	L.SetField(networkTable, "port_scan", L.NewFunction(module.luaPortScan))
	L.SetField(networkTable, "telnet", L.NewFunction(module.luaTelnet))
	
	// DNS operations
	L.SetField(networkTable, "dns_lookup", L.NewFunction(module.luaDNSLookup))
	L.SetField(networkTable, "reverse_dns", L.NewFunction(module.luaReverseDNS))
	L.SetField(networkTable, "mx_lookup", L.NewFunction(module.luaMXLookup))
	
	// Network information
	L.SetField(networkTable, "interfaces", L.NewFunction(module.luaInterfaces))
	L.SetField(networkTable, "public_ip", L.NewFunction(module.luaPublicIP))
	L.SetField(networkTable, "local_ip", L.NewFunction(module.luaLocalIP))
	
	// Network utilities
	L.SetField(networkTable, "traceroute", L.NewFunction(module.luaTraceroute))
	L.SetField(networkTable, "whois", L.NewFunction(module.luaWhois))
	L.SetField(networkTable, "ssl_check", L.NewFunction(module.luaSSLCheck))
	
	// Speed tests
	L.SetField(networkTable, "bandwidth_test", L.NewFunction(module.luaBandwidthTest))
	L.SetField(networkTable, "latency_test", L.NewFunction(module.luaLatencyTest))
	
	// Register the network table globally
	L.SetGlobal("network", networkTable)
}

// Connectivity tests
func (n *NetworkModule) luaPing(L *lua.LState) int {
	host := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	// Parse options
	count := 3
	timeout := 5
	
	if countVal := options.RawGetString("count"); countVal != lua.LNil {
		count = int(countVal.(lua.LNumber))
	}
	if timeoutVal := options.RawGetString("timeout"); timeoutVal != lua.LNil {
		timeout = int(timeoutVal.(lua.LNumber))
	}
	
	// Execute ping command
	cmd := exec.Command("ping", "-c", strconv.Itoa(count), "-W", strconv.Itoa(timeout*1000), host)
	output, err := cmd.Output()
	
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Parse ping results
	result := L.NewTable()
	lines := strings.Split(string(output), "\n")
	
	L.SetField(result, "success", lua.LBool(true))
	L.SetField(result, "output", lua.LString(string(output)))
	
	// Extract statistics if available
	for _, line := range lines {
		if strings.Contains(line, "packets transmitted") {
			L.SetField(result, "statistics", lua.LString(line))
		}
		if strings.Contains(line, "min/avg/max") {
			L.SetField(result, "timing", lua.LString(line))
		}
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (n *NetworkModule) luaPortCheck(L *lua.LState) int {
	host := L.CheckString(1)
	port := L.CheckInt(2)
	timeout := time.Duration(L.OptInt(3, 5)) * time.Second
	
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer conn.Close()
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString("port is open"))
	return 2
}

func (n *NetworkModule) luaPortScan(L *lua.LState) int {
	host := L.CheckString(1)
	portsArg := L.CheckAny(2)
	timeout := time.Duration(L.OptInt(3, 2)) * time.Second
	
	var ports []int
	
	// Parse ports argument
	if portsArg.Type() == lua.LTTable {
		portsTable := portsArg.(*lua.LTable)
		portsTable.ForEach(func(_, v lua.LValue) {
			if num, ok := v.(lua.LNumber); ok {
				ports = append(ports, int(num))
			}
		})
	} else if portsArg.Type() == lua.LTNumber {
		ports = append(ports, int(portsArg.(lua.LNumber)))
	}
	
	results := L.NewTable()
	openPorts := L.NewTable()
	closedPorts := L.NewTable()
	
	openCount := 1
	closedCount := 1
	
	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
		if err == nil {
			conn.Close()
			openPorts.RawSetInt(openCount, lua.LNumber(port))
			openCount++
		} else {
			closedPorts.RawSetInt(closedCount, lua.LNumber(port))
			closedCount++
		}
	}
	
	L.SetField(results, "open", openPorts)
	L.SetField(results, "closed", closedPorts)
	L.SetField(results, "total_scanned", lua.LNumber(len(ports)))
	
	L.Push(results)
	return 1
}

func (n *NetworkModule) luaTelnet(L *lua.LState) int {
	host := L.CheckString(1)
	port := L.CheckInt(2)
	timeout := time.Duration(L.OptInt(3, 10)) * time.Second
	
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer conn.Close()
	
	// Read initial response
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	
	result := L.NewTable()
	L.SetField(result, "connected", lua.LBool(true))
	
	if err == nil && bytesRead > 0 {
		L.SetField(result, "banner", lua.LString(string(buffer[:bytesRead])))
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// DNS operations
func (n *NetworkModule) luaDNSLookup(L *lua.LState) int {
	hostname := L.CheckString(1)
	recordType := L.OptString(2, "A")
	
	var ips []string
	var err error
	
	switch strings.ToUpper(recordType) {
	case "A":
		ips, err = net.LookupHost(hostname)
	case "CNAME":
		cname, err := net.LookupCNAME(hostname)
		if err == nil {
			ips = []string{cname}
		}
	case "TXT":
		txtRecords, err := net.LookupTXT(hostname)
		if err == nil {
			ips = txtRecords
		}
	default:
		ips, err = net.LookupHost(hostname)
	}
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	for i, ip := range ips {
		result.RawSetInt(i+1, lua.LString(ip))
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (n *NetworkModule) luaReverseDNS(L *lua.LState) int {
	ip := L.CheckString(1)
	
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	for i, host := range hosts {
		result.RawSetInt(i+1, lua.LString(host))
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (n *NetworkModule) luaMXLookup(L *lua.LState) int {
	domain := L.CheckString(1)
	
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	for i, mx := range mxRecords {
		record := L.NewTable()
		L.SetField(record, "host", lua.LString(mx.Host))
		L.SetField(record, "priority", lua.LNumber(mx.Pref))
		result.RawSetInt(i+1, record)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Network information
func (n *NetworkModule) luaInterfaces(L *lua.LState) int {
	interfaces, err := net.Interfaces()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	result := L.NewTable()
	
	for i, iface := range interfaces {
		ifaceTable := L.NewTable()
		L.SetField(ifaceTable, "name", lua.LString(iface.Name))
		L.SetField(ifaceTable, "mtu", lua.LNumber(iface.MTU))
		L.SetField(ifaceTable, "hardware_addr", lua.LString(iface.HardwareAddr.String()))
		L.SetField(ifaceTable, "flags", lua.LString(iface.Flags.String()))
		
		// Get addresses
		addrs, err := iface.Addrs()
		if err == nil {
			addrTable := L.NewTable()
			for j, addr := range addrs {
				addrTable.RawSetInt(j+1, lua.LString(addr.String()))
			}
			L.SetField(ifaceTable, "addresses", addrTable)
		}
		
		result.RawSetInt(i+1, ifaceTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (n *NetworkModule) luaPublicIP(L *lua.LState) int {
	// Try multiple services for reliability
	services := []string{
		"https://api.ipify.org",
		"https://icanhazip.com",
		"https://ifconfig.me",
	}
	
	for _, service := range services {
		cmd := exec.Command("curl", "-s", "--connect-timeout", "5", service)
		output, err := cmd.Output()
		if err == nil {
			ip := strings.TrimSpace(string(output))
			if net.ParseIP(ip) != nil {
				L.Push(lua.LString(ip))
				return 1
			}
		}
	}
	
	L.Push(lua.LNil)
	L.Push(lua.LString("failed to get public IP"))
	return 2
}

func (n *NetworkModule) luaLocalIP(L *lua.LState) int {
	// Get local IP by connecting to a remote address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer conn.Close()
	
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	L.Push(lua.LString(localAddr.IP.String()))
	return 1
}

// Network utilities
func (n *NetworkModule) luaTraceroute(L *lua.LState) int {
	host := L.CheckString(1)
	maxHops := L.OptInt(2, 30)
	
	cmd := exec.Command("traceroute", "-m", strconv.Itoa(maxHops), host)
	output, err := cmd.Output()
	
	result := L.NewTable()
	
	if err != nil {
		L.SetField(result, "success", lua.LBool(false))
		L.SetField(result, "error", lua.LString(err.Error()))
	} else {
		L.SetField(result, "success", lua.LBool(true))
		L.SetField(result, "output", lua.LString(string(output)))
		
		// Parse hops
		lines := strings.Split(string(output), "\n")
		hops := L.NewTable()
		hopCount := 1
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.Contains(line, "traceroute") {
				hops.RawSetInt(hopCount, lua.LString(line))
				hopCount++
			}
		}
		L.SetField(result, "hops", hops)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

func (n *NetworkModule) luaWhois(L *lua.LState) int {
	domain := L.CheckString(1)
	
	cmd := exec.Command("whois", domain)
	output, err := cmd.Output()
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(output)))
	return 1
}

func (n *NetworkModule) luaSSLCheck(L *lua.LState) int {
	url := L.CheckString(1)
	
	// Extract hostname from URL
	hostname := strings.TrimPrefix(url, "https://")
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.Split(hostname, "/")[0]
	hostname = strings.Split(hostname, ":")[0]
	
	cmd := exec.Command("openssl", "s_client", "-connect", hostname+":443", "-servername", hostname)
	cmd.Stdin = strings.NewReader("Q\n")
	output, err := cmd.Output()
	
	result := L.NewTable()
	
	if err != nil {
		L.SetField(result, "valid", lua.LBool(false))
		L.SetField(result, "error", lua.LString(err.Error()))
	} else {
		outputStr := string(output)
		L.SetField(result, "valid", lua.LBool(true))
		L.SetField(result, "output", lua.LString(outputStr))
		
		// Extract certificate info
		if strings.Contains(outputStr, "Verify return code: 0 (ok)") {
			L.SetField(result, "verified", lua.LBool(true))
		} else {
			L.SetField(result, "verified", lua.LBool(false))
		}
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// Speed tests
func (n *NetworkModule) luaBandwidthTest(L *lua.LState) int {
	server := L.OptString(1, "speedtest.net")
	
	cmd := exec.Command("curl", "-s", "-w", "@-", "-o", "/dev/null", "http://"+server)
	cmd.Stdin = strings.NewReader("%{speed_download}")
	
	output, err := cmd.Output()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	speed := strings.TrimSpace(string(output))
	L.Push(lua.LString(speed + " bytes/sec"))
	return 1
}

func (n *NetworkModule) luaLatencyTest(L *lua.LState) int {
	host := L.CheckString(1)
	count := L.OptInt(2, 5)
	
	var totalTime time.Duration
	var successCount int
	
	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", host+":80", 5*time.Second)
		if err == nil {
			conn.Close()
			totalTime += time.Since(start)
			successCount++
		}
	}
	
	result := L.NewTable()
	
	if successCount > 0 {
		avgLatency := totalTime / time.Duration(successCount)
		L.SetField(result, "average_ms", lua.LNumber(avgLatency.Milliseconds()))
		L.SetField(result, "success_rate", lua.LNumber(float64(successCount)/float64(count)*100))
	} else {
		L.SetField(result, "average_ms", lua.LNumber(0))
		L.SetField(result, "success_rate", lua.LNumber(0))
	}
	
	L.SetField(result, "total_tests", lua.LNumber(count))
	L.SetField(result, "successful_tests", lua.LNumber(successCount))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}