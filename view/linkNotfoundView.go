package view


import (
	"os"
	"golang.org/x/term"
	"fmt"
	// "strings"
	// "github.com/acarl005/stripansi"
)

func LinkNotFoundView(keywordStyle string, helpStyle string) string {

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}

	message := fmt.Sprintf("\n\n\n\n  No available links are available for %s\n\n\n%s", keywordStyle, helpStyle)

	return centerText(message, w, h)
}
