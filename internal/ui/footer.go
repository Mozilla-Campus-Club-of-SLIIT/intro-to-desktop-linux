package ui

import "fmt"

func (m *RootModel) footerView() string {
	return FooterStyle.Width(m.width).Render(fmt.Sprintf(" [q] Quit | [?] Help | %dx%d", m.width, m.height))
}
