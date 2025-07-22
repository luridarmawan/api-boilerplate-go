@echo off
REM Script to generate Swagger documentation with conditional modules
REM Usage: scripts\generate-swagger.bat [options]
REM Options:
REM   --show-audit    Include Audit module in documentation
REM   --show-access   Include Access module in documentation
REM   --show-all      Include both Audit and Access modules
REM   --help          Show this help message

set SHOW_AUDIT=false
set SHOW_ACCESS=false

:parse_args
if "%~1"=="" goto generate_docs
if "%~1"=="--show-audit" (
    set SHOW_AUDIT=true
    shift
    goto parse_args
)
if "%~1"=="--show-access" (
    set SHOW_ACCESS=true
    shift
    goto parse_args
)
if "%~1"=="--show-all" (
    set SHOW_AUDIT=true
    set SHOW_ACCESS=true
    shift
    goto parse_args
)
if "%~1"=="--help" (
    echo Usage: %0 [options]
    echo Options:
    echo   --show-audit    Include Audit module in documentation
    echo   --show-access   Include Access module in documentation
    echo   --show-all      Include both Audit and Access modules
    echo   --help          Show this help message
    exit /b 0
)
echo Unknown option: %~1
echo Use --help for usage information
exit /b 1

:generate_docs
echo Generating Swagger documentation...
echo Show Audit: %SHOW_AUDIT%
echo Show Access: %SHOW_ACCESS%

REM Create backup directory
if not exist .swagger-backup mkdir .swagger-backup

REM Backup original files
copy internal\modules\audit\audit_handler.go .swagger-backup\ >nul
copy internal\modules\access\access_handler.go .swagger-backup\ >nul

REM Process files based on options
if "%SHOW_AUDIT%"=="false" (
    echo Hiding Audit module from documentation...
    powershell -Command "(Get-Content internal\modules\audit\audit_handler.go) | Where-Object { $_ -notmatch 'SWAGGER_AUDIT_START' -and $_ -notmatch 'SWAGGER_AUDIT_END' -and -not ($_ -match 'SWAGGER_AUDIT_START' -or $inBlock) } | Set-Content internal\modules\audit\audit_handler.go"
)

if "%SHOW_ACCESS%"=="false" (
    echo Hiding Access module from documentation...
    powershell -Command "(Get-Content internal\modules\access\access_handler.go) | Where-Object { $_ -notmatch 'SWAGGER_ACCESS_START' -and $_ -notmatch 'SWAGGER_ACCESS_END' -and -not ($_ -match 'SWAGGER_ACCESS_START' -or $inBlock) } | Set-Content internal\modules\access\access_handler.go"
)

REM Generate Swagger docs
echo Generating documentation...
swag init -g cmd/api/main.go -o docs

REM Restore original files
echo Restoring original files...
copy .swagger-backup\audit_handler.go internal\modules\audit\ >nul
copy .swagger-backup\access_handler.go internal\modules\access\ >nul

REM Cleanup
rmdir /s /q .swagger-backup

echo Swagger documentation generated successfully!
echo Access documentation at: http://localhost:3000/docs