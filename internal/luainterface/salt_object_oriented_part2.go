package luainterface

import (
	lua "github.com/yuin/gopher-lua"
)

// Targeting and configuration methods
func (mod *ObjectOrientedSaltModule) saltTarget(L *lua.LState) int {
	client := L.CheckTable(1)
	target := L.CheckString(2)
	targetType := L.OptString(3, "glob")
	
	// Create a new target object that inherits from the client
	targetObj := L.NewTable()
	
	// Copy all methods from client to target
	client.ForEach(func(key, value lua.LValue) {
		if value.Type() == lua.LTFunction {
			targetObj.RawSet(key, value)
		}
	})
	
	// Set target-specific properties
	targetObj.RawSetString("_target", lua.LString(target))
	targetObj.RawSetString("_target_type", lua.LString(targetType))
	targetObj.RawSetString("_is_target", lua.LBool(true))
	
	L.Push(targetObj)
	return 1
}

func (mod *ObjectOrientedSaltModule) saltWithTimeout(L *lua.LState) int {
	client := L.CheckTable(1)
	timeout := L.CheckInt(2)
	
	// Clone the client with new timeout
	newClient := L.NewTable()
	client.ForEach(func(key, value lua.LValue) {
		newClient.RawSet(key, value)
	})
	
	newClient.RawSetString("_timeout", lua.LNumber(timeout))
	
	L.Push(newClient)
	return 1
}

func (mod *ObjectOrientedSaltModule) saltWithRetries(L *lua.LState) int {
	client := L.CheckTable(1)
	retries := L.CheckInt(2)
	
	// Clone the client with new retry count
	newClient := L.NewTable()
	client.ForEach(func(key, value lua.LValue) {
		newClient.RawSet(key, value)
	})
	
	newClient.RawSetString("_retries", lua.LNumber(retries))
	
	L.Push(newClient)
	return 1
}

func (mod *ObjectOrientedSaltModule) saltWithPillar(L *lua.LState) int {
	client := L.CheckTable(1)
	pillar := L.CheckTable(2)
	
	// Clone the client with pillar data
	newClient := L.NewTable()
	client.ForEach(func(key, value lua.LValue) {
		newClient.RawSet(key, value)
	})
	
	newClient.RawSetString("_pillar", pillar)
	
	L.Push(newClient)
	return 1
}

func (mod *ObjectOrientedSaltModule) saltWithGrains(L *lua.LState) int {
	client := L.CheckTable(1)
	grains := L.CheckTable(2)
	
	// Clone the client with grains data
	newClient := L.NewTable()
	client.ForEach(func(key, value lua.LValue) {
		newClient.RawSet(key, value)
	})
	
	newClient.RawSetString("_grains", grains)
	
	L.Push(newClient)
	return 1
}

// Helper to get target from client
func (mod *ObjectOrientedSaltModule) getTarget(L *lua.LState) string {
	client := L.CheckTable(1)
	target := client.RawGetString("_target")
	if target == lua.LNil {
		return "*"
	}
	return target.String()
}

// Helper to get timeout from client
func (mod *ObjectOrientedSaltModule) getTimeout(L *lua.LState) int {
	client := L.CheckTable(1)
	timeout := client.RawGetString("_timeout")
	if timeout == lua.LNil {
		return 60
	}
	return int(lua.LVAsNumber(timeout))
}

// Helper to get retries from client
func (mod *ObjectOrientedSaltModule) getRetries(L *lua.LState) int {
	client := L.CheckTable(1)
	retries := client.RawGetString("_retries")
	if retries == lua.LNil {
		return 3
	}
	return int(lua.LVAsNumber(retries))
}

// Basic execution functions
func (mod *ObjectOrientedSaltModule) saltCmd(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	args := []string{"salt", target, module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltRun(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	args := []string{"salt-run", module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltExecute(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	args := []string{"salt", target, module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltBatch(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	batchSize := L.CheckString(2)
	module := L.CheckString(3)
	function := L.CheckString(4)
	
	args := []string{"salt", "--batch-size=" + batchSize, target, module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 5; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltAsync(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	args := []string{"salt", target, module + "." + function, "--async", "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Connection and testing functions
func (mod *ObjectOrientedSaltModule) saltPing(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "test.ping", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltTest(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	testType := L.CheckString(2)
	
	args := []string{"salt", target, "test." + testType, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltVersion(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "test.version", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStatus(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.uptime", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Key management functions
func (mod *ObjectOrientedSaltModule) saltKeyList(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	keyType := L.OptString(2, "all")
	
	args := []string{"salt-key", "-L", keyType, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltKeyAccept(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	keyPattern := L.CheckString(2)
	args := []string{"salt-key", "-a", keyPattern, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltKeyReject(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	keyPattern := L.CheckString(2)
	args := []string{"salt-key", "-r", keyPattern, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltKeyDelete(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	keyPattern := L.CheckString(2)
	args := []string{"salt-key", "-d", keyPattern, "-y", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltKeyFinger(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	keyPattern := L.OptString(2, "*")
	args := []string{"salt-key", "-f", keyPattern, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltKeyGen(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	keyName := L.CheckString(2)
	keySize := L.OptString(3, "2048")
	
	args := []string{"salt-key", "--gen-keys=" + keyName, "--keysize=" + keySize}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// State management functions
func (mod *ObjectOrientedSaltModule) saltStateApply(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
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
	
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStateHighstate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	opts := L.OptTable(2, L.NewTable())
	
	args := []string{"salt", target, "state.highstate", "--out=json"}
	
	// Add test mode if specified
	if test := opts.RawGetString("test"); lua.LVAsBool(test) {
		args = append(args, "test=True")
	}
	
	result, err := mod.executeSaltCommand(L, timeout*5, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStateTest(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	sls := L.OptString(2, "")
	
	var args []string
	if sls != "" {
		args = []string{"salt", target, "state.apply", sls, "test=True", "--out=json"}
	} else {
		args = []string{"salt", target, "state.highstate", "test=True", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStateShowSls(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	sls := L.CheckString(2)
	
	args := []string{"salt", target, "state.show_sls", sls, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStateShowTop(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "state.show_top", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStateShowLowstate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "state.show_lowstate", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStateSingle(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	fun := L.CheckString(2)
	name := L.CheckString(3)
	
	args := []string{"salt", target, "state.single", fun, name, "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStateTemplate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	template := L.CheckString(2)
	
	args := []string{"salt", target, "state.template", template, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*3, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Grains management functions
func (mod *ObjectOrientedSaltModule) saltGrainsGet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	key := L.OptString(2, "")
	
	var args []string
	if key != "" {
		args = []string{"salt", target, "grains.get", key, "--out=json"}
	} else {
		args = []string{"salt", target, "grains.items", "--out=json"}
	}
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGrainsSet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	args := []string{"salt", target, "grains.setval", key, value, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGrainsAppend(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	args := []string{"salt", target, "grains.append", key, value, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGrainsRemove(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	args := []string{"salt", target, "grains.remove", key, value, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGrainsDelkey(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	key := L.CheckString(2)
	
	args := []string{"salt", target, "grains.delkey", key, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGrainsItems(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "grains.items", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Pillar management functions
func (mod *ObjectOrientedSaltModule) saltPillarGet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	key := L.CheckString(2)
	
	args := []string{"salt", target, "pillar.get", key, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPillarItems(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "pillar.items", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPillarShow(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "pillar.show_pillar", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPillarRefresh(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "saltutil.refresh_pillar", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}