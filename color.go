package main

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

// Colors
var (
	Colors = []*Color{Danger, Blue, Green, Red}

	Danger = &Color{
		Colorizer: color.New(color.FgWhite).Add(color.Bold).Add(color.BgRed).Sprint,
		ConfigKey: "Danger",
	}

	Blue = &Color{
		Colorizer: color.New(color.FgBlue).Sprint,
		ConfigKey: "Blue",
	}

	Green = &Color{
		Colorizer: color.New(color.FgGreen).Sprint,
		ConfigKey: "Green",
	}

	Red = &Color{
		Colorizer: color.New(color.FgRed).Add(color.Bold).Sprint,
		ConfigKey: "Red",
	}
)

// struct Color holds the configuration for each handled color.
type Color struct {
	// RE is a slice of regexps that, when matched, are colorized
	RE []*regexp.Regexp
	// Colorizer is a Sprint() method of a color from fatih/color
	Colorizer func(...any) string
	// ConfigKey is the key name (under Colorize) in the conf file
	// where REs are listed
	ConfigKey string
}

// c.addRE turns a slice of strings into regexps and attaches them
// to c.RE. Strings are surrounded by parentheses before being
// compiled into regexps.
func (c *Color) addRE(REs []string) {
	for _, re := range REs {
		re = fmt.Sprintf("(%s)", re)
		c.RE = append(c.RE, regexp.MustCompile(re))
	}
}

// c.colorize applies c.Colorizer() for each regexp in c.RE to
// the provided string, and returns a string fully colorized for
// the particular color.
func (c *Color) colorize(s string) string {
	for _, re := range c.RE {
		s = re.ReplaceAllString(s, c.Colorizer("$1"))
	}

	return s
}

// loadColors adds the provided regexps into each color's Color struct.
func loadColors(conf *muxytailConf) {
	for _, clr := range Colors {
		clr.addRE(conf.Colorize[clr.ConfigKey])
	}
}
