You are an expert software architect. Create a comprehensive project overview document specifically for this project:

**PROJECT FOCUS: {{PROJECT_DESCRIPTION}}**

Framework: {{FRAMEWORK}}

IMPORTANT: Everything in this document must be specifically tailored to "{{PROJECT_DESCRIPTION}}". Do not create generic {{FRAMEWORK}} documentation. Instead, create project-specific documentation that explains how this {{FRAMEWORK}} application will implement the described functionality.

Please create a detailed project overview in Markdown format that includes:

1. **Project Mission & Goals** - Specifically what "{{PROJECT_DESCRIPTION}}" aims to achieve
2. **Core Value Proposition** - Why this specific application ({{PROJECT_DESCRIPTION}}) is valuable
3. **Architecture Overview** - How {{FRAMEWORK}} will be structured to support "{{PROJECT_DESCRIPTION}}"
4. **Key Components** - Specific modules/components needed for "{{PROJECT_DESCRIPTION}}"
5. **Feature Implementation Plan** - How the core features of "{{PROJECT_DESCRIPTION}}" will be built
6. **Technical Approach** - {{FRAMEWORK}}-specific decisions made for this use case
7. **Database Design** - Data models and relationships specific to "{{PROJECT_DESCRIPTION}}"
8. **API Design** - Endpoints and data flow specific to this application
9. **Integration Points** - External services needed for "{{PROJECT_DESCRIPTION}}"

Focus on the specific requirements and implementation details for "{{PROJECT_DESCRIPTION}}". Include concrete examples, specific API endpoints, database tables, and features that relate directly to the project description.

**CRITICAL: Testing Directory Convention**
When providing examples of commands or testing procedures, always use the `/testing` directory:
- Use: `atempo create {{FRAMEWORK}} testing/project-name` for test projects
- Example: `atempo create {{FRAMEWORK}} testing/my-task-api`
- Never suggest creating test projects in the root directory

This should read like documentation for the specific application described, not generic {{FRAMEWORK}} documentation.