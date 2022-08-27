package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/k8s"
	"github.com/tty2/kubic/pkg/ui/components/deployments"
	"github.com/tty2/kubic/pkg/ui/components/help"
	"github.com/tty2/kubic/pkg/ui/components/namespaces"
	"github.com/tty2/kubic/pkg/ui/components/pods"
	"github.com/tty2/kubic/pkg/ui/components/tabs"
	"github.com/tty2/kubic/pkg/ui/shared"
	"github.com/tty2/kubic/pkg/ui/shared/themes"
	"golang.org/x/term"
)

type components struct {
	tabs        tea.Model
	namespaces  tea.Model
	deployments tea.Model
	pods        tea.Model
	help        tea.Model
}

type MainModel struct {
	components components
	app        *shared.App
}

func New(k8sClient *k8s.Client, theme themes.Theme) (tea.Model, error) {
	app := shared.NewApp(theme)
	model := MainModel{
		app: app,
		components: components{
			tabs: tabs.New(app, shared.GetTabItems()),
			help: help.New(app),
		},
	}

	ns, err := namespaces.New(app, k8sClient)
	if err != nil {
		return nil, err
	}
	model.components.namespaces = ns

	dep, err := deployments.New(app, k8sClient)
	if err != nil {
		return nil, err
	}
	model.components.deployments = dep

	pod, err := pods.New(app, k8sClient)
	if err != nil {
		return nil, err
	}
	model.components.pods = pod

	model.app.GUI.ScreenWidth, model.app.GUI.ScreenHeight, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}
	model.app.ResizeAreas()

	return &model, nil
}

func (model *MainModel) Init() tea.Cmd {
	return nil
}

func (model *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd = model.keyEventHandle(msg)
	case tea.WindowSizeMsg:
		model.onWindowSizeChanged(msg)
	}
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model *MainModel) View() string {
	s := strings.Builder{}

	// tabs
	s.WriteString(model.components.tabs.View())
	s.WriteString("\n")

	// content
	switch model.app.CurrentTab {
	case shared.NamespacesTab:
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, model.components.namespaces.View()))
	case shared.DeploymentsTab:
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, model.components.deployments.View()))
	case shared.PodsTab:
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, model.components.pods.View()))
	}

	// help
	s.WriteString(lipgloss.PlaceVertical(
		model.app.GUI.Areas.HelpBar.Height,
		lipgloss.Bottom, model.components.help.View()))

	return s.String()
}

func (model *MainModel) keyEventHandle(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, model.app.KeyMap.Quit):
		return tea.Quit
	case key.Matches(msg, model.app.KeyMap.Help):
		model.components.help.Update(msg)

		return nil
	case key.Matches(msg, model.app.KeyMap.Tab, model.app.KeyMap.ShiftTab):
		model.components.tabs.Update(msg)

		return nil
	default:
		return model.componentsKeyEventHandle(msg)
	}
}

func (model *MainModel) componentsKeyEventHandle(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch model.app.CurrentTab {
	case shared.NamespacesTab:
		_, cmd = model.components.namespaces.Update(msg)
	case shared.DeploymentsTab:
		_, cmd = model.components.deployments.Update(msg)
	case shared.PodsTab:
		_, cmd = model.components.pods.Update(msg)
	}

	return cmd
}

func (model *MainModel) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	model.app.GUI.ScreenWidth = msg.Width
	model.app.GUI.ScreenHeight = msg.Height
	model.app.ResizeAreas()
}
