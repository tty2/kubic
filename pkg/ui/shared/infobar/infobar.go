package infobar

import (
	"fmt"

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

func (m *Model) View() string {
	// log.Println(m.width, m.height)
	style := lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Width(m.width).
		MaxWidth(m.width)

	return style.Copy().Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		lipgloss.NewStyle().Height(percentHeight).
			PaddingTop(percentPaddingTop).
			Render(fmt.Sprintf("%d%%", int(m.viewport.ScrollPercent()*100))),
	))
}

func (m *Model) SetContent(data string) {
	m.data = data
	m.viewport.SetContent(data)
}

func (m *Model) SetWH(w, h int) {
	// log.Println("set", w, h)
	m.width = w
	m.height = h - percentHeight - percentPaddingTop
	m.viewport.Width = w
	m.viewport.Height = m.height
}
