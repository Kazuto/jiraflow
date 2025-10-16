package config

import (
	"fmt"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s' (value: %v): %s", e.Field, e.Value, e.Message)
}

// ValidationResult holds the result of configuration validation
type ValidationResult struct {
	Errors   []ValidationError
	Warnings []string
	Fixed    bool // indicates if any values were fixed with defaults
}

// IsValid returns true if there are no validation errors
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// ValidateAndFix validates the configuration and fixes invalid values with defaults
func ValidateAndFix(config *Config) *ValidationResult {
	if config == nil {
		return &ValidationResult{
			Errors: []ValidationError{
				{Field: "config", Value: nil, Message: "configuration cannot be nil"},
			},
		}
	}

	result := &ValidationResult{}
	defaults := GetDefaultConfig()

	// Validate and fix max_branch_length
	if config.MaxBranchLength < 10 || config.MaxBranchLength > 200 {
		result.Warnings = append(result.Warnings, 
			fmt.Sprintf("max_branch_length %d is out of range (10-200), using default %d", 
				config.MaxBranchLength, defaults.MaxBranchLength))
		config.MaxBranchLength = defaults.MaxBranchLength
		result.Fixed = true
	}

	// Validate and fix branch_types
	if len(config.BranchTypes) == 0 {
		result.Warnings = append(result.Warnings, 
			"branch_types is empty, using default branch types")
		config.BranchTypes = make(map[string]string)
		for k, v := range defaults.BranchTypes {
			config.BranchTypes[k] = v
		}
		result.Fixed = true
	} else {
		// Check for empty keys or values in branch_types
		for key, value := range config.BranchTypes {
			if key == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "branch_types",
					Value:   key,
					Message: "branch type key cannot be empty",
				})
			}
			if value == "" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("branch type value for key '%s' is empty, using default", key))
				if defaultValue, exists := defaults.BranchTypes[key]; exists {
					config.BranchTypes[key] = defaultValue
				} else {
					config.BranchTypes[key] = key + "/"
				}
				result.Fixed = true
			}
		}
	}

	// Validate and fix default_branch_type
	if config.DefaultBranchType == "" {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("default_branch_type is empty, using default '%s'", defaults.DefaultBranchType))
		config.DefaultBranchType = defaults.DefaultBranchType
		result.Fixed = true
	} else if _, exists := config.BranchTypes[config.DefaultBranchType]; !exists {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("default_branch_type '%s' does not exist in branch_types, using default '%s'", 
				config.DefaultBranchType, defaults.DefaultBranchType))
		config.DefaultBranchType = defaults.DefaultBranchType
		result.Fixed = true
	}

	// Validate and fix sanitization settings
	if config.Sanitization.Separator == "" {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("sanitization separator is empty, using default '%s'", defaults.Sanitization.Separator))
		config.Sanitization.Separator = defaults.Sanitization.Separator
		result.Fixed = true
	}

	// Validate separator length (should be reasonable)
	if len(config.Sanitization.Separator) > 5 {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("sanitization separator '%s' is too long (>5 chars), using default '%s'", 
				config.Sanitization.Separator, defaults.Sanitization.Separator))
		config.Sanitization.Separator = defaults.Sanitization.Separator
		result.Fixed = true
	}

	// Validate separator doesn't contain problematic characters
	problematicChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	for _, char := range problematicChars {
		if strings.Contains(config.Sanitization.Separator, char) {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("sanitization separator '%s' contains problematic character '%s', using default '%s'", 
					config.Sanitization.Separator, char, defaults.Sanitization.Separator))
			config.Sanitization.Separator = defaults.Sanitization.Separator
			result.Fixed = true
			break
		}
	}

	return result
}

// ValidateStrict performs strict validation without fixing values
func ValidateStrict(config *Config) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate max_branch_length
	if config.MaxBranchLength < 10 || config.MaxBranchLength > 200 {
		return &ValidationError{
			Field:   "max_branch_length",
			Value:   config.MaxBranchLength,
			Message: "must be between 10 and 200",
		}
	}

	// Validate branch_types
	if len(config.BranchTypes) == 0 {
		return &ValidationError{
			Field:   "branch_types",
			Value:   config.BranchTypes,
			Message: "cannot be empty",
		}
	}

	for key, value := range config.BranchTypes {
		if key == "" {
			return &ValidationError{
				Field:   "branch_types",
				Value:   key,
				Message: "branch type key cannot be empty",
			}
		}
		if value == "" {
			return &ValidationError{
				Field:   "branch_types",
				Value:   fmt.Sprintf("key '%s'", key),
				Message: "branch type value cannot be empty",
			}
		}
	}

	// Validate default_branch_type
	if config.DefaultBranchType == "" {
		return &ValidationError{
			Field:   "default_branch_type",
			Value:   config.DefaultBranchType,
			Message: "cannot be empty",
		}
	}

	if _, exists := config.BranchTypes[config.DefaultBranchType]; !exists {
		return &ValidationError{
			Field:   "default_branch_type",
			Value:   config.DefaultBranchType,
			Message: "must exist in branch_types",
		}
	}

	// Validate sanitization settings
	if config.Sanitization.Separator == "" {
		return &ValidationError{
			Field:   "sanitization.separator",
			Value:   config.Sanitization.Separator,
			Message: "cannot be empty",
		}
	}

	if len(config.Sanitization.Separator) > 5 {
		return &ValidationError{
			Field:   "sanitization.separator",
			Value:   config.Sanitization.Separator,
			Message: "cannot be longer than 5 characters",
		}
	}

	// Check for problematic characters in separator
	problematicChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	for _, char := range problematicChars {
		if strings.Contains(config.Sanitization.Separator, char) {
			return &ValidationError{
				Field:   "sanitization.separator",
				Value:   config.Sanitization.Separator,
				Message: fmt.Sprintf("cannot contain problematic character '%s'", char),
			}
		}
	}

	return nil
}