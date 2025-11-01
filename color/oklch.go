package color

import (
	"fmt"
	"math"

	"github.com/Nadim147c/material/v2/num"
)

// OkLch represents a color in the OkLCh color space, a cylindrical
// transformation of OkLab designed for perceptual uniformity and
// intuitive manipulation of lightness, chroma, and hue.
//
// OkLCh provides a more perceptually uniform representation than RGB
// and improves color blending, interpolation, and prediction of hue
// and lightness. The model was introduced by Björn Ottosson in 2020.
type OkLch struct {
	// L is the perceptual lightness component, ranging from 0 (black)
	// to 100 (reference white).
	L float64
	// C is the chroma (color intensity or saturation) component,
	// ranging from 0 (neutral gray) to 100 (maximum vividness).
	C float64
	// H is the hue angle in degrees, ranging from 0 to 360, where
	// 0° = red, 120° = green, and 240° = blue.
	H float64
}

// NewOkLch create a OkLch model from l,c,h values
func NewOkLch(l, c, h float64) OkLch {
	return OkLch{l, c, h}
}

// OkLchFromOkLab create a OkLch model from a OkLab model
func OkLchFromOkLab(ok OkLab) OkLch {
	l, a, b := ok.Values()

	chroma := math.Hypot(a, b)

	theta := math.Atan2(b, a) // radian
	hue := num.NormalizeDegree(num.Degree(theta))
	return NewOkLch(l, chroma, hue)
}

// ToOkLab convert OkLch model to OkLab color model.
func (ok OkLch) ToOkLab() OkLab {
	l, chroma, hue := ok.Values()

	theta := num.Radian(num.NormalizeDegree(hue))

	a := chroma * math.Cos(theta)
	b := chroma * math.Sin(theta)
	return NewOkLab(l, a, b)
}

// ToXYZ convert OkLab model to XYZ color model.
func (ok OkLch) ToXYZ() XYZ {
	return ok.ToOkLab().ToXYZ()
}

// ToARGB convert OkLab model to ARGB color model.
func (ok OkLch) ToARGB() ARGB {
	return ok.ToXYZ().ToARGB()
}

// String returns a formatted string representation of OkLab color.
func (ok OkLch) String() string {
	return fmt.Sprintf("OkLch(%.4f, %.4f, %.4f)", ok.L, ok.C, ok.H)
}

// Values returns L, a, b values of OkLab Model
func (ok OkLch) Values() (float64, float64, float64) {
	return ok.L, ok.C, ok.H
}
