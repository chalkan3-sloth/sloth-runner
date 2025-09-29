package ai

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// FailurePredictor provides AI-powered failure prediction
type FailurePredictor struct {
	models map[string]PredictionModel
}

// PredictionModel represents a model for predicting failures
type PredictionModel interface {
	Predict(context PredictionContext) float64
	Update(execution *TaskExecution) error
}

// PredictionContext provides context for failure prediction
type PredictionContext struct {
	TaskName        string
	Command         string
	History         []*TaskExecution
	SystemResources SystemResources
}

// FailurePrediction represents a failure prediction result
type FailurePrediction struct {
	TaskName           string                 `json:"task_name"`
	Command            string                 `json:"command"`
	FailureProbability float64                `json:"failure_probability"`
	Confidence         float64                `json:"confidence"`
	RiskFactors        []RiskFactor           `json:"risk_factors"`
	Recommendations    []string               `json:"recommendations"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// RiskFactor represents a specific risk factor
type RiskFactor struct {
	Factor      string  `json:"factor"`
	Impact      float64 `json:"impact"`
	Description string  `json:"description"`
}

// NewFailurePredictor creates a new failure predictor
func NewFailurePredictor() *FailurePredictor {
	predictor := &FailurePredictor{
		models: make(map[string]PredictionModel),
	}
	
	// Register built-in models
	predictor.models["historical"] = &HistoricalModel{}
	predictor.models["resource"] = &ResourceModel{}
	predictor.models["pattern"] = &PatternModel{}
	predictor.models["time"] = &TimeBasedModel{}
	
	return predictor
}

// PredictFailure predicts the likelihood of task failure
func (fp *FailurePredictor) PredictFailure(context PredictionContext) *FailurePrediction {
	predictions := make(map[string]float64)
	
	// Get predictions from all models
	for name, model := range fp.models {
		probability := model.Predict(context)
		predictions[name] = probability
	}
	
	// Calculate ensemble prediction (weighted average)
	weights := map[string]float64{
		"historical": 0.4,
		"resource":   0.25,
		"pattern":    0.2,
		"time":       0.15,
	}
	
	finalProbability := 0.0
	for name, weight := range weights {
		if prob, exists := predictions[name]; exists {
			finalProbability += prob * weight
		}
	}
	
	// Analyze risk factors
	riskFactors := fp.analyzeRiskFactors(context, predictions)
	
	// Generate recommendations
	recommendations := fp.generateRecommendations(context, riskFactors)
	
	// Calculate confidence based on available data
	confidence := fp.calculateConfidence(context)
	
	return &FailurePrediction{
		TaskName:           context.TaskName,
		Command:            context.Command,
		FailureProbability: finalProbability,
		Confidence:         confidence,
		RiskFactors:        riskFactors,
		Recommendations:    recommendations,
		Metadata: map[string]interface{}{
			"model_predictions": predictions,
			"generated_at":      time.Now(),
			"history_size":      len(context.History),
		},
	}
}

// UpdateModel updates prediction models with new execution data
func (fp *FailurePredictor) UpdateModel(execution *TaskExecution) error {
	for name, model := range fp.models {
		if err := model.Update(execution); err != nil {
			return fmt.Errorf("failed to update model %s: %w", name, err)
		}
	}
	return nil
}

// analyzeRiskFactors analyzes various risk factors
func (fp *FailurePredictor) analyzeRiskFactors(context PredictionContext, predictions map[string]float64) []RiskFactor {
	factors := make([]RiskFactor, 0)
	
	// Historical failure rate
	if len(context.History) > 0 {
		failures := 0
		for _, exec := range context.History {
			if !exec.Success {
				failures++
			}
		}
		failureRate := float64(failures) / float64(len(context.History))
		
		if failureRate > 0.2 {
			factors = append(factors, RiskFactor{
				Factor:      "high_historical_failure_rate",
				Impact:      failureRate,
				Description: fmt.Sprintf("Task has %.1f%% historical failure rate", failureRate*100),
			})
		}
	}
	
	// System resource pressure
	if context.SystemResources.CPUUsage > 80 {
		factors = append(factors, RiskFactor{
			Factor:      "high_cpu_usage",
			Impact:      context.SystemResources.CPUUsage / 100,
			Description: fmt.Sprintf("High CPU usage: %.1f%%", context.SystemResources.CPUUsage),
		})
	}
	
	if context.SystemResources.MemoryUsage > 90 {
		factors = append(factors, RiskFactor{
			Factor:      "high_memory_usage",
			Impact:      context.SystemResources.MemoryUsage / 100,
			Description: fmt.Sprintf("High memory usage: %.1f%%", context.SystemResources.MemoryUsage),
		})
	}
	
	// Command complexity
	complexity := fp.calculateCommandComplexity(context.Command)
	if complexity > 0.7 {
		factors = append(factors, RiskFactor{
			Factor:      "complex_command",
			Impact:      complexity,
			Description: "Command is complex and may be error-prone",
		})
	}
	
	// Time-based factors
	now := time.Now()
	if now.Hour() < 6 || now.Hour() > 22 {
		factors = append(factors, RiskFactor{
			Factor:      "off_hours_execution",
			Impact:      0.3,
			Description: "Executing during off-hours may have higher failure risk",
		})
	}
	
	// Sort by impact
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].Impact > factors[j].Impact
	})
	
	return factors
}

// generateRecommendations generates actionable recommendations
func (fp *FailurePredictor) generateRecommendations(context PredictionContext, riskFactors []RiskFactor) []string {
	recommendations := make([]string, 0)
	
	for _, factor := range riskFactors {
		switch factor.Factor {
		case "high_historical_failure_rate":
			recommendations = append(recommendations, 
				"Consider reviewing and improving task implementation",
				"Add more comprehensive error handling",
				"Implement retry mechanisms with exponential backoff")
			
		case "high_cpu_usage":
			recommendations = append(recommendations,
				"Wait for CPU usage to decrease before executing",
				"Consider running task with lower priority",
				"Schedule task for off-peak hours")
			
		case "high_memory_usage":
			recommendations = append(recommendations,
				"Free up memory before executing task",
				"Consider breaking task into smaller chunks",
				"Monitor for memory leaks in task implementation")
			
		case "complex_command":
			recommendations = append(recommendations,
				"Break complex command into simpler steps",
				"Add intermediate validation points",
				"Implement comprehensive logging for debugging")
			
		case "off_hours_execution":
			recommendations = append(recommendations,
				"Consider scheduling task during business hours",
				"Ensure monitoring systems are active",
				"Have incident response procedures ready")
		}
	}
	
	// Remove duplicates
	uniqueRecs := make([]string, 0)
	seen := make(map[string]bool)
	for _, rec := range recommendations {
		if !seen[rec] {
			uniqueRecs = append(uniqueRecs, rec)
			seen[rec] = true
		}
	}
	
	return uniqueRecs
}

// calculateConfidence calculates confidence in the prediction
func (fp *FailurePredictor) calculateConfidence(context PredictionContext) float64 {
	confidence := 0.5 // Base confidence
	
	// More history = higher confidence
	if len(context.History) > 0 {
		historyFactor := math.Min(float64(len(context.History))/20.0, 1.0)
		confidence += historyFactor * 0.3
	}
	
	// Recent executions = higher confidence
	recentExecutions := 0
	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	for _, exec := range context.History {
		if exec.Timestamp.After(cutoff) {
			recentExecutions++
		}
	}
	if recentExecutions > 0 {
		recentFactor := math.Min(float64(recentExecutions)/10.0, 1.0)
		confidence += recentFactor * 0.2
	}
	
	// Cap confidence at 0.95
	return math.Min(confidence, 0.95)
}

// calculateCommandComplexity calculates the complexity of a command
func (fp *FailurePredictor) calculateCommandComplexity(command string) float64 {
	complexity := 0.0
	
	// Length factor
	if len(command) > 100 {
		complexity += 0.2
	}
	
	// Pipe complexity
	pipes := strings.Count(command, "|")
	complexity += float64(pipes) * 0.1
	
	// Redirection complexity
	redirections := strings.Count(command, ">") + strings.Count(command, "<")
	complexity += float64(redirections) * 0.05
	
	// Special characters
	specials := strings.Count(command, "&") + strings.Count(command, ";") + 
	           strings.Count(command, "&&") + strings.Count(command, "||")
	complexity += float64(specials) * 0.1
	
	// Network operations
	if strings.Contains(command, "http") || strings.Contains(command, "curl") || 
	   strings.Contains(command, "wget") {
		complexity += 0.3
	}
	
	// File operations
	if strings.Contains(command, "rm") || strings.Contains(command, "mv") || 
	   strings.Contains(command, "cp") {
		complexity += 0.2
	}
	
	return math.Min(complexity, 1.0)
}

// HistoricalModel predicts based on historical failure patterns
type HistoricalModel struct{}

func (hm *HistoricalModel) Predict(context PredictionContext) float64 {
	if len(context.History) == 0 {
		return 0.5 // No history, neutral prediction
	}
	
	// Calculate recent failure rate (last 10 executions)
	recentHistory := context.History
	if len(recentHistory) > 10 {
		recentHistory = recentHistory[:10]
	}
	
	failures := 0
	for _, exec := range recentHistory {
		if !exec.Success {
			failures++
		}
	}
	
	return float64(failures) / float64(len(recentHistory))
}

func (hm *HistoricalModel) Update(execution *TaskExecution) error {
	// Historical model is stateless - no update needed
	return nil
}

// ResourceModel predicts based on system resource availability
type ResourceModel struct{}

func (rm *ResourceModel) Predict(context PredictionContext) float64 {
	resources := context.SystemResources
	
	// High resource usage increases failure probability
	cpuFactor := resources.CPUUsage / 100.0
	memoryFactor := resources.MemoryUsage / 100.0
	loadFactor := math.Min(resources.LoadAvg/4.0, 1.0) // Assume 4-core system
	
	// Weighted combination
	probability := (cpuFactor*0.4 + memoryFactor*0.4 + loadFactor*0.2)
	
	// Apply sigmoid function to smooth the curve
	return 1.0 / (1.0 + math.Exp(-5*(probability-0.5)))
}

func (rm *ResourceModel) Update(execution *TaskExecution) error {
	// Resource model is stateless - no update needed
	return nil
}

// PatternModel predicts based on command patterns
type PatternModel struct{}

func (pm *PatternModel) Predict(context PredictionContext) float64 {
	command := strings.ToLower(context.Command)
	
	// Risk patterns
	riskPatterns := map[string]float64{
		"rm -rf":     0.8, // Dangerous deletion
		"sudo":       0.3, // Elevated privileges
		"curl":       0.4, // Network dependency
		"wget":       0.4, // Network dependency
		"git":        0.2, // External dependency
		"docker":     0.3, // Container complexity
		"kubernetes": 0.4, // Orchestration complexity
		"npm":        0.3, // Package management
		"pip":        0.3, // Package management
		"make":       0.2, // Build system
	}
	
	maxRisk := 0.0
	for pattern, risk := range riskPatterns {
		if strings.Contains(command, pattern) {
			if risk > maxRisk {
				maxRisk = risk
			}
		}
	}
	
	return maxRisk
}

func (pm *PatternModel) Update(execution *TaskExecution) error {
	// Pattern model is stateless - no update needed
	return nil
}

// TimeBasedModel predicts based on time patterns
type TimeBasedModel struct{}

func (tm *TimeBasedModel) Predict(context PredictionContext) float64 {
	now := time.Now()
	
	// Higher failure probability during off-hours
	hour := now.Hour()
	if hour < 6 || hour > 22 {
		return 0.3
	}
	
	// Weekend execution might have higher risk
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return 0.2
	}
	
	// Monday morning might have higher risk due to weekend changes
	if now.Weekday() == time.Monday && hour < 10 {
		return 0.25
	}
	
	return 0.1 // Normal business hours
}

func (tm *TimeBasedModel) Update(execution *TaskExecution) error {
	// Time-based model is stateless - no update needed
	return nil
}