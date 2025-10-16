package branch

import (
	"fmt"
	"regexp"
	"strings"
)

// GeneratorConfig holds configuration for branch name generation
type GeneratorConfig struct {
	MaxBranchLength int
	Separator       string
	Lowercase       bool
	RemoveUmlauts   bool
}

// BranchInfo represents information needed to generate a branch name
type BranchInfo struct {
	Type       string
	TicketID   string
	Title      string
	BaseBranch string
	FullName   string
}

// Generator interface defines branch name generation operations
type Generator interface {
	GenerateName(info BranchInfo) string
	ValidateName(name string) error
}

// BranchGenerator implements the Generator interface
type BranchGenerator struct {
	sanitizer Sanitizer
}

// NewBranchGenerator creates a new BranchGenerator instance
func NewBranchGenerator(sanitizer Sanitizer) *BranchGenerator {
	return &BranchGenerator{
		sanitizer: sanitizer,
	}
}

// GenerateName generates a branch name from the provided BranchInfo
// Format: type/ticket-title
func (g *BranchGenerator) GenerateName(info BranchInfo) string {
	return g.GenerateNameWithConfig(info, GeneratorConfig{
		MaxBranchLength: 60, // Default max length
		Separator:       "-",
		Lowercase:       true,
		RemoveUmlauts:   true,
	})
}

// GenerateNameWithConfig generates a branch name with specific configuration
func (g *BranchGenerator) GenerateNameWithConfig(info BranchInfo, config GeneratorConfig) string {
	if info.Type == "" || info.TicketID == "" {
		return ""
	}

	// Start with the ticket title, use ticket ID if title is empty
	title := info.Title
	if title == "" {
		title = info.TicketID
	}

	// Calculate available length for title part
	// Format: type/ticket-title, so we need to account for type, slash, ticket, and separator
	prefixLength := len(info.Type) + 1 + len(info.TicketID) + len(config.Separator) // type + "/" + ticket + separator
	availableTitleLength := config.MaxBranchLength - prefixLength

	// Ensure we have at least some space for the title
	if availableTitleLength < 1 {
		availableTitleLength = 10 // Minimum title length
	}

	// Sanitize the title using the sanitizer
	sanitizedTitle := g.sanitizer.Sanitize(title, SanitizationOptions{
		Separator:     config.Separator,
		Lowercase:     config.Lowercase,
		RemoveUmlauts: config.RemoveUmlauts,
		MaxLength:     availableTitleLength,
	})

	// Create the branch name in format: type/ticket-title
	branchName := fmt.Sprintf("%s/%s%s%s", info.Type, info.TicketID, config.Separator, sanitizedTitle)

	// Final length check and truncation if needed
	if len(branchName) > config.MaxBranchLength {
		// Truncate from the title part while preserving the format
		maxTitleLength := config.MaxBranchLength - prefixLength
		if maxTitleLength > 0 {
			truncatedTitle := sanitizedTitle
			if len(sanitizedTitle) > maxTitleLength {
				truncatedTitle = sanitizedTitle[:maxTitleLength]
				// Remove trailing separator if truncation created one
				truncatedTitle = strings.TrimSuffix(truncatedTitle, config.Separator)
			}
			branchName = fmt.Sprintf("%s/%s%s%s", info.Type, info.TicketID, config.Separator, truncatedTitle)
		}
	}

	return branchName
}

// ValidateName validates a branch name according to Git naming rules
func (g *BranchGenerator) ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	// Git branch name validation rules
	invalidPatterns := []string{
		`^\.`,           // Cannot start with dot
		`\.$`,           // Cannot end with dot
		`\.\.`,          // Cannot contain double dots
		`^/`,            // Cannot start with slash
		`/$`,            // Cannot end with slash
		`//`,            // Cannot contain double slashes
		`\s`,            // Cannot contain spaces
		`[\x00-\x1f\x7f]`, // Cannot contain control characters
		`[~^:?*\[]`,     // Cannot contain special Git characters
	}

	for _, pattern := range invalidPatterns {
		matched, err := regexp.MatchString(pattern, name)
		if err != nil {
			return fmt.Errorf("error validating branch name: %v", err)
		}
		if matched {
			return fmt.Errorf("invalid branch name: contains invalid pattern %s", pattern)
		}
	}

	return nil
}

// GeneratorConfigFromAppConfig creates a GeneratorConfig from application config
func GeneratorConfigFromAppConfig(maxLength int, separator string, lowercase, removeUmlauts bool) GeneratorConfig {
	return GeneratorConfig{
		MaxBranchLength: maxLength,
		Separator:       separator,
		Lowercase:       lowercase,
		RemoveUmlauts:   removeUmlauts,
	}
}