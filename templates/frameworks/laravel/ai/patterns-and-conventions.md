# {{project}} - Patterns and Conventions

## Architecture Patterns

### 1. Clean Architecture Implementation

#### Layer Separation
```php
// Controller Layer - Handle HTTP requests
class UserController extends Controller
{
    public function __construct(private UserService $userService) {}
    
    public function store(StoreUserRequest $request): JsonResponse
    {
        $user = $this->userService->createUser($request->validated());
        return new JsonResponse(new UserResource($user), 201);
    }
}

// Service Layer - Business logic
class UserService
{
    public function __construct(private UserRepository $userRepository) {}
    
    public function createUser(array $data): User
    {
        // Business logic here
        return $this->userRepository->create($data);
    }
}

// Repository Layer - Data access
class UserRepository
{
    public function create(array $data): User
    {
        return User::create($data);
    }
}
```

### 2. Service Layer Pattern

#### Service Organization
```php
// app/Services/UserService.php
class UserService
{
    public function __construct(
        private UserRepository $userRepository,
        private EmailService $emailService
    ) {}
    
    public function createUser(array $data): User
    {
        DB::beginTransaction();
        
        try {
            $user = $this->userRepository->create($data);
            $this->emailService->sendWelcomeEmail($user);
            
            DB::commit();
            return $user;
        } catch (Exception $e) {
            DB::rollBack();
            throw new UserCreationException('Failed to create user', 0, $e);
        }
    }
}
```

### 3. Repository Pattern

#### Repository Implementation
```php
// app/Repositories/UserRepository.php
interface UserRepositoryInterface
{
    public function find(int $id): ?User;
    public function create(array $data): User;
    public function update(User $user, array $data): User;
    public function delete(User $user): bool;
}

class UserRepository implements UserRepositoryInterface
{
    public function find(int $id): ?User
    {
        return User::find($id);
    }
    
    public function create(array $data): User
    {
        return User::create($data);
    }
    
    public function update(User $user, array $data): User
    {
        $user->update($data);
        return $user->fresh();
    }
    
    public function delete(User $user): bool
    {
        return $user->delete();
    }
}
```

### 4. Command Pattern

#### Artisan Commands
```php
// app/Console/Commands/ProcessUserData.php
class ProcessUserData extends Command
{
    protected $signature = 'users:process {--batch=100}';
    protected $description = 'Process user data in batches';
    
    public function handle(UserService $userService): int
    {
        $batchSize = $this->option('batch');
        
        $this->info("Processing users in batches of {$batchSize}");
        
        $bar = $this->output->createProgressBar(User::count());
        
        User::chunk($batchSize, function ($users) use ($userService, $bar) {
            foreach ($users as $user) {
                $userService->processUser($user);
                $bar->advance();
            }
        });
        
        $bar->finish();
        $this->info('User processing completed');
        
        return Command::SUCCESS;
    }
}
```

## Naming Conventions

### 1. File and Class Naming

#### Controllers
```php
// Good
class UserController extends Controller          // UserController.php
class Api\UserController extends Controller     // Api/UserController.php
class Admin\UserController extends Controller   // Admin/UserController.php

// Bad
class UsersController extends Controller
class UserCtrl extends Controller
class UserMgr extends Controller
```

#### Models
```php
// Good
class User extends Model                         // User.php
class BlogPost extends Model                     // BlogPost.php
class OrderItem extends Model                    // OrderItem.php

// Bad
class Users extends Model
class blog_post extends Model
class orderitem extends Model
```

#### Services
```php
// Good
class UserService                                // UserService.php
class PaymentService                             // PaymentService.php
class EmailNotificationService                  // EmailNotificationService.php

// Bad
class UserSvc
class PaymentHandler
class EmailSender
```

### 2. Method Naming

#### Controller Methods
```php
// RESTful methods
public function index()      // GET /users
public function show($id)    // GET /users/{id}
public function create()     // GET /users/create
public function store()      // POST /users
public function edit($id)    // GET /users/{id}/edit
public function update($id)  // PUT/PATCH /users/{id}
public function destroy($id) // DELETE /users/{id}
```

#### Service Methods
```php
// Clear, action-oriented names
public function createUser(array $data): User
public function updateUser(User $user, array $data): User
public function deleteUser(User $user): bool
public function findActiveUsers(): Collection
public function sendWelcomeEmail(User $user): void
```

### 3. Variable and Property Naming

#### Variables
```php
// Good
$user = User::find($id);
$activeUsers = User::where('active', true)->get();
$userCount = User::count();

// Bad
$u = User::find($id);
$data = User::where('active', true)->get();
$cnt = User::count();
```

#### Properties
```php
class User extends Model
{
    // Good
    protected $fillable = ['name', 'email', 'password'];
    protected $hidden = ['password', 'remember_token'];
    protected $casts = ['email_verified_at' => 'datetime'];
    
    // Bad
    protected $guarded = [];
    protected $dates = ['email_verified_at'];
}
```

## Error Handling Patterns

### 1. Custom Exception Classes

