package color

import "testing"

func TestFromARGB(t *testing.T) {
	tests := []struct {
		name       string
		a, r, g, b uint8
		want       Color
	}{
		{
			name: "Black fully opaque",
			a:    0xFF,
			r:    0x00,
			g:    0x00,
			b:    0x00,
			want: Color(0xFF000000),
		},
		{
			name: "White fully opaque",
			a:    0xFF,
			r:    0xFF,
			g:    0xFF,
			b:    0xFF,
			want: Color(0xFFFFFFFF),
		},
		{
			name: "Red fully opaque",
			a:    0xFF,
			r:    0xFF,
			g:    0x00,
			b:    0x00,
			want: Color(0xFFFF0000),
		},
		{
			name: "Green fully opaque",
			a:    0xFF,
			r:    0x00,
			g:    0xFF,
			b:    0x00,
			want: Color(0xFF00FF00),
		},
		{
			name: "Blue fully opaque",
			a:    0xFF,
			r:    0x00,
			g:    0x00,
			b:    0xFF,
			want: Color(0xFF0000FF),
		},
		{
			name: "Semi-transparent purple",
			a:    0x80,
			r:    0x80,
			g:    0x00,
			b:    0x80,
			want: Color(0x80800080),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromARGB(tt.a, tt.r, tt.g, tt.b); got != tt.want {
				t.Errorf("FromARGB(%#x, %#x, %#x, %#x) = %#x, want %#x",
					tt.a, tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestColor_Components(t *testing.T) {
	tests := []struct {
		name       string
		color      Color
		a, r, g, b uint8
	}{
		{
			name:  "Black fully opaque",
			color: Color(0xFF000000),
			a:     0xFF,
			r:     0x00,
			g:     0x00,
			b:     0x00,
		},
		{
			name:  "White fully opaque",
			color: Color(0xFFFFFFFF),
			a:     0xFF,
			r:     0xFF,
			g:     0xFF,
			b:     0xFF,
		},
		{
			name:  "Red fully opaque",
			color: Color(0xFFFF0000),
			a:     0xFF,
			r:     0xFF,
			g:     0x00,
			b:     0x00,
		},
		{
			name:  "Semi-transparent teal",
			color: Color(0x8000FFFF),
			a:     0x80,
			r:     0x00,
			g:     0xFF,
			b:     0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.Alpha(); got != tt.a {
				t.Errorf("Color(%#x).Alpha() = %#x, want %#x", tt.color, got, tt.a)
			}
			if got := tt.color.Red(); got != tt.r {
				t.Errorf("Color(%#x).Red() = %#x, want %#x", tt.color, got, tt.r)
			}
			if got := tt.color.Green(); got != tt.g {
				t.Errorf("Color(%#x).Green() = %#x, want %#x", tt.color, got, tt.g)
			}
			if got := tt.color.Blue(); got != tt.b {
				t.Errorf("Color(%#x).Blue() = %#x, want %#x", tt.color, got, tt.b)
			}
		})
	}
}

func TestFromHex(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    Color
		wantErr bool
	}{
		{
			name:    "6-digit hex without #",
			hex:     "FF0000",
			want:    Color(0xFFFF0000),
			wantErr: false,
		},
		{
			name:    "6-digit hex with #",
			hex:     "#00FF00",
			want:    Color(0xFF00FF00),
			wantErr: false,
		},
		{
			name:    "8-digit hex (RRGGBBAA)",
			hex:     "#0000FFAA",
			want:    Color(0xAA0000FF),
			wantErr: false,
		},
		{
			name:    "3-digit hex (RGB)",
			hex:     "#F00",
			want:    Color(0xFFFF0000),
			wantErr: false,
		},
		{
			name:    "4-digit hex (RGBA)",
			hex:     "#F00A",
			want:    Color(0xAAFF0000),
			wantErr: false,
		},
		{
			name:    "Invalid characters",
			hex:     "#XY00ZZ",
			want:    Color(0),
			wantErr: true,
		},
		{
			name:    "Invalid length",
			hex:     "#12345",
			want:    Color(0),
			wantErr: true,
		},
		{
			name:    "Empty string",
			hex:     "",
			want:    Color(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromHex(tt.hex)
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
		color Color
		want  string
	}{
		{
			name:  "Black",
			color: Color(0x000000),
			want:  "#000000",
		},
		{
			name:  "White",
			color: Color(0xFFFFFFFF),
			want:  "#FFFFFF",
		},
		{
			name:  "Red",
			color: Color(0xFF0000),
			want:  "#FF0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.HexRGB(); got != tt.want {
				t.Errorf("Color(%#x).HexRGB() = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}

func TestColor_HexARGB(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{
			name:  "Black fully opaque",
			color: Color(0xFF000000),
			want:  "#FF000000",
		},
		{
			name:  "White fully opaque",
			color: Color(0xFFFFFFFF),
			want:  "#FFFFFFFF",
		},
		{
			name:  "Red semi-transparent",
			color: Color(0x80FF0000),
			want:  "#80FF0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.HexARGB(); got != tt.want {
				t.Errorf("Color(%#x).HexARGB() = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}

func TestColor_HexRGBA(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{
			name:  "Black fully opaque",
			color: Color(0xFF000000),
			want:  "#000000FF",
		},
		{
			name:  "White fully opaque",
			color: Color(0xFFFFFFFF),
			want:  "#FFFFFFFF",
		},
		{
			name:  "Blue semi-transparent",
			color: Color(0x800000FF),
			want:  "#0000FF80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.HexRGBA(); got != tt.want {
				t.Errorf("Color(%#x).HexRGBA() = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}
