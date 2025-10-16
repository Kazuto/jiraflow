package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfigManager implements the ConfigManager interface for file-based configuration
type FileConfigManager struct {
	configPath string
}

// NewFileConfigManager creates a new FileConfigManager instance
func NewFileConfigManager() *FileConfigManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is not accessible
		homeDir = "."
	}
	
	configPath := filepath.Join(homeDir, ".config", "jiraflow", "jiraflow.yaml")
	return &FileConfigManager{
		configPath: configPath,
	}
}

// GetConfigPath returns the path to the configuration file
func (m *FileConfigManager) GetConfigPath() string {
	return m.configPath
}

// Load reads and parses the configuration file
func (m *FileConfigManager) Load() (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// Create default configuration if file doesn't exist
		if err := m.CreateDefault(); err != nil {
			return nil, fmt.Errorf("failed to create default configuration: %w", err)
		}
	}

	// Read the configuration file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", m.configPath, err)
	}

	// Parse YAML content
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %s: %w", m.configPath, err)
	}

	// Validate the loaded configuration
	if err := m.Validate(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// CreateDefault creates the configuration directory and file with default values
func (m *FileConfigManager) CreateDefault() error {
	// Create the directory structure
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create configuration directory %s: %w", configDir, err)
	}

	// Create YAML content with comments
	yamlContent := `# Maximum branch name length (default: 60)
max_branch_length: 60

# Default branch type if not specified (default: feature)
default_branch_type: feature

# Branch type prefixes
branch_types:
  feature: "feature/"
  hotfix: "hotfix/"
  refactor: "refactor/"
  support: "support/"

# Character replacements for branch name sanitization
sanitization:
  # Replace spaces and special characters with this (default: -)
  separator: "-"
  # Convert to lowercase (default: true)
  lowercase: true
  # Remove German umlauts (äöüÄÖÜß) (default: false)
  remove_umlauts: false
`

	// Write to file
	if err := os.WriteFile(m.configPath, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write default configuration to %s: %w", m.configPath, err)
	}

	return nil
}

// Validate checks if the configuration values are valid
func (m *FileConfigManager) Validate(config *Config) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate max_branch_length
	if config.MaxBranchLength < 10 || config.MaxBranchLength > 200 {
		return fmt.Errorf("max_branch_length must be between 10 and 200, got %d", config.MaxBranchLength)
	}

	// Validate branch_types
	if len(config.BranchTypes) == 0 {
		return fmt.Errorf("branch_types cannot be empty")
	}

	for key, value := range config.BranchTypes {
		if key == "" {
			return fmt.Errorf("branch type key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("branch type value for key '%s' cannot be empty", key)
		}
	}

	// Validate default_branch_type exists in branch_types
	if _, exists := config.BranchTypes[config.DefaultBranchType]; !exists {
		return fmt.Errorf("default_branch_type '%s' must exist in branch_types", config.DefaultBranchType)
	}

	// Validate sanitization settings
	if config.Sanitization.Separator == "" {
		return fmt.Errorf("sanitization separator cannot be empty")
	}

	return nil
}