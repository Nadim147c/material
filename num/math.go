package num

import (
	"cmp"
	"math"

	"golang.org/x/exp/constraints"
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

// SignNum compares two ordered values n.
// It returns -1 if n < 0, 1 if n > 0, and 0 if n == 0.
func SignNum[T constraints.Signed | float64 | float32](n T) T {
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

// NormalizeAngle takes an angle in degrees and normalizes it to the range 0-360
func NormalizeAngle(angle float64) float64 {
	normalized := math.Mod(angle, 360)
	if normalized < 0 {
		normalized += 360
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
