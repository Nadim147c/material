package color

import "testing"

func TestOkLabRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		argb ARGB
	}{
		{"black", 0xFF000000},
		{"white", 0xFFFFFFFF},
		{"red", 0xFFFF0000},
		{"green", 0xFF00FF00},
		{"blue", 0xFF0000FF},
		{"yellow", 0xFFFFFF00},
		{"cyan", 0xFF00FFFF},
		{"magenta", 0xFFFF00FF},
		{"gray", 0xFF808080},
		{"random", 0xFF123456},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oklab := tt.argb.ToXYZ().ToOkLab()
			t.Log(oklab)
			got := oklab.ToXYZ().ToARGB()
			if tt.argb != got {
				t.Errorf(
					"Failed OkLab round-trip! Source %s, Returned %s",
					tt.argb.String(),
					got.String(),
				)
			}
		})
	}
}
