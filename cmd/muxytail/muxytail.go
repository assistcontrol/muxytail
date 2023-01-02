package muxytail

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/assistcontrol/muxytail/caddy"
	"github.com/assistcontrol/muxytail/color"

	"github.com/nxadm/tail"
	"golang.org/x/term"
)

// Default path to the config file. Can be overridden with --config.
const defaultConfigFile = "/usr/local/etc/muxytail.yaml"

var tailConfig = tail.Config{
	Location: &tail.SeekInfo{
		Whence: io.SeekEnd,
	},
	MustExist: true,
	Follow:    true,
	ReOpen:    true,
}

// Run is essentially main(), whereas the real main() is a stub.
func Run() {
	configFile := flag.String("config", defaultConfigFile, "config file location")
	flag.Parse()

	config := loadConfig(*configFile)
	loadColors(config)
	caddy.LoadColors(config.Caddy)

	// Watch for Enter
	separatorChannel := make(chan string)
	go watchStdin(separatorChannel)

	// Each file sends log lines to logChannel
	logChannel := make(chan string)
	for _, path := range config.Files {
		go watchFile(path, logChannel)
	}

	for {
		select {
		case s := <-logChannel:
			fmt.Println(s)
		case s := <-separatorChannel:
			fmt.Println(s)
		}
	}
}

// watchFile tails a given path and sends all new lines up the provided
// channel after formatting.
func watchFile(path string, c chan<- string) {
	t, err := tail.TailFile(path, tailConfig)
	if err != nil {
		log.Fatalln(err)
	}

	for line := range t.Lines {
		go func(s string) {
			c <- format(s)
		}(line.Text)
	}
}

// watchStdin sends a separator up the provided channel whenever Enter
// is pressed. It uses password mode so that the input doesn't echo.
func watchStdin(c chan<- string) {
	for {
		if _, err := term.ReadPassword(int(os.Stdin.Fd())); err != nil {
			log.Fatalln("ReadPassword:", err)
		}

		c <- separator()
	}
}

// format is the master function for log line manipulation. It returns
// its string argument converted, formatted, and colorized.
func format(s string) string {
	if s, converted := caddy.Convert(s); converted {
		return s
	}

	// Not a caddy string. Apply regexp colorizers.
	s = color.Colorize(s)

	return s
}

// separator returns a string containing a horizontal line the width of the
// terminal.
func separator() string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalln("term.GetSize:", err)
	}

	return color.ColorizeSeparator(strings.Repeat("─", width))
}
