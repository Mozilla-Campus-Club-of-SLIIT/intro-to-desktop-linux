package ui

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine"
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	width       int
	height      int
	environment string

	content     ContentModel
	leaderboard LeaderboardModel
	resizeState ResizeState
}

func (m *RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	resizeCmd := m.resizeState.HandleResizeUpdate(msg)
	if resizeCmd != nil {
		cmd = tea.Batch(cmd, resizeCmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m *RootModel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "Initializing..."
	}

	if m.resizeState.IsResizing {
		return engine.ClearAllImages() + resizeView(m.width, m.height)
	}

	return m.render()
}

func Bootstrap(env string) error {
	m := &RootModel{
		environment: env,
		content:     ContentModel{environment: env},
		leaderboard: LeaderboardModel{environment: env},
	}

	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}
