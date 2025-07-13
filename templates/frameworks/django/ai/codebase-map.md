# {{project}} - Codebase Map

## Project Structure Overview

```
{{project}}/
├── app/                           # Application Core (Laravel)
│   ├── Console/                   # Artisan commands
│   ├── Exceptions/                # Exception handlers
│   ├── Http/                      # HTTP layer
│   │   ├── Controllers/           # Request handlers
│   │   ├── Middleware/            # Request/response filters
│   │   ├── Requests/              # Form request validation
│   │   └── Resources/             # API response formatting
│   ├── Models/                    # Eloquent models
│   ├── Providers/                 # Service providers
│   └── Services/                  # Business logic layer
├── bootstrap/                     # Application bootstrapping
├── config/                        # Configuration files
├── database/                      # Database components
│   ├── factories/                 # Model factories
│   ├── migrations/                # Database migrations
│   └── seeders/                   # Database seeders
├── infra/                         # Infrastructure files
│   └── docker/                    # Docker configuration
├── public/                        # Public web assets
├── resources/                     # Views, assets, lang files
│   ├── css/                       # Stylesheets
│   ├── js/                        # JavaScript files
│   ├── lang/                      # Localization files
│   └── views/                     # Blade templates
├── routes/                        # Application routes
├── storage/                       # Storage directories
├── tests/                         # Test suites
│   ├── Feature/                   # Feature tests
│   └── Unit/                      # Unit tests
├── vendor/                        # Composer dependencies
├── .env                           # Environment configuration
├── artisan                        # Artisan CLI
├── composer.json                  # PHP dependencies
├── docker-compose.yml             # Docker services
└── package.json                   # Node.js dependencies
```

## Core Application Components

### HTTP Layer (`app/Http/`)

#### Controllers (`app/Http/Controllers/`)
**Purpose**: Handle HTTP requests and coordinate application flow
**Pattern**: Thin controllers that delegate to services
**Key Files**:
- `Controller.php`: Base controller with common functionality
- `HomeController.php`: Homepage and dashboard logic
- `Auth/`: Authentication controllers
- `API/`: API endpoint controllers

#### Middleware (`app/Http/Middleware/`)
**Purpose**: Filter HTTP requests and responses
**Pattern**: Chain of responsibility for request processing
**Key Files**:
- `Authenticate.php`: Authentication verification
- `VerifyCsrfToken.php`: CSRF protection
- `TrustProxies.php`: Proxy configuration

#### Requests (`app/Http/Requests/`)
**Purpose**: Validate and authorize incoming requests
**Pattern**: Form request validation with custom rules
**Structure**:
```php
class CreateUserRequest extends FormRequest
{
    public function authorize(): bool
    public function rules(): array
    public function messages(): array
}
```

#### Resources (`app/Http/Resources/`)
**Purpose**: Transform models into API responses
**Pattern**: Resource transformation layer
**Structure**:
```php
class UserResource extends JsonResource
{
    public function toArray($request): array
}
```

### Models (`app/Models/`)
**Purpose**: Eloquent ORM models for database interaction
**Pattern**: Active Record pattern with relationships
**Key Concepts**:
- Model relationships (hasMany, belongsTo, etc.)
- Accessors and mutators
- Scopes and query builders
- Model events and observers

**Example Structure**:
```php
class User extends Authenticatable
{
    protected $fillable = [];
    protected $hidden = [];
    protected $casts = [];
    
    // Relationships
    public function posts(): HasMany
    
    // Accessors/Mutators
    public function getFullNameAttribute(): string
    
    // Scopes
    public function scopeActive(Builder $query): Builder
}
```

### Services (`app/Services/`)
**Purpose**: Business logic and complex operations
**Pattern**: Service layer pattern for reusable business logic
**Structure**:
```php
class UserService
{
    public function createUser(array $data): User
    public function updateUser(User $user, array $data): User
    public function deleteUser(User $user): bool
}
```

### Providers (`app/Providers/`)
**Purpose**: Service container registration and bootstrapping
**Key Files**:
- `AppServiceProvider.php`: General application services
- `RouteServiceProvider.php`: Route registration
- `AuthServiceProvider.php`: Authentication and authorization

## Database Layer

### Migrations (`database/migrations/`)
**Purpose**: Database schema version control
**Pattern**: Up/down migration methods
**Naming Convention**: `YYYY_MM_DD_HHMMSS_create_table_name.php`
**Structure**:
```php
class CreateUsersTable extends Migration
{
    public function up(): void
    public function down(): void
}
```

### Seeders (`database/seeders/`)
**Purpose**: Database data population
**Pattern**: Seeder classes for test and initial data
**Structure**:
```php
class UserSeeder extends Seeder
{
    public function run(): void
}
```

### Factories (`database/factories/`)
**Purpose**: Generate model instances for testing
**Pattern**: Factory pattern for test data generation
**Structure**:
```php
class UserFactory extends Factory
{
    public function definition(): array
}
```

## Configuration Layer

### Environment Configuration (`.env`)
**Purpose**: Environment-specific settings
**Key Sections**:
- Database connection settings
- Cache and session configuration
- Mail and queue settings
- Third-party API keys

