# Công Cụ Sao Lưu Cơ Sở Dữ Liệu PostgreSQL

[English](README.md) | [Tiếng Việt](README_vi.md)

## Cấu Trúc Dự Án
Dự án được tổ chức như sau:

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

Repository này chứa ứng dụng Go để tạo các bản sao lưu cơ sở dữ liệu PostgreSQL tự động. Công cụ kết nối với cơ sở dữ liệu PostgreSQL, xác minh khả năng tương thích và tạo bản sao lưu của các schema được chỉ định.

## Tính năng

- Cấu hình biến môi trường qua file `.env`
- Thu thập thông tin xác thực tương tác khi thiếu biến môi trường
- Phát hiện phiên bản PostgreSQL và kiểm tra khả năng tương thích
- Tự động cài đặt PostgreSQL nếu cần
- Xử lý sao lưu đồng thời
- Sao lưu theo schema cụ thể

## Cài đặt

### Biến Môi Trường

Tạo file `.env` trong thư mục gốc với cấu hình sau:

```
DB_HOST=host_cơ_sở_dữ_liệu
DB_PORT=cổng_cơ_sở_dữ_liệu
DB_DATABASE=tên_cơ_sở_dữ_liệu
DB_USERNAME=tên_người_dùng
DB_PASSWORD=mật_khẩu
```

Nếu thiếu bất kỳ biến nào, ứng dụng sẽ nhắc bạn nhập tương tác.

## Cách sử dụng

1. Clone repository
2. Cấu hình file `.env`
3. Chạy ứng dụng:

```
go run main.go
```

Ứng dụng sẽ:
1. Kết nối với cơ sở dữ liệu sử dụng thông tin đã cung cấp
2. Kiểm tra phiên bản PostgreSQL trên máy chủ
3. Xác minh nếu công cụ PostgreSQL tương thích đã được cài đặt trên hệ thống
4. Cài đặt công cụ PostgreSQL nếu cần
5. Thực hiện sao lưu cơ sở dữ liệu

## Tạo Bản Sao Lưu Tùy Chỉnh

Để tạo bản sao lưu cho các schema bổ sung, bạn có thể sử dụng hàm `BackupDatabase()`:

```go
func BackupDatabase(creds *model.DatabaseCredentials, version, schema, backupDir string) error
```

### Tham số:
- `creds`: Thông tin xác thực cơ sở dữ liệu
- `version`: Phiên bản PostgreSQL (được phát hiện tự động)
- `schema`: Tên schema cần sao lưu (ví dụ: "public", "<tên_schema_tùy_chỉnh>")
- `backupDir`: Thư mục lưu trữ file sao lưu

### Ví dụ sử dụng:

```go
func BackupDatabaseNewSchema(creds *model.DatabaseCredentials, version string) error {
    return BackupDatabase(creds, version, "new_schema", model.BackupDirNewSchema)
}
```

Thêm vào quy trình sao lưu trong `PerformDatabaseBackups()`:

```go
go func() {
    defer wg.Done()
    if err := BackupDatabaseCustomSchema(creds, addPathVersion); err != nil {
        errChan <- fmt.Errorf("lỗi khi sao lưu schema tùy chỉnh: %v", err)
    }
}()
```

Và nhớ thay đổi goroutine trong hàm `PerformDatabaseBackups()` để bao gồm sao lưu schema mới:

```go
func PerformDatabaseBackups(creds *model.DatabaseCredentials, version string) error {
    var wg sync.WaitGroup
    errChan := make(chan error, 2)

    wg.Add(2)
    go func() {
        defer wg.Done()
        if err := BackupDatabase(creds, version, "public", model.BackupDirPublic); err != nil {
            errChan <- fmt.Errorf("lỗi khi sao lưu schema public: %v", err)
        }
    }()
    go func() {
        defer wg.Done()
        if err := BackupDatabaseCustomSchema(creds, version, "new_schema", model.BackupDirNewSchema); err != nil {
            errChan <- fmt.Errorf("lỗi khi sao lưu schema mới: %v", err)
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

## Yêu cầu

- Hệ điều hành Windows
- Go 1.23.1 trở lên
- Kết nối internet (để cài đặt PostgreSQL nếu cần)

## Đóng góp
Đóng góp luôn được hoan nghênh! Vui lòng mở issue hoặc gửi pull request cho bất kỳ cải tiến hoặc sửa lỗi nào.
```