package ai

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// TaskOptimizer provides AI-powered task optimization
type TaskOptimizer struct {
	optimizers map[string]OptimizerFunc
}

// OptimizerFunc represents a function that can optimize commands
type OptimizerFunc func(context OptimizationContext) []Optimization

// OptimizationContext provides context for optimization
type OptimizationContext struct {
	Command         string
	History         []*TaskExecution
	SystemResources SystemResources
	SimilarTasks    []*TaskExecution
	Options         map[string]interface{}
}

// NewTaskOptimizer creates a new task optimizer
func NewTaskOptimizer() *TaskOptimizer {
	optimizer := &TaskOptimizer{
		optimizers: make(map[string]OptimizerFunc),
	}
	
	// Register built-in optimizers
	optimizer.registerBuiltinOptimizers()
	
	return optimizer
}

// GenerateOptimizations generates optimization suggestions for a command
func (to *TaskOptimizer) GenerateOptimizations(context OptimizationContext) *OptimizationSuggestion {
	optimizations := make([]Optimization, 0)
	
	// Apply all registered optimizers
	for name, optimizer := range to.optimizers {
		opts := optimizer(context)
		for _, opt := range opts {
			opt.Metadata = map[string]interface{}{
				"optimizer": name,
				"generated_at": time.Now(),
			}
			optimizations = append(optimizations, opt)
		}
	}
	
	// Sort optimizations by impact (highest first)
	sort.Slice(optimizations, func(i, j int) bool {
		return optimizations[i].Impact > optimizations[j].Impact
	})
	
	// Generate optimized command
	optimizedCommand := to.applyOptimizations(context.Command, optimizations)
	
	// Calculate confidence score and expected speedup
	confidence := to.calculateConfidence(context, optimizations)
	speedup := to.calculateExpectedSpeedup(context, optimizations)
	resourceSavings := to.calculateResourceSavings(optimizations)
	
	return &OptimizationSuggestion{
		OriginalCommand:  context.Command,
		OptimizedCommand: optimizedCommand,
		Optimizations:    optimizations,
		ConfidenceScore:  confidence,
		ExpectedSpeedup:  speedup,
		ResourceSavings:  resourceSavings,
		Rationale:        to.generateRationale(optimizations),
		Metadata: map[string]interface{}{
			"generated_at":     time.Now(),
			"optimization_count": len(optimizations),
			"system_resources": context.SystemResources,
		},
	}
}

// UpdateModel updates the optimizer model with new execution data
func (to *TaskOptimizer) UpdateModel(execution *TaskExecution) error {
	// In a real implementation, this would update ML models
	// For now, we'll just log the update
	return nil
}

// registerBuiltinOptimizers registers built-in optimization functions
func (to *TaskOptimizer) registerBuiltinOptimizers() {
	to.optimizers["parallel"] = to.optimizeParallelization
	to.optimizers["memory"] = to.optimizeMemoryUsage
	to.optimizers["compiler"] = to.optimizeCompilerFlags
	to.optimizers["cache"] = to.optimizeCaching
	to.optimizers["network"] = to.optimizeNetworkOperations
	to.optimizers["io"] = to.optimizeIOOperations
	to.optimizers["resource"] = to.optimizeResourceAllocation
}

