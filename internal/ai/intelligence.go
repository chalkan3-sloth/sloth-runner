package ai

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/state"
)

// TaskIntelligence provides AI-powered task optimization and learning
type TaskIntelligence struct {
	stateManager  *state.StateManager
	learningStore *LearningStore
	optimizer     *TaskOptimizer
	predictor     *FailurePredictor
}

// TaskExecution represents a task execution record for learning
type TaskExecution struct {
	TaskName        string                 `json:"task_name"`
	Command         string                 `json:"command"`
	Parameters      map[string]interface{} `json:"parameters"`
	ExecutionTime   time.Duration          `json:"execution_time"`
	Success         bool                   `json:"success"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	SystemResources SystemResources        `json:"system_resources"`
	Timestamp       time.Time              `json:"timestamp"`
	Optimizations   []string               `json:"optimizations,omitempty"`
}

// SystemResources captures system state during task execution
type SystemResources struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskIO      float64 `json:"disk_io"`
	NetworkIO   float64 `json:"network_io"`
	LoadAvg     float64 `json:"load_avg"`
}

// OptimizationSuggestion represents AI-generated optimization recommendations
type OptimizationSuggestion struct {
	OriginalCommand   string                 `json:"original_command"`
	OptimizedCommand  string                 `json:"optimized_command"`
	Optimizations     []Optimization         `json:"optimizations"`
	ConfidenceScore   float64                `json:"confidence_score"`
	ExpectedSpeedup   float64                `json:"expected_speedup"`
	ResourceSavings   map[string]float64     `json:"resource_savings"`
	Rationale         string                 `json:"rationale"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// Optimization represents a specific optimization applied
type Optimization struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Impact      float64                `json:"impact"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewTaskIntelligence creates a new AI task intelligence system
func NewTaskIntelligence(stateManager *state.StateManager) *TaskIntelligence {
	return &TaskIntelligence{
		stateManager:  stateManager,
		learningStore: NewLearningStore(stateManager),
		optimizer:     NewTaskOptimizer(),
		predictor:     NewFailurePredictor(),
	}
}

// RecordExecution records a task execution for learning
func (ti *TaskIntelligence) RecordExecution(ctx context.Context, execution *TaskExecution) error {
	slog.Info("Recording task execution for AI learning",
		"task", execution.TaskName,
		"success", execution.Success,
		"duration", execution.ExecutionTime)

	// Store execution data
	if err := ti.learningStore.StoreExecution(ctx, execution); err != nil {
		return fmt.Errorf("failed to store execution: %w", err)
	}

	// Update learning models
	if err := ti.updateLearningModels(ctx, execution); err != nil {
		slog.Warn("Failed to update learning models", "error", err)
	}

	return nil
}

// OptimizeCommand generates AI-powered command optimizations
func (ti *TaskIntelligence) OptimizeCommand(ctx context.Context, command string, options map[string]interface{}) (*OptimizationSuggestion, error) {
	slog.Info("Generating AI optimization suggestions", "command", command)

	// Get historical data
	history, err := ti.learningStore.GetTaskHistory(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to get task history: %w", err)
	}

	// Get system resources
	systemResources := ti.getCurrentSystemResources()

	// Find similar tasks
	similarTasks, err := ti.learningStore.FindSimilarTasks(ctx, command, 10)
	if err != nil {
		slog.Warn("Failed to find similar tasks", "error", err)
		similarTasks = []*TaskExecution{}
	}

	// Generate optimization suggestions
	suggestion := ti.optimizer.GenerateOptimizations(OptimizationContext{
		Command:         command,
		History:         history,
		SystemResources: systemResources,
		SimilarTasks:    similarTasks,
		Options:         options,
	})

	slog.Info("Generated optimization suggestion",
		"confidence", suggestion.ConfidenceScore,
		"expected_speedup", suggestion.ExpectedSpeedup)

	return suggestion, nil
}

// PredictFailure predicts the likelihood of task failure
func (ti *TaskIntelligence) PredictFailure(ctx context.Context, taskName string, command string) (*FailurePrediction, error) {
	slog.Info("Predicting task failure probability", "task", taskName)

	// Get historical data
	history, err := ti.learningStore.GetTaskHistory(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to get task history: %w", err)
	}

	// Get current system state
	systemResources := ti.getCurrentSystemResources()

	// Generate prediction
	prediction := ti.predictor.PredictFailure(PredictionContext{
		TaskName:        taskName,
		Command:         command,
		History:         history,
		SystemResources: systemResources,
	})

	slog.Info("Generated failure prediction",
		"task", taskName,
		"probability", prediction.FailureProbability,
		"confidence", prediction.Confidence)

	return prediction, nil
}

// GetTaskHistory retrieves historical execution data for a task
func (ti *TaskIntelligence) GetTaskHistory(ctx context.Context, command string) ([]*TaskExecution, error) {
	return ti.learningStore.GetTaskHistory(ctx, command)
}

// FindSimilarTasks finds tasks similar to the given command
func (ti *TaskIntelligence) FindSimilarTasks(ctx context.Context, command string, limit int) ([]*TaskExecution, error) {
	return ti.learningStore.FindSimilarTasks(ctx, command, limit)
}

// GetTaskStats retrieves aggregated statistics for a task
func (ti *TaskIntelligence) GetTaskStats(ctx context.Context, taskName string) (*TaskStats, error) {
	return ti.learningStore.GetTaskStats(ctx, taskName)
}

// updateLearningModels updates the AI models with new execution data
func (ti *TaskIntelligence) updateLearningModels(ctx context.Context, execution *TaskExecution) error {
	// Update optimizer models
	if err := ti.optimizer.UpdateModel(execution); err != nil {
		return fmt.Errorf("failed to update optimizer model: %w", err)
	}

	// Update predictor models
	if err := ti.predictor.UpdateModel(execution); err != nil {
		return fmt.Errorf("failed to update predictor model: %w", err)
	}

	return nil
}

// getCurrentSystemResources gets current system resource usage
func (ti *TaskIntelligence) getCurrentSystemResources() SystemResources {
	// In a real implementation, this would collect actual system metrics
	// For now, return mock data
	return SystemResources{
		CPUUsage:    50.0,
		MemoryUsage: 60.0,
		DiskIO:      10.0,
		NetworkIO:   5.0,
		LoadAvg:     1.5,
	}
}

// LearningMode represents different AI learning modes
type LearningMode string

const (
	LearningModeAdaptive    LearningMode = "adaptive"
	LearningModeAggressive  LearningMode = "aggressive"
	LearningModeConservative LearningMode = "conservative"
	LearningModeExperimental LearningMode = "experimental"
)

// AITaskConfig represents AI configuration for a task
type AITaskConfig struct {
	Enabled         bool                   `json:"enabled"`
	LearningMode    LearningMode           `json:"learning_mode"`
	OptimizationLevel int                  `json:"optimization_level"` // 1-10
	FailurePrediction bool                 `json:"failure_prediction"`
	AutoOptimize     bool                  `json:"auto_optimize"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// DefaultAITaskConfig returns default AI configuration
func DefaultAITaskConfig() AITaskConfig {
	return AITaskConfig{
		Enabled:           true,
		LearningMode:      LearningModeAdaptive,
		OptimizationLevel: 5,
		FailurePrediction: true,
		AutoOptimize:      false,
		Metadata:          make(map[string]interface{}),
	}
}