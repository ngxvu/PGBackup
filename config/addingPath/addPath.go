package addingPath

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func AddPath(version string) error {
	if !isAdmin() {
		err := runAsAdmin()
		if err != nil {
			return fmt.Errorf("error requesting admin privileges: %v", err)
		}
		return nil
	}

	customPath := "C:\\Program Files\\PostgreSQL\\" + version + "\\bin"

	err := addToSystemPath(customPath)
	if err != nil {
		return fmt.Errorf("error adding custom path to system PATH: %v", err)
	}

	mess := fmt.Sprintf("Custom path added to system PATH: %s\n", customPath)
	log.Println(mess)

	return nil
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

	cmd := exec.Command("powershell", "Start-Process", exe, "-Verb", "runAs")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
