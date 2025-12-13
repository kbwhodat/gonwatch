package main

import (
	"log"
	"os"
	"gonwatch/models"
	tea "github.com/charmbracelet/bubbletea"

)

func main() {
	// Initialize debug mode if requested
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		f, err := tea.LogToFile("/tmp/debug.log", "debug")
		if err != nil {
			log.Printf("Failed to initialize debug logging: %v", err)
		} else {
			defer f.Close()
		}
	}

	// Start the TUI application
	p := tea.NewProgram(models.ChoiceModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
