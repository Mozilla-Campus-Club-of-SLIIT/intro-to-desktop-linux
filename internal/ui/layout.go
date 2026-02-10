package ui

import "github.com/charmbracelet/lipgloss"

func (m *RootModel) render() string {
	// Footer height is 1 cell
	footerHeight := 1
	mainHeight := m.height - footerHeight

	// Sidebar is 25% (4/16), Content is 75% (12/16)
	sidebarWidth := m.width / 4
	contentWidth := m.width - sidebarWidth

	sidebar := BorderStyle.
		Width(sidebarWidth - 2). // -2 for borders
		Height(mainHeight - 2).  // -2 for borders
		Render(m.leaderboardView())

	content := BorderStyle.
		Width(contentWidth - 2). // -2 for borders
		Height(mainHeight - 2).  // -2 for borders
		Render(m.contentView())

	footer := m.footerView()

	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, content, sidebar)
	return lipgloss.JoinVertical(lipgloss.Left, mainArea, footer)
}
