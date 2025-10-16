package models

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jiraflow/internal/jira"
	"jiraflow/internal/tui/components"
)

// InputField represents the current active input field
type InputField int

const (
	FieldTicketNumber InputField = iota
	FieldTitle
)

// InputFormModel handles ticket number and title input
type InputFormModel struct {
	ticketInput    textinput.Model
	titleInput     textinput.Model
	currentField   InputField
	ticketNumber   string
	ticketTitle    string
	width          int
	height         int
	keyMap         InputFormKeyMap
	jiraClient     jira.JiraClient
	
	// Validation state
	ticketValid    bool
	ticketError    string
	titleFetching  bool
	titleFetched   bool
	titleError     string
	
	// Form completion state
	completed      bool
}

// InputFormKeyMap defines key bindings for the input form
type InputFormKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Back     key.Binding
	Submit   key.Binding
}

// DefaultInputFormKeyMap returns the default key bindings
func DefaultInputFormKeyMap() InputFormKeyMap {
	return InputFormKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "previous field"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "next field"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous field"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit form"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Submit: key.NewBinding(
			key.WithKeys("ctrl+enter"),
			key.WithHelp("ctrl+enter", "submit form"),
		),
	}
}

// NewInputFormModel creates a new input form model
func NewInputFormModel(jiraClient jira.JiraClient) InputFormModel {
	// Create ticket number input
	ticketInput := textinput.New()
	ticketInput.Placeholder = "e.g., JIRA-123, PROJ-456"
	ticketInput.Focus()
	ticketInput.CharLimit = 50
	ticketInput.Width = 30
	
	// Create title input
	titleInput := textinput.New()
	titleInput.Placeholder = "Enter title or leave empty to fetch from Jira"
	titleInput.CharLimit = 100
	titleInput.Width = 50
	
	return InputFormModel{
		ticketInput:  ticketInput,
		titleInput:   titleInput,
		currentField: FieldTicketNumber,
		keyMap:       DefaultInputFormKeyMap(),
		jiraClient:   jiraClient,
		ticketValid:  false,
	}
}

