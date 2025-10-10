package schemes

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

func NewFruitSalad(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.DynamicScheme {
	return dynamic.NewDynamicScheme(
		sourceColor, dynamic.VariantFruitSalad, construst, isDark, platform, version,
		nil, nil, nil, nil, nil, nil,
	)
}
