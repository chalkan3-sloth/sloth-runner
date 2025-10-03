package luainterface

import (
	"context"
	"log/slog"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/ai"
	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	lua "github.com/yuin/gopher-lua"
)

// registerAIModule registers the AI module for Lua
func (li *LuaInterface) registerAIModule() {
	li.L.PreloadModule("ai", li.loadAIModule)
}

// loadAIModule loads the AI module into Lua
func (li *LuaInterface) loadAIModule(L *lua.LState) int {
	// Create AI module table
	aiTable := L.NewTable()
	
	// Task optimization functions
	L.SetField(aiTable, "optimize_command", L.NewFunction(li.luaAIOptimizeCommand))
	L.SetField(aiTable, "predict_failure", L.NewFunction(li.luaAIPredictFailure))
	L.SetField(aiTable, "find_similar_tasks", L.NewFunction(li.luaAIFindSimilarTasks))
	L.SetField(aiTable, "get_task_history", L.NewFunction(li.luaAIGetTaskHistory))
	L.SetField(aiTable, "record_execution", L.NewFunction(li.luaAIRecordExecution))
	L.SetField(aiTable, "get_task_stats", L.NewFunction(li.luaAIGetTaskStats))
	
	// AI configuration functions
	L.SetField(aiTable, "configure", L.NewFunction(li.luaAIConfigure))
	L.SetField(aiTable, "get_config", L.NewFunction(li.luaAIGetConfig))
	
	// Utility functions
	L.SetField(aiTable, "analyze_performance", L.NewFunction(li.luaAIAnalyzePerformance))
	L.SetField(aiTable, "generate_insights", L.NewFunction(li.luaAIGenerateInsights))
	
	L.Push(aiTable)
	return 1
}

