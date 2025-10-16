package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jiraflow/internal/config"
	"jiraflow/internal/git"
	"jiraflow/internal/tui/components"
	"jiraflow/internal/tui/models"
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
	state       AppState
	config      *config.Config
	git         git.GitRepository
	typeModel   models.TypeSelectorModel
	branchModel models.BranchSelectorModel
	inputModel  models.InputFormModel
	err         error
	width       int
	height      int
	
	// State data
	selectedType   string
	selectedBranch string
	ticketNumber   string
	ticketTitle    string
	finalBranch    string
}

// keyMap defines the key bindings for the application
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
	Search key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
	}
}

var keys = DefaultKeyMap()

// NewAppModel creates a new AppModel instance
func NewAppModel(cfg *config.Config, gitRepo git.GitRepository) *AppModel {
	// Initialize branch selector with Git branches
	var branchModel models.BranchSelectorModel
	if branches, err := gitRepo.GetBranchesWithInfo(); err == nil {
		branchModel = models.NewBranchSelectorModel(branches)
	} else {
		// Fallback to empty branch selector if Git operations fail
		branchModel = models.NewBranchSelectorModel([]git.BranchInfo{})
	}
	
	// Initialize input form model (Jira client will be set later if available)
	inputModel := models.NewInputFormModel(nil)
	
	return &AppModel{
		state:       StateTypeSelection,
		config:      cfg,
		git:         gitRepo,
		branchModel: branchModel,
		inputModel:  inputModel,
	}
}

// RunTUI starts the TUI application
func RunTUI(cfg *config.Config, gitRepo git.GitRepository) error {
	model := NewAppModel(cfg, gitRepo)
	
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	
	_, err := p.Run()
	return err
}

// Init initializes the TUI application
func (m AppModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles TUI events and state transitions
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update component sizes
		m.branchModel.SetSize(msg.Width, msg.Height-6) // Leave space for header/footer
		m.inputModel.SetSize(msg.Width, msg.Height-6)
		
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Back):
			return m.handleBack()
		}

		// Handle state-specific key events
		switch m.state {
		case StateTypeSelection:
			return m.updateTypeSelection(msg)
		case StateBranchSelection:
			return m.updateBranchSelection(msg)
		case StateTicketInput:
			return m.updateTicketInput(msg)
		case StateTitleInput:
			return m.updateTitleInput(msg)
		case StateConfirmation:
			return m.updateConfirmation(msg)
		case StateComplete:
			return m.updateComplete(msg)
		}
	}

	return m, nil
}

// handleBack handles the back navigation
func (m AppModel) handleBack() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateTypeSelection:
		return m, tea.Quit
	case StateBranchSelection:
		m.state = StateTypeSelection
	case StateTicketInput:
		m.state = StateBranchSelection
	case StateTitleInput:
		m.state = StateBranchSelection // Skip back to branch selection since title is in ticket form
	case StateConfirmation:
		m.state = StateTicketInput
	case StateComplete:
		m.state = StateConfirmation
	}
	return m, nil
}

// State-specific update methods (placeholder implementations)
func (m AppModel) updateTypeSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Placeholder - will be implemented in task 4.3
	if key.Matches(msg, keys.Enter) {
		m.selectedType = "feature" // Default for now
		m.state = StateBranchSelection
	}
	return m, nil
}

func (m AppModel) updateBranchSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	// Update the branch selector model
	m.branchModel, cmd = m.branchModel.Update(msg)
	
	// Check if a branch was selected
	if m.branchModel.HasSelection() {
		m.selectedBranch = m.branchModel.GetSelected()
		m.state = StateTicketInput
		return m, cmd
	}
	
	// Handle back navigation
	if key.Matches(msg, keys.Back) {
		m.state = StateTypeSelection
		return m, cmd
	}
	
	return m, cmd
}

func (m AppModel) updateTicketInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	// Update the input form model
	m.inputModel, cmd = m.inputModel.Update(msg)
	
	// Check if form was completed
	if m.inputModel.HasCompleted() {
		m.ticketNumber = m.inputModel.GetTicketNumber()
		m.ticketTitle = m.inputModel.GetTicketTitle()
		m.state = StateConfirmation // Skip title input since it's handled in the form
		return m, cmd
	}
	
	// Handle back navigation
	if key.Matches(msg, keys.Back) {
		m.state = StateBranchSelection
		return m, cmd
	}
	
	return m, cmd
}

func (m AppModel) updateTitleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// This state is now handled within the ticket input form
	// Redirect to confirmation if we somehow end up here
	m.state = StateConfirmation
	return m, nil
}

func (m AppModel) updateConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Placeholder - will be implemented in task 4.5
	if key.Matches(msg, keys.Enter) {
		m.finalBranch = fmt.Sprintf("%s/%s-%s", m.selectedType, m.ticketNumber, "sample-title")
		m.state = StateComplete
	}
	return m, nil
}

func (m AppModel) updateComplete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Placeholder - will be implemented in task 4.5
	if key.Matches(msg, keys.Enter) {
		return m, tea.Quit
	}
	return m, nil
}

