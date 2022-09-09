package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/tty2/kubic/pkg/domain"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Set *kubernetes.Clientset
}

func New(configPath string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Set: clientSet,
	}, nil
}

func (c *Client) GetNamespaces(ctx context.Context) ([]domain.Namespace, error) {
	apiResp, err := c.Set.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	ns := make([]domain.Namespace, len(apiResp.Items))
	for i := range apiResp.Items {
		ns[i].Name = apiResp.Items[i].Name
		ns[i].Status = string(apiResp.Items[i].Status.Phase)

		age := time.Now().Unix() - apiResp.Items[i].GetCreationTimestamp().Unix()
		ns[i].Age = ageToString(age)
	}

	return ns, nil
}

func (c *Client) GetDeployments(ctx context.Context, namespace string) ([]domain.Deployment, error) {
	apiResp, err := c.Set.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	deps := make([]domain.Deployment, len(apiResp.Items))
	for i := range apiResp.Items {
		deps[i].Name = apiResp.Items[i].Name
		deps[i].Ready = fmt.Sprintf("%d/%d", apiResp.Items[i].Status.ReadyReplicas, apiResp.Items[i].Status.Replicas)
		deps[i].UpdatedReplicas = int(apiResp.Items[i].Status.UpdatedReplicas)
		deps[i].AvailableReplicas = int(apiResp.Items[i].Status.AvailableReplicas)
		deps[i].ReadyReplicas = int(apiResp.Items[i].Status.ReadyReplicas)
		deps[i].Tolerations = len(apiResp.Items[i].Spec.Template.Spec.Tolerations)

		age := time.Now().Unix() - apiResp.Items[i].GetCreationTimestamp().Unix()
		deps[i].Age = ageToString(age)

		deps[i].Labels = apiResp.Items[i].Labels
		deps[i].Created = apiResp.Items[i].ObjectMeta.CreationTimestamp.Time

		// populate meta
		deps[i].Meta.Strategy = string(apiResp.Items[i].Spec.Strategy.Type)
		deps[i].Meta.DNSPolicy = string(apiResp.Items[i].Spec.Template.Spec.DNSPolicy)
		deps[i].Meta.RestartPolicy = string(apiResp.Items[i].Spec.Template.Spec.RestartPolicy)
		deps[i].Meta.SchedulerName = apiResp.Items[i].Spec.Template.Spec.SchedulerName
		deps[i].Meta.TerminationGracePeriodSeconds = *apiResp.Items[i].Spec.Template.Spec.TerminationGracePeriodSeconds
		deps[i].Meta.Containers = toDomainContainers(apiResp.Items[i].Spec.Template.Spec.Containers)
	}

	return deps, nil
}

func (c *Client) GetPods(ctx context.Context, namespace string) ([]domain.Pod, error) {
	apiResp, err := c.Set.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods := make([]domain.Pod, len(apiResp.Items))
	for i := range apiResp.Items {
		pods[i].Name = apiResp.Items[i].Name
		pods[i].Ready = getReadyOfListCont(apiResp.Items[i].Status.ContainerStatuses)
		pods[i].Status = string(apiResp.Items[i].Status.Phase)
		pods[i].Restarts = getRestartsCount(apiResp.Items[i].Status.ContainerStatuses)

		age := time.Now().Unix() - apiResp.Items[i].GetCreationTimestamp().Unix()
		pods[i].Age = ageToString(age)

		// populate meta
		pods[i].Meta.Created = apiResp.Items[i].CreationTimestamp.Time
		pods[i].Meta.Labels = apiResp.Items[i].Labels
		pods[i].Meta.Owners = toDomainOwnerInfoList(apiResp.Items[i].OwnerReferences)

		// populate spec
		pods[i].Spec.DNSPolicy = string(apiResp.Items[i].Spec.DNSPolicy)
		pods[i].Spec.RestartPolicy = string(apiResp.Items[i].Spec.RestartPolicy)
		pods[i].Spec.SchedulerName = apiResp.Items[i].Spec.SchedulerName
		pods[i].Spec.TerminationGracePeriodSeconds = *apiResp.Items[i].Spec.TerminationGracePeriodSeconds
		pods[i].Spec.Containers = toDomainContainers(apiResp.Items[i].Spec.Containers)

		// populate status info
		pods[i].StatusInfo.Phase = string(apiResp.Items[i].Status.Phase)
		pods[i].StatusInfo.QosClass = string(apiResp.Items[i].Status.QOSClass)
		pods[i].StatusInfo.HostIP = apiResp.Items[i].Status.HostIP
		pods[i].StatusInfo.PodIP = apiResp.Items[i].Status.PodIP
		pods[i].StatusInfo.PodIPs = podIPsToDomainList(apiResp.Items[i].Status.PodIPs)
		pods[i].StatusInfo.Conditions = conditionsToDomainList(apiResp.Items[i].Status.Conditions)
	}

	return pods, nil
}

