package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	sshpkg "github.com/chalkan3-sloth/sloth-runner/internal/ssh"
)

var (
	// SSH command group
	sshCmd = &cobra.Command{
		Use:   "ssh",
		Short: "Manage SSH connection profiles",
		Long:  `Manage SSH connection profiles for remote execution. Profiles store connection details securely in a local SQLite database. Passwords are NEVER stored.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// SSH add command
	sshAddCmd = &cobra.Command{
		Use:   "add <profile-name> --host <host> --user <user> --key <key-path>",
		Short: "Add a new SSH profile",
		Long: `Add a new SSH connection profile to the database.

Security Note: Only connection metadata is stored. Passwords are NEVER saved.
For password authentication, use --ssh-password-stdin when executing commands.`,
		Args: cobra.ExactArgs(1),
		RunE: runSSHAdd,
	}

	// SSH list command
	sshListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all SSH profiles",
		Long:  `List all configured SSH profiles with their connection details.`,
		RunE:  runSSHList,
	}

	// SSH show command
	sshShowCmd = &cobra.Command{
		Use:   "show <profile-name>",
		Short: "Show details of an SSH profile",
		Long:  `Display detailed information about a specific SSH profile.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSSHShow,
	}

	// SSH update command
	sshUpdateCmd = &cobra.Command{
		Use:   "update <profile-name> [flags]",
		Short: "Update an existing SSH profile",
		Long:  `Update the configuration of an existing SSH profile.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSSHUpdate,
	}

	// SSH remove command
	sshRemoveCmd = &cobra.Command{
		Use:     "remove <profile-name>",
		Aliases: []string{"rm", "delete", "del"},
		Short:   "Remove an SSH profile",
		Long:    `Remove an SSH profile from the database.`,
		Args:    cobra.ExactArgs(1),
		RunE:    runSSHRemove,
	}

	// SSH test command
	sshTestCmd = &cobra.Command{
		Use:   "test <profile-name>",
		Short: "Test SSH connectivity",
		Long:  `Test SSH connectivity for a profile without executing commands.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runSSHTest,
	}
)

func init() {
	// Add SSH commands to root
	rootCmd.AddCommand(sshCmd)

	// Add subcommands
	sshCmd.AddCommand(sshAddCmd)
	sshCmd.AddCommand(sshListCmd)
	sshCmd.AddCommand(sshShowCmd)
	sshCmd.AddCommand(sshUpdateCmd)
	sshCmd.AddCommand(sshRemoveCmd)
	sshCmd.AddCommand(sshTestCmd)

	// SSH add flags
	sshAddCmd.Flags().String("host", "", "Target hostname or IP address (required)")
	sshAddCmd.Flags().String("user", "", "SSH username (required)")
	sshAddCmd.Flags().Int("port", 22, "SSH port")
	sshAddCmd.Flags().String("key", "", "Path to SSH private key (required)")
	sshAddCmd.Flags().String("description", "", "Profile description")
	sshAddCmd.Flags().Int("timeout", 30, "Connection timeout in seconds")
	sshAddCmd.Flags().Int("keepalive", 60, "Keepalive interval in seconds")
	sshAddCmd.Flags().Bool("no-strict-host-checking", false, "Disable strict host key checking")

	sshAddCmd.MarkFlagRequired("host")
	sshAddCmd.MarkFlagRequired("user")
	sshAddCmd.MarkFlagRequired("key")

	// SSH list flags
	sshListCmd.Flags().String("format", "table", "Output format: table, json, yaml, csv")
	sshListCmd.Flags().String("filter", "", "Filter expression (e.g., 'host=192.168.*')")

	// SSH show flags
	sshShowCmd.Flags().String("format", "text", "Output format: text, json, yaml")

	// SSH update flags
	sshUpdateCmd.Flags().String("host", "", "New hostname or IP address")
	sshUpdateCmd.Flags().String("user", "", "New SSH username")
	sshUpdateCmd.Flags().Int("port", 0, "New SSH port")
	sshUpdateCmd.Flags().String("key", "", "New path to SSH private key")
	sshUpdateCmd.Flags().String("description", "", "New profile description")
	sshUpdateCmd.Flags().Int("timeout", 0, "New connection timeout in seconds")
	sshUpdateCmd.Flags().Int("keepalive", 0, "New keepalive interval in seconds")
	sshUpdateCmd.Flags().Bool("no-strict-host-checking", false, "Disable strict host key checking")

	// SSH remove flags
	sshRemoveCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	// SSH test flags
	sshTestCmd.Flags().Bool("ssh-password-stdin", false, "Read SSH password from stdin")
}

