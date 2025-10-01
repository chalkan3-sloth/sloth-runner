package luainterface

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// PulumiModule provides Pulumi infrastructure management functionality
type PulumiModule struct{}

// NewPulumiModule creates a new PulumiModule
func NewPulumiModule() *PulumiModule {
	return &PulumiModule{}
}

// PulumiClient represents a Pulumi client with context
type PulumiClient struct {
	module     *PulumiModule
	workdir    string
	backend    string
	loggedIn   bool
	currentStack string
}

// Loader returns the Lua loader for the pulumi module
func (mod *PulumiModule) Loader(L *lua.LState) int {
	// Create pulumi module table
	pulumiTable := L.NewTable()
	
	// Factory method: pulumi.login(backend, options) returns client
	L.SetField(pulumiTable, "login", L.NewFunction(func(L *lua.LState) int {
		backend := L.CheckString(1)
		options := L.OptTable(2, L.NewTable())
		
		// Create pulumi client object
		pulumiClient := L.NewUserData()
		client := &PulumiClient{
			module:  mod,
			backend: backend,
			workdir: "",
		}
		
		// Handle login options
		localLogin := true
		if localOpt := options.RawGetString("login_local"); localOpt != lua.LNil {
			localLogin = lua.LVAsBool(localOpt)
		}
		
		// Perform login
		var err error
		if localLogin {
			err = mod.pulumiCommand("login", "--local")
		} else {
			if backend == "urllogin" {
				err = mod.pulumiCommand("login")
			} else {
				err = mod.pulumiCommand("login", backend)
			}
		}
		
		if err != nil {
			// Return error table instead of client
			errorTable := L.NewTable()
			errorTable.RawSetString("error", lua.LString(err.Error()))
			pulumiClient.Value = errorTable
		} else {
			client.loggedIn = true
			pulumiClient.Value = client
		}
		
		// Create metatable for pulumi client with fluent methods
		pulumiMt := L.NewTypeMetatable("PulumiClient")
		L.SetField(pulumiMt, "__index", L.NewFunction(func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			method := L.CheckString(2)
			
			// Check if this is an error object
			if errorObj, ok := ud.Value.(*lua.LTable); ok {
				if method == "error" {
					if errMsg := errorObj.RawGetString("error"); errMsg != lua.LNil {
						L.Push(errMsg)
						return 1
					}
				}
				L.Push(lua.LNil)
				return 1
			}
			
			client, ok := ud.Value.(*PulumiClient)
			if !ok {
				L.ArgError(1, "PulumiClient expected")
				return 0
			}
			
			switch method {
			case "stack":
				L.Push(L.NewFunction(client.stack))
			case "preview":
				L.Push(L.NewFunction(client.preview))
			case "up":
				L.Push(L.NewFunction(client.up))
			case "destroy":
				L.Push(L.NewFunction(client.destroy))
			case "refresh":
				L.Push(L.NewFunction(client.refresh))
			case "set_config":
				L.Push(L.NewFunction(client.setConfig))
			case "get_config":
				L.Push(L.NewFunction(client.getConfig))
			case "set_workdir":
				L.Push(L.NewFunction(client.setWorkdir))
			default:
				L.Push(lua.LNil)
			}
			return 1
		}))
		
		L.SetMetatable(pulumiClient, pulumiMt)
		L.Push(pulumiClient)
		return 1
	}))
	
	L.Push(pulumiTable)
	return 1
}

// pulumiCommand executes a pulumi command
func (mod *PulumiModule) pulumiCommand(args ...string) error {
	cmd := exec.Command("pulumi", args...)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pulumi %s failed: %s", strings.Join(args, " "), stderr.String())
	}
	
	return nil
}

// pulumiCommandWithOutput executes a pulumi command and returns output
func (mod *PulumiModule) pulumiCommandWithOutput(args ...string) (string, error) {
	cmd := exec.Command("pulumi", args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pulumi %s failed: %s", strings.Join(args, " "), stderr.String())
	}
	
	return stdout.String(), nil
}

// setWorkdir sets the working directory for the client
func (client *PulumiClient) setWorkdir(L *lua.LState) int {
	workdir := L.CheckString(2)
	client.workdir = workdir
	
	L.Push(lua.LBool(true))
	return 1
}

