package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jiraflow/internal/config"
	"jiraflow/internal/errors"
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
	state            AppState
	config           *config.Config
	git              git.GitRepository
	typeModel        models.TypeSelectorModel
	branchModel      models.BranchSelectorModel
	inputModel       models.InputFormModel
	confirmationModel models.ConfirmationModel
	completionModel  models.CompletionModel
	err              error
	errorHandler     *errors.ErrorHandler
	degradationHandler *errors.DegradationHandler
	width            int
	height           int
	
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
	// Initialize type selector model
	typeModel := models.NewTypeSelectorModel(cfg)
	
	// Initialize branch selector with Git branches
	var branchModel models.BranchSelectorModel
	if branches, err := gitRepo.GetBranchesWithInfo(); err == nil {
		branchModel = models.NewBranchSelectorModel(branches)
	} else {
		// Fallback to empty branch selector if Git operations fail
		branchModel = models.NewBranchSelectorModel([]git.BranchInfo{})
		// Note: Error will be handled gracefully during runtime
	}
	
	// Initialize input form model (Jira client will be set later if available)
	inputModel := models.NewInputFormModel(nil)
	
	// Initialize confirmation and completion models
	confirmationModel := models.NewConfirmationModel()
	completionModel := models.NewCompletionModel()
	
	return &AppModel{
		state:              StateTypeSelection,
		config:             cfg,
		git:                gitRepo,
		typeModel:          typeModel,
		branchModel:        branchModel,
		inputModel:         inputModel,
		confirmationModel:  confirmationModel,
		completionModel:    completionModel,
		errorHandler:       errors.NewErrorHandler(),
		degradationHandler: errors.NewDegradationHandler(),
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
	
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI application error: %w", err)
	}
	
	// Check if the final model has an error state
	if appModel, ok := finalModel.(AppModel); ok {
		if appModel.err != nil {
			return appModel.err
		}
		
		// Check if user quit without completing the workflow
		if appModel.state != StateComplete {
			return fmt.Errorf("user cancelled operation")
		}
	}
	
	return nil
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
		m.typeModel.SetSize(msg.Width, msg.Height-6) // Leave space for header/footer
		m.branchModel.SetSize(msg.Width, msg.Height-6)
		m.inputModel.SetSize(msg.Width, msg.Height-6)
		m.confirmationModel.SetSize(msg.Width, msg.Height-6)
		m.completionModel.SetSize(msg.Width, msg.Height-6)
		
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

// State-specific update methods
func (m AppModel) updateTypeSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	// Update the type selector model
	m.typeModel, cmd = m.typeModel.Update(msg)
	
	// Check if a type was selected
	if m.typeModel.HasSelection() {
		m.selectedType = m.typeModel.GetSelected()
		m.state = StateBranchSelection
		return m, cmd
	}
	
	return m, cmd
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
	var cmd tea.Cmd
	
	// Update the confirmation model
	updatedConfirmation, confirmCmd := m.confirmationModel.Update(msg)
	m.confirmationModel = updatedConfirmation
	cmd = confirmCmd
	
	// Check if user confirmed
	if m.confirmationModel.HasConfirmed() {
		// Generate final branch name
		m.finalBranch = m.generateBranchName()
		
		// Set confirmation data
		m.confirmationModel.SetData(
			m.selectedType,
			m.selectedBranch,
			m.ticketNumber,
			m.ticketTitle,
			m.finalBranch,
		)
		
		// Attempt to create the branch
		if err := m.createBranch(); err != nil {
			// Set error state in completion model
			m.completionModel.SetError(err.Error())
		} else {
			// Set success state in completion model
			m.completionModel.SetSuccess(m.finalBranch, m.selectedBranch)
		}
		
		m.state = StateComplete
		return m, cmd
	}
	
	return m, cmd
}

func (m AppModel) updateComplete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	// Update the completion model
	updatedCompletion, completeCmd := m.completionModel.Update(msg)
	m.completionModel = updatedCompletion
	cmd = completeCmd
	
	// Check if user wants to exit
	if m.completionModel.ShouldExit() {
		return m, tea.Quit
	}
	
	return m, cmd
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
		stateText = "Step 1/4: Select Branch Type"
	case StateBranchSelection:
		stateText = "Step 2/4: Select Base Branch"
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
	help := m.renderContextualHelp()
	
	// Add a separator line above the help
	separator := strings.Repeat("─", m.width)
	separatorStyle := lipgloss.NewStyle().Foreground(components.ColorMuted)
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		separatorStyle.Render(separator),
		help,
	)
}

