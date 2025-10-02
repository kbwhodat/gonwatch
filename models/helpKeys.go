package models

import (
	"github.com/charmbracelet/bubbles/key"
)

var GoBackKeys = key.NewBinding(
    key.WithKeys("←", "h"), // keys that trigger this action
    key.WithHelp("←/h", "back"), // text to display in help menu
)

var GoForwardKeys = key.NewBinding(
    key.WithKeys("→", "l", "enter"), // keys that trigger this action
    key.WithHelp("→/l/enter", "forward"), // text to display in help menu
)
