#!/bin/bash

# Script to generate Swagger documentation with conditional modules
# Usage: ./scripts/generate-swagger.sh [options]
# Options:
#   --show-audit    Include Audit module in documentation
#   --show-access   Include Access module in documentation
#   --show-all      Include both Audit and Access modules
#   --help          Show this help message

SHOW_AUDIT=false
SHOW_ACCESS=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --show-audit)
            SHOW_AUDIT=true
            shift
            ;;
        --show-access)
            SHOW_ACCESS=true
            shift
            ;;
        --show-all)
            SHOW_AUDIT=true
            SHOW_ACCESS=true
            shift
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --show-audit    Include Audit module in documentation"
            echo "  --show-access   Include Access module in documentation"
            echo "  --show-all      Include both Audit and Access modules"
            echo "  --help          Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo "Generating Swagger documentation..."
echo "Show Audit: $SHOW_AUDIT"
echo "Show Access: $SHOW_ACCESS"

# Create backup directory
mkdir -p .swagger-backup

# Backup original files
cp internal/modules/audit/audit_handler.go .swagger-backup/
cp internal/modules/access/access_handler.go .swagger-backup/

# Process files based on options
if [ "$SHOW_AUDIT" = false ]; then
    echo "Hiding Audit module from documentation..."
    sed -i '/SWAGGER_AUDIT_START/,/SWAGGER_AUDIT_END/d' internal/modules/audit/audit_handler.go
fi

if [ "$SHOW_ACCESS" = false ]; then
    echo "Hiding Access module from documentation..."
    sed -i '/SWAGGER_ACCESS_START/,/SWAGGER_ACCESS_END/d' internal/modules/access/access_handler.go
fi

# Generate Swagger docs
echo "Generating documentation..."
swag init -g cmd/api/main.go -o docs

# Restore original files
echo "Restoring original files..."
cp .swagger-backup/audit_handler.go internal/modules/audit/
cp .swagger-backup/access_handler.go internal/modules/access/

# Cleanup
rm -rf .swagger-backup

echo "Swagger documentation generated successfully!"
echo "Access documentation at: http://localhost:3000/docs"