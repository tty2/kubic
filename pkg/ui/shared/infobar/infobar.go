package infobar

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	percentPaddingTop = 1
	percentHeight     = 2
)

type Model struct {
	data     string
	width    int
	height   int
	viewport viewport.Model
}

func New() *Model {
	return &Model{}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Down):
			m.viewport.HalfViewDown()

		case key.Matches(msg, keys.Up):
			m.viewport.HalfViewUp()
		}
	}

	return m, nil
}

func (m *Model) ResetView() {
	m.viewport.GotoTop()
}

func (m *Model) View() string {
	style := lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Width(m.width).
		MaxWidth(m.width)

	var continueRead string
	if m.viewport.ScrollPercent() < 1 {
		continueRead = "..."
	}
	return style.Copy().Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		lipgloss.NewStyle().Height(percentHeight).
			Render(continueRead),
	))
}

func (m *Model) SetContent(data string) {
	m.data = data
	m.viewport.SetContent(data)
}

func (m *Model) SetWH(w, h int) {
	m.width = w
	m.height = h - percentHeight
	m.viewport.Width = w
	m.viewport.Height = m.height - percentPaddingTop
}
