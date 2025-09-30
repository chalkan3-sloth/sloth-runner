package luainterface

import (
	lua "github.com/yuin/gopher-lua"
)

// User management functions
func (mod *ComprehensiveSaltModule) saltUserAdd(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	opts := L.OptTable(3, L.NewTable())
	
	args := []string{"salt", target, "user.add", username}
	
	// Add optional parameters
	if home := opts.RawGetString("home"); home != lua.LNil {
		args = append(args, "home="+home.String())
	}
	if shell := opts.RawGetString("shell"); shell != lua.LNil {
		args = append(args, "shell="+shell.String())
	}
	if uid := opts.RawGetString("uid"); uid != lua.LNil {
		args = append(args, "uid="+uid.String())
	}
	if gid := opts.RawGetString("gid"); gid != lua.LNil {
		args = append(args, "gid="+gid.String())
	}
	
	args = append(args, "--out=json")
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserDelete(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	
	args := []string{"salt", target, "user.delete", username, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserInfo(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	
	args := []string{"salt", target, "user.info", username, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserList(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "user.list_users", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserChuid(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	uid := L.CheckString(3)
	
	args := []string{"salt", target, "user.chuid", username, uid, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserChgid(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	gid := L.CheckString(3)
	
	args := []string{"salt", target, "user.chgid", username, gid, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserChshell(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	shell := L.CheckString(3)
	
	args := []string{"salt", target, "user.chshell", username, shell, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserChhome(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	home := L.CheckString(3)
	
	args := []string{"salt", target, "user.chhome", username, home, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltUserPrimaryGroup(L *lua.LState) int {
	target := L.CheckString(1)
	username := L.CheckString(2)
	group := L.CheckString(3)
	
	args := []string{"salt", target, "user.chgid", username, group, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Group management functions
func (mod *ComprehensiveSaltModule) saltGroupAdd(L *lua.LState) int {
	target := L.CheckString(1)
	groupname := L.CheckString(2)
	gid := ""
	if L.GetTop() > 2 {
		gid = L.CheckString(3)
	}
	
	args := []string{"salt", target, "group.add", groupname}
	if gid != "" {
		args = append(args, gid)
	}
	args = append(args, "--out=json")
	
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGroupDelete(L *lua.LState) int {
	target := L.CheckString(1)
	groupname := L.CheckString(2)
	
	args := []string{"salt", target, "group.delete", groupname, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGroupInfo(L *lua.LState) int {
	target := L.CheckString(1)
	groupname := L.CheckString(2)
	
	args := []string{"salt", target, "group.info", groupname, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGroupList(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "group.getent", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGroupAdduser(L *lua.LState) int {
	target := L.CheckString(1)
	groupname := L.CheckString(2)
	username := L.CheckString(3)
	
	args := []string{"salt", target, "group.adduser", groupname, username, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGroupDeluser(L *lua.LState) int {
	target := L.CheckString(1)
	groupname := L.CheckString(2)
	username := L.CheckString(3)
	
	args := []string{"salt", target, "group.deluser", groupname, username, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGroupMembers(L *lua.LState) int {
	target := L.CheckString(1)
	groupname := L.CheckString(2)
	
	args := []string{"salt", target, "group.members", groupname, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Network management functions
func (mod *ComprehensiveSaltModule) saltNetworkInterface(L *lua.LState) int {
	target := L.CheckString(1)
	interface_name := L.CheckString(2)
	
	args := []string{"salt", target, "network.interface", interface_name, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltNetworkInterfaces(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "network.interfaces", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltNetworkPing(L *lua.LState) int {
	target := L.CheckString(1)
	host := L.CheckString(2)
	
	args := []string{"salt", target, "network.ping", host, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltNetworkTraceroute(L *lua.LState) int {
	target := L.CheckString(1)
	host := L.CheckString(2)
	
	args := []string{"salt", target, "network.traceroute", host, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltNetworkNetstat(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "network.netstat", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltNetworkArp(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "network.arp", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// System information functions
func (mod *ComprehensiveSaltModule) saltSystemInfo(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "status.all_status", "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltSystemUptime(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "status.uptime", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltSystemReboot(L *lua.LState) int {
	target := L.CheckString(1)
	delay := "1"
	if L.GetTop() > 1 {
		delay = L.CheckString(2)
	}
	
	args := []string{"salt", target, "system.reboot", delay, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltSystemShutdown(L *lua.LState) int {
	target := L.CheckString(1)
	delay := "1"
	if L.GetTop() > 1 {
		delay = L.CheckString(2)
	}
	
	args := []string{"salt", target, "system.shutdown", delay, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltSystemHalt(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "system.halt", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltSystemHostname(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "network.get_hostname", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltSystemSetHostname(L *lua.LState) int {
	target := L.CheckString(1)
	hostname := L.CheckString(2)
	
	args := []string{"salt", target, "network.mod_hostname", hostname, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Disk and mount management functions
func (mod *ComprehensiveSaltModule) saltDiskUsage(L *lua.LState) int {
	target := L.CheckString(1)
	path := "/"
	if L.GetTop() > 1 {
		path = L.CheckString(2)
	}
	
	args := []string{"salt", target, "disk.usage", path, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltDiskStats(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "disk.percent", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltMountActive(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "mount.active", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltMountFstab(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "mount.fstab", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltMountMount(L *lua.LState) int {
	target := L.CheckString(1)
	name := L.CheckString(2)
	device := L.CheckString(3)
	
	args := []string{"salt", target, "mount.mount", name, device, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltMountUmount(L *lua.LState) int {
	target := L.CheckString(1)
	name := L.CheckString(2)
	
	args := []string{"salt", target, "mount.umount", name, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltMountRemount(L *lua.LState) int {
	target := L.CheckString(1)
	name := L.CheckString(2)
	
	args := []string{"salt", target, "mount.remount", name, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Process management functions
func (mod *ComprehensiveSaltModule) saltProcessList(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "ps.pgrep", ".*", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltProcessInfo(L *lua.LState) int {
	target := L.CheckString(1)
	pid := L.CheckString(2)
	
	args := []string{"salt", target, "ps.proc_info", pid, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltProcessKill(L *lua.LState) int {
	target := L.CheckString(1)
	pid := L.CheckString(2)
	signal := "TERM"
	if L.GetTop() > 2 {
		signal = L.CheckString(3)
	}
	
	args := []string{"salt", target, "ps.kill_pid", pid, signal, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltProcessKillall(L *lua.LState) int {
	target := L.CheckString(1)
	name := L.CheckString(2)
	signal := "TERM"
	if L.GetTop() > 2 {
		signal = L.CheckString(3)
	}
	
	args := []string{"salt", target, "ps.killall", name, signal, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltProcessPkill(L *lua.LState) int {
	target := L.CheckString(1)
	pattern := L.CheckString(2)
	signal := "TERM"
	if L.GetTop() > 2 {
		signal = L.CheckString(3)
	}
	
	args := []string{"salt", target, "ps.pkill", pattern, signal, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Cron management functions
func (mod *ComprehensiveSaltModule) saltCronList(L *lua.LState) int {
	target := L.CheckString(1)
	user := "root"
	if L.GetTop() > 1 {
		user = L.CheckString(2)
	}
	
	args := []string{"salt", target, "cron.list_tab", user, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCronSet(L *lua.LState) int {
	target := L.CheckString(1)
	user := L.CheckString(2)
	minute := L.CheckString(3)
	hour := L.CheckString(4)
	daymonth := L.CheckString(5)
	month := L.CheckString(6)
	dayweek := L.CheckString(7)
	cmd := L.CheckString(8)
	
	args := []string{"salt", target, "cron.set_job", user, minute, hour, daymonth, month, dayweek, cmd, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCronDelete(L *lua.LState) int {
	target := L.CheckString(1)
	user := L.CheckString(2)
	cmd := L.CheckString(3)
	
	args := []string{"salt", target, "cron.rm_job", user, cmd, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCronRawCron(L *lua.LState) int {
	target := L.CheckString(1)
	user := "root"
	if L.GetTop() > 1 {
		user = L.CheckString(2)
	}
	
	args := []string{"salt", target, "cron.raw_cron", user, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Archive operations
func (mod *ComprehensiveSaltModule) saltArchiveGunzip(L *lua.LState) int {
	target := L.CheckString(1)
	gzipfile := L.CheckString(2)
	
	args := []string{"salt", target, "archive.gunzip", gzipfile, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltArchiveGzip(L *lua.LState) int {
	target := L.CheckString(1)
	sourcefile := L.CheckString(2)
	
	args := []string{"salt", target, "archive.gzip", sourcefile, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltArchiveTar(L *lua.LState) int {
	target := L.CheckString(1)
	tarfile := L.CheckString(2)
	sources := L.CheckString(3)
	
	args := []string{"salt", target, "archive.tar", "zcf", tarfile, sources, "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltArchiveUntar(L *lua.LState) int {
	target := L.CheckString(1)
	tarfile := L.CheckString(2)
	dest := L.CheckString(3)
	
	args := []string{"salt", target, "archive.tar", "zxf", tarfile, "-C", dest, "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltArchiveUnzip(L *lua.LState) int {
	target := L.CheckString(1)
	zipfile := L.CheckString(2)
	dest := L.CheckString(3)
	
	args := []string{"salt", target, "archive.unzip", zipfile, dest, "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltArchiveZip(L *lua.LState) int {
	target := L.CheckString(1)
	zipfile := L.CheckString(2)
	sources := L.CheckString(3)
	
	args := []string{"salt", target, "archive.zip", zipfile, sources, "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// These are placeholder implementations for the advanced functions
// In a real implementation, these would have full functionality

func (mod *ComprehensiveSaltModule) saltCloudListNodes(L *lua.LState) int {
	provider := ""
	if L.GetTop() > 0 {
		provider = L.CheckString(1)
	}
	
	var args []string
	if provider != "" {
		args = []string{"salt-cloud", "-Q", "--list-providers", provider, "--out=json"}
	} else {
		args = []string{"salt-cloud", "-Q", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCloudCreate(L *lua.LState) int {
	profile := L.CheckString(1)
	name := L.CheckString(2)
	
	args := []string{"salt-cloud", "-p", profile, name, "--out=json"}
	result, err := mod.executeSaltCommand(600, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCloudDestroy(L *lua.LState) int {
	name := L.CheckString(1)
	
	args := []string{"salt-cloud", "-d", name, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCloudAction(L *lua.LState) int {
	action := L.CheckString(1)
	name := L.CheckString(2)
	
	args := []string{"salt-cloud", "-a", action, name, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCloudFunction(L *lua.LState) int {
	function := L.CheckString(1)
	provider := L.CheckString(2)
	
	args := []string{"salt-cloud", "-f", function, provider, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCloudMap(L *lua.LState) int {
	mapfile := L.CheckString(1)
	
	args := []string{"salt-cloud", "-m", mapfile, "--out=json"}
	result, err := mod.executeSaltCommand(600, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCloudProfile(L *lua.LState) int {
	args := []string{"salt-cloud", "--list-profiles", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltCloudProvider(L *lua.LState) int {
	args := []string{"salt-cloud", "--list-providers", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Simplified implementations for remaining functions
func (mod *ComprehensiveSaltModule) saltEventSend(L *lua.LState) int {
	tag := L.CheckString(1)
	data := L.CheckString(2)
	
	args := []string{"salt-run", "event.send", tag, data, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltEventListen(L *lua.LState) int {
	tag := "*"
	if L.GetTop() > 0 {
		tag = L.CheckString(1)
	}
	
	args := []string{"salt-run", "state.event", "pretty=True", "tag=" + tag, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltEventFire(L *lua.LState) int {
	data := L.CheckString(1)
	tag := L.CheckString(2)
	target := L.CheckString(3)
	
	args := []string{"salt", target, "event.fire", data, tag, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltEventFireMaster(L *lua.LState) int {
	data := L.CheckString(1)
	tag := L.CheckString(2)
	target := L.CheckString(3)
	
	args := []string{"salt", target, "event.fire_master", data, tag, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltOrchestrate(L *lua.LState) int {
	sls := L.CheckString(1)
	opts := L.OptTable(2, L.NewTable())
	
	args := []string{"salt-run", "state.orchestrate", sls, "--out=json"}
	
	// Add pillar data if specified
	if pillar := opts.RawGetString("pillar"); pillar != lua.LNil {
		args = append(args, "pillar='"+pillar.String()+"'")
	}
	
	result, err := mod.executeSaltCommand(600, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltRunner(L *lua.LState) int {
	module := L.CheckString(1)
	function := L.CheckString(2)
	
	args := []string{"salt-run", module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltWheel(L *lua.LState) int {
	module := L.CheckString(1)
	function := L.CheckString(2)
	
	args := []string{"salt-wheel", module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}