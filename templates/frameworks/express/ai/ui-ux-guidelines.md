# {{project}} API Design & Response Guidelines

## Core API Design Principles

### 1. RESTful Design Standards
- **Resource-Based URLs**: Use nouns, not verbs (`/users` not `/getUsers`)
- **HTTP Methods**: Use appropriate HTTP verbs (GET, POST, PUT, DELETE)
- **Status Codes**: Return meaningful HTTP status codes
- **Stateless**: Each request contains all necessary information

### 2. Consistent Response Format
All API responses follow a standardized JSON structure:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {},
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### 3. Error Response Standards
Consistent error response format across all endpoints:

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

## HTTP Status Codes

### Success Responses
- **200 OK**: Successful GET, PUT, PATCH requests
- **201 Created**: Successful POST requests that create resources
- **204 No Content**: Successful DELETE requests

### Client Error Responses
- **400 Bad Request**: Invalid request data or parameters
- **401 Unauthorized**: Authentication required or invalid credentials
- **403 Forbidden**: Valid authentication but insufficient permissions
- **404 Not Found**: Resource does not exist
- **409 Conflict**: Resource conflict (duplicate email, etc.)
- **422 Unprocessable Entity**: Valid JSON but semantically incorrect
- **429 Too Many Requests**: Rate limit exceeded

### Server Error Responses
- **500 Internal Server Error**: Unexpected server error
- **502 Bad Gateway**: Invalid response from upstream server
- **503 Service Unavailable**: Server temporarily unavailable

## API Endpoint Conventions

### Resource Naming
```javascript
// ✅ Good: Plural nouns for collections
GET    /api/users
POST   /api/users
GET    /api/users/:id
PUT    /api/users/:id
DELETE /api/users/:id

// ✅ Good: Nested resources
GET    /api/users/:id/posts
POST   /api/users/:id/posts

// ✅ Good: Action endpoints when needed
POST   /api/users/:id/activate
POST   /api/auth/login
POST   /api/auth/refresh
```

### Query Parameters
```javascript
// Pagination
GET /api/users?page=1&limit=20

// Filtering
GET /api/users?role=admin&active=true

// Sorting
GET /api/users?sort=-createdAt,name

// Field selection
GET /api/users?fields=name,email,role

// Search
GET /api/users?search=john&searchFields=name,email
```

## Response Data Formats

### Single Resource Response
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "64f123abc456def789",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "createdAt": "2024-01-15T10:30:00.000Z",
    "updatedAt": "2024-01-15T10:30:00.000Z"
  },
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### Collection Response with Pagination
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "64f123abc456def789",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "pages": 8,
    "hasNext": true,
    "hasPrev": false
  },
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### Creation Response
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "64f123abc456def789",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "createdAt": "2024-01-15T10:30:00.000Z"
  },
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

## Error Response Examples

### Validation Error (400)
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Please provide a valid email address",
      "code": "INVALID_EMAIL",
      "value": "invalid-email"
    },
    {
      "field": "password",
      "message": "Password must be at least 8 characters",
      "code": "MIN_LENGTH",
      "value": "123"
    }
  ],
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### Authentication Error (401)
```json
{
  "success": false,
  "message": "Authentication required",
  "error": {
    "code": "UNAUTHORIZED",
    "details": "No valid authentication token provided"
  },
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### Not Found Error (404)
```json
{
  "success": false,
  "message": "User not found",
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "resource": "user",
    "id": "64f123abc456def789"
  },
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### Server Error (500)
```json
{
  "success": false,
  "message": "Internal server error",
  "error": {
    "code": "INTERNAL_ERROR",
    "requestId": "req_64f123abc456def789"
  },
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

## Request/Response Headers

### Required Headers
```javascript
// Request headers
{
  "Content-Type": "application/json",
  "Authorization": "Bearer <jwt-token>",
  "Accept": "application/json"
}

// Response headers
{
  "Content-Type": "application/json",
  "X-Request-ID": "req_64f123abc456def789",
  "X-Rate-Limit-Remaining": "99",
  "X-Rate-Limit-Reset": "1642248600"
}
```

## API Documentation Standards

### Endpoint Documentation Format
```javascript
/**
 * @route   GET /api/users/:id
 * @desc    Get user by ID
 * @access  Private (requires authentication)
 * @params  {string} id - User ID (MongoDB ObjectId)
 * @returns {object} User object
 * @throws  {400} Invalid user ID format
 * @throws  {401} Unauthorized - No valid token
 * @throws  {404} User not found
 * @throws  {500} Internal server error
 * 
 * @example
 * GET /api/users/64f123abc456def789
 * Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
 * 
 * Response:
 * {
 *   "success": true,
 *   "data": {
 *     "id": "64f123abc456def789",
 *     "name": "John Doe",
 *     "email": "john@example.com"
 *   }
 * }
 */
```

## Client Integration Examples

### JavaScript/Fetch Example
```javascript
// GET request
const getUser = async (id) => {
  try {
    const response = await fetch(`/api/users/${id}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    const result = await response.json();
    
    if (!result.success) {
      throw new Error(result.message);
    }
    
    return result.data;
  } catch (error) {
    console.error('Error fetching user:', error);
    throw error;
  }
};

// POST request
const createUser = async (userData) => {
  try {
    const response = await fetch('/api/users', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(userData)
    });
    
    const result = await response.json();
    
    if (!result.success) {
      if (result.errors) {
        // Handle validation errors
        result.errors.forEach(error => {
          console.error(`${error.field}: ${error.message}`);
        });
      }
      throw new Error(result.message);
    }
    
    return result.data;
  } catch (error) {
    console.error('Error creating user:', error);
    throw error;
  }
};
```

### cURL Examples
```bash
# GET request
curl -X GET \
  http://localhost:3000/api/users/64f123abc456def789 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  -H 'Content-Type: application/json'

# POST request
curl -X POST \
  http://localhost:3000/api/users \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123"
  }'

# GET with query parameters
curl -X GET \
  "http://localhost:3000/api/users?page=1&limit=10&role=admin" \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
```

## Health Check Endpoint

### Standard Health Check Response
```json
GET /health

{
  "status": "OK",
  "timestamp": "2024-01-15T10:30:00.000Z",
  "uptime": 3600.25,
  "environment": "development",
  "version": "1.0.0",
  "services": {
    "database": "connected",
    "redis": "connected",
    "external_api": "operational"
  }
}
```

## Performance Considerations

### Response Time Guidelines
- **API responses**: < 200ms for simple queries
- **Complex queries**: < 1000ms with proper indexing
- **File uploads**: Progress indicators for > 5MB files
- **Background jobs**: Return job ID for long-running tasks

### Caching Headers
```javascript
// Cache control headers
res.set({
  'Cache-Control': 'public, max-age=300', // 5 minutes
  'ETag': '"abc123def456"',
  'Last-Modified': new Date().toUTCString()
});
```

These API design guidelines ensure consistent, predictable, and developer-friendly APIs that follow industry best practices and provide excellent developer experience for {{project}} consumers.