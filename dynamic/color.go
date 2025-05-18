package dynamic

import (
	"math"

	"github.com/Nadim147c/goyou/color"
	"github.com/Nadim147c/goyou/contrast"
	"github.com/Nadim147c/goyou/palettes"
)

// Function type definitions
type (
	DynamicSchemeFn func(s DynamicScheme) any
	TonalPaletteFn  func(s DynamicScheme) palettes.TonalPalette
	ToneFn          func(s DynamicScheme) float64
	DynamicColorFn  func(s DynamicScheme) DynamicColor
	ToneDeltaPairFn func(s DynamicScheme) ToneDeltaPair
)

// DynamicColor represents a color in a dynamic color scheme
type DynamicColor struct {
	Name             string
	Palette          TonalPaletteFn
	Tone             ToneFn
	IsBackground     bool
	Background       DynamicColorFn
	SecondBackground DynamicColorFn
	ToneDeltaPair    ToneDeltaPairFn
	ContrastCurve    *contrastCurve
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
func FromPalette(name string, palette TonalPaletteFn, tone ToneFn) DynamicColor {
	return DynamicColor{
		name,
		palette,
		tone,
		false, // isBackground
		nil,   // background
		nil,   // secondBackground
		nil,   // contrastCurve
		nil,   // toneDeltaPair
	}
}

// GetArgb returns the ARGB value for the DynamicColor in the given scheme
func (dc DynamicColor) GetArgb(scheme DynamicScheme) color.ARGB {
	p := dc.Palette(scheme)
	color := p.Get(dc.GetTone(scheme))
	return color
}

// GetHct returns the HCT color for the DynamicColor in the given scheme
func (dc DynamicColor) GetHct(scheme DynamicScheme) color.Hct {
	return dc.GetArgb(scheme).ToHct()
}

// GetTone returns the tone for the DynamicColor in the given scheme
func (dc DynamicColor) GetTone(scheme DynamicScheme) float64 {
	decreasingContrast := scheme.ContrastLevel < 0

	// Case 1: dual foreground, pair of colors with delta constraint.
	if dc.ToneDeltaPair != nil {
		toneDeltaPair := dc.ToneDeltaPair(scheme)
		roleA := toneDeltaPair.RoleA
		roleB := toneDeltaPair.RoleB
		delta := toneDeltaPair.Delta
		polarity := toneDeltaPair.Polarity
		stayTogether := toneDeltaPair.StayTogether

		bg := dc.Background(scheme)
		bgTone := bg.GetTone(scheme)

		aIsNearer := (polarity == Nearer ||
			(polarity == Lighter && !scheme.IsDark) ||
			(polarity == Darker && scheme.IsDark))
		var nearer, farther DynamicColor
		if aIsNearer {
			nearer = roleA
			farther = roleB
		} else {
			nearer = roleB
			farther = roleA
		}
		amNearer := dc.Name == nearer.Name
		expansionDir := 1.0
		if !scheme.IsDark {
			expansionDir = -1.0
		}

		// 1st round: solve to min, each
		nContrast := nearer.ContrastCurve.Get(scheme.ContrastLevel)
		fContrast := farther.ContrastCurve.Get(scheme.ContrastLevel)

		// If a color is good enough, it is not adjusted.
		// Initial and adjusted tones for `nearer`
		nInitialTone := nearer.Tone(scheme)
		nTone := nInitialTone
		if contrast.RatioOfTones(bgTone, nInitialTone) < nContrast {
			nTone = ForegroundTone(bgTone, nContrast)
		}

		// Initial and adjusted tones for `farther`
		fInitialTone := farther.Tone(scheme)
		fTone := fInitialTone
		if contrast.RatioOfTones(bgTone, fInitialTone) < fContrast {
			fTone = ForegroundTone(bgTone, fContrast)
		}

		if decreasingContrast {
			// If decreasing contrast, adjust color to the "bare minimum"
			// that satisfies contrast.
			nTone = ForegroundTone(bgTone, nContrast)
			fTone = ForegroundTone(bgTone, fContrast)
		}

		if (fTone-nTone)*expansionDir < delta {
			// 2nd round: expand farther to match delta.
			fTone = math.Min(math.Max(nTone+delta*expansionDir, 0.0), 100.0)
			if (fTone-nTone)*expansionDir < delta {
				// 3rd round: contract nearer to match delta.
				nTone = math.Min(math.Max(fTone-delta*expansionDir, 0.0), 100.0)
			}
		}

		// Avoids the 50-59 awkward zone.
		if 50 <= nTone && nTone < 60 {
			// If `nearer` is in the awkward zone, move it away, together with `farther`.
			if expansionDir > 0 {
				nTone = 60
				fTone = math.Max(fTone, nTone+delta*expansionDir)
			} else {
				nTone = 49
				fTone = math.Min(fTone, nTone+delta*expansionDir)
			}
		} else if 50 <= fTone && fTone < 60 {
			if stayTogether {
				// Fixes both, to avoid two colors on opposite sides of the "awkward zone".
				if expansionDir > 0 {
					nTone = 60
					fTone = math.Max(fTone, nTone+delta*expansionDir)
				} else {
					nTone = 49
					fTone = math.Min(fTone, nTone+delta*expansionDir)
				}
			} else {
				// Not required to stay together; fixes just one.
				if expansionDir > 0 {
					fTone = 60
				} else {
					fTone = 49
				}
			}
		}

		// Returns `nTone` if this color is `nearer`, otherwise `fTone`.
		if amNearer {
			return nTone
		}
		return fTone
	} else {
		// Case 2: No contrast pair; just solve for itself.
		answer := dc.Tone(scheme)

		if dc.Background == nil {
			return answer // No adjustment for colors with no background.
		}

		bgTone := dc.Background(scheme).GetTone(scheme)
		desiredRatio := dc.ContrastCurve.Get(scheme.ContrastLevel)

		if contrast.RatioOfTones(bgTone, answer) < desiredRatio {
			// Rough improvement.
			answer = ForegroundTone(bgTone, desiredRatio)
		}

		if decreasingContrast {
			answer = ForegroundTone(bgTone, desiredRatio)
		}

		if dc.IsBackground && 50 <= answer && answer < 60 {
			// Must adjust
			if contrast.RatioOfTones(49, bgTone) >= desiredRatio {
				answer = 49
			} else {
				answer = 60
			}
		}

		if dc.SecondBackground != nil {
			// Case 3: Adjust for dual backgrounds.
			bgTone1 := dc.Background(scheme).GetTone(scheme)
			bgTone2 := dc.SecondBackground(scheme).GetTone(scheme)

			upper := math.Max(bgTone1, bgTone2)
			lower := math.Min(bgTone1, bgTone2)

			if contrast.RatioOfTones(upper, answer) >= desiredRatio &&
				contrast.RatioOfTones(lower, answer) >= desiredRatio {
				return answer
			}

			// The darkest light tone that satisfies the desired ratio,
			// or -1 if such ratio cannot be reached.
			lightOption := contrast.Lighter(upper, desiredRatio)

			// The lightest dark tone that satisfies the desired ratio,
			// or -1 if such ratio cannot be reached.
			darkOption := contrast.Darker(lower, desiredRatio)

			// Tones suitable for the foreground.
			availables := []float64{}
			if lightOption != -1 {
				availables = append(availables, lightOption)
			}
			if darkOption != -1 {
				availables = append(availables, darkOption)
			}

			prefersLight := TonePrefersLightForeground(bgTone1) ||
				TonePrefersLightForeground(bgTone2)
			if prefersLight {
				if lightOption < 0 {
					return 100
				}
				return lightOption
			}
			if len(availables) == 1 {
				return availables[0]
			}
			if darkOption < 0 {
				return 0
			}
			return darkOption
		}

		return answer
	}
}
