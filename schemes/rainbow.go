package schemes

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

// NewRainbow create a playful theme - the source color's hue does not appear in
// the theme.
func NewRainbow(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor,
		dynamic.VariantRainbow,
		construst,
		isDark,
		platform,
		version,
	)
}
