# My API - Modular REST API with Go Fiber

REST API yang dibangun dengan Go Fiber menggunakan arsitektur modular yang sangat scalable. Setiap modul bersifat self-contained dan dapat ditambahkan dengan mudah tanpa mengubah kode yang sudah ada.

## Struktur Proyek

```
apiserver/
├── cmd/
│   └── api/
│       └── main.go              # Entry point aplikasi
├── configs/
│   └── config.go                # Fungsi untuk load .env
├── internal/
│   ├── database/
│   │   ├── database.go          # Inisialisasi koneksi DB (GORM)
│   │   └── seeder.go            # Database seeder untuk test data
│   ├── middleware/
│   │   ├── auth.go              # Middleware untuk autentikasi Bearer Token
│   │   ├── permission.go        # Middleware untuk permission checking
│   │   └── audit.go             # Middleware untuk audit logging
│   └── modules/
│       ├── access/              # Modul untuk autentikasi & user access
│       │   ├── access_handler.go
│       │   ├── access_model.go
│       │   ├── access_repository.go
│       │   └── access_route.go
│       ├── permission/          # Modul untuk manajemen permissions
│       │   ├── permission_handler.go
│       │   ├── permission_model.go
│       │   ├── permission_repository.go
│       │   └── permission_route.go
│       ├── group/               # Modul untuk manajemen groups
│       │   ├── group_handler.go
│       │   ├── group_model.go
│       │   ├── group_repository.go
│       │   └── group_route.go
│       ├── audit/               # Modul untuk audit logging
│       │   ├── audit_handler.go
│       │   ├── audit_model.go
│       │   ├── audit_repository.go
│       │   └── audit_route.go
│       └── example/             # Modul kedua untuk demonstrasi
│           ├── example_handler.go
│           ├── example_model.go
│           ├── example_repository.go
│           └── example_route.go
├── docs/                        # Folder untuk file swagger yang digenerasi
├── .env.example                 # Contoh file environment
├── go.mod
├── go.sum
└── README.md
```

## Mengapa Struktur Ini Modular?

1. **Feature-Based Structure**: Setiap modul (access, permission, group, example) memiliki folder terpisah dengan semua komponen yang dibutuhkan
2. **Self-Contained**: Setiap modul memiliki model, repository, handler, dan route sendiri
3. **Dependency Injection**: Repository dan handler diinisialisasi di main.go dan di-inject ke modul
4. **Interface-Based**: Repository menggunakan interface sehingga mudah untuk testing dan swapping implementation
5. **RBAC System**: Sistem Role-Based Access Control yang terintegrasi dengan middleware permission
6. **Easy Scaling**: Menambah modul baru hanya perlu:
   - Buat folder baru di `internal/modules/`
   - Buat 4 file: model, repository, handler, route
   - Register route di main.go dengan permission middleware yang sesuai

## Prerequisites

- Go 1.21+
- PostgreSQL
- Swaggo CLI untuk generate dokumentasi

Install Swaggo:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Setup dan Menjalankan

1. **Clone dan setup dependencies:**
```bash
git clone <repository>
cd apiserver
go mod tidy
```

2. **Setup database:**
```bash
# Buat database PostgreSQL
createdb my_api_db

# Copy dan edit file environment
cp .env.example .env
# Edit .env sesuai konfigurasi database Anda
```

3. **Generate dokumentasi Swagger:**
```bash
swag init -g cmd/api/main.go -o docs
```

4. **Jalankan server:**
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:3000`

## API Documentation

Dokumentasi Swagger tersedia di: `http://localhost:3000/docs`

## Endpoints

### Access
- `GET /v1/profile` - Get user profile (Requires: profile:read)

### Examples
- `GET /v1/examples` - Get all active examples (Requires: examples:read)
- `POST /v1/examples` - Create new example (Requires: examples:create)
- `GET /v1/examples/:id` - Get example by ID (Requires: examples:read)
- `PUT /v1/examples/:id` - Update example (Requires: examples:update)
- `DELETE /v1/examples/:id` - Soft delete example (Requires: examples:delete)
- `POST /v1/examples/:id/restore` - Restore deleted example (Requires: examples:update)
- `GET /v1/examples/deleted` - Get all deleted examples (Requires: examples:read)

### Permissions Management
- `GET /v1/permissions` - Get all permissions (Requires: permissions:manage)
- `POST /v1/permissions` - Create new permission (Requires: permissions:manage)
- `GET /v1/permissions/:id` - Get permission by ID (Requires: permissions:manage)
- `DELETE /v1/permissions/:id` - Delete permission (Requires: permissions:manage)

