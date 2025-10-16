package errors

import (
	"testing"
)

func TestConfigError(t *testing.T) {
	tests := []struct {
		name        string
		field       string
		value       interface{}
		message     string
		recoverable bool
		wantType    ErrorType
		wantError   string
	}{
		{
			name:        "config error with field",
			field:       "max_branch_length",
			value:       5,
			message:     "value too small",
			recoverable: true,
			wantType:    ErrorTypeConfig,
			wantError:   "configuration error in field 'max_branch_length': value too small",
		},
		{
			name:        "config error without field",
			field:       "",
			value:       nil,
			message:     "general config error",
			recoverable: false,
			wantType:    ErrorTypeConfig,
			wantError:   "configuration error: general config error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewConfigError(tt.field, tt.value, tt.message, tt.recoverable)
			
			if err.Type() != tt.wantType {
				t.Errorf("ConfigError.Type() = %v, want %v", err.Type(), tt.wantType)
			}
			
			if err.Error() != tt.wantError {
				t.Errorf("ConfigError.Error() = %v, want %v", err.Error(), tt.wantError)
			}
			
			if err.IsRecoverable() != tt.recoverable {
				t.Errorf("ConfigError.IsRecoverable() = %v, want %v", err.IsRecoverable(), tt.recoverable)
			}
			
			// Test suggestions
			suggestions := err.Suggestions()
			if len(suggestions) == 0 {
				t.Error("ConfigError.Suggestions() returned empty slice")
			}
		})
	}
}

func TestGitError(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		message     string
		recoverable bool
		wantType    ErrorType
	}{
		{
			name:        "git branch error",
			operation:   "branch",
			message:     "not a git repository",
			recoverable: false,
			wantType:    ErrorTypeGit,
		},
		{
			name:        "git checkout error",
			operation:   "checkout",
			message:     "branch does not exist",
			recoverable: true,
			wantType:    ErrorTypeGit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewGitError(tt.operation, tt.message, tt.recoverable)
			
			if err.Type() != tt.wantType {
				t.Errorf("GitError.Type() = %v, want %v", err.Type(), tt.wantType)
			}
			
			if err.IsRecoverable() != tt.recoverable {
				t.Errorf("GitError.IsRecoverable() = %v, want %v", err.IsRecoverable(), tt.recoverable)
			}
			
			// Test user message
			userMsg := err.UserMessage()
			if userMsg == "" {
				t.Error("GitError.UserMessage() returned empty string")
			}
			
			// Test suggestions
			suggestions := err.Suggestions()
			if len(suggestions) == 0 {
				t.Error("GitError.Suggestions() returned empty slice")
			}
		})
	}
}

func TestJiraError(t *testing.T) {
	tests := []struct {
		name        string
		ticketID    string
		message     string
		recoverable bool
		wantType    ErrorType
	}{
		{
			name:        "jira ticket not found",
			ticketID:    "PROJ-123",
			message:     "ticket not found",
			recoverable: true,
			wantType:    ErrorTypeJira,
		},
		{
			name:        "jira authentication error",
			ticketID:    "PROJ-456",
			message:     "authentication failed",
			recoverable: true,
			wantType:    ErrorTypeJira,
		},
		{
			name:        "jira cli not found",
			ticketID:    "",
			message:     "jira CLI not found",
			recoverable: true,
			wantType:    ErrorTypeJira,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewJiraError(tt.ticketID, tt.message, tt.recoverable)
			
			if err.Type() != tt.wantType {
				t.Errorf("JiraError.Type() = %v, want %v", err.Type(), tt.wantType)
			}
			
			if err.IsRecoverable() != tt.recoverable {
				t.Errorf("JiraError.IsRecoverable() = %v, want %v", err.IsRecoverable(), tt.recoverable)
			}
			
			// Test user message
			userMsg := err.UserMessage()
			if userMsg == "" {
				t.Error("JiraError.UserMessage() returned empty string")
			}
			
			// Test suggestions
			suggestions := err.Suggestions()
			if len(suggestions) == 0 {
				t.Error("JiraError.Suggestions() returned empty slice")
			}
		})
	}
}

func TestTUIError(t *testing.T) {
	err := NewTUIError("branch_selector", "component failed to render", true)
	
	if err.Type() != ErrorTypeTUI {
		t.Errorf("TUIError.Type() = %v, want %v", err.Type(), ErrorTypeTUI)
	}
	
	if !err.IsRecoverable() {
		t.Error("TUIError.IsRecoverable() = false, want true")
	}
	
	userMsg := err.UserMessage()
	if userMsg == "" {
		t.Error("TUIError.UserMessage() returned empty string")
	}
	
	suggestions := err.Suggestions()
	if len(suggestions) == 0 {
		t.Error("TUIError.Suggestions() returned empty slice")
	}
}

func TestGeneralError(t *testing.T) {
	err := NewGeneralError("something went wrong", false)
	
	if err.Type() != ErrorTypeGeneral {
		t.Errorf("GeneralError.Type() = %v, want %v", err.Type(), ErrorTypeGeneral)
	}
	
	if err.IsRecoverable() {
		t.Error("GeneralError.IsRecoverable() = true, want false")
	}
	
	userMsg := err.UserMessage()
	if userMsg != "something went wrong" {
		t.Errorf("GeneralError.UserMessage() = %v, want %v", userMsg, "something went wrong")
	}
}

func TestErrorTypeString(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		want      string
	}{
		{ErrorTypeConfig, "Configuration"},
		{ErrorTypeGit, "Git"},
		{ErrorTypeJira, "Jira"},
		{ErrorTypeTUI, "Interface"},
		{ErrorTypeGeneral, "General"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.errorType.String(); got != tt.want {
				t.Errorf("ErrorType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}