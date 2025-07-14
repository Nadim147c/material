package material

import (
	"embed"
	"image/jpeg"
	"testing"

	"github.com/Nadim147c/material/dynamic"
)

//go:embed quantizer/gophar.jpg
var gophar embed.FS

func TestGenerateFromImage(t *testing.T) {
	file, err := gophar.Open("quantizer/gophar.jpg")
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("failed to decode image: %v", err)
	}

	colors, err := GenerateFromImage(img, dynamic.Expressive, true, 0, dynamic.Phone, dynamic.V2021)
	if err != nil {
		t.Fatalf("failed to generate colors: %v", err)
	}

	for key, value := range colors {
		t.Log(key, value)
	}
}
