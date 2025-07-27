#!/bin/bash
# Module Generator Script for Unix/Linux/Mac
# Usage: ./generate-module.sh <module-name> [--with-permissions]

echo 🔧 Module Generator
echo =============================

if [ -z "$1" ]; then
    echo "💡 Usage: ./generate-module.sh <module-name> [--with-permissions]"
    echo "Example: ./generate-module.sh product"
    echo "Example: ./generate-module.sh product --with-permissions"
    exit 1
fi

echo "🔒 Generating module: $1"
if [ "$2" = "--with-permissions" ]; then
    echo "With permissions enabled"
    go run tools/module-generator/main.go "$1" --with-permissions
else
    go run tools/module-generator/main.go "$1"
fi