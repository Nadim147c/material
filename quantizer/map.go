package quantizer

import "slices"

// QuantizeMap takes a slice of []color.Color and returns Quantized
func QuantizeMap(input pixels) QuantizedMap {
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
