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

func NewAddAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-agent <group-name> <agent-name> [agent-name...]",
		Short: "Add agents to a group",
		Long:  `Add one or more agents to an existing group.`,
		Example: `  # Add single agent
  sloth-runner group add-agent production-web server-01

  # Add multiple agents
  sloth-runner group add-agent production-web server-01 server-02 server-03`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupName := args[0]
			agentNames := args[1:]
			return addAgents(groupName, agentNames)
		},
	}

	return cmd
}

func addAgents(groupName string, agentNames []string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	payload := map[string]interface{}{
		"agent_names": agentNames,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/agent-groups/%s/agents", apiURL, groupName),
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

	fmt.Printf("✅ Added %d agent(s) to group '%s'\n", len(agentNames), groupName)
	return nil
}
