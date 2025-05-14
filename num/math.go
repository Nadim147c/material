package num

import "cmp"

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

// Lerp is The linear interpolation function.
func Lerp(start float64, stop float64, amount float64) float64 {
	return (1.0-amount)*start + amount*stop
}
