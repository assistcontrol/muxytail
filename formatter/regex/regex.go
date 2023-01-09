package regex

import (
	"fmt"
	"regexp"

	"github.com/assistcontrol/muxytail/color"
	"github.com/assistcontrol/muxytail/config"
)

// regexList is a slice of known reColor structs. It is created
// by regex.New(). New colors are added by .AddColor().
type colorList []*reColor

// struct reColor holds the configuration for a RE-based
// colorizer. reColor.RE is the slice of regexps that, when
// matched, are colorized by the reColor.Colorizer function.
type reColor struct {
	RE        []*regexp.Regexp
	Colorizer color.Colorizer
}

// rc.addREs turns a slice of strings into regexps and appends
// them to the RE field. The string arguments are converted to
// regexps by surrounding them in parentheses.
func (rc *reColor) addREs(REs []string) {
	for _, re := range REs {
		re = fmt.Sprintf("(%s)", re)
		rc.RE = append(rc.RE, regexp.MustCompile(re))
	}
}

// res.Format is the main function that handles colorization. It
// applies each regeistered color to all matches of each RE.
func (colors colorList) Format(in string) (string, bool) {
	out := in

	for _, clr := range colors {
		for _, re := range clr.RE {
			out = re.ReplaceAllString(out, clr.Colorizer("$1"))
		}
	}

	return out, in != out
}

// New creates a new regex formatter. It takes a config.REConfig
// and returns a formatter that colorizes matches of each regexp.
func New(conf config.REConfig) colorList {
	colors := make(colorList, 0, len(conf))

	for colorString, REs := range conf {
		rc := &reColor{
			Colorizer: color.GenerateColorizer(colorString),
		}
		rc.addREs(REs)

		colors = append(colors, rc)
	}

	return colors
}
