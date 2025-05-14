package color

import (
	"math"

	"github.com/Nadim147c/goyou/num"
)

var (
	Cat16Matrix = num.NewMatrix3(
		0.401288, 0.650173, -0.051461,
		-0.250268, 1.204414, 0.045854,
		-0.002079, 0.048952, 0.953127,
	)

	Cat16InvMatrix = num.NewMatrix3(
		1.86206786, -1.01125463, 0.14918678,
		0.38752654, 0.62144744, -0.00897399,
		-0.0158415, -0.03412294, 1.04996444,
	)
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

func NewCam16(hue, chroma, j, q, m, s, jstar, astar, bstar float64) *Cam16 {
	return &Cam16{hue, chroma, j, q, m, s, jstar, astar, bstar}
}

func Cam16FromColor(color Color) *Cam16 {
	return Cam16FromColorInViewingCondition(color, &DefaultViewingConditions)
}

func Cam16FromColorInViewingCondition(color Color, vc *ViewingConditions) *Cam16 {
	// Get XYZ color model
	x, y, z := color.ToXYZ().Values()

	// Convert XYZ to 'cone'/'rgb' responses
	rC, gC, bC := Cat16Matrix.MultiplyXYZ(x, y, z).Values()

	// RGBD of viewing condition
	rD, gD, bD := vc.RgbD.Values()

	// Discount illuminant.
	rD *= rC
	gD *= gC
	bD *= bC

	// Chromatic adaptation.
	rAF := math.Pow((vc.Fl*math.Abs(rD))/100.0, 0.42)
	gAF := math.Pow((vc.Fl*math.Abs(gD))/100.0, 0.42)
	bAF := math.Pow((vc.Fl*math.Abs(bD))/100.0, 0.42)
	rA := num.SignNum(rD) * 400.0 * rAF / (rAF + 27.13)
	gA := num.SignNum(gD) * 400.0 * gAF / (gAF + 27.13)
	bA := num.SignNum(bD) * 400.0 * bAF / (bAF + 27.13)

	// Redness-greenness
	a := (11*rA + -12*gA + bA) / 11
	b := (rA + gA - 2*bA) / 9
	u := (20*rA + 20*gA + 21*bA) / 20
	p2 := (40*rA + 20*gA + bA) / 20

	radians := math.Atan2(b, a)
	degrees := num.Degree(radians)
	hue := num.NormalizeAngle(degrees)
	hueRadians := num.Radian(hue)
	ac := p2 * vc.Nbb

	j := 100 * math.Pow(ac/vc.Aw, vc.C*vc.Z)
	q := (4 / vc.C) * math.Sqrt(j/100) * (vc.Aw + 4) * vc.FlRoot

	huePrime := hue
	if hue < 20.14 {
		huePrime = hue + 360
	}
	eHue := 0.25 * (math.Cos((huePrime*math.Pi)/180.0+2.0) + 3.8)

	p1 := (50000.0 / 13.0) * eHue * vc.Nc * vc.Ncb
	t := (p1 * math.Sqrt(a*a+b*b)) / (u + 0.305)

	alpha := math.Pow(t, 0.9) * math.Pow(1.64-math.Pow(0.29, vc.N), 0.73)

	chroma := alpha * math.Sqrt(j/100.0)
	m := chroma * vc.FlRoot
	s := 50.0 * math.Sqrt((alpha*vc.C)/(vc.Aw+4.0))

	jstar := ((1.0 + 100.0*0.007) * j) / (1.0 + 0.007*j)
	mstar := (1.0 / 0.0228) * math.Log(1.0+0.0228*m)
	astar := mstar * math.Cos(hueRadians)
	bstar := mstar * math.Sin(hueRadians)

	return NewCam16(hue, chroma, j, q, m, s, jstar, astar, bstar)
}

func (c *Cam16) ToColor() Color {
	return c.Viewed(&DefaultViewingConditions)
}

// Viewed converts a CAM16 color to an ARGB integer based on
// the given viewing conditions
func (c *Cam16) Viewed(vc *ViewingConditions) Color {
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
	rC := num.SignNum(rA) * (100.0 / vc.Fl) *
		math.Pow(rCBase, 1.0/0.42)
	gCBase := math.Max(0, (27.13*math.Abs(gA))/(400.0-math.Abs(gA)))
	gC := num.SignNum(gA) * (100.0 / vc.Fl) *
		math.Pow(gCBase, 1.0/0.42)
	bCBase := math.Max(0, (27.13*math.Abs(bA))/(400.0-math.Abs(bA)))
	bC := num.SignNum(bA) * (100.0 / vc.Fl) *
		math.Pow(bCBase, 1.0/0.42)

	rF := rC / vc.RgbD[0]
	gF := gC / vc.RgbD[1]
	bF := bC / vc.RgbD[2]

	x, y, z := Cat16InvMatrix.MultiplyXYZ(rF, gF, bF).Values()
	return FromXYZ(x, y, z)
}
