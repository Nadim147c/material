package color

import (
	"math"
	"testing"
)

func TestCam(t *testing.T) {
	const tolerance = 0.001

	tests := []struct {
		name      string
		argb      Color
		expected  Cam16
		roundTrip bool
	}{
		{
			name: "Red",
			argb: 0xffff0000,
			expected: Cam16{
				Hue:    27.408,
				Chroma: 113.357,
				J:      46.445,
				M:      89.494,
				S:      91.889,
				Q:      105.988,
			},
			roundTrip: true,
		},
		{
			name: "Green",
			argb: 0xff00ff00,
			expected: Cam16{
				Hue:    142.139,
				Chroma: 108.410,
				J:      79.331,
				M:      85.587,
				S:      78.604,
				Q:      138.520,
			},
			roundTrip: true,
		},
		{
			name: "Blue",
			argb: 0xff0000ff,
			expected: Cam16{
				Hue:    282.788,
				Chroma: 87.230,
				J:      25.465,
				M:      68.867,
				S:      93.674,
				Q:      78.481,
			},
			roundTrip: true,
		},
		{
			name: "White",
			argb: 0xffffffff,
			expected: Cam16{
				Hue:    209.492,
				Chroma: 2.869,
				J:      100.0,
				M:      2.265,
				S:      12.068,
				Q:      155.521,
			},
			roundTrip: false,
		},
		{
			name: "Black",
			argb: 0xff000000,
			expected: Cam16{
				Hue:    0.0,
				Chroma: 0.0,
				J:      0.0,
				M:      0.0,
				S:      0.0,
				Q:      0.0,
			},
			roundTrip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cam16FromColor(tt.argb)

			if math.Abs(c.Hue-tt.expected.Hue) > tolerance {
				t.Errorf("Hue = %f; want %f", c.Hue, tt.expected.Hue)
			}
			if math.Abs(c.Chroma-tt.expected.Chroma) > tolerance {
				t.Errorf("Chroma = %f; want %f", c.Chroma, tt.expected.Chroma)
			}
			if math.Abs(c.J-tt.expected.J) > tolerance {
				t.Errorf("J = %f; want %f", c.J, tt.expected.J)
			}
			if math.Abs(c.M-tt.expected.M) > tolerance {
				t.Errorf("M = %f; want %f", c.M, tt.expected.M)
			}
			if math.Abs(c.S-tt.expected.S) > tolerance {
				t.Errorf("S = %f; want %f", c.S, tt.expected.S)
			}
			if math.Abs(c.Q-tt.expected.Q) > tolerance {
				t.Errorf("Q = %f; want %f", c.Q, tt.expected.Q)
			}

			if tt.roundTrip {
				roundTripped := c.ToColor()
				if roundTripped != tt.argb {
					t.Errorf("Round-trip = %s; want %s",
						roundTripped.HexRGBA()+roundTripped.AnsiBg("  "),
						tt.argb.HexRGBA()+tt.argb.AnsiBg("  "))
				}
			}
		})
	}
}
