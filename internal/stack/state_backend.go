//go:build cgo
// +build cgo

package stack

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// StateBackend provides Pulumi/Terraform-like state management
// This integrates stack management with workflow state tracking
type StateBackend struct {
	sm *StackManager
}

// NewStateBackend creates a new state backend
func NewStateBackend(dbPath string) (*StateBackend, error) {
	sm, err := NewStackManager(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create stack manager: %w", err)
	}

	backend := &StateBackend{sm: sm}

	// Initialize extended schema for state management
	if err := backend.initExtendedSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize extended schema: %w", err)
	}

	return backend, nil
}

// initExtendedSchema adds state management tables
func (sb *StateBackend) initExtendedSchema() error {
	schema := `
	-- State Versions Table (for rollback support)
	CREATE TABLE IF NOT EXISTS state_versions (
		id TEXT PRIMARY KEY,
		stack_id TEXT NOT NULL,
		version INTEGER NOT NULL,
		state_snapshot TEXT NOT NULL, -- JSON snapshot of entire stack state
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		description TEXT,
		checksum TEXT,
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE,
		UNIQUE(stack_id, version)
	);

	-- Drift Detection Table
	CREATE TABLE IF NOT EXISTS drift_detections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stack_id TEXT NOT NULL,
		resource_id TEXT NOT NULL,
		detected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expected_state TEXT, -- JSON
		actual_state TEXT, -- JSON
		drift_fields TEXT, -- JSON array of drifted field names
		is_drifted BOOLEAN DEFAULT 0,
		resolution_status TEXT DEFAULT 'pending', -- pending, resolved, ignored
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE,
		FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
	);

	-- State Locks Table (for concurrent access control)
	CREATE TABLE IF NOT EXISTS state_locks (
		stack_id TEXT PRIMARY KEY,
		lock_id TEXT NOT NULL,
		operation TEXT NOT NULL,
		who TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL,
		info TEXT, -- Additional lock information
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);

	-- State Backups Table
	CREATE TABLE IF NOT EXISTS state_backups (
		id TEXT PRIMARY KEY,
		stack_id TEXT NOT NULL,
		version INTEGER NOT NULL,
		backup_path TEXT NOT NULL,
		backup_type TEXT, -- manual, automatic, pre-operation
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		size_bytes INTEGER,
		compressed BOOLEAN DEFAULT 1,
		description TEXT,
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);

	-- Stack Tags Table (for organization)
	CREATE TABLE IF NOT EXISTS stack_tags (
		stack_id TEXT NOT NULL,
		tag TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (stack_id, tag),
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);

	-- Resource Dependencies Table (explicit dependency tracking)
	CREATE TABLE IF NOT EXISTS resource_dependencies (
		resource_id TEXT NOT NULL,
		depends_on_id TEXT NOT NULL,
		dependency_type TEXT DEFAULT 'explicit', -- explicit, implicit
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (resource_id, depends_on_id),
		FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
		FOREIGN KEY (depends_on_id) REFERENCES resources(id) ON DELETE CASCADE
	);

	-- State Activity Log (audit trail)
	CREATE TABLE IF NOT EXISTS state_activity (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stack_id TEXT NOT NULL,
		activity_type TEXT NOT NULL, -- create, update, delete, lock, unlock, backup, restore, rollback
		resource_id TEXT,
		details TEXT, -- JSON
		user TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);

	-- State Events Table (for event persistence)
	CREATE TABLE IF NOT EXISTS state_events (
		id TEXT PRIMARY KEY,
		event_type TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		source TEXT NOT NULL,
		stack_id TEXT,
		stack_name TEXT,
		data TEXT, -- JSON
		severity TEXT NOT NULL,
		FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_state_versions_stack ON state_versions(stack_id);
	CREATE INDEX IF NOT EXISTS idx_state_versions_version ON state_versions(stack_id, version);
	CREATE INDEX IF NOT EXISTS idx_drift_stack ON drift_detections(stack_id);
	CREATE INDEX IF NOT EXISTS idx_drift_resource ON drift_detections(resource_id);
	CREATE INDEX IF NOT EXISTS idx_drift_status ON drift_detections(resolution_status);
	CREATE INDEX IF NOT EXISTS idx_state_backups_stack ON state_backups(stack_id);
	CREATE INDEX IF NOT EXISTS idx_stack_tags_tag ON stack_tags(tag);
	CREATE INDEX IF NOT EXISTS idx_resource_deps_resource ON resource_dependencies(resource_id);
	CREATE INDEX IF NOT EXISTS idx_resource_deps_depends ON resource_dependencies(depends_on_id);
	CREATE INDEX IF NOT EXISTS idx_state_activity_stack ON state_activity(stack_id);
	CREATE INDEX IF NOT EXISTS idx_state_activity_time ON state_activity(created_at);
	CREATE INDEX IF NOT EXISTS idx_state_events_stack ON state_events(stack_id);
	CREATE INDEX IF NOT EXISTS idx_state_events_type ON state_events(event_type);
	CREATE INDEX IF NOT EXISTS idx_state_events_time ON state_events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_state_events_severity ON state_events(severity);
	`

	_, err := sb.sm.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create extended schema: %w", err)
	}

	return nil
}

