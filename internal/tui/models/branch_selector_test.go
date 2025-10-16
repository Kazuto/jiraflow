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

func TestBranchSelectorModel_KeyboardNavigation(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
		{Name: "feature/test", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test up/down navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	// Test down navigation
	updatedModel, _ := model.Update(downMsg)
	// Should not crash and maintain valid state
	if len(updatedModel.filteredItems) != 3 {
		t.Error("Navigation should not affect filtered items")
	}

	// Test up navigation
	updatedModel, _ = updatedModel.Update(upMsg)
	if len(updatedModel.filteredItems) != 3 {
		t.Error("Navigation should not affect filtered items")
	}

	// Test 'k' and 'j' vim-style navigation
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}

	updatedModel, _ = model.Update(jMsg)
	if len(updatedModel.filteredItems) != 3 {
		t.Error("Vim navigation should not affect filtered items")
	}

	updatedModel, _ = updatedModel.Update(kMsg)
	if len(updatedModel.filteredItems) != 3 {
		t.Error("Vim navigation should not affect filtered items")
	}
}

func TestBranchSelectorModel_SearchKeyboardHandling(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
		{Name: "feature/auth", IsCurrent: false, IsRemote: false},
		{Name: "feature/ui", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test entering search mode with '/'
	searchMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
	updatedModel, _ := model.Update(searchMsg)

	if !updatedModel.searching {
		t.Error("Expected to enter search mode after pressing '/'")
	}

	// Test typing in search mode
	fMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}
	updatedModel, _ = updatedModel.Update(fMsg)

	eMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")}
	updatedModel, _ = updatedModel.Update(eMsg)

	aMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	updatedModel, _ = updatedModel.Update(aMsg)

	// Should filter to feature branches
	if len(updatedModel.filteredItems) != 2 {
		t.Errorf("Expected 2 filtered items for 'fea', got %d", len(updatedModel.filteredItems))
	}

	// Test clearing search with Ctrl+U
	clearMsg := tea.KeyMsg{Type: tea.KeyCtrlU}
	updatedModel, _ = updatedModel.Update(clearMsg)

	if updatedModel.searchInput.Value() != "" {
		t.Error("Expected search input to be cleared after Ctrl+U")
	}

	if len(updatedModel.filteredItems) != 4 {
		t.Errorf("Expected all 4 items after clearing search, got %d", len(updatedModel.filteredItems))
	}

	// Test exiting search mode with Enter
	updatedModel.searching = true
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = updatedModel.Update(enterMsg)

	if updatedModel.searching {
		t.Error("Expected to exit search mode after pressing Enter")
	}
}

func TestBranchSelectorModel_RealTimeFiltering(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
		{Name: "feature/auth-service", IsCurrent: false, IsRemote: false},
		{Name: "feature/user-interface", IsCurrent: false, IsRemote: false},
		{Name: "hotfix/critical-bug", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test progressive filtering
	testCases := []struct {
		searchTerm    string
		expectedCount int
		description   string
	}{
		{"", 5, "empty search should show all branches"},
		{"feature", 2, "search 'feature' should match feature branches"},
		{"auth", 1, "search 'auth' should match auth-service"},
		{"main", 1, "search 'main' should match main branch"},
		{"hotfix", 1, "search 'hotfix' should match hotfix branch"},
		{"nonexistent", 0, "search 'nonexistent' should match nothing"},
		{"dev", 1, "search 'dev' should match develop branch"},
		{"interface", 1, "search 'interface' should match user-interface"},
	}

	for _, tc := range testCases {
		model.updateFilter(tc.searchTerm)
		if len(model.filteredItems) != tc.expectedCount {
			t.Errorf("%s: expected %d items, got %d", tc.description, tc.expectedCount, len(model.filteredItems))
		}
	}
}

func TestBranchSelectorModel_SearchResultsHandling(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test search with results
	model.updateFilter("main")
	if !model.searchResults.HasResults {
		t.Error("Expected search results to have results for 'main'")
	}

	// Test search with no results
	model.updateFilter("nonexistent")
	if model.searchResults.HasResults {
		t.Error("Expected search results to have no results for 'nonexistent'")
	}

	// Test that list selection resets when filtering
	model.updateFilter("main")
	if len(model.filteredItems) > 0 {
		// List should reset selection to first item
		currentItem, ok := model.GetCurrentItem()
		if ok && currentItem.name != "main" {
			t.Error("Expected selection to reset to first filtered item")
		}
	}
}

func TestBranchSelectorModel_WindowSizeUpdate(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test window size message handling
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, _ := model.Update(sizeMsg)

	if updatedModel.width != 120 {
		t.Errorf("Expected width to be 120, got %d", updatedModel.width)
	}

	if updatedModel.height != 40 {
		t.Errorf("Expected height to be 40, got %d", updatedModel.height)
	}

	// Verify that search input width is updated
	expectedSearchWidth := 120 - 10
	if updatedModel.searchInput.Width != expectedSearchWidth {
		t.Errorf("Expected search input width to be %d, got %d", expectedSearchWidth, updatedModel.searchInput.Width)
	}
}

func TestBranchSelectorModel_StateManagement(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test initial state
	if model.HasSelection() {
		t.Error("Expected no selection initially")
	}

	if model.GetSelected() != "" {
		t.Error("Expected empty selected branch initially")
	}

	// Test making a selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)

	if !updatedModel.HasSelection() {
		t.Error("Expected selection after Enter key")
	}

	selected := updatedModel.GetSelected()
	if selected == "" {
		t.Error("Expected non-empty selected branch after Enter key")
	}

	// Test getting current item
	currentItem, ok := model.GetCurrentItem()
	if !ok {
		t.Error("Expected to get current item")
	}

	if currentItem.name == "" {
		t.Error("Expected current item to have a name")
	}
}

func TestBranchSelectorModel_ViewRendering(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "main", IsCurrent: true, IsRemote: false},
		{Name: "develop", IsCurrent: false, IsRemote: false},
	}

	model := NewBranchSelectorModel(branches)

	// Test basic view rendering
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Test view with search term
	model.updateFilter("main")
	viewWithSearch := model.View()
	if viewWithSearch == "" {
		t.Error("Expected non-empty view with search")
	}

	// Test view with no search results
	model.updateFilter("nonexistent")
	viewNoResults := model.View()
	if viewNoResults == "" {
		t.Error("Expected non-empty view even with no results")
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