package schemes

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

// NewFruitSalad creates a playful theme - the source color's hue does not
// appear in the theme.
func NewFruitSalad(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor,
		dynamic.VariantFruitSalad,
		construst,
		isDark,
		platform,
		version,
	)
}