// StateSnapshot represents a point-in-time snapshot of stack state
type StateSnapshot struct {
	Version       int                    `json:"version"`
	StackState    *StackState            `json:"stack_state"`
	Resources     []*Resource            `json:"resources"`
	Checksum      string                 `json:"checksum"`
	CreatedAt     time.Time              `json:"created_at"`
	CreatedBy     string                 `json:"created_by"`
	Description   string                 `json:"description"`
}

// DriftInfo represents drift detection information
type DriftInfo struct {
	ID               int64                  `json:"id"`
	StackID          string                 `json:"stack_id"`
	ResourceID       string                 `json:"resource_id"`
	DetectedAt       time.Time              `json:"detected_at"`
	ExpectedState    map[string]interface{} `json:"expected_state"`
	ActualState      map[string]interface{} `json:"actual_state"`
	DriftedFields    []string               `json:"drifted_fields"`
	IsDrifted        bool                   `json:"is_drifted"`
	ResolutionStatus string                 `json:"resolution_status"`
}

// StateLock represents a state lock
type StateLock struct {
	StackID   string    `json:"stack_id"`
	LockID    string    `json:"lock_id"`
	Operation string    `json:"operation"`
	Who       string    `json:"who"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Info      string    `json:"info"`
}

// CreateSnapshot creates a new state snapshot (version)
func (sb *StateBackend) CreateSnapshot(stackID, createdBy, description string) (int, error) {
	// Get current stack state (this method handles its own locking)
	stack, err := sb.sm.GetStack(stackID)
	if err != nil {
		return 0, fmt.Errorf("failed to get stack: %w", err)
	}

	// Get all resources (this method handles its own locking)
	resources, err := sb.sm.ListResources(stackID)
	if err != nil {
		return 0, fmt.Errorf("failed to get resources: %w", err)
	}

	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	// Get next version number
	var version int
	err = sb.sm.db.QueryRow(`
		SELECT COALESCE(MAX(version), 0) + 1
		FROM state_versions
		WHERE stack_id = ?
	`, stackID).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("failed to get next version: %w", err)
	}

	// Create snapshot
	snapshot := StateSnapshot{
		Version:       version,
		StackState:    stack,
		Resources:     resources,
		CreatedAt:     time.Now(),
		CreatedBy:     createdBy,
		Description:   description,
	}

	// Serialize snapshot
	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// Calculate checksum
	checksum := calculateChecksum(snapshotJSON)

	// Insert version
	versionID := fmt.Sprintf("%s-v%d", stackID, version)
	_, err = sb.sm.db.Exec(`
		INSERT INTO state_versions (id, stack_id, version, state_snapshot, created_by, description, checksum)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, versionID, stackID, version, string(snapshotJSON), createdBy, description, checksum)

	if err != nil {
		return 0, fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Log activity
	sb.logActivity(stackID, "snapshot", "", fmt.Sprintf("Created version %d", version), createdBy)

	return version, nil
}

