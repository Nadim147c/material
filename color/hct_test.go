package color

import (
	"math"
	"testing"
)

func TestHct(t *testing.T) {
	originalLstar := 60.0
	originalHue := 180.0
	originalChroma := 40.0

	color := solveToARGB(originalHue, originalChroma, originalLstar)
	lstar2 := color.LStar()

	if math.Abs(lstar2-originalLstar) > 0.1 {
		t.Errorf("L* round-trip mismatch: got %.2f, want %.2f %s", lstar2, originalLstar, color.String())
	}
}

func TestHctRoundTrip(t *testing.T) {
	for _, tt := range ColorTestCases {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.ARGB.ToHct().ToARGB(); got != tt.ARGB {
				t.Errorf("Color(%s) Round Trip = %v, want %v", tt.ARGB.HexRGBA(), got.String(), tt.ARGB.String())
			}
		})
	}
}
