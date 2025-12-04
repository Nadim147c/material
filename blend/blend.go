package blend

import (
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/num"
)

// Harmonize returns a color whose hue is shifted toward the hue of sourceColor.
// The hue of designColor is rotated by up to 15 degrees toward sourceColor,
// while preserving its original chroma and tone.
func Harmonize(designColor color.ARGB, sourceColor color.ARGB) color.ARGB {
	fromHct := designColor.ToHct()
	toHct := sourceColor.ToHct()

	differenceDegrees := num.DifferenceDegrees(fromHct.Hue, toHct.Hue)
	rotationDegrees := min(differenceDegrees*0.5, 15.0)
	rotation := num.RotationDirection(fromHct.Hue, toHct.Hue)
	outputHue := num.NormalizeDegree(fromHct.Hue + rotationDegrees*rotation)

	return color.NewHct(outputHue, fromHct.Chroma, fromHct.Tone).ToARGB()
}

// HctHueDirect returns a color with its hue blended toward another color in HCT
// space. The chroma and tone of from are preserved. The amount must be in [0.0,
// 1.0], where 0.0 yields from and 1.0 yields to. Returns Hct representation of
// new color.
func HctHueDirect(from color.ARGB, to color.ARGB, amount float64) color.Hct {
	ucs := Cam16Ucs(from, to, amount)

	ucsCam := ucs.ToHct()
	fromCam := from.ToHct()

	return color.NewHct(ucsCam.Hue, fromCam.Chroma, from.ToCam16().J)
}

// HctHue returns a color with its hue blended toward another color in HCT
// space. The chroma and tone of from are preserved. The amount must be in [0.0,
// 1.0], where 0.0 yields from and 1.0 yields to.
func HctHue(from color.ARGB, to color.ARGB, amount float64) color.ARGB {
	return HctHueDirect(from, to, amount).ToARGB()
}

// Cam16Ucs returns a color interpolated between two colors in CAM16-UCS space.
// All perceptual attributes (hue, chroma, tone) are blended linearly based on
// amount, where 0.0 yields from and 1.0 yields to.
func Cam16Ucs(from color.ARGB, to color.ARGB, amount float64) color.ARGB {
	fromCam := from.ToCam16()
	toCam := to.ToCam16()
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
