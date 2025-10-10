package dynamic

import (
	"math"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/contrast"
	"github.com/Nadim147c/material/palettes"
)

// Function type definitions
type (
	SchemeFunc        func(s *Scheme) any
	TonalPaletteFunc  func(s *Scheme) palettes.TonalPalette
	ToneFunc          func(s *Scheme) float64
	ChromaMultiplier  func(s *Scheme) float64
	ColorFunc         func(s *Scheme) *Color
	ToneDeltaPairFunc func(s *Scheme) *ToneDeltaPair
	ContrastCurveFunc func(s *Scheme) *ContrastCurve
)

// Color represents a color in a dynamic color scheme
type Color struct {
	Name             string
	Palette          TonalPaletteFunc
	Tone             ToneFunc
	ChromaMultiplier ChromaMultiplier
	IsBackground     bool
	Background       ColorFunc
	SecondBackground ColorFunc
	ToneDeltaPair    ToneDeltaPairFunc
	ContrastCurve    ContrastCurveFunc
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

func GetInitialToneFromBackground(background ColorFunc) ToneFunc {
	if background == nil {
		return func(s *Scheme) float64 {
			return 50
		}
	}
	return func(s *Scheme) float64 { return background(s).GetTone(s) }
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
func FromPalette(name string, palette TonalPaletteFunc, tone ToneFunc) *Color {
	return &Color{
		name, palette, tone, nil,
		false, // isBackground
		nil,   // background
		nil,   // secondBackground
		nil,   // contrastCurve
		nil,   // toneDeltaPair
	}
}

// GetArgb returns the ARGB value for the DynamicColor in the given scheme
func (dc *Color) GetArgb(scheme *Scheme) color.ARGB {
	return dc.GetHct(scheme).ToARGB()
}

// GetHct returns the HCT color for the DynamicColor in the given scheme
func (dc *Color) GetHct(scheme *Scheme) color.Hct {
	if scheme.Version == Version2025 {
		return ColorCalculation2025.GetHct(scheme, dc)
	}
	return ColorCalculation2021.GetHct(scheme, dc)
}

func (dc *Color) GetTone(scheme *Scheme) float64 {
	if scheme.Version == Version2025 {
		return ColorCalculation2025.GetTone(scheme, dc)
	}
	return ColorCalculation2021.GetTone(scheme, dc)
}
