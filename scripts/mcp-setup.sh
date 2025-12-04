#!/bin/bash
# MCP Setup Script for LogiTrackPro
# This script helps set up MCP servers for the project

set -e

echo "🚀 Setting up MCP servers for LogiTrackPro..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker is not installed. Some MCPs may not work."
fi

# Create .cursor directory if it doesn't exist
mkdir -p .cursor

echo "✅ MCP configuration created at .cursor/mcp.json"
echo ""
echo "📋 Available MCP Servers:"
echo "  1. PostgreSQL MCP - Database operations"
echo "  2. Filesystem MCP - File operations"
echo "  3. Git MCP - Version control"
echo "  4. Docker MCP - Container management"
echo "  5. Brave Search MCP - Web search (optional)"
echo ""
echo "⚠️  Note: You may need to:"
echo "  - Set BRAVE_API_KEY environment variable for Brave Search"
echo "  - Ensure PostgreSQL is running for Postgres MCP"
echo "  - Restart Cursor after configuration changes"
echo ""
echo "✨ Setup complete!"

