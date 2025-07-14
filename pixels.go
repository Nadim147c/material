package material

import (
	"errors"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/palettes"
	"github.com/Nadim147c/material/quantizer"
	"github.com/Nadim147c/material/score"
)

var ErrNoColorFound = errors.New("no color found")

// Colors is key and color
type Colors = map[string]color.ARGB

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
		return Colors{}, ErrNoColorFound
	}

	scored := score.Score(quantized, score.ScoreOptions{Desired: 4, Fallback: score.FallbackColor})
	if len(scored) == 0 {
		return Colors{}, ErrNoColorFound
	}

	primary := palettes.NewFromARGB(scored[0])

	scheme := dynamic.NewDynamicScheme(
		scored[0].ToHct(), variant, constrast, dark,
		platform, version, primary,
		nil, nil, nil, nil, nil,
	)

	dcs := scheme.ToColorMap()

	colorMap := map[string]color.ARGB{}
	for key, value := range dcs {
		if value != nil {
			colorMap[key] = value.GetArgb(scheme)
		}
	}
	return colorMap, nil
}
