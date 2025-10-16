package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorHandler_HandleError(t *testing.T) {
	handler := NewErrorHandler()

	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "nil error",
			err:      nil,
			wantCode: 0,
		},
		{
			name:     "config error",
			err:      NewConfigError("test", nil, "test error", true),
			wantCode: 2,
		},
		{
			name:     "git error",
			err:      NewGitError("branch", "test error", true),
			wantCode: 3,
		},
		{
			name:     "jira error",
			err:      NewJiraError("PROJ-123", "test error", true),
			wantCode: 4,
		},
		{
			name:     "tui error",
			err:      NewTUIError("component", "test error", true),
			wantCode: 5,
		},
		{
			name:     "general error",
			err:      NewGeneralError("test error", true),
			wantCode: 1,
		},
		{
			name:     "standard error",
			err:      errors.New("standard error"),
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := handler.HandleError(tt.err)
			if code != tt.wantCode {
				t.Errorf("ErrorHandler.HandleError() = %v, want %v", code, tt.wantCode)
			}
		})
	}
}

func TestErrorHandler_FormatErrorForTUI(t *testing.T) {
	handler := NewErrorHandler()

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
		{
			name: "jiraflow error",
			err:  NewConfigError("test", nil, "test error", true),
			want: "Configuration Error",
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: "standard error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.FormatErrorForTUI(tt.err)
			if tt.want != "" && !strings.Contains(result, tt.want) {
				t.Errorf("ErrorHandler.FormatErrorForTUI() = %v, want to contain %v", result, tt.want)
			}
			if tt.want == "" && result != "" {
				t.Errorf("ErrorHandler.FormatErrorForTUI() = %v, want empty string", result)
			}
		})
	}
}

func TestErrorHandler_FormatWarningForTUI(t *testing.T) {
	handler := NewErrorHandler()
	
	result := handler.FormatWarningForTUI("test warning")
	if !strings.Contains(result, "test warning") {
		t.Errorf("ErrorHandler.FormatWarningForTUI() = %v, want to contain 'test warning'", result)
	}
	if !strings.Contains(result, "⚠️") {
		t.Errorf("ErrorHandler.FormatWarningForTUI() = %v, want to contain warning icon", result)
	}
}

func TestErrorHandler_FormatSuccessForTUI(t *testing.T) {
	handler := NewErrorHandler()
	
	result := handler.FormatSuccessForTUI("test success")
	if !strings.Contains(result, "test success") {
		t.Errorf("ErrorHandler.FormatSuccessForTUI() = %v, want to contain 'test success'", result)
	}
	if !strings.Contains(result, "✅") {
		t.Errorf("ErrorHandler.FormatSuccessForTUI() = %v, want to contain success icon", result)
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		errorType   ErrorType
		recoverable bool
		wantType    ErrorType
	}{
		{
			name:        "nil error",
			err:         nil,
			errorType:   ErrorTypeConfig,
			recoverable: true,
			wantType:    ErrorTypeConfig, // Won't be used since err is nil
		},
		{
			name:        "already jiraflow error",
			err:         NewConfigError("test", nil, "test", true),
			errorType:   ErrorTypeGit, // Should be ignored
			recoverable: false,        // Should be ignored
			wantType:    ErrorTypeConfig,
		},
		{
			name:        "standard error wrapped as config",
			err:         errors.New("test error"),
			errorType:   ErrorTypeConfig,
			recoverable: true,
			wantType:    ErrorTypeConfig,
		},
		{
			name:        "standard error wrapped as git",
			err:         errors.New("test error"),
			errorType:   ErrorTypeGit,
			recoverable: false,
			wantType:    ErrorTypeGit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.errorType, tt.recoverable)
			
			if tt.err == nil {
				if result != nil {
					t.Errorf("WrapError() = %v, want nil for nil input", result)
				}
				return
			}
			
			if result == nil {
				t.Error("WrapError() = nil, want non-nil for non-nil input")
				return
			}
			
			if result.Type() != tt.wantType {
				t.Errorf("WrapError().Type() = %v, want %v", result.Type(), tt.wantType)
			}
		})
	}
}

func TestIsRecoverableError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "recoverable jiraflow error",
			err:  NewConfigError("test", nil, "test", true),
			want: true,
		},
		{
			name: "non-recoverable jiraflow error",
			err:  NewConfigError("test", nil, "test", false),
			want: false,
		},
		{
			name: "standard error",
			err:  errors.New("test"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRecoverableError(tt.err); got != tt.want {
				t.Errorf("IsRecoverableError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetErrorType(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrorType
	}{
		{
			name: "config error",
			err:  NewConfigError("test", nil, "test", true),
			want: ErrorTypeConfig,
		},
		{
			name: "git error",
			err:  NewGitError("branch", "test", true),
			want: ErrorTypeGit,
		},
		{
			name: "standard error",
			err:  errors.New("test"),
			want: ErrorTypeGeneral,
		},
		{
			name: "nil error",
			err:  nil,
			want: ErrorTypeGeneral,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetErrorType(tt.err); got != tt.want {
				t.Errorf("GetErrorType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDegradationHandler(t *testing.T) {
	handler := NewDegradationHandler()

	t.Run("HandleJiraDegradation", func(t *testing.T) {
		err := NewJiraError("PROJ-123", "jira CLI not found", true)
		result := handler.HandleJiraDegradation(err)
		
		if !strings.Contains(result, "Jira") {
			t.Errorf("HandleJiraDegradation() = %v, want to contain 'Jira'", result)
		}
	})

	t.Run("HandleGitDegradation", func(t *testing.T) {
		err := NewGitError("branch", "not a git repository", false)
		result := handler.HandleGitDegradation(err)
		
		if result == "" {
			t.Error("HandleGitDegradation() returned empty string")
		}
	})

	t.Run("HandleConfigDegradation", func(t *testing.T) {
		err := NewConfigError("test", nil, "test error", true)
		result := handler.HandleConfigDegradation(err)
		
		if !strings.Contains(result, "Configuration") && !strings.Contains(result, "default") {
			t.Errorf("HandleConfigDegradation() = %v, want to contain 'Configuration' or 'default'", result)
		}
	})
}