package caddy

import (
	"strings"

	termcolor "github.com/gookit/color"
)

// struct ColorConfig is the caddy-specific color
// struct for the caddy formatter. Each member is
// a color string of format "#FG_RGB" or "#FG|#BG".
type ColorConfig struct {
	Bracket     string `json:"bracket"`
	Host        string `json:"host"`
	StatusOK    string `json:"status_ok"`
	StatusError string `json:"status_error"`
	StatusOther string `json:"status_other"`
	URL         string `json:"url"`
}

// colorize holds the colorizer functions for the
// formatter. It is populated by LoadColors.
var colorize colorizerFunctions

// struct colorizerFunctions holds colorization
// functions for each colorizable thing.
type colorizerFunctions struct {
	Bracket                            colorizerFunc
	Host                               colorizerFunc
	StatusOK, StatusError, StatusOther colorizerFunc
	URL                                colorizerFunc
}

// colorizerFunc is a function (generally from another
// package) that colorizes its arguments.
type colorizerFunc func(...any) string

// LoadColors sets the colors for the caddy formatter.
// It is passed a ColorConfig, and every struct member
// is passed to generateColor().
func LoadColors(c ColorConfig) {
	colorize = colorizerFunctions{
		Bracket:     generateColor(c.Bracket),
		Host:        generateColor(c.Host),
		StatusOK:    generateColor(c.StatusOK),
		StatusError: generateColor(c.StatusError),
		StatusOther: generateColor(c.StatusOther),
		URL:         generateColor(c.URL),
	}
}

// generateColor returns a colorization function based
// on the string argument..
func generateColor(s string) colorizerFunc {
	parts := strings.Split(s, "|")

	tclr := termcolor.HEXStyle(parts[0])
	if len(parts) == 2 {
		tclr = termcolor.HEXStyle(parts[0], parts[1])
	}

	return tclr.Sprint
}
