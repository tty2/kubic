package deployments

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/domain"
	"github.com/tty2/kubic/pkg/ui/shared"
	"github.com/tty2/kubic/pkg/ui/shared/elements/divider"
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
	app  *shared.App
	list list.Model
	repo deploymentsRepo
	mu   sync.Mutex
}

func New(app *shared.App, repo deploymentsRepo) (*Model, error) {
	m := Model{
		repo: repo,
		app:  app,
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

	return &m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

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
	m.list.SetHeight(m.app.GUI.Areas.MainContent.Height - tableHeaderHeight)

	m.mu.Lock()
	defer m.mu.Unlock()

	s.WriteString(m.app.Styles.ListRightBorder.Render(m.list.View()))

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
			Name:      deps[i].Name,
			Ready:     deps[i].Ready,
			UpToDate:  deps[i].UpToDate,
			Available: deps[i].Available,
			Labels:    deps[i].Labels,
			Age:       deps[i].Age,
		}
	}

	m.list.SetItems(items)
}

func (m *Model) renderInfo() string {
	item := m.list.SelectedItem()
	d, ok := item.(*deployment)
	if !ok {
		return ""
	}

	var info strings.Builder
	info.WriteString("Name")
	info.WriteString("\n")
	info.WriteString(d.Name)
	info.WriteString("\n")
	info.WriteString("Labels")

	for k, v := range d.Labels {
		info.WriteString(k)
		info.WriteString(": ")
		info.WriteString(v)
		info.WriteString("\n")
	}

	return info.String()
}
