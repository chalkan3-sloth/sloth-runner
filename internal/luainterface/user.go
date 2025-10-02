package luainterface

import (
	"fmt"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// UserModule provides functions for Linux/Unix user and group management.
type UserModule struct{}

// NewUserModule creates a new UserModule.
func NewUserModule() *UserModule {
	return &UserModule{}
}

// Loader is the module loader function.
func (u *UserModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), u.exports())
	L.Push(mod)
	return 1
}

func (u *UserModule) exports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		// User management
		"create":           u.createUser,
		"delete":           u.deleteUser,
		"modify":           u.modifyUser,
		"exists":           u.userExists,
		"get_info":         u.getUserInfo,
		"list":             u.listUsers,
		"lock":             u.lockUser,
		"unlock":           u.unlockUser,
		"is_locked":        u.isUserLocked,
		"set_password":     u.setPassword,
		"expire_password":  u.expirePassword,
		"set_shell":        u.setShell,
		"set_home":         u.setHomeDir,
		"get_uid":          u.getUID,
		"get_gid":          u.getGID,
		"get_groups":       u.getUserGroups,
		"add_to_group":     u.addUserToGroup,
		"remove_from_group": u.removeUserFromGroup,
		"set_primary_group": u.setPrimaryGroup,
		"get_home":         u.getHomeDir,
		"get_shell":        u.getShell,
		"get_comment":      u.getComment,
		"set_comment":      u.setComment,
		"is_system_user":   u.isSystemUser,
		"get_current":      u.getCurrentUser,
		
		// Group management
		"group_create":     u.createGroup,
		"group_delete":     u.deleteGroup,
		"group_exists":     u.groupExists,
		"group_get_info":   u.getGroupInfo,
		"group_list":       u.listGroups,
		"group_get_gid":    u.getGroupGID,
		"group_members":    u.getGroupMembers,
		"group_add_member": u.addGroupMember,
		"group_remove_member": u.removeGroupMember,
		
		// Advanced features
		"set_expiry":       u.setAccountExpiry,
		"get_last_login":   u.getLastLogin,
		"get_failed_logins": u.getFailedLogins,
		"validate_username": u.validateUsername,
		"is_root":          u.isRoot,
		"run_as":           u.runAs,
	}
}

// needsSudo checks if we need sudo for user management commands
func (u *UserModule) needsSudo() bool {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return false
	}
	
	currentUser, err := user.Current()
	if err != nil {
		return true
	}
	
	return currentUser.Uid != "0"
}

// parseUserOptions parses optional table argument
func (u *UserModule) parseUserOptions(L *lua.LState, idx int) map[string]string {
	options := make(map[string]string)
	
	val := L.Get(idx)
	if val.Type() != lua.LTTable {
		return options
	}
	
	tbl := val.(*lua.LTable)
	tbl.ForEach(func(k, v lua.LValue) {
		if k.Type() == lua.LTString && v.Type() == lua.LTString {
			options[k.String()] = v.String()
		}
	})
	
	return options
}

