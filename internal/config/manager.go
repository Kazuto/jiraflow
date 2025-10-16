package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"jiraflow/internal/errors"
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
			return nil, errors.NewConfigError("", nil, fmt.Sprintf("failed to create default configuration: %v", err), true)
		}
	}

	// Read the configuration file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, errors.NewConfigError("", m.configPath, fmt.Sprintf("failed to read configuration file: %v", err), true)
	}

	// Parse YAML content
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.NewConfigError("", m.configPath, fmt.Sprintf("failed to parse YAML configuration: %v", err), true)
	}

	// Validate and fix the loaded configuration
	result := ValidateAndFix(&config)
	if !result.IsValid() {
		// If there are validation errors that couldn't be fixed, return the first error
		return nil, &result.Errors[0]
	}

	// Log warnings if any values were fixed
	if result.Fixed && len(result.Warnings) > 0 {
		fmt.Fprintf(os.Stderr, "Configuration warnings (values were automatically corrected):\n")
		for _, warning := range result.Warnings {
			fmt.Fprintf(os.Stderr, "  - %s\n", warning)
		}
	}

	return &config, nil
}

// CreateDefault creates the configuration directory and file with default values
func (m *FileConfigManager) CreateDefault() error {
	// Create the directory structure
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return errors.NewConfigError("", configDir, fmt.Sprintf("failed to create configuration directory: %v", err), false)
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
	if err := os.WriteFile(m.configPath, []byte(yamlContent), 0600); err != nil {
		return errors.NewConfigError("", m.configPath, fmt.Sprintf("failed to write default configuration: %v", err), false)
	}

	return nil
}

// Validate checks if the configuration values are valid using strict validation
func (m *FileConfigManager) Validate(config *Config) error {
	return ValidateStrict(config)
}