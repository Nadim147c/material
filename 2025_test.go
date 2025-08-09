package material

import (
	"testing"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/schemes"
)

func Test2025(t *testing.T) {
	input := color.ARGBFromHexMust("#ff0000").ToHct()
	scheme := schemes.NewTonalSpot(input, true, 0, dynamic.Phone, dynamic.V2025)

	t.Run("Test Primary", func(t *testing.T) {
		color := scheme.MaterialColor.Primary().GetArgb(scheme)
		t.Log(color)
	})

	t.Run("Test OnSecondary", func(t *testing.T) {
		color := scheme.ToColorMap()
		for k, v := range color {
			t.Log(k, v.GetArgb(scheme))
		}
	})
}
