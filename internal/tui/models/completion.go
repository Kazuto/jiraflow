package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jiraflow/internal/tui/components"
)

// CompletionState represents the state of the completion screen
type CompletionState int

const (
	CompletionSuccess CompletionState = iota
	CompletionError
)

// CompletionModel represents the completion screen model
type CompletionModel struct {
	width  int
	height int
	
	// State
	state CompletionState
	
	// Data
	branchName   string
	baseBranch   string
	errorMessage string
	
	// Control
	shouldExit bool
}

// NewCompletionModel creates a new completion model
func NewCompletionModel() CompletionModel {
	return CompletionModel{
		state: CompletionSuccess,
	}
}

// SetSuccess sets the completion screen to success state
func (m *CompletionModel) SetSuccess(branchName, baseBranch string) {
	m.state = CompletionSuccess
	m.branchName = branchName
	m.baseBranch = baseBranch
	m.errorMessage = ""
}

// SetError sets the completion screen to error state
func (m *CompletionModel) SetError(errorMessage string) {
	m.state = CompletionError
	m.errorMessage = errorMessage
}

// SetSize sets the size of the completion screen
func (m *CompletionModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init initializes the completion model
func (m CompletionModel) Init() tea.Cmd {
	return nil
}

// Update handles events for the completion screen
func (m CompletionModel) Update(msg tea.Msg) (CompletionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "q", "ctrl+c":
			m.shouldExit = true
			return m, tea.Quit
		}
	}
	
	return m, nil
}

