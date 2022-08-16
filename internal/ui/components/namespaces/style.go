package namespaces

import (
	"github.com/charmbracelet/lipgloss"
)

const indent = 2

// nolint:gochecknoglobals // used only in this package
var (
	subtle            = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	signActiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6aa84f"))
	signInactiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#c32a2a"))
	listStyle         = lipgloss.NewStyle().
				MarginRight(indent).PaddingRight(indent)

	normalItem = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})

	selectedItem = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})
)
