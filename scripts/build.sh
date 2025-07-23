#!/bin/bash

# Build script for different environments
# Usage: ./scripts/build.sh [dev|staging|prod]

ENVIRONMENT=${1:-dev}

echo "Building for environment: $ENVIRONMENT"

# Copy environment-specific config
if [ -f ".env.$ENVIRONMENT" ]; then
    echo "Using .env.$ENVIRONMENT configuration"
    cp ".env.$ENVIRONMENT" ".env"
else
    echo "Warning: .env.$ENVIRONMENT not found, using default .env"
fi

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g cmd/api/main.go -o docs

# Build the application
echo "Building Go application..."
go build -o bin/apiserver cmd/api/main.go

echo "Build completed for $ENVIRONMENT environment"
echo "Binary: bin/apiserver"
echo "Swagger docs: docs/"
echo "Environment config: .env (copied from .env.$ENVIRONMENT)"