// View renders the completion screen
func (m CompletionModel) View() string {
	var sections []string
	
	switch m.state {
	case CompletionSuccess:
		sections = append(sections, m.renderSuccess()...)
	case CompletionError:
		sections = append(sections, m.renderError()...)
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderSuccess renders the success completion screen
func (m CompletionModel) renderSuccess() []string {
	var sections []string
	
	// Success icon and title
	successIcon := lipgloss.NewStyle().
		Foreground(components.ColorSuccess).
		Bold(true).
		Render("✅")
	
	title := components.TitleStyle.
		Foreground(components.ColorSuccess).
		Render("Branch Created Successfully!")
	
	titleLine := lipgloss.JoinHorizontal(lipgloss.Center, successIcon, " ", title)
	sections = append(sections, titleLine)
	sections = append(sections, "")
	
	// Success details in a styled box
	details := m.renderSuccessDetails()
	detailsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(components.ColorSuccess).
		Padding(1, 2).
		Width(m.width - 8).
		Render(details)
	sections = append(sections, detailsBox)
	sections = append(sections, "")
	
	// Next steps
	nextStepsTitle := lipgloss.NewStyle().
		Foreground(components.ColorSecondary).
		Bold(true).
		Render("Next Steps:")
	sections = append(sections, nextStepsTitle)
	
	nextSteps := []string{
		"• Your new branch has been created and checked out",
		"• You can now start working on your feature",
		"• Remember to push your branch when ready: git push -u origin " + m.branchName,
	}
	
	for _, step := range nextSteps {
		stepText := lipgloss.NewStyle().
			Foreground(components.ColorMuted).
			Render(step)
		sections = append(sections, stepText)
	}
	
	sections = append(sections, "")
	
	// Exit instructions with enhanced help
	help := m.renderSuccessHelp()
	sections = append(sections, help)
	
	return sections
}

// renderError renders the error completion screen
func (m CompletionModel) renderError() []string {
	var sections []string
	
	// Error icon and title
	errorIcon := lipgloss.NewStyle().
		Foreground(components.ColorError).
		Bold(true).
		Render("❌")
	
	title := components.TitleStyle.
		Foreground(components.ColorError).
		Render("Branch Creation Failed")
	
	titleLine := lipgloss.JoinHorizontal(lipgloss.Center, errorIcon, " ", title)
	sections = append(sections, titleLine)
	sections = append(sections, "")
	
	// Error details in a styled box
	errorDetails := m.renderErrorDetails()
	errorBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(components.ColorError).
		Padding(1, 2).
		Width(m.width - 8).
		Render(errorDetails)
	sections = append(sections, errorBox)
	sections = append(sections, "")
	
	// Troubleshooting tips
	tipsTitle := lipgloss.NewStyle().
		Foreground(components.ColorWarning).
		Bold(true).
		Render("Troubleshooting Tips:")
	sections = append(sections, tipsTitle)
	
	tips := []string{
		"• Check if you have the necessary permissions",
		"• Ensure you're in a valid Git repository",
		"• Verify the base branch exists and is accessible",
		"• Make sure the branch name doesn't already exist",
	}
	
	for _, tip := range tips {
		tipText := lipgloss.NewStyle().
			Foreground(components.ColorMuted).
			Render(tip)
		sections = append(sections, tipText)
	}
	
	sections = append(sections, "")
	
	// Exit instructions with enhanced help
	help := m.renderErrorHelp()
	sections = append(sections, help)
	
	return sections
}

// renderSuccessDetails renders the success details
func (m CompletionModel) renderSuccessDetails() string {
	var details []string
	
	// Branch created
	branchLabel := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Width(15).
		Render("Created:")
	branchValue := components.SuccessStyle.Render(m.branchName)
	details = append(details, fmt.Sprintf("%s %s", branchLabel, branchValue))
	
	// Base branch
	baseLabel := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Width(15).
		Render("From:")
	baseValue := components.SelectedStyle.Render(m.baseBranch)
	details = append(details, fmt.Sprintf("%s %s", baseLabel, baseValue))
	
	// Status
	statusLabel := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Width(15).
		Render("Status:")
	statusValue := components.SuccessStyle.Render("✓ Active and checked out")
	details = append(details, fmt.Sprintf("%s %s", statusLabel, statusValue))
	
	return strings.Join(details, "\n")
}

// renderErrorDetails renders the error details
func (m CompletionModel) renderErrorDetails() string {
	var details []string
	
	// Error message
	errorLabel := lipgloss.NewStyle().
		Foreground(components.ColorMuted).
		Width(15).
		Render("Error:")
	errorValue := components.ErrorStyle.Render(m.errorMessage)
	details = append(details, fmt.Sprintf("%s %s", errorLabel, errorValue))
	
	// Attempted branch (if available)
	if m.branchName != "" {
		branchLabel := lipgloss.NewStyle().
			Foreground(components.ColorMuted).
			Width(15).
			Render("Attempted:")
		branchValue := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Render(m.branchName)
		details = append(details, fmt.Sprintf("%s %s", branchLabel, branchValue))
	}
	
	return strings.Join(details, "\n")
}

// ShouldExit returns true if the application should exit
func (m CompletionModel) ShouldExit() bool {
	return m.shouldExit
}

// Reset resets the completion model state
func (m *CompletionModel) Reset() {
	m.shouldExit = false
	m.state = CompletionSuccess
	m.branchName = ""
	m.baseBranch = ""
	m.errorMessage = ""
}

// renderSuccessHelp renders help text for the success screen
func (m CompletionModel) renderSuccessHelp() string {
	var sections []string
	
	// Main help
	mainHelp := []string{
		"enter exit application",
		"q quit",
	}
	
	mainHelpText := strings.Join(mainHelp, " • ")
	sections = append(sections, components.HelpStyle.Render(mainHelpText))
	
	// Context help
	contextHelp := []string{
		"Branch created and checked out successfully",
		"Ready to start development",
	}
	
	contextStyle := components.HelpStyle.
		Foreground(components.ColorMuted).
		Faint(true)
	sections = append(sections, contextStyle.Render(strings.Join(contextHelp, " • ")))
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderErrorHelp renders help text for the error screen
func (m CompletionModel) renderErrorHelp() string {
	var sections []string
	
	// Main help
	mainHelp := []string{
		"enter exit application",
		"q quit",
	}
	
	mainHelpText := strings.Join(mainHelp, " • ")
	sections = append(sections, components.HelpStyle.Render(mainHelpText))
	
	// Context help
	contextHelp := []string{
		"Branch creation failed",
		"Check troubleshooting tips above",
	}
	
	contextStyle := components.HelpStyle.
		Foreground(components.ColorMuted).
		Faint(true)
	sections = append(sections, contextStyle.Render(strings.Join(contextHelp, " • ")))
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// GetState returns the current completion state
func (m CompletionModel) GetState() CompletionState {
	return m.state
}