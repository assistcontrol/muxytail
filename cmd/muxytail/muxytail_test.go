package muxytail

import (
	"testing"

	"github.com/assistcontrol/muxytail/formatter"
)

type mockFormatter struct {
	output string
	ok     bool
}

func (m *mockFormatter) Format(input string) (string, bool) {
	return m.output, m.ok
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		formatters formatter.List
		expected   string
	}{
		{
			name:       "No formatters",
			input:      "test log line",
			formatters: formatter.List{},
			expected:   "test log line",
		},
		{
			name:  "First formatter succeeds",
			input: "test log line",
			formatters: formatter.List{
				&mockFormatter{output: "formatted by first", ok: true},
				&mockFormatter{output: "formatted by second", ok: true},
			},
			expected: "formatted by first",
		},
		{
			name:  "Second formatter succeeds",
			input: "test log line",
			formatters: formatter.List{
				&mockFormatter{output: "unchanged", ok: false},
				&mockFormatter{output: "formatted by second", ok: true},
			},
			expected: "formatted by second",
		},
		{
			name:  "No formatter succeeds",
			input: "test log line",
			formatters: formatter.List{
				&mockFormatter{output: "unchanged", ok: false},
				&mockFormatter{output: "unchanged", ok: false},
			},
			expected: "test log line",
		},
		{
			name:  "Empty input string",
			input: "",
			formatters: formatter.List{
				&mockFormatter{output: "formatted empty", ok: true},
			},
			expected: "formatted empty",
		},
		{
			name:  "Multiple formatters, none succeed",
			input: "another test log line",
			formatters: formatter.List{
				&mockFormatter{output: "unchanged", ok: false},
				&mockFormatter{output: "unchanged", ok: false},
				&mockFormatter{output: "unchanged", ok: false},
			},
			expected: "another test log line",
		},
		{
			name:  "Multiple formatters, last one succeeds",
			input: "yet another test log line",
			formatters: formatter.List{
				&mockFormatter{output: "unchanged", ok: false},
				&mockFormatter{output: "unchanged", ok: false},
				&mockFormatter{output: "formatted by last", ok: true},
			},
			expected: "formatted by last",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := format(tt.input, tt.formatters)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
