package dynamic

import "github.com/Nadim147c/goyou/num"

// contrastCurve represents a curve that provides contrast values based on contrast level
type contrastCurve struct {
	low, normal, medium, high float64
}

// ContrastCurve creates a new ContrastCurve with the specified values
// low: Value for contrast level -1.0
// normal: Value for contrast level 0.0
// medium: Value for contrast level 0.5
// high: Value for contrast level 1.0
func ContrastCurve(low, normal, medium, high float64) *contrastCurve {
	return &contrastCurve{low, normal, medium, high}
}

// Get returns the value at a given contrast level
// contrastLevel: The contrast level. 0.0 is the default (normal); -1.0 is the lowest; 1.0 is the highest.
// return: The value. For contrast ratios, a number between 1.0 and 21.0.
func (c *contrastCurve) Get(contrastLevel float64) float64 {
	if contrastLevel <= -1.0 {
		return c.low
	} else if contrastLevel < 0.0 {
		return num.Lerp(c.low, c.normal, (contrastLevel-(-1))/1)
	} else if contrastLevel < 0.5 {
		return num.Lerp(c.normal, c.medium, contrastLevel/0.5)
	} else if contrastLevel < 1.0 {
		return num.Lerp(c.medium, c.high, (contrastLevel-0.5)/0.5)
	} else {
		return c.high
	}
}
