package caddy

import "testing"

func Test_caddyUserAgent_String(t *testing.T) {
	tests := []struct {
		name string
		cUA  caddyUserAgent
		want string
	}{
		{
			name: "Empty User Agent",
			cUA:  caddyUserAgent{},
			want: "-",
		},
		{
			name: "Single User Agent",
			cUA:  caddyUserAgent{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"},
			want: "Chrome 58 │ Windows 10",
		},
		{
			name: "Multiple User Agents",
			cUA:  caddyUserAgent{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15"},
			want: "Safari 14 │ macOS 10",
		},
		{
			name: "Bot User Agent",
			cUA:  caddyUserAgent{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
			want: "BOT: Googlebot 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cUA.String(); got != tt.want {
				t.Errorf("caddyUserAgent.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_major(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Single digit version",
			input: "1",
			want:  "1",
		},
		{
			name:  "Major and minor version",
			input: "1.2",
			want:  "1",
		},
		{
			name:  "Major, minor, and patch version",
			input: "1.2.3",
			want:  "1",
		},
		{
			name:  "Complex version string",
			input: "10.20.30",
			want:  "10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := major(tt.input); got != tt.want {
				t.Errorf("major() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_longest(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "Empty slice",
			input: []string{},
			want:  "",
		},
		{
			name:  "Single element",
			input: []string{"one"},
			want:  "one",
		},
		{
			name:  "Multiple elements with different lengths",
			input: []string{"short", "longer", "longest"},
			want:  "longest",
		},
		{
			name:  "Multiple elements with same length",
			input: []string{"same", "size", "test"},
			want:  "same",
		},
		{
			name:  "Mixed empty and non-empty strings",
			input: []string{"", "non-empty", ""},
			want:  "non-empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longest(tt.input); got != tt.want {
				t.Errorf("longest() = %v, want %v", got, tt.want)
			}
		})
	}
}
