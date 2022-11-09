/*
Package config keeps all config files.
*/
package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"k8s.io/client-go/util/homedir"
)

// SearchConfig for application.
// nolint lll // we need all of tags here. If we add CR here we'll catch structtag: *** key:"value" pairs not separated by spaces (govet)
type Config struct {
	KubeConfigPath string `short:"c" long:"config" env:"KUBIC_KUBERNETES_CONFIG_PATH" description:"kubernetes config file path"`
	ThemePath      string `short:"t" long:"theme" env:"KUBIC_THEME_FILE_PATH" default:"./style.json" description:"theme file path"`
	UpdateInterval int    `short:"u" long:"update_interval" env:"KUBIC_UPDATE_INTERVAL" default:"3" description:"update interval in seconds"`
	LogTail        int64  `short:"l" long:"log_tail" env:"KUBIC_LOG_TAIL_LINES" default:"100" description:"log tail lines"`
}

// New creates a new config.
func New() (Config, error) {
	var config Config

	parser := flags.NewParser(&config, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		return Config{}, fmt.Errorf("couldn't setup config: %w", err)
	}

	if config.KubeConfigPath == "" {
		home := homedir.HomeDir()
		if home == "" {
			return Config{}, errors.New("setup file for kubernetes config or add it to ~/.kube/config")
		}
		config.KubeConfigPath = filepath.Join(home, ".kube", "config")
	}

	return config, nil
}
