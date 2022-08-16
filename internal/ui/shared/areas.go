package shared

const (
	TabsBarHeight = 3
	HelpBarHeight = 4
	FullScreen    = -1
)

type uiArea struct {
	Height int
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
		},
		MainContent: uiArea{},
		HelpBar: uiArea{
			Height: HelpBarHeight,
		},
	}

	return &ua
}
