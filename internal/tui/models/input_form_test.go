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

func TestInputFormModel_KeyboardNavigation(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test arrow key navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ = model.Update(downMsg)
	if model.currentField != FieldTitle {
		t.Errorf("Expected Down arrow to move to FieldTitle, got %v", model.currentField)
	}

	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ = model.Update(upMsg)
	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected Up arrow to move to FieldTicketNumber, got %v", model.currentField)
	}

	// Test vim-style navigation
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	model, _ = model.Update(jMsg)
	if model.currentField != FieldTitle {
		t.Errorf("Expected 'j' to move to FieldTitle, got %v", model.currentField)
	}

	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	model, _ = model.Update(kMsg)
	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected 'k' to move to FieldTicketNumber, got %v", model.currentField)
	}
}

func TestInputFormModel_TextInput(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test typing in ticket field
	model.currentField = FieldTicketNumber
	
	// Type "JIRA-123"
	chars := []rune("JIRA-123")
	for _, char := range chars {
		charMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
		model, _ = model.Update(charMsg)
	}

	if model.ticketInput.Value() != "JIRA-123" {
		t.Errorf("Expected ticket input to be 'JIRA-123', got '%s'", model.ticketInput.Value())
	}

	if !model.ticketValid {
		t.Error("Expected ticket to be valid after typing valid format")
	}

	// Test backspace
	backspaceMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	model, _ = model.Update(backspaceMsg)
	model, _ = model.Update(backspaceMsg)
	model, _ = model.Update(backspaceMsg)

	// After 3 backspaces from "JIRA-123", we should have "JIRA-"
	expectedValue := "JIRA-"
	if model.ticketInput.Value() != expectedValue {
		t.Errorf("Expected ticket input to be '%s' after backspace, got '%s'", expectedValue, model.ticketInput.Value())
	}
}

func TestInputFormModel_FieldFocusManagement(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test initial focus
	if !model.ticketInput.Focused() {
		t.Error("Expected ticket input to be focused initially")
	}

	if model.titleInput.Focused() {
		t.Error("Expected title input to not be focused initially")
	}

	// Test focus change
	model.FocusTitleField()
	if model.currentField != FieldTitle {
		t.Errorf("Expected current field to be FieldTitle, got %v", model.currentField)
	}

	if model.ticketInput.Focused() {
		t.Error("Expected ticket input to not be focused after focus change")
	}

	if !model.titleInput.Focused() {
		t.Error("Expected title input to be focused after focus change")
	}

	// Test focus back to ticket field
	model.FocusTicketField()
	if model.currentField != FieldTicketNumber {
		t.Errorf("Expected current field to be FieldTicketNumber, got %v", model.currentField)
	}

	if !model.ticketInput.Focused() {
		t.Error("Expected ticket input to be focused after focus change")
	}

	if model.titleInput.Focused() {
		t.Error("Expected title input to not be focused after focus change")
	}
}

func TestInputFormModel_WindowSizeHandling(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test window size update
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	model, _ = model.Update(sizeMsg)

	if model.width != 120 {
		t.Errorf("Expected width to be 120, got %d", model.width)
	}

	if model.height != 40 {
		t.Errorf("Expected height to be 40, got %d", model.height)
	}

	// Test that input widths are updated appropriately
	expectedTicketWidth := min(30, 120-10)
	expectedTitleWidth := min(50, 120-10)

	if model.ticketInput.Width != expectedTicketWidth {
		t.Errorf("Expected ticket input width to be %d, got %d", expectedTicketWidth, model.ticketInput.Width)
	}

	if model.titleInput.Width != expectedTitleWidth {
		t.Errorf("Expected title input width to be %d, got %d", expectedTitleWidth, model.titleInput.Width)
	}
}

func TestInputFormModel_JiraIntegration(t *testing.T) {
	// Test with available Jira client
	mockJira := &MockJiraClient{
		available: true,
		titles: map[string]string{
			"JIRA-123": "Test Feature Implementation",
			"PROJ-456": "Bug Fix for Login Issue",
		},
	}

	model := NewInputFormModel(mockJira)

	// Set valid ticket number
	model.ticketInput.SetValue("JIRA-123")
	model.validateTicketNumber("JIRA-123")

	// Simulate title fetch message
	fetchMsg := FetchTitleMsg{
		TicketID: "JIRA-123",
		Title:    "Test Feature Implementation",
		Error:    "",
	}

	model, _ = model.Update(fetchMsg)

	if model.titleInput.Value() != "Test Feature Implementation" {
		t.Errorf("Expected title to be fetched and set, got '%s'", model.titleInput.Value())
	}

	if !model.titleFetched {
		t.Error("Expected titleFetched to be true after successful fetch")
	}

	if model.titleError != "" {
		t.Errorf("Expected no title error after successful fetch, got '%s'", model.titleError)
	}
}

