#!/bin/bash
# MCP Setup Script - Copies example config to .cursor directory

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CURSOR_DIR="$PROJECT_ROOT/.cursor"
EXAMPLE_CONFIG="$PROJECT_ROOT/mcp-config.example.json"
TARGET_CONFIG="$CURSOR_DIR/mcp.json"

echo "🚀 Setting up MCP configuration for LogiTrackPro..."
echo ""

# Create .cursor directory if it doesn't exist
if [ ! -d "$CURSOR_DIR" ]; then
    echo "📁 Creating .cursor directory..."
    mkdir -p "$CURSOR_DIR"
fi

# Check if config already exists
if [ -f "$TARGET_CONFIG" ]; then
    echo "⚠️  MCP config already exists at $TARGET_CONFIG"
    read -p "Overwrite? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Cancelled. Keeping existing config."
        exit 0
    fi
fi

# Copy example config
if [ -f "$EXAMPLE_CONFIG" ]; then
    echo "📋 Copying MCP configuration..."
    cp "$EXAMPLE_CONFIG" "$TARGET_CONFIG"
    
    # Update filesystem path to use actual project root
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|/Users/dankovartem/repos/LogiTrackPro|$PROJECT_ROOT|g" "$TARGET_CONFIG"
    else
        # Linux
        sed -i "s|/Users/dankovartem/repos/LogiTrackPro|$PROJECT_ROOT|g" "$TARGET_CONFIG"
    fi
    
    echo "✅ MCP configuration created at $TARGET_CONFIG"
else
    echo "❌ Example config not found: $EXAMPLE_CONFIG"
    exit 1
fi

echo ""
echo "📋 Configured MCP Servers:"
echo "  ✅ PostgreSQL MCP - Database operations"
echo "  ✅ Filesystem MCP - File operations"
echo "  ✅ Git MCP - Version control"
echo "  ✅ Docker MCP - Container management"
echo "  ⚠️  Brave Search MCP - Requires BRAVE_API_KEY env var"
echo ""
echo "⚠️  Next steps:"
echo "  1. Review and update $TARGET_CONFIG if needed"
echo "  2. Set BRAVE_API_KEY environment variable (optional)"
echo "  3. Restart Cursor to load MCP servers"
echo ""
echo "✨ Setup complete!"

