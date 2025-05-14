package color

import (
	"math"
	"testing"
)

func TestHctLimitedToSRGB(t *testing.T) {
	// Ensures that the HCT class can only represent sRGB colors.
	// An impossibly high chroma is used.
	hct := NewHct(120.0, 200.0, 50.0)
	argb := hct.ToColor()

	// The hue, chroma, and tone members of hct should actually
	// represent the sRGB color.
	cam := argb.ToCam16()
	t.Log(cam.ToColor().AnsiBg("  "))
	if !floatEquals(cam.Hue, hct.Hue) {
		t.Errorf("Expected hue %f, got %f", hct.Hue, cam.Hue)
	}
	if !floatEquals(cam.Chroma, hct.Chroma) {
		t.Errorf("Expected chroma %f, got %f", hct.Chroma, cam.Chroma)
	}
	if !floatEquals(cam.J, hct.Tone) {
		t.Errorf("Expected tone %f, got %f", hct.Tone, cam.J)
	}
}

func TestHctTruncatesColors(t *testing.T) {
	// Ensures that HCT truncates colors.
	hct := NewHct(120.0, 60.0, 50.0)
	chroma := hct.Chroma

	if chroma >= 60.0 {
		t.Errorf("Expected chroma < 60.0, got %f", chroma)
	}

	// The new chroma should be lower than the original.
	hct.SetTone(180.0)
	if hct.Chroma >= chroma {
		t.Errorf("Expected new chroma < %f, got %f", chroma, hct.Chroma)
	}
}

func isOnBoundary(rgbComponent uint8) bool {
	return rgbComponent == 0 || rgbComponent == 255
}

func colorIsOnBoundary(argb Color) bool {
	return isOnBoundary(argb.Red()) || isOnBoundary(argb.Green()) || isOnBoundary(argb.Blue())
}

func TestHctCorrectness(t *testing.T) {
	hues := []int{15, 45, 75, 105, 135, 165, 195, 225, 255, 285, 315, 345}
	chromas := []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	tones := []int{20, 30, 40, 50, 60, 70, 80}

	for _, hue := range hues {
		for _, chroma := range chromas {
			for _, tone := range tones {
				t.Run("", func(t *testing.T) {
					testHctCorrectness(t, hue, chroma, tone)
				})
			}
		}
	}
}

func testHctCorrectness(t *testing.T, hue, chroma, tone int) {
	color := NewHct(float64(hue), float64(chroma), float64(tone))

	if chroma > 0 {
		if math.Abs(color.Hue-float64(hue)) > 4.0 {
			t.Errorf("Expected hue near %d, got %f", hue, color.Hue)
		}
	}

	if color.Chroma >= float64(chroma)+2.5 {
		t.Errorf("Expected chroma < %f, got %f", float64(chroma)+2.5, color.Chroma)
	}

	if color.Chroma < float64(chroma)-2.5 {
		if !colorIsOnBoundary(color.ToColor()) {
			t.Errorf("Expected color to be on boundary when chroma is significantly reduced")
		}
	}

	if math.Abs(color.Tone-float64(tone)) > 0.5 {
		t.Errorf("Expected tone near %d, got %f", tone, color.Tone)
	}
}

// floatEquals checks if two floats are equal within a small epsilon
func floatEquals(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}
