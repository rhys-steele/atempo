You are an expert software architect specializing in {{FRAMEWORK}}. Create patterns and conventions specifically for building this application:

**PROJECT FOCUS: {{PROJECT_DESCRIPTION}}**

Framework: {{FRAMEWORK}}

IMPORTANT: These patterns and conventions must be specifically designed for "{{PROJECT_DESCRIPTION}}". Focus on architectural decisions, naming conventions, and code patterns that make sense for this particular use case, not generic {{FRAMEWORK}} patterns.

Please create a detailed patterns and conventions document in Markdown format that includes:

1. **Application-Specific Architecture** - Patterns chosen specifically for "{{PROJECT_DESCRIPTION}}"
2. **Domain-Driven Code Organization** - How to structure code around this application's domain
3. **Naming Conventions for This Use Case** - Model names, method names specific to "{{PROJECT_DESCRIPTION}}"
4. **Business Logic Patterns** - How to implement core features of "{{PROJECT_DESCRIPTION}}"
5. **Data Flow for This Application** - How data moves through "{{PROJECT_DESCRIPTION}}"
6. **Error Handling for This Domain** - Error patterns specific to this use case
7. **Security Patterns for This Application** - Security considerations for "{{PROJECT_DESCRIPTION}}"
8. **Performance Patterns for This Use Case** - Optimization specific to this application type
9. **Testing Patterns for This Domain** - How to test "{{PROJECT_DESCRIPTION}}" features
10. **API Design for This Application** - Endpoint patterns and data structures for this use case
11. **Database Patterns for This Domain** - Data modeling specific to "{{PROJECT_DESCRIPTION}}"
12. **Integration Patterns** - How to integrate external services for this use case

For each pattern, provide:
- Concrete examples using domain concepts from "{{PROJECT_DESCRIPTION}}"
- Actual class names, method names, and file structures for this use case
- Specific {{FRAMEWORK}} implementations for this application's needs
- Real code examples showing implementation
- Rationale for why this pattern fits "{{PROJECT_DESCRIPTION}}"

Example: If this is a task management API, include patterns like:
- TaskController, UserService, TaskRepository naming
- Task assignment and notification patterns
- Authorization patterns for task access
- Data validation patterns for task creation

This should be an architectural guide specifically for building "{{PROJECT_DESCRIPTION}}" using {{FRAMEWORK}}, not generic framework documentation.