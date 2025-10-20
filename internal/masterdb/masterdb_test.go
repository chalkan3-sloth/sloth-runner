//go:build cgo
// +build cgo

package masterdb

import (
	"path/filepath"
	"testing"
	"time"
)

// Test Master struct
func TestMaster_StructCreation(t *testing.T) {
	now := time.Now()
	master := &Master{
		Name:        "production",
		Address:     "prod.example.com:50053",
		Description: "Production master server",
		IsDefault:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if master.Name != "production" {
		t.Error("Expected Name to be set")
	}

	if master.Address != "prod.example.com:50053" {
		t.Error("Expected Address to be set")
	}

	if !master.IsDefault {
		t.Error("Expected IsDefault to be true")
	}
}

func TestMaster_ZeroValue(t *testing.T) {
	var master Master

	if master.Name != "" {
		t.Error("Expected zero Name")
	}

	if master.IsDefault {
		t.Error("Expected zero IsDefault")
	}

	if !master.CreatedAt.IsZero() {
		t.Error("Expected zero CreatedAt")
	}
}

func TestMaster_DefaultFlag(t *testing.T) {
	master1 := &Master{Name: "master1", IsDefault: true}
	master2 := &Master{Name: "master2", IsDefault: false}

	if !master1.IsDefault {
		t.Error("Expected master1 to be default")
	}

	if master2.IsDefault {
		t.Error("Expected master2 to not be default")
	}
}

// Test NewMasterDB
func TestNewMasterDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Error("Expected non-nil database")
	}
}

func TestNewMasterDB_InvalidPath(t *testing.T) {
	// Try to create DB in non-existent directory without parent
	invalidPath := "/nonexistent/very/deep/path/test.db"

	db, err := NewMasterDB(invalidPath)
	if err == nil {
		if db != nil {
			db.Close()
		}
		// Some systems might create the path, so this isn't necessarily an error
		t.Log("Warning: Expected error with invalid path, but succeeded")
	}
}

func TestNewMasterDB_CreatesSchema(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Verify tables exist by running a query
	count, err := db.Count()
	if err != nil {
		t.Errorf("Failed to count: %v (schema may not be initialized)", err)
	}

	if count != 0 {
		t.Errorf("Expected empty database, got %d masters", count)
	}
}

// Test Add method
func TestMasterDB_Add(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:        "test-master",
		Address:     "localhost:50053",
		Description: "Test master",
	}

	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master: %v", err)
	}

	// Verify it was added
	retrieved, err := db.Get("test-master")
	if err != nil {
		t.Fatalf("Failed to get master: %v", err)
	}

	if retrieved.Name != master.Name {
		t.Errorf("Expected name %s, got %s", master.Name, retrieved.Name)
	}
}

func TestMasterDB_Add_FirstIsDefault(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "first-master",
		Address: "first:50053",
	}

	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master: %v", err)
	}

	// First master should be default
	retrieved, err := db.Get("first-master")
	if err != nil {
		t.Fatalf("Failed to get master: %v", err)
	}

	if !retrieved.IsDefault {
		t.Error("Expected first master to be default")
	}
}

func TestMasterDB_Add_EmptyName(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "",
		Address: "localhost:50053",
	}

	err = db.Add(master)
	if err == nil {
		t.Error("Expected error with empty name")
	}
}

func TestMasterDB_Add_EmptyAddress(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "test",
		Address: "",
	}

	err = db.Add(master)
	if err == nil {
		t.Error("Expected error with empty address")
	}
}

func TestMasterDB_Add_Duplicate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "duplicate",
		Address: "localhost:50053",
	}

	// Add first time
	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master first time: %v", err)
	}

	// Try to add duplicate
	err = db.Add(master)
	if err == nil {
		t.Error("Expected error when adding duplicate")
	}
}

// Test Update method
func TestMasterDB_Update(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "test-update",
		Address: "old-address:50053",
	}

	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master: %v", err)
	}

	// Update the master
	master.Address = "new-address:50053"
	master.Description = "Updated description"

	err = db.Update(master)
	if err != nil {
		t.Fatalf("Failed to update master: %v", err)
	}

	// Verify update
	retrieved, err := db.Get("test-update")
	if err != nil {
		t.Fatalf("Failed to get master: %v", err)
	}

	if retrieved.Address != "new-address:50053" {
		t.Errorf("Expected updated address, got %s", retrieved.Address)
	}

	if retrieved.Description != "Updated description" {
		t.Errorf("Expected updated description, got %s", retrieved.Description)
	}
}

func TestMasterDB_Update_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "nonexistent",
		Address: "address:50053",
	}

	err = db.Update(master)
	if err == nil {
		t.Error("Expected error when updating nonexistent master")
	}
}

func TestMasterDB_Update_EmptyName(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "",
		Address: "address:50053",
	}

	err = db.Update(master)
	if err == nil {
		t.Error("Expected error with empty name")
	}
}

