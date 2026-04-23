package score

import (
	"testing"

	"github.com/Nadim147c/material/v3/color"
)

func TestScoring(t *testing.T) {
	t.Run("prioritizes chroma", func(t *testing.T) {
		colorsToPopulation := map[color.ARGB]int{
			color.ARGB(0xff000000): 1,
			color.ARGB(0xffffffff): 1,
			color.ARGB(0xff0000ff): 1,
		}

		ranked := Score(colorsToPopulation, WithFilter())

		if len(ranked) != 1 {
			t.Errorf("Expected 1 color, got %d", len(ranked))
		}
		if ranked[0] != color.ARGB(0xff0000ff) {
			t.Errorf("Expected 0xff0000ff, got %v", ranked[0])
		}
	})

	t.Run("prioritizes chroma when proportions equal", func(t *testing.T) {
		colorsToPopulation := map[color.ARGB]int{
			color.ARGB(0xffff0000): 1,
			color.ARGB(0xff00ff00): 1,
			color.ARGB(0xff0000ff): 1,
		}

		ranked := Score(colorsToPopulation)

		if len(ranked) != 3 {
			t.Errorf("Expected 3 colors, got %d", len(ranked))
		}
		if ranked[0] != color.ARGB(0xffff0000) {
			t.Errorf("Expected 0xffff0000, got %v", ranked[0])
		}
		if ranked[1] != color.ARGB(0xff00ff00) {
			t.Errorf("Expected 0xff00ff00, got %v", ranked[1])
		}
		if ranked[2] != color.ARGB(0xff0000ff) {
			t.Errorf("Expected 0xff0000ff, got %v", ranked[2])
		}
	})

	t.Run("generates gBlue when no colors available", func(t *testing.T) {
		colorsToPopulation := map[color.ARGB]int{
			color.ARGB(0xff000000): 1,
		}

		ranked := Score(colorsToPopulation, WithFilter())

		if len(ranked) != 1 {
			t.Errorf("Expected 1 color, got %d", len(ranked))
		}
		if ranked[0] != GoogleBlue {
			t.Errorf("Expected 0xff4285f4, got %v", ranked[0])
		}
	})

	t.Run("dedupes nearby hues", func(t *testing.T) {
		colorsToPopulation := map[color.ARGB]int{
			color.ARGB(0xff008772): 1, // H 180 C 42 T 50
			color.ARGB(0xff318477): 1, // H 184 C 35 T 50
		}

		ranked := Score(colorsToPopulation)

		if len(ranked) != 1 {
			t.Errorf("Expected 1 color, got %d", len(ranked))
		}
		if ranked[0] != color.ARGB(0xff008772) {
			t.Errorf("Expected 0xff008772, got %v", ranked[0])
		}
	})

	t.Run("maximizes hue distance", func(t *testing.T) {
		colorsToPopulation := map[color.ARGB]int{
			color.ARGB(0xff008772): 1, // H 180 C 42 T 50
			color.ARGB(0xff008587): 1, // H 198 C 50 T 50
			color.ARGB(0xff007ebc): 1, // H 245 C 50 T 50
		}

		ranked := Score(colorsToPopulation, WithLimit(2))

		if len(ranked) != 2 {
			t.Errorf("Expected 2 colors, got %d", len(ranked))
		}
		if ranked[0] != color.ARGB(0xff007ebc) {
			t.Errorf("Expected 0xff007ebc, got %v", ranked[0])
		}
		if ranked[1] != color.ARGB(0xff008772) {
			t.Errorf("Expected 0xff008772, got %v", ranked[1])
		}
	})

	t.Run("passes generated scenario one", func(t *testing.T) {
		colorsToPopulation := map[color.ARGB]int{
			color.ARGB(0xff7ea16d): 67,
			color.ARGB(0xffd8ccae): 67,
			color.ARGB(0xff835c0d): 49,
		}

		ranked := Score(
			colorsToPopulation,
			WithFallback(0xff8d3819),
			WithLimit(3),
		)

		if len(ranked) != 3 {
			t.Errorf("Expected 3 colors, got %d", len(ranked))
		}
		if ranked[0] != color.ARGB(0xff7ea16d) {
			t.Errorf("Expected 0xff7ea16d, got %v", ranked[0])
		}
		if ranked[1] != color.ARGB(0xffd8ccae) {
			t.Errorf("Expected 0xffd8ccae, got %v", ranked[1])
		}
		if ranked[2] != color.ARGB(0xff835c0d) {
			t.Errorf("Expected 0xff835c0d, got %v", ranked[2])
		}
	})

	t.Run("passes generated scenario two", func(t *testing.T) {
		pink := color.ARGBFromHexMust("#D33881")
		purple := color.ARGBFromHexMust("#3205CC")
		blue := color.ARGBFromHexMust("#0B48CF")
		sand := color.ARGBFromHexMust("#A08F5D")
		colorsToPopulation := map[color.ARGB]int{pink: 14, purple: 77, blue: 36, sand: 81}

		ranked := Score(
			colorsToPopulation,
			WithFallback(0xff7d772b),
			WithFilter(),
		)

		// blue gets filter out since, its similar to purple
		colors := []color.ARGB{purple, sand, pink}

		if len(ranked) != len(colors) {
			t.Errorf("Expected %d colors, got %d", len(colors), len(ranked))
		}

		for i, got := range ranked {
			expected := colors[i]
			if expected != got {
				t.Errorf("Expected %a, got %a", expected, got)
			}
		}
	})
}
