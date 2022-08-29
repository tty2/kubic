package deployments

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
	"github.com/tty2/kubic/pkg/ui/shared/infobar"
)

type focused int

const (
	listInFocus focused = iota
	infoInFocus
)

type deploymentsRepo interface {
	GetDeployments(ctx context.Context, namespace string) ([]domain.Deployment, error)
}

// Model for deployments.
// Mutex is necessary here.
// We must synchronize UpdateList function call and View function call on update namespaces.
// In order to make user interface faster on update namespace we call update callbacks in another goroutine.
// namespaces/model.go package Model.setActive() function has go m.app.OnUpdateNamespace()
// If user switch tab faster than k8s makes call to update list, user will get outdated list.
// Mutex helps us to wait for k8s response and update list before view.
type Model struct {
	app     *shared.App
	list    list.Model
	repo    deploymentsRepo
	mu      sync.Mutex
	focused focused
	infobar *infobar.Model
}

func New(app *shared.App, repo deploymentsRepo) (*Model, error) {
	m := Model{
		repo:    repo,
		app:     app,
		infobar: infobar.New(),
	}

	itemsModel := list.New([]list.Item{}, &deployment{
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

	m.infobar.SetWH(app.GUI.ScreenWidth-lipgloss.Width(getHeader()), app.GUI.Areas.MainContent.Height-tableHeaderHeight)

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

			return m, cmd
		}
	}

	if m.listInFocus() {
		m.list, cmd = m.list.Update(msg)
		m.infobar.SetContent(m.getCurrentDeployment().renderInfo())
	} else {
		_, cmd = m.infobar.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
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
	m.setContentHeight()

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

	deps, err := m.repo.GetDeployments(context.Background(), m.app.CurrentNamespace)
	if err != nil {
		return
	}

	items := make([]list.Item, len(deps))
	for i := range deps {
		items[i] = &deployment{
			Name:              deps[i].Name,
			Ready:             deps[i].Ready,
			UpdatedReplicas:   deps[i].UpdatedReplicas,
			AvailableReplicas: deps[i].AvailableReplicas,
			ReadyReplicas:     deps[i].ReadyReplicas,
			Tolerations:       deps[i].Tolerations,
			Labels:            deps[i].Labels,
			Age:               deps[i].Age,
			Meta:              deps[i].Meta,
		}
	}

	m.list.SetItems(items)
}

func (m *Model) changeFocusRight() {
	if m.listInFocus() {
		m.focused = infoInFocus
	}
}

func (m *Model) changeFocusLeft() {
	if m.infoInFocus() {
		m.focused = listInFocus
	}
}

func (m *Model) listInFocus() bool {
	return m.focused == listInFocus
}

func (m *Model) infoInFocus() bool {
	return m.focused == infoInFocus
}

func (m *Model) resetFocus() {
	m.focused = listInFocus
	m.list.ResetSelected()
}

func (m *Model) renderInfoBar() string {
	infoData := m.infobar.View()

	if !m.infoInFocus() {
		infoData = m.app.Styles.InactiveText.Render(infoData)
	}

	info := lipgloss.JoinVertical(lipgloss.Left,
		m.renderInfoTitleBar(),
		infoData,
	)

	return m.app.Styles.InitStyle.Copy().MarginLeft(m.app.Styles.TextLeftMargin).Render(info)
}

func (m *Model) getCurrentDeployment() *deployment {
	item := m.list.SelectedItem()
	dep, ok := item.(*deployment)
	if !ok {
		return nil
	}

	return dep
}

func (m *Model) renderInfoTitleBar() string {
	tabs := getInfoTabs()
	titles := make([]string, len(tabs))
	for i := range tabs {
		if m.infoInFocus() {
			titles[i] = m.app.Styles.ActiveInfoTab.Render(tabs[i])
		} else {
			titles[i] = m.app.Styles.InactiveInfoTab.Render(tabs[i])
		}
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

func getInfoTabs() []string {
	return []string{"Info"}
}

func (m *Model) setInfoContent() {
	dep := m.getCurrentDeployment()
	if dep == nil {
		return
	}
	dep.Styles = m.app.Styles
	m.infobar.SetContent(
		m.getCurrentDeployment().renderInfo(),
	)
}

func (m *Model) setContentHeight() {
	m.infobar.SetWH(
		m.app.GUI.ScreenWidth-lipgloss.Width(getHeader()),
		m.app.GUI.Areas.MainContent.Height-tableHeaderHeight,
	)
	m.list.SetHeight(m.app.GUI.Areas.MainContent.Height - tableHeaderHeight)
}
