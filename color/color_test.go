package color

import "math"

func almostEqual[T float64 | float32](a, b T) bool {
	return math.Abs(float64(a-b)) <= 0.01
}

func sameXYZ(a, b XYZ) bool {
	return almostEqual(a.X, b.X) && almostEqual(a.Y, b.Y) && almostEqual(a.Z, b.Z)
}

func sameLab(a, b Lab) bool {
	return almostEqual(a.L, b.L) && almostEqual(a.A, b.A) && almostEqual(a.B, b.B)
}

type ColorTestCase struct {
	Name string
	// Stored as 0xAARRGGBB
	ARGB ARGB
	HEX  string
	XYZ  XYZ
	Lab  Lab
}

var ColorTestCases = []ColorTestCase{
	{
		Name: "White",
		ARGB: 0xFFFFFFFF,
		HEX:  "FFFFFF",
		XYZ:  XYZ{95.047, 100.00, 108.88},
		Lab:  Lab{100.00, 0.00, 0.00},
	},
	{
		Name: "Red",
		ARGB: 0xFFFF0000,
		HEX:  "FF0000",
		XYZ:  XYZ{41.233, 21.26, 1.932},
		Lab:  Lab{53.232, 80.087, 67.202},
	},
	{
		Name: "Black",
		ARGB: 0xFF000000,
		HEX:  "000000",
		XYZ:  XYZ{0, 0, 0},
		Lab:  Lab{0.00, 0.00, 0.00},
	},
	{
		Name: "Green",
		ARGB: 0xFF00FF00,
		HEX:  "00FF00",
		XYZ:  XYZ{35.757, 71.515, 11.9192},
		Lab:  Lab{87.73461, -86.18431, 83.1791},
	},
	{
		Name: "Gray 200",
		ARGB: 0xFFC8C8C8,
		HEX:  "C8C8C8",
		XYZ:  XYZ{54.8972, 57.75804, 62.8886},
		Lab:  Lab{80.60408, 0.00, 0.00},
	},
}
