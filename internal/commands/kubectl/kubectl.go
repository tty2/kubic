package kubectl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/tty2/kubic/internal/domain"
)

func GetNamespaces() ([]domain.Namespace, error) {
	nn := []domain.Namespace{}

	cmd := exec.Command("kubectl", "get", "namespace", "-o", "json")
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can't execute %s: %s", cmd.String(), cmdErr.String())
	}

	err = json.Unmarshal(output, &nn)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal: %v", err)
	}

	return nn, nil
}

func GetDeployments(namespace string) ([]domain.Deployment, error) {
	dd := []domain.Deployment{}

	cmd := exec.Command("kubectl", "get", "deployments", "--namespace", namespace, "-o", "json")
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can't execute %s: %s", cmd.String(), cmdErr.String())
	}

	err = json.Unmarshal(output, &dd)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal: %v", err)
	}

	return dd, nil
}

func DeploymentInfo(namespace, name string) (domain.Deployment, error) {
	var dep domain.Deployment

	cmd := exec.Command("kubectl", "get", "deployment", "--namespace", namespace, name, "-o", "json")
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	output, err := cmd.Output()
	if err != nil {
		return dep, fmt.Errorf("can't execute %s: %s", cmd.String(), cmdErr.String())
	}

	err = json.Unmarshal(output, &dep)
	if err != nil {
		return dep, fmt.Errorf("can't unmarshal: %v", err)
	}

	return dep, nil
}

func GetPods(namespace string) ([]domain.Pod, error) {
	pp := []domain.Pod{}

	cmd := exec.Command("kubectl", "get", "pods", "--namespace", namespace, "-o", "json")
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can't execute %s: %s", cmd.String(), cmdErr.String())
	}

	err = json.Unmarshal(output, &pp)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal: %v", err)
	}

	return pp, nil
}

func PodInfo(namespace, name string) (domain.Pod, error) {
	var pod domain.Pod

	cmd := exec.Command("kubectl", "get", "pod", "--namespace", namespace, name, "-o", "json")
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	output, err := cmd.Output()
	if err != nil {
		return pod, fmt.Errorf("can't execute %s: %s", cmd.String(), cmdErr.String())
	}

	err = json.Unmarshal(output, &pod)
	if err != nil {
		return pod, fmt.Errorf("can't unmarshal: %v", err)
	}

	return pod, nil
}

func PodsLog(namespace, name string) ([]byte, error) {
	cmd := exec.Command("kubectl", "logs", "--namespace", namespace, name)
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can't execute %s: %s", cmd.String(), cmdErr.String())
	}

	return output, nil
}
