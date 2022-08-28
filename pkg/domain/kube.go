package domain

type Namespace struct {
	Name   string
	Status string
	Age    string
}

type Deployment struct {
	Name              string
	Ready             string
	UpdatedReplicas   int
	AvailableReplicas int
	ReadyReplicas     int
	Tolerations       int
	Age               string
	Labels            map[string]string
	Meta              DeploymentMeta
}

type DeploymentMeta struct {
	Strategy                      string
	DNSPolicy                     string
	RestartPolicy                 string
	SchedulerName                 string
	TerminationGracePeriodSeconds int64
	Containers                    []Container
}

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

type Pod struct {
	Name     string
	Ready    string
	Status   string
	Restarts int
	Age      string
	Labels   map[string]string
}
