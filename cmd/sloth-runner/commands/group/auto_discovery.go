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

type AutoDiscoveryConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	GroupID     string            `json:"group_id"`
	Rules       []GroupRule       `json:"rules"`
	Schedule    string            `json:"schedule"`
	Enabled     bool              `json:"enabled"`
	LastRun     int64             `json:"last_run"`
	NextRun     int64             `json:"next_run"`
	Tags        map[string]string `json:"tags,omitempty"`
	CreatedAt   int64             `json:"created_at"`
}

func NewAutoDiscoveryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auto-discovery",
		Aliases: []string{"ad", "discovery"},
		Short:   "Manage auto-discovery configurations",
		Long:    `Create, list, run, and delete auto-discovery configurations for automatic group population.`,
	}

	cmd.AddCommand(newAutoDiscoveryListCmd())
	cmd.AddCommand(newAutoDiscoveryCreateCmd())
	cmd.AddCommand(newAutoDiscoveryRunCmd())
	cmd.AddCommand(newAutoDiscoveryEnableCmd())
	cmd.AddCommand(newAutoDiscoveryDisableCmd())
	cmd.AddCommand(newAutoDiscoveryDeleteCmd())

	return cmd
}

func newAutoDiscoveryListCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all auto-discovery configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listAutoDiscovery(outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	return cmd
}

func newAutoDiscoveryCreateCmd() *cobra.Command {
	var groupID string
	var schedule string
	var rules []string
	var tags []string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "create <config-name>",
		Short: "Create a new auto-discovery configuration",
		Long: `Create a new auto-discovery configuration for automatic group population.

Schedule format: cron expression (e.g., "*/5 * * * *" for every 5 minutes)

Rules format: type:operator:value
  Types: tag_match, name_pattern, status
  Operators: equals, contains, regex`,
		Example: `  # Create auto-discovery for web servers
  sloth-runner group auto-discovery create web-discovery \
    --group production-web \
    --schedule "*/10 * * * *" \
    --rule "tag_match:equals:web" \
    --enabled

  # Create auto-discovery with multiple rules
  sloth-runner group auto-discovery create db-discovery \
    --group production-db \
    --schedule "0 * * * *" \
    --rule "tag_match:equals:database" \
    --rule "name_pattern:contains:db"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAutoDiscovery(args[0], groupID, schedule, rules, tags, enabled)
		},
	}

	cmd.Flags().StringVarP(&groupID, "group", "g", "", "Target group ID (required)")
	cmd.Flags().StringVarP(&schedule, "schedule", "s", "", "Cron schedule (required)")
	cmd.Flags().StringArrayVarP(&rules, "rule", "r", []string{}, "Discovery rules (type:operator:value)")
	cmd.Flags().StringArrayVarP(&tags, "tag", "t", []string{}, "Tags (key:value)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "Enable immediately")

	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("schedule")

	return cmd
}

func newAutoDiscoveryRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <config-id>",
		Short: "Run auto-discovery manually",
		Long:  `Manually trigger an auto-discovery configuration to update the target group.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAutoDiscovery(args[0])
		},
	}

	return cmd
}

func newAutoDiscoveryEnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable <config-id>",
		Short: "Enable auto-discovery configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateAutoDiscoveryStatus(args[0], true)
		},
	}

	return cmd
}

func newAutoDiscoveryDisableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable <config-id>",
		Short: "Disable auto-discovery configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateAutoDiscoveryStatus(args[0], false)
		},
	}

	return cmd
}

func newAutoDiscoveryDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <config-id>",
		Short: "Delete auto-discovery configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Print("Are you sure you want to delete this configuration? (yes/no): ")
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "yes" && confirm != "y" {
					fmt.Println("Cancelled")
					return nil
				}
			}
			return deleteAutoDiscovery(args[0])
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	return cmd
}

func listAutoDiscovery(outputFormat string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/agent-groups/auto-discovery", apiURL))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Configs []AutoDiscoveryConfig `json:"configs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	configs := response.Configs

	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(configs)
	}

	// Text format
	if len(configs) == 0 {
		fmt.Println("No auto-discovery configurations found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tGROUP\tSCHEDULE\tENABLED\tRULES")
	fmt.Fprintln(w, "--\t----\t-----\t--------\t-------\t-----")

	for _, c := range configs {
		enabled := "No"
		if c.Enabled {
			enabled = "Yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n",
			c.ID,
			c.Name,
			c.GroupID,
			c.Schedule,
			enabled,
			len(c.Rules),
		)
	}

	w.Flush()
	return nil
}

func createAutoDiscovery(name, groupID, schedule string, rules, tags []string, enabled bool) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	// Parse rules
	var parsedRules []GroupRule
	for _, rule := range rules {
		parts := parseRule(rule)
		if len(parts) != 3 {
			return fmt.Errorf("invalid rule format: %s (expected type:operator:value)", rule)
		}
		parsedRules = append(parsedRules, GroupRule{
			Type:     parts[0],
			Operator: parts[1],
			Value:    parts[2],
		})
	}

	// Parse tags
	parsedTags := make(map[string]string)
	for _, tag := range tags {
		parts := parseTag(tag)
		if len(parts) != 2 {
			return fmt.Errorf("invalid tag format: %s (expected key:value)", tag)
		}
		parsedTags[parts[0]] = parts[1]
	}

	payload := map[string]interface{}{
		"name":     name,
		"group_id": groupID,
		"schedule": schedule,
		"rules":    parsedRules,
		"tags":     parsedTags,
		"enabled":  enabled,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/agent-groups/auto-discovery", apiURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result AutoDiscoveryConfig
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("✅ Auto-discovery configuration '%s' created successfully (ID: %s)\n", name, result.ID)
	return nil
}

func runAutoDiscovery(configID string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/agent-groups/auto-discovery/%s/run", apiURL, configID),
		"application/json",
		bytes.NewBuffer([]byte("{}")),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("✅ Auto-discovery run completed\n")
	if matchedCount, ok := result["matched_agents"].(float64); ok {
		fmt.Printf("   Matched agents: %.0f\n", matchedCount)
	}
	return nil
}

func updateAutoDiscoveryStatus(configID string, enabled bool) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	payload := map[string]interface{}{
		"enabled": enabled,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/api/v1/agent-groups/auto-discovery/%s", apiURL, configID),
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

	status := "enabled"
	if !enabled {
		status = "disabled"
	}
	fmt.Printf("✅ Auto-discovery configuration '%s' %s\n", configID, status)
	return nil
}

func deleteAutoDiscovery(configID string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/agent-groups/auto-discovery/%s", apiURL, configID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✅ Auto-discovery configuration '%s' deleted successfully\n", configID)
	return nil
}
