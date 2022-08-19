package namespaces

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/domain"
	"github.com/tty2/kubic/pkg/ui/shared"
	"github.com/tty2/kubic/pkg/ui/shared/elements/divider"
)

const indent = 2

type namespacesRepo interface {
	GetNamespaces(ctx context.Context) ([]domain.Namespace, error)
}

type Model struct {
	app  *shared.App
	list list.Model
	repo namespacesRepo
}

func New(app *shared.App, repo namespacesRepo) (*Model, error) {
	m := Model{
		repo: repo,
		app:  app,
	}

	itemsModel := list.New([]list.Item{}, &namespace{
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
	m.setActive()

	return &m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.app.KeyMap.Select) {
			m.setActive()
		}
	}

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
		strings.Repeat(" ", shared.Max(len(minColumnGap),
			m.app.GUI.ScreenWidth-lipgloss.Width(header)-lipgloss.Width(m.app.CurrentNamespace)-tableHeaderHorizontalMargin)),
		m.app.CurrentNamespace)
	s.WriteString(m.app.Styles.InactiveText.Render(header))
	s.WriteString("\n")
	s.WriteString(divider.HorizontalLine(m.app.GUI.ScreenWidth, m.app.Styles.InactiveText))
	s.WriteString("\n")
	m.list.SetHeight(m.app.GUI.Areas.MainContent.Height - tableHeaderHeight)
	s.WriteString(m.app.Styles.InitStyle.Copy().MarginRight(indent).PaddingRight(indent).Render(m.list.View()))

	return s.String()
}

func (m *Model) setActive() {
	selected := m.list.Index()
	items := m.list.Items()
	for i := range items {
		s, ok := items[i].(*namespace)
		if !ok {
			return
		}
		if i == selected {
			s.Active = true
			m.app.CurrentNamespace = s.Name

			continue
		}
		s.Active = false
	}
}

func (m *Model) UpdateList() {
	ns, err := m.repo.GetNamespaces(context.Background())
	if err != nil {
		log.Fatalf("can't get namespaces: %v", err)
	}

	items := make([]list.Item, len(ns))
	for i := range ns {
		n := namespace{
			Name:   ns[i].Name,
			Status: ns[i].Status,
			Age:    ns[i].Age,
		}
		if i == 0 {
			n.Active = true
		}

		items[i] = &n
	}

	m.list.SetItems(items)
}
