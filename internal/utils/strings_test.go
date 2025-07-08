package utils

import "testing"

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{"item exists", []string{"apple", "banana", "cherry"}, "banana", true},
		{"item does not exist", []string{"apple", "banana", "cherry"}, "grape", false},
		{"empty slice", []string{}, "apple", false},
		{"empty item", []string{"apple", "banana", ""}, "", true},
		{"item is empty string in non-empty slice", []string{"apple", "banana"}, "", false},
		{"case sensitive", []string{"Apple", "Banana"}, "apple", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Contains(test.slice, test.item)
			if result != test.expected {
				t.Errorf("Contains(%v, %q) = %t, want %t", test.slice, test.item, result, test.expected)
			}
		})
	}
}

func TestContainsFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flag     string
		expected bool
	}{
		{"flag exists", []string{"--verbose", "--build", "--force"}, "--build", true},
		{"flag does not exist", []string{"--verbose", "--force"}, "--build", false},
		{"empty args", []string{}, "--build", false},
		{"partial match should not work", []string{"--building"}, "--build", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ContainsFlag(test.args, test.flag)
			if result != test.expected {
				t.Errorf("ContainsFlag(%v, %q) = %t, want %t", test.args, test.flag, result, test.expected)
			}
		})
	}
}
