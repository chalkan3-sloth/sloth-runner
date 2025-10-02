package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/agent"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	lua "github.com/yuin/gopher-lua"
)

// FactsModule provides access to agent collected facts/system information
type FactsModule struct {
	masterAddr string
}

// NewFactsModule creates a new facts module instance
func NewFactsModule(masterAddr string) *FactsModule {
	return &FactsModule{
		masterAddr: masterAddr,
	}
}

// getAgentFacts retrieves facts for an agent from the master
func (m *FactsModule) getAgentFacts(agentName string) (*agent.SystemInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.Dial(m.masterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master: %w", err)
	}
	defer conn.Close()

	registryClient := pb.NewAgentRegistryClient(conn)
	resp, err := registryClient.GetAgentInfo(ctx, &pb.GetAgentInfoRequest{
		AgentName: agentName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get agent info: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get agent info: %s", resp.Message)
	}

	agentInfo := resp.GetAgentInfo()
	if agentInfo.GetSystemInfoJson() == "" {
		return nil, fmt.Errorf("no system info available for agent: %s", agentName)
	}

	var sysInfo agent.SystemInfo
	if err := json.Unmarshal([]byte(agentInfo.GetSystemInfoJson()), &sysInfo); err != nil {
		return nil, fmt.Errorf("failed to parse system info: %w", err)
	}

	return &sysInfo, nil
}

// Register registers the facts module with the Lua state
func (m *FactsModule) Register(L *lua.LState) {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get_all":       m.luaGetAll,
		"get_hostname":  m.luaGetHostname,
		"get_platform":  m.luaGetPlatform,
		"get_memory":    m.luaGetMemory,
		"get_disk":      m.luaGetDisk,
		"get_network":   m.luaGetNetwork,
		"get_packages":  m.luaGetPackages,
		"get_package":   m.luaGetPackage,
		"get_services":  m.luaGetServices,
		"get_service":   m.luaGetService,
		"get_users":     m.luaGetUsers,
		"get_user":      m.luaGetUser,
		"get_processes": m.luaGetProcesses,
		"get_mounts":    m.luaGetMounts,
		"get_uptime":    m.luaGetUptime,
		"get_load":      m.luaGetLoad,
		"get_kernel":    m.luaGetKernel,
		"query":         m.luaQuery,
	})

	L.SetGlobal("facts", mod)
}

// luaGetAll retrieves all facts from an agent
// Usage: facts.get_all({ agent = "agent-name" })
func (m *FactsModule) luaGetAll(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	agentName := getStringField(L, opts, "agent", "")
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Convert facts to Lua table
	result := m.factsToLuaTable(L, facts)
	L.Push(result)
	return 1
}

// luaGetHostname gets the hostname
// Usage: facts.get_hostname({ agent = "agent-name" })
func (m *FactsModule) luaGetHostname(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(facts.Hostname))
	return 1
}

// luaGetPlatform gets platform information
// Usage: facts.get_platform({ agent = "agent-name" })
func (m *FactsModule) luaGetPlatform(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	platformTable := L.NewTable()
	platformTable.RawSetString("os", lua.LString(facts.Platform))
	platformTable.RawSetString("family", lua.LString(facts.PlatformFamily))
	platformTable.RawSetString("version", lua.LString(facts.PlatformVersion))
	platformTable.RawSetString("architecture", lua.LString(facts.Architecture))
	platformTable.RawSetString("kernel", lua.LString(facts.Kernel))
	platformTable.RawSetString("kernel_version", lua.LString(facts.KernelVersion))
	platformTable.RawSetString("virtualization", lua.LString(facts.Virtualization))

	L.Push(platformTable)
	return 1
}

// luaGetMemory gets memory information
// Usage: facts.get_memory({ agent = "agent-name" })
func (m *FactsModule) luaGetMemory(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if facts.Memory == nil {
		L.Push(lua.LNil)
		return 1
	}

	memTable := L.NewTable()
	memTable.RawSetString("total", lua.LNumber(facts.Memory.Total))
	memTable.RawSetString("available", lua.LNumber(facts.Memory.Available))
	memTable.RawSetString("used", lua.LNumber(facts.Memory.Used))
	memTable.RawSetString("used_percent", lua.LNumber(facts.Memory.UsedPercent))
	memTable.RawSetString("free", lua.LNumber(facts.Memory.Free))
	memTable.RawSetString("cached", lua.LNumber(facts.Memory.Cached))
	memTable.RawSetString("buffers", lua.LNumber(facts.Memory.Buffers))

	L.Push(memTable)
	return 1
}

