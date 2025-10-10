//go:build cgo
// +build cgo

package stack

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// WorkflowStateTracker tracks workflow executions in state backend
type WorkflowStateTracker struct {
	backend *StateBackend
}

// NewWorkflowStateTracker creates a new workflow state tracker
func NewWorkflowStateTracker(backend *StateBackend) *WorkflowStateTracker {
	return &WorkflowStateTracker{backend: backend}
}

// TrackWorkflowStart tracks the start of a workflow execution
func (wst *WorkflowStateTracker) TrackWorkflowStart(stackID, workflowName, runID string, metadata map[string]interface{}) error {
	wst.backend.sm.mu.Lock()
	defer wst.backend.sm.mu.Unlock()

	// Create snapshot before workflow execution
	_, err := wst.backend.CreateSnapshot(stackID, "system", fmt.Sprintf("Pre-execution snapshot for %s", workflowName))
	if err != nil {
		return fmt.Errorf("failed to create pre-execution snapshot: %w", err)
	}

	// Log activity
	wst.backend.logActivity(stackID, "workflow_started", "", fmt.Sprintf("Workflow %s started (run: %s)", workflowName, runID), "system")

	return nil
}

// TrackWorkflowEnd tracks the end of a workflow execution
func (wst *WorkflowStateTracker) TrackWorkflowEnd(stackID, workflowName, runID string, success bool, duration time.Duration, outputs map[string]interface{}) error {
	// Create post-execution snapshot
	status := "completed"
	if !success {
		status = "failed"
	}

	_, err := wst.backend.CreateSnapshot(stackID, "system", fmt.Sprintf("Post-execution snapshot: %s (%s)", workflowName, status))
	if err != nil {
		return fmt.Errorf("failed to create post-execution snapshot: %w", err)
	}

	// Update stack outputs
	stack, err := wst.backend.sm.GetStack(stackID)
	if err != nil {
		return err
	}

	stack.Outputs = outputs
	stack.LastDuration = duration
	stack.ExecutionCount++

	if err := wst.backend.sm.UpdateStack(stack); err != nil {
		return err
	}

	// Log activity
	wst.backend.sm.mu.Lock()
	wst.backend.logActivity(stackID, "workflow_completed", "", fmt.Sprintf("Workflow %s %s (duration: %v)", workflowName, status, duration), "system")
	wst.backend.sm.mu.Unlock()

	return nil
}

// StateEncryption handles state encryption
type StateEncryption struct {
	key []byte
}

// NewStateEncryption creates a new state encryption handler
func NewStateEncryption(password, salt string) *StateEncryption {
	key := pbkdf2.Key([]byte(password), []byte(salt), 100000, 32, sha256.New)
	return &StateEncryption{key: key}
}

