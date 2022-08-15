package shared

const (
	TabsBarHeight = 3
	HelpBarHeight = 2
	FullScreen    = -1
)

type uiArea struct {
	Height int
	Coords Coords
}

type uiAreas struct {
	TabBar      uiArea
	MainContent uiArea
	HelpBar     uiArea
}

func initAreas() *uiAreas {
	ua := uiAreas{
		TabBar: uiArea{
			Height: TabsBarHeight,
			Coords: Coords{
				X1: 0,
				Y1: 0,
				X2: FullScreen,
				Y2: TabsBarHeight,
			},
		},
		MainContent: uiArea{
			Coords: Coords{
				X1: 0,
				Y1: TabsBarHeight + 1,
				X2: FullScreen,
				// can't set Y2 without screen height
			},
		},
		HelpBar: uiArea{
			Height: HelpBarHeight,
			Coords: Coords{
				X1: 0,
				X2: FullScreen,
				// can't set Y1, Y2 without screen height
			},
		},
	}

	return &ua
}

func (ua *uiArea) GetWidth() int {
	return ua.Coords.X2 - ua.Coords.X1
}

func (ua *uiArea) GetHeight() int {
	return ua.Coords.Y2 - ua.Coords.Y1
}
