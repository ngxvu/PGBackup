# PostgreSQL Database Backup Tool

[Tiếng Việt 🇻🇳](README_vi.md)

## Project Layout
The project is organized as follows:

```
backup/
├── backupFunc/                           # Contains backup execution functions
│   └── backupFunc.go                     # Core backup functionality
├── backups/                              # Backup files directory (this directory will be created automatically)
├── config/                               # Configuration management
│   ├── addingPath/                       # PostgreSQL path utilities
│   ├── checkPsqlLatestVersion/           # Version checking utilities
│   ├── checkPsqlVersionExistOnWindows/   # Install verification
│   ├── dbconfig/                         # Database configuration utilities
│   ├── downloadPsqlInstaller/            # PostgreSQL installer download utilities
│   ├── getCurrentFolderPath/             # Current folder path utilities
│   └── installPsql/                      # PostgreSQL installation utilities
├── model/                                # Data structures and constants
├── .env                                  # Environment variables (this file needs to be created, read the README for details)
├── .gitignore                            # Git ignore file
├── main.go                               # Application entry point
└── README.md                             # Project documentation
```

This repository contains a Go application for creating automated PostgreSQL database backups. The tool connects to PostgreSQL databases, verifies compatibility, and creates backups of specified schemas.

## Features

- Environment variable configuration via `.env` file
- Interactive credential collection when environment variables are missing
- PostgreSQL version detection and compatibility checking
- Automatic PostgreSQL installation if needed
- Concurrent backup processing
- Schema-specific backups

## Requirements
- Windows operating system
- Go 1.23.1 or higher
- PostgreSQL client tools (if not installed, the application will prompt for installation)

## How this tool works
- The application checks for the PostgreSQL client tools if it is the latest version.
- If the tools are not found, it will prompt you to install them, or update them if they are not the latest version.
- The application will then connect to the PostgreSQL database using the provided credentials.
- It use the `pg_dump` command to create backups of the specified schemas.
- The backups are stored in the `backups/` directory, organized by schema and version.
- The application will create the necessary directories if they do not exist.

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
        if err := BackupDatabase(creds, version, "public"); err != nil {
            errChan <- fmt.Errorf("error backing up public schema: %v", err)
        }
    }()
    go func() {
        defer wg.Done()
        if err := BackupDatabaseCustomSchema(creds, version, "new_schema"); err != nil {
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

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.