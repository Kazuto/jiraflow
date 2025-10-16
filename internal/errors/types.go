package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the category of error
type ErrorType int

const (
	ErrorTypeConfig ErrorType = iota
	ErrorTypeGit
	ErrorTypeJira
	ErrorTypeTUI
	ErrorTypeGeneral
)

// String returns the string representation of ErrorType
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeConfig:
		return "Configuration"
	case ErrorTypeGit:
		return "Git"
	case ErrorTypeJira:
		return "Jira"
	case ErrorTypeTUI:
		return "Interface"
	case ErrorTypeGeneral:
		return "General"
	default:
		return "Unknown"
	}
}

// JiraFlowError is the base error interface for all application errors
type JiraFlowError interface {
	error
	Type() ErrorType
	UserMessage() string
	Suggestions() []string
	IsRecoverable() bool
}

// ConfigError represents configuration-related errors
type ConfigError struct {
	Field       string
	Value       interface{}
	Message     string
	Recoverable bool
}

func (e ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("configuration error in field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

func (e ConfigError) Type() ErrorType {
	return ErrorTypeConfig
}

func (e ConfigError) UserMessage() string {
	if e.Field != "" {
		return fmt.Sprintf("Configuration issue with '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("Configuration issue: %s", e.Message)
}

func (e ConfigError) Suggestions() []string {
	suggestions := []string{}
	
	switch e.Field {
	case "max_branch_length":
		suggestions = append(suggestions, "Set max_branch_length to a value between 10 and 200")
		suggestions = append(suggestions, "Edit ~/.config/jiraflow/jiraflow.yaml to fix this setting")
	case "branch_types":
		suggestions = append(suggestions, "Ensure all branch types have non-empty string values")
		suggestions = append(suggestions, "Check the branch_types section in your config file")
	case "sanitization.separator":
		suggestions = append(suggestions, "Use a single character or short string for separator")
		suggestions = append(suggestions, "Common separators: '-', '_', '.'")
	default:
		suggestions = append(suggestions, "Check your configuration file at ~/.config/jiraflow/jiraflow.yaml")
		suggestions = append(suggestions, "Delete the config file to regenerate with defaults")
	}
	
	return suggestions
}

func (e ConfigError) IsRecoverable() bool {
	return e.Recoverable
}

// NewConfigError creates a new ConfigError
func NewConfigError(field string, value interface{}, message string, recoverable bool) *ConfigError {
	return &ConfigError{
		Field:       field,
		Value:       value,
		Message:     message,
		Recoverable: recoverable,
	}
}

// GitError represents Git operation errors
type GitError struct {
	Operation   string
	Message     string
	Recoverable bool
}

func (e GitError) Error() string {
	return fmt.Sprintf("git %s: %s", e.Operation, e.Message)
}

func (e GitError) Type() ErrorType {
	return ErrorTypeGit
}

func (e GitError) UserMessage() string {
	switch e.Operation {
	case "branch":
		if strings.Contains(e.Message, "not a git repository") {
			return "This directory is not a Git repository"
		}
		if strings.Contains(e.Message, "already exists") {
			return "A branch with this name already exists"
		}
		return fmt.Sprintf("Branch operation failed: %s", e.Message)
	case "checkout":
		return fmt.Sprintf("Failed to switch branches: %s", e.Message)
	default:
		return fmt.Sprintf("Git operation failed: %s", e.Message)
	}
}

func (e GitError) Suggestions() []string {
	suggestions := []string{}
	
	switch e.Operation {
	case "branch":
		if strings.Contains(e.Message, "not a git repository") {
			suggestions = append(suggestions, "Navigate to a Git repository directory")
			suggestions = append(suggestions, "Initialize a Git repository with 'git init'")
		} else if strings.Contains(e.Message, "already exists") {
			suggestions = append(suggestions, "Use a different ticket number or title")
			suggestions = append(suggestions, "Delete the existing branch if it's no longer needed")
		} else if strings.Contains(e.Message, "does not exist") {
			suggestions = append(suggestions, "Check that the base branch exists locally")
			suggestions = append(suggestions, "Fetch the latest branches with 'git fetch'")
		} else {
			suggestions = append(suggestions, "Ensure you have proper Git permissions")
			suggestions = append(suggestions, "Check that the base branch exists and is accessible")
		}
	case "checkout":
		suggestions = append(suggestions, "Ensure the branch exists")
		suggestions = append(suggestions, "Commit or stash any uncommitted changes")
	default:
		suggestions = append(suggestions, "Check your Git repository status")
		suggestions = append(suggestions, "Ensure you have proper Git permissions")
	}
	
	return suggestions
}

func (e GitError) IsRecoverable() bool {
	return e.Recoverable
}

// NewGitError creates a new GitError
func NewGitError(operation, message string, recoverable bool) *GitError {
	return &GitError{
		Operation:   operation,
		Message:     message,
		Recoverable: recoverable,
	}
}

// JiraError represents Jira integration errors
type JiraError struct {
	TicketID    string
	Message     string
	Recoverable bool
}

func (e JiraError) Error() string {
	if e.TicketID != "" {
		return fmt.Sprintf("jira error for ticket %s: %s", e.TicketID, e.Message)
	}
	return fmt.Sprintf("jira error: %s", e.Message)
}

func (e JiraError) Type() ErrorType {
	return ErrorTypeJira
}

func (e JiraError) UserMessage() string {
	if strings.Contains(e.Message, "not found") {
		return "Jira CLI is not installed or not in PATH"
	}
	if strings.Contains(e.Message, "authentication") {
		return "Jira authentication failed"
	}
	if strings.Contains(e.Message, "not found") && e.TicketID != "" {
		return fmt.Sprintf("Ticket %s was not found", e.TicketID)
	}
	return fmt.Sprintf("Jira integration issue: %s", e.Message)
}

func (e JiraError) Suggestions() []string {
	suggestions := []string{}
	
	if strings.Contains(e.Message, "not found") && e.TicketID == "" {
		suggestions = append(suggestions, "Install Jira CLI from https://github.com/ankitpokhrel/jira-cli")
		suggestions = append(suggestions, "Ensure 'jira' command is in your PATH")
		suggestions = append(suggestions, "You can still use JiraFlow by entering ticket titles manually")
	} else if strings.Contains(e.Message, "authentication") {
		suggestions = append(suggestions, "Run 'jira init' to configure your Jira credentials")
		suggestions = append(suggestions, "Check your Jira server URL and credentials")
	} else if strings.Contains(e.Message, "not found") && e.TicketID != "" {
		suggestions = append(suggestions, fmt.Sprintf("Verify that ticket %s exists in your Jira instance", e.TicketID))
		suggestions = append(suggestions, "Check the ticket ID format (e.g., PROJ-123)")
		suggestions = append(suggestions, "You can proceed by entering the title manually")
	} else {
		suggestions = append(suggestions, "Check your Jira CLI configuration")
		suggestions = append(suggestions, "You can continue without Jira integration")
	}
	
	return suggestions
}

func (e JiraError) IsRecoverable() bool {
	return e.Recoverable
}

// NewJiraError creates a new JiraError
func NewJiraError(ticketID, message string, recoverable bool) *JiraError {
	return &JiraError{
		TicketID:    ticketID,
		Message:     message,
		Recoverable: recoverable,
	}
}

// TUIError represents TUI-related errors
type TUIError struct {
	Component   string
	Message     string
	Recoverable bool
}

func (e TUIError) Error() string {
	if e.Component != "" {
		return fmt.Sprintf("TUI error in %s: %s", e.Component, e.Message)
	}
	return fmt.Sprintf("TUI error: %s", e.Message)
}

func (e TUIError) Type() ErrorType {
	return ErrorTypeTUI
}

func (e TUIError) UserMessage() string {
	return fmt.Sprintf("Interface issue: %s", e.Message)
}

func (e TUIError) Suggestions() []string {
	suggestions := []string{}
	
	if strings.Contains(e.Message, "terminal") {
		suggestions = append(suggestions, "Ensure your terminal supports the required features")
		suggestions = append(suggestions, "Try resizing your terminal window")
	} else {
		suggestions = append(suggestions, "Try restarting the application")
		suggestions = append(suggestions, "Use non-interactive mode as a fallback")
	}
	
	return suggestions
}

func (e TUIError) IsRecoverable() bool {
	return e.Recoverable
}

// NewTUIError creates a new TUIError
func NewTUIError(component, message string, recoverable bool) *TUIError {
	return &TUIError{
		Component:   component,
		Message:     message,
		Recoverable: recoverable,
	}
}

// GeneralError represents general application errors
type GeneralError struct {
	Message     string
	Recoverable bool
}

func (e GeneralError) Error() string {
	return e.Message
}

func (e GeneralError) Type() ErrorType {
	return ErrorTypeGeneral
}

func (e GeneralError) UserMessage() string {
	return e.Message
}

func (e GeneralError) Suggestions() []string {
	return []string{
		"Try running the command again",
		"Check the application logs for more details",
	}
}

func (e GeneralError) IsRecoverable() bool {
	return e.Recoverable
}

// NewGeneralError creates a new GeneralError
func NewGeneralError(message string, recoverable bool) *GeneralError {
	return &GeneralError{
		Message:     message,
		Recoverable: recoverable,
	}
}