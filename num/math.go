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
