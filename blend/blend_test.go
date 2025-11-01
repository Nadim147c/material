package blend

import (
	"testing"

	"github.com/Nadim147c/material/v2/color"
)

func TestHarmonize(t *testing.T) {
	RED := color.ARGB(0xffff0000)
	BLUE := color.ARGB(0xff0000ff)
	GREEN := color.ARGB(0xff00ff00)
	YELLOW := color.ARGB(0xffffff00)

	tests := []struct {
		name     string
		design   color.ARGB
		source   color.ARGB
		expected color.ARGB
	}{
		{"redToBlue", RED, BLUE, 0xfffb0057},
		{"redToGreen", RED, GREEN, 0xffd85600},
		{"redToYellow", RED, YELLOW, 0xffd85600},
		{"blueToGreen", BLUE, GREEN, 0xff0047a3},
		{"blueToRed", BLUE, RED, 0xff5700dc},
		{"blueToYellow", BLUE, YELLOW, 0xff0047a3},
		{"greenToBlue", GREEN, BLUE, 0xff00fc94},
		{"greenToRed", GREEN, RED, 0xffb1f000},
		{"greenToYellow", GREEN, YELLOW, 0xffb1f000},
		{"yellowToBlue", YELLOW, BLUE, 0xffebffba},
		{"yellowToGreen", YELLOW, GREEN, 0xffebffba},
		{"yellowToRed", YELLOW, RED, 0xfffff6e3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Harmonize(tt.design, tt.source)
			if got != tt.expected {
				t.Errorf(
					"Harmonize(%s, %s) = %s; want %s",
					tt.design.String(),
					tt.source.String(),
					got.String(),
					tt.expected.String(),
				)
			}
		})
	}
}
