package color

import (
	"math"

	"github.com/Nadim147c/material/v2/num"
)

// LCHab represents a color in the CIE LCHab color space, a cylindrical
// transformation of the CIE 1976 (L*, a*, b*) color space (CIELAB). It defines
// a color using lightness (L), chroma (C), and hue (H), providing a more
// perceptually uniform and intuitive model for color manipulation compared to
// the Cartesian Lab model.
type LCHab struct {
	L float64 `json:"l"`
	C float64 `json:"c"`
	H float64 `json:"h"`
}

var _ Model = (*LCHab)(nil)

// NewLCHab creates a new LCHab color from the given lightness (L),
// chroma (C), and hue (H) values.
func NewLCHab(l, c, h float64) LCHab {
	return LCHab{l, c, h}
}

// LchFromLab converts a color from the CIELAB color space to the
// cylindrical LCHab representation.
func LchFromLab(c Lab) LCHab {
	l, a, b := c.Values()
	chroma := math.Hypot(a, b)
	theta := math.Atan2(b, a) // radians
	hue := num.NormalizeDegree(num.Degree(theta))
	return NewLCHab(l, chroma, hue)
}

// ToLab converts the LCHab color back to the CIELAB color space.
func (c LCHab) ToLab() Lab {
	l, chroma, hue := c.Values()
	theta := num.Radian(num.NormalizeDegree(hue))
	a := chroma * math.Cos(theta)
	b := chroma * math.Sin(theta)
	return NewLab(l, a, b)
}

// ToXYZ converts the LCHab color to the CIE 1931 XYZ color space.
func (c LCHab) ToXYZ() XYZ {
	return c.ToLab().ToXYZ()
}

// ToARGB converts the LCHab color to the ARGB color model.
func (c LCHab) ToARGB() ARGB {
	return c.ToXYZ().ToARGB()
}

// String returns a formatted string representation of LCH(ab) color.
func (c LCHab) String() string {
	return modelString("LCHab", c)
}

// Values returns the individual components (L, C, H) of the LCHab color.
func (c LCHab) Values() (float64, float64, float64) {
	return c.L, c.C, c.H
}
