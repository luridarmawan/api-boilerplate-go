#!/bin/bash
# Module Generator Script for Unix/Linux/Mac
# Usage: ./generate-module.sh <module-name>

if [ -z "$1" ]; then
    echo "Usage: ./generate-module.sh <module-name>"
    echo "Example: ./generate-module.sh product"
    exit 1
fi

echo "Generating module: $1"
go run tools/module-generator/main.go "$1"