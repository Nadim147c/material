package dynamic

import (
	"fmt"

	"github.com/Nadim147c/material/num"
	"github.com/Nadim147c/material/palettes"
)

// tMaxC
// Paramters:
//   - lowerBound: 0
//   - upperBound: 100
//   - chromaMultiplier: 1
func tMaxC(palette palettes.TonalPalette, params ...float64) float64 {
	lowerBound := 0.0
	upperBound := 100.0
	chromaMultiplier := 1.0
	if len(params) > 2 {
		lowerBound = params[0]
		upperBound = params[1]
		chromaMultiplier = params[2]
	} else if len(params) > 1 {
		lowerBound = params[0]
		upperBound = params[1]
	} else if len(params) > 0 {
		lowerBound = params[0]
	}
	answer := FindDesiredChromaByTone(palette.Hue, palette.Chroma*chromaMultiplier, 100, true)
	return num.Clamp(lowerBound, upperBound, answer)
}

// tMaxC
// Paramters:
//   - lowerBound: 0
//   - upperBound: 100
func tMinC(palette palettes.TonalPalette, bounds ...float64) float64 {
	lowerBound := 0.0
	upperBound := 100.0
	if len(bounds) > 1 {
		lowerBound = bounds[0]
		upperBound = bounds[1]
	} else if len(bounds) > 0 {
		lowerBound = bounds[0]
	}

	answer := FindDesiredChromaByTone(palette.Hue, palette.Chroma, 0, false)

	return num.Clamp(lowerBound, upperBound, answer)
}

// GetCurve returns the contrast curve for a given default contrast.
func GetCurve(defaultContrast float64) *ContrastCurve {
	switch defaultContrast {
	case 1.5:
		return NewContrastCurve(1.5, 1.5, 3, 4.5)
	case 3:
		return NewContrastCurve(3, 3, 4.5, 7)
	case 4.5:
		return NewContrastCurve(4.5, 4.5, 7, 11)
	case 6:
		return NewContrastCurve(6, 6, 7, 11)
	case 7:
		return NewContrastCurve(7, 7, 11, 21)
	case 9:
		return NewContrastCurve(9, 9, 11, 21)
	case 11:
		return NewContrastCurve(11, 11, 21, 21)
	case 21:
		return NewContrastCurve(21, 21, 21, 21)
	default:
		return NewContrastCurve(defaultContrast, defaultContrast, 7, 21)
	}
}

func validateExtendedColor(
	originalColor *DynamicColor,
	specVersion Version,
	extendedColor *DynamicColor,
) {
	if originalColor.Name != extendedColor.Name {
		panic(fmt.Sprintf(
			"Attempting to extend color %s with color %s of different name for spec version %v.",
			originalColor.Name, extendedColor.Name, specVersion,
		))
	}
	if originalColor.IsBackground != extendedColor.IsBackground {
		panic(fmt.Sprintf(
			"Attempting to extend color %s as a %s with color %s as a %s for spec version %v.",
			originalColor.Name,
			boolToBackgroundForeground(originalColor.IsBackground),
			extendedColor.Name,
			boolToBackgroundForeground(extendedColor.IsBackground),
			specVersion,
		))
	}
}

func boolToBackgroundForeground(isBackground bool) string {
	if isBackground {
		return "background"
	}
	return "foreground"
}

