package blend

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/num"
)

// Harmonize adjusts the hue of designColor to be closer to the hue of
// sourceColor. It rotates the hue of designColor towards sourceColor by up to
// 15 degrees. The chroma and tone of designColor are preserved.
func Harmonize(designColor color.ARGB, sourceColor color.ARGB) color.ARGB {
	fromHct := designColor.ToHct()
	toHct := sourceColor.ToHct()
	differenceDegrees := num.DifferenceDegrees(fromHct.Hue, toHct.Hue)
	rotationDegrees := min(differenceDegrees*0.5, 15.0)
	rotation := num.RotationDirection(fromHct.Hue, toHct.Hue)
	outputHue := num.NormalizeDegree(fromHct.Hue + rotationDegrees*rotation)
	return color.NewHct(outputHue, fromHct.Chroma, fromHct.Tone).ToARGB()
}

// HctHue blends the hue of from towards the hue of to in HCT color space. The
// chroma and tone of from are preserved. amount must be between 0.0 and 1.0.
func HctHue(from color.ARGB, to color.ARGB, amount float64) color.ARGB {
	ucs := Cam16Ucs(from, to, amount)
	ucsCam := ucs.ToHct()
	fromCam := from.ToHct()
	blended := color.NewHct(ucsCam.Hue, fromCam.Chroma, from.ToCam().J)
	return blended.ToARGB()
}

// Cam16Ucs blends from towards to in CAM16-UCS color space. Hue, chroma, and
// tone will all change. amount must be between 0.0 and 1.0.
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
