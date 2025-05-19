package dynamic

import (
	"math"
	"strings"

	"github.com/Nadim147c/goyou/color"
	"github.com/Nadim147c/goyou/contrast"
	"github.com/Nadim147c/goyou/num"
)

type ColorCalculationDelegate interface {
	GetHct(scheme DynamicScheme, color DynamicColor) color.Hct
	GetTone(scheme DynamicScheme, color DynamicColor) float64
}

type ColorCalculationDelegateImpl2021 struct{}

var ColorCalculation2025 = ColorCalculationDelegateImpl2025{}

type ColorCalculationDelegateImpl2025 struct{}

var ColorCalculation2021 = ColorCalculationDelegateImpl2021{}

func (d ColorCalculationDelegateImpl2021) GetHct(scheme DynamicScheme, dc DynamicColor) color.Hct {
	tone := dc.GetTone(scheme)
	palette := dc.Palette(scheme)
	return palette.GetHct(tone)
}

func (d ColorCalculationDelegateImpl2021) GetTone(scheme DynamicScheme, dc DynamicColor) float64 {
	decreasingContrast := scheme.ContrastLevel < 0

	if dc.ToneDeltaPair != nil && dc.ToneDeltaPair(scheme) != nil {
		toneDeltaPair := dc.ToneDeltaPair(scheme)
		roleA := toneDeltaPair.RoleA
		roleB := toneDeltaPair.RoleB
		delta := toneDeltaPair.Delta
		polarity := toneDeltaPair.Polarity
		stayTogether := toneDeltaPair.StayTogether

		aIsNearer := polarity == ToneNearer ||
			(polarity == ToneLighter && !scheme.IsDark) ||
			(polarity == ToneDarker && scheme.IsDark)
		nearer := roleA
		farther := roleB
		if !aIsNearer {
			nearer = roleB
			farther = roleA
		}
		amNearer := dc.Name == nearer.Name
		expansionDir := -1.0
		if scheme.IsDark {
			expansionDir = 1.0
		}
		nTone := nearer.Tone(scheme)
		fTone := farther.Tone(scheme)

		if dc.Background != nil && nearer.ContrastCurve != nil && farther.ContrastCurve != nil {
			bg := dc.Background(scheme)
			nCurve := nearer.ContrastCurve(scheme)
			fCurve := farther.ContrastCurve(scheme)
			if bg != nil && nCurve != nil && fCurve != nil {
				bgTone := bg.GetTone(scheme)
				nContrast := nCurve.Get(scheme.ContrastLevel)
				fContrast := fCurve.Get(scheme.ContrastLevel)

				if contrast.RatioOfTones(bgTone, nTone) < nContrast {
					nTone = ForegroundTone(bgTone, nContrast)
				}
				if contrast.RatioOfTones(bgTone, fTone) < fContrast {
					fTone = ForegroundTone(bgTone, fContrast)
				}
				if decreasingContrast {
					nTone = ForegroundTone(bgTone, nContrast)
					fTone = ForegroundTone(bgTone, fContrast)
				}
			}
		}

		if (fTone-nTone)*expansionDir < delta {
			fTone = num.Clamp(0, 100, nTone+delta*expansionDir)
			if (fTone-nTone)*expansionDir < delta {
				nTone = num.Clamp(0, 100, fTone-delta*expansionDir)
			}
		}

		if nTone >= 50 && nTone < 60 {
			if expansionDir > 0 {
				nTone = 60
				fTone = math.Max(fTone, nTone+delta*expansionDir)
			} else {
				nTone = 49
				fTone = math.Min(fTone, nTone+delta*expansionDir)
			}
		} else if fTone >= 50 && fTone < 60 {
			if stayTogether {
				if expansionDir > 0 {
					nTone = 60
					fTone = math.Max(fTone, nTone+delta*expansionDir)
				} else {
					nTone = 49
					fTone = math.Min(fTone, nTone+delta*expansionDir)
				}
			} else {
				if expansionDir > 0 {
					fTone = 60
				} else {
					fTone = 49
				}
			}
		}

		if amNearer {
			return nTone
		}
		return fTone
	}

	answer := dc.Tone(scheme)
	bg := dc.Background
	contrastCurve := dc.ContrastCurve

	if bg == nil || bg(scheme) == nil || contrastCurve == nil || contrastCurve(scheme) == nil {
		return answer
	}

	bgTone := bg(scheme).GetTone(scheme)
	desiredRatio := contrastCurve(scheme).Get(scheme.ContrastLevel)

	if contrast.RatioOfTones(bgTone, answer) < desiredRatio || decreasingContrast {
		answer = ForegroundTone(bgTone, desiredRatio)
	}

	if dc.IsBackground && answer >= 50 && answer < 60 {
		if contrast.RatioOfTones(49, bgTone) >= desiredRatio {
			answer = 49
		} else {
			answer = 60
		}
	}

	if dc.SecondBackground == nil || dc.SecondBackground(scheme) == nil {
		return answer
	}

	bg1 := dc.Background(scheme)
	bg2 := dc.SecondBackground(scheme)
	bgTone1 := bg1.GetTone(scheme)
	bgTone2 := bg2.GetTone(scheme)
	upper := math.Max(bgTone1, bgTone2)
	lower := math.Min(bgTone1, bgTone2)

	if contrast.RatioOfTones(upper, answer) >= desiredRatio &&
		contrast.RatioOfTones(lower, answer) >= desiredRatio {
		return answer
	}

	lightOption := contrast.Lighter(upper, desiredRatio)
	darkOption := contrast.Darker(lower, desiredRatio)

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

func (d ColorCalculationDelegateImpl2025) GetHct(scheme DynamicScheme, dc DynamicColor) color.Hct {
	palette := dc.Palette(scheme)
	tone := dc.GetTone(scheme)
	chromaMultiplier := 1.0
	if dc.ChromaMultiplier != nil {
		chromaMultiplier = dc.ChromaMultiplier(scheme)
	}
	return color.NewHct(palette.Hue, palette.Chroma*chromaMultiplier, tone)
}

func (d ColorCalculationDelegateImpl2025) GetTone(scheme DynamicScheme, dc DynamicColor) float64 {
	toneDeltaPair := dc.ToneDeltaPair(scheme)
	if toneDeltaPair != nil {
		roleA := toneDeltaPair.RoleA
		roleB := toneDeltaPair.RoleB
		polarity := toneDeltaPair.Polarity
		constraint := toneDeltaPair.Constraint
		delta := toneDeltaPair.Delta

		absoluteDelta := delta
		if polarity == ToneDarker || (polarity == ToneRelativeLighter && scheme.IsDark) ||
			(polarity == ToneRelativeDarker && !scheme.IsDark) {
			absoluteDelta = -delta
		}

		amRoleA := dc.Name == roleA.Name
		selfRole := roleA
		refRole := roleB
		if !amRoleA {
			selfRole = roleB
			refRole = roleA
		}

		selfTone := selfRole.Tone(scheme)
		refTone := refRole.GetTone(scheme)
		relativeDelta := absoluteDelta
		if !amRoleA {
			relativeDelta *= -1
		}

		switch constraint {
		case "exact":
			selfTone = num.Clamp(0, 100, refTone+relativeDelta)
		case "nearer":
			if relativeDelta > 0 {
				selfTone = num.Clamp(0, 100, num.Clamp(refTone, refTone+relativeDelta, selfTone))
			} else {
				selfTone = num.Clamp(0, 100, num.Clamp(refTone+relativeDelta, refTone, selfTone))
			}
		case "farther":
			if relativeDelta > 0 {
				selfTone = num.Clamp(refTone+relativeDelta, 100, selfTone)
			} else {
				selfTone = num.Clamp(0, refTone+relativeDelta, selfTone)
			}
		}

		if dc.Background != nil && dc.ContrastCurve != nil {
			bg := dc.Background(scheme)
			cc := dc.ContrastCurve(scheme)
			if bg != nil && cc != nil {
				bgTone := bg.GetTone(scheme)
				desiredContrast := cc.Get(scheme.ContrastLevel)
				if contrast.RatioOfTones(bgTone, selfTone) < desiredContrast || scheme.ContrastLevel < 0 {
					selfTone = ForegroundTone(bgTone, desiredContrast)
				}
			}
		}

		if dc.IsBackground && !strings.HasSuffix(dc.Name, "_fixed_dim") {
			if selfTone >= 57 {
				selfTone = num.Clamp(65, 100, selfTone)
			} else {
				selfTone = num.Clamp(0, 49, selfTone)
			}
		}
		return selfTone
	}

	answer := dc.Tone(scheme)
	bg := dc.Background
	cc := dc.ContrastCurve

	if bg == nil || bg(scheme) == nil || cc == nil || cc(scheme) == nil {
		return answer
	}

	bgTone := bg(scheme).GetTone(scheme)
	desiredRatio := cc(scheme).Get(scheme.ContrastLevel)
	if contrast.RatioOfTones(bgTone, answer) < desiredRatio || scheme.ContrastLevel < 0 {
		answer = ForegroundTone(bgTone, desiredRatio)
	}

	if dc.IsBackground && !strings.HasSuffix(dc.Name, "_fixed_dim") {
		if answer >= 57 {
			answer = num.Clamp(65, 100, answer)
		} else {
			answer = num.Clamp(0, 49, answer)
		}
	}

	if dc.SecondBackground == nil || dc.SecondBackground(scheme) == nil {
		return answer
	}

	bg1 := dc.Background(scheme)
	bg2 := dc.SecondBackground(scheme)
	bgTone1 := bg1.GetTone(scheme)
	bgTone2 := bg2.GetTone(scheme)
	upper := math.Max(bgTone1, bgTone2)
	lower := math.Min(bgTone1, bgTone2)

	if contrast.RatioOfTones(upper, answer) >= desiredRatio &&
		contrast.RatioOfTones(lower, answer) >= desiredRatio {
		return answer
	}

	lightOption := contrast.Lighter(upper, desiredRatio)
	darkOption := contrast.Darker(lower, desiredRatio)

	availables := []float64{}
	if lightOption != -1 {
		availables = append(availables, lightOption)
	}
	if darkOption != -1 {
		availables = append(availables, darkOption)
	}

	prefersLight := TonePrefersLightForeground(bgTone1) || TonePrefersLightForeground(bgTone2)
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
