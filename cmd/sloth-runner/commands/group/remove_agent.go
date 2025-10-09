package group

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func NewRemoveAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-agent <group-name> <agent-name> [agent-name...]",
		Short: "Remove agents from a group",
		Long:  `Remove one or more agents from an existing group.`,
		Example: `  # Remove single agent
  sloth-runner group remove-agent production-web server-01

  # Remove multiple agents
  sloth-runner group remove-agent production-web server-01 server-02`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupName := args[0]
			agentNames := args[1:]
			return removeAgents(groupName, agentNames)
		},
	}

	return cmd
}

func removeAgents(groupName string, agentNames []string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	payload := map[string]interface{}{
		"agent_names": agentNames,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/api/v1/agent-groups/%s/agents", apiURL, groupName),
		bytes.NewBuffer(jsonData),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✅ Removed %d agent(s) from group '%s'\n", len(agentNames), groupName)
	return nil
}
