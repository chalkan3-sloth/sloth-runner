package ai

import (
	"context"
	"testing"
	"time"
)

// Test TaskExecution structure
func TestTaskExecution_Structure(t *testing.T) {
	exec := &TaskExecution{
		TaskName:      "test-task",
		Command:       "echo test",
		Parameters:    map[string]interface{}{"arg1": "value1"},
		ExecutionTime: 5 * time.Second,
		Success:       true,
		ErrorMessage:  "",
		Timestamp:     time.Now(),
	}

	if exec.TaskName != "test-task" {
		t.Error("Expected TaskName to be set")
	}

	if exec.Command != "echo test" {
		t.Error("Expected Command to be set")
	}

	if !exec.Success {
		t.Error("Expected Success to be true")
	}
}

func TestTaskExecution_WithError(t *testing.T) {
	exec := &TaskExecution{
		TaskName:     "failing-task",
		Command:      "exit 1",
		Success:      false,
		ErrorMessage: "command failed",
	}

	if exec.Success {
		t.Error("Expected Success to be false")
	}

	if exec.ErrorMessage == "" {
		t.Error("Expected ErrorMessage to be set")
	}
}

func TestTaskExecution_Parameters(t *testing.T) {
	params := map[string]interface{}{
		"timeout":  30,
		"retry":    true,
		"env_vars": []string{"VAR1=value1"},
	}

	exec := &TaskExecution{
		Parameters: params,
	}

	if len(exec.Parameters) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(exec.Parameters))
	}

	if exec.Parameters["timeout"] != 30 {
		t.Error("Expected timeout parameter")
	}
}

func TestTaskExecution_Optimizations(t *testing.T) {
	exec := &TaskExecution{
		Optimizations: []string{"parallel", "cache"},
	}

	if len(exec.Optimizations) != 2 {
		t.Errorf("Expected 2 optimizations, got %d", len(exec.Optimizations))
	}
}

// Test SystemResources structure
func TestSystemResources_Structure(t *testing.T) {
	resources := SystemResources{
		CPUUsage:    45.5,
		MemoryUsage: 2048.0,
		DiskIO:      100.5,
		NetworkIO:   50.2,
		LoadAvg:     1.5,
	}

	if resources.CPUUsage != 45.5 {
		t.Error("Expected CPUUsage to be set")
	}

	if resources.MemoryUsage != 2048.0 {
		t.Error("Expected MemoryUsage to be set")
	}

	if resources.LoadAvg != 1.5 {
		t.Error("Expected LoadAvg to be set")
	}
}

func TestSystemResources_ZeroValues(t *testing.T) {
	resources := SystemResources{}

	if resources.CPUUsage != 0 {
		t.Error("Expected default CPUUsage to be 0")
	}

	if resources.MemoryUsage != 0 {
		t.Error("Expected default MemoryUsage to be 0")
	}
}

func TestSystemResources_HighLoad(t *testing.T) {
	resources := SystemResources{
		CPUUsage:  95.5,
		LoadAvg:   8.0,
	}

	if resources.CPUUsage <= 90 {
		t.Error("Expected high CPU usage")
	}

	if resources.LoadAvg <= 5 {
		t.Error("Expected high load average")
	}
}

// Test OptimizationSuggestion structure
func TestOptimizationSuggestion_Structure(t *testing.T) {
	suggestion := &OptimizationSuggestion{
		OriginalCommand:  "echo hello",
		OptimizedCommand: "echo hello | tee output.txt",
		Optimizations:    []Optimization{},
		ConfidenceScore:  0.85,
		ExpectedSpeedup:  1.5,
		ResourceSavings:  map[string]float64{"cpu": 10.0},
		Rationale:        "Added output redirection",
	}

	if suggestion.OriginalCommand == "" {
		t.Error("Expected OriginalCommand to be set")
	}

	if suggestion.OptimizedCommand == "" {
		t.Error("Expected OptimizedCommand to be set")
	}

	if suggestion.ConfidenceScore != 0.85 {
		t.Error("Expected ConfidenceScore to be set")
	}

	if suggestion.ExpectedSpeedup != 1.5 {
		t.Error("Expected ExpectedSpeedup to be set")
	}
}

func TestOptimizationSuggestion_WithOptimizations(t *testing.T) {
	opts := []Optimization{
		{Type: "parallelization", Description: "Use parallel execution", Impact: 0.5},
		{Type: "caching", Description: "Enable caching", Impact: 0.3},
	}

	suggestion := &OptimizationSuggestion{
		Optimizations: opts,
	}

	if len(suggestion.Optimizations) != 2 {
		t.Errorf("Expected 2 optimizations, got %d", len(suggestion.Optimizations))
	}
}

