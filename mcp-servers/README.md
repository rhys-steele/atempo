# Atempo MCP Servers

Plug-and-play MCP (Model Context Protocol) servers for Laravel and Django projects created with Atempo. These servers provide framework-specific tools that integrate seamlessly with Claude Code.

## Quick Start

### 1. Setup MCP Servers
```bash
cd mcp-servers
./setup.sh
```

### 2. Auto-Installation (Recommended)
When you run `atempo create laravel:11` or `atempo create django:5`, the appropriate MCP server is automatically copied to your project and configured.

### 3. Manual Installation
If you need to add an MCP server to an existing project:

```bash
# For Laravel projects
cp -r mcp-servers/laravel/* your-laravel-project/ai/mcp-server/

# For Django projects  
cp -r mcp-servers/django/* your-django-project/ai/mcp-server/
```

## Laravel MCP Server

### Available Tools

#### Core Laravel Commands
- **`laravel_artisan`** - Run any Artisan command
- **`laravel_make`** - Generate components (controllers, models, migrations, etc.)
- **`laravel_routes`** - Display and filter routes
- **`laravel_serve`** - Start development server

#### Database Operations
- **`laravel_db`** - Migrations, seeding, and database operations
- **`laravel_tinker`** - Execute PHP code in Laravel environment

#### Testing & Development
- **`laravel_test`** - Run PHPUnit tests
- **`laravel_config`** - Manage configuration

#### Composer Integration
- **`composer_install`** - Install dependencies
- **`composer_require`** - Add new packages

### Example Usage in Claude Code
```
Use laravel_make to create a User controller with resource methods
Use laravel_db to run fresh migrations  
Use laravel_artisan to clear all caches
Use composer_require to install laravel/sanctum
```

## Django MCP Server

### Available Tools

#### Core Django Commands
- **`django_manage`** - Run any Django management command
- **`django_startapp`** - Create new Django apps
- **`django_migrations`** - Handle migrations (make, migrate, show, etc.)
- **`django_runserver`** - Start development server

#### Development Tools
- **`django_shell`** - Execute Python code in Django environment
- **`django_test`** - Run Django tests
- **`django_check`** - Check project for issues
- **`django_collectstatic`** - Collect static files

#### User Management
- **`django_createsuperuser`** - Create admin users
- **`django_dbshell`** - Access database shell

#### Celery Integration
- **`celery_worker`** - Start Celery workers
- **`celery_beat`** - Start Celery beat scheduler

#### Python Package Management
- **`pip_install`** - Install Python packages
- **`pip_freeze`** - List installed packages

### Example Usage in Claude Code
```
Use django_startapp to create a new "blog" app
Use django_migrations to create and apply migrations for the blog app
Use django_shell to test model queries
Use celery_worker to start background task processing
```

## Configuration

### Claude Code Integration

Add this to your Claude Code MCP settings (usually in `~/.config/claude-code/mcp.json`):

#### For Laravel Projects
```json
{
  "mcpServers": {
    "atempo-laravel": {
      "command": "node",
      "args": ["ai/mcp-server/index.js"],
      "cwd": "/path/to/your/laravel/project"
    }
  }
}
```

#### For Django Projects  
```json
{
  "mcpServers": {
    "atempo-django": {
      "command": "node", 
      "args": ["ai/mcp-server/index.js"],
      "cwd": "/path/to/your/django/project"
    }
  }
}
```

### Auto-Configuration

When using Atempo CLI, the MCP configuration is automatically:
1. Copied to your project's `ai/mcp-config.json`
2. Set up with the correct paths
3. Ready to import into Claude Code

## Docker Integration

Both MCP servers are designed to work with the Docker environments created by Atempo:

- **Laravel**: Commands run inside the `app` container
- **Django**: Commands run inside the `web` container
- **Background services**: Full support for Redis, Celery, databases

## Requirements

- **Node.js 18+** (for running the MCP servers)
- **Docker & Docker Compose** (for the project environments)
- **Claude Code** (for MCP integration)

## Troubleshooting

### MCP Server Not Starting
1. Check Node.js version: `node --version` (requires 18+)
2. Install dependencies: `cd ai/mcp-server && npm install`
3. Check Docker containers are running: `docker-compose ps`

### Commands Failing
1. Ensure you're in the project root directory
2. Verify Docker containers are running
3. Check container logs: `docker-compose logs app` (Laravel) or `docker-compose logs web` (Django)

### Permission Issues
```bash
# Fix container permissions
docker-compose exec app chown -R www:www /var/www  # Laravel
docker-compose exec web chown -R django:django /app  # Django
```

## Development

### Adding New Tools

To add new tools to an MCP server:

1. Add the tool definition to the `tools` array in `ListToolsRequestSchema`
2. Add a case in the `CallToolRequestSchema` handler
3. Implement the tool method
4. Test with Claude Code

### Custom Commands

Both servers support running arbitrary commands through their base methods:
- Laravel: `laravel_artisan` for any Artisan command
- Django: `django_manage` for any management command

## Examples

### Laravel Workflow
```
1. Use laravel_make to create a Product model with migration
2. Use laravel_db to run the migration  
3. Use laravel_make to create a ProductController with resource methods
4. Use laravel_routes to verify the new routes
5. Use laravel_test to run tests
```

### Django Workflow
```
1. Use django_startapp to create a "products" app
2. Use django_migrations to create models and migrate
3. Use django_createsuperuser to create an admin user
4. Use django_runserver to start the development server
5. Use django_test to run the test suite
```

These MCP servers make framework development with Claude Code incredibly powerful and streamlined!