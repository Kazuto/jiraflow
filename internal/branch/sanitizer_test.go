package branch

import (
	"testing"
)

func TestBranchSanitizer_Sanitize(t *testing.T) {
	sanitizer := NewBranchSanitizer()

	tests := []struct {
		name     string
		input    string
		options  SanitizationOptions
		expected string
	}{
		{
			name:  "basic sanitization",
			input: "Add user authentication",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "add-user-authentication",
		},
		{
			name:  "remove quotes and parentheses",
			input: `Fix "login" (issue) with: special chars`,
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "fix-login-issue-with-special-chars",
		},
		{
			name:  "replace space-hyphen-space",
			input: "Update - user profile - settings",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "update-user-profile-settings",
		},
		{
			name:  "remove double hyphens",
			input: "Fix--multiple--hyphens",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "fix-multiple-hyphens",
		},
		{
			name:  "length truncation",
			input: "This is a very long title that should be truncated",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     20,
			},
			expected: "this-is-a-very-long",
		},
		{
			name:  "remove trailing separator",
			input: "Title with trailing space ",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "title-with-trailing-space",
		},
		{
			name:  "custom separator",
			input: "Use custom separator",
			options: SanitizationOptions{
				Separator:     "_",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "use_custom_separator",
		},
		{
			name:  "remove umlauts",
			input: "Füge Benutzerverwaltung hinzü ß test",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: true,
				MaxLength:     0,
			},
			expected: "fuege-benutzerverwaltung-hinzue-ss-test",
		},
		{
			name:  "preserve case",
			input: "Keep Original Case",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     false,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "Keep-Original-Case",
		},
		{
			name:  "empty input",
			input: "",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.Sanitize(tt.input, tt.options)
			if result != tt.expected {
				t.Errorf("Sanitize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBranchSanitizer_removeUmlauts(t *testing.T) {
	sanitizer := NewBranchSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase umlauts",
			input:    "äöüß",
			expected: "aeoeuess",
		},
		{
			name:     "uppercase umlauts",
			input:    "ÄÖÜ",
			expected: "AeOeUe",
		},
		{
			name:     "mixed case with text",
			input:    "Müller Straße",
			expected: "Mueller Strasse",
		},
		{
			name:     "no umlauts",
			input:    "regular text",
			expected: "regular text",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.removeUmlauts(tt.input)
			if result != tt.expected {
				t.Errorf("removeUmlauts() = %v, want %v", result, tt.expected)
			}
		})
	}
}