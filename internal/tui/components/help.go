package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpSection represents a section of help text
type HelpSection struct {
	Title string
	Items []string
}

// KeyBinding represents a keyboard shortcut
type KeyBinding struct {
	Keys        []string
	Description string
	Global      bool
}

// HelpRenderer provides utilities for rendering help text
type HelpRenderer struct {
	width int
}

// NewHelpRenderer creates a new help renderer
func NewHelpRenderer(width int) *HelpRenderer {
	return &HelpRenderer{width: width}
}

// RenderKeyBindings renders a list of key bindings in a formatted way
func (h *HelpRenderer) RenderKeyBindings(bindings []KeyBinding) string {
	if len(bindings) == 0 {
		return ""
	}

	var sections []string
	var mainBindings []KeyBinding
	var globalBindings []KeyBinding

	// Separate main and global bindings
	for _, binding := range bindings {
		if binding.Global {
			globalBindings = append(globalBindings, binding)
		} else {
			mainBindings = append(mainBindings, binding)
		}
	}

	// Render main bindings
	if len(mainBindings) > 0 {
		mainHelp := h.renderBindingList(mainBindings)
		sections = append(sections, HelpStyle.Render(mainHelp))
	}

	// Render global bindings
	if len(globalBindings) > 0 {
		globalHelp := h.renderBindingList(globalBindings)
		globalStyle := HelpStyle.
			Foreground(ColorMuted).
			Faint(true)
		sections = append(sections, globalStyle.Render("Global: "+globalHelp))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderBindingList renders a list of bindings as a single line
func (h *HelpRenderer) renderBindingList(bindings []KeyBinding) string {
	var items []string
	for _, binding := range bindings {
		keys := strings.Join(binding.Keys, "/")
		item := keys + " " + binding.Description
		items = append(items, item)
	}
	return strings.Join(items, " • ")
}

// RenderHelpSections renders multiple help sections
func (h *HelpRenderer) RenderHelpSections(sections []HelpSection) string {
	var rendered []string

	for _, section := range sections {
		sectionText := h.renderHelpSection(section)
		if sectionText != "" {
			rendered = append(rendered, sectionText)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

// renderHelpSection renders a single help section
func (h *HelpRenderer) renderHelpSection(section HelpSection) string {
	if len(section.Items) == 0 {
		return ""
	}

	var parts []string

	// Add title if provided
	if section.Title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)
		parts = append(parts, titleStyle.Render(section.Title+":"))
	}

	// Add items
	itemText := strings.Join(section.Items, " • ")
	parts = append(parts, HelpStyle.Render(itemText))

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// RenderContextualHelp renders help text with context information
func (h *HelpRenderer) RenderContextualHelp(mainHelp string, contextInfo []string) string {
	var sections []string

	// Main help
	if mainHelp != "" {
		sections = append(sections, HelpStyle.Render(mainHelp))
	}

	// Context information
	if len(contextInfo) > 0 {
		contextText := strings.Join(contextInfo, " • ")
		contextStyle := HelpStyle.
			Foreground(ColorMuted).
			Faint(true)
		sections = append(sections, contextStyle.Render(contextText))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// RenderSeparator renders a separator line
func (h *HelpRenderer) RenderSeparator() string {
	if h.width <= 0 {
		return ""
	}
	
	separator := strings.Repeat("─", h.width)
	return lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render(separator)
}

// SetWidth sets the width for the help renderer
func (h *HelpRenderer) SetWidth(width int) {
	h.width = width
}

// GetWidth returns the current width
func (h *HelpRenderer) GetWidth() int {
	return h.width
}