package color

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	termcolor "github.com/gookit/color"
)

// map Colors is the list of known colors. Colors are added
// with AddColor(), and REs are added to that color with
// c.AddRE().
var Colors = make(map[string]*Color)

// struct Color holds the configuration for each handled termcolor.
type Color struct {
	// RE is a slice of regexps that, when matched, are colorized
	RE []*regexp.Regexp
	// Colorizer is a function (generally from another package)
	// that returns a colorized version of its string arguments.
	Colorizer func(...any) string
}

// AddColor takes a color string, creates a renderer for it,
// and adds it to the Colors map. It returns the Color so that
// c.AddRE() can be chained.
func AddColor(s string) *Color {
	parts := strings.Split(s, "|")

	tclr := termcolor.HEXStyle(parts[0])
	if len(parts) == 2 {
		tclr = termcolor.HEXStyle(parts[0], parts[1])
	}

	clr := &Color{
		Colorizer: tclr.Sprint,
	}

	Colors[s] = clr
	return clr
}

// ColorizeString returns the string argument in a specific color.
func ColorizeString(clr, s string) string {
	c, exists := Colors[clr]
	if !exists {
		log.Fatalln("Unknown color", clr)
	}

	return c.Colorizer(s)
}

// Colorize iterates over each Color and applies all of the regexps
// associated with that color to the supplied string.
func Colorize(s string) string {
	for _, clr := range Colors {
		s = clr.Colorize(s)
	}

	return s
}

// c.AddRE turns a slice of strings into regexps and attaches them
// to c.RE. Strings are surrounded by parentheses before being
// compiled into regexps.
func (c *Color) AddRE(REs []string) {
	for _, re := range REs {
		re = fmt.Sprintf("(%s)", re)
		c.RE = append(c.RE, regexp.MustCompile(re))
	}
}

// c.Colorize applies c.Colorizer() for each regexp in c.RE to
// the provided string, and returns a string fully colorized for
// the particular termcolor.
func (c *Color) Colorize(s string) string {
	for _, re := range c.RE {
		s = re.ReplaceAllString(s, c.Colorizer("$1"))
	}

	return s
}
