package quantizer

import (
	"fmt"
	"image/jpeg"
	"os"
	"slices"
	"testing"

	"github.com/Nadim147c/material/color"
)

func TestQuantizeWu(t *testing.T) {
	// Load the test image
	file, err := os.Open("/home/ephemeral/Pictures/Wallpapers/a23.jpeg")
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
			goColor := img.At(x, y)

			pixels = append(pixels, color.ARGBFromInterface(goColor))
		}
	}

	for i := range 10 {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result := QuantizeWu(pixels, 2)
			if len(result) == 0 {
				t.Fatal("QuantizeWu() returned no colors")
			}

			c1 := color.ARGB(0xFF201E30)
			c2 := color.ARGB(0xFF8A889C)

			if len(result) != 2 {
				t.Fatalf("Result: %v has unexpected number of color", result)
			}

			if !slices.Contains(result, c1) {
				t.Fatalf("result: %v doesn't contains %v", result, c1)
			}
			if !slices.Contains(result, c2) {
				t.Fatalf("result: %v doesn't contains %v", result, c1)
			}
		})
	}
}
