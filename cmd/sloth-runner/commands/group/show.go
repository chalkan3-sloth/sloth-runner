package group

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// NewShowCmd creates the show command
func NewShowCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "show <group-name>",
		Short: "Show details of an agent group",
		Long:  `Show detailed information about a specific agent group including members and metrics.`,
		Example: `  # Show group details
  sloth-runner group show production-web

  # Show in JSON format
  sloth-runner group show production-web --output json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupName := args[0]
			return showGroup(groupName, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")

	return cmd
}

func showGroup(groupName, outputFormat string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/agent-groups/%s", apiURL, groupName))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var group Group
	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(group)
	}

	// Text format
	fmt.Printf("Group: %s\n", group.Name)
	fmt.Printf("Description: %s\n", group.Description)
	fmt.Printf("Agent Count: %d\n", group.AgentCount)

	if len(group.Tags) > 0 {
		fmt.Println("\nTags:")
		for k, v := range group.Tags {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}

	if len(group.AgentNames) > 0 {
		fmt.Println("\nAgents:")
		for _, agent := range group.AgentNames {
			fmt.Printf("  • %s\n", agent)
		}
	} else {
		fmt.Println("\nNo agents in this group")
	}

	return nil
}
