package shared

import (
	"github.com/tty2/kubic/pkg/ui/shared/themes"
)

type App struct {
	CurrentNamespace string
	CurrentTab       TabItem
	Theme            *themes.Theme
	KeyMap           *KeyMap
	GUI              GUI
}

type GUI struct {
	ScreenHeight int
	ScreenWidth  int
	Areas        *uiAreas
}

func NewApp(theme *themes.Theme) *App {
	if theme == nil {
		theme = &themes.DefaultTheme
	}

	keyMap := GetKeyMaps()

	return &App{
		Theme:  theme,
		KeyMap: &keyMap,
		GUI: GUI{
			Areas: initAreas(),
		},
	}
}

func (app *App) ResizeAreas() {
	app.GUI.Areas.MainContent.Height =
		app.GUI.ScreenHeight - (app.GUI.Areas.TabBar.Height + app.GUI.Areas.HelpBar.Height)
}
