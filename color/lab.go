package color

import "math"

// Lab is the CIELAB color space, also referred to as L*a*b*, is a color space
// defined by the International Commission on Illumination (abbreviated CIE) in
// 1976. It expresses color as three values: L* for perceptual lightness and
// a* and b* for the four unique colors of human vision: red, green, blue and
// yellow. CIELAB was intended as a perceptually uniform space, where a given
// numerical change corresponds to a similar perceived change in color.
type Lab struct {
	L float64 `json:"l"`
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// LabFuncE is the threshold for linear vs. nonlinear transition. [Reference]
//
// [Reference]: http://www.brucelindbloom.com/index.html?LContinuity.html
const LabFuncE float64 = 216.0 / 24389.0

// LabFuncK s the constant used for linear approximation. [Reference]
//
// [Reference]: http://www.brucelindbloom.com/index.html?LContinuity.html
const LabFuncK float64 = 24389.0 / 27.0

// NewLab creates Lab (CIELAB) color model
func NewLab(l, a, b float64) Lab {
	return Lab{l, a, b}
}

// Values returns L, a, b values of LABColor color
func (c Lab) Values() (float64, float64, float64) {
	return c.L, c.A, c.B
}

// ToARGB returns Color (ARGB) from LabColor
func (c Lab) ToARGB() ARGB {
	return c.ToXYZ().ToARGB()
}

// ToXYZ return XYZColor from LabColor
func (c Lab) ToXYZ() XYZ {
	l, a, b := c.Values()

	fy := (l + 16.0) / 116.0
	fx := a/500.0 + fy
	fz := fy - b/200.0

	// Normalizied x,y,z value from LabInvFunc (Lab inverse function)
	nx, ny, nz := LabInvFunc(fx), LabInvFunc(fy), LabInvFunc(fz)

	// White WhitePointD65
	wx, wy, wz := WhitePointD65.Values()

	// Denormalized value from WhitePointD65
	x, y, z := nx*wx, ny*wy, nz*wz
	return XYZ{x, y, z}
}

// LStar returns the L* value of L*a*b* (LabColor)
func (c Lab) LStar() float64 {
	return c.L
}

// LuminanceY returns the Y value for XYZ color model from Lab color model
func (c Lab) LuminanceY() float64 {
	return YFromLstar(c.L)
}

// DistanceSquared returns square of distance between two color
func (c Lab) DistanceSquared(b Lab) float64 {
	return c.L*b.L + c.A*b.A + c.B*b.B
}

// YFromLstar converts an L* (perceptual luminance) value from the CIELAB color
// space to Y (relative luminance) in the XYZ color space.
//
// Both L* and Y represent luminance, but L* is perceptually uniform and Y is
// linear.
//
// lstar is the L* value in the CIELAB color space.
// It returns the corresponding Y value in the XYZ color space.
func YFromLstar(lstar float64) float64 {
	return 100.0 * LabInvFunc((lstar+16.0)/116.0)
}

// LstarFromY converts Y (relative luminance) in the XYZ color space to L*
// (perceptual luminance) in the CIELAB color space.
//
// Both Y and L* represent luminance, but Y is linear and L* is perceptually
// uniform.
//
// y is the Y value in the XYZ color space.
// It returns the corresponding L* value in the CIELAB color space.
func LstarFromY(y float64) float64 {
	return 116.0*LabFunc(y/100.0) - 16.0
}

// LabFunc is part of the conversion from XYZ to Lab color space. It applies a
// nonlinear transformation that approximates human vision perception.
func LabFunc(t float64) float64 {
	if t > LabFuncE {
		return math.Cbrt(t)
	}
	return (LabFuncK*t + 16) / 116
}

// LabInvFunc is the inverse of LabFunc, used when converting from Lab to XYZ.
// It reverses the nonlinear transformation.
func LabInvFunc(ft float64) float64 {
	ft3 := ft * ft * ft
	if ft3 > LabFuncE {
		// If cube is above threshold, return it directly
		return ft3
	}
	// Otherwise, reverse the linear approximation
	return (116*ft - 16) / LabFuncK
}
