package color

import (
	"fmt"
	"regexp"

	termcolor "github.com/fatih/color"
)

var Colors = map[string]*Color{
	"BoldRed": {
		Colorizer: termcolor.New(termcolor.FgRed).Add(termcolor.Bold).Sprint,
		ConfigKey: "BoldRed",
	},

	"Danger": {
		Colorizer: termcolor.New(termcolor.FgWhite).Add(termcolor.Bold).Add(termcolor.BgRed).Sprint,
		ConfigKey: "Danger",
	},

	"Blue": {
		Colorizer: termcolor.New(termcolor.FgBlue).Sprint,
		ConfigKey: "Blue",
	},

	"Green": {
		Colorizer: termcolor.New(termcolor.FgGreen).Sprint,
		ConfigKey: "Green",
	},

	"Yellow": {
		Colorizer: termcolor.New(termcolor.FgYellow).Sprint,
		ConfigKey: "Yellow",
	},

	"Red": {
		Colorizer: termcolor.New(termcolor.FgRed).Sprint,
		ConfigKey: "Red",
	},
}

// struct Color holds the configuration for each handled termcolor.
type Color struct {
	// RE is a slice of regexps that, when matched, are colorized
	RE []*regexp.Regexp
	// Colorizer is a Sprint() method of a color from fatih/color
	Colorizer func(...any) string
	// ConfigKey is the key name (under Colorize) in the conf file
	// where REs are listed
	ConfigKey string
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
