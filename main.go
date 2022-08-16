package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tty2/kubic/internal/commands/k8s"
	"github.com/tty2/kubic/internal/config"
	"github.com/tty2/kubic/internal/ui"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	k8sClient, err := k8s.New()
	if err != nil {
		return err
	}

	gui, err := ui.New(&config.Config{}, k8sClient)
	if err != nil {
		return err
	}

	app := tea.NewProgram(
		gui,
		tea.WithAltScreen(),
	)

	return app.Start()
}
