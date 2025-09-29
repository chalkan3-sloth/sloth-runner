package core

import (
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// FileCache provides efficient file caching with compression and expiration
type FileCache struct {
	mu          sync.RWMutex
	cache       map[string]*CacheEntry
	maxSize     int64
	currentSize int64
	baseDir     string
	compress    bool
}

type CacheEntry struct {
	Key        string
	Path       string
	Size       int64
	AccessTime time.Time
	CreateTime time.Time
	TTL        time.Duration
	Compressed bool
	Hash       string
}

// NewFileCache creates a new file cache
func NewFileCache(baseDir string, maxSizeMB int64, compress bool) (*FileCache, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	
	fc := &FileCache{
		cache:    make(map[string]*CacheEntry),
		maxSize:  maxSizeMB * 1024 * 1024,
		baseDir:  baseDir,
		compress: compress,
	}
	
	// Load existing cache entries
	fc.loadExistingEntries()
	
	// Start cleanup goroutine
	go fc.cleanupLoop()
	
	return fc, nil
}

// Put stores data in the cache
func (fc *FileCache) Put(key string, data []byte, ttl time.Duration) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	
	// Generate unique filename
	hash := fc.generateHash(key)
	filename := hash + ".cache"
	if fc.compress {
		filename += ".gz"
	}
	
	path := filepath.Join(fc.baseDir, filename)
	
	// Write data to file
	var err error
	var size int64
	
	if fc.compress {
		size, err = fc.writeCompressed(path, data)
	} else {
		size, err = fc.writeUncompressed(path, data)
	}
	
	if err != nil {
		return err
	}
	
	// Remove old entry if exists
	if oldEntry, exists := fc.cache[key]; exists {
		fc.currentSize -= oldEntry.Size
		os.Remove(oldEntry.Path)
	}
	
	// Check if we need to evict entries
	if fc.currentSize+size > fc.maxSize {
		fc.evictLRU(size)
	}
	
	// Add new entry
	entry := &CacheEntry{
		Key:        key,
		Path:       path,
		Size:       size,
		AccessTime: time.Now(),
		CreateTime: time.Now(),
		TTL:        ttl,
		Compressed: fc.compress,
		Hash:       hash,
	}
	
	fc.cache[key] = entry
	fc.currentSize += size
	
	return nil
}

// Get retrieves data from the cache
func (fc *FileCache) Get(key string) ([]byte, bool) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	
	entry, exists := fc.cache[key]
	if !exists {
		return nil, false
	}
	
	// Check if expired
	if entry.TTL > 0 && time.Since(entry.CreateTime) > entry.TTL {
		fc.removeEntry(key)
		return nil, false
	}
	
	// Check if file still exists
	if _, err := os.Stat(entry.Path); os.IsNotExist(err) {
		fc.removeEntry(key)
		return nil, false
	}
	
	// Read data
	var data []byte
	var err error
	
	if entry.Compressed {
		data, err = fc.readCompressed(entry.Path)
	} else {
		data, err = fc.readUncompressed(entry.Path)
	}
	
	if err != nil {
		fc.removeEntry(key)
		return nil, false
	}
	
	// Update access time
	entry.AccessTime = time.Now()
	
	return data, true
}

// Delete removes an entry from the cache
func (fc *FileCache) Delete(key string) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	
	fc.removeEntry(key)
}

// Clear removes all entries from the cache
func (fc *FileCache) Clear() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	
	for key := range fc.cache {
		fc.removeEntry(key)
	}
}

// Stats returns cache statistics
func (fc *FileCache) Stats() CacheStats {
	fc.mu.RLock()
	defer fc.mu.RUnlock()
	
	return CacheStats{
		Entries:     len(fc.cache),
		CurrentSize: fc.currentSize,
		MaxSize:     fc.maxSize,
		UsageRatio:  float64(fc.currentSize) / float64(fc.maxSize),
	}
}

type CacheStats struct {
	Entries     int
	CurrentSize int64
	MaxSize     int64
	UsageRatio  float64
}

// Helper methods

func (fc *FileCache) generateHash(key string) string {
	h := fnv.New64a()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func (fc *FileCache) writeCompressed(path string, data []byte) (int64, error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	
	writer := gzip.NewWriter(file)
	defer writer.Close()
	
	_, err = writer.Write(data)
	if err != nil {
		return 0, err
	}
	
	if err := writer.Close(); err != nil {
		return 0, err
	}
	
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	
	return stat.Size(), nil
}

func (fc *FileCache) writeUncompressed(path string, data []byte) (int64, error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	
	n, err := file.Write(data)
	if err != nil {
		return 0, err
	}
	
	return int64(n), nil
}

func (fc *FileCache) readCompressed(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	
	return io.ReadAll(reader)
}

func (fc *FileCache) readUncompressed(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fc *FileCache) removeEntry(key string) {
	entry, exists := fc.cache[key]
	if !exists {
		return
	}
	
	os.Remove(entry.Path)
	fc.currentSize -= entry.Size
	delete(fc.cache, key)
}

func (fc *FileCache) evictLRU(neededSize int64) {
	// Sort entries by access time (oldest first)
	var entries []*CacheEntry
	for _, entry := range fc.cache {
		entries = append(entries, entry)
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].AccessTime.Before(entries[j].AccessTime)
	})
	
	freedSize := int64(0)
	for _, entry := range entries {
		if fc.currentSize-freedSize+neededSize <= fc.maxSize {
			break
		}
		
		fc.removeEntry(entry.Key)
		freedSize += entry.Size
	}
}

