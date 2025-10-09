package group

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type Group struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	AgentNames  []string          `json:"agent_names"`
	Tags        map[string]string `json:"tags"`
	AgentCount  int               `json:"agent_count"`
	CreatedAt   int64             `json:"created_at"`
}

type GroupsResponse struct {
	Groups []*Group `json:"groups"`
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all agent groups",
		Long: `List all agent groups with their details.

Shows group name, description, agent count, and tags.`,
		Example: `  # List all groups
  sloth-runner group list

  # List in JSON format
  sloth-runner group list --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listGroups(outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table|json)")

	return cmd
}

func listGroups(outputFormat string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	resp, err := http.Get(apiURL + "/api/v1/agent-groups")
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result GroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tAGENTS\tDESCRIPTION\tTAGS")
	fmt.Fprintln(w, "----\t------\t-----------\t----")

	for _, group := range result.Groups {
		tags := formatTags(group.Tags)
		desc := group.Description
		if desc == "" {
			desc = "-"
		}
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}

		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
			group.Name,
			group.AgentCount,
			desc,
			tags,
		)
	}

	w.Flush()

	if len(result.Groups) == 0 {
		fmt.Println("\nNo groups found. Create one with 'sloth-runner group create'")
	}

	return nil
}

func formatTags(tags map[string]string) string {
	if len(tags) == 0 {
		return "-"
	}

	result := ""
	count := 0
	for k, v := range tags {
		if count > 0 {
			result += ","
		}
		result += fmt.Sprintf("%s=%s", k, v)
		count++
		if count >= 2 {
			remaining := len(tags) - count
			if remaining > 0 {
				result += fmt.Sprintf("... +%d more", remaining)
			}
			break
		}
	}
	return result
}
