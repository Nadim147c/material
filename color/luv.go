package color

import "math"

// Luv is a color space adopted by the International Commission on Illumination
// (CIE) in 1976, as a simple-to-compute transformation of the 1931 CIE XYZ
// color space, which attempted perceptual uniformity. It is extensively used
// for applications such as computer graphics which deal with colored lights.
// Although additive mixtures of different colored lights will fall on a line in
// CIELUV's uniform chromaticity diagram (called the CIE 1976 UCS), such
// additive mixtures will not, contrary to popular belief, fall along a line in
// the CIELUV color space unless the mixtures are constant in lightness.
type Luv struct {
	L float64 `json:"l"`
	U float64 `json:"a"`
	V float64 `json:"b"`
}

// NewLuv creates the CIELUV color model
func NewLuv(l, u, v float64) Luv {
	return Luv{l, u, v}
}

var luvDPrimeRef = math.NaN()

// LuvFromXYZ converts from CIE XYZ to CIE L*u*v* color space.
func LuvFromXYZ(c XYZ) Luv {
	x, y, z := c.Values()
	Xr, Yr, Zr := WhitePointD65.Values()

	yr := y / Yr

	// Compute L*
	var L float64
	if yr > CieE {
		L = 116*math.Cbrt(yr) - 16
	} else {
		L = CieK * yr
	}

	dPrime := x + 15*y + 3*z
	// Handle degenerate case when d'=0 to avoid division by zero
	if dPrime == 0 {
		return NewLuv(L, 0, 0)
	}

	// calculate luvDPrimeRef only once by checking if the float is a NaN
	if luvDPrimeRef != luvDPrimeRef {
		luvDPrimeRef = Xr + 15*Yr + 3*Zr
	}

	// Compute u', v' for the sample and reference white
	uPrime := (4 * x) / dPrime
	vPrime := (9 * y) / dPrime
	uPrimeRef := (4 * Xr) / luvDPrimeRef
	vPrimeRef := (9 * Yr) / luvDPrimeRef

	// Compute u* and v*
	u := 13 * L * (uPrime - uPrimeRef)
	v := 13 * L * (vPrime - vPrimeRef)

	return NewLuv(L, u, v)
}

var (
	luvU0 = math.NaN()
	luvV0 = math.NaN()
)

// ToXYZ converts CIELUV to CIEXYZ
func (c Luv) ToXYZ() XYZ {
	l, u, v := c.Values()
	Xr, Yr, Zr := WhitePointD65.Values()

	// Reference white u0', v0'
	if luvU0 != luvU0 || luvV0 != luvV0 {
		luvU0 = (4 * Xr) / (Xr + 15*Yr + 3*Zr)
		luvV0 = (9 * Yr) / (Xr + 15*Yr + 3*Zr)
	}

	// Compute Y
	var Y float64
	if l > (CieK * CieE) {
		Y = math.Pow((l+16)/116, 3) * Yr
	} else {
		Y = (l / CieK) * Yr
	}

	// Handle degenerate case when L=0 to avoid division by zero
	if l == 0 {
		return NewXYZ(0, 0, 0)
	}

	const cPrime = -1.0 / 3.0

	// Intermediate variables from formula
	a := ((52 * l / (u + 13*l*luvU0)) - 1) / 3.0
	b := -5 * Y
	d := Y * ((39 * l / (v + 13*l*luvV0)) - 5)

	// Compute X and Z
	X := (d - b) / (a - cPrime)
	Z := X*a + b

	return NewXYZ(X, Y, Z)
}

// ToLCHuv converts CIELUV to LCHuv
func (c Luv) ToLCHuv() LCHuv {
	return LchFromLuv(c)
}

// ToARGB converts CIELUV to ARGB
func (c Luv) ToARGB() ARGB {
	return c.ToXYZ().ToARGB()
}

// Values returns L, U, C values
func (c Luv) Values() (float64, float64, float64) {
	return c.L, c.U, c.V
}
