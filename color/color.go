package color

import (
	"fmt"
	"strings"

	termcolor "github.com/gookit/color"
)

// Colorizer is the signature for a colorizer function.
type Colorizer func(...any) string

// GenerateColorizer returns a function that colorizes a string argument.
func GenerateColorizer(s string) Colorizer {
	// If passed an empty string, just use fmt.Sprint to get uncolored output.
	if s == "" {
		return fmt.Sprint
	}

	// Split into [FG, BG]
	parts := strings.Split(s, "|")

	tclr := termcolor.HEXStyle(parts[0])
	if len(parts) == 2 {
		tclr = termcolor.HEXStyle(parts[0], parts[1])
	}

	return tclr.Sprint
}
