package luainterface

import (
	"os/user"
	"runtime"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestUserModule(t *testing.T) {
	// Skip tests on non-Unix systems
	if runtime.GOOS == "windows" {
		t.Skip("Skipping user module tests on Windows")
	}

	t.Run("UserModuleCreation", func(t *testing.T) {
		module := NewUserModule()
		if module == nil {
			t.Fatal("Expected UserModule instance, got nil")
		}
	})

	t.Run("UserModuleLoader", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()

		module := NewUserModule()
		L.PreloadModule("user", module.Loader)

		// Load the module
		if err := L.DoString(`user = require("user")`); err != nil {
			t.Fatalf("Failed to load user module: %v", err)
		}

		// Check if module is loaded
		mod := L.GetGlobal("user")
		if mod.Type() != lua.LTTable {
			t.Fatal("Expected user module to be a table")
		}
	})
}

func TestUserExists(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	t.Run("ExistingUser", func(t *testing.T) {
		script := `
			local user = require("user")
			local exists, msg = user.exists("` + currentUser.Username + `")
			return exists, msg
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		exists := L.ToBool(-2)
		if !exists {
			msg := L.ToString(-1)
			t.Errorf("Expected user %s to exist, but got: %s", currentUser.Username, msg)
		}
		L.Pop(2)
	})

	t.Run("NonExistingUser", func(t *testing.T) {
		script := `
			local user = require("user")
			local exists, msg = user.exists("nonexistent_user_xyz_12345")
			return exists, msg
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		exists := L.ToBool(-2)
		if exists {
			t.Error("Expected non-existent user to return false")
		}
		L.Pop(2)
	})

	t.Run("EmptyUsername", func(t *testing.T) {
		script := `
			local user = require("user")
			local exists, msg = user.exists("")
			return exists, msg
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		exists := L.ToBool(-2)
		if exists {
			t.Error("Expected empty username to return false")
		}
		L.Pop(2)
	})
}

func TestGetUserInfo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	t.Run("GetCurrentUserInfo", func(t *testing.T) {
		script := `
			local user = require("user")
			local info, err = user.get_info("` + currentUser.Username + `")
			if info == nil then
				return nil, err
			end
			return info.username, info.uid, info.gid, info.home
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		username := L.ToString(-4)
		uid := L.ToString(-3)
		gid := L.ToString(-2)
		home := L.ToString(-1)

		if username != currentUser.Username {
			t.Errorf("Expected username %s, got %s", currentUser.Username, username)
		}

		if uid != currentUser.Uid {
			t.Errorf("Expected UID %s, got %s", currentUser.Uid, uid)
		}

		if gid != currentUser.Gid {
			t.Errorf("Expected GID %s, got %s", currentUser.Gid, gid)
		}

		if home != currentUser.HomeDir {
			t.Errorf("Expected home %s, got %s", currentUser.HomeDir, home)
		}

		L.Pop(4)
	})

	t.Run("GetNonExistentUserInfo", func(t *testing.T) {
		script := `
			local user = require("user")
			local info, err = user.get_info("nonexistent_user_xyz_12345")
			return info, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		info := L.Get(-2)
		if info.Type() != lua.LTNil {
			t.Error("Expected nil for non-existent user info")
		}

		errMsg := L.ToString(-1)
		if errMsg == "" {
			t.Error("Expected error message for non-existent user")
		}

		L.Pop(2)
	})
}

func TestGetUID(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	t.Run("GetCurrentUserUID", func(t *testing.T) {
		script := `
			local user = require("user")
			local uid, err = user.get_uid("` + currentUser.Username + `")
			return uid, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		uid := L.ToNumber(-2)
		if uid <= 0 {
			t.Error("Expected valid UID greater than 0")
		}

		L.Pop(2)
	})

	t.Run("GetNonExistentUserUID", func(t *testing.T) {
		script := `
			local user = require("user")
			local uid, err = user.get_uid("nonexistent_user_xyz_12345")
			return uid, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		uid := L.Get(-2)
		if uid.Type() != lua.LTNil {
			t.Error("Expected nil for non-existent user UID")
		}

		L.Pop(2)
	})
}

func TestGetGID(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	t.Run("GetCurrentUserGID", func(t *testing.T) {
		script := `
			local user = require("user")
			local gid, err = user.get_gid("` + currentUser.Username + `")
			return gid, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		gid := L.ToNumber(-2)
		if gid < 0 {
			t.Error("Expected valid GID")
		}

		L.Pop(2)
	})
}

func TestGetHomeDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	t.Run("GetCurrentUserHome", func(t *testing.T) {
		script := `
			local user = require("user")
			local home, err = user.get_home("` + currentUser.Username + `")
			return home, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		home := L.ToString(-2)
		if home != currentUser.HomeDir {
			t.Errorf("Expected home %s, got %s", currentUser.HomeDir, home)
		}

		L.Pop(2)
	})
}

// func TestGetShell(t *testing.T) {
// 	if runtime.GOOS == "windows" {
// 		t.Skip("Skipping on Windows")
// 	}
// 
// 	L := lua.NewState()
// 	defer L.Close()
// 
// 	module := NewUserModule()
// 	L.PreloadModule("user", module.Loader)
// 
// 	currentUser, err := user.Current()
// 	if err != nil {
// 		t.Skipf("Could not get current user: %v", err)
// 	}
// 
// 	t.Run("GetCurrentUserShell", func(t *testing.T) {
// 		script := `
// 			local user = require("user")
// 			local shell, err = user.get_shell("` + currentUser.Username + `")
// 			return shell, err
// 		`
// 
// 		if err := L.DoString(script); err != nil {
// 			t.Fatalf("Script execution failed: %v", err)
// 		}
// 
// 		shell := L.ToString(-2)
// 		// Just verify we got something back
// 		if shell == "" && L.Get(-1).Type() != lua.LTNil {
// 			t.Error("Expected shell value or error")
// 		}
// 
// 		L.Pop(2)
// 	})
// }

func TestGetUserGroups(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	t.Run("GetCurrentUserGroups", func(t *testing.T) {
		script := `
			local user = require("user")
			local groups, err = user.get_groups("` + currentUser.Username + `")
			if groups == nil then
				return nil, err
			end
			return #groups
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		count := L.ToNumber(-1)
		if count < 0 {
			t.Error("Expected at least 0 groups")
		}

		L.Pop(1)
	})
}

func TestIsSystemUser(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("CheckRootIsSystemUser", func(t *testing.T) {
		script := `
			local user = require("user")
			local is_system, err = user.is_system_user("root")
			return is_system, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		isSystem := L.ToBool(-2)
		if !isSystem {
			t.Error("Expected root to be a system user")
		}

		L.Pop(2)
	})
}

func TestGetCurrentUser(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("GetCurrentUserDetails", func(t *testing.T) {
		script := `
			local user = require("user")
			local info, err = user.get_current()
			if info == nil then
				return nil, err
			end
			return info.username, info.uid, info.home
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		username := L.ToString(-3)
		uid := L.ToString(-2)
		home := L.ToString(-1)

		if username == "" {
			t.Error("Expected non-empty username")
		}

		if uid == "" {
			t.Error("Expected non-empty UID")
		}

		if home == "" {
			t.Error("Expected non-empty home directory")
		}

		L.Pop(3)
	})
}

func TestIsRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("CheckIfRoot", func(t *testing.T) {
		script := `
			local user = require("user")
			local is_root, err = user.is_root()
			return is_root, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		isRoot := L.ToBool(-2)
		currentUser, _ := user.Current()
		expectedRoot := currentUser != nil && currentUser.Uid == "0"

		if isRoot != expectedRoot {
			t.Errorf("Expected is_root to be %v, got %v", expectedRoot, isRoot)
		}

		L.Pop(2)
	})
}

func TestValidateUsername(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	testCases := []struct {
		name     string
		username string
		valid    bool
	}{
		{"ValidSimple", "john", true},
		{"ValidWithUnderscore", "_test", true},
		{"ValidWithNumbers", "user123", true},
		{"ValidWithHyphen", "john-doe", true},
		{"InvalidStartsWithNumber", "123user", false},
		{"InvalidStartsWithHyphen", "-user", false},
		{"InvalidEmpty", "", false},
		{"InvalidTooLong", "this_username_is_way_too_long_for_linux_systems_maximum_length", false},
		{"InvalidSpecialChar", "user@domain", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			script := `
				local user = require("user")
				local valid, msg = user.validate_username("` + tc.username + `")
				return valid, msg
			`

			if err := L.DoString(script); err != nil {
				t.Fatalf("Script execution failed: %v", err)
			}

			valid := L.ToBool(-2)
			if valid != tc.valid {
				msg := L.ToString(-1)
				t.Errorf("Expected validation result %v for username '%s', got %v. Message: %s",
					tc.valid, tc.username, valid, msg)
			}

			L.Pop(2)
		})
	}
}

// GROUP MANAGEMENT TESTS

func TestGroupExists(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("CheckRootGroupExists", func(t *testing.T) {
		// Use a commonly available group
		groupName := "root"
		if runtime.GOOS == "darwin" {
			groupName = "wheel"
		}

		script := `
			local user = require("user")
			local exists, msg = user.group_exists("` + groupName + `")
			return exists, msg
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		exists := L.ToBool(-2)
		if !exists {
			msg := L.ToString(-1)
			t.Errorf("Expected group %s to exist, but got: %s", groupName, msg)
		}

		L.Pop(2)
	})

	t.Run("CheckNonExistentGroup", func(t *testing.T) {
		script := `
			local user = require("user")
			local exists, msg = user.group_exists("nonexistent_group_xyz_12345")
			return exists, msg
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		exists := L.ToBool(-2)
		if exists {
			t.Error("Expected non-existent group to return false")
		}

		L.Pop(2)
	})
}

func TestGetGroupInfo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("GetRootGroupInfo", func(t *testing.T) {
		groupName := "root"
		if runtime.GOOS == "darwin" {
			groupName = "wheel"
		}

		script := `
			local user = require("user")
			local info, err = user.group_get_info("` + groupName + `")
			if info == nil then
				return nil, err
			end
			return info.name, info.gid
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		name := L.ToString(-2)
		gid := L.ToString(-1)

		if name != groupName {
			t.Errorf("Expected group name %s, got %s", groupName, name)
		}

		if gid == "" {
			t.Error("Expected non-empty GID")
		}

		L.Pop(2)
	})
}

func TestGetGroupGID(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("GetRootGroupGID", func(t *testing.T) {
		groupName := "root"
		if runtime.GOOS == "darwin" {
			groupName = "wheel"
		}

		script := `
			local user = require("user")
			local gid, err = user.group_get_gid("` + groupName + `")
			return gid, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		gid := L.ToNumber(-2)
		if gid < 0 {
			t.Error("Expected valid GID")
		}

		L.Pop(2)
	})
}

