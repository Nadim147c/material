package color

import (
	"fmt"
)

// Hct represents a color in the HCT color space (Hue, Chroma, Tone).
// HCT provides a perceptually accurate color measurement system that can also
// accurately render what colors will appear as in different lighting
// environments.
type Hct struct {
	Hue    float64
	Chroma float64
	Tone   float64
	argb   Color
}

// From creates an HCT color from the provided hue, chroma, and tone values.
//
// hue: 0 <= hue < 360; invalid values are corrected.
// chroma: 0 <= chroma < ?; Chroma has a different maximum for any given hue and
// tone.
// tone: 0 <= tone <= 100; invalid values are corrected.
func NewHct(hue, chroma, tone float64) *Hct {
	argb := SolveToColor(hue, chroma, tone)
	h := &Hct{}
	h.setInternalState(argb)
	return h
}

// HctFromColor creates an HCT color from the provided ARGB integer.
func HctFromColor(argb Color) *Hct {
	hct := &Hct{argb: argb}
	cam := argb.ToCam16()
	hct.Hue = cam.Hue
	hct.Chroma = cam.Chroma
	hct.Tone = argb.ToXYZ().LStar()
	return hct
}

// ToInt returns the ARGB representation of this color.
func (h *Hct) ToColor() Color {
	return h.argb
}

// SetHue sets the hue value of this color.
// newHue: 0 <= newHue < 360; invalid values are corrected.
// Chroma may decrease because chroma has a different maximum for any given
// hue and tone.
func (h *Hct) SetHue(newHue float64) {
	h.setInternalState(SolveToColor(newHue, h.Chroma, h.Tone))
}

// SetChroma sets the chroma value of this color.
// newChroma: 0 <= newChroma < ?
// Chroma may decrease because chroma has a different maximum for any given
// hue and tone.
func (h *Hct) SetChroma(newChroma float64) {
	h.setInternalState(SolveToColor(h.Hue, newChroma, h.Tone))
}

// SetTone sets the tone value of this color.
// newTone: 0 <= newTone <= 100; invalid values are corrected.
// Chroma may decrease because chroma has a different maximum for any given
// hue and tone.
func (h *Hct) SetTone(newTone float64) {
	h.setInternalState(SolveToColor(h.Hue, h.Chroma, newTone))
}

// SetValue sets a property of the Hct object by name.
func (h *Hct) SetValue(propertyName string, value float64) {
	switch propertyName {
	case "hue":
		h.SetHue(value)
	case "chroma":
		h.SetChroma(value)
	case "tone":
		h.SetTone(value)
	}
}

// String returns a string representation of the HCT color.
func (h *Hct) String() string {
	return fmt.Sprintf("HCT(%.0f, %.0f, %.0f)", h.Hue, h.Chroma, h.Tone)
}

// IsBlue determines if a hue is in the blue range.
func IsBlue(hue float64) bool {
	return hue >= 250 && hue < 270
}

// IsYellow determines if a hue is in the yellow range.
func IsYellow(hue float64) bool {
	return hue >= 105 && hue < 125
}

// IsCyan determines if a hue is in the cyan range.
func IsCyan(hue float64) bool {
	return hue >= 170 && hue < 207
}

// setInternalState updates the  state of the Hct object.
func (h *Hct) setInternalState(argb Color) {
	cam := argb.ToCam16()
	h.Hue = cam.Hue
	h.Chroma = cam.Chroma
	h.Tone = cam.J
	h.argb = argb
}

// InViewingConditions translates a color into different ViewingConditions.
//
// Colors change appearance. They look different with lights on versus off,
// the same color, as in hex code, on white looks different when on black.
// This is called color relativity, most famously explicated by Josef Albers
// in Interaction of Color.
//
// In color science, color appearance models can account for this and
// calculate the appearance of a color in different settings. HCT is based on
// CAM16, a color appearance model, and uses it to make these calculations.
//
// See ViewingConditions.Make for parameters affecting color appearance.
// TODO: complete this function
func (h *Hct) InViewingConditions(vc *ViewingConditions) *Hct {
	// 1. Use CAM16 to find XYZ coordinates of color in specified VC.
	// cam := h.argb.ToCam16()
	return nil
	// viewedInVc := cam.XyzInViewingConditions(vc)
	//
	// // 2. Create CAM16 of those XYZ coordinates in default VC.
	// recastInVc := Cam16_FromXyzInViewingConditions(
	// 	viewedInVc[0],
	// 	viewedInVc[1],
	// 	viewedInVc[2],
	// 	ViewingConditions_Make(),
	// )
	//
	// // 3. Create HCT from:
	// // - CAM16 using default VC with XYZ coordinates in specified VC.
	// // - L* converted from Y in XYZ coordinates in specified VC.
	// recastHct := From(
	// 	recastInVc.Hue(),
	// 	recastInVc.Chroma(),
	// 	LstarFromY(viewedInVc[1]),
	// )
	// return recastHct
}