func TestOptimizationSuggestion_ResourceSavings(t *testing.T) {
	savings := map[string]float64{
		"cpu":    15.5,
		"memory": 256.0,
		"disk":   1024.0,
	}

	suggestion := &OptimizationSuggestion{
		ResourceSavings: savings,
	}

	if len(suggestion.ResourceSavings) != 3 {
		t.Error("Expected 3 resource savings entries")
	}

	if suggestion.ResourceSavings["cpu"] != 15.5 {
		t.Error("Expected CPU savings")
	}
}

func TestOptimizationSuggestion_Metadata(t *testing.T) {
	metadata := map[string]interface{}{
		"generated_at": time.Now(),
		"version":      "1.0",
	}

	suggestion := &OptimizationSuggestion{
		Metadata: metadata,
	}

	if len(suggestion.Metadata) != 2 {
		t.Error("Expected 2 metadata entries")
	}
}

// Test Optimization structure
func TestOptimization_Structure(t *testing.T) {
	opt := Optimization{
		Type:        "parallelization",
		Description: "Execute tasks in parallel",
		Impact:      0.5,
		Metadata:    map[string]interface{}{"threads": 4},
	}

	if opt.Type != "parallelization" {
		t.Error("Expected Type to be set")
	}

	if opt.Description == "" {
		t.Error("Expected Description to be set")
	}

	if opt.Impact != 0.5 {
		t.Error("Expected Impact to be set")
	}
}

func TestOptimization_Types(t *testing.T) {
	types := []string{
		"parallelization",
		"caching",
		"compression",
		"indexing",
		"batching",
	}

	for _, optType := range types {
		opt := Optimization{
			Type: optType,
		}

		if opt.Type != optType {
			t.Errorf("Expected type %s, got %s", optType, opt.Type)
		}
	}
}

func TestOptimization_HighImpact(t *testing.T) {
	opt := Optimization{
		Impact: 0.9,
	}

	if opt.Impact <= 0.5 {
		t.Error("Expected high impact optimization")
	}
}

func TestOptimization_LowImpact(t *testing.T) {
	opt := Optimization{
		Impact: 0.1,
	}

	if opt.Impact >= 0.5 {
		t.Error("Expected low impact optimization")
	}
}

// Test TaskOptimizer
func TestNewTaskOptimizer(t *testing.T) {
	optimizer := NewTaskOptimizer()

	if optimizer == nil {
		t.Error("Expected non-nil optimizer")
	}

	if optimizer.optimizers == nil {
		t.Error("Expected optimizers map to be initialized")
	}
}

func TestTaskOptimizer_HasOptimizers(t *testing.T) {
	optimizer := NewTaskOptimizer()

	if len(optimizer.optimizers) == 0 {
		t.Error("Expected built-in optimizers to be registered")
	}
}

// Test FailurePredictor
func TestNewFailurePredictor(t *testing.T) {
	predictor := NewFailurePredictor()

	if predictor == nil {
		t.Error("Expected non-nil predictor")
	}

	if predictor.models == nil {
		t.Error("Expected models map to be initialized")
	}
}

func TestFailurePredictor_HasModels(t *testing.T) {
	predictor := NewFailurePredictor()

	expectedModels := []string{"historical", "resource", "pattern", "time"}

	for _, modelName := range expectedModels {
		if _, exists := predictor.models[modelName]; !exists {
			t.Errorf("Expected model '%s' to be registered", modelName)
		}
	}
}

// Test FailurePrediction structure
func TestFailurePrediction_Structure(t *testing.T) {
	prediction := &FailurePrediction{
		TaskName:           "test-task",
		Command:            "echo test",
		FailureProbability: 0.15,
		Confidence:         0.85,
		RiskFactors:        []RiskFactor{},
		Recommendations:    []string{"Add timeout", "Enable retries"},
	}

	if prediction.TaskName == "" {
		t.Error("Expected TaskName to be set")
	}

	if prediction.FailureProbability != 0.15 {
		t.Error("Expected FailureProbability to be set")
	}

	if prediction.Confidence != 0.85 {
		t.Error("Expected Confidence to be set")
	}

	if len(prediction.Recommendations) != 2 {
		t.Error("Expected 2 recommendations")
	}
}

func TestFailurePrediction_HighRisk(t *testing.T) {
	prediction := &FailurePrediction{
		FailureProbability: 0.85,
	}

	if prediction.FailureProbability <= 0.5 {
		t.Error("Expected high failure probability")
	}
}

