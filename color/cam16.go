package color

import (
	"math"

	"github.com/Nadim147c/material/v2/num"
)

// Cat16Matrix is the forward CAT16 (Chromatic Adaptation Transform) matrix.
//
// It converts linear RGB values into the CAM16 LMS (Long, Medium, Short)
// cone response domain. This step models how the human visual system
// adapts to different white points and viewing conditions.
//
// The CAT16 matrix is part of the CAM16 color appearance model used
// by the HCT color system.
var Cat16Matrix = num.NewMatrix3(
	0.401288, 0.650173, -0.051461,
	-0.250268, 1.204414, 0.045854,
	-0.002079, 0.048952, 0.953127,
)

// Cat16InvMatrix is the inverse of Cat16Matrix.
//
// It converts CAM16 LMS cone responses back into linear RGB values,
// reversing the chromatic adaptation process. This matrix is used
// when converting from the CAM16/HCT perceptual space back to RGB.
var Cat16InvMatrix = num.NewMatrix3(
	1.86206786, -1.01125463, 0.14918678,
	0.38752654, 0.62144744, -0.00897399,
	-0.0158415, -0.03412294, 1.04996444,
)

// Cam16 represents the CAM16 color model, which includes various dimensions
// for color representation. It can be constructed using any combination of
// three of the following dimensions: j or q, c, m, or s, and hue.
// It can also be constructed using the CAM16-UCS J, a, and b coordinates.
type Cam16 struct {
	// Hue represents the hue of the color
	Hue float64
	// Chroma represents the colorfulness or color intensity, similar to
	// saturation in HSL but more perceptually accurate
	Chroma float64
	// J represents the lightness of the color
	J float64
	// Q represents the brightness, which is the ratio of lightness to
	// the white point's lightness
	Q float64
	// M represents the colorfulness of the color
	M float64
	// S represents the saturation, which is the ratio of chroma to
	// the white point's chroma
	S float64
	// Jstar represents the CAM16-UCS J coordinate
	Jstar float64
	// Astar represents the CAM16-UCS a coordinate
	Astar float64
	// Bstar represents the CAM16-UCS b coordinate
	Bstar float64
}

// NewCam16 create a CAM16 color model from given values
func NewCam16(hue, chroma, j, q, m, s, jstar, astar, bstar float64) *Cam16 {
	return &Cam16{hue, chroma, j, q, m, s, jstar, astar, bstar}
}

// Cam16FromXyzInEnv create a Cam16 color In specific ViewingConditions
func Cam16FromXyzInEnv(xyz XYZ, env *Environmnet) *Cam16 {
	// Get XYZ color model
	x, y, z := xyz.Values()

	// Convert XYZ to 'cone'/'rgb' responses
	rC, gC, bC := Cat16Matrix.MultiplyXYZ(x, y, z).Values()

	// RGBD of viewing condition
	rD, gD, bD := env.RgbD.Values()

	// Discount illuminant.
	rD *= rC
	gD *= gC
	bD *= bC

	// Chromatic adaptation.
	rAF := math.Pow((env.Fl*math.Abs(rD))/100.0, 0.42)
	gAF := math.Pow((env.Fl*math.Abs(gD))/100.0, 0.42)
	bAF := math.Pow((env.Fl*math.Abs(bD))/100.0, 0.42)
	rA := num.Sign(rD) * 400.0 * rAF / (rAF + 27.13)
	gA := num.Sign(gD) * 400.0 * gAF / (gAF + 27.13)
	bA := num.Sign(bD) * 400.0 * bAF / (bAF + 27.13)

	// Redness-greenness
	a := (11*rA + -12*gA + bA) / 11
	b := (rA + gA - 2*bA) / 9
	u := (20*rA + 20*gA + 21*bA) / 20
	p2 := (40*rA + 20*gA + bA) / 20

	radians := math.Atan2(b, a)
	degrees := num.Degree(radians)
	hue := num.NormalizeDegree(degrees)
	hueRadians := num.Radian(hue)
	ac := p2 * env.Nbb

	j := 100 * math.Pow(ac/env.Aw, env.C*env.Z)
	q := (4 / env.C) * math.Sqrt(j/100) * (env.Aw + 4) * env.FlRoot

	huePrime := hue
	if hue < 20.14 {
		huePrime = hue + 360
	}
	eHue := 0.25 * (math.Cos((huePrime*math.Pi)/180.0+2.0) + 3.8)

	p1 := (50000.0 / 13.0) * eHue * env.Nc * env.Ncb
	t := (p1 * math.Sqrt(a*a+b*b)) / (u + 0.305)

	alpha := math.Pow(t, 0.9) * math.Pow(1.64-math.Pow(0.29, env.N), 0.73)

	chroma := alpha * math.Sqrt(j/100.0)
	m := chroma * env.FlRoot
	s := 50.0 * math.Sqrt((alpha*env.C)/(env.Aw+4.0))

	jstar := ((1.0 + 100.0*0.007) * j) / (1.0 + 0.007*j)
	mstar := (1.0 / 0.0228) * math.Log(1.0+0.0228*m)
	astar := mstar * math.Cos(hueRadians)
	bstar := mstar * math.Sin(hueRadians)

	return NewCam16(hue, chroma, j, q, m, s, jstar, astar, bstar)
}

