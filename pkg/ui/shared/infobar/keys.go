package infobar

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up   key.Binding
	Down key.Binding
}

// nolint gochecknoglobals: used here on purpose
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
}
