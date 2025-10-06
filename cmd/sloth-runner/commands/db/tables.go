package db

import (
	"database/sql"
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewTablesCommand creates the db tables command
func NewTablesCommand(ctx *commands.AppContext) *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "tables",
		Short: "List all tables in a database",
		Long: `List all tables in the specified sloth-runner database.

Example:
  sloth-runner db tables --db agents
  sloth-runner db tables --db hooks`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// Query for tables
			query := `
				SELECT name, type, sql
				FROM sqlite_master
				WHERE type IN ('table', 'view')
				AND name NOT LIKE 'sqlite_%'
				ORDER BY name
			`

			rows, err := db.Query(query)
			if err != nil {
				return fmt.Errorf("failed to query tables: %w", err)
			}
			defer rows.Close()

			// Collect table information
			type TableInfo struct {
				Name string
				Type string
				SQL  string
			}

			var tables []TableInfo
			for rows.Next() {
				var t TableInfo
				if err := rows.Scan(&t.Name, &t.Type, &t.SQL); err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
				tables = append(tables, t)
			}

			if len(tables) == 0 {
				pterm.Info.Println("No tables found")
				return nil
			}

			// Display tables
			pterm.DefaultSection.Printf("Tables in %s database", dbName)
			fmt.Println()

			tableData := [][]string{{"Name", "Type", "Row Count"}}
			for _, t := range tables {
				// Get row count for each table
				var count int
				countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", t.Name)
				if err := db.QueryRow(countQuery).Scan(&count); err != nil {
					count = -1 // Indicate error
				}

				countStr := fmt.Sprintf("%d", count)
				if count == -1 {
					countStr = "N/A"
				}

				tableData = append(tableData, []string{
					t.Name,
					t.Type,
					countStr,
				})
			}

			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			pterm.Info.Printf("Total: %d tables\n", len(tables))

			return nil
		},
	}

	cmd.Flags().StringVar(&dbName, "db", "agents", "Database to query: agents or hooks")

	return cmd
}
