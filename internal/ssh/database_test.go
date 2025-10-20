//go:build cgo
// +build cgo

package ssh

import (
	"path/filepath"
	"testing"
	"time"
)

// Test Profile struct
func TestProfile_Creation(t *testing.T) {
	now := time.Now()
	profile := &Profile{
		Name:               "test",
		Host:               "localhost",
		User:               "user",
		Port:               22,
		KeyPath:            "/path/to/key",
		Description:        "Test profile",
		CreatedAt:          now,
		UpdatedAt:          now,
		UseCount:           5,
		ConnectionTimeout:  30,
		KeepaliveInterval:  60,
		StrictHostChecking: true,
	}

	if profile.Name != "test" {
		t.Error("Expected Name to be set")
	}
	if profile.Port != 22 {
		t.Error("Expected Port to be 22")
	}
}

func TestProfile_ZeroValue(t *testing.T) {
	var profile Profile
	if profile.Port != 0 {
		t.Error("Expected zero Port")
	}
	if profile.StrictHostChecking {
		t.Error("Expected zero StrictHostChecking")
	}
}

// Test AuditLog struct
func TestAuditLog_Creation(t *testing.T) {
	log := &AuditLog{
		ID:               1,
		ProfileName:      "test",
		Action:           "execute",
		Command:          "ls -la",
		Timestamp:        time.Now(),
		Success:          true,
		DurationMs:       100,
		BytesTransferred: 1024,
	}

	if log.Action != "execute" {
		t.Error("Expected Action to be set")
	}
}

// Test NewDatabase
func TestNewDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Error("Expected non-nil database")
	}
}

// Test GetDefaultDatabasePath
func TestGetDefaultDatabasePath(t *testing.T) {
	path := GetDefaultDatabasePath()
	if path == "" {
		t.Error("Expected non-empty path")
	}
	if !filepath.IsAbs(path) && path != ".sloth-runner/ssh_profiles.db" {
		t.Error("Expected absolute path or relative default")
	}
}

// Test AddProfile
func TestDatabase_AddProfile(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "test-profile",
		Host: "localhost",
		User: "testuser",
		Port: 22,
	}

	err := db.AddProfile(profile)
	if err != nil {
		t.Fatalf("Failed to add profile: %v", err)
	}
}

// Test GetProfile
func TestDatabase_GetProfile(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "get-test",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	retrieved, err := db.GetProfile("get-test")
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}

	if retrieved.Name != "get-test" {
		t.Error("Profile name mismatch")
	}
}

func TestDatabase_GetProfile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	_, err := db.GetProfile("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent profile")
	}
}

// Test ListProfiles
func TestDatabase_ListProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profiles := []*Profile{
		{Name: "p1", Host: "host1", User: "user1", Port: 22},
		{Name: "p2", Host: "host2", User: "user2", Port: 22},
		{Name: "p3", Host: "host3", User: "user3", Port: 22},
	}

	for _, p := range profiles {
		db.AddProfile(p)
	}

	list, err := db.ListProfiles()
	if err != nil {
		t.Fatalf("Failed to list profiles: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 profiles, got %d", len(list))
	}
}

func TestDatabase_ListProfiles_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	list, err := db.ListProfiles()
	if err != nil {
		t.Fatalf("Failed to list profiles: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d profiles", len(list))
	}
}

// Test UpdateProfile
func TestDatabase_UpdateProfile(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "update-test",
		Host: "old-host",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	updates := map[string]interface{}{
		"host": "new-host",
	}

	err := db.UpdateProfile("update-test", updates)
	if err != nil {
		t.Fatalf("Failed to update profile: %v", err)
	}

	updated, _ := db.GetProfile("update-test")
	if updated.Host != "new-host" {
		t.Error("Host was not updated")
	}
}

func TestDatabase_UpdateProfile_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	updates := map[string]interface{}{"host": "new-host"}
	err := db.UpdateProfile("nonexistent", updates)
	if err == nil {
		t.Error("Expected error updating nonexistent profile")
	}
}

