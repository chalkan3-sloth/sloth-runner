package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type IncusModule struct {
	agentClient interface{} // placeholder for future agent integration
}

func NewIncusModule(agentClient interface{}) *IncusModule {
	return &IncusModule{
		agentClient: agentClient,
	}
}

// IncusInstance representa uma instância Incus
type IncusInstance struct {
	L           *lua.LState
	agentClient interface{}
	name        string
	config      map[string]interface{}
	profiles    []string
	devices     map[string]map[string]string
	target      string
}

// IncusImage representa uma imagem Incus
type IncusImage struct {
	L           *lua.LState
	agentClient interface{}
	alias       string
	source      string
	server      string
	target      string
}

// IncusNetwork representa uma rede Incus
type IncusNetwork struct {
	L           *lua.LState
	agentClient interface{}
	name        string
	networkType string
	config      map[string]interface{}
	target      string
}

// IncusProfile representa um perfil Incus
type IncusProfile struct {
	L           *lua.LState
	agentClient interface{}
	name        string
	description string
	config      map[string]interface{}
	devices     map[string]map[string]string
	target      string
}

// IncusStorage representa um storage pool Incus
type IncusStorage struct {
	L           *lua.LState
	agentClient interface{}
	name        string
	driver      string
	config      map[string]interface{}
	target      string
}

// IncusSnapshot representa um snapshot de instância
type IncusSnapshot struct {
	L           *lua.LState
	agentClient interface{}
	instance    string
	name        string
	stateful    bool
	target      string
}

// Register registra o módulo Incus no Lua
func (m *IncusModule) Register(L *lua.LState) {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"instance": m.createInstance,
		"image":    m.createImage,
		"network":  m.createNetwork,
		"profile":  m.createProfile,
		"storage":  m.createStorage,
		"snapshot": m.createSnapshot,
		"exec":     m.exec,
		"list":     m.list,
		"info":     m.info,
		"delete":   m.delete,
	})

	L.SetGlobal("incus", mod)
}

// createInstance cria um builder de instância Incus
func (m *IncusModule) createInstance(L *lua.LState) int {
	name := L.CheckString(1)

	instance := &IncusInstance{
		L:           L,
		agentClient: m.agentClient,
		name:        name,
		config:      make(map[string]interface{}),
		profiles:    []string{},
		devices:     make(map[string]map[string]string),
		target:      "",
	}

	ud := L.NewUserData()
	ud.Value = instance
	L.SetMetatable(ud, L.GetTypeMetatable("incus_instance"))
	L.Push(ud)
	return 1
}

// createImage cria um builder de imagem Incus
func (m *IncusModule) createImage(L *lua.LState) int {
	alias := L.CheckString(1)

	image := &IncusImage{
		L:           L,
		agentClient: m.agentClient,
		alias:       alias,
		server:      "https://images.linuxcontainers.org",
		target:      "",
	}

	ud := L.NewUserData()
	ud.Value = image
	L.SetMetatable(ud, L.GetTypeMetatable("incus_image"))
	L.Push(ud)
	return 1
}

// createNetwork cria um builder de rede Incus
func (m *IncusModule) createNetwork(L *lua.LState) int {
	name := L.CheckString(1)

	network := &IncusNetwork{
		L:           L,
		agentClient: m.agentClient,
		name:        name,
		networkType: "bridge",
		config:      make(map[string]interface{}),
		target:      "",
	}

	ud := L.NewUserData()
	ud.Value = network
	L.SetMetatable(ud, L.GetTypeMetatable("incus_network"))
	L.Push(ud)
	return 1
}

// createProfile cria um builder de perfil Incus
func (m *IncusModule) createProfile(L *lua.LState) int {
	name := L.CheckString(1)

	profile := &IncusProfile{
		L:           L,
		agentClient: m.agentClient,
		name:        name,
		config:      make(map[string]interface{}),
		devices:     make(map[string]map[string]string),
		target:      "",
	}

	ud := L.NewUserData()
	ud.Value = profile
	L.SetMetatable(ud, L.GetTypeMetatable("incus_profile"))
	L.Push(ud)
	return 1
}

