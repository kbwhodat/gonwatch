package main

import (
	"log"
	"os"
	"fmt"
	"gonwatch/models"
	tea "github.com/charmbracelet/bubbletea"

)

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "debug" {
			file, _ := os.Create("/tmp/debug.log")
			defer file.Close()
			f, err := tea.LogToFile("/tmp/debug.log", "debug")
			if err != nil {
				fmt.Println("fatal:", err)
				defer f.Close()
			}
		}
	}

	p := tea.NewProgram(models.ChoiceModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
