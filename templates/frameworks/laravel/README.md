# Laravel Template for Steele

This template creates a comprehensive Laravel development environment with Docker, AI context, and best practices built-in.

## What's Included

### <� Project Structure
- **Laravel 11.x** - Latest stable version
- **PHP 8.3** - Modern PHP with performance improvements
- **Composer** - Dependency management
- **Docker Compose** - Full development environment

### =3 Docker Services
- **App Container** - PHP-FPM with all necessary extensions
- **Nginx** - Web server with optimized Laravel configuration
- **MySQL 8.0** - Database with persistent storage
- **Redis** - Caching and session storage
- **Mailhog** - Email testing and debugging

### >� AI Integration
- **MCP Context** - Comprehensive Laravel context for Claude
- **Best Practices** - Pre-configured patterns and conventions
- **Code Templates** - Controller, model, and migration templates

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Composer (for local development)

### Installation with Steele
```bash
mkdir my-laravel-app && cd my-laravel-app
steele start laravel:11
```

### Manual Setup After Installation
1. **Copy environment configuration:**
   ```bash
   cp src/.env.example src/.env
   ```

2. **Update database configuration in `.env`:**
   ```env
   DB_CONNECTION=mysql
   DB_HOST=mysql
   DB_PORT=3306
   DB_DATABASE=laravel
   DB_USERNAME=laravel
   DB_PASSWORD=laravel
   ```

3. **Start Docker services:**
   ```bash
   docker-compose up -d
   ```

4. **Install dependencies and setup Laravel:**
   ```bash
   docker-compose exec app composer install
   docker-compose exec app php artisan key:generate
   docker-compose exec app php artisan migrate
   ```

## Development Workflow

### Docker Commands
```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f app

# Access app container
docker-compose exec app bash

# Run Artisan commands
docker-compose exec app php artisan migrate
docker-compose exec app php artisan make:controller UserController
```

### Available Services
- **Laravel App**: http://localhost:8000
- **Mailhog UI**: http://localhost:8025
- **MySQL**: localhost:3306
- **Redis**: localhost:6379

### Common Laravel Commands
```bash
# Create controllers
docker-compose exec app php artisan make:controller ProductController --resource

# Create models with migration
docker-compose exec app php artisan make:model Product -m

# Run migrations
docker-compose exec app php artisan migrate

# Run tests
docker-compose exec app php artisan test

# Clear caches
docker-compose exec app php artisan cache:clear
docker-compose exec app php artisan config:clear
docker-compose exec app php artisan view:clear
```

## AI Context Integration

The template includes a comprehensive AI context file at `ai/context.yaml` that provides:
- Laravel-specific patterns and conventions
- Common commands and workflows
- File templates for controllers and models
- Troubleshooting guides
- Best practices

This context helps AI assistants understand your Laravel project structure and provide better code suggestions.

## File Structure
```
project/
   src/                    # Laravel application
   ai/
      context.yaml       # AI context for Laravel
   infra/
      docker/            # Docker configuration
          Dockerfile
          docker-compose.yml
          nginx.conf
          local.ini
   steele.json           # Steele template configuration
   README.md
```

## Customization

### Adding New Services
Edit `docker-compose.yml` to add additional services like PostgreSQL, Elasticsearch, or custom containers.

### PHP Configuration
Modify `infra/docker/local.ini` to adjust PHP settings like memory limits and upload sizes.

### Nginx Configuration
Update `infra/docker/nginx.conf` for custom routing or SSL configuration.

## Troubleshooting

### Permission Issues
```bash
sudo chown -R $USER:$USER src/
docker-compose exec app chmod -R 775 storage bootstrap/cache
```

### Database Connection Issues
1. Ensure MySQL container is running: `docker-compose ps`
2. Check database credentials in `.env`
3. Verify database exists: `docker-compose exec mysql mysql -u laravel -p`

### Port Conflicts
If ports 8000, 3306, or 6379 are in use, modify the port mappings in `docker-compose.yml`.

## Production Considerations

This template is optimized for development. For production:
1. Use multi-stage Docker builds
2. Configure proper SSL certificates
3. Set up database backups
4. Implement proper logging and monitoring
5. Use environment-specific configuration