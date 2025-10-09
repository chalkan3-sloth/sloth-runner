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

type GroupTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Rules       []GroupRule       `json:"rules"`
	Tags        map[string]string `json:"tags,omitempty"`
	CreatedAt   int64             `json:"created_at"`
}

type GroupRule struct {
	Type     string `json:"type"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

func NewTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage group templates",
		Long:  `Create, list, apply, and delete group templates.`,
	}

	cmd.AddCommand(newTemplateListCmd())
	cmd.AddCommand(newTemplateCreateCmd())
	cmd.AddCommand(newTemplateApplyCmd())
	cmd.AddCommand(newTemplateDeleteCmd())

	return cmd
}

func newTemplateListCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all group templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTemplates(outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	return cmd
}

func newTemplateCreateCmd() *cobra.Command {
	var description string
	var rules []string
	var tags []string

	cmd := &cobra.Command{
		Use:   "create <template-name>",
		Short: "Create a new group template",
		Long: `Create a new group template with rules for automatic agent matching.

Rules format: type:operator:value
  Types: tag_match, name_pattern, status
  Operators: equals, contains, regex

Examples:
  tag_match:equals:production
  name_pattern:contains:web
  status:equals:active`,
		Example: `  # Create template with tag match rule
  sloth-runner group template create web-servers \
    --description "Web server template" \
    --rule "tag_match:equals:web" \
    --tag "env:production"

  # Create template with multiple rules
  sloth-runner group template create prod-db \
    --rule "tag_match:equals:database" \
    --rule "name_pattern:contains:prod"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTemplate(args[0], description, rules, tags)
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "", "Template description")
	cmd.Flags().StringArrayVarP(&rules, "rule", "r", []string{}, "Template rules (type:operator:value)")
	cmd.Flags().StringArrayVarP(&tags, "tag", "t", []string{}, "Template tags (key:value)")

	return cmd
}

func newTemplateApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply <template-id> <group-name>",
		Short: "Apply a template to create/update a group",
		Long:  `Apply a template's rules to create a new group or update an existing one.`,
		Example: `  # Apply template to create a new group
  sloth-runner group template apply web-servers production-web`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyTemplate(args[0], args[1])
		},
	}

	return cmd
}

func newTemplateDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <template-id>",
		Short: "Delete a group template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Print("Are you sure you want to delete this template? (yes/no): ")
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "yes" && confirm != "y" {
					fmt.Println("Cancelled")
					return nil
				}
			}
			return deleteTemplate(args[0])
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	return cmd
}

func listTemplates(outputFormat string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/agent-groups/templates", apiURL))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Templates []GroupTemplate `json:"templates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	templates := response.Templates

	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(templates)
	}

	// Text format
	if len(templates) == 0 {
		fmt.Println("No templates found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tRULES\tDESCRIPTION")
	fmt.Fprintln(w, "--\t----\t-----\t-----------")

	for _, t := range templates {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			t.ID,
			t.Name,
			len(t.Rules),
			t.Description,
		)
	}

	w.Flush()
	return nil
}

func createTemplate(name, description string, rules, tags []string) error {
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
		"name":        name,
		"description": description,
		"rules":       parsedRules,
		"tags":        parsedTags,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/agent-groups/templates", apiURL),
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

	var result struct {
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		TemplateID string `json:"template_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("✅ Template '%s' created successfully (ID: %s)\n", name, result.TemplateID)
	return nil
}

func applyTemplate(templateID, groupName string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	payload := map[string]interface{}{
		"group_id": groupName,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/agent-groups/templates/%s/apply", apiURL, templateID),
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

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("✅ Template applied successfully to group '%s'\n", groupName)
	if matchedCount, ok := result["matched_agents"].(float64); ok {
		fmt.Printf("   Matched agents: %.0f\n", matchedCount)
	}
	return nil
}

func deleteTemplate(templateID string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/agent-groups/templates/%s", apiURL, templateID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✅ Template '%s' deleted successfully\n", templateID)
	return nil
}

func parseRule(rule string) []string {
	var parts []string
	current := ""
	for _, c := range rule {
		if c == ':' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func parseTag(tag string) []string {
	var parts []string
	current := ""
	for _, c := range tag {
		if c == ':' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
