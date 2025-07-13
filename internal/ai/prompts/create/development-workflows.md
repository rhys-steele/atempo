You are an expert DevOps engineer and software architect. Create development workflows specifically for building this application:

**PROJECT FOCUS: {{PROJECT_DESCRIPTION}}**

Framework: {{FRAMEWORK}}

IMPORTANT: These workflows must be tailored to developing "{{PROJECT_DESCRIPTION}}". Include specific commands, processes, and practices needed for this particular application, not generic {{FRAMEWORK}} workflows.

**CRITICAL: Testing Directory Convention**
- ALL testing and experimental projects MUST be created in the `/testing` directory
- Use: `atempo create {{FRAMEWORK}} testing/project-name` for test projects
- Never create test projects in the root directory
- This keeps the workspace clean and organized

Please create a detailed development workflows document in Markdown format that includes:

1. **Project Setup for {{PROJECT_DESCRIPTION}}** - How to set up this specific application
2. **Development Process for Core Features** - Workflow for building "{{PROJECT_DESCRIPTION}}" features
3. **Database Development Workflow** - Managing schemas and data for this use case
4. **API Development Process** - Building and testing endpoints for "{{PROJECT_DESCRIPTION}}"
5. **Feature-Specific Testing** - Testing strategies for this application's functionality
6. **Local Development Commands** - Specific commands for working on "{{PROJECT_DESCRIPTION}}"
7. **Debugging This Application** - Common issues and solutions for this use case
8. **Code Organization Practices** - How to structure code for "{{PROJECT_DESCRIPTION}}"
9. **Performance Considerations** - Optimization specific to this application type
10. **Deployment & Environment Setup** - Deploying "{{PROJECT_DESCRIPTION}}" to production

For each workflow, include:
- Actual commands and examples for this specific project
- Step-by-step processes for implementing core features
- Specific testing approaches for this use case
- Real examples of development tasks (e.g., "Adding a new task endpoint")
- Configuration specific to this application's needs
- **Always use `testing/` directory for test projects and examples**

Example sections:
- "How to add a new API endpoint for task management"
- "Testing user authentication flows"  
- "Setting up real-time notifications"
- "Creating test projects with: `atempo create {{FRAMEWORK}} testing/my-test-app`"

This should be a practical handbook for developers working specifically on "{{PROJECT_DESCRIPTION}}", not a generic {{FRAMEWORK}} development guide.

**Remember**: Always demonstrate commands using the `testing/` directory for any test or example projects!