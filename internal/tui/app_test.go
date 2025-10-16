package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"jiraflow/internal/config"
	"jiraflow/internal/git"
)

// MockGitRepository for testing
type MockGitRepository struct {
	branches      []git.BranchInfo
	currentBranch string
	createError   error
	checkoutError error
}

func (m *MockGitRepository) GetBranchesWithInfo() ([]git.BranchInfo, error) {
	return m.branches, nil
}

func (m *MockGitRepository) GetLocalBranches() ([]string, error) {
	var names []string
	for _, branch := range m.branches {
		if !branch.IsRemote {
			names = append(names, branch.Name)
		}
	}
	return names, nil
}

func (m *MockGitRepository) GetCurrentBranch() (string, error) {
	return m.currentBranch, nil
}

func (m *MockGitRepository) CreateBranch(name, baseBranch string) error {
	return m.createError
}

func (m *MockGitRepository) CheckoutBranch(name string) error {
	return m.checkoutError
}

func (m *MockGitRepository) IsGitRepository() bool {
	return true
}

func (m *MockGitRepository) SearchBranches(searchTerm string) (git.BranchSearchResult, error) {
	branches, err := m.GetLocalBranches()
	if err != nil {
		return git.BranchSearchResult{}, err
	}
	return git.FilterBranchesRealtime(branches, searchTerm), nil
}

func TestNewAppModel(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
			{Name: "develop", IsCurrent: false, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	// Test initial state
	if model.state != StateTypeSelection {
		t.Errorf("Expected initial state to be StateTypeSelection, got %v", model.state)
	}

	if model.config != cfg {
		t.Error("Expected config to be set correctly")
	}

	// We can't directly compare interfaces, so just check it's not nil
	if model.git == nil {
		t.Error("Expected git repository to be set")
	}
}

func TestAppModel_StateTransitions(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	// Test transition from TypeSelection to BranchSelection
	model.selectedType = "feature"
	model.SetState(StateBranchSelection)
	if model.GetCurrentState() != StateBranchSelection {
		t.Errorf("Expected state to be StateBranchSelection, got %v", model.GetCurrentState())
	}

	// Test transition from BranchSelection to TicketInput
	model.selectedBranch = "main"
	model.SetState(StateTicketInput)
	if model.GetCurrentState() != StateTicketInput {
		t.Errorf("Expected state to be StateTicketInput, got %v", model.GetCurrentState())
	}

	// Test transition from TicketInput to Confirmation
	model.ticketNumber = "JIRA-123"
	model.ticketTitle = "Test Title"
	model.SetState(StateConfirmation)
	if model.GetCurrentState() != StateConfirmation {
		t.Errorf("Expected state to be StateConfirmation, got %v", model.GetCurrentState())
	}

	// Test transition to Complete
	model.SetState(StateComplete)
	if model.GetCurrentState() != StateComplete {
		t.Errorf("Expected state to be StateComplete, got %v", model.GetCurrentState())
	}
}

func TestAppModel_BackNavigation(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	// Test back navigation from BranchSelection to TypeSelection
	model.SetState(StateBranchSelection)
	updatedModel, _ := model.handleBack()
	appModel := updatedModel.(AppModel)
	if appModel.GetCurrentState() != StateTypeSelection {
		t.Errorf("Expected back from BranchSelection to go to TypeSelection, got %v", appModel.GetCurrentState())
	}

	// Test back navigation from TicketInput to BranchSelection
	model.SetState(StateTicketInput)
	updatedModel, _ = model.handleBack()
	appModel = updatedModel.(AppModel)
	if appModel.GetCurrentState() != StateBranchSelection {
		t.Errorf("Expected back from TicketInput to go to BranchSelection, got %v", appModel.GetCurrentState())
	}

	// Test back navigation from Confirmation to TicketInput
	model.SetState(StateConfirmation)
	updatedModel, _ = model.handleBack()
	appModel = updatedModel.(AppModel)
	if appModel.GetCurrentState() != StateTicketInput {
		t.Errorf("Expected back from Confirmation to go to TicketInput, got %v", appModel.GetCurrentState())
	}
}

func TestAppModel_KeyboardHandling(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	// Test quit key
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := model.Update(quitMsg)
	if cmd == nil {
		t.Error("Expected quit command to be returned")
	}

	// Test Ctrl+C quit
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd = model.Update(ctrlCMsg)
	if cmd == nil {
		t.Error("Expected quit command to be returned for Ctrl+C")
	}

	// Test Esc key for back navigation
	model.SetState(StateBranchSelection)
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(escMsg)
	appModel := updatedModel.(AppModel)
	if appModel.GetCurrentState() != StateTypeSelection {
		t.Errorf("Expected Esc to trigger back navigation, got state %v", appModel.GetCurrentState())
	}
}

func TestAppModel_WindowSizeHandling(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	// Test window size update
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(sizeMsg)
	appModel := updatedModel.(AppModel)

	if appModel.width != 100 {
		t.Errorf("Expected width to be 100, got %d", appModel.width)
	}

	if appModel.height != 50 {
		t.Errorf("Expected height to be 50, got %d", appModel.height)
	}
}

func TestAppModel_BranchSelectionUpdate(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
			{Name: "develop", IsCurrent: false, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)
	model.SetState(StateBranchSelection)

	// Simulate branch selection by setting the branch model state
	model.branchModel.Reset()
	
	// Test that branch selection updates work
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.updateBranchSelection(enterMsg)

	// The model should handle the update without errors
	// State might change if selection was made, which is expected behavior
	_ = updatedModel // Just verify no panic occurs
}

func TestAppModel_TicketInputUpdate(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)
	model.SetState(StateTicketInput)
	model.selectedType = "feature"
	model.selectedBranch = "main"

	// Test ticket input handling
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.updateTicketInput(enterMsg)

	// The model should handle the update without errors
	// State might change if form was completed, which is expected behavior
	_ = updatedModel // Just verify no panic occurs
}