func ExtendSpecVersion(
	originalColor *DynamicColor,
	specVersion Version,
	extendedColor *DynamicColor,
) *DynamicColor {
	validateExtendedColor(originalColor, specVersion, extendedColor)

	return &DynamicColor{
		Name: originalColor.Name,
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			if s.Version == specVersion {
				return extendedColor.Palette(s)
			}
			return originalColor.Palette(s)
		},
		Tone: func(s *DynamicScheme) float64 {
			if s.Version == specVersion {
				return extendedColor.Tone(s)
			}
			return originalColor.Tone(s)
		},
		IsBackground: originalColor.IsBackground,
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			var chromaMultiplier ChromaMultiplier
			if s.Version == specVersion {
				chromaMultiplier = extendedColor.ChromaMultiplier
			} else {
				chromaMultiplier = originalColor.ChromaMultiplier
			}
			if chromaMultiplier != nil {
				return chromaMultiplier(s)
			}
			return 1
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			var background DynamicColorFn
			if s.Version == specVersion {
				background = extendedColor.Background
			} else {
				background = originalColor.Background
			}
			if background != nil {
				return background(s)
			}
			return nil
		},
		SecondBackground: func(s *DynamicScheme) *DynamicColor {
			var secondBackground DynamicColorFn
			if s.Version == specVersion {
				secondBackground = extendedColor.SecondBackground
			} else {
				secondBackground = originalColor.SecondBackground
			}
			if secondBackground != nil {
				return secondBackground(s)
			}
			return nil
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			var contrastCurve ContrastCurveFn
			if s.Version == specVersion {
				contrastCurve = extendedColor.ContrastCurve
			} else {
				contrastCurve = originalColor.ContrastCurve
			}
			if contrastCurve != nil {
				return contrastCurve(s)
			}
			return nil
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			var toneDeltaPair ToneDeltaPairFn
			if s.Version == specVersion {
				toneDeltaPair = extendedColor.ToneDeltaPair
			} else {
				toneDeltaPair = originalColor.ToneDeltaPair
			}
			if toneDeltaPair != nil {
				return toneDeltaPair(s)
			}
			return nil
		},
	}
}

type MaterialSpec2025 struct {
	MaterialSpec2021
}

var _ MaterialColorSpec = (*MaterialSpec2025)(nil)

func (m MaterialSpec2025) Surface() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				if s.IsDark {
					return 4
				}
				if s.NeutralPalette.IsBlue() {
					return 99
				} else if s.Variant == Vibrant {
					return 97
				}
				return 98
			}
			return 0
		},
		IsBackground: true,
	})
	return ExtendSpecVersion(m.MaterialSpec2021.Surface(), V2025, color)
}

func (m MaterialSpec2025) SurfaceDim() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface_dim",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.IsDark {
				return 4
			}
			if s.NeutralPalette.IsYellow() {
				return 90
			} else if s.Variant == Vibrant {
				return 85
			}
			return 87
		},
		IsBackground: true,
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.IsDark {
				switch s.Variant {
				case Neutral:
					return 2.5
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsBlue() {
						return 2.7
					}
					return 1.75
				case Vibrant:
					return 1.36
				}
			}
			return 1
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SurfaceDim(), V2025, color)
}

func (m MaterialSpec2025) SurfaceBright() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface_bright",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.IsDark {
				return 18
			}
			if s.NeutralPalette.IsBlue() {
				return 99
			} else if s.Variant == Vibrant {
				return 97
			}
			return 98
		},
		IsBackground: true,
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if !s.IsDark {
				return 1
			}
			switch s.Variant {
			case Neutral:
				return 2.5
			case TonalSpot:
				return 1.7
			case Vibrant:
				return 1.36
			default:
				return 1
			}
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SurfaceBright(), V2025, color)
}

func (m MaterialSpec2025) SurfaceContainerLowest() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface_container_lowest",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.IsDark {
				return 0.0
			}
			return 100.0
		},
		IsBackground: true,
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SurfaceContainerLowest(), V2025, color)
}

func (m MaterialSpec2025) SurfaceContainerLow() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface_container_low",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform != Phone {
				return 15
			}
			if s.IsDark {
				return 6
			} else if s.NeutralPalette.IsYellow() {
				return 98
			} else if s.Variant == Vibrant {
				return 95
			}
			return 96
		},
		IsBackground: true,
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.Platform != Phone {
				return 1
			}
			switch s.Variant {
			case Neutral:
				return 1.3
			case TonalSpot:
				return 1.25
			case Vibrant:
				return 1.08
			default:
				return 1
			}
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SurfaceContainerLow(), V2025, color)
}

func (m MaterialSpec2025) SurfaceContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform != Phone {
				return 20
			}
			if s.IsDark {
				return 9
			}
			if s.NeutralPalette.IsYellow() {
				return 96
			} else if s.Variant == Vibrant {
				return 92
			}
			return 94
		},
		IsBackground: true,
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.Platform != Phone {
				return 1
			}
			switch s.Variant {
			case Neutral:
				return 1.6
			case TonalSpot:
				return 1.4
			case Expressive:
				if s.NeutralPalette.IsYellow() {
					return 1.6
				}
				return 1.3
			case Vibrant:
				return 1.15
			default:
				return 1
			}
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SurfaceContainer(), V2025, color)
}

