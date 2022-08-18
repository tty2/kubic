package tabs

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tty2/kubic/pkg/ui/shared"
)

type Model struct {
	app  *shared.App
	tabs []*tab
}

func New(app *shared.App, ti []shared.TabItem) *Model {
	m := Model{
		app: app,
	}

	for i := range ti {
		m.tabs = append(m.tabs, newTab(ti[i]))
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
	tabs := make([]string, 0, len(m.tabs))
	for i := range m.tabs {
		tabs = append(tabs,
			m.tabs[i].render(
				m.tabs[i].getID() == m.app.CurrentTab,
			))
	}

	return renderTabBar(m.app.GUI.ScreenWidth, tabs)
}

func (m *Model) next() {
	i := m.getNextTabIndex()
	m.app.CurrentTab = m.tabs[i].getID()
}

func (m *Model) prev() {
	i := m.getPrevTabIndex()
	m.app.CurrentTab = m.tabs[i].getID()
}

func (m *Model) getNextTabIndex() int {
	return (m.getCurrentTabIndex() + 1) % len(m.tabs)
}

func (m *Model) getPrevTabIndex() int {
	return (m.getCurrentTabIndex() - 1 + len(m.tabs)) % len(m.tabs)
}

func (m *Model) getCurrentTabIndex() int {
	for i := range m.tabs {
		if m.app.CurrentTab == m.tabs[i].getID() {
			return i
		}
	}

	return 0
}
