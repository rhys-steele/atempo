# {{project}} Development Workflows & Commands

## Quick Start Guide

### Prerequisites
- Docker and Docker Compose installed
- Node.js 18+ (for local development)
- Git for version control

### Initial Setup
```bash
# Clone or initialize project
cd {{project}}

# Start all services with Docker
docker-compose up -d

# View application logs
docker-compose logs -f app

# Check service status
docker-compose ps
```

### Development Environment Access
```bash
# Application health check
curl http://localhost:3000/health

# API endpoints
curl http://localhost:3000/api
curl http://localhost:3000/api/users

# MailHog web interface (email testing)
open http://localhost:8025

# MongoDB connection (if using MongoDB Compass)
mongodb://localhost:27017/{{project}}

# Redis CLI access
docker-compose exec redis redis-cli
```

## Core Development Workflows

### 1. API Development Workflow
```bash
# Create new API endpoint
mkdir -p src/routes/api
touch src/routes/api/products.js

# Create corresponding controller
touch src/controllers/productController.js

# Create service layer
touch src/services/productService.js

# Create data model
touch src/models/Product.js

# Write tests
touch src/tests/integration/products.test.js

# Test the new endpoint
npm test -- --grep "products"
```

### 2. Database Development
```bash
# Connect to MongoDB
docker-compose exec mongo mongosh {{project}}

# View collections
show collections

# Connect to Redis
docker-compose exec redis redis-cli

# View Redis keys
keys *
```

### 3. Testing Workflow
```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run specific test file
npm test -- src/tests/unit/userService.test.js

# Run integration tests only
npm test -- src/tests/integration/

# Generate test coverage
npm test -- --coverage
```

### 4. Code Quality Workflow
```bash
# Lint code
npm run lint

# Fix linting issues automatically
npm run lint:fix

# Format code
npm run format

# Check code formatting
npm run format:check
```

## Available Commands

### Development Commands
```bash
# Start development server (with nodemon)
npm run dev

# Start production server
npm start

# Install dependencies
npm install

# Update dependencies
npm update
```

### Docker Commands
```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# Restart specific service
docker-compose restart app

# View logs
docker-compose logs -f app

# Execute commands in container
docker-compose exec app npm test
docker-compose exec app npm run lint

# Build and restart after Dockerfile changes
docker-compose up -d --build app
```

### Database Commands
```bash
# MongoDB operations
docker-compose exec mongo mongosh {{project}}

# Backup MongoDB
docker-compose exec mongo mongodump --db {{project}} --out /backup

# Restore MongoDB
docker-compose exec mongo mongorestore --db {{project}} /backup/{{project}}

# Redis operations
docker-compose exec redis redis-cli
docker-compose exec redis redis-cli flushall
```

### Testing Commands
```bash
# Run all tests
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch

# Run specific test suite
npm test -- --testPathPattern=integration

# Run tests with verbose output
npm test -- --verbose
```

## API Development Patterns

### Creating New Endpoints
```javascript
// 1. Define route (src/routes/api/users.js)
const express = require('express');
const router = express.Router();
const { getUsers, createUser } = require('../../controllers/userController');

router.get('/', getUsers);
router.post('/', createUser);

module.exports = router;

// 2. Implement controller (src/controllers/userController.js)
const userService = require('../services/userService');

const getUsers = async (req, res) => {
  try {
    const users = await userService.getAllUsers();
    res.json({ success: true, data: users });
  } catch (error) {
    res.status(500).json({ success: false, error: error.message });
  }
};

// 3. Add business logic (src/services/userService.js)
const User = require('../models/User');

const getAllUsers = async () => {
  return await User.find().select('-password');
};
```

### Error Handling Pattern
```javascript
// Global error handler (src/middleware/errorHandler.js)
const errorHandler = (err, req, res, next) => {
  console.error('Error:', err);
  
  res.status(err.status || 500).json({
    success: false,
    error: process.env.NODE_ENV === 'production' 
      ? 'Internal server error' 
      : err.message,
    ...(process.env.NODE_ENV !== 'production' && { stack: err.stack })
  });
};
```

### Validation Pattern
```javascript
// Input validation (src/middleware/validation.js)
const { body, validationResult } = require('express-validator');

const validateUser = [
  body('email').isEmail().normalizeEmail(),
  body('name').trim().isLength({ min: 2, max: 50 }),
  (req, res, next) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ 
        success: false, 
        errors: errors.array() 
      });
    }
    next();
  }
];
```

## Testing Strategies

### Unit Testing
```javascript
// Example unit test (src/tests/unit/userService.test.js)
const userService = require('../../services/userService');
const User = require('../../models/User');

jest.mock('../../models/User');

describe('UserService', () => {
  test('getAllUsers should return all users', async () => {
    const mockUsers = [{ id: 1, name: 'John' }];
    User.find.mockResolvedValue(mockUsers);
    
    const result = await userService.getAllUsers();
    
    expect(result).toEqual(mockUsers);
    expect(User.find).toHaveBeenCalledWith();
  });
});
```

### Integration Testing
```javascript
// Example integration test (src/tests/integration/users.test.js)
const request = require('supertest');
const app = require('../../server');

describe('Users API', () => {
  test('GET /api/users should return users list', async () => {
    const response = await request(app)
      .get('/api/users')
      .expect(200);
      
    expect(response.body.success).toBe(true);
    expect(Array.isArray(response.body.data)).toBe(true);
  });
});
```

## Debugging & Troubleshooting

### Common Issues and Solutions

#### Port Already in Use
```bash
# Find process using port
lsof -i :3000

# Kill process
kill -9 <PID>

# Or use different port
PORT=3001 npm run dev
```

#### Node Modules Issues
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear npm cache
npm cache clean --force
```

#### Docker Issues
```bash
# Remove containers and volumes
docker-compose down -v

# Rebuild containers
docker-compose up -d --build

# Check container logs
docker-compose logs app
```

#### Database Connection Issues
```bash
# Check MongoDB connection
docker-compose exec mongo mongosh --eval "db.runCommand({ping: 1})"

# Check Redis connection
docker-compose exec redis redis-cli ping

# Verify environment variables
docker-compose exec app printenv | grep MONGODB_URL
```

### Debugging Techniques
```javascript
// Add debugging middleware
app.use((req, res, next) => {
  console.log(`${req.method} ${req.path}`, req.body);
  next();
});

// Use debugging in controllers
const debug = require('debug')('app:controller');
debug('Processing user request:', req.params);

// Environment-based logging
if (process.env.NODE_ENV === 'development') {
  console.log('Debug info:', data);
}
```

## Performance Optimization

### Code Optimization
```bash
# Profile application
npm install --save-dev clinic
clinic doctor -- node server.js

# Analyze bundle size
npm install --save-dev webpack-bundle-analyzer
```

### Database Optimization
```javascript
// MongoDB indexing
db.users.createIndex({ email: 1 }, { unique: true });
db.products.createIndex({ category: 1, price: 1 });

// Query optimization
User.find({ active: true }).select('name email').lean();
```

### Caching Strategies
```javascript
// Redis caching
const redis = require('./config/redis');

app.get('/api/users', async (req, res) => {
  const cacheKey = 'users:all';
  const cached = await redis.get(cacheKey);
  
  if (cached) {
    return res.json(JSON.parse(cached));
  }
  
  const users = await userService.getAllUsers();
  await redis.setex(cacheKey, 300, JSON.stringify(users));
  
  res.json(users);
});
```

This development workflow guide provides comprehensive coverage of development practices, testing strategies, and troubleshooting techniques for efficient {{project}} development.