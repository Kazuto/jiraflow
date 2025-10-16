package models

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewConfirmationModel(t *testing.T) {
	model := NewConfirmationModel()
	
	if model.confirmed {
		t.Error("Expected confirmed to be false initially")
	}
	
	if model.width != 0 || model.height != 0 {
		t.Error("Expected width and height to be 0 initially")
	}
}

func TestConfirmationModel_SetData(t *testing.T) {
	model := NewConfirmationModel()
	
	branchType := "feature"
	baseBranch := "main"
	ticketNumber := "JIRA-123"
	ticketTitle := "Test Feature"
	finalBranch := "feature/JIRA-123-test-feature"
	
	model.SetData(branchType, baseBranch, ticketNumber, ticketTitle, finalBranch)
	
	if model.branchType != branchType {
		t.Errorf("Expected branchType %s, got %s", branchType, model.branchType)
	}
	if model.baseBranch != baseBranch {
		t.Errorf("Expected baseBranch %s, got %s", baseBranch, model.baseBranch)
	}
	if model.ticketNumber != ticketNumber {
		t.Errorf("Expected ticketNumber %s, got %s", ticketNumber, model.ticketNumber)
	}
	if model.ticketTitle != ticketTitle {
		t.Errorf("Expected ticketTitle %s, got %s", ticketTitle, model.ticketTitle)
	}
	if model.finalBranch != finalBranch {
		t.Errorf("Expected finalBranch %s, got %s", finalBranch, model.finalBranch)
	}
}

func TestConfirmationModel_SetSize(t *testing.T) {
	model := NewConfirmationModel()
	
	width, height := 80, 24
	model.SetSize(width, height)
	
	if model.width != width {
		t.Errorf("Expected width %d, got %d", width, model.width)
	}
	if model.height != height {
		t.Errorf("Expected height %d, got %d", height, model.height)
	}
}

func TestConfirmationModel_Update(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		expectConfirmed bool
	}{
		{
			name:          "enter key confirms",
			key:           "enter",
			expectConfirmed: true,
		},
		{
			name:          "other keys don't confirm",
			key:           "esc",
			expectConfirmed: false,
		},
		{
			name:          "space key doesn't confirm",
			key:           " ",
			expectConfirmed: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewConfirmationModel()
			
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "enter" {
				keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
			}
			
			updatedModel, _ := model.Update(keyMsg)
			
			if updatedModel.HasConfirmed() != tt.expectConfirmed {
				t.Errorf("Expected confirmed to be %v, got %v", tt.expectConfirmed, updatedModel.HasConfirmed())
			}
		})
	}
}

func TestConfirmationModel_View(t *testing.T) {
	model := NewConfirmationModel()
	model.SetSize(80, 24)
	model.SetData("feature", "main", "JIRA-123", "Test Feature", "feature/JIRA-123-test-feature")
	
	view := model.View()
	
	// Check that the view contains expected elements
	expectedElements := []string{
		"Confirm Branch Creation",
		"Summary:",
		"feature",
		"main",
		"JIRA-123",
		"Test Feature",
		"feature/JIRA-123-test-feature",
		"enter create branch",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("Expected view to contain '%s', but it didn't.\nView: %s", element, view)
		}
	}
}

func TestConfirmationModel_ViewWithoutTitle(t *testing.T) {
	model := NewConfirmationModel()
	model.SetSize(80, 24)
	model.SetData("hotfix", "develop", "BUG-456", "", "hotfix/BUG-456")
	
	view := model.View()
	
	// Check that the view contains expected elements but not empty title
	expectedElements := []string{
		"hotfix",
		"develop", 
		"BUG-456",
		"hotfix/BUG-456",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(view, element) {
			t.Errorf("Expected view to contain '%s', but it didn't", element)
		}
	}
	
	// Should not show title section when title is empty
	if strings.Contains(view, "Title:") {
		t.Error("Expected view to not contain 'Title:' when title is empty")
	}
}

func TestConfirmationModel_HasConfirmed(t *testing.T) {
	model := NewConfirmationModel()
	
	if model.HasConfirmed() {
		t.Error("Expected HasConfirmed to return false initially")
	}
	
	// Simulate confirmation
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	
	if !updatedModel.HasConfirmed() {
		t.Error("Expected HasConfirmed to return true after enter key")
	}
}

func TestConfirmationModel_Reset(t *testing.T) {
	model := NewConfirmationModel()
	
	// Confirm first
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(keyMsg)
	
	if !updatedModel.HasConfirmed() {
		t.Error("Expected model to be confirmed before reset")
	}
	
	// Reset
	updatedModel.Reset()
	
	if updatedModel.HasConfirmed() {
		t.Error("Expected model to not be confirmed after reset")
	}
}

func TestConfirmationModel_GetFinalBranch(t *testing.T) {
	model := NewConfirmationModel()
	finalBranch := "feature/JIRA-123-test-feature"
	
	model.SetData("feature", "main", "JIRA-123", "Test Feature", finalBranch)
	
	if model.GetFinalBranch() != finalBranch {
		t.Errorf("Expected GetFinalBranch to return %s, got %s", finalBranch, model.GetFinalBranch())
	}
}