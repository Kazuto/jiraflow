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
// Implements comprehensive sanitization logic:
// 1. Handle German umlauts (if enabled)
// 2. Remove quotes, parentheses, colons, and other problematic characters
// 3. Replace " - " (space-hyphen-space) with just the separator
// 4. Replace remaining spaces with the configured separator
// 5. Remove consecutive separators (double-hyphen cleanup)
// 6. Remove special characters that might cause Git issues
// 7. Apply case conversion (if enabled)
// 8. Trim to specified length
// 9. Remove trailing separator if any
func (s *BranchSanitizer) Sanitize(input string, options SanitizationOptions) string {
	if input == "" {
		return ""
	}

	result := strings.TrimSpace(input)

	// 1. Handle German umlauts first if requested (before other character removal)
	if options.RemoveUmlauts {
		result = s.removeUmlauts(result)
	}

	// 2. Remove quotes, parentheses, colons, brackets, and other problematic characters
	// Expanded to include more special characters that can cause issues
	// Replace path separators with spaces first to preserve word boundaries
	result = strings.ReplaceAll(result, "/", " ")
	result = strings.ReplaceAll(result, "\\", " ")
	result = regexp.MustCompile(`["\(\)\[\]{}:;,<>?|*&^%$#@!~` + "`" + `]`).ReplaceAllString(result, "")

	// 3. Replace " - " (space-hyphen-space) and similar patterns with just the separator
	separator := options.Separator
	if separator == "" {
		separator = "-"
	}
	
	// Handle various hyphen-space combinations
	result = regexp.MustCompile(`\s*-\s*`).ReplaceAllString(result, separator)
	result = regexp.MustCompile(`\s*_\s*`).ReplaceAllString(result, separator)

	// 4. Replace remaining spaces and tabs with the specified separator
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, separator)

	// 5. Remove consecutive separators (double-hyphen cleanup)
	// This handles cases where multiple separators appear consecutively
	doublePattern := regexp.MustCompile(regexp.QuoteMeta(separator) + `{2,}`)
	result = doublePattern.ReplaceAllString(result, separator)

	// 6. Remove special characters that might cause Git issues
	// Keep only alphanumeric characters, the separator, and dots (for version numbers)
	allowedChars := `a-zA-Z0-9` + regexp.QuoteMeta(separator) + `\.`
	result = regexp.MustCompile(`[^`+allowedChars+`]`).ReplaceAllString(result, "")

	// 7. Apply case conversion if requested
	if options.Lowercase {
		result = strings.ToLower(result)
	}

	// 8. Trim to specified length (ensuring we don't break in the middle of a word if possible)
	if options.MaxLength > 0 && len(result) > options.MaxLength {
		result = result[:options.MaxLength]
		// Try to break at a separator to avoid cutting words
		if lastSep := strings.LastIndex(result, separator); lastSep > options.MaxLength/2 {
			result = result[:lastSep]
		}
	}

	// 9. Clean up leading and trailing separators
	result = strings.Trim(result, separator)

	// Final cleanup: ensure no leading dots (Git doesn't allow branches starting with dots)
	result = strings.TrimLeft(result, ".")

	return result
}

// removeUmlauts replaces German umlauts and other special characters with their ASCII equivalents
// Handles both standard German umlauts and additional European characters
func (s *BranchSanitizer) removeUmlauts(input string) string {
	replacements := map[string]string{
		// German umlauts
		"ä": "ae", "Ä": "Ae",
		"ö": "oe", "Ö": "Oe", 
		"ü": "ue", "Ü": "Ue",
		"ß": "ss",
		
		// Additional European characters that might appear in branch names
		"à": "a", "á": "a", "â": "a", "ã": "a", "å": "a", "æ": "ae",
		"À": "A", "Á": "A", "Â": "A", "Ã": "A", "Å": "A", "Æ": "Ae",
		"è": "e", "é": "e", "ê": "e", "ë": "e",
		"È": "E", "É": "E", "Ê": "E", "Ë": "E",
		"ì": "i", "í": "i", "î": "i", "ï": "i",
		"Ì": "I", "Í": "I", "Î": "I", "Ï": "I",
		"ò": "o", "ó": "o", "ô": "o", "õ": "o", "ø": "o",
		"Ò": "O", "Ó": "O", "Ô": "O", "Õ": "O", "Ø": "O",
		"ù": "u", "ú": "u", "û": "u",
		"Ù": "U", "Ú": "U", "Û": "U",
		"ç": "c", "Ç": "C",
		"ñ": "n", "Ñ": "N",
	}

	result := input
	for char, replacement := range replacements {
		result = strings.ReplaceAll(result, char, replacement)
	}

	return result
}