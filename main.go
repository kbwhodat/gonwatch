package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"gonwatch/models"
)

var version = "1.0.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version)
		return
	}

	f, err := tea.LogToFile("/tmp/debug.log", "debug")
	if err != nil {
		log.Printf("Failed to initialize debug logging: %v", err)
	} else {
		defer f.Close()
	}

	// Start the TUI application
	p := tea.NewProgram(models.ChoiceModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