// View renders the TUI interface
func (m AppModel) View() string {
	if m.err != nil {
		return m.renderError()
	}

	var content string
	
	switch m.state {
	case StateTypeSelection:
		content = m.renderTypeSelection()
	case StateBranchSelection:
		content = m.renderBranchSelection()
	case StateTicketInput:
		content = m.renderTicketInput()
	case StateTitleInput:
		content = m.renderTitleInput()
	case StateConfirmation:
		content = m.renderConfirmation()
	case StateComplete:
		content = m.renderComplete()
	}

	// Add header and footer
	header := m.renderHeader()
	footer := m.renderFooter()
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		content,
		"",
		footer,
	)
}

// renderHeader renders the application header
func (m AppModel) renderHeader() string {
	title := components.TitleStyle.Render("🚀 JiraFlow - Interactive Branch Creator")
	
	var stateText string
	switch m.state {
	case StateTypeSelection:
		stateText = "Step 1/5: Select Branch Type"
	case StateBranchSelection:
		stateText = "Step 2/5: Select Base Branch"
	case StateTicketInput:
		stateText = "Step 3/4: Enter Ticket Information"
	case StateTitleInput:
		stateText = "Step 3/4: Enter Ticket Information"
	case StateConfirmation:
		stateText = "Step 4/4: Confirm Branch Creation"
	case StateComplete:
		stateText = "Complete!"
	}
	
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(stateText)
	
	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle)
}

// renderFooter renders the application footer with help text
func (m AppModel) renderFooter() string {
	var helpKeys []string
	
	switch m.state {
	case StateTypeSelection:
		helpKeys = []string{"↑/↓ navigate", "enter select", "q quit"}
	case StateBranchSelection:
		helpKeys = []string{"↑/↓ navigate", "/ search", "enter select", "esc back", "q quit"}
	case StateTicketInput, StateTitleInput:
		helpKeys = []string{"type to input", "enter continue", "esc back", "q quit"}
	case StateConfirmation:
		helpKeys = []string{"enter create branch", "esc back", "q quit"}
	case StateComplete:
		helpKeys = []string{"enter exit", "q quit"}
	}
	
	help := strings.Join(helpKeys, " • ")
	return components.HelpStyle.Render(help)
}

// renderError renders error messages
func (m AppModel) renderError() string {
	return components.ErrorStyle.Render("Error: " + m.err.Error())
}

// State-specific render methods (placeholder implementations)
func (m AppModel) renderTypeSelection() string {
	content := "Select branch type:\n\n"
	content += "→ feature (default)\n"
	content += "  hotfix\n"
	content += "  refactor\n"
	content += "  support\n"
	
	return content
}

func (m AppModel) renderBranchSelection() string {
	var sections []string
	
	// Show selected type
	selectedType := fmt.Sprintf("Selected type: %s", components.SelectedStyle.Render(m.selectedType))
	sections = append(sections, selectedType)
	sections = append(sections, "")
	
	// Render the branch selector
	branchView := m.branchModel.View()
	sections = append(sections, branchView)
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m AppModel) renderTicketInput() string {
	var sections []string
	
	// Show selected type and branch
	selectedType := fmt.Sprintf("Selected type: %s", components.SelectedStyle.Render(m.selectedType))
	selectedBranch := fmt.Sprintf("Selected branch: %s", components.SelectedStyle.Render(m.selectedBranch))
	sections = append(sections, selectedType)
	sections = append(sections, selectedBranch)
	sections = append(sections, "")
	
	// Render the input form
	inputView := m.inputModel.View()
	sections = append(sections, inputView)
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m AppModel) renderTitleInput() string {
	// Title input is now handled within the ticket input form
	// This should not be reached, but provide fallback
	return m.renderTicketInput()
}

func (m AppModel) renderConfirmation() string {
	content := fmt.Sprintf("Selected type: %s\n", components.SelectedStyle.Render(m.selectedType))
	content += fmt.Sprintf("Selected branch: %s\n", components.SelectedStyle.Render(m.selectedBranch))
	content += fmt.Sprintf("Ticket number: %s\n", components.SelectedStyle.Render(m.ticketNumber))
	content += fmt.Sprintf("Title: %s\n\n", components.SelectedStyle.Render(m.ticketTitle))
	
	content += "Branch to create:\n"
	content += components.SelectedStyle.Render(m.finalBranch) + "\n\n"
	content += "Press Enter to create this branch"
	
	return content
}

func (m AppModel) renderComplete() string {
	content := "✅ Branch created successfully!\n\n"
	content += fmt.Sprintf("Created: %s\n", components.SelectedStyle.Render(m.finalBranch))
	content += fmt.Sprintf("From: %s\n\n", components.SelectedStyle.Render(m.selectedBranch))
	content += "Press Enter to exit"
	
	return content
}

// Helper methods for state management

// SetError sets an error state
func (m *AppModel) SetError(err error) {
	m.err = err
}

// ClearError clears the error state
func (m *AppModel) ClearError() {
	m.err = nil
}

// GetCurrentState returns the current application state
func (m AppModel) GetCurrentState() AppState {
	return m.state
}

// SetState sets the application state
func (m *AppModel) SetState(state AppState) {
	m.state = state
}

// GetSelectedData returns the currently selected data
func (m AppModel) GetSelectedData() (string, string, string, string) {
	return m.selectedType, m.selectedBranch, m.ticketNumber, m.ticketTitle
}

// SetSelectedData sets the selected data
func (m *AppModel) SetSelectedData(branchType, baseBranch, ticket, title string) {
	m.selectedType = branchType
	m.selectedBranch = baseBranch
	m.ticketNumber = ticket
	m.ticketTitle = title
}