### Config Files (`config/`)
**Purpose**: Application configuration
**Key Files**:
- `app.php`: Application settings
- `database.php`: Database connections
- `cache.php`: Cache configuration
- `mail.php`: Email settings

## Testing Architecture

### Test Organization
```
tests/
├── Feature/                   # End-to-end feature tests
│   ├── Auth/                  # Authentication tests
│   ├── API/                   # API endpoint tests
│   └── Web/                   # Web interface tests
├── Unit/                      # Unit tests
│   ├── Models/                # Model tests
│   ├── Services/              # Service tests
│   └── Helpers/               # Helper function tests
├── TestCase.php               # Base test class
└── CreatesApplication.php     # Application factory
```

### Testing Patterns
- **Feature Tests**: Test complete user workflows
- **Unit Tests**: Test individual components in isolation
- **Database Tests**: Use RefreshDatabase trait
- **HTTP Tests**: Test API endpoints and responses

## Infrastructure Layer

### Docker Configuration (`infra/docker/`)
**Purpose**: Containerized development environment
**Key Files**:
- `Dockerfile`: Application container definition
- `nginx.conf`: Nginx web server configuration
- `docker-compose.yml`: Service orchestration

### Docker Services
```yaml
services:
  app:          # Laravel application
  webserver:    # Nginx reverse proxy
  mysql:        # Database server
  redis:        # Cache and session store
  mailhog:      # Email testing
```

## Frontend Architecture

### Views (`resources/views/`)
**Purpose**: Blade templates for server-side rendering
**Structure**:
- `layouts/`: Base layouts and components
- `components/`: Reusable UI components
- `pages/`: Individual page templates

### Assets (`resources/`)
**Purpose**: Frontend assets and build pipeline
**Structure**:
- `css/`: Stylesheets (CSS/SCSS)
- `js/`: JavaScript files
- `lang/`: Localization files

## API Design

### RESTful Endpoints
**Pattern**: Resource-based URLs with HTTP methods
**Structure**:
```
GET    /api/users           # List users
POST   /api/users           # Create user
GET    /api/users/{id}      # Show user
PUT    /api/users/{id}      # Update user
DELETE /api/users/{id}      # Delete user
```

### Response Format
**Pattern**: Consistent JSON response structure
**Structure**:
```json
{
  "data": {},
  "meta": {
    "status": "success",
    "message": "Operation completed"
  }
}
```

## Security Architecture

### Authentication Flow
1. User submits credentials
2. Middleware verifies authentication
3. Controller processes request
4. Service performs business logic
5. Response returned with proper status

### Authorization Patterns
- **Policies**: Model-based authorization
- **Gates**: Simple boolean authorization
- **Middleware**: Route-level protection
- **Guards**: Authentication drivers

## Performance Optimization

### Caching Strategy
- **Redis**: Session storage and application cache
- **Database**: Query result caching
- **HTTP**: Response caching with ETags
- **Static Assets**: CDN integration

### Database Optimization
- **Indexes**: Proper indexing strategy
- **Relationships**: Eager loading to prevent N+1
- **Query Optimization**: Use DB::listen() for monitoring
- **Connection Pooling**: Efficient connection management

## Error Handling

### Exception Hierarchy
```php
App\Exceptions\Handler extends ExceptionHandler
├── ValidationException
├── AuthenticationException
├── AuthorizationException
└── CustomBusinessException
```

### Logging Strategy
- **Channels**: Separate logs for different concerns
- **Levels**: Debug, Info, Warning, Error, Critical
- **Context**: Include relevant request/user information
- **Monitoring**: Integration with error tracking services

## Integration Points

### External Services
- **Email**: SMTP/API integration
- **File Storage**: Local/S3/CDN integration
- **Payment**: Payment gateway integration
- **Analytics**: Tracking and reporting

### Internal Dependencies
- **Queues**: Background job processing
- **Cache**: Redis for performance
- **Database**: MySQL for persistence
- **Session**: Redis for session storage

## Development Patterns

### Design Patterns Used
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic encapsulation
- **Factory Pattern**: Object creation
- **Observer Pattern**: Event-driven architecture
- **Decorator Pattern**: Middleware implementation

### Code Organization Principles
- **Single Responsibility**: Each class has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Dependency Inversion**: Depend on abstractions, not concretions
- **DRY**: Don't Repeat Yourself
- **KISS**: Keep It Simple, Stupid

## Extension Points

### Custom Artisan Commands
**Location**: `app/Console/Commands/`
**Purpose**: Custom CLI commands for maintenance and operations

### Event Listeners
**Location**: `app/Listeners/`
**Purpose**: Handle application events and side effects

### Service Providers
**Location**: `app/Providers/`
**Purpose**: Register services and configure application

### Middleware
**Location**: `app/Http/Middleware/`
**Purpose**: Request/response filtering and modification

This codebase map provides a comprehensive guide to understanding the {{project}} application structure, patterns, and architectural decisions. Use this as a reference for navigating the codebase and implementing new features according to established patterns.