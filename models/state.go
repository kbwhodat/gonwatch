package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
)

type ModelStateSnapshot struct {
	TextInputState textinput.Model
	ListState      list.Model
	Mode           string
}

func (m *Model) restorePreviousState() (*Model) {
	if len(m.PreviousStates) == 0 {
		return m
	}

	lastIndex := len(m.PreviousStates) - 1
	prevState := m.PreviousStates[lastIndex]
	m.PreviousStates = m.PreviousStates[:lastIndex]

	// Restore the model state
	m.TextInput = prevState.TextInputState
	m.List = prevState.ListState
	m.Mode = prevState.Mode

	return m
}

func (m *Model) saveCurrentState() {
	snapshot := ModelStateSnapshot{
		TextInputState: m.TextInput,
		ListState:      m.List,
		Mode:           m.Mode,
	}
	m.PreviousStates = append(m.PreviousStates, snapshot)
}
