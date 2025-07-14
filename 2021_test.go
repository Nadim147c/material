package material

import (
	"embed"
	"image/jpeg"
	"testing"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/palettes"
	"github.com/Nadim147c/material/quantizer"
	"github.com/Nadim147c/material/score"
)

//go:embed quantizer/gophar.jpg
var gophar embed.FS

func Test2021(t *testing.T) {
	file, err := gophar.Open("quantizer/gophar.jpg")
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("failed to decode image: %v", err)
	}

	var pixels []color.ARGB
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			goColor := img.At(x, y)

			pixels = append(pixels, color.ARGBFromInterface(goColor))
		}
	}

	colors := quantizer.QuantizeCelebi(pixels, 10)
	scored := score.NewScore().ScoreColors(colors, score.ScoreOptions{Desired: 4})

	var primary, secondary, ternary *palettes.TonalPalette
	primary = palettes.NewFromARGB(scored[0])

	if len(scored) > 1 {
		secondary = palettes.NewFromARGB(scored[1])
	}
	if len(scored) > 2 {
		ternary = palettes.NewFromARGB(scored[2])
	}

	scheme := dynamic.NewDynamicScheme(
		scored[0].ToHct(), dynamic.Expressive, 0, true,
		dynamic.Phone, dynamic.V2021,
		primary, secondary, ternary,
		nil, nil, nil,
	)

	dcs := scheme.ToColorMap()

	for key, value := range dcs {
		if value != nil {
			c := value.GetArgb(scheme)
			t.Log(key, c)
		}
	}
}
