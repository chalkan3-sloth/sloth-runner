package users

import (
	"os/user"
	"testing"
)

func TestNewUserManager(t *testing.T) {
	manager := NewUserManager()
	if manager == nil {
		t.Fatal("NewUserManager() returned nil")
	}

	_, ok := manager.(*SystemUserManager)
	if !ok {
		t.Error("NewUserManager() did not return *SystemUserManager")
	}
}

func TestList(t *testing.T) {
	manager := NewUserManager()

	options := ListOptions{
		SystemUsers: false,
	}

	users, err := manager.List(options)
	if err != nil {
		// getent não existe no macOS, então skip o teste
		if err.Error() == "failed to list users: exec: \"getent\": executable file not found in $PATH" {
			t.Skip("getent not available (not a Linux system)")
		}
		t.Fatalf("List() failed: %v", err)
	}

	if users == nil {
		t.Fatal("List() returned nil users")
	}

	// Deve ter pelo menos um usuário regular
	if len(users) == 0 {
		t.Error("List() returned no users")
	}

	// Verifica campos básicos
	for _, u := range users {
		if u.Username == "" {
			t.Error("User has empty username")
		}
		if u.UID == "" {
			t.Error("User has empty UID")
		}
		if u.HomeDir == "" {
			t.Error("User has empty HomeDir")
		}
	}
}

func TestListWithSystemUsers(t *testing.T) {
	manager := NewUserManager()

	options := ListOptions{
		SystemUsers: true,
	}

	users, err := manager.List(options)
	if err != nil {
		// getent não existe no macOS
		if err.Error() == "failed to list users: exec: \"getent\": executable file not found in $PATH" {
			t.Skip("getent not available (not a Linux system)")
		}
		t.Fatalf("List() with system users failed: %v", err)
	}

	// Com system users deve ter mais usuários
	if len(users) == 0 {
		t.Error("List() with system users returned no users")
	}
}

func TestListWithFilter(t *testing.T) {
	manager := NewUserManager()

	// Pega um usuário para usar como filtro
	allUsers, _ := manager.List(ListOptions{SystemUsers: true})
	if len(allUsers) == 0 {
		t.Skip("No users available for filter test")
	}

	filterName := allUsers[0].Username

	options := ListOptions{
		SystemUsers: true,
		Filter:      filterName,
	}

	users, err := manager.List(options)
	if err != nil {
		t.Fatalf("List() with filter failed: %v", err)
	}

	// Deve encontrar pelo menos o usuário filtrado
	found := false
	for _, u := range users {
		if u.Username == filterName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("List() with filter did not find user %s", filterName)
	}
}

func TestInfo(t *testing.T) {
	// Pega usuário atual
	currentUser, err := user.Current()
	if err != nil {
		t.Skip("Cannot get current user")
	}

	manager := NewUserManager()

	detail, err := manager.Info(currentUser.Username)
	if err != nil {
		t.Fatalf("Info() failed: %v", err)
	}

	if detail == nil {
		t.Fatal("Info() returned nil detail")
	}

	if detail.Username != currentUser.Username {
		t.Errorf("Info() returned username %s, expected %s", detail.Username, currentUser.Username)
	}

	if detail.UID == "" {
		t.Error("Info() returned empty UID")
	}

	if detail.HomeDir == "" {
		t.Error("Info() returned empty HomeDir")
	}

	// Groups deve ter pelo menos um grupo
	if len(detail.Groups) == 0 {
		t.Error("Info() returned no groups")
	}
}

func TestInfoNonExistentUser(t *testing.T) {
	manager := NewUserManager()

	_, err := manager.Info("nonexistentuser12345")
	if err == nil {
		t.Error("Info() with non-existent user should return error")
	}
}

func TestListGroups(t *testing.T) {
	manager := NewUserManager()

	groups, err := manager.ListGroups()
	if err != nil {
		// getent não existe no macOS
		if err.Error() == "failed to list groups: exec: \"getent\": executable file not found in $PATH" {
			t.Skip("getent not available (not a Linux system)")
		}
		t.Fatalf("ListGroups() failed: %v", err)
	}

	if groups == nil {
		t.Fatal("ListGroups() returned nil")
	}

	// Deve ter pelo menos um grupo
	if len(groups) == 0 {
		t.Error("ListGroups() returned no groups")
	}

	// Verifica campos básicos
	for _, g := range groups {
		if g.Name == "" {
			t.Error("Group has empty name")
		}
		if g.GID == "" {
			t.Error("Group has empty GID")
		}
		// Members pode ser vazio, então não verificamos
	}
}

func TestUserInfoStructure(t *testing.T) {
	info := &UserInfo{
		Username: "testuser",
		UID:      "1000",
		GID:      "1000",
		HomeDir:  "/home/testuser",
		Shell:    "/bin/bash",
		FullName: "Test User",
	}

	if info.Username != "testuser" {
		t.Error("Username not set correctly")
	}
	if info.UID != "1000" {
		t.Error("UID not set correctly")
	}
	if info.HomeDir != "/home/testuser" {
		t.Error("HomeDir not set correctly")
	}
}

