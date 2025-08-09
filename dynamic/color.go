package dynamic

import (
	"math"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/contrast"
	"github.com/Nadim147c/material/palettes"
)

// Function type definitions
type (
	DynamicSchemeFn  func(s *DynamicScheme) any
	TonalPaletteFn   func(s *DynamicScheme) palettes.TonalPalette
	ToneFn           func(s *DynamicScheme) float64
	ChromaMultiplier func(s *DynamicScheme) float64
	DynamicColorFn   func(s *DynamicScheme) *DynamicColor
	ToneDeltaPairFn  func(s *DynamicScheme) *ToneDeltaPair
	ContrastCurveFn  func(s *DynamicScheme) *ContrastCurve
)

// DynamicColor represents a color in a dynamic color scheme
type DynamicColor struct {
	Name             string
	Palette          TonalPaletteFn
	Tone             ToneFn
	ChromaMultiplier ChromaMultiplier
	IsBackground     bool
	Background       DynamicColorFn
	SecondBackground DynamicColorFn
	ToneDeltaPair    ToneDeltaPairFn
	ContrastCurve    ContrastCurveFn
}

// ForegroundTone calculates a foreground tone that has sufficient contrast with a background tone
func ForegroundTone(bgTone, ratio float64) float64 {
	lighterTone := contrast.LighterUnsafe(bgTone, ratio)
	darkerTone := contrast.DarkerUnsafe(bgTone, ratio)
	lighterRatio := contrast.RatioOfTones(lighterTone, bgTone)
	darkerRatio := contrast.RatioOfTones(darkerTone, bgTone)
	preferLighter := TonePrefersLightForeground(bgTone)

	if preferLighter {
		negligibleDifference := (math.Abs(lighterRatio-darkerRatio) < 0.1 &&
			lighterRatio < ratio && darkerRatio < ratio)
		if lighterRatio >= ratio || lighterRatio >= darkerRatio || negligibleDifference {
			return lighterTone
		}
		return darkerTone
	} else {
		if darkerRatio >= ratio || darkerRatio >= lighterRatio {
			return darkerTone
		}
		return lighterTone
	}
}

func GetInitialToneFromBackground(background DynamicColorFn) ToneFn {
	if background == nil {
		return func(s *DynamicScheme) float64 {
			return 50
		}
	}
	return func(s *DynamicScheme) float64 { return background(s).GetTone(s) }
}

// EnableLightForeground adjusts a tone to enable light foreground if needed
func EnableLightForeground(tone float64) float64 {
	if TonePrefersLightForeground(tone) && !ToneAllowsLightForeground(tone) {
		return 49.0
	}
	return tone
}

// TonePrefersLightForeground determines if a tone prefers light foreground
func TonePrefersLightForeground(tone float64) bool {
	return math.Round(tone) < 60
}

// ToneAllowsLightForeground determines if a tone allows light foreground
func ToneAllowsLightForeground(tone float64) bool {
	return math.Round(tone) <= 49
}

// FromPalette creates a DynamicColor from a palette and tone function
func FromPalette(name string, palette TonalPaletteFn, tone ToneFn) *DynamicColor {
	return &DynamicColor{
		name, palette, tone, nil,
		false, // isBackground
		nil,   // background
		nil,   // secondBackground
		nil,   // contrastCurve
		nil,   // toneDeltaPair
	}
}

// GetArgb returns the ARGB value for the DynamicColor in the given scheme
func (dc *DynamicColor) GetArgb(scheme *DynamicScheme) color.ARGB {
	return dc.GetHct(scheme).ToARGB()
}

// GetHct returns the HCT color for the DynamicColor in the given scheme
func (dc *DynamicColor) GetHct(scheme *DynamicScheme) color.Hct {
	if scheme.Version == V2025 {
		return ColorCalculation2025.GetHct(scheme, dc)
	}
	return ColorCalculation2021.GetHct(scheme, dc)
}

func (dc *DynamicColor) GetTone(scheme *DynamicScheme) float64 {
	if scheme.Version == V2025 {
		return ColorCalculation2025.GetTone(scheme, dc)
	}
	return ColorCalculation2021.GetTone(scheme, dc)
}
