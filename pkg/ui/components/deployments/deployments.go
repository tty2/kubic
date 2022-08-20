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
	deploymentNameColumn      = "Name"
	deploymentReadyColumn     = "Ready"
	deploymentUpToDateColumn  = "UpToDate"
	deploymentAvailableColumn = "Available"
	deploymentAgeColumn       = "Age"
	minColumnGap              = "  "
	nameColumnLen             = 20
	tableHeaderHeight         = 3
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
func (v *deployment) FilterValue() string { return v.Name }
func (v *deployment) Height() int         { return 1 }
func (v *deployment) Spacing() int        { return 1 }
func (v *deployment) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (v *deployment) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	s, ok := listItem.(*deployment)
	if !ok {
		return
	}

	name := shared.GetTextWithLen(s.Name, nameColumnLen)

	var row strings.Builder
	row.WriteString(name)
	row.WriteString(" ")

	ready := fmt.Sprintf("%s%s", s.Ready, minColumnGap)
	row.WriteString(ready)
	row.WriteString(strings.Repeat(" ", lipgloss.Width(deploymentReadyColumn)+2-lipgloss.Width(ready)))

	upToDate := fmt.Sprintf("%d%s", s.UpToDate, minColumnGap)
	row.WriteString(upToDate)
	row.WriteString(strings.Repeat(" ", lipgloss.Width(deploymentUpToDateColumn)+2-lipgloss.Width(upToDate)))

	available := fmt.Sprintf("%d%s", s.Available, minColumnGap)
	row.WriteString(available)
	row.WriteString(strings.Repeat(" ", lipgloss.Width(deploymentAvailableColumn)+2-lipgloss.Width(available)))

	row.WriteString(s.Age)

	deploymentInfo := row.String()

	if index == m.Index() {
		fmt.Fprint(w, v.Styles.SelectedText.Render(deploymentInfo))
	} else {
		fmt.Fprint(w, v.Styles.MainText.Render(deploymentInfo))
	}
}

func getHeader() string {
	var header strings.Builder
	header.WriteString(minColumnGap)
	header.WriteString(deploymentNameColumn)
	header.WriteString(strings.Repeat(" ", nameColumnLen-len(deploymentNameColumn)+len(minColumnGap)-1))
	header.WriteString(deploymentReadyColumn)
	header.WriteString(minColumnGap)
	header.WriteString(deploymentUpToDateColumn)
	header.WriteString(minColumnGap)
	header.WriteString(deploymentAvailableColumn)
	header.WriteString(minColumnGap)
	header.WriteString(deploymentAgeColumn)

	return header.String()
}
