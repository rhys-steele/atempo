# {{project}} Express.js Codebase Map

## Directory Structure Overview

```
{{project}}/
├── src/                           # Application source code
│   ├── server.js                  # Main application entry point
│   ├── routes/                    # Route definitions
│   │   ├── index.js              # Main route aggregator
│   │   ├── api/                  # API route modules
│   │   │   ├── users.js          # User-related endpoints
│   │   │   └── auth.js           # Authentication endpoints
│   │   └── health.js             # Health check routes
│   ├── controllers/              # Request/response logic
│   │   ├── userController.js     # User operations
│   │   └── authController.js     # Authentication logic
│   ├── middleware/               # Custom middleware
│   │   ├── auth.js              # Authentication middleware
│   │   ├── validation.js        # Input validation
│   │   └── errorHandler.js      # Error handling
│   ├── services/                # Business logic layer
│   │   ├── userService.js       # User business logic
│   │   └── authService.js       # Authentication services
│   ├── models/                  # Data models
│   │   ├── User.js              # User model (Mongoose)
│   │   └── index.js             # Model exports
│   ├── utils/                   # Utility functions
│   │   ├── logger.js            # Logging utilities
│   │   ├── validators.js        # Custom validators
│   │   └── helpers.js           # Helper functions
│   ├── config/                  # Configuration files
│   │   ├── database.js          # Database configuration
│   │   ├── redis.js             # Redis configuration
│   │   └── app.js               # App configuration
│   └── tests/                   # Test files
│       ├── unit/                # Unit tests
│       ├── integration/         # Integration tests
│       └── fixtures/            # Test data
├── infra/                       # Infrastructure files
│   └── docker/                  # Docker configuration
│       ├── Dockerfile           # Node.js container setup
│       ├── docker-compose.yml   # Multi-service orchestration
│       └── .env.example         # Environment template
├── package.json                 # Node.js dependencies and scripts
├── .gitignore                   # Git ignore patterns
├── .eslintrc.js                 # ESLint configuration
├── .prettierrc                  # Prettier configuration
└── README.md                    # Project documentation
```

## File Descriptions & Responsibilities

### Application Entry Point
- **`src/server.js`**: Main application server
  - Express app initialization and middleware setup
  - Security middleware (helmet, cors, rate limiting)
  - Route registration and error handling
  - Server startup and graceful shutdown

### Route Layer
- **`src/routes/index.js`**: Main route aggregator
  - Aggregates all route modules
  - API versioning and route organization
  - Health check and root route definitions

- **`src/routes/api/users.js`**: User-related API endpoints
  - RESTful user operations (CRUD)
  - Route-specific middleware integration
  - Request validation and parameter handling

### Controller Layer
- **`src/controllers/userController.js`**: User request handling
  - HTTP request/response processing
  - Input validation and sanitization
  - Service layer delegation
  - Response formatting and error handling

### Middleware System
- **`src/middleware/auth.js`**: Authentication middleware
  - JWT token validation
  - User authentication and authorization
  - Protected route access control

- **`src/middleware/validation.js`**: Input validation middleware
  - Request data validation using express-validator
  - Custom validation rules
  - Error response formatting

- **`src/middleware/errorHandler.js`**: Global error handling
  - Centralized error processing
  - Error logging and monitoring
  - Consistent error response format

### Service Layer
- **`src/services/userService.js`**: User business logic
  - Complex user operations and business rules
  - Data transformation and processing
  - Integration with external services
  - Database interaction coordination

### Data Layer
- **`src/models/User.js`**: User data model (Mongoose)
  - MongoDB schema definition
  - Validation rules and constraints
  - Model methods and static functions
  - Relationship definitions

### Utility Functions
- **`src/utils/logger.js`**: Logging utilities
  - Winston logger configuration
  - Log formatting and transport setup
  - Different log levels and contexts

- **`src/utils/validators.js`**: Custom validation functions
  - Reusable validation logic
  - Complex validation rules
  - Data sanitization functions

### Configuration Management
- **`src/config/database.js`**: Database configuration
  - MongoDB connection setup
  - Connection pooling and retry logic
  - Environment-specific settings

- **`src/config/redis.js`**: Redis configuration
  - Redis client setup and connection
  - Caching and session configuration
  - Connection error handling

### Testing Infrastructure
- **`src/tests/unit/`**: Unit test files
  - Individual function and method testing
  - Mocked dependencies and isolated testing
  - Service and utility function tests

- **`src/tests/integration/`**: Integration test files
  - API endpoint testing with supertest
  - Database integration testing
  - Multi-service interaction testing

## Architecture Patterns

### Express.js Middleware Pattern
```javascript
// Middleware execution flow
app.use(helmet());                 // Security headers
app.use(cors());                   // CORS handling
app.use(express.json());           // Body parsing
app.use('/api', apiRoutes);        // Route mounting
app.use(errorHandler);             // Error handling
```

### Controller Pattern
```javascript
// Controller structure
const asyncHandler = (fn) => (req, res, next) => {
  Promise.resolve(fn(req, res, next)).catch(next);
};

const getUsers = asyncHandler(async (req, res) => {
  const users = await userService.getAllUsers();
  res.json({ success: true, data: users });
});
```

### Service Layer Pattern
```javascript
// Service layer abstraction
class UserService {
  async createUser(userData) {
    // Business logic
    const user = await User.create(userData);
    // Additional processing
    return user;
  }
}
```

## Key Relationships & Data Flow

### Request Processing Flow
1. **HTTP Request** → Express server receives request
2. **Middleware Stack** → Security, parsing, logging middleware
3. **Route Matching** → Express router matches URL pattern
4. **Controller** → Request handler processes input
5. **Service Layer** → Business logic execution
6. **Data Layer** → Database operations (Mongoose/MongoDB)
7. **Response** → JSON response sent to client

### Error Handling Flow
1. **Error Occurrence** → Error thrown in any layer
2. **Error Propagation** → Error bubbles up through middleware
3. **Global Handler** → Central error processing
4. **Logging** → Error logged with context
5. **Client Response** → Structured error response

### Authentication Flow
1. **Request** → Client sends request with token
2. **Auth Middleware** → Token validation and user extraction
3. **Authorization** → Role/permission checking
4. **Controller** → Protected resource access
5. **Response** → Authorized data or error response

## External Dependencies

### Runtime Dependencies
- **Express.js**: Web framework and HTTP server
- **Mongoose**: MongoDB object modeling
- **Redis**: Caching and session storage
- **Helmet**: Security middleware
- **CORS**: Cross-origin resource sharing

### Development Dependencies
- **Nodemon**: Development server with auto-restart
- **Jest**: Testing framework
- **Supertest**: HTTP assertion library
- **ESLint**: Code linting and quality
- **Prettier**: Code formatting

## Development & Extension Points

### Adding New Features
1. Create route module in `src/routes/api/`
2. Implement controller in `src/controllers/`
3. Add business logic in `src/services/`
4. Define data model in `src/models/`
5. Write tests in `src/tests/`

### Database Extensions
1. Define new Mongoose schemas in `src/models/`
2. Add database migrations or seed data
3. Update connection configuration
4. Implement model relationships

### Middleware Extensions
1. Create custom middleware in `src/middleware/`
2. Add middleware to application stack
3. Implement route-specific middleware
4. Handle middleware errors appropriately

### API Documentation
1. Use OpenAPI/Swagger for API documentation
2. Document request/response schemas
3. Provide example requests and responses
4. Include authentication requirements

This codebase demonstrates modern Express.js application architecture with clear separation of concerns, comprehensive error handling, and extensible design patterns for scalable API development.