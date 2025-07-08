package utils

import "testing"

func TestGetFrameworkLanguage(t *testing.T) {
	tests := []struct {
		framework string
		expected  string
	}{
		{"laravel", "php"},
		{"django", "python"},
		{"nextjs", "javascript"},
		{"react", "javascript"},
		{"vue", "javascript"},
		{"nuxt", "javascript"},
		{"fastapi", "python"},
		{"express", "javascript"},
		{"rails", "ruby"},
		{"spring", "java"},
		{"dotnet", "csharp"},
		{"unknown-framework", "unknown"},
		{"", "unknown"},
	}

	for _, test := range tests {
		t.Run(test.framework, func(t *testing.T) {
			result := GetFrameworkLanguage(test.framework)
			if result != test.expected {
				t.Errorf("GetFrameworkLanguage(%q) = %q, want %q", test.framework, result, test.expected)
			}
		})
	}
}
