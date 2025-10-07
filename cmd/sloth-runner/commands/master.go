package commands

import (
	"fmt"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/chalkan3-sloth/sloth-runner/internal/masterdb"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// MasterServerStarter is a function type that starts the master server
// This will be injected from main package
var MasterServerStarter func(port int) error

// NewMasterCommand creates the master command
func NewMasterCommand(ctx *AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Manage master servers",
		Long:  `The master command manages master server configurations, allowing you to register, select, and manage multiple master servers.`,
	}

	// Add subcommands
	cmd.AddCommand(newMasterAddCommand())
	cmd.AddCommand(newMasterListCommand())
	cmd.AddCommand(newMasterSelectCommand())
	cmd.AddCommand(newMasterShowCommand())
	cmd.AddCommand(newMasterUpdateCommand())
	cmd.AddCommand(newMasterRemoveCommand())
	cmd.AddCommand(newMasterStartCommand(ctx))

	return cmd
}

// newMasterAddCommand creates the master add command
func newMasterAddCommand() *cobra.Command {
	var description string

	cmd := &cobra.Command{
		Use:   "add <name> <address>",
		Short: "Add a new master server",
		Long: `Add a new master server configuration with a name and address.

The name must be unique and will be used to reference this master server.
The address should be in the format HOST:PORT (e.g., 192.168.1.29:50053).

Examples:
  sloth-runner master add production 192.168.1.29:50053
  sloth-runner master add staging 10.0.0.5:50053 --description "Staging environment"
  sloth-runner master add local localhost:50053`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			address := args[1]

			// Validate name
			if strings.Contains(name, ":") {
				return fmt.Errorf("master name cannot contain ':' - did you mean to use the address as second argument?")
			}

			// Validate address format
			if !strings.Contains(address, ":") {
				return fmt.Errorf("address must be in format HOST:PORT (e.g., 192.168.1.29:50053)")
			}

			// Open database
			db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			// Add master
			master := &masterdb.Master{
				Name:        name,
				Address:     address,
				Description: description,
			}

			if err := db.Add(master); err != nil {
				return err
			}

			pterm.Success.Printf("Master '%s' added successfully\n", name)
			pterm.Info.Printf("  Address: %s\n", address)
			if description != "" {
				pterm.Info.Printf("  Description: %s\n", description)
			}

			// Check if this is the default
			count, _ := db.Count()
			if count == 1 {
				pterm.Info.Println("  ⭐ Set as default master (first master added)")
			}

			pterm.Info.Println("\n💡 Use 'sloth-runner master select " + name + "' to make it the default master")

			return nil
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "", "Description of the master server")

	return cmd
}

// newMasterListCommand creates the master list command
func newMasterListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all registered master servers",
		Long:  `Lists all registered master server configurations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Open database
			db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			// Get all masters
			masters, err := db.List()
			if err != nil {
				return fmt.Errorf("failed to list masters: %w", err)
			}

			if len(masters) == 0 {
				pterm.Warning.Println("No master servers registered")
				pterm.Info.Println("\n💡 Add a master with: sloth-runner master add <name> <address>")
				return nil
			}

			// Display masters
			pterm.DefaultSection.Println("Registered Master Servers")
			fmt.Println()

			for _, master := range masters {
				if master.IsDefault {
					pterm.FgGreen.Print("⭐ ")
					pterm.Bold.Print(master.Name)
					pterm.FgGreen.Println(" (default)")
				} else {
					pterm.FgCyan.Println(master.Name)
				}
				pterm.FgGray.Printf("   Address: %s\n", master.Address)
				if master.Description != "" {
					pterm.FgGray.Printf("   Description: %s\n", master.Description)
				}
				pterm.FgGray.Printf("   Created: %s\n", master.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Println()
			}

			pterm.Info.Printf("Total: %d master(s)\n", len(masters))
			return nil
		},
	}
}

// newMasterSelectCommand creates the master select command
func newMasterSelectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "select <name>",
		Short: "Set the default master server",
		Long: `Sets a master server as the default.

The default master will be used by all commands that require a master server
unless explicitly overridden with --master flag.

Examples:
  sloth-runner master select production
  sloth-runner master select local`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Open database
			db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			// Set as default
			if err := db.SetDefault(name); err != nil {
				return err
			}

			// Get the master to display info
			master, _ := db.Get(name)

			pterm.Success.Printf("Master '%s' is now the default\n", name)
			if master != nil {
				pterm.Info.Printf("  Address: %s\n", master.Address)
			}
			pterm.Info.Println("\n💡 All commands will now use this master by default")

			return nil
		},
	}
}

// newMasterShowCommand creates the master show command
func newMasterShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [name]",
		Short: "Show master server details",
		Long: `Displays details of a master server.

