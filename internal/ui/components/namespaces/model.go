package namespaces

import (
	"context"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tty2/kubic/internal/domain"
	"github.com/tty2/kubic/internal/ui/shared"
	"github.com/tty2/kubic/internal/ui/shared/elements/divider"
	"github.com/tty2/kubic/internal/ui/shared/themes"
)

type namespacesRepo interface {
	GetNamespaces(ctx context.Context) ([]domain.Namespace, error)
}

type Model struct {
	repo  namespacesRepo
	state *shared.State
	list  list.Model
	keys  shared.KeyMap
}

func New(st *shared.State, repo namespacesRepo) (*Model, error) {
	m := Model{
		repo:  repo,
		state: st,
		keys:  shared.GetKeyMaps(),
	}

	itemsModel := list.New([]list.Item{}, &namespace{}, 0, 0)
	itemsModel.SetFilteringEnabled(false)
	itemsModel.SetShowFilter(false)
	itemsModel.SetShowTitle(false)
	itemsModel.SetShowStatusBar(false)
	itemsModel.SetShowHelp(false)
	m.list = itemsModel
	m.UpdateList()

	return &m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.keys.Select) {
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
	s.WriteString(themes.MainDocStyle.Foreground(themes.DefaultTheme.InactiveText).Render(getHeader()))
	s.WriteString("\n")
	s.WriteString(divider.HorizontalLine(m.state.ScreenWidth, themes.DefaultTheme.InactiveText))
	s.WriteString("\n")
	m.list.SetHeight(m.state.Areas.MainContent.Height - tableHeaderHeight)
	s.WriteString(listStyle.Render(m.list.View()))

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
		s.Active = i == selected
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
		m.state.Namespace = n.Name
	}

	m.list.SetItems(items)
}
