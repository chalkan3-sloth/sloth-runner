package group

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	var description string
	var tags []string

	cmd := &cobra.Command{
		Use:   "create <group-name>",
		Short: "Create a new agent group",
		Long: `Create a new agent group with a name and optional description and tags.

Tags should be specified as key=value pairs.`,
		Example: `  # Create a basic group
  sloth-runner group create production-web

  # Create with description
  sloth-runner group create production-web --description "Production web servers"

  # Create with tags
  sloth-runner group create production-web --tag environment=production --tag role=webserver`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupName := args[0]
			return createGroup(groupName, description, tags)
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "", "Group description")
	cmd.Flags().StringArrayVarP(&tags, "tag", "t", []string{}, "Tags in format key=value")

	return cmd
}

func createGroup(name, description string, tagsList []string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	// Parse tags
	tags := make(map[string]string)
	for _, tagStr := range tagsList {
		parts := strings.SplitN(tagStr, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid tag format '%s', expected key=value", tagStr)
		}
		tags[parts[0]] = parts[1]
	}

	payload := map[string]interface{}{
		"group_name":  name,
		"description": description,
		"tags":        tags,
		"agent_names": []string{},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(apiURL+"/api/v1/agent-groups", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✅ Group '%s' created successfully\n", name)

	return nil
}
