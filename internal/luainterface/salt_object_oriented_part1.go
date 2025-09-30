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

// SaltClient represents a Salt client object with all methods
type SaltClient struct {
	Config map[string]interface{}
}

// NewSaltClient creates a new SaltClient
func NewSaltClient(config map[string]interface{}) *SaltClient {
	if config == nil {
		config = make(map[string]interface{})
	}
	return &SaltClient{Config: config}
}

// ObjectOrientedSaltModule provides the Salt module as an object
type ObjectOrientedSaltModule struct{}

// NewObjectOrientedSaltModule creates a new object-oriented salt module
func NewObjectOrientedSaltModule() *ObjectOrientedSaltModule {
	return &ObjectOrientedSaltModule{}
}

// Loader returns the Lua loader for the object-oriented salt module
func (mod *ObjectOrientedSaltModule) Loader(L *lua.LState) int {
	// Create the main salt table/constructor
	saltConstructor := L.NewFunction(mod.saltConstructor)
	L.Push(saltConstructor)
	return 1
}

// saltConstructor creates a new Salt client object
func (mod *ObjectOrientedSaltModule) saltConstructor(L *lua.LState) int {
	// Parse configuration options
	config := make(map[string]interface{})
	
	if L.GetTop() > 0 {
		configTable := L.OptTable(1, L.NewTable())
		configTable.ForEach(func(key, value lua.LValue) {
			config[key.String()] = mod.luaValueToGo(value)
		})
	}
	
	// Create the Salt client object
	saltClient := L.NewTable()
	
	// Store configuration
	saltClient.RawSetString("config", mod.goValueToLua(L, config))
	
	// Core execution methods
	L.SetFuncs(saltClient, map[string]lua.LGFunction{
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
		
		// Docker integration
		"docker_ps":        mod.saltDockerPs,
		"docker_run":       mod.saltDockerRun,
		"docker_stop":      mod.saltDockerStop,
		"docker_start":     mod.saltDockerStart,
		"docker_restart":   mod.saltDockerRestart,
		"docker_build":     mod.saltDockerBuild,
		"docker_pull":      mod.saltDockerPull,
		"docker_push":      mod.saltDockerPush,
		"docker_images":    mod.saltDockerImages,
		"docker_remove":    mod.saltDockerRemove,
		"docker_inspect":   mod.saltDockerInspect,
		"docker_logs":      mod.saltDockerLogs,
		"docker_exec":      mod.saltDockerExec,
		
		// Git operations
		"git_clone":        mod.saltGitClone,
		"git_pull":         mod.saltGitPull,
		"git_checkout":     mod.saltGitCheckout,
		"git_add":          mod.saltGitAdd,
		"git_commit":       mod.saltGitCommit,
		"git_push":         mod.saltGitPush,
		"git_status":       mod.saltGitStatus,
		"git_log":          mod.saltGitLog,
		"git_reset":        mod.saltGitReset,
		"git_remote_get":   mod.saltGitRemoteGet,
		"git_remote_set":   mod.saltGitRemoteSet,
		
		// Database operations
		"mysql_query":      mod.saltMysqlQuery,
		"mysql_db_create":  mod.saltMysqlDbCreate,
		"mysql_db_remove":  mod.saltMysqlDbRemove,
		"mysql_user_create": mod.saltMysqlUserCreate,
		"mysql_user_remove": mod.saltMysqlUserRemove,
		"mysql_grant_add":  mod.saltMysqlGrantAdd,
		"mysql_grant_revoke": mod.saltMysqlGrantRevoke,
		"postgres_query":   mod.saltPostgresQuery,
		"postgres_db_create": mod.saltPostgresDbCreate,
		"postgres_db_remove": mod.saltPostgresDbRemove,
		"postgres_user_create": mod.saltPostgresUserCreate,
		"postgres_user_remove": mod.saltPostgresUserRemove,
		
		// Monitoring and metrics
		"status_loadavg":   mod.saltStatusLoadavg,
		"status_cpuinfo":   mod.saltStatusCpuinfo,
		"status_meminfo":   mod.saltStatusMeminfo,
		"status_diskusage": mod.saltStatusDiskusage,
		"status_netdev":    mod.saltStatusNetdev,
		"status_w":         mod.saltStatusW,
		"status_uptime":    mod.saltStatusUptime,
		
		// Configuration management
		"config_get":       mod.saltConfigGet,
		"config_option":    mod.saltConfigOption,
		"config_valid_fileproto": mod.saltConfigValidFileproto,
		"config_backup_mode": mod.saltConfigBackupMode,
		
		// API integration
		"api_login":        mod.saltApiLogin,
		"api_logout":       mod.saltApiLogout,
		"api_minions":      mod.saltApiMinions,
		"api_jobs":         mod.saltApiJobs,
		"api_stats":        mod.saltApiStats,
		"api_events":       mod.saltApiEvents,
		"api_hook":         mod.saltApiHook,
		
		// Template engines
		"template_jinja":   mod.saltTemplateJinja,
		"template_yaml":    mod.saltTemplateYaml,
		"template_json":    mod.saltTemplateJson,
		"template_mako":    mod.saltTemplateMako,
		"template_py":      mod.saltTemplatePy,
		"template_wempy":   mod.saltTemplateWempy,
		
		// Logging and debugging
		"log_error":        mod.saltLogError,
		"log_warning":      mod.saltLogWarning,
		"log_info":         mod.saltLogInfo,
		"log_debug":        mod.saltLogDebug,
		"debug_mode":       mod.saltDebugMode,
		"debug_profile":    mod.saltDebugProfile,
		
		// Beacons management (commented out until implemented)
		// "beacon_list":      mod.saltBeaconList,
		// "beacon_add":       mod.saltBeaconAdd,
		// "beacon_modify":    mod.saltBeaconModify,
		// "beacon_delete":    mod.saltBeaconDelete,
		// "beacon_enable":    mod.saltBeaconEnable,
		// "beacon_disable":   mod.saltBeaconDisable,
		// "beacon_save":      mod.saltBeaconSave,
		// "beacon_reset":     mod.saltBeaconReset,
		
		// Schedule management (commented out until implemented)
		// "schedule_list":    mod.saltScheduleList,
		// "schedule_add":     mod.saltScheduleAdd,
		// "schedule_modify":  mod.saltScheduleModify,
		// "schedule_delete":  mod.saltScheduleDelete,
		// "schedule_enable":  mod.saltScheduleEnable,
		// "schedule_disable": mod.saltScheduleDisable,
		// "schedule_run_job": mod.saltScheduleRunJob,
		// "schedule_save":    mod.saltScheduleSave,
		// "schedule_reload":  mod.saltScheduleReload,
		
		// Advanced features (commented out until implemented)
		// "vault_read":       mod.saltVaultRead,
		// "vault_write":      mod.saltVaultWrite,
		// "vault_delete":     mod.saltVaultDelete,
		// "vault_list":       mod.saltVaultList,
		// "x509_create_certificate": mod.saltX509CreateCertificate,
		// "x509_read_certificate": mod.saltX509ReadCertificate,
		// "ssh_run":          mod.saltSSHRun,
		// "ssh_state":        mod.saltSSHState,
		// "ssh_ping":         mod.saltSSHPing,
		// "ssh_copy":         mod.saltSSHCopy,
		// "proxy_list":       mod.saltProxyList,
		// "proxy_ping":       mod.saltProxyPing,
		// "proxy_conn_check": mod.saltProxyConnCheck,
		// "proxy_alive":      mod.saltProxyAlive,
		// "reactor_list":     mod.saltReactorList,
		// "reactor_add":      mod.saltReactorAdd,
		// "reactor_delete":   mod.saltReactorDelete,
		// "reactor_clear":    mod.saltReactorClear,
		// "cache_grains":     mod.saltCacheGrains,
		// "cache_pillar":     mod.saltCachePillar,
		// "cache_mine":       mod.saltCacheMine,
		// "cache_store":      mod.saltCacheStore,
		// "cache_fetch":      mod.saltCacheFetch,
		// "cache_flush":      mod.saltCacheFlush,
		// "syndic_list":      mod.saltSyndicList,
		// "syndic_refresh":   mod.saltSyndicRefresh,
		// "multi_master_setup": mod.saltMultiMasterSetup,
		// "multi_master_failover": mod.saltMultiMasterFailover,
		// "multi_master_status": mod.saltMultiMasterStatus,
		// "performance_profile": mod.saltPerformanceProfile,
		// "performance_test": mod.saltPerformanceTest,
		// "performance_benchmark": mod.saltPerformanceBenchmark,
		// "cache_performance": mod.saltCachePerformance,
		// "roster_list":      mod.saltRosterList,
		// "roster_add":       mod.saltRosterAdd,
		// "roster_remove":    mod.saltRosterRemove,
		// "roster_update":    mod.saltRosterUpdate,
		// "fileserver_list_envs": mod.saltFileserverListEnvs,
		// "fileserver_file_list": mod.saltFileserverFileList,
		// "fileserver_dir_list": mod.saltFileserverDirList,
		// "fileserver_symlink_list": mod.saltFileserverSymlinkList,
		// "fileserver_update": mod.saltFileserverUpdate,
		
		// Utility methods (commented out until implemented)
		// "helper_match":     mod.saltHelperMatch,
		// "helper_glob":      mod.saltHelperGlob,
		// "helper_timeout":   mod.saltHelperTimeout,
		// "helper_retry":     mod.saltHelperRetry,
		// "helper_env":       mod.saltHelperEnv,
		// "helper_which":     mod.saltHelperWhich,
		// "helper_random":    mod.saltHelperRandom,
		// "helper_base64":    mod.saltHelperBase64,
	})
	
	// Add targeting methods
	L.SetFuncs(saltClient, map[string]lua.LGFunction{
		"target":           mod.saltTarget,
		"with_timeout":     mod.saltWithTimeout,
		"with_retries":     mod.saltWithRetries,
		"with_pillar":      mod.saltWithPillar,
		"with_grains":      mod.saltWithGrains,
	})
	
	L.Push(saltClient)
	return 1
}

