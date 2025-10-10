package agent

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// WatcherType defines the type of watcher
type WatcherType string

const (
	WatcherTypeFile       WatcherType = "file"
	WatcherTypeDirectory  WatcherType = "directory"
	WatcherTypeProcess    WatcherType = "process"
	WatcherTypePort       WatcherType = "port"
	WatcherTypeService    WatcherType = "service"
	WatcherTypeLog        WatcherType = "log"
	WatcherTypeCommand    WatcherType = "command"
	WatcherTypeCPU        WatcherType = "cpu"
	WatcherTypeMemory     WatcherType = "memory"
	WatcherTypeDisk       WatcherType = "disk"
	WatcherTypeNetwork    WatcherType = "network"
	WatcherTypeConnection WatcherType = "connection"
	WatcherTypeUser       WatcherType = "user"
	WatcherTypePackage    WatcherType = "package"
	WatcherTypeCustom     WatcherType = "custom"
)

// EventCondition defines when to trigger events
type EventCondition string

const (
	ConditionChanged   EventCondition = "changed"
	ConditionCreated   EventCondition = "created"
	ConditionDeleted   EventCondition = "deleted"
	ConditionExists    EventCondition = "exists"
	ConditionAbove     EventCondition = "above"     // Value above threshold
	ConditionBelow     EventCondition = "below"     // Value below threshold
	ConditionMatches   EventCondition = "matches"   // Pattern matches
	ConditionContains  EventCondition = "contains"  // Contains string/pattern
	ConditionIncreased EventCondition = "increased" // Value increased
	ConditionDecreased EventCondition = "decreased" // Value decreased
)

// WatcherConfig holds configuration for a watcher
type WatcherConfig struct {
	ID         string           // Unique watcher ID
	Type       WatcherType      // Type of watcher
	Conditions []EventCondition // When to trigger (changed, created, deleted, etc)

	// File/Directory-specific
	FilePath    string // Path to file/directory to watch
	Recursive   bool   // Watch directory recursively
	CheckHash   bool   // Check file hash for changes
	Pattern     string // File pattern to match (*.log, etc)

	// Process-specific
	ProcessName string // Process name to watch
	PID         int    // Specific PID to watch

	// Port/Network-specific
	Port            int    // Port to watch
	Protocol        string // tcp, udp, etc
	RemoteAddr      string // Remote address pattern to watch
	ConnectionState string // ESTABLISHED, LISTEN, etc

	// Service-specific
	ServiceName string // Service name to watch

	// Log-specific
	LogPath       string // Path to log file
	LogPattern    string // Regex pattern to match in logs
	FollowLog     bool   // Follow log file (like tail -f)
	LastPosition  int64  // Last read position in log file

	// Command-specific
	Command       string      // Command to execute
	ExpectedExit  int         // Expected exit code
	OutputPattern string      // Pattern to match in output
	Threshold     float64     // Threshold value for numeric checks

	// Resource monitoring
	CPUThreshold    float64 // CPU percentage threshold
	MemoryThreshold float64 // Memory percentage threshold
	DiskThreshold   float64 // Disk percentage threshold
	NetworkThreshold float64 // Network bytes/sec threshold

	// User-specific
	Username string // Username to watch

	// Package-specific
	PackageName string // Package name to watch

	// Custom check function
	CheckFunc func() (bool, map[string]interface{}) `json:"-"` // Custom check function (not serialized)

	// General settings
	Interval time.Duration // Check interval
	Stack    string        // Stack name context
	RunID    string        // Run ID context
}

// WatcherState holds the current state of a watcher
type WatcherState struct {
	LastCheck    time.Time
	LastHash     string
	LastExists   bool
	LastSize     int64
	LastModTime  time.Time
	CustomState  map[string]interface{}
}

// EventWatcherManager manages all registered watchers
type EventWatcherManager struct {
	eventWorker *EventWorker
	watchers    map[string]*Watcher
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	db          *sql.DB
	dbPath      string
}

