/*
Package themes keeps default styles.
*/
package themes

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

const (
	textRightMargin  = 1
	textLeftMargin   = 2
	listRightPadding = 3
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
	// list border style
	ListRightBorder lipgloss.Style
	// margin
	TextRightMargin int
	TextLeftMargin  int
	// padding
	ListRightPadding int
}

// Theme is a struct to keep all the application styles.
type Theme struct {
	MainText      lipgloss.Color `json:"main-text"`
	SelectedText  lipgloss.Color `json:"selected-text"`
	InactiveText  lipgloss.Color `json:"inactive-text"`
	Borders       lipgloss.Color `json:"tab-borders"`
	NamespaceSign lipgloss.Color `json:"namespace-sign"`
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

// nolint gomnd: default values
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

		ListRightBorder: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(theme.InactiveText).
			PaddingRight(listRightPadding).
			MarginLeft(textLeftMargin),

		TextRightMargin: textRightMargin,
		TextLeftMargin:  textLeftMargin,
	}
}

func validHexColor(st string) bool {
	return validColor.MatchString(st)
}

func InitTheme(ph string) Theme {
	data, err := os.ReadFile(filepath.Clean(ph))
	if err != nil {
		initTheme(Theme{})
	}

	th := Theme{}

	err = json.Unmarshal(data, &th)
	if err != nil {
		initTheme(Theme{})
	}

	return initTheme(th)
}

func initTheme(th Theme) Theme {
	theme := defaultTheme

	if validHexColor(string(th.MainText)) {
		theme.MainText = th.MainText
	}
	if validHexColor(string(th.SelectedText)) {
		theme.SelectedText = th.SelectedText
	}
	if validHexColor(string(th.InactiveText)) {
		theme.InactiveText = th.InactiveText
	}
	if validHexColor(string(th.Borders)) {
		theme.Borders = th.Borders
	}
	if validHexColor(string(th.NamespaceSign)) {
		theme.NamespaceSign = th.NamespaceSign
	}

	return theme
}
