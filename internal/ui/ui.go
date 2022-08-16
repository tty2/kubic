package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/internal/commands/k8s"
	"github.com/tty2/kubic/internal/config"
	"github.com/tty2/kubic/internal/ui/components/help"
	"github.com/tty2/kubic/internal/ui/components/namespaces"
	"github.com/tty2/kubic/internal/ui/components/tabs"
	"github.com/tty2/kubic/internal/ui/shared"
	"golang.org/x/term"
)

type components struct {
	tabs       tea.Model
	namespaces tea.Model
	// deployments tea.Model
	// pods        tea.Model
	help tea.Model
}

type App struct {
	components components
	config     *config.Config
	state      *shared.State
	keys       shared.KeyMap
}

func New(cfg *config.Config, k8sClient *k8s.Client) (tea.Model, error) {
	st := shared.NewState(nil)
	app := App{
		config: cfg,
		state:  st,
		keys:   shared.GetKeyMaps(),
		components: components{
			tabs: tabs.New(st, shared.GetTabItems()),
			help: help.New(st),
		},
	}

	ns, err := namespaces.New(st, k8sClient)
	if err != nil {
		return nil, err
	}
	app.components.namespaces = ns

	app.state.ScreenWidth, app.state.ScreenHeight, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}
	app.state.ResizeAreas()

	return &app, nil
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd = a.keyEventHandle(msg)
	case tea.WindowSizeMsg:
		a.onWindowSizeChanged(msg)
	}
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	s := strings.Builder{}

	// tabs
	s.WriteString(a.components.tabs.View())
	s.WriteString("\n")

	// content
	if a.state.CurrentTab == shared.NamespacesTab {
		mainContent := lipgloss.JoinHorizontal(lipgloss.Top, a.components.namespaces.View())
		s.WriteString(mainContent)
		s.WriteString("\n")
	}

	// help
	s.WriteString(lipgloss.PlaceVertical(
		a.state.Areas.HelpBar.Height,
		lipgloss.Bottom, a.components.help.View()))

	return s.String()
}

func (a *App) keyEventHandle(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, a.keys.Quit):
		return tea.Quit
	case key.Matches(msg, a.keys.Help):
		a.components.help.Update(msg)

		return nil
	case key.Matches(msg, a.keys.Tab, a.keys.ShiftTab):
		a.components.tabs.Update(msg)

		return nil
	default:
		return a.componentsKeyEventHandle(msg)
	}
}

func (a *App) componentsKeyEventHandle(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch a.state.CurrentTab {
	case shared.NamespacesTab:
		_, cmd = a.components.namespaces.Update(msg)
		// case shared.SettingsTab:
		// 	_, cmd = a.components.settings.Update(msg)
	}

	return cmd
}

func (a *App) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	a.state.ScreenWidth = msg.Width
	a.state.ScreenHeight = msg.Height
	a.state.ResizeAreas()
}
