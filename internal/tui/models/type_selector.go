package models

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jiraflow/internal/config"
	"jiraflow/internal/tui/components"
)

// TypeItem represents a branch type in the list
type TypeItem struct {
	key         string
	displayName string
	description string
	isDefault   bool
}

// FilterValue returns the value to filter on
func (i TypeItem) FilterValue() string {
	return i.displayName
}

// Title returns the type name for display
func (i TypeItem) Title() string {
	if i.isDefault {
		return components.SelectedStyle.Render(i.displayName + " (default)")
	}
	return i.displayName
}

// Description returns the type description
func (i TypeItem) Description() string {
	return i.description
}

// TypeItemDelegate handles rendering of type items
type TypeItemDelegate struct{}

func (d TypeItemDelegate) Height() int                             { return 2 }
func (d TypeItemDelegate) Spacing() int                            { return 1 }
func (d TypeItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d TypeItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(TypeItem)
	if !ok {
		return
	}

	var str string
	var desc string

	if i.isDefault {
		str = i.displayName + " (default)"
	} else {
		str = i.displayName
	}

	desc = i.description

	fn := components.UnselectedStyle.Render
	descFn := components.HelpStyle.Render
	
	if index == m.Index() {
		fn = func(s ...string) string {
			return components.ListSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
		descFn = func(s ...string) string {
			return components.ListSelectedItemStyle.
				Faint(true).
				Render("  " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprintf(w, "%s\n%s", fn(str), descFn(desc))
}

// TypeSelectorModel handles branch type selection
type TypeSelectorModel struct {
	list      list.Model
	types     []TypeItem
	selected  string
	width     int
	height    int
	keyMap    TypeSelectorKeyMap
}

// TypeSelectorKeyMap defines key bindings for the type selector
type TypeSelectorKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding
}

// DefaultTypeSelectorKeyMap returns the default key bindings
func DefaultTypeSelectorKeyMap() TypeSelectorKeyMap {
	return TypeSelectorKeyMap{
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
	}
}

// NewTypeSelectorModel creates a new type selector model
func NewTypeSelectorModel(cfg *config.Config) TypeSelectorModel {
	// Define descriptions for each branch type
	typeDescriptions := map[string]string{
		"feature":  "New features and enhancements",
		"hotfix":   "Critical bug fixes for production",
		"refactor": "Code improvements without changing functionality",
		"support":  "Supporting changes like documentation or tooling",
	}

	// Convert config branch types to TypeItem
	var typeItems []TypeItem
	var listItems []list.Item

	for key, displayName := range cfg.BranchTypes {
		description := typeDescriptions[key]
		if description == "" {
			description = "Custom branch type"
		}

		item := TypeItem{
			key:         key,
			displayName: displayName,
			description: description,
			isDefault:   key == cfg.DefaultBranchType,
		}
		typeItems = append(typeItems, item)
		listItems = append(listItems, item)
	}

	// Create the list model
	l := list.New(listItems, TypeItemDelegate{}, 0, 0)
	l.Title = "Select Branch Type"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = components.TitleStyle
	l.Styles.PaginationStyle = components.HelpStyle
	l.Styles.HelpStyle = components.HelpStyle

	// Set initial selection to default type if available
	defaultIndex := 0
	for i, item := range typeItems {
		if item.isDefault {
			defaultIndex = i
			break
		}
	}
	l.Select(defaultIndex)

	return TypeSelectorModel{
		list:   l,
		types:  typeItems,
		keyMap: DefaultTypeSelectorKeyMap(),
	}
}

// Init initializes the type selector model
func (m TypeSelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles events for the type selector
func (m TypeSelectorModel) Update(msg tea.Msg) (TypeSelectorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width - 4)
		m.list.SetHeight(msg.Height - 6) // Leave space for title and help
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			selectedItem := m.list.SelectedItem()
			if typeItem, ok := selectedItem.(TypeItem); ok {
				m.selected = typeItem.key
			}
			return m, nil
		case key.Matches(msg, m.keyMap.Back):
			return m, nil
		default:
			// Update list navigation
			m.list, cmd = m.list.Update(msg)
		}
	}

	return m, cmd
}

// View renders the type selector interface
func (m TypeSelectorModel) View() string {
	var sections []string

	// Title section
	title := components.TitleStyle.Render("Select Branch Type")
	sections = append(sections, title)

	// Subtitle with instruction
	subtitle := components.SubtitleStyle.Render("Choose the type of branch you want to create:")
	sections = append(sections, subtitle)

	// Type list section
	listView := m.list.View()
	sections = append(sections, listView)

	// Help section
	help := m.renderHelp()
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderHelp renders the help text
func (m TypeSelectorModel) renderHelp() string {
	var sections []string
	
	// Main help
	mainHelp := []string{
		"↑/↓ or j/k navigate",
		"enter select type",
	}
	
	mainHelpText := strings.Join(mainHelp, " • ")
	sections = append(sections, components.HelpStyle.Render(mainHelpText))
	
	// Context help - show current selection and available types
	var contextHelp []string
	if currentItem, ok := m.GetCurrentItem(); ok {
		contextHelp = append(contextHelp, fmt.Sprintf("Current: %s", currentItem.displayName))
	}
	contextHelp = append(contextHelp, fmt.Sprintf("%d types available", len(m.types)))
	
	contextStyle := components.HelpStyle.
		Foreground(components.ColorMuted).
		Faint(true)
	sections = append(sections, contextStyle.Render(strings.Join(contextHelp, " • ")))
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// GetSelected returns the currently selected branch type key
func (m TypeSelectorModel) GetSelected() string {
	return m.selected
}

// HasSelection returns true if a type has been selected
func (m TypeSelectorModel) HasSelection() bool {
	return m.selected != ""
}

// GetCurrentItem returns the currently highlighted type item
func (m TypeSelectorModel) GetCurrentItem() (TypeItem, bool) {
	selectedItem := m.list.SelectedItem()
	if typeItem, ok := selectedItem.(TypeItem); ok {
		return typeItem, true
	}
	return TypeItem{}, false
}

// GetSelectedDisplayName returns the display name of the selected type
func (m TypeSelectorModel) GetSelectedDisplayName() string {
	for _, item := range m.types {
		if item.key == m.selected {
			return item.displayName
		}
	}
	return m.selected
}

// SetSize sets the dimensions of the component
func (m *TypeSelectorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetWidth(width - 4)
	m.list.SetHeight(height - 6)
}

// Reset resets the component state
func (m *TypeSelectorModel) Reset() {
	m.selected = ""
	// Reset to default selection
	defaultIndex := 0
	for i, item := range m.types {
		if item.isDefault {
			defaultIndex = i
			break
		}
	}
	m.list.Select(defaultIndex)
}

// GetAvailableTypes returns all available type items
func (m TypeSelectorModel) GetAvailableTypes() []TypeItem {
	return m.types
}