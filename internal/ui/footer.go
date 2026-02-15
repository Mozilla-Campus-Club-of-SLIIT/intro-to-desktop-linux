package ui

import "fmt"

func (m *RootModel) footerView(width int) string {
	return FooterStyle.Width(width).Render(fmt.Sprintf(" [q] Quit | [?] Help | %dx%d", m.width, m.height))
}
