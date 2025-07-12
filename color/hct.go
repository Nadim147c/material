package color

import (
	"fmt"
	"math"
)

// Hct represents a color in the HCT (Hue, Chroma, Tone) color space.
//
// HCT provides a perceptually accurate color measurement system that can
// accurately render colors in different lighting environments.
type Hct struct {
	Hue    float64
	Chroma float64
	Tone   float64
}

// Ensure Color implements the color.Color interface
var _ digitalColor = (*Hct)(nil)

// NewHct creates an HCT color from hue, chroma, and tone values.
//
// Params:
//   - hue: Hue angle in degrees (0-360, invalid values corrected).
//   - chroma: Colorfulness (0-max varies by hue/tone).
//   - tone: Lightness (0-100, invalid values corrected).
//
// Returns Hct - The constructed HCT color.
func NewHct(hue, chroma, tone float64) Hct {
	return solveToARGB(hue, chroma, tone).ToHct()
}

// ToARGB converts HCT color to ARGB representation.
// Returns ARGB - 32-bit packed color value.
func (h Hct) ToARGB() ARGB {
	return solveToARGB(h.Hue, h.Chroma, h.Tone)
}

// ToXYZ converts HCT color to CIE XYZ color space.
// Returns XYZ - CIE XYZ color representation.
func (h Hct) ToXYZ() XYZ {
	return h.ToARGB().ToXYZ()
}

// ToLab converts HCT color to CIE L*a*b* color space.
// Returns Lab - CIE L*a*b* color representation.
func (h Hct) ToLab() Lab {
	return h.ToARGB().ToXYZ().ToLab()
}

// ToHct returns the receiver (implements digitalColor interface).
// Returns Hct - The color itself.
func (h Hct) ToHct() Hct {
	return h
}

// ToCam converts HCT color to CAM16 color appearance model.
// Returns *Cam16 - Pointer to CAM16 color representation.
func (h Hct) ToCam() *Cam16 {
	return h.ToARGB().ToCam()
}

// RGBA implements color.Color interface for HCT.
// Returns (r, g, b, a) - 16-bit per channel values (0-65535).
func (h Hct) RGBA() (uint32, uint32, uint32, uint32) {
	return solveToARGB(h.Hue, h.Chroma, h.Tone).RGBA()
}

// String returns a formatted string representation of HCT color.
// Returns string - Formatted as "HCT(H, C, T)" with ANSI background.
func (h Hct) String() string {
	return fmt.Sprintf("HCT(%.4f, %.4f, %.4f) %s", h.Hue, h.Chroma, h.Tone, h.ToARGB().AnsiBg("  "))
}

// Hash generates a uint64 hash value for the HCT color.
// Returns uint64 - Efficient hash value for color comparison.
func (h Hct) Hash() uint64 {
	// Convert each float to bits and extract portions for the hash
	hueBits := math.Float64bits(h.Hue)
	chromaBits := math.Float64bits(h.Chroma)
	toneBits := math.Float64bits(h.Tone)

	// Create hash using FNV-1a inspired approach, but with direct bit operations
	// for better performance, combining all three components
	hash := uint64(14695981039346656037) // FNV offset basis

	// Mix in the hue bits
	hash ^= (hueBits & 0xFFFFFFFF)
	hash *= 1099511628211 // FNV prime

	// Mix in the chroma bits
	hash ^= (chromaBits & 0xFFFFFFFF)
	hash *= 1099511628211

	// Mix in the tone bits
	hash ^= (toneBits & 0xFFFFFFFF)
	hash *= 1099511628211

	return hash
}

// IsBlue checks if the hue falls in the blue range (250-270°).
// Returns bool - True if hue is in blue range.
func (h Hct) IsBlue() bool {
	return h.Hue >= 250 && h.Hue < 270
}

// IsYellow checks if the hue falls in the yellow range (105-125°).
// Returns bool - True if hue is in yellow range.
func (h Hct) IsYellow() bool {
	return h.Hue >= 105 && h.Hue < 125
}

// IsCyan checks if the hue falls in the cyan range (170-207°).
// Returns bool - True if hue is in cyan range.
func (h Hct) IsCyan() bool {
	return h.Hue >= 170 && h.Hue < 207
}

// IsBlue checks if given hue falls in blue range (250-270°).
//
// Params:
//   - hue: Hue angle in degrees to check.
//
// Returns bool - True if hue is in blue range.
func IsBlue(hue float64) bool {
	return hue >= 250 && hue < 270
}

// IsYellow checks if given hue falls in yellow range (105-125°).
//
// Params:
//   - hue: Hue angle in degrees to check.
//
// Returns bool - True if hue is in yellow range.
func IsYellow(hue float64) bool {
	return hue >= 105 && hue < 125
}

// IsCyan checks if given hue falls in cyan range (170-207°).
//
// Params:
//   - hue: Hue angle in degrees to check.
//
// Returns bool - True if hue is in cyan range.
func IsCyan(hue float64) bool {
	return hue >= 170 && hue < 207
}

// InViewingConditions adjusts color appearance for different environments. Uses
// CAM16 color appearance model to account for viewing conditions.
//
// Params:
//   - env: *Environment containing viewing condition parameters.
//
// Returns Hct - Color adjusted for specified viewing conditions.
func (h Hct) InViewingConditions(env *Environmnet) Hct {
	cam := h.ToARGB().ToCam()
	viewedInEnv := cam.Viewed(env)
	newCam := viewedInEnv.ToCam()
	return newCam.ToHct()
}
