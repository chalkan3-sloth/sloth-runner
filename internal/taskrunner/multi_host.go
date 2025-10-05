package taskrunner

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/chalkan3-sloth/sloth-runner/internal/types"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MultiHostResult holds the result of task execution on multiple hosts
type MultiHostResult struct {
	Host    string
	Success bool
	Output  string
	Error   error
}

// executeTaskOnMultipleHosts executes a task on multiple hosts in parallel
func (tr *TaskRunner) executeTaskOnMultipleHosts(ctx context.Context, t *types.Task, hosts []string, session *types.SharedSession, groupName string) ([]MultiHostResult, error) {
	var wg sync.WaitGroup
	results := make([]MultiHostResult, len(hosts))

	// Show header for multi-host execution
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
		Println(fmt.Sprintf("🚀 Executing task '%s' on %d hosts", t.Name, len(hosts)))
	pterm.Println()

	// Create a table for real-time status
	statusTable := pterm.TableData{
		{"Host", "Status"},
	}

	for _, host := range hosts {
		statusTable = append(statusTable, []string{host, pterm.Yellow("⏳ Pending")})
	}

	pterm.DefaultTable.WithHasHeader(true).WithData(statusTable).Render()
	pterm.Println()

	// Execute on each host in parallel
	for i, host := range hosts {
		wg.Add(1)
		go func(index int, hostAddr string) {
			defer wg.Done()

			result := MultiHostResult{
				Host: hostAddr,
			}

			// Try to resolve host if it's a name
			agentAddress := hostAddr
			if !strings.Contains(hostAddr, ":") {
				// Try to resolve agent name to address
				resolvedAddress, err := resolveAgentAddress(hostAddr)
				if err != nil {
					result.Error = fmt.Errorf("failed to resolve agent '%s': %w", hostAddr, err)
					results[index] = result
					pterm.Error.Printf("❌ Failed to resolve host %s: %v\n", hostAddr, err)
					return
				}
				agentAddress = resolvedAddress
			}

			// Connect to the agent
			pterm.Info.Printf("🔗 Connecting to %s...\n", agentAddress)
			conn, err := grpc.Dial(agentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				result.Error = fmt.Errorf("failed to connect: %w", err)
				results[index] = result
				pterm.Error.Printf("❌ Failed to connect to %s: %v\n", agentAddress, err)
				return
			}
			defer conn.Close()

			c := pb.NewAgentClient(conn)

			// Create a tarball of the workspace
			var buf bytes.Buffer
			if err := createTar(session.Workdir, &buf); err != nil {
				result.Error = fmt.Errorf("failed to create workspace tarball: %w", err)
				results[index] = result
				return
			}

			// Generate a script compatible with agent execution
			agentScript := tr.generateAgentScript(t, groupName)

			// Send the task and workspace to the agent
			r, err := c.ExecuteTask(ctx, &pb.ExecuteTaskRequest{
				TaskName:    t.Name,
				TaskGroup:   groupName,
				LuaScript:   agentScript,
				Workspace:   buf.Bytes(),
				User:        t.User,
			})

			if err != nil {
				result.Error = fmt.Errorf("failed to execute: %w", err)
				results[index] = result
				pterm.Error.Printf("❌ Failed on %s: %v\n", agentAddress, err)
				return
			}

			if !r.GetSuccess() {
				result.Success = false
				result.Output = r.GetOutput()
				result.Error = fmt.Errorf("task failed: %s", r.GetOutput())
				pterm.Error.Printf("❌ Task failed on %s\n", agentAddress)
			} else {
				result.Success = true
				result.Output = r.GetOutput()
				pterm.Success.Printf("✅ Success on %s\n", agentAddress)

				// Extract the updated workspace (only if successful)
				if err := extractTar(bytes.NewReader(r.GetWorkspace()), session.Workdir); err != nil {
					slog.Warn("Failed to extract workspace from host", "host", agentAddress, "error", err)
				}
			}

			results[index] = result
		}(i, host)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Display final results
	pterm.Println()
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
		Println("📊 Multi-Host Execution Results")
	pterm.Println()

	resultTable := pterm.TableData{
		{"Host", "Status", "Details"},
	}

	successCount := 0
	failureCount := 0

	for _, result := range results {
		status := pterm.Green("✅ Success")
		details := "Completed successfully"

		if !result.Success || result.Error != nil {
			status = pterm.Red("❌ Failed")
			failureCount++
			if result.Error != nil {
				details = result.Error.Error()
			} else if result.Output != "" {
				// Extract clean error message
				lines := strings.Split(result.Output, "\n")
				if len(lines) > 0 {
					details = strings.TrimSpace(lines[0])
				}
			}
			// Truncate long errors
			if len(details) > 60 {
				details = details[:57] + "..."
			}
		} else {
			successCount++
		}

		resultTable = append(resultTable, []string{result.Host, status, details})
	}

	pterm.DefaultTable.
		WithHasHeader(true).
		WithBoxed(true).
		WithData(resultTable).
		Render()

	pterm.Println()

	// Summary
	summaryColor := pterm.FgGreen
	summaryIcon := "✅"
	if failureCount > 0 {
		if successCount == 0 {
			summaryColor = pterm.FgRed
			summaryIcon = "❌"
		} else {
			summaryColor = pterm.FgYellow
			summaryIcon = "⚠️"
		}
	}

	pterm.DefaultBox.
		WithTitle(fmt.Sprintf("%s Execution Summary", summaryIcon)).
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(summaryColor)).
		Printfln(
			"Task:      %s\n"+
			"Total:     %d hosts\n"+
			"Success:   %s\n"+
			"Failed:    %s",
			pterm.Cyan(t.Name),
			len(hosts),
			pterm.Green(fmt.Sprintf("%d", successCount)),
			pterm.Red(fmt.Sprintf("%d", failureCount)),
		)
	pterm.Println()

	// Return error if any host failed
	if failureCount > 0 {
		return results, fmt.Errorf("%d out of %d hosts failed", failureCount, len(hosts))
	}

	return results, nil
}

// isMultiHost checks if delegate_to contains multiple hosts
func isMultiHost(delegateTo interface{}) (bool, []string) {
	switch v := delegateTo.(type) {
	case []string:
		return len(v) > 1, v
	case []interface{}:
		hosts := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				hosts = append(hosts, str)
			}
		}
		return len(hosts) > 1, hosts
	default:
		return false, nil
	}
}

// getHostsList extracts the list of hosts from delegate_to
func getHostsList(delegateTo interface{}) []string {
	switch v := delegateTo.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []interface{}:
		hosts := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				hosts = append(hosts, str)
			}
		}
		return hosts
	default:
		return nil
	}
}