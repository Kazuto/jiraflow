package models

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"jiraflow/internal/git"
	"jiraflow/internal/tui/components"
)

// BranchItem represents a branch in the list
type BranchItem struct {
	name      string
	isCurrent bool
}

// FilterValue returns the value to filter on
func (i BranchItem) FilterValue() string {
	return i.name
}

// Title returns the branch name for display
func (i BranchItem) Title() string {
	if i.isCurrent {
		return components.SelectedStyle.Render("* " + i.name + " (current)")
	}
	return i.name
}

// Description returns empty string as we don't need descriptions
func (i BranchItem) Description() string {
	return ""
}

// BranchItemDelegate handles rendering of branch items
type BranchItemDelegate struct{}

func (d BranchItemDelegate) Height() int                             { return 1 }
func (d BranchItemDelegate) Spacing() int                            { return 0 }
func (d BranchItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d BranchItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(BranchItem)
	if !ok {
		return
	}

	str := i.Title()

	fn := components.UnselectedStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return components.ListSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

// BranchSelectorModel handles branch selection with search functionality
type BranchSelectorModel struct {
	list           list.Model
	allBranches    []BranchItem
	filteredItems  []list.Item
	searchInput    textinput.Model
	searching      bool
	selected       string
	width          int
	height         int
	searchResults  git.BranchSearchResult
	keyMap         BranchSelectorKeyMap
}

// BranchSelectorKeyMap defines key bindings for the branch selector
type BranchSelectorKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Search key.Binding
	Clear  key.Binding
}

// DefaultBranchSelectorKeyMap returns the default key bindings
func DefaultBranchSelectorKeyMap() BranchSelectorKeyMap {
	return BranchSelectorKeyMap{
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
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Clear: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "clear search"),
		),
	}
}

// NewBranchSelectorModel creates a new branch selector model
func NewBranchSelectorModel(branches []git.BranchInfo) BranchSelectorModel {
	// Convert BranchInfo to BranchItem
	var branchItems []BranchItem
	var listItems []list.Item

	for _, branch := range branches {
		item := BranchItem{
			name:      branch.Name,
			isCurrent: branch.IsCurrent,
		}
		branchItems = append(branchItems, item)
		listItems = append(listItems, item)
	}

	// Create the list model
	l := list.New(listItems, BranchItemDelegate{}, 0, 0)
	l.Title = "Select Base Branch"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false) // We'll handle filtering ourselves
	l.Styles.Title = components.TitleStyle
	l.Styles.PaginationStyle = components.HelpStyle
	l.Styles.HelpStyle = components.HelpStyle

	// Create search input
	ti := textinput.New()
	ti.Placeholder = "Type to search branches..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 50

	return BranchSelectorModel{
		list:          l,
		allBranches:   branchItems,
		filteredItems: listItems,
		searchInput:   ti,
		searching:     false,
		keyMap:        DefaultBranchSelectorKeyMap(),
		searchResults: git.BranchSearchResult{HasResults: true},
	}
}

