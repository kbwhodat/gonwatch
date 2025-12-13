package view

import (
	"strings"
	"golang.org/x/term"
	"os"
)

func SelectView(selectedIndex int) string {
	// get terminal width directly within the view function
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	// fallback option
	if err != nil {
		w = 80
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

	return centerText(s.String(), w, h)
}

func TrendingSelectView(selectedIndex int) string {
	// get terminal width directly within the view function
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	// fallback option
	if err != nil {
		w = 80
		h = 24
	}

	ss := strings.Builder{}
	ss.WriteString("Select trending content to view.\n\n")

	for i := 0; i < len(TrendingChoiceList); i++ {
		if selectedIndex == i {
			ss.WriteString("(•) ")
		} else {
			ss.WriteString("( ) ")
		}
		ss.WriteString(TrendingChoiceList[i].choice)
		ss.WriteString("\n")
	}

	return centerText(ss.String(), w, h)
}
