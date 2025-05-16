package color

import "fmt"

// Hct represents a color in the HCT color space (Hue, Chroma, Tone).
// HCT provides a perceptually accurate color measurement system that can also
// accurately render what colors will appear as in different lighting
// environments.
type Hct struct {
	Hue    float64
	Chroma float64
	Tone   float64
}

// From creates an HCT color from the provided hue, chroma, and tone values.
//
// hue: 0 <= hue < 360; invalid values are corrected.
// chroma: 0 <= chroma < ?; Chroma has a different maximum for any given hue and
// tone.
// tone: 0 <= tone <= 100; invalid values are corrected.
func NewHct(hue, chroma, tone float64) *Hct {
	return solveToColor(hue, chroma, tone).ToHct()
}

// HctFromColor creates an HCT color from the provided ARGB integer.
func HctFromColor(argb Color) *Hct {
	cam := argb.ToCam16()
	return &Hct{cam.Hue, cam.Chroma, argb.LStar()}
}

// HctFromColor creates an HCT color from the provided ARGB integer.
func (h *Hct) ToCam16() *Cam16 {
	return Cam16FromJch(h.Tone, h.Chroma, h.Hue)
}

// ToInt returns the ARGB representation of this color.
func (h *Hct) ToColor() Color {
	return solveToColor(h.Hue, h.Chroma, h.Tone)
}

// String returns a string representation of the HCT color.
func (h *Hct) String() string {
	return fmt.Sprintf("HCT(%.4f, %.4f, %.4f) %s", h.Hue, h.Chroma, h.Tone, h.ToColor().AnsiBg("  "))
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
func (h *Hct) InViewingConditions(vc *Environmnet) *Hct {
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
