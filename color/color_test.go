package color

import (
	"testing"

	"github.com/gookit/color"
)

func TestGenerateColorizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		arg      string
		expected string
	}{
		{
			name:     "Empty string input",
			input:    "",
			arg:      "test",
			expected: "test",
		},
		{
			name:     "Foreground color only",
			input:    "#FF5733",
			arg:      "test",
			expected: color.HEXStyle("#FF5733").Sprint("test"),
		},
		{
			name:     "Foreground and background color",
			input:    "#FF5733|#333FFF",
			arg:      "test",
			expected: color.HEXStyle("#FF5733", "#333FFF").Sprint("test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorizer := GenerateColorizer(tt.input)
			result := colorizer(tt.arg)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
