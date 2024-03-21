package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enkdress/go-todo/tui/pkg/components"
)

func main() {
	kanban := components.InitialModel()

	if _, err := tea.NewProgram(kanban, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