func (fc *FileCache) loadExistingEntries() {
	files, err := os.ReadDir(fc.baseDir)
	if err != nil {
		return
	}
	
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".cache") && !strings.HasSuffix(file.Name(), ".cache.gz") {
			continue
		}
		
		path := filepath.Join(fc.baseDir, file.Name())
		info, err := file.Info()
		if err != nil {
			continue
		}
		
		// Create cache entry for existing file
		hash := strings.TrimSuffix(file.Name(), ".cache")
		hash = strings.TrimSuffix(hash, ".gz")
		
		entry := &CacheEntry{
			Key:        "", // Unknown, will be populated on first access
			Path:       path,
			Size:       info.Size(),
			AccessTime: info.ModTime(),
			CreateTime: info.ModTime(),
			Compressed: strings.HasSuffix(file.Name(), ".gz"),
			Hash:       hash,
		}
		
		fc.currentSize += entry.Size
	}
}

func (fc *FileCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		fc.cleanup()
	}
}

func (fc *FileCache) cleanup() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	
	now := time.Now()
	var toDelete []string
	
	for key, entry := range fc.cache {
		// Remove expired entries
		if entry.TTL > 0 && now.Sub(entry.CreateTime) > entry.TTL {
			toDelete = append(toDelete, key)
		}
		
		// Remove entries for non-existent files
		if _, err := os.Stat(entry.Path); os.IsNotExist(err) {
			toDelete = append(toDelete, key)
		}
	}
	
	for _, key := range toDelete {
		fc.removeEntry(key)
	}
}

// SecureRandom generates cryptographically secure random bytes
type SecureRandom struct{}

// NewSecureRandom creates a new secure random generator
func NewSecureRandom() *SecureRandom {
	return &SecureRandom{}
}

// Bytes generates n random bytes
func (sr *SecureRandom) Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

// String generates a random hex string of specified length
func (sr *SecureRandom) String(length int) (string, error) {
	bytes, err := sr.Bytes(length / 2)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Int generates a random integer in range [0, max)
func (sr *SecureRandom) Int(max int64) (int64, error) {
	if max <= 0 {
		return 0, nil
	}
	
	// Use rejection sampling to ensure uniform distribution
	b := make([]byte, 8)
	for {
		_, err := rand.Read(b)
		if err != nil {
			return 0, err
		}
		
		// Convert bytes to uint64
		n := uint64(0)
		for i := 0; i < 8; i++ {
			n = n<<8 + uint64(b[i])
		}
		
		// Check if within acceptable range
		maxUint := ^uint64(0)
		limit := maxUint - (maxUint % uint64(max))
		
		if n < limit {
			return int64(n % uint64(max)), nil
		}
	}
}

// PathUtil provides utility functions for path manipulation
type PathUtil struct{}

// SafeJoin safely joins path elements, preventing directory traversal
func (pu *PathUtil) SafeJoin(base string, elements ...string) (string, error) {
	path := filepath.Join(append([]string{base}, elements...)...)
	cleanPath := filepath.Clean(path)
	
	// Ensure the result is still within the base directory
	if !strings.HasPrefix(cleanPath, filepath.Clean(base)) {
		return "", NewSlothError("PATH_TRAVERSAL", 
			"Path traversal attempt detected", 
			SeverityHigh).WithDetail("base", base).WithDetail("path", path)
	}
	
	return cleanPath, nil
}

// EnsureDir ensures a directory exists, creating it if necessary
func (pu *PathUtil) EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsSubPath checks if child is a subdirectory of parent
func (pu *PathUtil) IsSubPath(parent, child string) bool {
	cleanParent := filepath.Clean(parent)
	cleanChild := filepath.Clean(child)
	
	return strings.HasPrefix(cleanChild, cleanParent+string(filepath.Separator)) ||
		   cleanChild == cleanParent
}

// TempDir creates a temporary directory with a specific pattern
func (pu *PathUtil) TempDir(pattern string) (string, error) {
	return os.MkdirTemp("", pattern)
}

// RemoveContents removes all contents of a directory but keeps the directory
func (pu *PathUtil) RemoveContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	
	return nil
}