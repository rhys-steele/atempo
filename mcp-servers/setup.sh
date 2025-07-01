#!/bin/bash

# Setup script for Atempo MCP servers
# This script installs dependencies for both Laravel and Django MCP servers

set -e

echo "🔧 Setting up Atempo MCP servers..."

# Function to install npm dependencies
install_server() {
    local server_name=$1
    local server_dir=$2
    
    echo "📦 Installing dependencies for $server_name MCP server..."
    cd "$server_dir"
    
    if command -v npm &> /dev/null; then
        npm install
        echo "✅ $server_name MCP server dependencies installed"
    else
        echo "❌ npm not found. Please install Node.js and npm first."
        exit 1
    fi
    
    cd - > /dev/null
}

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js not found. Please install Node.js (version 18 or higher) first."
    echo "   Visit: https://nodejs.org/"
    exit 1
fi

# Check Node.js version
NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
    echo "❌ Node.js version 18 or higher required. Current version: $(node --version)"
    exit 1
fi

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Install Laravel MCP server
if [ -d "$SCRIPT_DIR/laravel" ]; then
    install_server "Laravel" "$SCRIPT_DIR/laravel"
else
    echo "⚠️  Laravel MCP server directory not found"
fi

# Install Django MCP server
if [ -d "$SCRIPT_DIR/django" ]; then
    install_server "Django" "$SCRIPT_DIR/django"
else
    echo "⚠️  Django MCP server directory not found"
fi

echo ""
echo "🎉 MCP server setup complete!"
echo ""
echo "Next steps:"
echo "1. Copy the appropriate MCP server to your project's ai/mcp-server/ directory"
echo "2. Add the MCP configuration to your Claude Code settings"
echo "3. Use the framework-specific tools in Claude Code!"
echo ""
echo "For Laravel projects:"
echo "  cp -r $SCRIPT_DIR/laravel/* your-project/ai/mcp-server/"
echo ""
echo "For Django projects:"
echo "  cp -r $SCRIPT_DIR/django/* your-project/ai/mcp-server/"