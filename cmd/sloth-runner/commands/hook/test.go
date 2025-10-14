//go:build cgo
// +build cgo

package hook

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/internal/hooks"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewTestCommand creates the hook test command
func NewTestCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <hook-name>",
		Short: "Test a hook with mock event data",
		Long:  `Test an event hook by executing it with mock event data.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]
			eventData, _ := cmd.Flags().GetString("data")

			// Create repository
			repo, err := hooks.NewRepository()
			if err != nil {
				return fmt.Errorf("failed to create repository: %w", err)
			}
			defer repo.Close()

			// Get hook
			hook, err := repo.GetByName(hookName)
			if err != nil {
				return fmt.Errorf("failed to get hook: %w", err)
			}

			pterm.Info.Printf("Testing hook: %s\n", hook.Name)
			pterm.Info.Printf("Event type: %s\n", hook.EventType)

			// Parse event data if provided
			var data map[string]interface{}
			if eventData != "" {
				if err := json.Unmarshal([]byte(eventData), &data); err != nil {
					return fmt.Errorf("invalid event data JSON: %w", err)
				}
			} else {
				// Use default test data based on event type
				data = getTestEventData(hook.EventType)
			}

			// Create test event
			event := &hooks.Event{
				Type:      hook.EventType,
				Timestamp: time.Now(),
				Data:      data,
			}

			// Create executor and execute
			executor := hooks.NewExecutor(repo)
			result, err := executor.Execute(hook, event)

			// Display results
			fmt.Println()
			pterm.DefaultSection.Println("Test Results")

			if result != nil {
				if result.Success {
					pterm.Success.Println("Status: Success")
				} else {
					pterm.Error.Println("Status: Failed")
				}

				pterm.Info.Printf("Duration: %s\n", result.Duration)

				if result.Output != "" {
					fmt.Println()
					pterm.DefaultSection.Println("Output")
					fmt.Println(result.Output)
				}

				if result.Error != "" {
					fmt.Println()
					pterm.DefaultSection.Println("Error")
					pterm.Error.Println(result.Error)
				}
			}

			if err != nil {
				return fmt.Errorf("execution error: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("data", "d", "", "Custom event data as JSON")

	return cmd
}

// getTestEventData returns test data for different event types
func getTestEventData(eventType hooks.EventType) map[string]interface{} {
	switch eventType {
	case hooks.EventAgentRegistered:
		return map[string]interface{}{
			"agent": map[string]interface{}{
				"name":    "test-agent",
				"address": "192.168.1.100:50051",
				"tags":    []string{"test", "linux"},
				"version": "1.0.0",
				"system_info": map[string]interface{}{
					"os":     "Linux",
					"cpus":   4,
					"memory": 8192,
				},
			},
		}
	case hooks.EventAgentDisconnected:
		return map[string]interface{}{
			"agent": map[string]interface{}{
				"name":    "test-agent",
				"address": "192.168.1.100:50051",
			},
		}
	case hooks.EventTaskCompleted:
		return map[string]interface{}{
			"task": map[string]interface{}{
				"task_name":  "test-task",
				"agent_name": "test-agent",
				"status":     "completed",
				"exit_code":  0,
				"duration":   "5s",
			},
		}
	default:
		return map[string]interface{}{
			"test": true,
			"timestamp": time.Now().Unix(),
		}
	}
}
