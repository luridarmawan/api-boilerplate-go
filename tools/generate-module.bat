@echo off
REM Module Generator Script for Windows
REM Usage: generate-module.bat <module-name>

if "%1"=="" (
    echo Usage: generate-module.bat ^<module-name^>
    echo Example: generate-module.bat product
    exit /b 1
)

echo Generating module: %1
go run tools/module-generator/main.go %1