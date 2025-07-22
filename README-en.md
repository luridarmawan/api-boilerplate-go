# 🧩 REST API Modular Boilerplate (Go)

![Go Version](https://img.shields.io/badge/go-1.22+-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-brightgreen)
![Build](https://img.shields.io/badge/build-passing-success)
![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)

## 📑 Table of Contents

- [🧩 Features](#-features)
- [🚀 Quickstart](#-quickstart)
- [📂 Project Structure](#-project-structure)
- [🧱 Architecture](#-architecture)
- [📖 API Reference](#-api-reference)
- [🛡️ Security & Auth](#-security--auth)
- [📈 Audit Logging System](#-audit-logging-system)
- [🔄 Status Management System](#-status-management-system)
- [⏳ API Key Expiration](#-api-key-expiration-system)
- [🚦 API Rate Limiting](#-api-rate-limiting-system)
- [🧪 Testing](#-testing)
- [🧩 Add New Module](#-add-new-module)
- [🧰 Development Tools](#-development-tools)
- [🗃️ Seeder & Test Data](#-seeder--test-data)
- [📚 Swagger & Docs](#-swagger--docs)

## 🧩 Features

- ✅ **Modular Architecture**: Feature-based directory structure
- ✅ **Authentication System**: Bearer token with expiration & rate limiting
- ✅ **RBAC**: Role & permission based access
- ✅ **Database Layer**: PostgreSQL with GORM
- ✅ **Documentation**: Auto-generated Swagger docs
- ✅ **Configuration**: Environment variables with `.env`
- ✅ **Error Handling**: Centralized error handler
- ✅ **Audit Log**: Request/response logger by user
- ✅ **Status Management**: Logical delete via `status_id`
- ✅ **Seeder & Sample Data**: With UUID & timestamps
- ✅ **Hot Reload**: With [air](https://github.com/cosmtrek/air)
- ✅ **Testing with Postman & cURL**

## 🚀 Quickstart

```bash
git clone https://github.com/your-org/api-boilerplate-go.git
cd api-boilerplate-go
cp .env.example .env
go mod tidy
air
```

## 📂 Project Structure

```
.
├── cmd/server/main.go
├── internal
│   ├── config/
│   ├── db/
│   ├── middleware/
│   ├── modules/
│   │   ├── example/
│   │   ├── group/
│   │   └── user/
│   ├── router/
│   └── utils/
└── docs/
```

## 🧱 Architecture

Setiap modul memiliki:

- `model.go`: Schema dan validasi
- `repository.go`: Query builder (GORM)
- `handler.go`: Logic utama API
- `route.go`: Router module (Fiber)

```
client → route → handler → repository → db
```

## 📖 API Reference

- Swagger: [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)
- Contoh API:

```bash
curl -X GET http://localhost:3000/api/v1/example -H "Authorization: Bearer {token}"
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
- 30 request per menit per IP
- Bisa dikustom via ENV

## 🧪 Testing

- Postman Collection tersedia di folder `docs/`
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
- `users`
- `groups`
- `group_permissions`
- `examples`

## 📚 Swagger & Docs

Swagger doc otomatis dari comment Fiber:

```go
// @Summary Get all examples
// @Tags example
// @Produce json
// @Success 200 {object} model.Example
// @Router /example [get]
```

---

> MIT Licensed · Built with ❤️ by your-team