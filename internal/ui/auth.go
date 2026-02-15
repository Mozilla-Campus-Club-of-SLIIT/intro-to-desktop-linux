package ui

import tea "github.com/charmbracelet/bubbletea"

type AuthModel struct {
	width       int
	height      int
	environment string
}

func (m *AuthModel) Init() tea.Cmd {
	return nil
}

func (m *AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *AuthModel) View() string {
	return "auth"
}
