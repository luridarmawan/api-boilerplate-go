@echo off
REM Swagger Filter Script
REM Usage: scripts\filter-swagger.bat [tags-to-remove] [input-file] [output-file]

REM Default values
set TAGS_TO_REMOVE=%1
if "%TAGS_TO_REMOVE%"=="" set TAGS_TO_REMOVE=Access,Example,Permission

set INPUT_FILE=%2
if "%INPUT_FILE%"=="" set INPUT_FILE=docs\swagger.json

set OUTPUT_FILE=%3

echo 🔧 Swagger Endpoint Filter
echo ================================

REM Check if input file exists
if not exist "%INPUT_FILE%" (
    echo ❌ Error: Input file '%INPUT_FILE%' not found
    exit /b 1
)

REM Build the tool if not exists
set TOOL_DIR=tools\swagger-filter
if not exist "%TOOL_DIR%\swagger-filter.exe" (
    echo 🔨 Building swagger filter tool...
    cd "%TOOL_DIR%"
    go build -o swagger-filter.exe main.go
    if errorlevel 1 (
        echo ❌ Failed to build swagger filter tool
        exit /b 1
    )
    cd ..\..
    echo ✅ Tool built successfully
)

REM Run the filter tool
echo 🔍 Filtering swagger endpoints...
echo 📋 Tags to remove: %TAGS_TO_REMOVE%
echo 📁 Input file: %INPUT_FILE%

if "%OUTPUT_FILE%"=="" (
    echo 📁 Output file: %INPUT_FILE% ^(overwrite^)
    "%TOOL_DIR%\swagger-filter.exe" -input="%INPUT_FILE%" -remove-tags="%TAGS_TO_REMOVE%" -verbose
) else (
    echo 📁 Output file: %OUTPUT_FILE%
    "%TOOL_DIR%\swagger-filter.exe" -input="%INPUT_FILE%" -output="%OUTPUT_FILE%" -remove-tags="%TAGS_TO_REMOVE%" -verbose
)

if errorlevel 1 (
    echo ❌ Swagger filtering failed
    exit /b 1
) else (
    echo 🎉 Swagger filtering completed successfully!
)