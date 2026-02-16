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
	{key: "ctrl+e", label: "[ctrl+e] quit", action: tea.Quit},
	{key: "up", label: "[↑/↓] nav", action: nil},
	{key: "down", label: "", action: nil},
	{key: "left", label: "[←/→] move", action: nil},
	{key: "right", label: "", action: nil},
	{key: "enter", label: "[enter] exec", action: nil},
	{key: "ctrl+p", label: "[ctrl+p] menu", action: nil},
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