// Init initializes the branch selector model
func (m BranchSelectorModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles events for the branch selector
func (m BranchSelectorModel) Update(msg tea.Msg) (BranchSelectorModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width - 4)
		m.list.SetHeight(msg.Height - 8) // Leave space for search input and help
		m.searchInput.Width = msg.Width - 10
		return m, nil

	case tea.KeyMsg:
		// Handle search mode
		if m.searching {
			switch {
			case key.Matches(msg, m.keyMap.Back):
				m.searching = false
				m.searchInput.Blur()
				return m, nil
			case key.Matches(msg, m.keyMap.Enter):
				m.searching = false
				m.searchInput.Blur()
				return m, nil
			case key.Matches(msg, m.keyMap.Clear):
				m.searchInput.SetValue("")
				m.updateFilter("")
				return m, nil
			default:
				// Update search input
				m.searchInput, cmd = m.searchInput.Update(msg)
				cmds = append(cmds, cmd)
				
				// Update filter based on search input
				searchTerm := m.searchInput.Value()
				m.updateFilter(searchTerm)
				return m, tea.Batch(cmds...)
			}
		}

		// Handle navigation mode
		switch {
		case key.Matches(msg, m.keyMap.Search):
			m.searching = true
			m.searchInput.Focus()
			return m, textinput.Blink
		case key.Matches(msg, m.keyMap.Enter):
			if len(m.filteredItems) > 0 {
				selectedItem := m.list.SelectedItem()
				if branchItem, ok := selectedItem.(BranchItem); ok {
					m.selected = branchItem.name
				}
			}
			return m, nil
		case key.Matches(msg, m.keyMap.Back):
			return m, nil
		default:
			// Update list navigation
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// updateFilter updates the filtered list based on search term
func (m *BranchSelectorModel) updateFilter(searchTerm string) {
	// Extract branch names for filtering
	var branchNames []string
	for _, branch := range m.allBranches {
		branchNames = append(branchNames, branch.name)
	}

	// Filter branches using git package functionality
	m.searchResults = git.FilterBranchesRealtime(branchNames, searchTerm)

	// Convert filtered results back to list items
	var filteredItems []list.Item
	for _, branchName := range m.searchResults.Branches {
		// Find the corresponding BranchItem
		for _, branch := range m.allBranches {
			if branch.name == branchName {
				filteredItems = append(filteredItems, branch)
				break
			}
		}
	}

	m.filteredItems = filteredItems
	m.list.SetItems(filteredItems)

	// Reset selection to first item if available
	if len(filteredItems) > 0 {
		m.list.Select(0)
	}
}

// View renders the branch selector interface
func (m BranchSelectorModel) View() string {
	var sections []string

	// Title section
	title := components.TitleStyle.Render("Select Base Branch")
	sections = append(sections, title)

	// Search input section
	searchSection := m.renderSearchSection()
	sections = append(sections, searchSection)

	// Results summary
	if m.searchInput.Value() != "" {
		summary := m.renderSearchSummary()
		sections = append(sections, summary)
	}

	// Branch list section
	if len(m.filteredItems) == 0 && m.searchInput.Value() != "" {
		// Show "no results" message
		noResults := components.ErrorStyle.Render(m.searchResults.GetEmptySearchMessage())
		sections = append(sections, noResults)
	} else {
		// Show the list
		listView := m.list.View()
		sections = append(sections, listView)
	}

	// Help section
	help := m.renderHelp()
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderSearchSection renders the search input area
func (m BranchSelectorModel) renderSearchSection() string {
	searchLabel := "Search: "
	
	var searchStyle lipgloss.Style
	if m.searching {
		searchStyle = components.InputFocusedStyle
	} else {
		searchStyle = components.InputStyle
	}

	searchBox := searchStyle.Render(m.searchInput.View())
	
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		components.SubtitleStyle.Render(searchLabel),
		searchBox,
	)
}

// renderSearchSummary renders the search results summary
func (m BranchSelectorModel) renderSearchSummary() string {
	summary := m.searchResults.GetSearchSummary()
	if summary == "" {
		return ""
	}
	
	style := components.HelpStyle
	if !m.searchResults.HasResults {
		style = components.WarningStyle
	}
	
	return style.Render(summary)
}

// renderHelp renders the help text
func (m BranchSelectorModel) renderHelp() string {
	var sections []string
	
	// Main help based on current mode
	var mainHelp []string
	if m.searching {
		mainHelp = []string{
			"type to search branches",
			"enter/esc finish search",
			"ctrl+u clear search",
		}
	} else {
		mainHelp = []string{
			"↑/↓ or j/k navigate",
			"/ start search",
			"enter select branch",
		}
	}
	
	mainHelpText := strings.Join(mainHelp, " • ")
	sections = append(sections, components.HelpStyle.Render(mainHelpText))
	
	// Additional context help
	var contextHelp []string
	if !m.searching {
		if len(m.filteredItems) == 0 {
			contextHelp = append(contextHelp, "No branches available")
		} else {
			contextHelp = append(contextHelp, fmt.Sprintf("%d branches available", len(m.filteredItems)))
		}
	} else {
		searchTerm := m.searchInput.Value()
		if searchTerm != "" {
			contextHelp = append(contextHelp, fmt.Sprintf("Searching for: %s", searchTerm))
		} else {
			contextHelp = append(contextHelp, "Type to filter branches")
		}
	}
	
	if len(contextHelp) > 0 {
		contextStyle := components.HelpStyle.
			Foreground(components.ColorMuted).
			Faint(true)
		sections = append(sections, contextStyle.Render(strings.Join(contextHelp, " • ")))
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// GetSelected returns the currently selected branch name
func (m BranchSelectorModel) GetSelected() string {
	return m.selected
}

// HasSelection returns true if a branch has been selected
func (m BranchSelectorModel) HasSelection() bool {
	return m.selected != ""
}

// GetCurrentItem returns the currently highlighted branch item
func (m BranchSelectorModel) GetCurrentItem() (BranchItem, bool) {
	if len(m.filteredItems) == 0 {
		return BranchItem{}, false
	}
	
	selectedItem := m.list.SelectedItem()
	if branchItem, ok := selectedItem.(BranchItem); ok {
		return branchItem, true
	}
	
	return BranchItem{}, false
}

// SetSize sets the dimensions of the component
func (m *BranchSelectorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetWidth(width - 4)
	m.list.SetHeight(height - 8)
	m.searchInput.Width = width - 10
}

// Reset resets the component state
func (m *BranchSelectorModel) Reset() {
	m.selected = ""
	m.searching = false
	m.searchInput.SetValue("")
	m.searchInput.Blur()
	m.updateFilter("")
}

// IsSearching returns true if the component is in search mode
func (m BranchSelectorModel) IsSearching() bool {
	return m.searching
}

// GetSearchTerm returns the current search term
func (m BranchSelectorModel) GetSearchTerm() string {
	return m.searchInput.Value()
}

// GetBranchCount returns the number of available branches
func (m BranchSelectorModel) GetBranchCount() int {
	return len(m.filteredItems)
}