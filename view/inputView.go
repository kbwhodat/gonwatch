package view

import (
	"os"
	"golang.org/x/term"
)


func InputView(m string) string {

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}

	prompt := m + "\n\n\n(left/h: go back • esc: quit)"

	return  centerText(prompt, w, h)

}
