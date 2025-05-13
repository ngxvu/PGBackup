package checkPsqlVersionExistOnWindows

import (
	"backup/config/installPg"
	"backup/model"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// HandlePostgreSQLInstallation checks if PostgreSQL exists and installs if needed
func HandlePostgreSQLInstallation(connectionDBVersion string) (string, error) {
	// Check if PostgreSQL is already installed
	psqlVersionOnWindows, err := CheckPsqlVersionExistOnWindows()
	if err != nil {
		// PostgreSQL not found, install latest version
		return installPg.InstallLatestPostgreSQL()
	}

	// Compare installed version with connection version
	compareResult, err := CompareVersions(*psqlVersionOnWindows.LatestVersionWithMinor, connectionDBVersion)
	if err != nil {
		return "", fmt.Errorf("error comparing versions: %v", err)
	}

	if compareResult < 0 {
		// Installed version is lower than connection version, install latest
		return installPg.InstallLatestPostgreSQL()
	}

	// Use existing PostgreSQL version
	return *psqlVersionOnWindows.VersionMinor, nil
}

func CompareVersions(version1, version2 string) (int, error) {
	v1Parts := strings.Split(version1, ".")
	v2Parts := strings.Split(version2, ".")

	if len(v1Parts) != 2 || len(v2Parts) != 2 {
		return 0, fmt.Errorf("invalid version format")
	}

	v1Major, err := strconv.Atoi(v1Parts[0])
	if err != nil {
		return 0, err
	}
	v1Minor, err := strconv.Atoi(v1Parts[1])
	if err != nil {
		return 0, err
	}

	v2Major, err := strconv.Atoi(v2Parts[0])
	if err != nil {
		return 0, err
	}
	v2Minor, err := strconv.Atoi(v2Parts[1])
	if err != nil {
		return 0, err
	}

	if v1Major > v2Major {
		return 1, nil
	} else if v1Major < v2Major {
		return -1, nil
	} else {
		if v1Minor > v2Minor {
			return 1, nil
		} else if v1Minor < v2Minor {
			return -1, nil
		} else {
			return 0, nil
		}
	}
}

func CheckPsqlVersionExistOnWindows() (*model.PostgresqlVersion, error) {
	// Thực thi lệnh psql --version
	command := exec.Command("psql", "--version")
	output, err := command.CombinedOutput()
	if err != nil {
		mess := fmt.Sprintf("Postgresql chưa được cài đặt")
		log.Println(mess)
		return nil, err
	}

	// Chuyển output thành chuỗi
	versionStr := string(output)

	// Biểu thức chính quy để lấy phần version và minor version
	re := regexp.MustCompile(`(\d+)\.(\d+)`) // Ví dụ: "17.2"
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) < 3 {
		return nil, fmt.Errorf("could not extract version from psql output: %s", versionStr)
	}

	// Đặt giá trị cho LatestVersionWithMinor và VersionMinor
	latestVersion := fmt.Sprintf("%s.%s", matches[1], matches[2])
	versionMinor := matches[1] // Version Minor là phần trước dấu "."

	// PatchVersion chính là phần sau dấu "." (tức là "minor")
	patchVersion := matches[2] // Patch version là phần sau dấu chấm, ví dụ "2" từ "17.2"

	// Tạo PostgresqlVersion
	psqlVersion := &model.PostgresqlVersion{
		LatestVersionWithMinor: &latestVersion,
		VersionMinor:           &versionMinor,
		PatchVersion:           &patchVersion,
		PsqlUrl:                nil, // Có thể thêm URL nếu cần
	}

	return psqlVersion, nil
}
