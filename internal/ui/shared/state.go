package shared

import "github.com/tty2/kubic/internal/ui/shared/themes"

type State struct {
	Namespace    string
	ScreenHeight int
	ScreenWidth  int
	CurrentTab   TabItem
	Theme        *themes.Theme
	Areas        *uiAreas
}

func NewState(theme *themes.Theme) *State {
	if theme == nil {
		theme = &themes.DefaultTheme
	}

	return &State{
		Theme: theme,
		Areas: initAreas(),
	}
}

func (s *State) ResizeAreas() {
	s.Areas.TabBar.Coords.X2 = s.ScreenWidth
	s.Areas.MainContent.Coords.X2 = s.ScreenWidth
	s.Areas.HelpBar.Coords.X2 = s.ScreenWidth

	s.Areas.MainContent.Coords.Y2 = s.ScreenHeight - (s.Areas.TabBar.Height + s.Areas.HelpBar.Height)
	s.Areas.HelpBar.Coords.Y1 = s.ScreenHeight - s.Areas.HelpBar.Height - 1
	s.Areas.HelpBar.Coords.Y2 = s.ScreenHeight
}
