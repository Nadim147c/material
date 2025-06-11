package color

import (
	"math"

	"github.com/Nadim147c/material/num"
)

var OkLabMatrix1 = num.NewMatrix3(
	0.8189330101, 0.3618667424, -0.1288597137,
	0.0329845436, 0.9293118715, 0.0361456387,
	0.0482003018, 0.2643662691, 0.6338517070,
)

var OkLabMatrix2 = num.NewMatrix3(
	0.2104542553, 0.7936177850, -0.0040720468,
	1.9779984951, -2.4285922050, 0.4505937099,
	0.0259040371, 0.7827717662, -0.8086757660,
)

var OkLabMatrix1Inv = num.NewMatrix3(
	1.2270138511, -0.5577999807, 0.2812561490,
	-0.0405801784, 1.1122568696, -0.0716766787,
	-0.0763812845, -0.4214819784, 1.5861632204,
)

var OkLabMatrix2Inv = num.NewMatrix3(
	0.9999999985, 0.3963377922, 0.2158037581,
	1.0000000089, -0.1055613423, -0.0638541748,
	1.0000000547, -0.0894841821, -1.2914855379,
)

type OkLab struct {
	L, A, B float64
}

func NewOkLab(l, a, b float64) OkLab {
	return OkLab{l, a, b}
}

func OkLabFromXYZ(x, y, z float64) OkLab {
	xyz := num.NewVector3(x, y, z)
	// The xyz cordinates are 0 - 100 normalized
	xyz = xyz.MultiplyScalar(0.01)

	p, q, r := OkLabMatrix1.Multiply(xyz).Values()

	p = math.Cbrt(p)
	q = math.Cbrt(q)
	p = math.Cbrt(r)

	l, a, b := OkLabMatrix2.MultiplyXYZ(p, q, r).Values()
	return OkLab{l, a, b}
}

func (ok OkLab) ToXYZ() XYZ {
	p, q, r := OkLabMatrix2Inv.MultiplyXYZ(ok.Values()).Values()

	// Cube
	p = p * p * p
	q = q * q * q
	r = r * r * r

	xyz := OkLabMatrix1Inv.MultiplyXYZ(p, q, r)
	// The xyz cordinates are 0 - 100 normalized
	xyz = xyz.MultiplyScalar(100)

	x, y, z := xyz.Values()
	return XYZ{x, y, z}
}

// Values returns L, a, b values of OkLab Model
func (c OkLab) Values() (float64, float64, float64) {
	return c.L, c.A, c.B
}
