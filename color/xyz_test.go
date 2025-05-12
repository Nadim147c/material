package color

import "testing"

func TestXYZColor_ToRGB(t *testing.T) {
	for _, tt := range ColorTestCases {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.XYZ.ToARGB(); got != tt.ARGB {
				t.Errorf("XYZColor(%v).ToARGB() = %v, want %v", tt.XYZ, got, tt.ARGB)
			}
		})
	}
}

func TestXYZColor_ToLab(t *testing.T) {
	for _, tt := range ColorTestCases {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.XYZ.ToLab(); !almostEqualColor(got, tt.Lab) {
				t.Errorf("XYZColor(%v).ToLab() = %v, want %v", tt.XYZ, got, tt.Lab)
			}
		})
	}
}
