package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/assistcontrol/muxytail/caddy"
	"github.com/assistcontrol/muxytail/color"

	"github.com/nxadm/tail"
	"golang.org/x/term"
)

// Default path to the config file. Can be overridden with --config.
const defaultConfigFile = "/usr/local/etc/muxytail.yaml"

// Config
type muxytailConf struct {
	Files    []string
	Colorize map[string][]string
}

var tailConfig = tail.Config{
	Location: &tail.SeekInfo{
		Whence: io.SeekEnd,
	},
	MustExist: true,
	Follow:    true,
	ReOpen:    true,
}

func main() {
	configFile := flag.String("config", defaultConfigFile, "config file location")
	flag.Parse()

	config := loadConfig(*configFile)
	loadColors(config)

	signalChannel := make(chan os.Signal, 2)
	stdinChannel := make(chan bool)
	textChannel := make(chan string)

	// Trap SIGQUIT to ignore ^\
	signal.Notify(signalChannel, syscall.SIGQUIT)

	// Watch for Enter
	go watchStdin(stdinChannel)

	// Each file sends log lines to textChannel
	for _, path := range config.Files {
		go watchFile(path, textChannel)
	}

	for {
		select {
		case line := <-textChannel:
			fmt.Println(line)
		case <-stdinChannel:
			fmt.Println(separator())
		case <-signalChannel:
			// Ignore
		}
	}
}

// watchFile tails a given path and sends all new lines up the provided
// channel.
func watchFile(path string, c chan<- string) {
	t, err := tail.TailFile(path, tailConfig)
	if err != nil {
		log.Fatalln(err)
	}

	for line := range t.Lines {
		go func(txt string) {
			c <- Format(txt)
		}(line.Text)
	}
}

// watchStdin sends a message up the provided channel whenever Enter
// is pressed. It uses password mode so that the input doesn't echo.
func watchStdin(c chan bool) {
	for {
		if _, err := term.ReadPassword(int(os.Stdin.Fd())); err != nil {
			log.Fatalln("ReadPassword:", err)
		}

		c <- true
	}
}

// Format is the master function for log line manipulation. It returns
// its string argument converted, formatted, and colorized.
func Format(s string) string {
	if s, converted := caddy.Convert(s); converted {
		return s
	}

	// Not a caddy string. Iterate over each Color and apply all of
	// the regexps associated with that color.
	for _, clr := range color.Colors {
		s = clr.Colorize(s)
	}

	return s
}

// separator returns a string containing a horizontal line the width of the
// terminal.
func separator() string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalln("term.GetSize:", err)
	}

	return color.Colors["BoldRed"].Colorizer(strings.Repeat("─", width))
}