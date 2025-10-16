package branch

import (
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
			expected: "feature/STR-123-this-is-a-very-long-titl",
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
		name      string
		branchName string
		wantError bool
	}{
		{
			name:      "valid branch name",
			branchName: "feature/STR-123-add-user-auth",
			wantError: false,
		},
		{
			name:      "empty branch name",
			branchName: "",
			wantError: true,
		},
		{
			name:      "branch name with spaces",
			branchName: "feature/STR-123 add user auth",
			wantError: true,
		},
		{
			name:      "branch name starting with dot",
			branchName: ".feature/STR-123-add-user-auth",
			wantError: true,
		},
		{
			name:      "branch name with double dots",
			branchName: "feature/STR..123-add-user-auth",
			wantError: true,
		},
		{
			name:      "branch name with special characters",
			branchName: "feature/STR-123~add-user-auth",
			wantError: true,
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