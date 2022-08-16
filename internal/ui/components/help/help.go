package help

import (
	"strings"

	bbHelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tty2/kubic/internal/ui/shared"
	"github.com/tty2/kubic/internal/ui/shared/elements/divider"
	"github.com/tty2/kubic/internal/ui/shared/themes"
)

const fullHelpHeigh = 2

type Model struct {
	state *shared.State
	help  bbHelp.Model
}

func New(state *shared.State) *Model {
	help := bbHelp.NewModel()
	help.Styles = bbHelp.Styles{
		ShortDesc:      helpTextStyle.Copy(),
		FullDesc:       helpTextStyle.Copy(),
		ShortSeparator: helpTextStyle.Copy(),
		FullSeparator:  helpTextStyle.Copy(),
		FullKey:        helpTextStyle.Copy(),
		ShortKey:       helpTextStyle.Copy(),
		Ellipsis:       helpTextStyle.Copy(),
	}

	return &Model{
		state: state,
		help:  help,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.state.KeyMap.Help) {
			m.help.ShowAll = !m.help.ShowAll
			if m.help.ShowAll {
				m.state.Areas.HelpBar.Height += fullHelpHeigh
			} else {
				m.state.Areas.HelpBar.Height -= fullHelpHeigh
			}
			m.state.ResizeAreas()
		}
	}

	return m, nil
}

func (m *Model) View() string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(divider.HorizontalLine(m.state.ScreenWidth, themes.DefaultTheme.InactiveText))
	s.WriteString("\n")
	if m.help.ShowAll {
		s.WriteString(m.help.FullHelpView(m.state.KeyMap.FullHelp()))
	} else {
		s.WriteString(m.help.ShortHelpView(m.state.KeyMap.ShortHelp()))
	}
	m.help.Width = m.state.ScreenWidth

	return helpStyle.Copy().Height(m.state.Areas.HelpBar.Height).Render(s.String())
}
