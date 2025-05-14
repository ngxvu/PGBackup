package backupFunc

import (
	"backup/model"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// PerformDatabaseBackups runs all database backups concurrently
func PerformDatabaseBackups(creds *model.DatabaseCredentials, addPathVersion string) error {
	var wg sync.WaitGroup
	wgCount := 2
	errChan := make(chan error, wgCount)

	wg.Add(wgCount)

	go func() {
		defer wg.Done()
		if err := backupDatabasePublic(creds, addPathVersion); err != nil {
			errChan <- fmt.Errorf("error backing up public database: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := backupDatabaseNewSchema(creds, addPathVersion); err != nil {
			errChan <- fmt.Errorf("error backing up private database: %v", err)
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		return err
	}

	return nil
}

func BackupDatabase(creds *model.DatabaseCredentials, version, schema string) error {
	programFilesDir := "C:\\Program Files\\PostgreSQL\\" + version + "\\bin"

	backupDir, _ := createSchemaDir(schema)

	// Create backup directory if not exists
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating backup directory: %v", err)
	}

	// Replace dots with underscores in PG_HOST
	hostWithUnderscores := strings.ReplaceAll(creds.PgHost, ".", "_")

	// Combine PG_DATABASE and modified PG_HOST to create dataSource
	dataSource := fmt.Sprintf("%s_%s", creds.PgDatabase, hostWithUnderscores)

	// Define timestamp
	timestamp := time.Now().Format("2006_01_02_15_04_05")

	// Create backup file name
	backupFile := fmt.Sprintf("%s/%s-%s-dump.sql", backupDir, dataSource, timestamp)

	// Create command
	command := exec.Command(
		filepath.Join(programFilesDir, "pg_dump.exe"),
		fmt.Sprintf("--username=%s", creds.PgUser),
		fmt.Sprintf("--host=%s", creds.PgHost),
		fmt.Sprintf("--port=%s", creds.PgPort),
		fmt.Sprintf("--dbname=%s", creds.PgDatabase),
		fmt.Sprintf("--schema=%s", schema),
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

func createSchemaDir(schema string) (string, error) {
	schemaDir := fmt.Sprintf("./backups/%s", schema)
	return schemaDir, nil
}

func backupDatabasePublic(creds *model.DatabaseCredentials, version string) error {
	return BackupDatabase(creds, version, "public")
}

// to create new backup function, just copy the BackupDatabase function and change the schema name, for example:
func backupDatabaseNewSchema(creds *model.DatabaseCredentials, version string) error {
	return BackupDatabase(creds, version, "dblog")
}
