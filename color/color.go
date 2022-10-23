package color

import (
	"fmt"
	"log"
	"regexp"

	termcolor "github.com/fatih/color"
)

// map Colors is the list of known colors. After adding a color
// here, be sure to add a key of the same name in muxytail.conf.
var Colors = map[string]*Color{
	"Blue":    {Colorizer: termcolor.New(termcolor.FgBlue).Sprint},
	"Green":   {Colorizer: termcolor.New(termcolor.FgGreen).Sprint},
	"Yellow":  {Colorizer: termcolor.New(termcolor.FgYellow).Sprint},
	"Red":     {Colorizer: termcolor.New(termcolor.FgRed).Sprint},
	"BoldRed": {Colorizer: termcolor.New(termcolor.FgRed).Add(termcolor.Bold).Sprint},
	"Danger":  {Colorizer: termcolor.New(termcolor.FgWhite).Add(termcolor.Bold).Add(termcolor.BgRed).Sprint},
}

// struct Color holds the configuration for each handled termcolor.
type Color struct {
	// RE is a slice of regexps that, when matched, are colorized
	RE []*regexp.Regexp
	// Colorizer is a Sprint() method of a color from fatih/color
	Colorizer func(...any) string
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
