package schemes

import (
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
)

// NewNeutral crates a dynamic color theme that is near grayscale.
func NewNeutral(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor,
		dynamic.VariantNeutral,
		construst,
		isDark,
		platform,
		version,
	)
}