func TestFailurePrediction_LowRisk(t *testing.T) {
	prediction := &FailurePrediction{
		FailureProbability: 0.05,
	}

	if prediction.FailureProbability >= 0.5 {
		t.Error("Expected low failure probability")
	}
}

// Test RiskFactor structure
func TestRiskFactor_Structure(t *testing.T) {
	factor := RiskFactor{
		Factor:      "high_cpu",
		Impact:      0.6,
		Description: "CPU usage exceeds 80%",
	}

	if factor.Factor == "" {
		t.Error("Expected Factor to be set")
	}

	if factor.Impact != 0.6 {
		t.Error("Expected Impact to be set")
	}

	if factor.Description == "" {
		t.Error("Expected Description to be set")
	}
}

func TestRiskFactor_CommonFactors(t *testing.T) {
	factors := []string{
		"high_cpu",
		"low_memory",
		"network_latency",
		"disk_full",
		"timeout",
	}

	for _, factorName := range factors {
		factor := RiskFactor{
			Factor: factorName,
		}

		if factor.Factor != factorName {
			t.Errorf("Expected factor %s, got %s", factorName, factor.Factor)
		}
	}
}

// Test OptimizationContext structure
func TestOptimizationContext_Structure(t *testing.T) {
	ctx := OptimizationContext{
		Command:         "echo test",
		History:         []*TaskExecution{},
		SystemResources: SystemResources{},
		SimilarTasks:    []*TaskExecution{},
		Options:         map[string]interface{}{},
	}

	if ctx.Command == "" {
		t.Error("Expected Command to be set")
	}

	if ctx.History == nil {
		t.Error("Expected History to be initialized")
	}

	if ctx.Options == nil {
		t.Error("Expected Options to be initialized")
	}
}

func TestOptimizationContext_WithHistory(t *testing.T) {
	history := []*TaskExecution{
		{TaskName: "task1", Success: true},
		{TaskName: "task2", Success: false},
	}

	ctx := OptimizationContext{
		History: history,
	}

	if len(ctx.History) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(ctx.History))
	}
}

// Test PredictionContext structure
func TestPredictionContext_Structure(t *testing.T) {
	ctx := PredictionContext{
		TaskName:        "test-task",
		Command:         "echo test",
		History:         []*TaskExecution{},
		SystemResources: SystemResources{},
	}

	if ctx.TaskName == "" {
		t.Error("Expected TaskName to be set")
	}

	if ctx.Command == "" {
		t.Error("Expected Command to be set")
	}

	if ctx.History == nil {
		t.Error("Expected History to be initialized")
	}
}

func TestPredictionContext_WithResources(t *testing.T) {
	resources := SystemResources{
		CPUUsage:    75.0,
		MemoryUsage: 4096.0,
	}

	ctx := PredictionContext{
		SystemResources: resources,
	}

	if ctx.SystemResources.CPUUsage != 75.0 {
		t.Error("Expected CPUUsage to be set")
	}
}

// Test edge cases and special scenarios

func TestTaskExecution_EmptyParameters(t *testing.T) {
	exec := &TaskExecution{
		Parameters: map[string]interface{}{},
	}

	if exec.Parameters == nil {
		t.Error("Expected Parameters to be initialized")
	}

	if len(exec.Parameters) != 0 {
		t.Error("Expected empty Parameters")
	}
}

func TestTaskExecution_ZeroDuration(t *testing.T) {
	exec := &TaskExecution{
		ExecutionTime: 0,
	}

	if exec.ExecutionTime != 0 {
		t.Error("Expected zero ExecutionTime")
	}
}

func TestOptimizationSuggestion_NoOptimizations(t *testing.T) {
	suggestion := &OptimizationSuggestion{
		Optimizations: []Optimization{},
	}

	if len(suggestion.Optimizations) != 0 {
		t.Error("Expected no optimizations")
	}
}

func TestOptimizationSuggestion_ZeroConfidence(t *testing.T) {
	suggestion := &OptimizationSuggestion{
		ConfidenceScore: 0.0,
	}

	if suggestion.ConfidenceScore != 0.0 {
		t.Error("Expected zero confidence")
	}
}

func TestOptimizationSuggestion_HighConfidence(t *testing.T) {
	suggestion := &OptimizationSuggestion{
		ConfidenceScore: 0.95,
	}

	if suggestion.ConfidenceScore <= 0.9 {
		t.Error("Expected high confidence")
	}
}

func TestTaskExecution_LongRunning(t *testing.T) {
	exec := &TaskExecution{
		ExecutionTime: 30 * time.Minute,
	}

	if exec.ExecutionTime < 10*time.Minute {
		t.Error("Expected long execution time")
	}
}

