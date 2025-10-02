package view

import (
	"strings"
	"golang.org/x/term"
	"os"
)

func SelectView(selectedIndex int) string {
	// Get terminal width directly within the view function
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Handle the error or fallback to a default width
		w = 80 // Fallback to a default width if there's an error
		h = 24
	}

	s := strings.Builder{}
	s.WriteString("What do you want to watch?\n\n")

	for i := 0; i < len(ChoiceList); i++ {
		if selectedIndex == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(ChoiceList[i].choice)
		s.WriteString("\n")
	}

	// Use centerText function to center the content based on the current terminal width
	return centerText(s.String(), w, h)
}
