package luainterface

import (
	lua "github.com/yuin/gopher-lua"
)

// Network management functions
func (mod *ObjectOrientedSaltModule) saltNetworkInterface(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	interfaceName := L.CheckString(2)
	
	args := []string{"salt", target, "network.interface", interfaceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltNetworkInterfaces(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "network.interfaces", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltNetworkPing(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	host := L.CheckString(2)
	
	args := []string{"salt", target, "network.ping", host, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltNetworkTraceroute(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	host := L.CheckString(2)
	
	args := []string{"salt", target, "network.traceroute", host, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltNetworkNetstat(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "network.netstat", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltNetworkArp(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "network.arp", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// System information functions
func (mod *ObjectOrientedSaltModule) saltSystemInfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.all_status", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltSystemUptime(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.uptime", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltSystemReboot(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	delay := L.OptString(2, "1")
	
	args := []string{"salt", target, "system.reboot", delay, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltSystemShutdown(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	delay := L.OptString(2, "1")
	
	args := []string{"salt", target, "system.shutdown", delay, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltSystemHalt(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "system.halt", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltSystemHostname(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "network.get_hostname", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltSystemSetHostname(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	hostname := L.CheckString(2)
	
	args := []string{"salt", target, "network.mod_hostname", hostname, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Disk and mount management functions
func (mod *ObjectOrientedSaltModule) saltDiskUsage(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.OptString(2, "/")
	
	args := []string{"salt", target, "disk.usage", path, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDiskStats(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "disk.percent", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMountActive(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "mount.active", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMountFstab(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "mount.fstab", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMountMount(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	name := L.CheckString(2)
	device := L.CheckString(3)
	
	args := []string{"salt", target, "mount.mount", name, device, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMountUmount(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	name := L.CheckString(2)
	
	args := []string{"salt", target, "mount.umount", name, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMountRemount(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	name := L.CheckString(2)
	
	args := []string{"salt", target, "mount.remount", name, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Process management functions
func (mod *ObjectOrientedSaltModule) saltProcessList(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "ps.pgrep", ".*", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltProcessInfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pid := L.CheckString(2)
	
	args := []string{"salt", target, "ps.proc_info", pid, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltProcessKill(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pid := L.CheckString(2)
	signal := L.OptString(3, "TERM")
	
	args := []string{"salt", target, "ps.kill_pid", pid, signal, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltProcessKillall(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	name := L.CheckString(2)
	signal := L.OptString(3, "TERM")
	
	args := []string{"salt", target, "ps.killall", name, signal, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltProcessPkill(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pattern := L.CheckString(2)
	signal := L.OptString(3, "TERM")
	
	args := []string{"salt", target, "ps.pkill", pattern, signal, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Cron management functions
func (mod *ObjectOrientedSaltModule) saltCronList(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.OptString(2, "root")
	
	args := []string{"salt", target, "cron.list_tab", user, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCronSet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.CheckString(2)
	minute := L.CheckString(3)
	hour := L.CheckString(4)
	daymonth := L.CheckString(5)
	month := L.CheckString(6)
	dayweek := L.CheckString(7)
	cmd := L.CheckString(8)
	
	args := []string{"salt", target, "cron.set_job", user, minute, hour, daymonth, month, dayweek, cmd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCronDelete(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.CheckString(2)
	cmd := L.CheckString(3)
	
	args := []string{"salt", target, "cron.rm_job", user, cmd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCronRawCron(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.OptString(2, "root")
	
	args := []string{"salt", target, "cron.raw_cron", user, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Archive operations
func (mod *ObjectOrientedSaltModule) saltArchiveGunzip(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	gzipfile := L.CheckString(2)
	
	args := []string{"salt", target, "archive.gunzip", gzipfile, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltArchiveGzip(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	sourcefile := L.CheckString(2)
	
	args := []string{"salt", target, "archive.gzip", sourcefile, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltArchiveTar(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	tarfile := L.CheckString(2)
	sources := L.CheckString(3)
	
	args := []string{"salt", target, "archive.tar", "zcf", tarfile, sources, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltArchiveUntar(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	tarfile := L.CheckString(2)
	dest := L.CheckString(3)
	
	args := []string{"salt", target, "archive.tar", "zxf", tarfile, "-C", dest, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltArchiveUnzip(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	zipfile := L.CheckString(2)
	dest := L.CheckString(3)
	
	args := []string{"salt", target, "archive.unzip", zipfile, dest, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltArchiveZip(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	zipfile := L.CheckString(2)
	sources := L.CheckString(3)
	
	args := []string{"salt", target, "archive.zip", zipfile, sources, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Cloud operations
func (mod *ObjectOrientedSaltModule) saltCloudListNodes(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	provider := L.OptString(2, "")
	
	var args []string
	if provider != "" {
		args = []string{"salt-cloud", "-Q", "--list-providers", provider, "--out=json"}
	} else {
		args = []string{"salt-cloud", "-Q", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCloudCreate(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	profile := L.CheckString(2)
	name := L.CheckString(3)
	
	args := []string{"salt-cloud", "-p", profile, name, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*10, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCloudDestroy(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	name := L.CheckString(2)
	
	args := []string{"salt-cloud", "-d", name, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*5, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCloudAction(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	action := L.CheckString(2)
	name := L.CheckString(3)
	
	args := []string{"salt-cloud", "-a", action, name, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCloudFunction(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	function := L.CheckString(2)
	provider := L.CheckString(3)
	
	args := []string{"salt-cloud", "-f", function, provider, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCloudMap(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	mapfile := L.CheckString(2)
	
	args := []string{"salt-cloud", "-m", mapfile, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*10, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCloudProfile(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt-cloud", "--list-profiles", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltCloudProvider(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt-cloud", "--list-providers", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Event system
func (mod *ObjectOrientedSaltModule) saltEventSend(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	tag := L.CheckString(2)
	data := L.CheckString(3)
	
	args := []string{"salt-run", "event.send", tag, data, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltEventListen(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	tag := L.OptString(2, "*")
	
	args := []string{"salt-run", "state.event", "pretty=True", "tag=" + tag, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltEventFire(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	data := L.CheckString(2)
	tag := L.CheckString(3)
	
	args := []string{"salt", target, "event.fire", data, tag, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltEventFireMaster(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	data := L.CheckString(2)
	tag := L.CheckString(3)
	
	args := []string{"salt", target, "event.fire_master", data, tag, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Orchestration
func (mod *ObjectOrientedSaltModule) saltOrchestrate(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	sls := L.CheckString(2)
	opts := L.OptTable(3, L.NewTable())
	
	args := []string{"salt-run", "state.orchestrate", sls, "--out=json"}
	
	// Add pillar data if specified
	if pillar := opts.RawGetString("pillar"); pillar != lua.LNil {
		args = append(args, "pillar='"+pillar.String()+"'")
	}
	
	result, err := mod.executeSaltCommand(L, timeout*10, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltRunner(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	args := []string{"salt-run", module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltWheel(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	args := []string{"salt-wheel", module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}