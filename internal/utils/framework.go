package utils

// GetFrameworkLanguage returns the programming language for a given framework
func GetFrameworkLanguage(framework string) string {
	switch framework {
	case "laravel":
		return "php"
	case "django":
		return "python"
	case "nextjs", "react", "vue", "nuxt":
		return "javascript"
	case "fastapi":
		return "python"
	case "express":
		return "javascript"
	case "rails":
		return "ruby"
	case "spring":
		return "java"
	case "dotnet":
		return "csharp"
	default:
		return "unknown"
	}
}
