package pods

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/domain"
	"github.com/tty2/kubic/pkg/ui/shared"
	"github.com/tty2/kubic/pkg/ui/shared/elements/divider"
	"github.com/tty2/kubic/pkg/ui/shared/elements/infobar"
)

type focused int

const (
	listInFocus focused = iota
	infoInFocus
	logInFocus
)

type podsRepo interface {
	GetPods(ctx context.Context, namespace string) ([]domain.Pod, error)
}

// Model for pods.
// Mutex is necessary here.
// We must synchronize UpdateList function call and View function call on update namespaces.
// In order to make user interface faster on update namespace we call update callbacks in another goroutine.
// namespaces/model.go package Model.setActive() function has go m.app.OnUpdateNamespace()
// If user switch tab faster than k8s makes call to update list, user will get outdated list.
// Mutex helps us to wait for k8s response and update list before view.
type Model struct {
	app     *shared.App
	list    list.Model
	repo    podsRepo
	mu      sync.Mutex
	focused focused
	infobar *infobar.Model
}

func New(app *shared.App, repo podsRepo) (*Model, error) {
	m := Model{
		repo:    repo,
		app:     app,
		infobar: infobar.New(),
	}

	itemsModel := list.New([]list.Item{}, &pod{
		Styles: app.Styles,
	}, 0, 0)
	itemsModel.SetFilteringEnabled(false)
	itemsModel.SetShowFilter(false)
	itemsModel.SetShowTitle(false)
	itemsModel.SetShowStatusBar(false)
	itemsModel.SetShowHelp(false)
	itemsModel.Paginator.Type = paginator.Dots
	m.list = itemsModel
	m.UpdateList()
	m.app.AddUpdateNamespaceCallback(m.UpdateList)
	m.app.AddUpdateNamespaceCallback(m.resetFocus)
	m.app.AddUpdateNamespaceCallback(m.setInfoContent)

	m.setInfoBarHeight()

	return &m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.app.KeyMap.FocusRight):
			m.changeFocusRight()

			return m, cmd
		case key.Matches(msg, m.app.KeyMap.FocusLeft):
			m.changeFocusLeft()
			m.infobar.ResetView()

			return m, cmd
		}
	}

	if m.listInFocus() {
		m.list, cmd = m.list.Update(msg)
		m.setInfoContent()
	} else {
		_, cmd = m.infobar.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	m.setInfoBarHeight()

	var s strings.Builder
	s.WriteString("\n")
	header := getHeader()
	header = fmt.Sprintf("%s%s%s",
		header,
		strings.Repeat(" ", shared.Max(
			len(minColumnGap),
			m.app.GUI.ScreenWidth-lipgloss.Width(header)-lipgloss.Width(m.app.CurrentNamespace)-m.app.Styles.TextRightMargin)),
		m.app.CurrentNamespace)
	s.WriteString(m.app.Styles.InactiveText.Render(header))
	s.WriteString("\n")
	s.WriteString(divider.HorizontalLine(m.app.GUI.ScreenWidth, m.app.Styles.InactiveText))
	s.WriteString("\n")
	m.list.SetHeight(m.app.GUI.Areas.MainContent.Height - tableHeaderHeight)

	m.mu.Lock()
	defer m.mu.Unlock()

	s.WriteString(
		m.app.Styles.InitStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.app.Styles.ListRightBorder.Render(m.list.View()),
				m.renderInfoBar(),
			),
		))

	return s.String()
}

func (m *Model) UpdateList() {
	m.mu.Lock()
	defer m.mu.Unlock()

	pods, err := m.repo.GetPods(context.Background(), m.app.CurrentNamespace)
	if err != nil {
		return
	}

	items := make([]list.Item, len(pods))
	for i := range pods {
		items[i] = &pod{
			Name:       pods[i].Name,
			Ready:      pods[i].Ready,
			Status:     pods[i].Status,
			Restarts:   pods[i].Restarts,
			Age:        pods[i].Age,
			Meta:       pods[i].Meta,
			Spec:       pods[i].Spec,
			StatusInfo: pods[i].StatusInfo,
		}
	}

	m.list.SetItems(items)
}

func (m *Model) changeFocusRight() {
	switch m.focused {
	case listInFocus:
		m.focused = infoInFocus
	case infoInFocus:
		m.focused = logInFocus
	}
}

func (m *Model) changeFocusLeft() {
	switch m.focused {
	case logInFocus:
		m.focused = infoInFocus
	case infoInFocus:
		m.focused = listInFocus
	}
}

func (m *Model) listInFocus() bool {
	return m.focused == listInFocus
}

func (m *Model) infoInFocus() bool {
	return m.focused == infoInFocus
}

func (m *Model) logInFocus() bool {
	return m.focused == logInFocus
}

func (m *Model) resetFocus() {
	m.focused = listInFocus
	m.list.ResetSelected()
}

func (m *Model) renderInfoBar() string {
	var infoBarData string
	switch m.focused {
	case listInFocus:
		infoData := m.infobar.View()
		infoBarData = m.app.Styles.InactiveText.Render(infoData)
	case infoInFocus, logInFocus:
		infoBarData = m.infobar.View()
	}

	info := lipgloss.JoinVertical(lipgloss.Left,
		m.renderInfoBarTabs(),
		infoBarData,
	)

	return m.app.Styles.InitStyle.Copy().MarginLeft(m.app.Styles.TextLeftMargin).Render(info)
}

func (m *Model) getCurrentPod() *pod {
	item := m.list.SelectedItem()
	p, ok := item.(*pod)
	if !ok {
		return nil
	}

	return p
}

func (m *Model) renderInfoBarTabs() string {
	tabs := getInfoTabs()
	titles := make([]string, len(tabs))
	for i := range tabs {
		if m.focused == tabs[i] {
			titles[i] = m.app.Styles.ActiveInfoTab.Render(tabs[i].String())
			continue
		}
		titles[i] = m.app.Styles.InactiveInfoTab.Render(tabs[i].String())
	}

	titlesStr := lipgloss.JoinHorizontal(
		lipgloss.Top,
		titles...,
	)

	gap := m.app.Styles.InfoGap.Render(
		strings.Repeat(" ", shared.Max(0, m.app.GUI.ScreenWidth-lipgloss.Width(titlesStr))),
	)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, titlesStr, gap)
}

func (m *Model) setInfoContent() {
	dep := m.getCurrentPod()
	if dep == nil {
		m.infobar.SetContent("")

		return
	}
	dep.Styles = m.app.Styles
	m.infobar.SetContent(
		m.getCurrentPod().renderInfo(),
	)
}

func (m *Model) setInfoBarHeight() {
	m.infobar.SetWH(
		m.app.GUI.ScreenWidth-lipgloss.Width(getHeader()),
		m.app.GUI.Areas.MainContent.Height-tableHeaderHeight,
	)
	m.list.SetHeight(m.app.GUI.Areas.MainContent.Height - tableHeaderHeight)
}

func getInfoTabs() []focused {
	return []focused{
		infoInFocus,
		logInFocus,
	}
}

func (f focused) String() string {
	switch f {
	case infoInFocus:
		return "Info"
	case logInFocus:
		return "Logs"
	default:
		return ""
	}
}
