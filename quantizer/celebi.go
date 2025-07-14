package quantizer

func QuantizeCelebi(input pixels, maxColor int) QuantizedMap {
	wu := QuantizeWu(input, maxColor*5)
	colors := make(pixelsLab, len(wu))
	for i, c := range wu {
		colors[i] = c.ToLab()
	}
	return QuantizeWsMeans(input, colors, maxColor)
}
