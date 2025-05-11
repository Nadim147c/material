package color

import "math"

type LabColor [3]float64

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
func LabF(t float64) float64 {
	// Threshold for linear vs. nonlinear transition (~0.008856)
	e := 216.0 / 24389.0
	// Constant used for linear approximation (~903.3)
	kappa := 24389.0 / 27.0

	if t > e {
		// For values above the threshold, use the cube root
		return math.Cbrt(t)
	}
	// For lower values, use the linear approximation
	return (kappa*t + 16) / 116
}

// labInvF is the inverse of labF, used when converting from Lab to XYZ. It
// reverses the nonlinear transformation.
func LabInvF(ft float64) float64 {
	e := 216.0 / 24389.0    // Same threshold as in labF
	kappa := 24389.0 / 27.0 // Same constant

	ft3 := ft * ft * ft // Cube of the input

	if ft3 > e {
		// If cube is above threshold, return it directly
		return ft3
	}
	// Otherwise, reverse the linear approximation
	return (116*ft - 16) / kappa
}