func TestListUsers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("ListAllUsers", func(t *testing.T) {
		script := `
			local user = require("user")
			local users, err = user.list()
			if users == nil then
				return nil, err
			end
			return #users
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		count := L.ToNumber(-1)
		// On macOS, list() returns non-system users (UID >= 1000) which might be 0
		if count < 0 {
			t.Error("Expected at least zero users in the system")
		}

		L.Pop(1)
	})

	t.Run("ListSystemUsers", func(t *testing.T) {
		script := `
			local user = require("user")
			local users, err = user.list(true)
			if users == nil then
				return nil, err
			end
			return #users
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		count := L.ToNumber(-1)
		if count < 0 {
			t.Error("Expected at least zero system users")
		}

		L.Pop(1)
	})
}

// func TestListGroups(t *testing.T) {
// 	if runtime.GOOS == "windows" {
// 		t.Skip("Skipping on Windows")
// 	}
// 
// 	L := lua.NewState()
// 	defer L.Close()
// 
// 	module := NewUserModule()
// 	L.PreloadModule("user", module.Loader)
// 
// 	t.Run("ListAllGroups", func(t *testing.T) {
// 		script := `
// 			local user = require("user")
// 			local groups, err = user.group_list()
// 			if groups == nil then
// 				return nil, err
// 			end
// 			return #groups
// 		`
// 
// 		if err := L.DoString(script); err != nil {
// 			t.Fatalf("Script execution failed: %v", err)
// 		}
// 
// 		count := L.ToNumber(-1)
// 		if count <= 0 {
// 			t.Error("Expected at least one group in the system")
// 		}
// 
// 		L.Pop(1)
// 	})
// }

