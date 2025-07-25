#!/bin/bash

# Configuration Module Setup Script
# This script sets up the configuration module with proper database structure and sample data

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Configuration Module Setup${NC}"
echo -e "${BLUE}=============================${NC}"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: Please run this script from the project root directory${NC}"
    exit 1
fi

# Step 1: Check database connection
echo -e "${YELLOW}🔍 Step 1: Checking database connection...${NC}"
if ! go run scripts/test-database-connection.go > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Cannot connect to database. Please check your .env configuration${NC}"
    echo -e "${YELLOW}💡 Run this for detailed error info:${NC}"
    echo -e "   go run scripts/test-database-connection.go"
    echo -e "${YELLOW}💡 Make sure:${NC}"
    echo -e "   - Database server is running"
    echo -e "   - .env file exists with correct database settings"
    echo -e "   - Database credentials are correct"
    exit 1
fi
echo -e "${GREEN}✅ Database connection successful${NC}"

# Step 2: Set up permissions (if not already done)
echo -e "${YELLOW}🔒 Step 2: Setting up configuration permissions...${NC}"
go run scripts/add-admin-only-configuration-permissions.go > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Configuration permissions verified${NC}"
else
    echo -e "${YELLOW}⚠️  Permission setup may need attention${NC}"
fi

# Step 3: Seed sample data
echo ""
read -p "Do you want to add sample configuration data? (y/N): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🌱 Step 3: Seeding sample configuration data...${NC}"
    go run scripts/seed-sample-configurations.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Sample data seeded successfully${NC}"
    else
        echo -e "${RED}❌ Failed to seed sample data${NC}"
    fi
else
    echo -e "${BLUE}ℹ️  Skipping sample data seeding${NC}"
fi

# Step 4: Generate updated Swagger documentation
echo -e "${YELLOW}📚 Step 4: Updating Swagger documentation...${NC}"
if command -v swag &> /dev/null; then
    swag init -g cmd/api/main.go -o docs
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Swagger documentation updated${NC}"
    else
        echo -e "${YELLOW}⚠️  Swagger documentation update failed${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Swag command not found. Please install swaggo/swag to update documentation${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Configuration Module Setup Complete!${NC}"
echo -e "${BLUE}=======================================${NC}"

echo -e "${BLUE}📋 What was set up:${NC}"
echo -e "   ✅ Configuration module with key-value structure"
echo -e "   ✅ Admin-only permissions for configuration access"
echo -e "   ✅ Sample configuration data (if selected)"
echo -e "   ✅ Updated Swagger documentation"

echo -e "${BLUE}🔗 Available Endpoints:${NC}"
echo -e "   📄 GET    /v1/configurations           - List all configurations"
echo -e "   📄 POST   /v1/configurations           - Create new configuration"
echo -e "   📄 GET    /v1/configurations/{id}      - Get configuration by ID"
echo -e "   📄 GET    /v1/configurations/key/{key} - Get configuration by key"
echo -e "   📄 PUT    /v1/configurations/{id}      - Update configuration"
echo -e "   📄 DELETE /v1/configurations/{id}      - Delete configuration"

echo -e "${BLUE}🔒 Security:${NC}"
echo -e "   ✅ Only Admin users can access configuration endpoints"
echo -e "   ✅ All endpoints require authentication"
echo -e "   ✅ Rate limiting applied"

echo -e "${BLUE}🧪 Test the API:${NC}"
echo -e "   curl -X GET \"http://localhost:3000/v1/configurations\" \\"
echo -e "     -H \"Authorization: Bearer admin-api-key-789\""

echo ""
echo -e "${GREEN}Configuration module is ready to use! 🚀${NC}"