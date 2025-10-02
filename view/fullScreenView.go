package view


import (
	"os"
	"golang.org/x/term"
	"fmt"
	"strings"
	"github.com/acarl005/stripansi"
)

func FullscreenView(keywordStyle string, helpStyle string) string {

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}

	message := ""

	if strings.Compare(stripansi.Strip(keywordStyle), "Please wait, updating database with new media...") == 0 {
		message = fmt.Sprintf("\n\n\n\n%s\n\n\n%s", keywordStyle, helpStyle)

	} else {
		message = fmt.Sprintf("\n\n\n\n  You are now watching  %s\n\n\n%s", keywordStyle, helpStyle)
	}

	return centerText(message, w, h)
}
