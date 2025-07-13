# {{project}} - Development Workflows

## Quick Start Guide

### Initial Setup
```bash
# Clone and navigate to project
cd {{project}}

# Start Docker environment
docker-compose up -d

# Install PHP dependencies
docker-compose exec app composer install

# Generate application key
docker-compose exec app php artisan key:generate

# Run database migrations
docker-compose exec app php artisan migrate --seed

# Install and build frontend assets
docker-compose exec app npm install
docker-compose exec app npm run dev

# Access application
open http://localhost:8000
```

### Development Environment Status
```bash
# Check all services
docker-compose ps

# View logs
docker-compose logs -f app
docker-compose logs -f webserver
docker-compose logs -f mysql

# Check application health
docker-compose exec app php artisan route:list
docker-compose exec app php artisan config:cache
```

## Core Development Workflows

### 1. Feature Development Workflow

#### Step 1: Create Feature Branch
```bash
git checkout -b feature/new-feature-name
```

#### Step 2: Implement Feature
```bash
# Create migration if needed
docker-compose exec app php artisan make:migration create_feature_table

# Create model with factory and seeder
docker-compose exec app php artisan make:model FeatureModel -mfs

# Create controller
docker-compose exec app php artisan make:controller FeatureController --resource

# Create form request
docker-compose exec app php artisan make:request StoreFeatureRequest

# Create service class
# Manually create app/Services/FeatureService.php
```

#### Step 3: Testing
```bash
# Run all tests
docker-compose exec app php artisan test

# Run specific test suite
docker-compose exec app php artisan test --testsuite=Feature
docker-compose exec app php artisan test --testsuite=Unit

# Run with coverage
docker-compose exec app php artisan test --coverage
```

#### Step 4: Code Quality
```bash
# Format code
docker-compose exec app ./vendor/bin/pint

# Static analysis
docker-compose exec app ./vendor/bin/phpstan analyse

# Check for security vulnerabilities
docker-compose exec app composer audit
```

### 2. Database Development Workflow

#### Migration Management
```bash
# Create migration
docker-compose exec app php artisan make:migration create_table_name

# Run migrations
docker-compose exec app php artisan migrate

# Rollback migrations
docker-compose exec app php artisan migrate:rollback
docker-compose exec app php artisan migrate:rollback --step=3

# Reset database
docker-compose exec app php artisan migrate:fresh --seed

# Check migration status
docker-compose exec app php artisan migrate:status
```

#### Seeding and Factories
```bash
# Create seeder
docker-compose exec app php artisan make:seeder TableSeeder

# Create factory
docker-compose exec app php artisan make:factory ModelFactory

# Run seeders
docker-compose exec app php artisan db:seed
docker-compose exec app php artisan db:seed --class=SpecificSeeder

# Generate test data
docker-compose exec app php artisan tinker
>>> User::factory()->count(50)->create()
```

### 3. API Development Workflow

#### Creating API Endpoints
```bash
# Create API controller
docker-compose exec app php artisan make:controller Api/UserController --api

# Create API resource
docker-compose exec app php artisan make:resource UserResource
docker-compose exec app php artisan make:resource UserCollection

# Create form requests
docker-compose exec app php artisan make:request StoreUserRequest
docker-compose exec app php artisan make:request UpdateUserRequest
```

#### API Testing
```bash
# Test API endpoints
docker-compose exec app php artisan test --filter=Api

# Manual API testing with curl
curl -X GET http://localhost:8000/api/users \
  -H "Accept: application/json" \
  -H "Authorization: Bearer token"
```

### 4. Frontend Development Workflow

#### Asset Development
```bash
# Watch for changes (development)
docker-compose exec app npm run dev

# Build for production
docker-compose exec app npm run build

# Hot reload (if configured)
docker-compose exec app npm run hot
```

#### Blade Component Development
```bash
# Create Blade component
docker-compose exec app php artisan make:component UserCard

# Create view component
docker-compose exec app php artisan make:component UserList --view
```

## Testing Workflows

### 1. Unit Testing
```bash
# Create unit test
docker-compose exec app php artisan make:test UserServiceTest --unit

# Run unit tests
docker-compose exec app php artisan test --testsuite=Unit
```

### 2. Feature Testing
```bash
# Create feature test
docker-compose exec app php artisan make:test UserManagementTest

# Run feature tests
docker-compose exec app php artisan test --testsuite=Feature
```

### 3. Database Testing
```bash
# Test with fresh database
docker-compose exec app php artisan test --env=testing

# Test specific database operations
docker-compose exec app php artisan test --filter=database
```

## Docker Development Commands

### Container Management
```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose stop

# Restart specific service
docker-compose restart app

# Rebuild containers
docker-compose up --build -d

# View resource usage
docker-compose top
```

### Container Access
```bash
# Access app container
docker-compose exec app bash

# Access database
docker-compose exec mysql mysql -u root -p

# Access Redis
docker-compose exec redis redis-cli

# Run commands without entering container
docker-compose exec app php artisan --version
```

