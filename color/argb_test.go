package color

import "testing"

func TestColor_ToXYZ(t *testing.T) {
	for _, tt := range ColorTestCases {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.ARGB.ToXYZ(); !sameXYZ(got, tt.XYZ) {
				t.Errorf("Color(%s).ToXYZ() = %v, want %v", tt.ARGB.String(), got, tt.XYZ)
			}
		})
	}
}

func TestColor_RoundTrip(t *testing.T) {
	for _, tt := range ColorTestCases {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			if got := tt.ARGB.ToLab().ToARGB().ToLab().ToARGB(); got != tt.ARGB {
				t.Errorf("Color(%s) Round Trip = %v, want %v", tt.ARGB, got, tt.ARGB)
			}
		})
	}
}

func TestFromARGB(t *testing.T) {
	tests := []struct {
		name       string
		a, r, g, b uint8
		want       ARGB
	}{
		{
			name: "Black fully opaque",
			a:    0xFF, r: 0x00, g: 0x00, b: 0x00,
			want: ARGB(0xFF000000),
		},
		{
			name: "Semi-transparent purple",
			a:    0x80, r: 0x80, g: 0x00, b: 0x80,
			want: ARGB(0x80800080),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewARGB(tt.a, tt.r, tt.g, tt.b); got != tt.want {
				t.Errorf("FromARGB(%#x, %#x, %#x, %#x) = %#x, want %#x",
					tt.a, tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestColor_Components(t *testing.T) {
	tests := []struct {
		name       string
		color      ARGB
		a, r, g, b uint8
	}{
		{
			name:  "White fully opaque",
			color: ARGB(0xFFFFFFFF),
			a:     0xFF, r: 0xFF, g: 0xFF, b: 0xFF,
		},
		{
			name:  "Semi-transparent teal",
			color: ARGB(0x8000FFFF),
			a:     0x80, r: 0x00, g: 0xFF, b: 0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.Alpha(); got != tt.a {
				t.Errorf("Alpha = %#x, want %#x", got, tt.a)
			}
			if got := tt.color.Red(); got != tt.r {
				t.Errorf("Red = %#x, want %#x", got, tt.r)
			}
			if got := tt.color.Green(); got != tt.g {
				t.Errorf("Green = %#x, want %#x", got, tt.g)
			}
			if got := tt.color.Blue(); got != tt.b {
				t.Errorf("Blue = %#x, want %#x", got, tt.b)
			}
		})
	}
}

func TestFromHex(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    ARGB
		wantErr bool
	}{
		{
			name: "6-digit hex with #",
			hex:  "#00FF00", want: ARGB(0xFF00FF00), wantErr: false,
		},
		{
			name: "Invalid characters",
			hex:  "#GGGGGG", want: ARGB(0), wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ARGBFromHex(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("FromHex(%q) = %#x, want %#x", tt.hex, got, tt.want)
			}
		})
	}
}

func TestColor_HexRGB(t *testing.T) {
	tests := []struct {
		name  string
		color ARGB
		want  string
	}{
		{
			name: "Red", color: ARGB(0xFF0000), want: "#FF0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.HexRGB(); got != tt.want {
				t.Errorf("HexRGB = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColor_HexARGB(t *testing.T) {
	tests := []struct {
		name  string
		color ARGB
		want  string
	}{
		{
			name: "Red semi-transparent", color: ARGB(0x80FF0000), want: "#80FF0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.HexARGB(); got != tt.want {
				t.Errorf("HexARGB = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColor_HexRGBA(t *testing.T) {
	tests := []struct {
		name  string
		color ARGB
		want  string
	}{
		{
			name: "Blue semi-transparent", color: ARGB(0x800000FF), want: "#0000FF80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.HexRGBA(); got != tt.want {
				t.Errorf("HexRGBA = %q, want %q", got, tt.want)
			}
		})
	}
}
