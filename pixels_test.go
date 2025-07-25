package material

import (
	"testing"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

func TestGenerateFromPixels(t *testing.T) {
	pixels := []color.ARGB{
		0xff001100,
		0xff001100,
		0xff001111,
		0xff001133,
		0xff110000,
	}

	colors, err := GenerateFromPixels(pixels, dynamic.Expressive, true, 0, dynamic.Phone, dynamic.V2021)
	if err != nil {
		t.Fatalf("failed to generate colors: %v", err)
	}

	for key, value := range colors {
		t.Log(key, value)
	}
}
