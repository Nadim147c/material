package material

import (
	"fmt"
	"testing"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/schemes"
)

func TestMain(t *testing.T) {
	red := color.ARGBFromRGB(0xff, 0, 0)
	scheme := schemes.NewTonalSpot(red.ToHct(), true, 0.5, dynamic.Phone, dynamic.V2021)
	colors := scheme.ToColorMap()

	for name, color := range colors {
		if color != nil {
			fmt.Println(name, color.GetArgb(scheme).String())
		}
	}
}
