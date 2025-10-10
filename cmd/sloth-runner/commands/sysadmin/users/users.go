package users

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users",
		Short:   "User and group management",
		Long:    `Manage system users and groups: list, add, remove, modify, and control group membership.`,
		Aliases: []string{"user"},
	}

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List system users",
		Run: func(cmd *cobra.Command, args []string) {
			system, _ := cmd.Flags().GetBool("system")
			filter, _ := cmd.Flags().GetString("filter")
			group, _ := cmd.Flags().GetString("group")

			if err := runList(system, filter, group); err != nil {
				pterm.Error.Printf("Failed to list users: %v\n", err)
			}
		},
	}
	listCmd.Flags().BoolP("system", "s", false, "Include system users (UID < 1000)")
	listCmd.Flags().StringP("filter", "f", "", "Filter by username or full name")
	listCmd.Flags().StringP("group", "g", "", "Filter by group membership")
	cmd.AddCommand(listCmd)

	// info command
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show detailed user information",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("user")

			if username == "" {
				pterm.Error.Println("Username is required (use --user)")
				return
			}

			if err := runInfo(username); err != nil {
				pterm.Error.Printf("Failed to get user info: %v\n", err)
			}
		},
	}
	infoCmd.Flags().StringP("user", "u", "", "Username")
	cmd.AddCommand(infoCmd)

	// add command
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new user",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("user")
			fullName, _ := cmd.Flags().GetString("fullname")
			homeDir, _ := cmd.Flags().GetString("home")
			shell, _ := cmd.Flags().GetString("shell")
			groups, _ := cmd.Flags().GetStringSlice("groups")
			createHome, _ := cmd.Flags().GetBool("create-home")
			system, _ := cmd.Flags().GetBool("system")

			if username == "" {
				pterm.Error.Println("Username is required (use --user)")
				return
			}

			if err := runAdd(username, fullName, homeDir, shell, groups, createHome, system); err != nil {
				pterm.Error.Printf("Failed to add user: %v\n", err)
			}
		},
	}
	addCmd.Flags().StringP("user", "u", "", "Username")
	addCmd.Flags().StringP("fullname", "n", "", "Full name (GECOS)")
	addCmd.Flags().StringP("home", "d", "", "Home directory")
	addCmd.Flags().StringP("shell", "s", "/bin/bash", "Login shell")
	addCmd.Flags().StringSliceP("groups", "g", []string{}, "Additional groups")
	addCmd.Flags().BoolP("create-home", "m", true, "Create home directory")
	addCmd.Flags().Bool("system", false, "Create system user")
	cmd.AddCommand(addCmd)

	// remove command
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a user",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("user")
			removeHome, _ := cmd.Flags().GetBool("remove-home")

			if username == "" {
				pterm.Error.Println("Username is required (use --user)")
				return
			}

			if err := runRemove(username, removeHome); err != nil {
				pterm.Error.Printf("Failed to remove user: %v\n", err)
			}
		},
	}
	removeCmd.Flags().StringP("user", "u", "", "Username")
	removeCmd.Flags().BoolP("remove-home", "r", false, "Remove home directory")
	cmd.AddCommand(removeCmd)

	// modify command
	modifyCmd := &cobra.Command{
		Use:   "modify",
		Short: "Modify user properties",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("user")
			fullName, _ := cmd.Flags().GetString("fullname")
			homeDir, _ := cmd.Flags().GetString("home")
			shell, _ := cmd.Flags().GetString("shell")
			lock, _ := cmd.Flags().GetBool("lock")
			unlock, _ := cmd.Flags().GetBool("unlock")

			if username == "" {
				pterm.Error.Println("Username is required (use --user)")
				return
			}

			if err := runModify(username, fullName, homeDir, shell, lock, unlock); err != nil {
				pterm.Error.Printf("Failed to modify user: %v\n", err)
			}
		},
	}
	modifyCmd.Flags().StringP("user", "u", "", "Username")
	modifyCmd.Flags().StringP("fullname", "n", "", "New full name")
	modifyCmd.Flags().StringP("home", "d", "", "New home directory")
	modifyCmd.Flags().StringP("shell", "s", "", "New login shell")
	modifyCmd.Flags().Bool("lock", false, "Lock user account")
	modifyCmd.Flags().Bool("unlock", false, "Unlock user account")
	cmd.AddCommand(modifyCmd)

	// groups command
	groupsCmd := &cobra.Command{
		Use:   "groups",
		Short: "List all groups",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runGroups(); err != nil {
				pterm.Error.Printf("Failed to list groups: %v\n", err)
			}
		},
	}
	cmd.AddCommand(groupsCmd)

	// add-to-group command
	addGroupCmd := &cobra.Command{
		Use:   "add-to-group",
		Short: "Add user to a group",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("user")
			group, _ := cmd.Flags().GetString("group")

			if username == "" || group == "" {
				pterm.Error.Println("Username and group are required")
				return
			}

			if err := runAddToGroup(username, group); err != nil {
				pterm.Error.Printf("Failed to add user to group: %v\n", err)
			}
		},
	}
	addGroupCmd.Flags().StringP("user", "u", "", "Username")
	addGroupCmd.Flags().StringP("group", "g", "", "Group name")
	cmd.AddCommand(addGroupCmd)

	// remove-from-group command
	removeGroupCmd := &cobra.Command{
		Use:   "remove-from-group",
		Short: "Remove user from a group",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("user")
			group, _ := cmd.Flags().GetString("group")

			if username == "" || group == "" {
				pterm.Error.Println("Username and group are required")
				return
			}

			if err := runRemoveFromGroup(username, group); err != nil {
				pterm.Error.Printf("Failed to remove user from group: %v\n", err)
			}
		},
	}
	removeGroupCmd.Flags().StringP("user", "u", "", "Username")
	removeGroupCmd.Flags().StringP("group", "g", "", "Group name")
	cmd.AddCommand(removeGroupCmd)

	// docs subcommand
	cmd.AddCommand(&cobra.Command{
		Use:   "docs",
		Short: "Show detailed documentation (man page style)",
		Long:  `Display comprehensive documentation for the users command with examples and usage patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			showUsersDocs()
		},
	})

	return cmd
}

func showUsersDocs() {
	title := "SLOTH-RUNNER SYSADMIN USERS(1)"
	description := "sloth-runner sysadmin users - User and group management"
	synopsis := "sloth-runner sysadmin users [subcommand] [options]"

	options := [][]string{
		{"list", "List system users with filtering options."},
		{"info", "Show detailed information about a specific user."},
		{"add", "Create a new user with specified properties."},
		{"remove", "Delete a user and optionally remove home directory."},
		{"modify", "Modify user properties (shell, home, lock/unlock)."},
		{"groups", "List all system groups."},
		{"add-to-group", "Add a user to a group."},
		{"remove-from-group", "Remove a user from a group."},
		{"docs", "Show this documentation page."},
	}

	examples := [][]string{
		{
			"List all regular users",
			"sloth-runner sysadmin users list",
			"Shows users with UID >= 1000",
		},
		{
			"List including system users",
			"sloth-runner sysadmin user list --system",
			"Shows all users including system accounts",
		},
		{
			"Filter users by group",
			"sloth-runner sysadmin users list --group sudo",
			"Shows only users in the sudo group",
		},
		{
			"Show user details",
			"sloth-runner sysadmin users info --user john",
			"Displays detailed information about user john",
		},
		{
			"Add a new user",
			"sloth-runner sysadmin users add --user john --fullname \"John Doe\" --groups sudo,docker",
			"Creates user john with specified properties",
		},
		{
			"Remove a user",
			"sloth-runner sysadmin users remove --user john --remove-home",
			"Deletes user john and removes home directory",
		},
		{
			"Lock user account",
			"sloth-runner sysadmin users modify --user john --lock",
			"Locks john's account preventing login",
		},
		{
			"Change user shell",
			"sloth-runner sysadmin users modify --user john --shell /bin/zsh",
			"Changes john's default shell to zsh",
		},
		{
			"Add user to group",
			"sloth-runner sysadmin users add-to-group --user john --group docker",
			"Adds john to the docker group",
		},
	}

	seeAlso := []string{
		"sloth-runner sysadmin process - Process management",
		"sloth-runner sysadmin security - Security auditing",
	}

	showDocs(title, description, synopsis, options, examples, seeAlso)
}

// showDocs displays formatted documentation similar to man pages
func showDocs(title, description, synopsis string, options [][]string, examples [][]string, seeAlso []string) {
	// Header
	pterm.DefaultHeader.WithFullWidth().Println(title)
	fmt.Println()

	// Name and Description
	pterm.DefaultSection.Println("NAME")
	fmt.Printf("    %s\n\n", description)

	// Synopsis
	if synopsis != "" {
		pterm.DefaultSection.Println("SYNOPSIS")
		fmt.Printf("    %s\n\n", synopsis)
	}

	// Options
	if len(options) > 0 {
		pterm.DefaultSection.Println("OPTIONS")
		for _, opt := range options {
			if len(opt) >= 2 {
				pterm.FgCyan.Printf("    %s\n", opt[0])
				fmt.Printf("        %s\n\n", opt[1])
			}
		}
	}

	// Examples
	if len(examples) > 0 {
		pterm.DefaultSection.Println("EXAMPLES")
		for i, ex := range examples {
			if len(ex) >= 2 {
				pterm.FgYellow.Printf("    Example %d: %s\n", i+1, ex[0])
				pterm.FgGreen.Printf("    $ %s\n", ex[1])
				if len(ex) >= 3 {
					fmt.Printf("        %s\n", ex[2])
				}
				fmt.Println()
			}
		}
	}

	// See Also
	if len(seeAlso) > 0 {
		pterm.DefaultSection.Println("SEE ALSO")
		for _, item := range seeAlso {
			fmt.Printf("    • %s\n", item)
		}
		fmt.Println()
	}

	// Footer
	pterm.FgGray.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	pterm.FgGray.Println("Documentation generated for sloth-runner sysadmin v2.0")
	pterm.FgGray.Println("For more information: sloth-runner sysadmin users --help")
}

// runList lista usuários
func runList(system bool, filter string, group string) error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("System Users")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading users...")

	options := ListOptions{
		SystemUsers: system,
		Filter:      filter,
		Group:       group,
	}

	users, err := manager.List(options)
	if err != nil {
		spinner.Fail("Failed to load users")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d users", len(users)))
	pterm.Println()

	if len(users) == 0 {
		pterm.Info.Println("No users found")
		return nil
	}

	// Users table
	tableData := pterm.TableData{
		{"Username", "UID", "GID", "Home", "Shell", "Full Name"},
	}

	for _, u := range users {
		tableData = append(tableData, []string{
			u.Username,
			u.UID,
			u.GID,
			truncate(u.HomeDir, 30),
			truncate(u.Shell, 20),
			truncate(u.FullName, 30),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Info.Printf("Total users: %d\n", len(users))
	if filter != "" {
		pterm.Info.Printf("Filter: %s\n", filter)
	}
	if group != "" {
		pterm.Info.Printf("Group: %s\n", group)
	}

	return nil
}

// runInfo mostra informações detalhadas de um usuário
func runInfo(username string) error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("User Information")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading user details...")

	detail, err := manager.Info(username)
	if err != nil {
		spinner.Fail("Failed to load user")
		return err
	}

	spinner.Success("✅ User information loaded")
	pterm.Println()

	// Basic Info
	pterm.DefaultSection.Println("Basic Information")
	basicData := pterm.TableData{
		{"Property", "Value"},
		{"Username", detail.Username},
		{"UID", detail.UID},
		{"GID", detail.GID},
		{"Home Directory", detail.HomeDir},
	}
	pterm.DefaultTable.WithHasHeader().WithData(basicData).Render()
	pterm.Println()

	// Groups
	if len(detail.Groups) > 0 {
		pterm.DefaultSection.Println("Groups")
		fmt.Printf("    %s\n\n", strings.Join(detail.Groups, ", "))
	}

	// Account Status
	pterm.DefaultSection.Println("Account Status")
	statusData := pterm.TableData{
		{"Property", "Value"},
	}

	if detail.PasswordSet {
		statusData = append(statusData, []string{"Password", pterm.FgGreen.Sprint("Set")})
	} else {
		statusData = append(statusData, []string{"Password", pterm.FgRed.Sprint("Not Set")})
	}

	if detail.Locked {
		statusData = append(statusData, []string{"Account", pterm.FgRed.Sprint("Locked")})
	} else {
		statusData = append(statusData, []string{"Account", pterm.FgGreen.Sprint("Active")})
	}

	if detail.ExpiryDate != "" && detail.ExpiryDate != "never" {
		statusData = append(statusData, []string{"Expiry", detail.ExpiryDate})
	}

	pterm.DefaultTable.WithHasHeader().WithData(statusData).Render()
	pterm.Println()

	return nil
}

// runAdd adiciona um novo usuário
func runAdd(username, fullName, homeDir, shell string, groups []string, createHome, system bool) error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("Add User")
	pterm.Println()

	options := AddUserOptions{
		Username:   username,
		FullName:   fullName,
		HomeDir:    homeDir,
		Shell:      shell,
		Groups:     groups,
		CreateHome: createHome,
		System:     system,
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Creating user %s...", username))

	err := manager.Add(options)
	if err != nil {
		spinner.Fail("Failed to create user")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ User %s created successfully", username))
	pterm.Println()

	pterm.Info.Printf("Username: %s\n", username)
	if fullName != "" {
		pterm.Info.Printf("Full name: %s\n", fullName)
	}
	if len(groups) > 0 {
		pterm.Info.Printf("Groups: %s\n", strings.Join(groups, ", "))
	}

	return nil
}

// runRemove remove um usuário
func runRemove(username string, removeHome bool) error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("Remove User")
	pterm.Println()

	if removeHome {
		pterm.Warning.Printf("⚠️  This will remove user %s and their home directory\n\n", username)
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing user %s...", username))

	err := manager.Remove(username, removeHome)
	if err != nil {
		spinner.Fail("Failed to remove user")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ User %s removed successfully", username))
	pterm.Println()

	return nil
}

// runModify modifica um usuário
func runModify(username, fullName, homeDir, shell string, lock, unlock bool) error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("Modify User")
	pterm.Println()

	options := ModifyOptions{
		FullName: fullName,
		HomeDir:  homeDir,
		Shell:    shell,
		Lock:     lock,
		Unlock:   unlock,
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Modifying user %s...", username))

	err := manager.Modify(username, options)
	if err != nil {
		spinner.Fail("Failed to modify user")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ User %s modified successfully", username))
	pterm.Println()

	if lock {
		pterm.Info.Printf("Account locked: %s\n", username)
	}
	if unlock {
		pterm.Info.Printf("Account unlocked: %s\n", username)
	}

	return nil
}

// runGroups lista grupos
func runGroups() error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("System Groups")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Loading groups...")

	groups, err := manager.ListGroups()
	if err != nil {
		spinner.Fail("Failed to load groups")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ Found %d groups", len(groups)))
	pterm.Println()

	// Groups table
	tableData := pterm.TableData{
		{"Group Name", "GID", "Members"},
	}

	for _, g := range groups {
		members := "-"
		if len(g.Members) > 0 {
			members = strings.Join(g.Members, ", ")
			if len(members) > 50 {
				members = members[:47] + "..."
			}
		}

		tableData = append(tableData, []string{
			g.Name,
			g.GID,
			members,
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	pterm.Info.Printf("Total groups: %d\n", len(groups))

	return nil
}

// runAddToGroup adiciona usuário a grupo
func runAddToGroup(username, group string) error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("Add User to Group")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Adding %s to group %s...", username, group))

	err := manager.AddToGroup(username, group)
	if err != nil {
		spinner.Fail("Failed to add user to group")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ User %s added to group %s", username, group))
	pterm.Println()

	return nil
}

// runRemoveFromGroup remove usuário de grupo
func runRemoveFromGroup(username, group string) error {
	manager := NewUserManager()

	pterm.DefaultHeader.WithFullWidth().Println("Remove User from Group")
	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing %s from group %s...", username, group))

	err := manager.RemoveFromGroup(username, group)
	if err != nil {
		spinner.Fail("Failed to remove user from group")
		return err
	}

	spinner.Success(fmt.Sprintf("✅ User %s removed from group %s", username, group))
	pterm.Println()

	return nil
}

// truncate trunca uma string
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