### Groups Management
- `GET /v1/groups` - Get all groups (Requires: groups:manage)
- `POST /v1/groups` - Create new group (Requires: groups:manage)
- `GET /v1/groups/:id` - Get group by ID (Requires: groups:manage)
- `PUT /v1/groups/:id/permissions` - Update group permissions (Requires: groups:manage)
- `DELETE /v1/groups/:id` - Delete group (Requires: groups:manage)

### Audit Logs
- `GET /v1/audit-logs` - Get audit logs with filtering (Requires: audit:read)
- `GET /v1/audit-logs/:id` - Get detailed audit log by ID (Requires: audit:read)
- `DELETE /v1/audit-logs/cleanup?days=30` - Delete old audit logs (Requires: audit:manage)

### Health Check
- `GET /health` - Health check endpoint (No authentication required)
- `GET /version` - Get API version info (No authentication required)

## Contoh Penggunaan

### 1. Setup Test Data
Pertama, insert user dengan API key ke database:

```sql
INSERT INTO access (name, email, api_key, created_at, updated_at) 
VALUES ('John Doe', 'john@example.com', 'test-api-key-123', NOW(), NOW());
```

### 2. Test Endpoints dengan cURL

**Get Profile (Valid Token):**
```bash
curl -X GET "http://localhost:3000/v1/profile" \
  -H "Authorization: Bearer test-api-key-123"
```

**Get Examples (Valid Token):**
```bash
curl -X GET "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123"
```

**Create Example (Valid Token):**
```bash
curl -X POST "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Example",
    "description": "This is a test example"
  }'
```

**Test Invalid Token:**
```bash
curl -X GET "http://localhost:3000/v1/profile" \
  -H "Authorization: Bearer invalid-token"
```

**Test Missing Token:**
```bash
curl -X GET "http://localhost:3000/v1/profile"
```

### 3. Test Permission System dengan Different User Roles

**Admin User - Can Access Everything:**
```bash
# Get all permissions (Admin only)
curl -X GET "http://localhost:3000/v1/permissions" \
  -H "Authorization: Bearer admin-api-key-789"

# Create new permission (Admin only)
curl -X POST "http://localhost:3000/v1/permissions" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Delete Examples",
    "description": "Permission to delete examples",
    "resource": "examples",
    "action": "delete"
  }'

# Get all groups (Admin only)
curl -X GET "http://localhost:3000/v1/groups" \
  -H "Authorization: Bearer admin-api-key-789"

# Create new group (Admin only)
curl -X POST "http://localhost:3000/v1/groups" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Manager",
    "description": "Manager role with limited admin access",
    "permission_ids": [1, 2, 3]
  }'
```

**Editor User - Can Create/Read Examples:**
```bash
# Can create examples
curl -X POST "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Editor Example",
    "description": "Created by editor user"
  }'

# Can read examples
curl -X GET "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-123"

# CANNOT access permissions (will get 403 Forbidden)
curl -X GET "http://localhost:3000/v1/permissions" \
  -H "Authorization: Bearer test-api-key-123"
```

**Viewer User - Read Only Access:**
```bash
# Can read examples
curl -X GET "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-456"

# Can view profile
curl -X GET "http://localhost:3000/v1/profile" \
  -H "Authorization: Bearer test-api-key-456"

# CANNOT create examples (will get 403 Forbidden)
curl -X POST "http://localhost:3000/v1/examples" \
  -H "Authorization: Bearer test-api-key-456" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Viewer Example",
    "description": "This should fail"
  }'
```

### 4. Test Group Permission Management

**Update Group Permissions (Admin only):**
```bash
# Add more permissions to Editor group
curl -X PUT "http://localhost:3000/v1/groups/2/permissions" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2, 3, 4]
  }'

# Get specific group details
curl -X GET "http://localhost:3000/v1/groups/2" \
  -H "Authorization: Bearer admin-api-key-789"
```

## Menambah Modul Baru

Untuk menambah modul baru (misal: `product`):

1. Buat folder `internal/modules/product/`
2. Buat 4 file: `product_model.go`, `product_repository.go`, `product_handler.go`, `product_route.go`
3. Follow pattern yang sama seperti modul `example`
4. Register route di `main.go`:
```go
productRepo := product.NewRepository(db)
productHandler := product.NewHandler(productRepo)
product.RegisterProductRoutes(app, productHandler, authMiddleware)
```

## Fitur Utama

