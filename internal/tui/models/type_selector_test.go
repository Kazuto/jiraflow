package models

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"jiraflow/internal/config"
)

func TestNewTypeSelectorModel(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test that model is initialized correctly
	if len(model.types) != 4 {
		t.Errorf("Expected 4 types, got %d", len(model.types))
	}

	// Test that default type is marked correctly
	foundDefault := false
	for _, item := range model.types {
		if item.isDefault {
			foundDefault = true
			if item.key != cfg.DefaultBranchType {
				t.Errorf("Expected default type to be %s, got %s", cfg.DefaultBranchType, item.key)
			}
		}
	}
	if !foundDefault {
		t.Error("No default type found")
	}

	// Test that all expected types are present
	expectedTypes := []string{"feature", "hotfix", "refactor", "support"}
	typeKeys := make(map[string]bool)
	for _, item := range model.types {
		typeKeys[item.key] = true
	}

	for _, expected := range expectedTypes {
		if !typeKeys[expected] {
			t.Errorf("Expected type %s not found", expected)
		}
	}
}

func TestTypeSelectorModel_Update(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test window size update
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ := model.Update(sizeMsg)
	if updatedModel.width != 80 || updatedModel.height != 24 {
		t.Errorf("Window size not updated correctly")
	}

	// Test Enter key selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = model.Update(enterMsg)
	if !updatedModel.HasSelection() {
		t.Error("Expected selection after Enter key")
	}

	// The default selection should be the default branch type
	if updatedModel.GetSelected() != cfg.DefaultBranchType {
		t.Errorf("Expected selected type to be %s, got %s", cfg.DefaultBranchType, updatedModel.GetSelected())
	}
}

func TestTypeSelectorModel_Navigation(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test down arrow navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(downMsg)

	// Test up arrow navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = updatedModel.Update(upMsg)

	// Test that we can navigate without errors
	// The exact position depends on the list implementation,
	// so we just verify no panics occur and model is valid
	if len(updatedModel.types) == 0 {
		t.Error("Types lost during navigation")
	}
}

func TestTypeSelectorModel_GetCurrentItem(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test getting current item
	item, ok := model.GetCurrentItem()
	if !ok {
		t.Error("Expected to get current item")
	}

	// Should be the default type initially
	if !item.isDefault {
		t.Error("Expected current item to be default type")
	}
}

func TestTypeSelectorModel_Reset(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Make a selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = model.Update(enterMsg)

	if !model.HasSelection() {
		t.Error("Expected selection before reset")
	}

	// Reset the model
	model.Reset()

	if model.HasSelection() {
		t.Error("Expected no selection after reset")
	}
}

func TestTypeSelectorModel_SetSize(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test setting size
	model.SetSize(100, 30)

	if model.width != 100 || model.height != 30 {
		t.Errorf("Expected size 100x30, got %dx%d", model.width, model.height)
	}
}

func TestTypeSelectorModel_GetSelectedDisplayName(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Make a selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = model.Update(enterMsg)

	displayName := model.GetSelectedDisplayName()
	selectedKey := model.GetSelected()

	// Verify the display name matches the selected key
	expectedDisplayName := cfg.BranchTypes[selectedKey]
	if displayName != expectedDisplayName {
		t.Errorf("Expected display name %s, got %s", expectedDisplayName, displayName)
	}
}

func TestTypeSelectorModel_KeyboardNavigation(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test vim-style navigation (j/k keys)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	updatedModel, _ := model.Update(jMsg)

	// Should not crash and maintain valid state
	if len(updatedModel.types) != 4 {
		t.Error("Navigation should not affect available types")
	}

	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	updatedModel, _ = updatedModel.Update(kMsg)

	if len(updatedModel.types) != 4 {
		t.Error("Navigation should not affect available types")
	}

	// Test arrow key navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	updatedModel, _ = model.Update(downMsg)
	if len(updatedModel.types) != 4 {
		t.Error("Arrow navigation should not affect available types")
	}

	updatedModel, _ = updatedModel.Update(upMsg)
	if len(updatedModel.types) != 4 {
		t.Error("Arrow navigation should not affect available types")
	}
}

func TestTypeSelectorModel_SelectionHandling(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test initial selection state
	if model.HasSelection() {
		t.Error("Expected no selection initially")
	}

	// Test Enter key selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)

	if !updatedModel.HasSelection() {
		t.Error("Expected selection after Enter key")
	}

	selected := updatedModel.GetSelected()
	if selected == "" {
		t.Error("Expected non-empty selected type")
	}

	// Test that selected type is valid
	validTypes := []string{"feature", "hotfix", "refactor", "support"}
	isValid := false
	for _, validType := range validTypes {
		if selected == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		t.Errorf("Selected type '%s' is not valid", selected)
	}
}