func TestSystemResources_AllFieldsSet(t *testing.T) {
	resources := SystemResources{
		CPUUsage:    50.0,
		MemoryUsage: 1024.0,
		DiskIO:      200.0,
		NetworkIO:   100.0,
		LoadAvg:     2.5,
	}

	if resources.CPUUsage == 0 {
		t.Error("Expected CPUUsage to be non-zero")
	}
	if resources.MemoryUsage == 0 {
		t.Error("Expected MemoryUsage to be non-zero")
	}
	if resources.DiskIO == 0 {
		t.Error("Expected DiskIO to be non-zero")
	}
	if resources.NetworkIO == 0 {
		t.Error("Expected NetworkIO to be non-zero")
	}
	if resources.LoadAvg == 0 {
		t.Error("Expected LoadAvg to be non-zero")
	}
}

func TestOptimization_ZeroImpact(t *testing.T) {
	opt := Optimization{
		Impact: 0.0,
	}

	if opt.Impact != 0.0 {
		t.Error("Expected zero impact")
	}
}

func TestOptimization_NegativeImpact(t *testing.T) {
	opt := Optimization{
		Impact: -0.1,
	}

	if opt.Impact >= 0 {
		t.Error("Expected negative impact")
	}
}

func TestFailurePrediction_NoRiskFactors(t *testing.T) {
	prediction := &FailurePrediction{
		RiskFactors: []RiskFactor{},
	}

	if len(prediction.RiskFactors) != 0 {
		t.Error("Expected no risk factors")
	}
}

func TestFailurePrediction_MultipleRiskFactors(t *testing.T) {
	factors := []RiskFactor{
		{Factor: "high_cpu", Impact: 0.5},
		{Factor: "low_memory", Impact: 0.3},
		{Factor: "network_issues", Impact: 0.2},
	}

	prediction := &FailurePrediction{
		RiskFactors: factors,
	}

	if len(prediction.RiskFactors) != 3 {
		t.Errorf("Expected 3 risk factors, got %d", len(prediction.RiskFactors))
	}
}

func TestOptimizationContext_EmptyOptions(t *testing.T) {
	ctx := OptimizationContext{
		Options: map[string]interface{}{},
	}

	if len(ctx.Options) != 0 {
		t.Error("Expected empty options")
	}
}

func TestOptimizationContext_WithOptions(t *testing.T) {
	options := map[string]interface{}{
		"max_parallel": 4,
		"timeout":      30,
		"retry":        true,
	}

	ctx := OptimizationContext{
		Options: options,
	}

	if len(ctx.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(ctx.Options))
	}
}

func TestPredictionContext_EmptyHistory(t *testing.T) {
	ctx := PredictionContext{
		History: []*TaskExecution{},
	}

	if len(ctx.History) != 0 {
		t.Error("Expected empty history")
	}
}

func TestPredictionContext_WithHistory(t *testing.T) {
	history := []*TaskExecution{
		{TaskName: "task1", Success: true},
		{TaskName: "task2", Success: true},
		{TaskName: "task3", Success: false},
	}

	ctx := PredictionContext{
		History: history,
	}

	if len(ctx.History) != 3 {
		t.Errorf("Expected 3 history entries, got %d", len(ctx.History))
	}
}

func TestTaskExecution_CompleteRecord(t *testing.T) {
	exec := &TaskExecution{
		TaskName:      "complete-task",
		Command:       "ls -la",
		Parameters:    map[string]interface{}{"dir": "/home"},
		ExecutionTime: 2 * time.Second,
		Success:       true,
		ErrorMessage:  "",
		SystemResources: SystemResources{
			CPUUsage:    30.0,
			MemoryUsage: 512.0,
		},
		Timestamp:     time.Now(),
		Optimizations: []string{"cache"},
	}

	if exec.TaskName == "" {
		t.Error("TaskName should be set")
	}
	if exec.Command == "" {
		t.Error("Command should be set")
	}
	if len(exec.Parameters) == 0 {
		t.Error("Parameters should be set")
	}
	if !exec.Success {
		t.Error("Success should be true")
	}
	if exec.SystemResources.CPUUsage == 0 {
		t.Error("SystemResources should be set")
	}
}

// Test context usage (basic tests)
func TestWithContext(t *testing.T) {
	ctx := context.Background()

	if ctx == nil {
		t.Error("Expected non-nil context")
	}

	if ctx.Err() != nil {
		t.Error("Expected no error from background context")
	}
}

func TestWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if ctx.Err() == nil {
		t.Error("Expected error from cancelled context")
	}
}
