# Django Template for Atempo

This template creates a comprehensive Django development environment with Docker, AI context, and modern Python best practices built-in.

## What's Included

### 🏗️ Project Structure
- **Django 5.x** - Latest stable version with async support
- **Python 3.12** - Modern Python with performance improvements
- **PostgreSQL 16** - Robust relational database
- **Redis** - Fast in-memory caching and session storage

### 🐳 Docker Services
- **Web Container** - Django application with Python 3.12
- **PostgreSQL** - Primary database with persistent storage
- **Redis** - Caching, sessions, and Celery message broker
- **Celery Worker** - Background task processing
- **Celery Beat** - Scheduled task management
- **Mailhog** - Email testing and debugging

### 🧠 AI Integration
- **MCP Context** - Comprehensive Django context for Claude
- **Best Practices** - Pre-configured patterns and conventions
- **Code Templates** - Model, view, and form templates

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Python 3.12+ (for local development)

### Installation with Atempo
```bash
mkdir my-django-app && cd my-django-app
atempo create django:5
```

### Manual Setup After Installation
1. **Copy requirements and install dependencies:**
   ```bash
   cp infra/docker/requirements.txt src/
   ```

2. **Start Docker services:**
   ```bash
   docker-compose up -d
   ```

3. **Run Django setup commands:**
   ```bash
   docker-compose exec web python manage.py migrate
   docker-compose exec web python manage.py createsuperuser
   docker-compose exec web python manage.py collectstatic --noinput
   ```

## Development Workflow

### Docker Commands
```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f web

# Access web container
docker-compose exec web bash

# Run Django management commands
docker-compose exec web python manage.py migrate
docker-compose exec web python manage.py startapp myapp
docker-compose exec web python manage.py shell
```

### Available Services
- **Django App**: http://localhost:8000
- **Django Admin**: http://localhost:8000/admin
- **Mailhog UI**: http://localhost:8025
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

### Common Django Commands
```bash
# Create new app
docker-compose exec web python manage.py startapp myapp

# Make and run migrations
docker-compose exec web python manage.py makemigrations
docker-compose exec web python manage.py migrate

# Create superuser
docker-compose exec web python manage.py createsuperuser

# Run tests
docker-compose exec web python manage.py test

# Collect static files
docker-compose exec web python manage.py collectstatic

# Django shell
docker-compose exec web python manage.py shell

# Check for deployment issues
docker-compose exec web python manage.py check --deploy
```

### Background Tasks with Celery
```bash
# Monitor Celery worker
docker-compose logs -f worker

# Monitor Celery beat scheduler
docker-compose logs -f beat

# Access Celery shell
docker-compose exec worker celery -A config shell
```

## AI Context Integration

The template includes a comprehensive AI context file at `ai/context.yaml` that provides:
- Django-specific patterns and conventions
- Common commands and workflows
- File templates for models, views, and forms
- Best practices and security guidelines
- Troubleshooting guides

This context helps AI assistants understand your Django project structure and provide better code suggestions.

## Project Structure
```
project/
├── src/                          # Django application
│   ├── config/                   # Django settings
│   ├── apps/                     # Django apps
│   ├── templates/                # HTML templates
│   ├── static/                   # Static files
│   ├── media/                    # User uploads
│   ├── manage.py
│   └── requirements.txt
├── ai/
│   └── context.yaml             # AI context for Django
├── infra/
│   └── docker/                  # Docker configuration
│       ├── Dockerfile
│       ├── docker-compose.yml
│       └── requirements.txt
├── atempo.json                  # Atempo template configuration
└── README.md
```

## Environment Variables

Key environment variables you can configure:

```env
# Django Settings
DEBUG=1
SECRET_KEY=your-secret-key-here
ALLOWED_HOSTS=localhost,127.0.0.1

# Database
DATABASE_URL=postgresql://django:django@postgres:5432/django

# Cache
REDIS_URL=redis://redis:6379/0

# Email (using Mailhog for development)
EMAIL_HOST=mailhog
EMAIL_PORT=1025
```

## Customization

### Adding New Django Apps
```bash
docker-compose exec web python manage.py startapp myapp
```

Then add the app to `INSTALLED_APPS` in your Django settings.

### Database Configuration
The template uses PostgreSQL by default. To switch databases, update the `DATABASE_URL` environment variable and the docker-compose.yml file.

### Adding New Services
Edit `docker-compose.yml` to add additional services like Elasticsearch, MongoDB, or custom containers.

## Production Considerations

This template is optimized for development. For production:

1. **Security**: Generate a strong `SECRET_KEY` and set `DEBUG=False`
2. **Database**: Use managed database services with backups
3. **Static Files**: Configure proper static file serving (WhiteNoise is included)
4. **Logging**: Implement comprehensive logging and monitoring
5. **SSL/HTTPS**: Configure SSL certificates and secure headers
6. **Environment**: Use proper environment variable management
7. **Scaling**: Consider using Gunicorn/uWSGI and load balancers

## Testing

```bash
# Run all tests
docker-compose exec web python manage.py test

# Run specific app tests
docker-compose exec web python manage.py test myapp

# Run with coverage
docker-compose exec web coverage run --source='.' manage.py test
docker-compose exec web coverage report
```

## Troubleshooting

### Database Issues
```bash
# Check database connection
docker-compose exec web python manage.py dbshell

# Reset database
docker-compose down
docker volume rm django_postgres_data
docker-compose up -d
```

### Static Files Issues
```bash
# Collect static files
docker-compose exec web python manage.py collectstatic --clear --noinput
```

### Permission Issues
```bash
# Fix container permissions
docker-compose exec web chown -R django:django /app
```

### Port Conflicts
If ports 8000, 5432, or 6379 are in use, modify the port mappings in `docker-compose.yml`.