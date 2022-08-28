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
	Envs              []string
	Labels            map[string]string
}

type Pod struct {
	Name     string
	Ready    string
	Status   string
	Restarts int
	Age      string
	Labels   map[string]string
}