func (m MaterialSpec2025) SurfaceContainerHigh() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface_container_high",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				if s.IsDark {
					return 12
				} else {
					if s.NeutralPalette.IsYellow() {
						return 94
					} else if s.Variant == Vibrant {
						return 90
					} else {
						return 92
					}
				}
			} else {
				return 25
			}
		},
		IsBackground: true,
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 1.9
				case TonalSpot:
					return 1.5
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						return 1.95
					}
					return 1.45
				case Vibrant:
					return 1.22
				}
			}
			return 1
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SurfaceContainerHigh(), V2025, color)
}

func (m MaterialSpec2025) SurfaceContainerHighest() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "surface_container_highest",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.IsDark {
				return 15
			} else {
				if s.NeutralPalette.IsYellow() {
					return 92
				} else if s.Variant == Vibrant {
					return 88
				} else {
					return 90
				}
			}
		},
		IsBackground: true,
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			switch s.Variant {
			case Neutral:
				return 2.2
			case TonalSpot:
				return 1.7
			case Expressive:
				if s.NeutralPalette.IsYellow() {
					return 2.3
				}
				return 1.6
			case Vibrant:
				return 1.29
			default:
				return 1
			}
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SurfaceContainerHighest(), V2025, color)
}

func (m MaterialSpec2025) OnSurface() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_surface",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Variant == Vibrant {
				return tMaxC(s.NeutralPalette, 0, 100, 1.1)
			} else {
				// For all other variants, the initial tone should be the default
				// tone, which is the same as the background color.
				return GetInitialToneFromBackground(func(s *DynamicScheme) *DynamicColor {
					if s.Platform == Phone {
						return m.HighestSurface(s)
					} else {
						return m.SurfaceContainerHigh()
					}
				})(s)
			}
		},
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						if s.IsDark {
							return 3.0
						}
						return 2.3
					}
					return 1.6
				}
			}
			return 1
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainerHigh()
			}
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.IsDark {
				return GetCurve(11)
			}
			return GetCurve(9)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnSurface(), V2025, color)
}

func (m MaterialSpec2025) OnSurfaceVariant() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_surface_variant",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						if s.IsDark {
							return 3.0
						}
						return 2.3
					}
					return 1.6
				}
			}
			return 1
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainerHigh()
			}
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnSurfaceVariant(), V2025, color)
}

func (m MaterialSpec2025) Outline() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "outline",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						if s.IsDark {
							return 3.0
						}
						return 2.3
					} else {
						return 1.6
					}
				}
			}
			return 1
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainer()
			}
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(3)
			}
			return GetCurve(4.5)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.Outline(), V2025, color)
}

func (m MaterialSpec2025) OutlineVariant() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "outline_variant",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		ChromaMultiplier: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						if s.IsDark {
							return 3.0
						}
						return 2.3
					}
					return 1.6
				}
			}
			return 1
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(1.5)
			}
			return GetCurve(3)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OutlineVariant(), V2025, color)
}

func (m MaterialSpec2025) InverseSurface() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "inverse_surface",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.IsDark {
				return 98
			}
			return 4
		},
		IsBackground: true,
	})
	return ExtendSpecVersion(m.MaterialSpec2021.InverseSurface(), V2025, color)
}

func (m MaterialSpec2025) InverseOnSurface() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "inverse_on_surface",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.InverseSurface()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.InverseOnSurface(), V2025, color)
}

