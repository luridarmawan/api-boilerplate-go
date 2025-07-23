# Module Generator


## Penggunaan

### 1. Generate Modul Customer
```bash
go run tools/module-generator/main.go customer --with-permissions
```

**Hasil**: ✅ Berhasil tanpa error

### 2. Generate Modul Order
```bash
go run tools/module-generator/main.go order --with-permissions
```

**Hasil**: ✅ Berhasil tanpa error

### 3. Build Test
```bash
go build ./cmd/api -o api.exe
```

**Hasil**: ✅ Compile berhasil tanpa syntax error

### 4. Permission Script Test
```bash
go run scripts/add-customer-permissions.go
```

**Hasil**: ✅ Permission berhasil ditambahkan ke database

## File yang Dibuat

- `internal/modules/customer/` - Modul customer dibuat dengan template yang benar
- `scripts/add-customer-permissions.go` - Script permission customer
- `scripts/add-customer-permissions.sql` - Script SQL permission customer
- `test/customer-api-test.http` - File test HTTP customer

## Hasil Akhir

Setelah perbaikan:
- ✅ Module generator berfungsi dengan baik
- ✅ Template repository menghasilkan syntax yang benar
- ✅ Permission script berfungsi dengan baik
- ✅ Build aplikasi berhasil tanpa error
- ✅ Semua fitur `--with-permissions` berfungsi normal

## Catatan Penting

Sebelum menjalankan permission script, pastikan:
1. Database sudah running
2. Seeder sudah dijalankan untuk membuat tabel dan data awal:
   ```bash
   go run ./cmd/api --seed
   ```
3. Grup Admin, Editor, Viewer, dan General client sudah ada di database

## Testing

Untuk memastikan module generator bekerja dengan baik:

1. **Test generate modul baru**:
   ```bash
   go run tools/module-generator/main.go product --with-permissions
   ```

2. **Test compile**:
   ```bash
   go build -o api.exe ./cmd/api
   ```

3. **Test permission script**:
   ```bash
   go run scripts/add-product-permissions.go
   ```

Semua langkah di atas harus berhasil tanpa error.