func TestInputFormModel_JiraErrorHandling(t *testing.T) {
	// Test with Jira client that returns errors
	mockJira := &MockJiraClient{
		available: true,
		errors: map[string]string{
			"JIRA-404": "ticket not found",
		},
	}

	model := NewInputFormModel(mockJira)

	// Simulate title fetch error
	fetchMsg := FetchTitleMsg{
		TicketID: "JIRA-404",
		Title:    "",
		Error:    "ticket not found",
	}

	model, _ = model.Update(fetchMsg)

	if model.titleFetched {
		t.Error("Expected titleFetched to be false after error")
	}

	if model.titleError != "ticket not found" {
		t.Errorf("Expected title error to be 'ticket not found', got '%s'", model.titleError)
	}

	if model.titleInput.Value() != "" {
		t.Error("Expected title input to remain empty after fetch error")
	}
}

func TestInputFormModel_PreFillingData(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test pre-filling ticket number
	model.SetTicketNumber("PROJ-789")
	if model.ticketInput.Value() != "PROJ-789" {
		t.Errorf("Expected ticket input to be 'PROJ-789', got '%s'", model.ticketInput.Value())
	}

	if !model.ticketValid {
		t.Error("Expected ticket to be valid after setting valid ticket number")
	}

	// Test pre-filling title
	model.SetTitle("Pre-filled Title")
	if model.titleInput.Value() != "Pre-filled Title" {
		t.Errorf("Expected title input to be 'Pre-filled Title', got '%s'", model.titleInput.Value())
	}
}

func TestInputFormModel_StateQueries(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test initial state queries
	if model.HasCompleted() {
		t.Error("Expected form to not be completed initially")
	}

	if model.IsValid() {
		t.Error("Expected form to not be valid initially")
	}

	if model.GetCurrentField() != FieldTicketNumber {
		t.Errorf("Expected current field to be FieldTicketNumber, got %v", model.GetCurrentField())
	}

	// Set valid data and complete form
	model.ticketInput.SetValue("JIRA-123")
	model.validateTicketNumber("JIRA-123")
	model.completed = true

	if !model.HasCompleted() {
		t.Error("Expected form to be completed after setting completed flag")
	}

	if !model.IsValid() {
		t.Error("Expected form to be valid with valid ticket number")
	}
}

func TestInputFormModel_ViewRendering(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test basic view rendering
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Test view with valid data
	model.ticketInput.SetValue("JIRA-123")
	model.validateTicketNumber("JIRA-123")
	
	viewWithData := model.View()
	if viewWithData == "" {
		t.Error("Expected non-empty view with valid data")
	}

	// Test view with errors
	model.ticketInput.SetValue("invalid")
	model.validateTicketNumber("invalid")
	
	viewWithError := model.View()
	if viewWithError == "" {
		t.Error("Expected non-empty view with validation error")
	}
}

func TestInputFormModel_ComplexValidationScenarios(t *testing.T) {
	model := NewInputFormModel(nil)

	// Test edge cases for ticket validation
	edgeCases := []struct {
		input       string
		shouldBeValid bool
		description string
	}{
		{"A-1", true, "minimal valid format"},
		{"ABC123DEF-999", true, "alphanumeric project code"},
		{"PROJECT-0", true, "zero ticket number"},
		{"X-1234567890", true, "long ticket number"},
		{"a-123", false, "lowercase project code"},
		{"PROJECT-", false, "missing ticket number"},
		{"PROJECT", false, "missing dash and number"},
		{"-123", false, "missing project code"},
		{"", false, "empty input"},
		{"PROJECT-ABC", false, "non-numeric ticket number"},
	}

	for _, tc := range edgeCases {
		model.validateTicketNumber(tc.input)
		if model.ticketValid != tc.shouldBeValid {
			t.Errorf("For input '%s' (%s): expected valid=%v, got %v", 
				tc.input, tc.description, tc.shouldBeValid, model.ticketValid)
		}
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