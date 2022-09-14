package infobar

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/ui/shared/elements/viewport"
)

const (
	continueReadHeight  = 2
	continueReadPadding = 1
	continueRead        = "..."
)

type Model struct {
	width    int
	height   int
	viewport viewport.Model
}

func New() *Model {
	m := Model{}
	m.viewport.HorizontalStep = 5

	return &m
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, keys.Down):
			m.viewport.HalfViewDown()

		case key.Matches(msg, keys.Up):
			m.viewport.HalfViewUp()

		case key.Matches(msg, keys.Left):
			m.viewport.MoveLeft()

		case key.Matches(msg, keys.Right):
			m.viewport.MoveRight()
		}
	}

	return m, cmd
}

func (m *Model) ResetView() {
	m.viewport.GotoTop()
}

func (m *Model) View() string {
	style := lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height)

	if m.viewport.ScrollPercent() < 1 {
		m.viewport.Height = m.height - continueReadPadding

		return style.Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.viewport.View(),
			lipgloss.NewStyle().Height(continueReadHeight).Render(continueRead),
		))
	}

	m.viewport.Height = m.height

	return style.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
	))
}

func (m *Model) SetContent(data string) {
	m.viewport.SetContent(data)
}

func (m *Model) SetWH(w, h int) {
	m.width = w
	m.height = h - continueReadHeight
	m.viewport.Width = w
	m.viewport.Height = m.height - continueReadPadding
}

func (m *Model) ResetIndent() {
	m.viewport.ResetIndent()
}