// luaGetDisk gets disk information
// Usage: facts.get_disk({ agent = "agent-name", mountpoint = "/home" })
func (m *FactsModule) luaGetDisk(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	mountpoint := getStringField(L, opts, "mountpoint", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if facts.Disk == nil {
		L.Push(L.NewTable())
		return 1
	}

	// If specific mountpoint requested
	if mountpoint != "" {
		for _, disk := range facts.Disk {
			if disk.Mountpoint == mountpoint {
				diskTable := m.diskToLuaTable(L, disk)
				L.Push(diskTable)
				return 1
			}
		}
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("mountpoint %s not found", mountpoint)))
		return 2
	}

	// Return all disks
	disksArray := L.NewTable()
	for _, disk := range facts.Disk {
		diskTable := m.diskToLuaTable(L, disk)
		disksArray.Append(diskTable)
	}

	L.Push(disksArray)
	return 1
}

// luaGetNetwork gets network information
// Usage: facts.get_network({ agent = "agent-name", interface = "eth0" })
func (m *FactsModule) luaGetNetwork(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	interfaceName := getStringField(L, opts, "interface", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if facts.Network == nil {
		L.Push(L.NewTable())
		return 1
	}

	// If specific interface requested
	if interfaceName != "" {
		for _, net := range facts.Network {
			if net.Name == interfaceName {
				netTable := m.networkToLuaTable(L, net)
				L.Push(netTable)
				return 1
			}
		}
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("interface %s not found", interfaceName)))
		return 2
	}

	// Return all interfaces
	networksArray := L.NewTable()
	for _, net := range facts.Network {
		netTable := m.networkToLuaTable(L, net)
		networksArray.Append(netTable)
	}

	L.Push(networksArray)
	return 1
}

// luaGetPackages gets installed packages
// Usage: facts.get_packages({ agent = "agent-name" })
func (m *FactsModule) luaGetPackages(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if facts.Packages == nil {
		L.Push(L.NewTable())
		return 1
	}

	packagesTable := L.NewTable()
	packagesTable.RawSetString("manager", lua.LString(facts.Packages.Manager))
	packagesTable.RawSetString("installed_count", lua.LNumber(facts.Packages.InstalledCount))
	packagesTable.RawSetString("updates_available", lua.LNumber(facts.Packages.UpdatesAvailable))
	
	// Packages list
	pkgList := L.NewTable()
	for _, pkg := range facts.Packages.Packages {
		pkgTable := L.NewTable()
		pkgTable.RawSetString("name", lua.LString(pkg.Name))
		pkgTable.RawSetString("version", lua.LString(pkg.Version))
		pkgTable.RawSetString("architecture", lua.LString(pkg.Architecture))
		pkgTable.RawSetString("description", lua.LString(pkg.Description))
		pkgList.Append(pkgTable)
	}
	packagesTable.RawSetString("packages", pkgList)
	
	// Updates list
	updatesList := L.NewTable()
	for _, upd := range facts.Packages.Updates {
		updTable := L.NewTable()
		updTable.RawSetString("name", lua.LString(upd.Name))
		updTable.RawSetString("version", lua.LString(upd.Version))
		updTable.RawSetString("architecture", lua.LString(upd.Architecture))
		updTable.RawSetString("description", lua.LString(upd.Description))
		updatesList.Append(updTable)
	}
	packagesTable.RawSetString("updates", updatesList)

	L.Push(packagesTable)
	return 1
}

// luaGetPackage checks if a specific package is installed
// Usage: facts.get_package({ agent = "agent-name", name = "nginx" })
func (m *FactsModule) luaGetPackage(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	pkgName := getStringField(L, opts, "name", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}
	
	if pkgName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("package name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if facts.Packages == nil {
		L.Push(lua.LNil)
		return 1
	}

	// Search for package
	for _, pkg := range facts.Packages.Packages {
		if pkg.Name == pkgName {
			pkgTable := L.NewTable()
			pkgTable.RawSetString("name", lua.LString(pkg.Name))
			pkgTable.RawSetString("version", lua.LString(pkg.Version))
			pkgTable.RawSetString("architecture", lua.LString(pkg.Architecture))
			pkgTable.RawSetString("description", lua.LString(pkg.Description))
			pkgTable.RawSetString("installed", lua.LTrue)
			L.Push(pkgTable)
			return 1
		}
	}

	// Package not found
	notFoundTable := L.NewTable()
	notFoundTable.RawSetString("name", lua.LString(pkgName))
	notFoundTable.RawSetString("installed", lua.LFalse)
	L.Push(notFoundTable)
	return 1
}

// luaGetServices gets all services
// Usage: facts.get_services({ agent = "agent-name" })
func (m *FactsModule) luaGetServices(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	servicesArray := L.NewTable()
	for _, svc := range facts.Services {
		svcTable := L.NewTable()
		svcTable.RawSetString("name", lua.LString(svc.Name))
		svcTable.RawSetString("status", lua.LString(svc.Status))
		svcTable.RawSetString("state", lua.LString(svc.State))
		servicesArray.Append(svcTable)
	}

	L.Push(servicesArray)
	return 1
}

