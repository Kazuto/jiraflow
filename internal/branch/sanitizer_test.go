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
		// Enhanced test cases for comprehensive sanitization
		{
			name:  "comprehensive special character removal",
			input: `Fix [bug] {urgent} <critical> /path\to\file | pipe & ampersand * wildcard`,
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "fix-bug-urgent-critical-path-to-file-pipe-ampersand-wildcard",
		},
		{
			name:  "multiple consecutive separators cleanup",
			input: "Fix---multiple----separators-----here",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "fix-multiple-separators-here",
		},
		{
			name:  "mixed whitespace handling",
			input: "Handle\ttabs\nand\r\nline\nbreaks   and   spaces",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "handle-tabs-and-line-breaks-and-spaces",
		},
		{
			name:  "leading and trailing separators",
			input: "---leading and trailing---",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "leading-and-trailing",
		},
		{
			name:  "dots and version numbers preserved",
			input: "Update to version 2.1.3 release",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "update-to-version-2.1.3-release",
		},
		{
			name:  "leading dots removed (Git safety)",
			input: "...hidden file update",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "hidden-file-update",
		},
		{
			name:  "smart length truncation at separator",
			input: "This is a very long title that should be truncated intelligently",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     30,
			},
			expected: "this-is-a-very-long-title",
		},
		{
			name:  "underscore separator with mixed patterns",
			input: "Fix - issue_with - mixed_separators",
			options: SanitizationOptions{
				Separator:     "_",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     0,
			},
			expected: "fix_issue_with_mixed_separators",
		},
		{
			name:  "extended European characters with umlauts enabled",
			input: "Café résumé naïve piñata",
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: true,
				MaxLength:     0,
			},
			expected: "cafe-resume-naive-pinata",
		},
		{
			name:  "complex real-world example",
			input: `[URGENT] Fix "user login" (OAuth 2.0) - handle special chars & edge cases!`,
			options: SanitizationOptions{
				Separator:     "-",
				Lowercase:     true,
				RemoveUmlauts: false,
				MaxLength:     50,
			},
			expected: "urgent-fix-user-login-oauth-2.0-handle-special",
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
		// Enhanced test cases for extended European characters
		{
			name:     "French accented characters",
			input:    "café résumé naïve",
			expected: "cafe resume naive",
		},
		{
			name:     "Spanish characters",
			input:    "niño piñata",
			expected: "nino pinata",
		},
		{
			name:     "mixed European characters",
			input:    "Åse Björk Çağlar",
			expected: "Ase Bjoerk Cağlar", // ğ is not in our mapping, so it stays
		},
		{
			name:     "comprehensive character test",
			input:    "àáâãäåèéêëìíîïòóôõöùúûüçñ",
			expected: "aaaaaeaeeeeiiiioooooeuuuuecn", // ø->o, ü->ue from German mappings
		},
		{
			name:     "uppercase European characters",
			input:    "ÀÁÂÃÄÅÈÉÊËÌÍÎÏÒÓÔÕÖÙÚÛÜÇÑ",
			expected: "AAAAAeAEEEEIIIIOOOOOeUUUUeCN",
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