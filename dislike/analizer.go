package dislike

import "github.com/Nadim147c/goyou/color"

// IsDisliked determines if a color is considered disliked.
//
// Disliked is defined as a dark yellow-green that is not neutral.
// Specifically, a color with:
//   - Hue between 90 and 111 degrees
//   - Chroma greater than 16
//   - Tone less than 65
//
// Parameters:
//   - hct: A color in HCT format to be evaluated
//
// Returns:
//   - true if the color is considered disliked, false otherwise
func IsDisliked(hct *color.Hct) bool {
	huePasses := int(hct.Hue+0.5) >= 90 && int(hct.Hue+0.5) <= 111
	chromaPasses := int(hct.Chroma+0.5) > 16
	tonePasses := int(hct.Tone+0.5) < 65
	return huePasses && chromaPasses && tonePasses
}

// FixIfDisliked lightens a color if it is considered disliked.
//
// If the color is identified as disliked, this function creates a new color
// with the same hue and chroma but with a tone of 70.0, making it more visually pleasant.
//
// Parameters:
//   - hct: A color in HCT format to be evaluated and potentially modified
//
// Returns:
//   - A new color if the original is disliked, or the original color if it is acceptable
func FixIfDisliked(hct *color.Hct) *color.Hct {
	if IsDisliked(hct) {
		return color.NewHct(hct.Hue, hct.Chroma, 70.0)
	}
	return hct
}
