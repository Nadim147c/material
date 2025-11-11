package color

import "testing"

func almostEqualLuv(a, b Luv) bool {
	return almostEqual(a.L, b.L) &&
		almostEqual(a.U, b.U) &&
		almostEqual(a.V, b.V)
}

func TestLuvRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		argb ARGB
		luv  Luv
	}{
		{"black", 0xFF000000, NewLuv(0, 0, 0)},
		{"white", 0xFFFFFFFF, NewLuv(100, 0, 0)},
		{"red", 0xFFFF0000, NewLuv(53.23, 175, 37.75)},
		{"green", 0xFF00FF00, NewLuv(87.74, -83.08, 107.4)},
		{"blue", 0xFF0000FF, NewLuv(32.3, -9.38, -130.36)},
		{"yellow", 0xFFFFFF00, NewLuv(97.14, 7.7, 106.79)},
		{"cyan", 0xFF00FFFF, NewLuv(91.12, -70.45, -15.19)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			luv := tt.argb.ToXYZ().ToLuv()
			t.Log(luv)
			if !almostEqualLuv(tt.luv, luv) {
				t.Errorf("wrong Luv. want %v, got %v", tt.luv, luv)
			}
			got := luv.ToXYZ().ToARGB()
			if tt.argb != got {
				t.Errorf("failed Luv round-trip. Source %s, Returned %s",
					tt.argb.String(), got.String())
			}
		})
	}
}
