package color

func almostEqual[T float64 | float32](a, b T) bool {
	return (a - b) < 0.001
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
	XYZ  XYZColor
	Lab  LabColor
}

var ColorTestCases = []ColorTestCase{
	{
		Name: "White",
		ARGB: 0xFFFFFFFF,
		HEX:  "FFFFFF",
		XYZ:  [3]float64{95.047, 100.00, 108.88},
		Lab:  [3]float64{100.00, 0.00, 0.00},
	},
	{
		Name: "Red",
		ARGB: 0xFFFF0000,
		HEX:  "FF0000",
		XYZ:  [3]float64{41.245, 21.267, 1.93},
		Lab:  [3]float64{53.24, 80.09, 67.20},
	},
	{
		Name: "Black",
		ARGB: 0xFF000000,
		HEX:  "000000",
		XYZ:  [3]float64{0, 0, 0},
		Lab:  [3]float64{0.00, 0.00, 0.00},
	},
	{
		Name: "Green",
		ARGB: 0xFF00FF00,
		HEX:  "00FF00",
		XYZ:  [3]float64{35.757, 71.515, 11.9192},
		Lab:  [3]float64{87.73461, -86.18431, 83.1791},
	},
	{
		Name: "Gray 200",
		ARGB: 0xFFC8C8C8,
		HEX:  "C8C8C8",
		XYZ:  [3]float64{54.8972, 57.75804, 62.8886},
		Lab:  [3]float64{80.60408, 0.00, 0.00},
	},
}
