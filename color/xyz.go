package color

import (
	"math"

	"github.com/Nadim147c/goyou/num"
)

// XYZColor is color in XYZ color space
type XYZColor [3]float64

func NewXYZColor(x, y, z float64) XYZColor {
	return XYZColor{x, y, z}
}

// Values returns x, y, z values of XYZColor
func (c XYZColor) Values() (float64, float64, float64) {
	return c[0], c[1], c[2]
}

func (c XYZColor) ToARGB() Color {
	x, y, z := c.Values()

	// Get linear values of RGB chanels
	lr, lg, lb := XYZ_TO_SRGB.MultiplyXYZ(x, y, z).Values()

	r, g, b := Delinearized3(lr, lg, lb)
	return ColorFromRGB(r, g, b)
}

func (c XYZColor) ToLab() LabColor {
	x, y, z := c.Values()

	// x,y,z value of WhitePointD65 cordinate
	wx, wy, wz := WhitePointD65.Values()

	// Normalize x,y,z with WhitePointD65
	nx, ny, nz := x/wx, y/wy, z/wz

	fx, fy, fz := LabFunc(nx), LabFunc(ny), LabFunc(nz)

	l := (116.0 * fy) - 16
	a := 500.0 * (fx - fy)
	b := 200.0 * (fy - fz)
	return NewLabColor(l, a, b)
}

// Luminance returns the Y value of XYZColor
func (c XYZColor) Luminance() float64 {
	return c[1]
}

// LStar returns the L* value of L*a*b* (LabColor)
func (c XYZColor) LStar() float64 {
	return LstarFromY(c[1])
}

// Linearized takes component (uint8) that represents R/G/B channel.
// Returns 0.0 <= output <= 1.0, color channel converted to linear RGB space
func Linearized(component uint8) float64 {
	normalized := float64(num.Clamp(0, 0xFF, component)) / 0xFF
	if normalized <= 0.040449936 {
		return normalized / 12.92 * 100
	} else {
		return math.Pow((normalized+0.055)/1.055, 2.4) * 100
	}
}

// Linearized3 is like Linearized but takes 3 input and returns 3 output.
func Linearized3(x, y, z uint8) (float64, float64, float64) {
	return Linearized(x), Linearized(y), Linearized(z)
}

// Delinearized takes component (float64) that represents linear R/G/B channel.
// Component should be 0.0 < component < 1.0. Returns uint8 (0 <= n <= 255)
// representation of color component.
func Delinearized(component float64) uint8 {
	normalized := num.Clamp(0, 1, component/100)

	delinearized := 0.0
	if normalized <= 0.0031308 {
		delinearized = normalized * 12.92
	} else {
		delinearized = 1.055*math.Pow(normalized, 1.0/2.4) - 0.055
	}
	return num.Clamp(0, 0xFF, uint8(math.Round(delinearized*255.0)))
}

// Delinearized3 is like Delinearized but takes 3 input and returns 3 output.
func Delinearized3(x, y, z float64) (uint8, uint8, uint8) {
	return Delinearized(x), Delinearized(y), Delinearized(z)
}
