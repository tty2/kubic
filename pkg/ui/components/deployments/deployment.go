package deployments

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
	nameHeader         = "Name"
	readyHeader        = "Ready"
	upToDateHeader     = "UpToDate"
	availableHeader    = "Available"
	ageHeader          = "Age"
	minColumnGap       = "  "
	nameColumnLen      = 20
	readyColumnLen     = 7
	upToDateColumnLen  = len(upToDateHeader)
	availableColumnLen = len(availableHeader)
	tableHeaderHeight  = 3
)

type (
	deployment struct {
		Name      string
		Ready     string
		UpToDate  int
		Available int
		Age       string
		Labels    map[string]string
		Styles    *themes.Styles
	}
)

// FilterValue is used to set filter item and required for `list.Model` interface.
func (d *deployment) FilterValue() string { return d.Name }
func (d *deployment) Height() int         { return 1 }
func (d *deployment) Spacing() int        { return 1 }
func (d *deployment) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d *deployment) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	s, ok := listItem.(*deployment)
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

	upToDate := fmt.Sprintf("%d", s.UpToDate)
	row.WriteString(upToDate)
	row.WriteString(strings.Repeat(" ", upToDateColumnLen-lipgloss.Width(upToDate)))
	row.WriteString(minColumnGap)

	available := fmt.Sprintf("%d", s.Available)
	row.WriteString(available)
	row.WriteString(strings.Repeat(" ", availableColumnLen-lipgloss.Width(available)))
	row.WriteString(minColumnGap)

	row.WriteString(s.Age)

	deploymentInfo := row.String()

	if index == m.Index() {
		fmt.Fprint(w, d.Styles.SelectedText.Render(deploymentInfo))
	} else {
		fmt.Fprint(w, d.Styles.MainText.Render(deploymentInfo))
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

	header.WriteString(upToDateHeader)
	header.WriteString(strings.Repeat(" ", upToDateColumnLen-len(upToDateHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(availableHeader)
	header.WriteString(strings.Repeat(" ", availableColumnLen-len(availableHeader)))
	header.WriteString(minColumnGap)

	header.WriteString(ageHeader)

	return header.String()
}
