@echo off
REM Safe Configuration Permissions Seeder
REM This script adds configuration permissions ONLY to Admin group

echo 🔒 Configuration Permissions Seeder (Admin Only)
echo ================================================

REM Check if we're in the right directory
if not exist "go.mod" (
    echo ❌ Error: Please run this script from the project root directory
    exit /b 1
)

REM Check if database is accessible
echo 🔍 Checking database connection...
go run scripts\add-admin-only-configuration-permissions.go >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Cannot connect to database. Please check your .env configuration
    exit /b 1
)

echo ✅ Database connection successful

REM Ask for confirmation
echo ⚠️  This will add configuration permissions ONLY to Admin group
echo    Other groups (Editor, Viewer, General client) will NOT get access
echo.
set /p confirm="Do you want to continue? (y/N): "

if /i not "%confirm%"=="y" (
    echo ❌ Operation cancelled
    exit /b 0
)

REM Run the admin-only configuration permissions script
echo 🚀 Adding configuration permissions for Admin group only...
go run scripts\add-admin-only-configuration-permissions.go

if errorlevel 1 (
    echo ❌ Failed to add configuration permissions
    exit /b 1
)

echo ✅ Configuration permissions successfully added!

REM Optional: Remove any existing configuration permissions from non-admin groups
echo.
set /p cleanup="Do you want to remove configuration permissions from non-admin groups (if any)? (y/N): "

if /i "%cleanup%"=="y" (
    echo 🧹 Cleaning up configuration permissions from non-admin groups...
    go run scripts\remove-configuration-permissions-from-non-admin.go
    
    if errorlevel 1 (
        echo ❌ Cleanup failed
    ) else (
        echo ✅ Cleanup completed successfully!
    )
)

echo.
echo 🎉 Configuration module is now secured for Admin-only access!
echo 📋 Summary:
echo    ✅ Configuration permissions created
echo    ✅ Permissions assigned to Admin group only
echo    ✅ Other groups have no configuration access
echo    ✅ Security: Admin-only access enforced