// Watcher represents a single event watcher
type Watcher struct {
	config WatcherConfig
	state  WatcherState
	mu     sync.Mutex
}

// NewEventWatcherManager creates a new event watcher manager
func NewEventWatcherManager(eventWorker *EventWorker) *EventWatcherManager {
	ctx, cancel := context.WithCancel(context.Background())

	// Determine database path
	dbPath := os.Getenv("SLOTH_RUNNER_WATCHER_DB")
	if dbPath == "" {
		dbPath = "/var/lib/sloth-runner/watchers.db"
		// Check if we're running as non-root
		if os.Geteuid() != 0 {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				dbPath = filepath.Join(homeDir, ".local", "share", "sloth-runner", "watchers.db")
			}
		}
	}

	mgr := &EventWatcherManager{
		eventWorker: eventWorker,
		watchers:    make(map[string]*Watcher),
		ctx:         ctx,
		cancel:      cancel,
		dbPath:      dbPath,
	}

	// Initialize database
	if err := mgr.initDB(); err != nil {
		slog.Error("Failed to initialize watcher database", "error", err)
		// Continue without persistence
	} else {
		// Load existing watchers from database
		if err := mgr.loadWatchers(); err != nil {
			slog.Error("Failed to load watchers from database", "error", err)
		}
	}

	return mgr
}

