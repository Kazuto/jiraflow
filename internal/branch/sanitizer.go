package branch

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