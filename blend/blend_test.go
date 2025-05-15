package blend

import (
	"testing"

	"github.com/Nadim147c/goyou/color"
)

func TestHarmonize(t *testing.T) {
	tests := []struct {
		name     string
		from     color.Color
		to       color.Color
		expected color.Color
	}{
		{"RedToBlue", 0xffff0000, 0xff0000ff, 0xffFB0057},
		{"RedToGreen", 0xffff0000, 0xff00ff00, 0xffD85600},
		{"RedToYellow", 0xffff0000, 0xffffff00, 0xffD85600},
		{"BlueToGreen", 0xff0000ff, 0xff00ff00, 0xff0047A3},
		{"BlueToRed", 0xff0000ff, 0xffff0000, 0xff5700DC},
		{"BlueToYellow", 0xff0000ff, 0xffffff00, 0xff0047A3},
		{"GreenToBlue", 0xff00ff00, 0xff0000ff, 0xff00FC94},
		{"GreenToRed", 0xff00ff00, 0xffff0000, 0xffB1F000},
		{"GreenToYellow", 0xff00ff00, 0xffffff00, 0xffB1F000},
		{"YellowToBlue", 0xffffff00, 0xff0000ff, 0xffEBFFBA},
		{"YellowToGreen", 0xffffff00, 0xff00ff00, 0xffEBFFBA},
		{"YellowToRed", 0xffffff00, 0xffff0000, 0xffFFF6E3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Harmonize(tc.from, tc.to)
			if got != tc.expected {
				t.Errorf("Harmonize(%#x, %#x) = %#x; want %#x", uint32(tc.from), uint32(tc.to), uint32(got), uint32(tc.expected))
			}
		})
	}
}
