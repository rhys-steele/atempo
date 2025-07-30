# {{project}} Patterns & Conventions

## Code Architecture Patterns

### Express.js Middleware Pattern
The application follows Express.js middleware-based architecture with clear execution flow:

```javascript
// Middleware execution stack
app.use(helmet());                    // Security headers
app.use(cors());                      // CORS handling  
app.use(express.json());              // Body parsing
app.use(rateLimit(limitConfig));      // Rate limiting
app.use('/api', apiRoutes);           // Route mounting
app.use(errorHandler);                // Global error handling
```

### MVC Pattern Implementation
Model-View-Controller pattern adapted for API development:

```javascript
// Route Layer (View)
router.get('/users/:id', userController.getUser);

// Controller Layer
const getUser = async (req, res, next) => {
  try {
    const user = await userService.getUserById(req.params.id);
    res.json({ success: true, data: user });
  } catch (error) {
    next(error);
  }
};

// Service Layer (Business Logic)
const getUserById = async (id) => {
  const user = await User.findById(id);
  if (!user) throw new Error('User not found');
  return user;
};

// Model Layer
const userSchema = new mongoose.Schema({
  name: { type: String, required: true },
  email: { type: String, unique: true, required: true }
});
```

### Service Layer Pattern
Separation of business logic from controllers:

```javascript
// Service class structure
class UserService {
  async createUser(userData) {
    // Validation
    await this.validateUserData(userData);
    
    // Business logic
    const hashedPassword = await bcrypt.hash(userData.password, 10);
    
    // Data persistence
    const user = await User.create({
      ...userData,
      password: hashedPassword
    });
    
    // Post-processing
    await this.sendWelcomeEmail(user);
    
    return user;
  }
  
  async validateUserData(userData) {
    if (!userData.email) {
      throw new ValidationError('Email is required');
    }
    
    const existingUser = await User.findOne({ email: userData.email });
    if (existingUser) {
      throw new ConflictError('Email already exists');
    }
  }
}
```

## Naming Conventions

### File and Directory Naming
- **camelCase** for files: `userController.js`, `authService.js`
- **PascalCase** for models: `User.js`, `Product.js`
- **kebab-case** for multi-word routes: `user-profile.js`
- **Descriptive names**: Files clearly indicate their purpose

### JavaScript Naming Conventions
```javascript
// Variables and functions: camelCase
const userName = 'john_doe';
const getUserById = async (id) => { };

// Constants: UPPER_SNAKE_CASE
const MAX_LOGIN_ATTEMPTS = 5;
const DEFAULT_PAGE_SIZE = 20;

// Classes: PascalCase
class UserService { }
class ValidationError extends Error { }

// Private methods: leading underscore
class UserService {
  async getUser(id) { }
  async _validateId(id) { }  // Private method
}
```

### Route Naming
```javascript
// RESTful route conventions
GET    /api/users          // List users
POST   /api/users          // Create user
GET    /api/users/:id      // Get specific user
PUT    /api/users/:id      // Update user
DELETE /api/users/:id      // Delete user

// Nested resources
GET    /api/users/:id/posts     // User's posts
POST   /api/users/:id/posts     // Create user's post

// Action-based routes
POST   /api/users/:id/activate   // Activate user
POST   /api/auth/login           // User login
POST   /api/auth/logout          // User logout
```

## Error Handling Patterns

### Custom Error Classes
```javascript
// Base error class
class AppError extends Error {
  constructor(message, statusCode = 500) {
    super(message);
    this.statusCode = statusCode;
    this.isOperational = true;
    
    Error.captureStackTrace(this, this.constructor);
  }
}

// Specific error types
class ValidationError extends AppError {
  constructor(message) {
    super(message, 400);
  }
}

class NotFoundError extends AppError {
  constructor(resource) {
    super(`${resource} not found`, 404);
  }
}

class ConflictError extends AppError {
  constructor(message) {
    super(message, 409);
  }
}
```

