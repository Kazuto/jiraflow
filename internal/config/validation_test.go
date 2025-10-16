package config

import (
	"testing"
)

func TestValidateAndFix(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		expectFixed    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name:         "nil config",
			config:       nil,
			expectFixed:  false,
			expectErrors: 1,
		},
		{
			name: "valid config",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
					"hotfix":  "hotfix/",
				},
				Sanitization: SanitizationConfig{
					Separator:     "-",
					Lowercase:     true,
					RemoveUmlauts: false,
				},
			},
			expectFixed:    false,
			expectErrors:   0,
			expectWarnings: 0,
		},
		{
			name: "invalid max_branch_length - too small",
			config: &Config{
				MaxBranchLength:   5,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "invalid max_branch_length - too large",
			config: &Config{
				MaxBranchLength:   300,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "empty branch_types",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes:       map[string]string{},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "empty branch type value",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "",
					"hotfix":  "hotfix/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "empty branch type key",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"":        "empty/",
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:  false,
			expectErrors: 1,
		},
		{
			name: "invalid default_branch_type",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "nonexistent",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "empty separator",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "separator too long",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "------",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
		{
			name: "separator with problematic characters",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "/",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAndFix(tt.config)

			if result.Fixed != tt.expectFixed {
				t.Errorf("ValidateAndFix() fixed = %v, want %v", result.Fixed, tt.expectFixed)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("ValidateAndFix() errors = %d, want %d", len(result.Errors), tt.expectErrors)
			}

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("ValidateAndFix() warnings = %d, want %d", len(result.Warnings), tt.expectWarnings)
			}

			// If config was fixed and no errors, verify the config is now valid
			if tt.config != nil && result.Fixed && len(result.Errors) == 0 {
				strictResult := ValidateStrict(tt.config)
				if strictResult != nil {
					t.Errorf("ValidateAndFix() should have fixed config, but strict validation failed: %v", strictResult)
				}
			}
		})
	}
}

func TestValidateStrict(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "valid config",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
					"hotfix":  "hotfix/",
				},
				Sanitization: SanitizationConfig{
					Separator:     "-",
					Lowercase:     true,
					RemoveUmlauts: false,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid max_branch_length - too small",
			config: &Config{
				MaxBranchLength:   5,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid max_branch_length - too large",
			config: &Config{
				MaxBranchLength:   300,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
		},
		{
			name: "empty branch_types",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes:       map[string]string{},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
		},
		{
			name: "empty branch type key",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"":        "empty/",
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
		},
		{
			name: "empty branch type value",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid default_branch_type",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "nonexistent",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
		},
		{
			name: "empty separator",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "",
				},
			},
			wantErr: true,
		},
		{
			name: "separator too long",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "------",
				},
			},
			wantErr: true,
		},
		{
			name: "separator with problematic characters",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "/",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStrict(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStrict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "test_field",
		Value:   "test_value",
		Message: "test message",
	}

	expected := "validation error for field 'test_field' (value: test_value): test message"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestValidationResult_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		result *ValidationResult
		want   bool
	}{
		{
			name: "no errors",
			result: &ValidationResult{
				Errors: []ValidationError{},
			},
			want: true,
		},
		{
			name: "with errors",
			result: &ValidationResult{
				Errors: []ValidationError{
					{Field: "test", Value: "test", Message: "test"},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsValid(); got != tt.want {
				t.Errorf("ValidationResult.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}