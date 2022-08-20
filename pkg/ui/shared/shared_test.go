package shared

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"
)

func Test_Max(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("a > b", func(t *testing.T) {
		t.Parallel()

		rq.Equal(6, Max(6, 4))
	})
	t.Run("a < b", func(t *testing.T) {
		t.Parallel()

		rq.Equal(5, Max(1, 5))
	})
	t.Run("a = b", func(t *testing.T) {
		t.Parallel()

		rq.Equal(2, Max(2, 2))
	})
}

func Test_GetTextWithLen(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("short source", func(t *testing.T) {
		t.Parallel()

		res := GetTextWithLen("three", 7)
		rq.Len(res, 7)

		rq.Equal("three  ", res)
	})
	t.Run("equal source", func(t *testing.T) {
		t.Parallel()

		res := GetTextWithLen("three", 5)
		rq.Len(res, 5)

		rq.Equal("three", res)
	})
	t.Run("too long source", func(t *testing.T) {
		t.Parallel()

		res := GetTextWithLen("three", 3)
		rq.Equal(3, lipgloss.Width(res))

		rq.Equal(fmt.Sprintf("th%s", ellipsis), res)
	})
}
