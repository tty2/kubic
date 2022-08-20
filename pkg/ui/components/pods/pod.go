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
	podNameColumn     = "Name"
	podReadyColumn    = "Ready"
	podStatusColumn   = "Status "
	podRestartsColumn = "Restarts"
	podAgeColumn      = "Age"
	minColumnGap      = "  "
	nameColumnLen     = 20
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
	row.WriteString(" ")

	ready := fmt.Sprintf("%s%s", s.Ready, minColumnGap)
	row.WriteString(ready)
	row.WriteString(strings.Repeat(" ", lipgloss.Width(podReadyColumn)+2-lipgloss.Width(ready)))

	status := fmt.Sprintf("%s%s", s.Status, minColumnGap)
	row.WriteString(status)
	row.WriteString(strings.Repeat(" ", lipgloss.Width(podStatusColumn)+2-lipgloss.Width(status)))

	restarts := fmt.Sprintf("%d%s", s.Restarts, minColumnGap)
	row.WriteString(restarts)
	row.WriteString(strings.Repeat(" ", lipgloss.Width(podRestartsColumn)+2-lipgloss.Width(restarts)))

	row.WriteString(s.Age)

	podInfo := row.String()

	if index == m.Index() {
		fmt.Fprint(w, p.Styles.SelectedText.Render(podInfo))
	} else {
		fmt.Fprint(w, p.Styles.MainText.Render(podInfo))
	}
}

func getHeader() string {
	var header strings.Builder
	header.WriteString(minColumnGap)
	header.WriteString(podNameColumn)
	header.WriteString(strings.Repeat(" ", nameColumnLen-len(podNameColumn)+len(minColumnGap)-1))
	header.WriteString(podReadyColumn)
	header.WriteString(minColumnGap)
	header.WriteString(podStatusColumn)
	header.WriteString(minColumnGap)
	header.WriteString(podRestartsColumn)
	header.WriteString(minColumnGap)
	header.WriteString(podAgeColumn)

	return header.String()
}
