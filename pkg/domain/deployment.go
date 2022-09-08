package domain

import "time"

type Deployment struct {
	Name              string
	Ready             string
	UpdatedReplicas   int
	AvailableReplicas int
	ReadyReplicas     int
	Tolerations       int
	Age               string
	Labels            map[string]string
	Created           time.Time
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
