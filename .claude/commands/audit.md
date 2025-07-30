# Comprehensive Codebase Audit

Run a thorough codebase audit on all source files within this laravel project. This audit should ensure the project follows industry best practices, framework-specific conventions, and modern development standards.

## Audit Areas

### 1. Code Quality & Standards
- **Linting and Code Style**: Check adherence to php and laravel coding standards
- **Code Documentation**: Evaluate inline comments, docstrings, and README documentation
- **Code Organization**: Assess file structure, naming conventions, and module organization
- **Error Handling**: Review error handling patterns and exception management

### 2. Architecture & Design
- **Project Structure**: Evaluate directory organization and architectural patterns
- **Separation of Concerns**: Check for proper separation between business logic, data access, and presentation layers
- **Design Patterns**: Identify and evaluate use of appropriate design patterns
- **Dependency Management**: Review dependency injection and service organization

### 3. Security & Best Practices
- **Security Vulnerabilities**: Scan for common security issues and vulnerabilities
- **Configuration Management**: Check for proper handling of secrets and environment variables
- **Input Validation**: Review data validation and sanitization practices
- **Authentication & Authorization**: Evaluate security implementation patterns

### 4. Performance & Optimization
- **Performance Patterns**: Identify potential performance bottlenecks
- **Database Queries**: Review query optimization and N+1 query problems
- **Caching Strategy**: Evaluate caching implementation and opportunities
- **Resource Management**: Check for proper resource cleanup and memory management

### 5. Testing & Quality Assurance
- **Test Coverage**: Assess current test coverage and identify gaps
- **Test Quality**: Review test structure, naming, and maintainability
- **Testing Patterns**: Evaluate unit, integration, and end-to-end testing strategies
- **Continuous Integration**: Check CI/CD pipeline configuration

### 6. Framework-Specific Standards
- **laravel Conventions**: Ensure adherence to laravel best practices
- **Framework Features**: Evaluate proper use of framework-specific features
- **Package Management**: Review dependency management and versioning
- **Configuration**: Check framework configuration and environment setup

### 7. Maintainability & Documentation
- **Code Readability**: Assess code clarity and maintainability
- **Documentation Quality**: Review API documentation, setup guides, and development docs
- **Changelog & Versioning**: Check version management and change documentation
- **Development Environment**: Evaluate development setup and tooling

## Deliverables

Create a comprehensive audit report in `.me/my-app-codebase-audit-YYYY-MM-DD.md` that includes:

### Executive Summary
- Overall project health score
- Critical issues requiring immediate attention
- Key strengths and areas of excellence
- Recommended next steps

### Detailed Findings
For each audit area, provide:
- Current state assessment
- Specific issues identified with file/line references
- Best practice recommendations
- Priority level (Critical, High, Medium, Low)

### Action Plan
Break down remediation tasks into phases:

#### Phase 1: Critical Issues (Week 1)
- Security vulnerabilities
- Breaking changes
- Critical performance issues

#### Phase 2: High Priority Improvements (Week 2-3)
- Architecture improvements
- Test coverage gaps
- Documentation updates

#### Phase 3: Medium Priority Enhancements (Week 4-6)
- Code quality improvements
- Performance optimizations
- Framework best practices

#### Phase 4: Long-term Improvements (Ongoing)
- Advanced features
- Developer experience improvements
- Process optimizations

### Task Breakdown
For each phase, provide:
- Clear, actionable tasks
- Estimated time requirements
- Prerequisites and dependencies
- Success criteria
- Suitable for both AI agents and human developers

## Implementation Notes

- Focus on practical, actionable recommendations
- Provide specific code examples where applicable
- Include links to relevant documentation and resources
- Prioritize issues based on impact and effort required
- Consider the project's current maturity and development stage

Begin the audit now, examining all source files systematically.