// optimizeParallelization suggests parallelization improvements
func (to *TaskOptimizer) optimizeParallelization(context OptimizationContext) []Optimization {
	optimizations := make([]Optimization, 0)
	
	command := strings.ToLower(context.Command)
	
	// Detect build commands that can benefit from parallel execution
	if strings.Contains(command, "go build") && !strings.Contains(command, "-p") {
		optimizations = append(optimizations, Optimization{
			Type:        "parallelization",
			Description: "Add parallel build flag to utilize multiple CPU cores",
			Impact:      0.7,
			Metadata: map[string]interface{}{
				"flag": "-p",
				"suggested_workers": context.SystemResources.CPUUsage,
			},
		})
	}
	
	// Detect make commands that can be parallelized
	if strings.Contains(command, "make") && !strings.Contains(command, "-j") {
		workers := int(context.SystemResources.LoadAvg) + 1
		optimizations = append(optimizations, Optimization{
			Type:        "parallelization",
			Description: fmt.Sprintf("Add -j%d flag for parallel make execution", workers),
			Impact:      0.8,
			Metadata: map[string]interface{}{
				"flag": fmt.Sprintf("-j%d", workers),
				"workers": workers,
			},
		})
	}
	
	// Detect npm/yarn commands that can be parallelized
	if (strings.Contains(command, "npm") || strings.Contains(command, "yarn")) && 
	   strings.Contains(command, "install") && !strings.Contains(command, "--network-concurrency") {
		optimizations = append(optimizations, Optimization{
			Type:        "parallelization",
			Description: "Add network concurrency for faster package installation",
			Impact:      0.6,
			Metadata: map[string]interface{}{
				"flag": "--network-concurrency=16",
			},
		})
	}
	
	return optimizations
}

// optimizeMemoryUsage suggests memory optimization improvements
func (to *TaskOptimizer) optimizeMemoryUsage(context OptimizationContext) []Optimization {
	optimizations := make([]Optimization, 0)
	
	command := strings.ToLower(context.Command)
	
	// Optimize Java applications
	if strings.Contains(command, "java") && !strings.Contains(command, "-Xmx") {
		availableMemory := (100 - context.SystemResources.MemoryUsage) / 100 * 8 // Assume 8GB system
		suggestedMemory := int(availableMemory * 0.8) // Use 80% of available
		
		optimizations = append(optimizations, Optimization{
			Type:        "memory",
			Description: fmt.Sprintf("Set heap size to %dG for optimal memory usage", suggestedMemory),
			Impact:      0.5,
			Metadata: map[string]interface{}{
				"flag": fmt.Sprintf("-Xmx%dg", suggestedMemory),
				"available_memory": availableMemory,
			},
		})
	}
	
	// Optimize Go builds
	if strings.Contains(command, "go build") && context.SystemResources.MemoryUsage < 50 {
		optimizations = append(optimizations, Optimization{
			Type:        "memory",
			Description: "Enable compiler optimizations for better memory usage",
			Impact:      0.3,
			Metadata: map[string]interface{}{
				"flags": []string{"-ldflags", "-s -w"},
			},
		})
	}
	
	return optimizations
}

// optimizeCompilerFlags suggests compiler optimization improvements
func (to *TaskOptimizer) optimizeCompilerFlags(context OptimizationContext) []Optimization {
	optimizations := make([]Optimization, 0)
	
	command := strings.ToLower(context.Command)
	
	// Optimize C/C++ compilation
	if (strings.Contains(command, "gcc") || strings.Contains(command, "g++") || 
		strings.Contains(command, "clang")) && !strings.Contains(command, "-O") {
		
		optimizations = append(optimizations, Optimization{
			Type:        "compiler",
			Description: "Add optimization flags for faster execution",
			Impact:      0.6,
			Metadata: map[string]interface{}{
				"flags": []string{"-O2", "-march=native"},
			},
		})
	}
	
	// Optimize Rust compilation
	if strings.Contains(command, "cargo build") && !strings.Contains(command, "--release") {
		optimizations = append(optimizations, Optimization{
			Type:        "compiler",
			Description: "Use release mode for optimized builds",
			Impact:      0.8,
			Metadata: map[string]interface{}{
				"flag": "--release",
			},
		})
	}
	
	return optimizations
}

// optimizeCaching suggests caching improvements
func (to *TaskOptimizer) optimizeCaching(context OptimizationContext) []Optimization {
	optimizations := make([]Optimization, 0)
	
	command := strings.ToLower(context.Command)
	
	// Docker build caching
	if strings.Contains(command, "docker build") && !strings.Contains(command, "--cache-from") {
		optimizations = append(optimizations, Optimization{
			Type:        "caching",
			Description: "Enable Docker layer caching for faster builds",
			Impact:      0.7,
			Metadata: map[string]interface{}{
				"suggestion": "Use multi-stage builds and --cache-from flag",
			},
		})
	}
	
	// Node.js dependency caching
	if strings.Contains(command, "npm install") && !strings.Contains(command, "--cache") {
		optimizations = append(optimizations, Optimization{
			Type:        "caching",
			Description: "Enable npm cache for faster dependency installation",
			Impact:      0.5,
			Metadata: map[string]interface{}{
				"flag": "--cache ~/.npm",
			},
		})
	}
	
	return optimizations
}

