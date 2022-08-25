package tabs

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/ui/shared"
)

const (
	tabsLeftRightIndents = 2
)

type Model struct {
	app  *shared.App
	tabs []shared.TabItem
}

func New(app *shared.App, ti []shared.TabItem) *Model {
	m := Model{
		app:  app,
		tabs: ti,
	}

	return &m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.app.KeyMap.Tab) {
			m.next()
		} else if key.Matches(msg, m.app.KeyMap.ShiftTab) {
			m.prev()
		}
	}

	return m, nil
}

func (m *Model) View() string {
	titles := m.getTabsTitles()

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		titles...,
	)

	gap := m.app.Styles.TabsGap.Render(
		strings.Repeat(" ", shared.Max(0, m.app.GUI.ScreenWidth-lipgloss.Width(row)-tabsLeftRightIndents)),
	)

	return m.app.Styles.InitStyle.Render(lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap))
}

func (m *Model) next() {
	i := m.getNextTabIndex()
	m.app.CurrentTab = m.tabs[i]
}

func (m *Model) prev() {
	i := m.getPrevTabIndex()
	m.app.CurrentTab = m.tabs[i]
}

func (m *Model) getNextTabIndex() int {
	return (m.getCurrentTabIndex() + 1) % len(m.tabs)
}

func (m *Model) getPrevTabIndex() int {
	return (m.getCurrentTabIndex() - 1 + len(m.tabs)) % len(m.tabs)
}

func (m *Model) getCurrentTabIndex() int {
	for i := range m.tabs {
		if m.app.CurrentTab == m.tabs[i] {
			return i
		}
	}

	return 0
}

func (m *Model) getTabsTitles() []string {
	titles := make([]string, len(m.tabs))
	for i := range m.tabs {
		title := m.tabs[i].String()
		if m.tabs[i] == m.app.CurrentTab {
			titles[i] = m.app.Styles.ActiveTab.Render(title)
		} else {
			titles[i] = m.app.Styles.InactiveTab.Render(title)
		}
	}

	return titles
}
