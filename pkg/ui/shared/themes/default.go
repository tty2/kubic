/*
Package themes keeps default styles.
*/
package themes

import (
	"os"
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/tty2/kubic/pkg/css"
)

var validColor = regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")

type Styles struct {
	InitStyle lipgloss.Style
	// text
	MainText     lipgloss.Style
	SelectedText lipgloss.Style
	InactiveText lipgloss.Style
	HelpBar      lipgloss.Style
	// signs
	NamespaceSign lipgloss.Style
	// tabs
	Borders     lipgloss.Style
	InactiveTab lipgloss.Style
	ActiveTab   lipgloss.Style
	TabsGap     lipgloss.Style
	// margin
	TextRightMargin int
	TextLeftMargin  int
}

// Theme is a struct to keep all the application styles.
type Theme struct {
	MainText      lipgloss.Color
	SelectedText  lipgloss.Color
	InactiveText  lipgloss.Color
	Borders       lipgloss.Color
	NamespaceSign lipgloss.Color
}

// DefaultTheme is an application default theme.
// nolint:gochecknoglobals // global on purpose
var (
	defaultTheme = Theme{
		MainText:      lipgloss.Color("#E2E1ED"),
		SelectedText:  lipgloss.Color("#EE6FF8"), // #AD58B4
		InactiveText:  lipgloss.Color("#5C5C5C"),
		Borders:       lipgloss.Color("#7D56F4"),
		NamespaceSign: lipgloss.Color("#6aa84f"),
	}

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
)

func GetStyle(theme Theme) Styles {
	return Styles{
		MainText:     lipgloss.NewStyle().Foreground(theme.MainText),
		SelectedText: lipgloss.NewStyle().Foreground(theme.SelectedText),
		InactiveText: lipgloss.NewStyle().Foreground(theme.InactiveText),
		Borders:      lipgloss.NewStyle().Foreground(theme.Borders),
		HelpBar: lipgloss.NewStyle().BorderForeground(theme.InactiveText).
			BorderTop(true),
		// signs
		NamespaceSign: lipgloss.NewStyle().Foreground(theme.NamespaceSign),
		// tabs
		InactiveTab: lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(theme.Borders).
			Padding(0, 1),
		ActiveTab: lipgloss.NewStyle().
			Border(activeTabBorder, true).
			BorderForeground(theme.Borders).
			Padding(0, 1),
		TabsGap: lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderBottom(true).
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false).
			BorderForeground(theme.Borders).
			Padding(0, 1),

		TextRightMargin: 1,
		TextLeftMargin:  2,
	}
}

func validHexColor(st string) bool {
	return validColor.Match([]byte(st))
}

func InitTheme(ph string) Theme {
	data, err := os.ReadFile(ph)
	if err != nil {
		initTheme(nil)
	}

	styleSheet, err := css.Unmarshal([]byte(data))
	if err != nil {
		initTheme(nil)
	}

	return initTheme(styleSheet)
}

func initTheme(styleSheet map[css.Selector]map[string]string) Theme {
	theme := defaultTheme

	if styleSheet == nil {
		return theme
	}

	mainText := css.CSSStyle("color", styleSheet[".main-text"])
	if validHexColor(mainText) {
		theme.MainText = lipgloss.Color(mainText)
	}
	selectedText := css.CSSStyle("color", styleSheet[".selected-text"])
	if validHexColor(selectedText) {
		theme.SelectedText = lipgloss.Color(selectedText)
	}
	inactiveText := css.CSSStyle("color", styleSheet[".inactive-text"])
	if validHexColor(inactiveText) {
		theme.InactiveText = lipgloss.Color(inactiveText)
	}
	borders := css.CSSStyle("color", styleSheet[".tab-borders"])
	if validHexColor(borders) {
		theme.Borders = lipgloss.Color(borders)
	}
	sign := css.CSSStyle("color", styleSheet[".namespace-sign"])
	if validHexColor(sign) {
		theme.NamespaceSign = lipgloss.Color(sign)
	}

	return theme
}
