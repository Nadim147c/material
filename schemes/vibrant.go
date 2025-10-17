package schemes

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

// NewVibrant creates a dynamic color theme that maxes out colorfulness at each
// position in the Primary Tonal Palette.
func NewVibrant(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor,
		dynamic.VariantVibrant,
		construst,
		isDark,
		platform,
		version,
	)
}