func (m MaterialSpec2025) Primary() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "primary",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			switch s.Variant {
			case Neutral:
				if s.Platform == Phone {
					if s.IsDark {
						return 80.0
					}
					return 40.0
				} else {
					return 90
				}
			case TonalSpot:
				if s.Platform == Phone {
					if s.IsDark {
						return 80
					} else {
						return tMaxC(s.PrimaryPalette, 0, 10)
					}
				} else {
					return tMaxC(s.PrimaryPalette, 0, 90)
				}
			case Expressive:
				if s.PrimaryPalette.IsYellow() {
					return tMaxC(s.PrimaryPalette, 0, 25)
				} else if s.PrimaryPalette.IsCyan() {
					return tMaxC(s.PrimaryPalette, 0, 88)
				} else {
					return tMaxC(s.PrimaryPalette, 0, 98)
				}
			default: // VIBRANT
				if s.PrimaryPalette.IsCyan() {
					return tMaxC(s.PrimaryPalette, 0, 88)
				} else {
					return tMaxC(s.PrimaryPalette, 0, 98)
				}
			}
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainerHigh()
			}
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.PrimaryContainer(),
					m.Primary(),
					5,
					ToneRelativeDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.Primary(), V2025, color)
}

func (m MaterialSpec2025) PrimaryDim() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "primary_dim",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s *DynamicScheme) float64 {
			switch s.Variant {
			case Neutral:
				return 85
			case TonalSpot:
				return tMaxC(s.PrimaryPalette, 0, 90)
			default:
				return tMaxC(s.PrimaryPalette)
			}
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryDim(), m.Primary(), 5, ToneDarker, true)
		},
	})
	return color
}

func (m MaterialSpec2025) OnPrimary() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_primary",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.Primary()
			} else {
				return m.PrimaryDim()
			}
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnPrimary(), V2025, color)
}

func (m MaterialSpec2025) PrimaryContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "primary_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Watch {
				return 30
			} else if s.Variant == Neutral {
				if s.IsDark {
					return 30.0
				}
				return 90.0
			} else if s.Variant == TonalSpot {
				if s.IsDark {
					return tMinC(s.PrimaryPalette, 35, 93)
				} else {
					return tMaxC(s.PrimaryPalette, 0, 90)
				}
			} else if s.Variant == Expressive {
				if s.IsDark {
					return tMaxC(s.PrimaryPalette, 30, 93)
				} else {
					if s.PrimaryPalette.IsCyan() {
						return tMaxC(s.PrimaryPalette, 78, 88)
					} else {
						return tMaxC(s.PrimaryPalette, 78, 90)
					}
				}
			} else { // VIBRANT
				if s.IsDark {
					return tMinC(s.PrimaryPalette, 66, 93)
				} else {
					if s.PrimaryPalette.IsCyan() {
						return tMaxC(s.PrimaryPalette, 66, 88)
					} else {
						return tMaxC(s.PrimaryPalette, 66, 93)
					}
				}
			}
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{}
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return nil
			}
			return NewToneDeltaPair(
				m.PrimaryContainer(),
				m.PrimaryDim(),
				10,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.PrimaryContainer(), V2025, color)
}

func (m MaterialSpec2025) OnPrimaryContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_primary_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.PrimaryContainer()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnPrimaryContainer(), V2025, color)
}

func (m MaterialSpec2025) PrimaryFixed() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "primary_fixed",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			s.IsDark = false
			return m.PrimaryContainer().GetTone(s)
		},
		IsBackground: true,
	})
	return ExtendSpecVersion(m.MaterialSpec2021.PrimaryFixed(), V2025, color)
}

func (m MaterialSpec2025) PrimaryFixedDim() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "primary_fixed_dim",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			return m.PrimaryFixed().GetTone(s)
		},
		IsBackground: true,
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.PrimaryFixedDim(),
				m.PrimaryFixed(),
				5,
				ToneDarker,
				true,
				ConstraintExact,
			)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.PrimaryFixedDim(), V2025, color)
}

func (m MaterialSpec2025) OnPrimaryFixed() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_primary_fixed",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.PrimaryFixedDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnPrimaryFixed(), V2025, color)
}

func (m MaterialSpec2025) OnPrimaryFixedVariant() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_primary_fixed_variant",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.PrimaryFixedDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnPrimaryFixedVariant(), V2025, color)
}

func (m MaterialSpec2025) InversePrimary() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "inverse_primary",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			return tMaxC(s.PrimaryPalette)
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.InverseSurface()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.InversePrimary(), V2025, color)
}

