package db

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewQueryCommand creates the db query command
func NewQueryCommand(ctx *commands.AppContext) *cobra.Command {
	var (
		dbName string
		format string
	)

	cmd := &cobra.Command{
		Use:   "query <sql-query>",
		Short: "Execute a SQL query on a database",
		Long: `Execute a SQL query on sloth-runner databases.

Available databases:
  - agents: Agent registry database (.sloth-cache/agents.db)
  - hooks:  Hooks and events database (.sloth-cache/hooks.db)

Example:
  sloth-runner db query "SELECT * FROM agents" --db agents --format json
  sloth-runner db query "SELECT * FROM hooks WHERE enabled = 1" --db hooks
  sloth-runner db query "SELECT COUNT(*) as total FROM events" --db hooks --format table`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			// Get database path
			dbPath, err := getDBPath(dbName)
			if err != nil {
				return err
			}

			// Open database
			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			defer db.Close()

			// Execute query
			rows, err := db.Query(query)
			if err != nil {
				return fmt.Errorf("query failed: %w", err)
			}
			defer rows.Close()

			// Get column names
			columns, err := rows.Columns()
			if err != nil {
				return fmt.Errorf("failed to get columns: %w", err)
			}

			// Collect results
			var results []map[string]interface{}
			for rows.Next() {
				// Create slice for scanning
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}

				// Scan row
				if err := rows.Scan(valuePtrs...); err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}

				// Create map for this row
				row := make(map[string]interface{})
				for i, col := range columns {
					val := values[i]
					// Convert []byte to string
					if b, ok := val.([]byte); ok {
						row[col] = string(b)
					} else {
						row[col] = val
					}
				}
				results = append(results, row)
			}

			if err := rows.Err(); err != nil {
				return fmt.Errorf("error iterating rows: %w", err)
			}

			// Display results based on format
			return displayResults(results, columns, format)
		},
	}

	cmd.Flags().StringVar(&dbName, "db", "agents", "Database to query: agents or hooks")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format: table, json, csv")

	return cmd
}

// getDBPath returns the path to the specified database
func getDBPath(dbName string) (string, error) {
	cacheDir := filepath.Join(".", ".sloth-cache")

	switch dbName {
	case "agents":
		return filepath.Join(cacheDir, "agents.db"), nil
	case "hooks":
		return filepath.Join(cacheDir, "hooks.db"), nil
	default:
		return "", fmt.Errorf("unknown database: %s (use 'agents' or 'hooks')", dbName)
	}
}

// displayResults formats and displays query results
func displayResults(results []map[string]interface{}, columns []string, format string) error {
	if len(results) == 0 {
		pterm.Info.Println("No results")
		return nil
	}

	switch format {
	case "json":
		return displayJSON(results)
	case "csv":
		return displayCSV(results, columns)
	case "table":
		return displayTable(results, columns)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

// displayJSON outputs results as JSON
func displayJSON(results []map[string]interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

// displayCSV outputs results as CSV
func displayCSV(results []map[string]interface{}, columns []string) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header
	if err := writer.Write(columns); err != nil {
		return err
	}

	// Write rows
	for _, row := range results {
		record := make([]string, len(columns))
		for i, col := range columns {
			record[i] = fmt.Sprintf("%v", row[col])
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// displayTable outputs results as a formatted table
func displayTable(results []map[string]interface{}, columns []string) error {
	// Build table data
	tableData := [][]string{columns}

	for _, row := range results {
		record := make([]string, len(columns))
		for i, col := range columns {
			val := row[col]
			if val == nil {
				record[i] = "NULL"
			} else {
				record[i] = fmt.Sprintf("%v", val)
			}
		}
		tableData = append(tableData, record)
	}

	// Display table
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	// Show row count
	pterm.Info.Printf("Total rows: %d\n", len(results))

	return nil
}
