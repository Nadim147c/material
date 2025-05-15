package color

import "testing"

func TestHctRoundTrip(t *testing.T) {
	for _, tt := range ColorTestCases {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			if got := tt.ARGB.ToHct().ToColor(); got != tt.ARGB {
				t.Errorf("Color(%s) Round Trip = %v, want %v", tt.ARGB.HexRGBA(), got.String(), tt.ARGB.String())
			}
		})
	}
}
