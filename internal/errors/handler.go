package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ErrorHandler provides centralized error handling functionality
type ErrorHandler struct {
	// Styles for different error display contexts
	cliErrorStyle lipgloss.Style
	tuiErrorStyle lipgloss.Style
	warnStyle     lipgloss.Style
	suggestionStyle lipgloss.Style
}

// NewErrorHandler creates a new ErrorHandler with default styles
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		cliErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),
		tuiErrorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")),
		warnStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true),
		suggestionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Italic(true),
	}
}

// HandleError processes an error and returns appropriate exit code
func (h *ErrorHandler) HandleError(err error) int {
	if err == nil {
		return 0
	}

	// Check if it's a JiraFlowError
	if jfErr, ok := err.(JiraFlowError); ok {
		return h.handleJiraFlowError(jfErr)
	}

	// Handle standard errors
	fmt.Fprintf(os.Stderr, "%s\n", h.cliErrorStyle.Render("Error: "+err.Error()))
	return 1
}

// handleJiraFlowError handles JiraFlowError types with appropriate exit codes
func (h *ErrorHandler) handleJiraFlowError(err JiraFlowError) int {
	// Display user-friendly error message
	fmt.Fprintf(os.Stderr, "%s\n", h.cliErrorStyle.Render("❌ "+err.UserMessage()))
	
	// Display suggestions if available
	suggestions := err.Suggestions()
	if len(suggestions) > 0 {
		fmt.Fprintf(os.Stderr, "\n%s\n", h.warnStyle.Render("💡 Suggestions:"))
		for _, suggestion := range suggestions {
			fmt.Fprintf(os.Stderr, "  %s\n", h.suggestionStyle.Render("• "+suggestion))
		}
	}

	// Return appropriate exit code based on error type
	switch err.Type() {
	case ErrorTypeConfig:
		return 2
	case ErrorTypeGit:
		return 3
	case ErrorTypeJira:
		return 4
	case ErrorTypeTUI:
		return 5
	default:
		return 1
	}
}

// FormatErrorForTUI formats an error for display within the TUI
func (h *ErrorHandler) FormatErrorForTUI(err error) string {
	if err == nil {
		return ""
	}

	var content strings.Builder
	
	// Check if it's a JiraFlowError
	if jfErr, ok := err.(JiraFlowError); ok {
		// Error title with icon
		title := fmt.Sprintf("❌ %s Error", jfErr.Type().String())
		titleStyle := h.tuiErrorStyle
		content.WriteString(titleStyle.Bold(true).Render(title))
		content.WriteString("\n\n")
		
		// User-friendly message
		messageStyle := h.tuiErrorStyle
		content.WriteString(messageStyle.Bold(false).Render(jfErr.UserMessage()))
		
		// Suggestions
		suggestions := jfErr.Suggestions()
		if len(suggestions) > 0 {
			content.WriteString("\n\n")
			content.WriteString(h.suggestionStyle.Render("💡 Suggestions:"))
			content.WriteString("\n")
			for _, suggestion := range suggestions {
				content.WriteString(h.suggestionStyle.Render("  • " + suggestion))
				content.WriteString("\n")
			}
		}
		
		// Recovery information
		if jfErr.IsRecoverable() {
			content.WriteString("\n")
			content.WriteString(h.warnStyle.Render("ℹ️  You can try again or use alternative options."))
		}
	} else {
		// Standard error formatting
		content.WriteString(h.tuiErrorStyle.Render("❌ " + err.Error()))
	}
	
	return content.String()
}

// FormatWarningForTUI formats a warning message for TUI display
func (h *ErrorHandler) FormatWarningForTUI(message string) string {
	return h.warnStyle.Render("⚠️  " + message)
}

// FormatSuccessForTUI formats a success message for TUI display
func (h *ErrorHandler) FormatSuccessForTUI(message string) string {
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)
	return successStyle.Render("✅ " + message)
}

// WrapError wraps a standard error as a JiraFlowError
func WrapError(err error, errorType ErrorType, recoverable bool) JiraFlowError {
	if err == nil {
		return nil
	}

	// If it's already a JiraFlowError, return as-is
	if jfErr, ok := err.(JiraFlowError); ok {
		return jfErr
	}

	// Wrap as appropriate error type
	switch errorType {
	case ErrorTypeConfig:
		return NewConfigError("", nil, err.Error(), recoverable)
	case ErrorTypeGit:
		return NewGitError("unknown", err.Error(), recoverable)
	case ErrorTypeJira:
		return NewJiraError("", err.Error(), recoverable)
	case ErrorTypeTUI:
		return NewTUIError("", err.Error(), recoverable)
	default:
		return NewGeneralError(err.Error(), recoverable)
	}
}

// IsRecoverableError checks if an error is recoverable
func IsRecoverableError(err error) bool {
	if jfErr, ok := err.(JiraFlowError); ok {
		return jfErr.IsRecoverable()
	}
	return false
}

// GetErrorType returns the error type if it's a JiraFlowError
func GetErrorType(err error) ErrorType {
	if jfErr, ok := err.(JiraFlowError); ok {
		return jfErr.Type()
	}
	return ErrorTypeGeneral
}

// HandleGracefulDegradation handles graceful degradation for missing dependencies
type DegradationHandler struct {
	errorHandler *ErrorHandler
}

// NewDegradationHandler creates a new DegradationHandler
func NewDegradationHandler() *DegradationHandler {
	return &DegradationHandler{
		errorHandler: NewErrorHandler(),
	}
}

// HandleJiraDegradation handles graceful degradation when Jira CLI is unavailable
func (d *DegradationHandler) HandleJiraDegradation(err error) string {
	if jiraErr, ok := err.(*JiraError); ok {
		if strings.Contains(jiraErr.Message, "not found") {
			return d.errorHandler.FormatWarningForTUI(
				"Jira CLI not available. You can still create branches by entering titles manually.",
			)
		}
		if strings.Contains(jiraErr.Message, "authentication") {
			return d.errorHandler.FormatWarningForTUI(
				"Jira authentication failed. Continuing with manual title entry.",
			)
		}
	}
	return d.errorHandler.FormatWarningForTUI(
		"Jira integration unavailable. Continuing with manual title entry.",
	)
}

// HandleGitDegradation handles graceful degradation for Git issues
func (d *DegradationHandler) HandleGitDegradation(err error) string {
	if gitErr, ok := err.(*GitError); ok {
		if strings.Contains(gitErr.Message, "not a git repository") {
			return d.errorHandler.FormatErrorForTUI(err)
		}
	}
	return d.errorHandler.FormatWarningForTUI(
		"Git operation failed. Some features may be limited.",
	)
}

// HandleConfigDegradation handles graceful degradation for config issues
func (d *DegradationHandler) HandleConfigDegradation(err error) string {
	return d.errorHandler.FormatWarningForTUI(
		"Configuration issue detected. Using default values.",
	)
}