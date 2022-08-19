/*
Package divider keeps helpers to draw dividers.
*/
package divider

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HorizontalLine is a horizontal line.
func HorizontalLine(width int, st lipgloss.Style) string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom, strings.Repeat(st.Copy().Render("─"), width))
}
