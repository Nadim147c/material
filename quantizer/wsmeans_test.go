package quantizer

import (
	"image/jpeg"
	"os"
	"testing"

	"github.com/Nadim147c/material/color"
)

func TestQuantizeWsMeans(t *testing.T) {
	// Load the test image
	file, err := os.Open("./gophar.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("failed to decode image: %v", err)
	}

	// Convert image pixels to []color.Color (your type)
	var pixels []color.ARGB
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			goCol := img.At(x, y)
			pixels = append(pixels, color.ARGBFromInterface(goCol))
		}
	}

	result := QuantizeWsMeans(pixels, nil, 3)

	for color, count := range result {
		t.Logf("Cluster %s %s: %d", color.HexRGB(), color.AnsiBg("  "), count)
	}
}
