package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// MuxytailConf is the root data structure holding configuration.
type MuxytailConf struct {
	Files     []string        `yaml:"files"`
	Colorize  REConfig        `yaml:"colorize"`
	Separator SeparatorConfig `yaml:"separator"`
	Caddy     CaddyConfig     `yaml:"caddy"`
}

// struct ColorConfig is the caddy-specific color struct for the
// caddy formatter.
type CaddyConfig struct {
	Bracket     string `yaml:"bracket"`
	Host        string `yaml:"host"`
	StatusOK    string `yaml:"status_ok"`
	StatusError string `yaml:"status_error"`
	StatusOther string `yaml:"status_other"`
	URL         string `yaml:"url"`
}

// struct REConfig is a table of color strings that
// map to a slice of regexps. The matches of each
// regexps get colorized according to the string key.
type REConfig map[string][]string

// struct SeparatorConfig is the separator-specific configuration.
type SeparatorConfig struct {
	Color string `yaml:"color"`
}

// Load reads the config file and parses the YAML into a MuxytailConf.
func Load(path string) *MuxytailConf {
	// Read in the config file
	confBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("ReadFile:", err)
	}

	c, err := unmarshal(confBytes)
	if err != nil {
		log.Fatalln("unmarshal:", err)
	}

	return c
}

func unmarshal(data []byte) (*MuxytailConf, error) {
	var config MuxytailConf
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