// createStorage cria um builder de storage pool Incus
func (m *IncusModule) createStorage(L *lua.LState) int {
	name := L.CheckString(1)

	storage := &IncusStorage{
		L:           L,
		agentClient: m.agentClient,
		name:        name,
		driver:      "dir",
		config:      make(map[string]interface{}),
		target:      "",
	}

	ud := L.NewUserData()
	ud.Value = storage
	L.SetMetatable(ud, L.GetTypeMetatable("incus_storage"))
	L.Push(ud)
	return 1
}

// createSnapshot cria um builder de snapshot
func (m *IncusModule) createSnapshot(L *lua.LState) int {
	instance := L.CheckString(1)
	name := L.CheckString(2)

	snapshot := &IncusSnapshot{
		L:           L,
		agentClient: m.agentClient,
		instance:    instance,
		name:        name,
		stateful:    false,
		target:      "",
	}

	ud := L.NewUserData()
	ud.Value = snapshot
	L.SetMetatable(ud, L.GetTypeMetatable("incus_snapshot"))
	L.Push(ud)
	return 1
}

// exec executa um comando em uma instância
func (m *IncusModule) exec(L *lua.LState) int {
	opts := L.CheckTable(1)

	instance := ""
	command := ""
	target := ""
	_ = false // interactive unused
	user := ""
	group := ""
	cwd := ""
	env := make(map[string]string)

	opts.ForEach(func(k, v lua.LValue) {
		key := k.String()
		switch key {
		case "instance":
			instance = v.String()
		case "command":
			if tbl, ok := v.(*lua.LTable); ok {
				var cmds []string
				tbl.ForEach(func(_, v lua.LValue) {
					cmds = append(cmds, v.String())
				})
				command = strings.Join(cmds, " ")
			} else {
				command = v.String()
			}
		case "target":
			target = v.String()
		case "interactive":
			_ = lua.LVAsBool(v) // interactive not yet used
		case "user":
			user = v.String()
		case "group":
			group = v.String()
		case "cwd":
			cwd = v.String()
		case "env":
			if tbl, ok := v.(*lua.LTable); ok {
				tbl.ForEach(func(k, v lua.LValue) {
					env[k.String()] = v.String()
				})
			}
		}
	})

	if instance == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("instance name is required"))
		return 2
	}

	if command == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("command is required"))
		return 2
	}

	// Construir comando incus exec
	args := []string{"exec", instance}

	if user != "" {
		args = append(args, "--user", user)
	}
	if group != "" {
		args = append(args, "--group", group)
	}
	if cwd != "" {
		args = append(args, "--cwd", cwd)
	}
	for k, v := range env {
		args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, "--", "sh", "-c", command)

	// Build command with proper quoting for sh -c
	// Need to quote the command argument to sh -c
	cmdArgs := []string{}
	for i, arg := range args {
		if i == len(args)-1 {
			// Last argument is the command for sh -c, needs to be quoted
			cmdArgs = append(cmdArgs, fmt.Sprintf("'%s'", strings.ReplaceAll(arg, "'", "'\\''")))
		} else {
			cmdArgs = append(cmdArgs, arg)
		}
	}
	fullCmd := fmt.Sprintf("incus %s", strings.Join(cmdArgs, " "))

	result, err := executeCommand(m.agentClient, fullCmd, target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// list lista recursos Incus
func (m *IncusModule) list(L *lua.LState) int {
	opts := L.CheckTable(1)

	resourceType := "instances"
	target := ""
	format := "json"

	opts.ForEach(func(k, v lua.LValue) {
		key := k.String()
		switch key {
		case "type":
			resourceType = v.String()
		case "target":
			target = v.String()
		case "format":
			format = v.String()
		}
	})

	cmd := fmt.Sprintf("incus list --format=%s", format)
	switch resourceType {
	case "images":
		cmd = fmt.Sprintf("incus image list --format=%s", format)
	case "networks":
		cmd = fmt.Sprintf("incus network list --format=%s", format)
	case "profiles":
		cmd = fmt.Sprintf("incus profile list --format=%s", format)
	case "storage":
		cmd = fmt.Sprintf("incus storage list --format=%s", format)
	}

	result, err := executeCommand(m.agentClient, cmd, target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// info obtém informações sobre um recurso
func (m *IncusModule) info(L *lua.LState) int {
	opts := L.CheckTable(1)

	resourceType := "instance"
	name := ""
	target := ""

	opts.ForEach(func(k, v lua.LValue) {
		key := k.String()
		switch key {
		case "type":
			resourceType = v.String()
		case "name":
			name = v.String()
		case "target":
			target = v.String()
		}
	})

	if name == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("name is required"))
		return 2
	}

	cmd := ""
	switch resourceType {
	case "instance":
		cmd = fmt.Sprintf("incus info %s", name)
	case "image":
		cmd = fmt.Sprintf("incus image info %s", name)
	case "network":
		cmd = fmt.Sprintf("incus network show %s", name)
	case "profile":
		cmd = fmt.Sprintf("incus profile show %s", name)
	case "storage":
		cmd = fmt.Sprintf("incus storage show %s", name)
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unknown resource type: %s", resourceType)))
		return 2
	}

	result, err := executeCommand(m.agentClient, cmd, target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// delete deleta um recurso
func (m *IncusModule) delete(L *lua.LState) int {
	opts := L.CheckTable(1)

	resourceType := "instance"
	name := ""
	target := ""
	force := false

	opts.ForEach(func(k, v lua.LValue) {
		key := k.String()
		switch key {
		case "type":
			resourceType = v.String()
		case "name":
			name = v.String()
		case "target":
			target = v.String()
		case "force":
			force = lua.LVAsBool(v)
		}
	})

	if name == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("name is required"))
		return 2
	}

	cmd := ""
	switch resourceType {
	case "instance":
		if force {
			cmd = fmt.Sprintf("incus delete %s --force", name)
		} else {
			cmd = fmt.Sprintf("incus delete %s", name)
		}
	case "image":
		cmd = fmt.Sprintf("incus image delete %s", name)
	case "network":
		cmd = fmt.Sprintf("incus network delete %s", name)
	case "profile":
		cmd = fmt.Sprintf("incus profile delete %s", name)
	case "storage":
		cmd = fmt.Sprintf("incus storage delete %s", name)
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unknown resource type: %s", resourceType)))
		return 2
	}

	result, err := executeCommand(m.agentClient, cmd, target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// RegisterInstanceMetatable registra o metatable para IncusInstance
func RegisterInstanceMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("incus_instance")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"image":       instanceImage,
		"config":      instanceConfig,
		"profile":     instanceProfile,
		"device":      instanceDevice,
		"proxy":       instanceProxy,
		"delegate_to": instanceDelegateTo,
		"ephemeral":   instanceEphemeral,
		"start":       instanceStart,
		"create":      instanceCreate,
		"launch":      instanceLaunch,
		"stop":        instanceStop,
		"restart":     instanceRestart,
		"freeze":      instanceFreeze,
		"unfreeze":    instanceUnfreeze,
		"copy":        instanceCopy,
		"move":        instanceMove,
	}))
}

