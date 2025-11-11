package color

import (
	"math"

	"github.com/Nadim147c/material/v2/num"
)

// LCHuv represents a color in the CIE LCHuv color space, a cylindrical
// transformation of the CIE 1976 (L*, u*, v*) color space (CIELUV). It defines
// a color using lightness (L), chroma (C), and hue (H), which provides a more
// intuitive representation for color manipulation.
type LCHuv struct {
	L float64 `json:"l"`
	C float64 `json:"c"`
	H float64 `json:"h"`
}

// NewLCHuv creates a new LuvLch color from the given lightness (L), chroma
// (C), and hue (H) values.
func NewLCHuv(l, c, h float64) LCHuv {
	return LCHuv{l, c, h}
}

// LchFromLuv converts a color from the CIELUV color space to the cylindrical
// LCHuv representation (LuvLch).
func LchFromLuv(c Luv) LCHuv {
	l, u, v := c.Values()
	chroma := math.Hypot(u, v)
	theta := math.Atan2(v, u) // radians
	hue := num.NormalizeDegree(num.Degree(theta))
	return NewLCHuv(l, chroma, hue)
}

// ToLuv converts the LCHuv color (LuvLch) back to the CIELUV color space.
func (c LCHuv) ToLuv() Luv {
	l, chroma, hue := c.Values()
	theta := num.Radian(num.NormalizeDegree(hue))
	u := chroma * math.Cos(theta)
	v := chroma * math.Sin(theta)
	return NewLuv(l, u, v)
}

// ToXYZ converts the LCHuv color to the CIE 1931 XYZ color space.
func (c LCHuv) ToXYZ() XYZ {
	return c.ToLuv().ToXYZ()
}

// ToARGB converts the LCHuv color to the ARGB color model.
func (c LCHuv) ToARGB() ARGB {
	return c.ToXYZ().ToARGB()
}

// Values returns the individual components (L, C, H) of the LCHuv color.
func (c LCHuv) Values() (float64, float64, float64) {
	return c.L, c.C, c.H
}
