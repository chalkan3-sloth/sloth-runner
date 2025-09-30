package luainterface

import (
	lua "github.com/yuin/gopher-lua"
)

// Key management functions
func (mod *ComprehensiveSaltModule) saltKeyList(L *lua.LState) int {
	keyType := "all"
	if L.GetTop() > 0 {
		keyType = L.CheckString(1)
	}
	
	args := []string{"salt-key", "-L", keyType, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltKeyAccept(L *lua.LState) int {
	keyPattern := L.CheckString(1)
	args := []string{"salt-key", "-a", keyPattern, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltKeyReject(L *lua.LState) int {
	keyPattern := L.CheckString(1)
	args := []string{"salt-key", "-r", keyPattern, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltKeyDelete(L *lua.LState) int {
	keyPattern := L.CheckString(1)
	args := []string{"salt-key", "-d", keyPattern, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltKeyFinger(L *lua.LState) int {
	keyPattern := "*"
	if L.GetTop() > 0 {
		keyPattern = L.CheckString(1)
	}
	args := []string{"salt-key", "-f", keyPattern, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltKeyGen(L *lua.LState) int {
	keyName := L.CheckString(1)
	keySize := "2048"
	if L.GetTop() > 1 {
		keySize = L.CheckString(2)
	}
	
	args := []string{"salt-key", "--gen-keys=" + keyName, "--keysize=" + keySize}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// State management functions
func (mod *ComprehensiveSaltModule) saltStateApply(L *lua.LState) int {
	target := L.CheckString(1)
	sls := L.CheckString(2)
	opts := L.OptTable(3, L.NewTable())
	
	args := []string{"salt", target, "state.apply", sls, "--out=json"}
	
	// Add test mode if specified
	if test := opts.RawGetString("test"); lua.LVAsBool(test) {
		args = append(args, "test=True")
	}
	
	// Add pillar data if specified
	if pillar := opts.RawGetString("pillar"); pillar != lua.LNil {
		args = append(args, "pillar='"+pillar.String()+"'")
	}
	
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStateHighstate(L *lua.LState) int {
	target := L.CheckString(1)
	opts := L.OptTable(2, L.NewTable())
	
	args := []string{"salt", target, "state.highstate", "--out=json"}
	
	// Add test mode if specified
	if test := opts.RawGetString("test"); lua.LVAsBool(test) {
		args = append(args, "test=True")
	}
	
	result, err := mod.executeSaltCommand(600, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStateTest(L *lua.LState) int {
	target := L.CheckString(1)
	sls := ""
	if L.GetTop() > 1 {
		sls = L.CheckString(2)
	}
	
	var args []string
	if sls != "" {
		args = []string{"salt", target, "state.apply", sls, "test=True", "--out=json"}
	} else {
		args = []string{"salt", target, "state.highstate", "test=True", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStateShowSls(L *lua.LState) int {
	target := L.CheckString(1)
	sls := L.CheckString(2)
	
	args := []string{"salt", target, "state.show_sls", sls, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStateShowTop(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "state.show_top", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStateShowLowstate(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "state.show_lowstate", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStateSingle(L *lua.LState) int {
	target := L.CheckString(1)
	fun := L.CheckString(2)
	name := L.CheckString(3)
	
	args := []string{"salt", target, "state.single", fun, name, "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStateTemplate(L *lua.LState) int {
	target := L.CheckString(1)
	template := L.CheckString(2)
	
	args := []string{"salt", target, "state.template", template, "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Grains management functions
func (mod *ComprehensiveSaltModule) saltGrainsGet(L *lua.LState) int {
	target := L.CheckString(1)
	key := ""
	if L.GetTop() > 1 {
		key = L.CheckString(2)
	}
	
	var args []string
	if key != "" {
		args = []string{"salt", target, "grains.get", key, "--out=json"}
	} else {
		args = []string{"salt", target, "grains.items", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGrainsSet(L *lua.LState) int {
	target := L.CheckString(1)
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	args := []string{"salt", target, "grains.setval", key, value, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGrainsAppend(L *lua.LState) int {
	target := L.CheckString(1)
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	args := []string{"salt", target, "grains.append", key, value, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGrainsRemove(L *lua.LState) int {
	target := L.CheckString(1)
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	args := []string{"salt", target, "grains.remove", key, value, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGrainsDelkey(L *lua.LState) int {
	target := L.CheckString(1)
	key := L.CheckString(2)
	
	args := []string{"salt", target, "grains.delkey", key, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltGrainsItems(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "grains.items", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Pillar management functions
func (mod *ComprehensiveSaltModule) saltPillarGet(L *lua.LState) int {
	target := L.CheckString(1)
	key := L.CheckString(2)
	
	args := []string{"salt", target, "pillar.get", key, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPillarItems(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "pillar.items", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPillarShow(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "pillar.show_pillar", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPillarRefresh(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "saltutil.refresh_pillar", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// File operations
func (mod *ComprehensiveSaltModule) saltFileCopy(L *lua.LState) int {
	target := L.CheckString(1)
	src := L.CheckString(2)
	dst := L.CheckString(3)
	
	args := []string{"salt", target, "file.copy", src, dst, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileGet(L *lua.LState) int {
	target := L.CheckString(1)
	path := L.CheckString(2)
	
	args := []string{"salt", target, "cp.get_file", path, path, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileList(L *lua.LState) int {
	target := L.CheckString(1)
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.readdir", path, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileManage(L *lua.LState) int {
	target := L.CheckString(1)
	name := L.CheckString(2)
	
	args := []string{"salt", target, "file.managed", name, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileRecurse(L *lua.LState) int {
	target := L.CheckString(1)
	name := L.CheckString(2)
	source := L.CheckString(3)
	
	args := []string{"salt", target, "file.recurse", name, source, "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileTouch(L *lua.LState) int {
	target := L.CheckString(1)
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.touch", path, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileStats(L *lua.LState) int {
	target := L.CheckString(1)
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.stats", path, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileFind(L *lua.LState) int {
	target := L.CheckString(1)
	path := L.CheckString(2)
	
	args := []string{"salt", target, "file.find", path, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileReplace(L *lua.LState) int {
	target := L.CheckString(1)
	path := L.CheckString(2)
	pattern := L.CheckString(3)
	repl := L.CheckString(4)
	
	args := []string{"salt", target, "file.replace", path, pattern, repl, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltFileCheckHash(L *lua.LState) int {
	target := L.CheckString(1)
	path := L.CheckString(2)
	hashType := "md5"
	if L.GetTop() > 2 {
		hashType = L.CheckString(3)
	}
	
	args := []string{"salt", target, "file.get_hash", path, hashType, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Package management functions
func (mod *ComprehensiveSaltModule) saltPkgInstall(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.install", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(600, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgRemove(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.remove", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgUpgrade(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := ""
	if L.GetTop() > 1 {
		pkgName = L.CheckString(2)
	}
	
	var args []string
	if pkgName != "" {
		args = []string{"salt", target, "pkg.upgrade", pkgName, "--out=json"}
	} else {
		args = []string{"salt", target, "pkg.upgrade", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(1200, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgRefresh(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "pkg.refresh_db", "--out=json"}
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgList(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "pkg.list_pkgs", "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgVersion(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.version", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgAvailable(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := ""
	if L.GetTop() > 1 {
		pkgName = L.CheckString(2)
	}
	
	var args []string
	if pkgName != "" {
		args = []string{"salt", target, "pkg.available_version", pkgName, "--out=json"}
	} else {
		args = []string{"salt", target, "pkg.list_repo_pkgs", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgInfo(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.info", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgHold(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.hold", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltPkgUnhold(L *lua.LState) int {
	target := L.CheckString(1)
	pkgName := L.CheckString(2)
	
	args := []string{"salt", target, "pkg.unhold", pkgName, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Service management functions
func (mod *ComprehensiveSaltModule) saltServiceStart(L *lua.LState) int {
	target := L.CheckString(1)
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.start", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltServiceStop(L *lua.LState) int {
	target := L.CheckString(1)
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.stop", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltServiceRestart(L *lua.LState) int {
	target := L.CheckString(1)
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.restart", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltServiceReload(L *lua.LState) int {
	target := L.CheckString(1)
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.reload", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(120, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltServiceStatus(L *lua.LState) int {
	target := L.CheckString(1)
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.status", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltServiceEnable(L *lua.LState) int {
	target := L.CheckString(1)
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.enable", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltServiceDisable(L *lua.LState) int {
	target := L.CheckString(1)
	serviceName := L.CheckString(2)
	
	args := []string{"salt", target, "service.disable", serviceName, "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltServiceList(L *lua.LState) int {
	target := L.CheckString(1)
	
	args := []string{"salt", target, "service.get_all", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}