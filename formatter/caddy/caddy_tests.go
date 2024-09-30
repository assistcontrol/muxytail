package caddy

import (
	"testing"

	"github.com/assistcontrol/muxytail/config"
	termcolor "github.com/gookit/color"
)

// Format tests
func TestFormat(t *testing.T) {
	conf := config.CaddyConfig{
		Bracket:     "#FFFFFF",
		Host:        "#FF5733",
		StatusOK:    "#00FF00",
		StatusError: "#FF0000",
		StatusOther: "#FFFF00",
		URL:         "#0000FF",
	}

	clr := New(conf)

	tests := []struct {
		name     string
		input    string
		expected string
		success  bool
	}{
		{
			name:     "Valid log entry",
			input:    `{"Req":{"Remote":"127.0.0.1","Method":"GET","Proto":"HTTP/1.1","Headers":{"Referer":"http://example.com","UA":"Mozilla/5.0"}},"TS":"2023-10-01T12:00:00Z","Status":200,"URL":"/index.html"}`,
			expected: "127.0.0.1 [2023-10-01T12:00:00Z] /index.html (200) GET HTTP/1.1 http://example.com [Mozilla/5.0]",
			success:  true,
		},
		{
			name:     "Invalid log entry",
			input:    `{"Req":{"Remote":"127.0.0.1","Method":"GET","Proto":"HTTP/1.1","Headers":{"Referer":"http://example.com","UA":"Mozilla/5.0"}},"TS":"2023-10-01T12:00:00Z","Status":"OK","URL":"/index.html"}`,
			expected: `{"Req":{"Remote":"127.0.0.1","Method":"GET","Proto":"HTTP/1.1","Headers":{"Referer":"http://example.com","UA":"Mozilla/5.0"}},"TS":"2023-10-01T12:00:00Z","Status":"OK","URL":"/index.html"}`,
			success:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, success := clr.Format(tt.input)
			if success != tt.success {
				t.Errorf("expected success %v, got %v", tt.success, success)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// formatLog tests
func TestFormatLog(t *testing.T) {
	tests := []struct {
		name     string
		cLog     *caddyLog
		colorCfg config.CaddyConfig
		expected string
	}{
		{
			name: "Basic log entry",
			cLog: &caddyLog{
				Req: caddyReq{
					Remote: "127.0.0.1",
					Method: "GET",
					Proto:  "HTTP/1.1",
					Headers: caddyHeaders{
						Referer: []string{"http://example.com"},
						UA:      []string{"Mozilla/5.0"},
					},
				},
				TS:     1696156800.0,
				Status: 200,
			},
			colorCfg: config.CaddyConfig{
				Bracket:     "#FFFFFF",
				Host:        "#FF5733",
				StatusOK:    "#00FF00",
				StatusError: "#FF0000",
				StatusOther: "#FFFF00",
				URL:         "#0000FF",
			},
			expected: "127.0.0.1 [2023-10-01T12:00:00Z] http://example.com (200) GET HTTP/1.1 http://example.com [Mozilla/5.0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clr := New(tt.colorCfg)
			result := clr.formatLog(tt.cLog)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// colorizeStatus tests
func TestColorizeStatus(t *testing.T) {
	conf := config.CaddyConfig{
		StatusOK:    "#00FF00",
		StatusError: "#FF0000",
		StatusOther: "#FFFF00",
	}
	clr := New(conf)

	tests := []struct {
		name     string
		status   caddyStatus
		expected string
	}{
		{
			name:     "Status OK",
			status:   200,
			expected: termcolor.HEXStyle("#00FF00").Sprint("200"),
		},
		{
			name:     "Status Error",
			status:   404,
			expected: termcolor.HEXStyle("#FF0000").Sprint("404"),
		},
		{
			name:     "Status Other",
			status:   302,
			expected: termcolor.HEXStyle("#FFFF00").Sprint("302"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clr.colorizeStatus(tt.status)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// parse tests
func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid JSON log entry",
			input:   `{"Req": {"Remote": "127.0.0.1", "Method": "GET", "Proto": "HTTP/1.1", "Headers": {"Referer": "http://example.com", "UA": "Mozilla/5.0"}}, "TS": "2023-10-01T12:00:00Z", "Status": 200}`,
			wantErr: false,
		},
		{
			name:    "Invalid JSON log entry",
			input:   `{"Req": {"Remote": "127.0.0.1", "Method": "GET", "Proto": "HTTP/1.1", "Headers": {"Referer": "http://example.com", "UA": "Mozilla/5.0"}, "TS": "2023-10-01T12:00:00Z", "Status": 200}`,
			wantErr: true,
		},
		{
			name:    "Empty JSON log entry",
			input:   `{}`,
			wantErr: false,
		},
		{
			name:    "Malformed JSON log entry",
			input:   `{"Req": {"Remote": "127.0.0.1", "Method": "GET", "Proto": "HTTP/1.1", "Headers": {"Referer": "http://example.com", "UA": "Mozilla/5.0"}, "TS": "2023-10-01T12:00:00Z", "Status": 200`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
