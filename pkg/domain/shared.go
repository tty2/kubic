package domain

type Container struct {
	Name                   string
	Image                  string
	ImagePullPolicy        string
	TerminationMessagePath string
	ENVs                   []ContainerEnv
}

type ContainerEnv struct {
	Name  string
	Value string
}
