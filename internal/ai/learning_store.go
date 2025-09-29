package ai

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/state"
)

// LearningStore manages persistent storage for AI learning data
type LearningStore struct {
	stateAdapter *StateAdapter
}

// NewLearningStore creates a new learning store
func NewLearningStore(stateManager *state.StateManager) *LearningStore {
	return &LearningStore{
		stateAdapter: NewStateAdapter(stateManager),
	}
}

// StoreExecution stores a task execution record
func (ls *LearningStore) StoreExecution(ctx context.Context, execution *TaskExecution) error {
	// Generate unique key for this execution
	key := fmt.Sprintf("ai:execution:%s:%d", ls.hashCommand(execution.Command), execution.Timestamp.Unix())
	
	// Serialize execution data
	data, err := json.Marshal(execution)
	if err != nil {
		return fmt.Errorf("failed to marshal execution: %w", err)
	}

	// Store in state manager with TTL (keep for 30 days)
	if err := ls.stateAdapter.Set(ctx, key, string(data), 30*24*time.Hour); err != nil {
		return fmt.Errorf("failed to store execution: %w", err)
	}

	// Update command index
	if err := ls.updateCommandIndex(ctx, execution.Command, key); err != nil {
		return fmt.Errorf("failed to update command index: %w", err)
	}

	// Update task name index
	if err := ls.updateTaskIndex(ctx, execution.TaskName, key); err != nil {
		return fmt.Errorf("failed to update task index: %w", err)
	}

	return nil
}

// GetTaskHistory retrieves execution history for a command
func (ls *LearningStore) GetTaskHistory(ctx context.Context, command string) ([]*TaskExecution, error) {
	commandHash := ls.hashCommand(command)
	indexKey := fmt.Sprintf("ai:index:command:%s", commandHash)
	
	// Get execution keys from index
	indexData, err := ls.stateAdapter.Get(ctx, indexKey)
	if err != nil {
		return []*TaskExecution{}, nil // Return empty if no history
	}

	var executionKeys []string
	if err := json.Unmarshal([]byte(indexData), &executionKeys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index: %w", err)
	}

	// Retrieve executions
	executions := make([]*TaskExecution, 0, len(executionKeys))
	for _, key := range executionKeys {
		executionData, err := ls.stateAdapter.Get(ctx, key)
		if err != nil {
			continue // Skip missing executions
		}

		var execution TaskExecution
		if err := json.Unmarshal([]byte(executionData), &execution); err != nil {
			continue // Skip malformed executions
		}

		executions = append(executions, &execution)
	}

	// Sort by timestamp (newest first)
	sort.Slice(executions, func(i, j int) bool {
		return executions[i].Timestamp.After(executions[j].Timestamp)
	})

	return executions, nil
}

// FindSimilarTasks finds tasks with similar commands
func (ls *LearningStore) FindSimilarTasks(ctx context.Context, command string, limit int) ([]*TaskExecution, error) {
	// Get all command indices
	pattern := "ai:index:command:*"
	keys, err := ls.stateAdapter.Keys(ctx, pattern)
	if err != nil {
		return []*TaskExecution{}, nil
	}

	// Score commands by similarity
	similarities := make([]commandSimilarity, 0)
	targetTokens := ls.tokenizeCommand(command)

	for _, key := range keys {
		// Extract command hash from key
		parts := strings.Split(key, ":")
		if len(parts) != 4 {
			continue
		}
		commandHash := parts[3]

		// Get sample execution to extract command
		indexData, err := ls.stateAdapter.Get(ctx, key)
		if err != nil {
			continue
		}

		var executionKeys []string
		if err := json.Unmarshal([]byte(indexData), &executionKeys); err != nil {
			continue
		}

		if len(executionKeys) == 0 {
			continue
		}

		// Get first execution to extract command
		executionData, err := ls.stateAdapter.Get(ctx, executionKeys[0])
		if err != nil {
			continue
		}

		var execution TaskExecution
		if err := json.Unmarshal([]byte(executionData), &execution); err != nil {
			continue
		}

		// Calculate similarity
		similarity := ls.calculateSimilarity(targetTokens, ls.tokenizeCommand(execution.Command))
		if similarity > 0.3 { // Only include if reasonably similar
			similarities = append(similarities, commandSimilarity{
				hash:       commandHash,
				command:    execution.Command,
				similarity: similarity,
			})
		}
	}

	// Sort by similarity
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].similarity > similarities[j].similarity
	})

	// Get executions for most similar commands
	executions := make([]*TaskExecution, 0, limit)
	for i, sim := range similarities {
		if i >= limit {
			break
		}

		history, err := ls.GetTaskHistory(ctx, sim.command)
		if err != nil {
			continue
		}

		// Add most recent execution from each similar command
		if len(history) > 0 {
			executions = append(executions, history[0])
		}
	}

	return executions, nil
}

