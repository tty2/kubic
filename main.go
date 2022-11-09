package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tty2/kubic/pkg/config"
	"github.com/tty2/kubic/pkg/k8s"
	"github.com/tty2/kubic/pkg/ui"
	"github.com/tty2/kubic/pkg/ui/shared/themes"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	theme := themes.InitTheme(cfg.ThemePath)

	k8sClient, err := k8s.New(cfg.KubeConfigPath, cfg.LogTail)
	if err != nil {
		return err
	}

	gui, err := ui.New(k8sClient, theme, cfg.UpdateInterval)
	if err != nil {
		return err
	}

	app := tea.NewProgram(
		gui,
		tea.WithAltScreen(),
	)

	return app.Start()
}
