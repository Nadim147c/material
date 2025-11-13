package material

import (
	"fmt"
	"testing"

	"github.com/Nadim147c/material/v2/color"
)

func ExampleGenerate() {
	c := color.ARGBFromHexMust("#FF0000")
	colors, _ := Generate(
		FromColor(c),
		WithDark(true),
		WithContrast(0.5),
		WithVariant(VariantExpressive),
	)
	fmt.Println(colors.Primary)
}

func TestGenerate(t *testing.T) {
	pixels := []string{"#FF1100", "#11FF00", "#1111FF", "#3F3F11", "#007755"}

	colors, err := Generate(FromHexes(pixels), WithDark(true))
	if err != nil {
		t.Errorf("failed to generate colors: %v", err)
	}
	t.Log(colors)
}