// Test RemoveProfile
func TestDatabase_RemoveProfile(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "to-delete",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	err := db.RemoveProfile("to-delete")
	if err != nil {
		t.Fatalf("Failed to remove profile: %v", err)
	}

	_, err = db.GetProfile("to-delete")
	if err == nil {
		t.Error("Profile should be deleted")
	}
}

// Test AddAuditLog
func TestDatabase_AddAuditLog(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "audit-test",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	log := &AuditLog{
		ProfileName: "audit-test",
		Action:      "execute",
		Command:     "ls",
		Success:     true,
		DurationMs:  50,
	}

	err := db.AddAuditLog(log)
	if err != nil {
		t.Fatalf("Failed to add audit log: %v", err)
	}
}

// Test GetAuditLogs
func TestDatabase_GetAuditLogs(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "log-test",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	// Add multiple logs
	for i := 0; i < 5; i++ {
		log := &AuditLog{
			ProfileName: "log-test",
			Action:      "execute",
			Success:     true,
		}
		db.AddAuditLog(log)
	}

	logs, err := db.GetAuditLogs("log-test", 10)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(logs) != 5 {
		t.Errorf("Expected 5 logs, got %d", len(logs))
	}
}

func TestDatabase_GetAuditLogs_Limit(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "limit-test",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	// Add 10 logs
	for i := 0; i < 10; i++ {
		log := &AuditLog{
			ProfileName: "limit-test",
			Action:      "execute",
			Success:     true,
		}
		db.AddAuditLog(log)
	}

	logs, _ := db.GetAuditLogs("limit-test", 3)
	if len(logs) != 3 {
		t.Errorf("Expected 3 logs (limited), got %d", len(logs))
	}
}

// Test Profile with optional fields
func TestProfile_OptionalFields(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name:        "optional-test",
		Host:        "localhost",
		User:        "user",
		Port:        22,
		KeyPath:     "/path/to/key",
		Description: "Test description",
	}

	db.AddProfile(profile)

	retrieved, _ := db.GetProfile("optional-test")
	if retrieved.KeyPath != "/path/to/key" {
		t.Error("KeyPath not preserved")
	}
	if retrieved.Description != "Test description" {
		t.Error("Description not preserved")
	}
}

// Test Profile with LastUsed
func TestProfile_LastUsed(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "lastused-test",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	// Add execution log to trigger LastUsed update
	log := &AuditLog{
		ProfileName: "lastused-test",
		Action:      "execute",
		Success:     true,
	}

	db.AddAuditLog(log)

	retrieved, _ := db.GetProfile("lastused-test")
	if retrieved.UseCount != 1 {
		t.Errorf("Expected UseCount 1, got %d", retrieved.UseCount)
	}
}

// Test Update different fields
func TestDatabase_UpdateProfile_DifferentFields(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "update-fields",
		Host: "host",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	testCases := []struct {
		field string
		value interface{}
	}{
		{"host", "newhost"},
		{"user", "newuser"},
		{"port", 2222},
		{"key_path", "/new/key"},
		{"description", "New description"},
		{"connection_timeout", 60},
		{"keepalive_interval", 120},
		{"strict_host_checking", false},
	}

	for _, tc := range testCases {
		updates := map[string]interface{}{tc.field: tc.value}
		err := db.UpdateProfile("update-fields", updates)
		if err != nil {
			t.Errorf("Failed to update %s: %v", tc.field, err)
		}
	}
}

// Test Database Close
func TestDatabase_Close(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)

	err := db.Close()
	if err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

// Test Database persistence
func TestDatabase_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persist.db")

	db1, _ := NewDatabase(dbPath)
	profile := &Profile{
		Name: "persist",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db1.AddProfile(profile)
	db1.Close()

	// Reopen
	db2, _ := NewDatabase(dbPath)
	defer db2.Close()

	retrieved, err := db2.GetProfile("persist")
	if err != nil {
		t.Error("Profile did not persist")
	}

	if retrieved.Name != "persist" {
		t.Error("Data corruption after reopen")
	}
}

