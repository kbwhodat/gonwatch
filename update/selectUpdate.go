package update

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	"github.com/charmbracelet/bubbles/list"
)

func SelectUpdate(m list.Model, msg tea.Msg) (list.Model, tea.Cmd) {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	m.SetSize(w-2, h-2)

	m, cmd := m.Update(msg)

	return m, cmd
}

func SelectUpdateMsg(m list.Model) list.Item {
	selectItem := m.SelectedItem()

	return selectItem

}
