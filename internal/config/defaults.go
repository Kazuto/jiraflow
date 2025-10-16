package config

// GetDefaultConfig returns a configuration with sensible default values
func GetDefaultConfig() *Config {
	return &Config{
		MaxBranchLength:   60,
		DefaultBranchType: "feature",
		BranchTypes: map[string]string{
			"feature": "feature",
			"hotfix":  "hotfix",
			"refactor": "refactor",
			"support": "support",
		},
		Sanitization: SanitizationConfig{
			Separator:     "-",
			Lowercase:     true,
			RemoveUmlauts: false,
		},
	}
}