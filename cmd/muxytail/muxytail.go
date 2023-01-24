package muxytail

import (
	"flag"
	"fmt"
	"io"
	"log"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/assistcontrol/muxytail/config"
	"github.com/assistcontrol/muxytail/formatter"
	"github.com/assistcontrol/muxytail/formatter/caddy"
	"github.com/assistcontrol/muxytail/formatter/regex"
	"github.com/assistcontrol/muxytail/separator"

	"github.com/nxadm/tail"
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

	conf := config.Load(*configFile)
	formatters := formatter.List{
		caddy.New(conf.Caddy),
		regex.New(conf.Colorize),
	}

	// Watch for Enter
	separatorChannel := make(chan string)
	exitChannel := make(chan bool)
	go watchStdin(separatorChannel, exitChannel, separator.New(conf.Separator))

	// Each file sends log lines to logChannel
	logChannel := make(chan string)
	for _, path := range conf.Files {
		go watchFile(path, formatters, logChannel)
	}

	for {
		select {
		case s := <-logChannel:
			fmt.Println(s)
		case s := <-separatorChannel:
			fmt.Println(s)
		case <-exitChannel:
			return
		}
	}
}

// watchFile tails a given path and sends all new lines up the provided
// channel after formatting.
func watchFile(path string, formatters formatter.List, c chan<- string) {
	t, err := tail.TailFile(path, tailConfig)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = t.Stop(); err != nil {
			log.Fatalln(err)
		}
	}()

	for line := range t.Lines {
		go func(s string) {
			c <- format(s, formatters)
		}(line.Text)
	}
}

// watchStdin listens for keyboard events. On Enter, a separator is
// passed up the provided channel. On Esc or ^C, the program exits.
func watchStdin(sepCh chan<- string, exitCh chan<- bool, sep *separator.Separator) {
	onKey := func(key keys.Key) (bool, error) {
		switch key.Code {
		case keys.Enter:
			go func() {
				sepCh <- sep.Display()
			}()
		case keys.CtrlC:
			exitCh <- true
		case keys.RuneKey:
			if key.String() == "q" {
				exitCh <- true
			}
		}

		return false, nil // Continue listening
	}

	if err := keyboard.Listen(onKey); err != nil {
		log.Fatalln("keyboard.Listen:", err)
	}
}

// format applies its string argument sequentially to each formatter.
// Formatting stops after the first formatter that indicates
// successful formatting.
func format(in string, formatters formatter.List) string {
	for _, f := range formatters {
		if out, ok := f.Format(in); ok {
			return out
		}
	}

	// No formatter was successful
	return in
}
