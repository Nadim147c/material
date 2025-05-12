package quantizer

import (
	"testing"

	"github.com/Nadim147c/goyou/color"
)

func TestQuantizeWsMeans(t *testing.T) {
	input := []color.Color{
		color.FromRGB(255, 0, 0), // Red
		color.FromRGB(255, 0, 0), // Red
		color.FromRGB(255, 0, 0), // Red
		color.FromRGB(255, 0, 0), // Red
		color.FromRGB(255, 0, 0), // Red
		color.FromRGB(0, 255, 0), // Green
		color.FromRGB(0, 255, 0), // Green
		color.FromRGB(255, 0, 0), // Red
		color.FromRGB(0, 0, 255), // Blue
		color.FromRGB(0, 0, 255), // Blue
	}

	result := QuantizeWsMeans(input, nil, 3)

	for color, count := range result {
		t.Logf("Cluster #%s:  %d", color.HexRGB(), count)
	}
}
