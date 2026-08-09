package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/mahadiksahil60/dev-launcher/tui"
)

func main() {

	status := tui.CheckDependencies()

	if err := tui.DependencyError(status); err != nil {
		fmt.Println(err)
		fmt.Println()
		fmt.Println("NEXUS cannot start.")
		fmt.Println("Press Enter to exit.")

		fmt.Scanln()

		os.Exit(1)
	}

	model := tui.NewModel()

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
