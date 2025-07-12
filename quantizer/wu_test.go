package quantizer

import (
	"embed"
	"image/jpeg"
	"slices"
	"testing"

	"github.com/Nadim147c/material/color"
)

//go:embed gophar.jpg
var gophar embed.FS

func TestQuantizeWu(t *testing.T) {
	file, err := gophar.Open("gophar.jpg")
	if err != nil {
		panic(err)
	}

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

	result := QuantizeWu(pixels, 2)
	if len(result) == 0 {
		t.Fatal("QuantizeWu() returned no colors")
	}

	c1 := color.ARGBFromHexMust("#0A0D0E")
	c2 := color.ARGBFromHexMust("#79D0DB")

	if len(result) != 2 {
		t.Fatalf("Result: %v has unexpected number of color", result)
	}

	if !slices.Contains(result, c1) {
		t.Fatalf("result: %v doesn't contains %v", result, c1)
	}
	if !slices.Contains(result, c2) {
		t.Fatalf("result: %v doesn't contains %v", result, c1)
	}
}
