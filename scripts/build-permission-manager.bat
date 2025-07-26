@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul

REM Build script for permission-manager CLI tool (Windows)
REM This script builds the permission manager for multiple platforms

echo 🔨 Building Permission Manager CLI Tool...

REM Get the current directory
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..
set BUILD_DIR=%PROJECT_ROOT%\bin

REM Create build directory if it doesn't exist
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

REM Build information
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=dev

for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT=%%i
if "%COMMIT%"=="" set COMMIT=unknown

for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'"') do set BUILD_DATE=%%i

REM Build flags
set LDFLAGS=-X 'main.Version=%VERSION%' -X 'main.GitCommit=%COMMIT%' -X 'main.BuildDate=%BUILD_DATE%'

echo 📦 Version: %VERSION%
echo 📦 Commit: %COMMIT%
echo 📦 Build Date: %BUILD_DATE%

REM Change to project root
cd /d "%PROJECT_ROOT%"

REM Build for current platform (Windows)
echo 🏗️  Building for current platform...
go build -ldflags "%LDFLAGS%" -o "%BUILD_DIR%\permission-manager.exe" .\cmd\permission-manager

REM Build for multiple platforms
echo 🏗️  Building for multiple platforms...

REM Windows
echo   - Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o "%BUILD_DIR%\permission-manager-windows-amd64.exe" .\cmd\permission-manager

REM Linux
echo   - Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o "%BUILD_DIR%\permission-manager-linux-amd64" .\cmd\permission-manager

REM macOS
echo   - Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o "%BUILD_DIR%\permission-manager-darwin-amd64" .\cmd\permission-manager

REM macOS ARM64 (Apple Silicon)
echo   - Building for macOS (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o "%BUILD_DIR%\permission-manager-darwin-arm64" .\cmd\permission-manager

echo ✅ Build completed successfully!
echo.
echo 📁 Binaries created in: %BUILD_DIR%
echo    - permission-manager.exe (Windows)
echo    - permission-manager-windows-amd64.exe
echo    - permission-manager-linux-amd64
echo    - permission-manager-darwin-amd64
echo    - permission-manager-darwin-arm64
echo.
echo 🚀 Usage example:
echo    %BUILD_DIR%\permission-manager.exe 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create

pause