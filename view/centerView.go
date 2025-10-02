package view

import (
    "strings"
)

// Center text based on the terminal width
func centerText(s string, width, height int) string {
    lines := strings.Split(s, "\n")
    maxLineLen := 0
    for _, line := range lines {
        if len(line) > maxLineLen {
            maxLineLen = len(line)
        }
    }

    // Calculate horizontal padding
    horizontalPad := (width - maxLineLen) / 2
    if horizontalPad < 0 {
        horizontalPad = 0
    }

    // Calculate vertical padding (e.g., 1/3 of the terminal height)
    verticalPad := height / 3
    if verticalPad < 0 {
        verticalPad = 0
    }

    centered := strings.Builder{}
    // Add vertical padding at the top
    for i := 0; i < verticalPad; i++ {
        centered.WriteString("\n")
    }

    padStr := strings.Repeat(" ", horizontalPad)
    for _, line := range lines {
        // Add horizontal padding to each line
        centered.WriteString(padStr + line + "\n")
    }

    return centered.String()
}