- ✅ **Modular Architecture**: Feature-based structure
- ✅ **Authentication**: Bearer token middleware
- ✅ **Database**: PostgreSQL dengan GORM
- ✅ **Documentation**: Auto-generated Swagger docs
- ✅ **Configuration**: Environment variables dengan .env
- ✅ **Error Handling**: Centralized error handling
- ✅ **CORS**: Cross-origin resource sharing
- ✅ **Logging**: Request logging middleware
- ✅ **Health Check**: Basic health check endpoint

## Development

Untuk development, gunakan air untuk hot reload:

```bash
go install github.com/cosmtrek/air@latest
air
```


## INIT

Menambahkan data test:




### Opsi #1

Melalui seeder:
```bash
go run cmd/api/main.go -seed
```

### Opsi #2

Setelah server berjalan dan tabel ter-migrate, jalankan query SQL ini di PostgreSQL:

```sql
INSERT INTO access (name, email, api_key, created_at, updated_at)
VALUES
  ('John Doe', 'john@example.com', 'test-api-key-123', NOW(), NOW()),
  ('Jane Smith', 'jane@example.com', 'test-api-key-456', NOW(), NOW()),
  ('Admin User', 'admin@example.com', 'admin-api-key-789', NOW(), NOW());
```

## API Keys yang Tersedia untuk Testing:

### Admin User (Full Access)
- **API Key**: `admin-api-key-789`
- **Email**: admin@example.com
- **Group**: Admin
- **Permissions**: All permissions (create, read, update, delete examples + manage permissions & groups)

### Editor User (Limited Access)
- **API Key**: `test-api-key-123`
- **Email**: john@example.com
- **Group**: Editor
- **Permissions**: Create, read, update examples + view profile

### Viewer User (Read Only)
- **API Key**: `test-api-key-456`
- **Email**: jane@example.com
- **Group**: Viewer
- **Permissions**: Read examples + view profile

## 🔍 Audit Logging System

API ini dilengkapi dengan sistem audit logging yang komprehensif untuk mencatat semua aktivitas API:

### Fitur Audit Logging:
- ✅ **Automatic Logging**: Semua request/response dicatat otomatis
- ✅ **Request Details**: Method, path, headers, body payload
- ✅ **Response Details**: Status code, response body, response time
- ✅ **User Tracking**: User ID, email, dan API key (masked)
- ✅ **IP & User Agent**: Tracking untuk security analysis
- ✅ **Filtering & Search**: Filter berdasarkan user, method, path, status, tanggal
- ✅ **Pagination**: Support untuk large datasets
- ✅ **Cleanup**: Auto-delete old logs untuk maintenance

### Contoh Penggunaan Audit Logs:

**Get Audit Logs (Admin only):**
```bash
# Get all audit logs
curl -X GET "http://localhost:3000/v1/audit-logs" \
  -H "Authorization: Bearer admin-api-key-789"

# Filter by user email
curl -X GET "http://localhost:3000/v1/audit-logs?user_email=john@example.com" \
  -H "Authorization: Bearer admin-api-key-789"

# Filter by method and status code
curl -X GET "http://localhost:3000/v1/audit-logs?method=POST&status_code=201" \
  -H "Authorization: Bearer admin-api-key-789"

# Filter by date range
curl -X GET "http://localhost:3000/v1/audit-logs?date_from=2024-01-01&date_to=2024-01-31" \
  -H "Authorization: Bearer admin-api-key-789"

# With pagination
curl -X GET "http://localhost:3000/v1/audit-logs?limit=10&offset=20" \
  -H "Authorization: Bearer admin-api-key-789"
```

**Get Detailed Audit Log:**
```bash
# Get specific audit log with full request/response details
curl -X GET "http://localhost:3000/v1/audit-logs/123" \
  -H "Authorization: Bearer admin-api-key-789"
```

**Cleanup Old Logs (Admin only):**
```bash
# Delete logs older than 30 days
curl -X DELETE "http://localhost:3000/v1/audit-logs/cleanup?days=30" \
  -H "Authorization: Bearer admin-api-key-789"
```