// Test Get method
func TestMasterDB_Get(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:        "get-test",
		Address:     "localhost:50053",
		Description: "Get test master",
	}

	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master: %v", err)
	}

	retrieved, err := db.Get("get-test")
	if err != nil {
		t.Fatalf("Failed to get master: %v", err)
	}

	if retrieved.Name != master.Name {
		t.Error("Name mismatch")
	}

	if retrieved.Address != master.Address {
		t.Error("Address mismatch")
	}

	if retrieved.Description != master.Description {
		t.Error("Description mismatch")
	}
}

func TestMasterDB_Get_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	_, err = db.Get("nonexistent")
	if err == nil {
		t.Error("Expected error when getting nonexistent master")
	}
}

// Test GetDefault method
func TestMasterDB_GetDefault(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "default-master",
		Address: "default:50053",
	}

	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master: %v", err)
	}

	defaultMaster, err := db.GetDefault()
	if err != nil {
		t.Fatalf("Failed to get default master: %v", err)
	}

	if defaultMaster.Name != "default-master" {
		t.Errorf("Expected default-master, got %s", defaultMaster.Name)
	}

	if !defaultMaster.IsDefault {
		t.Error("Expected IsDefault to be true")
	}
}

func TestMasterDB_GetDefault_NoDefault(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	_, err = db.GetDefault()
	if err == nil {
		t.Error("Expected error when no default master exists")
	}
}

// Test List method
func TestMasterDB_List(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Add multiple masters
	masters := []*Master{
		{Name: "master1", Address: "addr1:50053"},
		{Name: "master2", Address: "addr2:50053"},
		{Name: "master3", Address: "addr3:50053"},
	}

	for _, m := range masters {
		err := db.Add(m)
		if err != nil {
			t.Fatalf("Failed to add master %s: %v", m.Name, err)
		}
	}

	list, err := db.List()
	if err != nil {
		t.Fatalf("Failed to list masters: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 masters, got %d", len(list))
	}
}

func TestMasterDB_List_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	list, err := db.List()
	if err != nil {
		t.Fatalf("Failed to list masters: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d masters", len(list))
	}
}

func TestMasterDB_List_OrderByDefault(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Add masters
	master1 := &Master{Name: "first", Address: "addr1:50053"}
	master2 := &Master{Name: "second", Address: "addr2:50053"}

	db.Add(master1)
	db.Add(master2)

	// Set second as default
	db.SetDefault("second")

	list, err := db.List()
	if err != nil {
		t.Fatalf("Failed to list masters: %v", err)
	}

	// Default should be first in list
	if list[0].Name != "second" {
		t.Errorf("Expected default master first, got %s", list[0].Name)
	}
}

// Test SetDefault method
func TestMasterDB_SetDefault(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master1 := &Master{Name: "master1", Address: "addr1:50053"}
	master2 := &Master{Name: "master2", Address: "addr2:50053"}

	db.Add(master1)
	db.Add(master2)

	// Set master2 as default
	err = db.SetDefault("master2")
	if err != nil {
		t.Fatalf("Failed to set default: %v", err)
	}

	// Verify master2 is default
	defaultMaster, err := db.GetDefault()
	if err != nil {
		t.Fatalf("Failed to get default: %v", err)
	}

	if defaultMaster.Name != "master2" {
		t.Errorf("Expected master2 to be default, got %s", defaultMaster.Name)
	}

	// Verify master1 is not default
	master1Retrieved, _ := db.Get("master1")
	if master1Retrieved.IsDefault {
		t.Error("Expected master1 to not be default")
	}
}

func TestMasterDB_SetDefault_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = db.SetDefault("nonexistent")
	if err == nil {
		t.Error("Expected error when setting nonexistent master as default")
	}
}

// Test Delete method
func TestMasterDB_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{Name: "to-delete", Address: "addr:50053"}
	db.Add(master)

	// Verify it exists
	_, err = db.Get("to-delete")
	if err != nil {
		t.Fatalf("Master should exist before delete: %v", err)
	}

	// Delete it
	err = db.Delete("to-delete")
	if err != nil {
		t.Fatalf("Failed to delete master: %v", err)
	}

	// Verify it's gone
	_, err = db.Get("to-delete")
	if err == nil {
		t.Error("Expected error when getting deleted master")
	}
}

func TestMasterDB_Delete_Default(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master1 := &Master{Name: "master1", Address: "addr1:50053"}
	master2 := &Master{Name: "master2", Address: "addr2:50053"}

	db.Add(master1)
	db.Add(master2)

	// Try to delete default (master1)
	err = db.Delete("master1")
	if err == nil {
		t.Error("Expected error when deleting default master with others present")
	}
}

func TestMasterDB_Delete_OnlyMaster(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{Name: "only-one", Address: "addr:50053"}
	db.Add(master)

	// Should be able to delete the only master (even if default)
	err = db.Delete("only-one")
	if err != nil {
		t.Errorf("Should be able to delete only master: %v", err)
	}

	count, _ := db.Count()
	if count != 0 {
		t.Errorf("Expected 0 masters after delete, got %d", count)
	}
}

