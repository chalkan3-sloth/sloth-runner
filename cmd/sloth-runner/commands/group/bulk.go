package group

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type BulkOperationResult struct {
	GroupID      string                     `json:"group_id"`
	TotalAgents  int                        `json:"total_agents"`
	SuccessCount int                        `json:"success_count"`
	FailureCount int                        `json:"failure_count"`
	Results      map[string]BulkAgentResult `json:"results"`
	DurationMs   int64                      `json:"duration_ms"`
}

type BulkAgentResult struct {
	AgentName  string `json:"agent_name"`
	Success    bool   `json:"success"`
	Output     string `json:"output,omitempty"`
	Error      string `json:"error,omitempty"`
	DurationMs int64  `json:"duration_ms"`
}

func NewBulkCmd() *cobra.Command {
	var command string
	var timeout int

	cmd := &cobra.Command{
		Use:   "bulk <group-name> <operation>",
		Short: "Execute bulk operation on all agents in a group",
		Long: `Execute an operation on all agents in a group simultaneously.

Operations:
  restart        - Restart all agents
  update         - Update all agents
  shutdown       - Shutdown all agents
  execute        - Execute a command (requires --command flag)`,
		Example: `  # Restart all agents in a group
  sloth-runner group bulk production-web restart

  # Execute command on all agents
  sloth-runner group bulk production-web execute --command "systemctl restart nginx"

  # Update all agents with custom timeout
  sloth-runner group bulk production-web update --timeout 600`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupName := args[0]
			operation := args[1]

			if operation == "execute" && command == "" {
				return fmt.Errorf("--command flag is required for 'execute' operation")
			}

			return executeBulkOperation(groupName, operation, command, timeout)
		},
	}

	cmd.Flags().StringVarP(&command, "command", "c", "", "Command to execute (required for 'execute' operation)")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 300, "Timeout in seconds")

	return cmd
}

func executeBulkOperation(groupName, operation, command string, timeout int) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	params := make(map[string]interface{})
	if operation == "execute" {
		params["command"] = command
	}

	// Translate CLI operation names to API operation names
	apiOperation := operation
	if operation == "execute" {
		apiOperation = "execute_command"
	}

	payload := map[string]interface{}{
		"group_id":  groupName,
		"operation": apiOperation,
		"params":    params,
		"timeout":   timeout,
	}

	jsonData, _ := json.Marshal(payload)

	fmt.Printf("⏳ Executing '%s' operation on group '%s'...\n", operation, groupName)

	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/agent-groups/bulk-operation", apiURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result BulkOperationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Print results
	successRate := float64(result.SuccessCount) / float64(result.TotalAgents) * 100

	fmt.Printf("\n✅ Bulk operation completed in %dms\n", result.DurationMs)
	fmt.Printf("Summary: %d/%d agents succeeded (%.1f%%)\n\n",
		result.SuccessCount, result.TotalAgents, successRate)

	// Table of results
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "AGENT\tSTATUS\tDURATION\tOUTPUT/ERROR")
	fmt.Fprintln(w, "-----\t------\t--------\t------------")

	for _, res := range result.Results {
		status := "✅ SUCCESS"
		output := res.Output
		if !res.Success {
			status = "❌ FAILED"
			output = res.Error
		}

		if len(output) > 60 {
			output = output[:57] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%dms\t%s\n",
			res.AgentName,
			status,
			res.DurationMs,
			output,
		)
	}

	w.Flush()

	return nil
}