### Async Error Handling
```javascript
// Async wrapper utility
const asyncHandler = (fn) => (req, res, next) => {
  Promise.resolve(fn(req, res, next)).catch(next);
};

// Usage in controllers
const getUser = asyncHandler(async (req, res) => {
  const user = await userService.getUserById(req.params.id);
  res.json({ success: true, data: user });
});

// Global error handler
const errorHandler = (err, req, res, next) => {
  let error = { ...err };
  error.message = err.message;

  // Log error
  console.error(err);

  // Mongoose validation error
  if (err.name === 'ValidationError') {
    const message = Object.values(err.errors).map(val => val.message);
    error = new ValidationError(message);
  }

  // Mongoose duplicate key
  if (err.code === 11000) {
    const message = 'Duplicate field value entered';
    error = new ConflictError(message);
  }

  res.status(error.statusCode || 500).json({
    success: false,
    error: error.message || 'Server Error'
  });
};
```

## Validation Patterns

### Input Validation with express-validator
```javascript
// Validation middleware
const { body, param, query, validationResult } = require('express-validator');

// User validation rules
const userValidation = {
  create: [
    body('name')
      .trim()
      .isLength({ min: 2, max: 50 })
      .withMessage('Name must be between 2 and 50 characters'),
    
    body('email')
      .isEmail()
      .normalizeEmail()
      .withMessage('Please provide a valid email'),
    
    body('password')
      .isLength({ min: 8 })
      .matches(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/)
      .withMessage('Password must contain at least 8 characters, one uppercase, one lowercase, and one number')
  ],
  
  update: [
    param('id').isMongoId().withMessage('Invalid user ID'),
    body('name').optional().trim().isLength({ min: 2, max: 50 }),
    body('email').optional().isEmail().normalizeEmail()
  ]
};

// Validation handler
const handleValidation = (req, res, next) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(400).json({
      success: false,
      errors: errors.array().map(err => ({
        field: err.param,
        message: err.msg,
        value: err.value
      }))
    });
  }
  next();
};

// Route usage
router.post('/users', userValidation.create, handleValidation, userController.createUser);
```

### Schema Validation with Mongoose
```javascript
// Mongoose schema with validation
const userSchema = new mongoose.Schema({
  name: {
    type: String,
    required: [true, 'Name is required'],
    trim: true,
    minlength: [2, 'Name must be at least 2 characters'],
    maxlength: [50, 'Name cannot exceed 50 characters']
  },
  
  email: {
    type: String,
    required: [true, 'Email is required'],
    unique: true,
    lowercase: true,
    validate: {
      validator: function(email) {
        return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
      },
      message: 'Please provide a valid email address'
    }
  },
  
  role: {
    type: String,
    enum: {
      values: ['user', 'admin', 'moderator'],
      message: 'Role must be either user, admin, or moderator'
    },
    default: 'user'
  },
  
  createdAt: {
    type: Date,
    default: Date.now
  }
});

// Pre-save middleware
userSchema.pre('save', async function(next) {
  if (!this.isModified('password')) return next();
  
  this.password = await bcrypt.hash(this.password, 12);
  next();
});
```

## Response Patterns

### Consistent API Response Format
```javascript
// Success response format
const successResponse = (res, data, message = 'Success', statusCode = 200) => {
  res.status(statusCode).json({
    success: true,
    message,
    data,
    timestamp: new Date().toISOString()
  });
};

// Error response format
const errorResponse = (res, message, statusCode = 500, errors = null) => {
  res.status(statusCode).json({
    success: false,
    message,
    errors,
    timestamp: new Date().toISOString()
  });
};

// Pagination response
const paginatedResponse = (res, data, pagination) => {
  res.json({
    success: true,
    data,
    pagination: {
      page: pagination.page,
      limit: pagination.limit,
      total: pagination.total,
      pages: Math.ceil(pagination.total / pagination.limit)
    },
    timestamp: new Date().toISOString()
  });
};
```

## Authentication Patterns

### JWT Authentication
```javascript
// JWT middleware
const jwt = require('jsonwebtoken');

const authenticateToken = (req, res, next) => {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1]; // Bearer TOKEN

  if (!token) {
    return res.status(401).json({ 
      success: false, 
      message: 'Access token required' 
    });
  }

  jwt.verify(token, process.env.JWT_SECRET, (err, user) => {
    if (err) {
      return res.status(403).json({ 
        success: false, 
        message: 'Invalid or expired token' 
      });
    }
    
    req.user = user;
    next();
  });
};

// Authorization middleware
const authorize = (...roles) => {
  return (req, res, next) => {
    if (!roles.includes(req.user.role)) {
      return res.status(403).json({
        success: false,
        message: 'Insufficient permissions'
      });
    }
    next();
  };
};

// Usage
router.get('/admin/users', 
  authenticateToken, 
  authorize('admin'), 
  userController.getAllUsers
);
```

