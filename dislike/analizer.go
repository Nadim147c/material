package dislike

import "github.com/Nadim147c/material/color"

// IsDisliked determines if a color is considered visually unpleasant. Disliked
// colors are dark yellow-greens with sufficient saturation.
//
// Params:
//   - hct: Color in HCT format to evaluate.
//
// Returns bool - true if color meets disliked criteria, false otherwise.
func IsDisliked(hct color.Hct) bool {
	huePasses := int(hct.Hue+0.5) >= 90 && int(hct.Hue+0.5) <= 111
	chromaPasses := int(hct.Chroma+0.5) > 16
	tonePasses := int(hct.Tone+0.5) < 65
	return huePasses && chromaPasses && tonePasses
}

// FixIfDisliked lightens a color if it is considered disliked. If the color is
// identified as disliked, this function creates a new color with the same hue
// and chroma but with a tone of 70.0, making it more visually pleasant.
//
// Parameters:
//   - hct: A color in HCT format to be evaluated and potentially modified
//
// Returns a new color if the original is disliked, or the original color if it
// is acceptable
func FixIfDisliked(hct color.Hct) color.Hct {
	if IsDisliked(hct) {
		return color.NewHct(hct.Hue, hct.Chroma, 70.0)
	}
	return hct
}
