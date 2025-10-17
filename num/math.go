package num

import (
	"cmp"
	"math"
)

// Clamp takes in a value and two thresholds. If the value is smaller than the
// low threshold, it returns the low threshold. If it's bigger than the high
// threshold it returns the high threshold. Otherwise it returns the value.
func Clamp[T cmp.Ordered](low, high, value T) T {
	switch {
	case value < low:
		return low
	case value > high:
		return high
	default:
		return value
	}
}

// SignCmp compares two ordered values a and b.
// It returns -1 if a < b, 1 if a > b, and 0 if a == b.
func SignCmp[T cmp.Ordered](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float64 | ~float32
}

// Sign compares two ordered values n.
// It returns -1 if n < 0, 1 if n > 0, and 0 if n == 0.
func Sign[T signed](n T) T {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	default:
		return 0
	}
}

// Lerp is The linear interpolation function.
func Lerp(start float64, stop float64, amount float64) float64 {
	return (1.0-amount)*start + amount*stop
}

// NormalizeDegree takes an angle in degrees and normalizes it to the range
// 0-360.
func NormalizeDegree(angle float64) float64 {
	normalized := math.Mod(angle, 360)
	if normalized < 0 {
		normalized += 360
	}
	return normalized
}

// NormalizeDegreeInt takes an angle in degrees and normalizes it to the range
// 0-360.
func NormalizeDegreeInt(angle int) int {
	normalized := angle % 360
	if normalized < 0 {
		normalized += 360
	}
	return normalized
}

// NormalizeRadian takes an angle in degrees and normalizes it to the range
// 0-360.
func NormalizeRadian(angle float64) float64 {
	twoPi := 2 * math.Pi
	normalized := math.Mod(angle, twoPi)
	if normalized < 0 {
		normalized += twoPi
	}
	return normalized
}

// Radian converts an angle in degrees to radians.
func Radian(deg float64) float64 {
	return (deg * math.Pi) / 180
}

// Degree converts an angle in radians to degrees.
func Degree(rad float64) float64 {
	return (rad * 180) / math.Pi
}

// RotationDirection calculates the optimal rotation direction between two
// angles.
//
// Given two angles 'from' and 'to' in degrees, it returns:
//
//	-1.0 for clockwise rotation (shorter path)
//	 1.0 for counter-clockwise rotation (shorter path)
//	 0.0 if the angles are identical
//
// The function considers all three possible paths (direct, +360°, -360°) and
// chooses the one with the smallest absolute angular distance.
func RotationDirection(from float64, to float64) float64 {
	a := to - from
	b := to - from + 360.0
	c := to - from - 360.0

	aAbs, bAbs, cAbs := math.Abs(a), math.Abs(b), math.Abs(c)
	if aAbs <= bAbs && aAbs <= cAbs {
		return Sign(a)
	} else if bAbs <= aAbs && bAbs <= cAbs {
		return Sign(b)
	}
	return Sign(c)
}

// DifferenceDegrees returns the shortest angular difference between two angles
// in degrees.
func DifferenceDegrees(a, b float64) float64 { return NormalizeDegree(a - b) }
