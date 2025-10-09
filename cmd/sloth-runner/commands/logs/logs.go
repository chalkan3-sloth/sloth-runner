package logs

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	_ "github.com/mattn/go-sqlite3"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Agent     string            `json:"agent,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
}

func NewLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs",
		Short:   "Manage and view logs",
		Long:    `View, search, and export logs from sloth-runner master and agents.`,
		Example: `  # Tail logs in real-time
  sloth-runner logs tail --follow

  # Search logs for errors
  sloth-runner logs search --query "error" --since 1h

  # Export logs to JSON
  sloth-runner logs export --format json --output /tmp/logs.json`,
	}

	cmd.AddCommand(newTailCmd())
	cmd.AddCommand(newSearchCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newRotateCmd())
	cmd.AddCommand(newLevelCmd())
	cmd.AddCommand(newRemoteCmd())

	return cmd
}

func newTailCmd() *cobra.Command {
	var follow bool
	var lines int
	var agent string
	var level string

	cmd := &cobra.Command{
		Use:   "tail",
		Short: "Tail logs in real-time",
		Long:  `Display the last N lines of logs and optionally follow new entries.`,
		Example: `  # Show last 10 lines
  sloth-runner logs tail

  # Follow logs in real-time
  sloth-runner logs tail --follow

  # Filter by agent
  sloth-runner logs tail --agent web-01 --follow

  # Filter by log level
  sloth-runner logs tail --level error --follow`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tailLogs(lines, follow, agent, level)
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().IntVarP(&lines, "lines", "n", 10, "Number of lines to show")
	cmd.Flags().StringVarP(&agent, "agent", "a", "", "Filter by agent name")
	cmd.Flags().StringVarP(&level, "level", "l", "", "Filter by log level (debug, info, warn, error)")

	return cmd
}

func newSearchCmd() *cobra.Command {
	var query string
	var since string
	var until string
	var agent string
	var level string
	var limit int

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search logs with filters",
		Long:  `Search through logs using text queries and filters.`,
		Example: `  # Search for errors in last hour
  sloth-runner logs search --query "error" --since 1h

  # Search in specific agent
  sloth-runner logs search --query "failed" --agent web-01

  # Search with time range
  sloth-runner logs search --query "timeout" --since 2h --until 1h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchLogs(query, since, until, agent, level, limit)
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	cmd.Flags().StringVar(&since, "since", "", "Show logs since (e.g., 1h, 30m, 24h)")
	cmd.Flags().StringVar(&until, "until", "", "Show logs until (e.g., 1h, 30m)")
	cmd.Flags().StringVarP(&agent, "agent", "a", "", "Filter by agent name")
	cmd.Flags().StringVarP(&level, "level", "l", "", "Filter by log level")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of results")

	return cmd
}

func newExportCmd() *cobra.Command {
	var format string
	var output string
	var since string
	var agent string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export logs to file",
		Long:  `Export logs in various formats (JSON, CSV, plain text).`,
		Example: `  # Export to JSON
  sloth-runner logs export --format json --output logs.json

  # Export to CSV
  sloth-runner logs export --format csv --output logs.csv --since 24h

  # Export specific agent logs
  sloth-runner logs export --agent web-01 --output web-01.log`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return exportLogs(format, output, since, agent)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json, csv)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (stdout if not specified)")
	cmd.Flags().StringVar(&since, "since", "", "Export logs since (e.g., 1h, 24h, 7d)")
	cmd.Flags().StringVarP(&agent, "agent", "a", "", "Filter by agent name")

	return cmd
}

func newRotateCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate log files manually",
		Long:  `Manually trigger log rotation to archive old logs.`,
		Example: `  # Rotate logs
  sloth-runner logs rotate

  # Force rotation
  sloth-runner logs rotate --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rotateLogs(force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force log rotation")

	return cmd
}

func newLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "level [debug|info|warn|error]",
		Short: "Change log level dynamically",
		Long:  `Change the logging level of the running master server.`,
		Example: `  # Set log level to debug
  sloth-runner logs level debug

  # Set log level to error
  sloth-runner logs level error`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setLogLevel(args[0])
		},
	}

	return cmd
}

// Implementation functions

func tailLogs(lines int, follow bool, agent, level string) error {
	logFile := getLogFilePath()

	file, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Read last N lines
	lastLines, err := readLastLines(file, lines)
	if err != nil {
		return fmt.Errorf("failed to read logs: %w", err)
	}

	// Print last lines with filters
	for _, line := range lastLines {
		if shouldPrintLine(line, agent, level) {
			fmt.Println(line)
		}
	}

	// Follow mode
	if follow {
		file.Seek(0, io.SeekEnd)
		reader := bufio.NewReader(file)

		fmt.Fprintf(os.Stderr, "\n==> Following logs (Ctrl+C to stop) <==\n\n")

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				return err
			}

			if shouldPrintLine(strings.TrimSpace(line), agent, level) {
				fmt.Print(line)
			}
		}
	}

	return nil
}

func searchLogs(query, since, until, agent, level string, limit int) error {
	logFile := getLogFilePath()

	file, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var sinceTime, untilTime time.Time
	if since != "" {
		sinceTime, err = parseSince(since)
		if err != nil {
			return fmt.Errorf("invalid since format: %w", err)
		}
	}
	if until != "" {
		untilTime, err = parseSince(until)
		if err != nil {
			return fmt.Errorf("invalid until format: %w", err)
		}
	}

	scanner := bufio.NewScanner(file)
	count := 0
	found := 0

	for scanner.Scan() && found < limit {
		line := scanner.Text()

		// Apply filters
		if !shouldPrintLine(line, agent, level) {
			continue
		}

		// Time range filter
		if !sinceTime.IsZero() || !untilTime.IsZero() {
			lineTime := extractTimestamp(line)
			if !lineTime.IsZero() {
				if !sinceTime.IsZero() && lineTime.Before(sinceTime) {
					continue
				}
				if !untilTime.IsZero() && lineTime.After(untilTime) {
					continue
				}
			}
		}

		// Query filter
		if query != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			continue
		}

		fmt.Println(line)
		found++
		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading logs: %w", err)
	}

	if found == 0 {
		fmt.Println("No logs found matching the criteria")
	} else {
		fmt.Fprintf(os.Stderr, "\nFound %d matching log entries\n", found)
	}

	return nil
}

func exportLogs(format, output, since, agent string) error {
	logFile := getLogFilePath()

	file, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Determine output writer
	var writer io.Writer = os.Stdout
	if output != "" {
		outFile, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()
		writer = outFile
	}

	var sinceTime time.Time
	if since != "" {
		sinceTime, err = parseSince(since)
		if err != nil {
			return fmt.Errorf("invalid since format: %w", err)
		}
	}

	scanner := bufio.NewScanner(file)
	var entries []LogEntry

	for scanner.Scan() {
		line := scanner.Text()

		if agent != "" && !strings.Contains(line, agent) {
			continue
		}

		lineTime := extractTimestamp(line)
		if !sinceTime.IsZero() && !lineTime.IsZero() && lineTime.Before(sinceTime) {
			continue
		}

		entry := parseLogLine(line)
		entries = append(entries, entry)
	}

	// Export based on format
	switch format {
	case "json":
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(entries)

	case "csv":
		fmt.Fprintln(writer, "Timestamp,Level,Agent,Message")
		for _, entry := range entries {
			fmt.Fprintf(writer, "%s,%s,%s,%q\n",
				entry.Timestamp.Format(time.RFC3339),
				entry.Level,
				entry.Agent,
				entry.Message,
			)
		}

	default: // text
		for _, entry := range entries {
			fmt.Fprintf(writer, "[%s] %s %s: %s\n",
				entry.Timestamp.Format("2006-01-02 15:04:05"),
				entry.Level,
				entry.Agent,
				entry.Message,
			)
		}
	}

	if output != "" {
		fmt.Printf("Logs exported to %s (%d entries)\n", output, len(entries))
	}

	return scanner.Err()
}

func rotateLogs(force bool) error {
	logFile := getLogFilePath()

	info, err := os.Stat(logFile)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if !force && info.Size() < 10*1024*1024 {
		fmt.Println("Log file is smaller than 10MB. Use --force to rotate anyway.")
		return nil
	}

	// Create rotated filename with timestamp
	rotatedFile := fmt.Sprintf("%s.%s", logFile, time.Now().Format("20060102-150405"))

	if err := os.Rename(logFile, rotatedFile); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Create new empty log file
	if _, err := os.Create(logFile); err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	fmt.Printf("Log file rotated: %s\n", rotatedFile)
	return nil
}

func setLogLevel(level string) error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[strings.ToLower(level)] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", level)
	}

	// TODO: Send signal to running master to change log level
	fmt.Printf("Log level would be set to: %s\n", level)
	fmt.Println("(Dynamic log level change not yet implemented - requires master API)")

	return nil
}

// Helper functions

func getLogFilePath() string {
	logDir := config.GetLogDir()
	return filepath.Join(logDir, "sloth-runner.log")
}

func readLastLines(file *os.File, n int) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)

	// Read all lines first (for simplicity)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return last N lines
	if len(lines) > n {
		return lines[len(lines)-n:], nil
	}
	return lines, nil
}

func shouldPrintLine(line, agent, level string) bool {
	if agent != "" && !strings.Contains(line, agent) {
		return false
	}

	if level != "" {
		levelUpper := strings.ToUpper(level)
		if !strings.Contains(line, levelUpper) {
			return false
		}
	}

	return true
}

func extractTimestamp(line string) time.Time {
	// Try to parse timestamp from various log formats
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006/01/02 15:04:05",
	}

	for _, format := range formats {
		if len(line) >= len(format) {
			if t, err := time.Parse(format, line[:len(format)]); err == nil {
				return t
			}
		}
	}

	return time.Time{}
}

func parseSince(since string) (time.Time, error) {
	var duration time.Duration
	var err error

	if strings.HasSuffix(since, "h") {
		hours := since[:len(since)-1]
		var h int
		fmt.Sscanf(hours, "%d", &h)
		duration = time.Duration(h) * time.Hour
	} else if strings.HasSuffix(since, "m") {
		minutes := since[:len(since)-1]
		var m int
		fmt.Sscanf(minutes, "%d", &m)
		duration = time.Duration(m) * time.Minute
	} else if strings.HasSuffix(since, "d") {
		days := since[:len(since)-1]
		var d int
		fmt.Sscanf(days, "%d", &d)
		duration = time.Duration(d) * 24 * time.Hour
	} else {
		return time.Time{}, fmt.Errorf("invalid duration format: %s (use 1h, 30m, 7d)", since)
	}

	return time.Now().Add(-duration), err
}

func parseLogLine(line string) LogEntry {
	entry := LogEntry{
		Timestamp: extractTimestamp(line),
		Message:   line,
		Fields:    make(map[string]string),
	}

	// Extract log level
	for _, level := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
		if strings.Contains(line, level) {
			entry.Level = level
			break
		}
	}

	// Extract agent name if present
	if strings.Contains(line, "agent=") {
		parts := strings.Split(line, "agent=")
		if len(parts) > 1 {
			agentPart := strings.Fields(parts[1])
			if len(agentPart) > 0 {
				entry.Agent = agentPart[0]
			}
		}
	}

	return entry
}

func newRemoteCmd() *cobra.Command {
	var agentName string
	var logPath string
	var lines int
	var follow bool
	var systemLog string

	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Fetch logs from remote agent systems",
		Long: `Fetch system logs from remote agents without SSH/interactive shell.
