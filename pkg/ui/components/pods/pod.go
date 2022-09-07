package pods

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/domain"
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

// nolint gochecknoglobals: used here on purpose
var boldText = lipgloss.NewStyle().Bold(true)

type (
	pod struct {
		Name       string
		Ready      string
		Status     string
		Restarts   int
		Age        string
		Meta       domain.PodMeta
		Spec       domain.PodSpec
		StatusInfo domain.PodStatusInfo
		Styles     *themes.Styles
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

func (p *pod) renderInfo() string {
	var info strings.Builder
	info.WriteString(boldText.Render("Name"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(p.Name)

	info.WriteString("\n")
	info.WriteString(boldText.Render("Created"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(p.Meta.Created.Format(shared.TimeFormat))
	info.WriteString("\n")
	info.WriteString(boldText.Render("Labels"))
	info.WriteString("\n")

	for k, v := range p.Meta.Labels {
		info.WriteString(minColumnGap)
		info.WriteString(k)
		info.WriteString(": ")
		info.WriteString(v)
		info.WriteString("\n")
	}

	info.WriteString(boldText.Render("Restarts"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(fmt.Sprint(p.Restarts))
	info.WriteString("\n")

	info.WriteString(renderSpec(p.Spec))

	info.WriteString(renderStatusInfo(p.StatusInfo))

	return info.String()
}

func renderSpec(spec domain.PodSpec) string {
	var info strings.Builder

	info.WriteString(boldText.Render("DNS policy"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(spec.DNSPolicy)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Restart policy"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(spec.RestartPolicy)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Scheduler"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(spec.SchedulerName)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Termination graceful period"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(fmt.Sprintf("%d sec", spec.TerminationGracePeriodSeconds))
	info.WriteString("\n")

	info.WriteString(renderContainersInfo(spec.Containers))

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
		info.WriteString(minColumnGap)
		if len(cc[i].ENVs) > 0 {
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

func renderStatusInfo(si domain.PodStatusInfo) string {
	var info strings.Builder

	info.WriteString(boldText.Render("Phase"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(si.Phase)
	info.WriteString("\n")

	info.WriteString(boldText.Render("QOS Class"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(si.QosClass)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Host IP"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(si.HostIP)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Pod IP"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(si.PodIP)
	info.WriteString("\n")

	info.WriteString(boldText.Render("Pod IPs"))
	info.WriteString("\n")
	for i := range si.PodIPs {
		info.WriteString(minColumnGap)
		info.WriteString(si.PodIPs[i])
		info.WriteString("\n")
	}

	info.WriteString(boldText.Render("Conditions"))
	info.WriteString("\n")
	info.WriteString(minColumnGap)
	info.WriteString(strings.Join(si.Conditions, ", "))
	return info.String()
}
