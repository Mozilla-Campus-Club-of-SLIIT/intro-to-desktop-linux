package ui

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	width       int
	height      int
	environment string

	activeModel tea.Model
}

func (m *RootModel) Init() tea.Cmd {
	if m.activeModel != nil {
		return m.activeModel.Init()
	}
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.activeModel != nil {
			m.activeModel, cmd = m.activeModel.Update(msg)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// Delegate all other messages to the active model.
	if m.activeModel != nil {
		m.activeModel, cmd = m.activeModel.Update(msg)
	}

	return m, cmd
}

func (m *RootModel) View() string {
	if m.activeModel != nil {
		return m.activeModel.View()
	}
	return "Initializing..."
}

func Bootstrap(env string) error {
	var activeModel tea.Model

	if auth.VerifyAuth() {
		activeModel = &DashboardModel{
			environment: env,
			content:     ContentModel{environment: env},
			leaderboard: LeaderboardModel{environment: env},
		}
	} else {
		activeModel = &AuthModel{
			environment: env,
		}
	}

	root := &RootModel{
		environment: env,
		activeModel: activeModel,
	}

	_, err := tea.NewProgram(root, tea.WithAltScreen()).Run()
	return err
}