Supports common system log files and journalctl output.`,
		Example: `  # Fetch syslog from agent
  sloth-runner sysadmin logs remote --agent do-sloth-runner-01 --system syslog

  # Fetch custom log file
  sloth-runner sysadmin logs remote --agent web-01 --path /var/log/nginx/error.log

  # Follow journalctl on agent
  sloth-runner sysadmin logs remote --agent db-01 --system journalctl --follow

  # Last 50 lines from custom log
  sloth-runner sysadmin logs remote --agent app-01 --path /var/log/app.log --lines 50`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if agentName == "" {
				return fmt.Errorf("agent name is required (use --agent flag)")
			}
			return fetchRemoteLogs(agentName, logPath, systemLog, lines, follow)
		},
	}

	cmd.Flags().StringVarP(&agentName, "agent", "a", "", "Agent name (required)")
	cmd.Flags().StringVarP(&logPath, "path", "p", "", "Custom log file path")
	cmd.Flags().StringVarP(&systemLog, "system", "s", "", "System log type (syslog, messages, journalctl, kern, auth)")
	cmd.Flags().IntVarP(&lines, "lines", "n", 50, "Number of lines to show")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")

	cmd.MarkFlagRequired("agent")

	return cmd
}

func fetchRemoteLogs(agentName, logPath, systemLog string, lines int, follow bool) error {
	// Build remote command based on parameters
	var remoteCmd string

	if logPath != "" {
		// Custom log file path
		if follow {
			remoteCmd = fmt.Sprintf("tail -f -n %d %s", lines, logPath)
		} else {
			remoteCmd = fmt.Sprintf("tail -n %d %s", lines, logPath)
		}
	} else if systemLog != "" {
		// System log shortcuts
		switch strings.ToLower(systemLog) {
		case "syslog":
			if follow {
				remoteCmd = fmt.Sprintf("tail -f -n %d /var/log/syslog 2>/dev/null || tail -f -n %d /var/log/messages", lines, lines)
			} else {
				remoteCmd = fmt.Sprintf("tail -n %d /var/log/syslog 2>/dev/null || tail -n %d /var/log/messages", lines, lines)
			}
		case "messages":
			if follow {
				remoteCmd = fmt.Sprintf("tail -f -n %d /var/log/messages", lines)
			} else {
				remoteCmd = fmt.Sprintf("tail -n %d /var/log/messages", lines)
			}
		case "journalctl":
			if follow {
				remoteCmd = fmt.Sprintf("journalctl -n %d -f", lines)
			} else {
				remoteCmd = fmt.Sprintf("journalctl -n %d --no-pager", lines)
			}
		case "kern":
			if follow {
				remoteCmd = fmt.Sprintf("tail -f -n %d /var/log/kern.log", lines)
			} else {
				remoteCmd = fmt.Sprintf("tail -n %d /var/log/kern.log", lines)
			}
		case "auth":
			if follow {
				remoteCmd = fmt.Sprintf("tail -f -n %d /var/log/auth.log 2>/dev/null || tail -f -n %d /var/log/secure", lines, lines)
			} else {
				remoteCmd = fmt.Sprintf("tail -n %d /var/log/auth.log 2>/dev/null || tail -n %d /var/log/secure", lines, lines)
			}
		default:
			return fmt.Errorf("unknown system log type: %s (use: syslog, messages, journalctl, kern, auth)", systemLog)
		}
	} else {
		return fmt.Errorf("either --path or --system must be specified")
	}

	// Execute remote command via agent
	fmt.Fprintf(os.Stderr, "Fetching logs from agent %s...\n", agentName)
	if follow {
		fmt.Fprintf(os.Stderr, "(Press Ctrl+C to stop)\n\n")
	}

	// Use the exec module to run command on remote agent
	output, err := executeRemoteCommand(agentName, remoteCmd)
	if err != nil {
		return fmt.Errorf("failed to fetch remote logs: %w", err)
	}

	fmt.Print(output)
	return nil
}

func executeRemoteCommand(agentName, command string) (string, error) {
	// Get agent address from database
	dbPath := config.GetAgentDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open agent database: %w", err)
	}
	defer db.Close()

	var address string
	err = db.QueryRow("SELECT address FROM agents WHERE name = ?", agentName).Scan(&address)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("agent not found: %s", agentName)
		}
		return "", fmt.Errorf("failed to query agent: %w", err)
	}

	// Connect to agent via gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		return "", fmt.Errorf("failed to connect to agent %s at %s: %w", agentName, address, err)
	}
	defer conn.Close()

	// Create agent client and execute command
	client := pb.NewAgentClient(conn)
	stream, err := client.RunCommand(ctx, &pb.RunCommandRequest{
		Command: command,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	// Collect all output
	var output strings.Builder
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error receiving output: %w", err)
		}

		// Append stdout and stderr chunks
		if resp.StdoutChunk != "" {
			output.WriteString(resp.StdoutChunk)
		}
		if resp.StderrChunk != "" {
			output.WriteString(resp.StderrChunk)
		}

		// Check for errors
		if resp.Error != "" {
			return output.String(), fmt.Errorf("remote command error: %s", resp.Error)
		}

		// Check if finished
		if resp.Finished {
			break
		}
	}

	return output.String(), nil
}
