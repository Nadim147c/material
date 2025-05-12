package color

import "math"

type LabColor [3]float64

const (
	// Threshold for linear vs. nonlinear transition. [Reference]
	//
	// [Reference]: http://www.brucelindbloom.com/index.html?LContinuity.html
	LabFuncE float64 = 216.0 / 24389.0
	// Constant used for linear approximation. [Reference]
	//
	// [Reference]: http://www.brucelindbloom.com/index.html?LContinuity.html
	LabFuncK float64 = 24389.0 / 27.0
)

func NewLabColor(l, a, b float64) LabColor {
	return LabColor{l, a, b}
}

// Values returns L, a, b values of LABColor color
func (c LabColor) Values() (float64, float64, float64) {
	return c[0], c[1], c[2]
}

// TODO: Implement c.ToARGB() Color Method

// labF is part of the conversion from XYZ to Lab color space.
// It applies a nonlinear transformation that approximates human vision
// perception.
func LabFunc(t float64) float64 {
	if t > LabFuncE {
		return math.Cbrt(t)
	}
	return (LabFuncK*t + 16) / 116
}

// labInvF is the inverse of labF, used when converting from Lab to XYZ. It
// reverses the nonlinear transformation.
func LabInvFunc(ft float64) float64 {
	ft3 := ft * ft * ft
	if ft3 > LabFuncE {
		// If cube is above threshold, return it directly
		return ft3
	}
	// Otherwise, reverse the linear approximation
	return (116*ft - 16) / LabFuncK
}
