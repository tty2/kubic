package namespaces

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/ui/shared"
	"github.com/tty2/kubic/pkg/ui/shared/themes"
)

const (
	active                = "✔"
	inactive              = " "
	namespaceNameColumn   = "Name"
	namespaceStatusColumn = "Status"
	namespaceAgeColumn    = "Age"
	minColumnGap          = "  "
	nameColumnLen         = 20
	tableHeaderHeight     = 3
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
func (n *namespace) FilterValue() string { return n.Name }
func (n *namespace) Height() int         { return 1 }
func (n *namespace) Spacing() int        { return 1 }
func (n *namespace) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (n *namespace) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	s, ok := listItem.(*namespace)
	if !ok {
		return
	}

	sign := inactive
	if s.Active {
		sign = n.Styles.NamespaceSign.Render(active)
	}

	name := shared.GetTextWithLen(s.Name, nameColumnLen)

	var row strings.Builder
	row.WriteString(name)
	row.WriteString(" ")

	status := fmt.Sprintf("%s%s", s.Status, minColumnGap)
	row.WriteString(status)
	// +3 is alignment: longest `Succeeded` status is 3 symbols longer than `Status` header
	// TODO: find a better solution
	row.WriteString(strings.Repeat(" ",
		lipgloss.Width(namespaceStatusColumn)+len(minColumnGap)+5-lipgloss.Width(status)))

	age := fmt.Sprintf("%s%s", s.Age, minColumnGap)
	row.WriteString(age)
	row.WriteString(strings.Repeat(" ", lipgloss.Width(namespaceAgeColumn)+len(minColumnGap)-lipgloss.Width(age)))

	rowInfo := row.String()

	if index == m.Index() {
		fmt.Fprintf(w, "%s %s", sign, n.Styles.SelectedText.Render(rowInfo))
	} else {
		fmt.Fprintf(w, "%s %s", sign, n.Styles.MainText.Render(rowInfo))
	}
}

func getHeader() string {
	var header strings.Builder
	header.WriteString(minColumnGap)
	header.WriteString(namespaceNameColumn)
	header.WriteString(strings.Repeat(" ", nameColumnLen-len(namespaceNameColumn)+len(minColumnGap)-1))
	// alignment: the longest `Terminating` status is 5 symbols longer than `Status` header
	// TODO: find a better solution
	header.WriteString(fmt.Sprintf("%s     ", namespaceStatusColumn))
	header.WriteString(minColumnGap)
	header.WriteString(namespaceAgeColumn)

	return header.String()
}
