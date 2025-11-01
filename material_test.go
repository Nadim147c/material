package material

import "testing"

func TestGenerate(t *testing.T) {
	pixels := []string{"#FF1100", "#11FF00", "#1111FF", "#3F3F11", "#007755"}

	colors, err := Generate(FromHexes(pixels), WithDark(true))
	if err != nil {
		t.Errorf("failed to generate colors: %v", err)
	}
	t.Log(colors)
}
