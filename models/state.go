package models

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type ModelStateSnapshot struct {
	TextInputState textinput.Model
	ListState      list.Model
	Mode           string
	Cursor         int
}

func (m *Model) restorePreviousState() *Model {
	if len(m.PreviousStates) == 0 {
		return m
	}

	lastIndex := len(m.PreviousStates) - 1
	prevState := m.PreviousStates[lastIndex]
	m.PreviousStates = m.PreviousStates[:lastIndex]

	m.TextInput = prevState.TextInputState
	m.List = prevState.ListState
	m.Mode = prevState.Mode
	m.Cursor = prevState.Cursor

	if m.width > 0 && m.height > 0 {
		m.List.SetSize(m.width-2, m.height-2)
	}

	return m
}

func (m *Model) saveCurrentState() {
	snapshot := ModelStateSnapshot{
		TextInputState: m.TextInput,
		ListState:      m.List,
		Mode:           m.Mode,
		Cursor:         m.Cursor,
	}
	m.PreviousStates = append(m.PreviousStates, snapshot)
}
