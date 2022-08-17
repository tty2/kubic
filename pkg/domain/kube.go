package domain

type Namespace struct {
	Name   string
	Status string
	Age    string
}

type Deployment struct {
	Name      string
	Ready     string
	UpToDate  int
	Available int
	Age       string
	Labels    map[string]string
}

type Pod struct {
	Name     string
	Ready    string
	Status   string
	Restarts int
	Age      string
	Labels   map[string]string
}