// luaGetService gets a specific service status
// Usage: facts.get_service({ agent = "agent-name", name = "nginx" })
func (m *FactsModule) luaGetService(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	svcName := getStringField(L, opts, "name", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}
	
	if svcName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("service name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Search for service
	for _, svc := range facts.Services {
		if svc.Name == svcName {
			svcTable := L.NewTable()
			svcTable.RawSetString("name", lua.LString(svc.Name))
			svcTable.RawSetString("status", lua.LString(svc.Status))
			svcTable.RawSetString("state", lua.LString(svc.State))
			L.Push(svcTable)
			return 1
		}
	}

	L.Push(lua.LNil)
	L.Push(lua.LString(fmt.Sprintf("service %s not found", svcName)))
	return 2
}

// luaGetUsers gets all users
// Usage: facts.get_users({ agent = "agent-name" })
func (m *FactsModule) luaGetUsers(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	usersArray := L.NewTable()
	for _, user := range facts.Users {
		userTable := L.NewTable()
		userTable.RawSetString("username", lua.LString(user.Username))
		userTable.RawSetString("uid", lua.LString(user.UID))
		userTable.RawSetString("gid", lua.LString(user.GID))
		userTable.RawSetString("home", lua.LString(user.Home))
		userTable.RawSetString("shell", lua.LString(user.Shell))
		usersArray.Append(userTable)
	}

	L.Push(usersArray)
	return 1
}

// luaGetUser gets a specific user
// Usage: facts.get_user({ agent = "agent-name", username = "root" })
func (m *FactsModule) luaGetUser(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	username := getStringField(L, opts, "username", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}
	
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("username is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Search for user
	for _, user := range facts.Users {
		if user.Username == username {
			userTable := L.NewTable()
			userTable.RawSetString("username", lua.LString(user.Username))
			userTable.RawSetString("uid", lua.LString(user.UID))
			userTable.RawSetString("gid", lua.LString(user.GID))
			userTable.RawSetString("home", lua.LString(user.Home))
			userTable.RawSetString("shell", lua.LString(user.Shell))
			L.Push(userTable)
			return 1
		}
	}

	L.Push(lua.LNil)
	L.Push(lua.LString(fmt.Sprintf("user %s not found", username)))
	return 2
}

// luaGetProcesses gets process information
// Usage: facts.get_processes({ agent = "agent-name" })
func (m *FactsModule) luaGetProcesses(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if facts.Processes == nil {
		L.Push(lua.LNil)
		return 1
	}

	procTable := L.NewTable()
	procTable.RawSetString("total", lua.LNumber(facts.Processes.Total))
	procTable.RawSetString("running", lua.LNumber(facts.Processes.Running))
	procTable.RawSetString("sleeping", lua.LNumber(facts.Processes.Sleeping))
	procTable.RawSetString("zombie", lua.LNumber(facts.Processes.Zombie))

	L.Push(procTable)
	return 1
}

// luaGetMounts gets mount information
// Usage: facts.get_mounts({ agent = "agent-name" })
func (m *FactsModule) luaGetMounts(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	mountsArray := L.NewTable()
	for _, mount := range facts.Mounts {
		mountTable := L.NewTable()
		mountTable.RawSetString("device", lua.LString(mount.Device))
		mountTable.RawSetString("mountpoint", lua.LString(mount.Mountpoint))
		mountTable.RawSetString("fstype", lua.LString(mount.FSType))
		mountTable.RawSetString("options", lua.LString(mount.Options))
		mountsArray.Append(mountTable)
	}

	L.Push(mountsArray)
	return 1
}

// luaGetUptime gets system uptime
// Usage: facts.get_uptime({ agent = "agent-name" })
func (m *FactsModule) luaGetUptime(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	uptimeTable := L.NewTable()
	uptimeTable.RawSetString("seconds", lua.LNumber(facts.Uptime))
	uptimeTable.RawSetString("boot_time", lua.LNumber(facts.BootTime))
	uptimeTable.RawSetString("timezone", lua.LString(facts.Timezone))

	L.Push(uptimeTable)
	return 1
}

// luaGetLoad gets load average
// Usage: facts.get_load({ agent = "agent-name" })
func (m *FactsModule) luaGetLoad(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	loadArray := L.NewTable()
	for _, load := range facts.LoadAverage {
		loadArray.Append(lua.LNumber(load))
	}

	L.Push(loadArray)
	return 1
}

// luaGetKernel gets kernel information
// Usage: facts.get_kernel({ agent = "agent-name" })
func (m *FactsModule) luaGetKernel(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	kernelTable := L.NewTable()
	kernelTable.RawSetString("name", lua.LString(facts.Kernel))
	kernelTable.RawSetString("version", lua.LString(facts.KernelVersion))

	L.Push(kernelTable)
	return 1
}

// luaQuery performs a JSON path query on facts
// Usage: facts.query({ agent = "agent-name", path = "$.memory.total" })
func (m *FactsModule) luaQuery(L *lua.LState) int {
	opts := L.CheckTable(1)
	agentName := getStringField(L, opts, "agent", "")
	path := getStringField(L, opts, "path", "")
	
	if agentName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("agent name is required"))
		return 2
	}
	
	if path == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("query path is required"))
		return 2
	}

	facts, err := m.getAgentFacts(agentName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Convert facts to JSON for querying
	jsonData, err := json.Marshal(facts)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to marshal facts: %v", err)))
		return 2
	}

	// Simple path parsing (basic implementation)
	// For production, consider using a proper JSON path library
	var result interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to unmarshal facts: %v", err)))
		return 2
	}

	// Return raw JSON for now - TODO: implement proper JSON path query
	L.Push(lua.LString(string(jsonData)))
	return 1
}

