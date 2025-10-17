package schemes

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

// NewTonalSpot creates a dynamic color theme with low to medium colorfulness
// and a Tertiary TonalPalette with a hue related to the source color.
//
// The default Material You theme on Android 12 and 13.
func NewTonalSpot(sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor,
		dynamic.VariantTonalSpot,
		construst,
		isDark,
		platform,
		version,
	)
}