// optimizeNetworkOperations suggests network optimization improvements
func (to *TaskOptimizer) optimizeNetworkOperations(context OptimizationContext) []Optimization {
	optimizations := make([]Optimization, 0)
	
	command := strings.ToLower(context.Command)
	
	// Git operations
	if strings.Contains(command, "git clone") && !strings.Contains(command, "--depth") {
		optimizations = append(optimizations, Optimization{
			Type:        "network",
			Description: "Use shallow clone to reduce network transfer",
			Impact:      0.6,
			Metadata: map[string]interface{}{
				"flag": "--depth=1",
			},
		})
	}
	
	// Package manager optimizations
	if strings.Contains(command, "apt-get") && !strings.Contains(command, "--no-install-recommends") {
		optimizations = append(optimizations, Optimization{
			Type:        "network",
			Description: "Skip recommended packages to reduce download size",
			Impact:      0.4,
			Metadata: map[string]interface{}{
				"flag": "--no-install-recommends",
			},
		})
	}
	
	return optimizations
}

// optimizeIOOperations suggests I/O optimization improvements
func (to *TaskOptimizer) optimizeIOOperations(context OptimizationContext) []Optimization {
	optimizations := make([]Optimization, 0)
	
	command := strings.ToLower(context.Command)
	
	// File operations
	if strings.Contains(command, "cp") && !strings.Contains(command, "-a") {
		optimizations = append(optimizations, Optimization{
			Type:        "io",
			Description: "Use archive mode for efficient file copying",
			Impact:      0.3,
			Metadata: map[string]interface{}{
				"flag": "-a",
			},
		})
	}
	
	// Compression operations
	if strings.Contains(command, "tar") && strings.Contains(command, "-z") && 
	   !strings.Contains(command, "--use-compress-program") {
		optimizations = append(optimizations, Optimization{
			Type:        "io",
			Description: "Use pigz for parallel compression",
			Impact:      0.5,
			Metadata: map[string]interface{}{
				"flag": "--use-compress-program=pigz",
			},
		})
	}
	
	return optimizations
}

// optimizeResourceAllocation suggests resource allocation improvements
func (to *TaskOptimizer) optimizeResourceAllocation(context OptimizationContext) []Optimization {
	optimizations := make([]Optimization, 0)
	
	// Analyze historical performance
	if len(context.History) > 0 {
		avgTime := time.Duration(0)
		for _, exec := range context.History {
			avgTime += exec.ExecutionTime
		}
		avgTime /= time.Duration(len(context.History))
		
		// If task consistently takes a long time, suggest resource optimizations
		if avgTime > 5*time.Minute {
			optimizations = append(optimizations, Optimization{
				Type:        "resource",
				Description: "Consider increasing allocated resources for this long-running task",
				Impact:      0.4,
				Metadata: map[string]interface{}{
					"avg_execution_time": avgTime.String(),
					"suggestion": "increase_cpu_memory",
				},
			})
		}
	}
	
	return optimizations
}

// applyOptimizations applies optimizations to a command
func (to *TaskOptimizer) applyOptimizations(originalCommand string, optimizations []Optimization) string {
	command := originalCommand
	
	for _, opt := range optimizations {
		if opt.Impact < 0.3 {
			continue // Skip low-impact optimizations
		}
		
		switch opt.Type {
		case "parallelization":
			if flag, ok := opt.Metadata["flag"].(string); ok {
				command = to.addFlag(command, flag)
			}
		case "memory":
			if flag, ok := opt.Metadata["flag"].(string); ok {
				command = to.addFlag(command, flag)
			}
		case "compiler":
			if flags, ok := opt.Metadata["flags"].([]string); ok {
				for _, flag := range flags {
					command = to.addFlag(command, flag)
				}
			}
		case "caching", "network", "io":
			if flag, ok := opt.Metadata["flag"].(string); ok {
				command = to.addFlag(command, flag)
			}
		}
	}
	
	return command
}

