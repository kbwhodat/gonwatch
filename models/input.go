package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"gonwatch/history"
)

func InputModel(m *Model) {
	ti := textinput.New()
	ti.Focus()

	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = "What " + m.Choice.choice + " would you like to watch?\n\n\n"

	m.TextInput = ti
	m.Mode = "input"

	m.searchHistory = history.GetSearches(m.Choice.choice)
	m.searchHistoryIndex = -1
	m.searchHistoryDraft = ""
}
