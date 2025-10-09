package group

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <group-name>",
		Short: "Delete an agent group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Print("Are you sure you want to delete this group? (yes/no): ")
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "yes" && confirm != "y" {
					fmt.Println("Cancelled")
					return nil
				}
			}
			return deleteGroup(args[0])
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	return cmd
}

func deleteGroup(groupName string) error {
	apiURL := os.Getenv("SLOTH_RUNNER_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/agent-groups/%s", apiURL, groupName), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✅ Group '%s' deleted successfully\n", groupName)
	return nil
}
