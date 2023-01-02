package muxytail

import (
	"log"
	"os"

	"github.com/assistcontrol/muxytail/caddy"
	"github.com/assistcontrol/muxytail/color"
	"gopkg.in/yaml.v3"
)

// Config
type muxytailConf struct {
	Files    []string
	Colorize map[string][]string
	Caddy    caddy.ColorConfig
}

// loadConfig reads the config file and parses the YAML. It
// returns a muxytailConf populated with the YAML data.
func loadConfig(path string) *muxytailConf {
	// Read in the config file
	confBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("ReadFile:", err)
	}

	// Unmarshal the config into a muxytailConf
	var config muxytailConf
	err = yaml.Unmarshal(confBytes, &config)
	if err != nil {
		log.Fatalln("YAML parse:", err)
	}

	return &config
}

// loadColors adds the provided regexps into each color's Color struct.
func loadColors(conf *muxytailConf) {
	for clr, REs := range conf.Colorize {
		color.AddColor(clr).AddRE(REs)
	}
}