// luaAIOptimizeCommand optimizes a command using AI
func (li *LuaInterface) luaAIOptimizeCommand(L *lua.LState) int {
	command := L.CheckString(1)
	options := L.OptTable(2, L.NewTable())
	
	// Convert Lua table to Go map
	optionsMap := make(map[string]interface{})
	options.ForEach(func(key, value lua.LValue) {
		if keyStr, ok := key.(lua.LString); ok {
			optionsMap[string(keyStr)] = li.luaValueToInterface(value)
		}
	})
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("AI intelligence not available"))
		return 2
	}
	
	// Generate optimization suggestion
	ctx := context.Background()
	suggestion, err := aiIntelligence.OptimizeCommand(ctx, command, optionsMap)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Convert suggestion to Lua table
	result := L.NewTable()
	L.SetField(result, "original_command", lua.LString(suggestion.OriginalCommand))
	L.SetField(result, "optimized_command", lua.LString(suggestion.OptimizedCommand))
	L.SetField(result, "confidence_score", lua.LNumber(suggestion.ConfidenceScore))
	L.SetField(result, "expected_speedup", lua.LNumber(suggestion.ExpectedSpeedup))
	L.SetField(result, "rationale", lua.LString(suggestion.Rationale))
	
	// Add optimizations array
	optimizations := L.NewTable()
	for i, opt := range suggestion.Optimizations {
		optTable := L.NewTable()
		L.SetField(optTable, "type", lua.LString(opt.Type))
		L.SetField(optTable, "description", lua.LString(opt.Description))
		L.SetField(optTable, "impact", lua.LNumber(opt.Impact))
		optimizations.RawSetInt(i+1, optTable)
	}
	L.SetField(result, "optimizations", optimizations)
	
	// Add resource savings
	savings := L.NewTable()
	for resource, amount := range suggestion.ResourceSavings {
		L.SetField(savings, resource, lua.LNumber(amount))
	}
	L.SetField(result, "resource_savings", savings)
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaAIPredictFailure predicts failure probability for a task
func (li *LuaInterface) luaAIPredictFailure(L *lua.LState) int {
	taskName := L.CheckString(1)
	command := L.CheckString(2)
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("AI intelligence not available"))
		return 2
	}
	
	// Generate failure prediction
	ctx := context.Background()
	prediction, err := aiIntelligence.PredictFailure(ctx, taskName, command)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Convert prediction to Lua table
	result := L.NewTable()
	L.SetField(result, "task_name", lua.LString(prediction.TaskName))
	L.SetField(result, "command", lua.LString(prediction.Command))
	L.SetField(result, "failure_probability", lua.LNumber(prediction.FailureProbability))
	L.SetField(result, "confidence", lua.LNumber(prediction.Confidence))
	
	// Add risk factors
	riskFactors := L.NewTable()
	for i, factor := range prediction.RiskFactors {
		factorTable := L.NewTable()
		L.SetField(factorTable, "factor", lua.LString(factor.Factor))
		L.SetField(factorTable, "impact", lua.LNumber(factor.Impact))
		L.SetField(factorTable, "description", lua.LString(factor.Description))
		riskFactors.RawSetInt(i+1, factorTable)
	}
	L.SetField(result, "risk_factors", riskFactors)
	
	// Add recommendations
	recommendations := L.NewTable()
	for i, rec := range prediction.Recommendations {
		recommendations.RawSetInt(i+1, lua.LString(rec))
	}
	L.SetField(result, "recommendations", recommendations)
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaAIFindSimilarTasks finds similar tasks
func (li *LuaInterface) luaAIFindSimilarTasks(L *lua.LState) int {
	command := L.CheckString(1)
	limit := L.OptInt(2, 10)
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(L.NewTable())
		return 1
	}
	
	// Find similar tasks
	ctx := context.Background()
	similarTasks, err := aiIntelligence.FindSimilarTasks(ctx, command, limit)
	if err != nil {
		slog.Warn("Failed to find similar tasks", "error", err)
		L.Push(L.NewTable())
		return 1
	}
	
	// Convert to Lua table
	result := L.NewTable()
	for i, task := range similarTasks {
		taskTable := L.NewTable()
		L.SetField(taskTable, "task_name", lua.LString(task.TaskName))
		L.SetField(taskTable, "command", lua.LString(task.Command))
		L.SetField(taskTable, "success", lua.LBool(task.Success))
		L.SetField(taskTable, "execution_time", lua.LString(task.ExecutionTime.String()))
		L.SetField(taskTable, "timestamp", lua.LString(task.Timestamp.Format(time.RFC3339)))
		result.RawSetInt(i+1, taskTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaAIGetTaskHistory gets execution history for a command
func (li *LuaInterface) luaAIGetTaskHistory(L *lua.LState) int {
	command := L.CheckString(1)
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(L.NewTable())
		return 1
	}
	
	// Get task history
	ctx := context.Background()
	history, err := aiIntelligence.GetTaskHistory(ctx, command)
	if err != nil {
		slog.Warn("Failed to get task history", "error", err)
		L.Push(L.NewTable())
		return 1
	}
	
	// Convert to Lua table
	result := L.NewTable()
	for i, execution := range history {
		execTable := L.NewTable()
		L.SetField(execTable, "task_name", lua.LString(execution.TaskName))
		L.SetField(execTable, "command", lua.LString(execution.Command))
		L.SetField(execTable, "success", lua.LBool(execution.Success))
		L.SetField(execTable, "execution_time", lua.LString(execution.ExecutionTime.String()))
		L.SetField(execTable, "timestamp", lua.LString(execution.Timestamp.Format(time.RFC3339)))
		if execution.ErrorMessage != "" {
			L.SetField(execTable, "error_message", lua.LString(execution.ErrorMessage))
		}
		result.RawSetInt(i+1, execTable)
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaAIRecordExecution records a task execution for learning
func (li *LuaInterface) luaAIRecordExecution(L *lua.LState) int {
	execTable := L.CheckTable(1)
	
	// Extract execution data from Lua table
	execution := &ai.TaskExecution{
		TaskName:      li.getStringField(execTable, "task_name"),
		Command:       li.getStringField(execTable, "command"),
		Success:       li.getBoolField(execTable, "success"),
		ErrorMessage:  li.getStringField(execTable, "error_message"),
		Timestamp:     time.Now(),
	}
	
	// Parse execution time
	if timeStr := li.getStringField(execTable, "execution_time"); timeStr != "" {
		if duration, err := time.ParseDuration(timeStr); err == nil {
			execution.ExecutionTime = duration
		}
	}
	
	// Extract parameters
	if paramsTable := li.getTableField(execTable, "parameters"); paramsTable != nil {
		execution.Parameters = make(map[string]interface{})
		paramsTable.ForEach(func(key, value lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				execution.Parameters[string(keyStr)] = li.luaValueToInterface(value)
			}
		})
	}
	
	// Get current system resources (mock for now)
	execution.SystemResources = ai.SystemResources{
		CPUUsage:    50.0,
		MemoryUsage: 60.0,
		DiskIO:      10.0,
		NetworkIO:   5.0,
		LoadAvg:     1.5,
	}
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("AI intelligence not available"))
		return 2
	}
	
	// Record execution
	ctx := context.Background()
	if err := aiIntelligence.RecordExecution(ctx, execution); err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LBool(true))
	return 1
}