// Helper functions

func (m *FactsModule) factsToLuaTable(L *lua.LState, facts *agent.SystemInfo) *lua.LTable {
	table := L.NewTable()
	
	table.RawSetString("hostname", lua.LString(facts.Hostname))
	table.RawSetString("platform", lua.LString(facts.Platform))
	table.RawSetString("platform_family", lua.LString(facts.PlatformFamily))
	table.RawSetString("platform_version", lua.LString(facts.PlatformVersion))
	table.RawSetString("architecture", lua.LString(facts.Architecture))
	table.RawSetString("cpus", lua.LNumber(facts.CPUs))
	table.RawSetString("kernel", lua.LString(facts.Kernel))
	table.RawSetString("kernel_version", lua.LString(facts.KernelVersion))
	table.RawSetString("virtualization", lua.LString(facts.Virtualization))
	table.RawSetString("uptime", lua.LNumber(facts.Uptime))
	table.RawSetString("timezone", lua.LString(facts.Timezone))
	table.RawSetString("boot_time", lua.LNumber(facts.BootTime))
	table.RawSetString("collected_at", lua.LString(facts.CollectedAt.String()))
	
	// Memory
	if facts.Memory != nil {
		memTable := L.NewTable()
		memTable.RawSetString("total", lua.LNumber(facts.Memory.Total))
		memTable.RawSetString("available", lua.LNumber(facts.Memory.Available))
		memTable.RawSetString("used", lua.LNumber(facts.Memory.Used))
		memTable.RawSetString("used_percent", lua.LNumber(facts.Memory.UsedPercent))
		memTable.RawSetString("free", lua.LNumber(facts.Memory.Free))
		memTable.RawSetString("cached", lua.LNumber(facts.Memory.Cached))
		memTable.RawSetString("buffers", lua.LNumber(facts.Memory.Buffers))
		table.RawSetString("memory", memTable)
	}
	
	// Load average
	loadArray := L.NewTable()
	for _, load := range facts.LoadAverage {
		loadArray.Append(lua.LNumber(load))
	}
	table.RawSetString("load_average", loadArray)
	
	return table
}

func (m *FactsModule) diskToLuaTable(L *lua.LState, disk *agent.DiskInfo) *lua.LTable {
	table := L.NewTable()
	table.RawSetString("device", lua.LString(disk.Device))
	table.RawSetString("mountpoint", lua.LString(disk.Mountpoint))
	table.RawSetString("fstype", lua.LString(disk.Fstype))
	table.RawSetString("total", lua.LNumber(disk.Total))
	table.RawSetString("used", lua.LNumber(disk.Used))
	table.RawSetString("free", lua.LNumber(disk.Free))
	table.RawSetString("used_percent", lua.LNumber(disk.UsedPercent))
	return table
}

func (m *FactsModule) networkToLuaTable(L *lua.LState, net *agent.NetworkInfo) *lua.LTable {
	table := L.NewTable()
	table.RawSetString("name", lua.LString(net.Name))
	table.RawSetString("mac", lua.LString(net.MAC))
	table.RawSetString("mtu", lua.LNumber(net.MTU))
	table.RawSetString("is_up", lua.LBool(net.IsUp))
	table.RawSetString("speed", lua.LNumber(net.Speed))
	
	addrsArray := L.NewTable()
	for _, addr := range net.Addresses {
		addrsArray.Append(lua.LString(addr))
	}
	table.RawSetString("addresses", addrsArray)
	
	return table
}

func getStringField(L *lua.LState, table *lua.LTable, key, defaultValue string) string {
	lv := table.RawGetString(key)
	if str, ok := lv.(lua.LString); ok {
		return string(str)
	}
	return defaultValue
}
