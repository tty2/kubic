package pods

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
	nameHeader        = "Name"
	readyHeader       = "Ready"
	statusHeader      = "Status"
	restartsHeader    = "Restarts"
	ageHeader         = "Age"
	minColumnGap      = "  "
	nameColumnLen     = 20
	readyColumnLen    = 7
	statusColumnLen   = 9 // the longest status `Succeeded`
	restartsColumnLen = len(restartsHeader)
	tableHeaderHeight = 3
)

type (
	pod struct {
		Name     string
		Ready    string
		Status   string
		Restarts int
		Age      string
		Labels   map[string]string
		Styles   *themes.Styles
	}
)

// FilterValue is used to set filter item and required for `list.Model` interface.
func (p *pod) FilterValue() string { return p.Name }
func (p *pod) Height() int         { return 1 }
func (p *pod) Spacing() int        { return 1 }
func (p *pod) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (p *pod) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	s, ok := listItem.(*pod)
	if !ok {
		return
	}

	name := shared.GetTextWithLen(s.Name, nameColumnLen)

	var row strings.Builder
	row.WriteString(name)
	row.WriteString(minColumnGap)

	row.WriteString(s.Ready)
	row.WriteString(strings.Repeat(" ", readyColumnLen-lipgloss.Width(s.Ready)))
	row.WriteString(minColumnGap)

	row.WriteString(s.Status)
	row.WriteString(strings.Repeat(" ", statusColumnLen-lipgloss.Width(s.Status)))
	row.WriteString(minColumnGap)

	restarts := fmt.Sprintf("%d", s.Restarts)
	row.WriteString(restarts)
	row.WriteString(strings.Repeat(" ", restartsColumnLen-lipgloss.Width(restarts)))
	row.WriteString(minColumnGap)

	row.WriteString(s.Age)

	rowString := row.String()

	if index == m.Index() {
		fmt.Fprint(w, p.Styles.SelectedText.Render(rowString))
	} else {
		fmt.Fprint(w, p.Styles.MainText.Render(rowString))
	}
}

func getHeader() string {
	var header strings.Builder
	header.WriteString(minColumnGap)

	header.WriteString(nameHeader)
	header.WriteString(strings.Repeat(" ", nameColumnLen-len(nameHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(readyHeader)
	header.WriteString(strings.Repeat(" ", readyColumnLen-len(readyHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(statusHeader)
	header.WriteString(strings.Repeat(" ", statusColumnLen-len(statusHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(restartsHeader)
	header.WriteString(strings.Repeat(" ", restartsColumnLen-len(restartsHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(ageHeader)

	return header.String()
}
