package color_test

import (
	"testing"

	"github.com/Nadim147c/goyou/color"
)

func TestColor_Red(t *testing.T) {
	c := color.FromARGB(0, 11, 0, 0)
	want := color.Color(11)
	if c.Red() != want {
		t.Fatalf("Color = %d, Color.Red() = %d, want = %d", c, c.Red(), want)
	}
}
