/*
Package tab keeps helpers to create tabs.
*/
package tabs

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/internal/ui/shared"
)

const (
	tabsLeftRightIndents = 2
)

type tab struct {
	shared.TabItem
}

func newTab(ti shared.TabItem) *tab {
	return &tab{
		ti,
	}
}

func (t *tab) getID() shared.TabItem {
	return t.TabItem
}

func (t *tab) render(active bool) string {
	if active {
		return activeTab.Render(t.String())
	}

	return inactiveTab.Render(t.String())
}

func renderTabBar(screenWidth int, titles []string) string {
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		titles...,
	)

	gap := tabGap.Render(strings.Repeat(" ", shared.Max(0, screenWidth-lipgloss.Width(row)-tabsLeftRightIndents)))

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap))
}
