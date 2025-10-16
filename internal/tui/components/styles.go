package components

import "github.com/charmbracelet/lipgloss"

// Color palette
const (
	ColorPrimary   = lipgloss.Color("86")  // Bright green
	ColorSecondary = lipgloss.Color("39")  // Bright blue
	ColorError     = lipgloss.Color("196") // Bright red
	ColorWarning   = lipgloss.Color("214") // Orange
	ColorMuted     = lipgloss.Color("241") // Gray
	ColorSuccess   = lipgloss.Color("46")  // Bright green
)

// Common styles for the TUI interface
var (
	// Title and header styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	// Interactive element styles
	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	FocusedStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	// Status and feedback styles
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	// Layout styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)

	ContentStyle = lipgloss.NewStyle().
			Padding(0, 2)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMuted).
			Padding(1, 2)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(ColorSecondary).
			Padding(0, 1)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	ListSelectedItemStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Background(lipgloss.Color("236")).
			Bold(true).
			Padding(0, 2)

	// Progress and state styles
	ProgressStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	CompletedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Strikethrough(true)
)