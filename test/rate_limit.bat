@echo off
REM Buat script untuk mengirim banyak request ke API
for /l %%i in (1,1,150) do (
  curl -s -X GET "http://localhost:3000/v1/examples" ^
    -H "Authorization: Bearer test-api-key-123"
  echo Request %%i
  timeout /t 0 /nobreak >nul
  ping -n 1 127.0.0.1 >nul
)
