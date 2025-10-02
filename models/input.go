package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	// "github.com/charmbracelet/bubbles/list"
)

func InputModel(m *Model) {
	ti := textinput.New()
	ti.Focus()

	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = "What " + m.Choice.choice + " would you like to watch?\n\n\n"

	// if m.Choice.choice == "movies" {
	// 	ti.Placeholder = "Avengers"
	// } else if m.Choice.choice == "live" {
	// 	ti.Placeholder = "boxing"
	// }

	m.TextInput = ti
	m.Mode = "input"
}
