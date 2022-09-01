package help

import (
	"strings"

	bbHelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tty2/kubic/pkg/ui/shared"
	"github.com/tty2/kubic/pkg/ui/shared/elements/divider"
)

const fullHelpHeigh = 2

type Model struct {
	app  *shared.App
	help bbHelp.Model
}

func New(app *shared.App) *Model {
	help := bbHelp.NewModel()
	help.Styles = bbHelp.Styles{
		ShortDesc:      app.Styles.InactiveText.Copy(),
		FullDesc:       app.Styles.InactiveText.Copy(),
		ShortSeparator: app.Styles.InactiveText.Copy(),
		FullSeparator:  app.Styles.InactiveText.Copy(),
		FullKey:        app.Styles.InactiveText.Copy(),
		ShortKey:       app.Styles.InactiveText.Copy(),
		Ellipsis:       app.Styles.InactiveText.Copy(),
	}

	return &Model{
		app:  app,
		help: help,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.app.KeyMap.Help) {
			m.help.ShowAll = !m.help.ShowAll
			if m.help.ShowAll {
				m.app.GUI.Areas.HelpBar.Height += fullHelpHeigh
			} else {
				m.app.GUI.Areas.HelpBar.Height -= fullHelpHeigh
			}
			m.app.ResizeAreas()
		}
	}

	return m, nil
}

func (m *Model) View() string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(divider.HorizontalLine(m.app.GUI.ScreenWidth, m.app.Styles.InactiveText))
	s.WriteString("\n")
	if m.help.ShowAll {
		if m.app.CurrentTab == shared.NamespacesTab {
			s.WriteString(m.help.FullHelpView(m.app.KeyMap.FullHelp()))
		} else {
			s.WriteString(m.help.FullHelpView(m.app.KeyMap.FullWithFocus()))
		}
	} else {
		if m.app.CurrentTab == shared.NamespacesTab {
			s.WriteString(m.help.ShortHelpView(m.app.KeyMap.ShortHelp()))
		} else {
			s.WriteString(m.help.ShortHelpView(m.app.KeyMap.ShortWithFocus()))
		}
	}
	m.help.Width = m.app.GUI.ScreenWidth

	return m.app.Styles.HelpBar.
		Height(m.app.GUI.Areas.HelpBar.Height).Render(s.String())
}
