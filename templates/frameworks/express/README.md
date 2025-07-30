# {{project}} - Express.js API Project

> A modern Express.js REST API built with Atempo, featuring Docker containerization, MongoDB/Redis integration, and comprehensive development tooling.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local development)

### Start Development Environment
```bash
# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View application logs
docker-compose logs -f app

# Test the API
curl http://localhost:3000/health
```

### Access Your Services
- **API Server**: http://localhost:3000
- **API Documentation**: http://localhost:3000/api
- **Health Check**: http://localhost:3000/health
- **MailHog (Email Testing)**: http://localhost:8025
- **MongoDB**: mongodb://localhost:27017/{{project}}
- **Redis**: redis://localhost:6379

## 📁 Project Structure

```
{{project}}/
├── src/                        # Application source code
│   ├── server.js              # Main application entry point
│   ├── routes/                # API route definitions
│   ├── controllers/           # Request/response handling
│   ├── services/              # Business logic layer
│   ├── models/                # Database models (Mongoose)
│   ├── middleware/            # Custom middleware
│   ├── utils/                 # Helper functions
│   ├── config/                # Configuration files
│   └── tests/                 # Test files
├── infra/                     # Infrastructure configuration
│   └── docker/                # Docker setup files
├── .ai/                       # AI context for development
└── package.json               # Node.js dependencies
```

## 🛠️ Development Commands

### Application Commands
```bash
# Start development server (with hot reload)
npm run dev

# Start production server
npm start

# Install dependencies
npm install

# Run tests
npm test

# Run tests in watch mode
npm run test:watch

# Lint code
npm run lint

# Format code
npm run format
```

### Docker Commands
```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f app

# Restart a service
docker-compose restart app

# Execute commands in container
docker-compose exec app npm test
docker-compose exec app npm run lint

# Access container shell
docker-compose exec app sh
```

### Database Commands
```bash
# Connect to MongoDB
docker-compose exec mongo mongosh {{project}}

# Connect to Redis CLI
docker-compose exec redis redis-cli

# View MongoDB collections
docker-compose exec mongo mongosh --eval "use {{project}}; show collections"

# View Redis keys
docker-compose exec redis redis-cli keys "*"
```

## 🌐 API Endpoints

### Core Endpoints
```bash
# Health check
GET /health

# API information
GET /api

# Example user endpoints
GET    /api/users       # List users
POST   /api/users       # Create user
GET    /api/users/:id   # Get user by ID
PUT    /api/users/:id   # Update user
DELETE /api/users/:id   # Delete user
```

### API Response Format
All API responses follow a consistent format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {},
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### Error Response Format
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Email is required",
      "code": "REQUIRED_FIELD"
    }
  ],
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

## 🔧 Configuration

### Environment Variables
Copy `.env.example` to `.env` and configure:

```bash
# Server Configuration
NODE_ENV=development
PORT=3000

# Database Configuration
MONGODB_URL=mongodb://mongo:27017/{{project}}
REDIS_URL=redis://redis:6379

# Authentication
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=24h

# Email Configuration (MailHog for development)
SMTP_HOST=mailhog
SMTP_PORT=1025
```

### Service Configuration
The application includes these services:
- **Express App**: Main API server on port 3000
- **MongoDB**: Document database on port 27017
- **Redis**: Cache and session store on port 6379
- **MailHog**: Email testing on ports 1025 (SMTP) and 8025 (Web UI)

## 🧪 Testing

### Running Tests
```bash
# Run all tests
npm test

# Run specific test file
npm test -- src/tests/unit/userService.test.js

# Run integration tests
npm test -- src/tests/integration/

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

### API Testing Examples
```bash
# Test health endpoint
curl http://localhost:3000/health

# Test API endpoint
curl http://localhost:3000/api

# Test user creation
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Test with authentication (if implemented)
curl -X GET http://localhost:3000/api/users \
  -H "Authorization: Bearer <jwt-token>"
```

## 🔒 Security Features

### Built-in Security
- **Helmet**: Security headers
- **CORS**: Cross-origin resource sharing configuration
- **Rate Limiting**: API request rate limiting
- **Input Validation**: Request data validation
- **Environment Variables**: Secure configuration management

### Security Best Practices
- Use environment variables for sensitive data
- Implement proper authentication and authorization
- Validate all input data
- Use HTTPS in production
- Keep dependencies updated

## 🚀 Deployment

### Production Build
```bash
# Build production image
docker build -t {{project}} .

# Run production container
docker run -d \
  --name {{project}} \
  -p 3000:3000 \
  -e NODE_ENV=production \
  -e MONGODB_URL="mongodb://your-mongodb-url" \
  {{project}}
```

### Environment Setup
1. Set `NODE_ENV=production`
2. Configure production database URLs
3. Set secure JWT secret
4. Configure SMTP settings for email
5. Set up monitoring and logging

## 📚 Documentation

### API Documentation
- Interactive API documentation available at `/api`
- Health check endpoint at `/health`
- Comprehensive error handling with meaningful messages

### Development Resources
- **AI Context**: `.ai/` directory contains comprehensive development context
- **Code Patterns**: See `.ai/patterns-and-conventions.md`
- **Development Workflows**: See `.ai/development-workflows.md`
- **Architecture Guide**: See `.ai/codebase-map.md`

## 🔧 Troubleshooting

### Common Issues

**Port Already in Use**
```bash
lsof -i :3000
kill -9 <PID>
```

**Docker Issues**
```bash
docker-compose down -v
docker-compose up -d --build
```

**Node Modules Issues**
```bash
rm -rf node_modules package-lock.json
npm install
```

**Database Connection Issues**
```bash
# Check MongoDB
docker-compose exec mongo mongosh --eval "db.runCommand({ping: 1})"

# Check Redis
docker-compose exec redis redis-cli ping
```

### Health Checks
```bash
# Check all services
docker-compose ps

# Test API health
curl -f http://localhost:3000/health || echo "API is down"

# Check container logs
docker-compose logs app
```

## 🤝 Contributing

### Development Workflow
1. Create feature branch
2. Write tests for new functionality
3. Implement feature following existing patterns
4. Run tests and linting
5. Submit pull request

### Code Standards
- Follow ESLint configuration
- Use Prettier for code formatting
- Write comprehensive tests
- Follow RESTful API conventions
- Document new endpoints

## 📝 License

This project is licensed under the MIT License.

---

**Built with Atempo** - AI-enhanced project scaffolding for modern development.