// Helper function to execute salt commands with enhanced error handling and features
func (mod *ObjectOrientedSaltModule) executeSaltCommand(L *lua.LState, timeout int, retries int, cmdArgs ...string) (map[string]interface{}, error) {
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

func (mod *ObjectOrientedSaltModule) executeSingleSaltCommand(timeout int, cmdArgs ...string) (map[string]interface{}, error) {
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
func (mod *ObjectOrientedSaltModule) goValueToLua(L *lua.LState, value interface{}) lua.LValue {
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

// Convert Lua values to Go values
func (mod *ObjectOrientedSaltModule) luaValueToGo(value lua.LValue) interface{} {
	switch v := value.(type) {
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		result := make(map[string]interface{})
		v.ForEach(func(key, val lua.LValue) {
			result[key.String()] = mod.luaValueToGo(val)
		})
		return result
	default:
		return value.String()
	}
}

// Helper function to return consistent Salt results
func (mod *ObjectOrientedSaltModule) returnSaltResult(L *lua.LState, result map[string]interface{}, err error) int {
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

// ObjectOrientedSaltLoader loads the object-oriented salt module for Lua
func ObjectOrientedSaltLoader(L *lua.LState) int {
	mod := NewObjectOrientedSaltModule()
	return mod.Loader(L)
}