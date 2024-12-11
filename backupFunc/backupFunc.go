package backupFunc

import (
	"PgDtaBseBckUp/model"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func BackupDatabase(creds *model.DatabaseCredentials, version string) error {

	programFilesDir := "C:\\Program Files\\PostgreSQL\\" + version + "\\bin"

	// Create backup directory if not exists
	if err := os.MkdirAll(model.BackupDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating backup directory: %v", err)
	}

	// Replace dots with underscores in PG_HOST
	hostWithUnderscores := strings.ReplaceAll(creds.PgHost, ".", "_")

	// Combine PG_DATABASE and modified PG_HOST to create dataSource
	dataSource := fmt.Sprintf("%s_%s", creds.PgDatabase, hostWithUnderscores)

	// Define timestamp
	timestamp := time.Now().Format("2006_01_02_15_04_05")

	// Create backup file name
	backupFile := fmt.Sprintf("%s/%s-%s-dump.sql", model.BackupDir, dataSource, timestamp)

	// Create command
	command := exec.Command(
		filepath.Join(programFilesDir, "pg_dump.exe"),
		fmt.Sprintf("--username=%s", creds.PgUser),
		fmt.Sprintf("--host=%s", creds.PgHost),
		fmt.Sprintf("--port=%s", creds.PgPort),
		fmt.Sprintf("--dbname=%s", creds.PgDatabase),
		"--format=plain",
		"--file", backupFile,
	)

	// Add password
	command.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", creds.PgPassword))

	// Execute command
	if err := command.Run(); err != nil {
		return fmt.Errorf("error during backup: %v", err)
	}

	return nil
}

func generateTableDDL(db *sql.DB) (string, error) {
	query := `
	SELECT table_name
	FROM information_schema.tables
	WHERE table_schema = 'public'
	AND table_type = 'BASE TABLE';
	`

	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying table names: %v", err)
	}
	defer rows.Close()

	var ddl string

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return "", fmt.Errorf("error scanning table name: %v", err)
		}

		// Generate DDL for each table
		tableDDL, err := generateSingleTableDDL(db, tableName)
		if err != nil {
			return "", fmt.Errorf("error generating DDL for table %s: %v", tableName, err)
		}
		ddl += tableDDL + "\n\n"
	}

	return ddl, nil
}

func generateSingleTableDDL(db *sql.DB, tableName string) (string, error) {
	query := fmt.Sprintf(`
	SELECT column_name, data_type, is_nullable, column_default
	FROM information_schema.columns
	WHERE table_name = '%s';
	`, tableName)

	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error querying columns for table %s: %v", tableName, err)
	}
	defer rows.Close()

	var ddl string
	ddl += fmt.Sprintf("CREATE TABLE %s (\n", tableName)

	columns := []string{}

	for rows.Next() {
		var columnName, dataType, isNullable, columnDefault sql.NullString
		if err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault); err != nil {
			return "", fmt.Errorf("error scanning column for table %s: %v", tableName, err)
		}

		columnDef := fmt.Sprintf("\t%s %s", columnName.String, dataType.String)
		if isNullable.String == "NO" {
			columnDef += " NOT NULL"
		}
		if columnDefault.Valid {
			columnDef += fmt.Sprintf(" DEFAULT %s", columnDefault.String)
		}
		columns = append(columns, columnDef)
	}

	ddl += fmt.Sprintf("%s\n);", fmt.Sprintf("%s", join(columns, ",\n")))

	return ddl, nil
}

func join(elements []string, delimiter string) string {
	result := ""
	for i, element := range elements {
		result += element
		if i < len(elements)-1 {
			result += delimiter
		}
	}
	return result
}