// Cam16FromJch constructs a Cam16 color from J (lightness), C (chroma), and
// H (hue angle in degrees), using DefaultViewingConditions viewing conditions.
//
// This is used when synthesizing a CAM16 color from HCT values or
// performing color space conversions into perceptual models.
func Cam16FromJch(j, c, h float64) *Cam16 {
	return Cam16FromJchInEnv(j, c, h, &DefaultEnviroment)
}

// Cam16FromJchInEnv constructs a Cam16 color from J (lightness), C (chroma),
// and H (hue angle in degrees), using the given viewing conditions.
//
// This is used when synthesizing a CAM16 color from HCT values or
// performing color space conversions into perceptual models.
func Cam16FromJchInEnv(j, c, h float64, env *Environmnet) *Cam16 {
	q := (4.0 / env.C) * math.Sqrt(j/100.0) * (env.Aw + 4.0) * env.FlRoot
	m := c * env.FlRoot

	alpha := c / math.Sqrt(j/100.0)
	s := 50.0 * math.Sqrt((alpha*env.C)/(env.Aw+4.0))

	hueRadians := h * math.Pi / 180.0

	jstar := ((1.0 + 100.0*0.007) * j) / (1.0 + 0.007*j)

	mstar := (1.0 / 0.0228) * math.Log(1.0+0.0228*m)
	astar := mstar * math.Cos(hueRadians)
	bstar := mstar * math.Sin(hueRadians)

	return NewCam16(h, c, j, q, m, s, jstar, astar, bstar)
}