// createUser creates a new user
func (u *UserModule) createUser(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	options := u.parseUserOptions(L, 2)
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "useradd")
	
	// Apply options
	if home, ok := options["home"]; ok {
		args = append(args, "-d", home)
	}
	if shell, ok := options["shell"]; ok {
		args = append(args, "-s", shell)
	}
	if uid, ok := options["uid"]; ok {
		args = append(args, "-u", uid)
	}
	if gid, ok := options["gid"]; ok {
		args = append(args, "-g", gid)
	}
	if groups, ok := options["groups"]; ok {
		args = append(args, "-G", groups)
	}
	if comment, ok := options["comment"]; ok {
		args = append(args, "-c", comment)
	}
	if _, ok := options["system"]; ok {
		args = append(args, "-r")
	}
	if _, ok := options["create_home"]; ok {
		args = append(args, "-m")
	} else if _, ok := options["no_create_home"]; ok {
		args = append(args, "-M")
	}
	if expiry, ok := options["expiry"]; ok {
		args = append(args, "-e", expiry)
	}
	
	args = append(args, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create user: %s\n%s", err, string(output))))
		return 2
	}
	
	// Set password if provided
	if password, ok := options["password"]; ok && password != "" {
		var passArgs []string
		if u.needsSudo() {
			passArgs = append(passArgs, "sudo")
		}
		
		passArgs = append(passArgs, "chpasswd")
		
		passCmd := exec.Command(passArgs[0], passArgs[1:]...)
		passCmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s\n", username, password))
		
		passOutput, passErr := passCmd.CombinedOutput()
		if passErr != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(fmt.Sprintf("User created but failed to set password: %s\n%s", passErr, string(passOutput))))
			return 2
		}
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// deleteUser deletes a user
func (u *UserModule) deleteUser(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	removeHome := L.ToBool(2)
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "userdel")
	if removeHome {
		args = append(args, "-r")
	}
	args = append(args, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to delete user: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// modifyUser modifies an existing user
func (u *UserModule) modifyUser(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	options := u.parseUserOptions(L, 2)
	if len(options) == 0 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No modifications specified"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod")
	
	if home, ok := options["home"]; ok {
		args = append(args, "-d", home)
		if _, moveOk := options["move_home"]; moveOk {
			args = append(args, "-m")
		}
	}
	if shell, ok := options["shell"]; ok {
		args = append(args, "-s", shell)
	}
	if uid, ok := options["uid"]; ok {
		args = append(args, "-u", uid)
	}
	if gid, ok := options["gid"]; ok {
		args = append(args, "-g", gid)
	}
	if groups, ok := options["groups"]; ok {
		args = append(args, "-G", groups)
	}
	if comment, ok := options["comment"]; ok {
		args = append(args, "-c", comment)
	}
	if expiry, ok := options["expiry"]; ok {
		args = append(args, "-e", expiry)
	}
	if _, ok := options["lock"]; ok {
		args = append(args, "-L")
	}
	if _, ok := options["unlock"]; ok {
		args = append(args, "-U")
	}
	
	args = append(args, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to modify user: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// userExists checks if a user exists
func (u *UserModule) userExists(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	_, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("User does not exist: %s", username)))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString("User exists"))
	return 2
}

// getUserInfo gets detailed user information
func (u *UserModule) getUserInfo(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	usr, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get user info: %s", err)))
		return 2
	}
	
	info := L.NewTable()
	info.RawSetString("username", lua.LString(usr.Username))
	info.RawSetString("uid", lua.LString(usr.Uid))
	info.RawSetString("gid", lua.LString(usr.Gid))
	info.RawSetString("name", lua.LString(usr.Name))
	info.RawSetString("home", lua.LString(usr.HomeDir))
	
	// Get additional info from getent
	cmd := exec.Command("getent", "passwd", username)
	output, err := cmd.CombinedOutput()
	if err == nil {
		fields := strings.Split(strings.TrimSpace(string(output)), ":")
		if len(fields) >= 7 {
			info.RawSetString("shell", lua.LString(fields[6]))
			info.RawSetString("comment", lua.LString(fields[4]))
		}
	}
	
	L.Push(info)
	L.Push(lua.LNil)
	return 2
}

// listUsers lists all users
func (u *UserModule) listUsers(L *lua.LState) int {
	systemOnly := L.ToBool(1)
	
	cmd := exec.Command("getent", "passwd")
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list users: %s", err)))
		return 2
	}
	
	users := L.NewTable()
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}
		
		username := fields[0]
		uid, _ := strconv.Atoi(fields[2])
		
		// Filter system users if requested
		if systemOnly && uid >= 1000 {
			continue
		}
		if !systemOnly && uid < 1000 {
			continue
		}
		
		userInfo := L.NewTable()
		userInfo.RawSetString("username", lua.LString(username))
		userInfo.RawSetString("uid", lua.LString(fields[2]))
		userInfo.RawSetString("gid", lua.LString(fields[3]))
		userInfo.RawSetString("comment", lua.LString(fields[4]))
		userInfo.RawSetString("home", lua.LString(fields[5]))
		userInfo.RawSetString("shell", lua.LString(fields[6]))
		
		users.Append(userInfo)
	}
	
	L.Push(users)
	L.Push(lua.LNil)
	return 2
}

