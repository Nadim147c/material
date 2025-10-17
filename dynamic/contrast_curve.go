package dynamic

import "github.com/Nadim147c/material/num"

// ContrastCurve represents a curve that provides contrast values based on
// contrast level
type ContrastCurve struct {
	low, normal, medium, high float64
}

// NewContrastCurve creates a new NewContrastCurve with the specified values
// low: Value for contrast level -1.0
// normal: Value for contrast level 0.0
// medium: Value for contrast level 0.5
// high: Value for contrast level 1.0
func NewContrastCurve(low, normal, medium, high float64) *ContrastCurve {
	return &ContrastCurve{low, normal, medium, high}
}

// Get returns the value at a given contrast level
// contrast: The contrast level. 0.0 is the default (normal); -1.0 is the
// lowest; 1.0 is the highest.
// return: The value. For contrast ratios, a number between 1.0 and 21.0.
func (c *ContrastCurve) Get(contrast float64) float64 {
	if contrast <= -1.0 {
		return c.low
	} else if contrast < 0.0 {
		return num.Lerp(c.low, c.normal, (contrast-(-1))/1)
	} else if contrast < 0.5 {
		return num.Lerp(c.normal, c.medium, contrast/0.5)
	} else if contrast < 1.0 {
		return num.Lerp(c.medium, c.high, (contrast-0.5)/0.5)
	}
	return c.high
}
