# MCP (Model Context Protocol) Setup Guide

This document describes the MCP servers configured for LogiTrackPro and how to use them.

## Overview

MCP servers provide AI assistants with structured access to project resources, enabling more intelligent and context-aware assistance.

## Configured MCP Servers

### 1. PostgreSQL MCP
**Purpose**: Database operations, queries, and migrations

**Features**:
- Execute SQL queries
- View table schemas
- Run migrations
- Database introspection

**Configuration**: Located in `.cursor/mcp.json`
- Connection string: `postgres://postgres:postgres@localhost:5432/logitrackpro`

**Usage Examples**:
- "Show me all warehouses in the database"
- "What's the schema for the plans table?"
- "Run a query to find customers with low inventory"

### 2. Filesystem MCP
**Purpose**: File operations and project navigation

**Features**:
- Read/write files
- List directories
- Search for files
- File manipulation

**Configuration**: Project root directory

**Usage Examples**:
- "Read the main.go file"
- "List all files in the handlers directory"
- "Create a new migration file"

### 3. Git MCP
**Purpose**: Version control operations

**Features**:
- View commit history
- Check git status
- Create branches
- View diffs

**Configuration**: Project root directory

**Usage Examples**:
- "Show me recent commits"
- "What files have changed?"
- "Create a new branch for feature X"

### 4. Docker MCP
**Purpose**: Container management

**Features**:
- View running containers
- Check container logs
- Manage Docker Compose services
- Container health checks

**Usage Examples**:
- "Show me running containers"
- "What are the logs for the backend service?"
- "Restart the optimizer service"

### 5. Brave Search MCP (Optional)
**Purpose**: Web search for documentation and solutions

**Features**:
- Search the web
- Find documentation
- Look up error solutions

**Configuration**: Requires `BRAVE_API_KEY` environment variable

## Setup Instructions

### 1. Initial Setup

Run the setup script:
```bash
chmod +x scripts/mcp-setup.sh
./scripts/mcp-setup.sh
```

### 2. Environment Variables

Create a `.env` file in the project root (if needed):
```bash
# For Brave Search MCP (optional)
BRAVE_API_KEY=your_api_key_here

# Database connection (if different from default)
POSTGRES_CONNECTION_STRING=postgres://user:pass@host:port/db
```

### 3. Verify Installation

1. Restart Cursor
2. Check MCP server status in Cursor settings
3. Test with a simple query: "List all files in the backend directory"

## Helper Scripts

### Database Helper (`scripts/db-helper.sh`)

Common database operations:
```bash
./scripts/db-helper.sh reset      # Reset database
./scripts/db-helper.sh backup      # Create backup
./scripts/db-helper.sh status      # Check status
./scripts/db-helper.sh tables      # List tables
./scripts/db-helper.sh stats       # Show statistics
```

### Docker Helper (`scripts/docker-helper.sh`)

Docker Compose operations:
```bash
./scripts/docker-helper.sh up          # Start services
./scripts/docker-helper.sh down        # Stop services
./scripts/docker-helper.sh logs        # View logs
./scripts/docker-helper.sh shell backend  # Open shell
./scripts/docker-helper.sh db-shell    # PostgreSQL shell
```

## MCP Capabilities by Use Case

### Development
- **Code Navigation**: Filesystem MCP
- **Database Queries**: PostgreSQL MCP
- **Version Control**: Git MCP
- **Container Management**: Docker MCP

### Debugging
- **View Logs**: Docker MCP
- **Database Inspection**: PostgreSQL MCP
- **File Search**: Filesystem MCP

### Documentation
- **Web Search**: Brave Search MCP
- **Code Reading**: Filesystem MCP
- **Git History**: Git MCP

## Troubleshooting

### MCP Server Not Connecting

1. **Check Node.js**: Ensure Node.js 18+ is installed
   ```bash
   node --version
   ```

2. **Check Configuration**: Verify `.cursor/mcp.json` exists and is valid JSON

3. **Check Services**: Ensure required services are running:
   - PostgreSQL (for Postgres MCP)
   - Docker (for Docker MCP)

4. **Restart Cursor**: MCP servers load on Cursor startup

### PostgreSQL MCP Issues

- Verify database is running: `docker-compose ps`
- Check connection string in `.cursor/mcp.json`
- Test connection: `psql postgres://postgres:postgres@localhost:5432/logitrackpro`

### Docker MCP Issues

- Ensure Docker is running: `docker ps`
- Check Docker Compose file exists
- Verify container names match configuration

## Best Practices

1. **Use MCPs for Repetitive Tasks**: Let AI handle database queries, file operations
2. **Verify Critical Operations**: Always review changes before committing
3. **Use Helper Scripts**: For complex operations, use provided scripts
4. **Keep Configuration Updated**: Update MCP config when project structure changes

## Security Notes

- Never commit `.cursor/mcp.json` with production credentials
- Use environment variables for sensitive data
- Review MCP server permissions
- Keep MCP servers updated

## Additional Resources

- [MCP Documentation](https://modelcontextprotocol.io)
- [Cursor MCP Guide](https://docs.cursor.com)
- Project-specific: See `README.md` for project setup

## Support

For issues with MCP setup:
1. Check this guide
2. Review Cursor logs
3. Verify service status
4. Check MCP server documentation