func TestTypeSelectorModel_StateTransitions(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test navigation through all types
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	
	// Navigate through all types
	for i := 0; i < len(model.types); i++ {
		currentItem, ok := model.GetCurrentItem()
		if !ok {
			t.Errorf("Expected to get current item at position %d", i)
		}
		
		if currentItem.key == "" {
			t.Errorf("Expected current item to have a key at position %d", i)
		}
		
		model, _ = model.Update(downMsg)
	}

	// Test selection at different positions
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = model.Update(enterMsg)

	if !model.HasSelection() {
		t.Error("Expected selection after navigation and Enter")
	}
}

func TestTypeSelectorModel_BackNavigation(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test Esc key (back navigation)
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(escMsg)

	// Model should handle back navigation gracefully
	// (The actual behavior depends on parent component handling)
	if len(updatedModel.types) != 4 {
		t.Error("Back navigation should not affect available types")
	}
}

func TestTypeSelectorModel_WindowSizeHandling(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test window size update
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := model.Update(sizeMsg)

	if updatedModel.width != 100 {
		t.Errorf("Expected width to be 100, got %d", updatedModel.width)
	}

	if updatedModel.height != 30 {
		t.Errorf("Expected height to be 30, got %d", updatedModel.height)
	}
}

func TestTypeSelectorModel_GetAvailableTypes(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	availableTypes := model.GetAvailableTypes()
	if len(availableTypes) != 4 {
		t.Errorf("Expected 4 available types, got %d", len(availableTypes))
	}

	// Test that all types have required fields
	for i, typeItem := range availableTypes {
		if typeItem.key == "" {
			t.Errorf("Type at index %d has empty key", i)
		}
		if typeItem.displayName == "" {
			t.Errorf("Type at index %d has empty display name", i)
		}
		if typeItem.description == "" {
			t.Errorf("Type at index %d has empty description", i)
		}
	}
}

func TestTypeSelectorModel_DefaultTypeHandling(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test that exactly one type is marked as default
	defaultCount := 0
	var defaultType TypeItem
	for _, typeItem := range model.types {
		if typeItem.isDefault {
			defaultCount++
			defaultType = typeItem
		}
	}

	if defaultCount != 1 {
		t.Errorf("Expected exactly 1 default type, got %d", defaultCount)
	}

	if defaultType.key != cfg.DefaultBranchType {
		t.Errorf("Expected default type to be '%s', got '%s'", cfg.DefaultBranchType, defaultType.key)
	}

	// Test that initial selection is the default type
	currentItem, ok := model.GetCurrentItem()
	if !ok {
		t.Error("Expected to get current item initially")
	}

	if !currentItem.isDefault {
		t.Error("Expected initial current item to be the default type")
	}
}

func TestTypeSelectorModel_ViewRendering(t *testing.T) {
	cfg := config.GetDefaultConfig()
	model := NewTypeSelectorModel(cfg)

	// Test basic view rendering
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Test view after navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ = model.Update(downMsg)
	
	viewAfterNav := model.View()
	if viewAfterNav == "" {
		t.Error("Expected non-empty view after navigation")
	}

	// Test view after selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = model.Update(enterMsg)
	
	viewAfterSelection := model.View()
	if viewAfterSelection == "" {
		t.Error("Expected non-empty view after selection")
	}
}

func TestTypeItem_Methods(t *testing.T) {
	item := TypeItem{
		key:         "feature",
		displayName: "feature",
		description: "New features and enhancements",
		isDefault:   true,
	}

	// Test FilterValue
	if item.FilterValue() != "feature" {
		t.Errorf("Expected FilterValue to be 'feature', got %s", item.FilterValue())
	}

	// Test Title for default item
	title := item.Title()
	if title == "feature" { // Should include "(default)" when rendered
		t.Error("Expected title to indicate default status")
	}

	// Test Description
	if item.Description() != "New features and enhancements" {
		t.Errorf("Expected description to be 'New features and enhancements', got %s", item.Description())
	}

	// Test non-default item
	nonDefaultItem := TypeItem{
		key:         "hotfix",
		displayName: "hotfix",
		description: "Critical bug fixes",
		isDefault:   false,
	}

	nonDefaultTitle := nonDefaultItem.Title()
	if nonDefaultTitle != "hotfix" {
		t.Errorf("Expected non-default title to be 'hotfix', got %s", nonDefaultTitle)
	}
}