package separator

import (
	"strings"
	"testing"

	"github.com/assistcontrol/muxytail/color"
	termcolor "github.com/gookit/color"
)

func TestSeparator_Display(t *testing.T) {
	tests := []struct {
		name      string
		colorizer color.Colorizer
		want      string
	}{
		{
			name:      "Default colorizer",
			colorizer: color.GenerateColorizer(""),
			want:      strings.Repeat(separatorChar, 80),
		},
		{
			name:      "Red colorizer",
			colorizer: color.GenerateColorizer("#ff0000"),
			want:      termcolor.HEXStyle("#ff0000").Sprint(strings.Repeat(separatorChar, 80)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sep := &Separator{
				Colorizer: tt.colorizer,
				Width:     80,
			}
			if got := sep.Display(); got != tt.want {
				t.Errorf("Separator.Display() = %v, want %v", got, tt.want)
			}
		})
	}
}
