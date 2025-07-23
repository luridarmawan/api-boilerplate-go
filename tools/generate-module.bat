@echo off
REM Module Generator Script for Windows
REM Usage: generate-module.bat <module-name> [--with-permissions]

if "%1"=="" (
    echo Usage: generate-module.bat ^<module-name^> [--with-permissions]
    echo Example: generate-module.bat product
    echo Example: generate-module.bat product --with-permissions
    exit /b 1
)

echo Generating module: %1
if "%2"=="--with-permissions" (
    echo With permissions enabled
    go run tools/module-generator/main.go %1 --with-permissions
) else (
    go run tools/module-generator/main.go %1
)