// Init initializes the input form model
func (m InputFormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles events for the input form
func (m InputFormModel) Update(msg tea.Msg) (InputFormModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ticketInput.Width = min(30, msg.Width-10)
		m.titleInput.Width = min(50, msg.Width-10)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Back):
			return m, nil
		case key.Matches(msg, m.keyMap.Enter) || key.Matches(msg, m.keyMap.Submit):
			if m.isFormValid() {
				m.completed = true
				m.ticketNumber = strings.TrimSpace(m.ticketInput.Value())
				m.ticketTitle = strings.TrimSpace(m.titleInput.Value())
			}
			return m, nil
		case key.Matches(msg, m.keyMap.Tab) || key.Matches(msg, m.keyMap.Down):
			m.nextField()
			return m, nil
		case key.Matches(msg, m.keyMap.ShiftTab) || key.Matches(msg, m.keyMap.Up):
			m.previousField()
			return m, nil
		}

		// Update the active input field
		switch m.currentField {
		case FieldTicketNumber:
			oldValue := m.ticketInput.Value()
			m.ticketInput, cmd = m.ticketInput.Update(msg)
			cmds = append(cmds, cmd)
			
			// Validate ticket number if it changed
			newValue := m.ticketInput.Value()
			if oldValue != newValue {
				m.validateTicketNumber(newValue)
				
				// Auto-fetch title if ticket is valid and title is empty
				if m.ticketValid && m.titleInput.Value() == "" && m.jiraClient != nil && m.jiraClient.IsAvailable() {
					m.titleFetching = true
					m.titleFetched = false
					m.titleError = ""
					// In a real implementation, this would trigger an async command
					// For now, we'll simulate it in the validation
					cmds = append(cmds, m.fetchTitleCmd(newValue))
				}
			}
			
		case FieldTitle:
			m.titleInput, cmd = m.titleInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	
	case FetchTitleMsg:
		m.titleFetching = false
		if msg.Error != "" {
			m.titleError = msg.Error
			m.titleFetched = false
		} else {
			m.titleInput.SetValue(msg.Title)
			m.titleFetched = true
			m.titleError = ""
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// FetchTitleMsg represents a message for title fetching results
type FetchTitleMsg struct {
	TicketID string
	Title    string
	Error    string
}

// fetchTitleCmd creates a command to fetch title from Jira
func (m InputFormModel) fetchTitleCmd(ticketID string) tea.Cmd {
	return func() tea.Msg {
		if m.jiraClient == nil || !m.jiraClient.IsAvailable() {
			return FetchTitleMsg{
				TicketID: ticketID,
				Error:    "Jira CLI not available",
			}
		}
		
		title, err := m.jiraClient.GetTicketTitle(ticketID)
		if err != nil {
			return FetchTitleMsg{
				TicketID: ticketID,
				Error:    err.Error(),
			}
		}
		
		return FetchTitleMsg{
			TicketID: ticketID,
			Title:    title,
		}
	}
}

// validateTicketNumber validates the ticket number format
func (m *InputFormModel) validateTicketNumber(value string) {
	value = strings.TrimSpace(value)
	
	if value == "" {
		m.ticketValid = false
		m.ticketError = ""
		return
	}
	
	// Validate ticket format (PROJECT-NUMBER)
	ticketRegex := regexp.MustCompile(`^[A-Z][A-Z0-9]*-\d+$`)
	if !ticketRegex.MatchString(value) {
		m.ticketValid = false
		m.ticketError = "Invalid format. Use PROJECT-123 format (e.g., JIRA-123)"
		return
	}
	
	m.ticketValid = true
	m.ticketError = ""
}

// nextField moves to the next input field
func (m *InputFormModel) nextField() {
	switch m.currentField {
	case FieldTicketNumber:
		m.currentField = FieldTitle
		m.ticketInput.Blur()
		m.titleInput.Focus()
	case FieldTitle:
		// Stay on title field (last field)
	}
}

// previousField moves to the previous input field
func (m *InputFormModel) previousField() {
	switch m.currentField {
	case FieldTicketNumber:
		// Stay on ticket field (first field)
	case FieldTitle:
		m.currentField = FieldTicketNumber
		m.titleInput.Blur()
		m.ticketInput.Focus()
	}
}

// isFormValid checks if the form is valid for submission
func (m InputFormModel) isFormValid() bool {
	return m.ticketValid && strings.TrimSpace(m.ticketInput.Value()) != ""
}

// View renders the input form interface
func (m InputFormModel) View() string {
	var sections []string

	// Title section
	title := components.TitleStyle.Render("Enter Ticket Information")
	sections = append(sections, title)

	// Subtitle with instruction
	subtitle := components.SubtitleStyle.Render("Fill in the ticket details below:")
	sections = append(sections, subtitle)

	// Ticket number field
	ticketSection := m.renderTicketField()
	sections = append(sections, ticketSection)

	// Title field
	titleSection := m.renderTitleField()
	sections = append(sections, titleSection)

	// Status messages
	if statusMsg := m.renderStatusMessages(); statusMsg != "" {
		sections = append(sections, statusMsg)
	}

	// Form validation summary
	if validationMsg := m.renderValidationSummary(); validationMsg != "" {
		sections = append(sections, validationMsg)
	}

	// Help section
	help := m.renderHelp()
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderTicketField renders the ticket number input field
func (m InputFormModel) renderTicketField() string {
	var sections []string

	// Field label
	label := "Ticket Number:"
	if m.currentField == FieldTicketNumber {
		label = components.FocusedStyle.Render("→ " + label)
	} else {
		label = components.UnselectedStyle.Render("  " + label)
	}
	sections = append(sections, label)

	// Input field
	var inputStyle lipgloss.Style
	if m.currentField == FieldTicketNumber {
		inputStyle = components.InputFocusedStyle
	} else {
		inputStyle = components.InputStyle
	}

	inputBox := inputStyle.Render(m.ticketInput.View())
	sections = append(sections, "  "+inputBox)

	// Validation message
	if m.ticketError != "" {
		errorMsg := components.ErrorStyle.Render("  ✗ " + m.ticketError)
		sections = append(sections, errorMsg)
	} else if m.ticketValid {
		successMsg := components.SuccessStyle.Render("  ✓ Valid ticket format")
		sections = append(sections, successMsg)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderTitleField renders the title input field
func (m InputFormModel) renderTitleField() string {
	var sections []string

	// Field label
	label := "Title (optional):"
	if m.currentField == FieldTitle {
		label = components.FocusedStyle.Render("→ " + label)
	} else {
		label = components.UnselectedStyle.Render("  " + label)
	}
	sections = append(sections, label)

	// Input field
	var inputStyle lipgloss.Style
	if m.currentField == FieldTitle {
		inputStyle = components.InputFocusedStyle
	} else {
		inputStyle = components.InputStyle
	}

	inputBox := inputStyle.Render(m.titleInput.View())
	sections = append(sections, "  "+inputBox)

	// Jira integration status
	if m.titleFetching {
		fetchingMsg := components.ProgressStyle.Render("  ⏳ Fetching title from Jira...")
		sections = append(sections, fetchingMsg)
	} else if m.titleFetched {
		fetchedMsg := components.SuccessStyle.Render("  ✓ Title fetched from Jira")
		sections = append(sections, fetchedMsg)
	} else if m.titleError != "" {
		errorMsg := components.WarningStyle.Render("  ⚠ " + m.titleError + " (manual input required)")
		sections = append(sections, errorMsg)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderStatusMessages renders any status messages
func (m InputFormModel) renderStatusMessages() string {
	var messages []string

	// Jira availability status
	if m.jiraClient != nil && m.jiraClient.IsAvailable() {
		if m.ticketValid && m.titleInput.Value() == "" {
			msg := components.HelpStyle.Render("💡 Leave title empty to auto-fetch from Jira")
			messages = append(messages, msg)
		}
	} else {
		msg := components.WarningStyle.Render("⚠ Jira CLI not available - manual title input required")
		messages = append(messages, msg)
	}

	if len(messages) > 0 {
		return lipgloss.JoinVertical(lipgloss.Left, messages...)
	}
	return ""
}

// renderValidationSummary renders form validation summary
func (m InputFormModel) renderValidationSummary() string {
	if !m.isFormValid() {
		var issues []string
		
		if !m.ticketValid || strings.TrimSpace(m.ticketInput.Value()) == "" {
			issues = append(issues, "Valid ticket number required")
		}
		
		if len(issues) > 0 {
			summary := "Form incomplete: " + strings.Join(issues, ", ")
			return components.ErrorStyle.Render("⚠ " + summary)
		}
	} else {
		return components.SuccessStyle.Render("✓ Form ready to submit")
	}
	
	return ""
}

// renderHelp renders the help text
func (m InputFormModel) renderHelp() string {
	var sections []string
	
	// Main help based on current field
	var mainHelp []string
	switch m.currentField {
	case FieldTicketNumber:
		mainHelp = []string{
			"type ticket number (PROJ-123 format)",
			"tab/↓ next field",
			"enter submit form",
		}
	case FieldTitle:
		mainHelp = []string{
			"type title or leave empty",
			"tab/↑ previous field", 
			"enter submit form",
		}
	}
	
	mainHelpText := strings.Join(mainHelp, " • ")
	sections = append(sections, components.HelpStyle.Render(mainHelpText))
	
	// Context help - show form status and Jira availability
	var contextHelp []string
	
	// Form validation status
	if m.isFormValid() {
		contextHelp = append(contextHelp, "✓ Form ready")
	} else {
		contextHelp = append(contextHelp, "Form incomplete")
	}
	
	// Jira status
	if m.jiraClient != nil && m.jiraClient.IsAvailable() {
		contextHelp = append(contextHelp, "Jira available")
	} else {
		contextHelp = append(contextHelp, "Jira unavailable")
	}
	
	contextStyle := components.HelpStyle.
		Foreground(components.ColorMuted).
		Faint(true)
	sections = append(sections, contextStyle.Render(strings.Join(contextHelp, " • ")))
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// GetTicketNumber returns the entered ticket number
func (m InputFormModel) GetTicketNumber() string {
	return m.ticketNumber
}

// GetTicketTitle returns the entered or fetched title
func (m InputFormModel) GetTicketTitle() string {
	return m.ticketTitle
}

// HasCompleted returns true if the form has been completed
func (m InputFormModel) HasCompleted() bool {
	return m.completed
}

// IsValid returns true if the form is valid
func (m InputFormModel) IsValid() bool {
	return m.isFormValid()
}

// GetCurrentField returns the currently active field
func (m InputFormModel) GetCurrentField() InputField {
	return m.currentField
}

// SetSize sets the dimensions of the component
func (m *InputFormModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.ticketInput.Width = min(30, width-10)
	m.titleInput.Width = min(50, width-10)
}

// Reset resets the component state
func (m *InputFormModel) Reset() {
	m.ticketInput.SetValue("")
	m.titleInput.SetValue("")
	m.currentField = FieldTicketNumber
	m.ticketNumber = ""
	m.ticketTitle = ""
	m.ticketValid = false
	m.ticketError = ""
	m.titleFetching = false
	m.titleFetched = false
	m.titleError = ""
	m.completed = false
	
	m.ticketInput.Focus()
	m.titleInput.Blur()
}

// SetTicketNumber sets the ticket number (for pre-filling)
func (m *InputFormModel) SetTicketNumber(ticket string) {
	m.ticketInput.SetValue(ticket)
	m.validateTicketNumber(ticket)
}

// SetTitle sets the title (for pre-filling)
func (m *InputFormModel) SetTitle(title string) {
	m.titleInput.SetValue(title)
}

// FocusTicketField focuses the ticket number field
func (m *InputFormModel) FocusTicketField() {
	m.currentField = FieldTicketNumber
	m.ticketInput.Focus()
	m.titleInput.Blur()
}

// FocusTitleField focuses the title field
func (m *InputFormModel) FocusTitleField() {
	m.currentField = FieldTitle
	m.titleInput.Focus()
	m.ticketInput.Blur()
}

// IsJiraAvailable returns true if Jira client is available
func (m InputFormModel) IsJiraAvailable() bool {
	return m.jiraClient != nil && m.jiraClient.IsAvailable()
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}