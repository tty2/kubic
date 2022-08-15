/*
Package help keeps the logic for bottom help bar.
*/
package help

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/internal/ui/shared"
	"github.com/tty2/kubic/internal/ui/shared/themes"
)

// nolint:gochecknoglobals // this is used only in this package
var (
	helpTextStyle = lipgloss.NewStyle().Foreground(themes.DefaultTheme.InactiveText)
	helpStyle     = lipgloss.NewStyle().
			Height(shared.HelpBarHeight).
			BorderForeground(themes.DefaultTheme.InactiveText).BorderTop(true)
)
