package luainterface

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// ComprehensiveSaltModule provides complete SaltStack integration
type ComprehensiveSaltModule struct{}

// NewComprehensiveSaltModule creates a new comprehensive salt module
func NewComprehensiveSaltModule() *ComprehensiveSaltModule {
	return &ComprehensiveSaltModule{}
}

// Loader returns the Lua loader for the comprehensive salt module
func (mod *ComprehensiveSaltModule) Loader(L *lua.LState) int {
	saltTable := L.NewTable()
	
	// Core execution functions
	L.SetFuncs(saltTable, map[string]lua.LGFunction{
		// Basic execution
		"cmd":              mod.saltCmd,
		"run":              mod.saltRun,
		"execute":          mod.saltExecute,
		"batch":            mod.saltBatch,
		"async":            mod.saltAsync,
		
		// Connection and testing
		"ping":             mod.saltPing,
		"test":             mod.saltTest,
		"version":          mod.saltVersion,
		"status":           mod.saltStatus,
		
		// Key management
		"key_list":         mod.saltKeyList,
		"key_accept":       mod.saltKeyAccept,
		"key_reject":       mod.saltKeyReject,
		"key_delete":       mod.saltKeyDelete,
		"key_finger":       mod.saltKeyFinger,
		"key_gen":          mod.saltKeyGen,
		
		// State management
		"state_apply":      mod.saltStateApply,
		"state_highstate":  mod.saltStateHighstate,
		"state_test":       mod.saltStateTest,
		"state_show_sls":   mod.saltStateShowSls,
		"state_show_top":   mod.saltStateShowTop,
		"state_show_lowstate": mod.saltStateShowLowstate,
		"state_single":     mod.saltStateSingle,
		"state_template":   mod.saltStateTemplate,
		
		// Grains management
		"grains_get":       mod.saltGrainsGet,
		"grains_set":       mod.saltGrainsSet,
		"grains_append":    mod.saltGrainsAppend,
		"grains_remove":    mod.saltGrainsRemove,
		"grains_delkey":    mod.saltGrainsDelkey,
		"grains_items":     mod.saltGrainsItems,
		
		// Pillar management
		"pillar_get":       mod.saltPillarGet,
		"pillar_items":     mod.saltPillarItems,
		"pillar_show":      mod.saltPillarShow,
		"pillar_refresh":   mod.saltPillarRefresh,
		
		// File operations
		"file_copy":        mod.saltFileCopy,
		"file_get":         mod.saltFileGet,
		"file_list":        mod.saltFileList,
		"file_manage":      mod.saltFileManage,
		"file_recurse":     mod.saltFileRecurse,
		"file_touch":       mod.saltFileTouch,
		"file_stats":       mod.saltFileStats,
		"file_find":        mod.saltFileFind,
		"file_replace":     mod.saltFileReplace,
		"file_check_hash":  mod.saltFileCheckHash,
		
		// Package management
		"pkg_install":      mod.saltPkgInstall,
		"pkg_remove":       mod.saltPkgRemove,
		"pkg_upgrade":      mod.saltPkgUpgrade,
		"pkg_refresh":      mod.saltPkgRefresh,
		"pkg_list":         mod.saltPkgList,
		"pkg_version":      mod.saltPkgVersion,
		"pkg_available":    mod.saltPkgAvailable,
		"pkg_info":         mod.saltPkgInfo,
		"pkg_hold":         mod.saltPkgHold,
		"pkg_unhold":       mod.saltPkgUnhold,
		
		// Service management
		"service_start":    mod.saltServiceStart,
		"service_stop":     mod.saltServiceStop,
		"service_restart":  mod.saltServiceRestart,
		"service_reload":   mod.saltServiceReload,
		"service_status":   mod.saltServiceStatus,
		"service_enable":   mod.saltServiceEnable,
		"service_disable":  mod.saltServiceDisable,
		"service_list":     mod.saltServiceList,
		
		// User management
		"user_add":         mod.saltUserAdd,
		"user_delete":      mod.saltUserDelete,
		"user_info":        mod.saltUserInfo,
		"user_list":        mod.saltUserList,
		"user_chuid":       mod.saltUserChuid,
		"user_chgid":       mod.saltUserChgid,
		"user_chshell":     mod.saltUserChshell,
		"user_chhome":      mod.saltUserChhome,
		"user_primary_group": mod.saltUserPrimaryGroup,
		
		// Group management
		"group_add":        mod.saltGroupAdd,
		"group_delete":     mod.saltGroupDelete,
		"group_info":       mod.saltGroupInfo,
		"group_list":       mod.saltGroupList,
		"group_adduser":    mod.saltGroupAdduser,
		"group_deluser":    mod.saltGroupDeluser,
		"group_members":    mod.saltGroupMembers,
		
		// Network management
		"network_interface": mod.saltNetworkInterface,
		"network_interfaces": mod.saltNetworkInterfaces,
		"network_ping":     mod.saltNetworkPing,
		"network_traceroute": mod.saltNetworkTraceroute,
		"network_netstat":  mod.saltNetworkNetstat,
		"network_arp":      mod.saltNetworkArp,
		
		// System information
		"system_info":      mod.saltSystemInfo,
		"system_uptime":    mod.saltSystemUptime,
		"system_reboot":    mod.saltSystemReboot,
		"system_shutdown":  mod.saltSystemShutdown,
		"system_halt":      mod.saltSystemHalt,
		"system_hostname":  mod.saltSystemHostname,
		"system_set_hostname": mod.saltSystemSetHostname,
		
		// Disk and mount management
		"disk_usage":       mod.saltDiskUsage,
		"disk_stats":       mod.saltDiskStats,
		"mount_active":     mod.saltMountActive,
		"mount_fstab":      mod.saltMountFstab,
		"mount_mount":      mod.saltMountMount,
		"mount_umount":     mod.saltMountUmount,
		"mount_remount":    mod.saltMountRemount,
		
		// Process management
		"process_list":     mod.saltProcessList,
		"process_info":     mod.saltProcessInfo,
		"process_kill":     mod.saltProcessKill,
		"process_killall":  mod.saltProcessKillall,
		"process_pkill":    mod.saltProcessPkill,
		
		// Cron management
		"cron_list":        mod.saltCronList,
		"cron_set":         mod.saltCronSet,
		"cron_delete":      mod.saltCronDelete,
		"cron_raw_cron":    mod.saltCronRawCron,
		
		// Archive operations
		"archive_gunzip":   mod.saltArchiveGunzip,
		"archive_gzip":     mod.saltArchiveGzip,
		"archive_tar":      mod.saltArchiveTar,
		"archive_untar":    mod.saltArchiveUntar,
		"archive_unzip":    mod.saltArchiveUnzip,
		"archive_zip":      mod.saltArchiveZip,
		
		// Salt-cloud integration
		"cloud_list_nodes": mod.saltCloudListNodes,
		"cloud_create":     mod.saltCloudCreate,
		"cloud_destroy":    mod.saltCloudDestroy,
		"cloud_action":     mod.saltCloudAction,
		"cloud_function":   mod.saltCloudFunction,
		"cloud_map":        mod.saltCloudMap,
		"cloud_profile":    mod.saltCloudProfile,
		"cloud_provider":   mod.saltCloudProvider,
		
		// Event system
		"event_send":       mod.saltEventSend,
		"event_listen":     mod.saltEventListen,
		"event_fire":       mod.saltEventFire,
		"event_fire_master": mod.saltEventFireMaster,
		
		// Orchestration
		"orchestrate":      mod.saltOrchestrate,
		"runner":           mod.saltRunner,
		"wheel":            mod.saltWheel,
		
		// Mine operations
		"mine_get":         mod.saltMineGet,
		"mine_send":        mod.saltMineSend,
		"mine_update":      mod.saltMineUpdate,
		"mine_delete":      mod.saltMineDelete,
		"mine_flush":       mod.saltMineFlush,
		"mine_valid":       mod.saltMineValid,
		
		// Job management
		"job_active":       mod.saltJobActive,
		"job_list":         mod.saltJobList,
		"job_lookup":       mod.saltJobLookup,
		"job_exit_success": mod.saltJobExitSuccess,
		"job_print":        mod.saltJobPrint,
		
		// Advanced features
		"orchestrate_sls":  mod.saltOrchestrateSls,
		"syndic_list":      mod.saltSyndicList,
		"reactor_list":     mod.saltReactorList,
		"cache_grains":     mod.saltCacheGrains,
		"ssh_run":          mod.saltSSHRun,
		"proxy_ping":       mod.saltProxyPing,
		"vault_read":       mod.saltVaultRead,
		"docker_ps":        mod.saltDockerPs,
		"git_clone":        mod.saltGitClone,
		"mysql_query":      mod.saltMysqlQuery,
		"status_loadavg":   mod.saltStatusLoadavg,
		"config_get":       mod.saltConfigGet,
		"api_client":       mod.saltApiClient,
		"template_jinja":   mod.saltTemplateJinja,
		"log_info":         mod.saltLogInfo,
		"beacon_list":      mod.saltBeaconList,
		"schedule_list":    mod.saltScheduleList,
	})
	
	L.Push(saltTable)
	return 1
}

