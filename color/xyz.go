package color

import (
	"math"

	"github.com/Nadim147c/material/num"
)

// XYZ is color in XYZ color space
type XYZ struct {
	X, Y, Z float64
}

// Ensure Color implements the color.Color interface
var _ digitalColor = (*XYZ)(nil)

func NewXYZ(x, y, z float64) XYZ {
	return XYZ{x, y, z}
}

func (c XYZ) ToARGB() ARGB {
	x, y, z := c.X, c.Y, c.Z

	// Get linear values of RGB chanels
	lr, lg, lb := XYZ_TO_SRGB.MultiplyXYZ(x, y, z).Values()

	r, g, b := Delinearized3(lr, lg, lb)
	return ARGBFromRGB(r, g, b)
}

// RGBA implements the color.Color interface.
// It returns r, g, b, a values in the 0-65535 range.
func (c XYZ) RGBA() (uint32, uint32, uint32, uint32) {
	return c.ToARGB().RGBA()
}

func (c XYZ) ToXYZ() XYZ {
	return c
}

func (c XYZ) ToLab() Lab {
	x, y, z := c.Values()

	// x,y,z value of WhitePointD65 cordinate
	wx, wy, wz := WhitePointD65.Values()

	// Normalize x,y,z with WhitePointD65
	nx, ny, nz := x/wx, y/wy, z/wz

	fx, fy, fz := LabFunc(nx), LabFunc(ny), LabFunc(nz)

	l := (116.0 * fy) - 16
	a := 500.0 * (fx - fy)
	b := 200.0 * (fy - fz)
	return NewLab(l, a, b)
}

func (c XYZ) ToCam() *Cam16 {
	return Cam16FromXyzInEnv(c, &DefaultEnviroment)
}

func (c XYZ) ToHct() Hct {
	cam := c.ToCam()
	return NewHct(cam.Hue, cam.Chroma, c.LStar())
}

// Values returns x, y, z values of XYZColor
func (c XYZ) Values() (float64, float64, float64) {
	return c.X, c.Y, c.Z
}

// Luminance returns the Y value of XYZColor
func (c XYZ) Luminance() float64 {
	return c.Y
}

// LStar returns the L* value of L*a*b* (LabColor)
func (c XYZ) LStar() float64 {
	return LstarFromY(c.Y)
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
