package quantizer

import (
	"fmt"
	"image/jpeg"
	"os"
	"slices"
	"testing"

	"github.com/Nadim147c/goyou/color"
)

func TestQuantizeWu(t *testing.T) {
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
	var pixels []color.Color
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			goCol := img.At(x, y)
			pixels = append(pixels, color.FromGoColor(goCol))
		}
	}

	for i := range 10 {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			result := QuantizeWu(pixels, 3)
			if len(result) == 0 {
				t.Fatal("QuantizeWu() returned no colors")
			}

			c1 := color.Color(0xFF2E3B3C)
			c2 := color.Color(0xFF82D4DE)

			if len(result) != 2 {
				t.Fatalf("Result: %x has unexpected number of color", result)
			}

			if !slices.Contains(result, c1) {
				t.Fatalf("result: %x doesn't contains %x", result, c1)
			}
			if !slices.Contains(result, c2) {
				t.Fatalf("result: %x doesn't contains %x", result, c1)
			}
		})
	}
}
