package branch

import (
	"regexp"
	"strings"
)

// SanitizationOptions defines options for branch name sanitization
type SanitizationOptions struct {
	Separator     string
	Lowercase     bool
	RemoveUmlauts bool
	MaxLength     int
}

// Sanitizer interface defines branch name sanitization operations
type Sanitizer interface {
	Sanitize(input string, options SanitizationOptions) string
}

// BranchSanitizer implements the Sanitizer interface
type BranchSanitizer struct{}

// NewBranchSanitizer creates a new BranchSanitizer instance
func NewBranchSanitizer() *BranchSanitizer {
	return &BranchSanitizer{}
}

// Sanitize sanitizes a string according to the provided options
// Based on the shell script sanitization rules:
// 1. Remove quotes, parentheses, and other special characters
// 2. Replace " - " (space-hyphen-space) with just "-"
// 3. Replace remaining spaces with hyphens
// 4. Remove double hyphens
// 5. Remove special characters that might cause issues
// 6. Trim to specified length
// 7. Remove trailing hyphen if any
func (s *BranchSanitizer) Sanitize(input string, options SanitizationOptions) string {
	if input == "" {
		return ""
	}

	result := input

	// Handle German umlauts first if requested (before other character removal)
	if options.RemoveUmlauts {
		result = s.removeUmlauts(result)
	}

	// 1. Remove quotes, parentheses, colons, and other special characters
	result = regexp.MustCompile(`["\(\):]`).ReplaceAllString(result, "")

	// 2. Replace " - " (space-hyphen-space) with just "-"
	result = strings.ReplaceAll(result, " - ", "-")

	// 3. Replace remaining spaces with the specified separator (default: hyphen)
	separator := options.Separator
	if separator == "" {
		separator = "-"
	}
	result = strings.ReplaceAll(result, " ", separator)

	// 4. Remove double separators
	doublePattern := regexp.MustCompile(regexp.QuoteMeta(separator) + `+`)
	result = doublePattern.ReplaceAllString(result, separator)

	// 5. Remove special characters that might cause issues (keep only alphanumeric and the separator)
	result = regexp.MustCompile(`[^a-zA-Z0-9` + regexp.QuoteMeta(separator) + `]`).ReplaceAllString(result, "")

	// Convert to lowercase if requested
	if options.Lowercase {
		result = strings.ToLower(result)
	}

	// 6. Trim to specified length
	if options.MaxLength > 0 && len(result) > options.MaxLength {
		result = result[:options.MaxLength]
	}

	// 7. Remove trailing separator if any
	result = strings.TrimSuffix(result, separator)

	return result
}

// removeUmlauts replaces German umlauts with their ASCII equivalents
func (s *BranchSanitizer) removeUmlauts(input string) string {
	replacements := map[string]string{
		"ä": "ae",
		"ö": "oe",
		"ü": "ue",
		"Ä": "Ae",
		"Ö": "Oe",
		"Ü": "Ue",
		"ß": "ss",
	}

	result := input
	for umlaut, replacement := range replacements {
		result = strings.ReplaceAll(result, umlaut, replacement)
	}

	return result
}