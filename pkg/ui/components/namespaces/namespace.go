package namespaces

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/ui/shared/themes"
)

const (
	active                      = "✔"
	inactive                    = " "
	namespaceNameColumn         = "Name"
	namespaceStatusColumn       = "Status"
	namespaceAgeColumn          = "Age"
	minColumnGap                = "  "
	nameColumnLen               = 20
	tableHeaderHeight           = 3
	tableHeaderHorizontalMargin = 4
)

type (
	namespace struct {
		Name   string
		Status string
		Age    string
		Active bool
		Styles *themes.Styles
	}
)

// FilterValue is used to set filter item and required for `list.Model` interface.
func (v *namespace) FilterValue() string { return v.Name }
func (v *namespace) Height() int         { return 1 }
func (v *namespace) Spacing() int        { return 1 }
func (v *namespace) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (v *namespace) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	s, ok := listItem.(*namespace)
	if !ok {
		return
	}

	var name string
	lenName := len(s.Name)
	switch {
	case lenName > nameColumnLen:
		name = fmt.Sprintf("%s…", s.Name[:nameColumnLen-1])
	case lenName < nameColumnLen:
		name = fmt.Sprintf("%s%s", s.Name, strings.Repeat(" ", nameColumnLen-lipgloss.Width(s.Name)))
	default:
		name = s.Name
	}

	sign := inactive
	if s.Active {
		sign = v.Styles.NamespaceSign.Render(active)
	}

	var row strings.Builder
	namespaceInfo := fmt.Sprintf("%s %s\t%s", name, s.Status, s.Age)
	if index == m.Index() {
		row.WriteString(fmt.Sprintf("%s %s", sign, v.Styles.SelectedText.Render(namespaceInfo)))
	} else {
		row.WriteString(fmt.Sprintf("%s %s", sign, v.Styles.MainText.Render(namespaceInfo)))
	}

	fmt.Fprint(w, row.String())
}

func getHeader() string {
	var header strings.Builder
	header.WriteString(minColumnGap)
	header.WriteString(namespaceNameColumn)
	header.WriteString(strings.Repeat(" ", nameColumnLen-len(namespaceNameColumn)+len(minColumnGap)-1))
	header.WriteString(namespaceStatusColumn)
	header.WriteString("\t")
	header.WriteString(namespaceAgeColumn)

	return header.String()
}