func (m MaterialSpec2025) Secondary() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "secondary",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Watch {
				if s.Variant == Neutral {
					return 90
				}
				return tMaxC(s.SecondaryPalette, 0, 90)
			} else if s.Variant == Neutral {
				if s.IsDark {
					return tMinC(s.SecondaryPalette, 0, 98)
				}
				return tMaxC(s.SecondaryPalette)
			} else if s.Variant == Vibrant {
				if s.IsDark {
					return tMaxC(s.SecondaryPalette, 0, 90)
				}
				return tMaxC(s.SecondaryPalette, 0, 98)
			} else { // EXPRESSIVE and TONAL_SPOT
				if s.IsDark {
					return 80
				}
				return tMaxC(s.SecondaryPalette)
			}
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.SecondaryContainer(),
					m.Secondary(),
					5,
					ToneRelativeLighter,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.Secondary(), V2025, color)
}

func (m MaterialSpec2025) SecondaryDim() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "secondary_dim",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Variant == Neutral {
				return 85
			}
			return tMaxC(s.SecondaryPalette, 0, 90)
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.SecondaryDim(),
				m.Secondary(),
				5,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
	})
	return color
}

func (m MaterialSpec2025) OnSecondary() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_secondary",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.Secondary()
			}
			return m.SecondaryDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnSecondary(), V2025, color)
}

func (m MaterialSpec2025) SecondaryContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "secondary_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Watch {
				return 30
			} else if s.Variant == Vibrant {
				if s.IsDark {
					return tMinC(s.SecondaryPalette, 30, 40)
				}
				return tMaxC(s.SecondaryPalette, 84, 90)
			} else if s.Variant == Expressive {
				if s.IsDark {
					return 15
				}
				return tMaxC(s.SecondaryPalette, 90, 95)
			} else {
				if s.IsDark {
					return 25
				}
				return 90
			}
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{}
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Watch {
				return NewToneDeltaPair(
					m.SecondaryContainer(),
					m.SecondaryDim(),
					10,
					ToneDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SecondaryContainer(), V2025, color)
}

func (m MaterialSpec2025) OnSecondaryContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_secondary_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.SecondaryContainer()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnSecondaryContainer(), V2025, color)
}

func (m MaterialSpec2025) SecondaryFixed() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "secondary_fixed",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			s.IsDark = false
			return m.SecondaryContainer().GetTone(s)
		},
		IsBackground: true,
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SecondaryFixed(), V2025, color)
}

func (m MaterialSpec2025) SecondaryFixedDim() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "secondary_fixed_dim",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			return m.SecondaryFixed().GetTone(s)
		},
		IsBackground: true,
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.SecondaryFixedDim(),
				m.SecondaryFixed(),
				5,
				ToneDarker,
				true,
				ConstraintExact,
			)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.SecondaryFixedDim(), V2025, color)
}

func (m MaterialSpec2025) OnSecondaryFixed() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "on_secondary_fixed",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.SecondaryFixedDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve { return GetCurve(7) },
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnSecondaryFixed(), V2025, color)
}

func (m MaterialSpec2025) OnSecondaryFixedVariant() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "on_secondary_fixed_variant",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.SecondaryFixedDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnSecondaryFixedVariant(), V2025, color)
}

func (m MaterialSpec2025) Tertiary() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "tertiary",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Watch {
				if s.Variant == TonalSpot {
					return tMaxC(s.TertiaryPalette, 0, 90)
				}
				return tMaxC(s.TertiaryPalette)
			} else if s.Variant == Expressive || s.Variant == Vibrant {
				limit := 100.0
				if s.TertiaryPalette.IsYellow() {
					limit = 88
				} else if s.IsDark {
					limit = 98
				}
				return tMaxC(s.TertiaryPalette, 0, limit)
			} else {
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 98)
				}
				return tMaxC(s.TertiaryPalette)
			}
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.TertiaryContainer(),
					m.Tertiary(),
					5,
					ToneRelativeLighter,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.Tertiary(), V2025, color)
}

func (m MaterialSpec2025) TertiaryDim() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "tertiary_dim",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s *DynamicScheme) float64 {
			if s.Variant == TonalSpot {
				return tMaxC(s.TertiaryPalette, 0, 90)
			}
			return tMaxC(s.TertiaryPalette)
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.TertiaryDim(),
				m.Tertiary(),
				5,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
	})
	return color
}

