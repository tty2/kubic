package shared

import (
	"github.com/tty2/kubic/pkg/ui/shared/themes"
)

type App struct {
	CurrentNamespace  string
	CurrentTab        TabItem
	Styles            *themes.Styles
	KeyMap            *KeyMap
	GUI               GUI
	updateNScallbacks []func()
}

type GUI struct {
	ScreenHeight int
	ScreenWidth  int
	Areas        *uiAreas
}

func NewApp(theme themes.Theme) *App {
	keyMap := GetKeyMaps()
	styles := themes.GetStyle(theme)

	return &App{
		Styles: &styles,
		KeyMap: &keyMap,
		GUI: GUI{
			Areas: initAreas(),
		},
	}
}

func (app *App) ResizeAreas() {
	app.GUI.Areas.MainContent.Height = app.GUI.ScreenHeight -
		(app.GUI.Areas.TabBar.Height + app.GUI.Areas.HelpBar.Height)
}

func (app *App) AddUpdateNamespaceCallback(fn func()) {
	app.updateNScallbacks = append(app.updateNScallbacks, fn)
}

func (app *App) OnUpdateNamespace() {
	for i := range app.updateNScallbacks {
		app.updateNScallbacks[i]()
	}
}
