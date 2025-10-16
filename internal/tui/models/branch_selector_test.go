package models

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"jiraflow/internal/git"
)

func TestNewBranchSelectorModel(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
		{Name: "feature/test", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test initial state
	if model.searching {
		t.Error("Expected searching to be false initially")
	}

	if model.selected != "" {
		t.Error("Expected no selection initially")
	}

	if len(model.allBranches) != 3 {
		t.Errorf("Expected 3 branches, got %d", len(model.allBranches))
	}

	if len(model.filteredItems) != 3 {
		t.Errorf("Expected 3 filtered items initially, got %d", len(model.filteredItems))
	}
}

func TestBranchSelectorModel_UpdateFilter(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
		{Name: "feature/test", IsCurrent: false, IsRemote: false},
		{Name: "feature/auth", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test filtering with "feature"
	model.updateFilter("feature")

	if len(model.filteredItems) != 2 {
		t.Errorf("Expected 2 filtered items for 'feature', got %d", len(model.filteredItems))
	}

	// Test filtering with "main"
	model.updateFilter("main")

	if len(model.filteredItems) != 1 {
		t.Errorf("Expected 1 filtered item for 'main', got %d", len(model.filteredItems))
	}

	// Test filtering with no matches
	model.updateFilter("nonexistent")

	if len(model.filteredItems) != 0 {
		t.Errorf("Expected 0 filtered items for 'nonexistent', got %d", len(model.filteredItems))
	}

	// Test clearing filter
	model.updateFilter("")

	if len(model.filteredItems) != 4 {
		t.Errorf("Expected 4 filtered items when clearing filter, got %d", len(model.filteredItems))
	}
}

func TestBranchSelectorModel_SearchMode(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test entering search mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
	updatedModel, _ := model.Update(msg)

	if !updatedModel.searching {
		t.Error("Expected searching to be true after pressing '/'")
	}

	// Test exiting search mode with Esc
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ = updatedModel.Update(msg)

	if updatedModel.searching {
		t.Error("Expected searching to be false after pressing Esc")
	}
}

func TestBranchSelectorModel_Selection(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Initially no selection
	if model.HasSelection() {
		t.Error("Expected no selection initially")
	}

	// Simulate Enter key press to select current item
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)

	if !updatedModel.HasSelection() {
		t.Error("Expected selection after pressing Enter")
	}

	selected := updatedModel.GetSelected()
	if selected == "" {
		t.Error("Expected selected branch name to be non-empty")
	}
}

func TestBranchItem_FilterValue(t *testing.T) {
	item := BranchItem{name: "test-branch", isCurrent: false}
	
	if item.FilterValue() != "test-branch" {
		t.Errorf("Expected FilterValue to return 'test-branch', got '%s'", item.FilterValue())
	}
}

func TestBranchItem_Title(t *testing.T) {
	// Test non-current branch
	item := BranchItem{name: "test-branch", isCurrent: false}
	title := item.Title()
	
	if title != "test-branch" {
		t.Errorf("Expected Title to return 'test-branch', got '%s'", title)
	}

	// Test current branch (should include indicator)
	currentItem := BranchItem{name: "main", isCurrent: true}
	currentTitle := currentItem.Title()
	
	// The title should contain the branch name and current indicator
	if !contains(currentTitle, "main") || !contains(currentTitle, "current") {
		t.Errorf("Expected current branch title to contain 'main' and 'current', got '%s'", currentTitle)
	}
}

func TestBranchSelectorModel_SetSize(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)
	model.SetSize(100, 50)

	if model.width != 100 {
		t.Errorf("Expected width to be 100, got %d", model.width)
	}

	if model.height != 50 {
		t.Errorf("Expected height to be 50, got %d", model.height)
	}
}

func TestBranchSelectorModel_Reset(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)
	
	// Set some state
	model.selected = "test"
	model.searching = true
	model.searchInput.SetValue("search")

	// Reset
	model.Reset()

	if model.selected != "" {
		t.Error("Expected selected to be empty after reset")
	}

	if model.searching {
		t.Error("Expected searching to be false after reset")
	}

	if model.searchInput.Value() != "" {
		t.Error("Expected search input to be empty after reset")
	}
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