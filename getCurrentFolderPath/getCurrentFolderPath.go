package getCurrentFolderPath

import "os"

func GetCurrentFolderPath() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return currentPath, nil
}