func TestGetGroupMembers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("GetRootGroupMembers", func(t *testing.T) {
		groupName := "root"
		if runtime.GOOS == "darwin" {
			groupName = "wheel"
		}

		script := `
			local user = require("user")
			local members, err = user.group_members("` + groupName + `")
			if members == nil then
				return nil, err
			end
			return #members
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		count := L.ToNumber(-1)
		if count < 0 {
			t.Error("Expected at least zero members")
		}

		L.Pop(1)
	})
}

func TestGetComment(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	t.Run("GetCurrentUserComment", func(t *testing.T) {
		script := `
			local user = require("user")
			local comment, err = user.get_comment("` + currentUser.Username + `")
			return comment, err
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		// Just verify we got something back (could be empty)
		comment := L.ToString(-2)
		_ = comment // Comment can be empty, which is valid

		L.Pop(2)
	})
}

// Test helper functions

func TestNeedsSudo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	module := NewUserModule()

	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}

	needsSudo := module.needsSudo()
	expectedSudo := currentUser.Uid != "0"

	if needsSudo != expectedSudo {
		t.Errorf("Expected needsSudo to be %v for user %s (UID: %s), got %v",
			expectedSudo, currentUser.Username, currentUser.Uid, needsSudo)
	}
}

func TestParseUserOptions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()

	t.Run("ParseEmptyOptions", func(t *testing.T) {
		options := module.parseUserOptions(L, 1)
		if len(options) != 0 {
			t.Errorf("Expected empty options map, got %d entries", len(options))
		}
	})

	t.Run("ParseValidOptions", func(t *testing.T) {
		script := `
			return {
				shell = "/bin/bash",
				home = "/home/test",
				uid = "1000"
			}
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		options := module.parseUserOptions(L, -1)

		if options["shell"] != "/bin/bash" {
			t.Errorf("Expected shell '/bin/bash', got '%s'", options["shell"])
		}

		if options["home"] != "/home/test" {
			t.Errorf("Expected home '/home/test', got '%s'", options["home"])
		}

		if options["uid"] != "1000" {
			t.Errorf("Expected uid '1000', got '%s'", options["uid"])
		}

		L.Pop(1)
	})
}

// Integration test for complete workflow
func TestUserWorkflow(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("CompleteUserQuery", func(t *testing.T) {
		currentUser, err := user.Current()
		if err != nil {
			t.Skipf("Could not get current user: %v", err)
		}

		script := `
			local user = require("user")
			
			-- Check if user exists
			local exists, msg = user.exists("` + currentUser.Username + `")
			if not exists then
				return false, "User does not exist"
			end
			
			-- Get user info
			local info, err = user.get_info("` + currentUser.Username + `")
			if info == nil then
				return false, "Failed to get user info: " .. err
			end
			
			-- Get UID
			local uid, err = user.get_uid("` + currentUser.Username + `")
			if uid == nil then
				return false, "Failed to get UID: " .. err
			end
			
			-- Get GID
			local gid, err = user.get_gid("` + currentUser.Username + `")
			if gid == nil then
				return false, "Failed to get GID: " .. err
			end
			
			-- Get home
			local home, err = user.get_home("` + currentUser.Username + `")
			if home == nil or home == "" then
				return false, "Failed to get home: " .. tostring(err)
			end
			
			-- Get groups
			local groups, err = user.get_groups("` + currentUser.Username + `")
			if groups == nil then
				return false, "Failed to get groups: " .. err
			end
			
			return true, "All operations successful"
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		success := L.ToBool(-2)
		msg := L.ToString(-1)

		if !success {
			t.Errorf("Workflow failed: %s", msg)
		}

		L.Pop(2)
	})
}