func instanceImage(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	image := L.CheckString(2)
	instance.config["image"] = image
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func instanceConfig(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	config := L.CheckTable(2)

	config.ForEach(func(k, v lua.LValue) {
		instance.config[k.String()] = v.String()
	})

	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func instanceProfile(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	profile := L.CheckString(2)
	instance.profiles = append(instance.profiles, profile)
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func instanceDevice(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	name := L.CheckString(2)
	config := L.CheckTable(3)

	deviceConfig := make(map[string]string)
	config.ForEach(func(k, v lua.LValue) {
		deviceConfig[k.String()] = v.String()
	})

	instance.devices[name] = deviceConfig
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func instanceProxy(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	deviceName := L.CheckString(2)
	listen := L.CheckString(3)
	connect := L.CheckString(4)

	// Construir comando: incus config device add <instance> <device_name> proxy listen=<listen> connect=<connect>
	cmd := fmt.Sprintf("incus config device add %s %s proxy listen=%s connect=%s",
		instance.name, deviceName, listen, connect)

	_, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		// Se já existe, não é erro (idempotência)
		if !strings.Contains(err.Error(), "already exists") {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}

	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func instanceDelegateTo(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	target := L.CheckString(2)
	instance.target = target
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func instanceEphemeral(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	ephemeral := true
	if L.GetTop() >= 2 {
		ephemeral = lua.LVAsBool(L.Get(2))
	}
	instance.config["ephemeral"] = fmt.Sprintf("%t", ephemeral)
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func instanceCreate(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	return executeInstanceOperation(L, instance, "create")
}

func instanceLaunch(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	return executeInstanceOperation(L, instance, "launch")
}

func instanceStart(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	cmd := fmt.Sprintf("incus start %s", instance.name)
	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func instanceStop(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	force := false
	timeout := 30

	if L.GetTop() >= 2 {
		if tbl, ok := L.Get(2).(*lua.LTable); ok {
			tbl.ForEach(func(k, v lua.LValue) {
				switch k.String() {
				case "force":
					force = lua.LVAsBool(v)
				case "timeout":
					if num, ok := v.(lua.LNumber); ok {
						timeout = int(num)
					}
				}
			})
		}
	}

	cmd := fmt.Sprintf("incus stop %s", instance.name)
	if force {
		cmd += " --force"
	}
	cmd += fmt.Sprintf(" --timeout %d", timeout)

	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func instanceRestart(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	force := false
	timeout := 30

	if L.GetTop() >= 2 {
		if tbl, ok := L.Get(2).(*lua.LTable); ok {
			tbl.ForEach(func(k, v lua.LValue) {
				switch k.String() {
				case "force":
					force = lua.LVAsBool(v)
				case "timeout":
					if num, ok := v.(lua.LNumber); ok {
						timeout = int(num)
					}
				}
			})
		}
	}

	cmd := fmt.Sprintf("incus restart %s", instance.name)
	if force {
		cmd += " --force"
	}
	cmd += fmt.Sprintf(" --timeout %d", timeout)

	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func instanceFreeze(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	cmd := fmt.Sprintf("incus pause %s", instance.name)
	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func instanceUnfreeze(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	cmd := fmt.Sprintf("incus start %s", instance.name)
	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func instanceCopy(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	newName := L.CheckString(2)
	
	cmd := fmt.Sprintf("incus copy %s %s", instance.name, newName)
	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func instanceMove(L *lua.LState) int {
	instance := checkIncusInstance(L, 1)
	newName := L.CheckString(2)
	
	cmd := fmt.Sprintf("incus move %s %s", instance.name, newName)
	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func executeInstanceOperation(L *lua.LState, instance *IncusInstance, operation string) int {
	image, _ := instance.config["image"].(string)
	if image == "" && operation != "create" {
		L.Push(lua.LNil)
		L.Push(lua.LString("image is required"))
		return 2
	}

	// For launch/init operations, image must come before instance name
	// Syntax: incus launch <image> <instance_name> [options]
	var args []string
	if operation == "launch" || operation == "init" {
		if image != "" {
			args = []string{operation, image, instance.name}
		} else {
			args = []string{operation, instance.name}
		}
	} else {
		args = []string{operation, instance.name}
		if image != "" {
			args = append(args, image)
		}
	}

	// Add profiles
	for _, profile := range instance.profiles {
		args = append(args, "-p", profile)
	}

	// Add config
	for k, v := range instance.config {
		if k == "image" || k == "ephemeral" {
			continue
		}
		args = append(args, "-c", fmt.Sprintf("%s=%v", k, v))
	}

	// Add ephemeral flag
	if ephStr, ok := instance.config["ephemeral"].(string); ok && ephStr == "true" {
		args = append(args, "--ephemeral")
	}

	// Add devices
	for name, deviceConfig := range instance.devices {
		for k, v := range deviceConfig {
			args = append(args, "-d", fmt.Sprintf("%s,%s=%s", name, k, v))
		}
	}

	cmd := fmt.Sprintf("incus %s", strings.Join(args, " "))
	result, err := executeCommand(instance.agentClient, cmd, instance.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// RegisterImageMetatable registra o metatable para IncusImage
func RegisterImageMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("incus_image")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"source":      imageSource,
		"server":      imageServer,
		"delegate_to": imageDelegateTo,
		"copy":        imageCopy,
		"export":      imageExport,
		"refresh":     imageRefresh,
	}))
}

func imageSource(L *lua.LState) int {
	image := checkIncusImage(L, 1)
	source := L.CheckString(2)
	image.source = source
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func imageServer(L *lua.LState) int {
	image := checkIncusImage(L, 1)
	server := L.CheckString(2)
	image.server = server
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func imageDelegateTo(L *lua.LState) int {
	image := checkIncusImage(L, 1)
	target := L.CheckString(2)
	image.target = target
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func imageCopy(L *lua.LState) int {
	image := checkIncusImage(L, 1)
	
	cmd := fmt.Sprintf("incus image copy %s:%s local: --alias=%s", image.server, image.source, image.alias)
	result, err := executeCommand(image.agentClient, cmd, image.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func imageExport(L *lua.LState) int {
	image := checkIncusImage(L, 1)
	filepath := L.CheckString(2)
	
	cmd := fmt.Sprintf("incus image export %s %s", image.alias, filepath)
	result, err := executeCommand(image.agentClient, cmd, image.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func imageRefresh(L *lua.LState) int {
	image := checkIncusImage(L, 1)
	
	cmd := fmt.Sprintf("incus image refresh %s", image.alias)
	result, err := executeCommand(image.agentClient, cmd, image.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// RegisterNetworkMetatable registra o metatable para IncusNetwork
func RegisterNetworkMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("incus_network")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"type":        networkType,
		"config":      networkConfig,
		"delegate_to": networkDelegateTo,
		"create":      networkCreate,
		"attach":      networkAttach,
		"detach":      networkDetach,
	}))
}

func networkType(L *lua.LState) int {
	network := checkIncusNetwork(L, 1)
	netType := L.CheckString(2)
	network.networkType = netType
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func networkConfig(L *lua.LState) int {
	network := checkIncusNetwork(L, 1)
	config := L.CheckTable(2)

	config.ForEach(func(k, v lua.LValue) {
		network.config[k.String()] = v.String()
	})

	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func networkDelegateTo(L *lua.LState) int {
	network := checkIncusNetwork(L, 1)
	target := L.CheckString(2)
	network.target = target
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func networkCreate(L *lua.LState) int {
	network := checkIncusNetwork(L, 1)
	
	args := []string{"network", "create", network.name}
	
	for k, v := range network.config {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}

	cmd := fmt.Sprintf("incus %s", strings.Join(args, " "))
	result, err := executeCommand(network.agentClient, cmd, network.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func networkAttach(L *lua.LState) int {
	network := checkIncusNetwork(L, 1)
	instance := L.CheckString(2)
	deviceName := L.OptString(3, "eth0")
	
	cmd := fmt.Sprintf("incus network attach %s %s %s", network.name, instance, deviceName)
	result, err := executeCommand(network.agentClient, cmd, network.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func networkDetach(L *lua.LState) int {
	network := checkIncusNetwork(L, 1)
	instance := L.CheckString(2)
	deviceName := L.OptString(3, "eth0")
	
	cmd := fmt.Sprintf("incus network detach %s %s %s", network.name, instance, deviceName)
	result, err := executeCommand(network.agentClient, cmd, network.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// RegisterProfileMetatable registra o metatable para IncusProfile
func RegisterProfileMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("incus_profile")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"description": profileDescription,
		"config":      profileConfig,
		"device":      profileDevice,
		"delegate_to": profileDelegateTo,
		"create":      profileCreate,
		"apply":       profileApply,
		"copy":        profileCopy,
	}))
}

func profileDescription(L *lua.LState) int {
	profile := checkIncusProfile(L, 1)
	description := L.CheckString(2)
	profile.description = description
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func profileConfig(L *lua.LState) int {
	profile := checkIncusProfile(L, 1)
	config := L.CheckTable(2)

	config.ForEach(func(k, v lua.LValue) {
		profile.config[k.String()] = v.String()
	})

	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func profileDevice(L *lua.LState) int {
	profile := checkIncusProfile(L, 1)
	name := L.CheckString(2)
	config := L.CheckTable(3)

	deviceConfig := make(map[string]string)
	config.ForEach(func(k, v lua.LValue) {
		deviceConfig[k.String()] = v.String()
	})

	profile.devices[name] = deviceConfig
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func profileDelegateTo(L *lua.LState) int {
	profile := checkIncusProfile(L, 1)
	target := L.CheckString(2)
	profile.target = target
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func profileCreate(L *lua.LState) int {
	profile := checkIncusProfile(L, 1)
	
	// Create profile first
	cmd := fmt.Sprintf("incus profile create %s", profile.name)
	_, err := executeCommand(profile.agentClient, cmd, profile.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Set description if provided
	if profile.description != "" {
		cmd = fmt.Sprintf("incus profile set %s description='%s'", profile.name, profile.description)
		_, err = executeCommand(profile.agentClient, cmd, profile.target)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}

	// Set config
	for k, v := range profile.config {
		cmd = fmt.Sprintf("incus profile set %s %s=%v", profile.name, k, v)
		_, err = executeCommand(profile.agentClient, cmd, profile.target)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}

	// Add devices
	for name, deviceConfig := range profile.devices {
		deviceJSON, _ := json.Marshal(deviceConfig)
		cmd = fmt.Sprintf("incus profile device add %s %s %s", profile.name, name, string(deviceJSON))
		_, err = executeCommand(profile.agentClient, cmd, profile.target)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}

	L.Push(lua.LString(fmt.Sprintf("Profile %s created successfully", profile.name)))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func profileApply(L *lua.LState) int {
	profile := checkIncusProfile(L, 1)
	instance := L.CheckString(2)
	
	cmd := fmt.Sprintf("incus profile add %s %s", instance, profile.name)
	result, err := executeCommand(profile.agentClient, cmd, profile.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func profileCopy(L *lua.LState) int {
	profile := checkIncusProfile(L, 1)
	newName := L.CheckString(2)
	
	cmd := fmt.Sprintf("incus profile copy %s %s", profile.name, newName)
	result, err := executeCommand(profile.agentClient, cmd, profile.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// RegisterStorageMetatable registra o metatable para IncusStorage
func RegisterStorageMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("incus_storage")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"driver":      storageDriver,
		"config":      storageConfig,
		"delegate_to": storageDelegateTo,
		"create":      storageCreate,
		"volume":      storageVolume,
	}))
}

func storageDriver(L *lua.LState) int {
	storage := checkIncusStorage(L, 1)
	driver := L.CheckString(2)
	storage.driver = driver
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func storageConfig(L *lua.LState) int {
	storage := checkIncusStorage(L, 1)
	config := L.CheckTable(2)

	config.ForEach(func(k, v lua.LValue) {
		storage.config[k.String()] = v.String()
	})

	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func storageDelegateTo(L *lua.LState) int {
	storage := checkIncusStorage(L, 1)
	target := L.CheckString(2)
	storage.target = target
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func storageCreate(L *lua.LState) int {
	storage := checkIncusStorage(L, 1)
	
	args := []string{"storage", "create", storage.name, storage.driver}
	
	for k, v := range storage.config {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}

	cmd := fmt.Sprintf("incus %s", strings.Join(args, " "))
	result, err := executeCommand(storage.agentClient, cmd, storage.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func storageVolume(L *lua.LState) int {
	storage := checkIncusStorage(L, 1)
	volumeName := L.CheckString(2)
	
	cmd := fmt.Sprintf("incus storage volume create %s %s", storage.name, volumeName)
	result, err := executeCommand(storage.agentClient, cmd, storage.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// RegisterSnapshotMetatable registra o metatable para IncusSnapshot
func RegisterSnapshotMetatable(L *lua.LState) {
	mt := L.NewTypeMetatable("incus_snapshot")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"stateful":    snapshotStateful,
		"delegate_to": snapshotDelegateTo,
		"create":      snapshotCreate,
		"restore":     snapshotRestore,
		"delete":      snapshotDelete,
	}))
}

func snapshotStateful(L *lua.LState) int {
	snapshot := checkIncusSnapshot(L, 1)
	stateful := true
	if L.GetTop() >= 2 {
		stateful = lua.LVAsBool(L.Get(2))
	}
	snapshot.stateful = stateful
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func snapshotDelegateTo(L *lua.LState) int {
	snapshot := checkIncusSnapshot(L, 1)
	target := L.CheckString(2)
	snapshot.target = target
	L.Push(L.Get(1))
	L.Push(lua.LNil) // Fluent API também retorna (self, nil)
	return 2
}

func snapshotCreate(L *lua.LState) int {
	snapshot := checkIncusSnapshot(L, 1)
	
	cmd := fmt.Sprintf("incus snapshot %s %s", snapshot.instance, snapshot.name)
	if snapshot.stateful {
		cmd += " --stateful"
	}

	result, err := executeCommand(snapshot.agentClient, cmd, snapshot.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func snapshotRestore(L *lua.LState) int {
	snapshot := checkIncusSnapshot(L, 1)
	
	cmd := fmt.Sprintf("incus restore %s %s", snapshot.instance, snapshot.name)
	result, err := executeCommand(snapshot.agentClient, cmd, snapshot.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

func snapshotDelete(L *lua.LState) int {
	snapshot := checkIncusSnapshot(L, 1)
	
	cmd := fmt.Sprintf("incus delete %s/%s", snapshot.instance, snapshot.name)
	result, err := executeCommand(snapshot.agentClient, cmd, snapshot.target)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	L.Push(lua.LNil) // Sempre retornar (result, nil) no sucesso
	return 2
}

// Helper functions
func checkIncusInstance(L *lua.LState, n int) *IncusInstance {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*IncusInstance); ok {
		return v
	}
	L.ArgError(n, "IncusInstance expected")
	return nil
}

func checkIncusImage(L *lua.LState, n int) *IncusImage {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*IncusImage); ok {
		return v
	}
	L.ArgError(n, "IncusImage expected")
	return nil
}

func checkIncusNetwork(L *lua.LState, n int) *IncusNetwork {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*IncusNetwork); ok {
		return v
	}
	L.ArgError(n, "IncusNetwork expected")
	return nil
}

func checkIncusProfile(L *lua.LState, n int) *IncusProfile {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*IncusProfile); ok {
		return v
	}
	L.ArgError(n, "IncusProfile expected")
	return nil
}

func checkIncusStorage(L *lua.LState, n int) *IncusStorage {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*IncusStorage); ok {
		return v
	}
	L.ArgError(n, "IncusStorage expected")
	return nil
}

func checkIncusSnapshot(L *lua.LState, n int) *IncusSnapshot {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*IncusSnapshot); ok {
		return v
	}
	L.ArgError(n, "IncusSnapshot expected")
	return nil
}

func executeCommand(agentClient interface{}, cmd string, target string) (string, error) {
	// When running on agent via delegate_to, target should be empty
	// Execute locally using bash
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Get environment variables from parent process or use defaults
	xdgRuntime := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntime == "" {
		xdgRuntime = "/run/user/1000"
	}

	dbusAddr := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	if dbusAddr == "" {
		dbusAddr = "unix:path=/run/user/1000/bus"
	}

	home := os.Getenv("HOME")
	if home == "" {
		home = "/home/chalkan3"
	}

	// Prepend export statements to the command to ensure environment variables
	// are available to incus within the bash subprocess
	// HOME is needed for Incus to find its configuration in ~/.config/incus/
	wrappedCmd := fmt.Sprintf("export HOME='%s' && export XDG_RUNTIME_DIR='%s' && export DBUS_SESSION_BUS_ADDRESS='%s' && %s",
		home, xdgRuntime, dbusAddr, cmd)

	execCmd := exec.CommandContext(ctx, "bash", "-c", wrappedCmd)

	// Also set the environment on the process for completeness
	execCmd.Env = os.Environ()

	output, err := execCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v - %s", err, string(output))
	}
	return string(output), nil
}