// stack manages Pulumi stacks
func (client *PulumiClient) stack(L *lua.LState) int {
	stackName := L.CheckString(2)
	options := L.OptTable(3, L.NewTable())
	
	// Change to working directory if set
	originalDir := ""
	if client.workdir != "" {
		var err error
		originalDir, err = os.Getwd()
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to get current directory: %v", err)))
			return 2
		}
		
		err = os.Chdir(client.workdir)
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to change directory: %v", err)))
			return 2
		}
		
		defer func() {
			if originalDir != "" {
				os.Chdir(originalDir)
			}
		}()
	}
	
	// Check if we should create the stack
	create := false
	if createOpt := options.RawGetString("create"); createOpt != lua.LNil {
		create = lua.LVAsBool(createOpt)
	}
	
	if create {
		// First, try to select existing stack
		selectCmd := exec.Command("pulumi", "stack", "select", stackName)
		selectCmd.Stderr = nil // Ignore stderr for this check
		selectErr := selectCmd.Run()
		
		if selectErr != nil {
			// Stack doesn't exist, try to create it
			createCmd := exec.Command("pulumi", "stack", "init", stackName)
			var createStderr bytes.Buffer
			createCmd.Stderr = &createStderr
			
			createErr := createCmd.Run()
			if createErr != nil {
				// If creation fails, check if it's because stack already exists
				if strings.Contains(createStderr.String(), "already exists") {
					// Stack exists but select failed, try select again
					retrySelectCmd := exec.Command("pulumi", "stack", "select", stackName)
					var retryStderr bytes.Buffer
					retrySelectCmd.Stderr = &retryStderr
					
					if retrySelectErr := retrySelectCmd.Run(); retrySelectErr != nil {
						L.Push(lua.LBool(false))
						L.Push(lua.LString(fmt.Sprintf("Failed to select existing stack: %s", retryStderr.String())))
						return 2
					}
				} else {
					L.Push(lua.LBool(false))
					L.Push(lua.LString(fmt.Sprintf("Failed to create stack: %s", createStderr.String())))
					return 2
				}
			}
		}
	} else {
		// Just try to select the stack
		selectCmd := exec.Command("pulumi", "stack", "select", stackName)
		var selectStderr bytes.Buffer
		selectCmd.Stderr = &selectStderr
		
		if selectErr := selectCmd.Run(); selectErr != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to select stack '%s': %s", stackName, selectStderr.String())))
			return 2
		}
	}
	
	client.currentStack = stackName
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("Stack selected: %s", stackName)))
	return 2
}

// setConfig sets configuration values for the current stack
func (client *PulumiClient) setConfig(L *lua.LState) int {
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	// Change to working directory if set
	originalDir := ""
	if client.workdir != "" {
		var err error
		originalDir, err = os.Getwd()
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to get current directory: %v", err)))
			return 2
		}
		
		err = os.Chdir(client.workdir)
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to change directory: %v", err)))
			return 2
		}
		
		defer func() {
			if originalDir != "" {
				os.Chdir(originalDir)
			}
		}()
	}
	
	err := client.module.pulumiCommand("config", "set", key, value)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(fmt.Sprintf("Config set: %s = %s", key, value)))
	return 2
}

// getConfig gets configuration values for the current stack
func (client *PulumiClient) getConfig(L *lua.LState) int {
	key := L.CheckString(2)
	
	// Change to working directory if set
	originalDir := ""
	if client.workdir != "" {
		var err error
		originalDir, err = os.Getwd()
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("Failed to get current directory: %v", err)))
			return 2
		}
		
		err = os.Chdir(client.workdir)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("Failed to change directory: %v", err)))
			return 2
		}
		
		defer func() {
			if originalDir != "" {
				os.Chdir(originalDir)
			}
		}()
	}
	
	output, err := client.module.pulumiCommandWithOutput("config", "get", key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(strings.TrimSpace(output)))
	L.Push(lua.LNil)
	return 2
}

// preview runs pulumi preview
func (client *PulumiClient) preview(L *lua.LState) int {
	options := L.OptTable(2, L.NewTable())
	
	// Change to working directory if set
	originalDir := ""
	if client.workdir != "" {
		var err error
		originalDir, err = os.Getwd()
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to get current directory: %v", err)))
			return 2
		}
		
		err = os.Chdir(client.workdir)
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to change directory: %v", err)))
			return 2
		}
		
		defer func() {
			if originalDir != "" {
				os.Chdir(originalDir)
			}
		}()
	}
	
	args := []string{"preview"}
	
	// Add additional options
	if diffOpt := options.RawGetString("diff"); lua.LVAsBool(diffOpt) {
		args = append(args, "--diff")
	}
	
	// Handle environment variables
	cmd := exec.Command("pulumi", args...)
	
	// Set environment variables if provided
	if envs := options.RawGetString("envs"); envs != lua.LNil {
		if envTable, ok := envs.(*lua.LTable); ok {
			var environ []string
			environ = append(environ, os.Environ()...)
			
			envTable.ForEach(func(key, value lua.LValue) {
				environ = append(environ, fmt.Sprintf("%s=%s", key.String(), value.String()))
			})
			
			cmd.Env = environ
		}
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("pulumi preview failed: %s", stderr.String())))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(stdout.String()))
	return 2
}