#### Exception Hierarchy
```php
// app/Exceptions/BaseException.php
abstract class BaseException extends Exception
{
    public function report(): void
    {
        Log::error($this->getMessage(), [
            'exception' => get_class($this),
            'file' => $this->getFile(),
            'line' => $this->getLine(),
            'trace' => $this->getTraceAsString()
        ]);
    }
}

// app/Exceptions/UserException.php
class UserException extends BaseException
{
    public static function notFound(int $id): self
    {
        return new self("User with ID {$id} not found", 404);
    }
    
    public static function emailAlreadyExists(string $email): self
    {
        return new self("User with email {$email} already exists", 409);
    }
}
```

### 2. Service Layer Error Handling

#### Try-Catch Patterns
```php
class UserService
{
    public function createUser(array $data): User
    {
        try {
            DB::beginTransaction();
            
            $user = User::create($data);
            
            // Additional operations
            $this->createUserProfile($user);
            $this->sendWelcomeEmail($user);
            
            DB::commit();
            return $user;
            
        } catch (ValidationException $e) {
            DB::rollBack();
            throw $e;
        } catch (Exception $e) {
            DB::rollBack();
            throw new UserCreationException(
                'Failed to create user: ' . $e->getMessage(),
                0,
                $e
            );
        }
    }
}
```

### 3. API Error Responses

#### Consistent Error Format
```php
// app/Exceptions/Handler.php
class Handler extends ExceptionHandler
{
    public function render($request, Exception $exception)
    {
        if ($request->wantsJson()) {
            return $this->handleApiException($request, $exception);
        }
        
        return parent::render($request, $exception);
    }
    
    private function handleApiException($request, Exception $exception): JsonResponse
    {
        $statusCode = $this->getStatusCode($exception);
        
        return response()->json([
            'error' => [
                'message' => $exception->getMessage(),
                'code' => $statusCode,
                'type' => class_basename($exception)
            ]
        ], $statusCode);
    }
}
```

## Validation Patterns

### 1. Form Request Validation

#### Custom Form Requests
```php
// app/Http/Requests/StoreUserRequest.php
class StoreUserRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true; // or implement authorization logic
    }
    
    public function rules(): array
    {
        return [
            'name' => 'required|string|max:255',
            'email' => 'required|email|unique:users,email',
            'password' => 'required|min:8|confirmed',
            'date_of_birth' => 'required|date|before:today'
        ];
    }
    
    public function messages(): array
    {
        return [
            'name.required' => 'The name field is required.',
            'email.unique' => 'This email address is already registered.',
            'password.min' => 'Password must be at least 8 characters.',
        ];
    }
    
    public function attributes(): array
    {
        return [
            'date_of_birth' => 'date of birth'
        ];
    }
}
```

### 2. Model Validation

#### Custom Validation Rules
```php
// app/Rules/StrongPassword.php
class StrongPassword implements Rule
{
    public function passes($attribute, $value): bool
    {
        return preg_match('/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/', $value);
    }
    
    public function message(): string
    {
        return 'The :attribute must contain at least one uppercase letter, one lowercase letter, one digit, and one special character.';
    }
}

// Usage in Form Request
public function rules(): array
{
    return [
        'password' => ['required', 'min:8', new StrongPassword()]
    ];
}
```

## Database Patterns

### 1. Migration Patterns

#### Schema Design
```php
// database/migrations/create_users_table.php
class CreateUsersTable extends Migration
{
    public function up(): void
    {
        Schema::create('users', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->string('email')->unique();
            $table->timestamp('email_verified_at')->nullable();
            $table->string('password');
            $table->rememberToken();
            $table->timestamps();
            
            // Indexes
            $table->index(['email', 'email_verified_at']);
        });
    }
    
    public function down(): void
    {
        Schema::dropIfExists('users');
    }
}
```

### 2. Model Relationships

#### Relationship Definitions
```php
class User extends Model
{
    // One-to-Many
    public function posts(): HasMany
    {
        return $this->hasMany(Post::class);
    }
    
    // Many-to-Many
    public function roles(): BelongsToMany
    {
        return $this->belongsToMany(Role::class)
                    ->withTimestamps()
                    ->withPivot('assigned_at');
    }
    
    // Has One Through
    public function latestPost(): HasOne
    {
        return $this->hasOne(Post::class)->latestOfMany();
    }
}
```

### 3. Query Scopes

#### Local and Global Scopes
```php
class User extends Model
{
    // Local Scope
    public function scopeActive(Builder $query): Builder
    {
        return $query->where('active', true);
    }
    
    public function scopeVerified(Builder $query): Builder
    {
        return $query->whereNotNull('email_verified_at');
    }
    
    // Usage: User::active()->verified()->get()
}

// Global Scope
class ActiveUserScope implements Scope
{
    public function apply(Builder $builder, Model $model): void
    {
        $builder->where('active', true);
    }
}
```

## Testing Patterns

### 1. Test Organization

