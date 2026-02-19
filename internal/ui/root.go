package ui

import (
	"fmt"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidth  = 115
	minHeight = 35
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
		if newModel, cmd := keyboardShortcuts(m, msg); cmd != nil {
			return newModel, cmd
		}

	case switchToDashboardMsg:
		if s, ok := m.activeModel.(interface{ Stop() }); ok {
			s.Stop()
		}
		m.activeModel = NewDashboardModel()

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
	if m.width > 0 && m.height > 0 && (m.width < minWidth || m.height < minHeight) {
		return fmt.Sprintf("Is this a terminal for ants? 🤏 \nGive me some space or I'm not showing you anything  😤 \nDecrease your font size or stretch the window! \n\n[Width: %d, Height: %d]", m.width, m.height)
	}

	if m.activeModel != nil {
		return m.activeModel.View()
	}
	return "Initializing... (If you see this for more than a few seconds, something went wrong. 🐧)"
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
