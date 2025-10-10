//go:build cgo
// +build cgo

package services

import (
	"fmt"
	"time"
)

// SysadminStackService handles sysadmin operations tracking
type SysadminStackService struct {
	stackService *StackService
}

// NewSysadminStackService creates a new sysadmin stack service
func NewSysadminStackService(stackService *StackService) *SysadminStackService {
	return &SysadminStackService{
		stackService: stackService,
	}
}

// TrackBackupOperation tracks a backup operation in the stack
func (s *SysadminStackService) TrackBackupOperation(backupID, description string, paths []string, sizeBytes int64) error {
	// Create a stack for this backup operation
	stackID := fmt.Sprintf("backup-%s", backupID)

	stack, err := s.stackService.GetStackByName(stackID)
	if err != nil {
		// Create new stack for backup
		metadata := map[string]interface{}{
			"operation_type": "backup",
			"backup_id":      backupID,
			"paths":          paths,
			"size_bytes":     sizeBytes,
			"description":    description,
			"timestamp":      time.Now(),
		}

		_, createErr := s.stackService.GetOrCreateStack(stackID, "backup-operation", "")
		if createErr != nil {
			return fmt.Errorf("failed to create backup stack: %w", createErr)
		}

		// Update the stack with metadata
		stack, err = s.stackService.GetStackByName(stackID)
		if err != nil {
			return err
		}
		stack.Metadata = metadata
		return s.stackService.GetManager().UpdateStack(stack)
	}

	return nil
}

// TrackDeploymentOperation tracks a deployment operation
func (s *SysadminStackService) TrackDeploymentOperation(version string, agents []string, strategy string, success bool, duration time.Duration) error {
	stackID := fmt.Sprintf("deployment-%s", version)

	metadata := map[string]interface{}{
		"operation_type": "deployment",
		"version":        version,
		"agents":         agents,
		"strategy":       strategy,
		"success":        success,
		"duration":       duration.String(),
		"timestamp":      time.Now(),
	}

	_, err := s.stackService.GetOrCreateStack(stackID, "deployment-operation", "")
	if err != nil {
		return fmt.Errorf("failed to create deployment stack: %w", err)
	}

	stack, err := s.stackService.GetStackByName(stackID)
	if err != nil {
		return err
	}
	stack.Metadata = metadata
	return s.stackService.GetManager().UpdateStack(stack)
}

// TrackMaintenanceOperation tracks a maintenance operation
func (s *SysadminStackService) TrackMaintenanceOperation(operationType, agent, description string, success bool) error {
	stackID := fmt.Sprintf("maintenance-%s-%d", operationType, time.Now().Unix())

	metadata := map[string]interface{}{
		"operation_type": "maintenance",
		"type":           operationType,
		"agent":          agent,
		"description":    description,
		"success":        success,
		"timestamp":      time.Now(),
	}

	_, err := s.stackService.GetOrCreateStack(stackID, "maintenance-operation", "")
	if err != nil {
		return fmt.Errorf("failed to create maintenance stack: %w", err)
	}

	stack, err := s.stackService.GetStackByName(stackID)
	if err != nil {
		return err
	}
	stack.Metadata = metadata
	return s.stackService.GetManager().UpdateStack(stack)
}

// GetSysadminOperationHistory retrieves operation history
func (s *SysadminStackService) GetSysadminOperationHistory() ([]map[string]interface{}, error) {
	stacks, err := s.stackService.ListStacks()
	if err != nil {
		return nil, err
	}

	operations := make([]map[string]interface{}, 0)
	for _, st := range stacks {
		if st.Metadata != nil {
			if opType, ok := st.Metadata["operation_type"]; ok {
				if opType == "backup" || opType == "deployment" || opType == "maintenance" {
					operations = append(operations, st.Metadata)
				}
			}
		}
	}

	return operations, nil
}
