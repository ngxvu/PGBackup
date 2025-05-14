package dbconfig

import (
	"backup/model"
	"bufio"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

func ScanCredsInformation() (*model.DatabaseCredentials, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, please enter database credentials manually.")
	}

	var databaseCredentials model.DatabaseCredentials

	databaseCredentials.PgHost = os.Getenv("DB_HOST")
	databaseCredentials.PgPort = os.Getenv("DB_PORT")
	databaseCredentials.PgDatabase = os.Getenv("DB_DATABASE")
	databaseCredentials.PgUser = os.Getenv("DB_USERNAME")
	databaseCredentials.PgPassword = os.Getenv("DB_PASSWORD")

	if databaseCredentials.PgHost != "" && databaseCredentials.PgPort != "" && databaseCredentials.PgDatabase != "" && databaseCredentials.PgUser != "" && databaseCredentials.PgPassword != "" {
		return &databaseCredentials, nil
	}

	reader := bufio.NewReader(os.Stdin)

	if databaseCredentials.PgHost == "" {
		fmt.Println("Nhập host: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading host: %v", err)
		}
		databaseCredentials.PgHost = strings.TrimSpace(input)
	}

	if databaseCredentials.PgPort == "" {
		fmt.Println("Nhập port: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading port: %v", err)
		}
		databaseCredentials.PgPort = strings.TrimSpace(input)
	}

	if databaseCredentials.PgDatabase == "" {
		fmt.Println("Nhập tên database: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading database name: %v", err)
		}
		databaseCredentials.PgDatabase = strings.TrimSpace(input)
	}

	if databaseCredentials.PgUser == "" {
		fmt.Println("Nhập tên user: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading username: %v", err)
		}
		databaseCredentials.PgUser = strings.TrimSpace(input)
	}

	if databaseCredentials.PgPassword == "" {
		fmt.Println("Nhập password: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading password: %v", err)
		}
		databaseCredentials.PgPassword = strings.TrimSpace(input)
	}

	return &databaseCredentials, nil
}

func CheckDatabaseConnection(creds *model.DatabaseCredentials) (*sql.DB, error) {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		creds.PgHost, creds.PgPort, creds.PgUser, creds.PgPassword, creds.PgDatabase)

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