// Encrypt encrypts data
func (se *StateEncryption) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(se.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts data
func (se *StateEncryption) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(se.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// StateCompression handles state compression
type StateCompression struct{}

// NewStateCompression creates a new state compression handler
func NewStateCompression() *StateCompression {
	return &StateCompression{}
}

// Compress compresses data using gzip
func (sc *StateCompression) Compress(data []byte) ([]byte, error) {
	var buf strings.Builder
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

// Decompress decompresses gzip data
func (sc *StateCompression) Decompress(data []byte) ([]byte, error) {
	reader := strings.NewReader(string(data))
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	decompressed, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

// ResourceGraph represents a dependency graph of resources
type ResourceGraph struct {
	backend *StateBackend
}

// NewResourceGraph creates a new resource graph
func NewResourceGraph(backend *StateBackend) *ResourceGraph {
	return &ResourceGraph{backend: backend}
}

// GraphNode represents a node in the resource graph
type GraphNode struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	State        string   `json:"state"`
	Dependencies []string `json:"dependencies"`
	Dependents   []string `json:"dependents"`
}

// BuildGraph builds a resource dependency graph for a stack
func (rg *ResourceGraph) BuildGraph(stackID string) (map[string]*GraphNode, error) {
	resources, err := rg.backend.sm.ListResources(stackID)
	if err != nil {
		return nil, err
	}

	graph := make(map[string]*GraphNode)

	// Create nodes
	for _, res := range resources {
		node := &GraphNode{
			ID:           res.ID,
			Type:         res.Type,
			Name:         res.Name,
			State:        res.State,
			Dependencies: []string{},
			Dependents:   []string{},
		}

		// Get dependencies from database
		deps, err := rg.backend.GetResourceDependencies(res.ID)
		if err == nil {
			node.Dependencies = deps
		}

		graph[res.ID] = node
	}

	// Populate dependents (reverse dependencies)
	for _, node := range graph {
		for _, depID := range node.Dependencies {
			if depNode, exists := graph[depID]; exists {
				depNode.Dependents = append(depNode.Dependents, node.ID)
			}
		}
	}

	return graph, nil
}

// TopologicalSort returns resources in dependency order
func (rg *ResourceGraph) TopologicalSort(stackID string) ([]*Resource, error) {
	graph, err := rg.BuildGraph(stackID)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	tempMark := make(map[string]bool)
	var sorted []string

	var visit func(string) error
	visit = func(nodeID string) error {
		if tempMark[nodeID] {
			return fmt.Errorf("circular dependency detected at %s", nodeID)
		}
		if visited[nodeID] {
			return nil
		}

		tempMark[nodeID] = true

		node := graph[nodeID]
		for _, depID := range node.Dependencies {
			if err := visit(depID); err != nil {
				return err
			}
		}

		tempMark[nodeID] = false
		visited[nodeID] = true
		sorted = append(sorted, nodeID)

		return nil
	}

	// Visit all nodes
	for nodeID := range graph {
		if !visited[nodeID] {
			if err := visit(nodeID); err != nil {
				return nil, err
			}
		}
	}

	// Get resources in sorted order
	resources := make([]*Resource, 0, len(sorted))
	for _, id := range sorted {
		res, err := rg.backend.sm.GetResource(id)
		if err == nil {
			resources = append(resources, res)
		}
	}

	return resources, nil
}

// AutoRemediation handles automatic drift remediation
type AutoRemediation struct {
	backend *StateBackend
}

// NewAutoRemediation creates a new auto-remediation handler
func NewAutoRemediation(backend *StateBackend) *AutoRemediation {
	return &AutoRemediation{backend: backend}
}

// RemediationStrategy represents how to remediate drift
type RemediationStrategy string

const (
	RemediationNone    RemediationStrategy = "none"
	RemediationNotify  RemediationStrategy = "notify"
	RemediationAutoFix RemediationStrategy = "auto_fix"
	RemediationRollback RemediationStrategy = "rollback"
)

// RemediationConfig configures auto-remediation behavior
type RemediationConfig struct {
	Strategy   RemediationStrategy
	MaxRetries int
	Webhooks   []string
}

// CheckAndRemediate checks for drift and applies remediation
func (ar *AutoRemediation) CheckAndRemediate(stackID string, config RemediationConfig) error {
	// Get drift information
	drifts, err := ar.backend.GetDriftInfo(stackID)
	if err != nil {
		return err
	}

	hasDrift := false
	for _, drift := range drifts {
		if drift.IsDrifted && drift.ResolutionStatus == "pending" {
			hasDrift = true
			break
		}
	}

	if !hasDrift {
		return nil
	}

	switch config.Strategy {
	case RemediationNotify:
		return ar.notifyDrift(stackID, drifts, config)
	case RemediationAutoFix:
		return ar.autoFixDrift(stackID, drifts, config)
	case RemediationRollback:
		return ar.rollbackDrift(stackID, config)
	default:
		return nil
	}
}

// notifyDrift sends notifications about drift
func (ar *AutoRemediation) notifyDrift(stackID string, drifts []*DriftInfo, config RemediationConfig) error {
	for _, webhook := range config.Webhooks {
		payload := map[string]interface{}{
			"event":    "drift_detected",
			"stack_id": stackID,
			"drifts":   drifts,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		if err := sendWebhook(webhook, payload); err != nil {
			// Log error but continue
			ar.backend.sm.mu.Lock()
			ar.backend.logActivity(stackID, "webhook_failed", "", fmt.Sprintf("Failed to send webhook: %v", err), "system")
			ar.backend.sm.mu.Unlock()
		}
	}

	return nil
}

// autoFixDrift attempts to automatically fix drift
func (ar *AutoRemediation) autoFixDrift(stackID string, drifts []*DriftInfo, config RemediationConfig) error {
	// Log remediation attempt
	ar.backend.sm.mu.Lock()
	ar.backend.logActivity(stackID, "auto_remediation", "", "Attempting automatic drift remediation", "system")
	ar.backend.sm.mu.Unlock()

	// For each drifted resource, attempt to reapply expected state
	for _, drift := range drifts {
		if !drift.IsDrifted || drift.ResolutionStatus != "pending" {
			continue
		}

		resource, err := ar.backend.sm.GetResource(drift.ResourceID)
		if err != nil {
			continue
		}

		// Update resource properties to expected state
		resource.Properties = drift.ExpectedState
		resource.State = "pending_remediation"

		if err := ar.backend.sm.UpdateResource(resource); err != nil {
			continue
		}

		// Mark drift as resolved
		ar.backend.sm.mu.Lock()
		ar.backend.sm.db.Exec(`
			UPDATE drift_detections
			SET resolution_status = 'resolved'
			WHERE id = ?
		`, drift.ID)
		ar.backend.sm.mu.Unlock()
	}

	// Log completion
	ar.backend.sm.mu.Lock()
	ar.backend.logActivity(stackID, "auto_remediation", "", fmt.Sprintf("Remediated %d drifted resources", len(drifts)), "system")
	ar.backend.sm.mu.Unlock()

	// Notify
	return ar.notifyDrift(stackID, drifts, config)
}

// rollbackDrift rolls back to last known good state
func (ar *AutoRemediation) rollbackDrift(stackID string, config RemediationConfig) error {
	// Get latest successful snapshot
	snapshots, err := ar.backend.ListSnapshots(stackID)
	if err != nil {
		return err
	}

	if len(snapshots) < 2 {
		return fmt.Errorf("no previous snapshots available for rollback")
	}

	// Find last successful snapshot
	var rollbackVersion int
	for _, snap := range snapshots {
		if snap.StackState.Status == "completed" {
			rollbackVersion = snap.Version
			break
		}
	}

	if rollbackVersion == 0 {
		return fmt.Errorf("no successful snapshot found for rollback")
	}

	// Perform rollback
	return ar.backend.RollbackToSnapshot(stackID, rollbackVersion, "auto-remediation")
}

// sendWebhook sends a webhook notification
func sendWebhook(url string, payload map[string]interface{}) error {
	// This would be implemented with actual HTTP client
	// For now, just log
	payloadJSON, _ := json.Marshal(payload)
	fmt.Printf("Webhook to %s: %s\n", url, string(payloadJSON))
	return nil
}

// StateChecksum calculates checksums for state verification
type StateChecksum struct{}

// NewStateChecksum creates a new state checksum calculator
func NewStateChecksum() *StateChecksum {
	return &StateChecksum{}
}

// Calculate calculates a SHA-256 checksum
func (sc *StateChecksum) Calculate(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Verify verifies a checksum
func (sc *StateChecksum) Verify(data []byte, expected string) bool {
	actual := sc.Calculate(data)
	return actual == expected
}

// StateBackendConfig represents advanced configuration
type StateBackendConfig struct {
	EncryptionPassword string
	EncryptionSalt     string
	CompressSnapshots  bool
	AutoRemediation    RemediationConfig
	RetentionDays      int
	MaxVersions        int
}

// ApplyConfig applies advanced configuration to state backend
func (sb *StateBackend) ApplyConfig(config StateBackendConfig) error {
	// Store config in metadata table
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = sb.sm.db.Exec(`
		CREATE TABLE IF NOT EXISTS state_config (
			key TEXT PRIMARY KEY,
			value TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	_, err = sb.sm.db.Exec(`
		INSERT OR REPLACE INTO state_config (key, value)
		VALUES ('backend_config', ?)
	`, string(configJSON))

	return err
}

// GetConfig retrieves advanced configuration
func (sb *StateBackend) GetConfig() (*StateBackendConfig, error) {
	sb.sm.mu.RLock()
	defer sb.sm.mu.RUnlock()

	var configJSON string
	err := sb.sm.db.QueryRow(`
		SELECT value FROM state_config WHERE key = 'backend_config'
	`).Scan(&configJSON)

	if err == sql.ErrNoRows {
		// Return default config
		return &StateBackendConfig{
			CompressSnapshots: false,
			RetentionDays:     30,
			MaxVersions:       10,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	var config StateBackendConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// PruneOldSnapshots removes old snapshots based on retention policy
func (sb *StateBackend) PruneOldSnapshots(stackID string, retentionDays int, maxVersions int) (int, error) {
	sb.sm.mu.Lock()
	defer sb.sm.mu.Unlock()

	// Delete snapshots older than retention days
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	result, err := sb.sm.db.Exec(`
		DELETE FROM state_versions
		WHERE stack_id = ? AND created_at < ?
	`, stackID, cutoffDate)

	if err != nil {
		return 0, err
	}

	deleted, _ := result.RowsAffected()

	// Keep only maxVersions latest snapshots
	if maxVersions > 0 {
		_, err = sb.sm.db.Exec(`
			DELETE FROM state_versions
			WHERE stack_id = ? AND version NOT IN (
				SELECT version FROM state_versions
				WHERE stack_id = ?
				ORDER BY version DESC
				LIMIT ?
			)
		`, stackID, stackID, maxVersions)
	}

	return int(deleted), nil
}
