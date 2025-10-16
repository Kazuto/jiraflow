package models

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewCompletionModel(t *testing.T) {
	model := NewCompletionModel()
	
	if model.state != CompletionSuccess {
		t.Error("Expected initial state to be CompletionSuccess")
	}
	
	if model.shouldExit {
		t.Error("Expected shouldExit to be false initially")
	}
	
	if model.width != 0 || model.height != 0 {
		t.Error("Expected width and height to be 0 initially")
	}
}

func TestCompletionModel_SetSuccess(t *testing.T) {
	model := NewCompletionModel()
	
	branchName := "feature/JIRA-123-test-feature"
	baseBranch := "main"
	
	model.SetSuccess(branchName, baseBranch)
	
	if model.state != CompletionSuccess {
		t.Error("Expected state to be CompletionSuccess")
	}
	if model.branchName != branchName {
		t.Errorf("Expected branchName %s, got %s", branchName, model.branchName)
	}
	if model.baseBranch != baseBranch {
		t.Errorf("Expected baseBranch %s, got %s", baseBranch, model.baseBranch)
	}
	if model.errorMessage != "" {
		t.Error("Expected errorMessage to be empty for success state")
	}
}

func TestCompletionModel_SetError(t *testing.T) {
	model := NewCompletionModel()
	
	errorMessage := "Failed to create branch: permission denied"
	
	model.SetError(errorMessage)
	
	if model.state != CompletionError {
		t.Error("Expected state to be CompletionError")
	}
	if model.errorMessage != errorMessage {
		t.Errorf("Expected errorMessage %s, got %s", errorMessage, model.errorMessage)
	}
}

func TestCompletionModel_SetSize(t *testing.T) {
	model := NewCompletionModel()
	
	width, height := 80, 24
	model.SetSize(width, height)
	
	if model.width != width {
		t.Errorf("Expected width %d, got %d", width, model.width)
	}
	if model.height != height {
		t.Errorf("Expected height %d, got %d", height, model.height)
	}
}

func TestCompletionModel_Update(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectShouldExit bool
	}{
		{
			name:           "enter key exits",
			key:            "enter",
			expectShouldExit: true,
		},
		{
			name:           "q key exits",
			key:            "q",
			expectShouldExit: true,
		},
		{
			name:           "ctrl+c exits",
			key:            "ctrl+c",
			expectShouldExit: true,
		},
		{
			name:           "other keys don't exit",
			key:            "esc",
			expectShouldExit: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewCompletionModel()
			
			var keyMsg tea.KeyMsg
			switch tt.key {
			case "enter":
				keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
			case "q":
				keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
			case "ctrl+c":
				keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
			default:
				keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			}
			
			updatedModel, _ := model.Update(keyMsg)
			
			if updatedModel.ShouldExit() != tt.expectShouldExit {
				t.Errorf("Expected shouldExit to be %v, got %v", tt.expectShouldExit, updatedModel.ShouldExit())
			}
		})
	}
}

func TestCompletionModel_ViewSuccess(t *testing.T) {
	model := NewCompletionModel()
	model.SetSize(80, 24)
	model.SetSuccess("feature/JIRA-123-test-feature", "main")
	
	view := model.View()
	
	// Check that the view contains expected success elements
	expectedElements := []string{
		"✅",
		"Branch Created Successfully",
		"feature/JIRA-123-test-feature",
		"main",
		"Next Steps:",
		"git push -u origin",
		"enter exit application",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("Expected success view to contain '%s', but it didn't.\nView: %s", element, view)
		}
	}
}

func TestCompletionModel_ViewError(t *testing.T) {
	model := NewCompletionModel()
	model.SetSize(80, 24)
	model.SetError("Permission denied")
	
	view := model.View()
	
	// Check that the view contains expected error elements
	expectedElements := []string{
		"❌",
		"Branch Creation Failed",
		"Permission denied",
		"Troubleshooting Tips:",
		"permissions",
		"enter exit application",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("Expected error view to contain '%s', but it didn't.\nView: %s", element, view)
		}
	}
}

func TestCompletionModel_ShouldExit(t *testing.T) {
	model := NewCompletionModel()
	
	if model.ShouldExit() {
		t.Error("Expected ShouldExit to return false initially")
	}
	
	// Simulate exit key
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	
	if !updatedModel.ShouldExit() {
		t.Error("Expected ShouldExit to return true after enter key")
	}
}

func TestCompletionModel_Reset(t *testing.T) {
	model := NewCompletionModel()
	
	// Set some state
	model.SetError("Test error")
	model.SetSuccess("test-branch", "main")
	
	// Trigger exit
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	
	if !updatedModel.ShouldExit() {
		t.Error("Expected model to have shouldExit true before reset")
	}
	
	// Reset
	updatedModel.Reset()
	
	if updatedModel.ShouldExit() {
		t.Error("Expected shouldExit to be false after reset")
	}
	if updatedModel.state != CompletionSuccess {
		t.Error("Expected state to be CompletionSuccess after reset")
	}
	if updatedModel.branchName != "" {
		t.Error("Expected branchName to be empty after reset")
	}
	if updatedModel.baseBranch != "" {
		t.Error("Expected baseBranch to be empty after reset")
	}
	if updatedModel.errorMessage != "" {
		t.Error("Expected errorMessage to be empty after reset")
	}
}

func TestCompletionModel_ErrorWithBranchName(t *testing.T) {
	model := NewCompletionModel()
	model.SetSize(80, 24)
	
	// Set error with branch name for context
	model.SetError("Branch already exists")
	model.branchName = "feature/JIRA-123-duplicate"
	
	view := model.View()
	
	// Should show the attempted branch name in error view
	if !strings.Contains(view, "feature/JIRA-123-duplicate") {
		t.Error("Expected error view to show attempted branch name")
	}
	if !strings.Contains(view, "Attempted:") {
		t.Error("Expected error view to show 'Attempted:' label")
	}
}