// renderContextualHelp renders context-sensitive help based on current screen
func (m AppModel) renderContextualHelp() string {
	helpRenderer := components.NewHelpRenderer(m.width)
	
	// Get key bindings for current state
	bindings := m.getKeyBindings()
	
	// Get context information
	contextInfo := m.getContextInfo()
	
	// Render using the help renderer
	if len(bindings) > 0 {
		return helpRenderer.RenderKeyBindings(bindings)
	}
	
	// Fallback to simple context help
	mainHelp := m.getMainHelpText()
	return helpRenderer.RenderContextualHelp(mainHelp, contextInfo)
}

// getMainHelpText returns the main help text for the current state
func (m AppModel) getMainHelpText() string {
	var helpKeys []string
	
	switch m.state {
	case StateTypeSelection:
		helpKeys = []string{
			"↑/↓ or j/k navigate",
			"enter select type",
		}
		
	case StateBranchSelection:
		// Check if branch selector is in search mode
		isSearching := m.branchModel.IsSearching()
		if isSearching {
			helpKeys = []string{
				"type to search branches",
				"enter finish search",
				"ctrl+u clear search",
				"esc cancel search",
			}
		} else {
			helpKeys = []string{
				"↑/↓ or j/k navigate",
				"/ start search",
				"enter select branch",
				"esc back to type selection",
			}
		}
		
	case StateTicketInput, StateTitleInput:
		currentField := m.inputModel.GetCurrentField()
		if currentField == models.FieldTicketNumber {
			helpKeys = []string{
				"type ticket number (e.g., PROJ-123)",
				"tab/↓ next field",
				"enter submit form",
				"esc back to branch selection",
			}
		} else {
			helpKeys = []string{
				"type title or leave empty for auto-fetch",
				"tab/↑ previous field",
				"enter submit form",
				"esc back to branch selection",
			}
		}
		
	case StateConfirmation:
		helpKeys = []string{
			"enter create branch",
			"esc back to ticket input",
		}
		
	case StateComplete:
		helpKeys = []string{
			"enter exit application",
		}
	}
	
	if len(helpKeys) > 0 {
		return strings.Join(helpKeys, " • ")
	}
	return ""
}

// getKeyBindings returns key bindings for the current state
func (m AppModel) getKeyBindings() []components.KeyBinding {
	var bindings []components.KeyBinding
	
	switch m.state {
	case StateTypeSelection:
		bindings = []components.KeyBinding{
			{Keys: []string{"↑/↓", "j/k"}, Description: "navigate"},
			{Keys: []string{"enter"}, Description: "select type"},
			{Keys: []string{"q"}, Description: "quit", Global: true},
			{Keys: []string{"ctrl+c"}, Description: "force quit", Global: true},
		}
		
	case StateBranchSelection:
		isSearching := m.branchModel.IsSearching()
		if isSearching {
			bindings = []components.KeyBinding{
				{Keys: []string{"type"}, Description: "search branches"},
				{Keys: []string{"enter"}, Description: "finish search"},
				{Keys: []string{"ctrl+u"}, Description: "clear search"},
				{Keys: []string{"esc"}, Description: "cancel search"},
				{Keys: []string{"q"}, Description: "quit", Global: true},
				{Keys: []string{"ctrl+c"}, Description: "force quit", Global: true},
			}
		} else {
			bindings = []components.KeyBinding{
				{Keys: []string{"↑/↓", "j/k"}, Description: "navigate"},
				{Keys: []string{"/"}, Description: "start search"},
				{Keys: []string{"enter"}, Description: "select branch"},
				{Keys: []string{"esc"}, Description: "back to type selection"},
				{Keys: []string{"q"}, Description: "quit", Global: true},
				{Keys: []string{"ctrl+c"}, Description: "force quit", Global: true},
			}
		}
		
	case StateTicketInput, StateTitleInput:
		currentField := m.inputModel.GetCurrentField()
		if currentField == models.FieldTicketNumber {
			bindings = []components.KeyBinding{
				{Keys: []string{"type"}, Description: "enter ticket number (PROJ-123)"},
				{Keys: []string{"tab", "↓"}, Description: "next field"},
				{Keys: []string{"enter"}, Description: "submit form"},
				{Keys: []string{"esc"}, Description: "back to branch selection"},
				{Keys: []string{"q"}, Description: "quit", Global: true},
				{Keys: []string{"ctrl+c"}, Description: "force quit", Global: true},
			}
		} else {
			bindings = []components.KeyBinding{
				{Keys: []string{"type"}, Description: "enter title or leave empty"},
				{Keys: []string{"tab", "↑"}, Description: "previous field"},
				{Keys: []string{"enter"}, Description: "submit form"},
				{Keys: []string{"esc"}, Description: "back to branch selection"},
				{Keys: []string{"q"}, Description: "quit", Global: true},
				{Keys: []string{"ctrl+c"}, Description: "force quit", Global: true},
			}
		}
		
	case StateConfirmation:
		bindings = []components.KeyBinding{
			{Keys: []string{"enter"}, Description: "create branch"},
			{Keys: []string{"esc"}, Description: "back to edit"},
			{Keys: []string{"q"}, Description: "quit", Global: true},
			{Keys: []string{"ctrl+c"}, Description: "force quit", Global: true},
		}
		
	case StateComplete:
		bindings = []components.KeyBinding{
			{Keys: []string{"enter"}, Description: "exit application"},
			{Keys: []string{"q"}, Description: "quit", Global: true},
			{Keys: []string{"ctrl+c"}, Description: "force quit", Global: true},
		}
	}
	
	return bindings
}

