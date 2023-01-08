package separator

import (
	"log"
	"os"
	"strings"

	"github.com/assistcontrol/muxytail/color"
	"github.com/assistcontrol/muxytail/config"
	"golang.org/x/term"
)

// separatorChar is repeated across the terminal.
const separatorChar = "─"

// struct separator holds the separator config.
type Separator struct {
	Colorizer color.Colorizer
}

// s.Display repeats s.Char across the terminal, and applies
// s.Colorizer to it.
func (s *Separator) Display() string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalln("term.GetSize:", err)
	}

	return s.Colorizer(strings.Repeat(separatorChar, width))
}

// New returns a separator struct that is capable of displaying
// a colorized line across the terminal.
func New(conf config.SeparatorConfig) *Separator {
	sep := &Separator{
		Colorizer: color.GenerateColorizer(conf.Color),
	}

	return sep
}
