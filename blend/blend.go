package blend

import (
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/num"
)

// Harmonize adjusts hue of designColor to be closer to sourceColor's hue.
//
// Params:
//   - designColor: Color to adjust.
//   - sourceColor: Color to harmonize towards.
//
// Returns color.ARGB - Adjusted color with preserved chroma and tone. Hue is
// rotated towards sourceColor by up to 15 degrees.
func Harmonize(designColor color.ARGB, sourceColor color.ARGB) color.ARGB {
	fromHct := designColor.ToHct()
	toHct := sourceColor.ToHct()
	differenceDegrees := num.DifferenceDegrees(fromHct.Hue, toHct.Hue)
	rotationDegrees := min(differenceDegrees*0.5, 15.0)
	rotation := num.RotationDirection(fromHct.Hue, toHct.Hue)
	outputHue := num.NormalizeDegree(fromHct.Hue + rotationDegrees*rotation)
	return color.NewHct(outputHue, fromHct.Chroma, fromHct.Tone).ToARGB()
}

// HctHue blends hue of from towards to in HCT space while preserving chroma.
//
// Params:
//   - from: Starting color.
//   - to: Target color.
//   - amount: Blend ratio (0.0-1.0, 0.0=from, 1.0=to).
//
// Returns color.ARGB - Color with blended hue but original chroma and tone.
// Panics if amount is outside [0.0, 1.0].
func HctHue(from color.ARGB, to color.ARGB, amount float64) color.ARGB {
	ucs := Cam16Ucs(from, to, amount)
	ucsCam := ucs.ToHct()
	fromCam := from.ToHct()
	blended := color.NewHct(ucsCam.Hue, fromCam.Chroma, from.ToCam().J)
	return blended.ToARGB()
}

// Cam16Ucs blends colors in CAM16-UCS uniform color space.
//
// Params:
//   - from: Starting color in CAM16-UCS.
//   - to: Target color in CAM16-UCS.
//   - amount: Interpolation factor (0.0-1.0).
//
// Returns color.ARGB - Fully blended color with interpolated attributes.
// Blends all color attributes (hue, chroma, tone) simultaneously.
func Cam16Ucs(from color.ARGB, to color.ARGB, amount float64) color.ARGB {
	fromCam := from.ToCam()
	toCam := to.ToCam()
	fromJ := fromCam.Jstar
	fromA := fromCam.Astar
	fromB := fromCam.Bstar
	toJ := toCam.Jstar
	toA := toCam.Astar
	toB := toCam.Bstar

	jstar := fromJ + (toJ-fromJ)*amount
	astar := fromA + (toA-fromA)*amount
	bstar := fromB + (toB-fromB)*amount
	return color.Cam16FromUcs(jstar, astar, bstar).ToARGB()
}