// up runs pulumi up
func (client *PulumiClient) up(L *lua.LState) int {
	options := L.OptTable(2, L.NewTable())
	
	// Change to working directory if set
	originalDir := ""
	if client.workdir != "" {
		var err error
		originalDir, err = os.Getwd()
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to get current directory: %v", err)))
			return 2
		}
		
		err = os.Chdir(client.workdir)
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to change directory: %v", err)))
			return 2
		}
		
		defer func() {
			if originalDir != "" {
				os.Chdir(originalDir)
			}
		}()
	}
	
	args := []string{"up"}
	
	// Add auto-approve if specified
	if autoApprove := options.RawGetString("auto_approve"); lua.LVAsBool(autoApprove) {
		args = append(args, "--yes")
	}
	
	// Handle environment variables
	cmd := exec.Command("pulumi", args...)
	
	// Set environment variables if provided
	if envs := options.RawGetString("envs"); envs != lua.LNil {
		if envTable, ok := envs.(*lua.LTable); ok {
			var environ []string
			environ = append(environ, os.Environ()...)
			
			envTable.ForEach(func(key, value lua.LValue) {
				environ = append(environ, fmt.Sprintf("%s=%s", key.String(), value.String()))
			})
			
			cmd.Env = environ
		}
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("pulumi up failed: %s", stderr.String())))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(stdout.String()))
	return 2
}

// destroy runs pulumi destroy
func (client *PulumiClient) destroy(L *lua.LState) int {
	options := L.OptTable(2, L.NewTable())
	
	// Change to working directory if set
	originalDir := ""
	if client.workdir != "" {
		var err error
		originalDir, err = os.Getwd()
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to get current directory: %v", err)))
			return 2
		}
		
		err = os.Chdir(client.workdir)
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to change directory: %v", err)))
			return 2
		}
		
		defer func() {
			if originalDir != "" {
				os.Chdir(originalDir)
			}
		}()
	}
	
	args := []string{"destroy"}
	
	// Add auto-approve if specified
	if autoApprove := options.RawGetString("auto_approve"); lua.LVAsBool(autoApprove) {
		args = append(args, "--yes")
	}
	
	// Handle environment variables
	cmd := exec.Command("pulumi", args...)
	
	// Set environment variables if provided
	if envs := options.RawGetString("envs"); envs != lua.LNil {
		if envTable, ok := envs.(*lua.LTable); ok {
			var environ []string
			environ = append(environ, os.Environ()...)
			
			envTable.ForEach(func(key, value lua.LValue) {
				environ = append(environ, fmt.Sprintf("%s=%s", key.String(), value.String()))
			})
			
			cmd.Env = environ
		}
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("pulumi destroy failed: %s", stderr.String())))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(stdout.String()))
	return 2
}

// refresh runs pulumi refresh
func (client *PulumiClient) refresh(L *lua.LState) int {
	options := L.OptTable(2, L.NewTable())
	
	// Change to working directory if set
	originalDir := ""
	if client.workdir != "" {
		var err error
		originalDir, err = os.Getwd()
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to get current directory: %v", err)))
			return 2
		}
		
		err = os.Chdir(client.workdir)
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(fmt.Sprintf("Failed to change directory: %v", err)))
			return 2
		}
		
		defer func() {
			if originalDir != "" {
				os.Chdir(originalDir)
			}
		}()
	}
	
	args := []string{"refresh"}
	
	// Add auto-approve if specified
	if autoApprove := options.RawGetString("auto_approve"); lua.LVAsBool(autoApprove) {
		args = append(args, "--yes")
	}
	
	// Handle environment variables
	cmd := exec.Command("pulumi", args...)
	
	// Set environment variables if provided
	if envs := options.RawGetString("envs"); envs != lua.LNil {
		if envTable, ok := envs.(*lua.LTable); ok {
			var environ []string
			environ = append(environ, os.Environ()...)
			
			envTable.ForEach(func(key, value lua.LValue) {
				environ = append(environ, fmt.Sprintf("%s=%s", key.String(), value.String()))
			})
			
			cmd.Env = environ
		}
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("pulumi refresh failed: %s", stderr.String())))
		return 2
	}
	
	L.Push(lua.LBool(true))
	L.Push(lua.LString(stdout.String()))
	return 2
}

// PulumiLoader is the global loader function
func PulumiLoader(L *lua.LState) int {
	return NewPulumiModule().Loader(L)
}