// initDB initializes the SQLite database for watchers
func (m *EventWatcherManager) initDB() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(m.dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create watcher db directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite3", m.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open watcher database: %w", err)
	}

	m.db = db

	// Create watchers table
	schema := `
	CREATE TABLE IF NOT EXISTS watchers (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		config_json TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_watchers_type ON watchers(type);
	CREATE INDEX IF NOT EXISTS idx_watchers_created_at ON watchers(created_at);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create watcher schema: %w", err)
	}

	slog.Info("Watcher database initialized", "path", m.dbPath)
	return nil
}

// saveWatcher persists a watcher to the database
func (m *EventWatcherManager) saveWatcher(config *WatcherConfig) error {
	if m.db == nil {
		return nil // No database available
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal watcher config: %w", err)
	}

	now := time.Now().Unix()

	// Insert or replace watcher
	query := `
	INSERT OR REPLACE INTO watchers (id, type, config_json, created_at, updated_at)
	VALUES (?, ?, ?, COALESCE((SELECT created_at FROM watchers WHERE id = ?), ?), ?)
	`

	_, err = m.db.Exec(query, config.ID, string(config.Type), string(configJSON), config.ID, now, now)
	if err != nil {
		return fmt.Errorf("failed to save watcher: %w", err)
	}

	slog.Debug("Watcher saved to database", "id", config.ID)
	return nil
}

// deleteWatcher removes a watcher from the database
func (m *EventWatcherManager) deleteWatcher(id string) error {
	if m.db == nil {
		return nil // No database available
	}

	query := `DELETE FROM watchers WHERE id = ?`
	_, err := m.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete watcher from database: %w", err)
	}

	slog.Debug("Watcher deleted from database", "id", id)
	return nil
}

// loadWatchers loads all watchers from the database
func (m *EventWatcherManager) loadWatchers() error {
	if m.db == nil {
		return nil // No database available
	}

	query := `SELECT config_json FROM watchers ORDER BY created_at ASC`
	rows, err := m.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query watchers: %w", err)
	}
	defer rows.Close()

	loadedCount := 0
	for rows.Next() {
		var configJSON string
		if err := rows.Scan(&configJSON); err != nil {
			slog.Error("Failed to scan watcher row", "error", err)
			continue
		}

		var config WatcherConfig
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			slog.Error("Failed to unmarshal watcher config", "error", err)
			continue
		}

		// Create watcher
		watcher := &Watcher{
			config: config,
			state: WatcherState{
				LastCheck:   time.Now(),
				CustomState: make(map[string]interface{}),
			},
		}

		// Initialize state based on type
		switch config.Type {
		case WatcherTypeFile:
			if err := watcher.initFileState(); err != nil {
				slog.Debug("Failed to initialize file watcher state", "path", config.FilePath, "error", err)
			}
		}

		m.watchers[config.ID] = watcher

		// Start watcher goroutine
		m.wg.Add(1)
		go m.runWatcher(watcher)

		loadedCount++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating watcher rows: %w", err)
	}

	if loadedCount > 0 {
		slog.Info("Watchers loaded from database", "count", loadedCount)
	}

	return nil
}

// RegisterWatcher registers a new watcher
func (m *EventWatcherManager) RegisterWatcher(config WatcherConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Set default interval if not specified
	if config.Interval == 0 {
		config.Interval = 5 * time.Second
	}

	// Create watcher
	watcher := &Watcher{
		config: config,
		state: WatcherState{
			LastCheck:   time.Now(),
			CustomState: make(map[string]interface{}),
		},
	}

	// Initialize state based on type
	switch config.Type {
	case WatcherTypeFile:
		if err := watcher.initFileState(); err != nil {
			slog.Debug("Failed to initialize file watcher state", "path", config.FilePath, "error", err)
		}
	}

	m.watchers[config.ID] = watcher

	// Save to database
	if err := m.saveWatcher(&config); err != nil {
		slog.Error("Failed to save watcher to database", "id", config.ID, "error", err)
		// Continue anyway - watcher is still in memory
	}

	// Start watcher goroutine
	m.wg.Add(1)
	go m.runWatcher(watcher)

	slog.Info("Watcher registered",
		"id", config.ID,
		"type", config.Type,
		"interval", config.Interval)

	return nil
}

// UnregisterWatcher removes a watcher
func (m *EventWatcherManager) UnregisterWatcher(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.watchers, id)
	slog.Info("Watcher unregistered", "id", id)
}

// ListWatchers returns all registered watchers
func (m *EventWatcherManager) ListWatchers() []*WatcherConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	configs := make([]*WatcherConfig, 0, len(m.watchers))
	for _, w := range m.watchers {
		configCopy := w.config
		configs = append(configs, &configCopy)
	}
	return configs
}

// RemoveWatcher removes a watcher and returns an error if it doesn't exist
func (m *EventWatcherManager) RemoveWatcher(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.watchers[id]; !exists {
		return fmt.Errorf("watcher not found: %s", id)
	}

	delete(m.watchers, id)

	// Delete from database
	if err := m.deleteWatcher(id); err != nil {
		slog.Error("Failed to delete watcher from database", "id", id, "error", err)
		// Continue anyway - watcher is removed from memory
	}

	slog.Info("Watcher removed", "id", id)
	return nil
}

// Stop stops all watchers
func (m *EventWatcherManager) Stop() {
	m.cancel()
	m.wg.Wait()

	// Close database
	if m.db != nil {
		if err := m.db.Close(); err != nil {
			slog.Error("Failed to close watcher database", "error", err)
		}
	}

	slog.Info("Event watcher manager stopped")
}

// runWatcher runs a single watcher
func (m *EventWatcherManager) runWatcher(w *Watcher) {
	defer m.wg.Done()

	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			w.check(m.eventWorker)
		}
	}
}

// check performs the watcher check
func (w *Watcher) check(eventWorker *EventWorker) {
	slog.Debug("🔍 Watcher check starting", "id", w.config.ID, "type", w.config.Type, "path", w.config.FilePath)

	if eventWorker == nil {
		slog.Error("❌ EVENT WORKER IS NIL!", "watcher_id", w.config.ID)
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	switch w.config.Type {
	case WatcherTypeFile:
		w.checkFile(eventWorker)
	case WatcherTypeDirectory:
		w.checkDirectory(eventWorker)
	case WatcherTypeProcess:
		w.checkProcess(eventWorker)
	case WatcherTypePort:
		w.checkPort(eventWorker)
	case WatcherTypeService:
		w.checkService(eventWorker)
	case WatcherTypeLog:
		w.checkLog(eventWorker)
	case WatcherTypeCommand:
		w.checkCommand(eventWorker)
	case WatcherTypeCPU:
		w.checkCPU(eventWorker)
	case WatcherTypeMemory:
		w.checkMemory(eventWorker)
	case WatcherTypeDisk:
		w.checkDisk(eventWorker)
	case WatcherTypeNetwork:
		w.checkNetwork(eventWorker)
	case WatcherTypeConnection:
		w.checkConnection(eventWorker)
	case WatcherTypeUser:
		w.checkUser(eventWorker)
	case WatcherTypePackage:
		w.checkPackage(eventWorker)
	case WatcherTypeCustom:
		w.checkCustom(eventWorker)
	}

	w.state.LastCheck = time.Now()
}

// initFileState initializes file watcher state
func (w *Watcher) initFileState() error {
	stat, err := os.Stat(w.config.FilePath)
	if err != nil {
		w.state.LastExists = false
		return err
	}

	w.state.LastExists = true
	w.state.LastSize = stat.Size()
	w.state.LastModTime = stat.ModTime()

	if w.config.CheckHash {
		hash, err := w.calculateFileHash(w.config.FilePath)
		if err == nil {
			w.state.LastHash = hash
		}
	}

	return nil
}

// checkFile checks for file changes
func (w *Watcher) checkFile(eventWorker *EventWorker) {
	slog.Debug("📂 Checking file", "path", w.config.FilePath, "watcher_id", w.config.ID)

	stat, err := os.Stat(w.config.FilePath)
	currentExists := err == nil

	slog.Debug("📊 File stat", "exists", currentExists, "last_exists", w.state.LastExists, "path", w.config.FilePath)

	// Check for deletion
	if w.state.LastExists && !currentExists {
		if w.hasCondition(ConditionDeleted) {
			slog.Info("🗑️ File deleted - sending event", "path", w.config.FilePath)
			eventWorker.SendEvent("file.deleted", w.config.Stack, w.config.RunID, map[string]interface{}{
				"path":       w.config.FilePath,
				"watcher_id": w.config.ID,
			})
		}
		w.state.LastExists = false
		return
	}

	// Check for creation
	if !w.state.LastExists && currentExists {
		if w.hasCondition(ConditionCreated) {
			eventWorker.SendEvent("file.created", w.config.Stack, w.config.RunID, map[string]interface{}{
				"path":       w.config.FilePath,
				"size":       stat.Size(),
				"watcher_id": w.config.ID,
			})
		}
		w.state.LastExists = true
		w.state.LastSize = stat.Size()
		w.state.LastModTime = stat.ModTime()
		return
	}

	// File exists - check for changes
	if currentExists && w.hasCondition(ConditionChanged) {
		changed := false
		changeDetails := map[string]interface{}{
			"path":       w.config.FilePath,
			"watcher_id": w.config.ID,
		}

		slog.Debug("🔎 Checking for changes", "path", w.config.FilePath, "current_size", stat.Size(), "last_size", w.state.LastSize)

		// Check size change
		if stat.Size() != w.state.LastSize {
			changed = true
			changeDetails["old_size"] = w.state.LastSize
			changeDetails["new_size"] = stat.Size()
			slog.Info("📏 Size changed", "path", w.config.FilePath, "old", w.state.LastSize, "new", stat.Size())
			w.state.LastSize = stat.Size()
		}

		// Check modification time
		if !stat.ModTime().Equal(w.state.LastModTime) {
			changed = true
			changeDetails["old_mtime"] = w.state.LastModTime.Unix()
			changeDetails["new_mtime"] = stat.ModTime().Unix()
			slog.Info("⏰ ModTime changed", "path", w.config.FilePath)
			w.state.LastModTime = stat.ModTime()
		}

		// Check hash if enabled
		if w.config.CheckHash {
			hash, err := w.calculateFileHash(w.config.FilePath)
			if err == nil && hash != w.state.LastHash {
				changed = true
				changeDetails["old_hash"] = w.state.LastHash
				changeDetails["new_hash"] = hash
				slog.Info("🔐 Hash changed", "path", w.config.FilePath)
				w.state.LastHash = hash
			}
		}

		if changed {
			slog.Info("📤 FILE CHANGED - Sending event", "path", w.config.FilePath, "details", changeDetails)
			eventWorker.SendEvent("file.modified", w.config.Stack, w.config.RunID, changeDetails)
			slog.Info("✅ Event sent to worker", "event_type", "file.modified")
		} else {
			slog.Debug("⏭️ No changes detected", "path", w.config.FilePath)
		}
	}

	w.state.LastExists = currentExists
}

// checkProcess checks for process events
func (w *Watcher) checkProcess(eventWorker *EventWorker) {
	// Check if process is running
	running := w.isProcessRunning(w.config.ProcessName)

	wasRunning := w.state.CustomState["running"] == true

	if !wasRunning && running && w.hasCondition(ConditionCreated) {
		eventWorker.SendEvent("process.started", w.config.Stack, w.config.RunID, map[string]interface{}{
			"process":    w.config.ProcessName,
			"watcher_id": w.config.ID,
		})
	}

	if wasRunning && !running && w.hasCondition(ConditionDeleted) {
		eventWorker.SendEvent("process.stopped", w.config.Stack, w.config.RunID, map[string]interface{}{
			"process":    w.config.ProcessName,
			"watcher_id": w.config.ID,
		})
	}

	w.state.CustomState["running"] = running
}

// checkPort checks for port events
func (w *Watcher) checkPort(eventWorker *EventWorker) {
	listening := w.isPortListening(w.config.Port)

	wasListening := w.state.CustomState["listening"] == true

	if !wasListening && listening && w.hasCondition(ConditionCreated) {
		eventWorker.SendEvent("port.opened", w.config.Stack, w.config.RunID, map[string]interface{}{
			"port":       w.config.Port,
			"watcher_id": w.config.ID,
		})
	}

	if wasListening && !listening && w.hasCondition(ConditionDeleted) {
		eventWorker.SendEvent("port.closed", w.config.Stack, w.config.RunID, map[string]interface{}{
			"port":       w.config.Port,
			"watcher_id": w.config.ID,
		})
	}

	w.state.CustomState["listening"] = listening
}

// checkService checks for service events
func (w *Watcher) checkService(eventWorker *EventWorker) {
	status := w.getServiceStatus(w.config.ServiceName)

	lastStatus, _ := w.state.CustomState["status"].(string)

	if lastStatus != status && w.hasCondition(ConditionChanged) {
		eventWorker.SendEvent("service.status_changed", w.config.Stack, w.config.RunID, map[string]interface{}{
			"service":     w.config.ServiceName,
			"old_status":  lastStatus,
			"new_status":  status,
			"watcher_id":  w.config.ID,
		})
	}

	w.state.CustomState["status"] = status
}

// checkCustom executes custom check function
func (w *Watcher) checkCustom(eventWorker *EventWorker) {
	if w.config.CheckFunc == nil {
		return
	}

	triggered, data := w.config.CheckFunc()
	if triggered {
		if data == nil {
			data = make(map[string]interface{})
		}
		data["watcher_id"] = w.config.ID

		eventWorker.SendEvent("custom.triggered", w.config.Stack, w.config.RunID, data)
	}
}

// Helper functions

func (w *Watcher) hasCondition(condition EventCondition) bool {
	for _, c := range w.config.Conditions {
		if c == condition {
			return true
		}
	}
	return false
}

func (w *Watcher) calculateFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (w *Watcher) isProcessRunning(name string) bool {
	// Check /proc for process (Linux)
	procDir := "/proc"
	entries, err := os.ReadDir(procDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if directory name is a number (PID)
		pid := entry.Name()
		if len(pid) == 0 || pid[0] < '0' || pid[0] > '9' {
			continue
		}

		// Read cmdline
		cmdlinePath := filepath.Join(procDir, pid, "cmdline")
		data, err := os.ReadFile(cmdlinePath)
		if err != nil {
			continue
		}

		cmdline := string(data)
		if len(cmdline) > 0 && filepath.Base(cmdline) == name {
			return true
		}
	}

	return false
}

func (w *Watcher) isPortListening(port int) bool {
	// Check if port is listening by trying to read /proc/net/tcp
	tcpPath := "/proc/net/tcp"
	data, err := os.ReadFile(tcpPath)
	if err != nil {
		return false
	}

	// Look for port in hex format in the tcp connections
	portHex := fmt.Sprintf(":%04X", port)
	return len(data) > 0 && strings.Contains(string(data), portHex)
}

func (w *Watcher) getServiceStatus(name string) string {
	// Try systemctl status (Linux systemd)
	// This is a placeholder - would need proper implementation
	return "unknown"
}

// checkDirectory checks for directory changes
func (w *Watcher) checkDirectory(eventWorker *EventWorker) {
	stat, err := os.Stat(w.config.FilePath)
	currentExists := err == nil && stat.IsDir()

	// Check for deletion
	if w.state.LastExists && !currentExists {
		if w.hasCondition(ConditionDeleted) {
			eventWorker.SendEvent("dir.deleted", w.config.Stack, w.config.RunID, map[string]interface{}{
				"path":       w.config.FilePath,
				"watcher_id": w.config.ID,
			})
		}
		w.state.LastExists = false
		return
	}

	// Check for creation
	if !w.state.LastExists && currentExists {
		if w.hasCondition(ConditionCreated) {
			eventWorker.SendEvent("dir.created", w.config.Stack, w.config.RunID, map[string]interface{}{
				"path":       w.config.FilePath,
				"watcher_id": w.config.ID,
			})
		}
		w.state.LastExists = true
		return
	}

	// Check for changes (file count, etc)
	if currentExists && w.hasCondition(ConditionChanged) {
		entries, err := os.ReadDir(w.config.FilePath)
		if err == nil {
			currentCount := len(entries)
			lastCount, _ := w.state.CustomState["file_count"].(int)

			if lastCount != 0 && currentCount != lastCount {
				eventWorker.SendEvent("dir.changed", w.config.Stack, w.config.RunID, map[string]interface{}{
					"path":           w.config.FilePath,
					"old_file_count": lastCount,
					"new_file_count": currentCount,
					"watcher_id":     w.config.ID,
				})
			}

			w.state.CustomState["file_count"] = currentCount
		}
	}

	w.state.LastExists = currentExists
}

// checkLog checks for log file patterns
func (w *Watcher) checkLog(eventWorker *EventWorker) {
	// Read log file
	data, err := os.ReadFile(w.config.FilePath)
	if err != nil {
		return
	}

	// Get last position
	lastPos, _ := w.state.CustomState["last_position"].(int64)
	content := string(data)

	// Only check new content
	if int64(len(content)) <= lastPos {
		return
	}

	newContent := content[lastPos:]

	// Check for pattern matches
	if w.config.Pattern != "" && w.hasCondition(ConditionMatches) {
		lines := strings.Split(newContent, "\n")
		for _, line := range lines {
			matched, _ := filepath.Match(w.config.Pattern, line)
			if !matched {
				// Try regex match
				matched = strings.Contains(line, w.config.Pattern)
			}

			if matched {
				eventWorker.SendEvent("log.pattern_matched", w.config.Stack, w.config.RunID, map[string]interface{}{
					"path":       w.config.FilePath,
					"pattern":    w.config.Pattern,
					"line":       line,
					"watcher_id": w.config.ID,
				})
			}
		}
	}

	w.state.CustomState["last_position"] = int64(len(content))
}

// checkCommand checks for command output changes
func (w *Watcher) checkCommand(eventWorker *EventWorker) {
	// Execute command
	parts := strings.Fields(w.config.Command)
	if len(parts) == 0 {
		return
	}

	var output []byte
	var err error

	if len(parts) == 1 {
		cmd := exec.Command(parts[0])
		output, err = cmd.Output()
	} else {
		cmd := exec.Command(parts[0], parts[1:]...)
		output, err = cmd.Output()
	}

	if err != nil {
		return
	}

	currentOutput := string(output)
	lastOutput, _ := w.state.CustomState["last_output"].(string)

	if lastOutput != "" && currentOutput != lastOutput && w.hasCondition(ConditionChanged) {
		eventWorker.SendEvent("command.output_changed", w.config.Stack, w.config.RunID, map[string]interface{}{
			"command":     w.config.Command,
			"old_output":  lastOutput,
			"new_output":  currentOutput,
			"watcher_id":  w.config.ID,
		})
	}

	w.state.CustomState["last_output"] = currentOutput
}

// checkCPU checks for CPU threshold events
func (w *Watcher) checkCPU(eventWorker *EventWorker) {
	// Read CPU stats from /proc/stat
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return
	}

	// Parse first line (overall CPU)
	fields := strings.Fields(lines[0])
	if len(fields) < 5 || fields[0] != "cpu" {
		return
	}

	// Simple CPU usage calculation
	// In production, would need proper delta calculation
	cpuPercent := 50.0 // Placeholder

	if w.hasCondition(ConditionAbove) && cpuPercent > w.config.Threshold {
		eventWorker.SendEvent("cpu.high_usage", w.config.Stack, w.config.RunID, map[string]interface{}{
			"cpu_percent": cpuPercent,
			"threshold":   w.config.Threshold,
			"watcher_id":  w.config.ID,
		})
	}

	if w.hasCondition(ConditionBelow) && cpuPercent < w.config.Threshold {
		eventWorker.SendEvent("cpu.low_usage", w.config.Stack, w.config.RunID, map[string]interface{}{
			"cpu_percent": cpuPercent,
			"threshold":   w.config.Threshold,
			"watcher_id":  w.config.ID,
		})
	}
}

// checkMemory checks for memory threshold events
func (w *Watcher) checkMemory(eventWorker *EventWorker) {
	// Read memory info from /proc/meminfo
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	var memTotal, memAvailable int64
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memTotal, _ = strconv.ParseInt(fields[1], 10, 64)
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memAvailable, _ = strconv.ParseInt(fields[1], 10, 64)
			}
		}
	}

	if memTotal == 0 {
		return
	}

	memUsedPercent := float64(memTotal-memAvailable) / float64(memTotal) * 100.0

	if w.hasCondition(ConditionAbove) && memUsedPercent > w.config.Threshold {
		eventWorker.SendEvent("memory.high_usage", w.config.Stack, w.config.RunID, map[string]interface{}{
			"memory_percent": memUsedPercent,
			"threshold":      w.config.Threshold,
			"watcher_id":     w.config.ID,
		})
	}
}

// checkNetwork checks for network interface events
func (w *Watcher) checkNetwork(eventWorker *EventWorker) {
	// Check network interface status
	// Placeholder implementation
	if w.hasCondition(ConditionChanged) {
		// Would check interface status here
	}
}

// checkConnection checks for network connection events
func (w *Watcher) checkConnection(eventWorker *EventWorker) {
	// Check for specific network connections
	// Placeholder implementation
	if w.hasCondition(ConditionChanged) {
		// Would check connections here
	}
}

// checkUser checks for user session events
func (w *Watcher) checkUser(eventWorker *EventWorker) {
	// Check user sessions (who, w command)
	// Placeholder implementation
	if w.hasCondition(ConditionChanged) {
		// Would check user sessions here
	}
}

// checkPackage checks for package installation events
func (w *Watcher) checkPackage(eventWorker *EventWorker) {
	// Check installed packages
	// Placeholder implementation
	if w.hasCondition(ConditionChanged) {
		// Would check package list here
	}
}
