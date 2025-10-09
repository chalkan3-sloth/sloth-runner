package group

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

type Webhook struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Events      []string          `json:"events"`
	Secret      string            `json:"secret,omitempty"`
	Enabled     bool              `json:"enabled"`
	Headers     map[string]string `json:"headers,omitempty"`
	CreatedAt   int64             `json:"created_at"`
}

type WebhookLog struct {
	ID           string `json:"id"`
	WebhookID    string `json:"webhook_id"`
	Event        string `json:"event"`
	StatusCode   int    `json:"status_code"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
	Timestamp    int64  `json:"timestamp"`
}

func NewWebhookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Manage webhooks for group events",
		Long:  `Create, list, and delete webhooks for group change notifications.`,
	}

	cmd.AddCommand(newWebhookListCmd())
	cmd.AddCommand(newWebhookCreateCmd())
	cmd.AddCommand(newWebhookEnableCmd())
	cmd.AddCommand(newWebhookDisableCmd())
	cmd.AddCommand(newWebhookDeleteCmd())
	cmd.AddCommand(newWebhookLogsCmd())

	return cmd
}

func newWebhookListCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listWebhooks(outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	return cmd
}

func newWebhookCreateCmd() *cobra.Command {
	var url string
	var events []string
	var secret string
	var headers []string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "create <webhook-name>",
		Short: "Create a new webhook",
		Long: `Create a new webhook for group event notifications.

Available events:
  group.created         - New group created
  group.updated         - Group modified
  group.deleted         - Group deleted
  group.agent_added     - Agent added to group
  group.agent_removed   - Agent removed from group
  bulk.operation_start  - Bulk operation started
  bulk.operation_end    - Bulk operation completed`,
		Example: `  # Create webhook for all events
  sloth-runner group webhook create slack-notify \
    --url "https://hooks.slack.com/services/YOUR/WEBHOOK/URL" \
    --event "group.created" \
    --event "group.deleted" \
    --enabled

  # Create webhook with secret and headers
  sloth-runner group webhook create discord-webhook \
    --url "https://discord.com/api/webhooks/YOUR/WEBHOOK" \
    --event "bulk.operation_end" \
    --secret "my-secret-key" \
    --header "Content-Type:application/json"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createWebhook(args[0], url, events, secret, headers, enabled)
		},
	}

	cmd.Flags().StringVarP(&url, "url", "u", "", "Webhook URL (required)")
	cmd.Flags().StringArrayVarP(&events, "event", "e", []string{}, "Events to trigger webhook (required)")
	cmd.Flags().StringVarP(&secret, "secret", "s", "", "Secret for webhook signature")
	cmd.Flags().StringArrayVar(&headers, "header", []string{}, "Custom headers (key:value)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "Enable immediately")

	cmd.MarkFlagRequired("url")
	cmd.MarkFlagRequired("event")

	return cmd
}

func newWebhookEnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable <webhook-id>",
		Short: "Enable a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateWebhookStatus(args[0], true)
		},
	}

	return cmd
}

func newWebhookDisableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable <webhook-id>",
		Short: "Disable a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateWebhookStatus(args[0], false)
		},
	}

	return cmd
}

func newWebhookDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <webhook-id>",
		Short: "Delete a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Print("Are you sure you want to delete this webhook? (yes/no): ")
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "yes" && confirm != "y" {
					fmt.Println("Cancelled")
					return nil
				}
			}
			return deleteWebhook(args[0])
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	return cmd
}

func newWebhookLogsCmd() *cobra.Command {
	var webhookID string
	var limit int

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View webhook execution logs",
		Long:  `View recent webhook execution logs with status and error information.`,
		Example: `  # View recent logs for all webhooks
  sloth-runner group webhook logs

  # View logs for specific webhook
  sloth-runner group webhook logs --webhook slack-notify

  # View last 50 logs
  sloth-runner group webhook logs --limit 50`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewWebhookLogs(webhookID, limit)
		},
	}

	cmd.Flags().StringVarP(&webhookID, "webhook", "w", "", "Filter by webhook ID")
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Number of logs to display")

	return cmd
}

func listWebhooks(outputFormat string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/agent-groups/webhooks", apiURL))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Webhooks []Webhook `json:"webhooks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	webhooks := response.Webhooks

	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(webhooks)
	}

	// Text format
	if len(webhooks) == 0 {
		fmt.Println("No webhooks found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tURL\tEVENTS\tENABLED")
	fmt.Fprintln(w, "--\t----\t---\t------\t-------")

	for _, wh := range webhooks {
		enabled := "No"
		if wh.Enabled {
			enabled = "Yes"
		}
		url := wh.URL
		if len(url) > 40 {
			url = url[:37] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			wh.ID,
			wh.Name,
			url,
			len(wh.Events),
			enabled,
		)
	}

	w.Flush()
	return nil
}

func createWebhook(name, url string, events []string, secret string, headers []string, enabled bool) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	// Parse headers
	parsedHeaders := make(map[string]string)
	for _, header := range headers {
		parts := parseTag(header)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s (expected key:value)", header)
		}
		parsedHeaders[parts[0]] = parts[1]
	}

	payload := map[string]interface{}{
		"name":    name,
		"url":     url,
		"events":  events,
		"enabled": enabled,
	}

	if secret != "" {
		payload["secret"] = secret
	}
	if len(parsedHeaders) > 0 {
		payload["headers"] = parsedHeaders
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/agent-groups/webhooks", apiURL),
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

	var result Webhook
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("✅ Webhook '%s' created successfully (ID: %s)\n", name, result.ID)
	return nil
}

func updateWebhookStatus(webhookID string, enabled bool) error {
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
		fmt.Sprintf("%s/api/v1/agent-groups/webhooks/%s", apiURL, webhookID),
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
	fmt.Printf("✅ Webhook '%s' %s\n", webhookID, status)
	return nil
}

func deleteWebhook(webhookID string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/agent-groups/webhooks/%s", apiURL, webhookID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✅ Webhook '%s' deleted successfully\n", webhookID)
	return nil
}

func viewWebhookLogs(webhookID string, limit int) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	url := fmt.Sprintf("%s/api/v1/agent-groups/webhooks/logs?limit=%d", apiURL, limit)
	if webhookID != "" {
		url = fmt.Sprintf("%s&webhook_id=%s", url, webhookID)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var logs []WebhookLog
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(logs) == 0 {
		fmt.Println("No webhook logs found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tWEBHOOK\tEVENT\tSTATUS\tERROR")
	fmt.Fprintln(w, "---------\t-------\t-----\t------\t-----")

	for _, log := range logs {
		timestamp := time.Unix(log.Timestamp, 0).Format("2006-01-02 15:04:05")
		status := fmt.Sprintf("✅ %d", log.StatusCode)
		errMsg := "-"

		if !log.Success {
			status = fmt.Sprintf("❌ %d", log.StatusCode)
			errMsg = log.ErrorMessage
			if len(errMsg) > 40 {
				errMsg = errMsg[:37] + "..."
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			timestamp,
			log.WebhookID,
			log.Event,
			status,
			errMsg,
		)
	}

	w.Flush()
	return nil
}
