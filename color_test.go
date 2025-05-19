package goyou

import (
	"fmt"
	"testing"
	"time"

	"github.com/Nadim147c/goyou/color"
	"github.com/Nadim147c/goyou/dynamic"
	"github.com/Nadim147c/goyou/schemes"
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

	t.Fatal("sadljsdlajdlaskjd")
}
