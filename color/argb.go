package color

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Nadim147c/goyou/num"
)

// Offset indiacates bit offset of Components in Color
const (
	blueOffset Color = iota * 8
	greenOffset
	redOffset
	alphaOffset
)

// Brightest is the max value of uint8 color
const Brightest = uint8(0xFF) // 255

// Color is an ARGB color packed into a uint32.
type Color uint32

// FromARGB creates a Color from xyz color space cordinates.
func FromXYZ(x, y, z float64) Color {
	return NewXYZColor(x, y, z).ToARGB()
}

// FromARGB creates a Color from individual 8-bit red, green, and blue
// components.
func FromRGB(r, g, b uint8) Color {
	return FromARGB(0xFF, r, g, b)
}

// FromARGB creates a Color from individual 8-bit alpha, red, green, and blue
// components.
func FromARGB(a, r, g, b uint8) Color {
	return Color(a)<<alphaOffset |
		Color(r)<<redOffset |
		Color(g)<<greenOffset |
		Color(b)<<blueOffset
}

// ToXYZ return XYZ color version for c.
func (c Color) ToXYZ() XYZColor {
	r, g, b := c.Red(), c.Green(), c.Blue()
	lr, lg, lb := Linearized(r), Linearized(g), Linearized(b)
	x, y, z := num.NewVector3(lr, lg, lb).MultiplyMatrix(SRGB_TO_XYZ).Values()
	return XYZColor{x, y, z}
}

// ToLab convert Color to LabColor
func (c Color) ToLab() LabColor {
	x, y, z := c.ToXYZ().Values()

	// x,y,z value of WhitePointD65 cordinate
	wx, wy, wz := WhitePointD65.Values()

	// Normalize x,y,z with WhitePointD65
	nx, ny, nz := x/wx, y/wy, z/wz

	// Apologies!! Some Magic function that I have no idea
	fx, fy, fz := LabF(nx), LabF(ny), LabF(nz)

	// Magic Transformations 🌟
	l := (116.0 * fy) - 16
	a := 500.0 * (fx - fy)
	b := 200.0 * (fy - fz)
	return NewLabColor(l, a, b)
}

func (c Color) Values() (uint8, uint8, uint8, uint8) {
	return c.Alpha(), c.Red(), c.Green(), c.Blue()
}

// Alpha returns the 8-bit alpha component of the color.
func (c Color) Alpha() uint8 {
	return uint8((c >> alphaOffset) & 0xFF)
}

// Red returns the 8-bit red component of the color.
func (c Color) Red() uint8 {
	return uint8((c >> redOffset) & 0xFF)
}

// Green returns the 8-bit green component of the color.
func (c Color) Green() uint8 {
	return uint8((c >> greenOffset) & 0xFF)
}

// Blue returns the 8-bit blue component of the color.
func (c Color) Blue() uint8 {
	return uint8((c >> blueOffset) & 0xFF)
}

// HexARGB return #RRGGBB represetation of the color
func (c Color) HexRGB() string {
	return fmt.Sprintf("#%02X%02X%02X", c.Red(), c.Green(), c.Blue())
}

// HexARGB return #AARRGGBB represetation of the color
func (c Color) HexARGB() string {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.Alpha(), c.Red(), c.Green(), c.Blue())
}

// HexRGBA return #RRGGBBAA represetation of the color
func (c Color) HexRGBA() string {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.Red(), c.Green(), c.Blue(), c.Alpha())
}

// FromHex parses a hex color string and returns a Color.
// Supports formats: #RGB, #RGBA, #RRGGBB, #RRGGBBAA
func FromHex(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")

	// Regex check if input is valid or not
	hexColorRegex := regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	if !hexColorRegex.MatchString(hex) {
		return 0, errors.New("invalid hex color format")
	}

	// Expand shorthand formats to full 6/8-char form
	switch len(hex) {
	case 3: // #RGB → #RRGGBB
		hex = fmt.Sprintf("%c%c%c%c%c%c",
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
		)
	case 4: // #RGBA → #RRGGBBAA
		hex = fmt.Sprintf("%c%c%c%c%c%c%c%c",
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
			hex[3], hex[3],
		)
	}

	var r, g, b, a uint8 = 0, 0, 0, 0xFF

	switch len(hex) {
	case 6:
		val, err := strconv.ParseUint(hex, 16, 24)
		if err != nil {
			return 0, fmt.Errorf("invalid hex color: %w", err)
		}
		r = uint8(val >> 16)
		g = uint8((val >> 8) & 0xFF)
		b = uint8(val & 0xFF)

	case 8:
		val, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid hex color: %w", err)
		}
		r = uint8(val >> 24)
		g = uint8((val >> 16) & 0xFF)
		b = uint8((val >> 8) & 0xFF)
		a = uint8(val & 0xFF)
	}

	return FromARGB(a, r, g, b), nil
}
