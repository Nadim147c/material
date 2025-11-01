package schemes

import (
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
)

// NewContent create A scheme that places the source color in
// `Scheme.PrimaryContainer`.
//
// Primary Container is the source color, adjusted for color relativity. It
// maintains constant appearance in light mode and dark mode. This adds ~5 tone
// in light mode, and subtracts ~5 tone in dark mode. Tertiary Container is the
// complement to the source color, using `TemperatureCache`. It also maintains
// constant appearance.
func NewContent(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor,
		dynamic.VariantContent,
		construst,
		isDark,
		platform,
		version,
	)
}
