package main

import (
	"PgDtaBseBckUp/addingPath"
	"PgDtaBseBckUp/backupFunc"
	"PgDtaBseBckUp/checkPsqlLatestVersion"
	"PgDtaBseBckUp/checkPsqlVersionExistOnWindows"
	"PgDtaBseBckUp/downloadPsqlInstaller"
	"PgDtaBseBckUp/getCurrentFolderPath"
	"PgDtaBseBckUp/installPg"
	"PgDtaBseBckUp/model"
	"bufio"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	// Phần logic chính của chương trình
	creds, err := ScanCredsInformation()
	if err != nil {
		log.Fatalf("Error scanning credentials: %v", err)
	}

	db, err := CheckDatabaseConnection(creds)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	mess := fmt.Sprintln("Database connection successful")
	log.Println(mess)

	serverPsqlDtaBseVersion, err := CheckDatabaseVersion(db)
	if err != nil {
		log.Fatalf("Error checking database version: %v", err)
	}

	parsedVersion, err := ParsePostgresqlVersion(*serverPsqlDtaBseVersion)
	if err != nil {
		fmt.Println("Error parsing version:", err)
		return
	}

	connectionDataBaseVersionWithMinor := *parsedVersion.VersionMinor + "." + *parsedVersion.PatchVersion

	var addPathVersion string

	psqlVersionExistOnWindows, err := checkPsqlVersionExistOnWindows.CheckPsqlVersionExistOnWindows()
	if err != nil {
		postgresqlLatestVersion, err := checkPsqlLatestVersion.CheckCurrentPostgresqlLatestVersion()
		if err != nil {
			log.Println(err)
			return
		}

		err = installIfNotExist(postgresqlLatestVersion)
		if err != nil {
			return
		}
		mess := fmt.Sprintf("PostgreSQL %s installed", *postgresqlLatestVersion.LatestVersionWithMinor)
		log.Println(mess)
		addPathVersion = *postgresqlLatestVersion.VersionMinor
	} else {
		addPathVersion = *psqlVersionExistOnWindows.VersionMinor
	}

	compareVersions, err := CompareVersions(*psqlVersionExistOnWindows.LatestVersionWithMinor, connectionDataBaseVersionWithMinor)
	if err != nil {
		mess = fmt.Sprintln("Error comparing versions:", err)
		log.Println(mess)
		return
	}
	if compareVersions < 0 {
		postgresqlLatestVersion, err := checkPsqlLatestVersion.CheckCurrentPostgresqlLatestVersion()
		if err != nil {
			log.Println(err)
			return
		}

		err = installIfNotExist(postgresqlLatestVersion)
		if err != nil {
			return
		}
		mess := fmt.Sprintf("PostgreSQL %s installed", *postgresqlLatestVersion.LatestVersionWithMinor)
		log.Println(mess)
		addPathVersion = *postgresqlLatestVersion.VersionMinor
	} else {
		addPathVersion = *psqlVersionExistOnWindows.VersionMinor
	}

	err = addingPath.AddPath(addPathVersion)
	if err != nil {
		log.Fatalf("Error adding custom path to system Path: %v", err)
	}

	err = backupFunc.BackupDatabase(creds, addPathVersion)
	if err != nil {
		log.Fatalf("Error backing up database: %v", err)
		return
	}

	mess = fmt.Sprintln("Backup successful")
	log.Println(mess)
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

	err = installPg.RunInstallerWithBatch(exeFile, installDir)
	if err != nil {
		log.Fatalf("Error running installer: %v", err)
	}

	mess = fmt.Sprintf("PostgreSQL installed to: %s\n", installDir)
	log.Println(mess)
	return nil
}

func ScanCredsInformation() (*model.DatabaseCredentials, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var databaseCredentials model.DatabaseCredentials

	databaseCredentials.PgHost = os.Getenv("DB_HOST")
	databaseCredentials.PgPort = os.Getenv("DB_PORT")
	databaseCredentials.PgDatabase = os.Getenv("DB_DATABASE")
	databaseCredentials.PgUser = os.Getenv("DB_USERNAME")
	databaseCredentials.PgPassword = os.Getenv("DB_PASSWORD")

	if databaseCredentials.PgHost != "" && databaseCredentials.PgPort != "" && databaseCredentials.PgDatabase != "" && databaseCredentials.PgUser != "" && databaseCredentials.PgPassword != "" {
		return &databaseCredentials, nil
	} else {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Nhập host: ")
		databaseCredentials.PgHost, _ = reader.ReadString('\n')
		databaseCredentials.PgHost = strings.TrimSpace(databaseCredentials.PgHost)

		fmt.Println("Nhập port: ")
		databaseCredentials.PgPort, _ = reader.ReadString('\n')
		databaseCredentials.PgPort = strings.TrimSpace(databaseCredentials.PgPort)

		fmt.Println("Nhập tên database: ")
		databaseCredentials.PgDatabase, _ = reader.ReadString('\n')
		databaseCredentials.PgDatabase = strings.TrimSpace(databaseCredentials.PgDatabase)

		fmt.Println("Nhập tên user: ")
		databaseCredentials.PgUser, _ = reader.ReadString('\n')
		databaseCredentials.PgUser = strings.TrimSpace(databaseCredentials.PgUser)

		fmt.Println("Nhập password: ")
		databaseCredentials.PgPassword, _ = reader.ReadString('\n')
		databaseCredentials.PgPassword = strings.TrimSpace(databaseCredentials.PgPassword)

		fmt.Println("Nhập URL: - nếu không có thì bỏ trống")
		databaseCredentials.PgUrl, _ = reader.ReadString('\n')
		databaseCredentials.PgUrl = strings.TrimSpace(databaseCredentials.PgUrl)
	}
	return &databaseCredentials, nil
}

func CheckDatabaseConnection(creds *model.DatabaseCredentials) (*sql.DB, error) {
	var connStr string
	if creds.PgUrl != "" {
		connStr = creds.PgUrl
	} else {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			creds.PgHost, creds.PgPort, creds.PgUser, creds.PgPassword, creds.PgDatabase)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}

func CheckDatabaseVersion(db *sql.DB) (*string, error) {
	var version string
	err := db.QueryRow("SHOW server_version").Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("error querying database version: %v", err)
	}
	return &version, nil
}

func ParsePostgresqlVersion(versionString string) (*model.PostgresqlVersion, error) {
	re := regexp.MustCompile(`(\d+\.\d+)`)
	matches := re.FindStringSubmatch(versionString)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid version string: %s", versionString)
	}

	latestVersionWithMinor := matches[1]
	versionParts := strings.Split(latestVersionWithMinor, ".")
	versionMinor := versionParts[0]
	patchVersion := versionParts[1]

	return &model.PostgresqlVersion{
		LatestVersionWithMinor: &latestVersionWithMinor,
		VersionMinor:           &versionMinor,
		PatchVersion:           &patchVersion,
	}, nil
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
