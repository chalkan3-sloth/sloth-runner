package stack

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/services"
	stackpkg "github.com/chalkan3-sloth/sloth-runner/internal/stack"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewDepsCommand creates the dependencies command
func NewDepsCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deps",
		Aliases: []string{"dependencies"},
		Short:   "View and manage resource dependencies",
		Long:    `View dependency graphs and analyze resource relationships.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewDepsShowCommand(ctx),
		NewDepsGraphCommand(ctx),
		NewDepsAnalyzeCommand(ctx),
	)

	return cmd
}

// NewDepsShowCommand shows dependencies for a resource
func NewDepsShowCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <stack-name> <resource-name>",
		Short: "Show dependencies for a specific resource",
		Long:  `Displays the dependency tree for a resource.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			resourceName := args[1]

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			resources, err := stackService.ListStackResources(stack.ID)
			if err != nil {
				return err
			}

			// Find the resource
			var targetResource *stackpkg.Resource
			for _, res := range resources {
				if res.Name == resourceName {
					targetResource = res
					break
				}
			}

			if targetResource == nil {
				return fmt.Errorf("resource '%s' not found in stack '%s'", resourceName, stackName)
			}

			pterm.DefaultHeader.WithFullWidth().Printfln("Dependencies: %s", resourceName)
			fmt.Println()

			if len(targetResource.Dependencies) == 0 {
				pterm.Info.Println("No dependencies")
				return nil
			}

			pterm.DefaultSection.Println("Direct Dependencies")
			for _, depID := range targetResource.Dependencies {
				// Find dependency name
				for _, res := range resources {
					if res.ID == depID {
						pterm.Info.Printf("  • %s (%s)\n", res.Name, res.Type)
						break
					}
				}
			}

			return nil
		},
	}

	return cmd
}

// NewDepsGraphCommand shows dependency graph
func NewDepsGraphCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph <stack-name>",
		Short: "Show dependency graph for all resources",
		Long:  `Visualizes the complete dependency graph for a stack.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			format, _ := cmd.Flags().GetString("format")

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			deps, err := stackService.GetStackResourceDependencies(stack.ID)
			if err != nil {
				return err
			}

			if format == "mermaid" {
				fmt.Println("```mermaid")
				fmt.Println("graph TD")
				for _, dep := range deps {
					fmt.Printf("    %s --> %s\n",
						strings.ReplaceAll(dep.ResourceID, "-", "_"),
						strings.ReplaceAll(dep.DependsOnID, "-", "_"),
					)
				}
				fmt.Println("```")
				return nil
			}

			if format == "dot" {
				fmt.Println("digraph dependencies {")
				for _, dep := range deps {
					fmt.Printf("  \"%s\" -> \"%s\";\n", dep.ResourceID, dep.DependsOnID)
				}
				fmt.Println("}")
				return nil
			}

			// ASCII tree format
			pterm.DefaultHeader.WithFullWidth().Printfln("Dependency Graph: %s", stackName)
			fmt.Println()

			if len(deps) == 0 {
				pterm.Info.Println("No dependencies defined")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "RESOURCE\t→\tDEPENDS ON")
			fmt.Fprintln(w, "--------\t \t----------")

			for _, dep := range deps {
				fmt.Fprintf(w, "%s\t→\t%s\n", dep.ResourceID, dep.DependsOnID)
			}

			w.Flush()
			fmt.Println()
			pterm.Success.Printf("Total: %d dependency link(s)\n", len(deps))

			return nil
		},
	}

	cmd.Flags().String("format", "tree", "Output format (tree, mermaid, dot)")

	return cmd
}

// NewDepsAnalyzeCommand analyzes dependencies
func NewDepsAnalyzeCommand(ctx *commands.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze <stack-name>",
		Short: "Analyze dependency health and detect issues",
		Long:  `Detects circular dependencies, orphans, and other dependency issues.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]

			stackService, err := services.NewStackService()
			if err != nil {
				return err
			}
			defer stackService.Close()

			stack, err := stackService.GetStackByName(stackName)
			if err != nil {
				return fmt.Errorf("stack '%s' not found: %w", stackName, err)
			}

			spinner, _ := pterm.DefaultSpinner.Start("Analyzing dependencies...")

			resources, err := stackService.ListStackResources(stack.ID)
			if err != nil {
				spinner.Fail("Failed to analyze dependencies")
				return err
			}

			deps, err := stackService.GetStackResourceDependencies(stack.ID)
			if err != nil {
				spinner.Fail("Failed to analyze dependencies")
				return err
			}

			// Analyze for issues
			orphans := make([]string, 0)
			circular := make([]string, 0)

			// Check for orphaned dependencies
			resourceIDs := make(map[string]bool)
			for _, res := range resources {
				resourceIDs[res.ID] = true
			}

			for _, dep := range deps {
				if !resourceIDs[dep.DependsOnID] {
					orphans = append(orphans, fmt.Sprintf("%s → %s (missing)", dep.ResourceID, dep.DependsOnID))
				}
			}

			// Simple circular dependency detection
			visited := make(map[string]bool)
			recStack := make(map[string]bool)

			var hasCycle func(string) bool
			hasCycle = func(resourceID string) bool {
				visited[resourceID] = true
				recStack[resourceID] = true

				for _, dep := range deps {
					if dep.ResourceID == resourceID {
						if !visited[dep.DependsOnID] {
							if hasCycle(dep.DependsOnID) {
								return true
							}
						} else if recStack[dep.DependsOnID] {
							circular = append(circular, fmt.Sprintf("%s ↔ %s", resourceID, dep.DependsOnID))
							return true
						}
					}
				}

				recStack[resourceID] = false
				return false
			}

			for _, res := range resources {
				if !visited[res.ID] {
					hasCycle(res.ID)
				}
			}

			hasIssues := len(orphans) > 0 || len(circular) > 0

			if !hasIssues {
				spinner.Success("✓ No dependency issues found")
				return nil
			}

			spinner.Warning(fmt.Sprintf("⚠ Found %d issue(s)", len(orphans)+len(circular)))
			fmt.Println()

			if len(orphans) > 0 {
				pterm.DefaultSection.Println("Orphaned Dependencies")
				for _, orphan := range orphans {
					pterm.Error.Printf("  • %s\n", orphan)
				}
				fmt.Println()
			}

			if len(circular) > 0 {
				pterm.DefaultSection.Println("Circular Dependencies")
				for _, cycle := range circular {
					pterm.Error.Printf("  • %s\n", cycle)
				}
				fmt.Println()
			}

			pterm.Info.Println("Use 'sloth-runner stack validate repair' to fix dependency issues")

			return nil
		},
	}

	return cmd
}
