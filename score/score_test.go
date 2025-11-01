package score

import (
	"testing"

	"github.com/Nadim147c/material/v2/color"
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
		a := color.ARGB(0xFFD33881)
		b := color.ARGB(0xFF3205CC)
		c := color.ARGB(0xFF0B48CF)
		d := color.ARGB(0xFFA08F5D)
		colorsToPopulation := map[color.ARGB]int{a: 14, b: 77, c: 36, d: 81}

		ranked := Score(
			colorsToPopulation,
			WithFallback(0xff7d772b),
			WithFilter(),
		)

		if len(ranked) != 4 {
			t.Errorf("Expected 3 colors, got %d", len(ranked))
		}

		if ranked[0] != b {
			t.Errorf("Expected %v, got %v", b, ranked[0])
		}
		if ranked[1] != d {
			t.Errorf("Expected %v, got %v", d, ranked[1])
		}
		if ranked[2] != c {
			t.Errorf("Expected %v, got %v", c, ranked[2])
		}
		if ranked[3] != a {
			t.Errorf("Expected %v, got %v", a, ranked[3])
		}
	})
}
