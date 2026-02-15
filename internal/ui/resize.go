package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resizeTimeoutMsg struct{ id int }

type ResizeState struct {
	IsResizing bool
	Counter    int
}

func (r *ResizeState) HandleResizeUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.IsResizing = true
		r.Counter++

		currentID := r.Counter

		return tea.Tick(100*time.Millisecond, func(_ time.Time) tea.Msg {
			return resizeTimeoutMsg{id: currentID}
		})

	case resizeTimeoutMsg:
		if msg.id == r.Counter {
			r.IsResizing = false
		}
	}
	return nil
}

func resizeView(width, height int) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			Foreground(ColorLemonYellow).
			Bold(true).
			Render("[ Resizing... ]"),
	)
}
