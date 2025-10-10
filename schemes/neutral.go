package schemes

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

func NewNeutral(
	sourceColor color.Hct,
	isDark bool,
	construst float64,
	platform dynamic.Platform,
	version dynamic.Version,
) *dynamic.Scheme {
	return dynamic.NewDynamicScheme(
		sourceColor, dynamic.VariantNeutral, construst, isDark, platform, version,
		nil, nil, nil, nil, nil, nil,
	)
}