// lockUser locks a user account
func (u *UserModule) lockUser(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-L", username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to lock user: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// unlockUser unlocks a user account
func (u *UserModule) unlockUser(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-U", username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to unlock user: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// isUserLocked checks if a user account is locked
func (u *UserModule) isUserLocked(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "passwd", "-S", username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to check lock status: %s", err)))
		return 2
	}
	
	// Check if output contains 'L' (locked) status
	outputStr := string(output)
	if strings.Contains(outputStr, " L ") || strings.Contains(outputStr, " LK ") {
		L.Push(lua.LTrue)
		L.Push(lua.LString("User is locked"))
		return 2
	}
	
	L.Push(lua.LFalse)
	L.Push(lua.LString("User is not locked"))
	return 2
}

// setPassword sets a user's password
func (u *UserModule) setPassword(L *lua.LState) int {
	username := L.ToString(1)
	password := L.ToString(2)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if password == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Password is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "chpasswd")
	
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s\n", username, password))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to set password: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString("Password set successfully"))
	return 2
}

// expirePassword expires a user's password
func (u *UserModule) expirePassword(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "passwd", "-e", username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to expire password: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// setShell sets the user's shell
func (u *UserModule) setShell(L *lua.LState) int {
	username := L.ToString(1)
	shell := L.ToString(2)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if shell == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Shell is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-s", shell, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to set shell: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// setHomeDir sets the user's home directory
func (u *UserModule) setHomeDir(L *lua.LState) int {
	username := L.ToString(1)
	homeDir := L.ToString(2)
	moveFiles := L.ToBool(3)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if homeDir == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Home directory is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-d", homeDir)
	if moveFiles {
		args = append(args, "-m")
	}
	args = append(args, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to set home directory: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// getUID gets the UID of a user
func (u *UserModule) getUID(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	usr, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get UID: %s", err)))
		return 2
	}
	
	uid, _ := strconv.Atoi(usr.Uid)
	L.Push(lua.LNumber(uid))
	L.Push(lua.LNil)
	return 2
}

// getGID gets the primary GID of a user
func (u *UserModule) getGID(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	usr, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get GID: %s", err)))
		return 2
	}
	
	gid, _ := strconv.Atoi(usr.Gid)
	L.Push(lua.LNumber(gid))
	L.Push(lua.LNil)
	return 2
}

// getUserGroups gets all groups a user belongs to
func (u *UserModule) getUserGroups(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	cmd := exec.Command("id", "-Gn", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get user groups: %s", err)))
		return 2
	}
	
	groups := L.NewTable()
	groupList := strings.Fields(strings.TrimSpace(string(output)))
	for _, group := range groupList {
		groups.Append(lua.LString(group))
	}
	
	L.Push(groups)
	L.Push(lua.LNil)
	return 2
}

// addUserToGroup adds a user to a group
func (u *UserModule) addUserToGroup(L *lua.LState) int {
	username := L.ToString(1)
	group := L.ToString(2)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if group == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-aG", group, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to add user to group: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// removeUserFromGroup removes a user from a group
func (u *UserModule) removeUserFromGroup(L *lua.LState) int {
	username := L.ToString(1)
	group := L.ToString(2)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if group == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "gpasswd", "-d", username, group)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove user from group: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// setPrimaryGroup sets the user's primary group
func (u *UserModule) setPrimaryGroup(L *lua.LState) int {
	username := L.ToString(1)
	group := L.ToString(2)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if group == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-g", group, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to set primary group: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// getHomeDir gets the user's home directory
func (u *UserModule) getHomeDir(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	usr, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get home directory: %s", err)))
		return 2
	}
	
	L.Push(lua.LString(usr.HomeDir))
	L.Push(lua.LNil)
	return 2
}

// getShell gets the user's shell
func (u *UserModule) getShell(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	cmd := exec.Command("getent", "passwd", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get shell: %s", err)))
		return 2
	}
	
	fields := strings.Split(strings.TrimSpace(string(output)), ":")
	if len(fields) < 7 {
		L.Push(lua.LNil)
		L.Push(lua.LString("Invalid passwd entry"))
		return 2
	}
	
	L.Push(lua.LString(fields[6]))
	L.Push(lua.LNil)
	return 2
}

// getComment gets the user's comment/GECOS field
func (u *UserModule) getComment(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	usr, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get comment: %s", err)))
		return 2
	}
	
	L.Push(lua.LString(usr.Name))
	L.Push(lua.LNil)
	return 2
}

// setComment sets the user's comment/GECOS field
func (u *UserModule) setComment(L *lua.LState) int {
	username := L.ToString(1)
	comment := L.ToString(2)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-c", comment, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to set comment: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// isSystemUser checks if a user is a system user (UID < 1000)
func (u *UserModule) isSystemUser(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	usr, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to check user: %s", err)))
		return 2
	}
	
	uid, _ := strconv.Atoi(usr.Uid)
	L.Push(lua.LBool(uid < 1000))
	L.Push(lua.LNil)
	return 2
}

// getCurrentUser gets the current user
func (u *UserModule) getCurrentUser(L *lua.LState) int {
	usr, err := user.Current()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get current user: %s", err)))
		return 2
	}
	
	info := L.NewTable()
	info.RawSetString("username", lua.LString(usr.Username))
	info.RawSetString("uid", lua.LString(usr.Uid))
	info.RawSetString("gid", lua.LString(usr.Gid))
	info.RawSetString("name", lua.LString(usr.Name))
	info.RawSetString("home", lua.LString(usr.HomeDir))
	
	L.Push(info)
	L.Push(lua.LNil)
	return 2
}

// GROUP MANAGEMENT FUNCTIONS

// createGroup creates a new group
func (u *UserModule) createGroup(L *lua.LState) int {
	groupname := L.ToString(1)
	if groupname == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	options := u.parseUserOptions(L, 2)
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "groupadd")
	
	if gid, ok := options["gid"]; ok {
		args = append(args, "-g", gid)
	}
	if _, ok := options["system"]; ok {
		args = append(args, "-r")
	}
	
	args = append(args, groupname)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create group: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// deleteGroup deletes a group
func (u *UserModule) deleteGroup(L *lua.LState) int {
	groupname := L.ToString(1)
	if groupname == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "groupdel", groupname)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to delete group: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// groupExists checks if a group exists
func (u *UserModule) groupExists(L *lua.LState) int {
	groupname := L.ToString(1)
	if groupname == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	_, err := user.LookupGroup(groupname)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Group does not exist: %s", groupname)))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString("Group exists"))
	return 2
}

// getGroupInfo gets group information
func (u *UserModule) getGroupInfo(L *lua.LState) int {
	groupname := L.ToString(1)
	if groupname == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	grp, err := user.LookupGroup(groupname)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get group info: %s", err)))
		return 2
	}
	
	info := L.NewTable()
	info.RawSetString("name", lua.LString(grp.Name))
	info.RawSetString("gid", lua.LString(grp.Gid))
	
	L.Push(info)
	L.Push(lua.LNil)
	return 2
}