// addFlag adds a flag to a command if it doesn't already exist
func (to *TaskOptimizer) addFlag(command, flag string) string {
	if strings.Contains(command, flag) {
		return command
	}
	return command + " " + flag
}

// calculateConfidence calculates confidence score for optimizations
func (to *TaskOptimizer) calculateConfidence(context OptimizationContext, optimizations []Optimization) float64 {
	if len(optimizations) == 0 {
		return 0.0
	}
	
	// Base confidence on number of historical executions
	historyConfidence := math.Min(float64(len(context.History))/10.0, 1.0)
	
	// Factor in system resource availability
	resourceConfidence := (100 - context.SystemResources.CPUUsage) / 100.0
	
	// Factor in optimization impact
	totalImpact := 0.0
	for _, opt := range optimizations {
		totalImpact += opt.Impact
	}
	impactConfidence := math.Min(totalImpact/float64(len(optimizations)), 1.0)
	
	// Weighted average
	return (historyConfidence*0.4 + resourceConfidence*0.3 + impactConfidence*0.3)
}

// calculateExpectedSpeedup calculates expected speedup from optimizations
func (to *TaskOptimizer) calculateExpectedSpeedup(context OptimizationContext, optimizations []Optimization) float64 {
	if len(optimizations) == 0 {
		return 1.0
	}
	
	totalSpeedup := 1.0
	for _, opt := range optimizations {
		// Convert impact to speedup (higher impact = more speedup)
		speedup := 1.0 + (opt.Impact * 0.5) // Max 50% speedup per optimization
		totalSpeedup *= speedup
	}
	
	// Cap maximum speedup at 5x
	return math.Min(totalSpeedup, 5.0)
}

// calculateResourceSavings calculates expected resource savings
func (to *TaskOptimizer) calculateResourceSavings(optimizations []Optimization) map[string]float64 {
	savings := map[string]float64{
		"cpu":     0.0,
		"memory":  0.0,
		"network": 0.0,
		"disk":    0.0,
	}
	
	for _, opt := range optimizations {
		switch opt.Type {
		case "parallelization":
			savings["cpu"] += opt.Impact * 0.3
		case "memory":
			savings["memory"] += opt.Impact * 0.4
		case "network":
			savings["network"] += opt.Impact * 0.5
		case "io":
			savings["disk"] += opt.Impact * 0.3
		case "caching":
			savings["network"] += opt.Impact * 0.6
			savings["disk"] += opt.Impact * 0.2
		}
	}
	
	// Cap savings at reasonable levels
	for key, value := range savings {
		savings[key] = math.Min(value, 0.8) // Max 80% savings
	}
	
	return savings
}

// generateRationale generates human-readable rationale for optimizations
func (to *TaskOptimizer) generateRationale(optimizations []Optimization) string {
	if len(optimizations) == 0 {
		return "No optimizations suggested for this command."
	}
	
	rationale := fmt.Sprintf("Applied %d optimizations: ", len(optimizations))
	
	types := make(map[string]int)
	for _, opt := range optimizations {
		types[opt.Type]++
	}
	
	explanations := make([]string, 0)
	for optType, count := range types {
		switch optType {
		case "parallelization":
			explanations = append(explanations, fmt.Sprintf("%d parallelization improvements", count))
		case "memory":
			explanations = append(explanations, fmt.Sprintf("%d memory optimizations", count))
		case "compiler":
			explanations = append(explanations, fmt.Sprintf("%d compiler optimizations", count))
		case "caching":
			explanations = append(explanations, fmt.Sprintf("%d caching improvements", count))
		case "network":
			explanations = append(explanations, fmt.Sprintf("%d network optimizations", count))
		case "io":
			explanations = append(explanations, fmt.Sprintf("%d I/O optimizations", count))
		}
	}
	
	return rationale + strings.Join(explanations, ", ")
}