// Test Port range
func TestProfile_PortRange(t *testing.T) {
	ports := []int{22, 80, 443, 8080, 65535}

	for _, port := range ports {
		profile := &Profile{
			Name: "port-test",
			Host: "localhost",
			User: "user",
			Port: port,
		}

		if profile.Port != port {
			t.Errorf("Port not set correctly: %d", port)
		}
	}
}

// Test Connection settings
func TestProfile_ConnectionSettings(t *testing.T) {
	profile := &Profile{
		Name:               "conn-test",
		Host:               "localhost",
		User:               "user",
		Port:               22,
		ConnectionTimeout:  30,
		KeepaliveInterval:  60,
		StrictHostChecking: true,
	}

	if profile.ConnectionTimeout != 30 {
		t.Error("ConnectionTimeout not set")
	}
	if profile.KeepaliveInterval != 60 {
		t.Error("KeepaliveInterval not set")
	}
	if !profile.StrictHostChecking {
		t.Error("StrictHostChecking not set")
	}
}

// Test AuditLog fields
func TestAuditLog_AllFields(t *testing.T) {
	log := &AuditLog{
		ID:               1,
		ProfileName:      "test",
		Action:           "execute",
		Command:          "ls -la",
		Timestamp:        time.Now(),
		Success:          true,
		ErrorMessage:     "",
		DurationMs:       100,
		BytesTransferred: 1024,
	}

	if log.ProfileName != "test" {
		t.Error("ProfileName not set")
	}
	if log.DurationMs != 100 {
		t.Error("DurationMs not set")
	}
	if log.BytesTransferred != 1024 {
		t.Error("BytesTransferred not set")
	}
}

// Test AuditLog with error
func TestAuditLog_WithError(t *testing.T) {
	log := &AuditLog{
		ProfileName:  "test",
		Action:       "execute",
		Success:      false,
		ErrorMessage: "Connection failed",
	}

	if log.Success {
		t.Error("Expected Success to be false")
	}
	if log.ErrorMessage != "Connection failed" {
		t.Error("ErrorMessage not set")
	}
}

// Test multiple actions
func TestDatabase_MultipleActions(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "actions-test",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	actions := []string{"execute", "connect", "disconnect", "upload", "download"}

	for _, action := range actions {
		log := &AuditLog{
			ProfileName: "actions-test",
			Action:      action,
			Success:     true,
		}
		db.AddAuditLog(log)
	}

	logs, _ := db.GetAuditLogs("actions-test", 10)
	if len(logs) != 5 {
		t.Errorf("Expected 5 action logs, got %d", len(logs))
	}
}

// Test UpdateProfile no changes
func TestDatabase_UpdateProfile_NoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "nochange",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	// Empty updates map
	err := db.UpdateProfile("nochange", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error with no fields to update")
	}
}

// Test special characters in profile name
func TestProfile_SpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	names := []string{
		"profile-with-dash",
		"profile_with_underscore",
		"profile123",
	}

	for i, name := range names {
		profile := &Profile{
			Name: name,
			Host: "localhost",
			User: "user",
			Port: 22 + i, // Different port for each to avoid UNIQUE constraint
		}

		err := db.AddProfile(profile)
		if err != nil {
			t.Errorf("Failed to add profile with name '%s': %v", name, err)
		}

		retrieved, err := db.GetProfile(name)
		if err != nil || retrieved.Name != name {
			t.Errorf("Failed to retrieve profile with name '%s'", name)
		}
	}
}

// Test timestamps
func TestProfile_Timestamps(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "timestamp-test",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	retrieved, _ := db.GetProfile("timestamp-test")

	if retrieved.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	if retrieved.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

// Test concurrent reads
func TestDatabase_ConcurrentReads(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, _ := NewDatabase(dbPath)
	defer db.Close()

	profile := &Profile{
		Name: "concurrent",
		Host: "localhost",
		User: "user",
		Port: 22,
	}

	db.AddProfile(profile)

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			_, err := db.GetProfile("concurrent")
			if err != nil {
				t.Errorf("Concurrent read failed: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
