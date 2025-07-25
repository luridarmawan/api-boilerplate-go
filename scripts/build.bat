@echo off
REM Build script for different environments
REM Usage: scripts\build.bat [dev|staging|prod]
chcp 65001 >nul

set ENVIRONMENT=%1
if "%ENVIRONMENT%"=="" set ENVIRONMENT=dev

echo 🔒 Building for environment: %ENVIRONMENT%

REM Generate Swagger documentation
echo 📋 Generating Swagger documentation...
swag init -g cmd/api/main.go -o docs

REM Build the application
echo 🎉 Building Go application...
go build -o bin/apiserver.exe cmd/api/main.go

echo 💡 Build completed for %ENVIRONMENT% environment
echo Binary: bin/apiserver.exe
echo Swagger docs: docs/