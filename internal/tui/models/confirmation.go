package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jiraflow/internal/tui/components"
)

// ConfirmationModel represents the confirmation screen model
type ConfirmationModel struct {
	width  int
	height int
	
	// Data to display
	branchType   string
	baseBranch   string
	ticketNumber string
	ticketTitle  string
	finalBranch  string
	
	// State
	confirmed bool
}

// NewConfirmationModel creates a new confirmation model
func NewConfirmationModel() ConfirmationModel {
	return ConfirmationModel{}
}

// SetData sets the data to be displayed in the confirmation screen
func (m *ConfirmationModel) SetData(branchType, baseBranch, ticketNumber, ticketTitle, finalBranch string) {
	m.branchType = branchType
	m.baseBranch = baseBranch
	m.ticketNumber = ticketNumber
	m.ticketTitle = ticketTitle
	m.finalBranch = finalBranch
}

// SetSize sets the size of the confirmation screen
func (m *ConfirmationModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init initializes the confirmation model
func (m ConfirmationModel) Init() tea.Cmd {
	return nil
}

// Update handles events for the confirmation screen
func (m ConfirmationModel) Update(msg tea.Msg) (ConfirmationModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.confirmed = true
			return m, nil
		}
	}
	
	return m, nil
}

// View renders the confirmation screen
func (m ConfirmationModel) View() string {
	var sections []string
	
	// Title
	title := components.TitleStyle.Render("📋 Confirm Branch Creation")
	sections = append(sections, title)
	sections = append(sections, "")
	
	// Summary section
	summaryTitle := lipgloss.NewStyle().
		Foreground(components.ColorSecondary).
		Bold(true).
		Render("Summary:")
	sections = append(sections, summaryTitle)
	sections = append(sections, "")
	
	// Details in a styled box
	details := m.renderDetails()
	detailsBox := components.BorderStyle.Copy().
		Width(m.width - 8).
		Render(details)
	sections = append(sections, detailsBox)
	sections = append(sections, "")
	
	// Final branch name highlight
	branchTitle := lipgloss.NewStyle().
		Foreground(components.ColorSecondary).
		Bold(true).
		Render("Branch to create:")
	sections = append(sections, branchTitle)
	
	branchName := lipgloss.NewStyle().
		Foreground(components.ColorPrimary).
		Bold(true).
		Background(lipgloss.Color("236")).
		Padding(0, 2).
		Render(m.finalBranch)
	sections = append(sections, branchName)
	sections = append(sections, "")
	
	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Italic(true).
		Render("Press Enter to create this branch, or Esc to go back")
	sections = append(sections, instructions)
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderDetails renders the detailed information in a formatted way
func (m ConfirmationModel) renderDetails() string {
	var details []string
	
	// Branch type
	typeLabel := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Width(15).
		Render("Type:")
	typeValue := components.SelectedStyle.Render(m.branchType)
	details = append(details, fmt.Sprintf("%s %s", typeLabel, typeValue))
	
	// Base branch
	baseLabel := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Width(15).
		Render("Base branch:")
	baseValue := components.SelectedStyle.Render(m.baseBranch)
	details = append(details, fmt.Sprintf("%s %s", baseLabel, baseValue))
	
	// Ticket number
	ticketLabel := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Width(15).
		Render("Ticket:")
	ticketValue := components.SelectedStyle.Render(m.ticketNumber)
	details = append(details, fmt.Sprintf("%s %s", ticketLabel, ticketValue))
	
	// Title (if provided)
	if m.ticketTitle != "" {
		titleLabel := lipgloss.NewStyle().
			Foreground(components.ColorMuted).
			Width(15).
			Render("Title:")
		titleValue := components.SelectedStyle.Render(m.ticketTitle)
		details = append(details, fmt.Sprintf("%s %s", titleLabel, titleValue))
	}
	
	return strings.Join(details, "\n")
}

// HasConfirmed returns true if the user has confirmed the branch creation
func (m ConfirmationModel) HasConfirmed() bool {
	return m.confirmed
}

// Reset resets the confirmation state
func (m *ConfirmationModel) Reset() {
	m.confirmed = false
}

// GetFinalBranch returns the final branch name
func (m ConfirmationModel) GetFinalBranch() string {
	return m.finalBranch
}