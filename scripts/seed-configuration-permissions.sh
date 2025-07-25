#!/bin/bash

# Safe Configuration Permissions Seeder
# This script adds configuration permissions ONLY to Admin group

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔒 Configuration Permissions Seeder (Admin Only)${NC}"
echo -e "${BLUE}================================================${NC}"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: Please run this script from the project root directory${NC}"
    exit 1
fi

# Check if database is accessible
echo -e "${YELLOW}🔍 Checking database connection...${NC}"
if ! go run scripts/add-admin-only-configuration-permissions.go > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Cannot connect to database. Please check your .env configuration${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Database connection successful${NC}"

# Ask for confirmation
echo -e "${YELLOW}⚠️  This will add configuration permissions ONLY to Admin group${NC}"
echo -e "${YELLOW}   Other groups (Editor, Viewer, General client) will NOT get access${NC}"
echo ""
read -p "Do you want to continue? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}❌ Operation cancelled${NC}"
    exit 0
fi

# Run the admin-only configuration permissions script
echo -e "${BLUE}🚀 Adding configuration permissions for Admin group only...${NC}"
go run scripts/add-admin-only-configuration-permissions.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Configuration permissions successfully added!${NC}"
    
    # Optional: Remove any existing configuration permissions from non-admin groups
    echo ""
    read -p "Do you want to remove configuration permissions from non-admin groups (if any)? (y/N): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}🧹 Cleaning up configuration permissions from non-admin groups...${NC}"
        go run scripts/remove-configuration-permissions-from-non-admin.go
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ Cleanup completed successfully!${NC}"
        else
            echo -e "${RED}❌ Cleanup failed${NC}"
        fi
    fi
    
    echo ""
    echo -e "${GREEN}🎉 Configuration module is now secured for Admin-only access!${NC}"
    echo -e "${BLUE}📋 Summary:${NC}"
    echo -e "   ✅ Configuration permissions created"
    echo -e "   ✅ Permissions assigned to Admin group only"
    echo -e "   ✅ Other groups have no configuration access"
    echo -e "   ✅ Security: Admin-only access enforced"
    
else
    echo -e "${RED}❌ Failed to add configuration permissions${NC}"
    exit 1
fi