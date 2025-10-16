package models

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// MockJiraClient for testing
type MockJiraClient struct {
	available bool
	titles    map[string]string
	errors    map[string]string
}

func (m *MockJiraClient) GetTicketTitle(ticketID string) (string, error) {
	if err, exists := m.errors[ticketID]; exists {
		return "", &MockJiraError{ticketID, err}
	}
	if title, exists := m.titles[ticketID]; exists {
		return title, nil
	}
	return "", &MockJiraError{ticketID, "ticket not found"}
}

func (m *MockJiraClient) IsAvailable() bool {
	return m.available
}

type MockJiraError struct {
	ticketID string
	message  string
}

func (e *MockJiraError) Error() string {
	return e.message
}

func TestNewInputFormModel(t *testing.T) {
	mockJira := &MockJiraClient{available: true}
	model := NewInputFormModel(mockJira)

	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected current field to be FieldTicketNumber, got %v", model.currentField)
	}

	if model.ticketValid {
		t.Error("Expected ticket to be invalid initially")
	}

	if model.completed {
		t.Error("Expected form to not be completed initially")
	}
}

func TestTicketValidation(t *testing.T) {
	model := NewInputFormModel(nil)

	tests := []struct {
		input    string
		expected bool
		errorMsg string
	}{
		{"", false, ""},
		{"JIRA-123", true, ""},
		{"PROJ-456", true, ""},
		{"ABC123-789", true, ""},
		{"jira-123", false, "Invalid format. Use PROJECT-123 format (e.g., JIRA-123)"},
		{"JIRA123", false, "Invalid format. Use PROJECT-123 format (e.g., JIRA-123)"},
		{"123-JIRA", false, "Invalid format. Use PROJECT-123 format (e.g., JIRA-123)"},
		{"JIRA-", false, "Invalid format. Use PROJECT-123 format (e.g., JIRA-123)"},
		{"-123", false, "Invalid format. Use PROJECT-123 format (e.g., JIRA-123)"},
	}

	for _, test := range tests {
		model.validateTicketNumber(test.input)
		if model.ticketValid != test.expected {
			t.Errorf("For input '%s', expected valid=%v, got %v", test.input, test.expected, model.ticketValid)
		}
		if test.errorMsg != "" && model.ticketError != test.errorMsg {
			t.Errorf("For input '%s', expected error '%s', got '%s'", test.input, test.errorMsg, model.ticketError)
		}
	}
}

func TestFieldNavigation(t *testing.T) {
	model := NewInputFormModel(nil)

	// Start on ticket field
	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected to start on FieldTicketNumber, got %v", model.currentField)
	}

	// Move to next field
	model.nextField()
	if model.currentField != FieldTitle {
		t.Errorf("Expected to move to FieldTitle, got %v", model.currentField)
	}

	// Try to move past last field (should stay on title)
	model.nextField()
	if model.currentField != FieldTitle {
		t.Errorf("Expected to stay on FieldTitle, got %v", model.currentField)
	}

	// Move back to previous field
	model.previousField()
	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected to move back to FieldTicketNumber, got %v", model.currentField)
	}

	// Try to move before first field (should stay on ticket)
	model.previousField()
	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected to stay on FieldTicketNumber, got %v", model.currentField)
	}
}

func TestFormValidation(t *testing.T) {
	model := NewInputFormModel(nil)

	// Form should be invalid initially
	if model.isFormValid() {
		t.Error("Expected form to be invalid initially")
	}

	// Set valid ticket number
	model.ticketInput.SetValue("JIRA-123")
	model.validateTicketNumber("JIRA-123")

	// Form should now be valid
	if !model.isFormValid() {
		t.Error("Expected form to be valid with valid ticket number")
	}

	// Clear ticket number
	model.ticketInput.SetValue("")
	model.validateTicketNumber("")

	// Form should be invalid again
	if model.isFormValid() {
		t.Error("Expected form to be invalid with empty ticket number")
	}
}

func TestKeyboardNavigation(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test Tab key navigation
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	if model.currentField != FieldTitle {
		t.Errorf("Expected Tab to move to FieldTitle, got %v", model.currentField)
	}

	// Test Shift+Tab navigation
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected Shift+Tab to move to FieldTicketNumber, got %v", model.currentField)
	}
}

func TestFormCompletion(t *testing.T) {
	model := NewInputFormModel(nil)

	// Set valid data
	model.ticketInput.SetValue("JIRA-123")
	model.titleInput.SetValue("Test Title")
	model.validateTicketNumber("JIRA-123")

	// Submit form
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !model.completed {
		t.Error("Expected form to be completed after Enter with valid data")
	}

	if model.GetTicketNumber() != "JIRA-123" {
		t.Errorf("Expected ticket number 'JIRA-123', got '%s'", model.GetTicketNumber())
	}

	if model.GetTicketTitle() != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", model.GetTicketTitle())
	}
}

func TestReset(t *testing.T) {
	model := NewInputFormModel(nil)

	// Set some data
	model.ticketInput.SetValue("JIRA-123")
	model.titleInput.SetValue("Test Title")
	model.validateTicketNumber("JIRA-123")
	model.completed = true

	// Reset
	model.Reset()

	if model.ticketInput.Value() != "" {
		t.Error("Expected ticket input to be empty after reset")
	}

	if model.titleInput.Value() != "" {
		t.Error("Expected title input to be empty after reset")
	}

	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected current field to be FieldTicketNumber after reset, got %v", model.currentField)
	}

	if model.completed {
		t.Error("Expected form to not be completed after reset")
	}

	if model.ticketValid {
		t.Error("Expected ticket to be invalid after reset")
	}
}