## Database Patterns

### Repository Pattern (Optional)
```javascript
// Base repository
class BaseRepository {
  constructor(model) {
    this.model = model;
  }

  async findById(id) {
    return await this.model.findById(id);
  }

  async findOne(filter) {
    return await this.model.findOne(filter);
  }

  async create(data) {
    return await this.model.create(data);
  }

  async update(id, data) {
    return await this.model.findByIdAndUpdate(id, data, { 
      new: true, 
      runValidators: true 
    });
  }

  async delete(id) {
    return await this.model.findByIdAndDelete(id);
  }

  async findWithPagination(filter = {}, options = {}) {
    const { page = 1, limit = 10, sort = { createdAt: -1 } } = options;
    const skip = (page - 1) * limit;

    const [data, total] = await Promise.all([
      this.model.find(filter)
        .sort(sort)
        .limit(limit)
        .skip(skip),
      this.model.countDocuments(filter)
    ]);

    return { data, total, page, limit };
  }
}

// User repository
class UserRepository extends BaseRepository {
  constructor() {
    super(User);
  }

  async findByEmail(email) {
    return await this.model.findOne({ email });
  }

  async findActiveUsers() {
    return await this.model.find({ active: true });
  }
}
```

### Query Optimization
```javascript
// Efficient queries
const getUsers = async (options = {}) => {
  const { 
    page = 1, 
    limit = 10, 
    sort = '-createdAt',
    fields = 'name email role',
    populate = null 
  } = options;

  let query = User.find({ active: true })
    .select(fields)
    .sort(sort)
    .limit(limit * 1)
    .skip((page - 1) * limit)
    .lean(); // Returns plain JavaScript objects instead of Mongoose documents

  if (populate) {
    query = query.populate(populate);
  }

  return await query;
};
```

## Testing Patterns

### Unit Test Structure
```javascript
// Test file structure
describe('UserService', () => {
  beforeEach(() => {
    // Setup before each test
  });

  afterEach(() => {
    // Cleanup after each test
  });

  describe('createUser', () => {
    it('should create a user with valid data', async () => {
      // Arrange
      const userData = { name: 'John', email: 'john@example.com' };
      
      // Act
      const result = await userService.createUser(userData);
      
      // Assert
      expect(result).toHaveProperty('id');
      expect(result.name).toBe(userData.name);
    });

    it('should throw error for duplicate email', async () => {
      // Arrange
      const userData = { name: 'John', email: 'existing@example.com' };
      
      // Act & Assert
      await expect(userService.createUser(userData))
        .rejects
        .toThrow('Email already exists');
    });
  });
});
```

### Integration Test Patterns
```javascript
// API integration tests
describe('Users API', () => {
  let app;

  beforeAll(async () => {
    app = require('../server');
    await connectTestDatabase();
  });

  afterAll(async () => {
    await cleanupTestDatabase();
  });

  beforeEach(async () => {
    await seedTestData();
  });

  it('should create a new user', async () => {
    const userData = {
      name: 'Test User',
      email: 'test@example.com',
      password: 'Password123'
    };

    const response = await request(app)
      .post('/api/users')
      .send(userData)
      .expect(201);

    expect(response.body.success).toBe(true);
    expect(response.body.data).toHaveProperty('id');
    expect(response.body.data.email).toBe(userData.email);
  });
});
```

## Security Patterns

### Security Best Practices
```javascript
// Security middleware setup
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');

// Basic security headers
app.use(helmet({
  contentSecurityPolicy: {
    directives: {
      defaultSrc: ["'self'"],
      styleSrc: ["'self'", "'unsafe-inline'"],
      scriptSrc: ["'self'"],
      imgSrc: ["'self'", "data:", "https:"]
    }
  }
}));

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // limit each IP to 100 requests per windowMs
  message: {
    success: false,
    message: 'Too many requests from this IP, please try again later.'
  }
});

// Input sanitization
const mongoSanitize = require('express-mongo-sanitize');
const xss = require('xss-clean');

app.use(mongoSanitize()); // Remove NoSQL injection
app.use(xss()); // Clean user input from malicious HTML
```

This comprehensive patterns and conventions guide ensures consistent, maintainable, and secure code throughout the {{project}} application.