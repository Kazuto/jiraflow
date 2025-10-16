package config

import (
	"os"
	"path/filepath"
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