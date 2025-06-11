package contrast

import (
	"math"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/num"
)

// RatioOfTones returns the contrast ratio between two tones, in the range
// [1, 21].
//
// toneA and toneB should be values between 0 and 100. Values outside this range
// will be clamped internally.
func RatioOfTones(toneA, toneB float64) float64 {
	toneA = num.Clamp(0, 100, toneA)
	toneB = num.Clamp(0, 100, toneB)
	return RatioOfYs(color.YFromLstar(toneA), color.YFromLstar(toneB))
}

// RatioOfYs returns the contrast ratio between two relative luminance values.
//
// The ratio is calculated as (lighter + 5) / (darker + 5).
func RatioOfYs(y1, y2 float64) float64 {
	lighter := max(y1, y2)
	darker := y2
	if lighter == y2 {
		darker = y1
	}
	return (lighter + 5.0) / (darker + 5.0)
}

// Lighter returns a tone greater than or equal to the given tone that satisfies
// the specified contrast ratio.
//
// The returned value is in the range [0, 100]. Returns -1 if the desired
// contrast ratio cannot be achieved with the input tone.
//
// tone must be in the range [0, 100]. Invalid values result in -1.
// ratio should be in the range [1, 21]; invalid values have undefined behavior.
func Lighter(tone float64, ratio float64) float64 {
	if tone < 0 || tone > 100 {
		return -1
	}

	darkY := color.YFromLstar(tone)
	lightY := ratio*(darkY+5.0) - 5.0
	realContrast := RatioOfYs(lightY, darkY)
	delta := math.Abs(realContrast - ratio)
	if realContrast < ratio && delta > 0.04 {
		return -1
	}

	// Ensure gamut mapping, which requires a 'range' on tone, will still result
	// in the correct ratio by brightening slightly.
	returnValue := color.LstarFromY(lightY) + 0.4
	if returnValue < 0 || returnValue > 100 {
		return -1
	}
	return returnValue
}

// Darker returns a tone less than or equal to the given tone that satisfies the
// specified contrast ratio.
//
// The returned value is in the range [0, 100]. Returns -1 if the desired
// contrast ratio cannot be achieved with the input tone.
//
// tone must be in the range [0, 100]. Invalid values result in -1.
// ratio should be in the range [1, 21]; invalid values have undefined behavior.
func Darker(tone, ratio float64) float64 {
	if tone < 0.0 || tone > 100.0 {
		return -1.0
	}

	lightY := color.YFromLstar(tone)
	darkY := ((lightY + 5.0) / ratio) - 5.0
	realContrast := RatioOfYs(lightY, darkY)

	delta := math.Abs(realContrast - ratio)
	if realContrast < ratio && delta > 0.04 {
		return -1
	}

	// Ensure gamut mapping, which requires a 'range' on tone, will still result
	// in the correct ratio by darkening slightly.
	returnValue := color.LstarFromY(darkY) - 0.4
	if returnValue < 0 || returnValue > 100 {
		return -1
	}
	return returnValue
}

// LighterUnsafe returns a tone greater than or equal to the given tone that
// attempts to satisfy the specified contrast ratio.
//
// The returned value is always in the range [0, 100]. Returns 100 if the
// desired contrast ratio cannot be achieved. Unlike Lighter, this function does
// not fail but may not always satisfy the desired contrast ratio.
func LighterUnsafe(tone, ratio float64) float64 {
	lighterSafe := Lighter(tone, ratio)
	if lighterSafe < 0 {
		return 100
	}
	return lighterSafe
}

// DarkerUnsafe returns a tone less than or equal to the given tone that
// attempts to satisfy the specified contrast ratio.
//
// The returned value is always in the range [0, 100]. Returns 0 if the desired
// contrast ratio cannot be achieved. Unlike Darker, this function does not fail
// but may not always satisfy the desired contrast ratio.
func DarkerUnsafe(tone, ratio float64) float64 {
	darkerSafe := Darker(tone, ratio)
	if darkerSafe < 0 {
		return 0
	}
	return darkerSafe
}
