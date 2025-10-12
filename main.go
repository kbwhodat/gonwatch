package main

import (
	"log"
	"os"
	"fmt"
	"gonwatch/models"
	tea "github.com/charmbracelet/bubbletea"

)

func main() {

	f, err := tea.LogToFile("/tmp/gonwatch/debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
		defer f.Close()
	}

	p := tea.NewProgram(models.ChoiceModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
