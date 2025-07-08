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

## Implementation Notes
- All CLI commands should follow these guidelines
- Status output should be clean and scannable
- Error messages should be helpful and actionable
- Success states should be brief and confident