If no name is provided, shows the current default master.

Examples:
  sloth-runner master show
  sloth-runner master show production`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Open database
			db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			var master *masterdb.Master
			if len(args) == 0 {
				// Show default
				master, err = db.GetDefault()
				if err != nil {
					pterm.Warning.Println("No default master configured")
					pterm.Info.Println("\n💡 Add a master with: sloth-runner master add <name> <address>")
					pterm.Info.Println("💡 Or select one with: sloth-runner master select <name>")
					return nil
				}
			} else {
				// Show specific master
				master, err = db.Get(args[0])
				if err != nil {
					return err
				}
			}

			// Display master info
			if master.IsDefault {
				pterm.DefaultHeader.WithFullWidth().Println("Default Master Server")
			} else {
				pterm.DefaultHeader.WithFullWidth().Println("Master Server")
			}
			fmt.Println()

			pterm.FgCyan.Print("Name: ")
			pterm.Bold.Println(master.Name)
			pterm.FgCyan.Print("Address: ")
			fmt.Println(master.Address)
			if master.Description != "" {
				pterm.FgCyan.Print("Description: ")
				fmt.Println(master.Description)
			}
			pterm.FgCyan.Print("Created: ")
			fmt.Println(master.CreatedAt.Format("2006-01-02 15:04:05"))
			pterm.FgCyan.Print("Updated: ")
			fmt.Println(master.UpdatedAt.Format("2006-01-02 15:04:05"))

			return nil
		},
	}
}

// newMasterUpdateCommand creates the master update command
func newMasterUpdateCommand() *cobra.Command {
	var description string

	cmd := &cobra.Command{
		Use:   "update <name> <new_address>",
		Short: "Update a master server's address",
		Long: `Updates the address of an existing master server.

Examples:
  sloth-runner master update production 192.168.1.30:50053
  sloth-runner master update staging 10.0.0.6:50053 --description "New staging server"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			newAddress := args[1]

			// Validate address format
			if !strings.Contains(newAddress, ":") {
				return fmt.Errorf("address must be in format HOST:PORT (e.g., 192.168.1.29:50053)")
			}

			// Open database
			db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			// Get old master info
			oldMaster, err := db.Get(name)
			if err != nil {
				return err
			}

			// Update master
			master := &masterdb.Master{
				Name:        name,
				Address:     newAddress,
				Description: description,
			}
			if description == "" {
				master.Description = oldMaster.Description
			}

			if err := db.Update(master); err != nil {
				return err
			}

			pterm.Success.Printf("Master '%s' updated successfully\n", name)
			pterm.Info.Printf("  Old address: %s\n", oldMaster.Address)
			pterm.Info.Printf("  New address: %s\n", newAddress)

			return nil
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "", "Update description")

	return cmd
}

// newMasterRemoveCommand creates the master remove command
func newMasterRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove a master server",
		Long: `Removes a master server configuration.

Cannot remove the default master if other masters exist.
Select a different master as default first.

Examples:
  sloth-runner master remove staging
  sloth-runner master rm old-master`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Open database
			db, err := masterdb.NewMasterDB(config.GetMastersDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			// Get master to confirm deletion
			master, err := db.Get(name)
			if err != nil {
				return err
			}

			// Confirm deletion
			pterm.Warning.Printf("Are you sure you want to remove master '%s' (%s)?\n", name, master.Address)
			confirm := pterm.DefaultInteractiveConfirm
			result, _ := confirm.Show()

			if !result {
				pterm.Info.Println("Cancelled")
				return nil
			}

			// Delete master
			if err := db.Delete(name); err != nil {
				return err
			}

			pterm.Success.Printf("Master '%s' removed successfully\n", name)

			return nil
		},
	}
}

// newMasterStartCommand creates the master start command for starting a master server
func newMasterStartCommand(ctx *AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the master gRPC server",
		Long:  `Starts the sloth-runner master gRPC server which manages agent registry.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")

			if MasterServerStarter == nil {
				return fmt.Errorf("master server starter not initialized")
			}

			return MasterServerStarter(port)
		},
	}

	cmd.Flags().IntP("port", "p", 50053, "Port for the master gRPC server")
	cmd.Flags().String("bind", "0.0.0.0", "Address to bind the master server")
	cmd.Flags().Bool("daemon", false, "Run master server as daemon")

	return cmd
}
