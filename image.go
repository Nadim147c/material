package material

import (
	"image"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

// GenerateFromImage colors from an image.Image
func GenerateFromImage(
	img image.Image,
	variant dynamic.Variant,
	dark bool,
	constrast float64,
	platform dynamic.Platform,
	version dynamic.Version,
) (Colors, error) {
	bounds := img.Bounds()
	pixels := make([]color.ARGB, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			argb := color.ARGBFromInterface(c)
			pixels = append(pixels, argb)
		}
	}

	return GenerateFromPixels(pixels, variant, dark, constrast, platform, version)
}
