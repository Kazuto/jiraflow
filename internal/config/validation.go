package config

import "fmt"

// ConfigError represents a configuration validation error
type ConfigError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("config validation error for field '%s': %s (value: %v)", 
		e.Field, e.Message, e.Value)
}

// ValidateConfig validates the configuration values
func ValidateConfig(cfg *Config) error {
	// Validate max_branch_length
	if cfg.MaxBranchLength < 10 || cfg.MaxBranchLength > 200 {
		return ConfigError{
			Field:   "max_branch_length",
			Value:   cfg.MaxBranchLength,
			Message: "must be between 10 and 200",
		}
	}

	// Validate branch_types
	if len(cfg.BranchTypes) == 0 {
		return ConfigError{
			Field:   "branch_types",
			Value:   cfg.BranchTypes,
			Message: "must contain at least one branch type",
		}
	}

	for key, value := range cfg.BranchTypes {
		if value == "" {
			return ConfigError{
				Field:   fmt.Sprintf("branch_types.%s", key),
				Value:   value,
				Message: "branch type value cannot be empty",
			}
		}
	}

	return nil
}