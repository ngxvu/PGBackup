package checkPsqlVersionExistOnWindows

import (
	"PgDtaBseBckUp/model"
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

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
