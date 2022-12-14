package shared

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	ellipsis   = "…"
	TimeFormat = "2006-01-02 15:04:05"
)

// Max returns max of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func GetTextWithLen(source string, length int) string {
	var name string
	lenName := len(source)
	switch {
	case lenName > length:
		name = fmt.Sprintf("%s%s", source[:length-1], ellipsis)
	case lenName < length:
		name = fmt.Sprintf("%s%s", source, strings.Repeat(" ", length-lipgloss.Width(source)))
	default:
		name = source
	}

	return name
}