func TestAppModel_ErrorHandling(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
		createError:   &MockGitError{"failed to create branch"},
	}

	model := NewAppModel(cfg, mockGit)

	// Test setting error
	testError := &MockGitError{"test error"}
	model.SetError(testError)

	if model.err != testError {
		t.Error("Expected error to be set correctly")
	}

	// Test clearing error
	model.ClearError()
	if model.err != nil {
		t.Error("Expected error to be cleared")
	}
}

func TestAppModel_SelectedDataManagement(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	// Test setting selected data
	model.SetSelectedData("feature", "main", "JIRA-123", "Test Title")

	branchType, baseBranch, ticket, title := model.GetSelectedData()
	if branchType != "feature" {
		t.Errorf("Expected branch type 'feature', got '%s'", branchType)
	}
	if baseBranch != "main" {
		t.Errorf("Expected base branch 'main', got '%s'", baseBranch)
	}
	if ticket != "JIRA-123" {
		t.Errorf("Expected ticket 'JIRA-123', got '%s'", ticket)
	}
	if title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", title)
	}
}

func TestAppModel_BranchNameGeneration(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	// Test branch name generation with title
	model.selectedType = "feature"
	model.ticketNumber = "JIRA-123"
	model.ticketTitle = "Test Feature Title"

	branchName := model.generateBranchName()
	expected := "feature/JIRA-123-test-feature-title"
	if branchName != expected {
		t.Errorf("Expected branch name '%s', got '%s'", expected, branchName)
	}

	// Test branch name generation without title
	model.ticketTitle = ""
	branchName = model.generateBranchName()
	expected = "feature/JIRA-123-jira-123"
	if branchName != expected {
		t.Errorf("Expected branch name '%s', got '%s'", expected, branchName)
	}
}

func TestAppModel_TitleSanitization(t *testing.T) {
	cfg := config.GetDefaultConfig()
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)

	tests := []struct {
		input    string
		expected string
	}{
		{"Simple Title", "feature/JIRA-123-simple-title"},
		{"Title With Special @#$ Characters!", "feature/JIRA-123-title-with-special-characters"},
		{"Multiple   Spaces", "feature/JIRA-123-multiple-spaces"},
		{"--Leading--And--Trailing--", "feature/JIRA-123-leading-and-trailing"},
		{"", "feature/JIRA-123-jira-123"},
		{"CamelCase Title", "feature/JIRA-123-camelcase-title"},
	}

	// Set up model state for branch name generation
	model.selectedType = "feature"
	model.ticketNumber = "JIRA-123"

	for _, test := range tests {
		model.ticketTitle = test.input
		result := model.generateBranchName()
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestAppModel_BranchNameLengthLimit(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.MaxBranchLength = 30 // Set a short limit for testing
	
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
	}

	model := NewAppModel(cfg, mockGit)
	model.selectedType = "feature"
	model.ticketNumber = "JIRA-123"
	model.ticketTitle = "This is a very long title that should be truncated"

	branchName := model.generateBranchName()
	
	if len(branchName) > cfg.MaxBranchLength {
		t.Errorf("Expected branch name length to be <= %d, got %d ('%s')", cfg.MaxBranchLength, len(branchName), branchName)
	}

	// Should still contain the type and ticket number
	if !contains(branchName, "feature") || !contains(branchName, "JIRA-123") {
		t.Errorf("Expected branch name to contain type and ticket number, got '%s'", branchName)
	}
}

func TestAppModel_BranchCreation(t *testing.T) {
	cfg := config.GetDefaultConfig()
	
	// Test successful branch creation
	mockGit := &MockGitRepository{
		branches: []git.BranchInfo{
			{Name: "main", IsCurrent: true, IsRemote: false},
		},
		currentBranch: "main",
		createError:   nil,
	}

	model := NewAppModel(cfg, mockGit)
	model.selectedBranch = "main"
	model.finalBranch = "feature/JIRA-123-test"

	err := model.createBranch()
	if err != nil {
		t.Errorf("Expected no error for successful branch creation, got %v", err)
	}

	// Test failed branch creation
	mockGit.createError = &MockGitError{"failed to create branch"}
	err = model.createBranch()
	if err == nil {
		t.Error("Expected error for failed branch creation")
	}
}

// MockGitError for testing
type MockGitError struct {
	message string
}

func (e *MockGitError) Error() string {
	return e.message
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}