# PostgreSQL Database Backup Tool

[Tiếng Việt](README_vi.md) | [English](README.md)

## Project Layout
The project is organized as follows:

```
backup/
├── backupFunc/                           # Chứa các hàm thực hiện sao lưu
│   └── backupFunc.go                     # Chức năng sao lưu chính
├── backups/                              # Thư mục lưu trữ file backup (thư mục này sẽ được tạo tự động)
├── config/                               # Quản lý cấu hình
│   ├── addingPath/                       # Tiện ích đường dẫn PostgreSQL
│   ├── checkPsqlLatestVersion/           # Tiện ích kiểm tra phiên bản
│   ├── checkPsqlVersionExistOnWindows/   # Xác minh cài đặt
│   ├── dbconfig/                         # Tiện ích cấu hình cơ sở dữ liệu
│   ├── downloadPsqlInstaller/            # Tiện ích tải xuống bộ cài đặt PostgreSQL
│   ├── getCurrentFolderPath/             # Tiện ích lấy đường dẫn thư mục hiện tại
│   └── installPsql/                      # Tiện ích cài đặt PostgreSQL
├── model/                                # Các cấu trúc dữ liệu và hằng số
├── .env                                  # Biến môi trường (cần tạo file này, xem README để biết chi tiết)
├── .gitignore                            # File Git ignore
├── main.go                               # Điểm vào ứng dụng
└── README.md                             # Tài liệu dự án
```

This repository contains a Go application for creating automated PostgreSQL database backups. The tool connects to PostgreSQL databases, verifies compatibility, and creates backups of specified schemas.

## Features

- Environment variable configuration via `.env` file
- Interactive credential collection when environment variables are missing
- PostgreSQL version detection and compatibility checking
- Automatic PostgreSQL installation if needed
- Concurrent backup processing
- Schema-specific backups

## Setup

### Environment Variables

Create a `.env` file in the root directory with the following configuration:

```
DB_HOST=your_database_host
DB_PORT=your_database_port
DB_DATABASE=your_database_name
DB_USERNAME=your_username
DB_PASSWORD=your_password
```

If any variables are missing, the application will prompt you for them interactively.

## Usage

1. Clone the repository
2. Configure your `.env` file
3. Run the application:

```
go run main.go
```

The application will:
1. Connect to your database using the provided credentials
2. Check the PostgreSQL version on your server
3. Verify if compatible PostgreSQL tools are installed on your system
4. Install PostgreSQL tools if needed
5. Perform database backups

```markdown
## Creating Custom Backups

To create backups for additional schemas, you can use the `BackupDatabase()` function:

```go
func BackupDatabase(creds *model.DatabaseCredentials, version, schema, backupDir string) error
```

### Parameters:
- `creds`: Database credentials
- `version`: PostgreSQL version (automatically detected)
- `schema`: Schema name to backup (e.g., "public", "<your_custom_schema_name>")
- `backupDir`: Directory to store backup files

### Example usage:

```go
func BackupDatabaseNewSchema(creds *model.DatabaseCredentials, version string) error {
    return BackupDatabase(creds, version, "new_schema", model.BackupDirNewSchema)
}
```

Add it to your backup process in `PerformDatabaseBackups()`:

```go
go func() {
    defer wg.Done()
    if err := BackupDatabaseCustomSchema(creds, addPathVersion); err != nil {
        errChan <- fmt.Errorf("error backing up custom schema: %v", err)
    }
}()
```

And remember to change the `go` routine in the `PerformDatabaseBackups()` function to include the new schema backup:

```go
func PerformDatabaseBackups(creds *model.DatabaseCredentials, version string) error {
    var wg sync.WaitGroup
    errChan := make(chan error, 2)

    wg.Add(2)
    go func() {
        defer wg.Done()
        if err := BackupDatabase(creds, version, "public", model.BackupDirPublic); err != nil {
            errChan <- fmt.Errorf("error backing up public schema: %v", err)
        }
    }()
    go func() {
        defer wg.Done()
        if err := BackupDatabaseCustomSchema(creds, version, "new_schema", model.BackupDirNewSchema); err != nil {
            errChan <- fmt.Errorf("error backing up new schema: %v", err)
        }
    }()

    wg.Wait()
    close(errChan)
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    return nil
}
```
```

## Requirements

- Windows operating system
- Go 1.23.1 or higher
- Internet connection (for PostgreSQL installation if needed)

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.