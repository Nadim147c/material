package color

func almostEqual[T float64 | float32](a, b T) bool {
	return (a - b) < 0.0001
}

func almostEqualColor[T XYZColor | LabColor](a, b T) bool {
	for i, c := range a {
		if !almostEqual(b[i], c) {
			return false
		}
	}
	return true
}

type ColorTestCase struct {
	Name string
	// Stored as 0xAARRGGBB
	ARGB Color
	HEX  string
	XYZ  [3]float64
	Lab  [3]float64
}

var ColorTestCases = []ColorTestCase{
	{
		Name: "White",
		ARGB: 0xFFFFFFFF,
		HEX:  "FFFFFF",
		XYZ:  [3]float64{0.9505, 1.0000, 1.0888},
		Lab:  [3]float64{100.00, 0.00, 0.00},
	},
	{
		Name: "Red",
		ARGB: 0xFFFF0000,
		HEX:  "FF0000",
		XYZ:  [3]float64{0.4125, 0.2127, 0.0193},
		Lab:  [3]float64{53.24, 80.09, 67.20},
	},
	{
		Name: "Black",
		ARGB: 0xFF000000,
		HEX:  "000000",
		XYZ:  [3]float64{0.0000, 0.0000, 0.0000},
		Lab:  [3]float64{0.00, 0.00, 0.00},
	},
	{
		Name: "Green",
		ARGB: 0xFF00FF00,
		HEX:  "00FF00",
		XYZ:  [3]float64{0.3576, 0.7152, 0.1192},
		Lab:  [3]float64{87.73, -86.18, 83.18},
	},
	{
		Name: "Gray 200",
		ARGB: 0xFFC8C8C8,
		HEX:  "C8C8C8",
		XYZ:  [3]float64{0.5490, 0.5776, 0.6289},
		Lab:  [3]float64{80.60, 0.00, 0.00},
	},
}
