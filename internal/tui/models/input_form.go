package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputFormModel handles ticket number and title input
type InputFormModel struct {
	ticketInput textinput.Model
	titleInput  textinput.Model
	focused     int
}

// Init initializes the input form model
func (m InputFormModel) Init() tea.Cmd {
	return nil
}

// Update handles events for the input form
func (m InputFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View renders the input form interface
func (m InputFormModel) View() string {
	return "Input Form - Coming Soon!"
}