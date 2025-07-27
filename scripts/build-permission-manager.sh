#!/bin/bash

# Build script for permission-manager CLI tool
# This script builds the permission manager for multiple platforms

set -e

echo "🔨 Building Permission Manager CLI Tool..."

# Get the current directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_ROOT/bin"

# Create build directory if it doesn't exist
mkdir -p "$BUILD_DIR"

# Build information
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X 'main.Version=$VERSION' -X 'main.GitCommit=$COMMIT' -X 'main.BuildDate=$BUILD_DATE'"

echo "📦 Version: $VERSION"
echo "📦 Commit: $COMMIT"
echo "📦 Build Date: $BUILD_DATE"

# Change to project root
cd "$PROJECT_ROOT"

# Build for current platform
echo "🏗️  Building for current platform..."
go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/permission-manager" ./cmd/permission-manager

# Build for multiple platforms
echo "🏗️  Building for multiple platforms..."

# Windows
echo "  - Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/permission-manager-windows-amd64.exe" ./cmd/permission-manager

# Linux
echo "  - Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/permission-manager-linux-amd64" ./cmd/permission-manager

# macOS
echo "  - Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/permission-manager-darwin-amd64" ./cmd/permission-manager

# macOS ARM64 (Apple Silicon)
echo "  - Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/permission-manager-darwin-arm64" ./cmd/permission-manager

echo "✅ Build completed successfully!"
echo ""
echo "📁 Binaries created in: $BUILD_DIR"
echo "   - permission-manager (current platform)"
echo "   - permission-manager-windows-amd64.exe"
echo "   - permission-manager-linux-amd64"
echo "   - permission-manager-darwin-amd64"
echo "   - permission-manager-darwin-arm64"
echo ""
echo "🚀 Usage example:"
echo "   $BUILD_DIR/permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create"