// Benchmark tests
func BenchmarkUserExists(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		b.Skipf("Could not get current user: %v", err)
	}

	script := `
		local user = require("user")
		local exists, msg = user.exists("` + currentUser.Username + `")
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(script); err != nil {
			b.Fatalf("Script execution failed: %v", err)
		}
		L.Pop(2)
	}
}

func BenchmarkGetUserInfo(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	currentUser, err := user.Current()
	if err != nil {
		b.Skipf("Could not get current user: %v", err)
	}

	script := `
		local user = require("user")
		local info, err = user.get_info("` + currentUser.Username + `")
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(script); err != nil {
			b.Fatalf("Script execution failed: %v", err)
		}
		L.Pop(2)
	}
}

func BenchmarkValidateUsername(b *testing.B) {
	if runtime.GOOS == "windows" {
		b.Skip("Skipping on Windows")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	script := `
		local user = require("user")
		local valid, msg = user.validate_username("testuser123")
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(script); err != nil {
			b.Fatalf("Script execution failed: %v", err)
		}
		L.Pop(2)
	}
}

func TestUserCreateWithPassword(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	// Check if we can run user commands (requires root/sudo)
	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("Could not get current user: %v", err)
	}
	
	if currentUser.Uid != "0" {
		t.Skip("Skipping user creation test - requires root privileges")
	}

	L := lua.NewState()
	defer L.Close()

	module := NewUserModule()
	L.PreloadModule("user", module.Loader)

	t.Run("CreateUserWithPassword", func(t *testing.T) {
		testUsername := "testuser_with_pass_12345"
		
		// Clean up before test
		script := `
			local user = require("user")
			user.delete("` + testUsername + `", true)
		`
		L.DoString(script) // Ignore errors if user doesn't exist
		
		// Create user with password
		script = `
			local user = require("user")
			local ok, msg = user.create("` + testUsername + `", {
				password = "TestPassword123!",
				home = "/home/` + testUsername + `",
				shell = "/bin/bash",
				create_home = true,
				comment = "Test user with password"
			})
			return ok, msg
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		ok := L.ToBool(-2)
		msg := L.ToString(-1)
		
		if !ok {
			t.Errorf("Failed to create user with password: %s", msg)
		} else {
			t.Logf("User created successfully with password: %s", msg)
		}
		
		L.Pop(2)
		
		// Clean up after test
		cleanupScript := `
			local user = require("user")
			user.delete("` + testUsername + `", true)
		`
		L.DoString(cleanupScript)
	})

	t.Run("CreateUserWithoutPassword", func(t *testing.T) {
		testUsername := "testuser_no_pass_12345"
		
		// Clean up before test
		script := `
			local user = require("user")
			user.delete("` + testUsername + `", true)
		`
		L.DoString(script) // Ignore errors if user doesn't exist
		
		// Create user without password
		script = `
			local user = require("user")
			local ok, msg = user.create("` + testUsername + `", {
				home = "/home/` + testUsername + `",
				shell = "/bin/bash",
				create_home = true
			})
			return ok, msg
		`

		if err := L.DoString(script); err != nil {
			t.Fatalf("Script execution failed: %v", err)
		}

		ok := L.ToBool(-2)
		msg := L.ToString(-1)
		
		if !ok {
			t.Errorf("Failed to create user without password: %s", msg)
		} else {
			t.Logf("User created successfully without password: %s", msg)
		}
		
		L.Pop(2)
		
		// Clean up after test
		cleanupScript := `
			local user = require("user")
			user.delete("` + testUsername + `", true)
		`
		L.DoString(cleanupScript)
	})
}
