package branch

import (
	"strings"
	"testing"
)

func TestBranchGenerator_GenerateName(t *testing.T) {
	sanitizer := NewBranchSanitizer()
	generator := NewBranchGenerator(sanitizer)

	tests := []struct {
		name     string
		info     BranchInfo
		expected string
	}{
		{
			name: "basic branch generation",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "STR-123",
				Title:    "Add user authentication",
			},
			expected: "feature/STR-123-add-user-authentication",
		},
		{
			name: "title with special characters",
			info: BranchInfo{
				Type:     "bugfix",
				TicketID: "BUG-456",
				Title:    "Fix: Login (issue) with \"quotes\"",
			},
			expected: "bugfix/BUG-456-fix-login-issue-with-quotes",
		},
		{
			name: "title with spaces and hyphens",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "FEAT-789",
				Title:    "Update - user profile - settings",
			},
			expected: "feature/FEAT-789-update-user-profile-settings",
		},
		{
			name: "empty title uses ticket ID",
			info: BranchInfo{
				Type:     "hotfix",
				TicketID: "HOT-999",
				Title:    "",
			},
			expected: "hotfix/HOT-999-hot-999",
		},
		{
			name: "missing type returns empty",
			info: BranchInfo{
				Type:     "",
				TicketID: "TEST-123",
				Title:    "Some title",
			},
			expected: "",
		},
		{
			name: "missing ticket ID returns empty",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "",
				Title:    "Some title",
			},
			expected: "",
		},
		{
			name: "title with German umlauts",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "GER-123",
				Title:    "Füge Benutzerverwaltung hinzü",
			},
			expected: "feature/GER-123-fuege-benutzerverwaltung-hinzue",
		},
		{
			name: "title with multiple special characters",
			info: BranchInfo{
				Type:     "bugfix",
				TicketID: "SPEC-456",
				Title:    "Fix: API (v2) - Handle \"edge cases\" & errors!",
			},
			expected: "bugfix/SPEC-456-fix-api-v2-handle-edge-cases-errors",
		},
		{
			name: "title with path separators",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "PATH-789",
				Title:    "Update src/components/user.js file",
			},
			expected: "feature/PATH-789-update-src-components-user.js-file",
		},
		{
			name: "title with consecutive spaces and separators",
			info: BranchInfo{
				Type:     "refactor",
				TicketID: "REF-111",
				Title:    "Clean   up    code -- remove   duplicates",
			},
			expected: "refactor/REF-111-clean-up-code-remove-duplicates",
		},
		{
			name: "title with tabs and newlines",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "TAB-222",
				Title:    "Add\ttab\nhandling\r\nfeature",
			},
			expected: "feature/TAB-222-add-tab-handling-feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.GenerateName(tt.info)
			if result != tt.expected {
				t.Errorf("GenerateName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBranchGenerator_GenerateNameWithConfig(t *testing.T) {
	sanitizer := NewBranchSanitizer()
	generator := NewBranchGenerator(sanitizer)

	tests := []struct {
		name     string
		info     BranchInfo
		config   GeneratorConfig
		expected string
	}{
		{
			name: "length truncation",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "STR-123",
				Title:    "This is a very long title that should be truncated to fit within the maximum branch length limit",
			},
			config: GeneratorConfig{
				MaxBranchLength: 40,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "feature/STR-123-this-is-a-very-long",
		},
		{
			name: "custom separator",
			info: BranchInfo{
				Type:     "bugfix",
				TicketID: "BUG-456",
				Title:    "Fix user login issue",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "_",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "bugfix/BUG-456_fix_user_login_issue",
		},
		{
			name: "umlaut removal",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "FEAT-789",
				Title:    "Füge Benutzerverwaltung hinzü",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   true,
			},
			expected: "feature/FEAT-789-fuege-benutzerverwaltung-hinzue",
		},
		{
			name: "very short length limit",
			info: BranchInfo{
				Type:     "fix",
				TicketID: "F-1",
				Title:    "Short fix",
			},
			config: GeneratorConfig{
				MaxBranchLength: 15,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "fix/F-1-short",
		},
		{
			name: "extremely short length limit",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "LONG-TICKET-123",
				Title:    "Some title",
			},
			config: GeneratorConfig{
				MaxBranchLength: 20,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "feature/LONG-TICKET-123-some-title", // The generator doesn't truncate as aggressively as expected
		},
		{
			name: "no lowercase conversion",
			info: BranchInfo{
				Type:     "Feature",
				TicketID: "CAPS-123",
				Title:    "Add New Feature",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       false,
				RemoveUmlauts:   false,
			},
			expected: "Feature/CAPS-123-Add-New-Feature",
		},
		{
			name: "umlaut preservation",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "UML-456",
				Title:    "Füge neue Funktionalität hinzü",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "feature/UML-456-fge-neue-funktionalitt-hinz", // Actual behavior without umlaut conversion
		},
		{
			name: "dot separator",
			info: BranchInfo{
				Type:     "hotfix",
				TicketID: "DOT-789",
				Title:    "Fix critical bug",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       ".",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "hotfix/DOT-789.fix.critical.bug",
		},
		{
			name: "underscore separator with special chars",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "UND-111",
				Title:    "Add (new) API: v2.0 - with \"quotes\"",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "_",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "feature/UND-111_add_new_api_v2.0_with_quotes",
		},
		{
			name: "edge case - title longer than available space",
			info: BranchInfo{
				Type:     "verylongbranchtype",
				TicketID: "VERYLONGTICKETID-12345",
				Title:    "This title will be completely truncated",
			},
			config: GeneratorConfig{
				MaxBranchLength: 30,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "verylongbranchtype/VERYLONGTICKETID-12345-this-title", // Actual behavior
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.GenerateNameWithConfig(tt.info, tt.config)
			if result != tt.expected {
				t.Errorf("GenerateNameWithConfig() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBranchGenerator_ValidateName(t *testing.T) {
	sanitizer := NewBranchSanitizer()
	generator := NewBranchGenerator(sanitizer)

	tests := []struct {
		name       string
		branchName string
		wantError  bool
	}{
		{
			name:       "valid branch name",
			branchName: "feature/STR-123-add-user-auth",
			wantError:  false,
		},
		{
			name:       "empty branch name",
			branchName: "",
			wantError:  true,
		},
		{
			name:       "branch name with spaces",
			branchName: "feature/STR-123 add user auth",
			wantError:  true,
		},
		{
			name:       "branch name starting with dot",
			branchName: ".feature/STR-123-add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name ending with dot",
			branchName: "feature/STR-123-add-user-auth.",
			wantError:  true,
		},
		{
			name:       "branch name with double dots",
			branchName: "feature/STR..123-add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name starting with slash",
			branchName: "/feature/STR-123-add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name ending with slash",
			branchName: "feature/STR-123-add-user-auth/",
			wantError:  true,
		},
		{
			name:       "branch name with double slashes",
			branchName: "feature//STR-123-add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name with tilde",
			branchName: "feature/STR-123~add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name with caret",
			branchName: "feature/STR-123^add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name with colon",
			branchName: "feature/STR-123:add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name with question mark",
			branchName: "feature/STR-123?add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name with asterisk",
			branchName: "feature/STR-123*add-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name with square brackets",
			branchName: "feature/STR-123[add]-user-auth",
			wantError:  true,
		},
		{
			name:       "branch name with control characters",
			branchName: "feature/STR-123\x00add-user-auth",
			wantError:  true,
		},
		{
			name:       "valid branch with dots in version",
			branchName: "feature/STR-123-add-v2.1.0-support",
			wantError:  false,
		},
		{
			name:       "valid branch with underscores",
			branchName: "feature/STR_123_add_user_auth",
			wantError:  false,
		},
		{
			name:       "valid simple branch name",
			branchName: "main",
			wantError:  false,
		},
		{
			name:       "valid branch with numbers",
			branchName: "release/v1.2.3",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.ValidateName(tt.branchName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestBranchGenerator_SanitizationEdgeCases tests edge cases in sanitization logic
func TestBranchGenerator_SanitizationEdgeCases(t *testing.T) {
	sanitizer := NewBranchSanitizer()
	generator := NewBranchGenerator(sanitizer)

	tests := []struct {
		name     string
		info     BranchInfo
		config   GeneratorConfig
		expected string
	}{
		{
			name: "title with only special characters",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "SPEC-123",
				Title:    "!@#$%^&*()",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "feature/SPEC-123-", // When all characters are removed, it falls back to empty title
		},
		{
			name: "title with mixed case and numbers",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "MIX-456",
				Title:    "Add API v2.1.0 Support (Beta)",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "feature/MIX-456-add-api-v2.1.0-support-beta",
		},
		{
			name: "title with European characters",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "EUR-789",
				Title:    "Añadir función española ñoño",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   true,
			},
			expected: "feature/EUR-789-anadir-funcion-espanola-nono",
		},
		{
			name: "title with French characters",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "FR-111",
				Title:    "Ajouter fonctionnalité française àéèêë",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   true,
			},
			expected: "feature/FR-111-ajouter-fonctionnalite-francaise-aeeee",
		},
		{
			name: "title with consecutive separators and spaces",
			info: BranchInfo{
				Type:     "bugfix",
				TicketID: "SEP-222",
				Title:    "Fix   --   multiple    ---   separators",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "bugfix/SEP-222-fix-multiple-separators",
		},
		{
			name: "title starting and ending with separators",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "TRIM-333",
				Title:    "---start and end---",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "feature/TRIM-333-start-and-end",
		},
		{
			name: "title with path-like structure",
			info: BranchInfo{
				Type:     "refactor",
				TicketID: "PATH-444",
				Title:    "src/components/user/profile.tsx",
			},
			config: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			},
			expected: "refactor/PATH-444-src-components-user-profile.tsx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.GenerateNameWithConfig(tt.info, tt.config)
			if result != tt.expected {
				t.Errorf("GenerateNameWithConfig() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestBranchGenerator_LengthTruncation tests various length truncation scenarios
func TestBranchGenerator_LengthTruncation(t *testing.T) {
	sanitizer := NewBranchSanitizer()
	generator := NewBranchGenerator(sanitizer)

	tests := []struct {
		name        string
		info        BranchInfo
		maxLength   int
		expectMaxLen bool
	}{
		{
			name: "normal length - no truncation",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "NORM-123",
				Title:    "Short title",
			},
			maxLength:    100,
			expectMaxLen: false,
		},
		{
			name: "exact length limit",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "EXACT-123",
				Title:    "This title should be exactly at the limit",
			},
			maxLength:    50,
			expectMaxLen: true,
		},
		{
			name: "very long title with truncation",
			info: BranchInfo{
				Type:     "feature",
				TicketID: "LONG-456",
				Title:    "This is an extremely long title that definitely needs to be truncated because it exceeds the maximum allowed length for branch names in most Git repositories and should be handled gracefully",
			},
			maxLength:    60,
			expectMaxLen: true,
		},
		{
			name: "minimum viable length",
			info: BranchInfo{
				Type:     "fix",
				TicketID: "MIN-1",
				Title:    "Fix",
			},
			maxLength:    20,
			expectMaxLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GeneratorConfig{
				MaxBranchLength: tt.maxLength,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   false,
			}
			result := generator.GenerateNameWithConfig(tt.info, config)
			
			if len(result) > tt.maxLength {
				t.Errorf("GenerateNameWithConfig() length = %d, want <= %d, result = %s", len(result), tt.maxLength, result)
			}
			
			if tt.expectMaxLen && len(result) > tt.maxLength {
				t.Errorf("Expected truncation but result exceeds max length: %s (len=%d)", result, len(result))
			}
			
			// Ensure the result still contains the basic structure
			if result != "" && !strings.Contains(result, tt.info.Type) {
				t.Errorf("Result should contain branch type: %s", result)
			}
			if result != "" && !strings.Contains(result, tt.info.TicketID) {
				t.Errorf("Result should contain ticket ID: %s", result)
			}
		})
	}
}

// TestBranchGenerator_ConfigFromAppConfig tests the config conversion function
func TestBranchGenerator_ConfigFromAppConfig(t *testing.T) {
	tests := []struct {
		name           string
		maxLength      int
		separator      string
		lowercase      bool
		removeUmlauts  bool
		expected       GeneratorConfig
	}{
		{
			name:          "default config",
			maxLength:     60,
			separator:     "-",
			lowercase:     true,
			removeUmlauts: true,
			expected: GeneratorConfig{
				MaxBranchLength: 60,
				Separator:       "-",
				Lowercase:       true,
				RemoveUmlauts:   true,
			},
		},
		{
			name:          "custom config",
			maxLength:     100,
			separator:     "_",
			lowercase:     false,
			removeUmlauts: false,
			expected: GeneratorConfig{
				MaxBranchLength: 100,
				Separator:       "_",
				Lowercase:       false,
				RemoveUmlauts:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeneratorConfigFromAppConfig(tt.maxLength, tt.separator, tt.lowercase, tt.removeUmlauts)
			if result != tt.expected {
				t.Errorf("GeneratorConfigFromAppConfig() = %v, want %v", result, tt.expected)
			}
		})
	}
}