### Log Management
```bash
# View all logs
docker-compose logs

# Follow logs for specific service
docker-compose logs -f app

# View last N lines
docker-compose logs --tail=50 app

# Clear logs
docker-compose logs --no-log-prefix app > /dev/null
```

## Laravel Artisan Commands

### Code Generation
```bash
# Models
php artisan make:model User -mfs  # with migration, factory, seeder
php artisan make:model Post -a    # with all related files

# Controllers
php artisan make:controller UserController --resource
php artisan make:controller Api/UserController --api

# Middleware
php artisan make:middleware CheckAge

# Requests
php artisan make:request StoreUserRequest

# Resources
php artisan make:resource UserResource

# Jobs
php artisan make:job ProcessPayment

# Events & Listeners
php artisan make:event UserRegistered
php artisan make:listener SendWelcomeEmail
```

### Application Management
```bash
# Clear caches
php artisan cache:clear
php artisan config:clear
php artisan route:clear
php artisan view:clear

# Optimize for production
php artisan optimize

# Generate application key
php artisan key:generate

# Storage link
php artisan storage:link
```

## Debugging and Troubleshooting

### Common Issues and Solutions

#### Database Connection Issues
```bash
# Check database connection
docker-compose exec app php artisan tinker
>>> DB::connection()->getPdo()

# Restart database service
docker-compose restart mysql

# Check database logs
docker-compose logs mysql
```

#### Cache Issues
```bash
# Clear all caches
docker-compose exec app php artisan cache:clear
docker-compose exec app php artisan config:clear
docker-compose exec app php artisan route:clear
docker-compose exec app php artisan view:clear
```

#### Permission Issues
```bash
# Fix storage permissions
docker-compose exec app chown -R www-data:www-data storage
docker-compose exec app chmod -R 775 storage bootstrap/cache
```

#### Performance Issues
```bash
# Enable query logging
docker-compose exec app php artisan tinker
>>> DB::listen(function($query) { dump($query->sql); });

# Check slow queries
docker-compose exec mysql mysql -u root -p
mysql> SET GLOBAL slow_query_log = 'ON';
mysql> SET GLOBAL long_query_time = 1;
```

## Code Quality and Standards

### Code Formatting
```bash
# Format code with Laravel Pint
docker-compose exec app ./vendor/bin/pint

# Check formatting without fixing
docker-compose exec app ./vendor/bin/pint --test
```

### Static Analysis
```bash
# Run PHPStan
docker-compose exec app ./vendor/bin/phpstan analyse

# Run with specific level
docker-compose exec app ./vendor/bin/phpstan analyse --level=5
```

### Security Scanning
```bash
# Check for security vulnerabilities
docker-compose exec app composer audit

# Update dependencies
docker-compose exec app composer update
```

## Git Workflow

### Branch Management
```bash
# Create feature branch
git checkout -b feature/user-authentication

# Create hotfix branch
git checkout -b hotfix/security-patch

# Create release branch
git checkout -b release/v1.2.0
```

### Commit Guidelines
```bash
# Commit message format
git commit -m "feat: add user authentication"
git commit -m "fix: resolve database connection issue"
git commit -m "docs: update API documentation"
git commit -m "refactor: improve user service"
```

### Pull Request Process
1. Create feature branch
2. Implement changes with tests
3. Run code quality checks
4. Submit pull request
5. Code review and approval
6. Merge to develop/main

## Deployment Workflow

### Production Deployment
```bash
# Build for production
docker-compose exec app npm run build

# Optimize Laravel
docker-compose exec app php artisan optimize
docker-compose exec app php artisan config:cache
docker-compose exec app php artisan route:cache
docker-compose exec app php artisan view:cache

# Run migrations
docker-compose exec app php artisan migrate --force

# Clear and warm up caches
docker-compose exec app php artisan cache:clear
docker-compose exec app php artisan queue:restart
```

### Health Checks
```bash
# Check application health
curl -f http://localhost:8000/health || exit 1

# Check database connectivity
docker-compose exec app php artisan tinker
>>> DB::connection()->getPdo()

# Check queue workers
docker-compose exec app php artisan queue:monitor
```

## Performance Monitoring

### Application Metrics
```bash
# Enable debug mode (development only)
APP_DEBUG=true

# Monitor query performance
docker-compose exec app php artisan tinker
>>> DB::listen(function($query) { 
    echo $query->sql . ' - ' . $query->time . 'ms' . PHP_EOL; 
});
```

### Resource Monitoring
```bash
# Monitor container resources
docker stats

# Check disk usage
docker system df

# Clean up unused resources
docker system prune -a
```

## Backup and Recovery

### Database Backup
```bash
# Create database backup
docker-compose exec mysql mysqldump -u root -p laravel > backup.sql

# Restore database
docker-compose exec -T mysql mysql -u root -p laravel < backup.sql
```

### File System Backup
```bash
# Backup storage directory
tar -czf storage-backup.tar.gz storage/

# Backup entire project
tar -czf project-backup.tar.gz . --exclude=vendor --exclude=node_modules
```

This comprehensive workflow guide provides all the essential commands and processes for developing, testing, and maintaining the {{project}} application. Use these workflows as a reference for consistent development practices.