func (m MaterialSpec2025) OnTertiary() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "on_tertiary",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.Tertiary()
			}
			return m.TertiaryDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnTertiary(), V2025, color)
}

func (m MaterialSpec2025) TertiaryContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "tertiary_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Watch {
				if s.Variant == TonalSpot {
					return tMaxC(s.TertiaryPalette, 0, 90)
				}
				return tMaxC(s.TertiaryPalette)
			}
			switch s.Variant {
			case Neutral:
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 93)
				}
				return tMaxC(s.TertiaryPalette, 0, 96)
			case TonalSpot:
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 93)
				}
				return tMaxC(s.TertiaryPalette)
			case Expressive:
				upper := 100.0
				if s.TertiaryPalette.IsCyan() {
					upper = 88
				} else if s.IsDark {
					upper = 93
				}
				return tMaxC(s.TertiaryPalette, 75, upper)
			case Vibrant:
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 93)
				}
				return tMaxC(s.TertiaryPalette, 72, 100)
			}
			return tMaxC(s.TertiaryPalette) // fallback
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{} // undefined
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Watch {
				return NewToneDeltaPair(
					m.TertiaryContainer(),
					m.TertiaryDim(),
					10,
					ToneDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.TertiaryContainer(), V2025, color)
}

func (m MaterialSpec2025) OnTertiaryContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name: "on_tertiary_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.TertiaryContainer()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnTertiaryContainer(), V2025, color)
}

func (m MaterialSpec2025) TertiaryFixed() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "tertiary_fixed",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s *DynamicScheme) float64 {
			temp := s
			temp.IsDark = false
			temp.ContrastLevel = 0
			return m.TertiaryContainer().Tone(temp)
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == "phone" {
				return m.HighestSurface(s)
			}
			return nil
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == "phone" && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	})

	return ExtendSpecVersion(m.MaterialSpec2021.TertiaryFixed(), V2025, color)
}

func (m MaterialSpec2025) OnTertiaryFixed() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_tertiary_fixed",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.TertiaryFixedDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnTertiaryFixed(), V2025, color)
}

func (m MaterialSpec2025) OnTertiaryFixedVariant() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_tertiary_fixed_variant",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.TertiaryFixedDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnTertiaryFixedVariant(), V2025, color)
}

func (m MaterialSpec2025) Error() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "error",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Phone {
				if s.IsDark {
					return tMinC(s.ErrorPalette, 0, 98)
				}
				return tMaxC(s.ErrorPalette)
			}
			return tMinC(s.ErrorPalette)
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.ErrorContainer(),
					m.Error(),
					5,
					ToneRelativeLighter,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.Error(), V2025, color)
}

func (m MaterialSpec2025) ErrorDim() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "error_dim",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone: func(s *DynamicScheme) float64 {
			return tMinC(s.ErrorPalette)
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.ErrorDim(),
				m.Error(),
				5,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
	})
	return color
}

func (m MaterialSpec2025) OnError() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_error",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.Error()
			}
			return m.ErrorDim()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnError(), V2025, color)
}

func (m MaterialSpec2025) ErrorContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "error_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone: func(s *DynamicScheme) float64 {
			if s.Platform == Watch {
				return 30
			}
			if s.IsDark {
				return tMinC(s.ErrorPalette, 30, 93)
			}
			return tMaxC(s.ErrorPalette, 0, 90)
		},
		IsBackground: true,
		Background: func(s *DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{} // or `return nil` if *DynamicColor is a pointer type
		},
		ToneDeltaPair: func(s *DynamicScheme) *ToneDeltaPair {
			if s.Platform == Watch {
				return NewToneDeltaPair(
					m.ErrorContainer(),
					m.ErrorDim(),
					10,
					ToneDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.ErrorContainer(), V2025, color)
}

func (m MaterialSpec2025) OnErrorContainer() *DynamicColor {
	color := DynamicColorFromPalette(&DynamicColor{
		Name:    "on_error_container",
		Palette: func(s *DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Background: func(s *DynamicScheme) *DynamicColor {
			return m.ErrorContainer()
		},
		ContrastCurve: func(s *DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
	})
	return ExtendSpecVersion(m.MaterialSpec2021.OnErrorContainer(), V2025, color)
}
