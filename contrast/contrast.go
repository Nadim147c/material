package contrast

import (
	"math"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/num"
)

// RatioOfTones calculates the contrast ratio between two tones.
//
// Params:
//   - toneA: First tone value (0-100).
//   - toneB: Second tone value (0-100).
//
// Returns float64 - Contrast ratio in range [1, 21]. Values outside [0,100] are
// clamped.
func RatioOfTones(toneA, toneB float64) float64 {
	toneA = num.Clamp(0, 100, toneA)
	toneB = num.Clamp(0, 100, toneB)
	return RatioOfYs(color.YFromLstar(toneA), color.YFromLstar(toneB))
}

// RatioOfYs calculates the contrast ratio between two relative luminance
// values.
//
// Params:
//   - y1: First relative luminance value.
//   - y2: Second relative luminance value.
//
// Returns float64 - Ratio computed as (lighter + 5) / (darker + 5).
func RatioOfYs(y1, y2 float64) float64 {
	lighter := max(y1, y2)
	darker := y2
	if lighter == y2 {
		darker = y1
	}
	return (lighter + 5.0) / (darker + 5.0)
}

// Lighter finds a lighter tone meeting the specified contrast ratio.
//
// Params:
//   - tone: Base tone value (0-100).
//   - ratio: Desired contrast ratio (1-21).
//
// Returns float64 - Lighter tone in [0,100]. Returns -1 if ratio cannot be met.
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

// Darker finds a darker tone meeting the specified contrast ratio.
//
// Params:
//   - tone: Base tone value (0-100).
//   - ratio: Desired contrast ratio (1-21).
//
// Returns float64 - Darker tone in [0,100]. Returns -1 if ratio cannot be met.
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

// LighterUnsafe finds a lighter tone attempting to meet contrast ratio.
//
// Unlike Lighter, always returns a valid tone even if ratio cannot be met.
// Params:
//   - tone: Base tone value (0-100).
//   - ratio: Desired contrast ratio (1-21).
//
// Returns float64 - Lighter tone in [0,100]. Returns 100 if ratio cannot be
// met.
func LighterUnsafe(tone, ratio float64) float64 {
	lighterSafe := Lighter(tone, ratio)
	if lighterSafe < 0 {
		return 100
	}
	return lighterSafe
}

// DarkerUnsafe finds a darker tone attempting to meet contrast ratio.
//
// Unlike Darker, always returns a valid tone even if ratio cannot be met.
// Params:
//   - tone: Base tone value (0-100).
//   - ratio: Desired contrast ratio (1-21).
//
// Returns float64 - Darker tone in [0,100]. Returns 0 if ratio cannot be met.
func DarkerUnsafe(tone, ratio float64) float64 {
	darkerSafe := Darker(tone, ratio)
	if darkerSafe < 0 {
		return 0
	}
	return darkerSafe
}
