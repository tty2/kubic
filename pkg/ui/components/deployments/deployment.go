package deployments

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/domain"
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

// nolint gochecknoglobals: used here on purpose
var boldText = lipgloss.NewStyle().Bold(true)

type (
	deployment struct {
		Name              string
		Ready             string
		UpdatedReplicas   int
		AvailableReplicas int
		ReadyReplicas     int
		Tolerations       int
		Age               string
		Labels            map[string]string
		Styles            *themes.Styles
		Created           time.Time
		Meta              domain.DeploymentMeta
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

	upToDate := fmt.Sprintf("%d", s.UpdatedReplicas)
	row.WriteString(upToDate)
	row.WriteString(strings.Repeat(" ", upToDateColumnLen-lipgloss.Width(upToDate)))
	row.WriteString(minColumnGap)

	available := fmt.Sprintf("%d", s.AvailableReplicas)
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

func (d *deployment) renderInfo() string {
	var info strings.Builder
	info.WriteString(boldText.Render("Name"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(d.Name)
	info.WriteString("\n")
	info.WriteString(boldText.Render("Created"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(d.Created.Format(shared.TimeFormat))
	info.WriteString("\n")
	info.WriteString(boldText.Render("Labels"))
	info.WriteString("\n")

	for k, v := range d.Labels {
		info.WriteString(minColumnGap)
		info.WriteString(k)
		info.WriteString(": ")
		info.WriteString(v)
		info.WriteString("\n")
	}

	info.WriteString(boldText.Render("Replicas"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(fmt.Sprintf("Available: %d\n", d.AvailableReplicas))
	info.WriteString(minColumnGap)
	info.WriteString(fmt.Sprintf("Ready: %d\n", d.ReadyReplicas))
	info.WriteString(minColumnGap)
	info.WriteString(fmt.Sprintf("Updated: %d\n", d.UpdatedReplicas))

	info.WriteString(boldText.Render("Tolerations"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(fmt.Sprintf("Total: %d\n", d.Tolerations))

	info.WriteString(renderMeta(d.Meta))
	info.WriteString(renderContainersInfo(d.Meta.Containers))

	return info.String()
}

func renderMeta(meta domain.DeploymentMeta) string {
	var info strings.Builder

	info.WriteString(boldText.Render("Strategy"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(meta.Strategy)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Restart Policy"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(meta.RestartPolicy)
	info.WriteString("\n")

	info.WriteString(boldText.Render("DNS Policy"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(meta.DNSPolicy)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Scheduler"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(meta.SchedulerName)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Termination Grace Period"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(fmt.Sprintf("%d seconds", meta.TerminationGracePeriodSeconds))
	info.WriteString("\n")

	return info.String()
}

func renderContainersInfo(cc []domain.Container) string {
	var info strings.Builder

	info.WriteString(boldText.Render("Containers"))
	info.WriteString("\n")
	for i := range cc {
		info.WriteString(minColumnGap)
		info.WriteString(fmt.Sprintf("Name: %s", cc[i].Name))
		info.WriteString("\n")
		info.WriteString(minColumnGap)
		info.WriteString(fmt.Sprintf("Image: %s", cc[i].Image))
		info.WriteString("\n")
		info.WriteString(minColumnGap)
		info.WriteString(fmt.Sprintf("Policy: %s", cc[i].ImagePullPolicy))
		info.WriteString("\n")
		info.WriteString(minColumnGap)
		info.WriteString(fmt.Sprintf("Termination message path: %s", cc[i].TerminationMessagePath))
		info.WriteString("\n")
		if len(cc[i].ENVs) > 0 {
			info.WriteString(minColumnGap)
			info.WriteString(boldText.Render("Envs"))
			info.WriteString("\n")
			for j := range cc[i].ENVs {
				info.WriteString(minColumnGap)
				info.WriteString(minColumnGap)
				info.WriteString(cc[i].ENVs[j].Name)
				info.WriteString("\n")
			}
		}
	}

	return info.String()
}
