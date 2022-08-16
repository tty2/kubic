package k8s

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/tty2/kubic/internal/domain"
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
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
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