### Audit Log Data Structure:
```json
{
  "id": 123,
  "user_id": 1,
  "user_email": "john@example.com",
  "api_key": "test-api****", // Masked for security
  "method": "POST",
  "path": "/v1/examples",
  "status_code": 201,
  "request_headers": "{\"Content-Type\":\"application/json\"}",
  "request_body": "{\"name\":\"Test\",\"description\":\"Test example\"}",
  "response_body": "{\"status\":\"success\",\"data\":{...}}",
  "response_time": 45, // milliseconds
  "ip_address": "192.168.1.100",
  "user_agent": "curl/7.68.0",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Security Features:
- 🔒 **Sensitive Data Protection**: API keys dan sensitive headers di-mask
- 🔒 **Permission-Based Access**: Hanya admin yang bisa akses audit logs
- 🔒 **Response Size Limiting**: Response body dibatasi untuk mencegah log yang terlalu besar
- 🔒 **Async Logging**: Logging dilakukan secara asynchronous untuk performa optimal

## 🗂️ Status Management System

API ini menggunakan sistem status management dengan field `status_id` di setiap tabel untuk mengelola lifecycle data:

### Status Values:
- **0** - Active (Default, data aktif)
- **1** - Inactive/Deleted (Soft deleted)
- **2** - Pending (Untuk future use)
- **3** - Suspended (Untuk future use)

### Fitur Status Management:
- ✅ **Soft Delete**: Data tidak benar-benar dihapus, hanya diubah status_id menjadi 1
- ✅ **Auto Filtering**: Semua query otomatis memfilter data aktif (status_id = 0)
- ✅ **Restore Capability**: Data yang di-soft delete bisa di-restore kembali
- ✅ **Audit Trail**: Perubahan status tercatat dalam audit logs
- ✅ **Consistent Implementation**: Semua tabel menggunakan pattern yang sama

### Contoh Penggunaan Status Management:

**Soft Delete Example:**
```bash
# Soft delete example (set status_id = 1)
curl -X DELETE "http://localhost:3000/v1/examples/1" \
  -H "Authorization: Bearer admin-api-key-789"
```

**Restore Example:**
```bash
# Restore deleted example (set status_id = 0)
curl -X POST "http://localhost:3000/v1/examples/1/restore" \
  -H "Authorization: Bearer admin-api-key-789"
```

**Get Deleted Examples:**
```bash
# View all soft deleted examples
curl -X GET "http://localhost:3000/v1/examples/deleted" \
  -H "Authorization: Bearer admin-api-key-789"
```

**Update Example:**
```bash
# Update existing example
curl -X PUT "http://localhost:3000/v1/examples/1" \
  -H "Authorization: Bearer test-api-key-123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Example",
    "description": "This example has been updated"
  }'
```

### Database Schema dengan Status ID:
```sql
-- Contoh struktur tabel dengan status_id
CREATE TABLE examples (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    status_id SMALLINT NOT NULL DEFAULT 1
);

-- Index untuk performa query
CREATE INDEX idx_examples_status_id ON examples(status_id);

-- Note: Default value is 1 (inactive), data must be explicitly set to 0 to be active
```

### Benefits Status Management:
- 🔄 **Data Recovery**: Data yang terhapus bisa dipulihkan
- 📊 **Analytics**: Bisa menganalisis data yang dihapus
- 🔍 **Audit Trail**: Jejak lengkap perubahan status
- ⚡ **Performance**: Query lebih cepat dengan proper indexing
- 🛡️ **Data Integrity**: Mencegah kehilangan data permanen





## 📚 Swagger Documentation Configuration

API ini memiliki sistem konfigurasi fleksibel untuk menampilkan/menyembunyikan module tertentu dari dokumentasi Swagger:

### **Default Behavior:**
- ✅ **Examples** - Selalu ditampilkan
- ✅ **Permissions** - Selalu ditampilkan
- ✅ **Groups** - Selalu ditampilkan
- ❌ **Audit** - Disembunyikan (internal use)
- ❌ **Access/Profile** - Disembunyikan (internal use)

### **Cara Mengontrol Dokumentasi:**

**1. Menggunakan Script (Recommended):**
```bash
# Generate docs tanpa Audit & Access (default)
./scripts/generate-swagger.sh

# Generate docs dengan Audit module
./scripts/generate-swagger.sh --show-audit

# Generate docs dengan Access module
./scripts/generate-swagger.sh --show-access

# Generate docs dengan semua module
./scripts/generate-swagger.sh --show-all

# Windows
scripts\generate-swagger.bat --show-all
```

**2. Menggunakan Environment Variables:**
```bash
# Set di .env file
SHOW_AUDIT_IN_SWAGGER=true
SHOW_ACCESS_IN_SWAGGER=true