// Helper function to execute salt commands with enhanced error handling and features
func (mod *ComprehensiveSaltModule) executeSaltCommand(timeout int, retries int, cmdArgs ...string) (map[string]interface{}, error) {
	var lastErr error
	
	for i := 0; i <= retries; i++ {
		result, err := mod.executeSingleSaltCommand(timeout, cmdArgs...)
		if err == nil {
			return result, nil
		}
		lastErr = err
		
		if i < retries {
			// Exponential backoff
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	
	return nil, fmt.Errorf("command failed after %d retries: %w", retries, lastErr)
}

func (mod *ComprehensiveSaltModule) executeSingleSaltCommand(timeout int, cmdArgs ...string) (map[string]interface{}, error) {
	// Check if salt command exists
	if _, err := exec.LookPath(cmdArgs[0]); err != nil {
		return nil, fmt.Errorf("salt command '%s' not found in PATH: %w", cmdArgs[0], err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	
	// Set environment variables
	cmd.Env = os.Environ()
	
	// Set working directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		if _, err := os.Stat(filepath.Join(homeDir, ".saltrc")); err == nil {
			cmd.Dir = homeDir
		}
	}
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)
	
	result := map[string]interface{}{
		"success":     err == nil,
		"stdout":      stdout.String(),
		"stderr":      stderr.String(),
		"ret_code":    0,
		"duration_ms": duration.Milliseconds(),
		"command":     strings.Join(cmdArgs, " "),
		"timeout":     timeout,
	}
	
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result["ret_code"] = exitError.ExitCode()
		}
		result["error"] = err.Error()
		
		// Try to parse stderr as JSON if available
		if stderr.Len() > 0 {
			result["stderr"] = stderr.String()
		}
		
		return result, fmt.Errorf("salt command failed: %w", err)
	}
	
	// Try to parse stdout as JSON if it looks like JSON
	stdoutStr := stdout.String()
	if strings.TrimSpace(stdoutStr) != "" {
		var jsonResult interface{}
		if err := json.Unmarshal([]byte(stdoutStr), &jsonResult); err == nil {
			result["returns"] = jsonResult
		} else {
			result["output"] = stdoutStr
		}
	}
	
	return result, nil
}

