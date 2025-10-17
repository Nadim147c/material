package material

import (
	"image/color"

	icolor "github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

// GenerateFromColors generate theme from a slice of Color interface
func GenerateFromColors(
	colors []color.Color,
	variant dynamic.Variant,
	dark bool,
	constrast float64,
	platform dynamic.Platform,
	version dynamic.Version,
) (Colors, error) {
	argbs := make([]icolor.ARGB, len(colors))
	for i, c := range colors {
		argbs[i] = icolor.ARGBFromInterface(c)
	}
	return GenerateFromPixels(
		argbs,
		variant,
		dark,
		constrast,
		platform,
		version,
	)
}
