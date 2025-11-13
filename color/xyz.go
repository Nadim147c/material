package color

import (
	"math"

	"github.com/Nadim147c/material/v2/num"
)

//revive:disable:var-naming

// RGB_TO_XYZ is the Matrix3 that is used for converting sRGB colors to XYZ
var RGB_TO_XYZ = num.NewMatrix3(
	0.41233895, 0.35762064, 0.18051042,
	0.2126, 0.7152, 0.0722,
	0.01932141, 0.11916382, 0.95034478,
)

// XYZ_TO_RGB is the Matrix3 that is used for converting XYZ colors to sRGB
var XYZ_TO_RGB = num.NewMatrix3(
	3.2413774792388685, -1.5376652402851851, -0.49885366846268053,
	-0.9691452513005321, 1.8758853451067872, 0.04156585616912061,
	0.05562093689691305, -0.20395524564742123, 1.0571799111220335,
)

//revive:enable:var-naming

// The XYZ color space, developed by the International Commission on
// Illumination (CIE) in 1931, is a mathematical model that represents all
// colors visible to the human eye based on three components: X, Y, and Z. It
// serves as a standard reference for color measurement and is widely used in
// various industries for color management and communication.
type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// NewXYZ creates a color in XYZ model.
func NewXYZ(x, y, z float64) XYZ {
	return XYZ{x, y, z}
}

// ToARGB converts c in XZY to ARGB color model.
func (c XYZ) ToARGB() ARGB {
	vec := num.NewVector(c)

	// Get linear values of RGB chanels
	lr, lg, lb := XYZ_TO_RGB.Multiply(vec).Values()

	return ARGBFromRGB(Delinearized3(lr, lg, lb))
}

// ToLinearRGB converts c in XZY to LinearRGB color model.
func (c XYZ) ToLinearRGB() LinearRGB {
	return LinearRGBFromXYZ(c)
}

//revive:disable:function-result-limit

// RGBA implements the color.Color interface. Returns r, g, b, a values in the
// 0-65535 range.
func (c XYZ) RGBA() (r, g, b, a uint32) {
	return c.ToARGB().RGBA()
}

//revive:enable:function-result-limit

// ToLab convets XYZ to CieLab color model
func (c XYZ) ToLab() Lab {
	x, y, z := c.Values()

	// x,y,z value of WhitePointD65 cordinate
	wx, wy, wz := WhitePointD65.Values()

	// Normalize x,y,z with WhitePointD65
	nx, ny, nz := x/wx, y/wy, z/wz

	fx, fy, fz := labFunc(nx), labFunc(ny), labFunc(nz)

	l := (116.0 * fy) - 16
	a := 500.0 * (fx - fy)
	b := 200.0 * (fy - fz)
	return NewLab(l, a, b)
}

// ToLuv convets XYZ to CIELUV color model
func (c XYZ) ToLuv() Luv {
	return LuvFromXYZ(c)
}

// ToLCHuv convets XYZ to LCHuv color model
func (c XYZ) ToLCHuv() LCHuv {
	return c.ToLuv().ToLCHuv()
}

// ToLCHab convets XYZ to LCHab color model
func (c XYZ) ToLCHab() LCHab {
	return c.ToLab().ToLCHab()
}

// ToOkLab convets XYZ to OkLan color model
func (c XYZ) ToOkLab() OkLab {
	return OkLabFromXYZ(c)
}

// ToCam16 converts XYZ to color appearance model (Cam16)
func (c XYZ) ToCam16() Cam16 {
	return Cam16FromXYZInEnv(c, DefaultEnvironment)
}

// ToHct convets XYZ to Hct (Hue, Chroma, Tone) model
func (c XYZ) ToHct() Hct {
	cam := c.ToCam16()
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
	}
	return math.Pow((normalized+0.055)/1.055, 2.4) * 100
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