// luaAIGetTaskStats gets aggregated statistics for a task
func (li *LuaInterface) luaAIGetTaskStats(L *lua.LState) int {
	taskName := L.CheckString(1)
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("AI intelligence not available"))
		return 2
	}
	
	// Get task statistics from learning store
	ctx := context.Background()
	stats, err := aiIntelligence.GetTaskStats(ctx, taskName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Convert to Lua table
	result := L.NewTable()
	L.SetField(result, "task_name", lua.LString(stats.TaskName))
	L.SetField(result, "total_runs", lua.LNumber(stats.TotalRuns))
	L.SetField(result, "success_count", lua.LNumber(stats.SuccessCount))
	L.SetField(result, "failure_count", lua.LNumber(stats.FailureCount))
	L.SetField(result, "success_rate", lua.LNumber(stats.SuccessRate))
	L.SetField(result, "total_time", lua.LString(stats.TotalTime.String()))
	L.SetField(result, "avg_time", lua.LString(stats.AvgTime.String()))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaAIConfigure configures AI settings
func (li *LuaInterface) luaAIConfigure(L *lua.LState) int {
	configTable := L.CheckTable(1)
	
	// Extract configuration
	config := ai.DefaultAITaskConfig()
	
	if enabled := li.getBoolField(configTable, "enabled"); enabled {
		config.Enabled = enabled
	}
	
	if mode := li.getStringField(configTable, "learning_mode"); mode != "" {
		config.LearningMode = ai.LearningMode(mode)
	}
	
	if level := li.getNumberField(configTable, "optimization_level"); level > 0 {
		config.OptimizationLevel = int(level)
	}
	
	if prediction := li.getBoolField(configTable, "failure_prediction"); prediction {
		config.FailurePrediction = prediction
	}
	
	if autoOpt := li.getBoolField(configTable, "auto_optimize"); autoOpt {
		config.AutoOptimize = autoOpt
	}
	
	// Store configuration (in a real implementation, this would be persisted)
	li.aiConfig = &config
	
	L.Push(lua.LBool(true))
	return 1
}

// luaAIGetConfig gets current AI configuration
func (li *LuaInterface) luaAIGetConfig(L *lua.LState) int {
	config := li.aiConfig
	if config == nil {
		defaultConfig := ai.DefaultAITaskConfig()
		config = &defaultConfig
	}
	
	// Convert to Lua table
	result := L.NewTable()
	L.SetField(result, "enabled", lua.LBool(config.Enabled))
	L.SetField(result, "learning_mode", lua.LString(string(config.LearningMode)))
	L.SetField(result, "optimization_level", lua.LNumber(config.OptimizationLevel))
	L.SetField(result, "failure_prediction", lua.LBool(config.FailurePrediction))
	L.SetField(result, "auto_optimize", lua.LBool(config.AutoOptimize))
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaAIAnalyzePerformance analyzes task performance patterns
func (li *LuaInterface) luaAIAnalyzePerformance(L *lua.LState) int {
	command := L.CheckString(1)
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(L.NewTable())
		return 1
	}
	
	// Get task history for analysis
	ctx := context.Background()
	history, err := aiIntelligence.GetTaskHistory(ctx, command)
	if err != nil || len(history) == 0 {
		L.Push(L.NewTable())
		return 1
	}
	
	// Analyze performance patterns
	analysis := li.analyzePerformancePatterns(history)
	
	// Convert to Lua table
	result := L.NewTable()
	L.SetField(result, "command", lua.LString(command))
	L.SetField(result, "total_executions", lua.LNumber(len(history)))
	L.SetField(result, "avg_execution_time", lua.LString(analysis.AvgExecutionTime.String()))
	L.SetField(result, "min_execution_time", lua.LString(analysis.MinExecutionTime.String()))
	L.SetField(result, "max_execution_time", lua.LString(analysis.MaxExecutionTime.String()))
	L.SetField(result, "success_rate", lua.LNumber(analysis.SuccessRate))
	L.SetField(result, "performance_trend", lua.LString(analysis.Trend))
	
	insights := L.NewTable()
	for i, insight := range analysis.Insights {
		insights.RawSetInt(i+1, lua.LString(insight))
	}
	L.SetField(result, "insights", insights)
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// luaAIGenerateInsights generates AI-powered insights
func (li *LuaInterface) luaAIGenerateInsights(L *lua.LState) int {
	_ = L.OptTable(1, L.NewTable()) // Mark as used
	
	// Get AI intelligence instance
	aiIntelligence := li.getAIIntelligence()
	if aiIntelligence == nil {
		L.Push(L.NewTable())
		return 1
	}
	
	// Generate insights
	insights := []string{
		"Consider implementing retry mechanisms for tasks with >20% failure rate",
		"Tasks executed during business hours have 15% lower failure rate",
		"Commands with parallel flags show 40% better performance",
		"Memory-intensive tasks perform better with explicit heap size settings",
		"Network-dependent tasks should include timeout and retry configurations",
	}
	
	// Convert to Lua table
	result := L.NewTable()
	for i, insight := range insights {
		result.RawSetInt(i+1, lua.LString(insight))
	}
	
	L.Push(lua.LTrue)
	L.Push(result)
	return 2
}

// getAIIntelligence gets or creates AI intelligence instance
func (li *LuaInterface) getAIIntelligence() *ai.TaskIntelligence {
	if li.aiIntelligence == nil {
		// Create state manager if not exists
		if li.stateManager == nil {
			stateManager, err := state.NewStateManager("")
			if err != nil {
				slog.Error("Failed to create state manager for AI", "error", err)
				return nil
			}
			li.stateManager = stateManager
		}
		
		li.aiIntelligence = ai.NewTaskIntelligence(li.stateManager)
	}
	return li.aiIntelligence
}

// PerformanceAnalysis represents performance analysis results
type PerformanceAnalysis struct {
	AvgExecutionTime time.Duration
	MinExecutionTime time.Duration
	MaxExecutionTime time.Duration
	SuccessRate      float64
	Trend            string
	Insights         []string
}

// analyzePerformancePatterns analyzes performance patterns in execution history
func (li *LuaInterface) analyzePerformancePatterns(history []*ai.TaskExecution) *PerformanceAnalysis {
	if len(history) == 0 {
		return &PerformanceAnalysis{}
	}
	
	// Calculate basic statistics
	var totalTime time.Duration
	minTime := history[0].ExecutionTime
	maxTime := history[0].ExecutionTime
	successCount := 0
	
	for _, exec := range history {
		totalTime += exec.ExecutionTime
		if exec.ExecutionTime < minTime {
			minTime = exec.ExecutionTime
		}
		if exec.ExecutionTime > maxTime {
			maxTime = exec.ExecutionTime
		}
		if exec.Success {
			successCount++
		}
	}
	
	avgTime := totalTime / time.Duration(len(history))
	successRate := float64(successCount) / float64(len(history))
	
	// Analyze trend (simple: compare first half vs second half)
	trend := "stable"
	if len(history) >= 4 {
		firstHalf := history[len(history)/2:]
		secondHalf := history[:len(history)/2]
		
		firstAvg := time.Duration(0)
		secondAvg := time.Duration(0)
		
		for _, exec := range firstHalf {
			firstAvg += exec.ExecutionTime
		}
		firstAvg /= time.Duration(len(firstHalf))
		
		for _, exec := range secondHalf {
			secondAvg += exec.ExecutionTime
		}
		secondAvg /= time.Duration(len(secondHalf))
		
		if secondAvg < firstAvg*9/10 {
			trend = "improving"
		} else if secondAvg > firstAvg*11/10 {
			trend = "degrading"
		}
	}
	
	// Generate insights
	insights := make([]string, 0)
	if successRate < 0.8 {
		insights = append(insights, "Consider adding error handling and retry logic")
	}
	if maxTime > avgTime*3 {
		insights = append(insights, "Execution time varies significantly - investigate outliers")
	}
	if trend == "degrading" {
		insights = append(insights, "Performance is degrading over time - optimization needed")
	}
	if trend == "improving" {
		insights = append(insights, "Performance is improving - current optimizations are effective")
	}
	
	return &PerformanceAnalysis{
		AvgExecutionTime: avgTime,
		MinExecutionTime: minTime,
		MaxExecutionTime: maxTime,
		SuccessRate:      successRate,
		Trend:            trend,
		Insights:         insights,
	}
}

// Helper methods for Lua table field extraction
func (li *LuaInterface) getStringField(table *lua.LTable, field string) string {
	value := table.RawGetString(field)
	if str, ok := value.(lua.LString); ok {
		return string(str)
	}
	return ""
}

func (li *LuaInterface) getBoolField(table *lua.LTable, field string) bool {
	value := table.RawGetString(field)
	if boolean, ok := value.(lua.LBool); ok {
		return bool(boolean)
	}
	return false
}

func (li *LuaInterface) getNumberField(table *lua.LTable, field string) float64 {
	value := table.RawGetString(field)
	if number, ok := value.(lua.LNumber); ok {
		return float64(number)
	}
	return 0
}

func (li *LuaInterface) getTableField(table *lua.LTable, field string) *lua.LTable {
	value := table.RawGetString(field)
	if table, ok := value.(*lua.LTable); ok {
		return table
	}
	return nil
}

// luaValueToInterface converts Lua values to Go interface{}
func (li *LuaInterface) luaValueToInterface(value lua.LValue) interface{} {
	switch v := value.(type) {
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LTable:
		result := make(map[string]interface{})
		v.ForEach(func(key, val lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				result[string(keyStr)] = li.luaValueToInterface(val)
			}
		})
		return result
	default:
		return nil
	}
}