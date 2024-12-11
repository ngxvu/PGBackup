package installPg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func RunInstallerWithBatch(exeFile, installDir string) error {
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
