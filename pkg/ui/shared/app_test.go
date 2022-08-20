package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tty2/kubic/pkg/ui/shared/themes"
)

func Test_ResizeAreas(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		app := NewApp(themes.Theme{})

		rq.Equal(0, app.GUI.Areas.MainContent.Height)

		app.GUI.ScreenHeight = 200
		app.ResizeAreas()

		rq.Equal(app.GUI.ScreenHeight-TabsBarHeight-HelpBarHeight, app.GUI.Areas.MainContent.Height)

		app.GUI.ScreenHeight = 100
		app.ResizeAreas()
		rq.Equal(app.GUI.ScreenHeight-TabsBarHeight-HelpBarHeight, app.GUI.Areas.MainContent.Height)
	})
}

func Test_UpdateCallbacks(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		var counter int

		app := NewApp(themes.Theme{})

		rq.Equal(0, counter)

		app.AddUpdateNamespaceCallback(func() {
			counter++
		})

		app.OnUpdateNamespace()
		rq.Equal(1, counter)
		app.OnUpdateNamespace()
		app.OnUpdateNamespace()
		rq.Equal(3, counter)
	})
}
