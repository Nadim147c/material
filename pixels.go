package material

import (
	"errors"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/quantizer"
	"github.com/Nadim147c/material/score"
)

// ErrNoColorFound indecates no color found in array
var ErrNoColorFound = errors.New("no color found")

// Colors is key and color
type Colors = map[string]color.ARGB

// GenerateFromPixels generate theme from a slice of ARGB colors
func GenerateFromPixels(
	pixels []color.ARGB,
	variant dynamic.Variant,
	dark bool,
	constrast float64,
	platform dynamic.Platform,
	version dynamic.Version,
) (Colors, error) {
	quantized := quantizer.QuantizeCelebi(pixels, 5)
	if len(quantized) == 0 {
		return nil, ErrNoColorFound
	}

	scored := score.Score(quantized)
	if len(scored) == 0 {
		return nil, ErrNoColorFound
	}

	scheme := dynamic.NewDynamicScheme(
		scored[0].ToHct(),
		variant,
		constrast,
		dark,
		platform,
		version,
	)
	dcs := scheme.ToColorMap()

	colorMap := Colors{}
	for key, value := range dcs {
		if value != nil {
			colorMap[key] = value.GetArgb(scheme)
		}
	}
	return colorMap, nil
}