// getContextInfo returns context information for the current state
func (m AppModel) getContextInfo() []string {
	var info []string
	
	switch m.state {
	case StateTypeSelection:
		if currentItem, ok := m.typeModel.GetCurrentItem(); ok {
			info = append(info, fmt.Sprintf("Current: %s", currentItem.Title()))
		}
		info = append(info, fmt.Sprintf("%d types available", len(m.typeModel.GetAvailableTypes())))
		
	case StateBranchSelection:
		if m.branchModel.IsSearching() {
			searchTerm := m.branchModel.GetSearchTerm()
			if searchTerm != "" {
				info = append(info, fmt.Sprintf("Searching: %s", searchTerm))
			} else {
				info = append(info, "Type to filter branches")
			}
		} else {
			info = append(info, fmt.Sprintf("%d branches available", m.branchModel.GetBranchCount()))
		}
		
	case StateTicketInput, StateTitleInput:
		if m.inputModel.IsValid() {
			info = append(info, "✓ Form ready")
		} else {
			info = append(info, "Form incomplete")
		}
		
		if m.inputModel.IsJiraAvailable() {
			info = append(info, "Jira available")
		} else {
			info = append(info, "Jira unavailable")
		}
		
	case StateConfirmation:
		info = append(info, "Review details before creating")
		info = append(info, "Branch will be created and checked out")
		
	case StateComplete:
		if m.completionModel.GetState() == models.CompletionSuccess {
			info = append(info, "Branch created successfully")
			info = append(info, "Ready to start development")
		} else {
			info = append(info, "Branch creation failed")
			info = append(info, "Check troubleshooting tips")
		}
	}
	
	return info
}

// getGlobalHelpText returns help text for globally available shortcuts (legacy method)
func (m AppModel) getGlobalHelpText() string {
	var globalKeys []string
	
	// Always available shortcuts
	globalKeys = append(globalKeys, "q quit")
	globalKeys = append(globalKeys, "ctrl+c force quit")
	
	// Add state-specific global shortcuts
	switch m.state {
	case StateTypeSelection:
		// No additional global shortcuts for first screen
	default:
		// All other screens can go back
		if m.state != StateComplete {
			// Don't show esc for complete screen as it's handled differently
			if m.state == StateConfirmation {
				globalKeys = append(globalKeys, "esc back")
			}
		}
	}
	
	return "Global: " + strings.Join(globalKeys, " • ")
}

// renderError renders error messages using the error handler
func (m AppModel) renderError() string {
	return m.errorHandler.FormatErrorForTUI(m.err)
}

