package db

import (
	"database/sql"
	"fmt"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewSchemaCommand creates the db schema command
func NewSchemaCommand(ctx *commands.AppContext) *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "schema <table-name>",
		Short: "Show the schema of a table",
		Long: `Display the schema (structure) of a table in the database.

Example:
  sloth-runner db schema agents --db agents
  sloth-runner db schema hooks --db hooks
  sloth-runner db schema events --db hooks`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tableName := args[0]

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

			// Get table schema
			query := `
				SELECT sql
				FROM sqlite_master
				WHERE type = 'table'
				AND name = ?
			`

			var createSQL string
			err = db.QueryRow(query, tableName).Scan(&createSQL)
			if err == sql.ErrNoRows {
				return fmt.Errorf("table '%s' not found in database '%s'", tableName, dbName)
			} else if err != nil {
				return fmt.Errorf("failed to get table schema: %w", err)
			}

			// Display schema
			pterm.DefaultSection.Printf("Schema for table: %s", tableName)
			fmt.Println()

			// Pretty print the CREATE TABLE statement
			pterm.DefaultBox.Println(createSQL)
			fmt.Println()

			// Get column information using PRAGMA table_info
			columnQuery := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
			rows, err := db.Query(columnQuery)
			if err != nil {
				return fmt.Errorf("failed to get column info: %w", err)
			}
			defer rows.Close()

			// Collect column information
			type ColumnInfo struct {
				CID          int
				Name         string
				Type         string
				NotNull      int
				DefaultValue sql.NullString
				PK           int
			}

			var columns []ColumnInfo
			for rows.Next() {
				var c ColumnInfo
				if err := rows.Scan(&c.CID, &c.Name, &c.Type, &c.NotNull, &c.DefaultValue, &c.PK); err != nil {
					return fmt.Errorf("failed to scan column info: %w", err)
				}
				columns = append(columns, c)
			}

			// Display column details in a table
			pterm.DefaultSection.Println("Column Details")
			fmt.Println()

			tableData := [][]string{{"Column", "Type", "Null", "Default", "PK"}}
			for _, c := range columns {
				notNull := "YES"
				if c.NotNull == 1 {
					notNull = "NO"
				}

				defaultVal := "NULL"
				if c.DefaultValue.Valid {
					defaultVal = c.DefaultValue.String
				}

				pk := ""
				if c.PK > 0 {
					pk = "✓"
				}

				tableData = append(tableData, []string{
					c.Name,
					c.Type,
					notNull,
					defaultVal,
					pk,
				})
			}

			pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

			// Get and display indexes
			indexQuery := `
				SELECT name, sql
				FROM sqlite_master
				WHERE type = 'index'
				AND tbl_name = ?
				AND sql IS NOT NULL
				ORDER BY name
			`

			indexRows, err := db.Query(indexQuery, tableName)
			if err != nil {
				return fmt.Errorf("failed to get indexes: %w", err)
			}
			defer indexRows.Close()

			var indexes []struct {
				Name string
				SQL  string
			}

			for indexRows.Next() {
				var idx struct {
					Name string
					SQL  string
				}
				if err := indexRows.Scan(&idx.Name, &idx.SQL); err != nil {
					return fmt.Errorf("failed to scan index: %w", err)
				}
				indexes = append(indexes, idx)
			}

			if len(indexes) > 0 {
				fmt.Println()
				pterm.DefaultSection.Println("Indexes")
				fmt.Println()

				for _, idx := range indexes {
					pterm.Info.Printf("• %s\n", idx.Name)
					fmt.Printf("  %s\n\n", idx.SQL)
				}
			}

			// Get row count
			var count int
			countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
			if err := db.QueryRow(countQuery).Scan(&count); err == nil {
				pterm.Success.Printf("Total rows: %d\n", count)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbName, "db", "agents", "Database to query: agents or hooks")

	return cmd
}
