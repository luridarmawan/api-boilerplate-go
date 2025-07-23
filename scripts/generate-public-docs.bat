@echo off
REM Generate Public API Documentation
REM This script generates swagger documentation and removes internal endpoints

echo 📚 Generating Public API Documentation
echo =====================================

REM Step 1: Generate full swagger documentation
echo 🔨 Step 1: Generating Swagger documentation...
swag init -g cmd/api/main.go -o docs

if errorlevel 1 (
    echo ❌ Failed to generate Swagger documentation
    exit /b 1
)

echo ✅ Swagger documentation generated

REM Step 2: Create backup of original
echo 🔨 Step 2: Creating backup...
copy docs\swagger.json docs\swagger-full.json >nul
echo ✅ Backup created: docs\swagger-full.json

REM Step 3: Filter out internal endpoints
echo 🔨 Step 3: Filtering internal endpoints...
call scripts\filter-swagger.bat "Access,Permission"

if errorlevel 1 (
    echo ❌ Failed to filter endpoints
    exit /b 1
)

REM Step 4: Create public version
echo 🔨 Step 4: Creating public documentation...
copy docs\swagger.json docs\swagger-public.json >nul
echo ✅ Public documentation created: docs\swagger-public.json

echo 🎉 Public API documentation generated successfully!
echo 📁 Files created:
echo    📄 docs\swagger.json ^(filtered^)
echo    📄 docs\swagger-full.json ^(complete^)
echo    📄 docs\swagger-public.json ^(public^)
echo 🌐 Access documentation at:
echo    🔗 http://localhost:3000/docs ^(filtered^)
echo    🔗 Use swagger-full.json for complete documentation