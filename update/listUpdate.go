package update

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

func ListUpdate(m list.Model, msg tea.Msg) (list.Model, tea.Cmd) {

	m, cmd := m.Update(msg)

	return m, cmd
}
