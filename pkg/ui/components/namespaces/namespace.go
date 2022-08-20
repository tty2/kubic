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
	active            = "✔"
	inactive          = " "
	nameHeader        = "Name"
	statusHeader      = "Status"
	ageHeader         = "Age"
	minColumnGap      = "  "
	nameColumnLen     = 20
	statusColumnLen   = 11 // the longest status `Terminating`
	minGapLen         = len(minColumnGap)
	tableHeaderHeight = 3
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
	row.WriteString(minColumnGap)

	row.WriteString(s.Status)
	row.WriteString(strings.Repeat(" ", statusColumnLen-lipgloss.Width(s.Status)))
	row.WriteString(minColumnGap)

	row.WriteString(s.Age)

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

	header.WriteString(nameHeader)
	header.WriteString(strings.Repeat(" ", nameColumnLen-len(nameHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(statusHeader)
	header.WriteString(strings.Repeat(" ", statusColumnLen-len(statusHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(ageHeader)

	return header.String()
}
