package config

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

type Config struct {
	KubeConfigPath string
}

func New() Config {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("cfg", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kube config file")
	} else {
		kubeConfig = flag.String("cfg", "", "absolute path to the kube config file")
	}
	flag.Parse()

	return Config{
		KubeConfigPath: *kubeConfig,
	}
}