// Viewed converts a CAM16 color to an ARGB integer based on
// the given viewing conditions
func (c *Cam16) Viewed(vc *Environmnet) XYZ {
	var alpha float64
	if c.Chroma == 0.0 || c.J == 0.0 {
		alpha = 0.0
	} else {
		alpha = c.Chroma / math.Sqrt(c.J/100.0)
	}

	t := math.Pow(alpha/math.Pow(1.64-math.Pow(0.29, vc.N), 0.73), 1.0/0.9)
	hRad := (c.Hue * math.Pi) / 180.0

	eHue := 0.25 * (math.Cos(hRad+2.0) + 3.8)
	ac := vc.Aw * math.Pow(c.J/100.0, 1.0/vc.C/vc.Z)
	p1 := eHue * (50000.0 / 13.0) * vc.Nc * vc.Ncb
	p2 := ac / vc.Nbb

	hSin := math.Sin(hRad)
	hCos := math.Cos(hRad)

	gamma := (23.0 * (p2 + 0.305) * t) / (23.0*p1 + 11.0*t*hCos + 108.0*t*hSin)

	a := gamma * hCos
	b := gamma * hSin

	rA := (460.0*p2 + 451.0*a + 288.0*b) / 1403.0
	gA := (460.0*p2 - 891.0*a - 261.0*b) / 1403.0
	bA := (460.0*p2 - 220.0*a - 6300.0*b) / 1403.0

	rCBase := math.Max(0, (27.13*math.Abs(rA))/(400.0-math.Abs(rA)))
	rC := num.Sign(rA) * (100.0 / vc.Fl) *
		math.Pow(rCBase, 1.0/0.42)
	gCBase := math.Max(0, (27.13*math.Abs(gA))/(400.0-math.Abs(gA)))
	gC := num.Sign(gA) * (100.0 / vc.Fl) *
		math.Pow(gCBase, 1.0/0.42)
	bCBase := math.Max(0, (27.13*math.Abs(bA))/(400.0-math.Abs(bA)))
	bC := num.Sign(bA) * (100.0 / vc.Fl) *
		math.Pow(bCBase, 1.0/0.42)

	rF := rC / vc.RgbD[0]
	gF := gC / vc.RgbD[1]
	bF := bC / vc.RgbD[2]

	x, y, z := Cat16InvMatrix.MultiplyXYZ(rF, gF, bF).Values()
	return XYZ{x, y, z}
}

// Cam16FromUcs creates a CAM16 color from UCS coordinates (jstar, astar,
// bstar).
// Uses the default viewing environment for conversion.
func Cam16FromUcs(jstar, astar, bstar float64) *Cam16 {
	return Cam16FromUcsInEnv(jstar, astar, bstar, &DefaultEnviroment)
}

// Cam16FromUcsInEnv creates a CAM16 color from UCS coordinates (jstar, astar,
// bstar)
// using the specified viewing environment for conversion.
func Cam16FromUcsInEnv(jstar, astar, bstar float64, env *Environmnet) *Cam16 {
	a := astar
	b := bstar
	m := math.Sqrt(a*a + b*b)
	M := (math.Exp(m*0.0228) - 1.0) / 0.0228
	c := M / env.FlRoot
	h := math.Atan2(b, a) * (180.0 / math.Pi)
	h = num.NormalizeDegree(h)
	j := jstar / (1 - (jstar-100)*0.007)
	return Cam16FromJchInEnv(j, c, h, env)
}

// ToHct converts the CAM16 color to HCT (Hue, Chroma, Tone) color space.
func (c *Cam16) ToHct() Hct {
	return NewHct(c.Hue, c.Chroma, c.J)
}

// ToXYZ converts the CAM16 color to CIE XYZ color space.
// Uses the default viewing environment for conversion.
func (c *Cam16) ToXYZ() XYZ {
	return c.Viewed(&DefaultEnviroment)
}

// ToLab converts the CAM16 color to CIE L*a*b* color space.
// Uses the default viewing environment for conversion.
func (c *Cam16) ToLab() Lab {
	return c.Viewed(&DefaultEnviroment).ToLab()
}

// ToARGB converts the CAM16 color to ARGB format.
// Uses the default viewing environment for conversion.
func (c *Cam16) ToARGB() ARGB {
	return c.Viewed(&DefaultEnviroment).ToARGB()
}

//revive:disable:function-result-limit

// RGBA implements the color.Color interface. Returns the red, green, blue, and
// alpha values in the 0-65535 range.
func (c *Cam16) RGBA() (red uint32, green uint32, blue uint32, alpha uint32) {
	return c.ToARGB().RGBA()
}

//revive:enable:function-result-limit

// Distance returns distance between to Cam16 color
func (c Cam16) Distance(other Cam16) float64 {
	dJ := c.Jstar - other.Jstar
	dA := c.Astar - other.Astar
	dB := c.Bstar - other.Bstar

	dEPrime := math.Sqrt(dJ*dJ + dA*dA + dB*dB)
	dE := 1.41 * math.Pow(dEPrime, 0.63)
	return dE
}
