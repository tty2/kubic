package k8s

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/tty2/kubic/pkg/domain"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	Set *kubernetes.Clientset
}

func New() (*Client, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("cfg", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kube config file")
	} else {
		kubeconfig = flag.String("cfg", "", "absolute path to the kube config file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Set: clientset,
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
		deps[i].UpToDate = int(apiResp.Items[i].Status.UpdatedReplicas)
		deps[i].Available = int(apiResp.Items[i].Status.AvailableReplicas)

		age := time.Now().Unix() - apiResp.Items[i].GetCreationTimestamp().Unix()
		deps[i].Age = ageToString(age)

		deps[i].Labels = apiResp.Items[i].Labels
	}

	return deps, nil
}

func (c *Client) DeploymentInfo(ctx context.Context, namespace, name string) (domain.Deployment, error) {
	return domain.Deployment{}, nil
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

		pods[i].Labels = apiResp.Items[i].Labels
	}

	return pods, nil
}

func (c *Client) PodInfo(ctx context.Context, namespace, name string) (domain.Pod, error) {
	return domain.Pod{}, nil
}

func PodsLog(ctx context.Context, namespace, name string) ([]byte, error) {
	return nil, nil
}

func ageToString(age int64) string {
	switch {
	case age < 60:
		return fmt.Sprintf("%ds", age)
	case age >= 60 && age < 3600:
		return fmt.Sprintf("%dm", age/60)
	case age >= 3600 && age < 86400:
		return fmt.Sprintf("%dh", age/60/60)
	case age >= 86400 && age < 2592000:
		return fmt.Sprintf("%dd", age/60/60/24)
	case age >= 2592000 && age < 31536000:
		return fmt.Sprintf("%dM", age/60/60/24/30)
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