// GetTaskStats retrieves aggregated statistics for a task
func (ls *LearningStore) GetTaskStats(ctx context.Context, taskName string) (*TaskStats, error) {
	indexKey := fmt.Sprintf("ai:index:task:%s", taskName)
	
	// Get execution keys from index
	indexData, err := ls.stateAdapter.Get(ctx, indexKey)
	if err != nil {
		return nil, fmt.Errorf("no stats found for task: %s", taskName)
	}

	var executionKeys []string
	if err := json.Unmarshal([]byte(indexData), &executionKeys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index: %w", err)
	}

	stats := &TaskStats{
		TaskName:     taskName,
		TotalRuns:    len(executionKeys),
		SuccessCount: 0,
		FailureCount: 0,
		TotalTime:    0,
		AvgTime:      0,
	}

	// Calculate statistics
	for _, key := range executionKeys {
		executionData, err := ls.stateAdapter.Get(ctx, key)
		if err != nil {
			continue
		}

		var execution TaskExecution
		if err := json.Unmarshal([]byte(executionData), &execution); err != nil {
			continue
		}

		stats.TotalTime += execution.ExecutionTime
		if execution.Success {
			stats.SuccessCount++
		} else {
			stats.FailureCount++
		}
	}

	if stats.TotalRuns > 0 {
		stats.AvgTime = stats.TotalTime / time.Duration(stats.TotalRuns)
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.TotalRuns)
	}

	return stats, nil
}

// updateCommandIndex updates the index for a command
func (ls *LearningStore) updateCommandIndex(ctx context.Context, command string, executionKey string) error {
	commandHash := ls.hashCommand(command)
	indexKey := fmt.Sprintf("ai:index:command:%s", commandHash)
	
	// Get existing index
	var executionKeys []string
	indexData, err := ls.stateAdapter.Get(ctx, indexKey)
	if err == nil {
		json.Unmarshal([]byte(indexData), &executionKeys)
	}

	// Add new execution key
	executionKeys = append(executionKeys, executionKey)

	// Keep only last 100 executions
	if len(executionKeys) > 100 {
		executionKeys = executionKeys[len(executionKeys)-100:]
	}

	// Store updated index
	data, err := json.Marshal(executionKeys)
	if err != nil {
		return err
	}

	return ls.stateAdapter.Set(ctx, indexKey, string(data), 30*24*time.Hour)
}

// updateTaskIndex updates the index for a task name
func (ls *LearningStore) updateTaskIndex(ctx context.Context, taskName string, executionKey string) error {
	indexKey := fmt.Sprintf("ai:index:task:%s", taskName)
	
	// Get existing index
	var executionKeys []string
	indexData, err := ls.stateAdapter.Get(ctx, indexKey)
	if err == nil {
		json.Unmarshal([]byte(indexData), &executionKeys)
	}

	// Add new execution key
	executionKeys = append(executionKeys, executionKey)

	// Keep only last 100 executions
	if len(executionKeys) > 100 {
		executionKeys = executionKeys[len(executionKeys)-100:]
	}

	// Store updated index
	data, err := json.Marshal(executionKeys)
	if err != nil {
		return err
	}

	return ls.stateAdapter.Set(ctx, indexKey, string(data), 30*24*time.Hour)
}

// hashCommand creates a hash for a command for indexing
func (ls *LearningStore) hashCommand(command string) string {
	// Normalize command (remove extra spaces, etc.)
	normalized := strings.Join(strings.Fields(command), " ")
	
	// Create MD5 hash
	hash := md5.Sum([]byte(normalized))
	return fmt.Sprintf("%x", hash)
}

// tokenizeCommand breaks a command into tokens for similarity comparison
func (ls *LearningStore) tokenizeCommand(command string) []string {
	// Simple tokenization - split by spaces and common separators
	tokens := strings.FieldsFunc(command, func(r rune) bool {
		return r == ' ' || r == '=' || r == '-' || r == '/' || r == '\\'
	})
	
	// Convert to lowercase for case-insensitive comparison
	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}
	
	return tokens
}

// calculateSimilarity calculates Jaccard similarity between two token sets
func (ls *LearningStore) calculateSimilarity(tokens1, tokens2 []string) float64 {
	if len(tokens1) == 0 && len(tokens2) == 0 {
		return 1.0
	}
	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	// Create sets
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, token := range tokens1 {
		set1[token] = true
	}
	for _, token := range tokens2 {
		set2[token] = true
	}

	// Calculate intersection and union
	intersection := 0
	union := make(map[string]bool)
	
	for token := range set1 {
		union[token] = true
		if set2[token] {
			intersection++
		}
	}
	for token := range set2 {
		union[token] = true
	}

	// Jaccard similarity = |intersection| / |union|
	return float64(intersection) / float64(len(union))
}

// TaskStats represents aggregated statistics for a task
type TaskStats struct {
	TaskName     string        `json:"task_name"`
	TotalRuns    int           `json:"total_runs"`
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
	SuccessRate  float64       `json:"success_rate"`
	TotalTime    time.Duration `json:"total_time"`
	AvgTime      time.Duration `json:"avg_time"`
}

// commandSimilarity represents a command with its similarity score
type commandSimilarity struct {
	hash       string
	command    string
	similarity float64
}