// listGroups lists all groups
func (u *UserModule) listGroups(L *lua.LState) int {
	cmd := exec.Command("getent", "group")
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to list groups: %s", err)))
		return 2
	}
	
	groups := L.NewTable()
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}
		
		groupInfo := L.NewTable()
		groupInfo.RawSetString("name", lua.LString(fields[0]))
		groupInfo.RawSetString("gid", lua.LString(fields[2]))
		
		if fields[3] != "" {
			members := L.NewTable()
			for _, member := range strings.Split(fields[3], ",") {
				members.Append(lua.LString(member))
			}
			groupInfo.RawSetString("members", members)
		}
		
		groups.Append(groupInfo)
	}
	
	L.Push(groups)
	L.Push(lua.LNil)
	return 2
}

// getGroupGID gets the GID of a group
func (u *UserModule) getGroupGID(L *lua.LState) int {
	groupname := L.ToString(1)
	if groupname == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	grp, err := user.LookupGroup(groupname)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get GID: %s", err)))
		return 2
	}
	
	gid, _ := strconv.Atoi(grp.Gid)
	L.Push(lua.LNumber(gid))
	L.Push(lua.LNil)
	return 2
}

// getGroupMembers gets all members of a group
func (u *UserModule) getGroupMembers(L *lua.LState) int {
	groupname := L.ToString(1)
	if groupname == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	cmd := exec.Command("getent", "group", groupname)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get group members: %s", err)))
		return 2
	}
	
	fields := strings.Split(strings.TrimSpace(string(output)), ":")
	if len(fields) < 4 {
		L.Push(L.NewTable())
		L.Push(lua.LNil)
		return 2
	}
	
	members := L.NewTable()
	if fields[3] != "" {
		for _, member := range strings.Split(fields[3], ",") {
			members.Append(lua.LString(strings.TrimSpace(member)))
		}
	}
	
	L.Push(members)
	L.Push(lua.LNil)
	return 2
}

