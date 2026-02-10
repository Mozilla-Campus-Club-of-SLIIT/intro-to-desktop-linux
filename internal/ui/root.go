package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	width       int
	height      int
	environment string
}

func (m *RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *RootModel) View() string {
	return m.render()
}

func Bootstrap(env string) error {
	m := &RootModel{
		environment: env,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
