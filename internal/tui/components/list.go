package components

import "github.com/charmbracelet/bubbles/list"

// NewList creates a new list component with default styling
func NewList(items []list.Item, title string) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	return l
}