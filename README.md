# 🧩 REST API Modular Boilerplate (Go)

![Go Version](https://img.shields.io/badge/go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-brightgreen)
![Build](https://img.shields.io/badge/build-passing-success)
![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)

## 📑 Table of Contents

- [🧩 Features](#-features)
- [🚀 Quickstart](#-quickstart)
- [📂 Project Structure](#-project-structure)
- [🧱 Architecture](#-architecture)
- [🛡️ Security & Auth](#-security--auth)
- [📈 Audit Logging System](#-audit-logging-system)
- [🔄 Status Management System](#-status-management-system)
- [⏳ API Key Expiration](#-api-key-expiration-system)
- [🚦 API Rate Limiting](#-api-rate-limiting-system)
- [🧪 Testing](#-testing)
- [🧩 Add New Module](#-add-new-module)
- [🧰 Development Tools](#-development-tools)
- [🗃️ Seeder & Test Data](#%EF%B8%8F-seeder--test-data)
- [📖 API Documentation](#-api-documentation)

## 🧩 Features

- ✅ **Modular Architecture**: Feature-based directory structure
- ✅ **Authentication System**: Bearer token with expiration & rate limiting
- ✅ **Database Layer**: PostgreSQL with GORM
- ✅ **RBAC**: Role-based access control with permission mapping
- ✅ **Documentation**: Auto-generated Swagger documentation
- ✅ **Configuration**: Environment-based setup using `.env`
- ✅ **Error Handling**: Centralized and consistent error responses
- ✅ **Audit Log**: Tracks user requests and responses
- ✅ **Status Management**: Soft deletion using `status_id`
- ✅ **Seeder & Sample Data**: Default test data for quick setup
- ✅ **Health Check**: Built-in endpoint to check server status
- ✅ **Expiration System**: Supports token/key expiration policy
- ✅ **Rate Limit**:  Controls number of requests per user/IP

## 🚀 Quickstart

```bash
git clone https://github.com/your-org/api-boilerplate-go.git apiserver
cd apiserver

# Copy dan edit file environment
cp .env.example .env

go mod tidy

go run cmd/api/main.go --seed
go run cmd/api/main.go
#air
```

## 📂 Project Structure

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

## 🧱 Architecture

Each module includes:

- `model.go`: Defines schema and data validation
- `repository.go`: Contains GORM-based query logic
- `handler.go`: Implements the core API logic
- `route.go`: Registers routes using Fiber

```
client → route → handler → repository → db
```

## Why a Modular Structure?

1. **Feature-Based Structure**: Each feature (e.g., `access`, `permission`, `group`, `example`) resides in its own folder containing all required components.
2. **Self-Contained**: Every module is independent and includes its own model, repository, handler, and route, making it easy to maintain and scale.
3. **Dependency Injection**: Repositories and handlers are initialized in `main.go` and injected into modules, ensuring loose coupling and better testability.
4. **Interface-Based**: Repositories are defined via interfaces, allowing for easier testing and implementation swapping
5. **RBAC System**: Role-Based Access Control is built-in and integrated with permission middleware for fine-grained access control
6. **Easy Scaling**: Adding a new module is straightforward:
   - Create a new folder under `internal/modules/`
   - Add four files: `model.go`, `repository.go`, `handler.go`, `route.go`
   - Register the route in `main.go` and attach appropriate permission middleware

## Prerequisites

- Go 1.24+
- PostgreSQL
- Swaggo CLI untuk generate dokumentasi

Install Swaggo:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🛡️ Security & Auth

- **Bearer Token** with expiration (`token_expired_at`)
- **Rate Limiting**: Protect brute force (30/min)
- **RBAC**: Per-role permission check

## 📈 Audit Logging System

Log detail:
- User ID
- Endpoint + Method
- Status Code
- Duration & Response
- Saved to PostgreSQL

## 🔄 Status Management System

Gunakan `status_id`:
- `1`: Active (default)
- `0`: Deleted (soft delete)

Contoh soft-delete:
```sql
UPDATE examples SET status_id = 0 WHERE id = 'uuid';
```

## ⏳ API Key Expiration System

Tambahkan `token_expired_at` di table `users`. Token akan ditolak jika waktu sekarang melewati tanggal expired.

Contoh:
```sql
UPDATE users SET token_expired_at = NOW() + INTERVAL '2 days';
```

## 🚦 API Rate Limiting System

Gunakan middleware `ratelimit` untuk membatasi:
- 120 request per menit per IP
- Bisa dikustom per masing-masing key

## 🧪 Testing

- Gunakan token user untuk testing autentikasi

## 🧩 Add New Module

1. Duplikat folder `example/` di `internal/modules/`
2. Ganti nama file & isi sesuai kebutuhan
3. Tambahkan router di `router/route.go`

## 🧰 Development Tools

- Hot reload dengan `air`
- Linter: `golangci-lint`
- Database: `PostgreSQL`

## 🗃️ Seeder & Test Data

Jalankan:

```bash
go run cmd/seed/main.go
```

Data awal akan disimpan ke:
- `access`
- `groups`
- `group_permissions`
- `examples`

## 📚 API Documentation

Auto-generated Swagger docs from Fiber comments:

```go
// @Summary Get all examples
// @Tags example
// @Produce json
// @Success 200 {object} model.Example
// @Router /example [get]
```

---

**How to build documentation:**
```bash
swag init -g cmd/api/main.go -o docs
```

API Documentation will available at [http://localhost:3000/docs](http://localhost:3000/docs)

- API Example:

```bash
curl -X GET http://localhost:3000/api/v1/example -H "Authorization: Bearer {token}"
```


> MIT Licensed · Built with ❤️ by [CARIK.id](https://carik.id) team