# Atau set saat runtime
SHOW_AUDIT_IN_SWAGGER=true swag init -g cmd/api/main.go -o docs
```

**3. Menggunakan Build Tags:**
```bash
# Manual dengan build tags
swag init -g cmd/api/main.go -o docs --tags "swagger_audit,swagger_access"
```

### **Use Cases:**

**🔒 Production (Default):**
- Sembunyikan Audit & Access untuk keamanan
- Hanya tampilkan public API endpoints

**🛠️ Development:**
- Tampilkan semua module untuk testing
- Full API documentation untuk developer

**📋 Documentation:**
- Selective display berdasarkan audience
- Internal vs External documentation

### **Benefits:**
- 🔐 **Security**: Sensitive endpoints tidak ter-expose di public docs
- 🎯 **Focused**: Documentation sesuai kebutuhan audience
- 🔄 **Flexible**: Easy toggle tanpa mengubah code
- 📱 **Maintainable**: Single source of truth untuk docs#
# 🔑 API Key Expiration System

API ini dilengkapi dengan sistem manajemen API key expiration yang komprehensif:

### Fitur Expiration:
- ✅ **Flexible Expiration**: API key bisa diatur untuk expired pada tanggal tertentu
- ✅ **Never Expires**: API key bisa diatur untuk tidak pernah expired (NULL value)
- ✅ **Auto Validation**: API key yang sudah expired otomatis ditolak
- ✅ **Management API**: Endpoint untuk mengatur dan menghapus expiration date
- ✅ **Permission Based**: Hanya admin dengan permission "access:manage" yang bisa mengatur

### Contoh Penggunaan:

**Set Expiration Date:**
```bash
# Set API key untuk expired dalam 30 hari
curl -X PUT "http://localhost:3000/v1/access/1/expired-date" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "expired_date": "2025-08-22T00:00:00Z"
  }'
```

**Remove Expiration (Never Expires):**
```bash
# Set API key untuk tidak pernah expired
curl -X DELETE "http://localhost:3000/v1/access/1/expired-date" \
  -H "Authorization: Bearer admin-api-key-789"
```

### Contoh Data Seeded:
- **Admin User**: Tidak pernah expired (NULL)
- **John Doe**: Expired dalam 3 bulan dari sekarang
- **Jane Smith**: Sudah expired (1 bulan yang lalu)

### Implementasi:
- Field `expired_date` di tabel `access` (nullable)
- Check expired date di auth middleware
- Endpoint management dengan permission control
- Validasi tanggal (harus di masa depan)

### Security Benefits:
- 🔒 **Temporary Access**: Bisa memberikan akses sementara
- 🔒 **Auto Revocation**: API key expired otomatis tanpa manual revoke
- 🔒 **Audit Trail**: Semua perubahan expiration tercatat di audit logs
- 🔒 **Granular Control**: Bisa mengatur expiration per user## 🚦 API R
ate Limiting System

API ini dilengkapi dengan sistem rate limiting yang komprehensif untuk mengontrol jumlah request per API key:

### Fitur Rate Limiting:
- ✅ **Per-User Limits**: Setiap API key memiliki rate limit sendiri
- ✅ **Default Limit**: 120 requests per menit (configurable)
- ✅ **Custom Limits**: Bisa diatur per user sesuai kebutuhan
- ✅ **Rate Limit Headers**: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
- ✅ **Graceful Rejection**: 429 Too Many Requests dengan pesan yang jelas

### Contoh Penggunaan:

**Update Rate Limit:**
```bash
# Set rate limit untuk user (300 requests per menit)
curl -X PUT "http://localhost:3000/v1/access/1/rate-limit" \
  -H "Authorization: Bearer admin-api-key-789" \
  -H "Content-Type: application/json" \
  -d '{
    "rate_limit": 300
  }'
```

**Response Headers:**
```
X-RateLimit-Limit: 120
X-RateLimit-Remaining: 119
X-RateLimit-Reset: 1627484861
```

### Contoh Data Seeded:
- **Admin User**: 1000 requests per menit
- **John Doe**: 120 requests per menit (default)
- **Jane Smith**: 60 requests per menit (limited)

### Implementasi:
- Field `rate_limit` di tabel `access`
- In-memory tracking untuk request timestamps
- Middleware untuk validasi rate limit
- Endpoint management dengan permission control

### Benefits:
- 🛡️ **DDoS Protection**: Mencegah abuse dari single client
- 💰 **Cost Control**: Membatasi resource usage
- 🎯 **Tiered Access**: Bisa memberikan limit berbeda per user tier
- 📊 **Usage Insights**: Monitoring request patterns
- 🔄 **Fair Usage**: Mencegah satu client menghabiskan resource