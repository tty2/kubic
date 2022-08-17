package shared

import (
	"github.com/tty2/kubic/pkg/ui/shared/themes"
)

type State struct {
	Namespace    string
	ScreenHeight int
	ScreenWidth  int
	CurrentTab   TabItem
	Theme        *themes.Theme
	Areas        *uiAreas
	KeyMap       *KeyMap
}

func NewState(theme *themes.Theme) *State {
	if theme == nil {
		theme = &themes.DefaultTheme
	}

	keyMap := GetKeyMaps()

	return &State{
		Theme:  theme,
		Areas:  initAreas(),
		KeyMap: &keyMap,
	}
}

func (s *State) ResizeAreas() {
	s.Areas.MainContent.Height = s.ScreenHeight - (s.Areas.TabBar.Height + s.Areas.HelpBar.Height)
}
