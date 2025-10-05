package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/modules"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "List and inspect available Lua modules",
	Long:  `The modules command provides information about built-in Lua modules available in sloth-runner.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var modulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available modules",
	Long:  `List all built-in Lua modules available in sloth-runner with their descriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName, _ := cmd.Flags().GetString("module")
		category, _ := cmd.Flags().GetString("category")

		docsArray := modules.GetAllModuleDocs()

		// Convert array to map for easier lookup
		docs := make(map[string]modules.ModuleDoc)
		for _, doc := range docsArray {
			docs[doc.Name] = doc
		}

		if moduleName != "" {
			return showModuleDetails(moduleName, docs)
		}

		return listModules(docs, category)
	},
}

func listModules(docs map[string]modules.ModuleDoc, categoryFilter string) error {
	pterm.DefaultHeader.WithFullWidth().Println("Sloth Runner - Available Modules")
	fmt.Println()

	// Convert map to slice and sort
	modList := make([]modules.ModuleDoc, 0, len(docs))
	for _, doc := range docs {
		modList = append(modList, doc)
	}
	sort.Slice(modList, func(i, j int) bool {
		return modList[i].Name < modList[j].Name
	})

	// Module table
	tableData := pterm.TableData{
		{"Module", "Description", "Functions"},
	}

	for _, mod := range modList {
		tableData = append(tableData, []string{
			pterm.FgCyan.Sprint(mod.Name),
			mod.Description,
			pterm.FgYellow.Sprintf("%d", len(mod.Functions)),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	fmt.Println()

	pterm.Info.Printf("Total modules: %d\n", len(docs))
	pterm.Info.Println("Use 'sloth-runner modules list --module <name>' for detailed information")
	fmt.Println()

	return nil
}

func showModuleDetails(moduleName string, docs map[string]modules.ModuleDoc) error {
	doc, exists := docs[moduleName]
	if !exists {
		return fmt.Errorf("module '%s' not found", moduleName)
	}

	pterm.DefaultHeader.WithFullWidth().Printf("Module: %s", doc.Name)
	fmt.Println()

	// Module info
	infoData := pterm.TableData{
		{"Property", "Value"},
		{"Name", pterm.FgCyan.Sprint(doc.Name)},
		{"Description", doc.Description},
		{"Functions", fmt.Sprintf("%d", len(doc.Functions))},
	}
	pterm.DefaultTable.WithHasHeader().WithData(infoData).Render()
	fmt.Println()

	// Functions
	pterm.DefaultSection.Println("Functions")
	fmt.Println()

	for i, fn := range doc.Functions {
		if i > 0 {
			fmt.Println(strings.Repeat("─", 80))
			fmt.Println()
		}

		// Function name
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).
			WithTextStyle(pterm.NewStyle(pterm.FgLightCyan)).
			Printf("%s", fn.Name)
		fmt.Println()

		// Description
		fmt.Printf("  %s\n\n", fn.Description)

		// Parameters
		if fn.Parameters != "" {
			pterm.FgYellow.Println("  Parameters:")
			fmt.Printf("    %s\n\n", fn.Parameters)
		}

		// Returns
		if fn.Returns != "" {
			pterm.FgYellow.Println("  Returns:")
			fmt.Printf("    %s\n\n", fn.Returns)
		}

		// Example
		if fn.Example != "" {
			pterm.FgYellow.Println("  Example:")
			// Add indentation to example code
			lines := strings.Split(fn.Example, "\n")
			for _, line := range lines {
				fmt.Printf("    %s\n", pterm.FgGray.Sprint(line))
			}
			fmt.Println()
		}
	}

	return nil
}

func init() {
	modulesCmd.AddCommand(modulesListCmd)
	modulesListCmd.Flags().StringP("module", "m", "", "Show detailed information for a specific module")
	modulesListCmd.Flags().StringP("category", "c", "", "Filter modules by category")
}