// Convert Go values to Lua values
func (mod *ComprehensiveSaltModule) goValueToLua(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case map[string]interface{}:
		table := L.NewTable()
		for k, val := range v {
			table.RawSetString(k, mod.goValueToLua(L, val))
		}
		return table
	case []interface{}:
		table := L.NewTable()
		for i, val := range v {
			table.RawSetInt(i+1, mod.goValueToLua(L, val))
		}
		return table
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

// Helper function to return consistent Salt results
func (mod *ComprehensiveSaltModule) returnSaltResult(L *lua.LState, result map[string]interface{}, err error) int {
	luaResult := L.NewTable()
	
	if result != nil {
		for k, v := range result {
			luaResult.RawSetString(k, mod.goValueToLua(L, v))
		}
	}
	
	if err != nil {
		luaResult.RawSetString("success", lua.LBool(false))
		luaResult.RawSetString("error", lua.LString(err.Error()))
	}
	
	L.Push(luaResult)
	return 1
}

// Basic execution functions
func (mod *ComprehensiveSaltModule) saltCmd(L *lua.LState) int {
	return mod.executeSaltFunction(L, "salt")
}

func (mod *ComprehensiveSaltModule) saltRun(L *lua.LState) int {
	return mod.executeSaltFunction(L, "salt-run")
}

func (mod *ComprehensiveSaltModule) saltExecute(L *lua.LState) int {
	return mod.executeSaltFunction(L, "salt")
}

func (mod *ComprehensiveSaltModule) saltBatch(L *lua.LState) int {
	target := L.CheckString(1)
	batchSize := L.CheckString(2)
	module := L.CheckString(3)
	function := L.CheckString(4)
	
	args := []string{"salt", "--batch-size=" + batchSize, target, module + "." + function, "--out=json"}
	
	// Add any additional arguments
	for i := 5; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(300, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltAsync(L *lua.LState) int {
	target := L.CheckString(1)
	module := L.CheckString(2)
	function := L.CheckString(3)
	
	args := []string{"salt", target, module + "." + function, "--async", "--out=json"}
	
	// Add any additional arguments
	for i := 4; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

// Generic function executor with comprehensive options
func (mod *ComprehensiveSaltModule) executeSaltFunction(L *lua.LState, command string) int {
	args := []string{command}
	
	// Get all arguments
	for i := 1; i <= L.GetTop(); i++ {
		arg := L.Get(i)
		switch arg.Type() {
		case lua.LTString:
			args = append(args, arg.String())
		case lua.LTTable:
			// Handle options table
			table := arg.(*lua.LTable)
			if output := table.RawGetString("output"); output != lua.LNil {
				args = append(args, "--out="+output.String())
			}
			if timeout := table.RawGetString("timeout"); timeout != lua.LNil {
				args = append(args, "--timeout="+timeout.String())
			}
			if batch := table.RawGetString("batch"); batch != lua.LNil {
				args = append(args, "--batch-size="+batch.String())
			}
			if async := table.RawGetString("async"); lua.LVAsBool(async) {
				args = append(args, "--async")
			}
			if test := table.RawGetString("test"); lua.LVAsBool(test) {
				args = append(args, "test=True")
			}
			if pillar := table.RawGetString("pillar"); pillar != lua.LNil {
				args = append(args, "pillar='"+pillar.String()+"'")
			}
		default:
			args = append(args, lua.LVAsString(arg))
		}
	}
	
	// Default options
	timeout := 60
	retries := 3
	
	result, err := mod.executeSaltCommand(timeout, retries, args...)
	
	// Convert result to Lua table
	luaResult := L.NewTable()
	if result != nil {
		for k, v := range result {
			luaResult.RawSetString(k, mod.goValueToLua(L, v))
		}
	}
	
	if err != nil {
		luaResult.RawSetString("success", lua.LBool(false))
		luaResult.RawSetString("error", lua.LString(err.Error()))
	}
	
	L.Push(luaResult)
	return 1
}

// Connection and testing functions
func (mod *ComprehensiveSaltModule) saltPing(L *lua.LState) int {
	target := L.CheckString(1)
	opts := L.OptTable(2, L.NewTable())
	
	args := []string{"salt", target, "test.ping", "--out=json"}
	
	// Add timeout if specified
	if timeout := opts.RawGetString("timeout"); timeout != lua.LNil {
		args = append(args, "--timeout="+timeout.String())
	}
	
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltTest(L *lua.LState) int {
	target := L.CheckString(1)
	testType := L.CheckString(2) // ping, version, fib, etc.
	
	args := []string{"salt", target, "test." + testType, "--out=json"}
	
	// Add any additional arguments
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckString(i))
	}
	
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltVersion(L *lua.LState) int {
	target := "*"
	if L.GetTop() > 0 {
		target = L.CheckString(1)
	}
	
	args := []string{"salt", target, "test.version", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ComprehensiveSaltModule) saltStatus(L *lua.LState) int {
	target := "*"
	if L.GetTop() > 0 {
		target = L.CheckString(1)
	}
	
	args := []string{"salt", target, "status.uptime", "--out=json"}
	result, err := mod.executeSaltCommand(60, 3, args...)
	return mod.returnSaltResult(L, result, err)
}