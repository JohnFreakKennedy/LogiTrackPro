#!/bin/bash
# Test MCP Server Connectivity
# This script helps verify that MCP servers are properly configured

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🧪 Testing MCP Server Configuration${NC}"
echo ""

# Check Node.js
echo -n "Checking Node.js... "
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo -e "${GREEN}✅ $NODE_VERSION${NC}"
else
    echo -e "${RED}❌ Not installed${NC}"
    exit 1
fi

# Check Docker
echo -n "Checking Docker... "
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | cut -d' ' -f3 | tr -d ',')
    echo -e "${GREEN}✅ $DOCKER_VERSION${NC}"
else
    echo -e "${YELLOW}⚠️  Not installed (Docker MCP won't work)${NC}"
fi

# Check PostgreSQL connection
echo -n "Checking PostgreSQL connection... "
if PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d logitrackpro -c "SELECT 1;" &> /dev/null; then
    echo -e "${GREEN}✅ Connected${NC}"
else
    echo -e "${YELLOW}⚠️  Cannot connect (PostgreSQL MCP won't work)${NC}"
    echo "   Make sure PostgreSQL is running: docker-compose up -d postgres"
fi

# Check MCP config file
echo -n "Checking MCP config file... "
CONFIG_FILE=".cursor/mcp.json"
if [ -f "$CONFIG_FILE" ]; then
    echo -e "${GREEN}✅ Found${NC}"
    
    # Validate JSON
    if command -v jq &> /dev/null; then
        if jq empty "$CONFIG_FILE" 2>/dev/null; then
            echo -e "   ${GREEN}✅ Valid JSON${NC}"
        else
            echo -e "   ${RED}❌ Invalid JSON${NC}"
        fi
    else
        echo -e "   ${YELLOW}⚠️  Install 'jq' to validate JSON${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Not found${NC}"
    echo "   Run: ./scripts/setup-mcp.sh"
fi

# Check environment variables
echo -n "Checking BRAVE_API_KEY... "
if [ -n "$BRAVE_API_KEY" ]; then
    echo -e "${GREEN}✅ Set${NC}"
else
    echo -e "${YELLOW}⚠️  Not set (Brave Search MCP won't work)${NC}"
fi

echo ""
echo -e "${BLUE}📋 Summary${NC}"
echo "  MCP servers will be available after restarting Cursor"
echo "  Check Cursor settings > MCP to verify server status"
echo ""

