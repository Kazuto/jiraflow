package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileConfigManager_CreateDefault(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Create a config manager with a custom path
	manager := &FileConfigManager{
		configPath: filepath.Join(tempDir, ".config", "jiraflow", "jiraflow.yaml"),
	}

	// Test creating default configuration
	err := manager.CreateDefault()
	if err != nil {
		t.Fatalf("CreateDefault() failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(manager.configPath); os.IsNotExist(err) {
		t.Fatalf("Configuration file was not created at %s", manager.configPath)
	}

	// Test loading the created configuration
	config, err := manager.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify default values
	if config.MaxBranchLength != 60 {
		t.Errorf("Expected MaxBranchLength to be 60, got %d", config.MaxBranchLength)
	}

	if config.DefaultBranchType != "feature" {
		t.Errorf("Expected DefaultBranchType to be 'feature', got %s", config.DefaultBranchType)
	}

	if len(config.BranchTypes) == 0 {
		t.Error("Expected BranchTypes to be populated")
	}
}

func TestFileConfigManager_Validate(t *testing.T) {
	manager := &FileConfigManager{}

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
			name: "invalid max_branch_length too small",
			config: &Config{
				MaxBranchLength:   5,
				DefaultBranchType: "feature",
				BranchTypes:       map[string]string{"feature": "feature"},
				Sanitization:      SanitizationConfig{Separator: "-"},
			},
			wantErr: true,
		},
		{
			name: "invalid max_branch_length too large",
			config: &Config{
				MaxBranchLength:   300,
				DefaultBranchType: "feature",
				BranchTypes:       map[string]string{"feature": "feature"},
				Sanitization:      SanitizationConfig{Separator: "-"},
			},
			wantErr: true,
		},
		{
			name: "empty branch_types",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes:       map[string]string{},
				Sanitization:      SanitizationConfig{Separator: "-"},
			},
			wantErr: true,
		},
		{
			name: "invalid default_branch_type",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "invalid",
				BranchTypes:       map[string]string{"feature": "feature"},
				Sanitization:      SanitizationConfig{Separator: "-"},
			},
			wantErr: true,
		},
		{
			name: "empty separator",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes:       map[string]string{"feature": "feature"},
				Sanitization:      SanitizationConfig{Separator: ""},
			},
			wantErr: true,
		},
		{
			name: "valid config",
			config: &Config{
				MaxBranchLength:   60,
				DefaultBranchType: "feature",
				BranchTypes:       map[string]string{"feature": "feature", "hotfix": "hotfix"},
				Sanitization:      SanitizationConfig{Separator: "-", Lowercase: true},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.Validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileConfigManager_CreateDefault_Scenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() string // returns config path
		expectError bool
		errorCheck  func(error) bool // optional error validation
	}{
		{
			name: "creates config in new directory",
			setupFunc: func() string {
				return filepath.Join(t.TempDir(), ".config", "jiraflow", "jiraflow.yaml")
			},
			expectError: false,
		},
		{
			name: "creates config when parent directory exists",
			setupFunc: func() string {
				tempDir := t.TempDir()
				configDir := filepath.Join(tempDir, ".config", "jiraflow")
				_ = os.MkdirAll(configDir, 0750)
				return filepath.Join(configDir, "jiraflow.yaml")
			},
			expectError: false,
		},
		{
			name: "handles existing config file",
			setupFunc: func() string {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, ".config", "jiraflow", "jiraflow.yaml")
				manager := &FileConfigManager{configPath: configPath}
				_ = manager.CreateDefault() // Create it first
				return configPath
			},
			expectError: false, // Should overwrite existing file
		},
		{
			name: "fails with read-only directory",
			setupFunc: func() string {
				tempDir := t.TempDir()
				readOnlyDir := filepath.Join(tempDir, "readonly")
				_ = os.MkdirAll(readOnlyDir, 0444) //nolint:gosec // Read-only permissions intentional for testing
				return filepath.Join(readOnlyDir, ".config", "jiraflow", "jiraflow.yaml")
			},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "failed to create configuration directory")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := tt.setupFunc()
			manager := &FileConfigManager{configPath: configPath}

			err := manager.CreateDefault()

			if tt.expectError {
				if err == nil {
					t.Errorf("CreateDefault() expected error but got none")
				} else if tt.errorCheck != nil && !tt.errorCheck(err) {
					t.Errorf("CreateDefault() error validation failed: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("CreateDefault() unexpected error: %v", err)
			}

			// Verify file was created and has correct content
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Errorf("Configuration file was not created at %s", configPath)
			}

			// Verify content can be loaded
			config, err := manager.Load()
			if err != nil {
				t.Errorf("Failed to load created config: %v", err)
			}

			// Verify default values
			defaults := GetDefaultConfig()
			if config.MaxBranchLength != defaults.MaxBranchLength {
				t.Errorf("MaxBranchLength = %d, want %d", config.MaxBranchLength, defaults.MaxBranchLength)
			}
			if config.DefaultBranchType != defaults.DefaultBranchType {
				t.Errorf("DefaultBranchType = %s, want %s", config.DefaultBranchType, defaults.DefaultBranchType)
			}
		})
	}
}

