package main

import (
	"backup/backupFunc"
	"backup/config/addingPath"
	"backup/config/checkPsqlLatestVersion"
	"backup/config/checkPsqlVersionExistOnWindows"
	"backup/config/dbconfig"
	"backup/model"
	"fmt"
	"log"
	"sync"
)

func main() {
	// Initialize logging and setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get database credentials
	creds, err := dbconfig.ScanCredsInformation()
	if err != nil {
		log.Fatalf("Error scanning credentials: %v", err)
		log.Fatalf("Không thể lấy thông tin đăng nhập postgresql")
	}

	// Connect to database
	db, err := dbconfig.CheckDatabaseConnection(creds)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Database connection successful")
	defer db.Close()

	// Get server PostgreSQL version
	serverVersion, err := checkPsqlLatestVersion.GetAndParseServerVersion(db)
	if err != nil {
		log.Fatalf("Error processing database version: %v", err)
	}
	connectionDBVersion := *serverVersion.VersionMinor + "." + *serverVersion.PatchVersion

	// Determine which PostgreSQL version to use for backup tools
	addPathVersion, err := checkPsqlVersionExistOnWindows.HandlePostgreSQLInstallation(connectionDBVersion)
	if err != nil {
		log.Fatalf("Error handling PostgreSQL installation: %v", err)
	}

	// Add PostgreSQL bin to PATH
	if err = addingPath.AddPath(addPathVersion); err != nil {
		log.Fatalf("Error adding PostgreSQL path to system Path: %v", err)
	}

	// Perform backups concurrently
	if err = PerformDatabaseBackups(creds, addPathVersion); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	log.Println("Backup successful")
}

// PerformDatabaseBackups runs all database backups concurrently
func PerformDatabaseBackups(creds *model.DatabaseCredentials, addPathVersion string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := backupFunc.BackupDatabasePublic(creds, addPathVersion); err != nil {
			errChan <- fmt.Errorf("error backing up public database: %v", err)
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
