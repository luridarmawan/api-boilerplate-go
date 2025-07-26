# Webhook Tester

Tool untuk menguji fungsi `CallWebhook` di `internal/utils/net.go`.

## Fitur

- **Server Mode**: Menjalankan HTTP server untuk menerima webhook
- **Test Mode**: Menguji fungsi CallWebhook dengan berbagai jenis data
- **Quick Mode**: Menjalankan server dan test secara otomatis

## Cara Penggunaan

### 1. Persiapan

Pastikan Go sudah terinstall. Tool ini sudah standalone dan tidak memerlukan dependency eksternal.

```bash
cd tools/webhook-tester
```

### 2. Mode Server (Menerima Webhook)

Jalankan server untuk menerima webhook:

```bash
# Linux/Mac
./test-webhook.sh server [port]

# Windows
test-webhook.bat server [port]

# Atau langsung dengan Go
go run main.go server [port]
```

Contoh:
```bash
# Server di port 8080 (default)
./test-webhook.sh server

# Server di port 9000
./test-webhook.sh server 9000
```

Server akan menampilkan semua webhook yang diterima dengan detail:
- Timestamp
- Headers HTTP
- Body JSON

### 3. Mode Test (Mengirim Webhook)

Test fungsi CallWebhook dengan URL webhook:

```bash
# Linux/Mac
./test-webhook.sh test <webhook-url>

# Windows
test-webhook.bat test <webhook-url>

# Atau langsung dengan Go
go run main.go test <webhook-url>
```

Contoh:
```bash
# Test dengan server lokal
./test-webhook.sh test http://localhost:8080/webhook

# Test dengan webhook.site
./test-webhook.sh test https://webhook.site/your-unique-id

# Test dengan ngrok
./test-webhook.sh test https://abc123.ngrok.io/webhook
```

### 4. Mode Quick (Test Cepat)

Menjalankan server dan test secara otomatis:

```bash
# Linux/Mac
./test-webhook.sh quick

# Windows
test-webhook.bat quick
```

## Test Cases

Tool ini akan mengirim 3 jenis data berbeda:

1. **Simple Test Data**: Struct dengan ID, message, timestamp, dan status
2. **Map Data**: Data dalam bentuk map dengan nested objects
3. **Array Data**: Array berisi multiple objects

## Contoh Output

### Server Mode
```
=== Webhook Received ===
Time: 2024-01-15 10:30:45
Headers:
  Content-Type: application/json
  User-Agent: Webhook-Client/1.0
Body: {"id":1,"message":"Hello from webhook tester!","timestamp":"2024-01-15T10:30:45Z","status":"success"}
=====================
```

### Test Mode
```
--- Test 1: Simple Test Data ---
Sending data:
{
  "id": 1,
  "message": "Hello from webhook tester!",
  "timestamp": "2024-01-15T10:30:45Z",
  "status": "success"
}
Webhook called (async)
```

## Tips Penggunaan

1. **Testing dengan External Services**: 
   - Gunakan [webhook.site](https://webhook.site) untuk testing cepat
   - Gunakan [ngrok](https://ngrok.com) untuk expose local server

2. **Debugging**:
   - Server mode menampilkan semua detail request
   - Fungsi CallWebhook berjalan asynchronous, jadi tunggu sebentar untuk melihat hasil

3. **Production Testing**:
   - Test dengan berbagai jenis payload
   - Test dengan URL yang tidak valid untuk melihat error handling
   - Test dengan server yang lambat untuk melihat timeout behavior

## Troubleshooting

1. **Port Already in Use**:
   - Gunakan port yang berbeda: `./test-webhook.sh server 9000`

3. **Permission Denied (Linux/Mac)**:
   ```bash
   chmod +x test-webhook.sh
   ```

## Customization

Anda bisa memodifikasi `main.go` untuk:
- Menambah test cases baru
- Mengubah format output
- Menambah authentication headers
- Menambah retry logic testing

## Catatan Penting

Tool ini menggunakan copy dari fungsi `CallWebhook` yang ada di `internal/utils/net.go`. Jika Anda mengubah implementasi asli, pastikan untuk menyinkronkan perubahan di tool tester ini juga.