func TestFileConfigManager_Load_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *FileConfigManager
		expectError bool
		errorCheck  func(error) bool
	}{
		{
			name: "loads existing valid config",
			setupFunc: func() *FileConfigManager {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "jiraflow.yaml")
				manager := &FileConfigManager{configPath: configPath}
				_ = manager.CreateDefault()
				return manager
			},
			expectError: false,
		},
		{
			name: "creates default when config missing",
			setupFunc: func() *FileConfigManager {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "nonexistent", "jiraflow.yaml")
				return &FileConfigManager{configPath: configPath}
			},
			expectError: false,
		},
		{
			name: "fails with invalid YAML",
			setupFunc: func() *FileConfigManager {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "jiraflow.yaml")
				
				// Create invalid YAML content
				invalidYAML := `
max_branch_length: not_a_number
branch_types:
  - invalid: structure
`
				_ = os.WriteFile(configPath, []byte(invalidYAML), 0600)
				return &FileConfigManager{configPath: configPath}
			},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "failed to parse YAML configuration")
			},
		},
		{
			name: "handles config with validation errors that can't be fixed",
			setupFunc: func() *FileConfigManager {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "jiraflow.yaml")
				
				// Create config with unfixable validation error (empty branch type key)
				invalidConfig := `
max_branch_length: 60
default_branch_type: feature
branch_types:
  "": "empty_key"
  feature: "feature/"
sanitization:
  separator: "-"
`
				_ = os.WriteFile(configPath, []byte(invalidConfig), 0600)
				return &FileConfigManager{configPath: configPath}
			},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "branch type key cannot be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.setupFunc()

			config, err := manager.Load()

			if tt.expectError {
				if err == nil {
					t.Errorf("Load() expected error but got none")
				} else if tt.errorCheck != nil && !tt.errorCheck(err) {
					t.Errorf("Load() error validation failed: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}

			if config == nil {
				t.Error("Load() returned nil config")
			}
		})
	}
}

func TestFileConfigManager_Load_FallbackMechanisms(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectWarnings bool
		validateFunc   func(*Config) error
	}{
		{
			name: "fixes invalid max_branch_length",
			configContent: `
max_branch_length: 5
default_branch_type: feature
branch_types:
  feature: "feature/"
sanitization:
  separator: "-"
`,
			expectWarnings: true,
			validateFunc: func(c *Config) error {
				if c.MaxBranchLength != 60 { // Should be fixed to default
					return fmt.Errorf("expected MaxBranchLength to be fixed to 60, got %d", c.MaxBranchLength)
				}
				return nil
			},
		},
		{
			name: "fixes empty branch_types",
			configContent: `
max_branch_length: 60
default_branch_type: feature
branch_types: {}
sanitization:
  separator: "-"
`,
			expectWarnings: true,
			validateFunc: func(c *Config) error {
				if len(c.BranchTypes) == 0 {
					return fmt.Errorf("expected BranchTypes to be populated with defaults")
				}
				return nil
			},
		},
		{
			name: "fixes invalid separator",
			configContent: `
max_branch_length: 60
default_branch_type: feature
branch_types:
  feature: "feature/"
sanitization:
  separator: "/"
`,
			expectWarnings: true,
			validateFunc: func(c *Config) error {
				if c.Sanitization.Separator != "-" {
					return fmt.Errorf("expected separator to be fixed to '-', got %s", c.Sanitization.Separator)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "jiraflow.yaml")
			
			// Write test config
			err := os.WriteFile(configPath, []byte(tt.configContent), 0600)
			if err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			manager := &FileConfigManager{configPath: configPath}
			
			config, err := manager.Load()
			if err != nil {
				t.Fatalf("Load() failed: %v", err)
			}

			if tt.validateFunc != nil {
				if err := tt.validateFunc(config); err != nil {
					t.Errorf("Config validation failed: %v", err)
				}
			}

			// Verify the loaded config passes strict validation
			if err := ValidateStrict(config); err != nil {
				t.Errorf("Fixed config should pass strict validation, but got: %v", err)
			}
		})
	}
}

func TestFileConfigManager_GetConfigPath(t *testing.T) {
	expectedPath := "/test/path/config.yaml"
	manager := &FileConfigManager{configPath: expectedPath}
	
	if got := manager.GetConfigPath(); got != expectedPath {
		t.Errorf("GetConfigPath() = %v, want %v", got, expectedPath)
	}
}

func TestNewFileConfigManager(t *testing.T) {
	manager := NewFileConfigManager()
	
	configPath := manager.GetConfigPath()
	if configPath == "" {
		t.Error("NewFileConfigManager() returned empty config path")
	}
	
	// Should contain the expected path structure
	if !strings.Contains(configPath, ".config/jiraflow/jiraflow.yaml") {
		t.Errorf("Config path should contain '.config/jiraflow/jiraflow.yaml', got: %s", configPath)
	}
}

func TestNewConfigManager(t *testing.T) {
	manager := NewConfigManager()
	if manager == nil {
		t.Error("NewConfigManager() returned nil")
	}
	
	// Should return a FileConfigManager
	if _, ok := manager.(*FileConfigManager); !ok {
		t.Error("NewConfigManager() should return a FileConfigManager")
	}
}