package luainterface

import (
	lua "github.com/yuin/gopher-lua"
)

// File operations
func (mod *ObjectOrientedSaltModule) saltFileCopy(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	src := L.CheckString(2)
	dst := L.CheckString(3)
	
	args := []string{"salt", target, "file.copy", src, dst, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileGet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	
	args := []string{"salt", target, "cp.get_file", path, path, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileList(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.readdir", path, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileManage(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	name := L.CheckString(2)
	
	args := []string{"salt", target, "file.managed", name, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileRecurse(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	name := L.CheckString(2)
	source := L.CheckString(3)
	
	args := []string{"salt", target, "file.recurse", name, source, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileTouch(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.touch", path, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileStats(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.stats", path, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileFind(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.find", path, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileReplace(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	pattern := L.CheckString(3)
	repl := L.CheckString(4)
	
	args := []string{"salt", target, "file.replace", path, pattern, repl, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltFileCheckHash(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	hashType := L.OptString(3, "md5")
	
	args := []string{"salt", target, "file.get_hash", path, hashType, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Package management functions
func (mod *ObjectOrientedSaltModule) saltPkgInstall(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.install", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*5, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgRemove(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.remove", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgUpgrade(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.OptString(2, "")
	
	var args []string
	if pkgName != "" {
		args = []string{"salt", target, "pkg.upgrade", pkgName, "--out=json"}
	} else {
		args = []string{"salt", target, "pkg.upgrade", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(L, timeout*10, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgRefresh(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "pkg.refresh_db", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgList(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "pkg.list_pkgs", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgVersion(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.version", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgAvailable(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.OptString(2, "")
	
	var args []string
	if pkgName != "" {
		args = []string{"salt", target, "pkg.available_version", pkgName, "--out=json"}
	} else {
		args = []string{"salt", target, "pkg.list_repo_pkgs", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgInfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.info", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgHold(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.hold", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPkgUnhold(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.unhold", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Service management functions
func (mod *ObjectOrientedSaltModule) saltServiceStart(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.start", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltServiceStop(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.stop", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltServiceRestart(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.restart", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltServiceReload(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.reload", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltServiceStatus(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.status", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltServiceEnable(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.enable", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltServiceDisable(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.disable", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltServiceList(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "service.get_all", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// User management functions
func (mod *ObjectOrientedSaltModule) saltUserAdd(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
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
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserDelete(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	username := L.CheckString(2)
	
	args := []string{"salt", target, "user.delete", username, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserInfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	username := L.CheckString(2)
	
	args := []string{"salt", target, "user.info", username, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserList(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "user.list_users", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserChuid(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	username := L.CheckString(2)
	uid := L.CheckString(3)
	
	args := []string{"salt", target, "user.chuid", username, uid, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserChgid(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	username := L.CheckString(2)
	gid := L.CheckString(3)
	
	args := []string{"salt", target, "user.chgid", username, gid, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserChshell(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	username := L.CheckString(2)
	shell := L.CheckString(3)
	
	args := []string{"salt", target, "user.chshell", username, shell, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserChhome(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	username := L.CheckString(2)
	home := L.CheckString(3)
	
	args := []string{"salt", target, "user.chhome", username, home, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltUserPrimaryGroup(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	username := L.CheckString(2)
	group := L.CheckString(3)
	
	args := []string{"salt", target, "user.chgid", username, group, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Group management functions
func (mod *ObjectOrientedSaltModule) saltGroupAdd(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	groupname := L.CheckString(2)
	gid := L.OptString(3, "")
	
	args := []string{"salt", target, "group.add", groupname}
	if gid != "" {
		args = append(args, gid)
	}
	args = append(args, "--out=json")
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGroupDelete(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	groupname := L.CheckString(2)
	
	args := []string{"salt", target, "group.delete", groupname, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGroupInfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	groupname := L.CheckString(2)
	
	args := []string{"salt", target, "group.info", groupname, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGroupList(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "group.getent", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGroupAdduser(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	groupname := L.CheckString(2)
	username := L.CheckString(3)
	
	args := []string{"salt", target, "group.adduser", groupname, username, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGroupDeluser(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	groupname := L.CheckString(2)
	username := L.CheckString(3)
	
	args := []string{"salt", target, "group.deluser", groupname, username, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGroupMembers(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	groupname := L.CheckString(2)
	
	args := []string{"salt", target, "group.members", groupname, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}