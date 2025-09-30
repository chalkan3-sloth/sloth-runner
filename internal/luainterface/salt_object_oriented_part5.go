package luainterface

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

// All remaining Salt functions with object-oriented implementation

// Mine operations
func (mod *ObjectOrientedSaltModule) saltMineGet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	function := L.CheckString(2)
	
	args := []string{"salt", target, "mine.get", "*", function, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMineSend(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	function := L.CheckString(2)
	
	args := []string{"salt", target, "mine.send", function, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMineUpdate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "mine.update", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMineDelete(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	function := L.CheckString(2)
	
	args := []string{"salt", target, "mine.delete", function, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMineFlush(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "mine.flush", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMineValid(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "mine.valid", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Job management
func (mod *ObjectOrientedSaltModule) saltJobActive(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt-run", "jobs.active", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltJobList(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt-run", "jobs.list_jobs", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltJobLookup(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	jid := L.CheckString(2)
	
	args := []string{"salt-run", "jobs.lookup_jid", jid, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltJobExitSuccess(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	jid := L.CheckString(2)
	
	args := []string{"salt-run", "jobs.exit_success", jid, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltJobPrint(L *lua.LState) int {
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	jid := L.CheckString(2)
	
	args := []string{"salt-run", "jobs.print_job", jid, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Docker integration
func (mod *ObjectOrientedSaltModule) saltDockerPs(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "dockerng.ps", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerRun(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	image := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.run", image, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerStop(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	container := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.stop", container, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerStart(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	container := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.start", container, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerRestart(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	container := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.restart", container, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerBuild(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	tag := L.CheckString(3)
	
	args := []string{"salt", target, "dockerng.build", path, "tag=" + tag, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*5, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerPull(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	image := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.pull", image, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*5, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerPush(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	image := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.push", image, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*5, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerImages(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "dockerng.images", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerRemove(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	container := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.rm", container, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerInspect(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	container := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.inspect", container, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerLogs(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	container := L.CheckString(2)
	
	args := []string{"salt", target, "dockerng.logs", container, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDockerExec(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	container := L.CheckString(2)
	cmd := L.CheckString(3)
	
	args := []string{"salt", target, "dockerng.exec", container, cmd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Git operations
func (mod *ObjectOrientedSaltModule) saltGitClone(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	repo := L.CheckString(2)
	dir := L.CheckString(3)
	
	args := []string{"salt", target, "git.clone", repo, dir, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*5, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitPull(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	
	args := []string{"salt", target, "git.pull", cwd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitCheckout(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	rev := L.CheckString(3)
	
	args := []string{"salt", target, "git.checkout", cwd, rev, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitAdd(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	filename := L.CheckString(3)
	
	args := []string{"salt", target, "git.add", cwd, filename, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitCommit(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	message := L.CheckString(3)
	
	args := []string{"salt", target, "git.commit", cwd, message, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitPush(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	
	args := []string{"salt", target, "git.push", cwd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout*2, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitStatus(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	
	args := []string{"salt", target, "git.status", cwd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitLog(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	
	args := []string{"salt", target, "git.log", cwd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitReset(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	opts := L.OptString(3, "")
	
	args := []string{"salt", target, "git.reset", cwd}
	if opts != "" {
		args = append(args, opts)
	}
	args = append(args, "--out=json")
	
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitRemoteGet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	
	args := []string{"salt", target, "git.remote_get", cwd, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltGitRemoteSet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	cwd := L.CheckString(2)
	name := L.CheckString(3)
	url := L.CheckString(4)
	
	args := []string{"salt", target, "git.remote_set", cwd, name, url, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Database operations - MySQL
func (mod *ObjectOrientedSaltModule) saltMysqlQuery(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	query := L.CheckString(2)
	
	args := []string{"salt", target, "mysql.query", "''", query, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMysqlDbCreate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	database := L.CheckString(2)
	
	args := []string{"salt", target, "mysql.db_create", database, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMysqlDbRemove(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	database := L.CheckString(2)
	
	args := []string{"salt", target, "mysql.db_remove", database, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMysqlUserCreate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.CheckString(2)
	host := L.CheckString(3)
	password := L.CheckString(4)
	
	args := []string{"salt", target, "mysql.user_create", user, host, password, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMysqlUserRemove(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.CheckString(2)
	host := L.CheckString(3)
	
	args := []string{"salt", target, "mysql.user_remove", user, host, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMysqlGrantAdd(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	grant := L.CheckString(2)
	database := L.CheckString(3)
	user := L.CheckString(4)
	host := L.CheckString(5)
	
	args := []string{"salt", target, "mysql.grant_add", grant, database, user, host, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltMysqlGrantRevoke(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	grant := L.CheckString(2)
	database := L.CheckString(3)
	user := L.CheckString(4)
	host := L.CheckString(5)
	
	args := []string{"salt", target, "mysql.grant_revoke", grant, database, user, host, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Database operations - PostgreSQL
func (mod *ObjectOrientedSaltModule) saltPostgresQuery(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	query := L.CheckString(2)
	
	args := []string{"salt", target, "postgres.psql_query", query, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPostgresDbCreate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	database := L.CheckString(2)
	
	args := []string{"salt", target, "postgres.db_create", database, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPostgresDbRemove(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	database := L.CheckString(2)
	
	args := []string{"salt", target, "postgres.db_remove", database, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPostgresUserCreate(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.CheckString(2)
	password := L.CheckString(3)
	
	args := []string{"salt", target, "postgres.user_create", user, password, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltPostgresUserRemove(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	user := L.CheckString(2)
	
	args := []string{"salt", target, "postgres.user_remove", user, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Monitoring and metrics
func (mod *ObjectOrientedSaltModule) saltStatusLoadavg(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.loadavg", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStatusCpuinfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.cpuinfo", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStatusMeminfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.meminfo", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStatusDiskusage(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.OptString(2, "/")
	
	args := []string{"salt", target, "status.diskusage", path, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStatusNetdev(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.netdev", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStatusW(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.w", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltStatusUptime(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "status.uptime", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Configuration management
func (mod *ObjectOrientedSaltModule) saltConfigGet(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	key := L.CheckString(2)
	
	args := []string{"salt", target, "config.get", key, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltConfigOption(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	option := L.CheckString(2)
	
	args := []string{"salt", target, "config.option", option, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltConfigValidFileproto(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	path := L.CheckString(2)
	
	args := []string{"salt", target, "config.valid_fileproto", path, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltConfigBackupMode(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "config.backup_mode", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Template engines - CLI based implementation
func (mod *ObjectOrientedSaltModule) saltTemplateJinja(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	template := L.CheckString(2)
	
	args := []string{"salt", target, "slsutil.renderer", template, "jinja", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltTemplateYaml(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	template := L.CheckString(2)
	
	args := []string{"salt", target, "slsutil.renderer", template, "yaml", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltTemplateJson(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	template := L.CheckString(2)
	
	args := []string{"salt", target, "slsutil.renderer", template, "json", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltTemplateMako(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	template := L.CheckString(2)
	
	args := []string{"salt", target, "slsutil.renderer", template, "mako", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltTemplatePy(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	template := L.CheckString(2)
	
	args := []string{"salt", target, "slsutil.renderer", template, "py", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltTemplateWempy(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	template := L.CheckString(2)
	
	args := []string{"salt", target, "slsutil.renderer", template, "wempy", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

// Logging and debugging
func (mod *ObjectOrientedSaltModule) saltLogError(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	message := L.CheckString(2)
	
	args := []string{"salt", target, "log.error", message, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltLogWarning(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	message := L.CheckString(2)
	
	args := []string{"salt", target, "log.warning", message, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltLogInfo(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	message := L.CheckString(2)
	
	args := []string{"salt", target, "log.info", message, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltLogDebug(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	message := L.CheckString(2)
	
	args := []string{"salt", target, "log.debug", message, "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}

func (mod *ObjectOrientedSaltModule) saltDebugMode(L *lua.LState) int {
	enabled := L.CheckBool(2)
	
	result := map[string]interface{}{
		"success":    true,
		"debug_mode": enabled,
		"message":    fmt.Sprintf("Debug mode %s", map[bool]string{true: "enabled", false: "disabled"}[enabled]),
	}
	return mod.returnSaltResult(L, result, nil)
}

func (mod *ObjectOrientedSaltModule) saltDebugProfile(L *lua.LState) int {
	target := mod.getTarget(L)
	timeout := mod.getTimeout(L)
	retries := mod.getRetries(L)
	
	args := []string{"salt", target, "sys.doc", "--out=json"}
	result, err := mod.executeSaltCommand(L, timeout, retries, args...)
	return mod.returnSaltResult(L, result, err)
}