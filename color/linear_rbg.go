package color

import "github.com/Nadim147c/material/v2/num"

// LinearRGB represents a color in the linear RGB color space.
//
// Each component (R, G, and B) is expressed as a float64 value in the range [0,
// 100], where 0 corresponds to no intensity and 100 corresponds to full
// intensity. Unlike gamma-corrected RGB, LinearRGB stores values in a linear
// light space, making it suitable for accurate color calculations and
// conversions.
type LinearRGB struct {
	R, G, B float64
}

// NewLinearRGB creates a linear RGB color model
func NewLinearRGB(r, g, b float64) LinearRGB {
	return LinearRGB{r, g, b}
}

// LinearRGBFromARGB to converts ARGB (sRGB) to linear RGB.
func LinearRGBFromARGB(c ARGB) LinearRGB {
	r, g, b := c.Red(), c.Green(), c.Blue()
	return NewLinearRGB(Linearized3(r, g, b))
}

// LinearRGBFromXYZ to converts ARGB (sRGB) to linear RGB.
func LinearRGBFromXYZ(c XYZ) LinearRGB {
	vec := num.NewVector(c)
	r, g, b := XYZ_TO_RGB.Multiply(vec).Values()
	return NewLinearRGB(r, g, b)
}

// ToARGB converts LinearRGB to ARGB.
func (c LinearRGB) ToARGB() ARGB {
	return ARGBFromLinearRGB(c.Values())
}

// ToXYZ converts LinearRGB to ARGB.
func (c LinearRGB) ToXYZ() XYZ {
	vec := num.NewVector(c)
	xyz := RGB_TO_XYZ.Multiply(vec)
	return NewXYZ(xyz.Values())
}

// Values returns R, G, B components of LinearRGB.
func (c LinearRGB) Values() (float64, float64, float64) {
	return c.R, c.G, c.B
}