// validateProfileName validates the profile name format
func validateProfileName(name string) error {
	if len(name) == 0 || len(name) > 50 {
		return fmt.Errorf("profile name must be between 1 and 50 characters")
	}

	// Must start with letter and contain only alphanumeric, hyphen, underscore
	match, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_-]*$`, name)
	if !match {
		return fmt.Errorf("profile name must start with a letter and contain only letters, numbers, hyphens, and underscores")
	}

	return nil
}

// runSSHAdd handles the ssh add command
func runSSHAdd(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	// Validate profile name
	if err := validateProfileName(profileName); err != nil {
		return err
	}

	// Get flags
	host, _ := cmd.Flags().GetString("host")
	user, _ := cmd.Flags().GetString("user")
	port, _ := cmd.Flags().GetInt("port")
	keyPath, _ := cmd.Flags().GetString("key")
	description, _ := cmd.Flags().GetString("description")
	timeout, _ := cmd.Flags().GetInt("timeout")
	keepalive, _ := cmd.Flags().GetInt("keepalive")
	noStrictHost, _ := cmd.Flags().GetBool("no-strict-host-checking")

	// Expand tilde in key path
	if strings.HasPrefix(keyPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			keyPath = strings.Replace(keyPath, "~", homeDir, 1)
		}
	}

	// Validate key file exists and has correct permissions
	if err := sshpkg.ValidateKeyFile(keyPath); err != nil {
		pterm.Warning.Printf("Key file validation failed: %v\n", err)
		pterm.Info.Println("Ensure the key file exists and has 600 permissions:")
		pterm.Printf("  chmod 600 %s\n", keyPath)
		return err
	}

	// Create profile
	profile := &sshpkg.Profile{
		Name:               profileName,
		Host:               host,
		User:               user,
		Port:               port,
		KeyPath:            keyPath,
		Description:        description,
		ConnectionTimeout:  timeout,
		KeepaliveInterval:  keepalive,
		StrictHostChecking: !noStrictHost,
	}

	// Open database
	db, err := sshpkg.NewDatabase(sshpkg.GetDefaultDatabasePath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Add profile
	if err := db.AddProfile(profile); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return fmt.Errorf("profile '%s' already exists or host/user/port combination is already in use", profileName)
		}
		return err
	}

	pterm.Success.Printf("SSH profile '%s' added successfully\n", profileName)
	pterm.Info.Println("\nProfile Details:")
	pterm.Printf("  Host: %s\n", host)
	pterm.Printf("  User: %s\n", user)
	pterm.Printf("  Port: %d\n", port)
	pterm.Printf("  Key:  %s\n", keyPath)

	if description != "" {
		pterm.Printf("  Description: %s\n", description)
	}

	pterm.Info.Println("\nTest the connection:")
	pterm.Printf("  sloth-runner ssh test %s\n", profileName)

	return nil
}

// runSSHList handles the ssh list command
func runSSHList(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	filter, _ := cmd.Flags().GetString("filter")

	// Open database
	db, err := sshpkg.NewDatabase(sshpkg.GetDefaultDatabasePath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get all profiles
	profiles, err := db.ListProfiles()
	if err != nil {
		return err
	}

	// Apply filter if provided
	if filter != "" {
		profiles = filterProfiles(profiles, filter)
	}

	if len(profiles) == 0 {
		pterm.Info.Println("No SSH profiles found")
		pterm.Info.Println("\nAdd a profile with:")
		pterm.Printf("  sloth-runner ssh add <name> --host <host> --user <user> --key <key-path>\n")
		return nil
	}

	// Display based on format
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(profiles)

	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		return encoder.Encode(profiles)

	case "csv":
		fmt.Println("NAME,HOST,USER,PORT,KEY_PATH,DESCRIPTION")
		for _, p := range profiles {
			fmt.Printf("%s,%s,%s,%d,%s,%s\n",
				p.Name, p.Host, p.User, p.Port, p.KeyPath, p.Description)
		}

	default: // table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tHOST\tUSER\tPORT\tKEY PATH")
		fmt.Fprintln(w, "----\t----\t----\t----\t--------")

		for _, p := range profiles {
			keyPath := p.KeyPath
			if len(keyPath) > 40 {
				keyPath = "..." + keyPath[len(keyPath)-37:]
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				p.Name, p.Host, p.User, p.Port, keyPath)
		}
		w.Flush()
	}

	return nil
}

// runSSHShow handles the ssh show command
func runSSHShow(cmd *cobra.Command, args []string) error {
	profileName := args[0]
	format, _ := cmd.Flags().GetString("format")

	// Open database
	db, err := sshpkg.NewDatabase(sshpkg.GetDefaultDatabasePath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get profile
	profile, err := db.GetProfile(profileName)
	if err != nil {
		return err
	}

	// Display based on format
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(profile)

	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		return encoder.Encode(profile)

	default: // text
		pterm.DefaultHeader.Printf("Profile: %s", profileName)
		fmt.Println()

		fmt.Printf("Host:           %s\n", profile.Host)
		fmt.Printf("User:           %s\n", profile.User)
		fmt.Printf("Port:           %d\n", profile.Port)
		fmt.Printf("Key Path:       %s\n", profile.KeyPath)

		if profile.Description != "" {
			fmt.Printf("Description:    %s\n", profile.Description)
		}

		fmt.Printf("Created:        %s\n", profile.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Last Modified:  %s\n", profile.UpdatedAt.Format("2006-01-02 15:04:05"))

		if profile.LastUsed != nil {
			fmt.Printf("Last Used:      %s\n", profile.LastUsed.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("Last Used:      Never\n")
		}

		fmt.Printf("Use Count:      %d\n", profile.UseCount)
		fmt.Printf("Timeout:        %ds\n", profile.ConnectionTimeout)
		fmt.Printf("Keepalive:      %ds\n", profile.KeepaliveInterval)

		if profile.StrictHostChecking {
			fmt.Printf("Strict Host:    Enabled\n")
		} else {
			fmt.Printf("Strict Host:    Disabled\n")
		}
	}

	return nil
}

// runSSHUpdate handles the ssh update command
func runSSHUpdate(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	// Build update map from changed flags
	updates := make(map[string]interface{})

	if cmd.Flags().Changed("host") {
		host, _ := cmd.Flags().GetString("host")
		updates["host"] = host
	}

	if cmd.Flags().Changed("user") {
		user, _ := cmd.Flags().GetString("user")
		updates["user"] = user
	}

	if cmd.Flags().Changed("port") {
		port, _ := cmd.Flags().GetInt("port")
		updates["port"] = port
	}

	if cmd.Flags().Changed("key") {
		keyPath, _ := cmd.Flags().GetString("key")

		// Expand tilde
		if strings.HasPrefix(keyPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				keyPath = strings.Replace(keyPath, "~", homeDir, 1)
			}
		}

		// Validate key file
		if err := sshpkg.ValidateKeyFile(keyPath); err != nil {
			return err
		}

		updates["key_path"] = keyPath
	}

	if cmd.Flags().Changed("description") {
		description, _ := cmd.Flags().GetString("description")
		updates["description"] = description
	}

	if cmd.Flags().Changed("timeout") {
		timeout, _ := cmd.Flags().GetInt("timeout")
		updates["connection_timeout"] = timeout
	}

	if cmd.Flags().Changed("keepalive") {
		keepalive, _ := cmd.Flags().GetInt("keepalive")
		updates["keepalive_interval"] = keepalive
	}

	if cmd.Flags().Changed("no-strict-host-checking") {
		noStrictHost, _ := cmd.Flags().GetBool("no-strict-host-checking")
		updates["strict_host_checking"] = !noStrictHost
	}

	if len(updates) == 0 {
		return fmt.Errorf("no updates specified")
	}

	// Open database
	db, err := sshpkg.NewDatabase(sshpkg.GetDefaultDatabasePath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Update profile
	if err := db.UpdateProfile(profileName, updates); err != nil {
		return err
	}

	pterm.Success.Printf("SSH profile '%s' updated successfully\n", profileName)

	// Show what was updated
	pterm.Info.Println("\nUpdated fields:")
	for field := range updates {
		displayField := strings.ReplaceAll(field, "_", " ")
		displayField = strings.Title(displayField)
		pterm.Printf("  • %s\n", displayField)
	}

	return nil
}

// runSSHRemove handles the ssh remove command
func runSSHRemove(cmd *cobra.Command, args []string) error {
	profileName := args[0]
	force, _ := cmd.Flags().GetBool("force")

	// Confirm deletion if not forced
	if !force {
		confirm := false
		prompt := &pterm.InteractiveConfirmPrinter{
			DefaultText: fmt.Sprintf("Are you sure you want to remove profile '%s'?", profileName),
			DefaultValue: false,
		}
		confirm, _ = prompt.Show()

		if !confirm {
			pterm.Info.Println("Removal cancelled")
			return nil
		}
	}

	// Open database
	db, err := sshpkg.NewDatabase(sshpkg.GetDefaultDatabasePath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Remove profile
	if err := db.RemoveProfile(profileName); err != nil {
		return err
	}

	pterm.Success.Printf("SSH profile '%s' removed successfully\n", profileName)
	return nil
}

// runSSHTest handles the ssh test command
func runSSHTest(cmd *cobra.Command, args []string) error {
	profileName := args[0]
	usePassword, _ := cmd.Flags().GetBool("ssh-password-stdin")

	// Open database
	db, err := sshpkg.NewDatabase(sshpkg.GetDefaultDatabasePath())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Read password if requested
	var password *string
	if usePassword {
		pwd, err := sshpkg.ReadPasswordFromStdin()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password = &pwd
	}

	// Create executor
	executor := sshpkg.NewExecutor(db)

	// Show testing message
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Testing connection to '%s'...", profileName))

	// Test connection
	if err := executor.TestConnection(profileName, password); err != nil {
		spinner.Fail(fmt.Sprintf("Connection test failed: %v", err))
		return err
	}

	spinner.Success("Connection test successful!")

	// Show profile details
	profile, _ := db.GetProfile(profileName)
	pterm.Info.Println("\nConnection Details:")
	pterm.Printf("  Host: %s\n", profile.Host)
	pterm.Printf("  User: %s\n", profile.User)
	pterm.Printf("  Port: %d\n", profile.Port)

	if usePassword {
		pterm.Printf("  Auth: Password\n")
	} else {
		pterm.Printf("  Auth: Key (%s)\n", profile.KeyPath)
	}

	return nil
}

// filterProfiles applies a simple filter to profiles
func filterProfiles(profiles []*sshpkg.Profile, filter string) []*sshpkg.Profile {
	// Simple implementation - can be enhanced
	parts := strings.Split(filter, "=")
	if len(parts) != 2 {
		return profiles
	}

	field := strings.TrimSpace(parts[0])
	pattern := strings.TrimSpace(parts[1])

	var filtered []*sshpkg.Profile
	for _, p := range profiles {
		match := false
		switch field {
		case "host":
			match = matchPattern(p.Host, pattern)
		case "user":
			match = matchPattern(p.User, pattern)
		case "name":
			match = matchPattern(p.Name, pattern)
		}

		if match {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

// matchPattern performs simple pattern matching
func matchPattern(text, pattern string) bool {
	// Convert pattern to regex (simple glob-like matching)
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = "^" + pattern + "$"

	match, _ := regexp.MatchString(pattern, text)
	return match
}