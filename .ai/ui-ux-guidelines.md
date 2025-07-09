# UI/UX Guidelines for Atempo CLI

## Core Design Principles

### 1. Modern & Elegant Interface
- Clean, minimal output with purposeful whitespace
- Professional typography using consistent spacing
- Subdued color palette with meaningful color coding
- Focus on functionality over decoration

### 2. Emoji Usage Policy
**CRITICAL: Minimal emoji usage only**
- ✅ **Allowed**: Status indicators (✓, ✗, •)
- ✅ **Allowed**: Coloured text for status indicating text (eg. red, orange, green)
- ❌ **Forbidden**: Decorative emojis (🔧, 🧪, 📝, 🎉, etc.)
- ❌ **Forbidden**: Multiple emojis in sequence
- ❌ **Forbidden**: Emojis as section headers or bullets

### 3. Status Indicators
Use simple, clean symbols:
```
✓ Success (green)
✗ Error (red)  
• Info (default)
- Warning (yellow)
```

### 4. Output Structure
```
Command Name
────────────────────────────────
[clean, minimal content]

Section Name:
  Key: Value
  Status: running
  
Next Steps:
  Run: command here
```

### 5. Text Formatting
- **Bold** for commands and important actions
- *Italic* for file paths and variables
- `code` for literal commands
- Consistent indentation (2 spaces)

### 6. Professional Language
- Concise, direct communication
- No casual expressions or excessive punctuation
- Technical accuracy without jargon
- Action-oriented instructions

## Standardized CLI Output Formatting

### 7. Command Execution Indicators
Use the `⎿` symbol (U+23BF) for indicating command execution:
```
⎿ Running: docker-compose -f docker-compose.yml up -d (in /path/to/project, timeout: 3m0s)
```

**Format**: `⎿ Running: {command} (in {path}, timeout: {duration})`

### 8. Container Status Output - Minimal Approach
Container status lines should show ONLY the final state with clean formatting:

```
 my-app-2-redis            running
 my-app-2-mailhog          running
 my-app-2-mysql            running
 my-app-2-app              running
 my-app-2-webserver        running
```

**Format Standards**:
- **Minimal output**: Only show final container status (running, built, removed)
- **No intermediate messages**: Skip Creating/Created/Starting/Stopping/Stopped messages
- **Clean service names**: Remove "Container" prefix for brevity
- **Consistent alignment**: Service names padded to 25 characters
- **Status text in lowercase**: running, built, removed, etc.
- **No verbose warnings**: Filter out platform warnings and version warnings
- **Leading space for indentation**

### 9. Network and Volume Status - Filtered Out
Network and volume creation/removal messages are filtered out for minimal output. Only show when necessary for debugging.

### 10. Service Build Status
Service build status uses the same alignment:
```
 my-app-2-laravel-app          built
```

### 11. Filtered Content
The following are automatically filtered out for clean output:
- **Docker build steps**: Internal load, exporting layers, etc.
- **Intermediate container states**: Creating, Created, Starting, Stopping, Stopped
- **Network/Volume operations**: Creating, Created, Removing, Removed
- **Platform warnings**: Architecture mismatch warnings
- **Version warnings**: docker-compose.yml version obsolete warnings
- **Progress indicators**: [+] Running x/y status updates

### 12. Status Color Coding (Future Enhancement)
- **Green**: running, started, created, built
- **Red**: stopped, removed, error states
- **Yellow**: creating, starting, stopping, removing (transitional states)
- **Default**: unknown or neutral states

### 13. Indentation Standards
- Primary content: No indentation
- Secondary content: 2 spaces
- Container/service status: 1 space
- Nested items: 4 spaces

## Examples

### ❌ Bad (Too Many Emojis)
```
🔧 DNS Setup for Custom Domains  
───────────────────────────────────────────────────

📝 This will configure macOS to resolve .local domains
🧪 Testing DNS resolution...
⚠️  DNS test failed. 
```

### ✅ Good (Clean & Professional)  
```
DNS Configuration
────────────────

System will configure macOS resolver for .local domains.
This enables custom domains like 'project.local'.

Configure DNS resolver? [y/N]: 

Creating resolver configuration...
✓ DNS resolver configured
✗ DNS test failed - restart browser or wait
```

### ✅ Good (Docker Command Output) - Minimal & Clean
```
⎿ Running: docker-compose -f docker-compose.yml up -d (in /path/to/project, timeout: 3m0s)
 my-app-2-redis            running
 my-app-2-mailhog          running
 my-app-2-mysql            running
 my-app-2-app              running
 my-app-2-webserver        running
```

### ❌ Bad (Verbose & Inconsistent Formatting)
```
→ Running: docker-compose -f docker-compose.yml up -d (in /path/to/project, timeout: 3m0s)
WARN[0000] /path/docker-compose.yml: version obsolete
[+] Running 1/1
 ! app Warning pull access denied...
[+] Building 7.0s (19/19) FINISHED
 => [internal] load local bake definitions
 => => reading from stdin 430B
 Network my-app-2_default  Creating
 Network my-app-2_default  Created
 Container my-app-2-redis  Creating
 Container my-app-2-redis  Created
 Container my-app-2-redis  Starting
 Container my-app-2-redis  Started
 Container my-app-2-mailhog  Creating
 Container my-app-2-mailhog  Created
 Container my-app-2-mailhog  Starting
 Container my-app-2-mailhog  Started
```

## Implementation Notes
- All CLI commands should follow these guidelines
- Status output should be clean and scannable
- Error messages should be helpful and actionable
- Success states should be brief and confident
- Use consistent symbols and formatting throughout the application
- Apply color coding consistently based on status types