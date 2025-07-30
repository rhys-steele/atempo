# {{project}} - Express.js API Project Overview

## Project Mission
{{project}} is a modern Express.js application designed for building robust REST APIs and web services. Built with Atempo, it provides a solid foundation for rapid API development with best practices, security, and scalability in mind.

## Core Value Proposition
- **Performance-First**: Built on Node.js with Express.js for high-performance APIs
- **Security-Ready**: Includes helmet, rate limiting, and security best practices
- **Developer Experience**: Hot reloading, comprehensive logging, and testing setup
- **Container-Native**: Docker-first development with multi-service orchestration
- **Production-Ready**: Environment configuration, error handling, and monitoring

## Architecture Overview

### Express.js Architecture
The project follows Express.js middleware-based architecture with clear separation of concerns:

```
src/
├── server.js                 # Application entry point
├── routes/                   # Route definitions and routing logic
│   ├── index.js             # Main route aggregator
│   ├── api/                 # API route modules
│   └── health.js            # Health check endpoints
├── controllers/             # Request/response handling logic
├── middleware/              # Custom middleware functions
├── services/                # Business logic layer
├── models/                  # Data models and database schemas
├── utils/                   # Helper functions and utilities
├── config/                  # Configuration management
└── tests/                   # Test files and test utilities
```

### Key Components

#### 1. Middleware System
- **Security Middleware**: Helmet for security headers, CORS configuration
- **Request Processing**: Body parsing, compression, rate limiting
- **Logging**: Morgan for HTTP request logging
- **Error Handling**: Global error handling middleware

#### 2. Route Organization
- **Modular Routes**: Express Router for organized route management
- **RESTful Design**: Following REST conventions for API endpoints
- **Middleware Integration**: Route-specific middleware for authentication, validation
- **API Versioning**: Support for API versioning strategies

#### 3. Database Integration
- **MongoDB**: Primary database with Mongoose ODM
- **Redis**: Caching and session storage
- **Connection Management**: Robust connection handling and retry logic
- **Schema Validation**: Mongoose schemas with validation rules

#### 4. Development Tools
- **Nodemon**: Auto-restart during development
- **ESLint + Prettier**: Code quality and formatting
- **Jest**: Testing framework with supertest for API testing
- **Docker**: Containerized development environment

#### 5. Production Features
- **Environment Configuration**: Comprehensive .env management
- **Health Checks**: Built-in health check endpoints
- **Error Handling**: Structured error responses and logging
- **Performance**: Compression and optimization middleware

## Current State

### ✅ Completed Features
- Express.js server setup with security middleware
- Docker containerization with multi-service stack
- Environment configuration management
- Basic routing and middleware structure
- Health check endpoints
- Development tooling (ESLint, Prettier, Jest)
- Redis and MongoDB integration setup

### 🔄 Ready for Development
- Custom route implementation
- Database model definitions
- Authentication and authorization
- API endpoint development
- Business logic implementation
- Comprehensive testing

### 📋 Common Extensions
- Authentication system (JWT, Passport)
- File upload handling (Multer)
- Real-time features (Socket.io)
- GraphQL integration (Apollo Server)
- API documentation (Swagger/OpenAPI)
- Monitoring and logging (Winston, Morgan)

## Technical Highlights

### Node.js + Express Implementation
- **Minimal Dependencies**: Core Express.js with essential middleware
- **Async/Await**: Modern asynchronous JavaScript patterns
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Performance**: Optimized with compression and caching strategies

### Container Strategy
- **Multi-Service**: Express app, MongoDB, Redis, MailHog
- **Development**: Hot reloading with volume mounting
- **Health Monitoring**: Container health checks and service dependencies
- **Networking**: Internal service communication with Docker networks

### API Design
- **RESTful**: Following REST conventions and HTTP semantics
- **JSON-First**: Structured JSON responses with consistent formatting
- **Validation**: Input validation and sanitization
- **Documentation**: API documentation and client integration examples

## Development Workflows

### Getting Started
```bash
# Start development environment
docker-compose up -d

# View logs
docker-compose logs -f app

# Access application
curl http://localhost:3000/health

# Run tests
docker-compose exec app npm test
```

### Development Commands
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test

# Lint and format code
npm run lint
npm run format
```

### API Testing
```bash
# Health check
curl http://localhost:3000/health

# API endpoint
curl http://localhost:3000/api

# Example user endpoint
curl http://localhost:3000/api/users
```

## Key Files for AI Understanding

### Entry Points
- `src/server.js`: Application entry point and middleware setup
- `src/routes/index.js`: Main route configuration

### Core Logic
- `src/controllers/`: Request/response handling
- `src/services/`: Business logic implementation
- `src/models/`: Data models and schemas
- `src/middleware/`: Custom middleware functions

### Configuration
- `.env`: Environment variables and configuration
- `package.json`: Dependencies and scripts
- `docker-compose.yml`: Multi-service container setup

### Development
- `tests/`: Test files and utilities
- `.eslintrc.js`: Code quality configuration
- `.prettierrc`: Code formatting rules

This overview provides the foundational understanding needed to work effectively with the {{project}} Express.js application, its architecture, and development practices.