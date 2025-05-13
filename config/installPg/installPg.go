package installPg

import (
	"backup/config/checkPsqlLatestVersion"
	"backup/config/downloadPsqlInstaller"
	"backup/config/getCurrentFolderPath"
	"backup/model"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// installLatestPostgreSQL installs the latest PostgreSQL version
func InstallLatestPostgreSQL() (string, error) {
	postgresqlLatestVersion, err := checkPsqlLatestVersion.CheckCurrentPostgresqlLatestVersion()
	if err != nil {
		return "", fmt.Errorf("error checking latest PostgreSQL version: %v", err)
	}

	if err = installIfNotExist(postgresqlLatestVersion); err != nil {
		return "", fmt.Errorf("installation failed: %v", err)
	}

	log.Printf("PostgreSQL %s installed", *postgresqlLatestVersion.LatestVersionWithMinor)
	return *postgresqlLatestVersion.VersionMinor, nil
}

func installIfNotExist(postgresqlLatestVersion *model.PostgresqlVersion) error {
	mess := fmt.Sprintf("Đang cài đặt phiên bản mới nhất: %s", *postgresqlLatestVersion.LatestVersionWithMinor)
	log.Println(mess)

	currentPath, err := getCurrentFolderPath.GetCurrentFolderPath()
	if err != nil {
		return err
	}

	// Create backup directory if not exists
	if err = os.MkdirAll(model.InstallersDir, os.ModePerm); err != nil {
		mess := fmt.Sprintf("Error creating installer directory: %v\n", err)
		log.Println(mess)
		return err
	}

	err = downloadPsqlInstaller.DownloadPsqlInstaller(&currentPath, postgresqlLatestVersion)
	if err != nil {
		return err
	}

	exeFile := filepath.Join(currentPath, model.InstallersDir, "psql_installer.exe")
	installDir := "C:\\Program Files\\PostgreSQL\\" + *postgresqlLatestVersion.VersionMinor

	err = runInstallerWithBatch(exeFile, installDir)
	if err != nil {
		log.Fatalf("Error running installer: %v", err)
	}

	mess = fmt.Sprintf("PostgreSQL installed to: %s\n", installDir)
	log.Println(mess)
	return nil
}

func runInstallerWithBatch(exeFile, installDir string) error {
	batchFile := "install.bat"
	content := fmt.Sprintf(`@echo off 
"%s" --mode unattended --prefix "%s"
`, exeFile, installDir)

	err := os.WriteFile(batchFile, []byte(content), 0644)
	if err != nil {
		return err
	}

	defer os.Remove(batchFile)

	cmd := exec.Command("cmd", "/c", batchFile)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.Remove(exeFile); err != nil {
		return fmt.Errorf("error removing exe file: %v", err)
	}

	mess := fmt.Sprintln("Installer removed after installation.")
	log.Println(mess)

	return nil
}
