package score

import (
	"testing"

	"github.com/Nadim147c/goyou/color"
)

func TestScoring(t *testing.T) {
	t.Run("prioritizes chroma", func(t *testing.T) {
		colorsToPopulation := map[color.Color]int{
			color.Color(0xff000000): 1,
			color.Color(0xffffffff): 1,
			color.Color(0xff0000ff): 1,
		}

		ranked := Score(colorsToPopulation, &ScoreOptions{Desired: 4, Filter: true})

		if len(ranked) != 1 {
			t.Errorf("Expected 1 color, got %d", len(ranked))
		}
		if ranked[0] != color.Color(0xff0000ff) {
			t.Errorf("Expected 0xff0000ff, got %v", ranked[0])
		}
	})

	t.Run("prioritizes chroma when proportions equal", func(t *testing.T) {
		colorsToPopulation := map[color.Color]int{
			color.Color(0xffff0000): 1,
			color.Color(0xff00ff00): 1,
			color.Color(0xff0000ff): 1,
		}

		ranked := Score(colorsToPopulation, &ScoreOptions{Desired: 4})

		if len(ranked) != 3 {
			t.Errorf("Expected 3 colors, got %d", len(ranked))
		}
		if ranked[0] != color.Color(0xffff0000) {
			t.Errorf("Expected 0xffff0000, got %v", ranked[0])
		}
		if ranked[1] != color.Color(0xff00ff00) {
			t.Errorf("Expected 0xff00ff00, got %v", ranked[1])
		}
		if ranked[2] != color.Color(0xff0000ff) {
			t.Errorf("Expected 0xff0000ff, got %v", ranked[2])
		}
	})

	t.Run("generates gBlue when no colors available", func(t *testing.T) {
		colorsToPopulation := map[color.Color]int{
			color.Color(0xff000000): 1,
		}

		ranked := Score(colorsToPopulation, &ScoreOptions{Desired: 4, Filter: true})

		if len(ranked) != 1 {
			t.Errorf("Expected 1 color, got %d", len(ranked))
		}
		if ranked[0] != FallbackColor {
			t.Errorf("Expected 0xff4285f4, got %v", ranked[0])
		}
	})

	t.Run("dedupes nearby hues", func(t *testing.T) {
		colorsToPopulation := map[color.Color]int{
			color.Color(0xff008772): 1, // H 180 C 42 T 50
			color.Color(0xff318477): 1, // H 184 C 35 T 50
		}

		ranked := Score(colorsToPopulation, &ScoreOptions{Desired: 4})

		if len(ranked) != 1 {
			t.Errorf("Expected 1 color, got %d", len(ranked))
		}
		if ranked[0] != color.Color(0xff008772) {
			t.Errorf("Expected 0xff008772, got %v", ranked[0])
		}
	})

	t.Run("maximizes hue distance", func(t *testing.T) {
		colorsToPopulation := map[color.Color]int{
			color.Color(0xff008772): 1, // H 180 C 42 T 50
			color.Color(0xff008587): 1, // H 198 C 50 T 50
			color.Color(0xff007ebc): 1, // H 245 C 50 T 50
		}

		ranked := Score(colorsToPopulation, &ScoreOptions{Desired: 2})

		if len(ranked) != 2 {
			t.Errorf("Expected 2 colors, got %d", len(ranked))
		}
		if ranked[0] != color.Color(0xff007ebc) {
			t.Errorf("Expected 0xff007ebc, got %v", ranked[0])
		}
		if ranked[1] != color.Color(0xff008772) {
			t.Errorf("Expected 0xff008772, got %v", ranked[1])
		}
	})

	t.Run("passes generated scenario one", func(t *testing.T) {
		colorsToPopulation := map[color.Color]int{
			color.Color(0xff7ea16d): 67,
			color.Color(0xffd8ccae): 67,
			color.Color(0xff835c0d): 49,
		}

		ranked := Score(colorsToPopulation, &ScoreOptions{Desired: 3, Fallback: 0xff8d3819, Filter: false})

		if len(ranked) != 3 {
			t.Errorf("Expected 3 colors, got %d", len(ranked))
		}
		if ranked[0] != color.Color(0xff7ea16d) {
			t.Errorf("Expected 0xff7ea16d, got %v", ranked[0])
		}
		if ranked[1] != color.Color(0xffd8ccae) {
			t.Errorf("Expected 0xffd8ccae, got %v", ranked[1])
		}
		if ranked[2] != color.Color(0xff835c0d) {
			t.Errorf("Expected 0xff835c0d, got %v", ranked[2])
		}
	})

	t.Run("passes generated scenario two", func(t *testing.T) {
		a := color.Color(0xFFD33881)
		b := color.Color(0xFF3205CC)
		c := color.Color(0xFF0B48CF)
		d := color.Color(0xFFA08F5D)
		colorsToPopulation := map[color.Color]int{a: 14, b: 77, c: 36, d: 81}

		ranked := Score(colorsToPopulation, &ScoreOptions{Desired: 4, Fallback: 0xff7d772b, Filter: true})

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
