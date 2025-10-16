package models

import tea "github.com/charmbracelet/bubbletea"

// TypeSelectorModel handles branch type selection
type TypeSelectorModel struct {
	types    []string
	selected int
}

// Init initializes the type selector model
func (m TypeSelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles events for the type selector
func (m TypeSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View renders the type selector interface
func (m TypeSelectorModel) View() string {
	return "Type Selector - Coming Soon!"
}