// State-specific render methods (placeholder implementations)
func (m AppModel) renderTypeSelection() string {
	return m.typeModel.View()
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
	// Generate the final branch name for display
	finalBranch := m.generateBranchName()
	
	// Create a copy of the confirmation model with the data set
	confirmationCopy := m.confirmationModel
	confirmationCopy.SetData(
		m.selectedType,
		m.selectedBranch,
		m.ticketNumber,
		m.ticketTitle,
		finalBranch,
	)
	
	return confirmationCopy.View()
}

func (m AppModel) renderComplete() string {
	return m.completionModel.View()
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

// generateBranchName generates the final branch name based on selected data
func (m AppModel) generateBranchName() string {
	// Basic branch name generation logic
	// Format: type/ticket-sanitized-title
	
	sanitizedTitle := m.sanitizeTitle(m.ticketTitle)
	
	var branchName string
	if sanitizedTitle != "" {
		branchName = fmt.Sprintf("%s/%s-%s", m.selectedType, m.ticketNumber, sanitizedTitle)
	} else {
		branchName = fmt.Sprintf("%s/%s", m.selectedType, m.ticketNumber)
	}
	
	// Apply length limit from config
	maxLength := m.config.MaxBranchLength
	if len(branchName) > maxLength {
		// Truncate while preserving the type and ticket number
		prefix := fmt.Sprintf("%s/%s", m.selectedType, m.ticketNumber)
		if len(prefix) < maxLength {
			remainingLength := maxLength - len(prefix) - 1 // -1 for the dash
			if remainingLength > 0 && sanitizedTitle != "" {
				truncatedTitle := sanitizedTitle
				if len(truncatedTitle) > remainingLength {
					truncatedTitle = truncatedTitle[:remainingLength]
				}
				branchName = fmt.Sprintf("%s-%s", prefix, truncatedTitle)
			} else {
				branchName = prefix
			}
		} else {
			// If even the prefix is too long, just use it as is
			branchName = prefix
		}
	}
	
	return branchName
}

// sanitizeTitle sanitizes the title for use in branch names
func (m AppModel) sanitizeTitle(title string) string {
	if title == "" {
		return ""
	}
	
	// Basic sanitization rules
	sanitized := strings.ToLower(title)
	
	// Replace spaces with hyphens
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	
	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range sanitized {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	sanitized = result.String()
	
	// Remove multiple consecutive hyphens
	for strings.Contains(sanitized, "--") {
		sanitized = strings.ReplaceAll(sanitized, "--", "-")
	}
	
	// Remove leading and trailing hyphens
	sanitized = strings.Trim(sanitized, "-")
	
	return sanitized
}

// createBranch creates the new Git branch
func (m AppModel) createBranch() error {
	// Create and checkout the new branch
	return m.git.CreateBranch(m.finalBranch, m.selectedBranch)
}

// handleGracefulDegradation handles errors gracefully and provides fallback options
func (m *AppModel) handleGracefulDegradation(err error) {
	if err == nil {
		return
	}

	// Handle different types of errors with graceful degradation
	if jfErr, ok := err.(errors.JiraFlowError); ok {
		switch jfErr.Type() {
		case errors.ErrorTypeJira:
			// Jira integration failed - continue with manual title entry
			// This is already handled in the input form model
			return
		case errors.ErrorTypeGit:
			// Git operations failed - this is more serious
			if !jfErr.IsRecoverable() {
				m.SetError(err)
				return
			}
			// For recoverable Git errors, show warning but continue
			return
		case errors.ErrorTypeConfig:
			// Config issues - use defaults and continue
			if jfErr.IsRecoverable() {
				// Configuration will fall back to defaults
				return
			}
			m.SetError(err)
			return
		}
	}

	// For other errors, set error state
	m.SetError(err)
}

// showDegradationWarning displays a warning about degraded functionality
func (m AppModel) showDegradationWarning(warningType errors.ErrorType) string {
	switch warningType {
	case errors.ErrorTypeJira:
		return m.errorHandler.FormatWarningForTUI("Jira integration unavailable - you can enter ticket titles manually")
	case errors.ErrorTypeGit:
		return m.errorHandler.FormatWarningForTUI("Some Git features may be limited")
	case errors.ErrorTypeConfig:
		return m.errorHandler.FormatWarningForTUI("Using default configuration values")
	default:
		return m.errorHandler.FormatWarningForTUI("Some features may be limited")
	}
}