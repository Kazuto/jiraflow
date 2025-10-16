package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// AppState represents the current state of the TUI application
type AppState int

const (
	StateTypeSelection AppState = iota
	StateBranchSelection
	StateTicketInput
	StateTitleInput
	StateConfirmation
	StateComplete
)

// AppModel represents the main TUI application model
type AppModel struct {
	state AppState
	err   error
}

// Init initializes the TUI application
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update handles TUI events and state transitions
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View renders the TUI interface
func (m AppModel) View() string {
	return "JiraFlow TUI - Coming Soon!"
}