func TestUserDetailStructure(t *testing.T) {
	detail := &UserDetail{
		UserInfo: &UserInfo{
			Username: "testuser",
			UID:      "1000",
		},
		Groups:      []string{"users", "sudo"},
		PasswordSet: true,
		Locked:      false,
		ExpiryDate:  "never",
	}

	if detail.Username != "testuser" {
		t.Error("Username not accessible through UserDetail")
	}
	if len(detail.Groups) != 2 {
		t.Error("Groups not set correctly")
	}
	if !detail.PasswordSet {
		t.Error("PasswordSet not set correctly")
	}
	if detail.Locked {
		t.Error("Locked should be false")
	}
}

func TestGroupInfoStructure(t *testing.T) {
	group := &GroupInfo{
		Name:    "testgroup",
		GID:     "1000",
		Members: []string{"user1", "user2"},
	}

	if group.Name != "testgroup" {
		t.Error("Name not set correctly")
	}
	if group.GID != "1000" {
		t.Error("GID not set correctly")
	}
	if len(group.Members) != 2 {
		t.Error("Members not set correctly")
	}
}

func TestAddUserOptionsStructure(t *testing.T) {
	options := AddUserOptions{
		Username:   "newuser",
		FullName:   "New User",
		HomeDir:    "/home/newuser",
		Shell:      "/bin/bash",
		Groups:     []string{"users", "docker"},
		CreateHome: true,
		System:     false,
	}

	if options.Username != "newuser" {
		t.Error("Username not set correctly")
	}
	if len(options.Groups) != 2 {
		t.Error("Groups not set correctly")
	}
	if !options.CreateHome {
		t.Error("CreateHome should be true")
	}
}

func TestModifyOptionsStructure(t *testing.T) {
	options := ModifyOptions{
		FullName:    "Updated Name",
		HomeDir:     "/new/home",
		Shell:       "/bin/zsh",
		Lock:        true,
		Unlock:      false,
		ExpireDate:  "2025-12-31",
		AddGroups:   []string{"newgroup"},
		RemoveGroup: "oldgroup",
	}

	if options.FullName != "Updated Name" {
		t.Error("FullName not set correctly")
	}
	if !options.Lock {
		t.Error("Lock should be true")
	}
	if len(options.AddGroups) != 1 {
		t.Error("AddGroups not set correctly")
	}
}

func TestListOptionsStructure(t *testing.T) {
	options := ListOptions{
		SystemUsers: true,
		Filter:      "test",
		Group:       "sudo",
	}

	if !options.SystemUsers {
		t.Error("SystemUsers should be true")
	}
	if options.Filter != "test" {
		t.Error("Filter not set correctly")
	}
	if options.Group != "sudo" {
		t.Error("Group not set correctly")
	}
}

func TestIsUserInGroup(t *testing.T) {
	manager := &SystemUserManager{}

	// Pega usuário atual
	currentUser, err := user.Current()
	if err != nil {
		t.Skip("Cannot get current user")
	}

	// Pega grupos do usuário atual
	currentUserObj, err := user.Lookup(currentUser.Username)
	if err != nil {
		t.Skip("Cannot lookup current user")
	}

	groupIDs, err := currentUserObj.GroupIds()
	if err != nil || len(groupIDs) == 0 {
		t.Skip("Cannot get user groups")
	}

	// Pega nome do primeiro grupo
	firstGroup, err := user.LookupGroupId(groupIDs[0])
	if err != nil {
		t.Skip("Cannot lookup group")
	}

	// Testa com grupo que usuário pertence
	inGroup, err := manager.isUserInGroup(currentUser.Username, firstGroup.Name)
	if err != nil {
		t.Fatalf("isUserInGroup() failed: %v", err)
	}

	if !inGroup {
		t.Errorf("isUserInGroup() returned false for user in group")
	}

	// Testa com grupo que usuário não pertence
	inGroup, err = manager.isUserInGroup(currentUser.Username, "nonexistentgroup12345")
	if err != nil {
		// Erro é aceitável para grupo inexistente
		return
	}

	if inGroup {
		t.Error("isUserInGroup() returned true for user not in group")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"this is a very long username", 10, "this is..."},
		{"exactly10!", 10, "exactly10!"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
		if len(got) > tt.maxLen {
			t.Errorf("truncate(%q, %d) returned string longer than maxLen", tt.input, tt.maxLen)
		}
	}
}

// Note: Tests for Add, Remove, Modify, AddToGroup, RemoveFromGroup
// are not included here as they require root/sudo permissions and
// would modify the actual system. These should be tested in a
// controlled environment or with mock implementations.

func TestManagerMethodsExist(t *testing.T) {
	var manager UserManager = &SystemUserManager{}

	// Verifica que todos os métodos da interface existem
	_ = manager.List
	_ = manager.Info
	_ = manager.Add
	_ = manager.Remove
	_ = manager.Modify
	_ = manager.ListGroups
	_ = manager.AddToGroup
	_ = manager.RemoveFromGroup
}
