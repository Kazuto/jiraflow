package config

import (
	"fmt"
	"strings"
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

func TestValidateAndFix_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		expectFixed    bool
		expectErrors   int
		expectWarnings int
		validateResult func(*Config, *ValidationResult) error
	}{
		{
			name: "multiple validation issues fixed",
			config: &Config{
				MaxBranchLength:   5,    // too small
				DefaultBranchType: "",   // empty
				BranchTypes:       map[string]string{}, // empty
				Sanitization: SanitizationConfig{
					Separator: "", // empty
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 4, // max_branch_length, default_branch_type, branch_types, separator
			validateResult: func(c *Config, r *ValidationResult) error {
				defaults := GetDefaultConfig()
				if c.MaxBranchLength != defaults.MaxBranchLength {
					return fmt.Errorf("MaxBranchLength not fixed to default")
				}
				if c.DefaultBranchType != defaults.DefaultBranchType {
					return fmt.Errorf("DefaultBranchType not fixed to default")
				}
				if len(c.BranchTypes) == 0 {
					return fmt.Errorf("BranchTypes not populated with defaults")
				}
				if c.Sanitization.Separator != defaults.Sanitization.Separator {
					return fmt.Errorf("Separator not fixed to default")
				}
				return nil
			},
		},
		{
			name: "branch type with empty value gets fixed",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
					"hotfix":  "", // empty value
					"custom":  "", // empty value for non-default type
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 2, // two empty values
			validateResult: func(c *Config, r *ValidationResult) error {
				if c.BranchTypes["hotfix"] == "" {
					return fmt.Errorf("hotfix branch type value should be fixed")
				}
				if c.BranchTypes["custom"] == "" {
					return fmt.Errorf("custom branch type value should be fixed")
				}
				return nil
			},
		},
		{
			name: "separator with multiple problematic characters",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "/<>", // multiple problematic chars
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
			validateResult: func(c *Config, r *ValidationResult) error {
				if c.Sanitization.Separator != "-" {
					return fmt.Errorf("separator should be fixed to default")
				}
				return nil
			},
		},
		{
			name: "default_branch_type references non-existent type",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "nonexistent",
				BranchTypes: map[string]string{
					"feature": "feature/",
					"hotfix":  "hotfix/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			expectFixed:    true,
			expectErrors:   0,
			expectWarnings: 1,
			validateResult: func(c *Config, r *ValidationResult) error {
				if c.DefaultBranchType != "feature" {
					return fmt.Errorf("default_branch_type should be fixed to 'feature'")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAndFix(tt.config)

			if result.Fixed != tt.expectFixed {
				t.Errorf("ValidateAndFix() fixed = %v, want %v", result.Fixed, tt.expectFixed)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("ValidateAndFix() errors = %d, want %d. Errors: %v", len(result.Errors), tt.expectErrors, result.Errors)
			}

			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("ValidateAndFix() warnings = %d, want %d. Warnings: %v", len(result.Warnings), tt.expectWarnings, result.Warnings)
			}

			if tt.validateResult != nil {
				if err := tt.validateResult(tt.config, result); err != nil {
					t.Errorf("Result validation failed: %v", err)
				}
			}

			// Verify that fixed config passes strict validation
			if tt.config != nil && result.Fixed && len(result.Errors) == 0 {
				if err := ValidateStrict(tt.config); err != nil {
					t.Errorf("Fixed config should pass strict validation: %v", err)
				}
			}
		})
	}
}

func TestValidateStrict_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		wantErr    bool
		errorCheck func(error) bool
	}{
		{
			name: "boundary values - min valid max_branch_length",
			config: &Config{
				MaxBranchLength:   10, // minimum valid
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: false,
		},
		{
			name: "boundary values - max valid max_branch_length",
			config: &Config{
				MaxBranchLength:   200, // maximum valid
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: false,
		},
		{
			name: "boundary values - just below min max_branch_length",
			config: &Config{
				MaxBranchLength:   9, // just below minimum
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "max_branch_length") && strings.Contains(err.Error(), "between 10 and 200")
			},
		},
		{
			name: "boundary values - just above max max_branch_length",
			config: &Config{
				MaxBranchLength:   201, // just above maximum
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "-",
				},
			},
			wantErr: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "max_branch_length") && strings.Contains(err.Error(), "between 10 and 200")
			},
		},
		{
			name: "separator at max length boundary",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "12345", // exactly 5 chars (max allowed)
				},
			},
			wantErr: false,
		},
		{
			name: "separator exceeds max length",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "123456", // 6 chars (too long)
				},
			},
			wantErr: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "sanitization.separator") && strings.Contains(err.Error(), "longer than 5 characters")
			},
		},
		{
			name: "all problematic separator characters",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes: map[string]string{
					"feature": "feature/",
				},
				Sanitization: SanitizationConfig{
					Separator: "\\", // backslash
				},
			},
			wantErr: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "sanitization.separator") && strings.Contains(err.Error(), "problematic character")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStrict(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStrict() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.errorCheck != nil && err != nil && !tt.errorCheck(err) {
				t.Errorf("ValidateStrict() error validation failed: %v", err)
			}
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	
	if config == nil {
		t.Fatal("GetDefaultConfig() returned nil")
	}

	// Verify default values match expected values
	if config.MaxBranchLength != 60 {
		t.Errorf("Default MaxBranchLength = %d, want 60", config.MaxBranchLength)
	}

	if config.DefaultBranchType != "feature" {
		t.Errorf("Default DefaultBranchType = %s, want 'feature'", config.DefaultBranchType)
	}

	expectedBranchTypes := []string{"feature", "hotfix", "refactor", "support"}
	for _, branchType := range expectedBranchTypes {
		if _, exists := config.BranchTypes[branchType]; !exists {
			t.Errorf("Default BranchTypes missing '%s'", branchType)
		}
	}

	if config.Sanitization.Separator != "-" {
		t.Errorf("Default Separator = %s, want '-'", config.Sanitization.Separator)
	}

	if !config.Sanitization.Lowercase {
		t.Error("Default Lowercase should be true")
	}

	if config.Sanitization.RemoveUmlauts {
		t.Error("Default RemoveUmlauts should be false")
	}

	// Verify the default config passes strict validation
	if err := ValidateStrict(config); err != nil {
		t.Errorf("Default config should pass strict validation: %v", err)
	}
}