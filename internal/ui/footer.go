package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type shortcut struct {
	key    string
	label  string
	action tea.Cmd
}

var shortcuts = []shortcut{
	{key: "tab", label: "[tab] autocomplete", action: nil},
	{key: "up", label: "[↑/↓] history", action: nil},
	{key: "down", label: "", action: nil},
	{key: "enter", label: "[enter] exec", action: nil},
	{key: "ctrl+shift+l", label: "[ctrl+L] logout", action: nil},
	{key: "ctrl+e", label: "[ctrl+e] quit", action: tea.Quit},
}

func keyboardShortcuts(m tea.Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	for _, s := range shortcuts {
		if s.key == msg.String() && s.action != nil {
			return m, s.action
		}
	}
	return m, nil
}

func (m *DashboardModel) footerView(width int) string {
	var items []string
	for _, s := range shortcuts {
		if s.label != "" {
			items = append(items, s.label)
		}
	}
	return FooterStyle.Width(width).Render(" " + strings.Join(items, "\t"))
}