// addGroupMember adds a member to a group
func (u *UserModule) addGroupMember(L *lua.LState) int {
	groupname := L.ToString(1)
	username := L.ToString(2)
	
	if groupname == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-aG", groupname, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to add member to group: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// removeGroupMember removes a member from a group
func (u *UserModule) removeGroupMember(L *lua.LState) int {
	groupname := L.ToString(1)
	username := L.ToString(2)
	
	if groupname == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Group name is required"))
		return 2
	}
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "gpasswd", "-d", username, groupname)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove member from group: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// ADVANCED FEATURES

// setAccountExpiry sets when an account expires
func (u *UserModule) setAccountExpiry(L *lua.LState) int {
	username := L.ToString(1)
	expiry := L.ToString(2) // Format: YYYY-MM-DD
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if expiry == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expiry date is required (format: YYYY-MM-DD)"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo")
	}
	
	args = append(args, "usermod", "-e", expiry, username)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to set account expiry: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// getLastLogin gets the last login time for a user
func (u *UserModule) getLastLogin(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	cmd := exec.Command("lastlog", "-u", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get last login: %s", err)))
		return 2
	}
	
	L.Push(lua.LString(string(output)))
	L.Push(lua.LNil)
	return 2
}

// getFailedLogins gets failed login attempts for a user
func (u *UserModule) getFailedLogins(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	cmd := exec.Command("faillog", "-u", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// faillog might not exist or require specific permissions
		L.Push(lua.LString(""))
		L.Push(lua.LString(fmt.Sprintf("Failed to get failed logins: %s", err)))
		return 2
	}
	
	L.Push(lua.LString(string(output)))
	L.Push(lua.LNil)
	return 2
}

// validateUsername validates if a username follows Linux conventions
func (u *UserModule) validateUsername(L *lua.LState) int {
	username := L.ToString(1)
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username cannot be empty"))
		return 2
	}
	
	// Basic validation rules
	if len(username) > 32 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username too long (max 32 characters)"))
		return 2
	}
	
	// First character must be a letter or underscore
	if !((username[0] >= 'a' && username[0] <= 'z') ||
		(username[0] >= 'A' && username[0] <= 'Z') ||
		username[0] == '_') {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username must start with a letter or underscore"))
		return 2
	}
	
	// Other characters must be alphanumeric, hyphen, or underscore
	for i := 1; i < len(username); i++ {
		c := username[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_') {
			L.Push(lua.LFalse)
			L.Push(lua.LString("Username contains invalid characters"))
			return 2
		}
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString("Valid username"))
	return 2
}

// isRoot checks if the current user is root
func (u *UserModule) isRoot(L *lua.LState) int {
	usr, err := user.Current()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to get current user: %s", err)))
		return 2
	}
	
	L.Push(lua.LBool(usr.Uid == "0"))
	L.Push(lua.LNil)
	return 2
}

// runAs runs a command as a different user
func (u *UserModule) runAs(L *lua.LState) int {
	username := L.ToString(1)
	command := L.ToString(2)
	
	if username == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Username is required"))
		return 2
	}
	
	if command == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Command is required"))
		return 2
	}
	
	var args []string
	if u.needsSudo() {
		args = append(args, "sudo", "-u", username)
	} else {
		args = append(args, "su", "-", username, "-c")
	}
	
	args = append(args, "sh", "-c", command)
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to run command as user: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}
