#!/bin/bash

# Swagger Filter Script
# Usage: ./scripts/filter-swagger.sh [tags-to-remove] [input-file] [output-file]

# Default values
TAGS_TO_REMOVE=${1:-"Access,Example,Permission"}
INPUT_FILE=${2:-"docs/swagger.json"}
OUTPUT_FILE=${3:-""}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Swagger Endpoint Filter${NC}"
echo -e "${BLUE}================================${NC}"

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo -e "${RED}❌ Error: Input file '$INPUT_FILE' not found${NC}"
    exit 1
fi

# Build the tool if not exists
TOOL_DIR="tools/swagger-filter"
if [ ! -f "$TOOL_DIR/swagger-filter" ]; then
    echo -e "${YELLOW}🔨 Building swagger filter tool...${NC}"
    cd "$TOOL_DIR"
    go build -o swagger-filter main.go
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Failed to build swagger filter tool${NC}"
        exit 1
    fi
    cd - > /dev/null
    echo -e "${GREEN}✅ Tool built successfully${NC}"
fi

# Run the filter tool
echo -e "${YELLOW}🔍 Filtering swagger endpoints...${NC}"
echo -e "📋 Tags to remove: ${TAGS_TO_REMOVE}"
echo -e "📁 Input file: ${INPUT_FILE}"

if [ -n "$OUTPUT_FILE" ]; then
    echo -e "📁 Output file: ${OUTPUT_FILE}"
    "$TOOL_DIR/swagger-filter" -input="$INPUT_FILE" -output="$OUTPUT_FILE" -remove-tags="$TAGS_TO_REMOVE" -verbose
else
    echo -e "📁 Output file: ${INPUT_FILE} (overwrite)"
    "$TOOL_DIR/swagger-filter" -input="$INPUT_FILE" -remove-tags="$TAGS_TO_REMOVE" -verbose
fi

if [ $? -eq 0 ]; then
    echo -e "${GREEN}🎉 Swagger filtering completed successfully!${NC}"
else
    echo -e "${RED}❌ Swagger filtering failed${NC}"
    exit 1
fi