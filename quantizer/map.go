package quantizer

import (
	"slices"

	"github.com/Nadim147c/material/v2/color"
)

// QuantizeMap takes a slice of []color.Color and returns Quantized
func QuantizeMap(input []color.ARGB) QuantizedMap {
	colors := make(QuantizedMap)
	for pixel := range slices.Values(input) {
		alpha := pixel.Alpha()
		if alpha < 0xFF {
			continue
		}
		colors[pixel]++
	}

	return colors
}
