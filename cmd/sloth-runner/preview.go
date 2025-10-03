package main

import (
"fmt"
"os"
"strings"
"time"

"github.com/pterm/pterm"
"github.com/spf13/cobra"
"gopkg.in/yaml.v2"

"github.com/chalkan3-sloth/sloth-runner/internal/luainterface"
"github.com/chalkan3-sloth/sloth-runner/internal/stack"
"github.com/chalkan3-sloth/sloth-runner/internal/types"
lua "github.com/yuin/gopher-lua"
)

var previewCmd = &cobra.Command{
Use:   "preview [stack-name]",
Short: "Preview task execution plan",
Args:  cobra.ExactArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
return runPreviewCommand(cmd, args)
},
}

func runPreviewCommand(cmd *cobra.Command, args []string) error {
stackName := args[0]
filePath, _ := cmd.Flags().GetString("file")
values, _ := cmd.Flags().GetString("values")

if filePath == "" {
return fmt.Errorf("--file is required")
}

stackManager, err := stack.NewStackManager("")
if err != nil {
return fmt.Errorf("failed to initialize stack manager: %w", err)
}
defer stackManager.Close()

var valuesTable *lua.LTable
if values != "" {
valuesData, err := os.ReadFile(values)
if err != nil {
return fmt.Errorf("failed to read values file: %w", err)
}

var valuesMap map[string]interface{}
if err := yaml.Unmarshal(valuesData, &valuesMap); err != nil {
return fmt.Errorf("failed to parse values file: %w", err)
}

tempL := lua.NewState()
defer tempL.Close()
valuesTable = mapToLuaTable(tempL, valuesMap)
}

taskGroups, err := luainterface.ParseLuaScript(cmd.Context(), filePath, valuesTable)
if err != nil {
return fmt.Errorf("failed to parse Lua script: %w", err)
}

if len(taskGroups) == 0 {
pterm.Warning.Println("No task groups found")
return nil
}

return showExecutionPlanPreview(stackName, filePath, taskGroups, stackManager)
}

func buildTaskTreePreview(taskNames []string, taskMap map[string]*types.Task, visited map[string]bool, stack *stack.StackState) []pterm.TreeNode {
var nodes []pterm.TreeNode

for _, taskName := range taskNames {
if visited[taskName] {
continue
}
visited[taskName] = true

task, exists := taskMap[taskName]
if !exists {
continue
}

statusIcon := "●"
statusColor := pterm.FgCyan
stateInfo := " (new)"

if stack != nil {
if res, ok := stack.TaskResults[taskName]; ok && res != nil {
statusIcon = "✓"
statusColor = pterm.FgGreen
stateInfo = " (completed)"
} else {
statusColor = pterm.FgYellow
}
}

taskDesc := fmt.Sprintf("%s %s%s", 
pterm.NewStyle(statusColor).Sprint(statusIcon),
pterm.Bold.Sprint(task.Name),
pterm.Gray(stateInfo))

var details []string
if task.Description != "" {
details = append(details, pterm.Gray("  "+task.Description))
}
if task.DelegateTo != nil {
details = append(details, pterm.Yellow(fmt.Sprintf("  -> %v", task.DelegateTo)))
}
if task.RunIf != "" {
details = append(details, pterm.Blue("  ? RunIf: "+task.RunIf))
}
if task.Retries > 0 {
details = append(details, pterm.Gray(fmt.Sprintf("  Retries: %d", task.Retries)))
}
if len(task.DependsOn) > 0 {
depStr := strings.Join(task.DependsOn, ", ")
details = append(details, pterm.Gray("  Depends: "+depStr))
}

if len(details) > 0 {
for _, d := range details {
taskDesc = taskDesc + "\n" + d
}
}

node := pterm.TreeNode{Text: taskDesc}

var dependents []string
for name, t := range taskMap {
for _, dep := range t.DependsOn {
if dep == taskName {
dependents = append(dependents, name)
break
}
}
}

if len(dependents) > 0 {
node.Children = buildTaskTreePreview(dependents, taskMap, visited, stack)
}

nodes = append(nodes, node)
}

return nodes
}

func showExecutionPlanPreview(stackName string, filePath string, taskGroups map[string]types.TaskGroup, stackManager *stack.StackManager) error {
pterm.DefaultHeader.WithFullWidth().
WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
Println("EXECUTION PLAN PREVIEW")

pterm.Info.Printfln("Stack: %s", pterm.Cyan(stackName))
pterm.Info.Printfln("File: %s", pterm.Cyan(filePath))
pterm.Println()

existingStack, err := stackManager.GetStackByName(stackName)
var previewStack *stack.StackState
if err == nil {
previewStack = existingStack
pterm.Info.Printfln("Existing stack found (ID: %s)", existingStack.ID)
pterm.Info.Printfln("Last updated: %s", existingStack.UpdatedAt.Format(time.RFC3339))
} else {
pterm.Warning.Println("Stack will be created on run")
}
pterm.Println()

totalTasks := 0
for groupName, group := range taskGroups {
pterm.DefaultSection.Println("Task Group: " + groupName)

if group.Description != "" {
pterm.Info.Println("  Description: " + group.Description)
}

if group.DelegateTo != nil {
delegateInfo := fmt.Sprintf("%v", group.DelegateTo)
pterm.Info.Println("  Delegate To: " + pterm.Yellow(delegateInfo))
}

pterm.Println()

taskMap := make(map[string]*types.Task)
for i := range group.Tasks {
task := &group.Tasks[i]
taskMap[task.Name] = task
}

rootText := fmt.Sprintf("%s (%d tasks)", pterm.Bold.Sprint(groupName), len(group.Tasks))
root := pterm.TreeNode{Text: rootText}

rootTasks := []string{}
for _, task := range group.Tasks {
if len(task.DependsOn) == 0 {
rootTasks = append(rootTasks, task.Name)
}
}

visited := make(map[string]bool)
root.Children = buildTaskTreePreview(rootTasks, taskMap, visited, previewStack)

tree := pterm.TreePrinter{Root: root}
tree.Render()
pterm.Println()

totalTasks = totalTasks + len(group.Tasks)
}

summaryMsg := fmt.Sprintf("Total Tasks: %s\nChanges: Will execute tasks based on state", pterm.Cyan(fmt.Sprintf("%d", totalTasks)))
pterm.DefaultBox.WithTitle("Execution Plan").WithTitleTopCenter().Println(summaryMsg)
pterm.Println()

return nil
}
