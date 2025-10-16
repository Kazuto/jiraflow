package config

// Config represents the application configuration structure
type Config struct {
	MaxBranchLength   int                    `yaml:"max_branch_length"`
	DefaultBranchType string                 `yaml:"default_branch_type"`
	BranchTypes       map[string]string      `yaml:"branch_types"`
	Sanitization      SanitizationConfig     `yaml:"sanitization"`
}

// SanitizationConfig holds sanitization-related settings
type SanitizationConfig struct {
	Separator     string `yaml:"separator"`
	Lowercase     bool   `yaml:"lowercase"`
	RemoveUmlauts bool   `yaml:"remove_umlauts"`
}

// ConfigManager interface defines configuration management operations
type ConfigManager interface {
	Load() (*Config, error)
	CreateDefault() error
	Validate(*Config) error
	GetConfigPath() string
}