func (c *Client) PodsLog(ctx context.Context, namespace, name string) []byte {
	data, err := c.Set.CoreV1().
		Pods(namespace).
		GetLogs(name, &corev1.PodLogOptions{}).
		Do(ctx).
		Raw()

	if err != nil {
		return []byte("")
	}

	return data
}

// nolint gomnd: numbers are obvious here
func ageToString(age int64) string {
	switch {
	case age < 60:
		return fmt.Sprintf("%ds", age)
	case age >= 60 && age < 3600:
		return fmt.Sprintf("%dm", age/60)
	case age >= 3600 && age < 86400:
		return fmt.Sprintf("%dh", age/60/60)
	case age >= 86400 && age < 2678400:
		return fmt.Sprintf("%dd", age/60/60/24)
	case age >= 2678400 && age < 31536000:
		return fmt.Sprintf("%dM", age/60/60/24/31)
	case age >= 31536000:
		return fmt.Sprintf("%dY", age/60/60/24/365)
	default:
		return "-"
	}
}

func getReadyOfListCont(ss []corev1.ContainerStatus) string {
	var number int
	var ready int
	for i := range ss {
		if ss[i].Ready {
			ready++
		}
		number++
	}

	return fmt.Sprintf("%d/%d", ready, number)
}

func getRestartsCount(ss []corev1.ContainerStatus) int {
	var number int32
	for i := range ss {
		number += ss[i].RestartCount
	}

	return int(number)
}

func getEnvs(envs []corev1.EnvVar) []domain.ContainerEnv {
	cEnvs := make([]domain.ContainerEnv, len(envs))
	for i := range envs {
		cEnvs[i].Name = envs[i].Name
		cEnvs[i].Value = envs[i].Value
	}

	return cEnvs
}

func toDomainContainers(cc []corev1.Container) []domain.Container {
	domainContainers := make([]domain.Container, len(cc))
	for i := range cc {
		domainContainers[i].Name = cc[i].Name
		domainContainers[i].Image = cc[i].Image
		domainContainers[i].ImagePullPolicy = string(cc[i].ImagePullPolicy)
		domainContainers[i].TerminationMessagePath = cc[i].TerminationMessagePath
		domainContainers[i].ENVs = getEnvs(cc[i].Env)
	}

	return domainContainers
}

func toDomainOwnerInfoList(oref []metav1.OwnerReference) []domain.OwnerInfo {
	resp := make([]domain.OwnerInfo, len(oref))
	for i := range oref {
		resp[i].Kind = oref[i].Kind
		resp[i].Name = oref[i].Name
	}

	return resp
}

func podIPsToDomainList(ips []corev1.PodIP) []string {
	resp := make([]string, len(ips))
	for i := range ips {
		resp[i] = ips[i].IP
	}

	return resp
}

func conditionsToDomainList(conds []corev1.PodCondition) []string {
	resp := make([]string, len(conds))
	for i := range conds {
		resp[i] = string(conds[i].Type)
	}

	return resp
}
