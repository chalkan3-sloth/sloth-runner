//go:build cgo
// +build cgo

package services

import (
	"fmt"
	"time"
)

// SchedulerStackService handles scheduler operations tracking
type SchedulerStackService struct {
	stackService *StackService
}

// NewSchedulerStackService creates a new scheduler stack service
func NewSchedulerStackService(stackService *StackService) *SchedulerStackService {
	return &SchedulerStackService{
		stackService: stackService,
	}
}

// TrackScheduledExecution tracks a scheduled execution
func (s *SchedulerStackService) TrackScheduledExecution(scheduleName, schedule, executionID string, success bool, duration time.Duration, errorMsg string) error {
	stackID := fmt.Sprintf("schedule-%s", scheduleName)

	// Create or get stack
	_, err := s.stackService.GetOrCreateStack(stackID, "scheduled-execution", "")
	if err != nil {
		return fmt.Errorf("failed to create scheduler stack: %w", err)
	}

	stack, err := s.stackService.GetStackByName(stackID)
	if err != nil {
		return err
	}

	// Initialize executions list if not exists
	if stack.Metadata == nil {
		stack.Metadata = make(map[string]interface{})
	}

	executions, _ := stack.Metadata["executions"].([]interface{})
	executions = append(executions, map[string]interface{}{
		"execution_id": executionID,
		"success":      success,
		"duration":     duration.String(),
		"error":        errorMsg,
		"timestamp":    time.Now(),
	})

	stack.Metadata["executions"] = executions
	stack.Metadata["schedule"] = schedule
	stack.Metadata["last_execution"] = time.Now()

	return s.stackService.GetManager().UpdateStack(stack)
}

// TrackScheduleChange tracks a schedule change
func (s *SchedulerStackService) TrackScheduleChange(scheduleName, newSchedule, action string) error {
	stackID := fmt.Sprintf("schedule-%s", scheduleName)

	_, err := s.stackService.GetOrCreateStack(stackID, "schedule-change", "")
	if err != nil {
		return fmt.Errorf("failed to create scheduler stack: %w", err)
	}

	stack, err := s.stackService.GetStackByName(stackID)
	if err != nil {
		return err
	}

	if stack.Metadata == nil {
		stack.Metadata = make(map[string]interface{})
	}

	changes, _ := stack.Metadata["changes"].([]interface{})
	changes = append(changes, map[string]interface{}{
		"new_schedule": newSchedule,
		"action":       action,
		"timestamp":    time.Now(),
	})

	stack.Metadata["changes"] = changes
	stack.Metadata["current_schedule"] = newSchedule

	return s.stackService.GetManager().UpdateStack(stack)
}

// GetSchedulerHistory retrieves scheduler history
func (s *SchedulerStackService) GetSchedulerHistory(limit int) ([]map[string]interface{}, error) {
	stacks, err := s.stackService.ListStacks()
	if err != nil {
		return nil, err
	}

	history := make([]map[string]interface{}, 0)
	count := 0

	for _, st := range stacks {
		if count >= limit {
			break
		}

		if st.Metadata != nil {
			// Check if stack has executions or changes
			if executions, ok := st.Metadata["executions"]; ok {
				if execs, ok := executions.([]interface{}); ok {
					for _, exec := range execs {
						if execMap, ok := exec.(map[string]interface{}); ok {
							history = append(history, execMap)
							count++
							if count >= limit {
								break
							}
						}
					}
				}
			}

			if changes, ok := st.Metadata["changes"]; ok {
				if chgs, ok := changes.([]interface{}); ok {
					for _, chg := range chgs {
						if chgMap, ok := chg.(map[string]interface{}); ok {
							history = append(history, chgMap)
							count++
							if count >= limit {
								break
							}
						}
					}
				}
			}
		}
	}

	return history, nil
}

// GetScheduleExecutionStats retrieves execution statistics
func (s *SchedulerStackService) GetScheduleExecutionStats() (map[string]interface{}, error) {
	stacks, err := s.stackService.ListStacks()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_executions": 0,
		"successful":       0,
		"failed":           0,
	}

	totalExecs := 0
	successfulExecs := 0
	failedExecs := 0

	for _, st := range stacks {
		if st.Metadata != nil {
			if executions, ok := st.Metadata["executions"]; ok {
				if execs, ok := executions.([]interface{}); ok {
					for _, exec := range execs {
						totalExecs++
						if execMap, ok := exec.(map[string]interface{}); ok {
							if success, ok := execMap["success"].(bool); ok {
								if success {
									successfulExecs++
								} else {
									failedExecs++
								}
							}
						}
					}
				}
			}
		}
	}

	stats["total_executions"] = totalExecs
	stats["successful"] = successfulExecs
	stats["failed"] = failedExecs

	return stats, nil
}
