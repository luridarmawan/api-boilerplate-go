@echo off
setlocal enabledelayedexpansion

echo === Webhook Tester ===
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed or not in PATH
    exit /b 1
)

REM Check arguments
if "%1"=="" (
    call :show_usage
    exit /b 1
)

if "%1"=="server" (
    set PORT=%2
    if "!PORT!"=="" set PORT=8080
    echo Starting webhook receiver server on port !PORT!...
    echo Webhook URL: http://localhost:!PORT!/webhook
    echo Press Ctrl+C to stop
    echo.
    go run main.go server !PORT!
    goto :eof
)

if "%1"=="test" (
    if "%2"=="" (
        echo Error: Please provide webhook URL
        echo Example: %0 test http://localhost:8080/webhook
        exit /b 1
    )
    
    set WEBHOOK_URL=%2
    echo Testing CallWebhook function...
    echo Target URL: !WEBHOOK_URL!
    echo.
    go run main.go test "!WEBHOOK_URL!"
    goto :eof
)

if "%1"=="quick" (
    echo Quick test mode: Starting server and running test
    echo.
    
    REM Start server in background
    echo Starting webhook receiver server on port 8080...
    start /b go run main.go server 8080
    
    REM Wait for server to start
    echo Waiting for server to start...
    timeout /t 3 /nobreak >nul
    
    REM Run test
    echo Running webhook test...
    go run main.go test http://localhost:8080/webhook
    
    echo.
    echo Test completed!
    echo Note: You may need to manually stop the server process
    goto :eof
)

echo Error: Unknown command '%1'
echo.
call :show_usage
exit /b 1

:show_usage
echo Usage:
echo   %0 server [port]           - Start webhook receiver server (default port: 8080)
echo   %0 test ^<webhook-url^>      - Test CallWebhook function
echo   %0 quick                   - Start server and run test (uses localhost:8080)
echo.
echo Examples:
echo   %0 server 9000                                    # Start server on port 9000
echo   %0 test http://localhost:8080/webhook             # Test with local server
echo   %0 test https://webhook.site/your-unique-id       # Test with external service
echo   %0 quick                                          # Quick test with local server
goto :eof