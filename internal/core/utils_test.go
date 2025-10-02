package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creation without compression
	cache, err := NewFileCache(tmpDir, 10, false)
	if err != nil {
		t.Fatalf("Failed to create file cache: %v", err)
	}

	if cache == nil {
		t.Error("Expected non-nil cache")
	}

	stats := cache.Stats()
	if stats.MaxSize != 10*1024*1024 {
		t.Errorf("Expected max size 10MB, got %d", stats.MaxSize)
	}

	// Test creation with compression
	cache2, err := NewFileCache(filepath.Join(tmpDir, "compressed"), 5, true)
	if err != nil {
		t.Fatalf("Failed to create compressed cache: %v", err)
	}

	if !cache2.compress {
		t.Error("Expected compression to be enabled")
	}
}

func TestFileCache_PutGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache, err := NewFileCache(tmpDir, 10, false)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test Put and Get
	testData := []byte("test data content")
	err = cache.Put("test-key", testData, 0)
	if err != nil {
		t.Fatalf("Failed to put data: %v", err)
	}

	retrieved, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached data")
	}

	if string(retrieved) != string(testData) {
		t.Errorf("Expected %s, got %s", string(testData), string(retrieved))
	}
}

func TestFileCache_Delete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache, err := NewFileCache(tmpDir, 10, false)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Put data
	err = cache.Put("test-key", []byte("data"), 0)
	if err != nil {
		t.Fatalf("Failed to put data: %v", err)
	}

	// Delete
	cache.Delete("test-key")

	// Verify deletion
	_, found := cache.Get("test-key")
	if found {
		t.Error("Expected data to be deleted")
	}
}

func TestFileCache_Clear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache, err := NewFileCache(tmpDir, 10, false)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Put multiple entries
	cache.Put("key1", []byte("data1"), 0)
	cache.Put("key2", []byte("data2"), 0)
	cache.Put("key3", []byte("data3"), 0)

	// Clear
	cache.Clear()

	// Verify all are gone
	stats := cache.Stats()
	if stats.Entries != 0 {
		t.Errorf("Expected 0 entries, got %d", stats.Entries)
	}
}

func TestFileCache_Stats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache, err := NewFileCache(tmpDir, 10, false)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Initial stats
	stats := cache.Stats()
	if stats.Entries != 0 {
		t.Errorf("Expected 0 entries initially, got %d", stats.Entries)
	}

	// Add data
	cache.Put("key1", []byte("data"), 0)

	// Check stats
	stats = cache.Stats()
	if stats.Entries != 1 {
		t.Errorf("Expected 1 entry, got %d", stats.Entries)
	}

	if stats.CurrentSize == 0 {
		t.Error("Expected current size > 0")
	}
}

func TestFileCache_Compression(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache, err := NewFileCache(tmpDir, 10, true)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Put large compressible data
	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = 'A' // Highly compressible
	}

	err = cache.Put("compressed-key", largeData, 0)
	if err != nil {
		t.Fatalf("Failed to put compressed data: %v", err)
	}

	// Retrieve
	retrieved, found := cache.Get("compressed-key")
	if !found {
		t.Error("Expected to find compressed data")
	}

	if len(retrieved) != len(largeData) {
		t.Errorf("Expected %d bytes, got %d", len(largeData), len(retrieved))
	}
}

func TestSecureRandom_Bytes(t *testing.T) {
	sr := NewSecureRandom()

	// Test basic functionality
	bytes, err := sr.Bytes(16)
	if err != nil {
		t.Fatalf("Failed to generate random bytes: %v", err)
	}

	if len(bytes) != 16 {
		t.Errorf("Expected 16 bytes, got %d", len(bytes))
	}

	// Test that it generates different bytes each time
	bytes2, err := sr.Bytes(16)
	if err != nil {
		t.Fatalf("Failed to generate random bytes: %v", err)
	}

	same := true
	for i := range bytes {
		if bytes[i] != bytes2[i] {
			same = false
			break
		}
	}

	if same {
		t.Error("Expected different random bytes on each call")
	}
}

func TestSecureRandom_String(t *testing.T) {
	sr := NewSecureRandom()

	// Test string generation
	str, err := sr.String(32)
	if err != nil {
		t.Fatalf("Failed to generate random string: %v", err)
	}

	if len(str) != 32 {
		t.Errorf("Expected length 32, got %d", len(str))
	}

	// Test uniqueness
	str2, err := sr.String(32)
	if err != nil {
		t.Fatalf("Failed to generate random string: %v", err)
	}

	if str == str2 {
		t.Error("Expected different random strings")
	}
}

func TestSecureRandom_Int(t *testing.T) {
	sr := NewSecureRandom()

	// Test basic functionality
	n, err := sr.Int(100)
	if err != nil {
		t.Fatalf("Failed to generate random int: %v", err)
	}

	if n < 0 || n >= 100 {
		t.Errorf("Expected value in [0, 100), got %d", n)
	}

	// Test with zero max
	n, err = sr.Int(0)
	if err != nil {
		t.Fatalf("Failed with zero max: %v", err)
	}

	if n != 0 {
		t.Errorf("Expected 0 with zero max, got %d", n)
	}
}

func TestPathUtil_SafeJoin(t *testing.T) {
	pu := &PathUtil{}

	// Test normal join
	path, err := pu.SafeJoin("/tmp", "test", "file.txt")
	if err != nil {
		t.Fatalf("Failed to join path: %v", err)
	}

	expected := filepath.Join("/tmp", "test", "file.txt")
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}

	// Test path traversal attempt
	_, err = pu.SafeJoin("/tmp", "..", "..", "etc", "passwd")
	if err == nil {
		t.Error("Expected error for path traversal attempt")
	}
}

func TestPathUtil_EnsureDir(t *testing.T) {
	pu := &PathUtil{}

	tmpDir, err := os.MkdirTemp("", "test_path_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testDir := filepath.Join(tmpDir, "nested", "dir")
	err = pu.EnsureDir(testDir)
	if err != nil {
		t.Fatalf("Failed to ensure dir: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(testDir)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected path to be a directory")
	}
}

func TestPathUtil_IsSubPath(t *testing.T) {
	pu := &PathUtil{}

	tests := []struct {
		parent   string
		child    string
		expected bool
	}{
		{"/tmp", "/tmp/test", true},
		{"/tmp", "/tmp/test/nested", true},
		{"/tmp", "/tmp", true},
		{"/tmp", "/var", false},
		{"/tmp", "/tmpfile", false},
	}

	for _, tt := range tests {
		result := pu.IsSubPath(tt.parent, tt.child)
		if result != tt.expected {
			t.Errorf("IsSubPath(%s, %s) = %v, expected %v", 
				tt.parent, tt.child, result, tt.expected)
		}
	}
}

func TestPathUtil_TempDir(t *testing.T) {
	pu := &PathUtil{}

	dir, err := pu.TempDir("test_pattern_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Verify directory exists
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Temp directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected temp path to be a directory")
	}
}

func TestPathUtil_RemoveContents(t *testing.T) {
	pu := &PathUtil{}

	// Create test directory structure
	tmpDir, err := os.MkdirTemp("", "test_remove_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some files and directories
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("content"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	// Remove contents
	err = pu.RemoveContents(tmpDir)
	if err != nil {
		t.Fatalf("Failed to remove contents: %v", err)
	}

	// Verify directory is empty but still exists
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected empty directory, got %d entries", len(entries))
	}

	// Verify directory itself still exists
	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("Directory was removed: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected path to still be a directory")
	}
}
