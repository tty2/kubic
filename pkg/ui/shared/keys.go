package shared

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// nolint gochecknoglobals: used here on purpose
var boldText = lipgloss.NewStyle().Bold(true)

type KeyMap struct {
	Tab        key.Binding
	ShiftTab   key.Binding
	Up         key.Binding
	Down       key.Binding
	PrevPage   key.Binding
	NextPage   key.Binding
	FocusRight key.Binding
	FocusLeft  key.Binding
	Select     key.Binding
	Help       key.Binding
	HelpShort  key.Binding
	Quit       key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Tab, k.Select}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.HelpShort, k.Quit, k.Tab},
		{k.Up, k.Down, k.PrevPage, k.NextPage},
		{k.Select},
	}
}

func (k KeyMap) ShortWithFocus() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Tab, k.FocusRight}
}

func (k KeyMap) FullWithFocus() [][]key.Binding {
	return [][]key.Binding{
		{k.HelpShort, k.Quit, k.Tab},
		{k.Up, k.Down, k.PrevPage, k.NextPage},
		{k.FocusLeft, k.FocusRight},
	}
}

// GetKeyMaps returns all the shortcuts available.
func GetKeyMaps() KeyMap {
	return KeyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp(boldText.Render("tab"), "next tab"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp(boldText.Render("Shift+tab"), "previous tab"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp(boldText.Render("↑/k"), "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp(boldText.Render("↓/j"), "move down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp(boldText.Render("←/h"), "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp(boldText.Render("→/l"), "next page"),
		),
		FocusRight: key.NewBinding(
			key.WithKeys(tea.KeyCtrlL.String(), tea.KeyCtrlRight.String()),
			key.WithHelp(boldText.Render("Ctrl+→/Ctrl+l"), "focus right"),
		),
		FocusLeft: key.NewBinding(
			key.WithKeys(tea.KeyCtrlH.String(), tea.KeyCtrlLeft.String()),
			key.WithHelp(boldText.Render("Ctrl+←/Ctrl+h"), "focus left"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp(boldText.Render("?"), "full help"),
		),
		HelpShort: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp(boldText.Render("?"), "short help"),
		),
		Select: key.NewBinding(
			key.WithKeys(tea.KeyEnter.String()),
			key.WithHelp(boldText.Render("Enter"), "select item"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp(boldText.Render("q"), "quit"),
		),
	}
}
