# Module Generator dengan Permission - Panduan Lengkap

## Overview

Module Generator telah diupgrade untuk dapat membuat modul sekaligus dengan permission-nya secara otomatis. Fitur ini sangat memudahkan pengembangan karena tidak perlu lagi membuat permission secara manual.

## Fitur Baru

### 1. Flag `--with-permissions`

Ketika menggunakan flag ini, generator akan membuat:
- ✅ **Module files** (model, repository, handler, route)
- ✅ **Permission script** (Go dan SQL)
- ✅ **Test HTTP file** dengan contoh request
- ✅ **Automatic permission assignment** ke grup yang sesuai

### 2. Permission yang Dibuat Otomatis

Untuk setiap modul, akan dibuat 4 permission:
- **Create** - Permission untuk membuat data baru
- **Read** - Permission untuk membaca/melihat data
- **Update** - Permission untuk mengupdate data
- **Delete** - Permission untuk menghapus data

### 3. Assignment ke Grup

Permission akan otomatis di-assign ke grup:
- **Admin**: Semua permission (create, read, update, delete)
- **Editor**: create, read, update (tidak termasuk delete)
- **Viewer**: read saja
- **General client**: read saja

## Cara Penggunaan

### Metode 1: Menggunakan Script Batch/Shell

#### Windows:
```bash
tools\generate-module.bat customer --with-permissions
```

#### Linux/Mac:
```bash
./tools/generate-module.sh customer --with-permissions
```

### Metode 2: Direct Go Command

```bash
go run tools/module-generator/main.go customer --with-permissions
```

## Contoh Penggunaan Lengkap

Mari kita buat modul `customer` dengan permission:

### 1. Generate Module dengan Permission

```bash
go run tools/module-generator/main.go customer --with-permissions
```

Output:
```
Created internal\modules\customer\customer_model.go
Created internal\modules\customer\customer_repository.go
Created internal\modules\customer\customer_handler.go
Created internal\modules\customer\customer_route.go
Created scripts\add-customer-permissions.go
Created scripts\add-customer-permissions.sql
Created test\customer-api-test.http
Module created successfully!

🔐 Permission script created!
Run the following command to add permissions for customer module:
go run scripts/add-customer-permissions.go
```

### 2. Update main.go

Tambahkan kode berikut ke `main.go`:

```go
// Import
"apiserver/internal/modules/customer"

// AutoMigrate
err := db.AutoMigrate(&access.User{}, &example.Example{}, &customer.Customer{}, ...)

// Initialize repository dan handler
customerRepo := customer.NewRepository(db)
customerHandler := customer.NewHandler(customerRepo)

// Register routes
customer.RegisterCustomerRoutes(app, customerHandler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)
```

### 3. Jalankan Permission Script

```bash
go run scripts/add-customer-permissions.go
```

Output:
```
Created permission: Create Customers with ID: 16
Created permission: Read Customers with ID: 17
Created permission: Update Customers with ID: 18
Created permission: Delete Customers with ID: 19
Assigned permission 19 to Admin group
Assigned permission 16 to Admin group
Assigned permission 17 to Admin group
Assigned permission 18 to Admin group
Assigned permission 16 to Editor group
Assigned permission 17 to Editor group
Assigned permission 18 to Editor group
Assigned permission 17 to Viewer group
Assigned permission 17 to General client group
Customer permissions have been added and assigned to groups!
```

### 4. Test API Endpoints

Gunakan file `test/customer-api-test.http` untuk testing:

```http
### Create a new customer (Editor permission required)
POST http://localhost:3000/v1/customers
Authorization: Bearer test-api-key-123
Content-Type: application/json

{
  "name": "Sample Customer",
  "description": "This is a sample customer for testing"
}

### Get all customers (Editor permission required)
GET http://localhost:3000/v1/customers
Authorization: Bearer test-api-key-123
```

## File yang Dibuat

### 1. Module Files
- `internal/modules/customer/customer_model.go`
- `internal/modules/customer/customer_repository.go`
- `internal/modules/customer/customer_handler.go`
- `internal/modules/customer/customer_route.go`

### 2. Permission Scripts
- `scripts/add-customer-permissions.go` - Script Go untuk menambah permission
- `scripts/add-customer-permissions.sql` - Script SQL alternatif

### 3. Test File
- `test/customer-api-test.http` - File test dengan berbagai skenario

## Keuntungan Menggunakan Fitur Ini

### 1. **Efisiensi Waktu**
- Tidak perlu membuat permission script manual
- Tidak perlu memikirkan assignment ke grup
- Langsung dapat test API dengan file HTTP yang sudah disiapkan

### 2. **Konsistensi**
- Semua modul mengikuti pattern permission yang sama
- Naming convention yang konsisten
- Assignment ke grup yang standar

### 3. **Mengurangi Error**
- Tidak ada typo dalam nama permission
- Tidak lupa assign permission ke grup
- Template yang sudah teruji

### 4. **Dokumentasi Otomatis**
- File test HTTP sebagai dokumentasi penggunaan
- Komentar yang jelas di setiap endpoint
- Contoh request yang lengkap

## Testing dengan API Key

### API Key `test-api-key-123` (Editor Group)

Dapat mengakses:
- ✅ POST /v1/customers (create)
- ✅ GET /v1/customers (read)
- ✅ GET /v1/customers/:id (read)
- ✅ PUT /v1/customers/:id (update)
- ✅ GET /v1/customers/deleted (read)
- ✅ POST /v1/customers/:id/restore (update)
- ❌ DELETE /v1/customers/:id (delete) - Akan error 403

### API Key `admin-api-key-789` (Admin Group)

Dapat mengakses semua endpoint termasuk DELETE.

## Troubleshooting

### 1. Permission Script Error

Jika ada error saat menjalankan permission script:
- Pastikan database sudah running
- Pastikan konfigurasi database benar
- Pastikan grup Admin, Editor, Viewer, General client sudah ada

### 2. API Endpoint Error 403

Jika mendapat error 403 Forbidden:
- Periksa apakah permission script sudah dijalankan
- Periksa apakah API key valid
- Periksa apakah user memiliki grup yang sesuai

### 3. Module Tidak Terdaftar

Jika endpoint tidak ditemukan:
- Pastikan sudah menambahkan kode ke main.go
- Pastikan sudah menjalankan AutoMigrate
- Restart aplikasi setelah perubahan

## Kesimpulan

Dengan fitur `--with-permissions`, proses pembuatan modul baru menjadi sangat efisien:

1. **1 Command** → Generate module + permission + test
2. **Update main.go** → Register module
3. **Run permission script** → Setup permission
4. **Ready to use** → API siap digunakan

Fitur ini sangat membantu dalam pengembangan API yang konsisten dan aman dengan sistem permission yang proper.