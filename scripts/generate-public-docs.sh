#!/bin/bash

# Generate Public API Documentation
# This script generates swagger documentation and removes internal endpoints

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📚 Generating Public API Documentation${NC}"
echo -e "${BLUE}=====================================${NC}"

# Step 1: Generate full swagger documentation
echo -e "${YELLOW}🔨 Step 1: Generating Swagger documentation...${NC}"
swag init -g cmd/api/main.go -o docs

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to generate Swagger documentation${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Swagger documentation generated${NC}"

# Step 2: Create backup of original
echo -e "${YELLOW}🔨 Step 2: Creating backup...${NC}"
cp docs/swagger.json docs/swagger-full.json
echo -e "${GREEN}✅ Backup created: docs/swagger-full.json${NC}"

# Step 3: Filter out internal endpoints
echo -e "${YELLOW}🔨 Step 3: Filtering internal endpoints...${NC}"
./scripts/filter-swagger.sh "Access,Permission"

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to filter endpoints${NC}"
    exit 1
fi

# Step 4: Create public version
echo -e "${YELLOW}🔨 Step 4: Creating public documentation...${NC}"
cp docs/swagger.json docs/swagger-public.json
echo -e "${GREEN}✅ Public documentation created: docs/swagger-public.json${NC}"

echo -e "${GREEN}🎉 Public API documentation generated successfully!${NC}"
echo -e "${BLUE}📁 Files created:${NC}"
echo -e "   📄 docs/swagger.json (filtered)"
echo -e "   📄 docs/swagger-full.json (complete)"
echo -e "   📄 docs/swagger-public.json (public)"
echo -e "${BLUE}🌐 Access documentation at:${NC}"
echo -e "   🔗 http://localhost:3000/docs (filtered)"
echo -e "   🔗 Use swagger-full.json for complete documentation"