#### Test Structure
```php
// tests/Feature/UserManagementTest.php
class UserManagementTest extends TestCase
{
    use RefreshDatabase;
    
    /** @test */
    public function it_can_create_a_new_user(): void
    {
        // Arrange
        $userData = [
            'name' => 'John Doe',
            'email' => 'john@example.com',
            'password' => 'password123',
            'password_confirmation' => 'password123'
        ];
        
        // Act
        $response = $this->postJson('/api/users', $userData);
        
        // Assert
        $response->assertStatus(201)
                 ->assertJsonStructure([
                     'data' => ['id', 'name', 'email', 'created_at']
                 ]);
                 
        $this->assertDatabaseHas('users', [
            'name' => 'John Doe',
            'email' => 'john@example.com'
        ]);
    }
}
```

### 2. Factory Patterns

#### Model Factories
```php
// database/factories/UserFactory.php
class UserFactory extends Factory
{
    public function definition(): array
    {
        return [
            'name' => $this->faker->name(),
            'email' => $this->faker->unique()->safeEmail(),
            'email_verified_at' => now(),
            'password' => Hash::make('password'),
            'remember_token' => Str::random(10),
        ];
    }
    
    public function unverified(): static
    {
        return $this->state(fn (array $attributes) => [
            'email_verified_at' => null,
        ]);
    }
    
    public function admin(): static
    {
        return $this->state(fn (array $attributes) => [
            'is_admin' => true,
        ]);
    }
}
```

### 3. Test Helpers

#### Custom Assertions
```php
// tests/TestCase.php
abstract class TestCase extends BaseTestCase
{
    use CreatesApplication;
    
    protected function assertValidationError($response, $field): void
    {
        $response->assertStatus(422)
                 ->assertJsonValidationErrors($field);
    }
    
    protected function actingAsUser(?User $user = null): TestCase
    {
        $user = $user ?? User::factory()->create();
        return $this->actingAs($user);
    }
    
    protected function actingAsAdmin(?User $admin = null): TestCase
    {
        $admin = $admin ?? User::factory()->admin()->create();
        return $this->actingAs($admin);
    }
}
```

## Security Patterns

### 1. Authorization Patterns

#### Policy Classes
```php
// app/Policies/UserPolicy.php
class UserPolicy
{
    public function view(User $user, User $targetUser): bool
    {
        return $user->id === $targetUser->id || $user->isAdmin();
    }
    
    public function update(User $user, User $targetUser): bool
    {
        return $user->id === $targetUser->id;
    }
    
    public function delete(User $user, User $targetUser): bool
    {
        return $user->isAdmin() && $user->id !== $targetUser->id;
    }
}
```

### 2. Input Sanitization

#### Mutators and Accessors
```php
class User extends Model
{
    // Mutator - sanitize input
    public function setEmailAttribute($value): void
    {
        $this->attributes['email'] = strtolower(trim($value));
    }
    
    // Accessor - format output
    public function getFullNameAttribute(): string
    {
        return "{$this->first_name} {$this->last_name}";
    }
    
    // Cast attributes
    protected $casts = [
        'email_verified_at' => 'datetime',
        'settings' => 'array',
        'is_admin' => 'boolean'
    ];
}
```

## Performance Patterns

### 1. Caching Strategies

#### Model Caching
```php
class UserService
{
    public function findUser(int $id): ?User
    {
        return Cache::remember("user.{$id}", now()->addHour(), function () use ($id) {
            return User::with(['profile', 'roles'])->find($id);
        });
    }
    
    public function clearUserCache(int $id): void
    {
        Cache::forget("user.{$id}");
    }
}
```

### 2. Database Optimization

#### Query Optimization
```php
class UserService
{
    public function getUsersWithPosts(): Collection
    {
        // Eager loading to prevent N+1 queries
        return User::with(['posts' => function ($query) {
            $query->select('id', 'user_id', 'title', 'created_at')
                  ->orderBy('created_at', 'desc');
        }])->get();
    }
    
    public function getActiveUsersCount(): int
    {
        // Use raw queries for better performance
        return DB::table('users')
                 ->where('active', true)
                 ->count();
    }
}
```

## Configuration Patterns

### 1. Service Configuration

#### Service Provider Pattern
```php
// app/Providers/UserServiceProvider.php
class UserServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->app->bind(UserRepositoryInterface::class, UserRepository::class);
        $this->app->bind(UserServiceInterface::class, UserService::class);
    }
    
    public function boot(): void
    {
        // Register observers
        User::observe(UserObserver::class);
        
        // Register policies
        Gate::policy(User::class, UserPolicy::class);
    }
}
```

### 2. Environment Configuration

#### Configuration Files
```php
// config/services.php
return [
    'email' => [
        'driver' => env('MAIL_DRIVER', 'smtp'),
        'host' => env('MAIL_HOST', 'localhost'),
        'port' => env('MAIL_PORT', 587),
        'from' => [
            'address' => env('MAIL_FROM_ADDRESS', 'hello@example.com'),
            'name' => env('MAIL_FROM_NAME', 'Example')
        ]
    ]
];
```

These patterns and conventions provide a solid foundation for building maintainable, scalable Laravel applications. Follow these guidelines to ensure consistency and quality across the {{project}} codebase.