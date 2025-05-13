package quantizer

import (
	"image/jpeg"
	"os"
	"testing"

	"github.com/Nadim147c/goyou/color"
)

func TestQuantizeWu(t *testing.T) {
	// Load the test image
	file, err := os.Open("./gophar.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("failed to decode image: %v", err)
	}

	// Convert image pixels to []color.Color (your type)
	var pixels []color.Color
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			goCol := img.At(x, y)
			pixels = append(pixels, color.FromGoColor(goCol))
		}
	}

	// Perform Wu quantization
	result := QuantizeWu(pixels, 3)

	if len(result) == 0 {
		t.Fatal("QuantizeWu() returned no colors")
	}

	for _, color := range result {
		t.Logf("Cluster %s %s", color.HexRGB(), color.AnsiBg("  "))
	}
}

func containsColor(colors []color.Color, target uint32) bool {
	for _, c := range colors {
		if c == color.Color(target) {
			return true
		}
	}
	return false
}

func TestWu_OneRed(t *testing.T) {
	pixels := []color.Color{color.Color(0xffff0000)}
	result := QuantizeWu(pixels, 256)
	if len(result) != 1 {
		t.Fatalf("Expected 1 color, got %d", len(result))
	}
	if result[0] != 0xffff0000 {
		t.Errorf("Expected 0xffff0000, got %s", result[0].HexRGBA())
	}
}

func TestWu_OneGreen(t *testing.T) {
	pixels := []color.Color{color.Color(0xff00ff00)}
	result := QuantizeWu(pixels, 256)
	if len(result) != 1 {
		t.Fatalf("Expected 1 color, got %d", len(result))
	}
	if result[0] != 0xff00ff00 {
		t.Errorf("Expected 0xff00ff00, got %s", result[0].HexRGBA())
	}
}

func TestWu_OneBlue(t *testing.T) {
	pixels := []color.Color{color.Color(0xff0000ff)}
	result := QuantizeWu(pixels, 256)
	if len(result) != 1 {
		t.Fatalf("Expected 1 color, got %d", len(result))
	}
	if result[0] != 0xff0000ff {
		t.Errorf("Expected 0xff0000ff, got %s", result[0].HexRGBA())
	}
}

func TestWu_FiveBlue(t *testing.T) {
	pixels := make([]color.Color, 5)
	for i := 0; i < 5; i++ {
		pixels[i] = color.Color(0xff0000ff)
	}
	result := QuantizeWu(pixels, 256)
	if len(result) != 1 {
		t.Fatalf("Expected 1 color, got %d", len(result))
	}
	if result[0] != 0xff0000ff {
		t.Errorf("Expected 0xff0000ff, got %s", result[0].HexRGBA())
	}
}

func TestWu_TwoRedThreeGreen(t *testing.T) {
	pixels := []color.Color{
		color.Color(0xffff0000), color.Color(0xffff0000), color.Color(0xffff0000),
		color.Color(0xff00ff00), color.Color(0xff00ff00),
	}
	result := QuantizeWu(pixels, 256)
	if len(result) != 2 {
		t.Fatalf("Expected 2 colors, got %d", len(result))
	}
	if !containsColor(result, 0xffff0000) || !containsColor(result, 0xff00ff00) {
		t.Errorf("Expected red and green in result, got %+v", result)
	}
}

func TestWu_RedGreenBlue(t *testing.T) {
	pixels := []color.Color{
		color.Color(0xffff0000),
		color.Color(0xff00ff00),
		color.Color(0xff0000ff),
	}
	result := QuantizeWu(pixels, 256)
	if len(result) != 3 {
		t.Fatalf("Expected 3 colors, got %d", len(result))
	}
	expected := []uint32{0xffff0000, 0xff00ff00, 0xff0000ff}
	for _, want := range expected {
		if !containsColor(result, want) {
			t.Errorf("Expected color %08x in result", want)
		}
	}
}

func TestWu_OneRedAndO(t *testing.T) {
	pixels := []color.Color{color.Color(0xff141216)}
	result := QuantizeWu(pixels, 256)
	if len(result) != 1 {
		t.Fatalf("Expected 1 color, got %d", len(result))
	}
	if result[0] != 0xff141216 {
		t.Errorf("Expected 0xff141216, got %s", result[0].HexRGBA())
	}
}

func TestWu_TestOnlyMisc(t *testing.T) {
	pixels := []color.Color{
		color.Color(0xff010203),
		color.Color(0xff665544),
		color.Color(0xff708090),
		color.Color(0xffc0ffee),
		color.Color(0xfffedcba),
	}
	result := QuantizeWu(pixels, 256)
	if len(result) == 0 {
		t.Fatal("Expected at least one color")
	}
}