// GetSnapshot retrieves a specific snapshot
func (sb *StateBackend) GetSnapshot(stackID string, version int) (*StateSnapshot, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	var snapshotJSON string
	err := sb.sm.db.QueryRow(`
		SELECT state_snapshot FROM state_versions
		WHERE stack_id = ? AND version = ?
	`, stackID, version).Scan(&snapshotJSON)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("snapshot not found: stack=%s version=%d", stackID, version)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	var snapshot StateSnapshot
	if err := json.Unmarshal([]byte(snapshotJSON), &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &snapshot, nil
}

// ListSnapshots lists all snapshots for a stack
func (sb *StateBackend) ListSnapshots(stackID string) ([]StateSnapshot, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	rows, err := sb.sm.db.Query(`
		SELECT state_snapshot, created_at, created_by, description
		FROM state_versions
		WHERE stack_id = ?
		ORDER BY version DESC
	`, stackID)

	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}
	defer rows.Close()

	var snapshots []StateSnapshot
	for rows.Next() {
		var snapshotJSON, createdBy, description string
		var createdAt time.Time

		if err := rows.Scan(&snapshotJSON, &createdAt, &createdBy, &description); err != nil {
			return nil, err
		}

		var snapshot StateSnapshot
		if err := json.Unmarshal([]byte(snapshotJSON), &snapshot); err != nil {
			continue // Skip invalid snapshots
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

// RollbackToSnapshot rolls back stack to a specific version
func (sb *StateBackend) RollbackToSnapshot(stackID string, version int, performedBy string) error {
	// Get the snapshot (this method handles its own locking)
	snapshot, err := sb.GetSnapshot(stackID, version)
	if err != nil {
		return err
	}

	// Create a backup of current state before rollback (handles its own locking)
	_, backupErr := sb.CreateSnapshot(stackID, "system", fmt.Sprintf("Pre-rollback backup from v%d", version))
	if backupErr != nil {
		return fmt.Errorf("failed to create pre-rollback backup: %w", backupErr)
	}

	// Update stack state (handles its own locking)
	snapshot.StackState.Status = "rolled_back"
	snapshot.StackState.Metadata["rollback_from_version"] = fmt.Sprintf("%d", snapshot.StackState.ExecutionCount)
	snapshot.StackState.Metadata["rollback_to_version"] = fmt.Sprintf("%d", version)

	if err := sb.sm.UpdateStack(snapshot.StackState); err != nil {
		return fmt.Errorf("failed to update stack: %w", err)
	}

	// Delete and recreate resources (simple approach)
	// In production, you'd want more sophisticated reconciliation
	currentResources, err := sb.sm.ListResources(stackID)
	if err != nil {
		return fmt.Errorf("failed to list current resources: %w", err)
	}

	// Delete current resources (each handles its own locking)
	for _, res := range currentResources {
		sb.sm.DeleteResource(res.ID)
	}

	// Recreate resources from snapshot (each handles its own locking)
	for _, res := range snapshot.Resources {
		if err := sb.sm.CreateResource(res); err != nil {
			return fmt.Errorf("failed to recreate resource: %w", err)
		}
	}

	// Log activity (lock is acquired inside logActivity call to db.Exec)
	sb.sm.mu.Lock()
	sb.logActivity(stackID, "rollback", "", fmt.Sprintf("Rolled back to version %d", version), performedBy)
	sb.sm.mu.Unlock()

	return nil
}

// DetectDrift checks for drift between expected and actual state
func (sb *StateBackend) DetectDrift(stackID, resourceID string, expectedState, actualState map[string]interface{}) error {
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	// Compare states and find drifted fields
	driftedFields := []string{}
	isDrifted := false

	for key, expectedVal := range expectedState {
		actualVal, exists := actualState[key]
		if !exists || fmt.Sprintf("%v", expectedVal) != fmt.Sprintf("%v", actualVal) {
			driftedFields = append(driftedFields, key)
			isDrifted = true
		}
	}

	// Check for extra fields in actual
	for key := range actualState {
		if _, exists := expectedState[key]; !exists {
			driftedFields = append(driftedFields, key)
			isDrifted = true
		}
	}

	expectedJSON, _ := json.Marshal(expectedState)
	actualJSON, _ := json.Marshal(actualState)
	driftFieldsJSON, _ := json.Marshal(driftedFields)

	_, err := sb.sm.db.Exec(`
		INSERT INTO drift_detections (
			stack_id, resource_id, expected_state, actual_state,
			drift_fields, is_drifted
		) VALUES (?, ?, ?, ?, ?, ?)
	`, stackID, resourceID, string(expectedJSON), string(actualJSON),
	   string(driftFieldsJSON), isDrifted)

	if err != nil {
		return fmt.Errorf("failed to record drift: %w", err)
	}

	if isDrifted {
		sb.logActivity(stackID, "drift_detected", resourceID,
			fmt.Sprintf("Drift detected in %d fields", len(driftedFields)), "system")
	}

	return nil
}

// GetDriftInfo retrieves drift information for a stack
func (sb *StateBackend) GetDriftInfo(stackID string) ([]*DriftInfo, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	rows, err := sb.sm.db.Query(`
		SELECT id, stack_id, resource_id, detected_at, expected_state,
		       actual_state, drift_fields, is_drifted, resolution_status
		FROM drift_detections
		WHERE stack_id = ?
		ORDER BY detected_at DESC
		LIMIT 100
	`, stackID)

	if err != nil {
		return nil, fmt.Errorf("failed to get drift info: %w", err)
	}
	defer rows.Close()

	var drifts []*DriftInfo
	for rows.Next() {
		var drift DriftInfo
		var expectedJSON, actualJSON, driftFieldsJSON string

		err := rows.Scan(&drift.ID, &drift.StackID, &drift.ResourceID, &drift.DetectedAt,
			&expectedJSON, &actualJSON, &driftFieldsJSON, &drift.IsDrifted, &drift.ResolutionStatus)

		if err != nil {
			continue
		}

		json.Unmarshal([]byte(expectedJSON), &drift.ExpectedState)
		json.Unmarshal([]byte(actualJSON), &drift.ActualState)
		json.Unmarshal([]byte(driftFieldsJSON), &drift.DriftedFields)

		drifts = append(drifts, &drift)
	}

	return drifts, nil
}

// LockState acquires a lock on the state
func (sb *StateBackend) LockState(stackID, lockID, operation, who string, duration time.Duration) error {
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	// Check if already locked (use UTC time comparison)
	var existingWho string
	var expiresAt time.Time
	err := sb.sm.db.QueryRow(`
		SELECT who, expires_at FROM state_locks
		WHERE stack_id = ?
	`, stackID).Scan(&existingWho, &expiresAt)

	if err == nil {
		// Lock exists, check if it's still valid
		if time.Now().Before(expiresAt) {
			return fmt.Errorf("state is already locked by %s", existingWho)
		}
		// Lock is expired, we can replace it
	}

	expiresAtTime := time.Now().Add(duration)

	// Use REPLACE to handle both insert and update
	_, err = sb.sm.db.Exec(`
		REPLACE INTO state_locks (stack_id, lock_id, operation, who, created_at, expires_at, info)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, stackID, lockID, operation, who, time.Now(), expiresAtTime, "")

	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	sb.logActivity(stackID, "lock", "", fmt.Sprintf("State locked for %s", operation), who)

	return nil
}

// UnlockState releases a state lock
func (sb *StateBackend) UnlockState(stackID, lockID string) error {
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	result, err := sb.sm.db.Exec(`
		DELETE FROM state_locks
		WHERE stack_id = ? AND lock_id = ?
	`, stackID, lockID)

	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("lock not found or already released")
	}

	sb.logActivity(stackID, "unlock", "", "State unlocked", "system")

	return nil
}

// GetLockInfo retrieves lock information for a stack
func (sb *StateBackend) GetLockInfo(stackID string) (*StateLock, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	var lock StateLock
	err := sb.sm.db.QueryRow(`
		SELECT stack_id, lock_id, operation, who, created_at, expires_at, COALESCE(info, '') as info
		FROM state_locks
		WHERE stack_id = ?
	`, stackID).Scan(&lock.StackID, &lock.LockID, &lock.Operation, &lock.Who, &lock.CreatedAt, &lock.ExpiresAt, &lock.Info)

	if err == sql.ErrNoRows {
		return nil, nil // No lock found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lock info: %w", err)
	}

	// Check if lock is still valid (not expired)
	if time.Now().After(lock.ExpiresAt) {
		return nil, nil // Lock is expired
	}

	return &lock, nil
}

// AddTag adds a tag to a stack
func (sb *StateBackend) AddTag(stackID, tag string) error {
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	_, err := sb.sm.db.Exec(`
		INSERT OR IGNORE INTO stack_tags (stack_id, tag)
		VALUES (?, ?)
	`, stackID, tag)

	return err
}

// GetTags retrieves tags for a stack
func (sb *StateBackend) GetTags(stackID string) ([]string, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	rows, err := sb.sm.db.Query(`
		SELECT tag FROM stack_tags
		WHERE stack_id = ?
		ORDER BY tag
	`, stackID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// RemoveTag removes a tag from a stack
func (sb *StateBackend) RemoveTag(stackID, tag string) error {
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	result, err := sb.sm.db.Exec(`
		DELETE FROM stack_tags
		WHERE stack_id = ? AND tag = ?
	`, stackID, tag)

	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("tag not found: %s", tag)
	}

	return nil
}

// AddResourceDependency records a dependency between resources
func (sb *StateBackend) AddResourceDependency(resourceID, dependsOnID, depType string) error {
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	_, err := sb.sm.db.Exec(`
		INSERT OR IGNORE INTO resource_dependencies (resource_id, depends_on_id, dependency_type)
		VALUES (?, ?, ?)
	`, resourceID, dependsOnID, depType)

	return err
}

// GetResourceDependencies retrieves dependencies for a resource
func (sb *StateBackend) GetResourceDependencies(resourceID string) ([]string, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	rows, err := sb.sm.db.Query(`
		SELECT depends_on_id FROM resource_dependencies
		WHERE resource_id = ?
	`, resourceID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []string
	for rows.Next() {
		var depID string
		if err := rows.Scan(&depID); err != nil {
			return nil, err
		}
		deps = append(deps, depID)
	}

	return deps, nil
}

// logActivity logs an activity (internal helper)
func (sb *StateBackend) logActivity(stackID, activityType, resourceID, details, user string) {
	detailsJSON, _ := json.Marshal(map[string]string{"message": details})

	sb.sm.db.Exec(`
		INSERT INTO state_activity (stack_id, activity_type, resource_id, details, user)
		VALUES (?, ?, ?, ?, ?)
	`, stackID, activityType, resourceID, string(detailsJSON), user)
}

// GetActivity retrieves activity log for a stack
func (sb *StateBackend) GetActivity(stackID string, limit int) ([]map[string]interface{}, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	query := `
		SELECT activity_type, resource_id, details, user, created_at
		FROM state_activity
		WHERE stack_id = ?
		ORDER BY created_at DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := sb.sm.db.Query(query, stackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []map[string]interface{}
	for rows.Next() {
		var activityType, resourceID, details, user string
		var createdAt time.Time

		if err := rows.Scan(&activityType, &resourceID, &details, &user, &createdAt); err != nil {
			continue
		}

		activity := map[string]interface{}{
			"type":        activityType,
			"resource_id": resourceID,
			"details":     details,
			"user":        user,
			"created_at":  createdAt,
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// GetStackManager returns the underlying stack manager
func (sb *StateBackend) GetStackManager() *StackManager {
	return sb.sm
}

// Close closes the state backend
func (sb *StateBackend) Close() error {
	return sb.sm.Close()
}

// calculateChecksum calculates a checksum for state verification
func calculateChecksum(data []byte) string {
	// Simple checksum for now - in production use SHA256
	sum := 0
	for _, b := range data {
		sum += int(b)
	}
	return fmt.Sprintf("%x", sum)
}
