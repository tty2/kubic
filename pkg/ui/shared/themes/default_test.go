package themes

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"
)

func Test_GetStyle(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		th := GetStyle(Theme{})
		rq.Equal(1, th.TextRightMargin)
		rq.Equal(2, th.TextLeftMargin)
	})
}

func Test_validHexColor(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rq.True(validHexColor("#000"))
		rq.True(validHexColor("#000000"))
		rq.True(validHexColor("#fff"))
		rq.True(validHexColor("#ffffff"))
		rq.True(validHexColor("#E2E1ED"))
		rq.True(validHexColor("#e2e1ed"))
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		rq.False(validHexColor("red"))
		rq.False(validHexColor("#0000"))
		rq.False(validHexColor("fff"))
		rq.False(validHexColor("#fffffff"))
		rq.False(validHexColor("#E2$1ED"))
		rq.False(validHexColor("e2e1ed"))
	})
}

func Test_initTheme(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("completely default", func(t *testing.T) {
		t.Parallel()

		th := initTheme(Theme{})

		rq.Equal(defaultTheme.MainText, th.MainText)
		rq.Equal(defaultTheme.SelectedText, th.SelectedText)
		rq.Equal(defaultTheme.InactiveText, th.InactiveText)
		rq.Equal(defaultTheme.Borders, th.Borders)
		rq.Equal(defaultTheme.NamespaceSign, th.NamespaceSign)
	})

	t.Run("completely custom", func(t *testing.T) {
		t.Parallel()

		th := initTheme(Theme{
			MainText:      lipgloss.Color("#000"),
			SelectedText:  lipgloss.Color("#000"),
			InactiveText:  lipgloss.Color("#000"),
			Borders:       lipgloss.Color("#000"),
			NamespaceSign: lipgloss.Color("#000"),
		})

		rq.NotEqual(defaultTheme.MainText, th.MainText)
		rq.NotEqual(defaultTheme.SelectedText, th.SelectedText)
		rq.NotEqual(defaultTheme.InactiveText, th.InactiveText)
		rq.NotEqual(defaultTheme.Borders, th.Borders)
		rq.NotEqual(defaultTheme.NamespaceSign, th.NamespaceSign)
	})

	t.Run("partially custom", func(t *testing.T) {
		t.Parallel()

		th := initTheme(Theme{
			MainText:     lipgloss.Color("#000"),
			InactiveText: lipgloss.Color("#000"),
		})

		rq.NotEqual(defaultTheme.MainText, th.MainText)
		rq.Equal(defaultTheme.SelectedText, th.SelectedText)
		rq.NotEqual(defaultTheme.InactiveText, th.InactiveText)
		rq.Equal(defaultTheme.Borders, th.Borders)
		rq.Equal(defaultTheme.NamespaceSign, th.NamespaceSign)
	})
}