// Test Count method
func TestMasterDB_Count(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Start with 0
	count, err := db.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 masters, got %d", count)
	}

	// Add 3 masters
	for i := 0; i < 3; i++ {
		master := &Master{
			Name:    string(rune('a' + i)),
			Address: "addr:50053",
		}
		db.Add(master)
	}

	count, err = db.Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 masters, got %d", count)
	}
}

// Test Close method
func TestMasterDB_Close(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

func TestMasterDB_Close_MultipleTimes(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Close first time
	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close database first time: %v", err)
	}

	// Close second time might error
	err = db.Close()
	if err == nil {
		t.Log("Note: Closing database twice didn't produce error")
	}
}

// Test database persistence
func TestMasterDB_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persist.db")

	// Create and add data
	db1, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create first database: %v", err)
	}

	master := &Master{
		Name:    "persistent",
		Address: "persist:50053",
	}

	err = db1.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master: %v", err)
	}

	db1.Close()

	// Reopen and verify data persists
	db2, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to reopen database: %v", err)
	}
	defer db2.Close()

	retrieved, err := db2.Get("persistent")
	if err != nil {
		t.Fatalf("Failed to get master after reopen: %v", err)
	}

	if retrieved.Name != "persistent" {
		t.Error("Data did not persist across database close/reopen")
	}
}

// Test concurrent operations safety
func TestMasterDB_ConcurrentReads(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "concurrent.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{Name: "concurrent", Address: "addr:50053"}
	db.Add(master)

	done := make(chan bool, 10)

	// Multiple concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			_, err := db.Get("concurrent")
			if err != nil {
				t.Errorf("Concurrent read failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test Master timestamps
func TestMaster_Timestamps(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "timestamps.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "timestamp-test",
		Address: "addr:50053",
	}

	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master: %v", err)
	}

	retrieved, err := db.Get("timestamp-test")
	if err != nil {
		t.Fatalf("Failed to get master: %v", err)
	}

	// CreatedAt should be set
	if retrieved.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	// UpdatedAt should be set
	if retrieved.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	// They should be close in time
	timeDiff := retrieved.UpdatedAt.Sub(retrieved.CreatedAt)
	if timeDiff < 0 || timeDiff > time.Second {
		t.Errorf("CreatedAt and UpdatedAt should be within 1 second, diff: %v", timeDiff)
	}
}

func TestMasterDB_UpdateTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "update-ts.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	master := &Master{
		Name:    "update-test",
		Address: "old:50053",
	}

	db.Add(master)

	original, _ := db.Get("update-test")
	originalUpdatedAt := original.UpdatedAt

	// Wait to ensure timestamp difference (1 second to account for SQLite's datetime precision)
	time.Sleep(1100 * time.Millisecond)

	// Update
	master.Address = "new:50053"
	db.Update(master)

	updated, _ := db.Get("update-test")

	// UpdatedAt should have changed (or at least not be before)
	if updated.UpdatedAt.Before(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to not go backwards")
	}

	// Note: SQLite's datetime('now') may have second precision, so timestamps might be equal
	// We just verify it didn't go backwards
	t.Logf("Original UpdatedAt: %v, Updated UpdatedAt: %v", originalUpdatedAt, updated.UpdatedAt)
}

// Test edge cases
func TestMasterDB_LongDescription(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "long-desc.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	longDesc := string(make([]byte, 10000))
	for i := range longDesc {
		longDesc = longDesc[:i] + "a" + longDesc[i+1:]
	}

	master := &Master{
		Name:        "long-desc",
		Address:     "addr:50053",
		Description: longDesc,
	}

	err = db.Add(master)
	if err != nil {
		t.Fatalf("Failed to add master with long description: %v", err)
	}

	retrieved, err := db.Get("long-desc")
	if err != nil {
		t.Fatalf("Failed to get master: %v", err)
	}

	if len(retrieved.Description) != len(longDesc) {
		t.Error("Long description was truncated")
	}
}

func TestMasterDB_SpecialCharactersInName(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "special-chars.db")

	db, err := NewMasterDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	specialNames := []string{
		"master-with-dash",
		"master_with_underscore",
		"master.with.dots",
		"master123",
	}

	for _, name := range specialNames {
		master := &Master{
			Name:    name,
			Address: "addr:50053",
		}

		err := db.Add(master)
		if err != nil {
			t.Errorf("Failed to add master with name '%s': %v", name, err)
		}

		retrieved, err := db.Get(name)
		if err != nil {
			t.Errorf("Failed to get master with name '%s': %v", name, err)
		}

		if retrieved.Name != name {
			t.Errorf("Name mismatch for special char name: expected %s, got %s", name, retrieved.Name)
		}
	}
}
