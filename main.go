package main

import (
	"backup/backupFunc"
	"backup/config/checkPsqlLatestVersion"
	"backup/config/checkPsqlVersionExistOnWindows"
	"backup/config/dbconfig"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// Initialize logging and setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	elevated := false
	for _, arg := range os.Args {
		if arg == "--elevated" {
			elevated = true
			break
		}
	}

	// Get database credentials
	creds, err := dbconfig.ScanCredsInformation()
	if err != nil {
		log.Fatalf("Error scanning credentials: %v", err)
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

	if !elevated {
		needsRelaunch, err := addPath(addPathVersion)
		if err != nil {
			log.Fatalf("Error adding PostgreSQL path to system Path: %v", err)
		}

		if needsRelaunch {
			log.Println("Application restarting with admin privileges... Backup will run in elevated process")
			log.Println("Backup successful")
			return
		}
	} else {
		customPath := "C:\\Program Files\\PostgreSQL\\" + addPathVersion + "\\bin"
		err := addToSystemPath(customPath)
		if err != nil {
			log.Fatalf("Error adding custom path to system PATH: %v", err)
		}
		log.Printf("Custom path added to system PATH: %s\n", customPath)
	}

	// Perform backups concurrently
	if err = backupFunc.PerformDatabaseBackups(creds, addPathVersion); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	log.Println("Backup successful")
}

func addPath(version string) (bool, error) {
	if !isAdmin() {
		log.Println("Requesting admin privileges to add PostgreSQL to PATH...")
		err := runAsAdmin()
		if err != nil {
			return false, fmt.Errorf("error requesting admin privileges: %v", err)
		}
		return true, nil
	}

	customPath := "C:\\Program Files\\PostgreSQL\\" + version + "\\bin"
	err := addToSystemPath(customPath)
	if err != nil {
		return false, fmt.Errorf("error adding custom path to system PATH: %v", err)
	}

	log.Printf("Custom path added to system PATH: %s\n", customPath)
	return false, nil
}

func addToSystemPath(path string) error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
		$currentPath = [Environment]::GetEnvironmentVariable('Path', 'Machine');
		if ($currentPath -notlike '*%s*') {
			[Environment]::SetEnvironmentVariable('Path', $currentPath + ';%s', 'Machine')
		}`, path, path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func runAsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Remove the -Wait flag so the original process can continue
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Start-Process "%s" -ArgumentList "--elevated" -Verb RunAs`, exe))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
