package schemes

import (
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
)

// NewExpressive creates a dynamic Color theme that is intentionally detached
// from the source color.
func NewExpressive(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor,
		dynamic.VariantExpressive,
		construst,
		isDark,
		platform,
		version,
	)
}
