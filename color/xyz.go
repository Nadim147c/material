package color

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
	lr, lg, lb := XYZ_TO_SRGB.MultiplyXYZ(x, y, z).Values()
	r, g, b := Delinearized(lr), Delinearized(lg), Delinearized(lb)
	return FromRGB(r, g, b)
}
