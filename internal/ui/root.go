package ui

import (
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	width  int
	height int

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

	case switchToDashboardMsg:
		m.activeModel = NewDashboardModel()

		// Send window size to new model so it can render properly
		if m.width > 0 && m.height > 0 {
			m.activeModel, cmd = m.activeModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
		return m, tea.Batch(cmd, m.activeModel.Init())

	}

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

func Bootstrap() error {
	var activeModel tea.Model

	if auth.VerifyAuth() {
		activeModel = NewDashboardModel()
	} else {
		activeModel = NewAuthModel()
	}

	root := &RootModel{
		activeModel: activeModel,
	}

	_, err := tea.NewProgram(root, tea.WithAltScreen()).Run()
	return err
}
