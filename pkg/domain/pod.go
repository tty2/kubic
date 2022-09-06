package domain

import "time"

type Pod struct {
	// list info
	Name     string
	Ready    string
	Status   string
	Restarts int
	Age      string
	// meta
	Meta PodMeta
	// spec
	Spec PodSpec
	// status
	StatusInfo PodStatusInfo
}

type PodMeta struct {
	Created time.Time
	Labels  map[string]string
	Owners  []OwnerInfo
}

type PodSpec struct {
	DNSPolicy                     string
	RestartPolicy                 string
	SchedulerName                 string
	TerminationGracePeriodSeconds int64
	Containers                    []Container
}

type PodStatusInfo struct {
	Phase      string
	QosClass   string
	HostIP     string
	PodIP      string
	PodIPs     []string
	Conditions []string
}

type OwnerInfo struct {
	Kind string
	Name string
}
