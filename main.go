package main

import (
	"fmt"
	"os"

	"jiraflow/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// Check for specific error types to provide appropriate exit codes
		switch {
		case isConfigError(err):
			fmt.Fprintf(os.Stderr, "Configuration Error: %v\n", err)
			os.Exit(2)
		case isGitError(err):
			fmt.Fprintf(os.Stderr, "Git Error: %v\n", err)
			os.Exit(3)
		case isJiraError(err):
			fmt.Fprintf(os.Stderr, "Jira Error: %v\n", err)
			os.Exit(4)
		case isUserCancelled(err):
			// User cancelled operation, exit gracefully
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

// isConfigError checks if the error is related to configuration
func isConfigError(err error) bool {
	return err != nil && (
		containsString(err.Error(), "configuration") ||
		containsString(err.Error(), "config") ||
		containsString(err.Error(), "yaml"))
}

// isGitError checks if the error is related to Git operations
func isGitError(err error) bool {
	return err != nil && (
		containsString(err.Error(), "git") ||
		containsString(err.Error(), "repository") ||
		containsString(err.Error(), "branch"))
}

// isJiraError checks if the error is related to Jira operations
func isJiraError(err error) bool {
	return err != nil && (
		containsString(err.Error(), "jira") ||
		containsString(err.Error(), "ticket"))
}

// isUserCancelled checks if the user cancelled the operation
func isUserCancelled(err error) bool {
	return err != nil && (
		containsString(err.Error(), "cancelled") ||
		containsString(err.Error(), "interrupted") ||
		containsString(err.Error(), "user quit"))
}

// containsString checks if a string contains a substring (case-insensitive)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      findSubstring(s, substr))))
}

// findSubstring performs a simple substring search
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}