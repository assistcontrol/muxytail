package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		fileData string
		expected *MuxytailConf
		wantErr  bool
	}{
		{
			name: "Valid config",
			fileData: `
files:
  - "/var/log/syslog"
colorize:
  "error": ["ERROR", "FATAL"]
separator:
  color: "#FF5733"
caddy:
  bracket: "[]"
  host: "localhost"
  status_ok: "green"
  status_error: "red"
  status_other: "yellow"
  url: "http://example.com"
`,
			expected: &MuxytailConf{
				Files: []string{"/var/log/syslog"},
				Colorize: REConfig{
					"error": {"ERROR", "FATAL"},
				},
				Separator: SeparatorConfig{
					Color: "#FF5733",
				},
				Caddy: CaddyConfig{
					Bracket:     "[]",
					Host:        "localhost",
					StatusOK:    "green",
					StatusError: "red",
					StatusOther: "yellow",
					URL:         "http://example.com",
				},
			},
			wantErr: false,
		},
		{
			name:     "Invalid YAML",
			fileData: `invalid_yaml: [`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load the config
			got, err := unmarshal([]byte(tt.fileData))
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (got == nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", got == nil, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalMuxytailConf(got, tt.expected) {
				t.Errorf("Load() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// Helper function to compare two MuxytailConf structs
func equalMuxytailConf(a, b *MuxytailConf) bool {
	if a == nil || b == nil {
		return a == b
	}
	return yamlEqual(a, b)
}

// Helper function to compare two structs using YAML marshaling
func yamlEqual(a, b interface{}) bool {
	ay, _ := yaml.Marshal(a)
	by, _ := yaml.Marshal(b)
	return string(ay) == string(by)
}
